package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MedicationHandler 用药管理处理器
type MedicationHandler struct {
	DB *gorm.DB
}

// GetPrescriptions 获取处方列表
func (h *MedicationHandler) GetPrescriptions(c *gin.Context) {
	status := c.Query("status")
	elderlyId := c.Query("elderlyId")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "20")

	// TODO: 查询处方列表

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取成功",
		"data": gin.H{
			"list": []gin.H{
				{
					"id":             1,
					"prescriptionNo": "RX2026030001",
					"elderlyId":      1,
					"elderlyName":    "张奶奶",
					"bedNumber":      "3号楼201",
					"doctorId":       1,
					"doctorName":     "王医生",
					"prescriptionDate": "2026-03-01",
					"medications": []gin.H{
						{"name": "降压药", "dosage": "1片", "frequency": "一日3次", "days": 30},
						{"name": "钙片", "dosage": "1片", "frequency": "一日1次", "days": 30},
					},
					"status":     "active",
					"startDate":  "2026-03-01",
					"endDate":    "2026-03-31",
				},
			},
			"total":    page,
			"page":     pageSize,
		},
	})
}

// CreatePrescription 创建处方
func (h *MedicationHandler) CreatePrescription(c *gin.Context) {
	var req struct {
		ElderlyID        uint                   `json:"elderlyId" binding:"required"`
		PrescriptionDate string                 `json:"prescriptionDate" binding:"required"`
		Medications      []map[string]interface{} `json:"medications" binding:"required"`
		StartDate        string                 `json:"startDate" binding:"required"`
		EndDate          string                 `json:"endDate" binding:"required"`
		Remarks          string                 `json:"remarks"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	// TODO: 创建处方

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id":             1,
			"prescriptionNo": "RX2026030002",
			"createdAt":      time.Now().Format("2006-01-02 15:04:05"),
		},
	})
}

// GetPrescriptionDetail 获取处方详情
func (h *MedicationHandler) GetPrescriptionDetail(c *gin.Context) {
	id := c.Param("id")

	// TODO: 查询处方详情

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取成功",
		"data": gin.H{
			"id":              id,
			"prescriptionNo":  "RX2026030001",
			"elderlyName":     "张奶奶",
			"bedNumber":       "3号楼201",
			"doctorName":      "王医生",
			"prescriptionDate": "2026-03-01",
			"medications": []gin.H{
				{
					"name":      "降压药",
					"dosage":    "1片",
					"frequency": "一日3次",
					"days":      30,
					"total":     90,
				},
			},
			"status":    "active",
			"startDate": "2026-03-01",
			"endDate":   "2026-03-31",
		},
	})
}

// UpdatePrescription 更新处方
func (h *MedicationHandler) UpdatePrescription(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Medications []map[string]interface{} `json:"medications" binding:"required"`
		EndDate     string                 `json:"endDate"`
		Remarks     string                 `json:"remarks"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	// TODO: 更新处方

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
	})
}

// GetTodayMedications 获取今日用药
func (h *MedicationHandler) GetTodayMedications(c *gin.Context) {
	date := c.DefaultQuery("date", "")

	// TODO: 查询今日用药计划

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取成功",
		"data": gin.H{
			"list": []gin.H{
				{
					"id":           1,
					"elderlyName":  "张奶奶",
					"bedNumber":    "3号楼201",
					"medicationName": "降压药",
					"dosage":       "1片",
					"frequency":    "一日3次",
					"plannedTimes": []string{"08:00", "12:00", "18:00"},
					"status":       "pending",
				},
				{
					"id":           2,
					"elderlyName":  "王爷爷",
					"bedNumber":    "3号楼202",
					"medicationName": "降糖药",
					"dosage":       "2片",
					"frequency":    "一日2次",
					"plannedTimes": []string{"08:00", "20:00"},
					"status":       "completed",
				},
			},
			"summary": gin.H{
				"total":     145,
				"completed": 123,
				"pending":   22,
			},
		},
	})
}

// ConfirmMedication 确认服药
func (h *MedicationHandler) ConfirmMedication(c *gin.Context) {
	var req struct {
		MedicationID int     `json:"medicationId" binding:"required"`
		ActualDose   int     `json:"actualDose" binding:"required"`
		Result       string  `json:"result" binding:"required,oneof=normal refused vomited"`
		Remarks      string  `json:"remarks"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	// TODO: 确认服药
	// 记录服药情况

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "确认成功",
	})
}

// SkipMedication 跳过服药
func (h *MedicationHandler) SkipMedication(c *gin.Context) {
	var req struct {
		MedicationID int    `json:"medicationId" binding:"required"`
		Reason       string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	// TODO: 跳过服药

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "已跳过",
	})
}
