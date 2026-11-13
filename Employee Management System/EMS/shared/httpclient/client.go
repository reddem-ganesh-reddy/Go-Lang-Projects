package httpclient

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"EMS/config"
	"EMS/shared/middleware"
)

// Client handles HTTP communication between services.
type Client struct {
	httpClient *http.Client
	maxRetries int
	retryDelay time.Duration
	breaker    *CircuitBreaker
}

// New creates an HTTP client using application configuration.
func New(cfg config.Config) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: time.Duration(
				cfg.HTTPClientTimeoutSeconds,
			) * time.Second,
		},
		maxRetries: cfg.HTTPClientMaxRetries,
		retryDelay: time.Duration(
			cfg.HTTPClientRetryDelayMillis,
		) * time.Millisecond,
		breaker: NewCircuitBreaker(cfg),
	}
}

// Do sends an HTTP request with timeout, retry and circuit breaker support.
func (c *Client) Do(
	ctx context.Context,
	method string,
	requestURL string,
	body []byte,
	headers map[string]string,
) (*http.Response, error) {

	if err := c.breaker.Allow(); err != nil {
		return nil, err
	}

	// Retry only methods that are normally safe to repeat.
	retryAllowed := method == http.MethodGet ||
		method == http.MethodPut ||
		method == http.MethodDelete ||
		method == http.MethodHead

	attempts := 1
	if retryAllowed {
		attempts += c.maxRetries
	}

	for attempt := 0; attempt < attempts; attempt++ {
		req, err := http.NewRequestWithContext(
			ctx,
			method,
			requestURL,
			bytes.NewReader(body),
		)
		if err != nil {
			return nil, fmt.Errorf("create HTTP request: %w", err)
		}

		for name, value := range headers {
			req.Header.Set(name, value)
		}

		// Propagate Request ID.
		requestID := middleware.RequestIDFromContext(ctx)
		if requestID != "" {
			req.Header.Set("X-Request-ID", requestID)
		}

		// Propagate the same Correlation ID.
		correlationID := middleware.FromContext(ctx)
		if correlationID != "" {
			req.Header.Set("X-Correlation-ID", correlationID)
		}

		resp, err := c.httpClient.Do(req)

		if err != nil {
			c.breaker.RecordFailure()

			if attempt == attempts-1 {
				return nil, fmt.Errorf("send HTTP request: %w", err)
			}

			if err := waitForRetry(ctx, c.retryDelay, attempt); err != nil {
				return nil, err
			}

			continue
		}

		// Retry server errors for retryable methods.
		if resp.StatusCode >= 500 {
			if attempt < attempts-1 {
				resp.Body.Close()

				if err := waitForRetry(ctx, c.retryDelay, attempt); err != nil {
					return nil, err
				}

				continue
			}

			c.breaker.RecordFailure()
			return resp, nil
		}

		// A response other than HTTP 5xx means the service responded.
		c.breaker.RecordSuccess()

		return resp, nil
	}

	return nil, fmt.Errorf("HTTP request failed")
}

// waitForRetry waits before sending another request.
func waitForRetry(
	ctx context.Context,
	delay time.Duration,
	attempt int,
) error {

	waitTime := time.Duration(attempt+1) * delay

	timer := time.NewTimer(waitTime)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}