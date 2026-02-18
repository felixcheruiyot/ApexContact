package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/livestreamify/backend/internal/config"
	"github.com/livestreamify/backend/internal/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type StreamHandler struct {
	cfg *config.Config
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewStreamHandler(cfg *config.Config, db *pgxpool.Pool, rdb *redis.Client) *StreamHandler {
	return &StreamHandler{cfg: cfg, db: db, rdb: rdb}
}

// Subscribe creates a subscription record after a successful payment.
// Called after payment confirmation — the frontend hits this with the payment ID.
func (h *StreamHandler) Subscribe(c *fiber.Ctx) error {
	eventID := c.Params("eventId")
	userID := c.Locals("user_id").(string)

	var req struct {
		PaymentID         string `json:"payment_id"`
		DeviceFingerprint string `json:"device_fingerprint"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	// Verify payment is confirmed and belongs to this user+event
	var payment domain.Payment
	err := h.db.QueryRow(context.Background(),
		`SELECT id, user_id, event_id, status FROM payments WHERE id=$1`, req.PaymentID,
	).Scan(&payment.ID, &payment.UserID, &payment.EventID, &payment.Status)

	if err != nil || payment.Status != domain.PaymentSuccess {
		return fiber.NewError(fiber.StatusPaymentRequired, "valid payment not found")
	}
	if payment.UserID.String() != userID || payment.EventID.String() != eventID {
		return fiber.NewError(fiber.StatusForbidden, "payment mismatch")
	}

	// Check for existing subscription (idempotent)
	var existing domain.Subscription
	_ = h.db.QueryRow(context.Background(),
		`SELECT id, stream_token FROM subscriptions WHERE user_id=$1 AND event_id=$2`,
		userID, eventID,
	).Scan(&existing.ID, &existing.StreamToken)

	if existing.StreamToken != "" {
		return c.JSON(domain.Response{Data: fiber.Map{"stream_token": existing.StreamToken}})
	}

	// Fetch event end time for token expiry
	var scheduledAt time.Time
	h.db.QueryRow(context.Background(), `SELECT scheduled_at FROM events WHERE id=$1`, eventID).Scan(&scheduledAt)
	expiresAt := scheduledAt.Add(6 * time.Hour) // stream token valid for event + 6h buffer

	sub := domain.Subscription{
		ID:                uuid.New(),
		UserID:            uuid.MustParse(userID),
		EventID:           uuid.MustParse(eventID),
		PaymentID:         payment.ID,
		StreamToken:       uuid.NewString(),
		DeviceFingerprint: req.DeviceFingerprint,
		ExpiresAt:         expiresAt,
		CreatedAt:         time.Now(),
	}

	_, err = h.db.Exec(context.Background(),
		`INSERT INTO subscriptions (id, user_id, event_id, payment_id, stream_token, device_fingerprint, expires_at, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		sub.ID, sub.UserID, sub.EventID, sub.PaymentID, sub.StreamToken,
		sub.DeviceFingerprint, sub.ExpiresAt, sub.CreatedAt,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to create subscription")
	}

	// Cache token → eventID mapping in Redis for fast lookup
	h.rdb.Set(context.Background(),
		fmt.Sprintf("token:%s", sub.StreamToken),
		eventID,
		time.Until(expiresAt),
	)

	return c.Status(fiber.StatusCreated).JSON(domain.Response{
		Data: fiber.Map{"stream_token": sub.StreamToken},
	})
}

// GetToken validates the stream token and returns the signed HLS playlist URL.
// Admins bypass the subscription check and get free access to any live event.
func (h *StreamHandler) GetToken(c *fiber.Ctx) error {
	eventID := c.Params("eventId")
	role, _ := c.Locals("role").(domain.UserRole)

	if role == domain.RoleAdmin {
		// Admins watch for free — skip token validation entirely.
		var hlsPath, streamKey string
		h.db.QueryRow(context.Background(),
			`SELECT hls_path, stream_key FROM events WHERE id=$1 AND status='live'`, eventID,
		).Scan(&hlsPath, &streamKey)

		if hlsPath == "" {
			return fiber.NewError(fiber.StatusNotFound, "stream not yet available — event may not have started")
		}

		hlsURL := fmt.Sprintf("%s/hls/%s/index.m3u8", h.cfg.FrontendURL, streamKey)
		return c.JSON(domain.Response{Data: fiber.Map{"hls_url": hlsURL}})
	}

	token := c.Locals("stream_token").(string) // set by AntiPiracy middleware

	// Verify token is for this event
	cachedEventID, err := h.rdb.Get(context.Background(), fmt.Sprintf("token:%s", token)).Result()
	if err != nil || cachedEventID != eventID {
		return fiber.NewError(fiber.StatusForbidden, "invalid stream token for this event")
	}

	// Fetch HLS path from DB
	var hlsPath, streamKey string
	h.db.QueryRow(context.Background(),
		`SELECT hls_path, stream_key FROM events WHERE id=$1 AND status='live'`, eventID,
	).Scan(&hlsPath, &streamKey)

	if hlsPath == "" {
		return fiber.NewError(fiber.StatusNotFound, "stream not yet available — event may not have started")
	}

	// Use the public frontend URL so the browser can reach the HLS stream via nginx.
	hlsURL := fmt.Sprintf("%s/hls/%s/index.m3u8?token=%s", h.cfg.FrontendURL, streamKey, token)

	return c.JSON(domain.Response{
		Data: fiber.Map{"hls_url": hlsURL},
	})
}

// IngestCallback is called by the Nginx-RTMP media server on stream start/stop.
// Nginx-RTMP sends application/x-www-form-urlencoded with fields: call, name, app, addr
func (h *StreamHandler) IngestCallback(c *fiber.Ctx) error {
	// "call" = "publish" on start, "publish_done" on stop
	call := c.FormValue("call")
	name := c.FormValue("name") // stream key

	if name == "" {
		return c.SendStatus(fiber.StatusOK) // ignore malformed calls
	}

	switch call {
	case "publish":
		hlsPath := fmt.Sprintf("/tmp/hls/%s", name)
		h.db.Exec(context.Background(),
			`UPDATE events SET status='live', hls_path=$1, updated_at=NOW() WHERE stream_key=$2`,
			hlsPath, name,
		)
	case "publish_done":
		h.db.Exec(context.Background(),
			`UPDATE events SET status='completed', updated_at=NOW() WHERE stream_key=$1`,
			name,
		)
	}

	return c.SendStatus(fiber.StatusOK)
}
