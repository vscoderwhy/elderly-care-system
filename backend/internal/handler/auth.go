package handler

import (
	"elderly-care-system/internal/service"
	"elderly-care-system/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	token, user, err := h.authService.Register(&req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{
		"token": token,
		"user":  user,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	token, user, err := h.authService.Login(&req)
	if err != nil {
		response.Error(c, 401, err.Error())
		return
	}

	response.Success(c, gin.H{
		"token": token,
		"user":  user,
	})
}

func (h *AuthHandler) WeChatLogin(c *gin.Context) {
	var req service.WeChatLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	token, user, err := h.authService.WeChatLogin(&req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{
		"token": token,
		"user":  user,
	})
}
