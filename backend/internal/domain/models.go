package domain

import (
	"time"

	"github.com/google/uuid"
)

// UserRole defines the access level of a user.
type UserRole string

const (
	RoleViewer    UserRole = "viewer"
	RolePromoter  UserRole = "promoter"
	RoleBroadcaster UserRole = "broadcaster"
	RoleAdmin     UserRole = "admin"
)

// EventStatus represents the lifecycle state of an event.
type EventStatus string

const (
	StatusScheduled EventStatus = "scheduled"
	StatusLive      EventStatus = "live"
	StatusCompleted EventStatus = "completed"
	StatusCancelled EventStatus = "cancelled"
)

// SportType categorises what kind of event is being streamed.
type SportType string

const (
	SportBoxing  SportType = "boxing"
	SportRacing  SportType = "racing"
)

// PaymentStatus tracks the M-Pesa payment lifecycle.
type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "pending"
	PaymentSuccess   PaymentStatus = "success"
	PaymentFailed    PaymentStatus = "failed"
	PaymentCancelled PaymentStatus = "cancelled"
)

// ─── User ─────────────────────────────────────────────────────────────────────

type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	FullName     string    `json:"full_name" db:"full_name"`
	Phone        string    `json:"phone" db:"phone"`
	Role         UserRole  `json:"role" db:"role"`
	IsLocked     bool      `json:"is_locked" db:"is_locked"`
	Age          *int      `json:"age,omitempty" db:"age"`
	Gender       string    `json:"gender,omitempty" db:"gender"`
	Country      string    `json:"country,omitempty" db:"country"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// ─── Event ────────────────────────────────────────────────────────────────────

type Event struct {
	ID           uuid.UUID   `json:"id" db:"id"`
	PromoterID   uuid.UUID   `json:"promoter_id" db:"promoter_id"`
	Title        string      `json:"title" db:"title"`
	Description  string      `json:"description" db:"description"`
	SportType    SportType   `json:"sport_type" db:"sport_type"`
	ScheduledAt  time.Time   `json:"scheduled_at" db:"scheduled_at"`
	Status       EventStatus `json:"status" db:"status"`
	Price        float64     `json:"price" db:"price"`
	Currency     string      `json:"currency" db:"currency"`
	ThumbnailURL string      `json:"thumbnail_url" db:"thumbnail_url"`
	StreamKey    string      `json:"-" db:"stream_key"`  // never expose to viewers
	HLSPath      string      `json:"-" db:"hls_path"`
	CreatedAt    time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at" db:"updated_at"`
}

// ─── Payment ──────────────────────────────────────────────────────────────────

type Payment struct {
	ID           uuid.UUID     `json:"id" db:"id"`
	UserID       uuid.UUID     `json:"user_id" db:"user_id"`
	EventID      uuid.UUID     `json:"event_id" db:"event_id"`
	Amount       float64       `json:"amount" db:"amount"`
	Currency     string        `json:"currency" db:"currency"`
	Status       PaymentStatus `json:"status" db:"status"`
	IntaSendRef  string        `json:"intasend_ref" db:"intasend_ref"`
	PhoneNumber  string        `json:"phone_number" db:"phone_number"`
	CreatedAt    time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at" db:"updated_at"`
}

// ─── Subscription / StreamToken ───────────────────────────────────────────────

type Subscription struct {
	ID                uuid.UUID `json:"id" db:"id"`
	UserID            uuid.UUID `json:"user_id" db:"user_id"`
	EventID           uuid.UUID `json:"event_id" db:"event_id"`
	PaymentID         uuid.UUID `json:"payment_id" db:"payment_id"`
	StreamToken       string    `json:"stream_token" db:"stream_token"`
	DeviceFingerprint string    `json:"-" db:"device_fingerprint"`
	IPLock            string    `json:"-" db:"ip_lock"`
	ActiveSessionID   string    `json:"-" db:"active_session_id"`
	ExpiresAt         time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
}

// ─── Stream Analytics ─────────────────────────────────────────────────────────

type StreamAnalytics struct {
	ID            uuid.UUID `json:"id" db:"id"`
	EventID       uuid.UUID `json:"event_id" db:"event_id"`
	ViewerCount   int       `json:"viewer_count" db:"viewer_count"`
	PeakViewers   int       `json:"peak_viewers" db:"peak_viewers"`
	TotalRevenue  float64   `json:"total_revenue" db:"total_revenue"`
	RecordedAt    time.Time `json:"recorded_at" db:"recorded_at"`
}

// ─── Fraud Flag ───────────────────────────────────────────────────────────────

type FraudFlag struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	UserID         uuid.UUID  `json:"user_id" db:"user_id"`
	SubscriptionID uuid.UUID  `json:"subscription_id" db:"subscription_id"`
	Reason         string     `json:"reason" db:"reason"`
	DetectedAt     time.Time  `json:"detected_at" db:"detected_at"`
	Resolved       bool       `json:"resolved" db:"resolved"`
	ResolvedAt     *time.Time `json:"resolved_at,omitempty" db:"resolved_at"`
}

// ─── API Response envelopes ───────────────────────────────────────────────────

type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
	Meta  *Meta       `json:"meta,omitempty"`
}

type Meta struct {
	Page    int `json:"page,omitempty"`
	PerPage int `json:"per_page,omitempty"`
	Total   int `json:"total,omitempty"`
}
