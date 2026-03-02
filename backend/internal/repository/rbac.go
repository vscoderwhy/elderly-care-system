package repository

import (
	"elderly-care-system/internal/model"
	"gorm.io/gorm"
)

type PermissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

func (r *PermissionRepository) Create(permission *model.Permission) error {
	return r.db.Create(permission).Error
}

func (r *PermissionRepository) Update(permission *model.Permission) error {
	return r.db.Save(permission).Error
}

func (r *PermissionRepository) Delete(id uint) error {
	return r.db.Delete(&model.Permission{}, id).Error
}

func (r *PermissionRepository) GetByID(id uint) (*model.Permission, error) {
	var permission model.Permission
	err := r.db.First(&permission, id).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *PermissionRepository) List(page, pageSize int) ([]model.Permission, int64, error) {
	var permissions []model.Permission
	var total int64

	err := r.db.Model(&model.Permission{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = r.db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&permissions).Error

	return permissions, total, err
}

func (r *PermissionRepository) ListByRole(roleID uint) ([]model.Permission, error) {
	var permissions []model.Permission
	err := r.db.Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error
	return permissions, err
}

type MenuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) *MenuRepository {
	return &MenuRepository{db: db}
}

func (r *MenuRepository) Create(menu *model.Menu) error {
	return r.db.Create(menu).Error
}

func (r *MenuRepository) Update(menu *model.Menu) error {
	return r.db.Save(menu).Error
}

func (r *MenuRepository) Delete(id uint) error {
	return r.db.Delete(&model.Menu{}, id).Error
}

func (r *MenuRepository) GetByID(id uint) (*model.Menu, error) {
	var menu model.Menu
	err := r.db.First(&menu, id).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

func (r *MenuRepository) List() ([]model.Menu, error) {
	var menus []model.Menu
	err := r.db.Where("parent_id IS NULL").Order("sort ASC").Find(&menus).Error
	if err != nil {
		return nil, err
	}

	for i := range menus {
		r.loadChildren(&menus[i])
	}

	return menus, nil
}

func (r *MenuRepository) loadChildren(menu *model.Menu) error {
	return r.db.Where("parent_id = ?", menu.ID).Order("sort ASC").Find(&menu.Children).Error
}

func (r *MenuRepository) ListByRole(roleID uint) ([]model.Menu, error) {
	var menus []model.Menu
	err := r.db.Joins("JOIN role_menus ON role_menus.menu_id = menus.id").
		Where("role_menus.role_id = ? AND menus.status = ?", roleID, "active").
		Order("menus.sort ASC").
		Find(&menus).Error
	if err != nil {
		return nil, err
	}
	return menus, nil
}

func (r *MenuRepository) AssignPermissions(roleID uint, menuIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除现有关联
		tx.Where("role_id = ?", roleID).Delete(&model.RoleMenu{})

		// 添加新关联
		for _, menuID := range menuIDs {
			err := tx.Create(&model.RoleMenu{
				RoleID: roleID,
				MenuID: menuID,
			}).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) GetByID(id uint) (*model.Role, error) {
	var role model.Role
	err := r.db.Preload("Permissions").Preload("Menus").First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) List() ([]model.Role, error) {
	var roles []model.Role
	err := r.db.Order("id ASC").Find(&roles).Error
	return roles, err
}

func (r *RoleRepository) Update(role *model.Role) error {
	return r.db.Save(role).Error
}

func (r *RoleRepository) AssignPermissions(roleID uint, permissionIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除现有权限关联
		tx.Where("role_id = ?", roleID).Delete(&model.RolePermission{})

		// 添加新权限关联
		for _, permissionID := range permissionIDs {
			err := tx.Create(&model.RolePermission{
				RoleID:       roleID,
				PermissionID: permissionID,
			}).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *RoleRepository) GetUserPermissions(userID uint) ([]model.Permission, error) {
	var permissions []model.Permission
	err := r.db.Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ?", userID).
		Find(&permissions).Error
	return permissions, err
}

func (r *RoleRepository) GetUserMenus(userID uint) ([]model.Menu, error) {
	var menus []model.Menu
	err := r.db.Joins("JOIN role_menus ON role_menus.menu_id = menus.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_menus.role_id").
		Where("user_roles.user_id = ? AND menus.status = ?", userID, "active").
		Order("menus.sort ASC").
		Find(&menus).Error
	return menus, err
}
