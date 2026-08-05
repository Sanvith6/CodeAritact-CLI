package services

import (
	"context"
	"log/slog"

	"github.com/acme-corp/rs-cli/internal/platform"
)

// PackageService handles package operations.
type PackageService struct {
	repo   *RepositoryService
	pm     platform.PackageManagerProvider
	logger *slog.Logger
}

// NewPackageService creates a new PackageService.
func NewPackageService(
	repo *RepositoryService,
	pm platform.PackageManagerProvider,
	logger *slog.Logger,
) *PackageService {
	return &PackageService{
		repo:   repo,
		pm:     pm,
		logger: logger,
	}
}

// Install installs a package into the current environment.
func (s *PackageService) Install(ctx context.Context, pkg, version string) error {
	// First ensure we have access and configure the package manager
	endpoint, token, err := s.repo.GetAuthInfo(ctx)
	if err != nil {
		return err
	}

	indexName := "internal" // This would usually come from config

	err = s.pm.ConfigureAuth(ctx, indexName, endpoint, token)
	if err != nil {
		return err
	}
	// Always try to clean up auth from the environment afterwards
	defer func() {
		_ = s.pm.ClearAuth(ctx, indexName)
	}()

	s.logger.Info("Installing package...", slog.String("package", pkg), slog.String("version", version))
	
	if err := s.pm.Install(ctx, pkg, version); err != nil {
		return err
	}
	
	s.logger.Info("Package installed successfully.")
	return nil
}

// Sync synchronizes dependencies in the current project.
func (s *PackageService) Sync(ctx context.Context) error {
	endpoint, token, err := s.repo.GetAuthInfo(ctx)
	if err != nil {
		return err
	}

	indexName := "internal"

	err = s.pm.ConfigureAuth(ctx, indexName, endpoint, token)
	if err != nil {
		return err
	}
	defer func() {
		_ = s.pm.ClearAuth(ctx, indexName)
	}()

	s.logger.Info("Synchronizing dependencies...")
	
	if err := s.pm.Sync(ctx); err != nil {
		return err
	}
	
	s.logger.Info("Dependencies synchronized successfully.")
	return nil
}
