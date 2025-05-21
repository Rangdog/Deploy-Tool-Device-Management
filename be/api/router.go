package api

import (
	"BE_Manage_device/api/handler"
	"BE_Manage_device/api/middleware"
	"BE_Manage_device/config"
	"BE_Manage_device/internal/domain/repository"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, userHandler *handler.UserHandler, LocationHandler *handler.LocationHandler, CategoriesHandler *handler.CategoriesHandler, DepartmentsHandler *handler.DepartmentsHandler, AssetsHandler *handler.AssetsHandler, RoleHandler *handler.RoleHandler, AssignmentHandler *handler.AssignmentHandler, session repository.UsersSessionRepository) {
	//users
	r.Use(middleware.CORSMiddleware())
	api := r.Group("/api")
	api.POST("/auth/register", userHandler.Register)
	api.POST("/auth/login", userHandler.Login)
	api.POST("/auth/refresh", userHandler.Refresh)
	api.GET("/activate", userHandler.Activate)
	api.POST("/user/forget-password", userHandler.CheckPasswordReset)
	api.PATCH("/user/password-reset", userHandler.ResetPassword)
	api.DELETE("/user/:email", userHandler.DeleteUser)

	api.Use(middleware.AuthMiddleware(config.AccessSecret, session))

	api.GET("/user/session", userHandler.Session)
	api.POST("/auth/logout", userHandler.Logout)
	api.GET("/users", userHandler.GetAllUser)
	api.PATCH("users/information", userHandler.UpdateInformationUser)
	api.PATCH("users/role", userHandler.UpdateRoleUser)
	//Locations
	api.POST("/locations", LocationHandler.Create)
	api.GET("/locations", LocationHandler.GetAll)
	api.DELETE("/locations/:id", LocationHandler.Delete)

	//Categories
	api.POST("/categories", CategoriesHandler.Create)
	api.GET("/categories", CategoriesHandler.GetAll)
	api.DELETE("/categories/:id", CategoriesHandler.Delete)

	//Departments
	api.POST("/departments", DepartmentsHandler.Create)
	api.GET("/departments", DepartmentsHandler.GetAll)
	api.DELETE("/departments/:id", DepartmentsHandler.Delete)

	//Assets
	api.POST("/assets", AssetsHandler.Create)
	api.GET("/assets/:id", AssetsHandler.GetAssetById)
	api.GET("/assets", AssetsHandler.GetAllAsset)
	api.GET("/assets/filter", AssetsHandler.FilterAsset)
	api.PUT("/assets/:id", AssetsHandler.Update)
	api.DELETE("/assets/:id", AssetsHandler.DeleteAsset)

	//Roles
	api.GET("/roles", RoleHandler.GetAllRole)

	//assignment
	api.POST("/assignment", AssignmentHandler.Create)
	api.PUT("/assignment/:id", AssignmentHandler.Update)
	api.GET("/assignment/filter", AssignmentHandler.FilterAsset)
}
