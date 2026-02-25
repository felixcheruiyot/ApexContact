package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/livestreamify/backend/internal/config"
	"github.com/livestreamify/backend/internal/domain"
	"github.com/livestreamify/backend/internal/integrations/intasend"
	"github.com/livestreamify/backend/internal/service"
	"github.com/redis/go-redis/v9"
)

type WithdrawalHandler struct {
	svc *service.WithdrawalService
	db  *pgxpool.Pool
}

func NewWithdrawalHandler(cfg *config.Config, db *pgxpool.Pool, rdb *redis.Client, notifSvc *service.NotificationService) *WithdrawalHandler {
	is := intasend.New(cfg.IntaSendPrivateKey, cfg.IntaSendBaseURL)
	svc := service.NewWithdrawalService(db, rdb, is, notifSvc)
	return &WithdrawalHandler{svc: svc, db: db}
}

// GET /api/v1/profile/balance
func (h *WithdrawalHandler) Balance(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	balance, err := h.svc.AvailableBalance(c.Context(), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch balance")
	}
	upcoming, err := h.svc.UpcomingBalance(c.Context(), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch balance")
	}
	return c.JSON(domain.Response{Data: fiber.Map{
		"balance":          balance,
		"upcoming_balance": upcoming,
		"currency":         "KES",
	}})
}

// GET /api/v1/profile/payout-account
func (h *WithdrawalHandler) GetPayoutAccount(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	account, err := h.svc.GetPayoutAccount(c.Context(), userID)
	if errors.Is(err, service.ErrAccountNotFound) {
		return c.JSON(domain.Response{Data: nil})
	}
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch payout account")
	}
	return c.JSON(domain.Response{Data: account})
}

type setPayoutAccountRequest struct {
	AccountType   string `json:"account_type"`
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
	BankName      string `json:"bank_name"`
}

// POST /api/v1/profile/payout-account
func (h *WithdrawalHandler) SetPayoutAccount(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req setPayoutAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if req.AccountType != "mpesa" && req.AccountType != "bank" {
		return fiber.NewError(fiber.StatusBadRequest, "account_type must be 'mpesa' or 'bank'")
	}
	if req.AccountNumber == "" || req.AccountName == "" {
		return fiber.NewError(fiber.StatusBadRequest, "account_number and account_name are required")
	}

	account, err := h.svc.SetPayoutAccount(c.Context(), userID, req.AccountType, req.AccountNumber, req.AccountName, req.BankName)
	if errors.Is(err, service.ErrAccountExists) {
		return fiber.NewError(fiber.StatusConflict, err.Error())
	}
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to set payout account")
	}
	return c.Status(fiber.StatusCreated).JSON(domain.Response{Data: account})
}

type initiateWithdrawalRequest struct {
	Amount float64 `json:"amount"`
}

// POST /api/v1/profile/withdrawals
func (h *WithdrawalHandler) Initiate(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req initiateWithdrawalRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	// Fetch user from DB
	var user domain.User
	err := h.db.QueryRow(c.Context(),
		`SELECT id, email, full_name, phone, role, is_locked, email_verified, created_at, updated_at
		 FROM users WHERE id = $1`, userID,
	).Scan(&user.ID, &user.Email, &user.FullName, &user.Phone, &user.Role,
		&user.IsLocked, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	withdrawal, err := h.svc.Initiate(c.Context(), &user, req.Amount)
	if errors.Is(err, service.ErrEmailNotVerified) {
		return fiber.NewError(fiber.StatusForbidden, err.Error())
	}
	if errors.Is(err, service.ErrAccountNotFound) {
		return fiber.NewError(fiber.StatusBadRequest, "set up a payout account first")
	}
	if errors.Is(err, service.ErrBelowMinimum) || errors.Is(err, service.ErrInsufficientBalance) || errors.Is(err, service.ErrActiveWithdrawal) {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to initiate withdrawal")
	}
	return c.Status(fiber.StatusCreated).JSON(domain.Response{Data: withdrawal})
}

type confirmWithdrawalRequest struct {
	OTP string `json:"otp"`
}

// POST /api/v1/profile/withdrawals/:id/confirm
func (h *WithdrawalHandler) Confirm(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	withdrawalID := c.Params("id")

	var req confirmWithdrawalRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	withdrawal, err := h.svc.Confirm(c.Context(), userID, withdrawalID, req.OTP)
	if errors.Is(err, service.ErrWithdrawalNotFound) {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	if errors.Is(err, service.ErrOTPExpired) || errors.Is(err, service.ErrOTPInvalid) {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to confirm withdrawal")
	}
	return c.JSON(domain.Response{Data: withdrawal})
}

// GET /api/v1/profile/withdrawals
func (h *WithdrawalHandler) History(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	withdrawals, err := h.svc.History(c.Context(), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch withdrawal history")
	}
	if withdrawals == nil {
		withdrawals = []domain.Withdrawal{}
	}
	return c.JSON(domain.Response{Data: withdrawals})
}
