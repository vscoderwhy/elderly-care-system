package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/elderly-care/internal/config"
)

// ValidationError 验证错误
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// BusinessError 业务错误
type BusinessError struct {
	Code    int
	Message string
}

func (e *BusinessError) Error() string {
	return fmt.Sprintf("code=%d message=%s", e.Code, e.Message)
}

type Middleware struct {
	jwtSecret string
	mu        sync.RWMutex
	lastClean time.Time
}

func NewMiddleware(cfg *config.Config) *Middleware {
	return &Middleware{
		jwtSecret: cfg.JWT.Secret,
	}
}

// CORS 跨域中间件
func (m *Middleware) CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// Auth JWT认证中间件
func (m *Middleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "请先登录",
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "认证格式错误",
			})
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := m.parseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "无效的token",
			})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文
		c.Set("userId", claims["userId"])
		c.Set("username", claims["username"])
		c.Set("role", claims["role"])

		c.Next()
	}
}

// Role 角色验证中间件
func (m *Middleware) Role(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString("role")

		// 管理员拥有所有权限
		if userRole == "admin" {
			c.Next()
			return
		}

		for _, role := range roles {
			if userRole == role {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "权限不足",
		})
		c.Abort()
	}
}

// Logger 日志中间件
func (m *Middleware) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		// 获取请求ID
		requestID := c.GetString("requestId")

		// 构建日志条目
		logEntry := fmt.Sprintf(
			"[%s] %s %s | status=%d | latency=%v | client_ip=%s | user_agent=%s",
			requestID,
			c.Request.Method,
			path,
			c.Writer.Status(),
			latency,
			c.ClientIP(),
			c.Request.UserAgent(),
		)

		if query != "" {
			logEntry += fmt.Sprintf(" | query=%s", query)
		}

		// 记录慢请求
		if latency > time.Second {
			logEntry += " | SLOW_REQUEST"
		}

		// 记录错误响应
		if c.Writer.Status() >= 400 {
			logEntry += fmt.Sprintf(" | ERROR_RESPONSE")
		}

		// TODO: 使用结构化日志库（如 zap, logrus）
		fmt.Println(logEntry)
	}
}

// RequestID 请求ID中间件
func (m *Middleware) RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 尝试从请求头获取请求ID
		requestID := c.GetHeader("X-Request-ID")

		// 如果没有则生成新的
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// 存储到上下文
		c.Set("requestId", requestID)

		// 设置到响应头
		c.Writer.Header().Set("X-Request-ID", requestID)

		c.Next()
	}
}

// RateLimiter 简单的限流中间件（基于内存）
// 生产环境建议使用 Redis
func (m *Middleware) RateLimiter(requestsPerMinute int) gin.HandlerFunc {
	type client struct {
		lastRequest time.Time
		requests    int
	}

	clients := make(map[string]*client)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		// 清理过期的客户端记录（每分钟清理一次）
		if now.Minute() != m.lastClean.Minute() {
			m.mu.Lock()
			for key, client := range clients {
				if now.Sub(client.lastRequest) > time.Minute {
					delete(clients, key)
				}
			}
			m.lastClean = now
			m.mu.Unlock()
		}

		m.mu.Lock()
		currentClient, exists := clients[ip]
		if !exists {
			clients[ip] = &client{
				lastRequest: now,
				requests:    1,
			}
			m.mu.Unlock()
			c.Next()
			return
		}

		// 如果距离上次请求超过1分钟，重置计数
		if now.Sub(currentClient.lastRequest) > time.Minute {
			currentClient.requests = 1
			currentClient.lastRequest = now
			m.mu.Unlock()
			c.Next()
			return
		}

		// 检查是否超过限制
		if currentClient.requests >= requestsPerMinute {
			m.mu.Unlock()
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}

		currentClient.requests++
		currentClient.lastRequest = now
		m.mu.Unlock()

		c.Next()
	}
}

// Security 安全头中间件
func (m *Middleware) Security() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Writer.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Writer.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		c.Next()
	}
}

// ErrorHandler 统一错误处理中间件
func (m *Middleware) ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			// 根据错误类型返回不同的响应
			switch e := err.Err.(type) {
			case *ValidationError:
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    400,
					"message": e.Message,
				})
			case *BusinessError:
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    e.Code,
					"message": e.Message,
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "服务器内部错误",
				})
			}

			return
		}
	}
}

// Recovery 恢复中间件
func (m *Middleware) Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID := c.GetString("requestId")

				// 记录panic详情
				panicLog := fmt.Sprintf(
					"[PANIC] requestId=%s | error=%v | stack=%s",
					requestID,
					err,
					fmt.Sprintf("%+v", err),
				)

				// TODO: 发送到日志系统或监控系统
				fmt.Println(panicLog)

				// 返回友好的错误信息
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":       500,
					"message":    "服务器内部错误",
					"requestId":  requestID,
					"timestamp":  time.Now().Unix(),
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// parseToken 解析JWT token
func (m *Middleware) parseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
