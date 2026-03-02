package repository

import (
	"elderly-care-system/internal/model"
	"gorm.io/gorm"
)

type CareRepository struct {
	db *gorm.DB
}

func NewCareRepository(db *gorm.DB) *CareRepository {
	return &CareRepository{db: db}
}

func (r *CareRepository) CreateRecord(record *model.CareRecord) error {
	return r.db.Create(record).Error
}

func (r *CareRepository) FindRecordsByElderlyID(elderlyID uint, offset, limit int) ([]model.CareRecord, int64, error) {
	var records []model.CareRecord
	var total int64

	query := r.db.Where("elderly_id = ?", elderlyID)

	err := query.Model(&model.CareRecord{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("CareItem").Preload("Staff").
		Order("recorded_at DESC").
		Offset(offset).Limit(limit).
		Find(&records).Error

	return records, total, err
}

func (r *CareRepository) FindMyTasks(staffID uint) ([]model.CareRecord, error) {
	var tasks []model.CareRecord
	err := r.db.Where("staff_id = ? AND DATE(recorded_at) = CURRENT_DATE", staffID).
		Preload("Elderly").
		Preload("CareItem").
		Find(&tasks).Error
	return tasks, err
}

func (r *CareRepository) ListCareItems() ([]model.CareItem, error) {
	var items []model.CareItem
	err := r.db.Find(&items).Error
	return items, err
}

func (r *CareRepository) CreateServiceRequest(req *model.ServiceRequest) error {
	return r.db.Create(req).Error
}

func (r *CareRepository) UpdateServiceRequest(id uint, status string, handlerID uint) error {
	return r.db.Model(&model.ServiceRequest{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       status,
			"handler_id":   handlerID,
			"completed_at": gorm.Expr("NOW()"),
		}).Error
}

// FindServiceRequests 查询服务请求列表
func (r *CareRepository) FindServiceRequests(status string, offset, limit int) ([]model.ServiceRequest, int64, error) {
	var requests []model.ServiceRequest
	var total int64

	query := r.db.Model(&model.ServiceRequest{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Elderly").Preload("Requester").Preload("Handler").
		Order("requested_at DESC").
		Offset(offset).Limit(limit).
		Find(&requests).Error

	return requests, total, err
}

// Health Record methods
func (r *CareRepository) CreateHealthRecord(record *model.HealthRecord) error {
	return r.db.Create(record).Error
}

func (r *CareRepository) FindHealthRecords(elderlyID uint, recordType string, offset, limit int) ([]model.HealthRecord, int64, error) {
	var records []model.HealthRecord
	var total int64

	query := r.db.Model(&model.HealthRecord{}).Where("elderly_id = ?", elderlyID)
	if recordType != "" {
		query = query.Where("record_type = ?", recordType)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Elderly").Preload("Recorder").
		Order("recorded_at DESC").
		Offset(offset).Limit(limit).
		Find(&records).Error

	return records, total, err
}

func (r *CareRepository) GetLatestHealthRecords(elderlyID uint) ([]model.HealthRecord, error) {
	var records []model.HealthRecord

	// 获取每种类型的最新一条记录
	recordTypes := []string{"blood_pressure", "blood_sugar", "temperature", "weight", "heart_rate"}

	for _, rt := range recordTypes {
		var record model.HealthRecord
		err := r.db.Where("elderly_id = ? AND record_type = ?", elderlyID, rt).
			Order("recorded_at DESC").
			First(&record).Error
		if err == nil {
			records = append(records, record)
		}
	}

	return records, nil
}

func (r *CareRepository) DeleteHealthRecord(id uint) error {
	return r.db.Delete(&model.HealthRecord{}, id).Error
}

// CountToday 统计今日护理次数
func (r *CareRepository) CountToday() (int64, error) {
	var count int64
	err := r.db.Model(&model.CareRecord{}).
		Where("DATE(recorded_at) = CURRENT_DATE").
		Count(&count).Error
	return count, err
}

// CountPendingServices 统计待处理服务请求数
func (r *CareRepository) CountPendingServices() (int64, error) {
	var count int64
	err := r.db.Model(&model.ServiceRequest{}).
		Where("status = ?", "pending").
		Count(&count).Error
	return count, err
}

// GetRecentRecords 获取最近护理记录
func (r *CareRepository) GetRecentRecords(limit int) ([]model.CareRecord, error) {
	var records []model.CareRecord
	err := r.db.Preload("Elderly").Preload("CareItem").Preload("Staff").
		Order("recorded_at DESC").
		Limit(limit).
		Find(&records).Error
	return records, err
}

// GetPendingServiceRequests 获取待处理的服务请求
func (r *CareRepository) GetPendingServiceRequests(limit int) ([]model.ServiceRequest, error) {
	var requests []model.ServiceRequest
	err := r.db.Preload("Elderly").
		Where("status = ?", "pending").
		Order("requested_at DESC").
		Limit(limit).
		Find(&requests).Error
	return requests, err
}

// CountByDate 统计指定日期的护理次数
func (r *CareRepository) CountByDate(date interface{}) (int64, error) {
	var count int64
	err := r.db.Model(&model.CareRecord{}).
		Where("DATE(recorded_at) = ?", date).
		Count(&count).Error
	return count, err
}
