package repository

import (
	"database/sql"
	"manpower/internal/domain"
)

type RequestRepo struct {
	db *sql.DB
}

func NewRequestRepo(db *sql.DB) *RequestRepo {
	return &RequestRepo{db: db}
}

func (r *RequestRepo) Create(req *domain.ManpowerRequest) error {
	query := `
			INSERT INTO manpower_requests 
			(doc_no, department_id, requested_by, position_title, num_required, reason_note)
			VALUES ($1,$2,$3,$4,$5,$6)
			RETURNING request_id`
	return r.db.QueryRow(query,
		req.DocNo, req.DepartmentID, req.RequestedBy,
		req.PositionTitle, req.NumRequired, req.Reason,
	).Scan(&req.ID)
}

func (r *RequestRepo) GetByID(id int) (*domain.ManpowerRequest, error) {
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

func (r *RequestRepo) List() ([]domain.ManpowerRequest, error) {
	rows, err := r.db.Query(`
		SELECT
			request_id,
			doc_no,
			doc_date,
			department_id,
			requested_by,
			position_title,
			position_level,
			num_required,
			reason_note,                 -- Reason
			origin_status::text,         -- OriginStatus
			hr_status::text,             -- HRStatus
			manager_status::text,     -- ManagerStatus
			overall_status::text,        -- OverallStatus
			'' AS remark,                -- ไม่มีคอลัมน์ remark ใน schema -> ให้ค่าว่างไว้ก่อน
			created_at,
			updated_at
		FROM manpower_requests
		ORDER BY doc_date DESC, request_id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []domain.ManpowerRequest
	for rows.Next() {
		var req domain.ManpowerRequest
		if err := rows.Scan(
			&req.ID,
			&req.DocNo,
			&req.DocDate,
			&req.DepartmentID,
			&req.RequestedBy,
			&req.PositionTitle,
			&req.PositionLevel,
			&req.NumRequired,
			&req.Reason,
			&req.OriginStatus,
			&req.HRStatus,
			&req.ManagerStatus,   // = management_status
			&req.OverallStatus,
			&req.Remark,          // ค่าว่างจาก '' AS remark
			&req.CreatedAt,
			&req.UpdatedAt,
		); err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return requests, nil
}

func (r *RequestRepo) ListPendingForManager() ([]domain.ManpowerRequest, error) {
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

func (r *RequestRepo) UpdateHRStatus(id int, status string) error {
	_, err := r.db.Exec(`UPDATE manpower_requests SET hr_status=$1 WHERE request_id=$2`, status, id)
	return err
}

func (r *RequestRepo) UpdateManagerStatus(id int, status string) error {
	_, err := r.db.Exec(`UPDATE manpower_requests SET manager_status=$1 WHERE request_id=$2`, status, id)
	return err
}

func (r *RequestRepo) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM manpower_requests WHERE request_id=$1`, id)
	return err
}

func (r *RequestRepo) SyncToRecruitment(id int) error {
	_, err := r.db.Exec(`INSERT INTO recruitment_sync (request_id, sync_status, response_message) VALUES ($1,'success','mock sync success')`, id)
	return err
}
