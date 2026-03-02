package model

import (
	"time"
)

// Schedule 排班表
type Schedule struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	StaffID   uint          `json:"staff_id" gorm:"not null;index"`
	Date      time.Time     `json:"date" gorm:"type:date;not null"`
	ShiftType string        `json:"shift_type" gorm:"size:20"` // 早班/中班/晚班
	StartTime string        `json:"start_time" gorm:"size:10"`
	EndTime   string        `json:"end_time" gorm:"size:10"`
	Status    string        `json:"status" gorm:"size:20;default:scheduled"` // scheduled/completed/cancelled
	Notes     string        `json:"notes" gorm:"type:text"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`

	Staff *User `json:"staff,omitempty" gorm:"foreignKey:StaffID"`
}

// ScheduleShift 班次记录
type ScheduleShift struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	ScheduleID   uint     `json:"schedule_id" gorm:"not null;index"`
	StaffID     uint     `json:"staff_id" gorm:"not null"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Status      string   `json:"status" gorm:"size:20;default:completed"` // completed/late/absent
	Notes       string  `json:"notes" gorm:"type:text"`
	CreatedAt  time.Time `json:"created_at"`

	Staff *User `json:"staff,omitempty" gorm:"foreignKey:StaffID"`
}

// TableName 指定表名
func (Schedule) TableName() string {
	return "schedules"
}

func (ScheduleShift) TableName() string {
	return "schedule_shifts"
}
