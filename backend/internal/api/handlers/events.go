package handlers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/livestreamify/backend/internal/config"
	"github.com/livestreamify/backend/internal/domain"
	"github.com/livestreamify/backend/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventHandler struct {
	cfg      *config.Config
	db       *pgxpool.Pool
	notifSvc *service.NotificationService
}

func NewEventHandler(cfg *config.Config, db *pgxpool.Pool, notifSvc *service.NotificationService) *EventHandler {
	return &EventHandler{cfg: cfg, db: db, notifSvc: notifSvc}
}

type createEventRequest struct {
	Title        string           `json:"title"`
	Description  string           `json:"description"`
	SportType    domain.SportType `json:"sport_type"`
	ScheduledAt  time.Time        `json:"scheduled_at"`
	Price        float64          `json:"price"`
	Currency     string           `json:"currency"`
	ThumbnailURL string           `json:"thumbnail_url"`
	TeaserHook   string           `json:"teaser_hook"`
	IsPublic     bool             `json:"is_public"`
}

func (h *EventHandler) List(c *fiber.Ctx) error {
	sport := c.Query("sport")
	// Split comma-separated statuses into a []string so pgx encodes it as a proper text[]
	statusList := strings.Split(c.Query("status", "scheduled,live"), ",")

	query := `SELECT id, promoter_id, title, description, sport_type, scheduled_at,
	           status, price, currency, thumbnail_url, event_type, teaser_hook,
	           is_public, created_at, updated_at
	          FROM events WHERE status = ANY($1::text[]) AND event_type = 'video'`
	args := []interface{}{statusList}

	if sport != "" {
		query += " AND sport_type = $2"
		args = append(args, sport)
	}
	query += " ORDER BY scheduled_at ASC"

	rows, err := h.db.Query(context.Background(), query, args...)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch events")
	}
	defer rows.Close()

	var events []domain.Event
	for rows.Next() {
		var e domain.Event
		if err := rows.Scan(&e.ID, &e.PromoterID, &e.Title, &e.Description, &e.SportType,
			&e.ScheduledAt, &e.Status, &e.Price, &e.Currency, &e.ThumbnailURL,
			&e.EventType, &e.TeaserHook, &e.IsPublic, &e.CreatedAt, &e.UpdatedAt); err != nil {
			continue
		}
		events = append(events, e)
	}

	return c.JSON(domain.Response{Data: events})
}

func (h *EventHandler) Get(c *fiber.Ctx) error {
	id := c.Params("id")

	var e domain.Event
	err := h.db.QueryRow(context.Background(),
		`SELECT id, promoter_id, title, description, sport_type, scheduled_at,
		        status, price, currency, thumbnail_url, event_type, teaser_hook,
		        is_public, hls_path, created_at, updated_at
		 FROM events WHERE id = $1`, id,
	).Scan(&e.ID, &e.PromoterID, &e.Title, &e.Description, &e.SportType,
		&e.ScheduledAt, &e.Status, &e.Price, &e.Currency, &e.ThumbnailURL,
		&e.EventType, &e.TeaserHook, &e.IsPublic, &e.HLSPath, &e.CreatedAt, &e.UpdatedAt)

	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "event not found")
	}
	e.StreamActive = e.HLSPath != ""

	return c.JSON(domain.Response{Data: e})
}

