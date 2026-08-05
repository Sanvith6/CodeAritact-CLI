package aws

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
	rsconfig "github.com/acme-corp/rs-cli/internal/config"
	"github.com/acme-corp/rs-cli/internal/platform"
)

// ProfileManager handles reading and writing the AWS config file (~/.aws/config).
type ProfileManager struct {
	cfg         *rsconfig.Config
	platformCfg *platform.Info
}

// NewProfileManager creates a new ProfileManager.
func NewProfileManager(cfg *rsconfig.Config, pInfo *platform.Info) *ProfileManager {
	return &ProfileManager{
		cfg:         cfg,
		platformCfg: pInfo,
	}
}

// configFilePath returns the path to the ~/.aws/config file.
func (m *ProfileManager) configFilePath() string {
	return filepath.Join(m.platformCfg.AWSConfigDir, "config")
}

// EnsureProfile creates the SSO profile if it doesn't exist.
func (m *ProfileManager) EnsureProfile() error {
	path := m.configFilePath()
	
	// Create directory if missing
	if err := os.MkdirAll(m.platformCfg.AWSConfigDir, 0700); err != nil {
		return fmt.Errorf("failed to create aws config directory: %w", err)
	}

	// Load or create ini file
	var cfg *ini.File
	var err error
	
	if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
		cfg = ini.Empty()
	} else {
		cfg, err = ini.Load(path)
		if err != nil {
			return fmt.Errorf("failed to parse aws config file: %w", err)
		}
	}

	sectionName := "profile " + m.cfg.AWS.ProfileName
	section, err := cfg.GetSection(sectionName)
	
	// If section doesn't exist, create it
	if err != nil {
		section, err = cfg.NewSection(sectionName)
		if err != nil {
			return fmt.Errorf("failed to create profile section: %w", err)
		}
	}

	// Set SSO values
	section.Key("sso_session").SetValue("rs-session")
	section.Key("sso_account_id").SetValue(m.cfg.AWS.SSO.AccountID)
	section.Key("sso_role_name").SetValue(m.cfg.AWS.SSO.RoleName)
	section.Key("region").SetValue(m.cfg.AWS.Region)
	section.Key("output").SetValue("json")

	// Create SSO session section
	ssoSectionName := "sso-session rs-session"
	ssoSection, err := cfg.GetSection(ssoSectionName)
	if err != nil {
		ssoSection, err = cfg.NewSection(ssoSectionName)
		if err != nil {
			return fmt.Errorf("failed to create sso session section: %w", err)
		}
	}

	ssoSection.Key("sso_start_url").SetValue(m.cfg.AWS.SSO.StartURL)
	ssoSection.Key("sso_region").SetValue(m.cfg.AWS.SSO.Region)
	ssoSection.Key("sso_registration_scopes").SetValue("sso:account:access")

	// Save file
	if err := cfg.SaveTo(path); err != nil {
		return fmt.Errorf("failed to save aws config file: %w", err)
	}

	return nil
}

// ValidateProfile checks if the existing profile matches the expected configuration.
func (m *ProfileManager) ValidateProfile() error {
	path := m.configFilePath()
	
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("aws config file does not exist")
	}

	cfg, err := ini.Load(path)
	if err != nil {
		return fmt.Errorf("failed to parse aws config file: %w", err)
	}

	sectionName := "profile " + m.cfg.AWS.ProfileName
	section, err := cfg.GetSection(sectionName)
	if err != nil {
		return fmt.Errorf("profile '%s' not found in config", m.cfg.AWS.ProfileName)
	}

	if section.Key("sso_account_id").String() != m.cfg.AWS.SSO.AccountID {
		return fmt.Errorf("account ID mismatch")
	}
	if section.Key("sso_role_name").String() != m.cfg.AWS.SSO.RoleName {
		return fmt.Errorf("role name mismatch")
	}

	return nil
}
