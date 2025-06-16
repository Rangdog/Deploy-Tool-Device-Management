package handler

import (
	"BE_Manage_device/constant"
	"BE_Manage_device/internal/domain/dto"
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/filter"
	"BE_Manage_device/internal/domain/service"
	"BE_Manage_device/pkg"
	"BE_Manage_device/pkg/utils"
	"bytes"
	"encoding/csv"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/phpdave11/gofpdf"
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
// @Tags Assets
// @Accept multipart/form-data
// @Produce json
// @Param assetName formData string true "Asset Name"
// @Param purchaseDate formData string true "Purchase Date (RFC3339 format, e.g. 2023-04-15T10:00:00Z)"
// @Param cost formData number true "Cost"
// @Param warrantExpiry formData string true "Warranty Expiry (RFC3339 format, e.g. 2023-12-31T23:59:59Z)"
// @Param serialNumber formData string true "Serial Number"
// @Param categoryId formData int64 true "Category ID"
// @Param departmentId formData int64 true "Department ID"
// @Param redirectUrl formData string true "redirect url"
// @Param file formData file true "File to upload"
// @Param image formData file true "Image to upload"
// @param Authorization header string true "Authorization"
// @Router /api/assets [post]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *AssetsHandler) Create(c *gin.Context) {
	defer pkg.PanicHandler(c)

	userId := utils.GetUserIdFromContext(c)

	assetName := c.PostForm("assetName")
	purchaseDateStr := c.PostForm("purchaseDate")
	costStr := c.PostForm("cost")
	warrantExpiryStr := c.PostForm("warrantExpiry")
	serialNumber := c.PostForm("serialNumber")
	categoryIdStr := c.PostForm("categoryId")
	departmentIdStr := c.PostForm("departmentId")
	url := c.PostForm("redirectUrl")

	purchaseDate, err := time.Parse(time.RFC3339, purchaseDateStr)

	if err != nil {
		log.Info("Error: ", err.Error())
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
	categoryId, err := strconv.ParseInt(categoryIdStr, 10, 64)
	if err != nil {
		pkg.PanicExeption(constant.InvalidRequest, "Invalid category_id format")
	}

	var departmentId int64
	if departmentIdStr != "" {
		val, err := strconv.ParseInt(departmentIdStr, 10, 64)
		if err != nil {
			pkg.PanicExeption(constant.InvalidRequest, "Invalid department_id format")
		}
		departmentId = val
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

	err = h.service.CheckPermissionForManager(userId, departmentId)
	if err != nil {
		pkg.PanicExeption(constant.InvalidRequest, err.Error())
		return
	}

	assetCreate, err := h.service.Create(
		userId,
		assetName,
		purchaseDate,
		warrantExpiry,
		serialNumber,
		image,
		file,
		categoryId,
		departmentId,
		url,
		cost,
	)

	if err != nil {
		pkg.PanicExeption(constant.InvalidRequest, "Failed to create asset")
	}
	asset, err := h.service.GetAssetById(userId, assetCreate.Id)
	if err != nil {
		log.Error("Happened error when get asset by id. Error", err)
		pkg.PanicExeption(constant.UnknownError, "Happened error when get asset by id")
	}
	assetResponse := dto.AssetResponse{
		ID:             asset.Id,
		AssetName:      asset.AssetName,
		PurchaseDate:   asset.PurchaseDate.Format("2006-01-02"),
		Cost:           asset.Cost,
		WarrantExpiry:  asset.WarrantExpiry.Format("2006-01-02"),
		Status:         asset.Status,
		SerialNumber:   asset.SerialNumber,
		FileAttachment: *asset.FileAttachment,
		ImageUpload:    *asset.ImageUpload,
		QrURL:          *asset.QrUrl,
		Category: dto.CategoryResponse{
			ID:           asset.Category.Id,
			CategoryName: asset.Category.CategoryName,
		},
		Department: dto.DepartmentResponse{
			ID:             asset.Department.Id,
			DepartmentName: asset.Department.DepartmentName,
			Location: dto.LocationResponse{
				ID:           asset.Department.Location.Id,
				LocationName: asset.Department.Location.LocationName,
			},
		},
	}
	if asset.OnwerUser != nil {
		assetResponse.Owner = dto.OwnerResponse{
			ID:        asset.OnwerUser.Id,
			FirstName: asset.OnwerUser.FirstName,
			LastName:  asset.OnwerUser.LastName,
			Email:     asset.OnwerUser.Email,
		}
	}
	c.JSON(http.StatusCreated, pkg.BuildReponseSuccess(http.StatusCreated, constant.Success, assetResponse))
}

// Asset godoc
// @Summary Update assets
// @Description Update assets
// @Tags Assets
// @Accept multipart/form-data
// @Produce json
// @Param		id	path		string				true	"id"
// @Param assetName formData string true "Asset Name"
// @Param purchaseDate formData string true "Purchase Date (RFC3339 format, e.g. 2023-04-15T10:00:00Z)"
// @Param cost formData number true "Cost"
// @Param warrantExpiry formData string true "Warranty Expiry (RFC3339 format, e.g. 2023-12-31T23:59:59Z)"
// @Param serialNumber formData string true "Serial Number"
// @Param categoryId formData int64 true "Category ID"
// @Param file formData file true "File to upload"
// @Param image formData file true "Image to upload"
// @param Authorization header string true "Authorization"
// @Router /api/assets/{id} [PUT]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *AssetsHandler) Update(c *gin.Context) {
	defer pkg.PanicHandler(c)
	idStr := c.Param("id")
	assetId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("Happened error when convert asset id to int64. Error", err)
		pkg.PanicExeption(constant.UnknownError)
	}

	userId := utils.GetUserIdFromContext(c)

	assetName := c.PostForm("assetName")
	purchaseDateStr := c.PostForm("purchaseDate")
	costStr := c.PostForm("cost")
	warrantExpiryStr := c.PostForm("warrantExpiry")
	serialNumber := c.PostForm("serialNumber")
	categoryIdStr := c.PostForm("categoryId")

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

	categoryId, err := strconv.ParseInt(categoryIdStr, 10, 64)
	if err != nil {
		pkg.PanicExeption(constant.InvalidRequest, "Invalid category_id format")
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
	assetCheck, err := h.service.GetAssetById(userId, assetId)
	if err != nil {
		pkg.PanicExeption(constant.InvalidRequest, "Failed to update asset")
	}
	err = h.service.CheckPermissionForManager(userId, assetCheck.DepartmentId)
	if err != nil {
		pkg.PanicExeption(constant.InvalidRequest, err.Error())
		return
	}
	assetUpdate, err := h.service.UpdateAsset(
		userId,
		assetId,
		assetName,
		purchaseDate,
		warrantExpiry,
		serialNumber,
		image,
		file,
		categoryId,
		cost,
	)
	if err != nil {
		pkg.PanicExeption(constant.InvalidRequest, "Failed to update asset")
	}
	asset, err := h.service.GetAssetById(userId, assetUpdate.Id)
	if err != nil {
		log.Error("Happened error when get asset by id. Error", err)
		pkg.PanicExeption(constant.UnknownError, "Happened error when get asset by id")
	}
	assetResponse := dto.AssetResponse{
		ID:             asset.Id,
		AssetName:      asset.AssetName,
		PurchaseDate:   asset.PurchaseDate.Format("2006-01-02"),
		Cost:           asset.Cost,
		WarrantExpiry:  asset.WarrantExpiry.Format("2006-01-02"),
		Status:         asset.Status,
		SerialNumber:   asset.SerialNumber,
		FileAttachment: *asset.FileAttachment,
		ImageUpload:    *asset.ImageUpload,
		QrURL:          *asset.QrUrl,
		Category: dto.CategoryResponse{
			ID:           asset.Category.Id,
			CategoryName: asset.Category.CategoryName,
		},
		Department: dto.DepartmentResponse{
			ID:             asset.Department.Id,
			DepartmentName: asset.Department.DepartmentName,
			Location: dto.LocationResponse{
				ID:           asset.Department.Location.Id,
				LocationName: asset.Department.Location.LocationName,
			},
		},
	}
	if asset.OnwerUser != nil {
		assetResponse.Owner = dto.OwnerResponse{
			ID:        asset.OnwerUser.Id,
			FirstName: asset.OnwerUser.FirstName,
			LastName:  asset.OnwerUser.LastName,
			Email:     asset.OnwerUser.Email,
		}
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(http.StatusOK, constant.Success, assetResponse))
}

// Asset godoc
// @Summary Get assets
// @Description Get assets
// @Tags Assets
// @Accept json
// @Produce json
// @Param		id	path		string				true	"id"
// @param Authorization header string true "Authorization"
// @Router /api/assets/{id} [GET]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *AssetsHandler) GetAssetById(c *gin.Context) {
	defer pkg.PanicHandler(c)
	userId := utils.GetUserIdFromContext(c)
	idStr := c.Param("id")
	assetId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("Happened error when convert assetId to int64. Error", err)
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when convert assetId to int64")
	}
	asset, err := h.service.GetAssetById(userId, assetId)
	if err != nil {
		log.Error("Happened error when get asset by id. Error", err)
		pkg.PanicExeption(constant.UnknownError, "Happened error when get asset by id")
	}
	assetResponse := dto.AssetResponse{
		ID:             asset.Id,
		AssetName:      asset.AssetName,
		PurchaseDate:   asset.PurchaseDate.Format("2006-01-02"),
		Cost:           asset.Cost,
		WarrantExpiry:  asset.WarrantExpiry.Format("2006-01-02"),
		Status:         asset.Status,
		SerialNumber:   asset.SerialNumber,
		FileAttachment: *asset.FileAttachment,
		ImageUpload:    *asset.ImageUpload,
		QrURL:          *asset.QrUrl,
		Category: dto.CategoryResponse{
			ID:           asset.Category.Id,
			CategoryName: asset.Category.CategoryName,
		},
		Department: dto.DepartmentResponse{
			ID:             asset.Department.Id,
			DepartmentName: asset.Department.DepartmentName,
			Location: dto.LocationResponse{
				ID:           asset.Department.Location.Id,
				LocationName: asset.Department.Location.LocationName,
			},
		},
	}
	if asset.OnwerUser != nil {
		assetResponse.Owner = dto.OwnerResponse{
			ID:        asset.OnwerUser.Id,
			FirstName: asset.OnwerUser.FirstName,
			LastName:  asset.OnwerUser.LastName,
			Email:     asset.OnwerUser.Email,
		}
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(http.StatusOK, constant.Success, assetResponse))
}

