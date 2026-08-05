package ui

import (
	"fmt"
	"strings"
	"time"
)

// ProgressRenderer simulates a progress bar or continuous progress for slow operations.
type ProgressRenderer struct {
	Total    int
	Current  int
	Message  string
	barWidth int
}

// NewProgress creates a new ProgressRenderer.
func NewProgress(total int, message string) *ProgressRenderer {
	return &ProgressRenderer{
		Total:    total,
		Message:  message,
		barWidth: 40,
	}
}

// Render returns a string representation of the progress bar.
func (p *ProgressRenderer) Render() string {
	percent := 0.0
	if p.Total > 0 {
		percent = float64(p.Current) / float64(p.Total)
	}

	if percent > 1.0 {
		percent = 1.0
	}

	filled := int(percent * float64(p.barWidth))
	empty := p.barWidth - filled

	bar := StyleBrand.Render(strings.Repeat("█", filled)) + StyleMuted.Render(strings.Repeat("░", empty))
	
	return fmt.Sprintf("%s [%s] %d/%d (%.0f%%)", p.Message, bar, p.Current, p.Total, percent*100)
}

// Example usage in an event loop or similar. In a real advanced CLI, 
// we might use charmbracelet/bubbletea for true async progress bars.
func (p *ProgressRenderer) PrintUpdate() {
	// \r clears the line and returns cursor to start.
	// We need enough spaces to clear previous longer lines.
	fmt.Printf("\r\033[K%s", p.Render())
}

// SimulateProgress is a helper for showing a fake progress bar for UX purposes.
func SimulateProgress(message string, duration time.Duration) {
	steps := 50
	stepDuration := duration / time.Duration(steps)
	
	p := NewProgress(steps, message)
	
	for i := 0; i <= steps; i++ {
		p.Current = i
		p.PrintUpdate()
		time.Sleep(stepDuration)
	}
	fmt.Println() // Newline at the end
}
