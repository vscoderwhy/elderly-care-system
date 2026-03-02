package model

import (
	"time"
)

type CareItem struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:50;not null"`
	Category    string    `json:"category" gorm:"size:30"` // 喂饭/翻身/清洁/用药/其他
	Unit        string    `json:"unit" gorm:"size:20"`     // 次/分钟
	Description string    `json:"description" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at"`
}

type CareStandard struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ElderlyID uint      `json:"elderly_id" gorm:"not null"`
	CareItemID uint     `json:"care_item_id" gorm:"not null"`
	Frequency int       `json:"frequency" gorm:"default:1"` // 每天次数
	Priority  int       `json:"priority" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`
}

type CareRecord struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	ElderlyID uint           `json:"elderly_id" gorm:"not null"`
	CareItemID uint          `json:"care_item_id" gorm:"not null"`
	StaffID   uint           `json:"staff_id" gorm:"not null"`
	Status    string         `json:"status" gorm:"size:20;default:completed"`
	Notes     string         `json:"notes" gorm:"type:text"`
	Images    string         `json:"images" gorm:"type:text"` // 图片URL数组
	RecordedAt time.Time     `json:"recorded_at"`
	CreatedAt time.Time      `json:"created_at"`

	Elderly  *Elderly  `json:"elderly,omitempty" gorm:"foreignKey:ElderlyID"`
	CareItem *CareItem `json:"care_item,omitempty" gorm:"foreignKey:CareItemID"`
	Staff    *User     `json:"staff,omitempty" gorm:"foreignKey:StaffID"`
}

type ServiceRequest struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	ElderlyID   uint           `json:"elderly_id" gorm:"not null"`
	Type        string         `json:"type" gorm:"size:30"` // 护理/送餐/打扫/其他
	Status      string         `json:"status" gorm:"size:20;default:pending"` // pending/processing/completed
	RequesterID uint           `json:"requester_id" gorm:"not null"`
	HandlerID   *uint          `json:"handler_id"`
	Notes       string         `json:"notes" gorm:"type:text"`
	RequestedAt time.Time      `json:"requested_at"`
	CompletedAt *time.Time     `json:"completed_at"`
	CreatedAt   time.Time      `json:"created_at"`

	Elderly   *Elderly `json:"elderly,omitempty" gorm:"foreignKey:ElderlyID"`
	Requester *User    `json:"requester,omitempty" gorm:"foreignKey:RequesterID"`
	Handler   *User    `json:"handler,omitempty" gorm:"foreignKey:HandlerID"`
}

// HealthRecord 健康记录
type HealthRecord struct {
	ID           uint       `json:"id" gorm:"primaryKey"`
	ElderlyID    uint       `json:"elderly_id" gorm:"not null;index"`
	RecordType   string     `json:"record_type" gorm:"size:30"`  // blood_pressure/blood_sugar/temperature/weight/heart_rate
	Value        string     `json:"value" gorm:"size:50"`        // 数值
	Unit         string     `json:"unit" gorm:"size:20"`         // 单位
	Value2       string     `json:"value2" gorm:"size:50"`       // 第二个值（如血压低压）
	Notes        string     `json:"notes" gorm:"type:text"`      // 备注
	RecordedAt   time.Time  `json:"recorded_at"`                 // 记录时间
	RecordedBy   uint       `json:"recorded_by"`                 // 记录人ID
	CreatedAt    time.Time  `json:"created_at"`

	Elderly      *Elderly   `json:"elderly,omitempty" gorm:"foreignKey:ElderlyID"`
	Recorder     *User      `json:"recorder,omitempty" gorm:"foreignKey:RecordedBy"`
}
