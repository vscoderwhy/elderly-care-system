package repository

import (
	"elderly-care-system/internal/model"
	"gorm.io/gorm"
	"time"
)

type ScheduleRepository struct {
	db *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

// Create 创建排班
func (r *ScheduleRepository) Create(schedule *model.Schedule) error {
	return r.db.Create(schedule).Error
}

// Update 更新排班
func (r *ScheduleRepository) Update(schedule *model.Schedule) error {
	return r.db.Save(schedule).Error
}

// Delete 删除排班
func (r *ScheduleRepository) Delete(id uint) error {
	return r.db.Delete(&model.Schedule{}, id).Error
}

// GetByID 根据ID获取排班
func (r *ScheduleRepository) GetByID(id uint) (*model.Schedule, error) {
	var schedule model.Schedule
	err := r.db.Preload("Staff").First(&schedule, id).Error
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

// ListByDate 根据日期范围获取排班列表
func (r *ScheduleRepository) ListByDate(startDate, endDate time.Time, page, pageSize int) ([]model.Schedule, int64, error) {
	var schedules []model.Schedule
	var total int64

	query := r.db.Model(&model.Schedule{}).Preload("Staff").
		Where("date >= ? AND date <= ?", startDate, endDate)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Order("date ASC").Offset(offset).Limit(pageSize).Find(&schedules).Error

	return schedules, total, err
}

// ListByStaff 根据员工获取排班列表
func (r *ScheduleRepository) ListByStaff(staffID uint, page, pageSize int) ([]model.Schedule, int64, error) {
	var schedules []model.Schedule
	var total int64

	query := r.db.Model(&model.Schedule{}).Where("staff_id = ?", staffID)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Order("date DESC").Offset(offset).Limit(pageSize).Find(&schedules).Error

	return schedules, total, err
}

// GetWeeklySchedule 获取周排班
func (r *ScheduleRepository) GetWeeklySchedule(staffID uint, weekStart time.Time) ([]model.Schedule, error) {
	var schedules []model.Schedule
	weekEnd := weekStart.AddDate(0, 0, 7)

	err := r.db.Where("staff_id = ? AND date >= ? AND date < ?", staffID, weekStart, weekEnd).
		Order("date ASC").
		Find(&schedules).Error

	return schedules, err
}
