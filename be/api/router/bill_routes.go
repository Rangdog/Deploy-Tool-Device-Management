package api

import (
	"BE_Manage_device/api/handler"
	"BE_Manage_device/api/middleware"
	"BE_Manage_device/config"
	repository "BE_Manage_device/internal/repository/user_session"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func registerBillsRoutes(api *gin.RouterGroup, h *handler.BillsHandler, session repository.UsersSessionRepository, db *gorm.DB) {
	api.Use(middleware.AuthMiddleware(config.AccessSecret, session))

	api.POST("/bills", middleware.RequirePermission([]string{"manage-taxonomy"}, nil, db), h.Create)
	api.GET("/bills/:billNumber", middleware.RequirePermission([]string{"manage-taxonomy"}, nil, db), h.GetByBillNumber)
	api.GET("/bills/filter", middleware.RequirePermission([]string{"manage-taxonomy"}, nil, db), h.FilterBill)
	api.GET("/bills-un-paid", middleware.RequirePermission([]string{"manage-taxonomy"}, nil, db), h.GetAllBillUnpaid)
	api.PATCH("/bills/:billNumber", middleware.RequirePermission([]string{"manage-taxonomy"}, nil, db), h.UpdatePaid)
}
