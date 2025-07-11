package handler

import (
	"BE_Manage_device/constant"
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
// @Param assetId formData string true "Asset ID"
// @Param description formData string false "Description"
// @Param statusBill formData string true "Description"
// @Param file formData file false "File to upload"
// @Param image formData file false "Image to upload"
// @param Authorization header string true "Authorization"
// @Router       /api/bills [POST]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *BillsHandler) Create(c *gin.Context) {
	defer pkg.PanicHandler(c)
	var status string
	assetIdStr := c.PostForm("assetId")
	description := c.PostForm("description")
	statusStr := c.PostForm("statusBill")
	assetId, err := utils.ParseStrToInt64(assetIdStr)
	if err != nil {
		log.Info("Error: ", err.Error())
		pkg.PanicExeption(constant.InvalidRequest, "Invalid assetId format")
	}
	// handler status string
	if statusStr == "Unpaid" || statusStr == "Paid" {
		status = statusStr
	} else {
		log.Info("Error: ", err.Error())
		pkg.PanicExeption(constant.InvalidRequest, "Invalid status format")
	}
	file, err := c.FormFile("file")
	if err != nil {
		file = nil
	}

	image, err := c.FormFile("image")
	if err != nil {
		image = nil
	}
	userId := utils.GetUserIdFromContext(c)
	bill, err := h.service.Create(userId, assetId, description, image, file, status)
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

// Bill godoc
// @Summary Get un paid bill
// @Description un paid bill
// @Tags Bills
// @Accept json
// @Produce json
// @param Authorization header string true "Authorization"
// @Router /api/bills-un-paid/ [GET]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *BillsHandler) GetAllBillUnpaid(c *gin.Context) {
	defer pkg.PanicHandler(c)
	userId := utils.GetUserIdFromContext(c)
	bills, err := h.service.GetAllBillUnpaid(userId)
	if err != nil {
		log.Error("Happened error when get bill unpaid. Error: ", err.Error())
		pkg.PanicExeption(constant.UnknownError, "Happened error when get bill unpaid. Error: "+err.Error())
	}
	billsRes := utils.ConvertBillsToResponsesArray(bills)
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(http.StatusOK, constant.Success, billsRes))
}

// Bill godoc
// @Summary Update paid bill
// @Description Update paid bill
// @Tags Bills
// @Accept json
// @Produce json
// @Param		billNumber	path		string				true	"billNumber"
// @param Authorization header string true "Authorization"
// @Router /api/bills/{billNumber} [PATCH]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *BillsHandler) UpdatePaid(c *gin.Context) {
	defer pkg.PanicHandler(c)
	billNumber := c.Param("billNumber")
	err := h.service.UpdatePaid(billNumber)
	if err != nil {
		log.Error("Happened error when update paid. Error", err.Error())
		pkg.PanicExeption(constant.UnknownError, "Happened error when update paid"+err.Error())
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccessNoData(http.StatusOK, constant.Success))
}
