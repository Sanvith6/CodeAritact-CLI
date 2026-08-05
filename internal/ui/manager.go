package ui

import (
	"sync"
)

// Manager is a centralized UI manager.
// It tracks interactive state and ensures multiple spinners don't conflict.
type Manager struct {
	nonInteractive bool
	mu             sync.Mutex
	activeSpinner  *Spinner
}

// NewManager creates a new UI manager.
func NewManager(nonInteractive bool) *Manager {
	return &Manager{
		nonInteractive: nonInteractive || !IsInteractive(nonInteractive),
	}
}

// StartSpinner starts a spinner, stopping any existing one.
func (m *Manager) StartSpinner(msg string) *Spinner {
	if m.nonInteractive {
		PrintInfo(msg)
		return nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.activeSpinner != nil {
		m.activeSpinner.Stop()
	}

	s := NewSpinner(msg)
	s.Start()
	m.activeSpinner = s
	return s
}

// StopSpinner stops the active spinner.
func (m *Manager) StopSpinner(success bool, msg string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.activeSpinner == nil {
		if msg != "" {
			if success {
				PrintSuccess(msg)
			} else {
				PrintError(msg)
			}
		}
		return
	}

	if msg != "" {
		if success {
			m.activeSpinner.StopSuccess(msg)
		} else {
			m.activeSpinner.StopError(msg)
		}
	} else {
		m.activeSpinner.Stop()
	}
	
	m.activeSpinner = nil
}
