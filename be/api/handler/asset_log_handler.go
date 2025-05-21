package handler

import (
	"BE_Manage_device/constant"
	"BE_Manage_device/internal/domain/service"
	"BE_Manage_device/pkg"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type AssetLogHandler struct {
	service *service.AssetLogService
}

func NewAssetLogHandler(service *service.AssetLogService) *AssetLogHandler {
	return &AssetLogHandler{service: service}
}

// Asset godoc
// @Summary Get assets log by id
// @Description Get assets log by id
// @Tags Assets log
// @Accept json
// @Produce json
// @Param		id	path		string				true	"id"
// @Router /api/assets-log/{id} [GET]
func (h *AssetLogHandler) GetLogByAssetId(c *gin.Context) {
	id := c.Param("id")
	defer pkg.PanicHandler(c)
	assetId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Error("Happened error when get id via path. Error", err)
		pkg.PanicExeption(constant.UnknownError, "Happened error when get id via path")
	}
	assetLogs, err := h.service.GetLogByAssetId(assetId)
	if err != nil {
		log.Error("Happened error when get asset log. Error", err)
		pkg.PanicExeption(constant.UnknownError, "Happened error when get asset log")
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(http.StatusOK, constant.Success, assetLogs))
}
