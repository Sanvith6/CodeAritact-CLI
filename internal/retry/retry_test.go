package retry

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDo_SuccessFirstTry(t *testing.T) {
	opts := DefaultOptions()
	attempts := 0

	err := Do(context.Background(), opts, func(ctx context.Context) error {
		attempts++
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, 1, attempts)
}

func TestDo_SuccessAfterRetries(t *testing.T) {
	opts := DefaultOptions()
	opts.BaseDelay = 10 * time.Millisecond // fast for tests
	attempts := 0

	err := Do(context.Background(), opts, func(ctx context.Context) error {
		attempts++
		if attempts < 3 {
			return errors.New("temporary error")
		}
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, 3, attempts)
}

func TestDo_MaxAttemptsReached(t *testing.T) {
	opts := DefaultOptions()
	opts.MaxAttempts = 2
	opts.BaseDelay = 10 * time.Millisecond
	attempts := 0

	err := Do(context.Background(), opts, func(ctx context.Context) error {
		attempts++
		return errors.New("permanent error")
	})

	assert.Error(t, err)
	assert.Equal(t, 2, attempts)
	assert.Equal(t, "permanent error", err.Error())
}

func TestDo_NotRetryable(t *testing.T) {
	opts := DefaultOptions()
	opts.MaxAttempts = 5
	opts.BaseDelay = 10 * time.Millisecond
	opts.RetryableCheck = func(err error) bool {
		return err.Error() != "fatal error"
	}

	attempts := 0
	err := Do(context.Background(), opts, func(ctx context.Context) error {
		attempts++
		return errors.New("fatal error")
	})

	assert.Error(t, err)
	assert.Equal(t, 1, attempts) // Should stop immediately
}

func TestDo_ContextCancellation(t *testing.T) {
	opts := DefaultOptions()
	opts.BaseDelay = 1 * time.Second // long delay to ensure it hits context cancel
	
	ctx, cancel := context.WithCancel(context.Background())
	
	attempts := 0
	
	// Start a goroutine to cancel the context shortly after start
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	err := Do(ctx, opts, func(c context.Context) error {
		attempts++
		return errors.New("temporary error")
	})

	assert.Error(t, err)
	assert.Equal(t, 1, attempts) // Should not proceed to second attempt due to cancellation
	assert.True(t, errors.Is(err, context.Canceled))
}
