package repository

import (
	"elderly-care-system/internal/model"
	"time"

	"gorm.io/gorm"
)

type InventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

// Category methods
func (r *InventoryRepository) CreateCategory(category *model.InventoryCategory) error {
	return r.db.Create(category).Error
}

func (r *InventoryRepository) GetCategory(id uint) (*model.InventoryCategory, error) {
	var category model.InventoryCategory
	err := r.db.Preload("Parent").First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *InventoryRepository) ListCategories() ([]model.InventoryCategory, error) {
	var categories []model.InventoryCategory
	err := r.db.Where("is_active = ?", true).Order("sort ASC, id ASC").Find(&categories).Error
	return categories, err
}

func (r *InventoryRepository) UpdateCategory(category *model.InventoryCategory) error {
	return r.db.Save(category).Error
}

func (r *InventoryRepository) DeleteCategory(id uint) error {
	return r.db.Delete(&model.InventoryCategory{}, id).Error
}

// Inventory methods
func (r *InventoryRepository) CreateInventory(inventory *model.Inventory) error {
	return r.db.Create(inventory).Error
}

func (r *InventoryRepository) GetInventory(id uint) (*model.Inventory, error) {
	var inventory model.Inventory
	err := r.db.Preload("Category").First(&inventory, id).Error
	if err != nil {
		return nil, err
	}
	return &inventory, nil
}

func (r *InventoryRepository) ListInventories(categoryID *uint, keyword string, status string) ([]model.Inventory, error) {
	var inventories []model.Inventory
	query := r.db.Preload("Category")
	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}
	if keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Order("id DESC").Find(&inventories).Error
	return inventories, err
}

func (r *InventoryRepository) UpdateInventory(inventory *model.Inventory) error {
	return r.db.Save(inventory).Error
}

func (r *InventoryRepository) DeleteInventory(id uint) error {
	return r.db.Delete(&model.Inventory{}, id).Error
}

func (r *InventoryRepository) UpdateQuantity(id uint, quantity float64, status string) error {
	return r.db.Model(&model.Inventory{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"quantity": quantity,
			"status":   status,
		}).Error
}

// GetLowStockInventories 获取库存不足的物资
func (r *InventoryRepository) GetLowStockInventories() ([]model.Inventory, error) {
	var inventories []model.Inventory
	err := r.db.Preload("Category").
		Where("quantity <= min_quantity").
		Where("min_quantity > 0").
		Order("quantity ASC").
		Find(&inventories).Error
	return inventories, err
}

// InventoryLog methods
func (r *InventoryRepository) CreateLog(log *model.InventoryLog) error {
	return r.db.Create(log).Error
}

func (r *InventoryRepository) GetInventoryLogs(inventoryID uint, logType string, limit int) ([]model.InventoryLog, error) {
	var logs []model.InventoryLog
	query := r.db.Preload("Operator").Where("inventory_id = ?", inventoryID)
	if logType != "" {
		query = query.Where("type = ?", logType)
	}
	err := query.Order("created_at DESC").Limit(limit).Find(&logs).Error
	return logs, err
}

// Purchase methods
func (r *InventoryRepository) CreatePurchase(purchase *model.InventoryPurchase) error {
	return r.db.Create(purchase).Error
}

func (r *InventoryRepository) GetPurchase(id uint) (*model.InventoryPurchase, error) {
	var purchase model.InventoryPurchase
	err := r.db.Preload("Items.Inventory").First(&purchase, id).Error
	if err != nil {
		return nil, err
	}
	return &purchase, nil
}

func (r *InventoryRepository) ListPurchases(status string, startDate, endDate time.Time) ([]model.InventoryPurchase, error) {
	var purchases []model.InventoryPurchase
	query := r.db.Preload("Items.Inventory")
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if !startDate.IsZero() {
		query = query.Where("created_at >= ?", startDate)
	}
	if !endDate.IsZero() {
		query = query.Where("created_at <= ?", endDate)
	}
	err := query.Order("created_at DESC").Find(&purchases).Error
	return purchases, err
}

func (r *InventoryRepository) UpdatePurchase(purchase *model.InventoryPurchase) error {
	return r.db.Save(purchase).Error
}

func (r *InventoryRepository) DeletePurchase(id uint) error {
	return r.db.Delete(&model.InventoryPurchase{}, id).Error
}

// GetInventoryStats 获取库存统计
func (r *InventoryRepository) GetInventoryStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 总数量
	var totalItems int64
	r.db.Model(&model.Inventory{}).Count(&totalItems)
	stats["total_items"] = totalItems

	// 低库存数量
	var lowStockCount int64
	r.db.Model(&model.Inventory{}).Where("quantity <= min_quantity").Where("min_quantity > 0").Count(&lowStockCount)
	stats["low_stock_count"] = lowStockCount

	// 缺货数量
	var outOfStockCount int64
	r.db.Model(&model.Inventory{}).Where("quantity = 0").Count(&outOfStockCount)
	stats["out_of_stock_count"] = outOfStockCount

	// 分类统计
	var categoryStats []struct {
		CategoryName string
		Count        int64
	}
	r.db.Model(&model.Inventory{}).
		Select("inventory_categories.name as category_name, count(*) as count").
		Joins("LEFT JOIN inventory_categories ON inventory_categories.id = inventories.category_id").
		Group("inventory_categories.name").
		Scan(&categoryStats)
	stats["category_stats"] = categoryStats

	return stats, nil
}
