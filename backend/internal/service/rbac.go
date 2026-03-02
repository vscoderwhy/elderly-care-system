package service

import (
	"elderly-care-system/internal/model"
	"elderly-care-system/internal/repository"
)

type RBACService struct {
	permissionRepo *repository.PermissionRepository
	menuRepo       *repository.MenuRepository
	roleRepo       *repository.RoleRepository
	userRepo       *repository.UserRepository
}

func NewRBACService(
	permissionRepo *repository.PermissionRepository,
	menuRepo *repository.MenuRepository,
	roleRepo *repository.RoleRepository,
	userRepo *repository.UserRepository,
) *RBACService {
	return &RBACService{
		permissionRepo: permissionRepo,
		menuRepo:       menuRepo,
		roleRepo:       roleRepo,
		userRepo:       userRepo,
	}
}

// === Permission Management ===

type PermissionRequest struct {
	Code        string `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

func (s *RBACService) CreatePermission(req *PermissionRequest) (*model.Permission, error) {
	permission := &model.Permission{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Status:      "active",
	}

	if err := s.permissionRepo.Create(permission); err != nil {
		return nil, err
	}

	return permission, nil
}

func (s *RBACService) ListPermissions(page, pageSize int) ([]model.Permission, int64, error) {
	return s.permissionRepo.List(page, pageSize)
}

func (s *RBACService) UpdatePermission(id uint, req *PermissionRequest) error {
	permission, err := s.permissionRepo.GetByID(id)
	if err != nil {
		return err
	}

	permission.Name = req.Name
	permission.Description = req.Description
	if req.Type != "" {
		permission.Type = req.Type
	}

	return s.permissionRepo.Update(permission)
}

func (s *RBACService) DeletePermission(id uint) error {
	return s.permissionRepo.Delete(id)
}

// === Menu Management ===

type MenuRequest struct {
	ParentID   *uint  `json:"parent_id"`
	Name       string `json:"name" binding:"required"`
	Path       string `json:"path"`
	Component  string `json:"component"`
	Icon       string `json:"icon"`
	Sort       int    `json:"sort"`
	Type       string `json:"type"`
	Permission string `json:"permission"`
	Visible    bool   `json:"visible"`
}

func (s *RBACService) CreateMenu(req *MenuRequest) (*model.Menu, error) {
	menu := &model.Menu{
		ParentID:   req.ParentID,
		Name:       req.Name,
		Path:       req.Path,
		Component:  req.Component,
		Icon:       req.Icon,
		Sort:       req.Sort,
		Type:       req.Type,
		Permission: req.Permission,
		Visible:    req.Visible,
		Status:     "active",
	}

	if err := s.menuRepo.Create(menu); err != nil {
		return nil, err
	}

	return menu, nil
}

func (s *RBACService) ListMenus() ([]model.Menu, error) {
	return s.menuRepo.List()
}

func (s *RBACService) UpdateMenu(id uint, req *MenuRequest) error {
	menu, err := s.menuRepo.GetByID(id)
	if err != nil {
		return err
	}

	menu.Name = req.Name
	menu.Path = req.Path
	menu.Component = req.Component
	menu.Icon = req.Icon
	menu.Sort = req.Sort
	menu.Permission = req.Permission
	menu.Visible = req.Visible

	if req.ParentID != nil {
		menu.ParentID = req.ParentID
	}

	return s.menuRepo.Update(menu)
}

func (s *RBACService) DeleteMenu(id uint) error {
	return s.menuRepo.Delete(id)
}

// === Role Management ===

type RolePermissionRequest struct {
	PermissionIDs []uint `json:"permission_ids"`
}

type RoleMenuRequest struct {
	MenuIDs []uint `json:"menu_ids"`
}

func (s *RBACService) GetRole(id uint) (*model.Role, error) {
	return s.roleRepo.GetByID(id)
}

func (s *RBACService) ListRoles() ([]model.Role, error) {
	return s.roleRepo.List()
}

func (s *RBACService) AssignPermissionsToRole(roleID uint, permissionIDs []uint) error {
	return s.roleRepo.AssignPermissions(roleID, permissionIDs)
}

func (s *RBACService) AssignMenusToRole(roleID uint, menuIDs []uint) error {
	return s.menuRepo.AssignPermissions(roleID, menuIDs)
}

// === User Management with Roles ===

func (s *RBACService) GetUserMenus(userID uint) ([]model.Menu, error) {
	menus, err := s.roleRepo.GetUserMenus(userID)
	if err != nil {
		return nil, err
	}

	// 构建菜单树
	return s.buildMenuTree(menus, nil), nil
}

func (s *RBACService) GetUserPermissions(userID uint) ([]model.Permission, error) {
	return s.roleRepo.GetUserPermissions(userID)
}

func (s *RBACService) HasPermission(userID uint, permissionCode string) bool {
	permissions, err := s.roleRepo.GetUserPermissions(userID)
	if err != nil {
		return false
	}

	for _, p := range permissions {
		if p.Code == permissionCode {
			return true
		}
	}
	return false
}

func (s *RBACService) AssignRoleToUser(userID uint, roleID uint) error {
	return s.userRepo.AssignRole(userID, roleID)
}

func (s *RBACService) RemoveRoleFromUser(userID uint, roleID uint) error {
	return s.userRepo.RemoveRole(userID, roleID)
}

// buildMenuTree 构建菜单树
func (s *RBACService) buildMenuTree(menus []model.Menu, parentID *uint) []model.Menu {
	var result []model.Menu

	for _, menu := range menus {
		// 检查是否是当前层级的菜单
		if (parentID == nil && menu.ParentID == nil) || (parentID != nil && menu.ParentID != nil && *menu.ParentID == *parentID) {
			// 递归获取子菜单
			menu.Children = s.buildMenuTree(menus, &menu.ID)
			result = append(result, menu)
		}
	}

	return result
}

// InitializeDefaultPermissions 初始化默认权限
func (s *RBACService) InitializeDefaultPermissions() error {
	permissions := []model.Permission{
		// 用户管理权限
		{Code: "user:list", Name: "查看用户", Type: "api"},
		{Code: "user:create", Name: "创建用户", Type: "api"},
		{Code: "user:update", Name: "更新用户", Type: "api"},
		{Code: "user:delete", Name: "删除用户", Type: "api"},
		// 角色管理权限
		{Code: "role:list", Name: "查看角色", Type: "api"},
		{Code: "role:create", Name: "创建角色", Type: "api"},
		{Code: "role:update", Name: "更新角色", Type: "api"},
		{Code: "role:delete", Name: "删除角色", Type: "api"},
		// 菜单管理权限
		{Code: "menu:list", Name: "查看菜单", Type: "api"},
		{Code: "menu:create", Name: "创建菜单", Type: "api"},
		{Code: "menu:update", Name: "更新菜单", Type: "api"},
		{Code: "menu:delete", Name: "删除菜单", Type: "api"},
		// 老人管理权限
		{Code: "elderly:list", Name: "查看老人", Type: "api"},
		{Code: "elderly:create", Name: "创建老人档案", Type: "api"},
		{Code: "elderly:update", Name: "更新老人档案", Type: "api"},
		{Code: "elderly:delete", Name: "删除老人档案", Type: "api"},
		// 护理管理权限
		{Code: "care:list", Name: "查看护理记录", Type: "api"},
		{Code: "care:create", Name: "创建护理记录", Type: "api"},
		// 用药管理权限
		{Code: "medication:list", Name: "查看用药", Type: "api"},
		{Code: "medication:create", Name: "创建用药记录", Type: "api"},
		{Code: "medication:update", Name: "更新用药", Type: "api"},
		{Code: "medication:delete", Name: "删除用药", Type: "api"},
		// 财务管理权限
		{Code: "finance:list", Name: "查看财务", Type: "api"},
		{Code: "finance:create", Name: "创建账单", Type: "api"},
		// 系统设置权限
		{Code: "system:settings", Name: "系统设置", Type: "api"},
	}

	for _, p := range permissions {
		s.permissionRepo.Create(&p)
	}

	return nil
}

// InitializeDefaultMenus 初始化默认菜单
func (s *RBACService) InitializeDefaultMenus() error {
	menus := []model.Menu{
		// 工作台
		{Name: "工作台", Path: "/dashboard", Icon: "DataBoard", Sort: 1, Type: "menu"},
		// 老人管理
		{Name: "老人管理", Path: "/elderly", Icon: "User", Sort: 2, Type: "menu", Permission: "elderly:list"},
		// 护理管理
		{Name: "护理记录", Path: "/care", Icon: "Notebook", Sort: 3, Type: "menu", Permission: "care:list"},
		// 用药管理
		{Name: "用药管理", Path: "/medications", Icon: "Medicine", Sort: 4, Type: "menu", Permission: "medication:list"},
		// 探视预约
		{Name: "探视预约", Path: "/visits", Icon: "Calendar", Sort: 5, Type: "menu"},
		// 智能预警
		{Name: "智能预警", Path: "/alerts", Icon: "Warning", Sort: 6, Type: "menu"},
		// 财务管理
		{Name: "财务管理", Path: "/bills", Icon: "Wallet", Sort: 7, Type: "menu", Permission: "finance:list"},
		// 数据导出
		{Name: "数据导出", Path: "/export", Icon: "Download", Sort: 8, Type: "menu"},
		// 员工管理
		{Name: "员工管理", Path: "/staff", Icon: "Avatar", Sort: 9, Type: "menu"},
		// 房间管理
		{Name: "房间管理", Path: "/rooms", Icon: "House", Sort: 10, Type: "menu"},
		// 系统管理
		{Name: "系统管理", Path: "/system", Icon: "Setting", Sort: 11, Type: "menu"},
	}

	for _, m := range menus {
		s.menuRepo.Create(&m)
	}

	return nil
}

// InitializeAdminRole 初始化管理员角色并分配所有权限
func (s *RBACService) InitializeAdminRole() error {
	// 获取所有权限和菜单
	permissions, _, _ := s.permissionRepo.List(1, 1000)
	menus, _ := s.menuRepo.List()

	// 为admin角色分配所有权限和菜单
	var permIDs []uint
	for _, p := range permissions {
		permIDs = append(permIDs, p.ID)
	}

	var menuIDs []uint
	for _, m := range menus {
		menuIDs = append(menuIDs, m.ID)
		s.assignMenuRecursive(&m, &menuIDs)
	}

	s.roleRepo.AssignPermissions(4, permIDs)
	s.menuRepo.AssignPermissions(4, menuIDs)

	return nil
}

func (s *RBACService) assignMenuRecursive(menu *model.Menu, menuIDs *[]uint) {
	*menuIDs = append(*menuIDs, menu.ID)
	for _, child := range menu.Children {
		s.assignMenuRecursive(&child, menuIDs)
	}
}

// CheckUserMenuAccess 检查用户是否有权访问指定菜单
func (s *RBACService) CheckUserMenuAccess(userID uint, menuPath string) bool {
	menus, err := s.roleRepo.GetUserMenus(userID)
	if err != nil {
		return false
	}

	for _, menu := range menus {
		if menu.Path == menuPath {
			return true
		}
	}
	return false
}

// CheckUserPermission 检查用户是否有指定权限
func (s *RBACService) CheckUserPermission(userID uint, permissionCode string) bool {
	return s.HasPermission(userID, permissionCode)
}

// GetUserRoles 获取用户的所有角色
func (s *RBACService) GetUserRoles(userID uint) ([]model.Role, error) {
	return s.userRepo.GetUserRoles(userID)
}

// UpdateUserRoles 更新用户角色
func (s *RBACService) UpdateUserRoles(userID uint, roleIDs []uint) error {
	// 获取当前角色
	currentRoles, err := s.userRepo.GetUserRoles(userID)
	if err != nil {
		return err
	}

	// 删除所有现有角色
	for _, role := range currentRoles {
		s.userRepo.RemoveRole(userID, role.ID)
	}

	// 添加新角色
	for _, roleID := range roleIDs {
		if err := s.userRepo.AssignRole(userID, roleID); err != nil {
			return err
		}
	}

	return nil
}
