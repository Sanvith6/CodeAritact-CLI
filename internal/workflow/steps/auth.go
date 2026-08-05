package steps

import (
	"github.com/acme-corp/rs-cli/internal/platform"
	"github.com/acme-corp/rs-cli/internal/services"
	"github.com/acme-corp/rs-cli/internal/workflow"
)

// CheckSSOProfileStep ensures the local AWS config has the right profile.
type CheckSSOProfileStep struct {
	auth platform.AuthProvider
}

func NewCheckSSOProfileStep(a platform.AuthProvider) *CheckSSOProfileStep {
	return &CheckSSOProfileStep{auth: a}
}

func (s *CheckSSOProfileStep) Name() string { return "Check AWS SSO profile" }
func (s *CheckSSOProfileStep) Status() workflow.Status { return workflow.StatusPending }
func (s *CheckSSOProfileStep) Run(ctx *workflow.StepContext) error {
	return s.auth.EnsureProfile(ctx.Ctx)
}
func (s *CheckSSOProfileStep) Rollback(ctx *workflow.StepContext) error { return nil }
func (s *CheckSSOProfileStep) Skip(ctx *workflow.StepContext) (bool, string) { return false, "" }

// EnsureSSOLoginStep ensures the user is authenticated.
type EnsureSSOLoginStep struct {
	auth *services.AuthService
}

func NewEnsureSSOLoginStep(a *services.AuthService) *EnsureSSOLoginStep {
	return &EnsureSSOLoginStep{auth: a}
}

func (s *EnsureSSOLoginStep) Name() string { return "Ensure AWS SSO login" }
func (s *EnsureSSOLoginStep) Status() workflow.Status { return workflow.StatusPending }
func (s *EnsureSSOLoginStep) Run(ctx *workflow.StepContext) error {
	return s.auth.EnsureAuthenticated(ctx.Ctx)
}
func (s *EnsureSSOLoginStep) Rollback(ctx *workflow.StepContext) error { return nil }
func (s *EnsureSSOLoginStep) Skip(ctx *workflow.StepContext) (bool, string) {
	if s.auth.IsAuthenticated(ctx.Ctx) {
		return true, "Session is already valid"
	}
	return false, ""
}
