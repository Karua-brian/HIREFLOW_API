package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID 		  uuid.UUID
	Email 	  string
	Password  string
	Role 	  string 	// "recruiter" "applicant" "admin"
	CreatedAt time.Time
	UpdatedAt time.Time
}


