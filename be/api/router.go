package api

import (
	"BE_Manage_device/api/handler"
	"BE_Manage_device/api/middleware"
	"BE_Manage_device/config"
	"BE_Manage_device/internal/domain/repository"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, userHandler *handler.UserHandler, LocationHandler *handler.LocationHandler, CategoriesHandler *handler.CategoriesHandler, DepartmentsHandler *handler.DepartmentsHandler, AssetsHandler *handler.AssetsHandler, RoleHandler *handler.RoleHandler, AssignmentHandler *handler.AssignmentHandler, AssetLogHandler *handler.AssetLogHandler, RequestTransferHandler *handler.RequestTransferHandler, MaintenanceSchedulesHandler *handler.MaintenanceSchedulesHandler, session repository.UsersSessionRepository) {
	//users
	r.Use(middleware.CORSMiddleware())
	api := r.Group("/api")
	api.POST("/auth/register", userHandler.Register)                  // đã check
	api.POST("/auth/login", userHandler.Login)                        //dã check
	api.POST("/auth/refresh", userHandler.Refresh)                    // đã check
	api.GET("/activate", userHandler.Activate)                        // đã check
	api.POST("/user/forget-password", userHandler.CheckPasswordReset) // đã check
	api.PATCH("/user/password-reset", userHandler.ResetPassword)      // đã check
	api.DELETE("/user/:email", userHandler.DeleteUser)

	api.Use(middleware.AuthMiddleware(config.AccessSecret, session))

	api.GET("/user/session", userHandler.Session)                     // đã check: nên chỉnh lại api response
	api.POST("/auth/logout", userHandler.Logout)                      // đã check
	api.GET("/users", userHandler.GetAllUser)                         // đã check:
	api.PATCH("users/information", userHandler.UpdateInformationUser) // đã check
	api.PATCH("users/role", userHandler.UpdateRoleUser)               // đã check
	//Locations
	api.POST("/locations", LocationHandler.Create)       // đã check
	api.GET("/locations", LocationHandler.GetAll)        // đã check
	api.DELETE("/locations/:id", LocationHandler.Delete) // đã check

	//Categories
	api.POST("/categories", CategoriesHandler.Create)       // đã check
	api.GET("/categories", CategoriesHandler.GetAll)        // đã check
	api.DELETE("/categories/:id", CategoriesHandler.Delete) // đã check

	//Departments
	api.POST("/departments", DepartmentsHandler.Create)       // đã check
	api.GET("/departments", DepartmentsHandler.GetAll)        // đã check
	api.DELETE("/departments/:id", DepartmentsHandler.Delete) // đã check

	//Assets
	api.POST("/assets", AssetsHandler.Create)            // đã check
	api.GET("/assets/:id", AssetsHandler.GetAssetById)   // đã check
	api.GET("/assets", AssetsHandler.GetAllAsset)        // đã check
	api.GET("/assets/filter", AssetsHandler.FilterAsset) // đã check
	api.PUT("/assets/:id", AssetsHandler.Update)         // đã check
	api.DELETE("/assets/:id", AssetsHandler.DeleteAsset)
	api.PATCH("/assets-retired/:id", AssetsHandler.UpdateAssetRetired)      // đã check
	api.GET("/assets/filter-dashboard", AssetsHandler.FilterAssetDashboard) // đã check

	//Roles
	api.GET("/roles", RoleHandler.GetAllRole) // đã check

	//Assignment
	api.POST("/assignments", AssignmentHandler.Create)
	api.PUT("/assignments/:id", AssignmentHandler.Update)            // đã check
	api.GET("/assignments/filter", AssignmentHandler.FilterAsset)    // đã check
	api.GET("/assignments/:id", AssignmentHandler.GetAssignmentById) // đã check

	//AssetsLog
	api.GET("/assets-log/:id", AssetLogHandler.GetLogByAssetId) // đã check

	//Request
	api.POST("/request-transfer", RequestTransferHandler.Create)                    // đã check
	api.POST("/request-transfer/accept/:id", RequestTransferHandler.Accept)         // đã check
	api.POST("/request-transfer/deny/:id", RequestTransferHandler.Deny)             // đã check
	api.GET("/request-transfer/:id", RequestTransferHandler.GetRequestTransferById) // đã check
	api.GET("/request-transfer/filter", RequestTransferHandler.FilterRequestTransfer)

	// Schedule maintenance
	api.POST("/maintenance", MaintenanceSchedulesHandler.Create)
	api.GET("/maintenance/:id", MaintenanceSchedulesHandler.GetAllMaintenanceSchedulesByAssetId)
}
