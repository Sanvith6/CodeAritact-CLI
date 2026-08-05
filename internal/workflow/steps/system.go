package steps

import (
	"github.com/acme-corp/rs-cli/internal/network"
	"github.com/acme-corp/rs-cli/internal/validator"
	"github.com/acme-corp/rs-cli/internal/workflow"
)

// ValidateMachineStep ensures required tools are installed.
type ValidateMachineStep struct {
	validator validator.Validator
}

func NewValidateMachineStep(v validator.Validator) *ValidateMachineStep {
	return &ValidateMachineStep{validator: v}
}

func (s *ValidateMachineStep) Name() string { return "Validate machine requirements" }
func (s *ValidateMachineStep) Status() workflow.Status { return workflow.StatusPending }
func (s *ValidateMachineStep) Run(ctx *workflow.StepContext) error {
	results, err := s.validator.CheckAll(ctx.Ctx)
	if err != nil {
		return err
	}
	
	for _, res := range results {
		if res.Error != nil {
			return res.Error
		}
	}
	return nil
}
func (s *ValidateMachineStep) Rollback(ctx *workflow.StepContext) error { return nil }
func (s *ValidateMachineStep) Skip(ctx *workflow.StepContext) (bool, string) { return false, "" }

// CheckNetworkStep ensures internet and AWS connectivity.
type CheckNetworkStep struct {
	network network.Checker
	region  string
}

func NewCheckNetworkStep(n network.Checker, region string) *CheckNetworkStep {
	return &CheckNetworkStep{network: n, region: region}
}

func (s *CheckNetworkStep) Name() string { return "Check network connectivity" }
func (s *CheckNetworkStep) Status() workflow.Status { return workflow.StatusPending }
func (s *CheckNetworkStep) Run(ctx *workflow.StepContext) error {
	if err := s.network.CheckInternet(ctx.Ctx); err != nil {
		return err
	}
	if err := s.network.CheckAWSEndpoint(ctx.Ctx, s.region); err != nil {
		return err
	}
	return nil
}
func (s *CheckNetworkStep) Rollback(ctx *workflow.StepContext) error { return nil }
func (s *CheckNetworkStep) Skip(ctx *workflow.StepContext) (bool, string) { return false, "" }
