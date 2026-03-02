package repository

import (
	"elderly-care-system/internal/model"
	"gorm.io/gorm"
	"time"
)

type VisitRepository struct {
	db *gorm.DB
}

func NewVisitRepository(db *gorm.DB) *VisitRepository {
	return &VisitRepository{db: db}
}

// Create 创建探视预约
func (r *VisitRepository) Create(visit *model.VisitAppointment) error {
	return r.db.Create(visit).Error
}

// Update 更新探视预约
func (r *VisitRepository) Update(visit *model.VisitAppointment) error {
	return r.db.Save(visit).Error
}

// Delete 删除探视预约
func (r *VisitRepository) Delete(id uint) error {
	return r.db.Delete(&model.VisitAppointment{}, id).Error
}

// GetByID 根据ID获取探视预约
func (r *VisitRepository) GetByID(id uint) (*model.VisitAppointment, error) {
	var visit model.VisitAppointment
	err := r.db.Preload("Elderly").First(&visit, id).Error
	if err != nil {
		return nil, err
	}
	return &visit, nil
}

// List 获取探视预约列表
func (r *VisitRepository) List(page, pageSize int) ([]model.VisitAppointment, int64, error) {
	var visits []model.VisitAppointment
	var total int64

	err := r.db.Model(&model.VisitAppointment{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = r.db.Preload("Elderly").
		Order("visit_date ASC, visit_time ASC").
		Offset(offset).Limit(pageSize).
		Find(&visits).Error

	return visits, total, err
}

// ListByElderly 获取老人的探视预约列表
func (r *VisitRepository) ListByElderly(elderlyID uint, page, pageSize int) ([]model.VisitAppointment, int64, error) {
	var visits []model.VisitAppointment
	var total int64

	query := r.db.Model(&model.VisitAppointment{}).Where("elderly_id = ?", elderlyID)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Where("elderly_id = ?", elderlyID).
		Order("visit_date DESC").
		Offset(offset).Limit(pageSize).
		Find(&visits).Error

	return visits, total, err
}

// ListByDateRange 获取指定日期范围内的探视预约
func (r *VisitRepository) ListByDateRange(startDate, endDate time.Time) ([]model.VisitAppointment, error) {
	var visits []model.VisitAppointment
	err := r.db.Preload("Elderly").
		Where("visit_date >= ? AND visit_date <= ?", startDate, endDate).
		Order("visit_date ASC, visit_time ASC").
		Find(&visits).Error
	return visits, err
}

// ListByStatus 根据状态获取探视预约列表
func (r *VisitRepository) ListByStatus(status string) ([]model.VisitAppointment, error) {
	var visits []model.VisitAppointment
	err := r.db.Preload("Elderly").
		Where("status = ?", status).
		Order("visit_date ASC, visit_time ASC").
		Find(&visits).Error
	return visits, err
}

// GetTodayVisits 获取今日探视预约
func (r *VisitRepository) GetTodayVisits() ([]model.VisitAppointment, error) {
	var visits []model.VisitAppointment
	today := time.Now().Format("2006-01-02")
	err := r.db.Preload("Elderly").
		Where("DATE(visit_date) = ? AND status IN ?", today, []string{"pending", "confirmed"}).
		Order("visit_time ASC").
		Find(&visits).Error
	return visits, err
}

// GetUpcomingVisits 获取即将到来的探视预约
func (r *VisitRepository) GetUpcomingVisits(days int) ([]model.VisitAppointment, error) {
	var visits []model.VisitAppointment
	today := time.Now()
	endDate := today.AddDate(0, 0, days)

	err := r.db.Preload("Elderly").
		Where("visit_date >= ? AND visit_date <= ? AND status IN ?", today, endDate, []string{"pending", "confirmed"}).
		Order("visit_date ASC, visit_time ASC").
		Find(&visits).Error
	return visits, err
}
