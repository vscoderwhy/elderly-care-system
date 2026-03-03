package service

import (
	"elderly-care-system/internal/model"
	"elderly-care-system/internal/repository"
	"fmt"
	"math"
	"time"
)

type AttendanceService struct {
	attendanceRepo *repository.AttendanceRepository
	userRepo       *repository.UserRepository
	careRepo       *repository.CareRepository
}

func NewAttendanceService(
	attendanceRepo *repository.AttendanceRepository,
	userRepo *repository.UserRepository,
	careRepo *repository.CareRepository,
) *AttendanceService {
	return &AttendanceService{
		attendanceRepo: attendanceRepo,
		userRepo:       userRepo,
		careRepo:       careRepo,
	}
}

// ClockIn 上班打卡
type ClockInRequest struct {
	StaffID   uint    `json:"staff_id" binding:"required"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Location  string  `json:"location"`
	Device    string  `json:"device"`
}

func (s *AttendanceService) ClockIn(req *ClockInRequest) (*model.Attendance, error) {
	// 检查今日是否已打卡
	clockIn, _, err := s.attendanceRepo.GetTodayAttendance(req.StaffID)
	if err != nil {
		return nil, err
	}
	if clockIn != nil {
		return nil, fmt.Errorf("今日已打卡")
	}

	now := time.Now()
	status := s.checkClockStatus(now, true)

	attendance := &model.Attendance{
		StaffID:   req.StaffID,
		Date:      now.Truncate(24 * time.Hour),
		Type:      "clock_in",
		Time:      now,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Location:  req.Location,
		Device:    req.Device,
		Status:    status,
	}

	if err := s.attendanceRepo.CreateAttendance(attendance); err != nil {
		return nil, err
	}

	return attendance, nil
}

// ClockOut 下班打卡
func (s *AttendanceService) ClockOut(req *ClockInRequest) (*model.Attendance, error) {
	// 检查今日是否已上班打卡
	clockIn, _, err := s.attendanceRepo.GetTodayAttendance(req.StaffID)
	if err != nil {
		return nil, err
	}
	if clockIn == nil {
		return nil, fmt.Errorf("请先上班打卡")
	}

	// 检查是否已下班打卡
	_, clockOut, err := s.attendanceRepo.GetTodayAttendance(req.StaffID)
	if err != nil {
		return nil, err
	}
	if clockOut != nil {
		return nil, fmt.Errorf("今日已完成下班打卡")
	}

	now := time.Now()
	status := s.checkClockStatus(now, false)

	attendance := &model.Attendance{
		StaffID:   req.StaffID,
		Date:      now.Truncate(24 * time.Hour),
		Type:      "clock_out",
		Time:      now,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Location:  req.Location,
		Device:    req.Device,
		Status:    status,
	}

	if err := s.attendanceRepo.CreateAttendance(attendance); err != nil {
		return nil, err
	}

	return attendance, nil
}

// checkClockStatus 检查打卡状态
func (s *AttendanceService) checkClockStatus(t time.Time, isIn bool) string {
	hour := t.Hour()
	minute := t.Minute()

	// 早班 8:00, 晚班 20:00
	if isIn {
		// 上班打卡
		if hour > 8 || (hour == 8 && minute > 10) {
			return "late"
		}
		return "normal"
	} else {
		// 下班打卡
		if hour < 17 || (hour == 17 && minute < 30) {
			return "early"
		}
		return "normal"
	}
}

// GetAttendanceStats 获取考勤统计
func (s *AttendanceService) GetAttendanceStats(staffID uint, year, month int) (map[string]interface{}, error) {
	stats, err := s.attendanceRepo.GetMonthlyAttendanceStats(staffID, year, month)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

// GeneratePerformance 生成绩效
func (s *AttendanceService) GeneratePerformance(staffID uint, year, month int) (*model.Performance, error) {
	// 检查是否已存在
	existing, err := s.attendanceRepo.GetPerformance(staffID, year, month)
	if err == nil && existing != nil {
		return existing, nil
	}

	// 获取考勤统计
	attendanceStats, err := s.attendanceRepo.GetMonthlyAttendanceStats(staffID, year, month)
	if err != nil {
		return nil, err
	}

	// 获取护理记录数
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	careRecords, _, err := s.careRepo.FindRecordsByDateRange(startDate, endDate)
	if err != nil {
		// 忽略错误，继续处理
		careRecords = []model.CareRecord{}
	}

	// 计算绩效
	performance := &model.Performance{
		StaffID:           staffID,
		Year:              year,
		Month:             month,
		ServiceScore:      85, // 默认服务分，可从评价系统获取
		ElderlySatisfaction: 85, // 默认满意度
		CareRecordsCount:  len(careRecords),
		AttendanceRate:    attendanceStats["attendance_rate"].(float64),
		OvertimeHours:     0, // 从排班计算
		TotalScore:        0, // 综合计算
	}

	// 计算综合评分
	performance.TotalScore = performance.ServiceScore*0.3 +
		performance.ElderlySatisfaction*0.3 +
		float64(performance.CareRecordsCount)*0.2 +
		performance.AttendanceRate*0.2

	if err := s.attendanceRepo.CreatePerformance(performance); err != nil {
		return nil, err
	}

	return performance, nil
}

// CalculateSalary 计算工资
func (s *AttendanceService) CalculateSalary(staffID uint, year, month int) (*model.Salary, error) {
	// 检查是否已存在
	existing, err := s.attendanceRepo.GetSalary(staffID, year, month)
	if err == nil && existing != nil {
		return existing, nil
	}

	// 获取绩效
	performance, err := s.GeneratePerformance(staffID, year, month)
	if err != nil {
		return nil, err
	}

	// 基础工资 (可从配置或员工信息获取)
	baseSalary := 3000.0

	// 绩效工资
	performanceSalary := (performance.TotalScore / 100) * 1000

	// 加班费
	overtimeSalary := performance.OvertimeHours * 50

	// 计算应发工资
	grossSalary := baseSalary + performanceSalary + overtimeSalary + performance.RewardAmount - performance.PenaltyAmount

	// 社保扣除 (约10%)
	socialInsurance := grossSalary * 0.1

	// 实发工资
	netSalary := grossSalary - socialInsurance

	salary := &model.Salary{
		StaffID:            staffID,
		Year:               year,
		Month:              month,
		BaseSalary:         baseSalary,
		PerformanceSalary: performanceSalary,
		OvertimeSalary:    overtimeSalary,
		Reward:             performance.RewardAmount,
		Penalty:            performance.PenaltyAmount,
		SocialInsurance:    socialInsurance,
		GrossSalary:        grossSalary,
		NetSalary:          netSalary,
		Status:             "pending",
	}

	if err := s.attendanceRepo.CreateSalary(salary); err != nil {
		return nil, err
	}

	return salary, nil
}

// GetStaffAttendance 获取员工考勤记录
func (s *AttendanceService) GetStaffAttendance(staffID uint, startDate, endDate time.Time) ([]model.Attendance, error) {
	return s.attendanceRepo.GetStaffAttendance(staffID, startDate, endDate)
}

// ListShiftRules 获取排班规则列表
func (s *AttendanceService) ListShiftRules() ([]model.ShiftRule, error) {
	return s.attendanceRepo.ListShiftRules()
}

// CreateShiftRule 创建排班规则
type ShiftRuleRequest struct {
	Name        string  `json:"name" binding:"required"`
	Type        string  `json:"type" binding:"required"`
	StartTime   string  `json:"start_time" binding:"required"`
	EndTime     string  `json:"end_time" binding:"required"`
	BreakTime   int     `json:"break_time"`
	Description string  `json:"description"`
}

func (s *AttendanceService) CreateShiftRule(req *ShiftRuleRequest) (*model.ShiftRule, error) {
	// 计算工作时长
	start, _ := time.Parse("15:04", req.StartTime)
	end, _ := time.Parse("15:04", req.EndTime)
	duration := end.Sub(start)
	workHours := duration.Hours() - float64(req.BreakTime)/60

	rule := &model.ShiftRule{
		Name:        req.Name,
		Type:        req.Type,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		BreakTime:   req.BreakTime,
		WorkHours:   math.Round(workHours*100) / 100,
		Description: req.Description,
		IsActive:    true,
	}

	if err := s.attendanceRepo.CreateShiftRule(rule); err != nil {
		return nil, err
	}

	return rule, nil
}
