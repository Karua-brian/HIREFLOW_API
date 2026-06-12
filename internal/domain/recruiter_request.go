package domain

import (
	"time"

	"github.com/google/uuid"
)

type RecruiterRequest struct {
	ID  		 	uuid.UUID  
	UserID 		    uuid.UUID
	CompanyName 	string
	CompanyWebsite  string
	Message 		string
	Status 			string // "pending", "approved", "rejected"
	Reason 			string
	CreatedAt 		time.Time
	UpdatedAt		time.Time
}