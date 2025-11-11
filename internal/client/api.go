package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

// APIClient represents the Fleeks API client
type APIClient struct {
	client   *resty.Client
	baseURL  string
	apiKey   string
	timeout  time.Duration
	wsDialer *websocket.Dialer
}

// NewAPIClient creates a new Fleeks API client
func NewAPIClient() *APIClient {
	baseURL := viper.GetString("api.base_url")
	if baseURL == "" {
		baseURL = "https://api.fleeks.dev"
	}

	timeout := viper.GetDuration("api.timeout")
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	client := resty.New().
		SetBaseURL(baseURL).
		SetTimeout(timeout).
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "fleeks-cli/1.0.0")

	// Configure TLS
	client.SetTLSClientConfig(&tls.Config{
		InsecureSkipVerify: false,
	})

	// WebSocket dialer
	wsDialer := &websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
	}

	return &APIClient{
		client:   client,
		baseURL:  baseURL,
		timeout:  timeout,
		wsDialer: wsDialer,
	}
}

// SetAPIKey sets the API key for authentication
func (c *APIClient) SetAPIKey(apiKey string) {
	c.apiKey = apiKey
	c.client.SetHeader("Authorization", fmt.Sprintf("Bearer %s", apiKey))
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	// Message maps to the standard `error` field returned by the API
	Message string `json:"error"`
	// Detail contains any additional message or context
	Detail string `json:"message,omitempty"`
	Code   int    `json:"code,omitempty"`
}

// Error implements the error interface for ErrorResponse
func (e *ErrorResponse) Error() string {
	if e == nil {
		return ""
	}
	if e.Code != 0 {
		return fmt.Sprintf("API Error %d: %s - %s", e.Code, e.Message, e.Detail)
	}
	if e.Detail != "" {
		return fmt.Sprintf("API Error: %s - %s", e.Message, e.Detail)
	}
	return fmt.Sprintf("API Error: %s", e.Message)
}

// GET makes a GET request to the API
func (c *APIClient) GET(endpoint string, result interface{}) error {
	resp, err := c.client.R().
		SetResult(result).
		SetError(&ErrorResponse{}).
		Get(endpoint)

	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	if !resp.IsSuccess() {
		if errResp, ok := resp.Error().(*ErrorResponse); ok {
			errResp.Code = resp.StatusCode()
			return errResp
		}
		return fmt.Errorf("request failed with status %d", resp.StatusCode())
	}

	return nil
}

// POST makes a POST request to the API
func (c *APIClient) POST(endpoint string, body interface{}, result interface{}) error {
	resp, err := c.client.R().
		SetBody(body).
		SetResult(result).
		SetError(&ErrorResponse{}).
		Post(endpoint)

	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	if !resp.IsSuccess() {
		if errResp, ok := resp.Error().(*ErrorResponse); ok {
			errResp.Code = resp.StatusCode()
			return errResp
		}
		return fmt.Errorf("request failed with status %d", resp.StatusCode())
	}

	return nil
}

// PUT makes a PUT request to the API
func (c *APIClient) PUT(endpoint string, body interface{}, result interface{}) error {
	resp, err := c.client.R().
		SetBody(body).
		SetResult(result).
		SetError(&ErrorResponse{}).
		Put(endpoint)

	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	if !resp.IsSuccess() {
		if errResp, ok := resp.Error().(*ErrorResponse); ok {
			errResp.Code = resp.StatusCode()
			return errResp
		}
		return fmt.Errorf("request failed with status %d", resp.StatusCode())
	}

	return nil
}

// DELETE makes a DELETE request to the API
func (c *APIClient) DELETE(endpoint string, result interface{}) error {
	resp, err := c.client.R().
		SetResult(result).
		SetError(&ErrorResponse{}).
		Delete(endpoint)

	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	if !resp.IsSuccess() {
		if errResp, ok := resp.Error().(*ErrorResponse); ok {
			errResp.Code = resp.StatusCode()
			return errResp
		}
		return fmt.Errorf("request failed with status %d", resp.StatusCode())
	}

	return nil
}

// WebSocketURL converts HTTP(S) URL to WebSocket URL
func (c *APIClient) WebSocketURL(path string) string {
	u, _ := url.Parse(c.baseURL)
	scheme := "ws"
	if u.Scheme == "https" {
		scheme = "wss"
	}
	return fmt.Sprintf("%s://%s%s", scheme, u.Host, path)
}

// ConnectWebSocket establishes a WebSocket connection
func (c *APIClient) ConnectWebSocket(path string) (*websocket.Conn, error) {
	wsURL := c.WebSocketURL(path)

	headers := http.Header{}
	if c.apiKey != "" {
		headers.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	}

	conn, resp, err := c.wsDialer.Dial(wsURL, headers)
	if err != nil {
		if resp != nil {
			return nil, fmt.Errorf("websocket dial failed with status %d: %w", resp.StatusCode, err)
		}
		return nil, fmt.Errorf("websocket dial failed: %w", err)
	}

	return conn, nil
}

// StreamMessage represents a streaming message from WebSocket
type StreamMessage struct {
	Type      string                 `json:"type"`
	Content   string                 `json:"content"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// StreamReader handles streaming responses
type StreamReader struct {
	conn    *websocket.Conn
	ctx     context.Context
	cancel  context.CancelFunc
	msgChan chan StreamMessage
	errChan chan error
}

// NewStreamReader creates a new stream reader
func (c *APIClient) NewStreamReader(path string) (*StreamReader, error) {
	conn, err := c.ConnectWebSocket(path)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	reader := &StreamReader{
		conn:    conn,
		ctx:     ctx,
		cancel:  cancel,
		msgChan: make(chan StreamMessage, 100),
		errChan: make(chan error, 1),
	}

	// Start reading messages
	go reader.readLoop()

	return reader, nil
}

func (sr *StreamReader) readLoop() {
	defer close(sr.msgChan)
	defer close(sr.errChan)

	for {
		select {
		case <-sr.ctx.Done():
			return
		default:
			var msg StreamMessage
			err := sr.conn.ReadJSON(&msg)
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					return
				}
				sr.errChan <- err
				return
			}

			sr.msgChan <- msg
		}
	}
}

// Messages returns the message channel
func (sr *StreamReader) Messages() <-chan StreamMessage {
	return sr.msgChan
}

// Errors returns the error channel
func (sr *StreamReader) Errors() <-chan error {
	return sr.errChan
}

// Close closes the stream reader
func (sr *StreamReader) Close() error {
	sr.cancel()
	return sr.conn.Close()
}

// HealthCheck performs a health check on the API
func (c *APIClient) HealthCheck() error {
	var result map[string]interface{}
	return c.GET("/health", &result)
}
