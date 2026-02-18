package handlers

import (
	"context"

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
