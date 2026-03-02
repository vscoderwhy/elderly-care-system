package main

import (
	"log"
	"elderly-care-system/internal/config"
	"elderly-care-system/internal/handler"
	"elderly-care-system/internal/middleware"
	"elderly-care-system/internal/repository"
	"elderly-care-system/internal/service"
	"elderly-care-system/pkg/database"
	"elderly-care-system/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	cfg := config.Load()

	// Init logger
	logger.Init(cfg.LogLevel)

	// Init database
	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	// Init redis
	redisClient := database.ConnectRedis(cfg.Redis)

	// Init repositories
	userRepo := repository.NewUserRepository(db)
	elderlyRepo := repository.NewElderlyRepository(db)
	careRepo := repository.NewCareRepository(db)
	billRepo := repository.NewBillRepository(db)
	roomRepo := repository.NewRoomRepository(db)
	scheduleRepo := repository.NewScheduleRepository(db)
	medicationRepo := repository.NewMedicationRepository(db)
	visitRepo := repository.NewVisitRepository(db)
	alertRepo := repository.NewAlertRepository(db)
	permissionRepo := repository.NewPermissionRepository(db)
	menuRepo := repository.NewMenuRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	// Init services
	authService := service.NewAuthService(userRepo, redisClient, cfg.JWT.Secret)
	userService := service.NewUserService(userRepo)
	elderlyService := service.NewElderlyService(elderlyRepo)
	careService := service.NewCareService(careRepo)
	billService := service.NewBillService(billRepo)
	roomService := service.NewRoomService(roomRepo)
	statsService := service.NewStatsService(elderlyRepo, careRepo, billRepo, userRepo, roomRepo)
	scheduleService := service.NewScheduleService(scheduleRepo, userRepo)
	medicationService := service.NewMedicationService(medicationRepo, elderlyRepo)
	visitService := service.NewVisitService(visitRepo, elderlyRepo)
	alertService := service.NewAlertService(alertRepo, elderlyRepo, medicationRepo, billRepo, careRepo, roomRepo)
	rbacService := service.NewRBACService(permissionRepo, menuRepo, roleRepo, userRepo)

	// Init handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	elderlyHandler := handler.NewElderlyHandler(elderlyService)
	careHandler := handler.NewCareHandler(careService)
	billHandler := handler.NewBillHandler(billService)
	roomHandler := handler.NewRoomHandler(roomService)
	statsHandler := handler.NewStatsHandler(statsService)
	scheduleHandler := handler.NewScheduleHandler(scheduleService)
	exportHandler := handler.NewExportHandler(elderlyRepo, careRepo, billRepo)
	medicationHandler := handler.NewMedicationHandler(medicationService)
	visitHandler := handler.NewVisitHandler(visitService)
	alertHandler := handler.NewAlertHandler(alertService)
	rbacHandler := handler.NewRBACHandler(rbacService)

	// Setup router
	router := gin.Default()

	// Middleware
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())

	// Public routes
	public := router.Group("/api/auth")
	{
		public.POST("/register", authHandler.Register)
		public.POST("/login", authHandler.Login)
		public.POST("/wechat-login", authHandler.WeChatLogin)
	}

	// Protected routes
	protected := router.Group("/api")
	protected.Use(middleware.Auth(cfg.JWT.Secret))
	{
		// Stats routes (统计)
		protected.GET("/stats/dashboard", statsHandler.GetDashboardStats)
		protected.GET("/stats/occupancy", statsHandler.GetBedOccupancy)
		protected.GET("/stats/care", statsHandler.GetCareStats)
		protected.GET("/stats/finance", statsHandler.GetFinanceStats)

		// User routes
		protected.GET("/user/profile", userHandler.GetProfile)
		protected.PUT("/user/profile", userHandler.UpdateProfile)
		protected.GET("/user/elderly-list", userHandler.GetElderlyList)

		// Staff routes (员工管理)
		protected.GET("/staff", userHandler.ListStaff)
		protected.POST("/staff", userHandler.CreateStaff)
		protected.PUT("/staff/:id", userHandler.UpdateStaff)
		protected.DELETE("/staff/:id", userHandler.DeleteStaff)

		// Elderly routes
		protected.GET("/elderly", elderlyHandler.List)
		protected.GET("/elderly/:id", elderlyHandler.Get)
		protected.POST("/elderly", elderlyHandler.Create)
		protected.PUT("/elderly/:id", elderlyHandler.Update)

		// Care routes
		protected.GET("/care/records", careHandler.ListRecords)
		protected.POST("/care/records", careHandler.CreateRecord)
		protected.GET("/care/my-tasks", careHandler.GetMyTasks)
		protected.GET("/care/items", careHandler.ListCareItems)

		// Health record routes (健康记录)
		protected.GET("/health/records", careHandler.ListHealthRecords)
		protected.GET("/health/records/latest/:elderly_id", careHandler.GetLatestHealthRecords)
		protected.POST("/health/records", careHandler.CreateHealthRecord)
		protected.DELETE("/health/records/:id", careHandler.DeleteHealthRecord)

		// Service request routes (服务呼叫)
		protected.GET("/service/requests", careHandler.ListServiceRequests)
		protected.POST("/service/requests", careHandler.CreateServiceRequest)
		protected.PUT("/service/requests/:id", careHandler.HandleServiceRequest)

		// Bill routes
		protected.GET("/bills", billHandler.List)
		protected.GET("/bills/:id", billHandler.Get)
		protected.POST("/bills/:id/pay", billHandler.Pay)

		// Room routes
		protected.GET("/rooms/buildings", roomHandler.ListBuildings)
		protected.GET("/rooms/buildings/:id", roomHandler.GetBuilding)
		protected.GET("/rooms", roomHandler.ListRooms)
		protected.GET("/rooms/stats", roomHandler.GetBedStats)
		protected.POST("/rooms/beds/:id/assign", roomHandler.AssignBed)
		protected.POST("/rooms/beds/:id/release", roomHandler.ReleaseBed)

		// Schedule routes (排班管理)
		protected.GET("/schedules", scheduleHandler.GetScheduleList)
		protected.POST("/schedules", scheduleHandler.CreateSchedule)
		protected.GET("/schedules/staff/:staff_id", scheduleHandler.GetStaffSchedule)
		protected.GET("/schedules/my", scheduleHandler.GetMySchedule)
		protected.PUT("/schedules/:id/status", scheduleHandler.UpdateScheduleStatus)
		protected.DELETE("/schedules/:id", scheduleHandler.DeleteSchedule)
		protected.GET("/schedules/stats/monthly", scheduleHandler.GetMonthlyStats)

		// Export routes (数据导出)
		protected.GET("/export/elderly", exportHandler.ExportElderlyList)
		protected.GET("/export/care-records", exportHandler.ExportCareRecords)
		protected.GET("/export/health-data", exportHandler.ExportHealthData)
		protected.GET("/export/finance", exportHandler.ExportFinance)

		// Medication routes (用药管理)
		protected.GET("/medications", medicationHandler.ListMedications)
		protected.POST("/medications", medicationHandler.CreateMedication)
		protected.PUT("/medications/:id", medicationHandler.UpdateMedication)
		protected.DELETE("/medications/:id", medicationHandler.DeleteMedication)
		protected.GET("/medications/alerts", medicationHandler.GetMedicationAlerts)
		protected.GET("/elderly/:id/medications", medicationHandler.ListElderlyMedications)
		protected.POST("/elderly/:id/medications", medicationHandler.CreateMedicationRecord)
		protected.GET("/elderly/:id/medications/today", medicationHandler.GetTodayMedications)
		protected.PUT("/medication-logs/:log_id/complete", medicationHandler.CompleteMedicationLog)
		protected.GET("/export/medication-records", medicationHandler.ExportMedicationRecords)

		// Visit routes (探视预约)
		protected.GET("/visits", visitHandler.ListVisits)
		protected.GET("/visits/:id", visitHandler.GetVisit)
		protected.POST("/visits", visitHandler.CreateVisit)
		protected.PUT("/visits/:id", visitHandler.UpdateVisit)
		protected.DELETE("/visits/:id", visitHandler.DeleteVisit)
		protected.PUT("/visits/:id/confirm", visitHandler.ConfirmVisit)
		protected.PUT("/visits/:id/cancel", visitHandler.CancelVisit)
		protected.PUT("/visits/:id/complete", visitHandler.CompleteVisit)
		protected.GET("/visits/today", visitHandler.GetTodayVisits)
		protected.GET("/visits/upcoming", visitHandler.GetUpcomingVisits)
		protected.GET("/visits/date-range", visitHandler.GetVisitsByDateRange)
		protected.GET("/elderly/:id/visits", visitHandler.ListElderlyVisits)

		// Alert routes (智能预警)
		protected.GET("/alerts", alertHandler.ListAlerts)
		protected.GET("/alerts/summary", alertHandler.GetAlertSummary)
		protected.GET("/alerts/active", alertHandler.GetActiveAlerts)
		protected.GET("/alerts/:id", alertHandler.GetAlert)
		protected.POST("/alerts", alertHandler.CreateAlert)
		protected.PUT("/alerts/:id/acknowledge", alertHandler.AcknowledgeAlert)
		protected.PUT("/alerts/:id/resolve", alertHandler.ResolveAlert)
		protected.POST("/alerts/check", alertHandler.CheckAlerts)
		protected.GET("/alert-rules", alertHandler.ListRules)
		protected.POST("/alert-rules", alertHandler.CreateRule)
		protected.PUT("/alert-rules/:id", alertHandler.UpdateRule)
		protected.DELETE("/alert-rules/:id", alertHandler.DeleteRule)

		// RBAC routes (权限管理)
		protected.POST("/system/init", rbacHandler.InitializeSystem)
		protected.GET("/system/menus", rbacHandler.ListMenus)
		protected.POST("/system/menus", rbacHandler.CreateMenu)
		protected.PUT("/system/menus/:id", rbacHandler.UpdateMenu)
		protected.DELETE("/system/menus/:id", rbacHandler.DeleteMenu)
		protected.GET("/system/permissions", rbacHandler.ListPermissions)
		protected.POST("/system/permissions", rbacHandler.CreatePermission)
		protected.PUT("/system/permissions/:id", rbacHandler.UpdatePermission)
		protected.DELETE("/system/permissions/:id", rbacHandler.DeletePermission)
		protected.GET("/system/roles", rbacHandler.ListRoles)
		protected.GET("/system/roles/:id", rbacHandler.GetRole)
		protected.PUT("/system/roles/:id/permissions", rbacHandler.AssignPermissionsToRole)
		protected.PUT("/system/roles/:id/menus", rbacHandler.AssignMenusToRole)
		protected.GET("/system/users", rbacHandler.ListUsers)
		protected.GET("/system/users/:id/roles", rbacHandler.GetUserRoles)
		protected.PUT("/system/users/:id/roles", rbacHandler.UpdateUserRoles)
		protected.POST("/system/users/:id/roles", rbacHandler.AssignRoleToUser)
		protected.DELETE("/system/users/:id/roles", rbacHandler.RemoveRoleFromUser)
		protected.GET("/user/menus", rbacHandler.GetUserMenus)
		protected.GET("/user/permissions", rbacHandler.GetUserPermissions)
	}

	// Start server
	addr := ":" + cfg.Port
	log.Printf("Server starting on %s", addr)
	router.Run(addr)
}
