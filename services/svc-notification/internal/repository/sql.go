package repository

import (
	"context"
	"ops-monorepo/services/svc-notification/internal/model"
	pg "ops-monorepo/shared-libs/storage/postgres"
	"time"

	"github.com/google/uuid"
)

type INotificationSQLRepository interface {
	CreateEmailLog(ctx context.Context, emailLog *model.EmailLog) error
	UpdateEmailLogStatus(ctx context.Context, id string, status string, errorMsg *string) error
	GetEmailLogByID(ctx context.Context, id string) (*model.EmailLog, error)
}

type notificationSQLRepository struct {
	db *pg.PostgresPgx
}

func NewNotificationRepository(db *pg.PostgresPgx) INotificationSQLRepository {
	return &notificationSQLRepository{
		db: db,
	}
}

func (r *notificationSQLRepository) CreateEmailLog(ctx context.Context, emailLog *model.EmailLog) error {
	if emailLog.ID == "" {
		emailLog.ID = uuid.New().String()
	}
	emailLog.CreatedAt = time.Now()
	emailLog.UpdatedAt = time.Now()

	query := `
		INSERT INTO email_logs (id, to_email, from_email, subject, body, is_html, status, error, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.Pool().Exec(ctx, query,
		emailLog.ID,
		emailLog.To,
		emailLog.From,
		emailLog.Subject,
		emailLog.Body,
		emailLog.IsHTML,
		emailLog.Status,
		emailLog.Error,
		emailLog.CreatedAt,
		emailLog.UpdatedAt,
	)

	return err
}

func (r *notificationSQLRepository) UpdateEmailLogStatus(ctx context.Context, id string, status string, errorMsg *string) error {
	query := `
		UPDATE email_logs 
		SET status = $2, error = $3, updated_at = $4
		WHERE id = $1
	`

	_, err := r.db.Pool().Exec(ctx, query, id, status, errorMsg, time.Now())
	return err
}

func (r *notificationSQLRepository) GetEmailLogByID(ctx context.Context, id string) (*model.EmailLog, error) {
	query := `
		SELECT id, to_email, from_email, subject, body, is_html, status, error, created_at, updated_at
		FROM email_logs
		WHERE id = $1
	`

	var emailLog model.EmailLog
	err := r.db.Pool().QueryRow(ctx, query, id).Scan(
		&emailLog.ID,
		&emailLog.To,
		&emailLog.From,
		&emailLog.Subject,
		&emailLog.Body,
		&emailLog.IsHTML,
		&emailLog.Status,
		&emailLog.Error,
		&emailLog.CreatedAt,
		&emailLog.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &emailLog, nil
}
