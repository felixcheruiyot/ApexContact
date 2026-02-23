package service

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/livestreamify/backend/internal/domain"
	"github.com/rs/zerolog/log"
)

// StartReminderScheduler polls every 5 minutes for events starting in ~45 minutes
// and dispatches reminder emails to all subscribed, verified users.
// It blocks until ctx is cancelled.
func StartReminderScheduler(ctx context.Context, db *pgxpool.Pool, notifSvc *NotificationService) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	// Fire once immediately so restarts don't miss a window.
	dispatchReminders(ctx, db, notifSvc)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			dispatchReminders(ctx, db, notifSvc)
		}
	}
}

func dispatchReminders(ctx context.Context, db *pgxpool.Pool, notifSvc *NotificationService) {
	rows, err := db.Query(ctx, `
		SELECT DISTINCT
			u.id, u.email, u.full_name,
			e.id, e.title, e.scheduled_at
		FROM events e
		JOIN subscriptions sub ON sub.event_id = e.id
		JOIN users u ON u.id = sub.user_id
		WHERE e.status = 'scheduled'
		  AND e.scheduled_at BETWEEN NOW() + INTERVAL '40 minutes'
		                         AND NOW() + INTERVAL '50 minutes'
		  AND u.email_verified = true
	`)
	if err != nil {
		log.Error().Err(err).Msg("reminder_scheduler: query failed")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			user  domain.User
			event domain.Event
		)
		if err := rows.Scan(
			&user.ID, &user.Email, &user.FullName,
			&event.ID, &event.Title, &event.ScheduledAt,
		); err != nil {
			log.Error().Err(err).Msg("reminder_scheduler: scan failed")
			continue
		}

		if err := notifSvc.SendEventReminder(ctx, &user, &event); err != nil {
			log.Error().
				Err(err).
				Str("user_id", user.ID.String()).
				Str("event_id", event.ID.String()).
				Msg("reminder_scheduler: send failed")
		}
	}

	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("reminder_scheduler: rows error")
	}
}
