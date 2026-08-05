package services

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/acme-corp/rs-cli/internal/platform"
)

// EnvironmentService handles developer machine environment configuration.
type EnvironmentService struct {
	pInfo  *platform.Info
	logger *slog.Logger
}

// NewEnvironmentService creates a new EnvironmentService.
func NewEnvironmentService(
	pInfo *platform.Info,
	logger *slog.Logger,
) *EnvironmentService {
	return &EnvironmentService{
		pInfo:  pInfo,
		logger: logger,
	}
}

// CheckShellConfiguration checks if the rs CLI is in the PATH and shell is configured.
func (s *EnvironmentService) CheckShellConfiguration(ctx context.Context) error {
	// A simplistic check to see if we're in PATH
	pathEnv := os.Getenv("PATH")
	
	rsBinDir := filepath.Join(s.pInfo.HomeDir, ".local", "bin")
	if s.pInfo.OS == "windows" {
		rsBinDir = filepath.Join(s.pInfo.HomeDir, "AppData", "Local", "Programs", "rs-cli", "bin")
	}

	if !strings.Contains(pathEnv, rsBinDir) {
		s.logger.Warn("CLI might not be in your PATH.", slog.String("dir", rsBinDir))
		// In a real CLI we would offer to append it to ~/.bashrc or similar
	}
	
	return nil
}

// EnsureConfigDir ensures the ~/.rs directory exists.
func (s *EnvironmentService) EnsureConfigDir(ctx context.Context) error {
	dir := s.pInfo.RSConfigDir
	
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory %s: %w", dir, err)
	}
	
	return nil
}
