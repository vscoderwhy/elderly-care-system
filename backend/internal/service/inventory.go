package service

import (
	"elderly-care-system/internal/model"
	"elderly-care-system/internal/repository"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type InventoryService struct {
	inventoryRepo *repository.InventoryRepository
	userRepo      *repository.UserRepository
}

func NewInventoryService(
	inventoryRepo *repository.InventoryRepository,
	userRepo *repository.UserRepository,
) *InventoryService {
	return &InventoryService{
		inventoryRepo: inventoryRepo,
		userRepo:      userRepo,
	}
}

// CategoryRequest 分类请求
type CategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description"`
	ParentID    *uint  `json:"parent_id"`
	Sort        int    `json:"sort"`
}

// CreateCategory 创建分类
func (s *InventoryService) CreateCategory(req *CategoryRequest) (*model.InventoryCategory, error) {
	category := &model.InventoryCategory{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		ParentID:    req.ParentID,
		Sort:        req.Sort,
		IsActive:    true,
	}

	if err := s.inventoryRepo.CreateCategory(category); err != nil {
		return nil, err
	}

	return category, nil
}

// ListCategories 获取分类列表
func (s *InventoryService) ListCategories() ([]model.InventoryCategory, error) {
	return s.inventoryRepo.ListCategories()
}

