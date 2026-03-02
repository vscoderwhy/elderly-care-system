package repository

import (
	"elderly-care-system/internal/model"
	"gorm.io/gorm"
	"time"
)

type MedicationRepository struct {
	db *gorm.DB
}

func NewMedicationRepository(db *gorm.DB) *MedicationRepository {
	return &MedicationRepository{db: db}
}

// Medication CRUD
func (r *MedicationRepository) Create(medication *model.Medication) error {
	return r.db.Create(medication).Error
}

func (r *MedicationRepository) Update(medication *model.Medication) error {
	return r.db.Save(medication).Error
}

func (r *MedicationRepository) Delete(id uint) error {
	return r.db.Delete(&model.Medication{}, id).Error
}

func (r *MedicationRepository) GetByID(id uint) (*model.Medication, error) {
	var medication model.Medication
	err := r.db.First(&medication, id).Error
	if err != nil {
		return nil, err
	}
	return &medication, nil
}

func (r *MedicationRepository) List(page, pageSize int) ([]model.Medication, int64, error) {
	var medications []model.Medication
	var total int64

	err := r.db.Model(&model.Medication{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = r.db.Where("status = ?", "active").
		Order("name ASC").
		Offset(offset).Limit(pageSize).
		Find(&medications).Error

	return medications, total, err
}

func (r *MedicationRepository) Search(keyword string, page, pageSize int) ([]model.Medication, int64, error) {
	var medications []model.Medication
	var total int64

	query := r.db.Model(&model.Medication{}).Where("status = ?", "active")
	if keyword != "" {
		query = query.Where("name LIKE ? OR generic_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Order("name ASC").Offset(offset).Limit(pageSize).Find(&medications).Error

	return medications, total, err
}

// MedicationRecord CRUD
func (r *MedicationRepository) CreateRecord(record *model.MedicationRecord) error {
	return r.db.Create(record).Error
}

func (r *MedicationRepository) UpdateRecord(record *model.MedicationRecord) error {
	return r.db.Save(record).Error
}

func (r *MedicationRepository) DeleteRecord(id uint) error {
	return r.db.Delete(&model.MedicationRecord{}, id).Error
}

func (r *MedicationRepository) GetRecordByID(id uint) (*model.MedicationRecord, error) {
	var record model.MedicationRecord
	err := r.db.Preload("Elderly").Preload("Medication").First(&record, id).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *MedicationRepository) ListRecordsByElderly(elderlyID uint, page, pageSize int) ([]model.MedicationRecord, int64, error) {
	var records []model.MedicationRecord
	var total int64

	query := r.db.Model(&model.MedicationRecord{}).Where("elderly_id = ?", elderlyID)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Preload("Medication").Order("created_at DESC").
		Offset(offset).Limit(pageSize).Find(&records).Error

	return records, total, err
}

func (r *MedicationRepository) ListActiveRecords(elderlyID uint) ([]model.MedicationRecord, error) {
	var records []model.MedicationRecord
	err := r.db.Where("elderly_id = ? AND status = ?", elderlyID, "active").
		Preload("Medication").
		Order("created_at DESC").
		Find(&records).Error
	return records, err
}

// MedicationLog CRUD
func (r *MedicationRepository) CreateLog(log *model.MedicationLog) error {
	return r.db.Create(log).Error
}

func (r *MedicationRepository) UpdateLog(log *model.MedicationLog) error {
	return r.db.Save(log).Error
}

func (r *MedicationRepository) ListLogsByRecord(recordID uint, date time.Time) ([]model.MedicationLog, error) {
	var logs []model.MedicationLog
	query := r.db.Where("record_id = ?", recordID)
	if !date.IsZero() {
		query = query.Where("DATE(scheduled_time) = ?", date.Format("2006-01-02"))
	}
	err := query.Preload("Staff").Order("scheduled_time ASC").Find(&logs).Error
	return logs, err
}

func (r *MedicationRepository) GetTodayMedicationLogs(elderlyID uint) ([]model.MedicationLog, error) {
	var logs []model.MedicationLog
	today := time.Now().Format("2006-01-02")

	err := r.db.Table("medication_logs").
		Select("medication_logs.*").
		Joins("LEFT JOIN medication_records ON medication_logs.record_id = medication_records.id").
		Where("medication_records.elderly_id = ? AND DATE(medication_logs.scheduled_time) = ?", elderlyID, today).
		Preload("Staff").
		Order("medication_logs.scheduled_time ASC").
		Find(&logs).Error

	return logs, err
}

// GetExpiringMedications 获取即将过期的药品
func (r *MedicationRepository) GetExpiringMedications(days int) ([]model.Medication, error) {
	var medications []model.Medication
	expiryDate := time.Now().AddDate(0, 0, days)

	err := r.db.Where("status = ? AND expiry_date IS NOT NULL AND expiry_date <= ?", "active", expiryDate).
		Order("expiry_date ASC").
		Find(&medications).Error

	return medications, err
}

// GetLowStockMedications 获取库存不足的药品
func (r *MedicationRepository) GetLowStockMedications() ([]model.Medication, error) {
	var medications []model.Medication
	err := r.db.Where("status = ? AND stock <= min_stock", "active").
		Order("stock ASC").
		Find(&medications).Error
	return medications, err
}
