package handler

import (
	"elderly-care-system/internal/service"
	"elderly-care-system/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type InventoryHandler struct {
	inventoryService *service.InventoryService
}

func NewInventoryHandler(inventoryService *service.InventoryService) *InventoryHandler {
	return &InventoryHandler{inventoryService: inventoryService}
}

// Category handlers

// CreateCategory 创建分类
func (h *InventoryHandler) CreateCategory(c *gin.Context) {
	var req service.CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request")
		return
	}

	category, err := h.inventoryService.CreateCategory(&req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, category)
}

// ListCategories 获取分类列表
func (h *InventoryHandler) ListCategories(c *gin.Context) {
	categories, err := h.inventoryService.ListCategories()
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, categories)
}

// Inventory handlers

// CreateInventory 创建物资
func (h *InventoryHandler) CreateInventory(c *gin.Context) {
	var req service.InventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request")
		return
	}

	inventory, err := h.inventoryService.CreateInventory(&req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, inventory)
}

// ListInventories 获取物资列表
func (h *InventoryHandler) ListInventories(c *gin.Context) {
	var categoryID *uint
	if catID := c.Query("category_id"); catID != "" {
		id, _ := strconv.ParseUint(catID, 10, 32)
		eid := uint(id)
		categoryID = &eid
	}

	keyword := c.Query("keyword")
	status := c.Query("status")

	inventories, err := h.inventoryService.ListInventories(categoryID, keyword, status)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, inventories)
}

// GetInventory 获取物资详情
func (h *InventoryHandler) GetInventory(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	inventory, err := h.inventoryService.GetInventory(uint(id))
	if err != nil {
		response.Error(c, 404, "Inventory not found")
		return
	}

	response.Success(c, inventory)
}

// StockIn 入库
func (h *InventoryHandler) StockIn(c *gin.Context) {
	var req service.StockInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request")
		return
	}

	req.OperatorID = c.GetUint("user_id")

	if err := h.inventoryService.StockIn(&req); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "入库成功"})
}

// StockOut 出库
func (h *InventoryHandler) StockOut(c *gin.Context) {
	var req service.StockOutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request")
		return
	}

	req.OperatorID = c.GetUint("user_id")

	if err := h.inventoryService.StockOut(&req); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "出库成功"})
}

// AdjustInventory 调整库存 (盘点)
func (h *InventoryHandler) AdjustInventory(c *gin.Context) {
	var req service.AdjustRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request")
		return
	}

	req.OperatorID = c.GetUint("user_id")

	if err := h.inventoryService.AdjustInventory(&req); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "调整成功"})
}

// GetInventoryLogs 获取库存变动记录
func (h *InventoryHandler) GetInventoryLogs(c *gin.Context) {
	inventoryID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	logType := c.Query("type")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	logs, err := h.inventoryService.GetInventoryLogs(uint(inventoryID), logType, limit)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, logs)
}

// GetLowStockItems 获取低库存物资
func (h *InventoryHandler) GetLowStockItems(c *gin.Context) {
	items, err := h.inventoryService.GetLowStockItems()
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, items)
}

// Purchase handlers

// CreatePurchase 创建采购单
func (h *InventoryHandler) CreatePurchase(c *gin.Context) {
	var req service.PurchaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request")
		return
	}

	purchase, err := h.inventoryService.CreatePurchase(&req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, purchase)
}

// ListPurchases 获取采购单列表
func (h *InventoryHandler) ListPurchases(c *gin.Context) {
	status := c.Query("status")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	purchases, err := h.inventoryService.ListPurchases(status, startDate, endDate)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, purchases)
}

// GetPurchase 获取采购单详情
func (h *InventoryHandler) GetPurchase(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	// 获取采购单列表并查找
	purchases, err := h.inventoryService.ListPurchases("", "", "")
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	for _, p := range purchases {
		if p.ID == uint(id) {
			response.Success(c, p)
			return
		}
	}

	response.Error(c, 404, "Purchase not found")
}

// ApprovePurchase 审核采购单
func (h *InventoryHandler) ApprovePurchase(c *gin.Context) {
	var req service.ApprovePurchaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request")
		return
	}

	req.ApprovedBy = c.GetUint("user_id")

	if err := h.inventoryService.ApprovePurchase(&req); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "审核成功"})
}

// ReceivePurchase 采购收货
func (h *InventoryHandler) ReceivePurchase(c *gin.Context) {
	var req service.ReceivePurchaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request")
		return
	}

	req.OperatorID = c.GetUint("user_id")

	if err := h.inventoryService.ReceivePurchase(&req); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "收货成功"})
}

// GetStats 获取库存统计
func (h *InventoryHandler) GetStats(c *gin.Context) {
	stats, err := h.inventoryService.GetInventoryStats()
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, stats)
}
