package handler

import (
	"BE_Manage_device/internal/domain/service"
	"BE_Manage_device/pkg/utils"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type SSEHandler struct {
	service *service.NotificationService
}

func NewSSEHandler(service *service.NotificationService) *SSEHandler {
	return &SSEHandler{service: service}
}

// User godoc
// @Summary      GetRole
// @Description  GetRole
// @Tags         Sockets
// @Accept       json
// @Produce      json
// @param Authorization header string true "Authorization"
// @Router       /api/sse [GET]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *SSEHandler) SSEHandle(c *gin.Context) {
	userId := utils.GetUserIdFromContext(c)
	userIdStr := strconv.Itoa(int(userId))
	msgChan := h.service.Register(userIdStr)
	defer h.service.UnRegister(userIdStr, msgChan)
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Flush()
	// Giữ connection tới khi client disconnect
	notify := c.Writer.CloseNotify()
	for {
		select {
		case msg := <-msgChan:
			fmt.Fprintf(c.Writer, "data: %s\n\n", msg)
			c.Writer.Flush()
		case <-notify:
			return // Client disconnect
		case <-time.After(30 * time.Second):
			// Gửi keep-alive ping (tránh timeout)
			fmt.Fprintf(c.Writer, ":\n\n")
			c.Writer.Flush()
		}
	}
}

func (h *SSEHandler) SendNotificationHandler(c *gin.Context) {
	userId := c.Param("userId")
	var req struct{ Message string }
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid"})
		return
	}
	h.service.Push(userId, req.Message)
	c.JSON(http.StatusOK, gin.H{"status": "sent"})
}
