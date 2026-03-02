package repository

import (
	"elderly-care-system/internal/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByPhone(phone string) (*model.User, error) {
	var user model.User
	err := r.db.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Roles").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) FindByOpenID(openid string) (*model.User, error) {
	var user model.User
	err := r.db.Where("open_id = ?", openid).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) AssignRole(userID uint, roleID uint) error {
	return r.db.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?) ON CONFLICT DO NOTHING", userID, roleID).Error
}

// FindByFamilyUserID 获取家属关联的老人列表
func (r *UserRepository) FindByFamilyUserID(userID uint) ([]model.Elderly, error) {
	var elderlyList []model.Elderly
	err := r.db.Table("elderly").
		Joins("JOIN elderly_family ON elderly.id = elderly_family.elderly_id").
		Where("elderly_family.user_id = ?", userID).
		Find(&elderlyList).Error
	return elderlyList, err
}

// ListStaff 获取员工列表
func (r *UserRepository) ListStaff(role string, page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := r.db.Model(&model.User{}).Preload("Roles")
	if role != "" {
		query = query.Joins("JOIN user_roles ON users.id = user_roles.user_id").
			Joins("JOIN roles ON user_roles.role_id = roles.id").
			Where("roles.name = ?", role)
	} else {
		// 只查询员工角色（非家属）
		query = query.Joins("JOIN user_roles ON users.id = user_roles.user_id").
			Joins("JOIN roles ON user_roles.role_id = roles.id").
			Where("roles.name != ?", "family")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Distinct().Offset(offset).Limit(pageSize).Find(&users).Error
	return users, total, err
}

// CreateStaff 创建员工
func (r *UserRepository) CreateStaff(user *model.User) error {
	return r.db.Create(user).Error
}

// DeleteStaff 删除员工
func (r *UserRepository) DeleteStaff(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

// CountStaff 统计在职员工数
func (r *UserRepository) CountStaff() (int64, error) {
	var count int64
	err := r.db.Model(&model.User{}).
		Joins("JOIN user_roles ON users.id = user_roles.user_id").
		Joins("JOIN roles ON user_roles.role_id = roles.id").
		Where("roles.name != ?", "family").
		Where("users.status = ?", "active").
		Count(&count).Error
	return count, err
}

// RemoveRole 移除用户角色
func (r *UserRepository) RemoveRole(userID uint, roleID uint) error {
	return r.db.Exec("DELETE FROM user_roles WHERE user_id = ? AND role_id = ?", userID, roleID).Error
}

// GetUserRoles 获取用户的所有角色
func (r *UserRepository) GetUserRoles(userID uint) ([]model.Role, error) {
	var roles []model.Role
	err := r.db.Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles).Error
	return roles, err
}
