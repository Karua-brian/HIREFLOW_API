package handlers

import (
	"job_board/internal/contextkeys"
	"job_board/internal/dto"
	"job_board/internal/service"
	"job_board/pkg/response"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type NotificationHandler interface {

	GetMyNotifications(w http.ResponseWriter, r *http.Request)

	MarkAsRead(w http.ResponseWriter, r *http.Request)
}

type notificationHandler struct {
	notificationService service.NotificationService
	logger *zap.Logger
}

func NewNotificationHandler(notificationService service.NotificationService, logger *zap.Logger) NotificationHandler {
	return &notificationHandler{
		notificationService: notificationService,
		logger: logger,
	}
}

func (h *notificationHandler) GetMyNotifications(w http.ResponseWriter, r *http.Request) {

	user, ok := contextkeys.UserFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	notifications, err := h.notificationService.GetMyNotifications(r.Context(), user.ID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to retrieve notifications")
		return
	}

	resp := dto.NotificationsListResponse{
		Notifications: make(
			[]dto.NotificationResponse,
			len(notifications),
		),
	}

	for i, n := range notifications {
		resp.Notifications[i] = 
			dto.NotificationResponse{
				ID: 		n.ID,
				Type: 		n.Type,
				Title: 		n.Title,
				Message: 	n.Message,
				Link: 		n.Link,
				IsRead: 	n.IsRead,
				CreatedAt:  n.CreatedAt	,
			}
	}

	response.JSON(w, http.StatusOK, resp)

}

func (h *notificationHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "id")

	id, err := uuid.Parse(idStr)

	if err != nil {
		response.Error(
			w,
			http.StatusBadRequest,
			"invalid notification id",
		)
		return
	}

	user, ok := contextkeys.UserFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	err = h.notificationService.MarkAsRead(r.Context(), id, user.ID)

	if err != nil {
		response.Error(
			w,
			http.StatusInternalServerError,
			"failed to mark notification as read",
		)
		return
	}

	response.JSON(
		w,
		http.StatusOK,
		map[string]string{
			"message": "notification marked as read",
		},
	)
}