// Asset godoc
// @Summary Get all assets
// @Description Get all assets
// @Tags Assets
// @Accept json
// @Produce json
// @param Authorization header string true "Authorization"
// @Router /api/assets [GET]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *AssetsHandler) GetAllAsset(c *gin.Context) {
	defer pkg.PanicHandler(c)
	assets, err := h.service.GetAllAsset()
	if err != nil {
		log.Error("Happened error when get all assets. Error", err)
		pkg.PanicExeption(constant.UnknownError, "Happened error when get all assets")
	}
	assetsResponse := []dto.AssetResponse{}
	for _, asset := range assets {
		assetResponse := dto.AssetResponse{
			ID:             asset.Id,
			AssetName:      asset.AssetName,
			PurchaseDate:   asset.PurchaseDate.Format("2006-01-02"),
			Cost:           asset.Cost,
			WarrantExpiry:  asset.WarrantExpiry.Format("2006-01-02"),
			Status:         asset.Status,
			SerialNumber:   asset.SerialNumber,
			FileAttachment: *asset.FileAttachment,
			ImageUpload:    *asset.ImageUpload,
			QrURL:          *asset.QrUrl,
			Category: dto.CategoryResponse{
				ID:           asset.Category.Id,
				CategoryName: asset.Category.CategoryName,
			},
			Department: dto.DepartmentResponse{
				ID:             asset.Department.Id,
				DepartmentName: asset.Department.DepartmentName,
				Location: dto.LocationResponse{
					ID:           asset.Department.Location.Id,
					LocationName: asset.Department.Location.LocationName,
				},
			},
		}
		if asset.OnwerUser != nil {
			assetResponse.Owner = dto.OwnerResponse{
				ID:        asset.OnwerUser.Id,
				FirstName: asset.OnwerUser.FirstName,
				LastName:  asset.OnwerUser.LastName,
				Email:     asset.OnwerUser.Email,
			}
		}
		assetsResponse = append(assetsResponse, assetResponse)
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(http.StatusOK, constant.Success, assetsResponse))
}

