package features

import (
	"sync"
	"github.com/acme-corp/rs-cli/internal/config"
)

// Feature is an enum of supported feature flags.
type Feature string

const (
	FeatureCodeArtifact Feature = "codeartifact"
	FeatureDocker       Feature = "docker"
	FeatureJava         Feature = "java"
	FeatureTerraform    Feature = "terraform"
	FeatureKubernetes   Feature = "kubernetes"
)

// Registry manages the enabled/disabled state of features.
type Registry struct {
	features map[Feature]bool
	mu       sync.RWMutex
}

// NewRegistry creates a new feature registry based on the provided configuration.
func NewRegistry(cfg config.FeaturesConfig) *Registry {
	return &Registry{
		features: map[Feature]bool{
			FeatureCodeArtifact: cfg.CodeArtifact,
			FeatureDocker:       cfg.Docker,
			FeatureJava:         cfg.Java,
			FeatureTerraform:    cfg.Terraform,
			FeatureKubernetes:   cfg.Kubernetes,
		},
	}
}

// IsEnabled checks if a feature is enabled.
func (r *Registry) IsEnabled(f Feature) bool {
	if f == "" {
		// Empty feature requirement means it's always enabled/checked
		return true
	}
	
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.features[f]
}

// Enable turns on a feature.
func (r *Registry) Enable(f Feature) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.features[f] = true
}

// Disable turns off a feature.
func (r *Registry) Disable(f Feature) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.features[f] = false
}

// All returns a copy of all current feature states.
func (r *Registry) All() map[Feature]bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	cp := make(map[Feature]bool, len(r.features))
	for k, v := range r.features {
		cp[k] = v
	}
	return cp
}
