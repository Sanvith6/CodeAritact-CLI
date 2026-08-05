package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoader_Defaults(t *testing.T) {
	loader := NewLoader()
	cfg, err := loader.Load("")
	require.NoError(t, err)

	assert.Equal(t, 1, cfg.Version)
	assert.Equal(t, "https://acme.awsapps.com/start", cfg.AWS.SSO.StartURL)
	assert.Equal(t, "info", cfg.CLI.LogLevel)
	assert.True(t, cfg.Features.CodeArtifact)
}

func TestLoader_EnvOverride(t *testing.T) {
	os.Setenv("RS_CLI_LOG_LEVEL", "debug")
	os.Setenv("RS_AWS_SSO_REGION", "eu-west-1")
	defer func() {
		os.Unsetenv("RS_CLI_LOG_LEVEL")
		os.Unsetenv("RS_AWS_SSO_REGION")
	}()

	loader := NewLoader()
	cfg, err := loader.Load("")
	require.NoError(t, err)

	assert.Equal(t, "debug", cfg.CLI.LogLevel)
	assert.Equal(t, "eu-west-1", cfg.AWS.SSO.Region)
}

func TestLoader_FileOverride(t *testing.T) {
	yamlContent := `
version: 1
cli:
  log_level: error
aws:
  profile_name: custom-profile
`
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")
	err := os.WriteFile(configFile, []byte(yamlContent), 0644)
	require.NoError(t, err)

	loader := NewLoader()
	cfg, err := loader.Load(configFile)
	require.NoError(t, err)

	assert.Equal(t, "error", cfg.CLI.LogLevel)
	assert.Equal(t, "custom-profile", cfg.AWS.ProfileName)
}

func TestLoader_Validation(t *testing.T) {
	loader := NewLoader()
	cfg := Defaults()
	cfg.AWS.SSO.StartURL = "" // Missing required field

	err := loader.Validate(cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "aws.sso.start_url")
}

func TestMigration(t *testing.T) {
	cfg := Config{Version: 0}
	migrated, err := Migrate(cfg)
	require.NoError(t, err)
	assert.Equal(t, 1, migrated.Version)
}
