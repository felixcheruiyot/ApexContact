package service

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/livestreamify/backend/internal/domain"
	"github.com/rs/zerolog/log"
)

// StartReminderScheduler polls every minute for event state transitions and
// every 5 minutes for reminder emails. It blocks until ctx is cancelled.
func StartReminderScheduler(ctx context.Context, db *pgxpool.Pool, notifSvc *NotificationService) {
	reminderTicker := time.NewTicker(5 * time.Minute)
	transitionTicker := time.NewTicker(1 * time.Minute)
	defer reminderTicker.Stop()
	defer transitionTicker.Stop()

	// Fire once immediately so restarts don't miss a window.
	autoTransitionEvents(ctx, db)
	dispatchReminders(ctx, db, notifSvc)

	for {
		select {
		case <-ctx.Done():
			return
		case <-transitionTicker.C:
			autoTransitionEvents(ctx, db)
		case <-reminderTicker.C:
			dispatchReminders(ctx, db, notifSvc)
		}
	}
}

// autoTransitionEvents moves scheduled events to live once their scheduled time has arrived.
func autoTransitionEvents(ctx context.Context, db *pgxpool.Pool) {
	_, err := db.Exec(ctx, `
		UPDATE events
		SET status = 'live', updated_at = NOW()
		WHERE status = 'scheduled'
		  AND scheduled_at <= NOW()
	`)
	if err != nil {
		log.Error().Err(err).Msg("auto_transition: failed to transition events to live")
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
