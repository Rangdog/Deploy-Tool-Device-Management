package main

import (
	"BE_Manage_device/api"
	"BE_Manage_device/api/handler"
	"BE_Manage_device/cmd/server/docs"
	"BE_Manage_device/config"
	"BE_Manage_device/internal/database"
	"BE_Manage_device/internal/domain/service"
	"BE_Manage_device/pkg/utils"
	"log"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	config.LoadEnv()
	db := config.ConnectToDB()
	userRepository := database.NewPostgreSQLUserRepository(db)
	userSessionRepository := database.NewPostgreSQLUserSessionRepository(db)
	locationRepository := database.NewPostgreSQLLocationRepository(db)
	categoriesRepository := database.NewPostgreSQLCategoriesRepository(db)
	departmentRepository := database.NewPostgreSQLDepartmentsRepository(db)
	assetsRepository := database.NewPostgreSQLAssetsRepository(db)
	assetsLogRepository := database.NewPostgreSQLAssetsLogRepository(db)
	roleRepository := database.NewPostgreSQLRoleRepository(db)
	userRBACRepository := database.NewPostgreSQLUserRBACRepository(db)
	assignmentRepository := database.NewPostgreSQLAssignmentRepository(db)
	requestTransferRepository := database.NewPostgreSQLRequestTransferRepository(db)
	emailService := service.NewEmailService(config.SmtpPasswd)
	maintenanceRepository := database.NewPostgreSQLMaintenanceSchedulesRepository(db)
	notificationRepository := database.NewPostgreSQLNotificationRepository(db)
	//User
	userService := service.NewUserService(userRepository, emailService, userSessionRepository, roleRepository, assetsRepository, userRBACRepository)
	userHandler := handler.NewUserHandler(userService)
	//Location
	locationService := service.NewLocationService(locationRepository)
	locationHandler := handler.NewLocationHandler(locationService)
	//Categories
	categoriesService := service.NewCategoriesService(categoriesRepository)
	categoriesHandler := handler.NewCategoriesHandler(categoriesService)
	//Department
	departmentService := service.NewDepartmentsService(departmentRepository)
	departmentHandler := handler.NewDepartmentsHandler(departmentService)
	//SSE
	notificationsService := service.NewNotificationService(notificationRepository)
	SSeHandler := handler.NewSSEHandler(notificationsService)
	//Assets
	assetsService := service.NewAssetsService(assetsRepository, assetsLogRepository, roleRepository, userRBACRepository, userRepository, assignmentRepository, departmentRepository, notificationsService)
	assetsHandler := handler.NewAssetsHandler(assetsService)
	//Role
	roleService := service.NewRoleService(roleRepository)
	roleHandler := handler.NewRoleHandler(roleService)
	//Assignment
	assignmentService := service.NewAssignmentService(assignmentRepository, assetsLogRepository, assetsRepository, departmentRepository, userRepository, notificationsService)
	assignmentHandler := handler.NewAssignmentHandler(assignmentService)
	//AssetLog
	assetLogService := service.NewAssetLogService(assetsLogRepository, userRepository, roleRepository, assetsRepository)
	assetLogHandler := handler.NewAssetLogHandler(assetLogService)
	//Request Transfer
	requestTransferService := service.NewRequestTransferService(requestTransferRepository, assignmentService, userRepository, assetsRepository)
	requestTransferHandler := handler.NewRequestTransferHandler(requestTransferService)
	//Maintenance
	MaintenanceService := service.NewMaintenanceSchedulesService(maintenanceRepository, assetsRepository, userRepository, notificationsService)
	maintenanceHandler := handler.NewMaintenanceSchedulesHandler(MaintenanceService)

	// Notification
	notificationsHandler := handler.NewNotificationHandler(notificationsService)

	docs.SwaggerInfo.Title = "API Tool device manage"
	docs.SwaggerInfo.Description = "App Tool device manage"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = config.BASE_URL_BACKEND_FOR_SWAGGER
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	r := gin.Default()
	api.SetupRoutes(r, userHandler, locationHandler, categoriesHandler, departmentHandler, assetsHandler, roleHandler, assignmentHandler, assetLogHandler, requestTransferHandler, maintenanceHandler, SSeHandler, notificationsHandler, userSessionRepository, db)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	c := cron.New(cron.WithLocation(time.FixedZone("Asia/Ho_Chi_Minh", 7*3600)))

	_, err := c.AddFunc("0 8 * * *", func() {
		log.Println("🔔 Running maintenance notification check at 8:00 AM")
		utils.CheckAndSenMaintenanceNotification(db, emailService, assetsRepository, userRepository, notificationsService, assetsLogRepository)
	})
	if err != nil {
		log.Fatalf("❌ Failed to schedule cron job: %v", err)
	}
	_, err = c.AddFunc("0 8 * * *", func() {
		log.Println("🔔 Running warranty notification check at 8:00 AM")
		utils.SendEmailsForWarrantyExpiry(db, emailService, notificationsService, assetsRepository, userRepository)
	})
	if err != nil {
		log.Fatalf("❌ Failed to schedule cron job: %v", err)
	}

	_, err = c.AddFunc("0 9 * * *", func() {
		log.Println("🔔 Running update status when finish maintenance check at 8:00 AM")
		utils.UpdateStatusWhenFinishMaintenance(db, assetsRepository, userRepository, notificationsService, assetsLogRepository)
	})
	if err != nil {
		log.Fatalf("❌ Failed to schedule cron job: %v", err)
	}

	c.Start()

	if err := r.Run(config.Port); err != nil {
		log.Fatal("failed to run server:", err)
	}

}
