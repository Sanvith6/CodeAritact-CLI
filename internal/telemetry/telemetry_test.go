package telemetry

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNoOpTracker(t *testing.T) {
	tracker := NewNoOpTracker()
	
	// Should not panic or error
	tracker.Track(Event{
		Name: "test_event",
		Properties: map[string]string{"key": "value"},
		Timestamp: time.Now(),
	})
	
	err := tracker.Flush(context.Background())
	assert.NoError(t, err)
	
	err = tracker.Close()
	assert.NoError(t, err)
}