// Asset godoc
// @Summary Delete assets
// @Description Delete assets
// @Tags Assets
// @Accept json
// @Produce json
// @Param		id	path		string				true	"id"
// @param Authorization header string true "Authorization"
// @Router /api/assets/{id} [Delete]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *AssetsHandler) DeleteAsset(c *gin.Context) {
	defer pkg.PanicHandler(c)
	userId := utils.GetUserIdFromContext(c)
	idStr := c.Param("id")
	assetId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("Happened error when convert assetId to int64. Error", err)
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when convert assetId to int64")
	}
	assetCheck, err := h.service.GetAssetById(userId, assetId)
	if err != nil {
		pkg.PanicExeption(constant.InvalidRequest, "Failed to update asset")
	}
	err = h.service.CheckPermissionForManager(userId, assetCheck.DepartmentId)
	if err != nil {
		pkg.PanicExeption(constant.InvalidRequest, err.Error())
		return
	}
	err = h.service.DeleteAsset(userId, assetId)
	if err != nil {
		pkg.PanicExeption(constant.UnknownError, "Happened error when delete assets")
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccessNoData(http.StatusOK, constant.Success))
}

// Asset godoc
// @Summary Retired assets
// @Description Retired assets
// @Tags Assets
// @Accept json
// @Produce json
// @Param        ResidualValue   body    dto.RetiredAssetRequest   true  "Data"
// @Param		id	path		string				true	"id"
// @Router /api/assets-retired/{id} [PATCH]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *AssetsHandler) UpdateAssetRetired(c *gin.Context) {
	defer pkg.PanicHandler(c)
	userId := utils.GetUserIdFromContext(c)
	idStr := c.Param("id")
	assetId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("Happened error when convert assetId to int64. Error", err)
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when convert assetId to int64")
	}
	var request dto.RetiredAssetRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error("Happened error when mapping request from FE. Error", err)
		pkg.PanicExeption(constant.InvalidRequest)
	}
	assetCheck, err := h.service.GetAssetById(userId, assetId)
	if err != nil {
		pkg.PanicExeption(constant.InvalidRequest, "Failed to update asset")
	}
	err = h.service.CheckPermissionForManager(userId, assetCheck.DepartmentId)
	if err != nil {
		pkg.PanicExeption(constant.InvalidRequest, err.Error())
		return
	}
	asset, err := h.service.UpdateAssetRetired(userId, assetId, request.ResidualValue)
	if err != nil {
		pkg.PanicExeption(constant.UnknownError, "Happened error when retired assets")
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(http.StatusOK, constant.Success, asset))
}

