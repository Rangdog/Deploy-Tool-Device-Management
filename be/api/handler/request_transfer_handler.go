package handler

import (
	"BE_Manage_device/constant"
	"BE_Manage_device/internal/domain/dto"
	"BE_Manage_device/internal/domain/service"
	"BE_Manage_device/pkg"
	"BE_Manage_device/pkg/utils"
	"net/http"

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
// @Param		id	path		int				true	"project_id"
// @Router       /api/request-transfer/{id} [POST]
func (h *RequestTransferHandler) Create(c *gin.Context) {
	defer pkg.PanicHandler(c)
	userId := utils.GetUserIdFromContext(c)
	var request dto.CreateRequestTransferRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error("Happened error when mapping request from FE. Error", err)
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when mapping request from FE.")
	}
	requestTransfer, err := h.service.Create(userId, request.AssetId, request.DepartmentId)
	if err != nil {
		log.Error("Happened error when create request transfer. Error", err)
		pkg.PanicExeption(constant.UnknownError, "Happened error when create request transfer")
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(http.StatusOK, constant.Success, requestTransfer))
}
