package uv

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	rsconfig "github.com/acme-corp/rs-cli/internal/config"
	"github.com/acme-corp/rs-cli/internal/platform"
)

// UVPackageManager implements platform.PackageManagerProvider for uv.
type UVPackageManager struct {
	cfg   *rsconfig.Config
	pInfo *platform.Info
}

// NewUVPackageManager creates a new UVPackageManager.
func NewUVPackageManager(cfg *rsconfig.Config, pInfo *platform.Info) *UVPackageManager {
	return &UVPackageManager{
		cfg:   cfg,
		pInfo: pInfo,
	}
}

// Name returns the provider name.
func (u *UVPackageManager) Name() string {
	return "uv"
}

// ConfigureAuth sets up the required environment variables for uv to authenticate,
// and optionally writes to uv.toml for global index configuration.
func (u *UVPackageManager) ConfigureAuth(ctx context.Context, indexName, repoURL, token string) error {
	safeIndexName := strings.ToUpper(strings.ReplaceAll(indexName, "-", "_"))
	
	userEnvVar := fmt.Sprintf("UV_INDEX_%s_USERNAME", safeIndexName)
	passEnvVar := fmt.Sprintf("UV_INDEX_%s_PASSWORD", safeIndexName)
	
	os.Setenv(userEnvVar, "aws")
	os.Setenv(passEnvVar, token)
	
	return u.ensureGlobalIndexConfig(indexName, repoURL, token)
}

// ClearAuth removes the authentication from the current process environment.
func (u *UVPackageManager) ClearAuth(ctx context.Context, indexName string) error {
	safeIndexName := strings.ToUpper(strings.ReplaceAll(indexName, "-", "_"))
	os.Unsetenv(fmt.Sprintf("UV_INDEX_%s_USERNAME", safeIndexName))
	os.Unsetenv(fmt.Sprintf("UV_INDEX_%s_PASSWORD", safeIndexName))
	return nil
}

func (u *UVPackageManager) ensureGlobalIndexConfig(indexName, repoURL, token string) error {
	configPath := filepath.Join(u.pInfo.UVConfigDir, "uv.toml")
	
	if err := os.MkdirAll(u.pInfo.UVConfigDir, 0700); err != nil {
		return err
	}
	
	// Ensure the URL ends with "simple/" for PyPI index compatibility
	if !strings.HasSuffix(repoURL, "simple/") {
		if !strings.HasSuffix(repoURL, "/") {
			repoURL += "/"
		}
		repoURL += "simple/"
	}
	
	// Inject credentials into the URL so native uv commands work directly in the shell
	authURL := strings.Replace(repoURL, "https://", fmt.Sprintf("https://aws:%s@", token), 1)

	var content string
	data, err := os.ReadFile(configPath)
	if err == nil {
		content = string(data)
	}
	
	// If it already exists, replace it, otherwise append
	indexConfig := fmt.Sprintf("\n[[index]]\nname = %q\nurl = %q\n", indexName, authURL)

	if strings.Contains(content, fmt.Sprintf("name = %q", indexName)) {
		// Extremely simple replacement of the entire file for now to keep it clean,
		// in a real app we'd parse TOML or use regex.
		// For now, let's just write fresh to ensure token updates.
		content = indexConfig
	} else {
		if content != "" && !strings.HasSuffix(content, "\n") {
			content += "\n"
		}
		content += indexConfig
	}
	
	return os.WriteFile(configPath, []byte(content), 0644)
}
