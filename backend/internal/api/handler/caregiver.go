package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CaregiverHandler 护工相关处理器
type CaregiverHandler struct {
	DB *gorm.DB
}

// LoginRequest 护工登录请求
type LoginRequest struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required,len=6"`
}

// WechatLoginRequest 微信登录请求
type WechatLoginRequest struct {
	Code     string `json:"code" binding:"required"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// CaregiverLoginResponse 登录响应
type CaregiverLoginResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	Avatar      string `json:"avatar"`
	EmployeeNo  string `json:"employeeNo"`
	Department  string `json:"department"`
	Position    string `json:"position"`
	Token       string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

// Login 护工手机号登录
func (h *CaregiverHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	// TODO: 验证验证码
	// 这里应该验证短信验证码是否正确

	// TODO: 查询护工信息
	// 根据手机号查询护工，如果不存在则返回错误

	// TODO: 生成JWT token
	// 使用jwt包生成token和refreshToken

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "登录成功",
		"data": CaregiverLoginResponse{
			ID:          "CG001",
			Name:        "张护工",
			Phone:       req.Phone,
			Avatar:      "",
			EmployeeNo:  "EMP001",
			Department:  "护理部",
			Position:    "护理员",
			Token:       "mock_token_" + strconv.FormatInt(time.Now().Unix(), 10),
			RefreshToken: "mock_refresh_token",
		},
	})
}

// WechatLogin 护工微信登录
func (h *CaregiverHandler) WechatLogin(c *gin.Context) {
	var req WechatLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	// TODO: 用code换取微信openid和session_key
	// 调用微信API获取用户信息

	// TODO: 查询或创建护工账号

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "登录成功",
		"data": CaregiverLoginResponse{
			ID:          "CG002",
			Name:        req.Nickname,
			Phone:       "",
			Avatar:      req.Avatar,
			EmployeeNo:  "EMP002",
			Department:  "护理部",
			Position:    "护理员",
			Token:       "mock_wechat_token_" + strconv.FormatInt(time.Now().Unix(), 10),
			RefreshToken: "mock_refresh_token",
		},
	})
}

// GetInfo 获取护工信息
func (h *CaregiverHandler) GetInfo(c *gin.Context) {
	// 从上下文获取用户ID
	userID := c.GetString("user_id")

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取成功",
		"data": gin.H{
			"id":          userID,
			"name":        "张护工",
			"phone":       "138****8888",
			"avatar":      "",
			"employeeNo":  "EMP001",
			"department":  "护理部",
			"position":    "护理员",
			"hireDate":    "2023-03-15",
			"status":      "active",
		},
	})
}

// GetTasks 获取护工任务列表
func (h *CaregiverHandler) GetTasks(c *gin.Context) {
	// 获取查询参数
	status := c.DefaultQuery("status", "")
	date := c.DefaultQuery("date", "")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	// TODO: 从数据库查询任务列表

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取成功",
		"data": gin.H{
			"list": []gin.H{
				{
					"id":          1,
					"elderlyName": "张奶奶",
					"bedNumber":   "3号楼201",
					"taskType":    "日常护理",
					"description": "协助洗漱、测量血压、协助用餐",
					"scheduledTime": "08:00",
					"status":      "pending",
					"priority":    "normal",
				},
				{
					"id":          2,
					"elderlyName": "李爷爷",
					"bedNumber":   "2号楼105",
					"taskType":    "健康监测",
					"description": "测量体温、血压、血糖",
					"scheduledTime": "09:00",
					"status":      "in_progress",
					"priority":    "important",
				},
			},
			"total":    15,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// CompleteTask 完成任务
func (h *CaregiverHandler) CompleteTask(c *gin.Context) {
	taskID := c.Param("id")

	var req struct {
		Notes     string   `json:"notes"`
		Images    []string `json:"images"`
		EvalScore int      `json:"evalScore"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	// TODO: 更新任务状态为已完成
	// 保存护理记录

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "任务已完成",
		"data": gin.H{
			"id":     taskID,
			"status": "completed",
		},
	})
}

// GetSchedule 获取排班信息
func (h *CaregiverHandler) GetSchedule(c *gin.Context) {
	weekStart := c.Query("weekStart")

	// TODO: 查询排班信息

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取成功",
		"data": gin.H{
			"weekStart": weekStart,
			"schedules": []gin.H{
				{
					"date": "2026-03-04",
					"shift": "day",
				},
				{
					"date": "2026-03-05",
					"shift": "morning",
				},
			},
		},
	})
}

// ClockIn 上班打卡
func (h *CaregiverHandler) ClockIn(c *gin.Context) {
	var req struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Address   string  `json:"address"`
		Photo     string  `json:"photo"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	// TODO: 保存打卡记录

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "打卡成功",
		"data": gin.H{
			"clockInTime": time.Now().Format("2006-01-02 15:04:05"),
			"address":     req.Address,
		},
	})
}

// ClockOut 下班打卡
func (h *CaregiverHandler) ClockOut(c *gin.Context) {
	var req struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Address   string  `json:"address"`
		Photo     string  `json:"photo"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	// TODO: 更新打卡记录

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "打卡成功",
		"data": gin.H{
			"clockOutTime": time.Now().Format("2006-01-02 15:04:05"),
			"address":      req.Address,
		},
	})
}

// GetAttendance 获取考勤记录
func (h *CaregiverHandler) GetAttendance(c *gin.Context) {
	month := c.Query("month")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	// TODO: 查询考勤记录

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取成功",
		"data": gin.H{
			"list": []gin.H{
				{
					"date":         "2026-03-01",
					"clockInTime":  "08:00",
					"clockOutTime": "17:30",
					"workHours":    9.5,
					"status":       "normal",
				},
			},
			"total":    22,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}
