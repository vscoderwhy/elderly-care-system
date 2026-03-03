package repository

import (
	"elderly-care-system/internal/model"
	"time"

	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

// CreatePaymentOrder 创建支付订单
func (r *PaymentRepository) CreatePaymentOrder(order *model.PaymentOrder) error {
	return r.db.Create(order).Error
}

// GetByOrderNo 根据订单号获取订单
func (r *PaymentRepository) GetByOrderNo(orderNo string) (*model.PaymentOrder, error) {
	var order model.PaymentOrder
	err := r.db.Where("order_no = ?", orderNo).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// UpdateOrderStatus 更新订单状态
func (r *PaymentRepository) UpdateOrderStatus(id uint, status string, transactionID string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if transactionID != "" {
		updates["transaction_id"] = transactionID
	}
	if status == "paid" {
		now := time.Now()
		updates["paid_at"] = &now
	}
	return r.db.Model(&model.PaymentOrder{}).Where("id = ?", id).Updates(updates).Error
}

// GetUserOrders 获取用户订单列表
func (r *PaymentRepository) GetUserOrders(userID uint, page, pageSize int) ([]model.PaymentOrder, int64, error) {
	var orders []model.PaymentOrder
	var total int64

	err := r.db.Model(&model.PaymentOrder{}).Where("user_id = ?", userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("Bill").Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&orders).Error

	return orders, total, err
}

// CreateRefund 创建退款记录
func (r *PaymentRepository) CreateRefund(refund *model.RefundRecord) error {
	return r.db.Create(refund).Error
}

// UpdateRefundStatus 更新退款状态
func (r *PaymentRepository) UpdateRefundStatus(id uint, status string, refundID string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if refundID != "" {
		updates["refund_id"] = refundID
	}
	if status == "success" {
		now := time.Now()
		updates["processed_at"] = &now
	}
	return r.db.Model(&model.RefundRecord{}).Where("id = ?", id).Updates(updates).Error
}

// GetPaymentByID 获取支付订单详情
func (r *PaymentRepository) GetPaymentByID(id uint) (*model.PaymentOrder, error) {
	var order model.PaymentOrder
	err := r.db.Preload("Bill").Preload("User").First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}
