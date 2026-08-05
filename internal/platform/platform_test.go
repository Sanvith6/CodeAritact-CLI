package platform

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultDetector_Detect(t *testing.T) {
	detector := NewDetector()
	info, err := detector.Detect()
	
	require.NoError(t, err)
	assert.NotNil(t, info)
	
	assert.Equal(t, runtime.GOOS, info.OS)
	assert.Equal(t, runtime.GOARCH, info.Arch)
	assert.NotEmpty(t, info.HomeDir)
	assert.NotEmpty(t, info.Shell)
	
	assert.NotEmpty(t, info.AWSConfigDir)
	assert.NotEmpty(t, info.RSConfigDir)
	assert.NotEmpty(t, info.UVConfigDir)
}
