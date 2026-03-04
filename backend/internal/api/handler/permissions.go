package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/elderly-care/internal/models"
	"gorm.io/gorm"
)

type PermissionsHandler struct {
	db *gorm.DB
}

func NewPermissionsHandler(db *gorm.DB) *PermissionsHandler {
	return &PermissionsHandler{db: db}
}

// ========== 角色管理 ==========

// ListRoles 获取角色列表
func (h *PermissionsHandler) ListRoles(c *gin.Context) {
	var roles []models.Role

	if err := h.db.Find(&roles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取角色列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    roles,
	})
}

// CreateRole 创建角色
func (h *PermissionsHandler) CreateRole(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		DisplayName string `json:"displayName" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	// 检查角色名是否已存在
	var count int64
	if err := h.db.Model(&models.Role{}).Where("name = ?", req.Name).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "检查角色失败",
		})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "角色名已存在",
		})
		return
	}

	role := models.Role{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Status:      "active",
	}

	if err := h.db.Create(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建角色失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    role,
	})
}

// UpdateRole 更新角色
func (h *PermissionsHandler) UpdateRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的角色ID",
		})
		return
	}

	var req struct {
		DisplayName string `json:"displayName"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	var role models.Role
	if err := h.db.First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "角色不存在",
		})
		return
	}

	updates := map[string]interface{}{}
	if req.DisplayName != "" {
		updates["display_name"] = req.DisplayName
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	if err := h.db.Model(&role).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新角色失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    role,
	})
}

// DeleteRole 删除角色
func (h *PermissionsHandler) DeleteRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的角色ID",
		})
		return
	}

	// 检查是否有用户使用此角色
	var count int64
	if err := h.db.Model(&models.UserRole{}).Where("role_id = ?", id).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "检查角色使用情况失败",
		})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "该角色正在使用中，无法删除",
		})
		return
	}

	if err := h.db.Delete(&models.Role{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除角色失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}

// ========== 权限管理 ==========

// ListPermissions 获取权限列表
func (h *PermissionsHandler) ListPermissions(c *gin.Context) {
	var permissions []models.Permission

	module := c.Query("module")
	if module != "" {
		if err := h.db.Where("module = ?", module).Find(&permissions).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "获取权限列表失败",
			})
			return
		}
	} else {
		if err := h.db.Find(&permissions).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "获取权限列表失败",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    permissions,
	})
}

// GetRolePermissions 获取角色的权限列表
func (h *PermissionsHandler) GetRolePermissions(c *gin.Context) {
	idStr := c.Param("id")
	roleID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的角色ID",
		})
		return
	}

	var rolePermissions []models.RolePermission
	if err := h.db.Where("role_id = ?", roleID).Find(&rolePermissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取角色权限失败",
		})
		return
	}

	var permissionIDs []uint
	for _, rp := range rolePermissions {
		permissionIDs = append(permissionIDs, rp.PermissionID)
	}

	var permissions []models.Permission
	if len(permissionIDs) > 0 {
		if err := h.db.Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "获取权限详情失败",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    permissions,
	})
}

// AssignRolePermission 为角色分配权限
func (h *PermissionsHandler) AssignRolePermission(c *gin.Context) {
	idStr := c.Param("id")
	roleID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的角色ID",
		})
		return
	}

	var req struct {
		PermissionIDs []uint `json:"permissionIds" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	// 删除现有的角色权限
	if err := h.db.Where("role_id = ?", roleID).Delete(&models.RolePermission{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除旧权限失败",
		})
		return
	}

	// 添加新的角色权限
	for _, permissionID := range req.PermissionIDs {
		rolePermission := models.RolePermission{
			RoleID:       uint(roleID),
			PermissionID: permissionID,
		}
		if err := h.db.Create(&rolePermission).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "添加权限失败",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "分配权限成功",
	})
}

// ========== 用户角色管理 ==========

// GetUserRoles 获取用户的角色列表
func (h *PermissionsHandler) GetUserRoles(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的用户ID",
		})
		return
	}

	var userRoles []models.UserRole
	if err := h.db.Where("user_id = ?", userID).Find(&userRoles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取用户角色失败",
		})
		return
	}

	var roleIDs []uint
	for _, ur := range userRoles {
		roleIDs = append(roleIDs, ur.RoleID)
	}

	var roles []models.Role
	if len(roleIDs) > 0 {
		if err := h.db.Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "获取角色详情失败",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    roles,
	})
}

// AssignUserRole 为用户分配角色
func (h *PermissionsHandler) AssignUserRole(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的用户ID",
		})
		return
	}

	var req struct {
		RoleIDs []uint `json:"roleIds" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	// 删除现有的用户角色
	if err := h.db.Where("user_id = ?", userID).Delete(&models.UserRole{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除旧角色失败",
		})
		return
	}

	// 添加新的用户角色
	for _, roleID := range req.RoleIDs {
		userRole := models.UserRole{
			UserID: uint(userID),
			RoleID: roleID,
		}
		if err := h.db.Create(&userRole).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "添加角色失败",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "分配角色成功",
	})
}

// ========== 操作日志 ==========

// CreateOperationLog 创建操作日志
func (h *PermissionsHandler) CreateOperationLog(c *gin.Context, module, action string, statusCode int) {
	userID, _ := c.Get("userId")
	username, _ := c.Get("username")

	log := models.SystemLog{
		UserID:   userID.(uint),
		Username: username.(string),
		Module:   module,
		Action:   action,
		Method:   c.Request.Method,
		Path:     c.Request.URL.Path,
		IP:       c.ClientIP(),
		Status:   statusCode,
		Latency:  int(time.Since(time.Now()).Milliseconds()),
	}

	h.db.Create(&log)
}

// GetOperationLogs 获取操作日志列表
func (h *PermissionsHandler) GetOperationLogs(c *gin.Context) {
	var logs []models.SystemLog

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	username := c.Query("username")
	module := c.Query("module")
	action := c.Query("action")

	query := h.db.Model(&models.SystemLog{})

	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if module != "" {
		query = query.Where("module = ?", module)
	}
	if action != "" {
		query = query.Where("action = ?", action)
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取日志列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    logs,
		"total":   total,
		"page":    page,
		"pageSize": pageSize,
	})
}
