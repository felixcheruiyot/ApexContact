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

// sendMoneyInitiateRequest is the body for the send-money initiate call.
type sendMoneyInitiateRequest struct {
	Currency     string            `json:"currency"`
	Transactions []sendMoneyTxn    `json:"transactions"`
}

type sendMoneyTxn struct {
	Name      string  `json:"name"`
	Account   string  `json:"account"`
	Amount    float64 `json:"amount"`
	Narrative string  `json:"narrative"`
}

type sendMoneyInitiateResponse struct {
	ID           string `json:"id"`
	Nonce        string `json:"nonce"`
	Transactions []struct {
		TrackingID string `json:"tracking_id"`
		Status     string `json:"status"`
	} `json:"transactions"`
}

// SendMoney executes a two-step IntaSend bulk transfer (initiate → approve).
// Returns the transaction tracking_id on success.
func (c *Client) SendMoney(name, account string, amount float64, currency, narrative string) (string, error) {
	// Step 1: Initiate
	initiateBody, err := json.Marshal(sendMoneyInitiateRequest{
		Currency: currency,
		Transactions: []sendMoneyTxn{
			{Name: name, Account: account, Amount: amount, Narrative: narrative},
		},
	})
	if err != nil {
		return "", fmt.Errorf("intasend sendmoney: marshal initiate: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/api/v1/send-money/initiate/", bytes.NewReader(initiateBody))
	if err != nil {
		return "", fmt.Errorf("intasend sendmoney: build initiate request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.secretKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("intasend sendmoney: initiate http: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var e errorResponse
		json.NewDecoder(resp.Body).Decode(&e)
		return "", fmt.Errorf("intasend sendmoney: initiate api error %d: %v", resp.StatusCode, e.Detail)
	}

	var initiated sendMoneyInitiateResponse
	if err := json.NewDecoder(resp.Body).Decode(&initiated); err != nil {
		return "", fmt.Errorf("intasend sendmoney: decode initiate response: %w", err)
	}
	if len(initiated.Transactions) == 0 {
		return "", fmt.Errorf("intasend sendmoney: no transactions in initiate response")
	}
	trackingID := initiated.Transactions[0].TrackingID

	// Step 2: Approve
	approveBody, err := json.Marshal(map[string]string{"nonce": initiated.Nonce})
	if err != nil {
		return "", fmt.Errorf("intasend sendmoney: marshal approve: %w", err)
	}

	approveURL := fmt.Sprintf("%s/api/v1/send-money/%s/approve/", c.baseURL, initiated.ID)
	req2, err := http.NewRequest(http.MethodPost, approveURL, bytes.NewReader(approveBody))
	if err != nil {
		return "", fmt.Errorf("intasend sendmoney: build approve request: %w", err)
	}
	req2.Header.Set("Authorization", "Bearer "+c.secretKey)
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Accept", "application/json")

	resp2, err := c.http.Do(req2)
	if err != nil {
		return "", fmt.Errorf("intasend sendmoney: approve http: %w", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode >= 400 {
		var e errorResponse
		json.NewDecoder(resp2.Body).Decode(&e)
		return "", fmt.Errorf("intasend sendmoney: approve api error %d: %v", resp2.StatusCode, e.Detail)
	}

	return trackingID, nil
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
