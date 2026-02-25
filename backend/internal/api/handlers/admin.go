package handlers

import (
	"context"
	"time"

	"github.com/livestreamify/backend/internal/config"
	"github.com/livestreamify/backend/internal/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type AdminHandler struct {
	cfg *config.Config
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewAdminHandler(cfg *config.Config, db *pgxpool.Pool, rdb *redis.Client) *AdminHandler {
	return &AdminHandler{cfg: cfg, db: db, rdb: rdb}
}

func (h *AdminHandler) ListUsers(c *fiber.Ctx) error {
	rows, err := h.db.Query(context.Background(),
		`SELECT id, email, full_name, phone, role, is_locked, created_at FROM users ORDER BY created_at DESC`)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch users")
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		rows.Scan(&u.ID, &u.Email, &u.FullName, &u.Phone, &u.Role, &u.IsLocked, &u.CreatedAt)
		users = append(users, u)
	}

	return c.JSON(domain.Response{Data: users})
}

func (h *AdminHandler) ListFraudFlags(c *fiber.Ctx) error {
	rows, err := h.db.Query(context.Background(),
		`SELECT id, user_id, subscription_id, reason, detected_at, resolved
		 FROM fraud_flags WHERE resolved = false ORDER BY detected_at DESC`)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch fraud flags")
	}
	defer rows.Close()

	var flags []domain.FraudFlag
	for rows.Next() {
		var f domain.FraudFlag
		rows.Scan(&f.ID, &f.UserID, &f.SubscriptionID, &f.Reason, &f.DetectedAt, &f.Resolved)
		flags = append(flags, f)
	}

	return c.JSON(domain.Response{Data: flags})
}

func (h *AdminHandler) LockUser(c *fiber.Ctx) error {
	userID := c.Params("id")

	_, err := h.db.Exec(context.Background(),
		`UPDATE users SET is_locked=true, updated_at=NOW() WHERE id=$1`, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to lock user")
	}

	// Invalidate all active sessions for this user
	h.rdb.Set(context.Background(), "blocklist:"+userID, "1", 0)

	return c.JSON(domain.Response{Data: "user account locked"})
}

func (h *AdminHandler) UnlockUser(c *fiber.Ctx) error {
	userID := c.Params("id")

	_, err := h.db.Exec(context.Background(),
		`UPDATE users SET is_locked=false, updated_at=NOW() WHERE id=$1`, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to unlock user")
	}

	h.rdb.Del(context.Background(), "blocklist:"+userID)

	return c.JSON(domain.Response{Data: "user account unlocked"})
}

