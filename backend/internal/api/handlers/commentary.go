package handlers

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/livestreamify/backend/internal/config"
	"github.com/livestreamify/backend/internal/domain"
	lk "github.com/livestreamify/backend/internal/integrations/livekit"
	"github.com/redis/go-redis/v9"
)

// CommentaryHandler handles audio commentary lobby endpoints.
type CommentaryHandler struct {
	cfg *config.Config
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewCommentaryHandler(cfg *config.Config, db *pgxpool.Pool, rdb *redis.Client) *CommentaryHandler {
	return &CommentaryHandler{cfg: cfg, db: db, rdb: rdb}
}

// ─── List ──────────────────────────────────────────────────────────────────────

// List returns paginated commentary lobbies (scheduled, live, or completed).
func (h *CommentaryHandler) List(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}
	limit := 20
	offset := (page - 1) * limit

	status := c.Query("status", "")
	sport := c.Query("sport", "")

	args := []any{}
	where := []string{"event_type = 'commentary'", "status != 'draft'", "status != 'cancelled'", "status != 'declined'"}
	idx := 1

	if status != "" {
		where = append(where, fmt.Sprintf("status = $%d", idx))
		args = append(args, status)
		idx++
	}
	if sport != "" {
		where = append(where, fmt.Sprintf("sport_type = $%d", idx))
		args = append(args, sport)
		idx++
	}

	whereClause := strings.Join(where, " AND ")

	var total int
	countQ := fmt.Sprintf("SELECT COUNT(*) FROM events WHERE %s", whereClause)
	if err := h.db.QueryRow(context.Background(), countQ, args...).Scan(&total); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to count lobbies")
	}

	args = append(args, limit, offset)
	q := fmt.Sprintf(`
		SELECT id, promoter_id, title, description, sport_type, scheduled_at,
		       status, price, currency, thumbnail_url, teaser_hook, event_type,
		       review_note, created_at, updated_at
		FROM events
		WHERE %s
		ORDER BY
		  CASE WHEN status = 'live' THEN 0 WHEN status = 'scheduled' THEN 1 ELSE 2 END,
		  scheduled_at ASC
		LIMIT $%d OFFSET $%d`, whereClause, idx, idx+1)

	rows, err := h.db.Query(context.Background(), q, args...)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to query lobbies")
	}
	defer rows.Close()

	events := []domain.Event{}
	for rows.Next() {
		var e domain.Event
		if err := rows.Scan(
			&e.ID, &e.PromoterID, &e.Title, &e.Description, &e.SportType, &e.ScheduledAt,
			&e.Status, &e.Price, &e.Currency, &e.ThumbnailURL, &e.TeaserHook, &e.EventType,
			&e.ReviewNote, &e.CreatedAt, &e.UpdatedAt,
		); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to scan lobby")
		}
		events = append(events, e)
	}

	return c.JSON(domain.Response{
		Data: events,
		Meta: &domain.Meta{Page: page, PerPage: limit, Total: total},
	})
}

// ─── Get ───────────────────────────────────────────────────────────────────────

