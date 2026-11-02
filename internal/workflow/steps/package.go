package steps

import (
	"fmt"
	"os"

	"github.com/acme-corp/rs-cli/internal/platform"
	"github.com/acme-corp/rs-cli/internal/runner"
	"github.com/acme-corp/rs-cli/internal/services"
	"github.com/acme-corp/rs-cli/internal/workflow"
)

// ConfigureUVStep configures uv using the CodeArtifact token.
type ConfigureUVStep struct {
	repo *services.RepositoryService
	pm   platform.PackageManagerProvider
}

func NewConfigureUVStep(r *services.RepositoryService, p platform.PackageManagerProvider) *ConfigureUVStep {
	return &ConfigureUVStep{repo: r, pm: p}
}

func (s *ConfigureUVStep) Name() string { return "Configure uv package manager" }
func (s *ConfigureUVStep) Status() workflow.Status { return workflow.StatusPending }
func (s *ConfigureUVStep) Run(ctx *workflow.StepContext) error {
	endpoint, token, err := s.repo.GetAuthInfo(ctx.Ctx)
	if err != nil {
		return err
	}
	indexName := "internal"
	return s.pm.ConfigureAuth(ctx.Ctx, indexName, endpoint, token)
}
func (s *ConfigureUVStep) Rollback(ctx *workflow.StepContext) error { return nil }
func (s *ConfigureUVStep) Skip(ctx *workflow.StepContext) (bool, string) { return false, "" }

// VerifyValidationPackageStep ensures we can find the validation package.
type VerifyValidationPackageStep struct {
	repo *services.RepositoryService
}

func NewVerifyValidationPackageStep(r *services.RepositoryService) *VerifyValidationPackageStep {
	return &VerifyValidationPackageStep{repo: r}
}

func (s *VerifyValidationPackageStep) Name() string { return "Verify validation package" }
func (s *VerifyValidationPackageStep) Status() workflow.Status { return workflow.StatusPending }
func (s *VerifyValidationPackageStep) Run(ctx *workflow.StepContext) error {
	// The repo service can just fetch the package metadata from the REST API or using AWS SDK.
	// Since we are not wrapping uv for search anymore, we can just use the repository service's search method.
	pkgName := "company-web" // This should come from config, we'll fix it in bootstrap.go
	if v, ok := ctx.Ctx.Value("validation_package").(string); ok && v != "" {
		pkgName = v
	}

	err := s.repo.VerifyPackage(ctx.Ctx, pkgName)
	if err != nil {
		return fmt.Errorf("failed to verify validation package %s: %w", pkgName, err)
	}
	return nil
}
func (s *VerifyValidationPackageStep) Rollback(ctx *workflow.StepContext) error { return nil }
func (s *VerifyValidationPackageStep) Skip(ctx *workflow.StepContext) (bool, string) { return false, "" }

// SyncDependenciesStep runs `uv sync` to install/update dependencies.
type SyncDependenciesStep struct {
	runner runner.Runner
}

func NewSyncDependenciesStep(r runner.Runner) *SyncDependenciesStep {
	return &SyncDependenciesStep{runner: r}
}

func (s *SyncDependenciesStep) Name() string { return "Sync project dependencies" }
func (s *SyncDependenciesStep) Status() workflow.Status { return workflow.StatusPending }
func (s *SyncDependenciesStep) Run(ctx *workflow.StepContext) error {
	hasPyProject := false
	if _, err := os.Stat("pyproject.toml"); err == nil {
		hasPyProject = true
	}

	hasRequirements := false
	if _, err := os.Stat("requirements.txt"); err == nil {
		hasRequirements = true
	}

	if !hasPyProject && !hasRequirements {
		return fmt.Errorf("this file not found")
	}

	// If pyproject.toml exists, use uv sync (the default behavior)
	if hasPyProject {
		_, err := s.runner.Run(ctx.Ctx, "uv", "sync")
		if err != nil {
			return fmt.Errorf("not found for such package names")
		}
		return nil
	}

	// If requirements.txt exists, use uv pip install
	if hasRequirements {
		// Ensure a .venv exists first since pip install needs one
		if _, err := os.Stat(".venv"); os.IsNotExist(err) {
			_, err := s.runner.Run(ctx.Ctx, "uv", "venv")
			if err != nil {
				return fmt.Errorf("failed to create virtual environment: %w", err)
			}
		}
		
		_, err := s.runner.Run(ctx.Ctx, "uv", "pip", "install", "-r", "requirements.txt")
		if err != nil {
			return fmt.Errorf("not found for such package names")
		}
		return nil
	}

	return nil
}
func (s *SyncDependenciesStep) Rollback(ctx *workflow.StepContext) error { return nil }
func (s *SyncDependenciesStep) Skip(ctx *workflow.StepContext) (bool, string) {
	return false, ""
}
