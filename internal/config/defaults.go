package config

// Defaults returns a Config populated with default values.
func Defaults() *Config {
	return &Config{
		Version: 1,
		AWS: AWSConfig{
			SSO: SSOConfig{
				StartURL:  "https://acme.awsapps.com/start",
				Region:    "us-east-1",
				AccountID: "123456789012",
				RoleName:  "DeveloperAccess",
			},
			ProfileName: "rs-dev",
			Region:      "us-east-1",
		},
		Artifact: ArtifactConfig{
			Domain:      "acme",
			DomainOwner: "123456789012",
			Repository:  "python-internal",
			Region:      "us-east-1",
		},
		Package: PackageConfig{
			UVIndexName:       "internal",
			ValidationPackage: "acme-core",
			InstallStrategy:   "add",
		},
		CLI: CLIConfig{
			LogLevel:       "info",
			NoColor:        false,
			NonInteractive: false,
		},
		Features: FeaturesConfig{
			CodeArtifact: true,
			Docker:       false,
			Java:         false,
			Terraform:    false,
			Kubernetes:   false,
		},
	}
}
