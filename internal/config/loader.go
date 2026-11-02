package config

import (
	"fmt"
	"strings"
)

// Loader handles configuration retrieval.
type Loader struct {}

// NewLoader creates a new config loader.
func NewLoader() *Loader {
	return &Loader{}
}

// Load loads the compiled configuration directly from the generated struct.
func (l *Loader) Load() (*Config, error) {
	if DefaultConfig == nil {
		return nil, fmt.Errorf("configuration not compiled: DefaultConfig is nil")
	}

	if err := l.Validate(DefaultConfig); err != nil {
		return nil, err
	}
	
	return DefaultConfig, nil
}

// Validate checks the configuration for required fields.
func (l *Loader) Validate(cfg *Config) error {
	var missing []string

	if cfg.AWS.SSO.StartURL == "" {
		missing = append(missing, "AWS SSO StartURL")
	}
	if cfg.AWS.SSO.Region == "" {
		missing = append(missing, "AWS SSO Region")
	}
	if cfg.AWS.SSO.AccountID == "" {
		missing = append(missing, "AWS SSO AccountID")
	}
	if cfg.AWS.SSO.RoleName == "" {
		missing = append(missing, "AWS SSO RoleName")
	}
	if cfg.AWS.ProfileName == "" {
		missing = append(missing, "AWS ProfileName")
	}

	if cfg.Features.CodeArtifact {
		if cfg.Artifact.Domain == "" {
			missing = append(missing, "Artifact Domain")
		}
		if cfg.Artifact.DomainOwner == "" {
			missing = append(missing, "Artifact DomainOwner")
		}
		if cfg.Artifact.Repository == "" {
			missing = append(missing, "Artifact Repository")
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("invalid embedded configuration. Missing required fields: %s", strings.Join(missing, ", "))
	}

	return nil
}
