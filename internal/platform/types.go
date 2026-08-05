package platform

import "time"

// Identity represents the authenticated developer's identity.
type Identity struct {
	AccountID     string
	UserID        string
	ARN           string
	RoleName      string
	PermissionSet string
	Provider      string // "aws", "azure", etc.
}

// Package represents a package in a repository.
type Package struct {
	Name      string
	Namespace string
	Format    string
	Origin    string
}

// Version represents a specific version of a package.
type Version struct {
	Version string
	Status  string
	Updated time.Time
}

// PackageDetail contains extended metadata for a package.
type PackageDetail struct {
	Package
	LatestVersion string
	Description   string
	Versions      []Version
}
