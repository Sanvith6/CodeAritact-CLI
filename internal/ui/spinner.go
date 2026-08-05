package ui

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
)

// Spinner wraps the briandowns/spinner to provide a consistent experience.
type Spinner struct {
	s *spinner.Spinner
}

// NewSpinner creates a new terminal spinner.
func NewSpinner(message string) *Spinner {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond) // Dot spinner
	s.Color("blue", "bold")
	s.Suffix = " " + message
	return &Spinner{s: s}
}

// Start starts the spinner.
func (s *Spinner) Start() {
	s.s.Start()
}

// Stop stops the spinner and clears it.
func (s *Spinner) Stop() {
	s.s.Stop()
}

// StopMessage stops the spinner and leaves a final message.
func (s *Spinner) StopMessage(msg string) {
	s.s.FinalMSG = msg + "\n"
	s.s.Stop()
}

// StopSuccess stops the spinner with a success message.
func (s *Spinner) StopSuccess(msg string) {
	s.s.FinalMSG = fmt.Sprintf("%s %s\n", IconCheck, msg)
	s.s.Stop()
}

// StopError stops the spinner with an error message.
func (s *Spinner) StopError(msg string) {
	s.s.FinalMSG = fmt.Sprintf("%s %s\n", IconCross, StyleError.Render(msg))
	s.s.Stop()
}

// UpdateMessage updates the text next to the spinner.
func (s *Spinner) UpdateMessage(msg string) {
	s.s.Suffix = " " + msg
}
