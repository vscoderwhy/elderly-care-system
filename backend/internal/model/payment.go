package model

import (
	"time"
	"gorm.io/gorm"
)

// PaymentOrder 支付订单
type PaymentOrder struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	OrderNo     string         `json:"order_no" gorm:"uniqueIndex;size:50;not null"`
	UserID      uint           `json:"user_id" gorm:"not null;index"`
	User        *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	BillID      *uint          `json:"bill_id" gorm:"index"`
	Bill        *Bill          `json:"bill,omitempty" gorm:"foreignKey:BillID"`
	Amount      float64        `json:"amount" gorm:"not null"` // 金额(分)
	PaymentMethod string         `json:"payment_method" gorm:"size:20"` // wechat/alipay
	Status      string         `json:"status" gorm:"size:20;default:'pending'"` // pending/paid/failed/refunded
	TransactionID string         `json:"transaction_id" gorm:"size:100"`
	PrepayID    string         `json:"prepay_id" gorm:"size:100"`
	NotifyTime  *time.Time     `json:"notify_time"`
	PaidAt      *time.Time     `json:"paid_at"`
	RefundedAt  *time.Time     `json:"refunded_at"`
	Remark      string         `json:"remark" gorm:"size:500"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// RefundRecord 退款记录
type RefundRecord struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	PaymentOrderID uint          `json:"payment_order_id" gorm:"not null;index"`
	PaymentOrder  *PaymentOrder  `json:"payment_order,omitempty" gorm:"foreignKey:PaymentOrderID"`
	RefundNo      string         `json:"refund_no" gorm:"uniqueIndex;size:50;not null"`
	Amount        float64        `json:"amount" gorm:"not null"` // 退款金额(分)
	Reason        string         `json:"reason" gorm:"size:500"`
	Status        string         `json:"status" gorm:"size:20;default:'pending'"` // pending/success/failed
	RefundID      string         `json:"refund_id" gorm:"size:100"`
	ProcessedAt   *time.Time     `json:"processed_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}
