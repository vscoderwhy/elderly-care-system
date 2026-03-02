package service

import (
	"elderly-care-system/internal/model"
	"elderly-care-system/internal/repository"
)

type ElderlyService struct {
	elderlyRepo *repository.ElderlyRepository
}

func NewElderlyService(elderlyRepo *repository.ElderlyRepository) *ElderlyService {
	return &ElderlyService{elderlyRepo: elderlyRepo}
}

func (s *ElderlyService) List(page, pageSize int) ([]model.Elderly, int64, error) {
	offset := (page - 1) * pageSize
	return s.elderlyRepo.List(offset, pageSize)
}

func (s *ElderlyService) Get(id uint) (*model.Elderly, error) {
	return s.elderlyRepo.FindByID(id)
}

func (s *ElderlyService) Create(req *CreateElderlyRequest) (*model.Elderly, error) {
	elderly := &model.Elderly{
		Name:             req.Name,
		Gender:           req.Gender,
		IDCard:           req.IDCard,
		Phone:            req.Phone,
		EmergencyContact: req.EmergencyContact,
		EmergencyPhone:   req.EmergencyPhone,
		CareLevel:        req.CareLevel,
		HealthStatus:     "{}", // 默认空JSON对象
	}

	if err := s.elderlyRepo.Create(elderly); err != nil {
		return nil, err
	}

	return elderly, nil
}

func (s *ElderlyService) Update(id uint, req *UpdateElderlyRequest) error {
	elderly, err := s.elderlyRepo.FindByID(id)
	if err != nil {
		return err
	}

	elderly.Name = req.Name
	elderly.Gender = req.Gender
	elderly.Phone = req.Phone
	elderly.EmergencyContact = req.EmergencyContact
	elderly.EmergencyPhone = req.EmergencyPhone
	elderly.CareLevel = req.CareLevel

	return s.elderlyRepo.Update(elderly)
}

func (s *ElderlyService) Delete(id uint) error {
	return s.elderlyRepo.Delete(id)
}

type CreateElderlyRequest struct {
	Name             string `json:"name" binding:"required"`
	Gender           string `json:"gender" binding:"required"`
	IDCard           string `json:"id_card"`
	Phone            string `json:"phone"`
	EmergencyContact string `json:"emergency_contact"`
	EmergencyPhone   string `json:"emergency_phone"`
	CareLevel        int    `json:"care_level"`
}

type UpdateElderlyRequest struct {
	Name             string `json:"name"`
	Gender           string `json:"gender"`
	Phone            string `json:"phone"`
	EmergencyContact string `json:"emergency_contact"`
	EmergencyPhone   string `json:"emergency_phone"`
	CareLevel        int    `json:"care_level"`
}