func (h *AdminHandler) UpdateUserRole(c *fiber.Ctx) error {
	userID := c.Params("id")

	var req struct {
		Role domain.UserRole `json:"role"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	switch req.Role {
	case domain.RoleMember, domain.RoleAdmin:
		// valid
	default:
		return fiber.NewError(fiber.StatusBadRequest, "invalid role")
	}

	tag, err := h.db.Exec(context.Background(),
		`UPDATE users SET role=$1, updated_at=NOW() WHERE id=$2`, req.Role, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to update user role")
	}
	if tag.RowsAffected() == 0 {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	return c.JSON(domain.Response{Data: "role updated"})
}

type adminUpdateEventRequest struct {
	Title        string             `json:"title"`
	Description  string             `json:"description"`
	SportType    domain.SportType   `json:"sport_type"`
	ScheduledAt  time.Time          `json:"scheduled_at"`
	Price        float64            `json:"price"`
	Currency     string             `json:"currency"`
	ThumbnailURL string             `json:"thumbnail_url"`
	Status       domain.EventStatus `json:"status"`
}

type adminEventItem struct {
	domain.Event
	PromoterEmail string `json:"promoter_email"`
	PromoterName  string `json:"promoter_name"`
}

func (h *AdminHandler) ListEvents(c *fiber.Ctx) error {
	statusFilter := c.Query("status")

	query := `SELECT e.id, e.promoter_id, e.title, e.description, e.sport_type, e.scheduled_at,
		        e.status, e.price, e.currency, e.thumbnail_url, e.event_type, e.review_note,
		        e.created_at, e.updated_at,
		        COALESCE(u.email, '') AS promoter_email,
		        COALESCE(u.full_name, '') AS promoter_name
		 FROM events e
		 LEFT JOIN users u ON u.id = e.promoter_id`
	var args []interface{}
	if statusFilter != "" {
		query += " WHERE e.status=$1"
		args = append(args, statusFilter)
	}
	query += " ORDER BY e.created_at DESC"

	rows, err := h.db.Query(context.Background(), query, args...)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch events")
	}
	defer rows.Close()

	var events []adminEventItem
	for rows.Next() {
		var item adminEventItem
		if err := rows.Scan(
			&item.ID, &item.PromoterID, &item.Title, &item.Description, &item.SportType,
			&item.ScheduledAt, &item.Status, &item.Price, &item.Currency, &item.ThumbnailURL,
			&item.EventType, &item.ReviewNote, &item.CreatedAt, &item.UpdatedAt,
			&item.PromoterEmail, &item.PromoterName,
		); err != nil {
			continue
		}
		events = append(events, item)
	}

	return c.JSON(domain.Response{Data: events})
}

func (h *AdminHandler) ApproveEvent(c *fiber.Ctx) error {
	id := c.Params("id")

	tag, err := h.db.Exec(context.Background(),
		`UPDATE events SET status='scheduled', review_note='', updated_at=NOW()
		 WHERE id=$1 AND status='pending_review'`, id,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to approve event")
	}
	if tag.RowsAffected() == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "event not found or not pending review")
	}

	return c.JSON(domain.Response{Data: "event approved"})
}

func (h *AdminHandler) RequestEdits(c *fiber.Ctx) error {
	id := c.Params("id")

	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if req.Reason == "" {
		return fiber.NewError(fiber.StatusBadRequest, "reason is required")
	}

	tag, err := h.db.Exec(context.Background(),
		`UPDATE events SET status='draft', review_note=$1, updated_at=NOW()
		 WHERE id=$2 AND status='pending_review'`,
		req.Reason, id,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to request edits")
	}
	if tag.RowsAffected() == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "event not found or not pending review")
	}

	return c.JSON(domain.Response{Data: "edits requested"})
}

func (h *AdminHandler) DeclineEvent(c *fiber.Ctx) error {
	id := c.Params("id")

	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if req.Reason == "" {
		return fiber.NewError(fiber.StatusBadRequest, "reason is required")
	}

	tag, err := h.db.Exec(context.Background(),
		`UPDATE events SET status='declined', review_note=$1, updated_at=NOW()
		 WHERE id=$2 AND status='pending_review'`,
		req.Reason, id,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to decline event")
	}
	if tag.RowsAffected() == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "event not found or not pending review")
	}

	return c.JSON(domain.Response{Data: "event declined"})
}

func (h *AdminHandler) UpdateEvent(c *fiber.Ctx) error {
	id := c.Params("id")

	var req adminUpdateEventRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	tag, err := h.db.Exec(context.Background(),
		`UPDATE events SET title=$1, description=$2, sport_type=$3, scheduled_at=$4,
		  price=$5, currency=$6, thumbnail_url=$7, status=$8, updated_at=NOW()
		 WHERE id=$9`,
		req.Title, req.Description, req.SportType, req.ScheduledAt,
		req.Price, req.Currency, req.ThumbnailURL, req.Status, id,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to update event")
	}
	if tag.RowsAffected() == 0 {
		return fiber.NewError(fiber.StatusNotFound, "event not found")
	}

	return c.JSON(domain.Response{Data: "event updated"})
}

func (h *AdminHandler) PlatformAnalytics(c *fiber.Ctx) error {
	var stats struct {
		TotalUsers      int     `json:"total_users"`
		TotalEvents     int     `json:"total_events"`
		LiveEvents      int     `json:"live_events"`
		TotalRevenue    float64 `json:"total_revenue"`
		ActiveViewers   int     `json:"active_viewers"`
		FraudFlagsOpen  int     `json:"fraud_flags_open"`
	}

	h.db.QueryRow(context.Background(), `SELECT COUNT(*) FROM users`).Scan(&stats.TotalUsers)
	h.db.QueryRow(context.Background(), `SELECT COUNT(*) FROM events`).Scan(&stats.TotalEvents)
	h.db.QueryRow(context.Background(), `SELECT COUNT(*) FROM events WHERE status='live'`).Scan(&stats.LiveEvents)
	h.db.QueryRow(context.Background(), `SELECT COALESCE(SUM(amount),0) FROM payments WHERE status='success'`).Scan(&stats.TotalRevenue)
	h.db.QueryRow(context.Background(), `SELECT COUNT(*) FROM fraud_flags WHERE resolved=false`).Scan(&stats.FraudFlagsOpen)

	return c.JSON(domain.Response{Data: stats})
}
