package repository

import (
	"elderly-care-system/internal/model"
	"gorm.io/gorm"
)

type ElderlyRepository struct {
	db *gorm.DB
}

func NewElderlyRepository(db *gorm.DB) *ElderlyRepository {
	return &ElderlyRepository{db: db}
}

func (r *ElderlyRepository) Create(elderly *model.Elderly) error {
	return r.db.Create(elderly).Error
}

func (r *ElderlyRepository) FindByID(id uint) (*model.Elderly, error) {
	var elderly model.Elderly
	err := r.db.Preload("Bed").Preload("Bed.Room").Preload("Families").First(&elderly, id).Error
	if err != nil {
		return nil, err
	}
	return &elderly, nil
}

func (r *ElderlyRepository) List(offset, limit int) ([]model.Elderly, int64, error) {
	var elderly []model.Elderly
	var total int64

	err := r.db.Model(&model.Elderly{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("Bed").Offset(offset).Limit(limit).Find(&elderly).Error
	return elderly, total, err
}

func (r *ElderlyRepository) Update(elderly *model.Elderly) error {
	return r.db.Save(elderly).Error
}

func (r *ElderlyRepository) Delete(id uint) error {
	return r.db.Delete(&model.Elderly{}, id).Error
}

func (r *ElderlyRepository) FindByFamilyUserID(userID uint) ([]model.Elderly, error) {
	var elderly []model.Elderly
	err := r.db.Joins("JOIN elderly_families ON elderly_families.elderly_id = elderlies.id").
		Where("elderly_families.user_id = ?", userID).
		Preload("Bed").
		Find(&elderly).Error
	return elderly, err
}

// Count 获取老人总数
func (r *ElderlyRepository) Count() (int, error) {
	var count int64
	err := r.db.Model(&model.Elderly{}).Count(&count).Error
	return int(count), err
}