// Asset godoc
// @Summary Get all assets with filter
// @Description Get all assets have permission
// @Tags Assets
// @Accept json
// @Produce json
// @Param        asset   query    filter.AssetFilter   false  "filter asset"
// @param Authorization header string true "Authorization"
// @Router /api/assets/filter [GET]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *AssetsHandler) FilterAsset(c *gin.Context) {
	defer pkg.PanicHandler(c)
	var filter filter.AssetFilter
	userId := utils.GetUserIdFromContext(c)
	if err := c.ShouldBindQuery(&filter); err != nil {
		log.Error("Happened error when mapping query to filter. Error", err)
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when mapping query to filter")
	}
	data, err := h.service.Filter(userId, filter.AssetName, filter.Status, filter.CategoryId, filter.Cost, filter.SerialNumber, filter.Email, filter.DepartmentId)
	if err != nil {
		log.Error("Happened error when filter asset. Error", err)
		pkg.PanicExeption(constant.UnknownError, "Happened error when filter asset")
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(http.StatusOK, constant.Success, data))
}

// Asset godoc
// @Summary Get dashboard
// @Description Get dashboard
// @Tags Assets
// @Accept json
// @Produce json
// @Param        asset   query    filter.AssetFilterDashboard   false  "filter asset"
// @param Authorization header string true "Authorization"
// @Router /api/assets/filter-dashboard [GET]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *AssetsHandler) FilterAssetDashboard(c *gin.Context) {
	defer pkg.PanicHandler(c)
	var filter filter.AssetFilterDashboard
	userId := utils.GetUserIdFromContext(c)
	if err := c.ShouldBindQuery(&filter); err != nil {
		log.Error("Happened error when mapping query to filter. Error", err)
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when mapping query to filter")
	}
	summary, assets, err := h.service.ApplyFilterDashBoard(userId, filter.CategoryId, filter.DepartmentId, filter.Status, filter.Export)
	if err != nil {
		log.Error("Happened error when filter asset. Error", err)
		pkg.PanicExeption(constant.UnknownError, "Happened error when filter asset")
	}
	if filter.Export != nil {
		user, err := h.service.GetUserById(userId)
		if err != nil {
			log.Error("Happened error when filter asset. Error", err)
			pkg.PanicExeption(constant.UnknownError, "Happened error when filter asset")
		}
		if user.Role.Slug != "viewer" {
			if *filter.Export == "csv" {
				data, _ := GenerateCSV(assets)
				c.Header("Content-Disposition", "attachment; filename=assets.csv")
				c.Data(http.StatusOK, "text/csv", data)
				return
			} else if *filter.Export == "pdf" {
				data, _ := GeneratePDF(assets)
				c.Header("Content-Disposition", "attachment; filename=assets.pdf")
				c.Data(http.StatusOK, "application/pdf", data)
				return
			}
		} else {
			if user.CanExport {
				if *filter.Export == "csv" {
					data, _ := GenerateCSV(assets)
					c.Header("Content-Disposition", "attachment; filename=assets.csv")
					c.Data(http.StatusOK, "text/csv", data)
					return
				} else if *filter.Export == "pdf" {
					data, _ := GeneratePDF(assets)
					c.Header("Content-Disposition", "attachment; filename=assets.pdf")
					c.Data(http.StatusOK, "application/pdf", data)
					return
				}
			}
		}
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(http.StatusOK, constant.Success, summary))
}

