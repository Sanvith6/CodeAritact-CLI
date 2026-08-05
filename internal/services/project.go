package services

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/acme-corp/rs-cli/internal/platform"
)

// ProjectService handles project-level operations.
type ProjectService struct {
	pm     platform.PackageManagerProvider
	logger *slog.Logger
}

// NewProjectService creates a new ProjectService.
func NewProjectService(
	pm platform.PackageManagerProvider,
	logger *slog.Logger,
) *ProjectService {
	return &ProjectService{
		pm:     pm,
		logger: logger,
	}
}

// Initialize initializes a new project in the current directory.
func (s *ProjectService) Initialize(ctx context.Context) error {
	// A basic initialization might involve creating a venv and a pyproject.toml stub
	
	s.logger.Info("Initializing project...")
	
	if err := s.pm.CreateVenv(ctx, ".venv"); err != nil {
		return fmt.Errorf("failed to create virtual environment: %w", err)
	}
	
	// Create a dummy pyproject.toml if it doesn't exist
	if _, err := os.Stat("pyproject.toml"); os.IsNotExist(err) {
		content := `[project]
name = "my-project"
version = "0.1.0"
description = "Add your description here"
readme = "README.md"
requires-python = ">=3.9"
dependencies = []
`
		if err := os.WriteFile("pyproject.toml", []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create pyproject.toml: %w", err)
		}
		s.logger.Info("Created pyproject.toml")
	}

	s.logger.Info("Project initialized successfully.")
	return nil
}

// Verify ensures the project is valid.
func (s *ProjectService) Verify(ctx context.Context) error {
	if _, err := os.Stat("pyproject.toml"); os.IsNotExist(err) {
		return fmt.Errorf("no pyproject.toml found. Are you in a project directory?")
	}
	return nil
}
