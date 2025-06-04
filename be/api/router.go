package api

import (
	"BE_Manage_device/api/handler"
	"BE_Manage_device/api/middleware"
	"BE_Manage_device/config"
	"BE_Manage_device/internal/domain/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, userHandler *handler.UserHandler, LocationHandler *handler.LocationHandler, CategoriesHandler *handler.CategoriesHandler, DepartmentsHandler *handler.DepartmentsHandler, AssetsHandler *handler.AssetsHandler, RoleHandler *handler.RoleHandler, AssignmentHandler *handler.AssignmentHandler, AssetLogHandler *handler.AssetLogHandler, RequestTransferHandler *handler.RequestTransferHandler, MaintenanceSchedulesHandler *handler.MaintenanceSchedulesHandler, SSEHandler *handler.SSEHandler, NotificationHandler *handler.NotificationHandler, session repository.UsersSessionRepository, db *gorm.DB) {
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
	api.POST("/notify/:userId", SSEHandler.SendNotificationHandler)
	api.Use(middleware.AuthMiddleware(config.AccessSecret, session))

	api.GET("/user/department/:department_id", userHandler.GetAllUserOfDepartment)
	api.GET("/user/session", userHandler.Session)                                                                           // đã check: nên chỉnh lại api response
	api.POST("/auth/logout", userHandler.Logout)                                                                            // đã check
	api.GET("/users", userHandler.GetAllUser)                                                                               // đã check:
	api.PATCH("users/information", userHandler.UpdateInformationUser)                                                       // đã check
	api.PATCH("users/role", middleware.RequirePermission([]string{"role-assignment"}, nil, db), userHandler.UpdateRoleUser) // đã check
	api.PATCH("/user/department", middleware.RequirePermission([]string{"user-management"}, nil, db), userHandler.UpdateDepartment)
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
	api.PATCH("/assets-retired/:id", middleware.RequirePermission([]string{"Update lifecycle"}, nil, db), AssetsHandler.UpdateAssetRetired) // đã check
	api.GET("/assets/filter-dashboard", AssetsHandler.FilterAssetDashboard)                                                                 // đã check
	api.GET("/assets/request-transfer", AssetsHandler.GetAssetsByCateOfDepartment)

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
	api.POST("/request-transfer", RequestTransferHandler.Create)                      // đã check
	api.PATCH("/request-transfer/confirm/:id", RequestTransferHandler.Accept)         // đã check
	api.PATCH("/request-transfer/deny/:id", RequestTransferHandler.Deny)              // đã check
	api.GET("/request-transfer/:id", RequestTransferHandler.GetRequestTransferById)   // đã check
	api.GET("/request-transfer/filter", RequestTransferHandler.FilterRequestTransfer) // đã check

	// Schedule maintenance
	api.POST("/maintenance-schedules", MaintenanceSchedulesHandler.Create)
	api.GET("/maintenance-schedules/:id", MaintenanceSchedulesHandler.GetAllMaintenanceSchedulesByAssetId)
	api.PATCH("/maintenance-schedules/:id", MaintenanceSchedulesHandler.Update)
	api.DELETE("/maintenance-schedules/:id", MaintenanceSchedulesHandler.Delete)
	api.GET("/maintenance-schedules", MaintenanceSchedulesHandler.GetAllMaintenanceSchedules)

	// SSE
	api.GET("/sse", SSEHandler.SSEHandle)

	//Notifications
	api.GET("/notifications", NotificationHandler.GetNotificationsByUserId)
	//Notifications
	api.GET("/notifications/:id", NotificationHandler.UpdateStatusToSeen)
}
