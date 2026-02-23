package mailer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/livestreamify/backend/internal/config"
)

// Mailer sends transactional emails.
type Mailer interface {
	Send(to, subject, htmlBody string) error
}

// PostmarkMailer delivers email via the Postmark HTTP API.
type PostmarkMailer struct {
	serverToken string
	from        string
	fromName    string
	http        *http.Client
}

// NewPostmarkMailer creates a PostmarkMailer wired from config.
func NewPostmarkMailer(cfg *config.Config) *PostmarkMailer {
	return &PostmarkMailer{
		serverToken: cfg.PostmarkServerToken,
		from:        cfg.EmailFrom,
		fromName:    cfg.EmailFromName,
		http:        &http.Client{Timeout: 15 * time.Second},
	}
}

type postmarkRequest struct {
	From          string `json:"From"`
	To            string `json:"To"`
	Subject       string `json:"Subject"`
	HtmlBody      string `json:"HtmlBody"`
	MessageStream string `json:"MessageStream"`
}

type postmarkError struct {
	ErrorCode int    `json:"ErrorCode"`
	Message   string `json:"Message"`
}

// Send delivers an HTML email to the given address via Postmark.
func (m *PostmarkMailer) Send(to, subject, htmlBody string) error {
	from := m.from
	if m.fromName != "" {
		from = fmt.Sprintf("%s <%s>", m.fromName, m.from)
	}

	payload := postmarkRequest{
		From:          from,
		To:            to,
		Subject:       subject,
		HtmlBody:      htmlBody,
		MessageStream: "outbound",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("postmark: marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.postmarkapp.com/email", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("postmark: build request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Postmark-Server-Token", m.serverToken)

	resp, err := m.http.Do(req)
	if err != nil {
		return fmt.Errorf("postmark: http: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var e postmarkError
		json.NewDecoder(resp.Body).Decode(&e)
		return fmt.Errorf("postmark: api error %d (code %d): %s", resp.StatusCode, e.ErrorCode, e.Message)
	}

	return nil
}
