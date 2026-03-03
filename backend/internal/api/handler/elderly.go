package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ElderlyHandler 老人相关处理器
type ElderlyHandler struct {
	DB *gorm.DB
}

// ListRequest 老人列表请求
type ListRequest struct {
	Name       string `form:"name"`
	CareLevel  string `form:"careLevel"`
	BedNumber  string `form:"bedNumber"`
	Status     string `form:"status"`
	Department string `form:"department"`
	Page       int    `form:"page" binding:"min=1"`
	PageSize   int    `form:"pageSize" binding:"min=1,max=100"`
}

// CreateRequest 创建老人请求
type CreateRequest struct {
	Name        string   `json:"name" binding:"required"`
	Gender      string   `json:"gender" binding:"required,oneof=男 女"`
	IDCard      string   `json:"idCard" binding:"required"`
	Birthday    string   `json:"birthday" binding:"required"`
	Phone       string   `json:"phone"`
	Address     string   `json:"address"`
	CareLevel   string   `json:"careLevel" binding:"required,oneof=一级 二级 三级 特级"`
	BedNumber   string   `json:"bedNumber" binding:"required"`
	AdmitDate   string   `json:"admitDate" binding:"required"`
	EmergencyContact string `json:"emergencyContact" binding:"required"`
	EmergencyPhone   string `json:"emergencyPhone" binding:"required"`
	Photo       string   `json:"photo"`
	Remarks     string   `json:"remarks"`
}

// List 获取老人列表
func (h *ElderlyHandler) List(c *gin.Context) {
	var req ListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	// TODO: 构建查询条件
	// query := h.DB.Model(&models.Elderly{})
	// if req.Name != "" {
	//     query = query.Where("name LIKE ?", "%"+req.Name+"%")
	// }
	// if req.CareLevel != "" {
	//     query = query.Where("care_level = ?", req.CareLevel)
	// }
	// ... 其他筛选条件

	// TODO: 分页查询
	// var total int64
	// query.Count(&total)
	// var elderlyList []models.Elderly
	// query.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&elderlyList)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取成功",
		"data": gin.H{
			"list": []gin.H{
				{
					"id":          1,
					"name":        "张奶奶",
					"gender":      "女",
					"age":         78,
					"careLevel":   "二级",
					"bedNumber":   "3号楼201",
					"admitDate":   "2023-03-15",
					"status":      "在院",
					"photo":       "",
					"phone":       "138****8888",
					"emergencyContact": "张先生",
					"emergencyPhone":   "13900139000",
				},
				{
					"id":          2,
					"name":        "李爷爷",
					"gender":      "男",
					"age":         82,
					"careLevel":   "一级",
					"bedNumber":   "2号楼105",
					"admitDate":   "2023-05-20",
					"status":      "在院",
					"photo":       "",
					"phone":       "137****7777",
					"emergencyContact": "李女士",
					"emergencyPhone":   "13800138000",
				},
			},
			"total":    248,
			"page":     req.Page,
			"pageSize": req.PageSize,
		},
	})
}

// Create 创建老人档案
func (h *ElderlyHandler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	// TODO: 验证身份证号
	// TODO: 验证床位号是否可用
	// TODO: 创建老人档案
	// elderly := models.Elderly{
	//     Name: req.Name,
	//     Gender: req.Gender,
	//     ...
	// }
	// h.DB.Create(&elderly)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id":        1,
			"createdAt": time.Now().Format("2006-01-02 15:04:05"),
		},
	})
}

// Detail 获取老人详情
func (h *ElderlyHandler) Detail(c *gin.Context) {
	id := c.Param("id")

	// TODO: 查询老人详情
	// var elderly models.Elderly
	// h.DB.First(&elderly, id)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取成功",
		"data": gin.H{
			"id":          id,
			"name":        "张奶奶",
			"gender":      "女",
			"age":         78,
			"idCard":      "3201**********1234",
			"birthday":    "1946-05-15",
			"careLevel":   "二级",
			"bedNumber":   "3号楼201",
			"admitDate":   "2023-03-15",
			"status":      "在院",
			"photo":       "",
			"phone":       "138****8888",
			"address":     "江苏省南京市玄武区",
			"emergencyContact": "张先生",
			"emergencyPhone":   "13900139000",
			"remarks":     "老人有高血压，需定期监测",
			"healthScore": 85,
			"stayDays":    720,
		},
	})
}

// Update 更新老人信息
func (h *ElderlyHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	// TODO: 更新老人信息
	// var elderly models.Elderly
	// h.DB.First(&elderly, id)
	// h.DB.Model(&elderly).Updates(req)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
		"data": gin.H{
			"id":         id,
			"updatedAt": time.Now().Format("2006-01-02 15:04:05"),
		},
	})
}

// Delete 删除老人档案
func (h *ElderlyHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	// TODO: 软删除老人档案
	// h.DB.Delete(&models.Elderly{}, id)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
	})
}

// GetHealthRecords 获取健康记录
func (h *ElderlyHandler) GetHealthRecords(c *gin.Context) {
	id := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	// TODO: 查询健康记录
	// var records []models.HealthData
	// h.DB.Where("elderly_id = ?", id).Find(&records)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取成功",
		"data": gin.H{
			"list": []gin.H{
				{
					"id":         1,
					"recordDate": "2026-03-03",
					"recordTime": "08:30",
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

// GetCareRecords 获取护理记录
func (h *ElderlyHandler) GetCareRecords(c *gin.Context) {
	id := c.Param("id")
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
					"careType":   "日常护理",
					"content":    "协助洗漱、测量血压、协助用餐",
					"careTime":   "2026-03-03 08:00",
					"nurseName":  "赵护士",
					"images":     []string{"/static/care1.jpg"},
					"evaluation": 5,
				},
			},
			"total":    245,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// GetBills 获取费用账单
func (h *ElderlyHandler) GetBills(c *gin.Context) {
	id := c.Param("id")
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	// TODO: 查询费用账单

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取成功",
		"data": gin.H{
			"list": []gin.H{
				{
					"id":          1,
					"billNo":      "B202603001",
					"billType":    "床位费",
					"amount":      3500,
					"period":      "2026年3月",
					"billDate":    "2026-03-01",
					"dueDate":     "2026-03-10",
					"status":      "unpaid",
				},
			},
			"total":    48,
			"page":     page,
			"pageSize": pageSize,
			"summary": gin.H{
				"unpaidTotal":  3200,
				"paidTotal":    8500,
				"overdueTotal": 0,
			},
		},
	})
}

// GetFamily 获取家属信息
func (h *ElderlyHandler) GetFamily(c *gin.Context) {
	id := c.Param("id")

	// TODO: 查询家属信息

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取成功",
		"data": gin.H{
			"list": []gin.H{
				{
					"id":         1,
					"name":       "张先生",
					"relation":   "子女",
					"phone":      "13900139000",
					"wechat":     "wx123456",
					"isPrimary":  true,
				},
				{
					"id":         2,
					"name":       "李女士",
					"relation":   "子女",
					"phone":      "13800138000",
					"isPrimary":  false,
				},
			},
		},
	})
}

// AddFamily 添加家属
func (h *ElderlyHandler) AddFamily(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name      string `json:"name" binding:"required"`
		Relation  string `json:"relation" binding:"required"`
		Phone     string `json:"phone" binding:"required"`
		Wechat    string `json:"wechat"`
		IsPrimary bool   `json:"isPrimary"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	// TODO: 添加家属信息

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "添加成功",
		"data": gin.H{
			"id":        1,
			"createdAt": time.Now().Format("2006-01-02 15:04:05"),
		},
	})
}
