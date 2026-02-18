package handlers

import (
	"context"
	"time"

	"github.com/livestreamify/backend/internal/api/middleware"
	"github.com/livestreamify/backend/internal/config"
	"github.com/livestreamify/backend/internal/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	cfg *config.Config
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewAuthHandler(cfg *config.Config, db *pgxpool.Pool, rdb *redis.Client) *AuthHandler {
	return &AuthHandler{cfg: cfg, db: db, rdb: rdb}
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

	user := &domain.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: string(hash),
		FullName:     req.FullName,
		Phone:        req.Phone,
		Role:         domain.RoleViewer,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err = h.db.Exec(context.Background(),
		`INSERT INTO users (id, email, password_hash, full_name, phone, role, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		user.ID, user.Email, user.PasswordHash, user.FullName, user.Phone, user.Role, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusConflict, "email already registered")
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
		`SELECT id, email, password_hash, full_name, phone, role, is_locked, created_at, updated_at
		 FROM users WHERE email = $1`, req.Email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FullName, &user.Phone, &user.Role,
		&user.IsLocked, &user.CreatedAt, &user.UpdatedAt)

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
