package service

import (
	"elderly-care-system/internal/model"
	"elderly-care-system/internal/repository"
	"elderly-care-system/pkg/jwt"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	redis     *redis.Client
	jwtSecret string
}

func NewAuthService(userRepo *repository.UserRepository, redis *redis.Client, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		redis:     redis,
		jwtSecret: jwtSecret,
	}
}

type RegisterRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname"`
}

type LoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type WeChatLoginRequest struct {
	Code string `json:"code" binding:"required"`
}

func (s *AuthService) Register(req *RegisterRequest) (string, *model.User, error) {
	// Check if user exists
	_, err := s.userRepo.FindByPhone(req.Phone)
	if err == nil {
		return "", nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil, err
	}

	user := &model.User{
		Phone:    req.Phone,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
	}

	if err := s.userRepo.Create(user); err != nil {
		return "", nil, err
	}

	// Assign default role (family)
	s.userRepo.AssignRole(user.ID, 1) // Assuming role 1 is family

	token, err := jwt.GenerateToken(user.ID, user.Phone, s.jwtSecret, 24*7)
	return token, user, err
}

func (s *AuthService) Login(req *LoginRequest) (string, *model.User, error) {
	user, err := s.userRepo.FindByPhone(req.Phone)
	if err != nil {
		return "", nil, errors.New("user not found")
	}

	// 优先使用 PasswordHash，如果为空则使用 Password
	passwordHash := user.PasswordHash
	if passwordHash == "" {
		passwordHash = user.Password
	}

	// 调试日志
	fmt.Printf("DEBUG: Phone=%s, PasswordHash length=%d\n", user.Phone, len(passwordHash))

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password))
	if err != nil {
		fmt.Printf("DEBUG: Password comparison failed: %v\n", err)
		return "", nil, errors.New("invalid password")
	}

	token, err := jwt.GenerateToken(user.ID, user.Phone, s.jwtSecret, 24*7)
	return token, user, err
}

func (s *AuthService) WeChatLogin(req *WeChatLoginRequest) (string, *model.User, error) {
	// TODO: Implement WeChat login flow
	// 1. Exchange code for openid and session_key
	// 2. Find or create user by openid
	// 3. Generate JWT token

	return "", nil, errors.New("wechat login not implemented yet")
}

func (s *AuthService) Logout(token string) error {
	// TODO: Add token to Redis blacklist
	return nil
}
