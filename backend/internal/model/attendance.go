package model

import (
	"time"
	"gorm.io/gorm"
)

// Attendance 考勤记录
type Attendance struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	StaffID   uint           `json:"staff_id" gorm:"not null;index"`
	Staff     *User          `json:"staff,omitempty" gorm:"foreignKey:StaffID"`
	Date      time.Time      `json:"date" gorm:"type:date;not null;index"`
	Type      string         `json:"type" gorm:"size:20;not null"` // clock_in/clock_out
	Time      time.Time      `json:"time" gorm:"not null"`
	Latitude  float64        `json:"latitude"`
	Longitude float64        `json:"longitude"`
	Location  string         `json:"location" gorm:"size:200"`
	Device    string         `json:"device" gorm:"size:100"`
	Status    string         `json:"status" gorm:"size:20;default:'normal'"` // normal/late/early/overtime
	Notes     string         `json:"notes" gorm:"size:500"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Performance 绩效记录
type Performance struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	StaffID    uint           `json:"staff_id" gorm:"not null;index"`
	Staff      *User          `json:"staff,omitempty" gorm:"foreignKey:StaffID"`
	Year       int            `json:"year" gorm:"not null;index"`
	Month      int            `json:"month" gorm:"not null;index"`
	// 服务质量评分
	ServiceScore     float64 `json:"service_score" gorm:"default:0"`      // 服务质量 0-100
	ElderlySatisfaction float64 `json:"elderly_satisfaction" gorm:"default:0"` // 老人满意度 0-100
	// 工作量统计
	CareRecordsCount  int     `json:"care_records_count" gorm:"default:0"`   // 护理记录数
	ServiceRequests   int     `json:"service_requests" gorm:"default:0"`     // 服务请求数
	AttendanceRate    float64 `json:"attendance_rate" gorm:"default:100"`    // 出勤率
	OvertimeHours     float64 `json:"overtime_hours" gorm:"default:0"`       // 加班时长
	// 奖惩记录
	RewardAmount      float64 `json:"reward_amount" gorm:"default:0"`       // 奖励金额
	PenaltyAmount     float64 `json:"penalty_amount" gorm:"default:0"`       // 扣罚金额
	// 综合评分
	TotalScore        float64 `json:"total_score" gorm:"default:0"`         // 综合评分
	Rank              int     `json:"rank"`                                  // 排名
	Notes             string  `json:"notes" gorm:"size:500"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`
}

// Salary 工资记录
type Salary struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	StaffID      uint           `json:"staff_id" gorm:"not null;index"`
	Staff        *User          `json:"staff,omitempty" gorm:"foreignKey:StaffID"`
	Year         int            `json:"year" gorm:"not null;index"`
	Month        int            `json:"month" gorm:"not null;index"`
	// 基本工资
	BaseSalary   float64        `json:"base_salary"`
	// 绩效工资
	PerformanceSalary float64   `json:"performance_salary"`
	// 加班费
	OvertimeSalary float64      `json:"overtime_salary"`
	// 奖金
	Reward        float64        `json:"reward"`
	// 扣款
	Penalty       float64        `json:"penalty"`
	// 社保扣除
	SocialInsurance float64      `json:"social_insurance"`
	// 应发工资
	GrossSalary   float64        `json:"gross_salary"`
	// 实发工资
	NetSalary     float64        `json:"net_salary"`
	// 状态
	Status        string         `json:"status" gorm:"size:20;default:'pending'"` // pending/paid
	PaidAt        *time.Time     `json:"paid_at"`
	PaidBy        uint           `json:"paid_by"`
	Notes         string         `json:"notes" gorm:"size:500"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

// ShiftRule 排班规则
type ShiftRule struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:50;not null"`
	Type        string         `json:"type" gorm:"size:20;not null"` // day/night/rotation
	StartTime   string         `json:"start_time" gorm:"size:10;not null"` // HH:MM
	EndTime     string         `json:"end_time" gorm:"size:10;not null"`   // HH:MM
	BreakTime   int            `json:"break_time"` // 休息时长(分钟)
	WorkHours   float64        `json:"work_hours"` // 工作时长
	Description string         `json:"description" gorm:"size:200"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	SortOrder   int            `json:"sort_order" gorm:"default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
