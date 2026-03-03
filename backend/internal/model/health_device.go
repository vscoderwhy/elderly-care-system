package model

import (
	"time"
	"gorm.io/gorm"
)

// HealthDevice 智能健康设备
type HealthDevice struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:100;not null"`
	Type        string         `json:"type" gorm:"size:50;not null"` // blood_pressure_meter/blood_glucose_meter/thermometer/weight_scale/smart_band
	Brand       string         `json:"brand" gorm:"size:100"`
	Model       string         `json:"model" gorm:"size:100"`
	SerialNo    string         `json:"serial_no" gorm:"size:100;uniqueIndex"`
	MacAddress  string         `json:"mac_address" gorm:"size:50"`
	ElderlyID   *uint          `json:"elderly_id" gorm:"index"`
	Elderly     *Elderly       `json:"elderly,omitempty" gorm:"foreignKey:ElderlyID"`
	Status      string         `json:"status" gorm:"size:20;default:'active'"` // active/inactive
	LastConnect *time.Time     `json:"last_connect"`
	BatteryLevel int           `json:"battery_level"` // 0-100
	Settings    string         `json:"settings" gorm:"type:text"` // JSON格式的设备配置
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// DeviceData 设备数据记录
type DeviceData struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	DeviceID    uint           `json:"device_id" gorm:"not null;index"`
	Device      *HealthDevice  `json:"device,omitempty" gorm:"foreignKey:DeviceID"`
	ElderlyID   uint           `json:"elderly_id" gorm:"not null;index"`
	Elderly     *Elderly       `json:"elderly,omitempty" gorm:"foreignKey:ElderlyID"`
	DataType    string         `json:"data_type" gorm:"size:50;not null"` // blood_pressure/blood_sugar/temperature/weight/heart_rate/blood_oxygen
	DataValue   string         `json:"data_value" gorm:"type:text;not null"` // JSON格式的测量数据
	Unit        string         `json:"unit" gorm:"size:20"`
	MeasuredAt  time.Time      `json:"measured_at" gorm:"not null;index"`
	IsAbnormal  bool           `json:"is_abnormal" gorm:"default:false"`
	AlertLevel string         `json:"alert_level" gorm:"size:20"` // normal/warning/danger
	Notes       string         `json:"notes" gorm:"size:500"`
	SyncedAt   time.Time      `json:"synced_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// HealthAlertRule 健康预警规则
type HealthAlertRule struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:100;not null"`
	DataType    string         `json:"data_type" gorm:"size:50;not null"`
	MinValue    float64       `json:"min_value"`
	MaxValue    float64       `json:"max_value"`
	Condition   string         `json:"condition" gorm:"size:20"` // above/below/between
	AlertLevel  string         `json:"alert_level" gorm:"size:20;default:'warning'"` // warning/danger
	Action      string         `json:"action" gorm:"size:500"` // 触发后的操作
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
