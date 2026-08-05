package runner

import (
	"context"
)

// MockRunner is a mock implementation of Runner for testing.
type MockRunner struct {
	RunFunc        func(ctx context.Context, name string, args ...string) (Output, error)
	RunWithEnvFunc func(ctx context.Context, env []string, name string, args ...string) (Output, error)
	RunSilentFunc  func(ctx context.Context, name string, args ...string) (Output, error)
	LookPathFunc   func(name string) (string, error)
	
	// Calls tracks the commands that were executed
	Calls [][]string
}

func (m *MockRunner) recordCall(name string, args ...string) {
	call := append([]string{name}, args...)
	m.Calls = append(m.Calls, call)
}

func (m *MockRunner) Run(ctx context.Context, name string, args ...string) (Output, error) {
	m.recordCall(name, args...)
	if m.RunFunc != nil {
		return m.RunFunc(ctx, name, args...)
	}
	return Output{ExitCode: 0}, nil
}

func (m *MockRunner) RunWithEnv(ctx context.Context, env []string, name string, args ...string) (Output, error) {
	m.recordCall(name, args...)
	if m.RunWithEnvFunc != nil {
		return m.RunWithEnvFunc(ctx, env, name, args...)
	}
	return Output{ExitCode: 0}, nil
}

func (m *MockRunner) RunSilent(ctx context.Context, name string, args ...string) (Output, error) {
	m.recordCall(name, args...)
	if m.RunSilentFunc != nil {
		return m.RunSilentFunc(ctx, name, args...)
	}
	return Output{ExitCode: 0}, nil
}

func (m *MockRunner) LookPath(name string) (string, error) {
	if m.LookPathFunc != nil {
		return m.LookPathFunc(name)
	}
	return "/mock/path/" + name, nil
}
