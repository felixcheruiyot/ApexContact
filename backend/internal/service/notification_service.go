package service

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/livestreamify/backend/internal/domain"
	"github.com/livestreamify/backend/internal/integrations/mailer"
	"github.com/redis/go-redis/v9"
)

//go:embed templates/verification.html
var verificationTmplRaw string

//go:embed templates/reminder.html
var reminderTmplRaw string

//go:embed templates/withdrawal_otp.html
var withdrawalOTPTmplRaw string

// NotificationService sends transactional notifications and deduplicates via Redis.
type NotificationService struct {
	db                *pgxpool.Pool
	rdb               *redis.Client
	mailer            mailer.Mailer
	verifyTmpl        *template.Template
	reminderTmpl      *template.Template
	withdrawalOTPTmpl *template.Template
	appURL            string
}

// NewNotificationService parses email templates and returns a ready service.
func NewNotificationService(db *pgxpool.Pool, rdb *redis.Client, m mailer.Mailer, appURL string) (*NotificationService, error) {
	vt, err := template.New("verification.html").Parse(verificationTmplRaw)
	if err != nil {
		return nil, fmt.Errorf("parse verification template: %w", err)
	}
	rt, err := template.New("reminder.html").Parse(reminderTmplRaw)
	if err != nil {
		return nil, fmt.Errorf("parse reminder template: %w", err)
	}
	wt, err := template.New("withdrawal_otp.html").Parse(withdrawalOTPTmplRaw)
	if err != nil {
		return nil, fmt.Errorf("parse withdrawal OTP template: %w", err)
	}
	return &NotificationService{
		db:                db,
		rdb:               rdb,
		mailer:            m,
		verifyTmpl:        vt,
		reminderTmpl:      rt,
		withdrawalOTPTmpl: wt,
		appURL:            appURL,
	}, nil
}

// SendVerificationEmail renders and sends the email verification message.
func (s *NotificationService) SendVerificationEmail(ctx context.Context, user *domain.User, token string) error {
	data := struct {
		FullName        string
		VerificationURL string
	}{
		FullName:        user.FullName,
		VerificationURL: s.appURL + "/api/v1/auth/verify-email?token=" + token,
	}

	var buf bytes.Buffer
	if err := s.verifyTmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("render verification email: %w", err)
	}

	return s.mailer.Send(user.Email, "Verify Your Email — Live Streamify", buf.String())
}

// SendEventReminder renders and sends a 45-minute event reminder.
// It skips sending if a reminder was already sent for this user+event pair.
func (s *NotificationService) SendEventReminder(ctx context.Context, user *domain.User, event *domain.Event) error {
	userID := user.ID.String()
	eventID := event.ID.String()

	if s.hasBeenSent(ctx, userID, "event_reminder", eventID) {
		return nil
	}

	data := struct {
		FullName   string
		EventTitle string
		EventURL   string
		StartTime  string
	}{
		FullName:   user.FullName,
		EventTitle: event.Title,
		EventURL:   s.appURL + "/events/" + eventID,
		StartTime:  event.ScheduledAt.Format("January 2, 2006 at 3:04 PM MST"),
	}

	var buf bytes.Buffer
	if err := s.reminderTmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("render reminder email: %w", err)
	}

	if err := s.mailer.Send(user.Email, "Your Event Starts in 45 Minutes — "+event.Title, buf.String()); err != nil {
		return err
	}

	return s.markSent(ctx, userID, "event_reminder", eventID)
}

// SendWithdrawalOTP renders and sends the withdrawal OTP email.
func (s *NotificationService) SendWithdrawalOTP(ctx context.Context, user *domain.User, otp, amount, currency string) error {
	data := struct {
		FullName string
		OTP      string
		Amount   string
		Currency string
	}{
		FullName: user.FullName,
		OTP:      otp,
		Amount:   amount,
		Currency: currency,
	}

	var buf bytes.Buffer
	if err := s.withdrawalOTPTmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("render withdrawal OTP email: %w", err)
	}

	return s.mailer.Send(user.Email, "Withdrawal Verification Code — Live Streamify", buf.String())
}

// SendRaw sends an arbitrary HTML email. Used for service-level notifications
// where no template is required (e.g. withdrawal failure alerts).
func (s *NotificationService) SendRaw(to, subject, htmlBody string) error {
	return s.mailer.Send(to, subject, htmlBody)
}

// hasBeenSent checks the Redis dedup key to prevent double-sends.
func (s *NotificationService) hasBeenSent(ctx context.Context, userID, notifType, refID string) bool {
	key := fmt.Sprintf("notif:%s:%s:%s", notifType, userID, refID)
	val, _ := s.rdb.Get(ctx, key).Result()
	return val == "1"
}

// markSent sets the Redis dedup key and appends a row to notification_log.
func (s *NotificationService) markSent(ctx context.Context, userID, notifType, refID string) error {
	key := fmt.Sprintf("notif:%s:%s:%s", notifType, userID, refID)
	s.rdb.Set(ctx, key, "1", 2*time.Hour)

	_, err := s.db.Exec(ctx,
		`INSERT INTO notification_log (user_id, notification_type, reference_id)
		 VALUES ($1, $2, $3)`,
		userID, notifType, refID,
	)
	return err
}
