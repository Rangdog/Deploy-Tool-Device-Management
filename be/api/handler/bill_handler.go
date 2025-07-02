package handler

import (
	"BE_Manage_device/constant"
	"BE_Manage_device/internal/domain/dto"
	"BE_Manage_device/internal/domain/filter"
	service "BE_Manage_device/internal/service/bill"
	"BE_Manage_device/pkg"
	"BE_Manage_device/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type BillsHandler struct {
	service service.BillsService
}

func NewBillHandler(service *service.BillsService) *BillsHandler {
	return &BillsHandler{service: *service}
}

// User godoc
// @Summary      Create bill
// @Description  Create bill
// @Tags         Bills
// @Accept       json
// @Produce      json
// @Param        bill   body    dto.BillCreateRequest   true  "Data"
// @param Authorization header string true "Authorization"
// @Router       /api/bills [POST]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *BillsHandler) Create(c *gin.Context) {
	defer pkg.PanicHandler(c)
	var request dto.BillCreateRequest
	userId := utils.GetUserIdFromContext(c)
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error("Happened error when mapping request from FE. Error", err)
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when mapping request from FE.")
	}
	bill, err := h.service.Create(userId, request.AssetId, request.Description)
	if err != nil {
		log.Error("Happened error when create bill. Error", err.Error())
		pkg.PanicExeption(constant.UnknownError, "Happened error when create bill. Error: "+err.Error())
	}
	billResponse := utils.ConvertBillToResponse(bill)
	c.JSON(http.StatusCreated, pkg.BuildReponseSuccess(http.StatusCreated, constant.Success, billResponse))
}

// Bill godoc
// @Summary      Get bill by bill number
// @Description   Get bill by bill number
// @Tags         Bills
// @Accept       json
// @Produce      json
// @Param		billNumber	path		string				true	"billNumber"
// @param Authorization header string true "Authorization"
// @Router       /api/bills/{billNumber} [GET]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h BillsHandler) GetByBillNumber(c *gin.Context) {
	defer pkg.PanicHandler(c)
	billNumber := c.Param("billNumber")
	bill, err := h.service.GetByBillNumber(billNumber)
	if err != nil {
		log.Error("Happened error when get bill by bill number. Error", err.Error())
		pkg.PanicExeption(constant.UnknownError, "Happened error when get bill by bill number. Error: "+err.Error())
	}
	billResponse := utils.ConvertBillToResponse(bill)
	c.JSON(http.StatusCreated, pkg.BuildReponseSuccess(http.StatusCreated, constant.Success, billResponse))
}

// Bill godoc
// @Summary Get all bill with filter
// @Description Get all bill have permission
// @Tags Bills
// @Accept json
// @Produce json
// @Param        bill   query    filter.BillFilter   false  "filter bill"
// @param Authorization header string true "Authorization"
// @Router /api/bills/filter [GET]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *BillsHandler) FilterBill(c *gin.Context) {
	defer pkg.PanicHandler(c)
	var filter filter.BillFilter
	userId := utils.GetUserIdFromContext(c)
	if err := c.ShouldBindQuery(&filter); err != nil {
		log.Error("Happened error when mapping query to filter. Error", err)
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when mapping query to filter")
	}
	data, err := h.service.Filter(userId, filter.BillNumber, filter.Status, filter.CategoryId)
	if err != nil {
		log.Error("Happened error when filter bill. Error: ", err.Error())
		pkg.PanicExeption(constant.UnknownError, "Happened error when filter bill. Error: "+err.Error())
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(http.StatusOK, constant.Success, data))
}
