package domain

// ApplicationEvent represents an event that occurs when a user applies to a job
type ApplicationEvent struct {
	JobID 	int64
	UserID  int64
}