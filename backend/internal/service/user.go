package service

import (
	"elderly-care-system/internal/model"
	"elderly-care-system/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetProfile(userID uint) (*model.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *UserService) UpdateProfile(userID uint, updates map[string]interface{}) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	// Update allowed fields only
	if nickname, ok := updates["nickname"]; ok {
		user.Nickname = nickname.(string)
	}
	if avatar, ok := updates["avatar"]; ok {
		user.Avatar = avatar.(string)
	}

	return s.userRepo.Update(user)
}

func (s *UserService) GetElderlyList(userID uint) ([]model.Elderly, error) {
	return s.userRepo.FindByFamilyUserID(userID)
}

// ListStaff 获取员工列表
func (s *UserService) ListStaff(role string, page, pageSize int) ([]model.User, int64, error) {
	return s.userRepo.ListStaff(role, page, pageSize)
}

// CreateStaff 创建员工
func (s *UserService) CreateStaff(req *CreateStaffRequest) (*model.User, error) {
	// 检查手机号是否已存在
	existing, _ := s.userRepo.FindByPhone(req.Phone)
	if existing != nil {
		return nil, ErrPhoneExists
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Phone:    req.Phone,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
		Status:   "active",
	}

	if err := s.userRepo.CreateStaff(user); err != nil {
		return nil, err
	}

	// 分配角色
	if req.RoleID > 0 {
		s.userRepo.AssignRole(user.ID, req.RoleID)
	}

	return user, nil
}

// UpdateStaff 更新员工
func (s *UserService) UpdateStaff(id uint, updates map[string]interface{}) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}

	if nickname, ok := updates["nickname"]; ok {
		user.Nickname = nickname.(string)
	}
	if avatar, ok := updates["avatar"]; ok {
		user.Avatar = avatar.(string)
	}
	if status, ok := updates["status"]; ok {
		user.Status = status.(string)
	}

	return s.userRepo.Update(user)
}

// DeleteStaff 删除员工
func (s *UserService) DeleteStaff(id uint) error {
	return s.userRepo.DeleteStaff(id)
}

type CreateStaffRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
	Avatar   string `json:"avatar"`
	RoleID   uint   `json:"role_id"`
}

var ErrPhoneExists = &PhoneExistsError{}

type PhoneExistsError struct{}

func (e *PhoneExistsError) Error() string {
	return "phone number already exists"
}
