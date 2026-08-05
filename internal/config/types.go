package config

// Config represents the root configuration structure for the RS Platform CLI.
type Config struct {
	Version int `yaml:"version" mapstructure:"version"` // Schema version (currently 1)

	AWS      AWSConfig      `yaml:"aws"      mapstructure:"aws"`
	Artifact ArtifactConfig `yaml:"artifact" mapstructure:"artifact"`
	Package  PackageConfig  `yaml:"package"  mapstructure:"package"`
	CLI      CLIConfig      `yaml:"cli"      mapstructure:"cli"`
	Features FeaturesConfig `yaml:"features" mapstructure:"features"`
}

// AWSConfig holds AWS-specific configuration.
type AWSConfig struct {
	SSO         SSOConfig `yaml:"sso"          mapstructure:"sso"`
	ProfileName string    `yaml:"profile_name" mapstructure:"profile_name"`
	Region      string    `yaml:"region"       mapstructure:"region"`
}

// SSOConfig holds AWS SSO configuration parameters.
type SSOConfig struct {
	StartURL  string `yaml:"start_url"  mapstructure:"start_url"`
	Region    string `yaml:"region"     mapstructure:"region"`
	AccountID string `yaml:"account_id" mapstructure:"account_id"`
	RoleName  string `yaml:"role_name"  mapstructure:"role_name"`
}

// ArtifactConfig holds AWS CodeArtifact configuration.
type ArtifactConfig struct {
	Domain      string `yaml:"domain"       mapstructure:"domain"`
	DomainOwner string `yaml:"domain_owner" mapstructure:"domain_owner"`
	Repository  string `yaml:"repository"   mapstructure:"repository"`
	Region      string `yaml:"region"       mapstructure:"region"`
}

// PackageConfig holds package manager (e.g., uv) configuration.
type PackageConfig struct {
	UVIndexName       string `yaml:"uv_index_name"      mapstructure:"uv_index_name"`
	ValidationPackage string `yaml:"validation_package" mapstructure:"validation_package"`
	InstallStrategy   string `yaml:"install_strategy"   mapstructure:"install_strategy"` // "add" or "pip-install"
}

// CLIConfig holds general CLI operational configuration.
type CLIConfig struct {
	LogLevel       string `yaml:"log_level"       mapstructure:"log_level"`
	NoColor        bool   `yaml:"no_color"        mapstructure:"no_color"`
	NonInteractive bool   `yaml:"non_interactive" mapstructure:"non_interactive"`
}

// FeaturesConfig contains feature flags for enabling/disabling capabilities.
type FeaturesConfig struct {
	CodeArtifact bool `yaml:"codeartifact" mapstructure:"codeartifact"`
	Docker       bool `yaml:"docker"       mapstructure:"docker"`
	Java         bool `yaml:"java"         mapstructure:"java"`
	Terraform    bool `yaml:"terraform"    mapstructure:"terraform"`
	Kubernetes   bool `yaml:"kubernetes"   mapstructure:"kubernetes"`
}
