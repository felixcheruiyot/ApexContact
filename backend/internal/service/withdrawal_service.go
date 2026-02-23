package service

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/livestreamify/backend/internal/domain"
	"github.com/livestreamify/backend/internal/integrations/intasend"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

const (
	PromoterCut          = 0.70
	MinWithdrawalAmount  = 100.0
	withdrawalOTPTTL     = 10 * time.Minute
)

var (
	ErrAccountExists        = errors.New("payout account already exists — contact support to change it")
	ErrAccountNotFound      = errors.New("payout account not found")
	ErrWithdrawalNotFound   = errors.New("withdrawal not found")
	ErrOTPExpired           = errors.New("OTP expired, please request a new withdrawal")
	ErrOTPInvalid           = errors.New("invalid OTP")
	ErrActiveWithdrawal     = errors.New("you already have an active withdrawal pending confirmation")
	ErrInsufficientBalance  = errors.New("insufficient available balance")
	ErrBelowMinimum         = errors.New("amount is below the minimum withdrawal of KES 100")
	ErrEmailNotVerified     = errors.New("email address must be verified before withdrawing funds")
)

// WithdrawalService manages payout accounts, balances, and the OTP-gated
// withdrawal flow.
type WithdrawalService struct {
	db       *pgxpool.Pool
	rdb      *redis.Client
	intasend *intasend.Client
	notifSvc *NotificationService
}

func NewWithdrawalService(
	db *pgxpool.Pool,
	rdb *redis.Client,
	is *intasend.Client,
	notifSvc *NotificationService,
) *WithdrawalService {
	return &WithdrawalService{db: db, rdb: rdb, intasend: is, notifSvc: notifSvc}
}

// AvailableBalance returns 70% of total successful ticket sales for events
// owned by userID, minus the sum of all completed withdrawals.
func (s *WithdrawalService) AvailableBalance(ctx context.Context, userID string) (float64, error) {
	const query = `
		SELECT
		  COALESCE((
		    SELECT SUM(p.amount)
		    FROM events e
		    JOIN subscriptions sub ON sub.event_id = e.id
		    JOIN payments p ON sub.payment_id = p.id
		    WHERE e.promoter_id = $1 AND p.status = 'success'
		  ), 0) * $2
		  - COALESCE((
		    SELECT SUM(amount)
		    FROM withdrawals
		    WHERE user_id = $1 AND status = 'completed'
		  ), 0)
	`
	var balance float64
	err := s.db.QueryRow(ctx, query, userID, PromoterCut).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("AvailableBalance: %w", err)
	}
	return balance, nil
}

// GetPayoutAccount returns the user's registered payout account or
// ErrAccountNotFound if none has been set up.
func (s *WithdrawalService) GetPayoutAccount(ctx context.Context, userID string) (*domain.PayoutAccount, error) {
	var a domain.PayoutAccount
	err := s.db.QueryRow(ctx,
		`SELECT id, user_id, account_type, account_number, account_name, bank_name, created_at
		 FROM payout_accounts WHERE user_id = $1`,
		userID,
	).Scan(&a.ID, &a.UserID, &a.AccountType, &a.AccountNumber, &a.AccountName, &a.BankName, &a.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrAccountNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("GetPayoutAccount: %w", err)
	}
	return &a, nil
}

// SetPayoutAccount registers a payout account for the user.
// Returns ErrAccountExists if the user already has one (UNIQUE constraint).
func (s *WithdrawalService) SetPayoutAccount(
	ctx context.Context,
	userID, accountType, accountNumber, accountName, bankName string,
) (*domain.PayoutAccount, error) {
	// Pre-check in application layer before hitting the DB constraint.
	existing, err := s.GetPayoutAccount(ctx, userID)
	if err == nil && existing != nil {
		return nil, ErrAccountExists
	}
	if err != nil && !errors.Is(err, ErrAccountNotFound) {
		return nil, err
	}

	var a domain.PayoutAccount
	err = s.db.QueryRow(ctx,
		`INSERT INTO payout_accounts (user_id, account_type, account_number, account_name, bank_name)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, user_id, account_type, account_number, account_name, bank_name, created_at`,
		userID, accountType, accountNumber, accountName, bankName,
	).Scan(&a.ID, &a.UserID, &a.AccountType, &a.AccountNumber, &a.AccountName, &a.BankName, &a.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("SetPayoutAccount: %w", err)
	}
	return &a, nil
}

