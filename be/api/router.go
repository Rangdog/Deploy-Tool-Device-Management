package api

import (
	"BE_Manage_device/api/handler"
	"BE_Manage_device/api/middleware"
	"BE_Manage_device/config"
	"BE_Manage_device/internal/domain/repository"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, userHandler *handler.UserHandler, session repository.UsersSessionRepository) {
	//users
	r.Use(middleware.CORSMiddleware())
	api := r.Group("/api")
	api.POST("/auth/register", userHandler.Register)
	api.POST("/auth/login", userHandler.Login)
	api.POST("/auth/refresh", userHandler.Refresh)
	api.GET("/activate", userHandler.Activate)
	api.POST("/user/password-reset", userHandler.CheckPasswordReset)
	api.PATCH("/user/forget-password", userHandler.ResetPassword)
	api.DELETE("/user/:email", userHandler.DeleteUser)
	api.Use(middleware.AuthMiddleware(config.AccessSecret, session))
	api.GET("/user/session", userHandler.Session)
	api.POST("/auth/logout", userHandler.Logout)
}
