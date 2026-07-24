package service

import (
	"context"
	"job_board/internal/domain"
	"job_board/internal/repository"

	"github.com/google/uuid"
)

type NotificationService interface {

	CreateNotification(ctx context.Context, notification *domain.Notification) error

	NewRecruiterRequest(ctx context.Context, request *domain.RecruiterRequest) error

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

func (s *notificationService) NewRecruiterRequestNotification(ctx context.Context, req *domain.RecruiterRequest) error {

	adminID := uuid.MustParse("admin_id")

	notification := &domain.Notification{
		UserID: adminID,
		Type: "recruiter_request",
		Title: "New Recruiter Request",
		Message: req.CompanyName + " submitted recruiter access request",
		Link: "/admin/recruiter-requests",
	}

	return s.notificationRepo.CreateNotification(ctx, notification)
}

func (s *notificationService) GetMyNotifications(ctx context.Context, userID uuid.UUID) ([]domain.Notification, error) {
	
	return s.notificationRepo.GetUserNotifications(ctx, userID)
}

func (s *notificationService) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {

	return s.notificationRepo.MarkAllAsRead(ctx, userID)
}