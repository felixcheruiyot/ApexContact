package domain

import (
	"time"

	"github.com/google/uuid"
)

// UserRole defines the access level of a user.
type UserRole string

const (
	RoleMember UserRole = "member"
	RoleAdmin  UserRole = "admin"
)

// EventStatus represents the lifecycle state of an event.
type EventStatus string

const (
	StatusDraft         EventStatus = "draft"
	StatusPendingReview EventStatus = "pending_review"
	StatusScheduled     EventStatus = "scheduled"
	StatusLive          EventStatus = "live"
	StatusCompleted     EventStatus = "completed"
	StatusCancelled     EventStatus = "cancelled"
	StatusDeclined      EventStatus = "declined"
)

// SportType categorises what kind of event is being streamed.
type SportType string

const (
	SportSales     SportType = "sales"
	SportMentoring SportType = "mentoring"
	SportEducation SportType = "education"
	SportBusiness  SportType = "business"
	SportLegal     SportType = "legal"
	SportFitness   SportType = "fitness"
	SportVisa      SportType = "visa"
	SportMusic     SportType = "music"
	SportGaming    SportType = "gaming"
	SportCooking   SportType = "cooking"
	SportCommunity SportType = "community"
	SportOther     SportType = "other"
)

var validSportTypes = map[string]bool{
	"sales": true, "mentoring": true, "education": true, "business": true,
	"legal": true, "fitness": true, "visa": true, "music": true,
	"gaming": true, "cooking": true, "community": true, "other": true,
}

func IsValidSportType(s string) bool {
	return validSportTypes[s]
}

// EventType defines how an event is broadcast.
//
//   - video       — commercial OBS/RTMP stream → HLS delivery (one-to-many)
//   - audio_video — interactive LiveKit room with camera + mic (two-way)
//   - audio       — interactive LiveKit room with mic only (two-way)
//
// The legacy value "commentary" is treated as "audio_video" in all queries
// during the migration window.
type EventType string

const (
	EventTypeVideo      EventType = "video"
	EventTypeAudioVideo EventType = "audio_video"
	EventTypeAudio      EventType = "audio"

	// EventTypeCommentary is the legacy name for audio_video; kept for
	// backward-compat with any un-migrated rows or callers.
	EventTypeCommentary EventType = "commentary"
)

// LobbyRole defines a participant's speaking permissions in a commentary lobby.
type LobbyRole string

const (
	LobbyRoleHost     LobbyRole = "host"
	LobbyRoleSpeaker  LobbyRole = "speaker"
	LobbyRoleListener LobbyRole = "listener"
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
	ID            uuid.UUID `json:"id" db:"id"`
	Email         string    `json:"email" db:"email"`
	PasswordHash  string    `json:"-" db:"password_hash"`
	FullName      string    `json:"full_name" db:"full_name"`
	Phone         string    `json:"phone" db:"phone"`
	Role          UserRole  `json:"role" db:"role"`
	IsLocked      bool      `json:"is_locked" db:"is_locked"`
	EmailVerified bool      `json:"email_verified" db:"email_verified"`
	Age           *int      `json:"age,omitempty" db:"age"`
	Gender        string    `json:"gender,omitempty" db:"gender"`
	Country       string    `json:"country,omitempty" db:"country"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
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
	StreamActive bool        `json:"stream_active" db:"-"` // computed: hls_path != ""
	ReviewNote   string      `json:"review_note" db:"review_note"`
	EventType    EventType   `json:"event_type" db:"event_type"`
	LiveKitRoom  string      `json:"-" db:"livekit_room"`
	TeaserHook   string      `json:"teaser_hook" db:"teaser_hook"`
	IsPublic     bool        `json:"is_public" db:"is_public"`
	CreatedAt    time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at" db:"updated_at"`
}

// ─── Commentary Lobby ──────────────────────────────────────────────────────────

type LobbyParticipant struct {
	ID       uuid.UUID `json:"id" db:"id"`
	EventID  uuid.UUID `json:"event_id" db:"event_id"`
	UserID   uuid.UUID `json:"user_id" db:"user_id"`
	Nickname string    `json:"nickname" db:"nickname"`
	Role     LobbyRole `json:"role" db:"role"`
	JoinedAt time.Time `json:"joined_at" db:"joined_at"`
}

type LobbyMessage struct {
	ID          uuid.UUID `json:"id" db:"id"`
	EventID     uuid.UUID `json:"event_id" db:"event_id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	Nickname    string    `json:"nickname" db:"nickname"`
	Content     string    `json:"content" db:"content"`
	MessageType string    `json:"message_type" db:"message_type"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// ChatEvent is the envelope for real-time WebSocket messages.
type ChatEvent struct {
	Type      string    `json:"type"`    // "message" | "reaction" | "joined" | "left" | "speaker_granted" | "speaker_revoked"
	Nickname  string    `json:"nickname"`
	Content   string    `json:"content,omitempty"`
	UserID    string    `json:"user_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
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

// ─── Withdrawal / Payout ──────────────────────────────────────────────────────

type WithdrawalStatus string

const (
	WithdrawalPendingOTP WithdrawalStatus = "pending_otp"
	WithdrawalProcessing WithdrawalStatus = "processing"
	WithdrawalCompleted  WithdrawalStatus = "completed"
	WithdrawalFailed     WithdrawalStatus = "failed"
)

type PayoutAccount struct {
	ID            uuid.UUID `json:"id" db:"id"`
	UserID        uuid.UUID `json:"user_id" db:"user_id"`
	AccountType   string    `json:"account_type" db:"account_type"`
	AccountNumber string    `json:"account_number" db:"account_number"`
	AccountName   string    `json:"account_name" db:"account_name"`
	BankName      string    `json:"bank_name" db:"bank_name"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

type Withdrawal struct {
	ID              uuid.UUID        `json:"id" db:"id"`
	UserID          uuid.UUID        `json:"user_id" db:"user_id"`
	PayoutAccountID uuid.UUID        `json:"payout_account_id" db:"payout_account_id"`
	Amount          float64          `json:"amount" db:"amount"`
	Currency        string           `json:"currency" db:"currency"`
	Status          WithdrawalStatus `json:"status" db:"status"`
	IntaSendRef     string           `json:"intasend_ref,omitempty" db:"intasend_ref"`
	FailureReason   string           `json:"failure_reason,omitempty" db:"failure_reason"`
	CreatedAt       time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at" db:"updated_at"`
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
