package handler

import (
	"elderly-care-system/internal/model"
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

// ListAll 获取所有账单（新方法）
func (h *BillHandler) ListAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	bills, total, err := h.billService.ListAll(page, pageSize)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":      bills,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (h *BillHandler) List(c *gin.Context) {
	elderlyIDStr := c.Query("elderly_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var bills []model.Bill
	var total int64
	var err error

	if elderlyIDStr != "" && elderlyIDStr != "0" {
		// 按老人筛选
		elderlyID, _ := strconv.ParseUint(elderlyIDStr, 10, 32)
		bills, total, err = h.billService.List(uint(elderlyID), page, pageSize)
	} else {
		// 获取所有账单
		bills, total, err = h.billService.ListAll(page, pageSize)
	}

	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":      bills,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
