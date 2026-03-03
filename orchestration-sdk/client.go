package orchestrationsdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client is the Orchestration SDK client
type Client struct {
	baseURL    string
	httpClient *http.Client
	retries    int
	retryWait  time.Duration
	authToken  string
	apiKey     string
	userAgent  string
	headers    map[string]string
	debug      bool
}

// NewClient creates a new Orchestration SDK client with the given options
func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		baseURL: DefaultBaseURL,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		retries:   DefaultRetries,
		retryWait: DefaultRetryWait,
		userAgent: DefaultUserAgent,
		headers:   make(map[string]string),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// doRequest performs an HTTP request with retries
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	// Marshal body once and store the bytes for reuse in retries
	var jsonBody []byte
	if body != nil {
		var err error
		jsonBody, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	fullURL := c.baseURL + path

	var lastErr error
	for attempt := 0; attempt <= c.retries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(c.retryWait):
			}
		}

		// Create a fresh reader for each attempt to avoid EOF issues on retries
		var bodyReader io.Reader
		if jsonBody != nil {
			bodyReader = bytes.NewReader(jsonBody)
		}

		req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		c.setHeaders(req)

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request failed: %w", err)
			continue
		}

		err = c.handleResponse(resp, result)
		if err != nil {
			if IsRetryable(err) {
				lastErr = err
				continue
			}
			return err
		}

		return nil
	}

	return lastErr
}

// setHeaders sets the request headers
func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	if c.authToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.authToken)
	}

	if c.apiKey != "" {
		req.Header.Set("X-API-Key", c.apiKey)
	}

	for key, value := range c.headers {
		req.Header.Set(key, value)
	}
}

// handleResponse handles the HTTP response
func (c *Client) handleResponse(resp *http.Response, result interface{}) error {
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return c.parseErrorResponse(resp.StatusCode, bodyBytes)
	}

	if result != nil && len(bodyBytes) > 0 {
		// Try to parse as APIResponse first
		var apiResp APIResponse
		if err := json.Unmarshal(bodyBytes, &apiResp); err == nil && apiResp.Data != nil {
			// Re-marshal the data field and unmarshal into result
			dataBytes, err := json.Marshal(apiResp.Data)
			if err != nil {
				return fmt.Errorf("failed to marshal data: %w", err)
			}
			if err := json.Unmarshal(dataBytes, result); err != nil {
				return fmt.Errorf("failed to unmarshal data: %w", err)
			}
			return nil
		}

		// Fallback: try to unmarshal directly
		if err := json.Unmarshal(bodyBytes, result); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}
	}

	return nil
}

// parseErrorResponse parses an error response from the API
func (c *Client) parseErrorResponse(statusCode int, body []byte) error {
	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err == nil && apiResp.Error != "" {
		apiErr := &APIError{
			StatusCode: statusCode,
			Code:       fmt.Sprintf("%d", apiResp.Code),
			Message:    apiResp.Message,
			Details:    apiResp.Error,
		}

		// Check if retryable
		if statusCode == 429 || statusCode >= 500 {
			return NewRetryableError(apiErr, 5)
		}

		return apiErr
	}

	// Generic error based on status code
	switch statusCode {
	case 400:
		return NewAPIError(statusCode, "bad_request", "Bad request")
	case 401:
		return ErrUnauthorized
	case 403:
		return ErrForbidden
	case 404:
		return ErrNotFound
	case 409:
		return ErrConflict
	case 429:
		return NewRetryableError(NewAPIError(statusCode, "rate_limit", "Rate limit exceeded"), 5)
	default:
		if statusCode >= 500 {
			return NewRetryableError(NewAPIError(statusCode, "server_error", "Server error"), 5)
		}
		return NewAPIError(statusCode, "unknown", string(body))
	}
}

// Ping checks if the API is reachable
func (c *Client) Ping(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/health", nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API health check failed with status: %d", resp.StatusCode)
	}

	return nil
}

// SetBaseURL updates the base URL (useful for testing)
func (c *Client) SetBaseURL(baseURL string) {
	c.baseURL = baseURL
}

// SetAuthToken updates the auth token
func (c *Client) SetAuthToken(token string) {
	c.authToken = token
}

