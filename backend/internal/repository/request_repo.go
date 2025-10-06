package repository

import (
	"database/sql"
	"manpower/internal/domain"
)

// 1) ประกาศ interface (สัญญา)
type RequestRepo interface {
	Create(req *domain.ManpowerRequest) error
	GetByID(id int) (*domain.ManpowerRequest, error)
	List() ([]domain.ManpowerRequest, error)
	ListPendingForManager() ([]domain.ManpowerRequest, error)
}

// 2) ทำ struct (implementation) ให้ชื่อไม่ชนกับ interface
type requestRepo struct {
	db *sql.DB
}

// 3) Constructor คืนค่าเป็น "interface" ไม่ใช่ struct
func NewRequestRepo(db *sql.DB) RequestRepo {
	return &requestRepo{db: db}
}

// 4) เมธอดทั้งหมด ต้องผูกกับ *requestRepo (ตัวเล็ก)
func (r *requestRepo) Create(req *domain.ManpowerRequest) error {
	query := `
		INSERT INTO manpower_requests 
		(doc_no, department_id, requested_by, position_title, num_required, reason)
		VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING request_id`
	return r.db.QueryRow(query,
		req.DocNo, req.DepartmentID, req.RequestedBy,
		req.PositionTitle, req.NumRequired, req.Reason,
	).Scan(&req.ID)
}

func (r *requestRepo) GetByID(id int) (*domain.ManpowerRequest, error) {
	var req domain.ManpowerRequest
	query := `
		SELECT request_id, doc_no, department_id, requested_by, position_title, num_required, reason 
		FROM manpower_requests WHERE request_id=$1`
	err := r.db.QueryRow(query, id).Scan(
		&req.ID, &req.DocNo, &req.DepartmentID, &req.RequestedBy,
		&req.PositionTitle, &req.NumRequired, &req.Reason,
	)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *requestRepo) List() ([]domain.ManpowerRequest, error) {
	rows, err := r.db.Query(`
		SELECT request_id, doc_no, department_id, requested_by, position_title, num_required, reason
		FROM manpower_requests`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []domain.ManpowerRequest
	for rows.Next() {
		var req domain.ManpowerRequest
		if err := rows.Scan(
			&req.ID, &req.DocNo, &req.DepartmentID, &req.RequestedBy,
			&req.PositionTitle, &req.NumRequired, &req.Reason,
		); err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}
	return requests, nil
}

func (r *requestRepo) ListPendingForManager() ([]domain.ManpowerRequest, error) {
	rows, err := r.db.Query(`
		SELECT request_id, doc_no, position_title, num_required, overall_status
		FROM manpower_requests
		WHERE manager_status = 'pending'
		ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.ManpowerRequest
	for rows.Next() {
		var req domain.ManpowerRequest
		if err := rows.Scan(&req.ID, &req.DocNo, &req.PositionTitle, &req.NumRequired, &req.OverallStatus); err != nil {
			return nil, err
		}
		result = append(result, req)
	}
	return result, nil
}
