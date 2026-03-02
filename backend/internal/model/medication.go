package model

import (
	"time"
)

// Medication 药品信息
type Medication struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	Name            string    `json:"name" gorm:"size:100;not null"`
	GenericName     string    `json:"generic_name" gorm:"size:100"` // 通用名
	Specification   string    `json:"specification" gorm:"size:50"` // 规格
	Unit            string    `json:"unit" gorm:"size:20"` // 单位(片/瓶/盒)
	Stock           int       `json:"stock" gorm:"default:0"` // 库存
	MinStock       int       `json:"min_stock" gorm:"default:10"` // 最低库存预警
	ExpiryDate     *time.Time `json:"expiry_date"` // 过期日期
	Manufacturer     string    `json:"manufacturer" gorm:"size:100"` // 生产厂家
	UsageInstructions string `json:"usage_instructions" gorm:"type:text"` // 用法说明
	Status         string    `json:"status" gorm:"size:20;default:active"` // active/inactive
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// MedicationRecord 用药记录
type MedicationRecord struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	ElderlyID     uint      `json:"elderly_id" gorm:"not null;index"`
	MedicationID  uint      `json:"medication_id" gorm:"not null;index"`
	Dosage        string    `json:"dosage" gorm:"size:50"` // 剂量
	Frequency    string    `json:"frequency" gorm:"size:50"` // 频率
	StartDate    time.Time `json:"start_date"` // 开始日期
	EndDate      *time.Time `json:"end_date"` // 结束日期
	Status        string    `json:"status" gorm:"size:20;default:active"` // active/completed/discontinued
	Notes        string    `json:"notes" gorm:"type:text"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Elderly    *Elderly    `json:"elderly,omitempty" gorm:"foreignKey:ElderlyID"`
	Medication *Medication `json:"medication,omitempty" gorm:"foreignKey:MedicationID"`
}

// MedicationLog 用药日志
type MedicationLog struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	RecordID        uint      `json:"record_id" gorm:"not null;index"` // 关联用药记录
	StaffID         uint      `json:"staff_id" gorm:"not null"` // 执行人
	ScheduledTime   time.Time `json:"scheduled_time"` // 计划时间
	ActualTime      *time.Time `json:"actual_time"` // 实际执行时间
	Status          string    `json:"status" gorm:"size:20;default:pending"` // pending/completed/missed
	Notes           string    `json:"notes" gorm:"type:text"`
	CreatedAt       time.Time `json:"created_at"`

	Staff *User `json:"staff,omitempty" gorm:"foreignKey:StaffID"`
}

// TableName 指定表名
func (Medication) TableName() string {
	return "medications"
}

func (MedicationRecord) TableName() string {
	return "medication_records"
}

func (MedicationLog) TableName() string {
	return "medication_logs"
}