// Get returns a single commentary lobby with participant count.
func (h *CommentaryHandler) Get(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid event id")
	}

	var e domain.Event
	err = h.db.QueryRow(context.Background(), `
		SELECT id, promoter_id, title, description, sport_type, scheduled_at,
		       status, price, currency, thumbnail_url, teaser_hook, event_type,
		       review_note, created_at, updated_at
		FROM events
		WHERE id = $1 AND event_type = 'commentary'`, id).Scan(
		&e.ID, &e.PromoterID, &e.Title, &e.Description, &e.SportType, &e.ScheduledAt,
		&e.Status, &e.Price, &e.Currency, &e.ThumbnailURL, &e.TeaserHook, &e.EventType,
		&e.ReviewNote, &e.CreatedAt, &e.UpdatedAt,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "lobby not found")
	}

	var participantCount int
	_ = h.db.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM lobby_participants WHERE event_id = $1`, id,
	).Scan(&participantCount)

	return c.JSON(domain.Response{
		Data: fiber.Map{
			"event":             e,
			"participant_count": participantCount,
		},
	})
}

// ─── Create ────────────────────────────────────────────────────────────────────

type createCommentaryRequest struct {
	Title        string  `json:"title"`
	TeaserHook   string  `json:"teaser_hook"`
	Description  string  `json:"description"`
	SportType    string  `json:"sport_type"`
	ScheduledAt  string  `json:"scheduled_at"`
	Price        float64 `json:"price"`
	ThumbnailURL string  `json:"thumbnail_url"`
}

// Create lets any authenticated user create a commentary lobby immediately (no review).
func (h *CommentaryHandler) Create(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req createCommentaryRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if strings.TrimSpace(req.Title) == "" {
		return fiber.NewError(fiber.StatusBadRequest, "title is required")
	}
	if req.SportType != "boxing" && req.SportType != "racing" {
		return fiber.NewError(fiber.StatusBadRequest, "sport_type must be boxing or racing")
	}

	// JavaScript toISOString() includes milliseconds (RFC3339Nano), try both formats.
	scheduledAt, err := time.Parse(time.RFC3339, req.ScheduledAt)
	if err != nil {
		scheduledAt, err = time.Parse(time.RFC3339Nano, req.ScheduledAt)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "scheduled_at must be an ISO 8601 date-time string")
		}
	}

	promoterID, err := uuid.Parse(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid user id")
	}

	// Commentary lobbies skip review → go straight to scheduled.
	// stream_key must be unique even though commentary doesn't use RTMP.
	streamKey := "commentary-" + uuid.New().String()

	var e domain.Event
	err = h.db.QueryRow(context.Background(), `
		INSERT INTO events
		  (promoter_id, title, teaser_hook, description, sport_type, scheduled_at,
		   status, price, currency, thumbnail_url, stream_key, hls_path,
		   event_type, livekit_room, review_note)
		VALUES ($1,$2,$3,$4,$5,$6,'scheduled',$7,'KES',$8,$9,'','commentary','','')
		RETURNING id, promoter_id, title, description, sport_type, scheduled_at,
		          status, price, currency, thumbnail_url, teaser_hook, event_type,
		          review_note, created_at, updated_at`,
		promoterID, req.Title, req.TeaserHook, req.Description, req.SportType,
		scheduledAt, req.Price, req.ThumbnailURL, streamKey,
	).Scan(
		&e.ID, &e.PromoterID, &e.Title, &e.Description, &e.SportType, &e.ScheduledAt,
		&e.Status, &e.Price, &e.Currency, &e.ThumbnailURL, &e.TeaserHook, &e.EventType,
		&e.ReviewNote, &e.CreatedAt, &e.UpdatedAt,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to create lobby: "+err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(domain.Response{Data: e})
}

// ─── Start ─────────────────────────────────────────────────────────────────────

// Start transitions a lobby to live and sets up its LiveKit room name.
func (h *CommentaryHandler) Start(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid event id")
	}

	// Verify ownership
	var promoterID string
	if err := h.db.QueryRow(context.Background(),
		`SELECT promoter_id FROM events WHERE id = $1 AND event_type = 'commentary'`, id,
	).Scan(&promoterID); err != nil {
		return fiber.NewError(fiber.StatusNotFound, "lobby not found")
	}
	if promoterID != userID {
		return fiber.NewError(fiber.StatusForbidden, "only the creator can start this lobby")
	}

	roomName := fmt.Sprintf("commentary-%s", id.String())
	_, err = h.db.Exec(context.Background(),
		`UPDATE events SET status = 'live', livekit_room = $1, updated_at = NOW() WHERE id = $2`,
		roomName, id,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to start lobby")
	}

	return c.JSON(domain.Response{Data: fiber.Map{"status": "live", "room": roomName}})
}

// ─── End ───────────────────────────────────────────────────────────────────────

// End marks a lobby as completed.
func (h *CommentaryHandler) End(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid event id")
	}

	var promoterID string
	if err := h.db.QueryRow(context.Background(),
		`SELECT promoter_id FROM events WHERE id = $1 AND event_type = 'commentary'`, id,
	).Scan(&promoterID); err != nil {
		return fiber.NewError(fiber.StatusNotFound, "lobby not found")
	}
	if promoterID != userID {
		return fiber.NewError(fiber.StatusForbidden, "only the creator can end this lobby")
	}

	_, err = h.db.Exec(context.Background(),
		`UPDATE events SET status = 'completed', updated_at = NOW() WHERE id = $1`, id,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to end lobby")
	}

	return c.JSON(domain.Response{Data: fiber.Map{"status": "completed"}})
}

// ─── Join ──────────────────────────────────────────────────────────────────────

type joinRequest struct {
	Nickname string `json:"nickname"`
}

// Join subscribes a user to a commentary lobby (creates payment subscription if needed).
func (h *CommentaryHandler) Join(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid event id")
	}

	var req joinRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if strings.TrimSpace(req.Nickname) == "" {
		return fiber.NewError(fiber.StatusBadRequest, "nickname is required")
	}

	// Fetch event
	var price float64
	var eventStatus string
	var promoterIDStr string
	if err := h.db.QueryRow(context.Background(),
		`SELECT price, status, promoter_id FROM events WHERE id = $1 AND event_type = 'commentary'`, id,
	).Scan(&price, &eventStatus, &promoterIDStr); err != nil {
		return fiber.NewError(fiber.StatusNotFound, "lobby not found")
	}

	if eventStatus == "completed" || eventStatus == "cancelled" {
		return fiber.NewError(fiber.StatusBadRequest, "this lobby is no longer accepting participants")
	}

	uid, _ := uuid.Parse(userID)
	role := domain.LobbyRoleListener
	if promoterIDStr == userID {
		role = domain.LobbyRoleHost
	}

	// For paid lobbies, verify an approved payment/subscription exists
	if price > 0 && promoterIDStr != userID {
		var count int
		_ = h.db.QueryRow(context.Background(),
			`SELECT COUNT(*) FROM subscriptions WHERE user_id = $1 AND event_id = $2`, uid, id,
		).Scan(&count)
		if count == 0 {
			return fiber.NewError(fiber.StatusPaymentRequired, "payment required to join this lobby")
		}
	}

	// For free lobbies or the creator, ensure a subscription row exists (zero-amount)
	if price == 0 || promoterIDStr == userID {
		var subCount int
		_ = h.db.QueryRow(context.Background(),
			`SELECT COUNT(*) FROM subscriptions WHERE user_id = $1 AND event_id = $2`, uid, id,
		).Scan(&subCount)
		if subCount == 0 {
			nilUUID := uuid.Nil
			streamToken := uuid.New().String()
			_, _ = h.db.Exec(context.Background(), `
				INSERT INTO subscriptions
				  (user_id, event_id, payment_id, stream_token, device_fingerprint, ip_lock,
				   active_session_id, expires_at)
				VALUES ($1, $2, $3, $4, '', '', '', NOW() + INTERVAL '30 days')
				ON CONFLICT DO NOTHING`,
				uid, id, nilUUID, streamToken,
			)
		}
	}

	// Upsert lobby_participants
	_, err = h.db.Exec(context.Background(), `
		INSERT INTO lobby_participants (event_id, user_id, nickname, role)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (event_id, user_id) DO UPDATE SET nickname = EXCLUDED.nickname`,
		id, uid, req.Nickname, role,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to join lobby")
	}

	return c.JSON(domain.Response{Data: fiber.Map{
		"joined":   true,
		"nickname": req.Nickname,
		"role":     role,
	}})
}

// ─── GetToken ──────────────────────────────────────────────────────────────────

// GetToken issues a LiveKit JWT for a participant who has joined the lobby.
func (h *CommentaryHandler) GetToken(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid event id")
	}

	// Verify subscription/membership
	uid, _ := uuid.Parse(userID)
	var subCount int
	_ = h.db.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM subscriptions WHERE user_id = $1 AND event_id = $2`, uid, id,
	).Scan(&subCount)
	if subCount == 0 {
		return fiber.NewError(fiber.StatusForbidden, "you have not joined this lobby")
	}

	// Fetch participant role + livekit_room
	var role domain.LobbyRole
	var roomName string
	err = h.db.QueryRow(context.Background(), `
		SELECT lp.role, e.livekit_room
		FROM lobby_participants lp
		JOIN events e ON e.id = lp.event_id
		WHERE lp.user_id = $1 AND lp.event_id = $2`, uid, id,
	).Scan(&role, &roomName)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "participant record not found")
	}
	if roomName == "" {
		return fiber.NewError(fiber.StatusBadRequest, "lobby is not live yet")
	}

	canPublish := role == domain.LobbyRoleHost || role == domain.LobbyRoleSpeaker
	isAdmin := role == domain.LobbyRoleHost

	token, err := lk.GenerateToken(
		h.cfg.LiveKitAPIKey,
		h.cfg.LiveKitAPISecret,
		roomName,
		userID,
		canPublish,
		isAdmin,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to generate token")
	}

	return c.JSON(domain.Response{Data: fiber.Map{
		"token":    token,
		"room":     roomName,
		"livekit_url": h.cfg.LiveKitURL,
	}})
}

