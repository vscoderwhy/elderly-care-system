package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CareHandler 护理相关处理器
type CareHandler struct {
	DB *gorm.DB
}

// GetTasks 获取护理任务列表
func (h *CareHandler) GetTasks(c *gin.Context) {
	status := c.Query("status")
	taskType := c.Query("taskType")
	nurseID := c.Query("nurseId")
	date := c.Query("date")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	// TODO: 构建查询条件并查询任务

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取成功",
		"data": gin.H{
			"list": []gin.H{
				{
					"id":            1,
					"elderlyId":     1,
					"elderlyName":   "张奶奶",
					"bedNumber":     "3号楼201",
					"taskType":      "日常护理",
					"description":   "协助洗漱、测量血压、体温、心率",
					"scheduledDate": "2026-03-04",
					"scheduledTime": "08:00",
					"nurseId":       1,
					"nurseName":     "赵护士",
					"status":        "pending",
					"priority":      "normal",
					"createdAt":     "2026-03-03 16:00:00",
				},
				{
					"id":            2,
					"elderlyId":     2,
					"elderlyName":   "李爷爷",
					"bedNumber":     "2号楼105",
					"taskType":      "健康监测",
					"description":   "测量血压、血糖、体温",
					"scheduledDate": "2026-03-04",
					"scheduledTime": "09:00",
					"nurseId":       2,
					"nurseName":     "李护士",
					"status":        "in_progress",
					"priority":      "important",
					"createdAt":     "2026-03-03 16:00:00",
				},
				{
					"id":            3,
					"elderlyId":     3,
					"elderlyName":   "王奶奶",
					"bedNumber":     "3号楼202",
					"taskType":      "康复训练",
					"description":   "上肢关节活动训练30分钟",
					"scheduledDate": "2026-03-04",
					"scheduledTime": "14:00",
					"nurseId":       3,
					"nurseName":     "陈康复师",
					"status":        "completed",
					"priority":      "normal",
					"completedAt":   "2026-03-04 14:35:00",
				},
			},
			"total":    45,
			"page":     page,
			"pageSize": pageSize,
			"summary": gin.H{
				"todayTotal":     45,
				"pending":        12,
				"inProgress":     8,
				"completed":      25,
			},
		},
	})
}

// CreateTask 创建护理任务
func (h *CareHandler) CreateTask(c *gin.Context) {
	var req struct {
		ElderlyID     uint   `json:"elderlyId" binding:"required"`
		TaskType      string `json:"taskType" binding:"required"`
		Description   string `json:"description" binding:"required"`
		ScheduledDate string `json:"scheduledDate" binding:"required"`
		ScheduledTime string `json:"scheduledTime" binding:"required"`
		NurseID       uint   `json:"nurseId" binding:"required"`
		Priority      string `json:"priority" binding:"oneof=normal important urgent"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	// TODO: 创建护理任务

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id":        1,
			"createdAt": time.Now().Format("2006-01-02 15:04:05"),
		},
	})
}

// GetTaskDetail 获取任务详情
func (h *CareHandler) GetTaskDetail(c *gin.Context) {
	id := c.Param("id")

	// TODO: 查询任务详情

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取成功",
		"data": gin.H{
			"id":            id,
			"elderlyName":   "张奶奶",
			"bedNumber":     "3号楼201",
			"taskType":      "日常护理",
			"description":   "协助洗漱、测量血压、体温、心率、协助用餐",
			"scheduledDate": "2026-03-04",
			"scheduledTime": "08:00",
			"nurseName":     "赵护士",
			"status":        "pending",
			"priority":      "normal",
			"createdAt":     "2026-03-03 16:00:00",
			"elderly": gin.H{
				"name":      "张奶奶",
				"age":       78,
				"careLevel": "二级",
				"phone":     "138****8888",
			},
		},
	})
}

// UpdateTaskStatus 更新任务状态
func (h *CareHandler) UpdateTaskStatus(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Status string `json:"status" binding:"required,oneof=pending in_progress completed cancelled"`
		Notes  string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	// TODO: 更新任务状态
	// 如果状态为completed，记录完成时间

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
		"data": gin.H{
			"id":         id,
			"status":     req.Status,
			"updatedAt":  time.Now().Format("2006-01-02 15:04:05"),
		},
	})
}

// GetRecords 获取护理记录列表
func (h *CareHandler) GetRecords(c *gin.Context) {
	elderlyID := c.Query("elderlyId")
	careType := c.Query("careType")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	// TODO: 查询护理记录

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取成功",
		"data": gin.H{
			"list": []gin.H{
				{
					"id":         1,
					"elderlyId":  1,
					"elderlyName": "张奶奶",
					"bedNumber":  "3号楼201",
					"careType":   "日常护理",
					"content":    "协助老人起床洗漱，测量血压130/85mmHg、心率72次/分、体温36.4℃，协助用餐，老人状态良好",
					"careTime":   "2026-03-04 08:30",
					"nurseId":    1,
					"nurseName":  "赵护士",
					"images":     []string{"/uploads/care1.jpg", "/uploads/care2.jpg"},
					"evaluation": 5,
					"tags":       []string{"服务热情", "专业细致"},
					"createdAt":  "2026-03-04 08:45:00",
				},
				{
					"id":         2,
					"elderlyId":  2,
					"elderlyName": "李爷爷",
					"bedNumber":  "2号楼105",
					"careType":   "康复训练",
					"content":    "上肢关节活动训练30分钟，包括屈伸、旋转等动作，老人配合度良好",
					"careTime":   "2026-03-04 14:00",
					"nurseId":    3,
					"nurseName":  "陈康复师",
					"images":     []string{},
					"evaluation": 4,
					"createdAt":  "2026-03-04 14:35:00",
				},
			},
			"total":    156,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// CreateRecord 创建护理记录
func (h *CareHandler) CreateRecord(c *gin.Context) {
	var req struct {
		ElderlyID  uint     `json:"elderlyId" binding:"required"`
		CareType   string   `json:"careType" binding:"required"`
		Content    string   `json:"content" binding:"required"`
		CareTime   string   `json:"careTime" binding:"required"`
		Images     []string `json:"images"`
		Evaluation int      `json:"evaluation" binding:"min=0,max=5"`
		Tags       []string `json:"tags"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	// TODO: 创建护理记录
	// 自动关联到当前登录的护工

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id":        1,
			"createdAt": time.Now().Format("2006-01-02 15:04:05"),
		},
	})
}

