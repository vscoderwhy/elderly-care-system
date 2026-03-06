package repository

import (
	"elderly-care-system/internal/model"
	"gorm.io/gorm"
)

type BillRepository struct {
	db *gorm.DB
}

func NewBillRepository(db *gorm.DB) *BillRepository {
	return &BillRepository{db: db}
}

func (r *BillRepository) Create(bill *model.Bill) error {
	return r.db.Create(bill).Error
}

func (r *BillRepository) FindByID(id uint) (*model.Bill, error) {
	var bill model.Bill
	err := r.db.Preload("Elderly").Preload("Items").Preload("Payments").First(&bill, id).Error
	if err != nil {
		return nil, err
	}
	return &bill, nil
}

func (r *BillRepository) FindByElderlyID(elderlyID uint, offset, limit int) ([]model.Bill, int64, error) {
	var bills []model.Bill
	var total int64

	query := r.db.Where("elderly_id = ?", elderlyID)

	err := query.Model(&model.Bill{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Payments").Order("created_at DESC").Offset(offset).Limit(limit).Find(&bills).Error
	return bills, total, err
}

// List 获取所有账单（不按老人筛选）
func (r *BillRepository) List(offset, limit int) ([]model.Bill, int64, error) {
	var bills []model.Bill
	var total int64

	query := r.db.Where("deleted_at IS NULL")

	err := query.Model(&model.Bill{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Elderly").Preload("Payments").Order("created_at DESC").Offset(offset).Limit(limit).Find(&bills).Error
	return bills, total, err
}

func (r *BillRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&model.Bill{}).Where("id = ?", id).Update("status", status).Error
}

func (r *BillRepository) CreatePayment(payment *model.Payment) error {
	return r.db.Create(payment).Error
}

// CountPending 统计待缴费账单数
func (r *BillRepository) CountPending() (int64, error) {
	var count int64
	err := r.db.Model(&model.Bill{}).
		Where("status = ?", "unpaid").
		Count(&count).Error
	return count, err
}

// FinanceStatData 财务统计数据
type FinanceStatData struct {
	TotalIncome   float64
	MonthlyIncome float64
	PendingAmount float64
}

// GetFinanceStats 获取财务统计
func (r *BillRepository) GetFinanceStats() (*FinanceStatData, error) {
	stats := &FinanceStatData{}

	// 总收入（已支付金额）
	var totalIncome float64
	r.db.Model(&model.Payment{}).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalIncome)
	stats.TotalIncome = totalIncome

	// 本月收入
	var monthlyIncome float64
	r.db.Model(&model.Payment{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("EXTRACT(MONTH FROM paid_at) = EXTRACT(MONTH FROM CURRENT_DATE)").
		Where("EXTRACT(YEAR FROM paid_at) = EXTRACT(YEAR FROM CURRENT_DATE)").
		Scan(&monthlyIncome)
	stats.MonthlyIncome = monthlyIncome

	// 待收金额
	var pendingAmount float64
	r.db.Model(&model.Bill{}).
		Select("COALESCE(SUM(total_amount), 0)").
		Where("status = ?", "unpaid").
		Scan(&pendingAmount)
	stats.PendingAmount = pendingAmount

	return stats, nil
}
