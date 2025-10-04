package handlers

import (
	"backend/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type ManpowerHandler struct {
	DB *sql.DB
}

// POST /manpower_requests
func (h *ManpowerHandler) CreateRequest(w http.ResponseWriter, r *http.Request) {
	var req models.ManpowerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// กำหนด doc_no (จริง ๆ ควรมี generator เช่น PQ24110013)
	if req.DocNo == "" {
		req.DocNo = "AUTO-" + time.Now().Format("20060102150405")
	}

	query := `
        INSERT INTO manpower_requests 
        (doc_no, doc_date, due_date, department_id, requested_by, 
         position_title, position_level, num_required, reason, origin_status)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
        RETURNING request_id, created_at, updated_at
    `
	err := h.DB.QueryRow(
		query,
		req.DocNo, req.DocDate, req.DueDate,
		req.DepartmentID, req.RequestedBy,
		req.PositionTitle, req.PositionLevel, req.NumRequired, req.Reason,
		"submitted",
	).Scan(&req.RequestID, &req.CreatedAt, &req.UpdatedAt)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(req)
}
