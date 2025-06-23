package main

import (
	"BE_Manage_device/api/handler"
	api "BE_Manage_device/api/router"
	"BE_Manage_device/cmd/server/docs"
	"BE_Manage_device/config"
	"BE_Manage_device/internal/repository"
	"BE_Manage_device/internal/service"
	cronjob "BE_Manage_device/pkg/cron_job"
	"log"

	"github.com/gin-contrib/pprof"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	config.LoadEnv()
	db := config.ConnectToDB()
	config.InitRedis()
	repos := repository.NewRepository(db)
	services := service.NewServices(repos, config.SmtpPasswd)
	//User
	userHandler := handler.NewUserHandler(services.User)
	//Location
	locationHandler := handler.NewLocationHandler(services.Location)
	//Categories
	categoriesHandler := handler.NewCategoriesHandler(services.Categories)
	//Department
	departmentHandler := handler.NewDepartmentsHandler(services.Department)
	//SSE
	SSeHandler := handler.NewSSEHandler(services.Notification)
	//Assets
	assetsHandler := handler.NewAssetsHandler(services.Assets)
	//Role
	roleHandler := handler.NewRoleHandler(services.Role)
	//Assignment
	assignmentHandler := handler.NewAssignmentHandler(services.Assignment)
	//AssetLog
	assetLogHandler := handler.NewAssetLogHandler(services.AssetLog)
	//Request Transfer
	requestTransferHandler := handler.NewRequestTransferHandler(services.RequestTransfer)
	//Maintenance
	maintenanceHandler := handler.NewMaintenanceSchedulesHandler(services.MaintenanceSchedules)
	// Notification
	notificationsHandler := handler.NewNotificationHandler(services.Notification)
	//CronjobTest
	cronJobTestHandler := handler.NewCronJobTestHandler(db, services.Email, repos.Assets, repos.User, services.Notification, repos.AssetsLog)
	docs.SwaggerInfo.Title = "API Tool device manage"
	docs.SwaggerInfo.Description = "App Tool device manage"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = config.BASE_URL_BACKEND_FOR_SWAGGER
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	r := gin.Default()
	pprof.Register(r)
	api.SetupRoutes(r, userHandler, locationHandler, categoriesHandler, departmentHandler, assetsHandler, roleHandler, assignmentHandler, assetLogHandler, requestTransferHandler, maintenanceHandler, SSeHandler, notificationsHandler, cronJobTestHandler, repos.UserSession, db)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	cronjob.InitCronJobs(db, services.Email, repos.Assets, repos.User, services.Notification, repos.AssetsLog)

	if err := r.Run(config.Port); err != nil {
		log.Fatal("failed to run server:", err)
	}

}
