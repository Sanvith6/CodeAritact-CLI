package platform

import (
	"context"
)

// PackageManagerProvider defines the interface for package managers.
// This decouples the CLI from tools like uv, pip, poetry, or pipenv.
type PackageManagerProvider interface {
	// Name returns the provider name (e.g., "uv").
	Name() string
	
	// ConfigureAuth configures the package manager to authenticate with the repository.
	// This usually involves setting environment variables or writing a config file.
	ConfigureAuth(ctx context.Context, indexName, repoURL, token string) error
	
	// ClearAuth removes the authentication configuration.
	ClearAuth(ctx context.Context, indexName string) error
	
	// Install installs a package into the current environment.
	Install(ctx context.Context, pkg, version string) error
	
	// Sync synchronizes the project dependencies based on its lockfile.
	Sync(ctx context.Context) error
	
	// CreateVenv creates a new virtual environment at the specified path.
	CreateVenv(ctx context.Context, path string) error
	
	// IsVenvValid checks if a valid virtual environment exists at the specified path.
	IsVenvValid(ctx context.Context, path string) bool
}
