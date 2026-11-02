package platform

import (
	"context"
)

// PackageManagerProvider defines the interface for package managers.
type PackageManagerProvider interface {
	// Name returns the provider name (e.g., "uv").
	Name() string
	
	// ConfigureAuth configures the package manager to authenticate with the repository.
	ConfigureAuth(ctx context.Context, indexName, repoURL, token string) error
	
	// ClearAuth removes the authentication configuration.
	ClearAuth(ctx context.Context, indexName string) error
}
