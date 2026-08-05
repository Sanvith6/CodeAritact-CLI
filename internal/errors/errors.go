package errors

import (
	"fmt"
)

// ErrorCode is a machine-readable string representing the error type.
type ErrorCode string

const (
	ErrCodeNetwork         ErrorCode = "NETWORK_UNREACHABLE"
	ErrCodeToolMissing     ErrorCode = "TOOL_NOT_FOUND"
	ErrCodeToolVersion     ErrorCode = "TOOL_VERSION_MISMATCH"
	ErrCodeSSOExpired      ErrorCode = "SSO_SESSION_EXPIRED"
	ErrCodeSSOLogin        ErrorCode = "SSO_LOGIN_FAILED"
	ErrCodeSSOProfile      ErrorCode = "SSO_PROFILE_INVALID"
	ErrCodeArtifactAuth    ErrorCode = "ARTIFACT_AUTH_FAILED"
	ErrCodeArtifactRepo    ErrorCode = "ARTIFACT_REPO_NOT_FOUND"
	ErrCodeArtifactPkg     ErrorCode = "ARTIFACT_PACKAGE_NOT_FOUND"
	ErrCodeProjectNotFound ErrorCode = "PROJECT_NOT_FOUND"
	ErrCodeVenvInvalid     ErrorCode = "VENV_INVALID"
	ErrCodeConfigInvalid   ErrorCode = "CONFIG_INVALID"
	ErrCodePermission      ErrorCode = "PERMISSION_DENIED"
	ErrCodeUnknown         ErrorCode = "UNKNOWN_ERROR"
)

// CLIError represents an error with user-facing context.
type CLIError struct {
	Code       ErrorCode
	Message    string
	Cause      error
	Suggestion string
	DocsURL    string
}

// Error implements the error interface.
func (e *CLIError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap implements the error unwrap interface (Go 1.13+).
func (e *CLIError) Unwrap() error {
	return e.Cause
}

// Constructors

// New constructs a new CLIError.
func New(code ErrorCode, msg string) *CLIError {
	return &CLIError{
		Code:    code,
		Message: msg,
	}
}

// Wrap wraps an existing error into a CLIError.
func Wrap(code ErrorCode, err error, msg string) *CLIError {
	return &CLIError{
		Code:    code,
		Message: msg,
		Cause:   err,
	}
}

// WithSuggestion adds a suggestion to the error.
func (e *CLIError) WithSuggestion(suggestion string) *CLIError {
	e.Suggestion = suggestion
	return e
}

// WithDocs adds a docs URL to the error.
func (e *CLIError) WithDocs(docsURL string) *CLIError {
	e.DocsURL = docsURL
	return e
}

// FormatForUser returns a string formatted for user display, handling colour
// wrapping gracefully if passed through the UI package later.
func FormatForUser(err error) string {
	if err == nil {
		return ""
	}

	cliErr, ok := err.(*CLIError)
	if !ok {
		return fmt.Sprintf("Error: %v\nSuggestion: Please check the logs for more details.", err)
	}

	out := fmt.Sprintf("Error: %s\n", cliErr.Message)
	if cliErr.Suggestion != "" {
		out += fmt.Sprintf("Suggestion: %s\n", cliErr.Suggestion)
	}
	if cliErr.DocsURL != "" {
		out += fmt.Sprintf("Docs: %s\n", cliErr.DocsURL)
	}
	out += fmt.Sprintf("Code: %s", cliErr.Code)
	
	return out
}
