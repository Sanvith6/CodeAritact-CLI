package ui

import (
	"os"

	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/erikgeiser/promptkit/textinput"
)

// Confirm asks a yes/no question.
func Confirm(message string, defaultVal bool) (bool, error) {
	input := confirmation.New(message, confirmation.Yes)
	if !defaultVal {
		input = confirmation.New(message, confirmation.No)
	}

	ready, err := input.RunPrompt()
	if err != nil {
		return false, err
	}

	return ready, nil
}

// AskForString asks for text input.
func AskForString(message string, defaultVal string) (string, error) {
	input := textinput.New(message)
	input.InitialValue = defaultVal
	input.Placeholder = defaultVal

	val, err := input.RunPrompt()
	if err != nil {
		return "", err
	}
	
	if val == "" {
		return defaultVal, nil
	}

	return val, nil
}

// EnsureInteractive checks if we're connected to a terminal and allows interactive prompts.
// In CI environments, this should return false.
func IsInteractive(nonInteractiveFlag bool) bool {
	if nonInteractiveFlag {
		return false
	}
	
	fileInfo, _ := os.Stdout.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}
