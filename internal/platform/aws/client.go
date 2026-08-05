package aws

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/codeartifact"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	
	rsconfig "github.com/acme-corp/rs-cli/internal/config"
)

// ClientFactory manages AWS SDK v2 clients.
type ClientFactory struct {
	cfg        *rsconfig.Config
	awsCfg     *aws.Config
	mu         sync.Mutex
	sts        *sts.Client
	sso        *sso.Client
	ssooidc    *ssooidc.Client
	codeartifact *codeartifact.Client
}

// NewClientFactory creates a new factory for AWS clients.
func NewClientFactory(cfg *rsconfig.Config) *ClientFactory {
	return &ClientFactory{
		cfg: cfg,
	}
}

// loadConfig lazily loads the AWS SDK config based on our profile and region.
func (f *ClientFactory) loadConfig(ctx context.Context) (aws.Config, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.awsCfg != nil {
		return *f.awsCfg, nil
	}

	opts := []func(*config.LoadOptions) error{
		config.WithSharedConfigProfile(f.cfg.AWS.ProfileName),
		config.WithRegion(f.cfg.AWS.Region),
	}

	awsCfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return aws.Config{}, fmt.Errorf("failed to load AWS configuration: %w", err)
	}

	f.awsCfg = &awsCfg
	return awsCfg, nil
}

// STSClient returns a configured STS client.
func (f *ClientFactory) STSClient(ctx context.Context) (*sts.Client, error) {
	f.mu.Lock()
	if f.sts != nil {
		f.mu.Unlock()
		return f.sts, nil
	}
	f.mu.Unlock()

	cfg, err := f.loadConfig(ctx)
	if err != nil {
		return nil, err
	}

	client := sts.NewFromConfig(cfg)
	
	f.mu.Lock()
	f.sts = client
	f.mu.Unlock()
	
	return client, nil
}

// SSOClient returns a configured SSO client.
func (f *ClientFactory) SSOClient(ctx context.Context) (*sso.Client, error) {
	f.mu.Lock()
	if f.sso != nil {
		f.mu.Unlock()
		return f.sso, nil
	}
	f.mu.Unlock()

	cfg, err := f.loadConfig(ctx)
	if err != nil {
		return nil, err
	}

	client := sso.NewFromConfig(cfg)
	
	f.mu.Lock()
	f.sso = client
	f.mu.Unlock()
	
	return client, nil
}

// SSOOIDCClient returns a configured SSOOIDC client.
func (f *ClientFactory) SSOOIDCClient(ctx context.Context) (*ssooidc.Client, error) {
	f.mu.Lock()
	if f.ssooidc != nil {
		f.mu.Unlock()
		return f.ssooidc, nil
	}
	f.mu.Unlock()

	cfg, err := f.loadConfig(ctx)
	if err != nil {
		return nil, err
	}

	client := ssooidc.NewFromConfig(cfg)
	
	f.mu.Lock()
	f.ssooidc = client
	f.mu.Unlock()
	
	return client, nil
}

// CodeArtifactClient returns a configured CodeArtifact client.
func (f *ClientFactory) CodeArtifactClient(ctx context.Context) (*codeartifact.Client, error) {
	f.mu.Lock()
	if f.codeartifact != nil {
		f.mu.Unlock()
		return f.codeartifact, nil
	}
	f.mu.Unlock()

	cfg, err := f.loadConfig(ctx)
	if err != nil {
		return nil, err
	}

	// CodeArtifact clients typically need the specific region of the repository
	cfg.Region = f.cfg.Artifact.Region
	
	client := codeartifact.NewFromConfig(cfg)
	
	f.mu.Lock()
	f.codeartifact = client
	f.mu.Unlock()
	
	return client, nil
}
