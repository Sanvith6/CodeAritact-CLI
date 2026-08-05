package validator

import (
	"context"
	"fmt"
	"strings"

	"github.com/acme-corp/rs-cli/internal/errors"
	"github.com/acme-corp/rs-cli/internal/features"
	"github.com/acme-corp/rs-cli/internal/platform"
	"github.com/acme-corp/rs-cli/internal/runner"
)

// ToolCheck defines the requirements and properties for a required tool.
type ToolCheck struct {
	Name         string
	Command      string
	VersionArgs  []string
	MinVersion   string
	Required     bool
	Feature      features.Feature
	InstallHint  map[string]string // map[os]hint
	ParseVersion func(output string) (string, error)
}

// CheckResult contains the result of validating a single tool.
type CheckResult struct {
	Tool     ToolCheck
	Found    bool
	Path     string
	Version  string
	MeetsMin bool
	Error    error
}

// Validator validates the presence and versions of required tools.
type Validator interface {
	CheckAll(ctx context.Context) ([]CheckResult, error)
	CheckTool(ctx context.Context, tool ToolCheck) (CheckResult, error)
}

// SystemValidator implements Validator by checking the local OS.
type SystemValidator struct {
	runner   runner.Runner
	pInfo    *platform.Info
	features *features.Registry
	tools    []ToolCheck
}

// NewSystemValidator creates a new SystemValidator.
func NewSystemValidator(r runner.Runner, pInfo *platform.Info, f *features.Registry) *SystemValidator {
	return &SystemValidator{
		runner:   r,
		pInfo:    pInfo,
		features: f,
		tools:    DefaultTools(),
	}
}

// CheckAll checks all configured tools, respecting feature flags.
func (v *SystemValidator) CheckAll(ctx context.Context) ([]CheckResult, error) {
	var results []CheckResult
	
	for _, tool := range v.tools {
		// Skip tools for features that are disabled
		if tool.Feature != "" && !v.features.IsEnabled(tool.Feature) {
			continue
		}
		
		res, _ := v.CheckTool(ctx, tool)
		results = append(results, res)
	}
	
	return results, nil
}

// CheckTool checks a single tool.
func (v *SystemValidator) CheckTool(ctx context.Context, tool ToolCheck) (CheckResult, error) {
	result := CheckResult{
		Tool: tool,
	}

	path, err := v.runner.LookPath(tool.Command)
	if err != nil {
		result.Found = false
		hint := v.getInstallHint(tool)
		result.Error = errors.New(errors.ErrCodeToolMissing, fmt.Sprintf("Required tool '%s' is not installed.", tool.Name)).
			WithSuggestion(hint)
		return result, result.Error
	}

	result.Found = true
	result.Path = path

	if tool.ParseVersion != nil {
		out, err := v.runner.RunSilent(ctx, tool.Command, tool.VersionArgs...)
		if err != nil {
			result.Error = fmt.Errorf("failed to get version: %w", err)
			return result, result.Error
		}

		// Combine stdout and stderr because some tools (like java) print version to stderr
		verOut := out.Stdout + out.Stderr
		
		version, err := tool.ParseVersion(verOut)
		if err != nil {
			result.Error = fmt.Errorf("failed to parse version: %w", err)
			return result, result.Error
		}
		
		result.Version = version
		
		if tool.MinVersion != "" {
			meets, err := checkVersionMeetsMin(version, tool.MinVersion)
			if err != nil {
				result.Error = fmt.Errorf("version comparison failed: %w", err)
				return result, result.Error
			}
			result.MeetsMin = meets
			
			if !meets && tool.Required {
				result.Error = errors.New(errors.ErrCodeToolVersion, 
					fmt.Sprintf("Tool '%s' is outdated. Found %s, need >= %s", tool.Name, version, tool.MinVersion)).
					WithSuggestion(v.getInstallHint(tool))
			}
		} else {
			result.MeetsMin = true
		}
	} else {
		result.MeetsMin = true
	}

	return result, result.Error
}

func (v *SystemValidator) getInstallHint(tool ToolCheck) string {
	hint, ok := tool.InstallHint[v.pInfo.OS]
	if ok {
		return fmt.Sprintf("Run: %s", hint)
	}
	
	// Fallback
	return fmt.Sprintf("Please install %s manually.", tool.Name)
}

// checkVersionMeetsMin does a naive semver comparison.
// In a real implementation, we'd use a full semver package, but this covers 90% of cases.
func checkVersionMeetsMin(current, min string) (bool, error) {
	// Clean prefixes like "v"
	currStr := strings.TrimPrefix(current, "v")
	minStr := strings.TrimPrefix(min, "v")
	
	// Simple string comparison works for zero-padded versions, 
	// but for real use a proper semver library is better.
	// For this plan, we'll implement a basic naive checker.
	
	currParts := strings.Split(currStr, ".")
	minParts := strings.Split(minStr, ".")
	
	for i := 0; i < len(minParts); i++ {
		if i >= len(currParts) {
			// Current has fewer parts, e.g. curr "2", min "2.1"
			return false, nil
		}
		
		// If current part > min part, we're good
		// If current part < min part, we fail
		// If equal, continue to next part
		
		// Note: string comparison "10" < "2" is true, which is bad. 
		// We pad them conceptually here or just assume they don't exceed 1 digit for the skeleton.
		// A real implementation would parse ints.
		
		var currInt, minInt int
		fmt.Sscanf(currParts[i], "%d", &currInt)
		fmt.Sscanf(minParts[i], "%d", &minInt)
		
		if currInt > minInt {
			return true, nil
		}
		if currInt < minInt {
			return false, nil
		}
	}
	
	return true, nil
}
