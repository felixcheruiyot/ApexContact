package handlers

import (
	"context"
	"time"

	"github.com/livestreamify/backend/internal/config"
	"github.com/livestreamify/backend/internal/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProfileHandler struct {
	cfg *config.Config
	db  *pgxpool.Pool
}

func NewProfileHandler(cfg *config.Config, db *pgxpool.Pool) *ProfileHandler {
	return &ProfileHandler{cfg: cfg, db: db}
}

func (h *ProfileHandler) Get(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var u domain.User
	err := h.db.QueryRow(context.Background(),
		`SELECT id, email, full_name, phone, role, is_locked, age, gender, country, created_at, updated_at
		 FROM users WHERE id = $1`, userID,
	).Scan(&u.ID, &u.Email, &u.FullName, &u.Phone, &u.Role, &u.IsLocked,
		&u.Age, &u.Gender, &u.Country, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	return c.JSON(domain.Response{Data: u})
}

type updateProfileRequest struct {
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Age      *int   `json:"age"`
	Gender   string `json:"gender"`
	Country  string `json:"country"`
}

// MySubscriptions returns events the user has paid for.
type subscriptionWithEvent struct {
	SubscriptionID string    `json:"subscription_id"`
	EventID        string    `json:"event_id"`
	Title          string    `json:"title"`
	ThumbnailURL   string    `json:"thumbnail_url"`
	ScheduledAt    time.Time `json:"scheduled_at"`
	Status         string    `json:"status"`
	EventType      string    `json:"event_type"`
	Price          float64   `json:"price"`
	Currency       string    `json:"currency"`
	CreatedAt      time.Time `json:"created_at"`
}

func (h *ProfileHandler) MySubscriptions(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	rows, err := h.db.Query(context.Background(),
		`SELECT s.id, e.id, e.title, e.thumbnail_url, e.scheduled_at,
		        e.status, e.event_type, e.price, e.currency, s.created_at
		 FROM subscriptions s
		 JOIN events e ON e.id = s.event_id
		 WHERE s.user_id = $1
		 ORDER BY s.created_at DESC`, userID,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch subscriptions")
	}
	defer rows.Close()

	var results []subscriptionWithEvent
	for rows.Next() {
		var sw subscriptionWithEvent
		if err := rows.Scan(&sw.SubscriptionID, &sw.EventID, &sw.Title, &sw.ThumbnailURL,
			&sw.ScheduledAt, &sw.Status, &sw.EventType, &sw.Price, &sw.Currency, &sw.CreatedAt); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to scan subscription")
		}
		results = append(results, sw)
	}
	if results == nil {
		results = []subscriptionWithEvent{}
	}
	return c.JSON(domain.Response{Data: results})
}

func (h *ProfileHandler) Update(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req updateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	_, err := h.db.Exec(context.Background(),
		`UPDATE users SET full_name=$1, phone=$2, age=$3, gender=$4, country=$5, updated_at=NOW()
		 WHERE id=$6`,
		req.FullName, req.Phone, req.Age, req.Gender, req.Country, userID,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to update profile")
	}

	return c.JSON(domain.Response{Data: "profile updated"})
}
