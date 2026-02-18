package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/apexcontact/backend/internal/config"
	"github.com/apexcontact/backend/internal/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type PaymentHandler struct {
	cfg *config.Config
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewPaymentHandler(cfg *config.Config, db *pgxpool.Pool, rdb *redis.Client) *PaymentHandler {
	return &PaymentHandler{cfg: cfg, db: db, rdb: rdb}
}

type initiatePaymentRequest struct {
	EventID     string `json:"event_id"`
	PhoneNumber string `json:"phone_number"` // E.164 format e.g. +254712345678
}

// Initiate triggers the IntaSend M-Pesa STK Push.
func (h *PaymentHandler) Initiate(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req initiatePaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
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

	// TODO: Call IntaSend API to initiate STK Push
	// intasendRef, err := h.intasendClient.StkPush(req.PhoneNumber, event.Price, event.Currency, payment.ID.String())
	// For now, store placeholder ref
	intasendRef := fmt.Sprintf("PENDING-%s", payment.ID.String()[:8])
	h.db.Exec(context.Background(),
		`UPDATE payments SET intasend_ref=$1 WHERE id=$2`, intasendRef, payment.ID,
	)

	// Cache payment status in Redis for fast polling
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
func (h *PaymentHandler) Callback(c *fiber.Ctx) error {
	var body struct {
		InvoiceID string `json:"invoice_id"`
		State     string `json:"state"` // "COMPLETE" | "FAILED" | "CANCELLED"
		Value     string `json:"value"`
		Account   string `json:"account"`
		APIRef    string `json:"api_ref"` // our payment UUID
	}
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid callback body")
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
		return c.SendStatus(fiber.StatusOK) // ignore unknown states
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
