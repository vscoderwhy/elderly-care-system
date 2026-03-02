package model

import (
	"time"

	"gorm.io/gorm"
)

// Permission 权限
type Permission struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Code        string         `json:"code" gorm:"uniqueIndex;size:100;not null"` // 权限编码，如: user:create, user:delete
	Name        string         `json:"name" gorm:"size:100;not null"`             // 权限名称
	Description string         `json:"description" gorm:"type:text"`               // 权限描述
	Type        string         `json:"type" gorm:"size:20;default:'menu'"`         // menu/button/api
	Status      string         `json:"status" gorm:"size:20;default:'active'"`      // active/inactive
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// Menu 菜单
type Menu struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	ParentID   *uint          `json:"parent_id"`                     // 父菜单ID
	Name       string         `json:"name" gorm:"size:50;not null"`  // 菜单名称
	Path       string         `json:"path" gorm:"size:200"`           // 路由路径
	Component  string         `json:"component" gorm:"size:200"`     // 组件路径
	Icon       string         `json:"icon" gorm:"size:100"`           // 图标
	Sort       int            `json:"sort" gorm:"default:0"`          // 排序
	Type       string         `json:"type" gorm:"size:20;default:'menu'"` // menu/directory/button
	Permission string         `json:"permission" gorm:"size:100"`     // 权限标识
	Status     string         `json:"status" gorm:"size:20;default:'active'"` // active/inactive
	Visible    bool           `json:"visible" gorm:"default:true"`    // 是否可见
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	Children   []Menu         `json:"children,omitempty" gorm:"foreignKey:ParentID"`
}

// RolePermission 角色权限关联
type RolePermission struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	RoleID       uint      `json:"role_id" gorm:"not null"`
	PermissionID uint      `json:"permission_id" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at"`
}

// RoleMenu 角色菜单关联
type RoleMenu struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	RoleID    uint      `json:"role_id" gorm:"not null"`
	MenuID    uint      `json:"menu_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定表名
func (Permission) TableName() string {
	return "permissions"
}

func (Menu) TableName() string {
	return "menus"
}

func (RolePermission) TableName() string {
	return "role_permissions"
}

func (RoleMenu) TableName() string {
	return "role_menus"
}
