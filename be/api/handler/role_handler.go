package handler

import (
	"BE_Manage_device/constant"
	"BE_Manage_device/internal/domain/service"
	"BE_Manage_device/pkg"
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
// @Tags         roles
// @Accept       json
// @Produce      json
// @Router       /api/roles [GET]
func (h *RoleHandler) GetAllRole(c *gin.Context) {
	defer pkg.PanicHandler(c)
	roles := h.service.GetAllRole()
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(http.StatusOK, constant.Success, roles))
}
