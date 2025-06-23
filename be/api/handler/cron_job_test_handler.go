package handler

import (
	"BE_Manage_device/constant"
	asset_log "BE_Manage_device/internal/repository/asset_log"
	asset "BE_Manage_device/internal/repository/assets"
	user "BE_Manage_device/internal/repository/user"
	emailS "BE_Manage_device/internal/service/email"
	notificationS "BE_Manage_device/internal/service/notification"
	"net/http"

	"BE_Manage_device/pkg"
	"BE_Manage_device/pkg/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CronJobTestHandler struct {
	db                   *gorm.DB
	emailService         *emailS.EmailService
	assetsRepository     asset.AssetsRepository
	userRepository       user.UserRepository
	notificationsService *notificationS.NotificationService
	assetsLogRepository  asset_log.AssetsLogRepository
}

func NewCronJobTestHandler(db *gorm.DB, emailService *emailS.EmailService, assetsRepository asset.AssetsRepository, userRepository user.UserRepository, notificationsService *notificationS.NotificationService, assetsLogRepository asset_log.AssetsLogRepository) *CronJobTestHandler {
	return &CronJobTestHandler{db: db, emailService: emailService, assetsRepository: assetsRepository, userRepository: userRepository, notificationsService: notificationsService, assetsLogRepository: assetsLogRepository}
}

// Cron godoc
// @Summary      CheckAndSenMaintenanceNotification
// @Description  CheckAndSenMaintenanceNotification
// @Tags         Cron
// @Accept       json
// @Produce      json
// @Router       /api/CheckAndSenMaintenanceNotification [GET]
func (h *CronJobTestHandler) CheckAndSenMaintenanceNotification(c *gin.Context) {
	defer pkg.PanicHandler(c)
	utils.CheckAndSenMaintenanceNotification(h.db, h.emailService, h.assetsRepository, h.userRepository, h.notificationsService, h.assetsLogRepository)
	c.JSON(http.StatusCreated, pkg.BuildReponseSuccessNoData(http.StatusCreated, constant.Success))
}

// Cron godoc
// @Summary      SendEmailsForWarrantyExpiry
// @Description  SendEmailsForWarrantyExpiry
// @Tags         Cron
// @Accept       json
// @Produce      json
// @Router       /api/SendEmailsForWarrantyExpiry [GET]
func (h *CronJobTestHandler) SendEmailsForWarrantyExpiry(c *gin.Context) {
	defer pkg.PanicHandler(c)
	utils.SendEmailsForWarrantyExpiry(h.db, h.emailService, h.notificationsService, h.assetsRepository, h.userRepository)
	c.JSON(http.StatusCreated, pkg.BuildReponseSuccessNoData(http.StatusCreated, constant.Success))
}

// Cron godoc
// @Summary      UpdateStatusWhenFinishMaintenance
// @Description  UpdateStatusWhenFinishMaintenance
// @Tags         Cron
// @Accept       json
// @Produce      json
// @Router       /api/UpdateStatusWhenFinishMaintenance [GET]
func (h *CronJobTestHandler) UpdateStatusWhenFinishMaintenance(c *gin.Context) {
	defer pkg.PanicHandler(c)
	utils.UpdateStatusWhenFinishMaintenance(h.db, h.assetsRepository, h.userRepository, h.notificationsService, h.assetsLogRepository)
	c.JSON(http.StatusCreated, pkg.BuildReponseSuccessNoData(http.StatusCreated, constant.Success))
}
