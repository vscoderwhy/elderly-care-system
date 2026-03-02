package service

import (
	"elderly-care-system/internal/model"
	"testing"
)

// TestUserModel 测试用户模型
func TestUserModel(t *testing.T) {
	user := &model.User{
		Phone:    "13800138000",
		Nickname: "测试用户",
		Status:   "active",
	}

	if user.Phone != "13800138000" {
		t.Errorf("Expected phone 13800138000, got %s", user.Phone)
	}

	if user.Status != "active" {
		t.Errorf("Expected status active, got %s", user.Status)
	}
}

// TestRoleModel 测试角色模型
func TestRoleModel(t *testing.T) {
	role := &model.Role{
		Name:        "admin",
		Description: "管理员角色",
	}

	if role.Name != "admin" {
		t.Errorf("Expected role name admin, got %s", role.Name)
	}
}
