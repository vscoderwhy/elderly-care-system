package model

import (
	"time"

	"gorm.io/gorm"
)

// VisitAppointment 探视预约
type VisitAppointment struct {
	ID           uint           `json:"id" gorm:"primarykey"`
	ElderlyID    uint           `json:"elderly_id" gorm:"not null;index"`
	Elderly      *Elderly       `json:"elderly,omitempty" gorm:"foreignKey:ElderlyID"`
	VisitorName  string         `json:"visitor_name" gorm:"size:100;not null"`
	VisitorPhone string         `json:"visitor_phone" gorm:"size:20;not null"`
	Relationship string         `json:"relationship" gorm:"size:50;not null"`
	VisitDate    time.Time      `json:"visit_date" gorm:"not null;index"`
	VisitTime    string         `json:"visit_time" gorm:"size:50;not null"` // 例如: "09:00-10:00"
	VisitorCount int            `json:"visitor_count" gorm:"default:1"`
	Status       string         `json:"status" gorm:"size:20;default:'pending'"` // pending, confirmed, completed, cancelled
	Notes        string         `json:"notes" gorm:"size:500"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 指定表名
func (VisitAppointment) TableName() string {
	return "visit_appointments"
}
