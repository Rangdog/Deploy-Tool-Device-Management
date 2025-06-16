package handler

import (
	"BE_Manage_device/constant"
	"BE_Manage_device/internal/domain/dto"
	"BE_Manage_device/internal/domain/service"
	"BE_Manage_device/pkg"
	"BE_Manage_device/pkg/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type MaintenanceSchedulesHandler struct {
	service *service.MaintenanceSchedulesService
}

func NewMaintenanceSchedulesHandler(service *service.MaintenanceSchedulesService) *MaintenanceSchedulesHandler {
	return &MaintenanceSchedulesHandler{service: service}
}

// Maintenance Schedules godoc
// @Summary      Create maintenanceSchedules
// @Description  Create maintenanceSchedules
// @Tags         MaintenanceSchedules
// @Accept       json
// @Produce      json
// @Param        MaintenanceSchedules   body    dto.CreateMaintenanceSchedulesRequest   true  "Data"
// @param Authorization header string true "Authorization"
// @Router       /api/maintenance-schedules [POST]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *MaintenanceSchedulesHandler) Create(c *gin.Context) {
	defer pkg.PanicHandler(c)
	userId := utils.GetUserIdFromContext(c)
	var request dto.CreateMaintenanceSchedulesRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error("Happened error when mapping request from FE. Error", err)
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when mapping request from FE.")
	}
	now := time.Now()
	if !request.StartDate.After(now) {
		log.Error("Happened error start date <= now .")
		pkg.PanicExeption(constant.InvalidRequest, "Happened error start date can't in past.")
	}
	if !request.EndDate.After(request.StartDate) {
		log.Error("Happened error start date >= end date .")
		pkg.PanicExeption(constant.InvalidRequest, "Happened error start date > end date.")
	}
	maintenance, err := h.service.Create(userId, request.AssetId, request.StartDate, request.EndDate)
	if err != nil {
		log.Error("Happened error when create maintenance. Error", err)
		pkg.PanicExeption(constant.UnknownError, "Happened error when create maintenance.")
	}
	MaintenanceScheduleRes := utils.ConvertMaintenanceSchedulesToResponses(maintenance)
	c.JSON(http.StatusCreated, pkg.BuildReponseSuccess(http.StatusCreated, constant.Success, MaintenanceScheduleRes))
}

// Maintenance Schedules godoc
// @Summary      Get maintenanceSchedules by assetId
// @Description  Get maintenanceSchedules
// @Tags         MaintenanceSchedules
// @Accept       json
// @Produce      json
// @Param		id	path		int				true	"asset_id"
// @param Authorization header string true "Authorization"
// @Router       /api/maintenance-schedules/{id} [GET]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *MaintenanceSchedulesHandler) GetAllMaintenanceSchedulesByAssetId(c *gin.Context) {
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
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when get maintenance by assetId.")
	}
	MaintenanceScheduleRes := utils.ConvertMaintenanceSchedulesToResponsesArray(maintenances)
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(http.StatusOK, constant.Success, MaintenanceScheduleRes))
}

// Maintenance Schedules godoc
// @Summary      Update maintenanceSchedules by id
// @Description  Update maintenanceSchedules by id
// @Tags         MaintenanceSchedules
// @Accept       json
// @Produce      json
// @Param        MaintenanceSchedules   body    dto.UpdateMaintenanceSchedulesRequest   true  "Data"
// @Param		id	path		int				true	"maintenance_id"
// @param Authorization header string true "Authorization"
// @Router       /api/maintenance-schedules/{id} [PATCH]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *MaintenanceSchedulesHandler) Update(c *gin.Context) {
	defer pkg.PanicHandler(c)
	userId := utils.GetUserIdFromContext(c)
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("Happened error when convert project id to int64. Error", err)
		pkg.PanicExeption(constant.InvalidRequest)
	}
	var request dto.UpdateMaintenanceSchedulesRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error("Happened error when mapping request from FE. Error", err)
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when mapping request from FE.")
	}
	now := time.Now()
	if !request.StartDate.After(now) {
		log.Error("Start date must be in the future.")
		pkg.PanicExeption(constant.InvalidRequest, "Start date must be after now.")
	}
	if !request.EndDate.After(request.StartDate) {
		log.Error("End date must be after start date.")
		pkg.PanicExeption(constant.InvalidRequest, "End date must be after start date.")
	}
	maintenance, err := h.service.Update(userId, id, request.StartDate, request.EndDate)
	if err != nil {
		log.Error("Happened error when create maintenance. Error", err)
		pkg.PanicExeption(constant.UnknownError, "Happened error when update maintenance.")
	}
	MaintenanceScheduleRes := utils.ConvertMaintenanceSchedulesToResponses(maintenance)
	c.JSON(http.StatusCreated, pkg.BuildReponseSuccess(http.StatusCreated, constant.Success, MaintenanceScheduleRes))
}

// Maintenance Schedules godoc
// @Summary      Delete maintenanceSchedules by id
// @Description  Delete maintenanceSchedules by id
// @Tags         MaintenanceSchedules
// @Accept       json
// @Produce      json
// @Param		id	path		int				true	"maintenance_id"
// @param Authorization header string true "Authorization"
// @Router       /api/maintenance-schedules/{id} [DELETE]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *MaintenanceSchedulesHandler) Delete(c *gin.Context) {
	defer pkg.PanicHandler(c)
	userId := utils.GetUserIdFromContext(c)
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("Happened error when convert project id to int64. Error", err)
		pkg.PanicExeption(constant.InvalidRequest)
	}
	err = h.service.Delete(userId, id)
	if err != nil {
		log.Error("Happened error when create maintenance. Error", err)
		pkg.PanicExeption(constant.UnknownError, "Happened error when delete maintenance.")
		c.JSON(http.StatusCreated, pkg.BuildReponseSuccessNoData(http.StatusCreated, constant.Success))
	}
}

// Maintenance Schedules godoc
// @Summary      Get maintenanceSchedules
// @Description  Get maintenanceSchedules
// @Tags         MaintenanceSchedules
// @Accept       json
// @Produce      json
// @param Authorization header string true "Authorization"
// @Router       /api/maintenance-schedules [GET]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *MaintenanceSchedulesHandler) GetAllMaintenanceSchedules(c *gin.Context) {
	defer pkg.PanicHandler(c)
	maintenances, err := h.service.GetAllMaintenanceSchedules()
	if err != nil {
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when get maintenance.")
	}
	MaintenanceScheduleRes := utils.ConvertMaintenanceSchedulesToResponsesArray(maintenances)
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(http.StatusOK, constant.Success, MaintenanceScheduleRes))
}
