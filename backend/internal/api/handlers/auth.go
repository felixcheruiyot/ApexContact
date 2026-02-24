package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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

// GoogleCallback handles POST /api/v1/auth/google/callback.
// The SPA sends the authorization code it received from Google; this handler
// exchanges it for an access token, fetches the user's profile, upserts the
// user in the DB, and returns an application JWT.
func (h *AuthHandler) GoogleCallback(c *fiber.Ctx) error {
	if h.cfg.GoogleClientID == "" || h.cfg.GoogleClientSecret == "" {
		return fiber.NewError(fiber.StatusNotImplemented, "Google OAuth is not configured on this server")
	}

	var req struct {
		Code        string `json:"code"`
		RedirectURI string `json:"redirect_uri"`
	}
	if err := c.BodyParser(&req); err != nil || req.Code == "" || req.RedirectURI == "" {
		return fiber.NewError(fiber.StatusBadRequest, "code and redirect_uri are required")
	}

	// Exchange authorization code → access token
	accessToken, err := exchangeGoogleCode(req.Code, req.RedirectURI, h.cfg.GoogleClientID, h.cfg.GoogleClientSecret)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "failed to exchange Google authorization code")
	}

	// Fetch verified profile from Google
	gInfo, err := fetchGoogleUserInfo(accessToken)
	if err != nil || gInfo.Email == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "failed to retrieve Google account info")
	}

	// Upsert: create account on first sign-in, or update email_verified on subsequent ones.
	// password_hash is set to a non-bcrypt placeholder so the email/password login path
	// always rejects Google-only accounts gracefully.
	var user domain.User
	err = h.db.QueryRow(context.Background(), `
		INSERT INTO users (id, email, password_hash, full_name, phone, role, email_verified, created_at, updated_at)
		VALUES ($1, $2, 'GOOGLE_OAUTH', $3, '', $4, true, NOW(), NOW())
		ON CONFLICT (email) DO UPDATE
			SET email_verified = true,
			    full_name = CASE WHEN users.full_name = '' THEN EXCLUDED.full_name ELSE users.full_name END,
			    updated_at = NOW()
		RETURNING id, email, password_hash,
		          COALESCE(full_name, ''), COALESCE(phone, ''),
		          role, is_locked, email_verified, created_at, updated_at`,
		uuid.New(), gInfo.Email, gInfo.Name, domain.RoleMember,
	).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FullName, &user.Phone,
		&user.Role, &user.IsLocked, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to sign in with Google")
	}

	if user.IsLocked {
		return fiber.NewError(fiber.StatusForbidden, "account suspended")
	}

	token, err := h.generateToken(&user)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to generate token")
	}

	return c.JSON(domain.Response{
		Data: fiber.Map{"token": token, "user": user},
	})
}

// ── Google OAuth helpers ───────────────────────────────────────────────────────

type googleTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type googleUserInfo struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
}

func exchangeGoogleCode(code, redirectURI, clientID, clientSecret string) (string, error) {
	resp, err := http.PostForm("https://oauth2.googleapis.com/token", url.Values{
		"code":          {code},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"redirect_uri":  {redirectURI},
		"grant_type":    {"authorization_code"},
	})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("google token endpoint returned %d", resp.StatusCode)
	}

	var t googleTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return "", err
	}
	if t.AccessToken == "" {
		return "", fmt.Errorf("empty access token from Google")
	}
	return t.AccessToken, nil
}

func fetchGoogleUserInfo(accessToken string) (*googleUserInfo, error) {
	req, err := http.NewRequest(http.MethodGet, "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google userinfo returned %d", resp.StatusCode)
	}

	var info googleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}
