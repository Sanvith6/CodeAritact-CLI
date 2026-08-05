package runner

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockRunner(t *testing.T) {
	mock := &MockRunner{
		RunFunc: func(ctx context.Context, name string, args ...string) (Output, error) {
			return Output{Stdout: "mocked", ExitCode: 0}, nil
		},
	}

	out, err := mock.Run(context.Background(), "echo", "hello")
	require.NoError(t, err)
	assert.Equal(t, "mocked", out.Stdout)
	assert.Equal(t, 0, out.ExitCode)
	
	assert.Len(t, mock.Calls, 1)
	assert.Equal(t, []string{"echo", "hello"}, mock.Calls[0])
}

func TestTypedRunners(t *testing.T) {
	mock := &MockRunner{}
	
	aws := NewAWSRunner(mock)
	_, _ = aws.Run(context.Background(), "aws", "sso", "login")
	
	assert.Len(t, mock.Calls, 1)
	// Should inject --output json and --no-cli-pager
	assert.Equal(t, []string{"aws", "sso", "login", "--output", "json", "--no-cli-pager"}, mock.Calls[0])
	
	uv := NewUVRunner(mock)
	_, _ = uv.Run(context.Background(), "uv", "venv")
	
	assert.Len(t, mock.Calls, 2)
	assert.Equal(t, []string{"uv", "venv"}, mock.Calls[1])
}
