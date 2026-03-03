package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FinanceHandler 财务相关处理器
type FinanceHandler struct {
	DB *gorm.DB
}

// GetBills 获取账单列表
func (h *FinanceHandler) GetBills(c *gin.Context) {
	status := c.Query("status")
	billType := c.Query("billType")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	// TODO: 构建查询条件并查询账单

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取成功",
		"data": gin.H{
			"list": []gin.H{
				{
					"id":          1,
					"billNo":      "B202603001",
					"elderlyId":   1,
					"elderlyName": "张奶奶",
					"bedNumber":   "3号楼201",
					"billType":    "床位费",
					"amount":      3500,
					"period":      "2026年3月",
					"billDate":    "2026-03-01",
					"dueDate":     "2026-03-10",
					"status":      "unpaid",
					"details": []gin.H{
						{"name": "床位费", "amount": 2800},
						{"name": "护理费", "amount": 500},
						{"name": "伙食费", "amount": 200},
					},
				},
				{
					"id":          2,
					"billNo":      "B202603002",
					"elderlyId":   1,
					"elderlyName": "张奶奶",
					"bedNumber":   "3号楼201",
					"billType":    "护理费",
					"amount":      1800,
					"period":      "2026年3月",
					"billDate":    "2026-03-01",
					"dueDate":     "2026-03-10",
					"status":      "paid",
					"paymentDate": "2026-03-05",
					"paymentMethod": "微信支付",
				},
			},
			"total":    156,
			"page":     page,
			"pageSize": pageSize,
			"summary": gin.H{
				"totalAmount":    528600,
				"paidAmount":     485200,
				"unpaidAmount":   43400,
				"overdueAmount":  0,
			},
		},
	})
}

// CreateBill 创建账单
func (h *FinanceHandler) CreateBill(c *gin.Context) {
	var req struct {
		ElderlyID   uint                    `json:"elderlyId" binding:"required"`
		BillType    string                  `json:"billType" binding:"required"`
		Period      string                  `json:"period" binding:"required"`
		BillDate    string                  `json:"billDate" binding:"required"`
		DueDate     string                  `json:"dueDate" binding:"required"`
		Amount      float64                 `json:"amount" binding:"required,gt=0"`
		Details     []map[string]interface{} `json:"details" binding:"required"`
		Remarks     string                  `json:"remarks"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	// TODO: 创建账单

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id":        1,
			"billNo":    "B202603003",
			"createdAt": time.Now().Format("2006-01-02 15:04:05"),
		},
	})
}

// GetBillDetail 获取账单详情
func (h *FinanceHandler) GetBillDetail(c *gin.Context) {
	id := c.Param("id")

	// TODO: 查询账单详情

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取成功",
		"data": gin.H{
			"id":          id,
			"billNo":      "B202603001",
			"elderlyName": "张奶奶",
			"bedNumber":   "3号楼201",
			"billType":    "床位费",
			"amount":      3500,
			"period":      "2026年3月",
			"billDate":    "2026-03-01",
			"dueDate":     "2026-03-10",
			"status":      "unpaid",
			"details": []gin.H{
				{"name": "床位费（双人间）", "description": "3号楼201", "amount": 2800},
				{"name": "护理费", "description": "二级护理", "amount": 500},
				{"name": "伙食费", "description": "标准餐", "amount": 200},
			},
			"remarks": "",
		},
	})
}

// GetPayments 获取支付记录
func (h *FinanceHandler) GetPayments(c *gin.Context) {
	billID := c.Query("billId")
	paymentMethod := c.Query("paymentMethod")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	// TODO: 查询支付记录

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取成功",
		"data": gin.H{
			"list": []gin.H{
				{
					"id":            1,
					"paymentNo":     "P2026030001",
					"billNo":        "B202603002",
					"elderlyName":   "张奶奶",
					"amount":        1800,
					"paymentMethod": "微信支付",
					"paymentTime":   "2026-03-05 14:30:25",
					"transactionId": "wx1234567890",
					"operator":      "系统",
				},
			},
			"total":    142,
			"page":     page,
			"pageSize": pageSize,
			"summary": gin.H{
				"wechatAmount":  285600,
				"alipayAmount":  199600,
				"bankAmount":    0,
				"cashAmount":    0,
			},
		},
	})
}

// CreatePayment 创建支付
func (h *FinanceHandler) CreatePayment(c *gin.Context) {
	var req struct {
		BillID        uint    `json:"billId" binding:"required"`
		PaymentMethod string  `json:"paymentMethod" binding:"required,oneof=wechat alipay bank cash"`
		Amount        float64 `json:"amount" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	// TODO: 创建支付记录
	// 如果是微信/支付宝，返回支付参数

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"paymentNo":     "P2026030002",
			"paymentParams": gin.H{
				"timeStamp": strconv.FormatInt(time.Now().Unix(), 10),
				"nonceStr":   "random_string",
				"package":    "prepay_id=xxx",
				"signType":   "MD5",
				"paySign":    "sign",
			},
		},
	})
}

// Refund 退款
func (h *FinanceHandler) Refund(c *gin.Context) {
	var req struct {
		PaymentID uint    `json:"paymentId" binding:"required"`
		Amount    float64 `json:"amount" binding:"required,gt=0"`
		Reason    string  `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	// TODO: 处理退款
	// 调用微信/支付宝退款接口

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "退款申请已提交",
		"data": gin.H{
			"refundId":   "R2026030001",
			"refundNo":   "RF2026030001",
			"status":     "processing",
			"createdAt":  time.Now().Format("2006-01-02 15:04:05"),
		},
	})
}

// GetStatistics 获取财务统计
func (h *FinanceHandler) GetStatistics(c *gin.Context) {
	startDate := c.DefaultQuery("startDate", "")
	endDate := c.DefaultQuery("endDate", "")

	// TODO: 查询财务统计数据

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取成功",
		"data": gin.H{
			"summary": gin.H{
				"totalRevenue":   528600,
				"paidAmount":     485200,
				"unpaidAmount":    43400,
				"refundAmount":    0,
				"collectionRate":  91.8,
			},
			"byType": []gin.H{
				{"type": "床位费", "amount": 280000, "percent": 52.98},
				{"type": "护理费", "amount": 150000, "percent": 28.38},
				{"type": "伙食费", "amount": 80000, "percent": 15.13},
				{"type": "其他", "amount": 18600, "percent": 3.51},
			},
			"byMonth": []gin.H{
				{"month": "2026-01", "amount": 175000},
				{"month": "2026-02", "amount": 175000},
				{"month": "2026-03", "amount": 178600},
			},
			"overdue": gin.H{
				"count":  0,
				"amount": 0,
			},
		},
	})
}

// Export 导出财务报表
func (h *FinanceHandler) Export(c *gin.Context) {
	reportType := c.Query("type")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	// TODO: 生成Excel文件
	// 使用excelize库生成报表

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "导出成功",
		"data": gin.H{
			"downloadUrl": "/downloads/finance_report_" + strconv.FormatInt(time.Now().Unix(), 10) + ".xlsx",
		},
	})
}
