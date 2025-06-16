package handler

import (
	"BE_Manage_device/constant"
	"BE_Manage_device/internal/domain/filter"
	"BE_Manage_device/internal/domain/service"
	"BE_Manage_device/pkg"
	"BE_Manage_device/pkg/utils"
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
// @Param		asset_id	path		string				true	"id"
// @Param        asset   query    filter.AssetLogFilter   false  "filter asset"
// @param Authorization header string true "Authorization"
// @Router /api/assets-log/{asset_id} [GET]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *AssetLogHandler) GetLogByAssetId(c *gin.Context) {
	defer pkg.PanicHandler(c)
	userId := utils.GetUserIdFromContext(c)
	id := c.Param("id")
	assetId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Error("Happened error when get id via path. Error", err)
		pkg.PanicExeption(constant.UnknownError, "Happened error when get id via path")
	}
	var filter filter.AssetLogFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		log.Error("Happened error when mapping query to filter. Error", err)
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when mapping query to filter")
	}
	asset, err := h.service.GetAssetById(userId, assetId)
	if err != nil {
		log.Error("Happened error when get asset by id. Error", err)
		pkg.PanicExeption(constant.UnknownError, "Happened error when get asset by id")
	}
	err = h.service.CheckPermissionForManager(userId, asset.DepartmentId)
	if err != nil {
		pkg.PanicExeption(constant.InvalidRequest, err.Error())
		return
	}
	assetLogs, err := h.service.Filter(userId, assetId, filter.Action, filter.StartTime, filter.EndTime)
	if err != nil {
		log.Error("Happened error when get asset log. Error", err)
		pkg.PanicExeption(constant.UnknownError, "Happened error when get asset log")
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(http.StatusOK, constant.Success, assetLogs))
}
