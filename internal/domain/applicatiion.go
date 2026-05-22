package domain

import "time"

// Application represents a user applying to a job
type Application struct {
	ID		  int64
	JobID     int64
	UserID    int64
	CreatedAt time.Time
}