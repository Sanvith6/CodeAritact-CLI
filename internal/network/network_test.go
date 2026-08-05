package network

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPChecker_CheckInternet(t *testing.T) {
	// Note: This test requires internet access. In a real CI environment,
	// you might want to mock the HTTP client or skip if no internet.
	checker := NewHTTPChecker()
	err := checker.CheckInternet(context.Background())
	assert.NoError(t, err)
}

func TestHTTPChecker_CheckAWSEndpoint(t *testing.T) {
	// Note: This test requires internet access and DNS resolution for AWS.
	checker := NewHTTPChecker()
	err := checker.CheckAWSEndpoint(context.Background(), "us-east-1")
	assert.NoError(t, err)
}
