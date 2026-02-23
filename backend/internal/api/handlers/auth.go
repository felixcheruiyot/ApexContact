package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/livestreamify/backend/internal/api/middleware"
	"github.com/livestreamify/backend/internal/config"
	"github.com/livestreamify/backend/internal/domain"
	"github.com/livestreamify/backend/internal/service"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	cfg      *config.Config
	db       *pgxpool.Pool
	rdb      *redis.Client
	notifSvc *service.NotificationService
}

func NewAuthHandler(cfg *config.Config, db *pgxpool.Pool, rdb *redis.Client, notifSvc *service.NotificationService) *AuthHandler {
	return &AuthHandler{cfg: cfg, db: db, rdb: rdb, notifSvc: notifSvc}
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req registerRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to hash password")
	}

	verificationToken, err := generateVerificationToken()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to generate verification token")
	}
	tokenExpiresAt := time.Now().Add(24 * time.Hour)

	user := &domain.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: string(hash),
		FullName:     req.FullName,
		Phone:        req.Phone,
		Role:         domain.RoleMember,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err = h.db.Exec(context.Background(),
		`INSERT INTO users (id, email, password_hash, full_name, phone, role, created_at, updated_at,
		                    verification_token, verification_token_expires_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		user.ID, user.Email, user.PasswordHash, user.FullName, user.Phone, user.Role,
		user.CreatedAt, user.UpdatedAt, verificationToken, tokenExpiresAt,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusConflict, "email already registered")
	}

	// Send verification email asynchronously — do not block registration on email delivery.
	if h.notifSvc != nil {
		go func() {
			_ = h.notifSvc.SendVerificationEmail(context.Background(), user, verificationToken)
		}()
	}

	token, err := h.generateToken(user)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to generate token")
	}

	return c.Status(fiber.StatusCreated).JSON(domain.Response{
		Data: fiber.Map{"token": token, "user": user},
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req loginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	var user domain.User
	err := h.db.QueryRow(context.Background(),
		`SELECT id, email, password_hash, full_name, phone, role, is_locked, email_verified, created_at, updated_at
		 FROM users WHERE email = $1`, req.Email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FullName, &user.Phone, &user.Role,
		&user.IsLocked, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
	}

	if user.IsLocked {
		return fiber.NewError(fiber.StatusForbidden, "account suspended")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
	}

	token, err := h.generateToken(&user)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to generate token")
	}

	return c.JSON(domain.Response{
		Data: fiber.Map{"token": token, "user": user},
	})
}

// VerifyEmail handles GET /api/v1/auth/verify-email?token=<token>.
func (h *AuthHandler) VerifyEmail(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		return fiber.NewError(fiber.StatusBadRequest, "missing token")
	}

	var userID string
	err := h.db.QueryRow(context.Background(),
		`UPDATE users
		 SET email_verified = true,
		     verification_token = NULL,
		     verification_token_expires_at = NULL
		 WHERE verification_token = $1
		   AND verification_token_expires_at > NOW()
		 RETURNING id`,
		token,
	).Scan(&userID)

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid or expired verification token")
	}

	return c.JSON(domain.Response{Data: fiber.Map{"message": "email verified successfully"}})
}

// ResendVerification handles POST /api/v1/auth/resend-verification (auth required).
// Rate-limited to one send per 5 minutes per user.
func (h *AuthHandler) ResendVerification(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	rateLimitKey := "resend_verify:" + userID
	if h.rdb.Exists(context.Background(), rateLimitKey).Val() > 0 {
		return fiber.NewError(fiber.StatusTooManyRequests, "please wait before requesting another verification email")
	}

	var user domain.User
	err := h.db.QueryRow(context.Background(),
		`SELECT id, email, full_name, email_verified FROM users WHERE id = $1`, userID,
	).Scan(&user.ID, &user.Email, &user.FullName, &user.EmailVerified)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	if user.EmailVerified {
		return fiber.NewError(fiber.StatusBadRequest, "email already verified")
	}

	newToken, err := generateVerificationToken()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to generate token")
	}

	_, err = h.db.Exec(context.Background(),
		`UPDATE users SET verification_token = $1, verification_token_expires_at = $2 WHERE id = $3`,
		newToken, time.Now().Add(24*time.Hour), userID,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to update verification token")
	}

	if h.notifSvc != nil {
		go func() {
			_ = h.notifSvc.SendVerificationEmail(context.Background(), &user, newToken)
		}()
	}

	h.rdb.Set(context.Background(), rateLimitKey, "1", 5*time.Minute)

	return c.JSON(domain.Response{Data: fiber.Map{"message": "verification email sent"}})
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	// TODO: implement refresh token rotation
	return c.JSON(domain.Response{Data: "refresh not yet implemented"})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// Token invalidation via Redis blocklist
	userID := c.Locals("user_id").(string)
	h.rdb.Set(context.Background(), "blocklist:"+userID, "1", 24*time.Hour)
	return c.JSON(domain.Response{Data: "logged out"})
}

func (h *AuthHandler) generateToken(user *domain.User) (string, error) {
	claims := middleware.Claims{
		UserID: user.ID.String(),
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.cfg.JWTSecret))
}

// generateVerificationToken creates a cryptographically random 64-char hex token.
func generateVerificationToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
