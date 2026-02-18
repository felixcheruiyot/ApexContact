package middleware

import (
	"strings"

	"github.com/livestreamify/backend/internal/config"
	"github.com/livestreamify/backend/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gofiber/fiber/v2"
)

type Claims struct {
	UserID string          `json:"user_id"`
	Role   domain.UserRole `json:"role"`
	jwt.RegisteredClaims
}

// RequireAuth validates the JWT bearer token and attaches claims to context.
func RequireAuth(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			return fiber.NewError(fiber.StatusUnauthorized, "missing or invalid authorization header")
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "unexpected signing method")
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid or expired token")
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("role", claims.Role)
		return c.Next()
	}
}

// RequireMediaServerKey validates requests from the Nginx-RTMP media server.
// nginx-rtmp cannot send Authorization headers, so we use a shared secret
// passed as the ?key= query parameter.
func RequireMediaServerKey(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if cfg.MediaServerKey == "" || c.Query("key") != cfg.MediaServerKey {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid media server key")
		}
		return c.Next()
	}
}

// RequireRole checks that the authenticated user has one of the allowed roles.
func RequireRole(roles ...string) fiber.Handler {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(domain.UserRole)
		if !ok || !allowed[string(role)] {
			return fiber.NewError(fiber.StatusForbidden, "insufficient permissions")
		}
		return c.Next()
	}
}
