package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPError is returned when an external service responds with a non-2xx status.
// Mirrors Java ServiceRequestRepository HttpClientErrorException handling.
type HTTPError struct {
	StatusCode int
	Body       string
	URL        string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("http call failed: url=%s status=%d body=%s", e.URL, e.StatusCode, e.Body)
}

// Client performs JSON HTTP calls to DIGIT external services.
type Client struct {
	http *http.Client
}

// New creates a Client with the given timeout.
func New(timeout time.Duration) *Client {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return &Client{
		http: &http.Client{Timeout: timeout},
	}
}

// PostJSON sends a JSON POST and decodes the response into target.
// target may be nil to discard the response body.
func (c *Client) PostJSON(ctx context.Context, url string, payload any, target any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return &HTTPError{
			StatusCode: resp.StatusCode,
			Body:       string(respBody),
			URL:        url,
		}
	}

	if target == nil || len(respBody) == 0 {
		return nil
	}
	if err := json.Unmarshal(respBody, target); err != nil {
		return fmt.Errorf("unmarshal response: %w", err)
	}
	return nil
}

// PostJSONMap posts JSON and returns the response as a generic map (Java Map.class parity).
func (c *Client) PostJSONMap(ctx context.Context, url string, payload any) (map[string]any, error) {
	var out map[string]any
	if err := c.PostJSON(ctx, url, payload, &out); err != nil {
		return nil, err
	}
	return out, nil
}
