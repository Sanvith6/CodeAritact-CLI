package validator

import (
	"regexp"
	"strings"

	"github.com/acme-corp/rs-cli/internal/features"
)

// DefaultTools returns the standard list of tools required by the platform.
func DefaultTools() []ToolCheck {
	return []ToolCheck{
		{
			Name:        "AWS CLI v2",
			Command:     "aws",
			VersionArgs: []string{"--version"},
			MinVersion:  "2.10.0",
			Required:    true,
			InstallHint: map[string]string{
				"darwin":  "brew install awscli",
				"windows": "winget install Amazon.AWSCLI",
				"linux":   "curl \"https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip\" -o \"awscliv2.zip\" && unzip awscliv2.zip && sudo ./aws/install",
			},
			ParseVersion: func(output string) (string, error) {
				// aws-cli/2.15.30 Python/3.11.8 ...
				re := regexp.MustCompile(`aws-cli/([0-9]+\.[0-9]+\.[0-9]+)`)
				matches := re.FindStringSubmatch(output)
				if len(matches) > 1 {
					return matches[1], nil
				}
				return "", nil // Couldn't parse, assume it's weird
			},
		},
		{
			Name:        "uv",
			Command:     "uv",
			VersionArgs: []string{"--version"},
			MinVersion:  "0.1.0",
			Required:    true,
			InstallHint: map[string]string{
				"darwin":  "curl -LsSf https://astral.sh/uv/install.sh | sh",
				"windows": "powershell -ExecutionPolicy ByPass -c \"irm https://astral.sh/uv/install.ps1 | iex\"",
				"linux":   "curl -LsSf https://astral.sh/uv/install.sh | sh",
			},
			ParseVersion: func(output string) (string, error) {
				// uv 0.1.20
				parts := strings.Split(strings.TrimSpace(output), " ")
				if len(parts) >= 2 {
					return parts[1], nil
				}
				return "", nil
			},
		},
		{
			Name:        "Git",
			Command:     "git",
			VersionArgs: []string{"--version"},
			MinVersion:  "2.0.0",
			Required:    true,
			InstallHint: map[string]string{
				"darwin":  "brew install git",
				"windows": "winget install Git.Git",
				"linux":   "sudo apt-get install git",
			},
			ParseVersion: func(output string) (string, error) {
				// git version 2.39.3 (Apple Git-146)
				parts := strings.Split(strings.TrimSpace(output), " ")
				if len(parts) >= 3 {
					return parts[2], nil
				}
				return "", nil
			},
		},
		{
			Name:        "Python",
			Command:     "python", // Might be python3 on some systems, this is a simplified check
			VersionArgs: []string{"--version"},
			MinVersion:  "3.9.0",
			Required:    true,
			InstallHint: map[string]string{
				"darwin":  "brew install python",
				"windows": "winget install Python.Python.3.11",
				"linux":   "sudo apt-get install python3",
			},
			ParseVersion: func(output string) (string, error) {
				// Python 3.11.8
				parts := strings.Split(strings.TrimSpace(output), " ")
				if len(parts) >= 2 {
					return parts[1], nil
				}
				return "", nil
			},
		},
		{
			Name:        "Docker",
			Command:     "docker",
			VersionArgs: []string{"--version"},
			MinVersion:  "20.10.0",
			Required:    true,
			Feature:     features.FeatureDocker,
			InstallHint: map[string]string{
				"darwin":  "brew install --cask docker",
				"windows": "winget install Docker.DockerDesktop",
				"linux":   "sudo apt-get install docker-ce",
			},
			ParseVersion: func(output string) (string, error) {
				// Docker version 25.0.3, build 4debf41
				parts := strings.Split(strings.TrimSpace(output), " ")
				if len(parts) >= 3 {
					return strings.TrimSuffix(parts[2], ","), nil
				}
				return "", nil
			},
		},
	}
}
