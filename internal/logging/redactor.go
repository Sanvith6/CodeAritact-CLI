package logging

import (
	"context"
	"log/slog"
	"strings"
)

// RedactingHandler is an slog.Handler that redacts sensitive values.
type RedactingHandler struct {
	handler slog.Handler
}

// NewRedactingHandler creates a new RedactingHandler.
func NewRedactingHandler(h slog.Handler) *RedactingHandler {
	return &RedactingHandler{
		handler: h,
	}
}

// Enabled implements slog.Handler.
func (h *RedactingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

// Handle implements slog.Handler.
func (h *RedactingHandler) Handle(ctx context.Context, r slog.Record) error {
	newRecord := slog.NewRecord(r.Time, r.Level, r.Message, r.PC)
	
	r.Attrs(func(a slog.Attr) bool {
		newRecord.AddAttrs(redactAttr(a))
		return true
	})

	return h.handler.Handle(ctx, newRecord)
}

// WithAttrs implements slog.Handler.
func (h *RedactingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	redactedAttrs := make([]slog.Attr, len(attrs))
	for i, a := range attrs {
		redactedAttrs[i] = redactAttr(a)
	}
	return NewRedactingHandler(h.handler.WithAttrs(redactedAttrs))
}

// WithGroup implements slog.Handler.
func (h *RedactingHandler) WithGroup(name string) slog.Handler {
	return NewRedactingHandler(h.handler.WithGroup(name))
}

func redactAttr(a slog.Attr) slog.Attr {
	// If it's a group, recursively redact
	if a.Value.Kind() == slog.KindGroup {
		attrs := a.Value.Group()
		redactedAttrs := make([]slog.Attr, len(attrs))
		for i, attr := range attrs {
			redactedAttrs[i] = redactAttr(attr)
		}
		return slog.Attr{Key: a.Key, Value: slog.GroupValue(redactedAttrs...)}
	}

	keyLower := strings.ToLower(a.Key)
	if strings.Contains(keyLower, "token") ||
		strings.Contains(keyLower, "secret") ||
		strings.Contains(keyLower, "password") ||
		strings.Contains(keyLower, "credential") ||
		strings.Contains(keyLower, "authorization") {
		return slog.String(a.Key, "[REDACTED]")
	}
	
	return a
}
