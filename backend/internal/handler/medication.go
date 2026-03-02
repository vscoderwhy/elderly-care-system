package handler

import (
	"elderly-care-system/internal/service"
	"elderly-care-system/pkg/export"
	"elderly-care-system/pkg/response"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type MedicationHandler struct {
	medicationService *service.MedicationService
}

func NewMedicationHandler(medicationService *service.MedicationService) *MedicationHandler {
	return &MedicationHandler{medicationService: medicationService}
}

// CreateMedication 创建药品
func (h *MedicationHandler) CreateMedication(c *gin.Context) {
	var req service.MedicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	medication, err := h.medicationService.CreateMedication(&req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, medication)
}

// ListMedications 获取药品列表
func (h *MedicationHandler) ListMedications(c *gin.Context) {
	keyword := c.Query("keyword")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	medications, total, err := h.medicationService.ListMedications(keyword, page, pageSize)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":       medications,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
	})
}

// UpdateMedication 更新药品
func (h *MedicationHandler) UpdateMedication(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.Error(c, 400, "Invalid request")
		return
	}

	if err := h.medicationService.UpdateMedication(uint(id), updates); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}

// DeleteMedication 删除药品
func (h *MedicationHandler) DeleteMedication(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := h.medicationService.DeleteMedication(uint(id)); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}

// CreateMedicationRecord 创建用药记录
func (h *MedicationHandler) CreateMedicationRecord(c *gin.Context) {
	var req service.MedicationRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	record, err := h.medicationService.CreateMedicationRecord(&req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, record)
}

// ListElderlyMedications 获取老人用药列表
func (h *MedicationHandler) ListElderlyMedications(c *gin.Context) {
	elderlyID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	records, total, err := h.medicationService.ListElderlyMedications(uint(elderlyID), page, pageSize)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":       records,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
	})
}

// GetTodayMedications 获取今日用药任务
func (h *MedicationHandler) GetTodayMedications(c *gin.Context) {
	elderlyID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	logs, err := h.medicationService.GetTodayMedications(uint(elderlyID))
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, logs)
}

// GetMedicationAlerts 获取用药预警
func (h *MedicationHandler) GetMedicationAlerts(c *gin.Context) {
	alerts, err := h.medicationService.GetMedicationAlerts()
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, alerts)
}

// CompleteMedicationLog 完成用药
func (h *MedicationHandler) CompleteMedicationLog(c *gin.Context) {
	logID, _ := strconv.ParseUint(c.Param("log_id"), 10, 32)

	userID, _ := c.Get("user_id")

	var req struct {
		Notes string `json:"notes"`
	}
	c.ShouldBindJSON(&req)

	if err := h.medicationService.CompleteMedicationLog(uint(logID), userID.(uint), req.Notes); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}

// ExportMedicationRecords 导出用药记录
func (h *MedicationHandler) ExportMedicationRecords(c *gin.Context) {
	elderlyIDStr := c.Query("elderly_id")
	_ = c.Query("start_date")
	_ = c.Query("end_date")

	// 获取所有用药记录或特定老人的用药记录
	_ = uint(0)
	if elderlyIDStr != "" {
		_, _ = strconv.ParseUint(elderlyIDStr, 10, 32)
	}

	// 这里应该调用service获取用药记录
	// 简化实现：返回空列表
	report := make([]export.MedicationReport, 0)

	data, err := export.ExportMedicationRecordsCSV(report)
	if err != nil {
		response.Error(c, 500, "生成报表失败")
		return
	}

	filename := fmt.Sprintf("用药记录_%s.csv", time.Now().Format("20060102150405"))
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Data(200, "text/csv; charset=utf-8", data)
}
