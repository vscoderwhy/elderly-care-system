package repository

import (
	"elderly-care-system/internal/model"
	"gorm.io/gorm"
	"time"
)

type AlertRepository struct {
	db *gorm.DB
}

func NewAlertRepository(db *gorm.DB) *AlertRepository {
	return &AlertRepository{db: db}
}

// Alert CRUD
func (r *AlertRepository) Create(alert *model.Alert) error {
	return r.db.Create(alert).Error
}

func (r *AlertRepository) Update(alert *model.Alert) error {
	return r.db.Save(alert).Error
}

func (r *AlertRepository) Delete(id uint) error {
	return r.db.Delete(&model.Alert{}, id).Error
}

func (r *AlertRepository) GetByID(id uint) (*model.Alert, error) {
	var alert model.Alert
	err := r.db.First(&alert, id).Error
	if err != nil {
		return nil, err
	}
	return &alert, nil
}

func (r *AlertRepository) List(page, pageSize int, status string) ([]model.Alert, int64, error) {
	var alerts []model.Alert
	var total int64

	query := r.db.Model(&model.Alert{})
	if status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&alerts).Error

	return alerts, total, err
}

func (r *AlertRepository) ListByType(alertType string, page, pageSize int) ([]model.Alert, int64, error) {
	var alerts []model.Alert
	var total int64

	query := r.db.Model(&model.Alert{}).Where("type = ?", alertType)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&alerts).Error

	return alerts, total, err
}

func (r *AlertRepository) ListByLevel(level string, page, pageSize int) ([]model.Alert, int64, error) {
	var alerts []model.Alert
	var total int64

	query := r.db.Model(&model.Alert{}).Where("level = ?", level)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&alerts).Error

	return alerts, total, err
}

func (r *AlertRepository) GetActiveAlerts() ([]model.Alert, error) {
	var alerts []model.Alert
	err := r.db.Where("status = ?", "active").
		Order("level DESC, created_at DESC").
		Find(&alerts).Error
	return alerts, err
}

func (r *AlertRepository) GetAlertsByDateRange(startDate, endDate time.Time) ([]model.Alert, error) {
	var alerts []model.Alert
	err := r.db.Where("created_at >= ? AND created_at <= ?", startDate, endDate).
		Order("created_at DESC").
		Find(&alerts).Error
	return alerts, err
}

// AcknowledgeAlert 确认预警
func (r *AlertRepository) AcknowledgeAlert(id uint, userID uint) error {
	now := time.Now()
	return r.db.Model(&model.Alert{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":           "acknowledged",
		"acknowledged_by":  userID,
		"acknowledged_at":  &now,
	}).Error
}

// ResolveAlert 解决预警
func (r *AlertRepository) ResolveAlert(id uint) error {
	now := time.Now()
	return r.db.Model(&model.Alert{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":      "resolved",
		"resolved_at": &now,
	}).Error
}

// AlertRule CRUD
func (r *AlertRepository) CreateRule(rule *model.AlertRule) error {
	return r.db.Create(rule).Error
}

func (r *AlertRepository) UpdateRule(rule *model.AlertRule) error {
	return r.db.Save(rule).Error
}

func (r *AlertRepository) DeleteRule(id uint) error {
	return r.db.Delete(&model.AlertRule{}, id).Error
}

func (r *AlertRepository) GetRuleByID(id uint) (*model.AlertRule, error) {
	var rule model.AlertRule
	err := r.db.First(&rule, id).Error
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

func (r *AlertRepository) ListRules(page, pageSize int) ([]model.AlertRule, int64, error) {
	var rules []model.AlertRule
	var total int64

	err := r.db.Model(&model.AlertRule{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = r.db.Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&rules).Error

	return rules, total, err
}

func (r *AlertRepository) ListEnabledRules() ([]model.AlertRule, error) {
	var rules []model.AlertRule
	err := r.db.Where("is_enabled = ?", true).Find(&rules).Error
	return rules, err
}

func (r *AlertRepository) GetRulesByType(alertType string) ([]model.AlertRule, error) {
	var rules []model.AlertRule
	err := r.db.Where("type = ?", alertType).Find(&rules).Error
	return rules, err
}
