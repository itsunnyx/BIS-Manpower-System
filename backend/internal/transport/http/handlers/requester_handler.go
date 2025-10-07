package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RequestHandler struct {
	DB *sql.DB
}

func NewRequestHandler(db *sql.DB) *RequestHandler {
	return &RequestHandler{DB: db}
}

func (h *RequestHandler) GetAllRequests(c *gin.Context) {
	rows, err := h.DB.Query(`SELECT request_id, doc_no, position_title, num_required, overall_status FROM manpower_requests`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
		return
	}
	defer rows.Close()

	var requests []map[string]interface{}
	for rows.Next() {
		var id int
		var docNo, title, status string
		var numRequired int
		rows.Scan(&id, &docNo, &title, &numRequired, &status)

		requests = append(requests, gin.H{
			"id":           id,
			"doc_no":       docNo,
			"title":        title,
			"num_required": numRequired,
			"status":       status,
		})
	}

	c.JSON(http.StatusOK, requests)
}

func (h *RequestHandler) CreateRequest(c *gin.Context) {
	var input struct {
		DocNo         string `json:"doc_no"`
		PositionTitle string `json:"position_title"`
		NumRequired   int    `json:"num_required"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	_, err := h.DB.Exec(
		`INSERT INTO manpower_requests (doc_no, doc_date, position_title, num_required, department_id, requested_by) 
		 VALUES ($1, CURRENT_DATE, $2, $3, 1, 1)`, // department_id, requested_by mock ไว้ก่อน
		input.DocNo, input.PositionTitle, input.NumRequired,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "insert failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "request created"})
}

func (h *RequestHandler) UpdateRequestStatus(c *gin.Context) {
	id := c.Param("id")
	_, err := h.DB.Exec(`UPDATE manpower_requests SET hr_status = 'approved' WHERE request_id = $1`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "request updated"})
}

func (h *RequestHandler) ApproveRequest(c *gin.Context) {
	id := c.Param("id")
	_, err := h.DB.Exec(`UPDATE manpower_requests SET manager_status = 'approved' WHERE request_id = $1`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "approval failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "request approved"})
}

func (h *RequestHandler) RejectRequest(c *gin.Context) {
	id := c.Param("id")
	_, err := h.DB.Exec(`UPDATE manpower_requests SET manager_status = 'approved' WHERE request_id = $1`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "approval failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "request rejected"})
}

func (h *RequestHandler) DeleteRequest(c *gin.Context) {
	id := c.Param("id")
	_, err := h.DB.Exec(`DELETE FROM manpower_requests WHERE request_id = $1`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "request deleted"})
}

func (h *RequestHandler) SyncToRecruitment(c *gin.Context) {
	id := c.Param("id")
	_, err := h.DB.Exec(`INSERT INTO recruitment_sync (request_id, sync_status, response_message) VALUES ($1,'success','mock sync success')`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "sync failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "synced to recruitment system"})
}
