package dto

import "github.com/google/uuid"

type CreateJobRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Location    string `json:"location,omitempty"`
	Company     string `json:"company_name"`
	Salary      string `json:"salary_range,omitempty"`
}

type ListJobsResponse struct {
	Jobs   []JobSummary `json:"jobs"`
	Limit  int          `json:"limit"`
	Offset int          `json:"offset"`
	Total  int64         `json:"total"`
}

type JobSummary struct {
	ID          uuid.UUID   `json:"id"`
	Title       string 		`json:"title"`
	Description string 		`json:"description"`
	Location    string 		`json:"location"`
	Company     string 		`json:"company_name"`
	Salary      string 		`json:"salary_range"`
}