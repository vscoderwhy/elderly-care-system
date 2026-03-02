package model

import (
	"time"
	"gorm.io/gorm"
)

type FeeItem struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:50;not null"`
	Amount      float64   `json:"amount" gorm:"type:decimal(10,2);not null"`
	Type        string    `json:"type" gorm:"size:30"` // 床位费/护理费/餐费/其他
	Cycle       string    `json:"cycle" gorm:"size:20"` // daily/monthly/once
	Description string    `json:"description" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at"`
}

type Bill struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	ElderlyID       uint           `json:"elderly_id" gorm:"not null"`
	BillNo          string         `json:"bill_no" gorm:"uniqueIndex;size:50"`
	TotalAmount     float64        `json:"total_amount" gorm:"type:decimal(10,2)"`
	Status          string         `json:"status" gorm:"size:20;default:unpaid"` // unpaid/paid/overdue
	BillPeriodStart *time.Time     `json:"bill_period_start"`
	BillPeriodEnd   *time.Time     `json:"bill_period_end"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`

	Elderly   *Elderly    `json:"elderly,omitempty" gorm:"foreignKey:ElderlyID"`
	Items     []BillItem  `json:"items,omitempty" gorm:"foreignKey:BillID"`
	Payments  []Payment   `json:"payments,omitempty" gorm:"foreignKey:BillID"`
}

type BillItem struct {
	ID          uint    `json:"id" gorm:"primaryKey"`
	BillID      uint    `json:"bill_id" gorm:"not null"`
	FeeItemID   uint    `json:"fee_item_id" gorm:"not null"`
	Quantity    int     `json:"quantity" gorm:"default:1"`
	Amount      float64 `json:"amount" gorm:"type:decimal(10,2)"`
	Description string  `json:"description" gorm:"type:text"`
}

type Payment struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	BillID        uint       `json:"bill_id" gorm:"not null"`
	Amount        float64    `json:"amount" gorm:"type:decimal(10,2)"`
	Method        string     `json:"method" gorm:"size:30"` // wechat/cash/transfer
	TransactionNo string     `json:"transaction_no" gorm:"size:100"`
	PaidAt        time.Time  `json:"paid_at"`
}
