package uv

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	rsconfig "github.com/acme-corp/rs-cli/internal/config"
	"github.com/acme-corp/rs-cli/internal/platform"
	"github.com/acme-corp/rs-cli/internal/runner"
)

// UVPackageManager implements platform.PackageManagerProvider for uv.
type UVPackageManager struct {
	cfg    *rsconfig.Config
	pInfo  *platform.Info
	runner runner.Runner
}

// NewUVPackageManager creates a new UVPackageManager.
func NewUVPackageManager(cfg *rsconfig.Config, pInfo *platform.Info, r runner.Runner) *UVPackageManager {
	return &UVPackageManager{
		cfg:    cfg,
		pInfo:  pInfo,
		runner: runner.NewUVRunner(r),
	}
}

// Name returns the provider name.
func (u *UVPackageManager) Name() string {
	return "uv"
}

// ConfigureAuth sets up the required environment variables for uv to authenticate,
// and optionally writes to uv.toml for global index configuration.
func (u *UVPackageManager) ConfigureAuth(ctx context.Context, indexName, repoURL, token string) error {
	// For uv, the idiomatic way to configure dynamic auth is via env vars.
	// We use the pattern UV_INDEX_<NAME>_USERNAME and UV_INDEX_<NAME>_PASSWORD
	
	safeIndexName := strings.ToUpper(strings.ReplaceAll(indexName, "-", "_"))
	
	userEnvVar := fmt.Sprintf("UV_INDEX_%s_USERNAME", safeIndexName)
	passEnvVar := fmt.Sprintf("UV_INDEX_%s_PASSWORD", safeIndexName)
	
	// Set it for the current process so subsequent commands (like uv sync) inherit it
	os.Setenv(userEnvVar, "aws")
	os.Setenv(passEnvVar, token)
	
	// Note: We avoid writing the token to uv.toml because it's a short-lived token (max 12h)
	// and writing it to disk is a security risk. Instead, we rely on the `rs` CLI wrapping `uv`
	// or developers running `rs login` to set it in their shell session.
	// For now, we ensure the index URL is known.
	
	return u.ensureGlobalIndexConfig(indexName, repoURL)
}

// ClearAuth removes the authentication from the current process environment.
func (u *UVPackageManager) ClearAuth(ctx context.Context, indexName string) error {
	safeIndexName := strings.ToUpper(strings.ReplaceAll(indexName, "-", "_"))
	os.Unsetenv(fmt.Sprintf("UV_INDEX_%s_USERNAME", safeIndexName))
	os.Unsetenv(fmt.Sprintf("UV_INDEX_%s_PASSWORD", safeIndexName))
	return nil
}

// Install installs a package.
func (u *UVPackageManager) Install(ctx context.Context, pkg, version string) error {
	pkgSpec := pkg
	if version != "" {
		pkgSpec = fmt.Sprintf("%s==%s", pkg, version)
	}

	var cmd string
	var args []string
	
	if u.cfg.Package.InstallStrategy == "pip-install" {
		cmd = "uv"
		args = []string{"pip", "install", pkgSpec}
	} else {
		cmd = "uv"
		args = []string{"add", pkgSpec}
	}

	_, err := u.runner.Run(ctx, cmd, args...)
	if err != nil {
		return fmt.Errorf("failed to install package %s: %w", pkg, err)
	}

	return nil
}

// Sync runs `uv sync`.
func (u *UVPackageManager) Sync(ctx context.Context) error {
	_, err := u.runner.Run(ctx, "uv", "sync")
	if err != nil {
		return fmt.Errorf("failed to sync dependencies: %w", err)
	}
	return nil
}

// CreateVenv creates a virtual environment using `uv venv`.
func (u *UVPackageManager) CreateVenv(ctx context.Context, path string) error {
	args := []string{"venv"}
	if path != "" && path != ".venv" {
		args = append(args, path)
	}
	
	_, err := u.runner.Run(ctx, "uv", args...)
	if err != nil {
		return fmt.Errorf("failed to create virtual environment: %w", err)
	}
	return nil
}

// IsVenvValid checks if a `.venv` directory exists and has a python executable.
func (u *UVPackageManager) IsVenvValid(ctx context.Context, path string) bool {
	if path == "" {
		path = ".venv"
	}
	
	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		return false
	}
	
	// Check for bin/python (Unix) or Scripts/python.exe (Windows)
	unixPython := filepath.Join(path, "bin", "python")
	winPython := filepath.Join(path, "Scripts", "python.exe")
	
	if _, err := os.Stat(unixPython); err == nil {
		return true
	}
	
	if _, err := os.Stat(winPython); err == nil {
		return true
	}
	
	return false
}

func (u *UVPackageManager) ensureGlobalIndexConfig(indexName, repoURL string) error {
	// Ensures the uv.toml file contains the index definition (without credentials).
	// Because TOML parsing in Go is complex and we want to keep dependencies low,
	// we use a simple approach: if we don't find the index URL in the file, we append it.
	
	configPath := filepath.Join(u.pInfo.UVConfigDir, "uv.toml")
	
	// Create directory if it doesn't exist
	if err := os.MkdirAll(u.pInfo.UVConfigDir, 0700); err != nil {
		return err
	}
	
	var content string
	data, err := os.ReadFile(configPath)
	if err == nil {
		content = string(data)
	}
	
	// Simple string matching to avoid adding it repeatedly
	if strings.Contains(content, repoURL) {
		return nil // Already configured
	}
	
	// Ensure we end with a newline before appending
	if content != "" && !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	
	// Append the index configuration
	indexConfig := fmt.Sprintf("\n[[index]]\nname = %q\nurl = %q\n", indexName, repoURL)
	content += indexConfig
	
	return os.WriteFile(configPath, []byte(content), 0644)
}
