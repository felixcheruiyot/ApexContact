package handlers

import (
	"context"
	"time"

	"github.com/apexcontact/backend/internal/config"
	"github.com/apexcontact/backend/internal/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventHandler struct {
	cfg *config.Config
	db  *pgxpool.Pool
}

func NewEventHandler(cfg *config.Config, db *pgxpool.Pool) *EventHandler {
	return &EventHandler{cfg: cfg, db: db}
}

type createEventRequest struct {
	Title        string           `json:"title"`
	Description  string           `json:"description"`
	SportType    domain.SportType `json:"sport_type"`
	ScheduledAt  time.Time        `json:"scheduled_at"`
	Price        float64          `json:"price"`
	Currency     string           `json:"currency"`
	ThumbnailURL string           `json:"thumbnail_url"`
}

func (h *EventHandler) List(c *fiber.Ctx) error {
	sport := c.Query("sport")
	status := c.Query("status", "scheduled,live")

	query := `SELECT id, promoter_id, title, description, sport_type, scheduled_at,
	           status, price, currency, thumbnail_url, created_at, updated_at
	          FROM events WHERE status = ANY($1::text[])`
	args := []interface{}{status}

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
			&e.CreatedAt, &e.UpdatedAt); err != nil {
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
		        status, price, currency, thumbnail_url, created_at, updated_at
		 FROM events WHERE id = $1`, id,
	).Scan(&e.ID, &e.PromoterID, &e.Title, &e.Description, &e.SportType,
		&e.ScheduledAt, &e.Status, &e.Price, &e.Currency, &e.ThumbnailURL,
		&e.CreatedAt, &e.UpdatedAt)

	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "event not found")
	}

	return c.JSON(domain.Response{Data: e})
}

func (h *EventHandler) Create(c *fiber.Ctx) error {
	promoterID := c.Locals("user_id").(string)

	var req createEventRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	event := domain.Event{
		ID:           uuid.New(),
		PromoterID:   uuid.MustParse(promoterID),
		Title:        req.Title,
		Description:  req.Description,
		SportType:    req.SportType,
		ScheduledAt:  req.ScheduledAt,
		Status:       domain.StatusScheduled,
		Price:        req.Price,
		Currency:     req.Currency,
		ThumbnailURL: req.ThumbnailURL,
		StreamKey:    uuid.NewString(), // auto-generate stream key
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err := h.db.Exec(context.Background(),
		`INSERT INTO events (id, promoter_id, title, description, sport_type, scheduled_at,
		  status, price, currency, thumbnail_url, stream_key, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
		event.ID, event.PromoterID, event.Title, event.Description, event.SportType,
		event.ScheduledAt, event.Status, event.Price, event.Currency, event.ThumbnailURL,
		event.StreamKey, event.CreatedAt, event.UpdatedAt,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to create event")
	}

	return c.Status(fiber.StatusCreated).JSON(domain.Response{Data: event})
}

func (h *EventHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var req createEventRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	_, err := h.db.Exec(context.Background(),
		`UPDATE events SET title=$1, description=$2, sport_type=$3, scheduled_at=$4,
		  price=$5, currency=$6, thumbnail_url=$7, updated_at=$8
		 WHERE id=$9`,
		req.Title, req.Description, req.SportType, req.ScheduledAt,
		req.Price, req.Currency, req.ThumbnailURL, time.Now(), id,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to update event")
	}

	return c.JSON(domain.Response{Data: "event updated"})
}
