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

// NotificationService sends transactional notifications and deduplicates via Redis.
type NotificationService struct {
	db           *pgxpool.Pool
	rdb          *redis.Client
	mailer       mailer.Mailer
	verifyTmpl   *template.Template
	reminderTmpl *template.Template
	appURL       string
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
	return &NotificationService{
		db:           db,
		rdb:          rdb,
		mailer:       m,
		verifyTmpl:   vt,
		reminderTmpl: rt,
		appURL:       appURL,
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

// SendWithdrawalOTP sends a 6-digit OTP email for withdrawal confirmation.
func (s *NotificationService) SendWithdrawalOTP(ctx context.Context, user *domain.User, otp, amount, currency string) error {
	body := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head><meta charset="UTF-8"></head>
<body style="margin:0;padding:0;background-color:#0A0A0F;font-family:Arial,sans-serif;">
  <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" style="background-color:#0A0A0F;">
    <tr><td align="center" style="padding:40px 16px;">
      <table role="presentation" width="600" cellpadding="0" cellspacing="0"
             style="max-width:600px;background-color:#141418;border-radius:12px;overflow:hidden;">
        <tr><td style="background-color:#E8002D;padding:20px 40px;">
          <p style="margin:0;font-size:22px;font-weight:700;letter-spacing:3px;color:#FFFFFF;text-transform:uppercase;">LIVE STREAMIFY</p>
        </td></tr>
        <tr><td style="padding:40px;">
          <h1 style="margin:0 0 16px;font-size:24px;font-weight:700;color:#FFFFFF;">Confirm Your Withdrawal</h1>
          <p style="margin:0 0 8px;font-size:15px;color:#A0A0B0;line-height:1.6;">Hi %s,</p>
          <p style="margin:0 0 24px;font-size:15px;color:#A0A0B0;line-height:1.6;">
            Enter the code below to confirm your withdrawal of
            <strong style="color:#FFFFFF;">%s %s</strong>. This code expires in 10 minutes.
          </p>
          <div style="text-align:center;margin:32px 0;">
            <span style="font-size:40px;font-weight:700;letter-spacing:12px;color:#FFFFFF;background-color:#1E1E26;padding:16px 28px;border-radius:8px;">%s</span>
          </div>
          <p style="margin:0;font-size:13px;color:#606070;line-height:1.6;">
            If you did not request this withdrawal, please contact support immediately.
          </p>
        </td></tr>
        <tr><td style="padding:24px 40px;border-top:1px solid #1E1E26;">
          <p style="margin:0;font-size:12px;color:#606070;text-align:center;">&copy; 2026 Live Streamify</p>
        </td></tr>
      </table>
    </td></tr>
  </table>
</body>
</html>`, user.FullName, currency, amount, otp)

	return s.mailer.Send(user.Email, "Your Withdrawal OTP — Live Streamify", body)
}

// SendEventUpdated emails a ticket holder informing them that event details have changed.
func (s *NotificationService) SendEventUpdated(ctx context.Context, user *domain.User, event *domain.Event) error {
	body := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head><meta charset="UTF-8"></head>
<body style="margin:0;padding:0;background-color:#0A0A0F;font-family:Arial,sans-serif;">
  <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" style="background-color:#0A0A0F;">
    <tr><td align="center" style="padding:40px 16px;">
      <table role="presentation" width="600" cellpadding="0" cellspacing="0"
             style="max-width:600px;background-color:#141418;border-radius:12px;overflow:hidden;">
        <tr><td style="background-color:#E8002D;padding:20px 40px;">
          <p style="margin:0;font-size:22px;font-weight:700;letter-spacing:3px;color:#FFFFFF;text-transform:uppercase;">LIVE STREAMIFY</p>
        </td></tr>
        <tr><td style="padding:40px;">
          <h1 style="margin:0 0 16px;font-size:24px;font-weight:700;color:#FFFFFF;">Event Updated</h1>
          <p style="margin:0 0 8px;font-size:15px;color:#A0A0B0;line-height:1.6;">Hi %s,</p>
          <p style="margin:0 0 24px;font-size:15px;color:#A0A0B0;line-height:1.6;">
            The event you have a ticket for has been updated by the organiser. Please review the latest details below.
          </p>
          <div style="background-color:#1E1E26;border-radius:8px;padding:20px;margin:0 0 24px;">
            <p style="margin:0 0 8px;font-size:18px;font-weight:700;color:#FFFFFF;">%s</p>
            <p style="margin:0;font-size:14px;color:#A0A0B0;">%s</p>
          </div>
          <a href="%s" style="display:inline-block;background-color:#E8002D;color:#FFFFFF;font-weight:700;
             font-size:15px;padding:14px 28px;border-radius:8px;text-decoration:none;">
            View Event Details
          </a>
          <p style="margin:24px 0 0;font-size:13px;color:#606070;line-height:1.6;">
            Your ticket remains valid. If you have questions, please contact the event organiser.
          </p>
        </td></tr>
        <tr><td style="padding:24px 40px;border-top:1px solid #1E1E26;">
          <p style="margin:0;font-size:12px;color:#606070;text-align:center;">&copy; 2026 Live Streamify</p>
        </td></tr>
      </table>
    </td></tr>
  </table>
</body>
</html>`,
		user.FullName,
		event.Title,
		event.ScheduledAt.Format("January 2, 2006 at 3:04 PM MST"),
		s.appURL+"/events/"+event.ID.String(),
	)

	return s.mailer.Send(user.Email, "Event Updated — "+event.Title, body)
}

// SendRaw sends a pre-rendered HTML email.
func (s *NotificationService) SendRaw(email, subject, body string) error {
	return s.mailer.Send(email, subject, body)
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
