package state

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManager(t *testing.T) {
	tmpDir := t.TempDir()
	stateFile := filepath.Join(tmpDir, "state.json")
	
	manager := NewManager(stateFile)
	
	// Test Load (should return empty state when missing)
	s, err := manager.Load()
	require.NoError(t, err)
	assert.False(t, s.BootstrapCompleted)
	assert.Equal(t, 1, s.Version)
	
	// Test Save
	s.BootstrapCompleted = true
	s.CurrentProfile = "rs-dev"
	s.CLIVersion = "1.0.0"
	
	err = manager.Save(s)
	require.NoError(t, err)
	
	// Test Load (should return saved state)
	s2, err := manager.Load()
	require.NoError(t, err)
	assert.True(t, s2.BootstrapCompleted)
	assert.Equal(t, "rs-dev", s2.CurrentProfile)
	
	// Test Update
	err = manager.Update(func(st *State) {
		st.CurrentOrg = "acme"
	})
	require.NoError(t, err)
	
	s3, err := manager.Load()
	require.NoError(t, err)
	assert.Equal(t, "acme", s3.CurrentOrg)
	assert.True(t, manager.IsBootstrapped())
	
	// Test Token Expiry
	assert.False(t, manager.IsTokenValid())
	
	err = manager.Update(func(st *State) {
		st.TokenExpiry = time.Now().Add(1 * time.Hour)
	})
	require.NoError(t, err)
	assert.True(t, manager.IsTokenValid())
	
	// Test Clear
	err = manager.Clear()
	require.NoError(t, err)
	
	s4, err := manager.Load()
	require.NoError(t, err)
	assert.False(t, s4.BootstrapCompleted)
}
