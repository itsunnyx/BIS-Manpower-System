package domain

import "time"

type Request struct {
    ID                 int       `json:"id"`
    RequestDate        time.Time `json:"request_date"`
    Department         string    `json:"department"`
    Division           string    `json:"division"`
    EmploymentType     string    `json:"employment_type"`
    ContractType       string    `json:"contract_type"`
    Reason             string    `json:"reason"`
    RequesterName      string    `json:"requester_name"`
    JobCode            string    `json:"job_code"`
    JobTitle           string    `json:"job_title"`
    AgeFrom            int       `json:"age_from"`
    AgeTo              int       `json:"age_to"`
    Gender             string    `json:"gender"`
    Nationality        string    `json:"nationality"`
    Experience         string    `json:"experience"`
    Education          string    `json:"education"`
    SpecialRequirements string   `json:"special_requirements"`
    CreatedAt          time.Time `json:"created_at"`
}
