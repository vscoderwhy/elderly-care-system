package handler

import (
	"elderly-care-system/internal/service"
	"elderly-care-system/pkg/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type VisitHandler struct {
	visitService *service.VisitService
}

func NewVisitHandler(visitService *service.VisitService) *VisitHandler {
	return &VisitHandler{visitService: visitService}
}

// CreateVisit 创建探视预约
func (h *VisitHandler) CreateVisit(c *gin.Context) {
	var req service.VisitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	visit, err := h.visitService.CreateVisit(&req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, visit)
}

// UpdateVisit 更新探视预约
func (h *VisitHandler) UpdateVisit(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req service.UpdateVisitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	if err := h.visitService.UpdateVisit(uint(id), &req); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}

// DeleteVisit 删除探视预约
func (h *VisitHandler) DeleteVisit(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := h.visitService.DeleteVisit(uint(id)); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetVisit 获取探视预约详情
func (h *VisitHandler) GetVisit(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	visit, err := h.visitService.GetVisit(uint(id))
	if err != nil {
		response.Error(c, 404, "探视预约不存在")
		return
	}

	response.Success(c, visit)
}

// ListVisits 获取探视预约列表
func (h *VisitHandler) ListVisits(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	visits, total, err := h.visitService.ListVisits(page, pageSize)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":      visits,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// ListElderlyVisits 获取老人的探视预约列表
func (h *VisitHandler) ListElderlyVisits(c *gin.Context) {
	elderlyID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	visits, total, err := h.visitService.ListElderlyVisits(uint(elderlyID), page, pageSize)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":      visits,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// ConfirmVisit 确认探视预约
func (h *VisitHandler) ConfirmVisit(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := h.visitService.ConfirmVisit(uint(id)); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}

// CancelVisit 取消探视预约
func (h *VisitHandler) CancelVisit(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := h.visitService.CancelVisit(uint(id)); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}

// CompleteVisit 完成探视
func (h *VisitHandler) CompleteVisit(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := h.visitService.CompleteVisit(uint(id)); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetTodayVisits 获取今日探视预约
func (h *VisitHandler) GetTodayVisits(c *gin.Context) {
	visits, err := h.visitService.GetTodayVisits()
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, visits)
}

// GetUpcomingVisits 获取即将到来的探视预约
func (h *VisitHandler) GetUpcomingVisits(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))

	visits, err := h.visitService.GetUpcomingVisits(days)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, visits)
}

// GetVisitsByDateRange 获取指定日期范围内的探视预约
func (h *VisitHandler) GetVisitsByDateRange(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	// 默认查询本月
	if startDate == "" {
		now := time.Now()
		startDate = now.Format("2006-01") + "-01"
	}
	if endDate == "" {
		now := time.Now()
		endDate = now.AddDate(0, 1, -1).Format("2006-01-02")
	}

	visits, err := h.visitService.GetVisitsByDateRange(startDate, endDate)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, visits)
}
