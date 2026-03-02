package service

import (
	"elderly-care-system/internal/model"
	"elderly-care-system/internal/repository"
	"errors"
	"time"
)

type ScheduleService struct {
	scheduleRepo *repository.ScheduleRepository
	userRepo     *repository.UserRepository
}

func NewScheduleService(
	scheduleRepo *repository.ScheduleRepository,
	userRepo *repository.UserRepository,
) *ScheduleService {
	return &ScheduleService{
		scheduleRepo: scheduleRepo,
		userRepo:     userRepo,
	}
}

type CreateScheduleRequest struct {
	StaffID   uint      `json:"staff_id" binding:"required"`
	Date      time.Time `json:"date" binding:"required"`
	ShiftType string    `json:"shift_type" binding:"required"` // 早班/中班/晚班
	StartTime string    `json:"start_time" binding:"required"`
	EndTime   string    `json:"end_time" binding:"required"`
	Notes     string    `json:"notes"`
}

type ScheduleResponse struct {
	ID        uint   `json:"id"`
	StaffID   uint   `json:"staff_id"`
	StaffName string `json:"staff_name"`
	Date      string `json:"date"`
	ShiftType  string `json:"shift_type"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	Status     string `json:"status"`
	Notes      string `json:"notes"`
}

// CreateSchedule 创建排班
func (s *ScheduleService) CreateSchedule(req *CreateScheduleRequest) (*model.Schedule, error) {
	// 验证员工存在
	staff, err := s.userRepo.FindByID(req.StaffID)
	if err != nil {
		return nil, errors.New("员工不存在")
	}

	// 检查该员工当天是否已有排班
	existing, _, _ := s.scheduleRepo.ListByStaff(req.StaffID, 1, 100)
	for _, e := range existing {
		if e.Date.Format("2006-01-02") == req.Date.Format("2006-01-02") {
			return nil, errors.New("该员工当天已有排班")
		}
	}

	schedule := &model.Schedule{
		StaffID:   req.StaffID,
		Date:      req.Date,
		ShiftType: req.ShiftType,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Notes:     req.Notes,
		Status:    "scheduled",
	}

	if err := s.scheduleRepo.Create(schedule); err != nil {
		return nil, err
	}

	// 加载员工信息
	schedule.Staff = staff

	return schedule, nil
}

// GetScheduleList 获取排班列表
func (s *ScheduleService) GetScheduleList(startDate, endDate time.Time, page, pageSize int) ([]ScheduleResponse, int64, error) {
	schedules, total, err := s.scheduleRepo.ListByDate(startDate, endDate, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]ScheduleResponse, len(schedules))
	for i, schedule := range schedules {
		staffName := ""
		if schedule.Staff != nil {
			staffName = schedule.Staff.Nickname
		}
		responses[i] = ScheduleResponse{
			ID:        schedule.ID,
			StaffID:   schedule.StaffID,
			StaffName: staffName,
			Date:      schedule.Date.Format("2006-01-02"),
			ShiftType: schedule.ShiftType,
			StartTime: schedule.StartTime,
			EndTime:   schedule.EndTime,
			Status:    schedule.Status,
			Notes:     schedule.Notes,
		}
	}

	return responses, total, nil
}

// GetStaffSchedule 获取员工排班
func (s *ScheduleService) GetStaffSchedule(staffID uint, page, pageSize int) ([]model.Schedule, int64, error) {
	return s.scheduleRepo.ListByStaff(staffID, page, pageSize)
}

// UpdateScheduleStatus 更新排班状态
func (s *ScheduleService) UpdateScheduleStatus(id uint, status string) error {
	schedule, err := s.scheduleRepo.GetByID(id)
	if err != nil {
		return err
	}

	schedule.Status = status
	return s.scheduleRepo.Update(schedule)
}

// DeleteSchedule 删除排班
func (s *ScheduleService) DeleteSchedule(id uint) error {
	return s.scheduleRepo.Delete(id)
}

// GetWeeklySchedule 获取周排班
func (s *ScheduleService) GetWeeklySchedule(staffID uint, weekStart time.Time) ([]model.Schedule, error) {
	return s.scheduleRepo.GetWeeklySchedule(staffID, weekStart)
}

// GetMonthlyStats 获取月度排班统计
func (s *ScheduleService) GetMonthlyStats(year int, month time.Month) (map[string]interface{}, error) {
	startOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	schedules, _, err := s.scheduleRepo.ListByDate(startOfMonth, endOfMonth, 1, 1000)
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_shifts":    len(schedules),
		"morning_shifts": 0,
		"afternoon_shifts": 0,
		"night_shifts":    0,
		"completed":       0,
		"cancelled":       0,
	}

	for _, s := range schedules {
		switch s.ShiftType {
		case "早班":
			stats["morning_shifts"] = stats["morning_shifts"].(int) + 1
		case "中班":
			stats["afternoon_shifts"] = stats["afternoon_shifts"].(int) + 1
		case "晚班":
			stats["night_shifts"] = stats["night_shifts"].(int) + 1
		}
		switch s.Status {
		case "completed":
			stats["completed"] = stats["completed"].(int) + 1
		case "cancelled":
			stats["cancelled"] = stats["cancelled"].(int) + 1
		}
	}

	return stats, nil
}