// Initiate validates preconditions and creates a pending_otp withdrawal,
// then sends an OTP email to the user.
func (s *WithdrawalService) Initiate(ctx context.Context, user *domain.User, amount float64) (*domain.Withdrawal, error) {
	if !user.EmailVerified {
		return nil, ErrEmailNotVerified
	}

	account, err := s.GetPayoutAccount(ctx, user.ID.String())
	if err != nil {
		return nil, err // ErrAccountNotFound or other
	}

	if amount < MinWithdrawalAmount {
		return nil, ErrBelowMinimum
	}

	balance, err := s.AvailableBalance(ctx, user.ID.String())
	if err != nil {
		return nil, err
	}
	if amount > balance {
		return nil, ErrInsufficientBalance
	}

	// Reject if there is already a pending or processing withdrawal.
	var activeCount int
	err = s.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM withdrawals
		 WHERE user_id = $1 AND status IN ('pending_otp','processing')`,
		user.ID,
	).Scan(&activeCount)
	if err != nil {
		return nil, fmt.Errorf("Initiate: check active withdrawals: %w", err)
	}
	if activeCount > 0 {
		return nil, ErrActiveWithdrawal
	}

	// Create the withdrawal row.
	var w domain.Withdrawal
	err = s.db.QueryRow(ctx,
		`INSERT INTO withdrawals (user_id, payout_account_id, amount, currency)
		 VALUES ($1, $2, $3, 'KES')
		 RETURNING id, user_id, payout_account_id, amount, currency, status,
		           intasend_ref, failure_reason, created_at, updated_at`,
		user.ID, account.ID, amount,
	).Scan(
		&w.ID, &w.UserID, &w.PayoutAccountID, &w.Amount, &w.Currency, &w.Status,
		&w.IntaSendRef, &w.FailureReason, &w.CreatedAt, &w.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("Initiate: insert withdrawal: %w", err)
	}

	otp, err := generateWithdrawalOTP()
	if err != nil {
		return nil, fmt.Errorf("Initiate: generate OTP: %w", err)
	}

	otpKey := withdrawalOTPKey(w.ID.String())
	if err := s.rdb.Set(ctx, otpKey, otp, withdrawalOTPTTL).Err(); err != nil {
		return nil, fmt.Errorf("Initiate: store OTP: %w", err)
	}

	amountStr := fmt.Sprintf("%.2f", amount)
	if err := s.notifSvc.SendWithdrawalOTP(ctx, user, otp, amountStr, "KES"); err != nil {
		log.Warn().Err(err).Str("withdrawal_id", w.ID.String()).Msg("failed to send withdrawal OTP email")
		// Non-fatal: the OTP is in Redis; user can retry the confirm step.
	}

	return &w, nil
}

// Confirm validates the OTP and dispatches the payout via IntaSend.
func (s *WithdrawalService) Confirm(ctx context.Context, userID, withdrawalID, otp string) (*domain.Withdrawal, error) {
	// Fetch the withdrawal — must be pending_otp and owned by this user.
	var w domain.Withdrawal
	var accountID uuid.UUID
	err := s.db.QueryRow(ctx,
		`SELECT w.id, w.user_id, w.payout_account_id, w.amount, w.currency, w.status,
		        w.intasend_ref, w.failure_reason, w.created_at, w.updated_at
		 FROM withdrawals w
		 WHERE w.id = $1 AND w.user_id = $2 AND w.status = 'pending_otp'`,
		withdrawalID, userID,
	).Scan(
		&w.ID, &w.UserID, &accountID, &w.Amount, &w.Currency, &w.Status,
		&w.IntaSendRef, &w.FailureReason, &w.CreatedAt, &w.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrWithdrawalNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("Confirm: fetch withdrawal: %w", err)
	}
	w.PayoutAccountID = accountID

	// Validate OTP from Redis.
	otpKey := withdrawalOTPKey(w.ID.String())
	storedOTP, err := s.rdb.Get(ctx, otpKey).Result()
	if err != nil {
		// Key missing = expired.
		return nil, ErrOTPExpired
	}

	if subtle.ConstantTimeCompare([]byte(storedOTP), []byte(otp)) != 1 {
		return nil, ErrOTPInvalid
	}

	// Mark as processing and delete the OTP key.
	_, err = s.db.Exec(ctx,
		`UPDATE withdrawals SET status = 'processing', updated_at = NOW() WHERE id = $1`,
		w.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("Confirm: mark processing: %w", err)
	}
	s.rdb.Del(ctx, otpKey)

	// Fetch payout account details for the IntaSend call.
	var acct domain.PayoutAccount
	err = s.db.QueryRow(ctx,
		`SELECT account_name, account_number FROM payout_accounts WHERE id = $1`,
		accountID,
	).Scan(&acct.AccountName, &acct.AccountNumber)
	if err != nil {
		return nil, fmt.Errorf("Confirm: fetch payout account: %w", err)
	}

	narrative := fmt.Sprintf("Live Streamify promoter payout - withdrawal %s", w.ID.String())
	trackingID, sendErr := s.intasend.SendMoney(
		acct.AccountName, acct.AccountNumber,
		w.Amount, w.Currency, narrative,
	)

	if sendErr != nil {
		failReason := sendErr.Error()
		log.Error().Err(sendErr).Str("withdrawal_id", w.ID.String()).Msg("intasend SendMoney failed")

		s.db.Exec(ctx,
			`UPDATE withdrawals SET status = 'failed', failure_reason = $1, updated_at = NOW() WHERE id = $2`,
			failReason, w.ID,
		)
		w.Status = domain.WithdrawalFailed
		w.FailureReason = failReason

		// Send failure notification — fetch user email.
		s.sendFailureEmail(ctx, userID, w.Amount, w.Currency)

		return &w, nil
	}

	// Success — mark completed.
	s.db.Exec(ctx,
		`UPDATE withdrawals SET status = 'completed', intasend_ref = $1, updated_at = NOW() WHERE id = $2`,
		trackingID, w.ID,
	)
	w.Status = domain.WithdrawalCompleted
	w.IntaSendRef = trackingID

	return &w, nil
}

// History returns all withdrawals for the given user, newest first.
func (s *WithdrawalService) History(ctx context.Context, userID string) ([]domain.Withdrawal, error) {
	rows, err := s.db.Query(ctx,
		`SELECT id, user_id, payout_account_id, amount, currency, status,
		        intasend_ref, failure_reason, created_at, updated_at
		 FROM withdrawals
		 WHERE user_id = $1
		 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("History: %w", err)
	}
	defer rows.Close()

	var results []domain.Withdrawal
	for rows.Next() {
		var w domain.Withdrawal
		if err := rows.Scan(
			&w.ID, &w.UserID, &w.PayoutAccountID, &w.Amount, &w.Currency, &w.Status,
			&w.IntaSendRef, &w.FailureReason, &w.CreatedAt, &w.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("History: scan: %w", err)
		}
		results = append(results, w)
	}
	return results, rows.Err()
}

