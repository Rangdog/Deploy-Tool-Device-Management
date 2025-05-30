package handler

import (
	"BE_Manage_device/constant"
	"BE_Manage_device/internal/domain/service"
	"BE_Manage_device/pkg"
	"BE_Manage_device/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type NotificationHandler struct {
	service *service.NotificationService
}

func NewNotificationHandler(service *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

// Notification godoc
// @Summary      Get notifications by user id
// @Description  Get notifications by user id
// @Tags         Notification
// @Accept       json
// @Produce      json
// @param Authorization header string true "Authorization"
// @Router       /api/notifications [GET]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *NotificationHandler) GetNotificationsByUserId(c *gin.Context) {
	userId := utils.GetUserIdFromContext(c)
	notifications, err := h.service.GetNotificationsByUserId(userId)
	if err != nil {
		log.Error("Happened error when get notification.")
		pkg.PanicExeption(constant.UnknownError, "Happened error when get notification.")
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(http.StatusOK, constant.Success, notifications))
}

// Notification godoc
// @Summary      Update notification
// @Description  Update notification
// @Tags         Notification
// @Accept       json
// @Produce      json
// @Param		id	path		string				true	"id"
// @param Authorization header string true "Authorization"
// @Router       /api/notifications/{id} [PUT]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *NotificationHandler) UpdateStatusToSeen(c *gin.Context) {
	idStr := c.Param("id")
	NotificationId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("Happened error when convert notification to int64. Error", err)
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when convert notification to int64")
	}
	err = h.service.UpdateStatusToSeen(NotificationId)
	if err != nil {
		log.Error("Happened error when update notification.")
		pkg.PanicExeption(constant.UnknownError, "Happened error when update notification.")
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccessNoData(http.StatusOK, constant.Success))
}
