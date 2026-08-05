package runner

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
)

// ExecRunner implements Runner using os/exec.
type ExecRunner struct {
	logger *slog.Logger
}

// NewExecRunner creates a new ExecRunner.
func NewExecRunner(logger *slog.Logger) *ExecRunner {
	if logger == nil {
		logger = slog.Default()
	}
	return &ExecRunner{
		logger: logger,
	}
}

// Run executes a command and returns its output.
func (r *ExecRunner) Run(ctx context.Context, name string, args ...string) (Output, error) {
	return r.run(ctx, false, nil, name, args...)
}

// RunWithEnv executes a command with additional environment variables.
func (r *ExecRunner) RunWithEnv(ctx context.Context, env []string, name string, args ...string) (Output, error) {
	return r.run(ctx, false, env, name, args...)
}

// RunSilent executes a command without logging its execution details.
func (r *ExecRunner) RunSilent(ctx context.Context, name string, args ...string) (Output, error) {
	return r.run(ctx, true, nil, name, args...)
}

// LookPath searches for an executable in the PATH.
func (r *ExecRunner) LookPath(name string) (string, error) {
	return exec.LookPath(name)
}

func (r *ExecRunner) run(ctx context.Context, silent bool, env []string, name string, args ...string) (Output, error) {
	if !silent {
		r.logger.Debug("executing command", 
			slog.String("command", name), 
			slog.Int("args_count", len(args)),
		)
	}

	cmd := exec.CommandContext(ctx, name, args...)
	
	if len(env) > 0 {
		cmd.Env = append(os.Environ(), env...)
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	
	out := Output{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
		ExitCode: cmd.ProcessState.ExitCode(),
	}

	if err != nil {
		if !silent {
			r.logger.Error("command failed", 
				slog.String("command", name),
				slog.String("error", err.Error()),
				slog.String("stderr", out.Stderr),
				slog.Int("exit_code", out.ExitCode),
			)
		}
		return out, fmt.Errorf("command execution failed: %w", err)
	}

	if !silent {
		r.logger.Debug("command succeeded", slog.String("command", name))
	}

	return out, nil
}
