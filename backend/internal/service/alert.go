package service

import (
	"elderly-care-system/internal/model"
	"elderly-care-system/internal/repository"
	"errors"
	"fmt"
)

type AlertService struct {
	alertRepo      *repository.AlertRepository
	elderlyRepo    *repository.ElderlyRepository
	medicationRepo *repository.MedicationRepository
	billRepo       *repository.BillRepository
	careRepo       *repository.CareRepository
	roomRepo       *repository.RoomRepository
}

func NewAlertService(
	alertRepo *repository.AlertRepository,
	elderlyRepo *repository.ElderlyRepository,
	medicationRepo *repository.MedicationRepository,
	billRepo *repository.BillRepository,
	careRepo *repository.CareRepository,
	roomRepo *repository.RoomRepository,
) *AlertService {
	return &AlertService{
		alertRepo:      alertRepo,
		elderlyRepo:    elderlyRepo,
		medicationRepo: medicationRepo,
		billRepo:       billRepo,
		careRepo:       careRepo,
		roomRepo:       roomRepo,
	}
}

// AlertRequest 创建预警请求
type AlertRequest struct {
	Type       string `json:"type" binding:"required"`
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content" binding:"required"`
	Level      string `json:"level"`
	EntityID   uint   `json:"entity_id"`
	EntityType string `json:"entity_type"`
}

