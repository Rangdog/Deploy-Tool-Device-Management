package handler

import (
	"BE_Manage_device/constant"
	"BE_Manage_device/internal/domain/dto"
	"BE_Manage_device/internal/domain/service"
	"BE_Manage_device/pkg"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type CategoriesHandler struct {
	service *service.CategoriesService
}

func NewCategoriesHandler(service *service.CategoriesService) *CategoriesHandler {
	return &CategoriesHandler{service: service}
}

// User godoc
// @Summary      Create categories
// @Description  Create categories
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        category   body    dto.CreateCategoryRequest   true  "Data"
// @Router       /api/categories [POST]
func (h *CategoriesHandler) Create(c *gin.Context) {
	defer pkg.PanicHandler(c)
	var request dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error("Happened error when mapping request from FE. Error", err)
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when mapping request from FE.")
	}
	location, err := h.service.Create(request.CategoryName)
	if err != nil {
		log.Error("Happened error when create category. Error", err)
		pkg.PanicExeption(constant.UnknownError)
	}
	c.JSON(http.StatusCreated, pkg.BuildReponseSuccess(http.StatusCreated, constant.Success, location))
}

// User godoc
// @Summary      Get all categories
// @Description  Get all categories
// @Tags         categories
// @Accept       json
// @Produce      json
// @Router       /api/categories [GET]
func (h *CategoriesHandler) GetAll(c *gin.Context) {
	defer pkg.PanicHandler(c)
	categories, err := h.service.GetAll()
	if err != nil {
		log.Error("Happened error when get all categories. Error", err)
		pkg.PanicExeption(constant.UnknownError, "Happened error when get all categories")
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(http.StatusOK, constant.Success, categories))
}

// User godoc
// @Summary      Delete category
// @Description   Delete category via id
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param		id	path		string				true	"id"
// @Router       /api/categories/{id} [DELETE]
func (h *CategoriesHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	IdConvert, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Error("Happened error when get id via path. Error", err)
		pkg.PanicExeption(constant.UnknownError, "Happened error when get id via path")
	}

	err = h.service.Delete(IdConvert)
	if err != nil {
		log.Error("Happened error when delete category. Error", err)
		pkg.PanicExeption(constant.UnknownError, err.Error())
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccessNoData(http.StatusOK, constant.Success))
}
