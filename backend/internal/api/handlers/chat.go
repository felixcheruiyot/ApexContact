package handlers

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/livestreamify/backend/internal/api/middleware"
	"github.com/livestreamify/backend/internal/config"
	"github.com/livestreamify/backend/internal/domain"
	"github.com/redis/go-redis/v9"
)

// ChatHandler manages WebSocket connections for lobby chat rooms.
type ChatHandler struct {
	cfg *config.Config
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewChatHandler(cfg *config.Config, db *pgxpool.Pool, rdb *redis.Client) *ChatHandler {
	return &ChatHandler{cfg: cfg, db: db, rdb: rdb}
}

type incomingMessage struct {
	Type    string `json:"type"`    // "message" | "reaction"
	Content string `json:"content"`
}

type wsError struct {
	Error string `json:"error"`
}

// Handle is the WebSocket handler for a commentary lobby chat room.
// Connection: WS /ws/commentary/:id/chat?token=<jwt>
func (h *ChatHandler) Handle(c *websocket.Conn) {
	writeErr := func(msg string) {
		_ = c.WriteJSON(wsError{Error: msg})
	}

	eventIDStr := c.Params("id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		writeErr("invalid event id")
		return
	}

	// ── Authenticate via query param token ────────────────────────────────────
	tokenStr := c.Query("token")
	if tokenStr == "" {
		writeErr("missing token")
		return
	}

	claims := &middleware.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(h.cfg.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		writeErr("invalid or expired token")
		return
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		writeErr("invalid user id in token")
		return
	}

	// ── Verify the user has joined this lobby ─────────────────────────────────
	var nickname string
	err = h.db.QueryRow(context.Background(),
		`SELECT nickname FROM lobby_participants WHERE event_id = $1 AND user_id = $2`,
		eventID, userID,
	).Scan(&nickname)
	if err != nil {
		writeErr("you have not joined this lobby")
		return
	}

	// ── Subscribe to Redis pub/sub channel ────────────────────────────────────
	channel := "lobby_chat:" + eventIDStr
	sub := h.rdb.Subscribe(context.Background(), channel)
	defer sub.Close()

	redisCh := sub.Channel()
	done := make(chan struct{})

	// Forward incoming Redis messages to this WebSocket connection
	go func() {
		defer close(done)
		for msg := range redisCh {
			if err := c.WriteMessage(websocket.TextMessage, []byte(msg.Payload)); err != nil {
				return
			}
		}
	}()

	// ── Announce arrival ──────────────────────────────────────────────────────
	joined := domain.ChatEvent{
		Type:      "joined",
		Nickname:  nickname,
		UserID:    claims.UserID,
		CreatedAt: time.Now().UTC(),
	}
	if raw, err := json.Marshal(joined); err == nil {
		h.rdb.Publish(context.Background(), channel, string(raw))
	}

	// ── Read loop ─────────────────────────────────────────────────────────────
	for {
		_, raw, err := c.ReadMessage()
		if err != nil {
			break
		}

		var inc incomingMessage
		if err := json.Unmarshal(raw, &inc); err != nil {
			continue
		}

		content := strings.TrimSpace(inc.Content)
		if content == "" {
			continue
		}

		now := time.Now().UTC()
		chatEvent := domain.ChatEvent{
			Nickname:  nickname,
			Content:   content,
			UserID:    claims.UserID,
			CreatedAt: now,
		}

		switch inc.Type {
		case "message":
			chatEvent.Type = "message"
			// Persist to DB
			_, _ = h.db.Exec(context.Background(), `
				INSERT INTO lobby_messages (event_id, user_id, nickname, content, message_type)
				VALUES ($1, $2, $3, $4, 'text')`,
				eventID, userID, nickname, content,
			)
			// Broadcast via Redis
			if payload, err := json.Marshal(chatEvent); err == nil {
				h.rdb.Publish(context.Background(), channel, string(payload))
			}

		case "reaction":
			chatEvent.Type = "reaction"
			// Reactions are ephemeral — publish only, not persisted
			if payload, err := json.Marshal(chatEvent); err == nil {
				h.rdb.Publish(context.Background(), channel, string(payload))
			}
		}
	}

	// ── Announce departure ────────────────────────────────────────────────────
	left := domain.ChatEvent{
		Type:      "left",
		Nickname:  nickname,
		UserID:    claims.UserID,
		CreatedAt: time.Now().UTC(),
	}
	if payload, err := json.Marshal(left); err == nil {
		h.rdb.Publish(context.Background(), channel, string(payload))
	}

	<-done
}
