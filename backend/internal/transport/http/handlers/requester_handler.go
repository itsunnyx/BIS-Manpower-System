package handlers

import (

	svc "manpower/internal/service"

	"database/sql"
	"net/http"
	"github.com/gin-gonic/gin"
)

type RequesterHandler struct {
	svc *svc.RequesterService
}

func NewRequesterHandler(s *svc.RequesterService) *RequesterHandler {
	return &RequesterHandler{svc: s}
}

func (h *RequesterHandler) GetManpowerRequests(c *gin.Context) {
	data, err := h.svc.ListRequests()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func CreateManpowerRequest(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			DocNo    string `json:"doc_no"`
			Title    string `json:"title"`
			Num      int    `json:"num"`
			Reason   string `json:"reason"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := db.Exec(
			`INSERT INTO manpower_requests (doc_no, doc_date, department_id, requested_by, position_title, num_required, reason)
			 VALUES ($1, CURRENT_DATE, 1, 1, $2, $3, $4)`,
			req.DocNo, req.Title, req.Num, req.Reason,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "request created"})
	}
}