// InventoryRequest 物资请求
type InventoryRequest struct {
	CategoryID  uint    `json:"category_id" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Code        string  `json:"code" binding:"required"`
	Spec        string  `json:"spec"`
	Unit        string  `json:"unit"`
	Quantity    float64 `json:"quantity"`
	MinQuantity float64 `json:"min_quantity"`
	MaxQuantity float64 `json:"max_quantity"`
	CostPrice   float64 `json:"cost_price"`
	SellPrice   float64 `json:"sell_price"`
	ExpiryDays  int     `json:"expiry_days"`
	Supplier    string  `json:"supplier"`
	Location    string  `json:"location"`
	Remark      string  `json:"remark"`
}

// CreateInventory 创建物资
func (s *InventoryService) CreateInventory(req *InventoryRequest) (*model.Inventory, error) {
	inventory := &model.Inventory{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Code:        req.Code,
		Spec:        req.Spec,
		Unit:        req.Unit,
		Quantity:    req.Quantity,
		MinQuantity: req.MinQuantity,
		MaxQuantity: req.MaxQuantity,
		CostPrice:   req.CostPrice,
		SellPrice:   req.SellPrice,
		ExpiryDays:  req.ExpiryDays,
		Supplier:    req.Supplier,
		Location:    req.Location,
		Remark:      req.Remark,
		Status:      s.calculateStatus(req.Quantity, req.MinQuantity),
	}

	if err := s.inventoryRepo.CreateInventory(inventory); err != nil {
		return nil, err
	}

	return inventory, nil
}

// calculateStatus 计算库存状态
func (s *InventoryService) calculateStatus(quantity, minQuantity float64) string {
	if quantity == 0 {
		return "out_of_stock"
	}
	if quantity <= minQuantity && minQuantity > 0 {
		return "low_stock"
	}
	return "normal"
}

// ListInventories 获取物资列表
func (s *InventoryService) ListInventories(categoryID *uint, keyword string, status string) ([]model.Inventory, error) {
	return s.inventoryRepo.ListInventories(categoryID, keyword, status)
}

// GetInventory 获取物资详情
func (s *InventoryService) GetInventory(id uint) (*model.Inventory, error) {
	return s.inventoryRepo.GetInventory(id)
}

// StockInRequest 入库请求
type StockInRequest struct {
	InventoryID uint    `json:"inventory_id" binding:"required"`
	Quantity    float64 `json:"quantity" binding:"required,gt=0"`
	Reason      string  `json:"reason"`
	BatchNo     string  `json:"batch_no"`
	ExpiryDate  string  `json:"expiry_date"`
	CostPrice   float64 `json:"cost_price"`
	OperatorID  uint    `json:"operator_id"`
}

// StockIn 入库
func (s *InventoryService) StockIn(req *StockInRequest) error {
	// 获取物资
	inventory, err := s.inventoryRepo.GetInventory(req.InventoryID)
	if err != nil {
		return err
	}

	// 解析过期日期
	var expiryDate *time.Time
	if req.ExpiryDate != "" {
		ed, err := time.Parse("2006-01-02", req.ExpiryDate)
		if err == nil {
			expiryDate = &ed
		}
	}

	// 记录变动
	log := &model.InventoryLog{
		InventoryID: req.InventoryID,
		Type:        "in",
		Quantity:    req.Quantity,
		BeforeQty:   inventory.Quantity,
		AfterQty:    inventory.Quantity + req.Quantity,
		Reason:      req.Reason,
		BatchNo:     req.BatchNo,
		ExpiryDate:  expiryDate,
		CostPrice:   req.CostPrice,
		OperatorID:  req.OperatorID,
	}

	// 更新库存
	newQuantity := inventory.Quantity + req.Quantity
	newStatus := s.calculateStatus(newQuantity, inventory.MinQuantity)

	if err := s.inventoryRepo.UpdateQuantity(req.InventoryID, newQuantity, newStatus); err != nil {
		return err
	}

	return s.inventoryRepo.CreateLog(log)
}

// StockOutRequest 出库请求
type StockOutRequest struct {
	InventoryID uint    `json:"inventory_id" binding:"required"`
	Quantity    float64 `json:"quantity" binding:"required,gt=0"`
	Reason      string  `json:"reason"`
	OperatorID  uint    `json:"operator_id"`
}

// StockOut 出库
func (s *InventoryService) StockOut(req *StockOutRequest) error {
	// 获取物资
	inventory, err := s.inventoryRepo.GetInventory(req.InventoryID)
	if err != nil {
		return err
	}

	// 检查库存
	if inventory.Quantity < req.Quantity {
		return fmt.Errorf("库存不足")
	}

	// 记录变动
	log := &model.InventoryLog{
		InventoryID: req.InventoryID,
		Type:        "out",
		Quantity:    req.Quantity,
		BeforeQty:   inventory.Quantity,
		AfterQty:    inventory.Quantity - req.Quantity,
		Reason:      req.Reason,
		OperatorID:  req.OperatorID,
	}

	// 更新库存
	newQuantity := inventory.Quantity - req.Quantity
	newStatus := s.calculateStatus(newQuantity, inventory.MinQuantity)

	if err := s.inventoryRepo.UpdateQuantity(req.InventoryID, newQuantity, newStatus); err != nil {
		return err
	}

	return s.inventoryRepo.CreateLog(log)
}

// AdjustRequest 调整库存请求
type AdjustRequest struct {
	InventoryID uint    `json:"inventory_id" binding:"required"`
	Quantity    float64 `json:"quantity" binding:"required"`
	Reason      string  `json:"reason" binding:"required"`
	OperatorID  uint    `json:"operator_id"`
}

// AdjustInventory 调整库存 (盘点)
func (s *InventoryService) AdjustInventory(req *AdjustRequest) error {
	// 获取物资
	inventory, err := s.inventoryRepo.GetInventory(req.InventoryID)
	if err != nil {
		return err
	}

	// 记录变动
	log := &model.InventoryLog{
		InventoryID: req.InventoryID,
		Type:        "adjust",
		Quantity:    req.Quantity - inventory.Quantity,
		BeforeQty:   inventory.Quantity,
		AfterQty:    req.Quantity,
		Reason:      req.Reason,
		OperatorID:  req.OperatorID,
	}

	// 更新库存
	newStatus := s.calculateStatus(req.Quantity, inventory.MinQuantity)

	if err := s.inventoryRepo.UpdateQuantity(req.InventoryID, req.Quantity, newStatus); err != nil {
		return err
	}

	return s.inventoryRepo.CreateLog(log)
}

// GetInventoryLogs 获取库存变动记录
func (s *InventoryService) GetInventoryLogs(inventoryID uint, logType string, limit int) ([]model.InventoryLog, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.inventoryRepo.GetInventoryLogs(inventoryID, logType, limit)
}

// GetLowStockItems 获取低库存物资
func (s *InventoryService) GetLowStockItems() ([]model.Inventory, error) {
	return s.inventoryRepo.GetLowStockInventories()
}

// PurchaseRequest 采购单请求
type PurchaseRequest struct {
	Supplier string                    `json:"supplier" binding:"required"`
	Items    []PurchaseItemRequest     `json:"items" binding:"required"`
	Remark   string                    `json:"remark"`
}

type PurchaseItemRequest struct {
	InventoryID uint    `json:"inventory_id" binding:"required"`
	Quantity    float64 `json:"quantity" binding:"required,gt=0"`
	CostPrice   float64 `json:"cost_price"`
	BatchNo     string  `json:"batch_no"`
	ExpiryDate  string  `json:"expiry_date"`
}

// CreatePurchase 创建采购单
func (s *InventoryService) CreatePurchase(req *PurchaseRequest) (*model.InventoryPurchase, error) {
	// 生成采购单号
	no := fmt.Sprintf("PO%s%s", time.Now().Format("20060102"), uuid.New().String()[:8])

	purchase := &model.InventoryPurchase{
		No:     no,
		Supplier: req.Supplier,
		Status: "pending",
		Remark: req.Remark,
	}

	var totalAmount float64
	for _, item := range req.Items {
		amount := item.Quantity * item.CostPrice
		totalAmount += amount
		purchase.Items = append(purchase.Items, model.InventoryPurchaseItem{
			InventoryID: item.InventoryID,
			Quantity:    item.Quantity,
			CostPrice:   item.CostPrice,
			Amount:      amount,
			BatchNo:     item.BatchNo,
		})
	}
	purchase.TotalAmount = totalAmount

	if err := s.inventoryRepo.CreatePurchase(purchase); err != nil {
		return nil, err
	}

	return purchase, nil
}

// ListPurchases 获取采购单列表
func (s *InventoryService) ListPurchases(status string, startDate, endDate string) ([]model.InventoryPurchase, error) {
	var start, end time.Time
	var err error
	if startDate != "" {
		start, err = time.Parse("2006-01-02", startDate)
	}
	if endDate != "" {
		end, err = time.Parse("2006-01-02", endDate)
	}
	if err != nil {
		return nil, err
	}
	return s.inventoryRepo.ListPurchases(status, start, end)
}

// ReceivePurchaseRequest 收货请求
type ReceivePurchaseRequest struct {
	PurchaseID uint `json:"purchase_id" binding:"required"`
	OperatorID uint `json:"operator_id" binding:"required"`
}

// ReceivePurchase 采购收货
func (s *InventoryService) ReceivePurchase(req *ReceivePurchaseRequest) error {
	purchase, err := s.inventoryRepo.GetPurchase(req.PurchaseID)
	if err != nil {
		return err
	}

	if purchase.Status != "approved" {
		return fmt.Errorf("采购单未审核")
	}

	now := time.Now()
	purchase.Status = "received"
	purchase.ReceivedBy = &req.OperatorID
	purchase.ReceivedAt = &now

	// 入库操作
	for _, item := range purchase.Items {
		stockInReq := &StockInRequest{
			InventoryID: item.InventoryID,
			Quantity:    item.Quantity,
			Reason:      fmt.Sprintf("采购收货: %s", purchase.No),
			BatchNo:     item.BatchNo,
			CostPrice:   item.CostPrice,
			OperatorID:  req.OperatorID,
		}
		if err := s.StockIn(stockInReq); err != nil {
			return err
		}
	}

	return s.inventoryRepo.UpdatePurchase(purchase)
}

// ApprovePurchaseRequest 审核请求
type ApprovePurchaseRequest struct {
	PurchaseID uint `json:"purchase_id" binding:"required"`
	ApprovedBy uint `json:"approved_by" binding:"required"`
}

// ApprovePurchase 审核采购单
func (s *InventoryService) ApprovePurchase(req *ApprovePurchaseRequest) error {
	purchase, err := s.inventoryRepo.GetPurchase(req.PurchaseID)
	if err != nil {
		return err
	}

	if purchase.Status != "pending" {
		return fmt.Errorf("采购单状态不正确")
	}

	now := time.Now()
	purchase.Status = "approved"
	purchase.ApprovedBy = &req.ApprovedBy
	purchase.ApprovedAt = &now

	return s.inventoryRepo.UpdatePurchase(purchase)
}

// GetInventoryStats 获取库存统计
func (s *InventoryService) GetInventoryStats() (map[string]interface{}, error) {
	return s.inventoryRepo.GetInventoryStats()
}
