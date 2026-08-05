package telemetry

import (
	"context"
	"time"
)

// Event represents a tracked action or state in the CLI.
type Event struct {
	Name       string
	Properties map[string]string
	Timestamp  time.Time
}

// Tracker defines the interface for recording telemetry events.
type Tracker interface {
	// Track records an event. It should not block the main execution flow.
	Track(event Event)
	
	// Flush forces any buffered events to be sent.
	Flush(ctx context.Context) error
	
	// Close cleans up any resources used by the tracker.
	Close() error
}
