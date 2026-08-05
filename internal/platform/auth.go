package platform

import (
	"context"
)

// AuthProvider defines the interface for authenticating users.
// This decouples the CLI from a specific identity provider (e.g., AWS SSO, Entra ID).
type AuthProvider interface {
	// Name returns the provider name (e.g., "aws-sso").
	Name() string
	
	// EnsureProfile ensures that the local environment has the required configuration profile.
	EnsureProfile(ctx context.Context) error
	
	// ValidateProfile verifies that the existing profile matches the expected configuration.
	ValidateProfile(ctx context.Context) error
	
	// Login initiates an interactive or programmatic login flow.
	Login(ctx context.Context) error
	
	// Logout removes the current session.
	Logout(ctx context.Context) error
	
	// IsSessionValid checks if the user is currently authenticated with a valid session.
	IsSessionValid(ctx context.Context) bool
	
	// GetIdentity returns information about the authenticated user.
	GetIdentity(ctx context.Context) (*Identity, error)
}
