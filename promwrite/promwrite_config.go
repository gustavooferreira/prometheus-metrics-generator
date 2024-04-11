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

	// headers represents the headers to be sent with every single request.
	// Extra headers can be sent when calling the Send() method.
	headers map[string][]string
}

// validate validates the config struct.
func (c *PrometheusRemoteWriterConfig) validate() error {
	if c.Endpoint == "" {
		return fmt.Errorf("prometheus endpoint cannot be nil")
	}

	if err := validateHTTPHeaders(c.headers); err != nil {
		return fmt.Errorf("failed validating headers: %w", err)
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

// WithHeaders sets the headers to be sent in all Prometheus Remote Write requests.
func WithHeaders(httpHeaders http.Header) PrometheusRemoteWriterConfigOption {
	return func(c *PrometheusRemoteWriterConfig) {
		c.headers = httpHeaders
	}
}
