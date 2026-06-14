package repository

import (
	"context"
	"database/sql"
	"job_board/internal/domain"

	"github.com/google/uuid"
)

type NotificationRepo interface {

	CreateNotification(ctx context.Context, notification *domain.Notification) error 

	GetUserNotifications(ctx context.Context, userID uuid.UUID) ([]domain.Notification, error)

	MarkAsRead(ctx context.Context, noticiationID, userID uuid.UUID) error
}

type PostgresNotificationRepo struct {
	db *sql.DB
}

func NewPostgresNotificationRepo(db *sql.DB) *PostgresNotificationRepo {
	return &PostgresNotificationRepo{db: db}
}

func (r *PostgresNotificationRepo) CreateNotification(ctx context.Context, notification *domain.Notification) error {

	query := `
	INSERT INTO notifications (user_id, type, title, message, link)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at
	`

	return r.db.QueryRowContext(
		ctx, 
		query,
		notification.UserID,
		notification.Type,
		notification.Title,
		notification.Message,
		notification.Link,
	).Scan(&notification.ID, &notification.CreatedAt)
}

func (r *PostgresNotificationRepo) GetUserNotifications(ctx context.Context, userID uuid.UUID) ([]domain.Notification, error) {

	query := `
	SELECT
		id,
		user_id,
		type,
		title,
		message,
		link,
		is_read,
		created_at
	FROM notifications
	WHERE user_id = $1
	ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		userID,
	)

	if err != nil {	
		return nil, err
	}

	defer rows.Close()

	var notifications []domain.Notification

	for rows.Next() {

		var n domain.Notification

		err := rows.Scan(
			&n.ID,
			&n.UserID,
			&n.Type,
			&n.Title,
			&n.Message,
			&n.Link,
			&n.IsRead,
			&n.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		notifications = append(
			notifications,
			n,
		)
	}

	return notifications, nil
}

func (r *PostgresNotificationRepo) MarkAsRead(ctx context.Context, notificationID, userID uuid.UUID) error {

	result, err := r.db.ExecContext(
		ctx,
		`
		UPDATE notifications
		SET is_read = true
		WHERE id = $1
		AND user_id = $2
		`,
		notificationID,
		userID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}