// CreateOrchestration creates a new orchestration
func (c *Client) CreateOrchestration(ctx context.Context, req *OrchestrationRequest) (*OrchestrationResponse, error) {
	var resp OrchestrationResponse
	err := c.doRequest(ctx, http.MethodPost, PathOrchestration, req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetOrchestrationStatus gets the status of an orchestration
func (c *Client) GetOrchestrationStatus(ctx context.Context, orchestrationID string) (*OrchestrationStatusResponse, error) {
	path := fmt.Sprintf(PathOrchestrationStatus, orchestrationID)
	var resp OrchestrationStatusResponse
	err := c.doRequest(ctx, http.MethodGet, path, nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListOrchestrations lists all orchestrations
func (c *Client) ListOrchestrations(ctx context.Context, opts *ListOptions) (*OrchestrationListResponse, error) {
	path := PathOrchestration
	if opts != nil {
		params := url.Values{}
		if opts.Page > 0 {
			params.Set("page", fmt.Sprintf("%d", opts.Page))
		}
		if opts.Limit > 0 {
			params.Set("limit", fmt.Sprintf("%d", opts.Limit))
		}
		if opts.SortBy != "" {
			params.Set("sort_by", opts.SortBy)
		}
		if opts.SortOrder != "" {
			params.Set("sort_order", opts.SortOrder)
		}
		if len(params) > 0 {
			path = path + "?" + params.Encode()
		}
	}

	var resp OrchestrationListResponse
	err := c.doRequest(ctx, http.MethodGet, path, nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// CancelOrchestration cancels an orchestration
func (c *Client) CancelOrchestration(ctx context.Context, orchestrationID string) error {
	path := fmt.Sprintf(PathOrchestrationCancel, orchestrationID)
	return c.doRequest(ctx, http.MethodPost, path, nil, nil)
}

// ExecuteCommand executes a command on a VM
func (c *Client) ExecuteCommand(ctx context.Context, vmID string, command []string) (*ExecuteCommandResponse, error) {
	path := fmt.Sprintf(PathExecuteCommand, vmID)
	req := &ExecuteCommandRequest{Command: command}
	var resp ExecuteCommandResponse
	err := c.doRequest(ctx, http.MethodPost, path, req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// StartVM starts a VM
func (c *Client) StartVM(ctx context.Context, vmID string) (*OrchestrationResponse, error) {
	path := fmt.Sprintf(PathStartVM, vmID)
	var resp OrchestrationResponse
	err := c.doRequest(ctx, http.MethodPost, path, nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// StopVM stops a VM
func (c *Client) StopVM(ctx context.Context, vmID string) (*OrchestrationResponse, error) {
	path := fmt.Sprintf(PathStopVM, vmID)
	var resp OrchestrationResponse
	err := c.doRequest(ctx, http.MethodPost, path, nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// RemoveVM removes a VM
func (c *Client) RemoveVM(ctx context.Context, vmID string) (*OrchestrationResponse, error) {
	path := fmt.Sprintf(PathRemoveVM, vmID)
	var resp OrchestrationResponse
	err := c.doRequest(ctx, http.MethodDelete, path, nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetTemplate retrieves a template by ID
func (c *Client) GetTemplate(ctx context.Context, templateID string) (*Template, error) {
	path := fmt.Sprintf(PathTemplateByID, templateID)
	var resp Template
	err := c.doRequest(ctx, http.MethodGet, path, nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListTemplates lists all templates with optional filtering
func (c *Client) ListTemplates(ctx context.Context, opts *TemplateListOptions) (*TemplateListResponse, error) {
	path := PathTemplates
	if opts != nil {
		params := url.Values{}
		if opts.Page > 0 {
			params.Set("page", fmt.Sprintf("%d", opts.Page))
		}
		if opts.Limit > 0 {
			params.Set("limit", fmt.Sprintf("%d", opts.Limit))
		}
		if opts.SortBy != "" {
			params.Set("sort_by", opts.SortBy)
		}
		if opts.SortOrder != "" {
			params.Set("sort_order", opts.SortOrder)
		}
		if opts.Type != "" {
			params.Set("type", string(opts.Type))
		}
		if opts.OSType != "" {
			params.Set("os_type", string(opts.OSType))
		}
		if opts.Architecture != "" {
			params.Set("architecture", opts.Architecture)
		}
		if opts.Name != "" {
			params.Set("name", opts.Name)
		}
		if len(params) > 0 {
			path = path + "?" + params.Encode()
		}
	}

	var resp TemplateListResponse
	err := c.doRequest(ctx, http.MethodGet, path, nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
