package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/acme-corp/rs-cli/internal/platform"
	"github.com/acme-corp/rs-cli/internal/state"
)

// RepositoryService manages interaction with package repositories.
type RepositoryService struct {
	auth     *AuthService
	provider platform.RepositoryProvider
	state    *state.Manager
	logger   *slog.Logger
}

// NewRepositoryService creates a new RepositoryService.
func NewRepositoryService(
	auth *AuthService,
	provider platform.RepositoryProvider,
	state *state.Manager,
	logger *slog.Logger,
) *RepositoryService {
	return &RepositoryService{
		auth:     auth,
		provider: provider,
		state:    state,
		logger:   logger,
	}
}

// ValidateAccess ensures the user has access to the repository.
func (s *RepositoryService) ValidateAccess(ctx context.Context) error {
	if err := s.auth.EnsureAuthenticated(ctx); err != nil {
		return fmt.Errorf("authentication required: %w", err)
	}

	return s.provider.ValidateAccess(ctx)
}

// GetAuthInfo returns the endpoint and a short-lived token for the repository.
func (s *RepositoryService) GetAuthInfo(ctx context.Context) (endpoint, token string, err error) {
	if err := s.auth.EnsureAuthenticated(ctx); err != nil {
		return "", "", fmt.Errorf("authentication required: %w", err)
	}

	endpoint, err = s.provider.GetEndpoint(ctx)
	if err != nil {
		return "", "", err
	}

	token, expiry, err := s.provider.GetAuthToken(ctx)
	if err != nil {
		return "", "", err
	}
	
	// Update state with token expiry if we care to cache it, though this is usually managed
	// by the provider or not needed to be persisted if we fetch on demand.
	s.logger.Debug("Retrieved auth token", slog.Time("expiry", expiry))

	return endpoint, token, nil
}

// VerifyPackage ensures the package exists in the repository.
func (s *RepositoryService) VerifyPackage(ctx context.Context, pkgName string) error {
	versions, err := s.provider.ListVersions(ctx, pkgName)
	if err != nil {
		return fmt.Errorf("failed to list versions for package %s: %w", pkgName, err)
	}
	if len(versions) == 0 {
		return fmt.Errorf("package %s not found in repository", pkgName)
	}
	return nil
}
