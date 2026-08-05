package runner

import (
	"context"
)

// Output contains the results of a command execution.
type Output struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

// Runner provides an interface for executing external commands safely.
type Runner interface {
	// Run executes a command and returns its output.
	Run(ctx context.Context, name string, args ...string) (Output, error)
	
	// RunWithEnv executes a command with additional environment variables.
	RunWithEnv(ctx context.Context, env []string, name string, args ...string) (Output, error)
	
	// RunSilent executes a command without logging its output (except on error).
	RunSilent(ctx context.Context, name string, args ...string) (Output, error)
	
	// LookPath searches for an executable named file in the directories named by the PATH environment variable.
	LookPath(name string) (string, error)
}
