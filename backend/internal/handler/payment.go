package handler

import (
	"elderly-care-system/internal/service"
	"elderly-care-system/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentService *service.PaymentService
}

func NewPaymentHandler(paymentService *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

// CreateOrder 创建支付订单
func (h *PaymentHandler) CreateOrder(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req service.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request")
		return
	}

	req.UserID = userID.(uint)

	order, err := h.paymentService.CreateOrder(&req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, order)
}

// GetPaymentParams 获取支付参数
func (h *PaymentHandler) GetPaymentParams(c *gin.Context) {
	orderNo := c.Query("order_no")
	if orderNo == "" {
		response.Error(c, 400, "Order number is required")
		return
	}

	params, err := h.paymentService.GetPaymentParams(orderNo)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, params)
}

// HandleNotify 处理支付通知(微信回调)
func (h *PaymentHandler) HandleNotify(c *gin.Context) {
	var notifyData map[string]interface{}
	if err := c.ShouldBindJSON(&notifyData); err != nil {
		c.JSON(200, gin.H{"code": "FAIL", "message": "Invalid request"})
		return
	}

	if err := h.paymentService.HandlePaymentNotify(notifyData); err != nil {
		c.JSON(200, gin.H{"code": "FAIL", "message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": "SUCCESS", "message": "OK"})
}

// Refund 申请退款
func (h *PaymentHandler) Refund(c *gin.Context) {
	var req service.RefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request")
		return
	}

	refund, err := h.paymentService.CreateRefund(&req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, refund)
}

// GetMyOrders 获取我的订单
func (h *PaymentHandler) GetMyOrders(c *gin.Context) {
	userID, _ := c.Get("user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	orders, total, err := h.paymentService.GetUserOrders(userID.(uint), page, pageSize)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":      orders,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
