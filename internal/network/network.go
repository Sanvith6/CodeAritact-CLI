package network

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/acme-corp/rs-cli/internal/errors"
)

// Checker interface for validating network connectivity.
type Checker interface {
	CheckInternet(ctx context.Context) error
	CheckAWSEndpoint(ctx context.Context, region string) error
}

// HTTPChecker implements Checker using standard HTTP client.
type HTTPChecker struct {
	client *http.Client
}

// NewHTTPChecker creates a new HTTPChecker.
func NewHTTPChecker() *HTTPChecker {
	return &HTTPChecker{
		client: &http.Client{
			Timeout: 5 * time.Second, // Fast timeout for responsiveness
		},
	}
}

// CheckInternet checks connectivity to a reliable external host.
func (c *HTTPChecker) CheckInternet(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, "https://www.google.com", nil)
	if err != nil {
		return err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return errors.Wrap(errors.ErrCodeNetwork, err, "Failed to connect to the internet.").
			WithSuggestion("Check your Wi-Fi, Ethernet, or VPN connection.")
	}
	defer res.Body.Close()

	if res.StatusCode >= 500 {
		return errors.New(errors.ErrCodeNetwork, fmt.Sprintf("Unexpected status code checking internet: %d", res.StatusCode))
	}

	return nil
}

// CheckAWSEndpoint checks connectivity to the regional STS endpoint.
func (c *HTTPChecker) CheckAWSEndpoint(ctx context.Context, region string) error {
	url := fmt.Sprintf("https://sts.%s.amazonaws.com", region)
	
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		return err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return errors.Wrap(errors.ErrCodeNetwork, err, fmt.Sprintf("Failed to connect to AWS STS in %s.", region)).
			WithSuggestion("Ensure your corporate VPN allows traffic to AWS endpoints.")
	}
	defer res.Body.Close()

	// It's normal for STS to return 4xx (like 403 Forbidden or 405 Method Not Allowed for HEAD)
	// We just care that the TCP/TLS connection succeeded and we got an HTTP response.
	return nil
}
