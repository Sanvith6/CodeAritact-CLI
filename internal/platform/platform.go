package platform

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

// Info contains platform-specific paths and metadata.
type Info struct {
	OS           string // "windows", "darwin", "linux"
	Arch         string // "amd64", "arm64", etc.
	HomeDir      string
	Shell        string
	IsWSL        bool
	AWSConfigDir string
	RSConfigDir  string
	UVConfigDir  string
}

// Detector is responsible for detecting platform details.
type Detector interface {
	Detect() (*Info, error)
}

// DefaultDetector provides standard OS-based detection.
type DefaultDetector struct{}

// NewDetector creates a new default platform detector.
func NewDetector() *DefaultDetector {
	return &DefaultDetector{}
}

// Detect gathers platform information.
func (d *DefaultDetector) Detect() (*Info, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		// Fallback for some weird environments
		usr, err := user.Current()
		if err != nil {
			return nil, fmt.Errorf("could not determine home directory: %w", err)
		}
		home = usr.HomeDir
	}

	info := &Info{
		OS:          runtime.GOOS,
		Arch:        runtime.GOARCH,
		HomeDir:     home,
		IsWSL:       detectWSL(),
		Shell:       detectShell(),
		RSConfigDir: filepath.Join(home, ".rs"),
	}

	// Platform specific paths
	switch info.OS {
	case "windows":
		info.AWSConfigDir = filepath.Join(home, ".aws")
		// uv defaults to %APPDATA%\uv on windows, but checks other places.
		// For our config needs (uv.toml), it's typically APPDATA.
		appData := os.Getenv("APPDATA")
		if appData != "" {
			info.UVConfigDir = filepath.Join(appData, "uv")
		} else {
			info.UVConfigDir = filepath.Join(home, "AppData", "Roaming", "uv")
		}
	case "darwin":
		info.AWSConfigDir = filepath.Join(home, ".aws")
		info.UVConfigDir = filepath.Join(home, ".config", "uv")
	default: // Linux and others
		info.AWSConfigDir = filepath.Join(home, ".aws")
		info.UVConfigDir = filepath.Join(home, ".config", "uv")
	}

	return info, nil
}

func detectWSL() bool {
	if runtime.GOOS != "linux" {
		return false
	}
	content, err := os.ReadFile("/proc/version")
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(content)), "microsoft")
}

func detectShell() string {
	if runtime.GOOS == "windows" {
		// Simple heuristic for Windows
		if ps := os.Getenv("PSModulePath"); ps != "" {
			return "powershell"
		}
		return "cmd"
	}
	
	shell := os.Getenv("SHELL")
	if shell != "" {
		return filepath.Base(shell)
	}
	
	return "bash" // fallback
}
