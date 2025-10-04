package repository

import (
    "database/sql"
    "backend/domain"
)

type RequestRepository interface {
    Save(req domain.Request) error
    FindAll() ([]domain.Request, error)
}

type requestRepo struct {
    db *sql.DB
}

func NewRequestRepository(db *sql.DB) RequestRepository {
    return &requestRepo{db: db}
}

func (r *requestRepo) Save(req domain.Request) error {
    _, err := r.db.Exec(`
        INSERT INTO requests 
        (request_date, department, division, employment_type, contract_type, reason, requester_name, job_code, job_title, 
         age_from, age_to, gender, nationality, experience, education, special_requirements, created_at)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`,
        req.RequestDate, req.Department, req.Division, req.EmploymentType, req.ContractType,
        req.Reason, req.RequesterName, req.JobCode, req.JobTitle,
        req.AgeFrom, req.AgeTo, req.Gender, req.Nationality,
        req.Experience, req.Education, req.SpecialRequirements,
        req.CreatedAt,
    )
    return err
}

func (r *requestRepo) FindAll() ([]domain.Request, error) {
    rows, err := r.db.Query(`SELECT id, request_date, department, division, employment_type, contract_type, reason,
        requester_name, job_code, job_title, age_from, age_to, gender, nationality, experience, education, special_requirements, created_at 
        FROM requests`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var requests []domain.Request
    for rows.Next() {
        var req domain.Request
        err = rows.Scan(&req.ID, &req.RequestDate, &req.Department, &req.Division,
            &req.EmploymentType, &req.ContractType, &req.Reason, &req.RequesterName,
            &req.JobCode, &req.JobTitle, &req.AgeFrom, &req.AgeTo,
            &req.Gender, &req.Nationality, &req.Experience, &req.Education,
            &req.SpecialRequirements, &req.CreatedAt)
        if err != nil {
            return nil, err
        }
        requests = append(requests, req)
    }
    return requests, nil
}
