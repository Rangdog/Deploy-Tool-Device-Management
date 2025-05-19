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

type DepartmentsHandler struct {
	service *service.DepartmentsService
}

func NewDepartmentsHandler(service *service.DepartmentsService) *DepartmentsHandler {
	return &DepartmentsHandler{service: service}
}

// User godoc
// @Summary      Create departments
// @Description  Create departments
// @Tags         departments
// @Accept       json
// @Produce      json
// @Param        department   body    dto.CreateDepartmentRequest   true  "Data"
// @Router       /api/departments [POST]
func (h *DepartmentsHandler) Create(c *gin.Context) {
	defer pkg.PanicHandler(c)
	var request dto.CreateDepartmentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error("Happened error when mapping request from FE. Error", err)
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when mapping request from FE.")
	}
	department, err := h.service.Create(request.DepartmentName, request.LocationId)
	if err != nil {
		log.Error("Happened error when create department. Error", err)
		pkg.PanicExeption(constant.UnknownError)
	}
	c.JSON(http.StatusCreated, pkg.BuildReponseSuccess(http.StatusCreated, constant.Success, department))
}

// User godoc
// @Summary      Get all departments
// @Description  Get all departments
// @Tags         departments
// @Accept       json
// @Produce      json
// @Router       /api/departments [GET]
func (h *DepartmentsHandler) GetAll(c *gin.Context) {
	defer pkg.PanicHandler(c)
	departments, err := h.service.GetAll()
	if err != nil {
		log.Error("Happened error when get all departments. Error", err)
		pkg.PanicExeption(constant.UnknownError, "Happened error when get all departments")
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(http.StatusOK, constant.Success, departments))
}

// User godoc
// @Summary      Delete department
// @Description   Delete department via id
// @Tags         departments
// @Accept       json
// @Produce      json
// @Param		id	path		string				true	"id"
// @Router       /api/departments/{id} [DELETE]
func (h *DepartmentsHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	IdConvert, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Error("Happened error when get id via path. Error", err)
		pkg.PanicExeption(constant.UnknownError, "Happened error when get id via path")
	}

	err = h.service.Delete(IdConvert)
	if err != nil {
		log.Error("Happened error when delete department. Error", err)
		pkg.PanicExeption(constant.UnknownError, err.Error())
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccessNoData(http.StatusOK, constant.Success))
}