func GenerateCSV(assets []*entity.Assets) ([]byte, error) {
	var b bytes.Buffer
	writer := csv.NewWriter(&b)
	writer.Write([]string{"ID", "Name", "Category", "Department", "Status"})
	for _, a := range assets {
		writer.Write([]string{
			strconv.FormatInt(a.Id, 10),
			a.AssetName, a.Category.CategoryName, a.Department.DepartmentName, a.Status,
		})
	}
	writer.Flush()
	return b.Bytes(), writer.Error()
}

func GeneratePDF(assets []*entity.Assets) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, "Asset Dashboard Report")
	pdf.Ln(12)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(10, 10, "ID")
	pdf.Cell(40, 10, "Name")
	pdf.Cell(30, 10, "Category")
	pdf.Cell(30, 10, "Department")
	pdf.Cell(30, 10, "Status")
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 10)

	for _, a := range assets {
		pdf.Cell(10, 10, strconv.FormatInt(a.Id, 10))
		pdf.Cell(40, 10, a.AssetName)
		pdf.Cell(30, 10, a.Category.CategoryName)
		pdf.Cell(30, 10, a.Department.DepartmentName)
		pdf.Cell(30, 10, a.Status)
		pdf.Ln(8)
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Asset godoc
// @Summary Get asset by category of department
// @Description Get asset by category of department
// @Tags Assets
// @Accept json
// @Produce json
// @Param        asset   query   dto.GetAssetsByCateOfDepartmentRequest   false  "category id and department id"
// @param Authorization header string true "Authorization"
// @Router /api/assets/request-transfer [GET]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *AssetsHandler) GetAssetsByCateOfDepartment(c *gin.Context) {
	defer pkg.PanicHandler(c)
	userId := utils.GetUserIdFromContext(c)
	var request dto.GetAssetsByCateOfDepartmentRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		log.Error("Happened error when mapping request from FE. Error", err)
		pkg.PanicExeption(constant.InvalidRequest)
	}
	assets, err := h.service.GetAssetsByCateOfDepartment(userId, request.CategoryId, request.DepartmentId)
	if err != nil {
		log.Error("Happened error when get assets. Error", err)
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when get assets")
	}
	assetsResponse := []dto.AssetResponse{}
	for _, asset := range assets {
		assetResponse := dto.AssetResponse{
			ID:             asset.Id,
			AssetName:      asset.AssetName,
			PurchaseDate:   asset.PurchaseDate.Format("2006-01-02"),
			Cost:           asset.Cost,
			WarrantExpiry:  asset.WarrantExpiry.Format("2006-01-02"),
			Status:         asset.Status,
			SerialNumber:   asset.SerialNumber,
			FileAttachment: *asset.FileAttachment,
			ImageUpload:    *asset.ImageUpload,
			QrURL:          *asset.QrUrl,
			Category: dto.CategoryResponse{
				ID:           asset.Category.Id,
				CategoryName: asset.Category.CategoryName,
			},
			Department: dto.DepartmentResponse{
				ID:             asset.Department.Id,
				DepartmentName: asset.Department.DepartmentName,
				Location: dto.LocationResponse{
					ID:           asset.Department.Location.Id,
					LocationName: asset.Department.Location.LocationName,
				},
			},
		}
		if asset.OnwerUser != nil {
			assetResponse.Owner = dto.OwnerResponse{
				ID:        asset.OnwerUser.Id,
				FirstName: asset.OnwerUser.FirstName,
				LastName:  asset.OnwerUser.LastName,
				Email:     asset.OnwerUser.Email,
			}
		}
		assetsResponse = append(assetsResponse, assetResponse)
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(http.StatusOK, constant.Success, assetsResponse))
}

