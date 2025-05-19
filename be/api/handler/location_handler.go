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

type LocationHandler struct {
	service *service.LocationService
}

func NewLocationHandler(service *service.LocationService) *LocationHandler {
	return &LocationHandler{service: service}
}

// User godoc
// @Summary      Create Location
// @Description  Create location
// @Tags         locations
// @Accept       json
// @Produce      json
// @Param        location   body    dto.CreateLocationRequest   true  "Data"
// @Router       /api/locations [POST]
func (h *LocationHandler) Create(c *gin.Context) {
	defer pkg.PanicHandler(c)
	var request dto.CreateLocationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error("Happened error when mapping request from FE. Error", err)
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when mapping request from FE.")
	}
	location, err := h.service.Create(request.LocationName)
	if err != nil {
		log.Error("Happened error when create location. Error", err)
		pkg.PanicExeption(constant.UnknownError)
	}
	c.JSON(http.StatusCreated, pkg.BuildReponseSuccess(http.StatusCreated, constant.Success, location))
}

// User godoc
// @Summary      Get all location
// @Description  Get all location
// @Tags         locations
// @Accept       json
// @Produce      json
// @Router       /api/locations [GET]
func (h *LocationHandler) GetAll(c *gin.Context) {
	defer pkg.PanicHandler(c)
	locations, err := h.service.GetAll()
	if err != nil {
		log.Error("Happened error when get all location. Error", err)
		pkg.PanicExeption(constant.UnknownError, "Happened error when get all location")
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(http.StatusOK, constant.Success, locations))
}

// User godoc
// @Summary      Delete location
// @Description   Delete location via id
// @Tags         locations
// @Accept       json
// @Produce      json
// @Param		id	path		string				true	"id"
// @Router       /api/locations/{id} [DELETE]
func (h *LocationHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	IdConvert, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Error("Happened error when get id via path. Error", err)
		pkg.PanicExeption(constant.UnknownError, "Happened error when get id via path")
	}

	err = h.service.Delete(IdConvert)
	if err != nil {
		log.Error("Happened error when delete location. Error", err)
		pkg.PanicExeption(constant.UnknownError, err.Error())
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccessNoData(http.StatusOK, constant.Success))
}
