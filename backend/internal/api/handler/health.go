package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthHandler 健康数据处理器
type HealthHandler struct {
	DB *gorm.DB
}

// HealthRecord 健康记录
type HealthRecord struct {
	BloodPressure string  `json:"bloodPressure"`
	HeartRate     int     `json:"heartRate"`
	Temperature   float64 `json:"temperature"`
	BloodSugar    float64 `json:"bloodSugar"`
	Weight        float64 `json:"weight"`
	SPO2          int     `json:"spo2"`
}

// GetRecords 获取健康记录列表
func (h *HealthHandler) GetRecords(c *gin.Context) {
	elderlyID := c.Query("elderlyId")
	recordType := c.Query("type")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	// TODO: 查询健康记录

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取成功",
		"data": gin.H{
			"list": []gin.H{
				{
					"id":          1,
					"elderlyId":   1,
					"elderlyName": "张奶奶",
					"recordDate":  "2026-03-04",
					"recordTime":  "08:30",
					"bloodPressure": "130/85",
					"heartRate":     72,
					"temperature":   36.4,
					"bloodSugar":    6.8,
					"weight":        55,
					"spo2":          98,
					"recorder":      "赵护士",
					"remarks":       "血压正常，继续监测",
				},
			},
			"total":    156,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// CreateRecord 创建健康记录
func (h *HealthHandler) CreateRecord(c *gin.Context) {
	var req struct {
		ElderlyID     uint          `json:"elderlyId" binding:"required"`
		RecordDate    string        `json:"recordDate" binding:"required"`
		RecordTime    string        `json:"recordTime" binding:"required"`
		BloodPressure string        `json:"bloodPressure"`
		HeartRate     int           `json:"heartRate"`
		Temperature   float64       `json:"temperature"`
		BloodSugar    float64       `json:"bloodSugar"`
		Weight        float64       `json:"weight"`
		SPO2          int           `json:"spo2"`
		Remarks       string        `json:"remarks"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	// TODO: 创建健康记录
	// 检查异常值并触发告警

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id":        1,
			"createdAt": time.Now().Format("2006-01-02 15:04:05"),
		},
	})
}

// GetRecordDetail 获取健康记录详情
func (h *HealthHandler) GetRecordDetail(c *gin.Context) {
	id := c.Param("id")

	// TODO: 查询健康记录详情

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取成功",
		"data": gin.H{
			"id":            id,
			"elderlyName":   "张奶奶",
			"recordDate":    "2026-03-04",
			"recordTime":    "08:30",
			"bloodPressure": "130/85",
			"heartRate":     72,
			"temperature":   36.4,
			"bloodSugar":    6.8,
			"weight":        55,
			"spo2":          98,
			"recorder":      "赵护士",
			"remarks":       "血压正常，继续监测",
		},
	})
}

// UpdateRecord 更新健康记录
func (h *HealthHandler) UpdateRecord(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		BloodPressure string  `json:"bloodPressure"`
		HeartRate     int     `json:"heartRate"`
		Temperature   float64 `json:"temperature"`
		BloodSugar    float64 `json:"bloodSugar"`
		Weight        float64 `json:"weight"`
		SPO2          int     `json:"spo2"`
		Remarks       string  `json:"remarks"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	// TODO: 更新健康记录

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
	})
}

// GetStatistics 获取健康统计数据
func (h *HealthHandler) GetStatistics(c *gin.Context) {
	elderlyID := c.Query("elderlyId")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	// TODO: 查询健康统计数据

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取成功",
		"data": gin.H{
			"summary": gin.H{
				"avgBloodPressure": "128/82",
				"avgHeartRate":     74,
				"avgTemperature":   36.5,
				"avgBloodSugar":    6.5,
				"avgWeight":        54.5,
			},
			"trend": gin.H{
				"bloodPressure": []gin.H{
					{"date": "2026-03-01", "value": "130/85"},
					{"date": "2026-03-02", "value": "128/82"},
					{"date": "2026-03-03", "value": "125/80"},
				},
				"heartRate": []gin.H{
					{"date": "2026-03-01", "value": 75},
					{"date": "2026-03-02", "value": 73},
					{"date": "2026-03-03", "value": 72},
				},
			},
		},
	})
}

// GetAbnormalData 获取异常数据
func (h *HealthHandler) GetAbnormalData(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	// TODO: 查询异常健康数据

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取成功",
		"data": gin.H{
			"list": []gin.H{
				{
					"id":          1,
					"elderlyName": "李爷爷",
					"bedNumber":   "2号楼105",
					"abnormalType": "血压偏高",
					"value":       "150/95mmHg",
					"recordTime":  "2026-03-04 09:00",
					"level":       "中度",
				},
			},
			"total":    8,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}
