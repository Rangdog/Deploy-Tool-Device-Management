package api

import (
	"BE_Manage_device/api/handler"
	"BE_Manage_device/api/middleware"
	"BE_Manage_device/config"
	repository "BE_Manage_device/internal/repository/user_session"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func registerUserRoutes(api *gin.RouterGroup, h *handler.UserHandler, session repository.UsersSessionRepository, db *gorm.DB) {
	api.Use(middleware.AuthMiddleware(config.AccessSecret, session))

	api.GET("/user/department/:department_id", h.GetAllUserOfDepartment)
	api.GET("/user/session", h.Session)
	api.POST("/auth/logout", h.Logout)
	api.GET("/users", h.GetAllUser)
	api.PATCH("/users/information", h.UpdateInformationUser)
	api.PATCH("/users/role", middleware.RequirePermission([]string{"role-assignment"}, nil, db), h.UpdateRoleUser)
	api.PATCH("/user/department", middleware.RequirePermission([]string{"user-management"}, nil, db), h.UpdateDepartment)
	api.PATCH("/user/head-department/:user_id", middleware.RequirePermission([]string{"user-management"}, nil, db), h.UpdateHeadDep)
	api.PATCH("/user/manager-department/:user_id", middleware.RequirePermission([]string{"user-management"}, nil, db), h.UpdateManagerDep)
	api.PATCH("/user/can-export/:user_id", middleware.RequirePermission([]string{"user-management"}, nil, db), h.UpdateCanExport)
	api.GET("/users/not-dep", h.GetUserNotHaveDep)
}
