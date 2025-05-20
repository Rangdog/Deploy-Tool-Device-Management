package handler

import (
	"BE_Manage_device/constant"
	"BE_Manage_device/internal/domain/service"
	"BE_Manage_device/pkg"
	"BE_Manage_device/pkg/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type AssetsHandler struct {
	service *service.AssetsService
}

func NewAssetsHandler(service *service.AssetsService) *AssetsHandler {
	return &AssetsHandler{service: service}
}

// Asset godoc
// @Summary Create assets
// @Description Create assets
// @Tags assets
// @Accept multipart/form-data
// @Produce json
// @Param asset_name formData string true "Asset Name"
// @Param purchase_date formData string true "Purchase Date (RFC3339 format, e.g. 2023-04-15T10:00:00Z)"
// @Param cost formData number true "Cost"
// @Param owner formData int64 false "Owner ID"
// @Param warrant_expiry formData string true "Warranty Expiry (RFC3339 format, e.g. 2023-12-31T23:59:59Z)"
// @Param status formData string true "Status"
// @Param serial_number formData string true "Serial Number"
// @Param category_id formData int64 true "Category ID"
// @Param department_id formData int64 false "Department ID"
// @Param file formData file true "File to upload"
// @Param image formData file true "Image to upload"
// @Router /api/assets [post]
func (h *AssetsHandler) Create(c *gin.Context) {
	defer pkg.PanicHandler(c)

	userId := utils.GetUserIdFromContext(c)

	assetName := c.PostForm("asset_name")
	purchaseDateStr := c.PostForm("purchase_date")
	costStr := c.PostForm("cost")
	ownerStr := c.PostForm("owner")
	warrantExpiryStr := c.PostForm("warrant_expiry")
	status := c.PostForm("status")
	serialNumber := c.PostForm("serial_number")
	categoryIdStr := c.PostForm("category_id")
	departmentIdStr := c.PostForm("department_id")

	purchaseDate, err := time.Parse(time.RFC3339, purchaseDateStr)
	if err != nil {
		pkg.PanicExeption(constant.InvalidRequest, "Invalid purchase_date format")
	}

	warrantExpiry, err := time.Parse(time.RFC3339, warrantExpiryStr)
	if err != nil {
		pkg.PanicExeption(constant.InvalidRequest, "Invalid warrant_expiry format")
	}

	cost, err := strconv.ParseFloat(costStr, 64)
	if err != nil {
		pkg.PanicExeption(constant.InvalidRequest, "Invalid cost format")
	}

	var owner *int64
	if ownerStr != "" {
		val, err := strconv.ParseInt(ownerStr, 10, 64)
		if err != nil {
			pkg.PanicExeption(constant.InvalidRequest, "Invalid owner format")
		}
		owner = &val
	}

	categoryId, err := strconv.ParseInt(categoryIdStr, 10, 64)
	if err != nil {
		pkg.PanicExeption(constant.InvalidRequest, "Invalid category_id format")
	}

	var departmentId *int64
	if departmentIdStr != "" {
		val, err := strconv.ParseInt(departmentIdStr, 10, 64)
		if err != nil {
			pkg.PanicExeption(constant.InvalidRequest, "Invalid department_id format")
		}
		departmentId = &val
	}

	file, err := c.FormFile("file")
	if err != nil {
		pkg.PanicExeption(constant.InvalidRequest, "File upload missing")
		return
	}

	image, err := c.FormFile("image")
	if err != nil {
		pkg.PanicExeption(constant.InvalidRequest, "Image upload missing")
		return
	}

	assetCreate, err := h.service.Create(
		userId,
		assetName,
		purchaseDate,
		cost,
		owner,
		warrantExpiry,
		status,
		serialNumber,
		image,
		file,
		categoryId,
		departmentId,
	)

	if err != nil {
		pkg.PanicExeption(constant.InvalidRequest, "Failed to create asset")
	}

	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(http.StatusOK, constant.Success, assetCreate))
}

// Asset godoc
// @Summary Get  assets
// @Description Get  assets
// @Tags assets
// @Accept json
// @Produce json
// @Param		id	path		string				true	"id"
// @Router /api/assets/{id} [GET]
func (h *AssetsHandler) GetAssetById(c *gin.Context) {
	defer pkg.PanicHandler(c)
	userId := utils.GetUserIdFromContext(c)
	idStr := c.Param("id")
	assetId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("Happened error when convert step id to int64. Error", err)
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when convert assetId to int64")
	}
	asset, err := h.service.GetAssetById(userId, assetId)
	if err != nil {
		log.Error("Happened error when get asset by id. Error", err)
		pkg.PanicExeption(constant.UnknownError, "Happened error when get asset by id")
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(http.StatusOK, constant.Success, asset))
}

// Asset godoc
// @Summary Get all assets
// @Description Get all assets
// @Tags assets
// @Accept json
// @Produce json
// @Router /api/assets [GET]
func (h *AssetsHandler) GetAllAsset(c *gin.Context) {
	defer pkg.PanicHandler(c)
	assets, err := h.service.GetAllAsset()
	if err != nil {
		log.Error("Happened error when get all assets. Error", err)
		pkg.PanicExeption(constant.UnknownError, "Happened error when get all assets")
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(http.StatusOK, constant.Success, assets))
}