func (h *EventHandler) Create(c *fiber.Ctx) error {
	promoterID := c.Locals("user_id").(string)

	var req createEventRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if req.ScheduledAt.Before(time.Now().Add(15 * time.Minute)) {
		return fiber.NewError(fiber.StatusBadRequest, "scheduled date must be at least 15 minutes in the future")
	}

	currency := req.Currency
	if currency == "" {
		currency = "KES"
	}

	event := domain.Event{
		ID:           uuid.New(),
		PromoterID:   uuid.MustParse(promoterID),
		Title:        req.Title,
		Description:  req.Description,
		SportType:    req.SportType,
		ScheduledAt:  req.ScheduledAt,
		Status:       domain.StatusDraft,
		Price:        req.Price,
		Currency:     currency,
		ThumbnailURL: req.ThumbnailURL,
		TeaserHook:   req.TeaserHook,
		IsPublic:     req.IsPublic,
		StreamKey:    uuid.NewString(), // auto-generate stream key
		EventType:    domain.EventTypeVideo,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err := h.db.Exec(context.Background(),
		`INSERT INTO events (id, promoter_id, title, description, sport_type, scheduled_at,
		  status, price, currency, thumbnail_url, teaser_hook, is_public, stream_key,
		  event_type, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)`,
		event.ID, event.PromoterID, event.Title, event.Description, event.SportType,
		event.ScheduledAt, event.Status, event.Price, event.Currency, event.ThumbnailURL,
		event.TeaserHook, event.IsPublic, event.StreamKey, event.EventType,
		event.CreatedAt, event.UpdatedAt,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to create event")
	}

	return c.Status(fiber.StatusCreated).JSON(domain.Response{Data: event})
}

func (h *EventHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	promoterID := c.Locals("user_id").(string)

	var req createEventRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if req.ScheduledAt.Before(time.Now().Add(15 * time.Minute)) {
		return fiber.NewError(fiber.StatusBadRequest, "scheduled date must be at least 15 minutes in the future")
	}

	tag, err := h.db.Exec(context.Background(),
		`UPDATE events SET title=$1, description=$2, sport_type=$3, scheduled_at=$4,
		  price=$5, currency=$6, thumbnail_url=$7, teaser_hook=$8, is_public=$9, updated_at=$10
		 WHERE id=$11 AND promoter_id=$12
		   AND (status='draft' OR (status='scheduled' AND scheduled_at > NOW()))`,
		req.Title, req.Description, req.SportType, req.ScheduledAt,
		req.Price, req.Currency, req.ThumbnailURL, req.TeaserHook, req.IsPublic,
		time.Now(), id, promoterID,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to update event")
	}
	if tag.RowsAffected() == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "event not found, not editable, or access denied")
	}

	// Asynchronously notify subscribers about the update
	go h.notifyEventUpdated(id)

	return c.JSON(domain.Response{Data: "event updated"})
}

// Discover returns up to 20 public events (upcoming + recently completed).
func (h *EventHandler) Discover(c *fiber.Ctx) error {
	q := c.Query("q")
	sport := c.Query("sport")

	query := `SELECT id, promoter_id, title, description, sport_type, scheduled_at,
	           status, price, currency, thumbnail_url, event_type, teaser_hook, is_public, created_at, updated_at
	          FROM events
	          WHERE is_public = true
	            AND status IN ('scheduled','live','completed')
	            AND (status != 'completed' OR scheduled_at >= NOW() - INTERVAL '7 days')`

	args := []interface{}{}
	argIdx := 1

	if sport != "" {
		argIdx++
		query += " AND sport_type = $" + itoa(argIdx)
		args = append(args, sport)
	}
	if q != "" {
		argIdx++
		query += " AND (title ILIKE $" + itoa(argIdx) + " OR description ILIKE $" + itoa(argIdx) + ")"
		args = append(args, "%"+q+"%")
	}

	query += " ORDER BY CASE status WHEN 'live' THEN 0 WHEN 'scheduled' THEN 1 ELSE 2 END, scheduled_at DESC LIMIT 20"

	rows, err := h.db.Query(context.Background(), query, args...)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch events")
	}
	defer rows.Close()

	var events []domain.Event
	for rows.Next() {
		var e domain.Event
		if err := rows.Scan(&e.ID, &e.PromoterID, &e.Title, &e.Description, &e.SportType,
			&e.ScheduledAt, &e.Status, &e.Price, &e.Currency, &e.ThumbnailURL,
			&e.EventType, &e.TeaserHook, &e.IsPublic, &e.CreatedAt, &e.UpdatedAt); err != nil {
			continue
		}
		events = append(events, e)
	}
	if events == nil {
		events = []domain.Event{}
	}
	return c.JSON(domain.Response{Data: events})
}

// notifyEventUpdated queries all ticket holders for the event and emails each one.
func (h *EventHandler) notifyEventUpdated(eventID string) {
	ctx := context.Background()

	var evt domain.Event
	err := h.db.QueryRow(ctx,
		`SELECT id, title, scheduled_at FROM events WHERE id = $1`, eventID,
	).Scan(&evt.ID, &evt.Title, &evt.ScheduledAt)
	if err != nil {
		return
	}

	rows, err := h.db.Query(ctx,
		`SELECT u.id, u.email, u.full_name
		 FROM subscriptions s
		 JOIN users u ON u.id = s.user_id
		 WHERE s.event_id = $1`, eventID,
	)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Email, &user.FullName); err != nil {
			continue
		}
		_ = h.notifSvc.SendEventUpdated(ctx, &user, &evt)
	}
}

func itoa(n int) string {
	return fmt.Sprintf("%d", n)
}
