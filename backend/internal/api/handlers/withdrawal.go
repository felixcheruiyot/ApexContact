package handlers

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/livestreamify/backend/internal/domain"
	"github.com/livestreamify/backend/internal/service"
)

type WithdrawalHandler struct {
	db  *pgxpool.Pool
	svc *service.WithdrawalService
}

func NewWithdrawalHandler(db *pgxpool.Pool, svc *service.WithdrawalService) *WithdrawalHandler {
	return &WithdrawalHandler{db: db, svc: svc}
}

// Balance returns the caller's available promoter balance.
//
// GET /api/v1/withdrawals/balance
func (h *WithdrawalHandler) Balance(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	balance, err := h.svc.AvailableBalance(c.Context(), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "could not calculate balance")
	}

	return c.JSON(domain.Response{
		Data: fiber.Map{
			"available_balance": balance,
			"currency":          "KES",
		},
	})
}

// GetPayoutAccount returns the caller's registered payout account.
//
// GET /api/v1/withdrawals/payout-account
func (h *WithdrawalHandler) GetPayoutAccount(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	acct, err := h.svc.GetPayoutAccount(c.Context(), userID)
	if errors.Is(err, service.ErrAccountNotFound) {
		return fiber.NewError(fiber.StatusNotFound, "no payout account registered")
	}
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "could not fetch payout account")
	}

	return c.JSON(domain.Response{Data: acct})
}

// SetPayoutAccount registers a payout account for the caller.
//
// POST /api/v1/withdrawals/payout-account
func (h *WithdrawalHandler) SetPayoutAccount(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var body struct {
		AccountType   string `json:"account_type"`
		AccountNumber string `json:"account_number"`
		AccountName   string `json:"account_name"`
		BankName      string `json:"bank_name"`
	}
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if body.AccountType == "" || body.AccountNumber == "" {
		return fiber.NewError(fiber.StatusBadRequest, "account_type and account_number are required")
	}
	if body.AccountType != "mpesa" && body.AccountType != "bank" {
		return fiber.NewError(fiber.StatusBadRequest, "account_type must be 'mpesa' or 'bank'")
	}

	acct, err := h.svc.SetPayoutAccount(c.Context(), userID,
		body.AccountType, body.AccountNumber, body.AccountName, body.BankName)
	if errors.Is(err, service.ErrAccountExists) {
		return fiber.NewError(fiber.StatusConflict, "payout account already exists — contact support to change it")
	}
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "could not save payout account")
	}

	return c.Status(fiber.StatusCreated).JSON(domain.Response{Data: acct})
}

// Initiate creates a pending withdrawal and sends an OTP to the user's email.
//
// POST /api/v1/withdrawals/initiate
func (h *WithdrawalHandler) Initiate(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var body struct {
		Amount float64 `json:"amount"`
	}
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if body.Amount <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "amount must be greater than zero")
	}

	user, err := h.fetchUser(c.Context(), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "could not load user")
	}

	w, err := h.svc.Initiate(c.Context(), user, body.Amount)
	switch {
	case errors.Is(err, service.ErrEmailNotVerified):
		return fiber.NewError(fiber.StatusForbidden, err.Error())
	case errors.Is(err, service.ErrAccountNotFound):
		return fiber.NewError(fiber.StatusBadRequest, "no payout account registered — add one first")
	case errors.Is(err, service.ErrBelowMinimum):
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	case errors.Is(err, service.ErrInsufficientBalance):
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	case errors.Is(err, service.ErrActiveWithdrawal):
		return fiber.NewError(fiber.StatusConflict, err.Error())
	case err != nil:
		return fiber.NewError(fiber.StatusInternalServerError, "could not initiate withdrawal")
	}

	return c.Status(fiber.StatusCreated).JSON(domain.Response{
		Data: fiber.Map{
			"withdrawal_id": w.ID,
			"message":       fmt.Sprintf("An OTP has been sent to %s. Enter it to confirm the withdrawal.", user.Email),
		},
	})
}

// Confirm validates the OTP and executes the payout.
//
// POST /api/v1/withdrawals/confirm
func (h *WithdrawalHandler) Confirm(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var body struct {
		WithdrawalID string `json:"withdrawal_id"`
		OTP          string `json:"otp"`
	}
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if body.WithdrawalID == "" || body.OTP == "" {
		return fiber.NewError(fiber.StatusBadRequest, "withdrawal_id and otp are required")
	}

	w, err := h.svc.Confirm(c.Context(), userID, body.WithdrawalID, body.OTP)
	switch {
	case errors.Is(err, service.ErrWithdrawalNotFound):
		return fiber.NewError(fiber.StatusNotFound, "withdrawal not found or already confirmed")
	case errors.Is(err, service.ErrOTPExpired):
		return fiber.NewError(fiber.StatusGone, err.Error())
	case errors.Is(err, service.ErrOTPInvalid):
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	case err != nil:
		return fiber.NewError(fiber.StatusInternalServerError, "could not confirm withdrawal")
	}

	return c.JSON(domain.Response{Data: w})
}

// History lists all withdrawals for the authenticated user.
//
// GET /api/v1/withdrawals/history
func (h *WithdrawalHandler) History(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	withdrawals, err := h.svc.History(c.Context(), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "could not fetch withdrawal history")
	}

	return c.JSON(domain.Response{Data: withdrawals})
}

// fetchUser loads a domain.User from the database by UUID string.
func (h *WithdrawalHandler) fetchUser(ctx context.Context, userID string) (*domain.User, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	var u domain.User
	err = h.db.QueryRow(ctx,
		`SELECT id, email, full_name, phone, role, is_locked, email_verified, created_at, updated_at
		 FROM users WHERE id = $1`,
		uid,
	).Scan(&u.ID, &u.Email, &u.FullName, &u.Phone, &u.Role, &u.IsLocked, &u.EmailVerified, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
