package platform

import (
	"context"
	"time"
)

// RepositoryProvider defines the interface for interacting with a package repository.
// This decouples the CLI from specific implementations like AWS CodeArtifact, JFrog, or GitHub Packages.
type RepositoryProvider interface {
	// Name returns the provider name (e.g., "aws-codeartifact").
	Name() string
	
	// GetAuthToken retrieves a temporary authentication token for the repository.
	GetAuthToken(ctx context.Context) (token string, expiry time.Time, err error)
	
	// GetEndpoint returns the base URL for the repository.
	GetEndpoint(ctx context.Context) (string, error)
	
	// ValidateAccess checks if the user has read access to the repository.
	ValidateAccess(ctx context.Context) error
	
	// SearchPackages searches for packages matching a query.
	SearchPackages(ctx context.Context, query string) ([]Package, error)
	
	// ListVersions lists all versions for a specific package.
	ListVersions(ctx context.Context, pkg string) ([]Version, error)
	
	// DescribePackage retrieves detailed metadata for a package.
	DescribePackage(ctx context.Context, pkg string) (*PackageDetail, error)
}
