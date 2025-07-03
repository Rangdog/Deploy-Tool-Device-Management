package api

import (
	"BE_Manage_device/api/handler"
	"BE_Manage_device/api/middleware"
	"BE_Manage_device/config"
	repository "BE_Manage_device/internal/repository/user_session"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func registerAssignmentRoutes(api *gin.RouterGroup, h *handler.AssignmentHandler, session repository.UsersSessionRepository, db *gorm.DB) {
	api.Use(middleware.AuthMiddleware(config.AccessSecret, session))

	api.POST("/assignments", middleware.RequirePermission([]string{"assign-assets"}, nil, db), h.Create)
	api.PUT("/assignments/:id", middleware.RequirePermission([]string{"assign-assets"}, []string{"full", "conditional"}, db), h.Update)              // đã check
	api.GET("/assignments/filter", middleware.RequirePermission([]string{"assign-assets"}, []string{"full", "conditional"}, db), h.FilterAssignment) // đã check
	api.GET("/assignments/:id", middleware.RequirePermission([]string{"assign-assets"}, []string{"full", "conditional"}, db), h.GetAssignmentById)   // đã check

}
