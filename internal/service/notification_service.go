package service

import (
	"context"
	"job_board/internal/domain"
	"job_board/internal/repository"

	"github.com/google/uuid"
)

type NotificationService interface {

	CreateNotification(ctx context.Context, notification *domain.Notification) error

	GetMyNotifications(ctx context.Context, userID uuid.UUID) ([]domain.Notification, error)

	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
}

type notificationService struct {
	notificationRepo repository.NotificationRepo
}

func NewNotificationService(notificationRepo repository.NotificationRepo) NotificationService {
	return &notificationService{notificationRepo: notificationRepo}
}
func (s *notificationService) CreateNotification(ctx context.Context, notification *domain.Notification) error {
	
	return s.notificationRepo.CreateNotification(ctx, notification)
}

func (s *notificationService) GetMyNotifications(ctx context.Context, userID uuid.UUID) ([]domain.Notification, error) {
	
	return s.notificationRepo.GetUserNotifications(ctx, userID)
}

func (s *notificationService) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {

	return s.notificationRepo.MarkAllAsRead(ctx, userID)
}