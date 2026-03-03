package service

import (
	"elderly-care-system/internal/config"
	"elderly-care-system/internal/model"
	"elderly-care-system/internal/repository"
	"encoding/json"
	"fmt"
	"time"

	"github.com/silenceper/wechat/v2/miniprogram"
	"github.com/silenceper/wechat/v2"
)

type PushNotificationService struct {
	pushRepo     *repository.PushNotificationRepository
	userRepo     *repository.UserRepository
	cfg          config.Config
	wechatApp    *miniprogram.MiniProgram
}

func NewPushNotificationService(
	pushRepo *repository.PushNotificationRepository,
	userRepo *repository.UserRepository,
	cfg config.Config,
) *PushNotificationService {
	s := &PushNotificationService{
		pushRepo: pushRepo,
		userRepo: userRepo,
		cfg:      cfg,
	}

	// 初始化微信小程序SDK
	if cfg.WeChat.AppID != "" {
		wechat := wechat.NewWechat()
		wechat.SetAppIDAndSecret(cfg.WeChat.AppID, cfg.WeChat.AppSecret)
		s.wechatApp = wechat.GetMiniProgram()
	}

	return s
}

// RegisterTokenRequest 注册推送Token请求
type RegisterTokenRequest struct {
	Token    string `json:"token" binding:"required"`
	Platform string `json:"platform" binding:"required"` // wechat, android, ios
	DeviceID string `json:"device_id"`
}

// RegisterToken 注册设备推送Token
func (s *PushNotificationService) RegisterToken(userID uint, req *RegisterTokenRequest) error {
	token := &model.PushToken{
		UserID:     userID,
		Token:      req.Token,
		Platform:   req.Platform,
		DeviceID:   req.DeviceID,
		IsActive:   true,
		LastUsedAt: &[]time.Time{time.Now()}[0],
	}

	return s.pushRepo.SaveToken(token)
}

// SendPushNotificationRequest 发送推送请求
type SendPushNotificationRequest struct {
	UserID  uint                   `json:"user_id" binding:"required"`
	Title   string                 `json:"title" binding:"required"`
	Content string                 `json:"content" binding:"required"`
	Type    string                 `json:"type"` // emergency, warning, reminder, info
	Data    map[string]interface{} `json:"data"`
	Platform string                `json:"platform"` // 指定平台，不指定则发所有平台
}

// SendPushNotification 发送推送通知
func (s *PushNotificationService) SendPushNotification(req *SendPushNotificationRequest) error {
	// 获取用户的推送Token
	tokens, err := s.pushRepo.GetUserTokens(req.UserID, req.Platform)
	if err != nil || len(tokens) == 0 {
		return fmt.Errorf("没有可用的推送Token")
	}

	// 序列化附加数据
	dataJSON := ""
	if req.Data != nil {
		if bytes, err := json.Marshal(req.Data); err == nil {
			dataJSON = string(bytes)
		}
	}

	var lastErr error
	for _, token := range tokens {
		// 创建推送记录
		notification := &model.PushNotification{
			UserID:   req.UserID,
			Title:    req.Title,
			Content:  req.Content,
			Type:     req.Type,
			Data:     dataJSON,
			Platform: token.Platform,
			Status:   "pending",
		}

		if err := s.pushRepo.CreateNotification(notification); err != nil {
			continue
		}

		// 根据平台发送推送
		var sendErr error
		switch token.Platform {
		case "wechat":
			sendErr = s.sendWeChatPush(token.Token, req.Title, req.Content, req.Data)
		case "android":
			sendErr = s.sendAndroidPush(token.Token, req.Title, req.Content, req.Data)
		case "ios":
			sendErr = s.sendIOSPush(token.Token, req.Title, req.Content, req.Data)
		default:
			sendErr = fmt.Errorf("不支持的平台: %s", token.Platform)
		}

		now := time.Now()
		if sendErr != nil {
			s.pushRepo.UpdateNotificationStatus(notification.ID, "failed", &now, sendErr.Error())
			lastErr = sendErr
		} else {
			s.pushRepo.UpdateNotificationStatus(notification.ID, "sent", &now, "")
		}
	}

	return lastErr
}

// sendWeChatPush 发送微信小程序订阅消息
func (s *PushNotificationService) sendWeChatPush(openID, title, content string, data map[string]interface{}) error {
	if s.wechatApp == nil {
		return fmt.Errorf("微信小程序未配置")
	}

	// 微信小程序订阅消息推送
	// 注意：需要用户先订阅消息模板
	subscribe := s.wechatApp.GetSubscribe()

	// 构造消息数据
	msgData := map[string]interface{}{
		"thing1": map[string]string{"value": title},
		"thing2": map[string]string{"value": content},
	}

	err := subscribe.Send(&miniprogram.Message{
		ToUser:     openID,
		TemplateID: s.cfg.WeChat.TemplateID, // 需要在配置中添加模板ID
		Page:       "pages/index/index",
		Data:       msgData,
	})

	return err
}

// sendAndroidPush 发送Android推送 (使用极光推送、个推等)
func (s *PushNotificationService) sendAndroidPush(token, title, content string, data map[string]interface{}) error {
	// TODO: 集成极光推送SDK
	// 这里是示例代码
	return fmt.Errorf("Android推送未实现")
}

// sendIOSPush 发送iOS推送 (使用APNs)
func (s *PushNotificationService) sendIOSPush(token, title, content string, data map[string]interface{}) error {
	// TODO: 集成APNs
	// 这里是示例代码
	return fmt.Errorf("iOS推送未实现")
}

// BroadcastToRole 向指定角色的所有用户发送广播
func (s *PushNotificationService) BroadcastToRole(role string, title, content string, data map[string]interface{}) error {
	// 获取该角色的所有用户
	users, err := s.userRepo.ListUsers(1, 1000, "", role)
	if err != nil {
		return err
	}

	for _, user := range users {
		req := &SendPushNotificationRequest{
			UserID:  user.ID,
			Title:   title,
			Content: content,
			Data:    data,
		}
		s.SendPushNotification(req)
	}

	return nil
}

// GetNotifications 获取用户的通知记录
func (s *PushNotificationService) GetNotifications(userID uint, limit int) ([]model.PushNotification, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.pushRepo.GetUserNotifications(userID, limit)
}

// UnregisterToken 注销推送Token
func (s *PushNotificationService) UnregisterToken(token string) error {
	return s.pushRepo.DeleteToken(token)
}
