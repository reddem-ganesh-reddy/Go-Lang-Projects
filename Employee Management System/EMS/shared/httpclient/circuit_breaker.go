package httpclient

import (
	"fmt"
	"sync"
	"time"

	"EMS/config"
)

// CircuitBreaker prevents repeated calls to an unavailable service.
type CircuitBreaker struct {
	mu sync.Mutex

	failureCount int
	maxFailures  int

	resetTimeout time.Duration
	openedAt     time.Time

	isOpen bool
}

// NewCircuitBreaker creates a circuit breaker using configuration.
func NewCircuitBreaker(cfg config.Config) *CircuitBreaker {
	return &CircuitBreaker{
		maxFailures: cfg.CircuitBreakerFailureThreshold,
		resetTimeout: time.Duration(
			cfg.CircuitBreakerResetTimeoutSeconds,
		) * time.Second,
	}
}

// Allow checks whether a request can be sent.
func (b *CircuitBreaker) Allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.isOpen {
		return nil
	}

	if time.Since(b.openedAt) >= b.resetTimeout {
		b.isOpen = false
		b.failureCount = 0
		return nil
	}

	return fmt.Errorf("circuit breaker is open; try again later")
}

// RecordSuccess resets the failure count.
func (b *CircuitBreaker) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.failureCount = 0
	b.isOpen = false
}

// RecordFailure records a failed request.
func (b *CircuitBreaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.failureCount++

	if b.failureCount >= b.maxFailures {
		b.isOpen = true
		b.openedAt = time.Now()
	}
}