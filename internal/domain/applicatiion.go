package domain

import (
	"time"

	"github.com/google/uuid"
)

// Application represents a user applying to a job
type Application struct {
	ID		  uuid.UUID
	JobID     uuid.UUID
	UserID    uuid.UUID
	CreatedAt time.Time
}