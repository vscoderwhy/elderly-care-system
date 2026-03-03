package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户表
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Password  string         `json:"-" gorm:"size:255;not null"`
	Name      string         `json:"name" gorm:"size:50"`
	Avatar    string         `json:"avatar" gorm:"size:255"`
	Phone     string         `json:"phone" gorm:"size:20"`
	Email     string         `json:"email" gorm:"size:100"`
	Role      string         `json:"role" gorm:"size:20;not null;default:'user'"` // admin, nurse, family
	Status    string         `json:"status" gorm:"size:20;not null;default:'active'"`
	LastLoginAt *time.Time    `json:"lastLoginAt"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Elderly 老人表
type Elderly struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Name         string         `json:"name" gorm:"size:50;not null"`
	Gender       string         `json:"gender" gorm:"size:10;not null"` // 男, 女
	Age          int            `json:"age"`
	BirthDate    time.Time      `json:"birthDate"`
	IDCard       string         `json:"idCard" gorm:"size:18;uniqueIndex"`
	Avatar       string         `json:"avatar" gorm:"size:255"`
	CareLevel    string         `json:"careLevel" gorm:"size:20;not null"` // level1, level2, level3, special
	BedNumber    string         `json:"bedNumber" gorm:"size:50;uniqueIndex"`
	CheckInDate  time.Time      `json:"checkInDate"`
	HealthStatus string         `json:"healthStatus" gorm:"type:text"`
	Status       string         `json:"status" gorm:"size:20;not null;default:'active'"` // active, leave, hospital, discharged
	Remark       string         `json:"remark" gorm:"type:text"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// ElderlyFamily 老人家属表
type ElderlyFamily struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	ElderlyID uint           `json:"elderlyId" gorm:"not null;index"`
	Name      string         `json:"name" gorm:"size:50;not null"`
	Relation string         `json:"relation" gorm:"size:20;not null"` // 父子, 夫妻, 子女等
	Phone     string         `json:"phone" gorm:"size:20;not null"`
	Email     string         `json:"email" gorm:"size:100"`
	Address   string         `json:"address" gorm:"size:255"`
	IsPrimary bool           `json:"isPrimary" gorm:"not null;default:false"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// CareRecord 护理记录表
type CareRecord struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	ElderlyID   uint           `json:"elderlyId" gorm:"not null;index"`
	CareType    string         `json:"careType" gorm:"size:50;not null"` // 日常护理, 康复训练, 健康监测, 医疗护理
	Content     string         `json:"content" gorm:"type:text;not null"`
	Result      string         `json:"result" gorm:"type:text"`
	Evaluation  int            `json:"evaluation"` // 1-5星评价
	Images      string         `json:"images" gorm:"type:text"` // JSON数组
	NurseID     uint           `json:"nurseId" gorm:"not null"`
	NurseName   string         `json:"nurseName" gorm:"size:50"`
	RecordTime  time.Time      `json:"recordTime" gorm:"not null"`
	Duration    int            `json:"duration"` // 分钟
	Status      string         `json:"status" gorm:"size:20;not null;default:'completed'"` // pending, completed
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// CareTask 护理任务表
type CareTask struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	TaskNo      string         `json:"taskNo" gorm:"size:50;uniqueIndex;not null"`
	ElderlyID   uint           `json:"elderlyId" gorm:"not null;index"`
	TaskType    string         `json:"taskType" gorm:"size:50;not null"`
	Title       string         `json:"title" gorm:"size:200;not null"`
	Content     string         `json:"content" gorm:"type:text"`
	NurseID     uint           `json:"nurseId" gorm:"not null"`
	NurseName   string         `json:"nurseName" gorm:"size:50"`
	PlanTime    time.Time      `json:"planTime" gorm:"not null;index"`
	ExecuteTime *time.Time     `json:"executeTime"`
	Duration    int            `json:"duration"` // 预计时长(分钟)
	Priority    string         `json:"priority" gorm:"size:20;not null;default:'normal'"` // normal, important, urgent
	Progress    int            `json:"progress" gorm:"not null;default:0"` // 0-100
	Status      string         `json:"status" gorm:"size:20;not null;default:'pending'"` // pending, in_progress, completed, overdue
	Result      string         `json:"result" gorm:"type:text"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// HealthData 健康数据表
