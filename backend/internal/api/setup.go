package api

import (
	"github.com/gin-gonic/gin"
	"github.com/elderly-care/internal/api/handler"
	"github.com/elderly-care/internal/middleware"
	"gorm.io/gorm"
)

func Setup(r *gin.Engine, db *gorm.DB, m *middleware.Middleware) {
	// 初始化handler
	h := handler.NewHandler(db)

	// API v1
	v1 := r.Group("/api/v1")
	{
		// 公开路由
		public := v1.Group("")
		{
			// 认证
			auth := public.Group("/auth")
			{
				auth.POST("/login", h.Login)
				auth.POST("/logout", h.Logout)
				auth.POST("/refresh", h.RefreshToken)
			}

			// WebSocket
			v1.GET("/ws", h.WebSocket)
		}

		// 需要认证的路由
		auth := v1.Group("")
		auth.Use(m.Auth())
		{
			// 用户
			users := auth.Group("/users")
			{
				users.GET("/me", h.GetCurrentUser)
				users.PUT("/me", h.UpdateCurrentUser)
				users.POST("/avatar", h.UploadAvatar)
			}

			// 管理员路由
			admin := auth.Group("")
			admin.Use(m.Role("admin"))
			{
				// 用户管理
				adminUsers := admin.Group("/users")
				{
					adminUsers.GET("", h.GetUserList)
					adminUsers.POST("", h.CreateUser)
					adminUsers.GET("/:id", h.GetUser)
					adminUsers.PUT("/:id", h.UpdateUser)
					adminUsers.DELETE("/:id", h.DeleteUser)
				}
			}

			// 老人管理
			elderly := auth.Group("/elderly")
			{
				elderly.GET("", h.GetElderlyList)
				elderly.POST("", h.CreateElderly)
				elderly.GET("/:id", h.GetElderly)
				elderly.PUT("/:id", h.UpdateElderly)
				elderly.DELETE("/:id", h.DeleteElderly)
				elderly.GET("/:id/family", h.GetElderlyFamily)
				elderly.POST("/:id/family", h.AddElderlyFamily)
				elderly.GET("/:id/health", h.GetElderlyHealth)
			}

			// 护理管理
			care := auth.Group("/care")
			{
				// 护理任务
				tasks := care.Group("/tasks")
				{
					tasks.GET("", h.GetCareTasks)
					tasks.POST("", h.CreateCareTask)
					tasks.GET("/:id", h.GetCareTask)
					tasks.PUT("/:id", h.UpdateCareTask)
					tasks.DELETE("/:id", h.DeleteCareTask)
					tasks.POST("/:id/start", h.StartCareTask)
					tasks.POST("/:id/complete", h.CompleteCareTask)
				}

				// 护理记录
				records := care.Group("/records")
				{
					records.GET("", h.GetCareRecords)
					records.POST("", h.CreateCareRecord)
					records.GET("/:id", h.GetCareRecord)
					records.PUT("/:id", h.UpdateCareRecord)
					records.DELETE("/:id", h.DeleteCareRecord)
				}
			}

			// 健康数据
			health := auth.Group("/health")
			{
				health.GET("/latest/:elderlyId", h.GetLatestHealth)
				health.GET("/trend", h.GetHealthTrend)
				health.POST("", h.CreateHealthData)
				health.GET("/:id", h.GetHealthData)
				health.PUT("/:id", h.UpdateHealthData)
			}

			// 财务管理
			finance := auth.Group("/finance")
			{
				// 账单
				bills := finance.Group("/bills")
				{
					bills.GET("", h.GetBills)
					bills.POST("", h.CreateBill)
					bills.GET("/:id", h.GetBill)
					bills.PUT("/:id", h.UpdateBill)
					bills.DELETE("/:id", h.DeleteBill)
					bills.POST("/:id/pay", h.PayBill)
					bills.GET("/:id/invoice", h.GetInvoice)
				}

				// 支付记录
				payments := finance.Group("/payments")
				{
					payments.GET("", h.GetPayments)
					payments.POST("/refund", h.RefundPayment)
				}
			}

			// 探视预约
			visits := auth.Group("/visits")
			{
				visits.GET("/appointments", h.GetVisitAppointments)
				visits.POST("/appointments", h.CreateAppointment)
				visits.PUT("/appointments/:id", h.UpdateAppointment)
				visits.DELETE("/appointments/:id", h.CancelAppointment)
				visits.GET("/slots", h.GetAvailableSlots)
			}

			// 消息通知
			notifications := auth.Group("/notifications")
			{
				notifications.GET("", h.GetNotifications)
				notifications.PUT("/:id/read", h.MarkAsRead)
				notifications.POST("/read-all", h.MarkAllAsRead)
				notifications.GET("/unread-count", h.GetUnreadCount)
			}

			// 文件上传
			upload := auth.Group("/upload")
			{
				upload.POST("/image", h.UploadImage)
				upload.POST("/file", h.UploadFile)
			}

			// 数据导出
			export := auth.Group("/export")
			{
				export.POST("/elderly", h.ExportElderlyList)
				export.POST("/care-records", h.ExportCareRecords)
				export.POST("/bills", h.ExportFinancialReport)
			}

			// 统计数据
			stats := auth.Group("/stats")
			{
				stats.GET("/home", h.GetHomeStats)
				stats.GET("/care", h.GetCareStats)
				stats.GET("/health/:elderlyId", h.GetHealthStats)
			}
		}
	}
}
