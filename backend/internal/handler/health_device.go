package handler

import (
	"elderly-care-system/internal/service"
	"elderly-care-system/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HealthDeviceHandler struct {
	healthDeviceService *service.HealthDeviceService
}

func NewHealthDeviceHandler(healthDeviceService *service.HealthDeviceService) *HealthDeviceHandler {
	return &HealthDeviceHandler{healthDeviceService: healthDeviceService}
}

// BindDevice 绑定设备
func (h *HealthDeviceHandler) BindDevice(c *gin.Context) {
	var req service.BindDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request")
		return
	}

	device, err := h.healthDeviceService.BindDevice(&req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, device)
}

// ReceiveDeviceData 接收设备数据
func (h *HealthDeviceHandler) ReceiveDeviceData(c *gin.Context) {
	var req service.DeviceDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request")
		return
	}

	data, err := h.healthDeviceService.ReceiveDeviceData(&req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, data)
}

// ListDevices 获取设备列表
func (h *HealthDeviceHandler) ListDevices(c *gin.Context) {
	elderlyIDStr := c.Query("elderly_id")
	var elderlyID *uint
	if elderlyIDStr != "" {
		id, _ := strconv.ParseUint(elderlyIDStr, 10, 32)
		eid := uint(id)
		elderlyID = &eid
	}

	devices, err := h.healthDeviceService.ListDevices(elderlyID)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, devices)
}

// GetDeviceTrends 获取设备数据趋势
func (h *HealthDeviceHandler) GetDeviceTrends(c *gin.Context) {
	deviceID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	dataType := c.Query("data_type")
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))

	trends, err := h.healthDeviceService.GetDeviceTrends(uint(deviceID), dataType, days)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, trends)
}

// GetAbnormalData 获取异常数据
func (h *HealthDeviceHandler) GetAbnormalData(c *gin.Context) {
	elderlyID, _ := strconv.ParseUint(c.Query("elderly_id"), 10, 32)
	hours, _ := strconv.Atoi(c.DefaultQuery("hours", "24"))

	data, err := h.healthDeviceService.GetAbnormalData(uint(elderlyID), hours)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, data)
}

// UnbindDevice 解绑设备
func (h *HealthDeviceHandler) UnbindDevice(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := h.healthDeviceService.UnbindDevice(uint(id)); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "设备已解绑"})
}
