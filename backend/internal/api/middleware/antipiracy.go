package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

// AntiPiracy enforces single-session streaming per token.
// If a token is already active on another session, the old session is invalidated.
func AntiPiracy(rdb *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Query("token")
		if token == "" {
			return fiber.NewError(fiber.StatusBadRequest, "stream token required")
		}

		userID, ok := c.Locals("user_id").(string)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "unauthenticated")
		}

		clientIP := c.IP()
		sessionKey := fmt.Sprintf("session:%s", token)
		ctx := context.Background()

		// Check for an existing active session
		existingSession, err := rdb.Get(ctx, sessionKey).Result()
		if err != nil && err != redis.Nil {
			return fiber.NewError(fiber.StatusInternalServerError, "session check failed")
		}

		newSession := fmt.Sprintf("%s:%s", userID, clientIP)

		if existingSession != "" && existingSession != newSession {
			// Different device/IP is trying to use the same token — flag it
			fraudKey := fmt.Sprintf("fraud:%s", token)
			rdb.Set(ctx, fraudKey, "duplicate_session", 24*time.Hour)
			return fiber.NewError(fiber.StatusForbidden, "stream token is already active on another device")
		}

		// Register or refresh this session (TTL extends on activity)
		rdb.Set(ctx, sessionKey, newSession, 30*time.Minute)

		c.Locals("stream_token", token)
		c.Locals("client_ip", clientIP)
		return c.Next()
	}
}
