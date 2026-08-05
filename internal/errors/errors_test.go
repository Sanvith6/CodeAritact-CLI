package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCLIError_Error(t *testing.T) {
	err := New(ErrCodeNetwork, "Network is down")
	assert.Equal(t, "[NETWORK_UNREACHABLE] Network is down", err.Error())

	errWithCause := Wrap(ErrCodeUnknown, errors.New("underlying"), "Something failed")
	assert.Equal(t, "[UNKNOWN_ERROR] Something failed: underlying", errWithCause.Error())
}

func TestCLIError_Unwrap(t *testing.T) {
	cause := errors.New("underlying problem")
	err := Wrap(ErrCodeConfigInvalid, cause, "Invalid config file")
	
	assert.True(t, errors.Is(err, cause))
	assert.Equal(t, cause, errors.Unwrap(err))
}

func TestCLIError_Builders(t *testing.T) {
	err := New(ErrCodeToolMissing, "Tool not found").
		WithSuggestion("Install tool x").
		WithDocs("https://docs.acme.com/tool")
		
	assert.Equal(t, "Install tool x", err.Suggestion)
	assert.Equal(t, "https://docs.acme.com/tool", err.DocsURL)
}

func TestFormatForUser(t *testing.T) {
	err := New(ErrCodeArtifactAuth, "Failed to authenticate").
		WithSuggestion("Run 'rs login' to refresh your session.")
		
	out := FormatForUser(err)
	assert.Contains(t, out, "Error: Failed to authenticate")
	assert.Contains(t, out, "Suggestion: Run 'rs login' to refresh your session.")
	assert.Contains(t, out, "Code: ARTIFACT_AUTH_FAILED")
	
	// Test standard error
	stdErr := errors.New("standard go error")
	out2 := FormatForUser(stdErr)
	assert.Contains(t, out2, "Error: standard go error")
}
