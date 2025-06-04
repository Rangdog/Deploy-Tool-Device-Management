package handler

import (
	"BE_Manage_device/constant"
	"BE_Manage_device/internal/domain/dto"
	"BE_Manage_device/internal/domain/filter"
	"BE_Manage_device/internal/domain/service"
	"BE_Manage_device/pkg"
	"BE_Manage_device/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type RequestTransferHandler struct {
	service *service.RequestTransferService
}

func NewRequestTransferHandler(service *service.RequestTransferService) *RequestTransferHandler {
	return &RequestTransferHandler{service: service}
}

// Request Transfer godoc
// @Summary      Request Transfer
// @Description  Request Transfer
// @Tags         RequestTransfer
// @Accept       json
// @Produce      json
// @Param        Request-Transfer  body    dto.CreateRequestTransferRequest   true  "Data"
// @param Authorization header string true "Authorization"
// @Router       /api/request-transfer [POST]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *RequestTransferHandler) Create(c *gin.Context) {
	defer pkg.PanicHandler(c)
	userId := utils.GetUserIdFromContext(c)
	var request dto.CreateRequestTransferRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error("Happened error when mapping request from FE. Error", err)
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when mapping request from FE.")
	}
	requestTransfer, err := h.service.Create(userId, request.CategoryId, request.Description)
	if err != nil {
		log.Error("Happened error when create request transfer. Error", err)
		pkg.PanicExeption(constant.UnknownError, "Happened error when create request transfer")
	}
	requestTransferResponse := utils.ConvertRequestTransferToResponse(requestTransfer)
	c.JSON(http.StatusCreated, pkg.BuildReponseSuccess(http.StatusCreated, constant.Success, requestTransferResponse))
}

// Request Transfer godoc
// @Summary      Accept Request Transfer
// @Description  Accept Request Transfer
// @Tags         RequestTransfer
// @Accept       json
// @Produce      json
// @Param        Asset   body    dto.ConfirmRequestTransferRequest   true  "Data"
// @Param		id	path		int				true	"request_transfer_id"
// @param Authorization header string true "Authorization"
// @Router       /api/request-transfer/confirm/{id} [PATCH]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *RequestTransferHandler) Accept(c *gin.Context) {
	defer pkg.PanicHandler(c)
	userId := utils.GetUserIdFromContext(c)
	var request dto.ConfirmRequestTransferRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error("Happened error when mapping request from FE. Error", err)
		pkg.PanicExeption(constant.InvalidRequest)
	}
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("Happened error when convert id to int64. Error", err)
		pkg.PanicExeption(constant.InvalidRequest)
	}
	requestTransfer, err := h.service.Accept(userId, id, request.AssetId)
	if err != nil {
		log.Error("Happened error when accept request transfer. Error", err)
		pkg.PanicExeption(constant.UnknownError, "Happened error when accept request transfer")
	}
	requestTransferResponse := utils.ConvertRequestTransferToResponse(requestTransfer)
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(http.StatusOK, constant.Success, requestTransferResponse))
}

// Request Transfer godoc
// @Summary      Deny Request Transfer
// @Description  Deny Request Transfer
// @Tags         RequestTransfer
// @Accept       json
// @Produce      json
// @Param		id	path		int				true	"request_transfer_id"
// @param Authorization header string true "Authorization"
// @Router       /api/request-transfer/deny/{id} [PATCH]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *RequestTransferHandler) Deny(c *gin.Context) {
	defer pkg.PanicHandler(c)
	userId := utils.GetUserIdFromContext(c)
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("Happened error when convert id to int64. Error", err)
		pkg.PanicExeption(constant.InvalidRequest)
	}
	requestTransfer, err := h.service.Deny(userId, id)
	if err != nil {
		log.Error("Happened error when deny request transfer. Error", err)
		pkg.PanicExeption(constant.UnknownError, "Happened error when deny request transfer")
	}
	requestTransferResponse := utils.ConvertRequestTransferToResponse(requestTransfer)
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(http.StatusOK, constant.Success, requestTransferResponse))
}

// Request transfer godoc
// @Summary      Get request transfer
// @Description  Get request transfer
// @Tags         RequestTransfer
// @Accept       json
// @Produce      json
// @Param		id	path		string				true	"id"
// @param Authorization header string true "Authorization"
// @Router       /api/request-transfer/{id} [GET]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *RequestTransferHandler) GetRequestTransferById(c *gin.Context) {
	defer pkg.PanicHandler(c)
	userId := utils.GetUserIdFromContext(c)
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("Happened error when convert id to int64. Error", err)
		pkg.PanicExeption(constant.InvalidRequest)
	}
	requestTransfer, err := h.service.GetRequestTransferById(userId, id)
	if err != nil {
		log.Error("Happened error when get request transfer. Error", err)
		pkg.PanicExeption(constant.UnknownError, "Happened error when get request transfer")
	}
	requestTransferResponse := utils.ConvertRequestTransferToResponse(requestTransfer)
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(http.StatusOK, constant.Success, requestTransferResponse))
}

// Request Transfer godoc
// @Summary Get all RequestTransfer with filter
// @Description Get all request transfer have permission
// @Tags RequestTransfer
// @Accept json
// @Produce json
// @Param        request_transfer   query    filter.RequestTransferFilter   false  "filter request transfer"
// @param Authorization header string true "Authorization"
// @Router /api/request-transfer/filter [GET]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *RequestTransferHandler) FilterRequestTransfer(c *gin.Context) {
	defer pkg.PanicHandler(c)
	var filter filter.RequestTransferFilter
	userId := utils.GetUserIdFromContext(c)
	if err := c.ShouldBindQuery(&filter); err != nil {
		log.Error("Happened error when mapping query to filter. Error", err)
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when mapping query to filter")
	}
	data, err := h.service.Filter(userId, filter.Status)
	if err != nil {
		log.Error("Happened error when filter request transfer. Error", err)
		pkg.PanicExeption(constant.UnknownError, "Happened error when filter request transfer")
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(http.StatusOK, constant.Success, data))
}
