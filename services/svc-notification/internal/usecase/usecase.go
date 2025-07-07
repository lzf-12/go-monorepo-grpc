package usecase

import (
	"context"
	"fmt"
	"ops-monorepo/services/svc-notification/config"
	"ops-monorepo/services/svc-notification/internal/model"
	"ops-monorepo/services/svc-notification/internal/repository"
	"ops-monorepo/shared-libs/logger"
	"time"

	"github.com/google/uuid"
	"gopkg.in/gomail.v2"
)

type INotificationUsecase interface {
	SendEmail(ctx context.Context, to, subject, body, from string, isHTML bool) (*model.EmailLog, error)
}

type notificationUsecase struct {
	log  logger.Logger
	repo repository.INotificationSQLRepository
	cfg  *config.Config
}

func NewNotificationUsecase(log logger.Logger, repo repository.INotificationSQLRepository, cfg *config.Config) INotificationUsecase {
	return &notificationUsecase{
		log:  log,
		repo: repo,
		cfg:  cfg,
	}
}

func (u *notificationUsecase) SendEmail(ctx context.Context, to, subject, body, from string, isHTML bool) (*model.EmailLog, error) {
	// Create email log entry
	emailLog := &model.EmailLog{
		ID:      uuid.New().String(),
		To:      to,
		From:    from,
		Subject: subject,
		Body:    body,
		IsHTML:  isHTML,
		Status:  "pending",
	}

	// If from is empty, use default from config
	if from == "" {
		emailLog.From = u.cfg.SMTP.From
	}

	// Save initial log
	if err := u.repo.CreateEmailLog(ctx, emailLog); err != nil {
		u.log.Errorf("failed to create email log: %v", err)
		return nil, fmt.Errorf("failed to create email log: %w", err)
	}

	// Send email
	m := gomail.NewMessage()
	m.SetHeader("From", emailLog.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)

	if isHTML {
		m.SetBody("text/html", body)
	} else {
		m.SetBody("text/plain", body)
	}

	d := gomail.NewDialer(u.cfg.SMTP.Host, u.cfg.SMTP.Port, u.cfg.SMTP.Username, u.cfg.SMTP.Password)

	if err := d.DialAndSend(m); err != nil {
		u.log.Errorf("failed to send email: %v", err)
		errorMsg := err.Error()
		u.repo.UpdateEmailLogStatus(ctx, emailLog.ID, "failed", &errorMsg)
		return emailLog, fmt.Errorf("failed to send email: %w", err)
	}

	// Update status to sent
	if err := u.repo.UpdateEmailLogStatus(ctx, emailLog.ID, "sent", nil); err != nil {
		u.log.Errorf("failed to update email log status: %v", err)
	}

	emailLog.Status = "sent"
	emailLog.UpdatedAt = time.Now()

	u.log.Infof("email sent successfully to %s with subject: %s", to, subject)
	return emailLog, nil
}