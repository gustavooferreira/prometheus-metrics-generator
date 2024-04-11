package promwrite

import (
	"fmt"
	"net/http"
	"time"
)

type PrometheusRemoteWriterConfig struct {
	// Endpoint represents the URL the client will send the samples to.
	// Ex: http://localhost:9090/api/v1/write
	Endpoint string

	// -------------------------------------------------
	// Unexported fields are set via a functional option
	// -------------------------------------------------

	// httpClient is the HTTP client used to send the samples to the Remote Write Server.
	httpClient *http.Client
}

// validate validates the config struct.
func (c *PrometheusRemoteWriterConfig) validate() error {
	if c.Endpoint == "" {
		return fmt.Errorf("prometheus endpoint cannot be nil")
	}

	return nil
}

// applyDefaults applies defaults to the fields set via functional options.
func (c *PrometheusRemoteWriterConfig) applyDefaults() {
	c.httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}
}

// applyFunctionalOptions applies the set of PrometheusRemoteWriterConfigOption onto the PrometheusRemoteWriterConfig.
func (c *PrometheusRemoteWriterConfig) applyFunctionalOptions(opts ...PrometheusRemoteWriterConfigOption) {
	for _, opt := range opts {
		opt(c)
	}
}

// Functional Options -----------------

type PrometheusRemoteWriterConfigOption func(c *PrometheusRemoteWriterConfig)

// WithHTTPClient sets the HTTP client to be used by the PrometheusRemoteWriter.
func WithHTTPClient(httpClient *http.Client) PrometheusRemoteWriterConfigOption {
	return func(c *PrometheusRemoteWriterConfig) {
		c.httpClient = httpClient
	}
}
