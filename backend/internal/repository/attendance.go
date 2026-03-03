package repository

import (
	"elderly-care-system/internal/model"
	"time"

	"gorm.io/gorm"
)

type AttendanceRepository struct {
	db *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) *AttendanceRepository {
	return &AttendanceRepository{db: db}
}

// CreateAttendance 创建考勤记录
func (r *AttendanceRepository) CreateAttendance(attendance *model.Attendance) error {
	return r.db.Create(attendance).Error
}

// GetStaffAttendance 获取员工考勤记录
func (r *AttendanceRepository) GetStaffAttendance(staffID uint, startDate, endDate time.Time) ([]model.Attendance, error) {
	var records []model.Attendance
	err := r.db.Where("staff_id = ? AND date >= ? AND date <= ?", staffID, startDate, endDate).
		Order("date DESC, time DESC").
		Find(&records).Error
	return records, err
}

// GetTodayAttendance 获取今日考勤
func (r *AttendanceRepository) GetTodayAttendance(staffID uint) (clockIn, clockOut *model.Attendance, err error) {
	today := time.Now().Truncate(24 * time.Hour)
	var records []model.Attendance
	err = r.db.Where("staff_id = ? AND date = ?", staffID, today).Find(&records).Error
	if err != nil {
		return nil, nil, err
	}
	for _, r := range records {
		if r.Type == "clock_in" {
			clockIn = &r
		} else if r.Type == "clock_out" {
			clockOut = &r
		}
	}
	return
}

// GetMonthlyAttendanceStats 获取月度考勤统计
func (r *AttendanceRepository) GetMonthlyAttendanceStats(staffID uint, year, month int) (map[string]interface{}, error) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1)

	type Result struct {
		TotalDays   int
		PresentDays int
		LateDays    int
		EarlyDays   int
	}

	var result Result
	err := r.db.Raw(`
		SELECT
			COUNT(DISTINCT date) as total_days,
			SUM(CASE WHEN status = 'normal' THEN 1 ELSE 0 END) as present_days,
			SUM(CASE WHEN status = 'late' THEN 1 ELSE 0 END) as late_days,
			SUM(CASE WHEN status = 'early' THEN 1 ELSE 0 END) as early_days
		FROM attendances
		WHERE staff_id = ? AND date >= ? AND date <= ?
	`, staffID, startDate, endDate).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	attendanceRate := 0.0
	if result.TotalDays > 0 {
		attendanceRate = float64(result.PresentDays) / float64(result.TotalDays) * 100
	}

	return map[string]interface{}{
		"total_days":     result.TotalDays,
		"present_days":   result.PresentDays,
		"late_days":      result.LateDays,
		"early_days":     result.EarlyDays,
		"attendance_rate": attendanceRate,
	}, nil
}

// Performance methods
func (r *AttendanceRepository) CreatePerformance(performance *model.Performance) error {
	return r.db.Create(performance).Error
}

func (r *AttendanceRepository) GetPerformance(staffID uint, year, month int) (*model.Performance, error) {
	var performance model.Performance
	err := r.db.Where("staff_id = ? AND year = ? AND month = ?", staffID, year, month).
		First(&performance).Error
	if err != nil {
		return nil, err
	}
	return &performance, nil
}

func (r *AttendanceRepository) ListPerformance(year, month int) ([]model.Performance, error) {
	var list []model.Performance
	query := r.db.Preload("Staff")
	if year > 0 && month > 0 {
		query = query.Where("year = ? AND month = ?", year, month)
	}
	err := query.Order("total_score DESC").Find(&list).Error
	return list, err
}

// Salary methods
func (r *AttendanceRepository) CreateSalary(salary *model.Salary) error {
	return r.db.Create(salary).Error
}

func (r *AttendanceRepository) GetSalary(staffID uint, year, month int) (*model.Salary, error) {
	var salary model.Salary
	err := r.db.Where("staff_id = ? AND year = ? AND month = ?", staffID, year, month).
		First(&salary).Error
	if err != nil {
		return nil, err
	}
	return &salary, nil
}

func (r *AttendanceRepository) ListSalaries(year, month int, page, pageSize int) ([]model.Salary, int64, error) {
	var salaries []model.Salary
	var total int64

	query := r.db.Preload("Staff").Model(&model.Salary{})
	if year > 0 {
		query = query.Where("year = ?", year)
	}
	if month > 0 {
		query = query.Where("month = ?", month)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("year DESC, month DESC").
		Find(&salaries).Error

	return salaries, total, err
}

func (r *AttendanceRepository) UpdateSalaryStatus(id uint, status string) error {
	return r.db.Model(&model.Salary{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// ShiftRule methods
func (r *AttendanceRepository) CreateShiftRule(rule *model.ShiftRule) error {
	return r.db.Create(rule).Error
}

func (r *AttendanceRepository) ListShiftRules() ([]model.ShiftRule, error) {
	var rules []model.ShiftRule
	err := r.db.Where("is_active = ?", true).
		Order("sort_order").
		Find(&rules).Error
	return rules, err
}

func (r *AttendanceRepository) GetShiftRule(id uint) (*model.ShiftRule, error) {
	var rule model.ShiftRule
	err := r.db.First(&rule, id).Error
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

func (r *AttendanceRepository) UpdateShiftRule(rule *model.ShiftRule) error {
	return r.db.Save(rule).Error
}

func (r *AttendanceRepository) DeleteShiftRule(id uint) error {
	return r.db.Delete(&model.ShiftRule{}, id).Error
}
