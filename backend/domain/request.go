package domain

type ManpowerRequest struct {
    ID            int    `json:"id"`
    DocNo         string `json:"doc_no"`
    DepartmentID  int    `json:"department_id"`
    RequestedBy   int    `json:"requested_by"`
    PositionTitle string `json:"position_title"`
    NumRequired   int    `json:"num_required"`
    Reason        string `json:"reason"`
}
