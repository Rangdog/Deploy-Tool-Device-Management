package api

import (
	"BE_Manage_device/api/handler"

	"github.com/gin-gonic/gin"
)

func registerCronJobTestRoutes(api *gin.RouterGroup, h *handler.CronJobTestHandler) {

	api.GET("/CheckAndSenMaintenanceNotification", h.CheckAndSenMaintenanceNotification)
	api.GET("/SendEmailsForWarrantyExpiry", h.SendEmailsForWarrantyExpiry)
	api.GET("/UpdateStatusWhenFinishMaintenance", h.UpdateStatusWhenFinishMaintenance)

}
