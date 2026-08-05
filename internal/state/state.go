package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// State represents the CLI's persisted state.
type State struct {
	Version            int       `json:"version"`
	BootstrapCompleted bool      `json:"bootstrap_completed"`
	BootstrapTimestamp time.Time `json:"bootstrap_timestamp,omitempty"`
	LastLogin          time.Time `json:"last_login,omitempty"`
	CurrentProfile     string    `json:"current_profile,omitempty"`
	CurrentRepository  string    `json:"current_repository,omitempty"`
	CurrentOrg         string    `json:"current_org,omitempty"`
	CurrentEnvironment string    `json:"current_environment,omitempty"`
	CLIVersion         string    `json:"cli_version"`
	TokenExpiry        time.Time `json:"token_expiry,omitempty"`
}

// Manager handles reading and writing the state file.
type Manager struct {
	path string
	mu   sync.Mutex
}

// NewManager creates a new state manager for the given file path.
func NewManager(path string) *Manager {
	return &Manager{
		path: path,
	}
}

// Load reads the state from disk. Returns a new empty State if the file doesn't exist.
func (m *Manager) Load() (*State, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.loadNoLock()
}

func (m *Manager) loadNoLock() (*State, error) {
	data, err := os.ReadFile(m.path)
	if err != nil {
		if os.IsNotExist(err) {
			return &State{Version: 1}, nil
		}
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("failed to parse state file: %w", err)
	}

	return &s, nil
}

// Save writes the state to disk atomically.
func (m *Manager) Save(s *State) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.saveNoLock(s)
}

func (m *Manager) saveNoLock(s *State) error {
	if err := os.MkdirAll(filepath.Dir(m.path), 0700); err != nil {
		return fmt.Errorf("failed to create state directory: %w", err)
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize state: %w", err)
	}

	// Write to temp file first for atomic write
	tmpFile := m.path + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0600); err != nil {
		return fmt.Errorf("failed to write temp state file: %w", err)
	}

	// Rename temp file to actual file
	if err := os.Rename(tmpFile, m.path); err != nil {
		_ = os.Remove(tmpFile)
		return fmt.Errorf("failed to rename state file: %w", err)
	}

	return nil
}

// Update provides atomic read-modify-write on the state.
func (m *Manager) Update(fn func(s *State)) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	s, err := m.loadNoLock()
	if err != nil {
		return err
	}

	fn(s)

	return m.saveNoLock(s)
}

// Clear completely deletes the state file.
func (m *Manager) Clear() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	err := os.Remove(m.path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete state file: %w", err)
	}
	return nil
}

// IsBootstrapped checks if the bootstrap process has been completed successfully.
func (m *Manager) IsBootstrapped() bool {
	s, err := m.Load()
	if err != nil {
		return false
	}
	return s.BootstrapCompleted
}

// IsTokenValid checks if the currently stored token is still valid.
func (m *Manager) IsTokenValid() bool {
	s, err := m.Load()
	if err != nil {
		return false
	}
	// Buffer of 5 minutes to avoid edge cases
	return time.Now().Add(5 * time.Minute).Before(s.TokenExpiry)
}
