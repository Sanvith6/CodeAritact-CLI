package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/acme-corp/rs-cli/internal/platform"
	"github.com/acme-corp/rs-cli/internal/state"
)

// AuthService orchestrates authentication workflows.
type AuthService struct {
	provider platform.AuthProvider
	state    *state.Manager
	logger   *slog.Logger
}

// NewAuthService creates a new AuthService.
func NewAuthService(provider platform.AuthProvider, state *state.Manager, logger *slog.Logger) *AuthService {
	return &AuthService{
		provider: provider,
		state:    state,
		logger:   logger,
	}
}

// EnsureAuthenticated ensures the user has a valid authenticated session.
func (s *AuthService) EnsureAuthenticated(ctx context.Context) error {
	if s.IsAuthenticated(ctx) {
		return nil
	}
	
	s.logger.Info("Session expired or missing, initiating login...")
	return s.Login(ctx, false)
}

// Login performs the login flow.
func (s *AuthService) Login(ctx context.Context, force bool) error {
	if !force && s.IsAuthenticated(ctx) {
		s.logger.Info("Already authenticated with a valid session.")
		return nil
	}

	if err := s.provider.Login(ctx); err != nil {
		return err
	}

	// Wait briefly for credentials to be fully available
	time.Sleep(1 * time.Second)

	identity, err := s.provider.GetIdentity(ctx)
	if err != nil {
		return err
	}

	// Update state
	err = s.state.Update(func(st *state.State) {
		st.LastLogin = time.Now()
	})
	if err != nil {
		s.logger.Warn("Failed to update login timestamp in state", slog.String("error", err.Error()))
	}

	s.logger.Info("Login successful", 
		slog.String("role", identity.RoleName),
		slog.String("account", identity.AccountID),
	)

	return nil
}

// Logout performs the logout flow.
func (s *AuthService) Logout(ctx context.Context) error {
	if err := s.provider.Logout(ctx); err != nil {
		return err
	}
	
	s.logger.Info("Logged out successfully.")
	return nil
}

// IsAuthenticated checks if the current session is valid.
func (s *AuthService) IsAuthenticated(ctx context.Context) bool {
	return s.provider.IsSessionValid(ctx)
}

// GetIdentity retrieves the current identity.
func (s *AuthService) GetIdentity(ctx context.Context) (*platform.Identity, error) {
	return s.provider.GetIdentity(ctx)
}
