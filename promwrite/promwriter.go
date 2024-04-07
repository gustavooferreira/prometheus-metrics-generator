// This package provides the logic to be able to send metrics using Prometheus Remote Write Protocol.
//
// Ref: https://prometheus.io/docs/concepts/remote_write_spec/
package promwrite

// PrometheusRemoteWriter represents the client that will send metrics to a Prometheus Remote Write enabled server.
type PrometheusRemoteWriter struct {
}
