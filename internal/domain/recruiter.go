package domain

import (
	"time"

	"github.com/google/uuid"
)

type RecruiterRequest struct {
	ID  		 	uuid.UUID  
	RecruiterID 	uuid.UUID
	CompanyName 	string
	CompanyWebsite  string
	Message 		string
	Status 			string // "pending", "approved", "rejected"
	CreatedAt 		time.Time
}