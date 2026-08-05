package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Loader handles loading and parsing configuration from files and environment.
type Loader struct {
	v *viper.Viper
}

// NewLoader creates a new config loader.
func NewLoader() *Loader {
	v := viper.New()
	
	// Environment variable settings
	v.SetEnvPrefix("RS")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Default values
	defaults := Defaults()
	v.SetDefault("version", defaults.Version)
	
	v.SetDefault("aws.sso.start_url", defaults.AWS.SSO.StartURL)
	v.SetDefault("aws.sso.region", defaults.AWS.SSO.Region)
	v.SetDefault("aws.sso.account_id", defaults.AWS.SSO.AccountID)
	v.SetDefault("aws.sso.role_name", defaults.AWS.SSO.RoleName)
	v.SetDefault("aws.profile_name", defaults.AWS.ProfileName)
	v.SetDefault("aws.region", defaults.AWS.Region)

	v.SetDefault("artifact.domain", defaults.Artifact.Domain)
	v.SetDefault("artifact.domain_owner", defaults.Artifact.DomainOwner)
	v.SetDefault("artifact.repository", defaults.Artifact.Repository)
	v.SetDefault("artifact.region", defaults.Artifact.Region)

	v.SetDefault("package.uv_index_name", defaults.Package.UVIndexName)
	v.SetDefault("package.validation_package", defaults.Package.ValidationPackage)
	v.SetDefault("package.install_strategy", defaults.Package.InstallStrategy)

	v.SetDefault("cli.log_level", defaults.CLI.LogLevel)
	v.SetDefault("cli.no_color", defaults.CLI.NoColor)
	v.SetDefault("cli.non_interactive", defaults.CLI.NonInteractive)

	v.SetDefault("features.codeartifact", defaults.Features.CodeArtifact)
	v.SetDefault("features.docker", defaults.Features.Docker)
	v.SetDefault("features.java", defaults.Features.Java)
	v.SetDefault("features.terraform", defaults.Features.Terraform)
	v.SetDefault("features.kubernetes", defaults.Features.Kubernetes)

	return &Loader{
		v: v,
	}
}

// Load loads configuration from the specified path or falls back to defaults/env.
func (l *Loader) Load(configPath string) (*Config, error) {
	if configPath != "" {
		l.v.SetConfigFile(configPath)
	} else {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			l.v.AddConfigPath(filepath.Join(homeDir, ".rs"))
			l.v.SetConfigName("config")
			l.v.SetConfigType("yaml")
		}
	}

	// Ignore ConfigFileNotFoundError, as we have defaults.
	if err := l.v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok && !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var cfg Config
	if err := l.v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Run migrations if needed
	cfg, err := Migrate(cfg)
	if err != nil {
		return nil, fmt.Errorf("config migration failed: %w", err)
	}

	if err := l.Validate(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Validate checks the configuration for required fields.
func (l *Loader) Validate(cfg *Config) error {
	var missing []string

	if cfg.AWS.SSO.StartURL == "" {
		missing = append(missing, "aws.sso.start_url")
	}
	if cfg.AWS.SSO.Region == "" {
		missing = append(missing, "aws.sso.region")
	}
	if cfg.AWS.ProfileName == "" {
		missing = append(missing, "aws.profile_name")
	}

	if cfg.Features.CodeArtifact {
		if cfg.Artifact.Domain == "" {
			missing = append(missing, "artifact.domain")
		}
		if cfg.Artifact.Repository == "" {
			missing = append(missing, "artifact.repository")
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required configuration fields: %s", strings.Join(missing, ", "))
	}

	return nil
}

// Viper returns the underlying Viper instance, useful for Cobra binding.
func (l *Loader) Viper() *viper.Viper {
	return l.v
}
