package dto

import (
	"time"

	"github.com/google/uuid"
)

type NotificationResponse struct {
	ID			uuid.UUID		`json:"id"`
	Type		string			`json:"type"`
	Title		string			`json:"title"`
	Message		string			`json:"message"`
	IsRead		bool			`json:"is_read"`
	CreatedAt	time.Time		`json:"created_at"`
}

type NotificationsListResponse struct {
	Notifications 		[]NotificationResponse		`json:"notifications"`
}