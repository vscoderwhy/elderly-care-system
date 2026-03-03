package model

import (
	"time"

	"gorm.io/gorm"
)

// InventoryCategory 库存物资分类
type InventoryCategory struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Name        string         `json:"name" gorm:"size:100;not null"`
	Code        string         `json:"code" gorm:"size:50;not null;uniqueIndex"`
	Description string         `json:"description" gorm:"size:500"`
	ParentID    *uint          `json:"parent_id"`
	Sort        int            `json:"sort" gorm:"default:0"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联
	Parent      *InventoryCategory `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children    []InventoryCategory `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Items       []Inventory         `json:"items,omitempty" gorm:"foreignKey:CategoryID"`
}

// Inventory 库存物资
type Inventory struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	CategoryID  uint           `json:"category_id" gorm:"not null;index"`
	Name        string         `json:"name" gorm:"size:200;not null"`
	Code        string         `json:"code" gorm:"size:50;not null;uniqueIndex"`
	Spec        string         `json:"spec" gorm:"size:200"` // 规格
	Unit        string         `json:"unit" gorm:"size:20"` // 单位
	Quantity    float64        `json:"quantity" gorm:"default:0"` // 当前数量
	MinQuantity float64        `json:"min_quantity" gorm:"default:0"` // 最小库存
	MaxQuantity float64        `json:"max_quantity" gorm:"default:0"` // 最大库存
	CostPrice   float64        `json:"cost_price" gorm:"default:0"` // 成本价
	SellPrice   float64        `json:"sell_price" gorm:"default:0"` // 销售价
	ExpiryDays  int            `json:"expiry_days" gorm:"default:0"` // 有效期(天)
	Supplier    string         `json:"supplier" gorm:"size:200"` // 供应商
	Location    string         `json:"location" gorm:"size:100"` // 存放位置
	Status      string         `json:"status" gorm:"size:20;default:'normal'"` // normal, low_stock, out_of_stock, expired
	Remark      string         `json:"remark" gorm:"size:500"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联
	Category    InventoryCategory `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
}

// InventoryLog 库存变动记录
type InventoryLog struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	InventoryID uint           `json:"inventory_id" gorm:"not null;index"`
	Type        string         `json:"type" gorm:"size:20;not null"` // in, out, adjust, transfer
	Quantity    float64        `json:"quantity" gorm:"not null"`
	BeforeQty   float64        `json:"before_qty"`
	AfterQty    float64        `json:"after_qty"`
	Reason      string         `json:"reason" gorm:"size:500"`
	BatchNo     string         `json:"batch_no" gorm:"size:100"` // 批次号
	ExpiryDate  *time.Time     `json:"expiry_date"`
	CostPrice   float64        `json:"cost_price"`
	OperatorID  uint           `json:"operator_id" gorm:"not null"`
	RelatedID   *uint          `json:"related_id"` // 关联单据ID
	RelatedType string         `json:"related_type" gorm:"size:50"` // purchase, usage, transfer, etc.
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联
	Inventory   Inventory      `json:"inventory,omitempty" gorm:"foreignKey:InventoryID"`
	Operator    User           `json:"operator,omitempty" gorm:"foreignKey:OperatorID"`
}

// InventoryPurchase 采购单
type InventoryPurchase struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	No          string         `json:"no" gorm:"size:50;not null;uniqueIndex"` // 采购单号
	Supplier    string         `json:"supplier" gorm:"size:200;not null"`
	TotalAmount float64        `json:"total_amount" gorm:"default:0"`
	Status      string         `json:"status" gorm:"size:20;default:'pending'"` // pending, approved, received, cancelled
	Remark      string         `json:"remark" gorm:"size:500"`
	ApprovedBy  *uint          `json:"approved_by"`
	ApprovedAt  *time.Time     `json:"approved_at"`
	ReceivedBy  *uint          `json:"received_by"`
	ReceivedAt  *time.Time     `json:"received_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联
	Items       []InventoryPurchaseItem `json:"items,omitempty" gorm:"foreignKey:PurchaseID"`
}

// InventoryPurchaseItem 采购单明细
type InventoryPurchaseItem struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	PurchaseID  uint           `json:"purchase_id" gorm:"not null"`
	InventoryID uint           `json:"inventory_id" gorm:"not null"`
	Quantity    float64        `json:"quantity" gorm:"not null"`
	CostPrice   float64        `json:"cost_price" gorm:"default:0"`
	Amount      float64        `json:"amount" gorm:"default:0"`
	BatchNo     string         `json:"batch_no" gorm:"size:100"`
	ExpiryDate  *time.Time     `json:"expiry_date"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联
	Purchase    InventoryPurchase `json:"purchase,omitempty" gorm:"foreignKey:PurchaseID"`
	Inventory   Inventory         `json:"inventory,omitempty" gorm:"foreignKey:InventoryID"`
}

// TableName 指定表名
func (InventoryCategory) TableName() string {
	return "inventory_categories"
}

func (Inventory) TableName() string {
	return "inventories"
}

func (InventoryLog) TableName() string {
	return "inventory_logs"
}

func (InventoryPurchase) TableName() string {
	return "inventory_purchases"
}

func (InventoryPurchaseItem) TableName() string {
	return "inventory_purchase_items"
}
