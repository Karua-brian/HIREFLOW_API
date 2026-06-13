package service

import (
	"context"
	"job_board/internal/domain"
	"job_board/internal/repository"

	"github.com/google/uuid"
)

type NotificationService interface {

	GetMyNotifications(ctx context.Context, userID uuid.UUID) ([]domain.Notification, error)

	MarkAsRead(ctx context.Context, notificationID, userID uuid.UUID) error
}

type notificationService struct {
	notificationRepo repository.NotificationRepo
}

func NewNotificationService(notificationRepo repository.NotificationRepo) NotificationService {
	return &notificationService{notificationRepo: notificationRepo}
}

func (s *notificationService) GetMyNotifications(ctx context.Context, userID uuid.UUID) ([]domain.Notification, error) {
	
	return s.notificationRepo.GetUserNotifications(ctx, userID)
}

func (s *notificationService) MarkAsRead(ctx context.Context, notificationID, userID uuid.UUID) error {

	return s.notificationRepo.MarkAsRead(ctx, notificationID, userID)
}