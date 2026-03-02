package service

import (
	"elderly-care-system/internal/model"
	"elderly-care-system/internal/repository"
	"errors"
	"time"
)

type MedicationService struct {
	medicationRepo *repository.MedicationRepository
	elderlyRepo   *repository.ElderlyRepository
}

func NewMedicationService(
	medicationRepo *repository.MedicationRepository,
	elderlyRepo *repository.ElderlyRepository,
) *MedicationService {
	return &MedicationService{
		medicationRepo: medicationRepo,
		elderlyRepo:   elderlyRepo,
	}
}

// MedicationRequest 药品创建请求
type MedicationRequest struct {
	Name              string `json:"name" binding:"required"`
	GenericName       string `json:"generic_name"`
	Specification     string `json:"specification"`
	Unit             string `json:"unit" binding:"required"`
	Stock            int    `json:"stock" binding:"min=1"`
	MinStock         int    `json:"min_stock"`
	ExpiryDate       *time.Time `json:"expiry_date"`
	Manufacturer       string `json:"manufacturer"`
	UsageInstructions string `json:"usage_instructions"`
}

// MedicationRecordRequest 用药记录创建请求
type MedicationRecordRequest struct {
	ElderlyID    uint   `json:"elderly_id" binding:"required"`
	MedicationID uint   `json:"medication_id" binding:"required"`
	Dosage       string `json:"dosage" binding:"required"`
	Frequency    string `json:"frequency" binding:"required"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     *time.Time `json:"end_date"`
	Notes       string `json:"notes"`
}

// MedicationLogRequest 用药日志请求
type MedicationLogRequest struct {
	RecordID      uint   `json:"record_id" binding:"required"`
	StaffID       uint   `json:"staff_id" binding:"required"`
	ScheduledTime time.Time `json:"scheduled_time" binding:"required"`
	Notes         string `json:"notes"`
}

// CreateMedication 创建药品
func (s *MedicationService) CreateMedication(req *MedicationRequest) (*model.Medication, error) {
	medication := &model.Medication{
		Name:              req.Name,
		GenericName:       req.GenericName,
		Specification:     req.Specification,
		Unit:             req.Unit,
		Stock:            req.Stock,
		MinStock:         req.MinStock,
		ExpiryDate:       req.ExpiryDate,
		Manufacturer:     req.Manufacturer,
		UsageInstructions: req.UsageInstructions,
		Status:           "active",
	}

	if err := s.medicationRepo.Create(medication); err != nil {
		return nil, err
	}

	return medication, nil
}

// UpdateMedication 更新药品
func (s *MedicationService) UpdateMedication(id uint, updates map[string]interface{}) error {
	medication, err := s.medicationRepo.GetByID(id)
	if err != nil {
		return err
	}

	if name, ok := updates["name"]; ok {
		medication.Name = name.(string)
	}
	if stock, ok := updates["stock"]; ok {
		medication.Stock = stock.(int)
	}
	if minStock, ok := updates["min_stock"]; ok {
		medication.MinStock = minStock.(int)
	}

	return s.medicationRepo.Update(medication)
}

// ListMedications 获取药品列表
func (s *MedicationService) ListMedications(keyword string, page, pageSize int) ([]model.Medication, int64, error) {
	if keyword != "" {
		return s.medicationRepo.Search(keyword, page, pageSize)
	}
	return s.medicationRepo.List(page, pageSize)
}

// DeleteMedication 删除药品（软删除）
func (s *MedicationService) DeleteMedication(id uint) error {
	medication, err := s.medicationRepo.GetByID(id)
	if err != nil {
		return err
	}
	medication.Status = "inactive"
	return s.medicationRepo.Update(medication)
}

// CreateMedicationRecord 创建用药记录
func (s *MedicationService) CreateMedicationRecord(req *MedicationRecordRequest) (*model.MedicationRecord, error) {
	// 验证老人存在
	_, err := s.elderlyRepo.FindByID(req.ElderlyID)
	if err != nil {
		return nil, errors.New("老人不存在")
	}

	// 验证药品存在
	medication, err := s.medicationRepo.GetByID(req.MedicationID)
	if err != nil {
		return nil, errors.New("药品不存在")
	}

	record := &model.MedicationRecord{
		ElderlyID:    req.ElderlyID,
		MedicationID: req.MedicationID,
		Dosage:       req.Dosage,
		Frequency:    req.Frequency,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		Notes:        req.Notes,
		Status:       "active",
	}

	if err := s.medicationRepo.CreateRecord(record); err != nil {
		return nil, err
	}

	// 生成用药提醒日志
	s.generateMedicationLogs(record, medication)

	return record, nil
}

// generateMedicationLogs 生成用药提醒日志
func (s *MedicationService) generateMedicationLogs(record *model.MedicationRecord, medication *model.Medication) {
	// 根据频率生成用药日志
	// 例如：每天3次
	startDate := record.StartDate
	days := 30 // 生成未来30天的日志

	for i := 0; i < days; i++ {
		scheduledDate := startDate.AddDate(0, 0, i)
		// 每天多次用药的时间点
		times := []int{8, 12, 18} // 8:00, 12:00, 18:00
		for _, hour := range times {
			scheduledTime := time.Date(
				scheduledDate.Year(),
				scheduledDate.Month(),
				scheduledDate.Day(),
				hour, 0, 0, 0, time.Local,
			)
			log := &model.MedicationLog{
				RecordID:      record.ID,
				ScheduledTime: scheduledTime,
				Status:         "pending",
			}
			s.medicationRepo.CreateLog(log)
		}
	}
}

// ListElderlyMedications 获取老人用药列表
func (s *MedicationService) ListElderlyMedications(elderlyID uint, page, pageSize int) ([]model.MedicationRecord, int64, error) {
	return s.medicationRepo.ListRecordsByElderly(elderlyID, page, pageSize)
}

// GetTodayMedications 获取今日用药任务
func (s *MedicationService) GetTodayMedications(elderlyID uint) ([]model.MedicationLog, error) {
	return s.medicationRepo.GetTodayMedicationLogs(elderlyID)
}

// CompleteMedicationLog 完成用药
func (s *MedicationService) CompleteMedicationLog(logID uint, staffID uint, notes string) error {
	// 获取日志并更新状态
	// TODO: 实现完整的用药完成逻辑
	return nil
}

// GetMedicationAlerts 获取用药预警
func (s *MedicationService) GetMedicationAlerts() (map[string]interface{}, error) {
	// 获取即将过期的药品
	expiring, err := s.medicationRepo.GetExpiringMedications(30)
	if err != nil {
		return nil, err
	}

	// 获取库存不足的药品
	lowStock, err := s.medicationRepo.GetLowStockMedications()
	if err != nil {
		return nil, err
	}

	alerts := map[string]interface{}{
		"expiring_medications": expiring,
		"low_stock_medications": lowStock,
		"expiring_count":        len(expiring),
		"low_stock_count":      len(lowStock),
	}

	return alerts, nil
}