// Asset godoc
// @Summary Get asset haven't maintenance schedules
// @Description Get asset haven't maintenance schedules
// @Tags Assets
// @Accept json
// @Produce json
// @param Authorization header string true "Authorization"
// @Router /api/assets/maintenance-schedules [GET]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *AssetsHandler) GetAllAssetNotHaveMaintenance(c *gin.Context) {
	defer pkg.PanicHandler(c)
	assets, err := h.service.GetAllAssetNotHaveMaintenance()
	if err != nil {
		log.Error("Happened error when get assets. Error", err)
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when get assets")
	}
	assetsResponse := []dto.AssetResponse{}
	for _, asset := range assets {
		assetResponse := dto.AssetResponse{
			ID:             asset.Id,
			AssetName:      asset.AssetName,
			PurchaseDate:   asset.PurchaseDate.Format("2006-01-02"),
			Cost:           asset.Cost,
			WarrantExpiry:  asset.WarrantExpiry.Format("2006-01-02"),
			Status:         asset.Status,
			SerialNumber:   asset.SerialNumber,
			FileAttachment: *asset.FileAttachment,
			ImageUpload:    *asset.ImageUpload,
			QrURL:          *asset.QrUrl,
			Category: dto.CategoryResponse{
				ID:           asset.Category.Id,
				CategoryName: asset.Category.CategoryName,
			},
			Department: dto.DepartmentResponse{
				ID:             asset.Department.Id,
				DepartmentName: asset.Department.DepartmentName,
				Location: dto.LocationResponse{
					ID:           asset.Department.Location.Id,
					LocationName: asset.Department.Location.LocationName,
				},
			},
		}
		if asset.OnwerUser != nil {
			assetResponse.Owner = dto.OwnerResponse{
				ID:        asset.OnwerUser.Id,
				FirstName: asset.OnwerUser.FirstName,
				LastName:  asset.OnwerUser.LastName,
				Email:     asset.OnwerUser.Email,
			}
		}
		assetsResponse = append(assetsResponse, assetResponse)
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(http.StatusOK, constant.Success, assetsResponse))
}
