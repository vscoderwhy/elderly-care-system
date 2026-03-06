package service

import (
	"elderly-care-system/internal/model"
	"elderly-care-system/internal/repository"
	"encoding/json"
	"fmt"
	"time"
)

type HealthDeviceService struct {
	deviceRepo    *repository.HealthDeviceRepository
	elderlyRepo   *repository.ElderlyRepository
	alertService  *repository.AlertRepository
}

func NewHealthDeviceService(
	deviceRepo *repository.HealthDeviceRepository,
	elderlyRepo *repository.ElderlyRepository,
	alertService *repository.AlertRepository,
) *HealthDeviceService {
	return &HealthDeviceService{
		deviceRepo:   deviceRepo,
		elderlyRepo:  elderlyRepo,
		alertService: alertService,
	}
}

// BindDeviceRequest 绑定设备请求
type BindDeviceRequest struct {
	ElderlyID   uint    `json:"elderly_id" binding:"required"`
	DeviceType  string  `json:"device_type" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	SerialNo    string  `json:"serial_no" binding:"required"`
	MacAddress string  `json:"mac_address"`
}

func (s *HealthDeviceService) BindDevice(req *BindDeviceRequest) (*model.HealthDevice, error) {
	// 检查老人是否存在
	_, err := s.elderlyRepo.FindByID(req.ElderlyID)
	if err != nil {
		return nil, fmt.Errorf("老人不存在")
	}

	// TODO: 检查serial_no是否已存在

	device := &model.HealthDevice{
		Name:       req.Name,
		Type:       req.DeviceType,
		SerialNo:   req.SerialNo,
		MacAddress: req.MacAddress,
		ElderlyID:  &req.ElderlyID,
		Status:     "active",
		IsActive:   true,
	}

	if err := s.deviceRepo.CreateDevice(device); err != nil {
		return nil, err
	}

	return device, nil
}

// ReceiveDeviceData 接收设备数据
type DeviceDataRequest struct {
	SerialNo   string                 `json:"serial_no" binding:"required"`
	DataType   string                 `json:"data_type" binding:"required"`
	DataValue  map[string]interface{} `json:"data_value" binding:"required"`
	MeasuredAt string                 `json:"measured_at"`
}

func (s *HealthDeviceService) ReceiveDeviceData(req *DeviceDataRequest) (*model.DeviceData, error) {
	// 查找设备
	devices, err := s.deviceRepo.ListDevices(nil)
	var device *model.HealthDevice
	for _, d := range devices {
		if d.SerialNo == req.SerialNo {
			device = &d
			break
		}
	}

	if device == nil {
		return nil, fmt.Errorf("设备未绑定")
	}

	// 解析测量时间
	measuredAt, err := time.Parse(time.RFC3339, req.MeasuredAt)
	if err != nil {
		return nil, fmt.Errorf("时间格式错误")
	}

	// 转换数据为JSON字符串
	dataValueJSON, _ := json.Marshal(req.DataValue)

	// 检查数据是否异常
	isAbnormal, alertLevel := s.checkDataAbnormal(*device.ElderlyID, req.DataType, req.DataValue)

	// 创建数据记录
	data := &model.DeviceData{
		DeviceID:   device.ID,
		ElderlyID:  *device.ElderlyID,
		DataType:   req.DataType,
		DataValue:  string(dataValueJSON),
		MeasuredAt: measuredAt,
		IsAbnormal: isAbnormal,
		AlertLevel: alertLevel,
		SyncedAt:   time.Now(),
	}

	if err := s.deviceRepo.CreateDeviceData(data); err != nil {
		return nil, err
	}

	// 如果数据异常，创建预警
	if isAbnormal && device.ElderlyID != nil {
		s.createHealthAlert(*device.ElderlyID, req.DataType, req.DataValue, alertLevel)
	}

	// 更新设备最后连接时间
	batteryLevel := 100 // 从数据中获取电池电量
	if val, ok := req.DataValue["battery"]; ok {
		if bl, ok := val.(float64); ok {
			batteryLevel = int(bl)
		}
	}
	s.deviceRepo.UpdateLastConnect(device.ID, batteryLevel)

	return data, nil
}

// checkDataAbnormal 检查数据是否异常
func (s *HealthDeviceService) checkDataAbnormal(elderlyID uint, dataType string, data map[string]interface{}) (bool, string) {
	// TODO: 从数据库加载预警规则进行检查
	// 这里先写死一些常见的阈值

	switch dataType {
	case "blood_pressure":
		if systolic, ok := data["systolic"].(float64); ok {
			if systolic >= 140 || systolic <= 90 {
				return true, "warning"
			}
		}
	case "blood_sugar":
		if value, ok := data["value"].(float64); ok {
			if value > 11.0 || value < 3.9 {
				return true, "warning"
			}
		}
	case "temperature":
		if value, ok := data["value"].(float64); ok {
			if value > 37.3 || value < 36.0 {
				return true, "warning"
			}
		}
	case "heart_rate":
		if value, ok := data["value"].(float64); ok {
			if value > 120 || value < 50 {
				return true, "warning"
			}
		}
	}

	return false, "normal"
}

// createHealthAlert 创建健康预警
func (s *HealthDeviceService) createHealthAlert(elderlyID uint, dataType string, data map[string]interface{}, level string) {
	var alertType string
	var description string

	switch dataType {
	case "blood_pressure":
		alertType = "blood_pressure_abnormal"
		if systolic, ok := data["systolic"].(float64); ok {
			description = fmt.Sprintf("血压异常: %.0f mmHg", systolic)
		}
	case "blood_sugar":
		alertType = "blood_sugar_abnormal"
		if value, ok := data["value"].(float64); ok {
			description = fmt.Sprintf("血糖异常: %.1f mmol/L", value)
		}
	case "temperature":
		alertType = "temperature_abnormal"
		if value, ok := data["value"].(float64); ok {
			description = fmt.Sprintf("体温异常: %.1f°C", value)
		}
	case "heart_rate":
		alertType = "heart_rate_abnormal"
		if value, ok := data["value"].(float64); ok {
			description = fmt.Sprintf("心率异常: %.0f bpm", value)
		}
	default:
		return
	}

	alert := &model.Alert{
		EntityID:   elderlyID,
		EntityType: "elderly",
		Type:       alertType,
		Level:      level,
		Title:      description,
		Content:    description,
		Status:     "active",
	}

	s.alertService.Create(alert)
}

// ListDevices 获取设备列表
func (s *HealthDeviceService) ListDevices(elderlyID *uint) ([]model.HealthDevice, error) {
	return s.deviceRepo.ListDevices(elderlyID)
}

// GetDeviceTrends 获取设备数据趋势
func (s *HealthDeviceService) GetDeviceTrends(deviceID uint, dataType string, days int) (map[string]interface{}, error) {
	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -days)

	data, err := s.deviceRepo.GetDeviceData(deviceID, dataType, startTime, endTime)
	if err != nil {
		return nil, err
	}

	// 分析数据趋势
	// TODO: 实现更复杂的趋势分析

	result := map[string]interface{}{
		"device_id": deviceID,
		"data_type": dataType,
		"period": map[string]interface{}{
			"start": startTime.Format("2006-01-02"),
			"end":   endTime.Format("2006-01-02"),
		},
		"data_points": len(data),
		"latest": func() any {
			if len(data) > 0 {
				return data[0]
			}
			return nil
		}(),
		"average":    nil,
		"max":        nil,
		"min":        nil,
	}

	return result, nil
}

// GetAbnormalData 获取异常数据
func (s *HealthDeviceService) GetAbnormalData(elderlyID uint, hours int) ([]model.DeviceData, error) {
	return s.deviceRepo.GetAbnormalData(elderlyID, hours)
}

// UnbindDevice 解绑设备
func (s *HealthDeviceService) UnbindDevice(id uint) error {
	return s.deviceRepo.DeleteDevice(id)
}
