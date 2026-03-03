package handler

import (
	"elderly-care-system/internal/service"
	"elderly-care-system/pkg/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AttendanceHandler struct {
	attendanceService *service.AttendanceService
}

func NewAttendanceHandler(attendanceService *service.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{attendanceService: attendanceService}
}

// ClockIn 上班打卡
func (h *AttendanceHandler) ClockIn(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req service.ClockInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request")
		return
	}

	req.StaffID = userID.(uint)

	attendance, err := h.attendanceService.ClockIn(&req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, attendance)
}

// ClockOut 下班打卡
func (h *AttendanceHandler) ClockOut(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req service.ClockInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request")
		return
	}

	req.StaffID = userID.(uint)

	attendance, err := h.attendanceService.ClockOut(&req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, attendance)
}

// GetTodayAttendance 获取今日打卡状态
func (h *AttendanceHandler) GetTodayAttendance(c *gin.Context) {
	// TODO: 实现获取今日打卡状态
	response.Success(c, gin.H{
		"clock_in":    nil,
		"clock_out":   nil,
		"can_clock_in": true,
		"can_clock_out": false,
	})
}

// GetAttendanceStats 获取考勤统计
func (h *AttendanceHandler) GetAttendanceStats(c *gin.Context) {
	userID, _ := c.Get("user_id")

	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(time.Now().Year())))
	month, _ := strconv.Atoi(c.DefaultQuery("month", strconv.Itoa(int(time.Now().Month()))))

	stats, err := h.attendanceService.GetAttendanceStats(userID.(uint), year, month)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, stats)
}

// GetPerformance 获取绩效
func (h *AttendanceHandler) GetPerformance(c *gin.Context) {
	userID, _ := c.Get("user_id")

	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(time.Now().Year())))
	month, _ := strconv.Atoi(c.DefaultQuery("month", strconv.Itoa(int(time.Now().Month()))))

	performance, err := h.attendanceService.GeneratePerformance(userID.(uint), year, month)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, performance)
}

// ListPerformance 列出绩效排行榜
func (h *AttendanceHandler) ListPerformance(c *gin.Context) {
	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(time.Now().Year())))
	month, _ := strconv.Atoi(c.DefaultQuery("month", strconv.Itoa(int(time.Now().Month()))))

	// TODO: 从service获取列表
	response.Success(c, gin.H{
		"list": []interface{}{},
		"year": year,
		"month": month,
	})
}

// GetSalary 获取工资
func (h *AttendanceHandler) GetSalary(c *gin.Context) {
	userID, _ := c.Get("user_id")

	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(time.Now().Year())))
	month, _ := strconv.Atoi(c.DefaultQuery("month", strconv.Itoa(int(time.Now().Month()))))

	salary, err := h.attendanceService.CalculateSalary(userID.(uint), year, month)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, salary)
}

// ListSalaries 工资列表(管理员)
func (h *AttendanceHandler) ListSalaries(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	year, _ := strconv.Atoi(c.DefaultQuery("year", "0"))
	month, _ := strconv.Atoi(c.DefaultQuery("month", "0"))

	// TODO: 从service获取列表
	response.Success(c, gin.H{
		"list":      []interface{}{},
		"total":     0,
		"page":      page,
		"page_size": pageSize,
		"year":      year,
		"month":     month,
	})
}

// ListShiftRules 获取排班规则
func (h *AttendanceHandler) ListShiftRules(c *gin.Context) {
	rules, err := h.attendanceService.ListShiftRules()
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, rules)
}

// CreateShiftRule 创建排班规则
func (h *AttendanceHandler) CreateShiftRule(c *gin.Context) {
	var req service.ShiftRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request")
		return
	}

	rule, err := h.attendanceService.CreateShiftRule(&req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, rule)
}

// GetMySchedule 获取我的排班
func (h *AttendanceHandler) GetMySchedule(c *gin.Context) {
	userID, _ := c.Get("user_id")

	startDate := c.DefaultQuery("start_date", time.Now().Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().AddDate(0, 0, 7).Format("2006-01-02"))

	// TODO: 从schedule repository获取
	response.Success(c, gin.H{
		"staff_id":   userID,
		"start_date": startDate,
		"end_date":   endDate,
		"schedules":  []interface{}{},
	})
}
