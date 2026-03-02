package model

import (
	"time"

	"gorm.io/gorm"
)

// Alert 智能预警
type Alert struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Type        string         `json:"type" gorm:"size:50;not null;index"` // medication_low, medication_expiry, health_abnormal, bill_overdue, bed_available, etc.
	Title       string         `json:"title" gorm:"size:200;not null"`
	Content     string         `json:"content" gorm:"size:1000"`
	Level       string         `json:"level" gorm:"size:20;default:'info'"` // info, warning, critical
	EntityID    uint           `json:"entity_id"` // 关联实体的ID（如老人ID、药品ID等）
	EntityType  string         `json:"entity_type" gorm:"size:50"` // elderly, medication, bill, etc.
	Status      string         `json:"status" gorm:"size:20;default:'active'"` // active, acknowledged, resolved
	AcknowledgedBy *uint      `json:"acknowledged_by,omitempty"`
	AcknowledgedAt *time.Time `json:"acknowledged_at,omitempty"`
	ResolvedAt     *time.Time `json:"resolved_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// AlertRule 预警规则
type AlertRule struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Name        string         `json:"name" gorm:"size:100;not null"`
	Type        string         `json:"type" gorm:"size:50;not null"`
	Description string         `json:"description" gorm:"size:500"`
	Condition   string         `json:"condition" gorm:"size:500;not null"` // JSON格式的条件
	Threshold   string         `json:"threshold" gorm:"size:200"` // 阈值
	Level       string         `json:"level" gorm:"size:20;default:'warning'"`
	IsEnabled   bool           `json:"is_enabled" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 指定表名
func (Alert) TableName() string {
	return "alerts"
}

func (AlertRule) TableName() string {
	return "alert_rules"
}
