package handler

import (
	"elderly-care-system/internal/service"
	"elderly-care-system/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RBACHandler struct {
	rbacService *service.RBACService
}

func NewRBACHandler(rbacService *service.RBACService) *RBACHandler {
	return &RBACHandler{rbacService: rbacService}
}

// === Permission Handlers ===

func (h *RBACHandler) CreatePermission(c *gin.Context) {
	var req service.PermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request")
		return
	}

	permission, err := h.rbacService.CreatePermission(&req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, permission)
}

func (h *RBACHandler) ListPermissions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "100"))

	permissions, total, err := h.rbacService.ListPermissions(page, pageSize)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":      permissions,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (h *RBACHandler) UpdatePermission(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req service.PermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request")
		return
	}

	if err := h.rbacService.UpdatePermission(uint(id), &req); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *RBACHandler) DeletePermission(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := h.rbacService.DeletePermission(uint(id)); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}

// === Menu Handlers ===

func (h *RBACHandler) CreateMenu(c *gin.Context) {
	var req service.MenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request")
		return
	}

	menu, err := h.rbacService.CreateMenu(&req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, menu)
}

func (h *RBACHandler) ListMenus(c *gin.Context) {
	menus, err := h.rbacService.ListMenus()
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, menus)
}

func (h *RBACHandler) UpdateMenu(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req service.MenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request")
		return
	}

	if err := h.rbacService.UpdateMenu(uint(id), &req); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *RBACHandler) DeleteMenu(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := h.rbacService.DeleteMenu(uint(id)); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}

// === Role Handlers ===

func (h *RBACHandler) GetRole(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	role, err := h.rbacService.GetRole(uint(id))
	if err != nil {
		response.Error(c, 404, "Role not found")
		return
	}

	response.Success(c, role)
}

func (h *RBACHandler) ListRoles(c *gin.Context) {
	roles, err := h.rbacService.ListRoles()
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, roles)
}

func (h *RBACHandler) AssignPermissionsToRole(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req service.RolePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request")
		return
	}

	if err := h.rbacService.AssignPermissionsToRole(uint(id), req.PermissionIDs); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *RBACHandler) AssignMenusToRole(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req service.RoleMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request")
		return
	}

	if err := h.rbacService.AssignMenusToRole(uint(id), req.MenuIDs); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}

// === User Role Handlers ===

func (h *RBACHandler) GetUserMenus(c *gin.Context) {
	userID, _ := c.Get("user_id")

	menus, err := h.rbacService.GetUserMenus(userID.(uint))
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, menus)
}

func (h *RBACHandler) GetUserPermissions(c *gin.Context) {
	userID, _ := c.Get("user_id")

	permissions, err := h.rbacService.GetUserPermissions(userID.(uint))
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, permissions)
}

func (h *RBACHandler) AssignRoleToUser(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req struct {
		RoleID uint `json:"role_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request")
		return
	}

	if err := h.rbacService.AssignRoleToUser(uint(id), req.RoleID); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *RBACHandler) RemoveRoleFromUser(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req struct {
		RoleID uint `json:"role_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request")
		return
	}

	if err := h.rbacService.RemoveRoleFromUser(uint(id), req.RoleID); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *RBACHandler) GetUserRoles(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	roles, err := h.rbacService.GetUserRoles(uint(id))
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, roles)
}

func (h *RBACHandler) UpdateUserRoles(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req struct {
		RoleIDs []uint `json:"role_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request")
		return
	}

	if err := h.rbacService.UpdateUserRoles(uint(id), req.RoleIDs); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}

// === System Management - Users List ===

func (h *RBACHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	// 这里需要调用UserService，暂时返回空列表
	response.Success(c, gin.H{
		"list":      []interface{}{},
		"total":     0,
		"page":      page,
		"page_size": pageSize,
		"keyword":   keyword,
	})
}

// === Initialization Handlers ===

func (h *RBACHandler) InitializeSystem(c *gin.Context) {
	// 初始化默认权限、菜单和管理员角色
	if err := h.rbacService.InitializeDefaultPermissions(); err != nil {
		response.Error(c, 500, "Failed to initialize permissions: "+err.Error())
		return
	}

	if err := h.rbacService.InitializeDefaultMenus(); err != nil {
		response.Error(c, 500, "Failed to initialize menus: "+err.Error())
		return
	}

	if err := h.rbacService.InitializeAdminRole(); err != nil {
		response.Error(c, 500, "Failed to initialize admin role: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "System initialized successfully",
	})
}
