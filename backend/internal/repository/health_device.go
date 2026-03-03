package repository

import (
	"elderly-care-system/internal/model"
	"time"

	"gorm.io/gorm"
)

type HealthDeviceRepository struct {
	db *gorm.DB
}

func NewHealthDeviceRepository(db *gorm.DB) *HealthDeviceRepository {
	return &HealthDeviceRepository{db: db}
}

// Device methods
func (r *HealthDeviceRepository) CreateDevice(device *model.HealthDevice) error {
	return r.db.Create(device).Error
}

func (r *HealthDeviceRepository) GetDevice(id uint) (*model.HealthDevice, error) {
	var device model.HealthDevice
	err := r.db.Preload("Elderly").First(&device, id).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *HealthDeviceRepository) ListDevices(elderlyID *uint) ([]model.HealthDevice, error) {
	var devices []model.HealthDevice
	query := r.db.Preload("Elderly").Where("is_active = ?", true)
	if elderlyID != nil {
		query = query.Where("elderly_id = ?", *elderlyID)
	}
	err := query.Order("created_at DESC").Find(&devices).Error
	return devices, err
}

func (r *HealthDeviceRepository) UpdateDevice(device *model.HealthDevice) error {
	return r.db.Save(device).Error
}

func (r *HealthDeviceRepository) DeleteDevice(id uint) error {
	return r.db.Delete(&model.HealthDevice{}, id).Error
}

func (r *HealthDeviceRepository) UpdateLastConnect(id uint, batteryLevel int) error {
	now := time.Now()
	return r.db.Model(&model.HealthDevice{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"last_connect": &now,
			"battery_level": batteryLevel,
		}).Error
}

// DeviceData methods
func (r *HealthDeviceRepository) CreateDeviceData(data *model.DeviceData) error {
	return r.db.Create(data).Error
}

func (r *HealthDeviceRepository) GetDeviceData(deviceID uint, dataType string, startTime, endTime time.Time) ([]model.DeviceData, error) {
	var dataList []model.DeviceData
	query := r.db.Where("device_id = ?", deviceID)
	if dataType != "" {
		query = query.Where("data_type = ?", dataType)
	}
	if !startTime.IsZero() {
		query = query.Where("measured_at >= ?", startTime)
	}
	if !endTime.IsZero() {
		query = query.Where("measured_at <= ?", endTime)
	}
	err := query.Order("measured_at DESC").Find(&dataList).Error
	return dataList, err
}

func (r *HealthDeviceRepository) GetLatestData(elderlyID uint, dataType string, limit int) ([]model.DeviceData, error) {
	var dataList []model.DeviceData
	query := r.db.Joins("JOIN health_devices ON health_devices.id = device_data.device_id").
		Where("health_devices.elderly_id = ?", elderlyID)
	if dataType != "" {
		query = query.Where("device_data.data_type = ?", dataType)
	}
	err := query.Order("device_data.measured_at DESC").
		Limit(limit).
		Find(&dataList).Error
	return dataList, err
}

// HealthAlertRule methods
func (r *HealthDeviceRepository) CreateAlertRule(rule *model.HealthAlertRule) error {
	return r.db.Create(rule).Error
}

func (r *HealthDeviceRepository) ListAlertRules() ([]model.HealthAlertRule, error) {
	var rules []model.HealthAlertRule
	err := r.db.Where("is_active = ?", true).Order("created_at DESC").Find(&rules).Error
	return rules, err
}

func (r *HealthDeviceRepository) UpdateAlertRule(rule *model.HealthAlertRule) error {
	return r.db.Save(rule).Error
}

func (r *HealthDeviceRepository) DeleteAlertRule(id uint) error {
	return r.db.Delete(&model.HealthAlertRule{}, id).Error
}

// GetAbnormalData 获取异常数据
func (r *HealthDeviceRepository) GetAbnormalData(elderlyID uint, hours int) ([]model.DeviceData, error) {
	var dataList []model.DeviceData
	startTime := time.Now().Add(-time.Duration(hours) * time.Hour)

	err := r.db.Joins("JOIN health_devices ON health_devices.id = device_data.device_id").
		Where("health_devices.elderly_id = ?", elderlyID).
		Where("device_data.measured_at >= ?", startTime).
		Where("device_data.is_abnormal = ?", true).
		Order("device_data.measured_at DESC").
		Find(&dataList).Error

	return dataList, err
}
