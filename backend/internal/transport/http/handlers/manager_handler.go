package handlers

import (
	svc "manpower/internal/service"
	"net/http"
	"github.com/gin-gonic/gin"
)

type ManagerHandler struct {
	svc *svc.ManagerService
}

func NewManagerHandler(s *svc.ManagerService) *ManagerHandler {
	return &ManagerHandler{svc: s}
}

func (h *ManagerHandler) ListPendingRequests(c *gin.Context) {
	data, err := h.svc.ListPendingRequests()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(data) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "No pending requests found",
			"data":    []any{},
		})
		return
	}

	c.JSON(http.StatusOK, data)
}
