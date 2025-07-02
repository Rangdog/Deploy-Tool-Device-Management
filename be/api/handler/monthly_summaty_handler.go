package handler

import (
	"BE_Manage_device/constant"
	"BE_Manage_device/internal/domain/filter"
	service "BE_Manage_device/internal/service/monthly_summary"
	"BE_Manage_device/pkg"
	"BE_Manage_device/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type MonthlySummaryHandler struct {
	service service.MonthlySummaryService
}

func NewMonthlySummry(service *service.MonthlySummaryService) *MonthlySummaryHandler {
	return &MonthlySummaryHandler{
		service: *service,
	}
}

// Bill godoc
// @Summary Get monthly summary with filter
// @Description Get monthly summary have permission
// @Tags Monthly Summary
// @Accept json
// @Produce json
// @Param        bill   query    filter.MonthlySummaryFilter   false  "filter MonthlySummaryFilter"
// @param Authorization header string true "Authorization"
// @Router /api/monthly-summary/filter [GET]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *MonthlySummaryHandler) Filter(c *gin.Context) {
	defer pkg.PanicHandler(c)
	var filter filter.MonthlySummaryFilter
	userId := utils.GetUserIdFromContext(c)
	if err := c.ShouldBindQuery(&filter); err != nil {
		log.Error("Happened error when mapping query to filter. Error", err)
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when mapping query to filter")
	}
	data, err := h.service.Filter(userId, filter.Month, filter.Year)
	if err != nil {
		log.Error("Happened error when filter monthly summary. Error: ", err.Error())
		pkg.PanicExeption(constant.UnknownError, "Happened error when filter monthly summary. Error: "+err.Error())
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(http.StatusOK, constant.Success, data))
}
