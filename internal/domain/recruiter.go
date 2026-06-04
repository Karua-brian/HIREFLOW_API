package domain

import "time"

type RecruiterRequest struct {
	ID  		 	int64  
	UserID 		 	int64
	CompanyName 	string
	CompanyWebsite  string
	Message 		string
	Status 			string // "pending", "approved", "rejected"
	CreatedAt 		time.Time
}