// GetRecordDetail 获取护理记录详情
func (h *CareHandler) GetRecordDetail(c *gin.Context) {
	id := c.Param("id")

	// TODO: 查询护理记录详情

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取成功",
		"data": gin.H{
			"id":         id,
			"elderlyName": "张奶奶",
			"bedNumber":  "3号楼201",
			"careType":   "日常护理",
			"content":    "协助老人起床洗漱，测量血压130/85mmHg、心率72次/分、体温36.4℃，协助用餐，老人状态良好",
			"careTime":   "2026-03-04 08:30",
			"nurseName":  "赵护士",
			"images":     []string{"/uploads/care1.jpg"},
			"evaluation": 5,
			"tags":       []string{"服务热情", "专业细致"},
			"createdAt":  "2026-03-04 08:45:00",
		},
	})
}

// GetAssessments 获取护理评估列表
func (h *CareHandler) GetAssessments(c *gin.Context) {
	elderlyID := c.Query("elderlyId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	// TODO: 查询护理评估

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取成功",
		"data": gin.H{
			"list": []gin.H{
				{
					"id":          1,
					"elderlyId":   1,
					"elderlyName": "张奶奶",
					"assessmentDate": "2026-03-01",
					"careLevel":      "二级",
					"adlScore":       75,
					"healthScore":    85,
					"nurseName":      "王护士长",
					"remarks":         "老人整体状况良好，建议维持当前护理等级",
					"createdAt":       "2026-03-01 10:00:00",
				},
			},
			"total":    48,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// CreateAssessment 创建护理评估
func (h *CareHandler) CreateAssessment(c *gin.Context) {
	var req struct {
		ElderlyID        uint    `json:"elderlyId" binding:"required"`
		AssessmentDate   string  `json:"assessmentDate" binding:"required"`
		CareLevel        string  `json:"careLevel" binding:"required"`
		ADLScore         int     `json:adlScore binding:"min=0,max=100"`
		HealthScore      int     `json:"healthScore binding:"min=0,max=100"`
		PhysicalCondition string `json:"physicalCondition"`
		MentalCondition   string `json:"mentalCondition"`
		SocialAbility     string `json:"socialAbility"`
		RiskFactors       []string `json:"riskFactors"`
		Remarks           string `json:"remarks"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	// TODO: 创建护理评估

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id":        1,
			"createdAt": time.Now().Format("2006-01-02 15:04:05"),
		},
	})
}

// GetStatistics 获取护理统计数据
func (h *CareHandler) GetStatistics(c *gin.Context) {
	startDate := c.DefaultQuery("startDate", "")
	endDate := c.DefaultQuery("endDate", "")

	// TODO: 查询护理统计数据

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取成功",
		"data": gin.H{
			"summary": gin.H{
				"totalTasks":      156,
				"completedTasks":  142,
				"pendingTasks":    12,
				"completionRate":  91.0,
				"totalRecords":    245,
				"avgEvaluation":   4.7,
			},
			"byType": []gin.H{
				{"type": "日常护理", "count": 85, "percent": 54.5},
				{"type": "康复训练", "count": 35, "percent": 22.4},
				{"type": "健康监测", "count": 28, "percent": 17.9},
				{"type": "医疗护理", "count": 8, "percent": 5.1},
			},
			"byNurse": []gin.H{
				{"nurseName": "赵护士", "completed": 45, "pending": 5},
				{"nurseName": "李护士", "completed": 42, "pending": 3},
				{"nurseName": "王护士", "completed": 38, "pending": 2},
			},
		},
	})
}
