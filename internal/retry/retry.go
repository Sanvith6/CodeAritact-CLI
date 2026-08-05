package retry

import (
	"context"
	"errors"
	"math/rand"
	"time"
)

// Options controls the retry behavior.
type Options struct {
	MaxAttempts    int
	BaseDelay      time.Duration
	MaxDelay       time.Duration
	Jitter         bool
	RetryableCheck func(error) bool
}

// DefaultOptions provides sensible defaults for AWS API calls and network operations.
func DefaultOptions() Options {
	return Options{
		MaxAttempts: 3,
		BaseDelay:   100 * time.Millisecond,
		MaxDelay:    2 * time.Second,
		Jitter:      true,
		RetryableCheck: func(err error) bool {
			// By default, assume all errors might be retryable if they aren't explicitly fatal.
			// Callers should provide specific implementations for their domains.
			return true
		},
	}
}

// Do executes the given function with exponential backoff and optional jitter.
func Do(ctx context.Context, opts Options, fn func(ctx context.Context) error) error {
	if opts.MaxAttempts <= 0 {
		opts.MaxAttempts = 1
	}

	var lastErr error
	delay := opts.BaseDelay

	for attempt := 1; attempt <= opts.MaxAttempts; attempt++ {
		err := fn(ctx)
		if err == nil {
			return nil
		}

		lastErr = err

		// If it's the last attempt or the error is not retryable, return immediately.
		if attempt == opts.MaxAttempts || (opts.RetryableCheck != nil && !opts.RetryableCheck(err)) {
			return lastErr
		}

		// Calculate next delay
		currentDelay := delay
		if opts.Jitter {
			// Add up to 50% jitter
			jitterAmount := time.Duration(rand.Float64() * float64(currentDelay) * 0.5)
			currentDelay += jitterAmount
		}

		select {
		case <-ctx.Done():
			// Context cancelled during retry backoff
			return errors.Join(lastErr, ctx.Err())
		case <-time.After(currentDelay):
			// Wait for next attempt
		}

		// Double the base delay for the next iteration
		delay *= 2
		if opts.MaxDelay > 0 && delay > opts.MaxDelay {
			delay = opts.MaxDelay
		}
	}

	return lastErr
}
