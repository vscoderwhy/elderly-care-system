package service

import (
	"elderly-care-system/internal/model"
	"elderly-care-system/internal/repository"
	"fmt"
	"math/rand"
	"time"
)

type PaymentService struct {
	paymentRepo    *repository.PaymentRepository
	billRepo       *repository.BillRepository
}

func NewPaymentService(
	paymentRepo *repository.PaymentRepository,
	billRepo *repository.BillRepository,
) *PaymentService {
	return &PaymentService{
		paymentRepo: paymentRepo,
		billRepo:    billRepo,
	}
}

// CreatePaymentRequest 创建支付请求
type CreatePaymentRequest struct {
	UserID        uint    `json:"user_id" binding:"required"`
	BillID        *uint   `json:"bill_id"`
	Amount        float64 `json:"amount" binding:"required"`
	PaymentMethod string  `json:"payment_method" binding:"required"` // wechat/alipay
	Remark        string  `json:"remark"`
}

func (s *PaymentService) CreateOrder(req *CreatePaymentRequest) (*model.PaymentOrder, error) {
	// 生成订单号
	orderNo := s.generateOrderNo()

	// 转换为分
	amountInCents := int64(req.Amount * 100)

	order := &model.PaymentOrder{
		OrderNo:       orderNo,
		UserID:        req.UserID,
		BillID:        req.BillID,
		Amount:        float64(amountInCents),
		PaymentMethod: req.PaymentMethod,
		Status:        "pending",
		Remark:        req.Remark,
	}

	if err := s.paymentRepo.CreatePaymentOrder(order); err != nil {
		return nil, err
	}

	return order, nil
}

// generateOrderNo 生成订单号
func (s *PaymentService) generateOrderNo() string {
	timestamp := time.Now().Format("20060102150405")
	random := fmt.Sprintf("%04d", rand.Intn(10000))
	return fmt.Sprintf("PAY%s%s", timestamp, random)
}

// GetPaymentParams 获取支付参数(微信支付)
func (s *PaymentService) GetPaymentParams(orderNo string) (map[string]interface{}, error) {
	order, err := s.paymentRepo.GetByOrderNo(orderNo)
	if err != nil {
		return nil, err
	}

	// TODO: 实际对接微信支付API
	// 这里返回模拟参数
	params := map[string]interface{}{
		"appId":     "your_wechat_appid",
		"timeStamp": time.Now().Unix(),
		"nonceStr":  s.generateNonceStr(),
		"package":   fmt.Sprintf("prepay_id=wx%s", s.generatePrepayID()),
		"signType":  "RSA",
		"paySign":   "mock_sign",
	}

	// 保存prepay_id
	order.PrepayID = fmt.Sprintf("wx%s", s.generatePrepayID())
	s.paymentRepo.UpdateOrderStatus(order.ID, "pending", "")

	return params, nil
}

// generateNonceStr 生成随机字符串
func (s *PaymentService) generateNonceStr() string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 32)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// generatePrepayID 生成预支付ID
func (s *PaymentService) generatePrepayID() string {
	return fmt.Sprintf("%d%s%d", time.Now().Unix(), s.generateNonceStr()[:16], rand.Intn(1000))
}

// HandlePaymentNotify 处理支付通知
func (s *PaymentService) HandlePaymentNotify(notifyData map[string]interface{}) error {
	orderNo := notifyData["out_trade_no"].(string)
	transactionID := notifyData["transaction_id"].(string)

	order, err := s.paymentRepo.GetByOrderNo(orderNo)
	if err != nil {
		return err
	}

	// 更新订单状态
	if order.Status != "paid" {
		now := time.Now()
		order.Status = "paid"
		order.TransactionID = transactionID
		order.PaidAt = &now
		order.NotifyTime = &now

		// 更新账单状态
		if order.BillID != nil {
			// TODO: 更新账单为已支付
		}
	}

	return nil
}

// Refund 申请退款
type RefundRequest struct {
	PaymentOrderID uint    `json:"payment_order_id" binding:"required"`
	Amount         float64 `json:"amount" binding:"required"`
	Reason         string  `json:"reason"`
}

func (s *PaymentService) CreateRefund(req *RefundRequest) (*model.RefundRecord, error) {
	// 获取支付订单
	order, err := s.paymentRepo.GetPaymentByID(req.PaymentOrderID)
	if err != nil {
		return nil, err
	}

	if order.Status != "paid" {
		return nil, fmt.Errorf("订单未支付，无法退款")
	}

	// 检查退款金额
	refundAmount := int64(req.Amount * 100)
	if refundAmount > int64(order.Amount) {
		return nil, fmt.Errorf("退款金额不能超过支付金额")
	}

	// 生成退款单号
	refundNo := s.generateRefundNo()

	refund := &model.RefundRecord{
		PaymentOrderID: req.PaymentOrderID,
		RefundNo:       refundNo,
		Amount:        float64(refundAmount),
		Reason:        req.Reason,
		Status:        "pending",
	}

	if err := s.paymentRepo.CreateRefund(refund); err != nil {
		return nil, err
	}

	// TODO: 调用微信退款API
	// 模拟退款成功
	s.paymentRepo.UpdateRefundStatus(refund.ID, "success", "wx_refund_"+refundNo)

	return refund, nil
}

// generateRefundNo 生成退款单号
func (s *PaymentService) generateRefundNo() string {
	timestamp := time.Now().Format("20060102150405")
	random := fmt.Sprintf("%04d", rand.Intn(10000))
	return fmt.Sprintf("REF%s%s", timestamp, random)
}

// GetUserOrders 获取用户订单
func (s *PaymentService) GetUserOrders(userID uint, page, pageSize int) ([]model.PaymentOrder, int64, error) {
	return s.paymentRepo.GetUserOrders(userID, page, pageSize)
}
