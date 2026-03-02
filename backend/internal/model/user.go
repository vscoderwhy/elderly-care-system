package model

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Phone     string         `json:"phone" gorm:"uniqueIndex;size:20"`
	Password  string         `json:"-" gorm:"size:255"`
	Nickname  string         `json:"nickname" gorm:"size:50"`
	Avatar    string         `json:"avatar" gorm:"size:255"`
	OpenID    string         `json:"-" gorm:"size:100;index"`
	Status    string         `json:"status" gorm:"size:20;default:active"` // active/inactive
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Roles     []Role         `json:"roles" gorm:"many2many:user_roles;"`
}

type Role struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"uniqueIndex;size:50"`
	Description string    `json:"description" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at"`
}

type UserRole struct {
	UserID uint `json:"user_id"`
	RoleID uint `json:"role_id"`
}

func (Role) TableName() string {
	return "roles"
}
