package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// VisitHandler 探视相关处理器
type VisitHandler struct {
	DB *gorm.DB
}

// GetAppointments 获取预约列表
func (h *VisitHandler) GetAppointments(c *gin.Context) {
	status := c.Query("status")
	elderlyId := c.Query("elderlyId")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	// TODO: 查询预约列表

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取成功",
		"data": gin.H{
			"list": []gin.H{
				{
					"id":            1,
					"appointmentNo": "VA202603001",
					"elderlyId":     1,
					"elderlyName":   "张奶奶",
					"bedNumber":     "3号楼201",
					"visitorName":   "王先生",
					"visitorPhone":  "13900139001",
					"relationship":  "子女",
					"visitType":     "现场探访",
					"visitDate":     "2026-03-05",
					"visitTime":     "09:00",
					"visitorCount":  2,
					"status":        "pending",
					"createdAt":     "2026-03-03 14:30:00",
				},
				{
					"id":            2,
					"appointmentNo": "VA202603002",
					"elderlyId":     2,
					"elderlyName":   "李爷爷",
					"bedNumber":     "2号楼105",
					"visitorName":   "赵女士",
					"visitorPhone":  "13800138002",
					"relationship":  "配偶",
					"visitType":     "视频探访",
					"visitDate":     "2026-03-05",
					"visitTime":     "14:00",
					"visitorCount":  1,
					"status":        "approved",
					"approvedAt":    "2026-03-03 15:00:00",
				},
			},
			"total":    45,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// CreateAppointment 创建预约
func (h *VisitHandler) CreateAppointment(c *gin.Context) {
	var req struct {
		ElderlyID    uint   `json:"elderlyId" binding:"required"`
		VisitorName  string `json:"visitorName" binding:"required"`
		VisitorPhone string `json:"visitorPhone" binding:"required"`
		Relationship string `json:"relationship" binding:"required"`
		VisitType    string `json:"visitType" binding:"required,oneof=现场探访 视频探访"`
		VisitDate    string `json:"visitDate" binding:"required"`
		VisitTime    string `json:"visitTime" binding:"required"`
		VisitorCount int    `json:"visitorCount" binding:"min=1,max=10"`
		Remarks      string `json:"remarks"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	// TODO: 创建预约
	// 检查时间段是否可用

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "预约成功",
		"data": gin.H{
			"id":            1,
			"appointmentNo": "VA202603003",
			"status":        "pending",
			"createdAt":     time.Now().Format("2006-01-02 15:04:05"),
		},
	})
}

// GetAppointmentDetail 获取预约详情
func (h *VisitHandler) GetAppointmentDetail(c *gin.Context) {
	id := c.Param("id")

	// TODO: 查询预约详情

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取成功",
		"data": gin.H{
			"id":            id,
			"appointmentNo": "VA202603001",
			"elderlyName":   "张奶奶",
			"bedNumber":     "3号楼201",
			"visitorName":   "王先生",
			"visitorPhone":  "13900139001",
			"relationship":  "子女",
			"visitType":     "现场探访",
			"visitDate":     "2026-03-05",
			"visitTime":     "09:00",
			"visitorCount":  2,
			"status":        "pending",
			"remarks":       "",
		},
	})
}

// UpdateAppointment 更新预约
func (h *VisitHandler) UpdateAppointment(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		VisitDate    string `json:"visitDate" binding:"required"`
		VisitTime    string `json:"visitTime" binding:"required"`
		VisitorCount int    `json:"visitorCount" binding:"min=1,max=10"`
		Remarks      string `json:"remarks"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	// TODO: 更新预约

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
	})
}

// ApproveAppointment 审核通过
func (h *VisitHandler) ApproveAppointment(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Remarks string `json:"remarks"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	// TODO: 审核通过
	// 发送通知给家属

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "已通过",
		"data": gin.H{
			"id":         id,
			"status":     "approved",
			"approvedAt": time.Now().Format("2006-01-02 15:04:05"),
		},
	})
}

// RejectAppointment 审核拒绝
func (h *VisitHandler) RejectAppointment(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	// TODO: 审核拒绝
	// 发送通知给家属

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "已拒绝",
		"data": gin.H{
			"id":        id,
			"status":    "rejected",
			"rejectReason": req.Reason,
		},
	})
}

// CancelAppointment 取消预约
func (h *VisitHandler) CancelAppointment(c *gin.Context) {
	id := c.Param("id")

	// TODO: 取消预约

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "已取消",
		"data": gin.H{
			"id":     id,
			"status": "cancelled",
		},
	})
}

// CompleteAppointment 完成探访
func (h *VisitHandler) CompleteAppointment(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		ActualTime string `json:"actualTime"`
		Remarks    string `json:"remarks"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	// TODO: 完成探访

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "已完成",
		"data": gin.H{
			"id":         id,
			"status":     "completed",
			"completedAt": time.Now().Format("2006-01-02 15:04:05"),
		},
	})
}
