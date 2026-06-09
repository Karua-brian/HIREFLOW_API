package domain

import (
	"time"

	"github.com/google/uuid"
)

// Job represents a business concept
// This is not tied to HTTP, JSON, or the database
// It represents how the business thinks abour a job posting
type Job struct {
	ID 				uuid.UUID       // unique identifier (set by database)
	Title 			string 		// Job title
	Description 	string 		// Detailed description
	Location 		string 		// Job location
	Company 		string 		// Company offering the job
	Salary 			string 		// Salary offered
	CreatedAt		time.Time	// When the job was created
}

