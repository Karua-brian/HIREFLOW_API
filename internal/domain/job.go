package domain

import "time"

// Job represents a business concept
// This is not tied to HTTP, JSON, or the database
// It represents how the business thinks abour a job posting
type Job struct {
	ID 				int64       // unique identifier (set by database)
	Title 			string 		// Job title
	Description 	string 		// Detailed description
	Company 		string 		// Company offering the job
	Location 		string 		// Job location
	Salary 			string 		// Salary offered
	CreatedAt		time.Time	// When the job was created
	CreatedBy		int64		// User ID of the recruiter who created it
}

