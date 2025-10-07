package handlers

import (
	// "database/sql"
	"manpower/internal/domain"
	svc "manpower/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RequestHandler struct {
	svc *svc.RequestService
}

func NewRequestHandler(s *svc.RequestService) *RequestHandler {
	return &RequestHandler{svc: s}
}

func (h *RequestHandler) GetManpowerRequests(c *gin.Context) {
	data, err := h.svc.ListRequests()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// สร้างใหม่
func (h *RequestHandler) CreateRequest(c *gin.Context) {
	var input struct {
		DocNo         string `json:"doc_no"`
		DepartmentID  int    `json:"department_id"`
		RequestedBy   int    `json:"requested_by"`
		PositionTitle string `json:"position_title"`
		NumRequired   int    `json:"num_required"`
		Reason        string `json:"reason"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	req := &domain.ManpowerRequest{
		DocNo:         input.DocNo,
		DepartmentID:  input.DepartmentID,
		RequestedBy:   input.RequestedBy,
		PositionTitle: input.PositionTitle,
		NumRequired:   input.NumRequired,
		Reason:        input.Reason,
	}

	if err := h.svc.CreateRequest(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "insert failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "request created", "id": req.ID})
}

// อัพเดทสถานะ HR
func (h *RequestHandler) UpdateRequestStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.svc.UpdateHRStatus(id, "approved"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "request updated"})
}

// approve
func (h *RequestHandler) ApproveRequest(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.svc.UpdateManagerStatus(id, "approved"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "approval failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "request approved"})
}

// ลบ
func (h *RequestHandler) DeleteRequest(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.svc.DeleteRequest(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "request deleted"})
}
