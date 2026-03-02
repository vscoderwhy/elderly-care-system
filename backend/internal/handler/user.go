package handler

import (
	"elderly-care-system/internal/service"
	"elderly-care-system/pkg/response"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	user, err := h.userService.GetProfile(userID.(uint))
	if err != nil {
		response.Error(c, 404, "User not found")
		return
	}

	response.Success(c, user)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.Error(c, 400, "Invalid request")
		return
	}

	if err := h.userService.UpdateProfile(userID.(uint), updates); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *UserHandler) GetElderlyList(c *gin.Context) {
	userID, _ := c.Get("user_id")

	elderlyList, err := h.userService.GetElderlyList(userID.(uint))
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, elderlyList)
}

// ListStaff 获取员工列表
func (h *UserHandler) ListStaff(c *gin.Context) {
	role := c.Query("role")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	staff, total, err := h.userService.ListStaff(role, page, pageSize)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":       staff,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
	})
}

// CreateStaff 创建员工
func (h *UserHandler) CreateStaff(c *gin.Context) {
	var req service.CreateStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	user, err := h.userService.CreateStaff(&req)
	if err != nil {
		if errors.Is(err, service.ErrPhoneExists) {
			response.Error(c, 400, "手机号已存在")
			return
		}
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, user)
}

// UpdateStaff 更新员工
func (h *UserHandler) UpdateStaff(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.Error(c, 400, "Invalid request")
		return
	}

	if err := h.userService.UpdateStaff(uint(id), updates); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}

// DeleteStaff 删除员工
func (h *UserHandler) DeleteStaff(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := h.userService.DeleteStaff(uint(id)); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}
