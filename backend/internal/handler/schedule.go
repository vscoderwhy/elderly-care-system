package handler

import (
	"elderly-care-system/internal/service"
	"elderly-care-system/pkg/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ScheduleHandler struct {
	scheduleService *service.ScheduleService
}

func NewScheduleHandler(scheduleService *service.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{scheduleService: scheduleService}
}

// CreateSchedule 创建排班
func (h *ScheduleHandler) CreateSchedule(c *gin.Context) {
	var req service.CreateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	schedule, err := h.scheduleService.CreateSchedule(&req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, schedule)
}

// GetScheduleList 获取排班列表
func (h *ScheduleHandler) GetScheduleList(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			response.Error(c, 400, "Invalid start_date format")
			return
		}
	} else {
		startDate = time.Now().AddDate(0, 0, -7)
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			response.Error(c, 400, "Invalid end_date format")
			return
		}
	} else {
		endDate = time.Now().AddDate(0, 0, 7)
	}

	schedules, total, err := h.scheduleService.GetScheduleList(startDate, endDate, page, pageSize)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":       schedules,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
	})
}

// GetStaffSchedule 获取员工排班
func (h *ScheduleHandler) GetStaffSchedule(c *gin.Context) {
	staffID, _ := strconv.ParseUint(c.Param("staff_id"), 10, 32)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	schedules, total, err := h.scheduleService.GetStaffSchedule(uint(staffID), page, pageSize)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":       schedules,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
	})
}

// GetMySchedule 获取我的排班
func (h *ScheduleHandler) GetMySchedule(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, 401, "Unauthorized")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	schedules, total, err := h.scheduleService.GetStaffSchedule(userID.(uint), page, pageSize)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":       schedules,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
	})
}

// UpdateScheduleStatus 更新排班状态
func (h *ScheduleHandler) UpdateScheduleStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request")
		return
	}

	if err := h.scheduleService.UpdateScheduleStatus(uint(id), req.Status); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}

// DeleteSchedule 删除排班
func (h *ScheduleHandler) DeleteSchedule(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := h.scheduleService.DeleteSchedule(uint(id)); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetMonthlyStats 获取月度排班统计
func (h *ScheduleHandler) GetMonthlyStats(c *gin.Context) {
	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(time.Now().Year())))
	month, _ := strconv.Atoi(c.DefaultQuery("month", strconv.Itoa(int(time.Now().Month()))))

	stats, err := h.scheduleService.GetMonthlyStats(year, time.Month(month))
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, stats)
}
