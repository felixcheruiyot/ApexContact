package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/livestreamify/backend/internal/config"
	"github.com/livestreamify/backend/internal/domain"
	"github.com/livestreamify/backend/internal/integrations/intasend"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type PaymentHandler struct {
	cfg      *config.Config
	db       *pgxpool.Pool
	rdb      *redis.Client
	intasend *intasend.Client
}

func NewPaymentHandler(cfg *config.Config, db *pgxpool.Pool, rdb *redis.Client) *PaymentHandler {
	return &PaymentHandler{
		cfg:      cfg,
		db:       db,
		rdb:      rdb,
		intasend: intasend.New(cfg.IntaSendPrivateKey, cfg.IntaSendBaseURL),
	}
}

type initiatePaymentRequest struct {
	EventID     string `json:"event_id"`
	PhoneNumber string `json:"phone_number"` // E.164 format e.g. 254712345678
}

// Initiate triggers the IntaSend M-Pesa STK Push.
func (h *PaymentHandler) Initiate(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req initiatePaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if req.PhoneNumber == "" {
		return fiber.NewError(fiber.StatusBadRequest, "phone_number is required")
	}

	// Fetch event price
	var event domain.Event
	err := h.db.QueryRow(context.Background(),
		`SELECT id, price, currency, title FROM events WHERE id=$1 AND status IN ('scheduled','live')`,
		req.EventID,
	).Scan(&event.ID, &event.Price, &event.Currency, &event.Title)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "event not found or not available")
	}

	payment := domain.Payment{
		ID:          uuid.New(),
		UserID:      uuid.MustParse(userID),
		EventID:     event.ID,
		Amount:      event.Price,
		Currency:    event.Currency,
		Status:      domain.PaymentPending,
		PhoneNumber: req.PhoneNumber,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err = h.db.Exec(context.Background(),
		`INSERT INTO payments (id, user_id, event_id, amount, currency, status, phone_number, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		payment.ID, payment.UserID, payment.EventID, payment.Amount, payment.Currency,
		payment.Status, payment.PhoneNumber, payment.CreatedAt, payment.UpdatedAt,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to create payment record")
	}

	// Initiate IntaSend STK Push — payment.ID is used as api_ref so the
	// webhook callback can look up and update the right payment record.
	invoiceID, err := h.intasend.StkPush(req.PhoneNumber, payment.Amount, payment.Currency, payment.ID.String())
	if err != nil {
		log.Error().Err(err).Str("payment_id", payment.ID.String()).Msg("intasend stk push failed")
		// Roll back the pending payment so the user can retry.
		h.db.Exec(context.Background(), `DELETE FROM payments WHERE id=$1`, payment.ID)
		return fiber.NewError(fiber.StatusBadGateway, "failed to initiate M-Pesa payment — please try again")
	}

	h.db.Exec(context.Background(),
		`UPDATE payments SET intasend_ref=$1 WHERE id=$2`, invoiceID, payment.ID,
	)

	// Cache payment status in Redis for fast polling (30-min TTL matches STK timeout).
	h.rdb.Set(context.Background(),
		fmt.Sprintf("payment:%s", payment.ID.String()),
		"pending",
		30*time.Minute,
	)

	return c.Status(fiber.StatusCreated).JSON(domain.Response{
		Data: fiber.Map{
			"payment_id":   payment.ID,
			"amount":       payment.Amount,
			"currency":     payment.Currency,
			"phone_number": payment.PhoneNumber,
			"message":      "STK Push sent to your phone. Enter your M-Pesa PIN to complete payment.",
		},
	})
}

// Callback is the IntaSend webhook handler — called when payment status changes.
// IntaSend posts: invoice_id, state (COMPLETE|FAILED|CANCELLED), value, account, api_ref
func (h *PaymentHandler) Callback(c *fiber.Ctx) error {
	var body struct {
		InvoiceID string `json:"invoice_id"`
		State     string `json:"state"`
		Value     string `json:"value"`
		Account   string `json:"account"`
		APIRef    string `json:"api_ref"` // our payment UUID
		Currency  string `json:"currency"`
		NetAmount string `json:"net_amount"`
	}
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid callback body")
	}

	if body.APIRef == "" || body.InvoiceID == "" {
		// Not a payment callback we recognise — acknowledge to stop retries.
		return c.SendStatus(fiber.StatusOK)
	}

	var status domain.PaymentStatus
	switch body.State {
	case "COMPLETE":
		status = domain.PaymentSuccess
	case "FAILED":
		status = domain.PaymentFailed
	case "CANCELLED":
		status = domain.PaymentCancelled
	default:
		return c.SendStatus(fiber.StatusOK) // pending/processing states — ignore
	}

	h.db.Exec(context.Background(),
		`UPDATE payments SET status=$1, intasend_ref=$2, updated_at=NOW() WHERE id=$3`,
		status, body.InvoiceID, body.APIRef,
	)

	h.rdb.Set(context.Background(),
		fmt.Sprintf("payment:%s", body.APIRef),
		string(status),
		10*time.Minute,
	)

	log.Info().
		Str("payment_id", body.APIRef).
		Str("invoice_id", body.InvoiceID).
		Str("state", body.State).
		Msg("payment webhook received")

	return c.SendStatus(fiber.StatusOK)
}

// Status allows the frontend to poll payment progress.
func (h *PaymentHandler) Status(c *fiber.Ctx) error {
	paymentID := c.Params("id")
	userID := c.Locals("user_id").(string)

	// Try Redis cache first
	cached, err := h.rdb.Get(context.Background(), fmt.Sprintf("payment:%s", paymentID)).Result()
	if err == nil {
		return c.JSON(domain.Response{Data: fiber.Map{"status": cached}})
	}

	// Fallback to DB
	var status domain.PaymentStatus
	err = h.db.QueryRow(context.Background(),
		`SELECT status FROM payments WHERE id=$1 AND user_id=$2`, paymentID, userID,
	).Scan(&status)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "payment not found")
	}

	return c.JSON(domain.Response{Data: fiber.Map{"status": status}})
}