// AlertRuleRequest 预警规则请求
type AlertRuleRequest struct {
	Name        string `json:"name" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Description string `json:"description"`
	Condition   string `json:"condition" binding:"required"`
	Threshold   string `json:"threshold"`
	Level       string `json:"level"`
}

// CreateAlert 创建预警
func (s *AlertService) CreateAlert(req *AlertRequest) (*model.Alert, error) {
	level := req.Level
	if level == "" {
		level = "info"
	}

	alert := &model.Alert{
		Type:       req.Type,
		Title:      req.Title,
		Content:    req.Content,
		Level:      level,
		EntityID:   req.EntityID,
		EntityType: req.EntityType,
		Status:     "active",
	}

	if err := s.alertRepo.Create(alert); err != nil {
		return nil, err
	}

	return alert, nil
}

// AcknowledgeAlert 确认预警
func (s *AlertService) AcknowledgeAlert(id uint, userID uint) error {
	alert, err := s.alertRepo.GetByID(id)
	if err != nil {
		return err
	}

	if alert.Status != "active" {
		return errors.New("预警已被处理")
	}

	return s.alertRepo.AcknowledgeAlert(id, userID)
}

// ResolveAlert 解决预警
func (s *AlertService) ResolveAlert(id uint) error {
	alert, err := s.alertRepo.GetByID(id)
	if err != nil {
		return err
	}

	if alert.Status == "resolved" {
		return errors.New("预警已解决")
	}

	return s.alertRepo.ResolveAlert(id)
}

// GetAlert 获取预警详情
func (s *AlertService) GetAlert(id uint) (*model.Alert, error) {
	return s.alertRepo.GetByID(id)
}

// ListAlerts 获取预警列表
func (s *AlertService) ListAlerts(page, pageSize int, status string) ([]model.Alert, int64, error) {
	return s.alertRepo.List(page, pageSize, status)
}

// GetActiveAlerts 获取活跃预警
func (s *AlertService) GetActiveAlerts() ([]model.Alert, error) {
	return s.alertRepo.GetActiveAlerts()
}

// CheckAllAlerts 检查所有预警规则并生成预警
func (s *AlertService) CheckAllAlerts() error {
	// 检查药品库存不足
	if err := s.checkLowStockMedications(); err != nil {
		return err
	}

	// 检查药品即将过期
	if err := s.checkExpiringMedications(); err != nil {
		return err
	}

	// 检查账单逾期未支付
	if err := s.checkOverdueBills(); err != nil {
		return err
	}

	// 检查床位可用情况
	if err := s.checkBedAvailability(); err != nil {
		return err
	}

	return nil
}

// checkLowStockMedications 检查库存不足的药品
func (s *AlertService) checkLowStockMedications() error {
	medications, err := s.medicationRepo.GetLowStockMedications()
	if err != nil {
		return err
	}

	for _, med := range medications {
		// 检查是否已存在相同的active预警
		existingAlerts, _, _ := s.alertRepo.ListByType("medication_low", 1, 100)
		alreadyAlerted := false
		for _, alert := range existingAlerts {
			if alert.EntityID == med.ID && alert.Status == "active" {
				alreadyAlerted = true
				break
			}
		}

		if !alreadyAlerted {
			_, _ = s.CreateAlert(&AlertRequest{
				Type:       "medication_low",
				Title:      "药品库存不足",
				Content:    fmt.Sprintf("药品【%s】当前库存为%d，低于最低库存%d", med.Name, med.Stock, med.MinStock),
				Level:      "warning",
				EntityID:   med.ID,
				EntityType: "medication",
			})
		}
	}

	return nil
}

// checkExpiringMedications 检查即将过期的药品
func (s *AlertService) checkExpiringMedications() error {
	days := 30
	medications, err := s.medicationRepo.GetExpiringMedications(days)
	if err != nil {
		return err
	}

	for _, med := range medications {
		existingAlerts, _, _ := s.alertRepo.ListByType("medication_expiry", 1, 100)
		alreadyAlerted := false
		for _, alert := range existingAlerts {
			if alert.EntityID == med.ID && alert.Status == "active" {
				alreadyAlerted = true
				break
			}
		}

		if !alreadyAlerted {
			_, _ = s.CreateAlert(&AlertRequest{
				Type:       "medication_expiry",
				Title:      "药品即将过期",
				Content:    fmt.Sprintf("药品【%s】将于%s过期", med.Name, med.ExpiryDate.Format("2006-01-02")),
				Level:      "warning",
				EntityID:   med.ID,
				EntityType: "medication",
			})
		}
	}

	return nil
}

// checkOverdueBills 检查逾期未支付的账单
func (s *AlertService) checkOverdueBills() error {
	// TODO: 需要BillRepository添加ListAllUnpaid方法来获取所有未支付账单
	// 暂时跳过此检查
	return nil
}

// checkBedAvailability 检查床位可用情况
func (s *AlertService) checkBedAvailability() error {
	stats, err := s.roomRepo.GetBedStats()
	if err != nil {
		return err
	}

	total := stats["total"]
	occupied := stats["occupied"]

	occupancyRate := 0.0
	if total > 0 {
		occupancyRate = float64(occupied) / float64(total) * 100
	}

	// 如果入住率超过90%，发出预警
	if occupancyRate > 90 {
		existingAlerts, _, _ := s.alertRepo.ListByType("bed_full", 1, 100)
		hasActiveAlert := false
		for _, alert := range existingAlerts {
			if alert.Status == "active" {
				hasActiveAlert = true
				break
			}
		}

		if !hasActiveAlert {
			_, _ = s.CreateAlert(&AlertRequest{
				Type:       "bed_full",
				Title:      "床位紧张",
				Content:    fmt.Sprintf("当前床位入住率已达%.1f%%，剩余空床位%d个", occupancyRate, total-occupied),
				Level:      "warning",
				EntityType: "room",
			})
		}
	}

	return nil
}

// GetAlertSummary 获取预警统计摘要
func (s *AlertService) GetAlertSummary() (map[string]interface{}, error) {
	activeAlerts, err := s.alertRepo.GetActiveAlerts()
	if err != nil {
		return nil, err
	}

	summary := map[string]interface{}{
		"total_active": len(activeAlerts),
		"critical":     0,
		"warning":      0,
		"info":         0,
		"by_type":      make(map[string]int),
	}

	for _, alert := range activeAlerts {
		switch alert.Level {
		case "critical":
			summary["critical"] = summary["critical"].(int) + 1
		case "warning":
			summary["warning"] = summary["warning"].(int) + 1
		case "info":
			summary["info"] = summary["info"].(int) + 1
		}

		byType := summary["by_type"].(map[string]int)
		byType[alert.Type]++
	}

	return summary, nil
}

// Alert Rule Management
func (s *AlertService) CreateRule(req *AlertRuleRequest) (*model.AlertRule, error) {
	level := req.Level
	if level == "" {
		level = "warning"
	}

	rule := &model.AlertRule{
		Name:        req.Name,
		Type:        req.Type,
		Description: req.Description,
		Condition:   req.Condition,
		Threshold:   req.Threshold,
		Level:       level,
		IsEnabled:   true,
	}

	if err := s.alertRepo.CreateRule(rule); err != nil {
		return nil, err
	}

	return rule, nil
}

func (s *AlertService) ListRules(page, pageSize int) ([]model.AlertRule, int64, error) {
	return s.alertRepo.ListRules(page, pageSize)
}

func (s *AlertService) UpdateRule(id uint, updates map[string]interface{}) error {
	rule, err := s.alertRepo.GetRuleByID(id)
	if err != nil {
		return err
	}

	if name, ok := updates["name"]; ok {
		rule.Name = name.(string)
	}
	if description, ok := updates["description"]; ok {
		rule.Description = description.(string)
	}
	if condition, ok := updates["condition"]; ok {
		rule.Condition = condition.(string)
	}
	if threshold, ok := updates["threshold"]; ok {
		rule.Threshold = threshold.(string)
	}
	if level, ok := updates["level"]; ok {
		rule.Level = level.(string)
	}
	if isEnabled, ok := updates["is_enabled"]; ok {
		rule.IsEnabled = isEnabled.(bool)
	}

	return s.alertRepo.UpdateRule(rule)
}

func (s *AlertService) DeleteRule(id uint) error {
	return s.alertRepo.DeleteRule(id)
}
