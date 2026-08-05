package workflow

import (
	"fmt"
	"log/slog"
	
	"github.com/acme-corp/rs-cli/internal/ui"
)

// Engine executes a series of workflow steps.
type Engine struct {
	steps     []Step
	uiManager *ui.Manager
	logger    *slog.Logger
}

// NewEngine creates a new workflow engine.
func NewEngine(steps []Step, uiManager *ui.Manager, logger *slog.Logger) *Engine {
	return &Engine{
		steps:     steps,
		uiManager: uiManager,
		logger:    logger,
	}
}

// Run executes all steps sequentially.
// If a step fails, it stops execution and attempts to rollback completed steps in reverse order.
func (e *Engine) Run(ctx *StepContext) error {
	var completedSteps []Step
	
	for i, step := range e.steps {
		stepName := step.Name()
		prefix := fmt.Sprintf("[%d/%d]", i+1, len(e.steps))
		
		// Check if it can be skipped
		skip, reason := step.Skip(ctx)
		if skip {
			if reason != "" {
				ui.PrintMuted(fmt.Sprintf("%s Skipping %s: %s", prefix, stepName, reason))
			} else {
				ui.PrintMuted(fmt.Sprintf("%s Skipping %s", prefix, stepName))
			}
			continue
		}

		e.logger.Info("Starting step", slog.String("step", stepName))
		
		// Start spinner
		spinnerMsg := fmt.Sprintf("%s %s", prefix, stepName)
		e.uiManager.StartSpinner(spinnerMsg)

		err := step.Run(ctx)
		if err != nil {
			e.uiManager.StopSpinner(false, fmt.Sprintf("%s %s failed", prefix, stepName))
			e.logger.Error("Step failed", slog.String("step", stepName), slog.String("error", err.Error()))
			
			// Rollback completed steps
			e.rollback(ctx, completedSteps)
			
			return fmt.Errorf("workflow failed at step '%s': %w", stepName, err)
		}

		e.uiManager.StopSpinner(true, fmt.Sprintf("%s %s", prefix, stepName))
		completedSteps = append(completedSteps, step)
	}

	return nil
}

// rollback attempts to revert the given steps in reverse order.
func (e *Engine) rollback(ctx *StepContext, completedSteps []Step) {
	if len(completedSteps) == 0 {
		return
	}
	
	ui.PrintWarning("Attempting to roll back changes...")
	
	for i := len(completedSteps) - 1; i >= 0; i-- {
		step := completedSteps[i]
		e.logger.Info("Rolling back step", slog.String("step", step.Name()))
		
		if err := step.Rollback(ctx); err != nil {
			// Best effort rollback. We log but continue trying to roll back earlier steps.
			e.logger.Error("Rollback failed for step", slog.String("step", step.Name()), slog.String("error", err.Error()))
			ui.PrintError(fmt.Sprintf("Failed to roll back step '%s': %v", step.Name(), err))
		}
	}
}
