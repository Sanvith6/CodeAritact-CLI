package validator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	
	"github.com/acme-corp/rs-cli/internal/config"
	"github.com/acme-corp/rs-cli/internal/features"
	"github.com/acme-corp/rs-cli/internal/platform"
	"github.com/acme-corp/rs-cli/internal/runner"
)

func TestCheckVersionMeetsMin(t *testing.T) {
	tests := []struct {
		current string
		min     string
		want    bool
	}{
		{"2.15.30", "2.10.0", true},
		{"2.1.0", "2.10.0", false},
		{"3.11.8", "3.9.0", true},
		{"1.0.0", "1.0.0", true},
		{"v1.2.3", "v1.2.3", true},
		{"v2.0.0", "v1.5.0", true},
	}

	for _, tc := range tests {
		got, _ := checkVersionMeetsMin(tc.current, tc.min)
		assert.Equal(t, tc.want, got, "current: %s, min: %s", tc.current, tc.min)
	}
}

func TestSystemValidator(t *testing.T) {
	mockRunner := &runner.MockRunner{
		RunSilentFunc: func(ctx context.Context, name string, args ...string) (runner.Output, error) {
			if name == "aws" {
				return runner.Output{Stdout: "aws-cli/2.15.30 Python/3.11.8 Windows/10 exe/AMD64 prompt/off"}, nil
			}
			return runner.Output{Stdout: "mock version 1.0.0"}, nil
		},
		LookPathFunc: func(name string) (string, error) {
			return "/usr/local/bin/" + name, nil
		},
	}
	
	pInfo := &platform.Info{OS: "darwin"}
	fRegistry := features.NewRegistry(config.FeaturesConfig{
		CodeArtifact: true,
	})
	
	v := NewSystemValidator(mockRunner, pInfo, fRegistry)
	
	// Test CheckTool
	awsTool := DefaultTools()[0] // aws
	res, err := v.CheckTool(context.Background(), awsTool)
	require.NoError(t, err)
	
	assert.True(t, res.Found)
	assert.Equal(t, "2.15.30", res.Version)
	assert.True(t, res.MeetsMin)
}
