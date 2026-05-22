package domain

import "time"

type User struct {
	ID 		  int64
	Email 	  string
	Password  string
	Role 	  string 	// "recruiter" "applicant" "admin"
	CreatedAt time.Time
	UpdatedAt time.Time
}