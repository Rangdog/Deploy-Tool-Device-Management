package handler

import (
	"BE_Manage_device/constant"
	"BE_Manage_device/internal/domain/dto"
	"BE_Manage_device/internal/domain/service"
	"BE_Manage_device/pkg"
	"BE_Manage_device/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type MaintenanceschedulesHandler struct {
	service *service.MaintenanceSchedulesService
}

func NewMaintenanceSchedulesHandler(service *service.MaintenanceSchedulesService) *MaintenanceschedulesHandler {
	return &MaintenanceschedulesHandler{service: service}
}

// Maintenance Schedules godoc
// @Summary      Create maintenanceSchedules
// @Description  Create maintenanceSchedules
// @Tags         MaintenanceScheduless
// @Accept       json
// @Produce      json
// @Param        MaintenanceSchedules   body    dto.CreateMaintenanceSchedulesRequest   true  "Data"
// @param Authorization header string true "Authorization"
// @Router       /api/maintenance-schedules [POST]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *MaintenanceschedulesHandler) Create(c *gin.Context) {
	defer pkg.PanicHandler(c)
	userId := utils.GetUserIdFromContext(c)
	var request dto.CreateMaintenanceSchedulesRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error("Happened error when mapping request from FE. Error", err)
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when mapping request from FE.")
	}
	maintenance, err := h.service.Create(userId, request.AssetId, request.StartDate, request.EndDate)
	if err != nil {
		log.Error("Happened error when create maintenance. Error", err)
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when create maintenance.")
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(http.StatusOK, constant.Success, maintenance))
}

// Maintenance Schedules godoc
// @Summary      Get maintenanceSchedules by assetId
// @Description  Get maintenanceSchedules
// @Tags         MaintenanceScheduless
// @Accept       json
// @Produce      json
// @Param		id	path		int				true	"asset_id"
// @param Authorization header string true "Authorization"
// @Router       /api/maintenance-schedules/{id} [GET]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *MaintenanceschedulesHandler) GetAllMaintenanceSchedulesByAssetId(c *gin.Context) {
	defer pkg.PanicHandler(c)
	userId := utils.GetUserIdFromContext(c)
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("Happened error when convert project id to int64. Error", err)
		pkg.PanicExeption(constant.InvalidRequest)
	}
	maintenances, err := h.service.GetAllMaintenanceSchedulesByAssetId(userId, id)
	if err != nil {
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when get maintenance.")
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(http.StatusOK, constant.Success, maintenances))
}
