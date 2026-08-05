package logging

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedactingHandler(t *testing.T) {
	var buf bytes.Buffer
	
	// Create JSON handler writing to buffer
	jsonHandler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	redactingHandler := NewRedactingHandler(jsonHandler)
	logger := slog.New(redactingHandler)

	logger.Info("Login attempt", 
		slog.String("user", "alice"),
		slog.String("aws_token", "SECRET-12345"),
		slog.String("authorization", "Bearer xyz"),
		slog.String("password", "p@ssword1"),
		slog.Group("credentials", 
			slog.String("secretAccessKey", "very-secret"),
		),
	)

	output := buf.String()
	
	// Safe values should be present
	assert.Contains(t, output, `"user":"alice"`)
	
	// Sensitive values should be redacted
	assert.NotContains(t, output, "SECRET-12345")
	assert.Contains(t, output, `"aws_token":"[REDACTED]"`)
	
	assert.NotContains(t, output, "Bearer xyz")
	assert.Contains(t, output, `"authorization":"[REDACTED]"`)
	
	assert.NotContains(t, output, "p@ssword1")
	assert.Contains(t, output, `"password":"[REDACTED]"`)
	
	// Nested groups should be redacted
	assert.NotContains(t, output, "very-secret")
	assert.Contains(t, output, `"secretAccessKey":"[REDACTED]"`)
}

func TestSetup(t *testing.T) {
	logger := Setup("debug", false)
	assert.NotNil(t, logger)
	
	ctx := context.Background()
	// Should be enabled at debug
	assert.True(t, logger.Enabled(ctx, slog.LevelDebug))
	assert.True(t, logger.Enabled(ctx, slog.LevelInfo))
}
