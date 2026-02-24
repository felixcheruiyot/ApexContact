package intasend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client is an HTTP client for the IntaSend payment API.
type Client struct {
	secretKey string
	baseURL   string
	http      *http.Client
}

// New creates an IntaSend client. secretKey is the server-side API secret;
// baseURL is either the sandbox (https://sandbox.intasend.com) or production URL.
func New(secretKey, baseURL string) *Client {
	return &Client{
		secretKey: secretKey,
		baseURL:   baseURL,
		http:      &http.Client{Timeout: 30 * time.Second},
	}
}

type stkPushRequest struct {
	PhoneNumber string  `json:"phone_number"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	APIRef      string  `json:"api_ref"`
	Method      string  `json:"method"`
}

// stkPushResponse mirrors the IntaSend STK Push creation response.
type stkPushResponse struct {
	Invoice struct {
		InvoiceID    string `json:"invoice_id"`
		State        string `json:"state"`
		Provider     string `json:"provider"`
		Value        json.Number `json:"value"`
		Account      string `json:"account"`
		APIRef       string `json:"api_ref"`
		FailedReason any    `json:"failed_reason"`
	} `json:"invoice"`
}

// errorResponse is returned by IntaSend on 4xx/5xx.
type errorResponse struct {
	Detail any `json:"detail"`
}

// StkPush initiates an M-Pesa STK Push request to the customer's phone.
// apiRef should be the internal payment UUID so the webhook can correlate it.
// Returns the InvoiceID assigned by IntaSend.
func (c *Client) StkPush(phoneNumber string, amount float64, currency, apiRef string) (string, error) {
	payload := stkPushRequest{
		PhoneNumber: phoneNumber,
		Amount:      amount,
		Currency:    currency,
		APIRef:      apiRef,
		Method:      "MPESA_STK_PUSH",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("intasend: marshal: %w", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		c.baseURL+"/api/v1/payment/mpesa-stk-push/",
		bytes.NewReader(body),
	)
	if err != nil {
		return "", fmt.Errorf("intasend: build request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.secretKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("intasend: http: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var e errorResponse
		json.NewDecoder(resp.Body).Decode(&e)
		return "", fmt.Errorf("intasend: api error %d: %v", resp.StatusCode, e.Detail)
	}

	var result stkPushResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("intasend: decode response: %w", err)
	}

	return result.Invoice.InvoiceID, nil
}
