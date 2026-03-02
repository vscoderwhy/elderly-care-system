package service

import (
	"elderly-care-system/internal/model"
	"elderly-care-system/internal/repository"
	"errors"
	"time"
)

type VisitService struct {
	visitRepo  *repository.VisitRepository
	elderlyRepo *repository.ElderlyRepository
}

func NewVisitService(
	visitRepo *repository.VisitRepository,
	elderlyRepo *repository.ElderlyRepository,
) *VisitService {
	return &VisitService{
		visitRepo:  visitRepo,
		elderlyRepo: elderlyRepo,
	}
}

// VisitRequest 探视预约请求
type VisitRequest struct {
	ElderlyID    uint   `json:"elderly_id" binding:"required"`
	VisitorName  string `json:"visitor_name" binding:"required"`
	VisitorPhone string `json:"visitor_phone" binding:"required"`
	Relationship string `json:"relationship" binding:"required"`
	VisitDate    string `json:"visit_date" binding:"required"`
	VisitTime    string `json:"visit_time" binding:"required"`
	VisitorCount int    `json:"visitor_count"`
	Notes        string `json:"notes"`
}

// UpdateVisitRequest 更新探视预约请求
type UpdateVisitRequest struct {
	VisitorName  string `json:"visitor_name"`
	VisitorPhone string `json:"visitor_phone"`
	Relationship string `json:"relationship"`
	VisitDate    string `json:"visit_date"`
	VisitTime    string `json:"visit_time"`
	VisitorCount int    `json:"visitor_count"`
	Status       string `json:"status"`
	Notes        string `json:"notes"`
}

// CreateVisit 创建探视预约
func (s *VisitService) CreateVisit(req *VisitRequest) (*model.VisitAppointment, error) {
	// 验证老人存在
	elderly, err := s.elderlyRepo.FindByID(req.ElderlyID)
	if err != nil {
		return nil, errors.New("老人不存在")
	}

	// 检查老人是否已入住（有床位）
	if elderly.BedID == nil {
		return nil, errors.New("老人未入住或已退住")
	}

	// 解析日期
	visitDate, err := time.Parse("2006-01-02", req.VisitDate)
	if err != nil {
		return nil, errors.New("日期格式错误")
	}

	// 检查预约日期不能是过去
	if visitDate.Before(time.Now().Truncate(24 * time.Hour)) {
		return nil, errors.New("不能预约过去的日期")
	}

	visitorCount := req.VisitorCount
	if visitorCount <= 0 {
		visitorCount = 1
	}

	visit := &model.VisitAppointment{
		ElderlyID:    req.ElderlyID,
		VisitorName:  req.VisitorName,
		VisitorPhone: req.VisitorPhone,
		Relationship: req.Relationship,
		VisitDate:    visitDate,
		VisitTime:    req.VisitTime,
		VisitorCount: visitorCount,
		Status:       "pending",
		Notes:        req.Notes,
	}

	if err := s.visitRepo.Create(visit); err != nil {
		return nil, err
	}

	return visit, nil
}

// UpdateVisit 更新探视预约
func (s *VisitService) UpdateVisit(id uint, req *UpdateVisitRequest) error {
	visit, err := s.visitRepo.GetByID(id)
	if err != nil {
		return err
	}

	// 如果修改日期，需要验证
	if req.VisitDate != "" {
		visitDate, err := time.Parse("2006-01-02", req.VisitDate)
		if err != nil {
			return errors.New("日期格式错误")
		}

		// 只有pending或confirmed状态的预约才能修改日期
		if visit.Status != "pending" && visit.Status != "confirmed" {
			return errors.New("当前状态不能修改日期")
		}

		visit.VisitDate = visitDate
	}

	if req.VisitorName != "" {
		visit.VisitorName = req.VisitorName
	}
	if req.VisitorPhone != "" {
		visit.VisitorPhone = req.VisitorPhone
	}
	if req.Relationship != "" {
		visit.Relationship = req.Relationship
	}
	if req.VisitTime != "" {
		visit.VisitTime = req.VisitTime
	}
	if req.VisitorCount > 0 {
		visit.VisitorCount = req.VisitorCount
	}
	if req.Status != "" {
		visit.Status = req.Status
	}
	if req.Notes != "" {
		visit.Notes = req.Notes
	}

	return s.visitRepo.Update(visit)
}

// DeleteVisit 删除探视预约
func (s *VisitService) DeleteVisit(id uint) error {
	visit, err := s.visitRepo.GetByID(id)
	if err != nil {
		return err
	}

	// 只能删除pending或cancelled状态的预约
	if visit.Status == "completed" {
		return errors.New("已完成的探视预约不能删除")
	}

	return s.visitRepo.Delete(id)
}

// GetVisit 获取探视预约详情
func (s *VisitService) GetVisit(id uint) (*model.VisitAppointment, error) {
	return s.visitRepo.GetByID(id)
}

// ListVisits 获取探视预约列表
func (s *VisitService) ListVisits(page, pageSize int) ([]model.VisitAppointment, int64, error) {
	return s.visitRepo.List(page, pageSize)
}

// ListElderlyVisits 获取老人的探视预约列表
func (s *VisitService) ListElderlyVisits(elderlyID uint, page, pageSize int) ([]model.VisitAppointment, int64, error) {
	return s.visitRepo.ListByElderly(elderlyID, page, pageSize)
}

// ConfirmVisit 确认探视预约
func (s *VisitService) ConfirmVisit(id uint) error {
	visit, err := s.visitRepo.GetByID(id)
	if err != nil {
		return err
	}

	if visit.Status != "pending" {
		return errors.New("只能确认待处理的预约")
	}

	visit.Status = "confirmed"
	return s.visitRepo.Update(visit)
}

// CancelVisit 取消探视预约
func (s *VisitService) CancelVisit(id uint) error {
	visit, err := s.visitRepo.GetByID(id)
	if err != nil {
		return err
	}

	if visit.Status == "completed" {
		return errors.New("已完成的探视不能取消")
	}

	if visit.Status == "cancelled" {
		return errors.New("探视预约已取消")
	}

	visit.Status = "cancelled"
	return s.visitRepo.Update(visit)
}

// CompleteVisit 完成探视
func (s *VisitService) CompleteVisit(id uint) error {
	visit, err := s.visitRepo.GetByID(id)
	if err != nil {
		return err
	}

	if visit.Status != "confirmed" {
		return errors.New("只能确认已确认的探视")
	}

	visit.Status = "completed"
	return s.visitRepo.Update(visit)
}

// GetTodayVisits 获取今日探视预约
func (s *VisitService) GetTodayVisits() ([]model.VisitAppointment, error) {
	return s.visitRepo.GetTodayVisits()
}

// GetUpcomingVisits 获取即将到来的探视预约
func (s *VisitService) GetUpcomingVisits(days int) ([]model.VisitAppointment, error) {
	return s.visitRepo.GetUpcomingVisits(days)
}

// GetVisitsByDateRange 获取指定日期范围内的探视预约
func (s *VisitService) GetVisitsByDateRange(startDate, endDate string) ([]model.VisitAppointment, error) {
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, errors.New("开始日期格式错误")
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, errors.New("结束日期格式错误")
	}

	return s.visitRepo.ListByDateRange(start, end)
}
