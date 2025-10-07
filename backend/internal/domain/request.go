package domain

import "time"

type ManpowerRequest struct {
	ID             int       `json:"id"`
	DocNo          string    `json:"doc_no"`
	DocDate        time.Time `json:"doc_date"`             // วันที่ออกเอกสาร
	DepartmentID   int       `json:"department_id"`
	RequestedBy    int       `json:"requested_by"`
	PositionTitle  string    `json:"position_title"`
	PositionLevel  string    `json:"position_level"`       // เช่น Officer, Supervisor, Manager
	NumRequired    int       `json:"num_required"`
	Reason         string    `json:"reason_note"`
	OriginStatus   string    `json:"origin_status"`        // สถานะฝั่งผู้ขอ
	HRStatus       string    `json:"hr_status"`            // สถานะฝั่ง HR
	ManagerStatus  string    `json:"manager_status"`       // สถานะฝั่งผู้อนุมัติ
	OverallStatus  string    `json:"overall_status"`       // สถานะรวม เช่น pending, approved, rejected
	Remark         string    `json:"remark"`               // หมายเหตุ / comment
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}