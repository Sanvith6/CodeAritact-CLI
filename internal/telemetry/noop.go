package telemetry

import (
	"context"
)

// NoOpTracker is a telemetry tracker that does nothing.
// It serves as the default implementation until an organization plugs in a real one.
type NoOpTracker struct{}

// NewNoOpTracker creates a new NoOpTracker.
func NewNoOpTracker() *NoOpTracker {
	return &NoOpTracker{}
}

// Track implements Tracker.Track but does nothing.
func (t *NoOpTracker) Track(event Event) {}

// Flush implements Tracker.Flush but does nothing.
func (t *NoOpTracker) Flush(ctx context.Context) error {
	return nil
}

// Close implements Tracker.Close but does nothing.
func (t *NoOpTracker) Close() error {
	return nil
}
