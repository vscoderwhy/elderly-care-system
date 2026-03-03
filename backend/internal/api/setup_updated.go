package api

import (
	 elderlyCareSystem "elderly-care-system/internal/api/handler"
	 "elderly-care-system/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes 配置所有路由
func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	// 初始化各个处理器
	commonHandler := &elderlyCareSystem.CommonHandler{DB: db}
	userHandler := &elderlyCareSystem.UserHandler{DB: db}
	elderlyHandler := &elderlyCareSystem.ElderlyHandler{DB: db}
	careHandler := &elderlyCareSystem.CareHandler{DB: db}
	healthHandler := &elderlyCareSystem.HealthHandler{DB: db}
	financeHandler := &elderlyCareSystem.FinanceHandler{DB: db}
	visitHandler := &elderlyCareSystem.VisitHandler{DB: db}
	notificationHandler := &elderlyCareSystem.NotificationHandler{DB: db}
	staffHandler := &elderlyCareSystem.StaffHandler{DB: db}
	caregiverHandler := &elderlyCareSystem.CaregiverHandler{DB: db}
	medicationHandler := &elderlyCareSystem.MedicationHandler{DB: db}
	mealHandler := &elderlyCareSystem.MealHandler{DB: db}
	equipmentHandler := &elderlyCareSystem.EquipmentHandler{DB: db}

	// CORS 中间件
	router.Use(middleware.CORS())

	// 健康检查
	router.GET("/health", commonHandler.HealthCheck)

	// 静态文件服务
	router.Static("/uploads", "./uploads")
	router.Static("/downloads", "./downloads")

	// API v1
	v1 := router.Group("/api/v1")
	{
		// ========== 公共路由 ==========
		public := v1.Group("")
		{
			// 家属登录
			public.POST("/auth/family/login", userHandler.FamilyLogin)
			public.POST("/auth/family/wechat-login", userHandler.FamilyWechatLogin)
			public.POST("/auth/family/sms-code", userHandler.SendSMSCode)

			// 护工登录
			public.POST("/auth/caregiver/login", caregiverHandler.Login)
			public.POST("/auth/caregiver/wechat-login", caregiverHandler.WechatLogin)

			// WebSocket
			public.GET("/ws", commonHandler.WebSocket)

			// 文件上传
			public.POST("/upload", commonHandler.UploadFile)

			// 数据导出
			public.POST("/export", commonHandler.ExportData)
		}

		// ========== 需要认证的路由 ==========
		protected := v1.Group("")
		protected.Use(middleware.Auth())
		{
			// ========== 用户相关 ==========
			user := protected.Group("/user")
			{
				user.GET("/info", userHandler.GetInfo)
				user.PUT("/info", userHandler.UpdateInfo)
				user.PUT("/password", userHandler.ChangePassword)
				user.POST("/logout", userHandler.Logout)
			}

			// ========== 老人管理 ==========
			elderly := protected.Group("/elderly")
			{
				elderly.GET("/list", elderlyHandler.List)
				elderly.POST("/create", elderlyHandler.Create)
				elderly.GET("/detail/:id", elderlyHandler.Detail)
				elderly.PUT("/update/:id", elderlyHandler.Update)
				elderly.DELETE("/delete/:id", elderlyHandler.Delete)
				elderly.GET("/:id/health", elderlyHandler.GetHealthRecords)
				elderly.GET("/:id/care", elderlyHandler.GetCareRecords)
				elderly.GET("/:id/bills", elderlyHandler.GetBills)
				elderly.GET("/:id/family", elderlyHandler.GetFamily)
				elderly.POST("/:id/family", elderlyHandler.AddFamily)
				elderly.GET("/:id/visits", elderlyHandler.GetVisits)
			}

			// ========== 护理管理 ==========
			care := protected.Group("/care")
			{
				care.GET("/tasks", careHandler.GetTasks)
				care.POST("/tasks", careHandler.CreateTask)
				care.GET("/tasks/:id", careHandler.GetTaskDetail)
				care.PUT("/tasks/:id/status", careHandler.UpdateTaskStatus)
				care.GET("/records", careHandler.GetRecords)
				care.POST("/records", careHandler.CreateRecord)
				care.GET("/records/:id", careHandler.GetRecordDetail)
				care.GET("/assessments", careHandler.GetAssessments)
				care.POST("/assessments", careHandler.CreateAssessment)
				care.GET("/statistics", careHandler.GetStatistics)
			}

			// ========== 健康数据 ==========
			health := protected.Group("/health")
			{
				health.GET("/records", healthHandler.GetRecords)
				health.POST("/records", healthHandler.CreateRecord)
				health.GET("/records/:id", healthHandler.GetRecordDetail)
				health.PUT("/records/:id", healthHandler.UpdateRecord)
				health.GET("/statistics", healthHandler.GetStatistics)
				health.GET("/abnormal", healthHandler.GetAbnormalData)
			}

			// ========== 财务管理 ==========
			finance := protected.Group("/finance")
			{
				finance.GET("/bills", financeHandler.GetBills)
				finance.POST("/bills", financeHandler.CreateBill)
				finance.GET("/bills/:id", financeHandler.GetBillDetail)
				finance.GET("/payments", financeHandler.GetPayments)
				finance.POST("/payments", financeHandler.CreatePayment)
				finance.POST("/refund", financeHandler.Refund)
				finance.GET("/statistics", financeHandler.GetStatistics)
				finance.GET("/export", financeHandler.Export)
			}

			// ========== 探视预约 ==========
			visits := protected.Group("/visits")
			{
				visits.GET("/appointments", visitHandler.GetAppointments)
				visits.POST("/appointments", visitHandler.CreateAppointment)
				visits.GET("/appointments/:id", visitHandler.GetAppointmentDetail)
				visits.PUT("/appointments/:id", visitHandler.UpdateAppointment)
				visits.POST("/appointments/:id/approve", visitHandler.ApproveAppointment)
				visits.POST("/appointments/:id/reject", visitHandler.RejectAppointment)
				visits.POST("/appointments/:id/cancel", visitHandler.CancelAppointment)
				visits.POST("/appointments/:id/complete", visitHandler.CompleteAppointment)
			}

			// ========== 消息通知 ==========
			notifications := protected.Group("/notifications")
			{
				notifications.GET("/list", notificationHandler.GetNotifications)
				notifications.GET("/unread", notificationHandler.GetUnreadCount)
				notifications.PUT("/:id/read", notificationHandler.MarkAsRead)
				notifications.PUT("/read-all", notificationHandler.MarkAllAsRead)
			}

			// ========== 员工管理（管理员） ==========
			staffGroup := protected.Group("/staff")
			staffGroup.Use(middleware.RoleAuth("admin"))
			{
				staffGroup.GET("/list", staffHandler.GetList)
				staffGroup.POST("/create", staffHandler.Create)
				staffGroup.GET("/detail/:id", staffHandler.GetDetail)
				staffGroup.PUT("/update/:id", staffHandler.Update)
				staffGroup.DELETE("/delete/:id", staffHandler.Delete)
				staffGroup.GET("/schedule", staffHandler.GetSchedule)
				staffGroup.POST("/schedule", staffHandler.CreateSchedule)
				staffGroup.GET("/attendance", staffHandler.GetAttendance)
				staffGroup.GET("/performance", staffHandler.GetPerformance)
			}

			// ========== 护工专属 ==========
			caregiver := protected.Group("/caregiver")
			caregiver.Use(middleware.RoleAuth("caregiver"))
			{
				caregiver.GET("/info", caregiverHandler.GetInfo)
				caregiver.GET("/tasks", caregiverHandler.GetTasks)
				caregiver.POST("/tasks/:id/complete", caregiverHandler.CompleteTask)
				caregiver.GET("/schedule", caregiverHandler.GetSchedule)
				caregiver.POST("/clock-in", caregiverHandler.ClockIn)
				caregiver.POST("/clock-out", caregiverHandler.ClockOut)
				caregiver.GET("/attendance", caregiverHandler.GetAttendance)
			}

			// ========== 用药管理 ==========
			medication := protected.Group("/medication")
			{
				medication.GET("/prescriptions", medicationHandler.GetPrescriptions)
				medication.POST("/prescriptions", medicationHandler.CreatePrescription)
				medication.GET("/prescriptions/:id", medicationHandler.GetPrescriptionDetail)
				medication.PUT("/prescriptions/:id", medicationHandler.UpdatePrescription)
				medication.GET("/today", medicationHandler.GetTodayMedications)
				medication.POST("/confirm", medicationHandler.ConfirmMedication)
				medication.POST("/skip", medicationHandler.SkipMedication)
			}

			// ========== 膳食管理 ==========
			meal := protected.Group("/meal")
			{
				meal.GET("/list", mealHandler.GetMeals)
				meal.POST("/create", mealHandler.CreateMeal)
				meal.GET("/detail/:id", mealHandler.GetMealDetail)
				meal.PUT("/update/:id", mealHandler.UpdateMeal)
				meal.DELETE("/delete/:id", mealHandler.DeleteMeal)
				meal.GET("/today", mealHandler.GetTodayMeals)
				meal.GET("/statistics", mealHandler.GetStatistics)
			}

			// ========== 设备管理 ==========
			equipment := protected.Group("/equipment")
			{
				equipment.GET("/list", equipmentHandler.GetEquipment)
				equipment.POST("/create", equipmentHandler.CreateEquipment)
				equipment.GET("/detail/:id", equipmentHandler.GetEquipmentDetail)
				equipment.PUT("/update/:id", equipmentHandler.UpdateEquipment)
				equipment.DELETE("/delete/:id", equipmentHandler.DeleteEquipment)
				equipment.GET("/categories", equipmentHandler.GetCategories)
				equipment.POST("/:id/maintenance", equipmentHandler.AddMaintenance)
				equipment.GET("/:id/maintenance", equipmentHandler.GetMaintenanceRecords)
			}
		}

		// ========== 管理员专属路由 ==========
		admin := v1.Group("/admin")
		admin.Use(middleware.Auth(), middleware.RoleAuth("admin"))
		{
			admin.GET("/dashboard", commonHandler.GetDashboard)
			admin.GET("/statistics", commonHandler.GetStatistics)
			admin.GET("/system/config", commonHandler.GetSystemConfig)
			admin.PUT("/system/config", commonHandler.UpdateSystemConfig)
			admin.GET("/system/logs", commonHandler.GetSystemLogs)
		}
	}
}
