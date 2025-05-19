package main

import (
	"BE_Manage_device/api"
	"BE_Manage_device/api/handler"
	"BE_Manage_device/config"
	"BE_Manage_device/internal/database"
	"BE_Manage_device/internal/domain/service"
	"log"

	"BE_Manage_device/cmd/server/docs"

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

	emailService := service.NewEmailService(config.SmtpPasswd)

	//User
	userService := service.NewUserService(userRepository, emailService, userSessionRepository)
	userHandler := handler.NewUserHandler(userService)
	//Location
	locationService := service.NewLocationService(locationRepository)
	locationHandler := handler.NewLocationHandler(locationService)
	//Categories
	categoriesService := service.NewCategoriesService(categoriesRepository)
	categoriesHandler := handler.NewCategoriesHandler(categoriesService)

	docs.SwaggerInfo.Title = "API Tool device manage"
	docs.SwaggerInfo.Description = "App Tool device manage"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = config.BASE_URL_BACKEND
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	r := gin.Default()
	api.SetupRoutes(r, userHandler, locationHandler, categoriesHandler, userSessionRepository)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(config.Port); err != nil {
		log.Fatal("failed to run server:", err)
	}

}