// ─── UpdateParticipant ─────────────────────────────────────────────────────────

type updateParticipantRequest struct {
	Role string `json:"role"` // "speaker" | "listener"
}

// UpdateParticipant grants or revokes speaker role for a participant.
func (h *CommentaryHandler) UpdateParticipant(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid event id")
	}
	targetID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user id")
	}

	// Only creator can change roles
	var promoterIDStr string
	if err := h.db.QueryRow(context.Background(),
		`SELECT promoter_id FROM events WHERE id = $1 AND event_type = 'commentary'`, id,
	).Scan(&promoterIDStr); err != nil {
		return fiber.NewError(fiber.StatusNotFound, "lobby not found")
	}
	if promoterIDStr != userID {
		return fiber.NewError(fiber.StatusForbidden, "only the host can change participant roles")
	}

	var req updateParticipantRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if req.Role != "speaker" && req.Role != "listener" {
		return fiber.NewError(fiber.StatusBadRequest, "role must be speaker or listener")
	}

	_, err = h.db.Exec(context.Background(),
		`UPDATE lobby_participants SET role = $1 WHERE event_id = $2 AND user_id = $3`,
		req.Role, id, targetID,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to update participant")
	}

	// Notify via Redis pub/sub so the chat hub can broadcast a role-change event
	eventType := "speaker_granted"
	if req.Role == "listener" {
		eventType = "speaker_revoked"
	}
	h.rdb.Publish(context.Background(),
		fmt.Sprintf("lobby_chat:%s", id.String()),
		fmt.Sprintf(`{"type":"%s","user_id":"%s","nickname":"","created_at":"%s"}`,
			eventType, targetID.String(), time.Now().UTC().Format(time.RFC3339)),
	)

	return c.JSON(domain.Response{Data: fiber.Map{"role": req.Role}})
}

