package runner

import (
	"context"
)

// AWSRunner wraps a base Runner and pre-configures it for AWS CLI commands.
type AWSRunner struct {
	base Runner
}

// NewAWSRunner creates a new AWSRunner.
func NewAWSRunner(base Runner) *AWSRunner {
	return &AWSRunner{base: base}
}

// Run executes an AWS CLI command with --output json and --no-cli-pager.
func (r *AWSRunner) Run(ctx context.Context, name string, args ...string) (Output, error) {
	// Name is expected to be "aws", we just append the default args
	fullArgs := append(args, "--output", "json", "--no-cli-pager")
	return r.base.Run(ctx, name, fullArgs...)
}

func (r *AWSRunner) RunWithEnv(ctx context.Context, env []string, name string, args ...string) (Output, error) {
	fullArgs := append(args, "--output", "json", "--no-cli-pager")
	return r.base.RunWithEnv(ctx, env, name, fullArgs...)
}

func (r *AWSRunner) RunSilent(ctx context.Context, name string, args ...string) (Output, error) {
	fullArgs := append(args, "--output", "json", "--no-cli-pager")
	return r.base.RunSilent(ctx, name, fullArgs...)
}

func (r *AWSRunner) LookPath(name string) (string, error) {
	return r.base.LookPath(name)
}

// UVRunner wraps a base Runner and pre-configures it for uv commands.
type UVRunner struct {
	base Runner
}

// NewUVRunner creates a new UVRunner.
func NewUVRunner(base Runner) *UVRunner {
	return &UVRunner{base: base}
}

func (r *UVRunner) Run(ctx context.Context, name string, args ...string) (Output, error) {
	return r.base.Run(ctx, name, args...)
}

func (r *UVRunner) RunWithEnv(ctx context.Context, env []string, name string, args ...string) (Output, error) {
	return r.base.RunWithEnv(ctx, env, name, args...)
}

func (r *UVRunner) RunSilent(ctx context.Context, name string, args ...string) (Output, error) {
	return r.base.RunSilent(ctx, name, args...)
}

func (r *UVRunner) LookPath(name string) (string, error) {
	return r.base.LookPath(name)
}

// GitRunner wraps a base Runner and pre-configures it for git commands.
type GitRunner struct {
	base Runner
}

// NewGitRunner creates a new GitRunner.
func NewGitRunner(base Runner) *GitRunner {
	return &GitRunner{base: base}
}

func (r *GitRunner) Run(ctx context.Context, name string, args ...string) (Output, error) {
	return r.base.Run(ctx, name, args...)
}

func (r *GitRunner) RunWithEnv(ctx context.Context, env []string, name string, args ...string) (Output, error) {
	return r.base.RunWithEnv(ctx, env, name, args...)
}

func (r *GitRunner) RunSilent(ctx context.Context, name string, args ...string) (Output, error) {
	return r.base.RunSilent(ctx, name, args...)
}

func (r *GitRunner) LookPath(name string) (string, error) {
	return r.base.LookPath(name)
}
