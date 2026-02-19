package handlers

import (
	"context"
	"time"

	"github.com/livestreamify/backend/internal/config"
	"github.com/livestreamify/backend/internal/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PromoterHandler struct {
	cfg *config.Config
	db  *pgxpool.Pool
}

func NewPromoterHandler(cfg *config.Config, db *pgxpool.Pool) *PromoterHandler {
	return &PromoterHandler{cfg: cfg, db: db}
}

// StreamKey returns the RTMP stream key for a specific event owned by the promoter.
// This is sensitive — only the promoter of that event can retrieve it.
func (h *PromoterHandler) StreamKey(c *fiber.Ctx) error {
	eventID := c.Params("eventId")
	promoterID := c.Locals("user_id").(string)

	var streamKey string
	err := h.db.QueryRow(context.Background(),
		`SELECT stream_key FROM events WHERE id=$1 AND promoter_id=$2`,
		eventID, promoterID,
	).Scan(&streamKey)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "event not found or access denied")
	}

	return c.JSON(domain.Response{
		Data: fiber.Map{
			"stream_key": streamKey,
			"rtmp_url":   "rtmp://<your-host>:1935/live",
			"push_to":    "rtmp://<your-host>:1935/live/" + streamKey,
		},
	})
}

func (h *PromoterHandler) MyEvents(c *fiber.Ctx) error {
	promoterID := c.Locals("user_id").(string)

	rows, err := h.db.Query(context.Background(),
		`SELECT id, title, description, sport_type, scheduled_at, status, price, currency, thumbnail_url, review_note, created_at
		 FROM events WHERE promoter_id=$1 ORDER BY scheduled_at DESC`, promoterID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch events")
	}
	defer rows.Close()

	var events []domain.Event
	for rows.Next() {
		var e domain.Event
		rows.Scan(&e.ID, &e.Title, &e.Description, &e.SportType, &e.ScheduledAt,
			&e.Status, &e.Price, &e.Currency, &e.ThumbnailURL, &e.ReviewNote, &e.CreatedAt)
		events = append(events, e)
	}

	return c.JSON(domain.Response{Data: events})
}

func (h *PromoterHandler) Submit(c *fiber.Ctx) error {
	eventID := c.Params("eventId")
	promoterID := c.Locals("user_id").(string)

	tag, err := h.db.Exec(context.Background(),
		`UPDATE events SET status='pending_review', updated_at=NOW()
		 WHERE id=$1 AND promoter_id=$2 AND status='draft'`,
		eventID, promoterID,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to submit event")
	}
	if tag.RowsAffected() == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "event not found, not in draft status, or access denied")
	}

	return c.JSON(domain.Response{Data: "event submitted for review"})
}

func (h *PromoterHandler) Analytics(c *fiber.Ctx) error {
	eventID := c.Params("eventId")
	promoterID := c.Locals("user_id").(string)

	// Verify ownership
	var count int
	h.db.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM events WHERE id=$1 AND promoter_id=$2`, eventID, promoterID,
	).Scan(&count)
	if count == 0 {
		return fiber.NewError(fiber.StatusForbidden, "event not found or access denied")
	}

	var stats struct {
		EventID      string  `json:"event_id"`
		TotalTickets int     `json:"total_tickets"`
		TotalRevenue float64 `json:"total_revenue"`
		PeakViewers  int     `json:"peak_viewers"`
	}
	stats.EventID = eventID

	h.db.QueryRow(context.Background(),
		`SELECT COUNT(*), COALESCE(SUM(p.amount),0)
		 FROM subscriptions s JOIN payments p ON s.payment_id=p.id
		 WHERE s.event_id=$1 AND p.status='success'`, eventID,
	).Scan(&stats.TotalTickets, &stats.TotalRevenue)

	h.db.QueryRow(context.Background(),
		`SELECT COALESCE(MAX(peak_viewers),0) FROM stream_analytics WHERE event_id=$1`, eventID,
	).Scan(&stats.PeakViewers)

	return c.JSON(domain.Response{Data: stats})
}

func (h *PromoterHandler) Revenue(c *fiber.Ctx) error {
	promoterID := c.Locals("user_id").(string)

	rows, err := h.db.Query(context.Background(),
		`SELECT
		   e.id, e.title, e.sport_type, e.status, e.scheduled_at,
		   e.price, e.currency,
		   COUNT(s.id)                         AS tickets,
		   COALESCE(SUM(p.amount), 0)          AS revenue,
		   COALESCE(MAX(sa.peak_viewers), 0)   AS peak_viewers
		 FROM events e
		 LEFT JOIN subscriptions s  ON e.id = s.event_id
		 LEFT JOIN payments p       ON s.payment_id = p.id AND p.status = 'success'
		 LEFT JOIN stream_analytics sa ON e.id = sa.event_id
		 WHERE e.promoter_id = $1
		 GROUP BY e.id, e.title, e.sport_type, e.status, e.scheduled_at, e.price, e.currency
		 ORDER BY e.scheduled_at DESC`, promoterID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch revenue data")
	}
	defer rows.Close()

	type eventRevenue struct {
		EventID      string    `json:"event_id"`
		Title        string    `json:"title"`
		SportType    string    `json:"sport_type"`
		Status       string    `json:"status"`
		ScheduledAt  time.Time `json:"scheduled_at"`
		Price        float64   `json:"price"`
		Currency     string    `json:"currency"`
		Tickets      int       `json:"tickets"`
		Revenue      float64   `json:"revenue"`
		PromoterCut  float64   `json:"promoter_cut"`
		PeakViewers  int       `json:"peak_viewers"`
	}
	var results []eventRevenue
	for rows.Next() {
		var er eventRevenue
		rows.Scan(&er.EventID, &er.Title, &er.SportType, &er.Status, &er.ScheduledAt,
			&er.Price, &er.Currency, &er.Tickets, &er.Revenue, &er.PeakViewers)
		er.PromoterCut = er.Revenue * 0.70
		results = append(results, er)
	}
	if results == nil {
		results = []eventRevenue{}
	}

	return c.JSON(domain.Response{Data: results})
}
