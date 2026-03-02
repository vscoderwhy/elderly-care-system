package handler

import (
	"elderly-care-system/internal/service"
	"elderly-care-system/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BillHandler struct {
	billService *service.BillService
}

func NewBillHandler(billService *service.BillService) *BillHandler {
	return &BillHandler{billService: billService}
}

func (h *BillHandler) List(c *gin.Context) {
	elderlyID, _ := strconv.ParseUint(c.Query("elderly_id"), 10, 32)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	bills, total, err := h.billService.List(uint(elderlyID), page, pageSize)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":  bills,
		"total": total,
		"page":  page,
		"page_size": pageSize,
	})
}

func (h *BillHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	bill, err := h.billService.Get(uint(id))
	if err != nil {
		response.Error(c, 404, "Bill not found")
		return
	}

	response.Success(c, bill)
}

func (h *BillHandler) Pay(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req struct {
		Amount        float64 `json:"amount" binding:"required"`
		Method        string  `json:"method" binding:"required"`
		TransactionNo string  `json:"transaction_no"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request")
		return
	}

	if err := h.billService.Pay(uint(id), req.Amount, req.Method, req.TransactionNo); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "Payment successful"})
}