type HealthData struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	ElderlyID    uint           `json:"elderlyId" gorm:"not null;index"`
	RecordDate   time.Time      `json:"recordDate" gorm:"not null;index"`
	RecordType   string         `json:"recordType" gorm:"size:20;not null"` // bloodPressure, heartRate, bloodSugar, temperature, weight
	Value        string         `json:"value" gorm:"size:50;not null"`
	Unit         string         `json:"unit" gorm:"size:20"`
	Systolic     int            `json:"systolic"`  // 收缩压
	Diastolic    int            `json:"diastolic"` // 舒张压
	RecordBy     uint           `json:"recordBy"`
	RecordByName string         `json:"recordByName" gorm:"size:50"`
	Remark       string         `json:"remark" gorm:"type:text"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// Bill 账单表
type Bill struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	BillNo      string         `json:"billNo" gorm:"size:50;uniqueIndex;not null"`
	ElderlyID   uint           `json:"elderlyId" gorm:"not null;index"`
	BillType    string         `json:"billType" gorm:"size:50;not null"` // bed, care, meal, medical, other
	Amount      float64        `json:"amount" gorm:"not null"`
	BillDate    time.Time      `json:"billDate" gorm:"not null;index"`
	DueDate     time.Time      `json:"dueDate" gorm:"not null"`
	PaidDate    *time.Time     `json:"paidDate"`
	PaymentMethod string       `json:"paymentMethod" gorm:"size:20"` // wechat, alipay, cash, pos, transfer
	Status      string         `json:"status" gorm:"size:20;not null;default:'unpaid'"` // unpaid, paid, overdue
	Remark      string         `json:"remark" gorm:"type:text"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// VisitAppointment 探视预约表
type VisitAppointment struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	ElderlyID     uint           `json:"elderlyId" gorm:"not null;index"`
	VisitorName   string         `json:"visitorName" gorm:"size:50;not null"`
	Relation      string         `json:"relation" gorm:"size:20;not null"`
	Phone         string         `json:"phone" gorm:"size:20;not null"`
	VisitDate     time.Time      `json:"visitDate" gorm:"not null;index"`
	VisitTime     string         `json:"visitTime" gorm:"size:20;not null"` // 09:00-10:00
	VisitorCount  int            `json:"visitorCount" gorm:"not null;default:1"`
	Purpose       string         `json:"purpose" gorm:"type:text"`
	Status        string         `json:"status" gorm:"size:20;not null;default:'pending'"` // pending, approved, rejected, completed, cancelled
	Remark        string         `json:"remark" gorm:"type:text"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

// Notification 消息通知表
type Notification struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"userId" gorm:"not null;index"`
	Title     string         `json:"title" gorm:"size:200;not null"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	Type      string         `json:"type" gorm:"size:20;not null"` // care, health, bill, visit, system
	Data      string         `json:"data" gorm:"type:text"` // JSON数据
	IsRead    bool           `json:"isRead" gorm:"not null;default:false"`
	ReadAt    *time.Time     `json:"readAt"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// SystemLog 系统日志表
type SystemLog struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"userId"`
	Username  string         `json:"username" gorm:"size:50"`
	Module    string         `json:"module" gorm:"size:50;not null;index"`
	Action    string         `json:"action" gorm:"size:50;not null"`
	Method    string         `json:"method" gorm:"size:10"` // GET, POST, PUT, DELETE
	Path      string         `json:"path" gorm:"size:255"`
	IP        string         `json:"ip" gorm:"size:50"`
	Status    int            `json:"status"` // 状态码
	Error     string         `json:"error" gorm:"type:text"`
	Latency   int            `json:"latency"` // 毫秒
	CreatedAt time.Time      `json:"createdAt" gorm:"index"`
}
