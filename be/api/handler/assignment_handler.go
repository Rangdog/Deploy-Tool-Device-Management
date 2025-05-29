package handler

import (
	"BE_Manage_device/constant"
	"BE_Manage_device/internal/domain/dto"
	"BE_Manage_device/internal/domain/filter"
	"BE_Manage_device/internal/domain/service"
	"BE_Manage_device/pkg"
	"BE_Manage_device/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type AssignmentHandler struct {
	service *service.AssignmentService
}

func NewAssignmentHandler(service *service.AssignmentService) *AssignmentHandler {
	return &AssignmentHandler{service: service}
}

func (h *AssignmentHandler) Create(c *gin.Context) {
	defer pkg.PanicHandler(c)
	userId := utils.GetUserIdFromContext(c)
	var request dto.AssignmentCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error("Happened error when mapping request from FE. Error", err)
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when mapping request from FE.")
	}
	assignment, err := h.service.Create(request.UserId, request.DepartmentId, userId, request.AssetId)
	if err != nil {
		log.Error("Happened error when create assignment. Error", err)
		pkg.PanicExeption(constant.UnknownError, "Happened error when create assignment.")
	}

	assignResponse := dto.AssignmentResponse{}
	assignResponse.Id = assignment.Id
	assignResponse.UserAssigned.Id = assignment.UserAssigned.Id
	assignResponse.UserAssigned.FirstName = assignment.UserAssigned.FirstName
	assignResponse.UserAssigned.LastName = assignment.UserAssigned.LastName
	assignResponse.UserAssigned.Email = assignment.UserAssigned.Email

	assignResponse.UserAssign.Id = assignment.UserAssign.Id
	assignResponse.UserAssign.FirstName = assignment.UserAssign.FirstName
	assignResponse.UserAssign.LastName = assignment.UserAssign.LastName
	assignResponse.UserAssign.Email = assignment.UserAssign.Email

	assignResponse.Asset.Id = assignment.Asset.Id
	assignResponse.Asset.AssetName = assignment.Asset.AssetName
	assignResponse.Asset.Status = assignment.Asset.Status
	assignResponse.Asset.FileAttachment = *assignment.Asset.FileAttachment
	assignResponse.Asset.ImageUpload = *assignment.Asset.ImageUpload

	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(http.StatusOK, constant.Success, assignResponse))
}

// Assignment godoc
// @Summary      Update assignment
// @Description  Update assignment
// @Tags         Assignments
// @Accept       json
// @Produce      json
// @Param        assignment   body    dto.AssignmentUpdateRequest   true  "Data"
// @Param		id	path		string				true	"id"
// @param Authorization header string true "Authorization"
// @Router       /api/assignments/{id} [PUT]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *AssignmentHandler) Update(c *gin.Context) {
	defer pkg.PanicHandler(c)
	userId := utils.GetUserIdFromContext(c)
	idStr := c.Param("id")
	assignmentId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("Happened error when convert assetId to int64. Error", err)
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when convert assetId to int64")
	}
	var request dto.AssignmentUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error("Happened error when mapping request from FE. Error", err)
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when mapping request from FE.")
	}
	assignmentUpdated, err := h.service.Update(userId, assignmentId, &request.UserId, &request.DepartmentId)
	if err != nil {
		log.Error("Happened error when update assignment. Error", err)
		pkg.PanicExeption(constant.UnknownError, "Happened error when update assignment.")
	}
	assignResponse := dto.AssignmentResponse{}
	assignResponse.Id = assignmentUpdated.Id
	assignResponse.UserAssigned.Id = assignmentUpdated.UserAssigned.Id
	assignResponse.UserAssigned.FirstName = assignmentUpdated.UserAssigned.FirstName
	assignResponse.UserAssigned.LastName = assignmentUpdated.UserAssigned.LastName
	assignResponse.UserAssigned.Email = assignmentUpdated.UserAssigned.Email

	assignResponse.UserAssign.Id = assignmentUpdated.UserAssign.Id
	assignResponse.UserAssign.FirstName = assignmentUpdated.UserAssign.FirstName
	assignResponse.UserAssign.LastName = assignmentUpdated.UserAssign.LastName
	assignResponse.UserAssign.Email = assignmentUpdated.UserAssign.Email

	assignResponse.Asset.Id = assignmentUpdated.Asset.Id
	assignResponse.Asset.AssetName = assignmentUpdated.Asset.AssetName
	assignResponse.Asset.Status = assignmentUpdated.Asset.Status
	assignResponse.Asset.FileAttachment = *assignmentUpdated.Asset.FileAttachment
	assignResponse.Asset.ImageUpload = *assignmentUpdated.Asset.ImageUpload
	assignResponse.Department.ID = *assignmentUpdated.DepartmentId
	assignResponse.Department.DepartmentName = assignmentUpdated.Department.DepartmentName
	assignResponse.Department.Location.ID = assignmentUpdated.Department.LocationId
	assignResponse.Department.Location.LocationName = assignmentUpdated.Department.Location.LocationName
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(http.StatusOK, constant.Success, assignResponse))
}

// Assignment godoc
// @Summary Get all assign with filter
// @Description Get all assign have permission
// @Tags Assignments
// @Accept json
// @Produce json
// @Param        assignment   query    filter.AssignmentFilter   false  "filter assignment"
// @param Authorization header string true "Authorization"
// @Router /api/assignments/filter [GET]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *AssignmentHandler) FilterAsset(c *gin.Context) {
	defer pkg.PanicHandler(c)
	var filter filter.AssignmentFilter
	userId := utils.GetUserIdFromContext(c)
	if err := c.ShouldBindQuery(&filter); err != nil {
		log.Error("Happened error when mapping query to filter. Error", err)
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when mapping query to filter")
	}
	data, err := h.service.Filter(userId, filter.EmailAssigned, filter.EmailAssign, filter.AssetName, filter.Page, filter.Limit)
	if err != nil {
		log.Error("Happened error when filter asset. Error", err)
		pkg.PanicExeption(constant.UnknownError, "Happened error when filter asset")
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(http.StatusOK, constant.Success, data))
}

// Assignment godoc
// @Summary      Get assignment
// @Description  Get assignment
// @Tags         Assignments
// @Accept       json
// @Produce      json
// @Param		id	path		string				true	"id"
// @param Authorization header string true "Authorization"
// @Router       /api/assignments/{id} [GET]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *AssignmentHandler) GetAssignmentById(c *gin.Context) {
	defer pkg.PanicHandler(c)
	userId := utils.GetUserIdFromContext(c)
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("Happened error when convert assignment id to int64. Error", err)
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when convert assignment id to int64")
	}
	assignment, err := h.service.GetAssignmentById(userId, id)
	if err != nil {
		log.Error("Happened error when get assignment. Error", err)
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when get assignment.")
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(http.StatusOK, constant.Success, assignment))
}
