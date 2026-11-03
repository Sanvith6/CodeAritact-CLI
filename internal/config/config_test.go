package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoader_Defaults(t *testing.T) {
	loader := NewLoader()
	cfg, err := loader.Load()
	require.NoError(t, err)

	assert.Equal(t, 1, cfg.Version)
	assert.Equal(t, "https://my-company.awsapps.com/start", cfg.AWS.SSO.StartURL)
	assert.Equal(t, "info", cfg.CLI.LogLevel)
	assert.True(t, cfg.Features.CodeArtifact)
}

func TestLoader_Validation(t *testing.T) {
	loader := NewLoader()
	cfg := &Config{
		AWS: AWSConfig{
			SSO: SSOConfig{}, // Missing StartURL
		},
	}

	err := loader.Validate(cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "AWS SSO StartURL")
}

func TestMigration(t *testing.T) {
	cfg := Config{Version: 0}
	migrated, err := Migrate(cfg)
	require.NoError(t, err)
	assert.Equal(t, 1, migrated.Version)
}
