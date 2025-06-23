package api

import (
	"BE_Manage_device/api/handler"
	"BE_Manage_device/api/middleware"
	repository "BE_Manage_device/internal/repository/user_session"

	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, userHandler *handler.UserHandler, LocationHandler *handler.LocationHandler, CategoriesHandler *handler.CategoriesHandler, DepartmentsHandler *handler.DepartmentsHandler, AssetsHandler *handler.AssetsHandler, RoleHandler *handler.RoleHandler, AssignmentHandler *handler.AssignmentHandler, AssetLogHandler *handler.AssetLogHandler, RequestTransferHandler *handler.RequestTransferHandler, MaintenanceSchedulesHandler *handler.MaintenanceSchedulesHandler, SSEHandler *handler.SSEHandler, NotificationHandler *handler.NotificationHandler, CronJobTestHandler *handler.CronJobTestHandler, session repository.UsersSessionRepository, db *gorm.DB) {
	//users
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.TimeoutMiddleware(5 * time.Second))
	api := r.Group("/api")
	registerCronJobTestRoutes(api, CronJobTestHandler)
	registerAuthRoutes(api, userHandler, SSEHandler)
	registerUserRoutes(api, userHandler, session, db)
	registerLocationsRoutes(api, LocationHandler, session, db)
	registerCategoriesRoutes(api, CategoriesHandler, session, db)
	registerDepartmentRoutes(api, DepartmentsHandler, session, db)
	registerAssetsRoutes(api, AssetsHandler, session, db)
	registerRoleRoutes(api, RoleHandler, session, db)
	registerAssignmentRoutes(api, AssignmentHandler, session, db)
	registerAssetLogsRoutes(api, AssetLogHandler, session, db)
	registerRequestTransferRoutes(api, RequestTransferHandler, session, db)
	registerMaintenanceSchedulesRoutes(api, MaintenanceSchedulesHandler, session, db)
	registerNotificationsRoutes(api, NotificationHandler, session, db)
	registerSSEHandlerRoutes(api, SSEHandler, session, db)
}
