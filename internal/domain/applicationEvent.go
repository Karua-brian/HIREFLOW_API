package domain

import "github.com/google/uuid"

// ApplicationEvent represents an event that occurs when a user applies to a job
type ApplicationEvent struct {
	JobID 	uuid.UUID
	UserID  uuid.UUID
}