package postal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	BaseURL string
	APIKey  string
	client  *http.Client
}

func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL: baseURL,
		APIKey:  apiKey,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendEmailRequest represents the request to send an email
type SendEmailRequest struct {
	To          []string          `json:"to"`
	From        string            `json:"from"`
	Subject     string            `json:"subject"`
	HTMLBody    string            `json:"html_body,omitempty"`
	PlainBody   string            `json:"plain_body,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	Attachments []Attachment      `json:"attachments,omitempty"`
}

type Attachment struct {
	Name        string `json:"name"`
	ContentType string `json:"content_type"`
	Data        string `json:"data"` // Base64 encoded
}

// SendEmailResponse represents the response from Postal
type SendEmailResponse struct {
	Status   string                 `json:"status"`
	Time     float64                `json:"time"`
	Flags    map[string]interface{} `json:"flags"`
	Data     SendEmailResponseData  `json:"data"`
	Messages []string               `json:"messages,omitempty"`
}

type SendEmailResponseData struct {
	MessageID string                 `json:"message_id"`
	Messages  map[string]MessageInfo `json:"messages"`
}

type MessageInfo struct {
	ID    int    `json:"id"`
	Token string `json:"token"`
}

// SendEmail sends an email via Postal
func (c *Client) SendEmail(req SendEmailRequest) (*SendEmailResponse, error) {
	payload := map[string]interface{}{
		"to":      req.To,
		"from":    req.From,
		"subject": req.Subject,
	}

	if req.HTMLBody != "" {
		payload["html_body"] = req.HTMLBody
	}
	if req.PlainBody != "" {
		payload["plain_body"] = req.PlainBody
	}
	if req.Headers != nil {
		payload["headers"] = req.Headers
	}
	if req.Attachments != nil {
		payload["attachments"] = req.Attachments
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.doRequest("POST", "/api/v1/send/message", body)
	if err != nil {
		return nil, err
	}

	var result SendEmailResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if result.Status != "success" {
		return nil, fmt.Errorf("postal API error: %v", result.Messages)
	}

	return &result, nil
}

// GetMessageDetails retrieves details about a sent message
type MessageDetails struct {
	ID               int        `json:"id"`
	Token            string     `json:"token"`
	Direction        string     `json:"direction"`
	MessageID        string     `json:"message_id"`
	To               string     `json:"to"`
	From             string     `json:"from"`
	Subject          string     `json:"subject"`
	Timestamp        time.Time  `json:"timestamp"`
	Status           string     `json:"status"`
	HeldUntil        *time.Time `json:"held_until,omitempty"`
	InspectionStatus string     `json:"inspection_status,omitempty"`
}

func (c *Client) GetMessage(messageID string) (*MessageDetails, error) {
	payload := map[string]interface{}{
		"id": messageID,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.doRequest("POST", "/api/v1/messages/message", body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Status string         `json:"status"`
		Data   MessageDetails `json:"data"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if result.Status != "success" {
		return nil, fmt.Errorf("postal API error")
	}

	return &result.Data, nil
}

// doRequest performs an HTTP request to Postal API
func (c *Client) doRequest(method, path string, body []byte) ([]byte, error) {
	url := c.BaseURL + path

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Server-API-Key", c.APIKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("postal API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}
