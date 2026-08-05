package steps

import (
	"github.com/acme-corp/rs-cli/internal/services"
	"github.com/acme-corp/rs-cli/internal/workflow"
)

// ValidateArtifactAccessStep ensures we can access CodeArtifact.
type ValidateArtifactAccessStep struct {
	repo *services.RepositoryService
}

func NewValidateArtifactAccessStep(r *services.RepositoryService) *ValidateArtifactAccessStep {
	return &ValidateArtifactAccessStep{repo: r}
}

func (s *ValidateArtifactAccessStep) Name() string { return "Validate package repository access" }
func (s *ValidateArtifactAccessStep) Status() workflow.Status { return workflow.StatusPending }
func (s *ValidateArtifactAccessStep) Run(ctx *workflow.StepContext) error {
	return s.repo.ValidateAccess(ctx.Ctx)
}
func (s *ValidateArtifactAccessStep) Rollback(ctx *workflow.StepContext) error { return nil }
func (s *ValidateArtifactAccessStep) Skip(ctx *workflow.StepContext) (bool, string) { return false, "" }

// VerifyProjectStep ensures the current directory is a valid project.
type VerifyProjectStep struct {
	project *services.ProjectService
}

func NewVerifyProjectStep(p *services.ProjectService) *VerifyProjectStep {
	return &VerifyProjectStep{project: p}
}

func (s *VerifyProjectStep) Name() string { return "Verify project configuration" }
func (s *VerifyProjectStep) Status() workflow.Status { return workflow.StatusPending }
func (s *VerifyProjectStep) Run(ctx *workflow.StepContext) error {
	return s.project.Verify(ctx.Ctx)
}
func (s *VerifyProjectStep) Rollback(ctx *workflow.StepContext) error { return nil }
func (s *VerifyProjectStep) Skip(ctx *workflow.StepContext) (bool, string) { return false, "" }
