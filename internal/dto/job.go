package dto

import ()

type CreateJobRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Company     string `json:"company"`
}

type ListJobsResponse struct {
	Jobs   []JobSummary `json:"jobs"`
	Limit  int          `json:"limit"`
	Offset int          `json:"offset"`
	Total  int64          `json:"total"`
}

type JobSummary struct {
	ID          int64    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Company     string `json:"company"`
}