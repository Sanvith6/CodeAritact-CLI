package features

import (
	"testing"
	"github.com/acme-corp/rs-cli/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestRegistry(t *testing.T) {
	cfg := config.FeaturesConfig{
		CodeArtifact: true,
		Docker:       false,
	}
	
	registry := NewRegistry(cfg)
	
	// Check initial state
	assert.True(t, registry.IsEnabled(FeatureCodeArtifact))
	assert.False(t, registry.IsEnabled(FeatureDocker))
	assert.True(t, registry.IsEnabled("")) // Empty feature is always enabled
	
	// Toggle states
	registry.Enable(FeatureDocker)
	assert.True(t, registry.IsEnabled(FeatureDocker))
	
	registry.Disable(FeatureCodeArtifact)
	assert.False(t, registry.IsEnabled(FeatureCodeArtifact))
	
	// Check all map
	all := registry.All()
	assert.True(t, all[FeatureDocker])
	assert.False(t, all[FeatureCodeArtifact])
}
