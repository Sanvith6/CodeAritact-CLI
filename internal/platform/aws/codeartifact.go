package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codeartifact"
	"github.com/aws/aws-sdk-go-v2/service/codeartifact/types"
	
	rsconfig "github.com/acme-corp/rs-cli/internal/config"
	"github.com/acme-corp/rs-cli/internal/errors"
	"github.com/acme-corp/rs-cli/internal/platform"
)

// CodeArtifactProvider implements platform.RepositoryProvider for AWS CodeArtifact.
type CodeArtifactProvider struct {
	cfg     *rsconfig.Config
	factory *ClientFactory
}

// NewCodeArtifactProvider creates a new CodeArtifactProvider.
func NewCodeArtifactProvider(cfg *rsconfig.Config) *CodeArtifactProvider {
	return &CodeArtifactProvider{
		cfg:     cfg,
		factory: NewClientFactory(cfg),
	}
}

// Name returns the provider name.
func (p *CodeArtifactProvider) Name() string {
	return "aws-codeartifact"
}

// GetAuthToken retrieves a temporary authorization token.
func (p *CodeArtifactProvider) GetAuthToken(ctx context.Context) (string, time.Time, error) {
	client, err := p.factory.CodeArtifactClient(ctx)
	if err != nil {
		return "", time.Time{}, err
	}

	// 12 hours is the max duration for CodeArtifact tokens
	input := &codeartifact.GetAuthorizationTokenInput{
		Domain:          aws.String(p.cfg.Artifact.Domain),
		DomainOwner:     aws.String(p.cfg.Artifact.DomainOwner),
		DurationSeconds: aws.Int64(43200),
	}

	res, err := client.GetAuthorizationToken(ctx, input)
	if err != nil {
		return "", time.Time{}, errors.Wrap(errors.ErrCodeArtifactAuth, err, "Failed to get CodeArtifact authorization token")
	}

	return aws.ToString(res.AuthorizationToken), aws.ToTime(res.Expiration), nil
}

// GetEndpoint retrieves the repository endpoint.
func (p *CodeArtifactProvider) GetEndpoint(ctx context.Context) (string, error) {
	client, err := p.factory.CodeArtifactClient(ctx)
	if err != nil {
		return "", err
	}

	input := &codeartifact.GetRepositoryEndpointInput{
		Domain:      aws.String(p.cfg.Artifact.Domain),
		DomainOwner: aws.String(p.cfg.Artifact.DomainOwner),
		Repository:  aws.String(p.cfg.Artifact.Repository),
		Format:      types.PackageFormatPypi,
	}

	res, err := client.GetRepositoryEndpoint(ctx, input)
	if err != nil {
		return "", errors.Wrap(errors.ErrCodeArtifactRepo, err, "Failed to get CodeArtifact repository endpoint")
	}

	return aws.ToString(res.RepositoryEndpoint), nil
}

// ValidateAccess checks if the user has access to the repository.
func (p *CodeArtifactProvider) ValidateAccess(ctx context.Context) error {
	client, err := p.factory.CodeArtifactClient(ctx)
	if err != nil {
		return err
	}

	input := &codeartifact.DescribeRepositoryInput{
		Domain:      aws.String(p.cfg.Artifact.Domain),
		DomainOwner: aws.String(p.cfg.Artifact.DomainOwner),
		Repository:  aws.String(p.cfg.Artifact.Repository),
	}

	_, err = client.DescribeRepository(ctx, input)
	if err != nil {
		return errors.Wrap(errors.ErrCodeArtifactRepo, err, "Failed to access CodeArtifact repository. Verify your permissions and repository name.")
	}

	return nil
}

// SearchPackages searches for packages matching the query prefix.
func (p *CodeArtifactProvider) SearchPackages(ctx context.Context, query string) ([]platform.Package, error) {
	client, err := p.factory.CodeArtifactClient(ctx)
	if err != nil {
		return nil, err
	}

	var packages []platform.Package
	var nextToken *string

	for {
		input := &codeartifact.ListPackagesInput{
			Domain:      aws.String(p.cfg.Artifact.Domain),
			DomainOwner: aws.String(p.cfg.Artifact.DomainOwner),
			Repository:  aws.String(p.cfg.Artifact.Repository),
			MaxResults:  aws.Int32(100), // Adjust as needed
		}

		if query != "" {
			input.PackagePrefix = aws.String(query)
		}
		
		if nextToken != nil {
			input.NextToken = nextToken
		}

		res, err := client.ListPackages(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("failed to list packages: %w", err)
		}

		for _, pkg := range res.Packages {
			packages = append(packages, platform.Package{
				Name:      aws.ToString(pkg.Package),
				Namespace: aws.ToString(pkg.Namespace),
				Format:    string(pkg.Format),
				Origin:    string(pkg.OriginConfiguration.Restrictions.Upstream), // Simplified
			})
		}

		nextToken = res.NextToken
		if nextToken == nil {
			break
		}
	}

	return packages, nil
}

// ListVersions lists all versions for a package.
func (p *CodeArtifactProvider) ListVersions(ctx context.Context, pkg string) ([]platform.Version, error) {
	client, err := p.factory.CodeArtifactClient(ctx)
	if err != nil {
		return nil, err
	}

	input := &codeartifact.ListPackageVersionsInput{
		Domain:      aws.String(p.cfg.Artifact.Domain),
		DomainOwner: aws.String(p.cfg.Artifact.DomainOwner),
		Repository:  aws.String(p.cfg.Artifact.Repository),
		Format:      types.PackageFormatPypi,
		Package:     aws.String(pkg),
		SortBy:      types.PackageVersionSortTypePublishedTime,
	}

	res, err := client.ListPackageVersions(ctx, input)
	if err != nil {
		return nil, errors.Wrap(errors.ErrCodeArtifactPkg, err, fmt.Sprintf("Failed to list versions for package %s", pkg))
	}

	var versions []platform.Version
	for _, v := range res.Versions {
		// Only include Published statuses for simplicity, or we could include everything
		if v.Status == types.PackageVersionStatusPublished {
			versions = append(versions, platform.Version{
				Version: aws.ToString(v.Version),
				Status:  string(v.Status),
				Updated: time.Time{}, // Time not available in basic list output
			})
		}
	}

	return versions, nil
}

// DescribePackage gets detailed info for a package.
func (p *CodeArtifactProvider) DescribePackage(ctx context.Context, pkg string) (*platform.PackageDetail, error) {
	// First check if we can list its versions
	versions, err := p.ListVersions(ctx, pkg)
	if err != nil {
		return nil, err
	}
	
	if len(versions) == 0 {
		return nil, fmt.Errorf("package %s not found or has no published versions", pkg)
	}
	
	// Since PyPI packages in CodeArtifact don't have rich description metadata via DescribePackage
	// we just assemble what we have.
	latest := versions[len(versions)-1].Version // Assuming SortBy PublishedTime puts newest last or first. CodeArtifact default is ascending.

	return &platform.PackageDetail{
		Package: platform.Package{
			Name:   pkg,
			Format: string(types.PackageFormatPypi),
		},
		LatestVersion: latest,
		Versions:      versions,
	}, nil
}
