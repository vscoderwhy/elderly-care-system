package service

import (
	"elderly-care-system/internal/model"
	"elderly-care-system/internal/repository"
	"time"
)

type CareService struct {
	careRepo *repository.CareRepository
}

func NewCareService(careRepo *repository.CareRepository) *CareService {
	return &CareService{careRepo: careRepo}
}

func (s *CareService) ListRecords(elderlyID uint, page, pageSize int) ([]model.CareRecord, int64, error) {
	offset := (page - 1) * pageSize
	return s.careRepo.FindRecordsByElderlyID(elderlyID, offset, pageSize)
}

func (s *CareService) CreateRecord(req *CreateCareRecordRequest, staffID uint) (*model.CareRecord, error) {
	record := &model.CareRecord{
		ElderlyID:  req.ElderlyID,
		CareItemID: req.CareItemID,
		StaffID:    staffID,
		Status:     "completed",
		Notes:      req.Notes,
		Images:     req.Images,
		RecordedAt: time.Now(),
	}

	if err := s.careRepo.CreateRecord(record); err != nil {
		return nil, err
	}

	return record, nil
}

func (s *CareService) GetMyTasks(staffID uint) ([]model.CareRecord, error) {
	return s.careRepo.FindMyTasks(staffID)
}

func (s *CareService) ListCareItems() ([]model.CareItem, error) {
	return s.careRepo.ListCareItems()
}

type CreateCareRecordRequest struct {
	ElderlyID  uint    `json:"elderly_id" binding:"required"`
	CareItemID uint    `json:"care_item_id" binding:"required"`
	Notes      string  `json:"notes"`
	Images     string  `json:"images"`
}

// ListServiceRequests 获取服务请求列表
func (s *CareService) ListServiceRequests(status string, page, pageSize int) ([]model.ServiceRequest, int64, error) {
	offset := (page - 1) * pageSize
	return s.careRepo.FindServiceRequests(status, offset, pageSize)
}

// CreateServiceRequest 创建服务请求
func (s *CareService) CreateServiceRequest(req *CreateServiceRequestReq, requesterID uint) (*model.ServiceRequest, error) {
	request := &model.ServiceRequest{
		ElderlyID:   req.ElderlyID,
		Type:        req.Type,
		Notes:       req.Notes,
		RequesterID: requesterID,
		Status:      "pending",
		RequestedAt: time.Now(),
	}

	if err := s.careRepo.CreateServiceRequest(request); err != nil {
		return nil, err
	}

	return request, nil
}

// HandleServiceRequest 处理服务请求
func (s *CareService) HandleServiceRequest(id uint, status string, handlerID uint) error {
	return s.careRepo.UpdateServiceRequest(id, status, handlerID)
}

type CreateServiceRequestReq struct {
	ElderlyID uint   `json:"elderly_id" binding:"required"`
	Type      string `json:"type" binding:"required"` // 护理/送餐/打扫/紧急/其他
	Notes     string `json:"notes"`
}

type HandleServiceRequestReq struct {
	Status string `json:"status" binding:"required"` // processing/completed
}

// Health Record methods
func (s *CareService) ListHealthRecords(elderlyID uint, recordType string, page, pageSize int) ([]model.HealthRecord, int64, error) {
	offset := (page - 1) * pageSize
	return s.careRepo.FindHealthRecords(elderlyID, recordType, offset, pageSize)
}

func (s *CareService) GetLatestHealthRecords(elderlyID uint) ([]model.HealthRecord, error) {
	return s.careRepo.GetLatestHealthRecords(elderlyID)
}

func (s *CareService) CreateHealthRecord(req *CreateHealthRecordReq, recordedBy uint) (*model.HealthRecord, error) {
	record := &model.HealthRecord{
		ElderlyID:  req.ElderlyID,
		RecordType: req.RecordType,
		Value:      req.Value,
		Unit:       req.Unit,
		Value2:     req.Value2,
		Notes:      req.Notes,
		RecordedAt: time.Now(),
		RecordedBy: recordedBy,
	}

	if err := s.careRepo.CreateHealthRecord(record); err != nil {
		return nil, err
	}

	return record, nil
}

func (s *CareService) DeleteHealthRecord(id uint) error {
	return s.careRepo.DeleteHealthRecord(id)
}

type CreateHealthRecordReq struct {
	ElderlyID  uint   `json:"elderly_id" binding:"required"`
	RecordType string `json:"record_type" binding:"required"`
	Value      string `json:"value" binding:"required"`
	Unit       string `json:"unit"`
	Value2     string `json:"value2"`
	Notes      string `json:"notes"`
}
