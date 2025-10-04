package models

import "time"

type ManpowerRequest struct {
    RequestID       int       `db:"request_id" json:"request_id"`
    DocNo           string    `db:"doc_no" json:"doc_no"`
    DocDate         time.Time `db:"doc_date" json:"doc_date"`
    DueDate         *time.Time `db:"due_date" json:"due_date,omitempty"`

    DepartmentID    int       `db:"department_id" json:"department_id"`
    RequestedBy     int       `db:"requested_by" json:"requested_by"`

    PositionTitle   string    `db:"position_title" json:"position_title"`
    PositionLevel   string    `db:"position_level" json:"position_level"`
    NumRequired     int       `db:"num_required" json:"num_required"`
    Reason          string    `db:"reason" json:"reason"`

    OriginStatus    string    `db:"origin_status" json:"origin_status"`
    ManagerStatus   string    `db:"manager_status" json:"manager_status"`
    HrStatus        string    `db:"hr_status" json:"hr_status"`
    ManagementStatus string   `db:"management_status" json:"management_status"`
    OverallStatus   string    `db:"overall_status" json:"overall_status"`

    CreatedAt       time.Time `db:"created_at" json:"created_at"`
    UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}
