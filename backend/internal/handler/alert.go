package handler

import (
	"elderly-care-system/internal/service"
	"elderly-care-system/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AlertHandler struct {
	alertService *service.AlertService
}

func NewAlertHandler(alertService *service.AlertService) *AlertHandler {
	return &AlertHandler{alertService: alertService}
}

// CreateAlert 创建预警
func (h *AlertHandler) CreateAlert(c *gin.Context) {
	var req service.AlertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	alert, err := h.alertService.CreateAlert(&req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, alert)
}

// AcknowledgeAlert 确认预警
func (h *AlertHandler) AcknowledgeAlert(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	userID, _ := c.Get("user_id")
	var uid uint
	if userID != nil {
		uid = userID.(uint)
	}

	if err := h.alertService.AcknowledgeAlert(uint(id), uid); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}

// ResolveAlert 解决预警
func (h *AlertHandler) ResolveAlert(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := h.alertService.ResolveAlert(uint(id)); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetAlert 获取预警详情
func (h *AlertHandler) GetAlert(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	alert, err := h.alertService.GetAlert(uint(id))
	if err != nil {
		response.Error(c, 404, "预警不存在")
		return
	}

	response.Success(c, alert)
}

// ListAlerts 获取预警列表
func (h *AlertHandler) ListAlerts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.DefaultQuery("status", "")

	alerts, total, err := h.alertService.ListAlerts(page, pageSize, status)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":      alerts,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetActiveAlerts 获取活跃预警
func (h *AlertHandler) GetActiveAlerts(c *gin.Context) {
	alerts, err := h.alertService.GetActiveAlerts()
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, alerts)
}

// CheckAlerts 手动触发预警检查
func (h *AlertHandler) CheckAlerts(c *gin.Context) {
	if err := h.alertService.CheckAllAlerts(); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "预警检查完成",
	})
}

// GetAlertSummary 获取预警统计摘要
func (h *AlertHandler) GetAlertSummary(c *gin.Context) {
	summary, err := h.alertService.GetAlertSummary()
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, summary)
}

// Alert Rule Management
func (h *AlertHandler) CreateRule(c *gin.Context) {
	var req service.AlertRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	rule, err := h.alertService.CreateRule(&req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, rule)
}

func (h *AlertHandler) ListRules(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	rules, total, err := h.alertService.ListRules(page, pageSize)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":      rules,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (h *AlertHandler) UpdateRule(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.Error(c, 400, "Invalid request")
		return
	}

	if err := h.alertService.UpdateRule(uint(id), updates); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *AlertHandler) DeleteRule(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := h.alertService.DeleteRule(uint(id)); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}
