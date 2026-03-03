package handler

import (
	"elderly-care-system/internal/service"
	"elderly-care-system/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PushNotificationHandler struct {
	pushService *service.PushNotificationService
}

func NewPushNotificationHandler(pushService *service.PushNotificationService) *PushNotificationHandler {
	return &PushNotificationHandler{pushService: pushService}
}

// RegisterToken 注册设备推送Token
func (h *PushNotificationHandler) RegisterToken(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req service.RegisterTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request")
		return
	}

	if err := h.pushService.RegisterToken(userID, &req); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "Token registered successfully"})
}

// SendNotification 发送推送通知
func (h *PushNotificationHandler) SendNotification(c *gin.Context) {
	var req service.SendPushNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request")
		return
	}

	if err := h.pushService.SendPushNotification(&req); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "Notification sent successfully"})
}

// GetNotifications 获取通知列表
func (h *PushNotificationHandler) GetNotifications(c *gin.Context) {
	userID := c.GetUint("user_id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	notifications, err := h.pushService.GetNotifications(userID, limit)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, notifications)
}

// UnregisterToken 注销推送Token
func (h *PushNotificationHandler) UnregisterToken(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request")
		return
	}

	if err := h.pushService.UnregisterToken(req.Token); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "Token unregistered successfully"})
}

// BroadcastNotification 广播通知 (管理员功能)
func (h *PushNotificationHandler) BroadcastNotification(c *gin.Context) {
	var req struct {
		Role    string                 `json:"role" binding:"required"`
		Title   string                 `json:"title" binding:"required"`
		Content string                 `json:"content" binding:"required"`
		Data    map[string]interface{} `json:"data"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request")
		return
	}

	if err := h.pushService.BroadcastToRole(req.Role, req.Title, req.Content, req.Data); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "Broadcast sent successfully"})
}
