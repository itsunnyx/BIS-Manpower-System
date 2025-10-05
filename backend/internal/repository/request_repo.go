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
    query := `INSERT INTO manpower_requests 
              (doc_no, department_id, requested_by, position_title, num_required, reason)
              VALUES ($1,$2,$3,$4,$5,$6) RETURNING request_id`
    return r.db.QueryRow(query,
        req.DocNo, req.DepartmentID, req.RequestedBy,
        req.PositionTitle, req.NumRequired, req.Reason,
    ).Scan(&req.ID)
}

func (r *RequestRepo) GetByID(id int) (*domain.ManpowerRequest, error) {
    var req domain.ManpowerRequest
    query := `SELECT request_id, doc_no, department_id, requested_by, position_title, num_required, reason 
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
    rows, err := r.db.Query(`SELECT request_id, doc_no, department_id, requested_by, position_title, num_required, reason FROM manpower_requests`)
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