// ─── Messages ──────────────────────────────────────────────────────────────────

// Messages returns paginated chat history for a lobby (for replay).
func (h *CommentaryHandler) Messages(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid event id")
	}

	page := c.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}
	limit := 50
	offset := (page - 1) * limit

	var total int
	_ = h.db.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM lobby_messages WHERE event_id = $1`, id,
	).Scan(&total)

	rows, err := h.db.Query(context.Background(), `
		SELECT id, event_id, user_id, nickname, content, message_type, created_at
		FROM lobby_messages
		WHERE event_id = $1
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3`, id, limit, offset)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch messages")
	}
	defer rows.Close()

	messages := []domain.LobbyMessage{}
	for rows.Next() {
		var m domain.LobbyMessage
		if err := rows.Scan(&m.ID, &m.EventID, &m.UserID, &m.Nickname, &m.Content, &m.MessageType, &m.CreatedAt); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to scan message")
		}
		messages = append(messages, m)
	}

	return c.JSON(domain.Response{
		Data: messages,
		Meta: &domain.Meta{Page: page, PerPage: limit, Total: total},
	})
}

// ─── SuggestNicknames ──────────────────────────────────────────────────────────

var adjectives = []string{
	"Silent", "Wild", "Ground", "Iron", "Blazing", "Shadow", "Thunder", "Neon",
	"Rapid", "Steel", "Midnight", "Golden", "Savage", "Cosmic", "Fearless", "Electric",
}

var nouns = []string{
	"Lion", "Hawk", "Fist", "Corner", "Fan", "Champion", "Falcon", "Wolf",
	"Tiger", "Eagle", "Panther", "Knight", "Bull", "Phoenix", "Viper", "Cobra",
}

// SuggestNicknames returns 4 auto-generated nickname suggestions.
func (h *CommentaryHandler) SuggestNicknames(c *fiber.Ctx) error {
	suggestions := make([]string, 4)
	seen := map[string]bool{}
	for i := 0; i < 4; {
		adj := adjectives[rand.Intn(len(adjectives))]
		noun := nouns[rand.Intn(len(nouns))]
		num := rand.Intn(90) + 10
		name := fmt.Sprintf("%s%s%d", adj, noun, num)
		if !seen[name] {
			suggestions[i] = name
			seen[name] = true
			i++
		}
	}
	return c.JSON(domain.Response{Data: suggestions})
}
