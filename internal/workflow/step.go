package workflow

import (
	"context"
)

// Status represents the state of a workflow step.
type Status string

const (
	StatusPending   Status = "pending"
	StatusRunning   Status = "running"
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
	StatusSkipped   Status = "skipped"
	StatusRolledBack Status = "rolled_back"
)

// StepContext provides contextual state to a workflow step.
// It allows steps to share data during a workflow execution.
type StepContext struct {
	Ctx    context.Context
	Values map[string]interface{}
}

// NewStepContext creates a new StepContext.
func NewStepContext(ctx context.Context) *StepContext {
	return &StepContext{
		Ctx:    ctx,
		Values: make(map[string]interface{}),
	}
}

// Step defines a discrete unit of work in a workflow.
type Step interface {
	// Name returns the name of the step.
	Name() string

	// Status returns the current status of the step.
	Status() Status

	// Run executes the step logic.
	Run(ctx *StepContext) error

	// Rollback attempts to revert the changes made by Run.
	Rollback(ctx *StepContext) error

	// Skip determines whether this step can be safely skipped based on current state.
	Skip(ctx *StepContext) (bool, string)
}