// sendFailureEmail fetches the user and sends a withdrawal failure notice.
func (s *WithdrawalService) sendFailureEmail(ctx context.Context, userID string, amount float64, currency string) {
	var email, fullName string
	err := s.db.QueryRow(ctx,
		`SELECT email, full_name FROM users WHERE id = $1`, userID,
	).Scan(&email, &fullName)
	if err != nil {
		log.Warn().Err(err).Str("user_id", userID).Msg("withdrawal failure email: could not fetch user")
		return
	}

	body := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head><meta charset="UTF-8"></head>
<body style="margin:0;padding:0;background-color:#0A0A0F;font-family:Arial,sans-serif;">
  <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" style="background-color:#0A0A0F;">
    <tr><td align="center" style="padding:40px 16px;">
      <table role="presentation" width="600" cellpadding="0" cellspacing="0"
             style="max-width:600px;background-color:#141418;border-radius:12px;overflow:hidden;">
        <tr><td style="background-color:#E8002D;padding:20px 40px;">
          <p style="margin:0;font-size:22px;font-weight:700;letter-spacing:3px;color:#FFFFFF;text-transform:uppercase;">LIVE STREAMIFY</p>
        </td></tr>
        <tr><td style="padding:40px;">
          <h1 style="margin:0 0 16px;font-size:24px;font-weight:700;color:#FFFFFF;">Withdrawal Could Not Be Processed</h1>
          <p style="margin:0 0 8px;font-size:15px;color:#A0A0B0;line-height:1.6;">Hi %s,</p>
          <p style="margin:0 0 24px;font-size:15px;color:#A0A0B0;line-height:1.6;">
            Unfortunately, your withdrawal of <strong style="color:#FFFFFF;">%s %.2f</strong> could not be
            processed at this time due to a payment provider error.
          </p>
          <p style="margin:0;font-size:15px;color:#A0A0B0;line-height:1.6;">
            Please contact our support team at
            <a href="mailto:support@livestreamify.com" style="color:#E8002D;">support@livestreamify.com</a>
            and we will resolve this as quickly as possible.
          </p>
        </td></tr>
        <tr><td style="padding:24px 40px;border-top:1px solid #1E1E26;">
          <p style="margin:0;font-size:12px;color:#606070;text-align:center;">&copy; 2026 Live Streamify</p>
        </td></tr>
      </table>
    </td></tr>
  </table>
</body>
</html>`, fullName, currency, amount)

	if err := s.notifSvc.SendRaw(email, "Withdrawal Failed — Please Contact Support", body); err != nil {
		log.Warn().Err(err).Str("user_id", userID).Msg("failed to send withdrawal failure email")
	}
}

func withdrawalOTPKey(withdrawalID string) string {
	return "withdraw_otp:" + withdrawalID
}

func generateWithdrawalOTP() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}
