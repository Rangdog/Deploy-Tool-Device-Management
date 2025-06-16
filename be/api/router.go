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
	api.GET("/user/session", userHandler.Session)                                                                            // đã check: nên chỉnh lại api response
	api.POST("/auth/logout", userHandler.Logout)                                                                             // đã check
	api.GET("/users", userHandler.GetAllUser)                                                                                // đã check:
	api.PATCH("/users/information", userHandler.UpdateInformationUser)                                                       // đã check
	api.PATCH("/users/role", middleware.RequirePermission([]string{"role-assignment"}, nil, db), userHandler.UpdateRoleUser) // đã check
	api.PATCH("/user/department", middleware.RequirePermission([]string{"user-management"}, nil, db), userHandler.UpdateDepartment)
	api.PATCH("/user/head-department/:user_id", middleware.RequirePermission([]string{"user-management"}, nil, db), userHandler.UpdateHeadDep)
	api.PATCH("/user/manager-department/:user_id", middleware.RequirePermission([]string{"user-management"}, nil, db), userHandler.UpdateManagerDep)
	api.PATCH("/user/can-export/:user_id", middleware.RequirePermission([]string{"user-management"}, nil, db), userHandler.UpdateCanExport)
	//Locations
	api.POST("/locations", middleware.RequirePermission([]string{"manage-taxonomy"}, nil, db), LocationHandler.Create)       // đã check
	api.GET("/locations", LocationHandler.GetAll)                                                                            // đã check
	api.DELETE("/locations/:id", middleware.RequirePermission([]string{"manage-taxonomy"}, nil, db), LocationHandler.Delete) // đã check

	//Categories
	api.POST("/categories", middleware.RequirePermission([]string{"manage-taxonomy"}, nil, db), CategoriesHandler.Create)       // đã check
	api.GET("/categories", CategoriesHandler.GetAll)                                                                            // đã check
	api.DELETE("/categories/:id", middleware.RequirePermission([]string{"manage-taxonomy"}, nil, db), CategoriesHandler.Delete) // đã check

	//Departments
	api.POST("/departments", middleware.RequirePermission([]string{"manage-taxonomy"}, nil, db), DepartmentsHandler.Create)       // đã check
	api.GET("/departments", DepartmentsHandler.GetAll)                                                                            // đã check
	api.DELETE("/departments/:id", middleware.RequirePermission([]string{"manage-taxonomy"}, nil, db), DepartmentsHandler.Delete) // đã check

	//Assets
	api.POST("/assets", middleware.RequirePermission([]string{"manage-assets"}, []string{"full", "limited"}, db), AssetsHandler.Create)    // đã check
	api.GET("/assets/:id", AssetsHandler.GetAssetById)                                                                                     // đã check
	api.GET("/assets", AssetsHandler.GetAllAsset)                                                                                          // đã check
	api.GET("/assets/filter", AssetsHandler.FilterAsset)                                                                                   // đã check
	api.PUT("/assets/:id", middleware.RequirePermission([]string{"manage-assets"}, []string{"full", "limited"}, db), AssetsHandler.Update) // đã check
	api.DELETE("/assets/:id", middleware.RequirePermission([]string{"manage-assets"}, nil, db), AssetsHandler.DeleteAsset)
	api.PATCH("/assets-retired/:id", middleware.RequirePermission([]string{"lifecycle-update", "manage-assets"}, nil, db), AssetsHandler.UpdateAssetRetired)      // đã check
	api.GET("/assets/filter-dashboard", middleware.RequirePermission([]string{"dashboards"}, []string{"full", "scoped"}, db), AssetsHandler.FilterAssetDashboard) // đã check
	api.GET("/assets/request-transfer", AssetsHandler.GetAssetsByCateOfDepartment)
	api.GET("/assets/maintenance-schedules", AssetsHandler.GetAllAssetNotHaveMaintenance)

	//Roles
	api.GET("/roles", RoleHandler.GetAllRole) // đã check

	//Assignment
	api.POST("/assignments", middleware.RequirePermission([]string{"assign-assets"}, nil, db), AssignmentHandler.Create)
	api.PUT("/assignments/:id", middleware.RequirePermission([]string{"assign-assets"}, nil, db), AssignmentHandler.Update)            // đã check
	api.GET("/assignments/filter", middleware.RequirePermission([]string{"assign-assets"}, nil, db), AssignmentHandler.FilterAsset)    // đã check
	api.GET("/assignments/:id", middleware.RequirePermission([]string{"assign-assets"}, nil, db), AssignmentHandler.GetAssignmentById) // đã check

	//AssetsLog
	api.GET("/assets-log/:id", middleware.RequirePermission([]string{"audit-logs"}, []string{"full", "partial"}, db), AssetLogHandler.GetLogByAssetId) // đã check

	//Request
	api.POST("/request-transfer", middleware.RequirePermission([]string{"transfer-assets"}, []string{"full", "can-request"}, db), RequestTransferHandler.Create) // đã check
	api.PATCH("/request-transfer/confirm/:id", middleware.RequirePermission([]string{"transfer-assets"}, nil, db), RequestTransferHandler.Accept)                // đã check
	api.PATCH("/request-transfer/deny/:id", middleware.RequirePermission([]string{"transfer-assets"}, nil, db), RequestTransferHandler.Deny)                     // đã check
	api.GET("/request-transfer/:id", middleware.RequirePermission([]string{"transfer-assets"}, nil, db), RequestTransferHandler.GetRequestTransferById)          // đã check
	api.GET("/request-transfer/filter", middleware.RequirePermission([]string{"transfer-assets"}, nil, db), RequestTransferHandler.FilterRequestTransfer)        // đã check

	// Schedule maintenance
	api.POST("/maintenance-schedules", middleware.RequirePermission([]string{"maintenance-logs"}, nil, db), MaintenanceSchedulesHandler.Create)
	api.GET("/maintenance-schedules/:id", middleware.RequirePermission([]string{"maintenance-logs"}, []string{"full", "view"}, db), MaintenanceSchedulesHandler.GetAllMaintenanceSchedulesByAssetId)
	api.PATCH("/maintenance-schedules/:id", middleware.RequirePermission([]string{"maintenance-logs"}, nil, db), MaintenanceSchedulesHandler.Update)
	api.DELETE("/maintenance-schedules/:id", middleware.RequirePermission([]string{"maintenance-logs"}, nil, db), MaintenanceSchedulesHandler.Delete)
	api.GET("/maintenance-schedules", middleware.RequirePermission([]string{"maintenance-logs"}, []string{"full", "view"}, db), MaintenanceSchedulesHandler.GetAllMaintenanceSchedules)

	// SSE
	api.GET("/sse", SSEHandler.SSEHandle)

	//Notifications
	api.GET("/notifications", NotificationHandler.GetNotificationsByUserId)
	//Notifications
	api.GET("/notifications/:id", NotificationHandler.UpdateStatusToSeen)
}
