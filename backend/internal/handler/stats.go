package handler

import (
	"elderly-care-system/internal/service"
	"elderly-care-system/pkg/response"

	"github.com/gin-gonic/gin"
)

type StatsHandler struct {
	statsService *service.StatsService
}

func NewStatsHandler(statsService *service.StatsService) *StatsHandler {
	return &StatsHandler{statsService: statsService}
}

// GetDashboardStats 获取Dashboard统计数据
func (h *StatsHandler) GetDashboardStats(c *gin.Context) {
	stats, err := h.statsService.GetDashboardStats()
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, stats)
}

// GetBedOccupancy 获取床位入住率
func (h *StatsHandler) GetBedOccupancy(c *gin.Context) {
	stats, err := h.statsService.GetBedOccupancy()
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, stats)
}

// GetCareStats 获取护理统计
func (h *StatsHandler) GetCareStats(c *gin.Context) {
	stats, err := h.statsService.GetCareStats()
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, stats)
}

// GetFinanceStats 获取财务统计
func (h *StatsHandler) GetFinanceStats(c *gin.Context) {
	stats, err := h.statsService.GetFinanceStats()
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, stats)
}
