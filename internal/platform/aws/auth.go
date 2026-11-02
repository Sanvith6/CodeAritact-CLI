package aws

import (
	"context"
	"fmt"
	"strings"
	"log/slog"

	rsconfig "github.com/acme-corp/rs-cli/internal/config"
	"github.com/acme-corp/rs-cli/internal/errors"
	"github.com/acme-corp/rs-cli/internal/platform"
	"github.com/acme-corp/rs-cli/internal/runner"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// SSOAuthProvider implements platform.AuthProvider for AWS IAM Identity Center.
type SSOAuthProvider struct {
	cfg     *rsconfig.Config
	pInfo   *platform.Info
	runner  runner.Runner
	logger  *slog.Logger
	factory *ClientFactory
	profile *ProfileManager
}

// NewSSOAuthProvider creates a new SSOAuthProvider.
func NewSSOAuthProvider(
	cfg *rsconfig.Config, 
	pInfo *platform.Info, 
	r runner.Runner, 
	logger *slog.Logger,
) *SSOAuthProvider {
	return &SSOAuthProvider{
		cfg:     cfg,
		pInfo:   pInfo,
		runner:  runner.NewAWSRunner(r),
		logger:  logger,
		factory: NewClientFactory(cfg),
		profile: NewProfileManager(cfg, pInfo),
	}
}

// Name returns the provider name.
func (a *SSOAuthProvider) Name() string {
	return "aws-sso"
}

// EnsureProfile ensures the local AWS config file is correct.
func (a *SSOAuthProvider) EnsureProfile(ctx context.Context) error {
	// Ensure the binary was compiled with valid configuration
	if a.cfg.AWS.SSO.StartURL == "" {
		return fmt.Errorf("RS CLI is not configured.\n\nRequired embedded configuration is missing.\n\nPlease contact the DevOps team.")
	}

	err := a.profile.EnsureProfile()
	if err != nil {
		return errors.Wrap(errors.ErrCodeSSOProfile, err, "Failed to configure AWS SSO profile")
	}
	return nil
}

// ValidateProfile checks if the profile is correct.
func (a *SSOAuthProvider) ValidateProfile(ctx context.Context) error {
	err := a.profile.ValidateProfile()
	if err != nil {
		return errors.Wrap(errors.ErrCodeSSOProfile, err, "Invalid AWS SSO profile")
	}
	return nil
}

// Login triggers the AWS CLI SSO login flow.
func (a *SSOAuthProvider) Login(ctx context.Context) error {
	// The AWS CLI inherently supports opening the browser for OIDC device code flow,
	// which is much better than implementing it ourselves.
	// We use --use-device-code to prevent the AWS CLI from starting a local HTTP server
	// for the callback, which triggers the Windows Firewall / UAC prompt.
	out, err := a.runner.Run(ctx, "aws", "sso", "login", "--profile", a.cfg.AWS.ProfileName, "--use-device-code")
	if err != nil {
		errMsg := "The organization's AWS SSO configuration appears to be invalid.\n\nPlease contact the DevOps team."
		
		if a.cfg.CLI.LogLevel == "debug" {
			return fmt.Errorf("%s\n\nDetailed error: %v\nOutput: %s", errMsg, err, strings.TrimSpace(out.Stderr))
		}
		
		return fmt.Errorf("%s", errMsg)
	}
	
	a.logger.Debug("SSO login output", slog.String("stdout", out.Stdout), slog.String("stderr", out.Stderr))
	return nil
}

// Logout triggers the AWS CLI SSO logout flow.
func (a *SSOAuthProvider) Logout(ctx context.Context) error {
	_, err := a.runner.Run(ctx, "aws", "sso", "logout", "--profile", a.cfg.AWS.ProfileName)
	if err != nil {
		return fmt.Errorf("failed to log out: %w", err)
	}
	return nil
}

// IsSessionValid checks if we can call sts:GetCallerIdentity.
func (a *SSOAuthProvider) IsSessionValid(ctx context.Context) bool {
	client, err := a.factory.STSClient(ctx)
	if err != nil {
		return false
	}

	_, err = client.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	
	if err != nil {
		a.logger.Debug("STS GetCallerIdentity failed, session invalid", slog.String("error", err.Error()))
		return false
	}
	
	return true
}

// GetIdentity retrieves the current AWS identity from STS.
func (a *SSOAuthProvider) GetIdentity(ctx context.Context) (*platform.Identity, error) {
	client, err := a.factory.STSClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get STS client: %w", err)
	}

	res, err := client.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, errors.Wrap(errors.ErrCodeSSOExpired, err, "Your AWS session has expired.\n\nRun:\n    rs bootstrap\n\nto authenticate again.")
	}

	// Parse ARN to get role name. 
	// Typical SSO ARN: arn:aws:sts::123456789012:assumed-role/AWSReservedSSO_RoleName_12345/user@domain.com
	roleName := "Unknown"
	permissionSet := "Unknown"
	
	if res.Arn != nil {
		parts := strings.Split(*res.Arn, "/")
		if len(parts) >= 2 {
			rolePart := parts[1]
			// e.g. AWSReservedSSO_DeveloperAccess_b2c3d4e5
			if strings.HasPrefix(rolePart, "AWSReservedSSO_") {
				roleParts := strings.Split(rolePart, "_")
				if len(roleParts) >= 2 {
					roleName = roleParts[1]
					permissionSet = roleName
				} else {
					roleName = rolePart
				}
			} else {
				roleName = rolePart
			}
		}
	}

	return &platform.Identity{
		AccountID:     aws.ToString(res.Account),
		UserID:        aws.ToString(res.UserId),
		ARN:           aws.ToString(res.Arn),
		RoleName:      roleName,
		PermissionSet: permissionSet,
		Provider:      "aws-sso",
	}, nil
}
