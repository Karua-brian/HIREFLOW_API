package dto

import "github.com/google/uuid"

type CreateJobRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Company     string `json:"company"`
	Location    string `json:"location,omitempty"`
	Salary      string `json:"salary,omitempty"`
}

type ListJobsResponse struct {
	Jobs   []JobSummary `json:"jobs"`
	Limit  int          `json:"limit"`
	Offset int          `json:"offset"`
	Total  int64          `json:"total"`
}

type JobSummary struct {
	ID          uuid.UUID    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Company     string `json:"company"`
	Location    string `json:"location"`
	Salary      string `json:"salary"`
}