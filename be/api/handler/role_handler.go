package handler

import (
	"BE_Manage_device/constant"
	"BE_Manage_device/internal/domain/service"
	"BE_Manage_device/pkg"
	"BE_Manage_device/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	service *service.RoleService
}

func NewRoleHandler(service *service.RoleService) *RoleHandler {
	return &RoleHandler{service: service}
}

// User godoc
// @Summary      GetRole
// @Description  GetRole
// @Tags         Roles
// @Accept       json
// @Produce      json
// @param Authorization header string true "Authorization"
// @Router       /api/roles [GET]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *RoleHandler) GetAllRole(c *gin.Context) {
	defer pkg.PanicHandler(c)
	roles := h.service.GetAllRole()
	rolesResponse := utils.ConvertRolesToResponsesArray(roles)
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(http.StatusOK, constant.Success, rolesResponse))
}
