package model

import (
	"time"

	"gorm.io/gorm"
)

// PushNotification 推送通知记录
type PushNotification struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	UserID      uint           `json:"user_id" gorm:"not null;index"`
	Title       string         `json:"title" gorm:"size:200;not null"`
	Content     string         `json:"content" gorm:"size:1000"`
	Type        string         `json:"type" gorm:"size:50"` // emergency, warning, reminder, info
	Data        string         `json:"data" gorm:"size:1000"` // JSON格式的附加数据
	Platform    string         `json:"platform" gorm:"size:20"` // wechat, android, ios
	Status      string         `json:"status" gorm:"size:20;default:'pending'"` // pending, sent, failed
	SentAt      *time.Time     `json:"sent_at"`
	FailedReason string        `json:"failed_reason" gorm:"size:500"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联
	User        User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// PushToken 设备推送Token
type PushToken struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	UserID      uint           `json:"user_id" gorm:"not null;index"`
	Token       string         `json:"token" gorm:"size:500;not null;index"`
	Platform    string         `json:"platform" gorm:"size:20;not null"` // wechat, android, ios
	DeviceID    string         `json:"device_id" gorm:"size:200"` // 设备唯一标识
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	LastUsedAt  *time.Time     `json:"last_used_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联
	User        User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName 指定表名
func (PushNotification) TableName() string {
	return "push_notifications"
}

func (PushToken) TableName() string {
	return "push_tokens"
}
