// This package provides the logic to be able to send metrics using the Prometheus Remote Write Protocol.
//
// Ref: https://prometheus.io/docs/concepts/remote_write_spec/
package promwrite

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math"
	"net/http"
	"regexp"
	"sort"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/prometheus/prometheus/prompb"
)

var (
	metricNameRegex = regexp.MustCompile(`^[a-zA-Z_:][a-zA-Z0-9_:]*$`)
	labelNameRegex  = regexp.MustCompile(`[a-zA-Z_]([a-zA-Z0-9_])*.`)

	// staleMarker is the last value of a time series and signals the time series will not be appended to anymore.
	// Prometheus code: https://pkg.go.dev/github.com/prometheus/prometheus/pkg/value#pkg-constants
	// Spec Ref:
	//  Stale markers MUST be signalled by the special NaN value 0x7ff0000000000002.
	//  This value MUST NOT be used otherwise.
	staleMarker = math.Float64frombits(0x7ff0000000000002)
)

// PrometheusRemoteWriter represents the client that will send metrics to a Prometheus Remote Write enabled server.
type PrometheusRemoteWriter struct {
	cfg PrometheusRemoteWriterConfig
}

// NewPrometheusRemoteWriter creates a new instance of PrometheusRemoteWriter.
func NewPrometheusRemoteWriter(cfg PrometheusRemoteWriterConfig, opts ...PrometheusRemoteWriterConfigOption) (*PrometheusRemoteWriter, error) {
	cfg.applyDefaults()

	cfg.applyFunctionalOptions(opts...)

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("error validating prometheus remote writer configuration: %w", err)
	}

	promRemoteWriter := &PrometheusRemoteWriter{
		cfg: cfg,
	}

	return promRemoteWriter, nil
}

// Send sends HTTP requests to Prometheus Remote Write compatible API endpoint including Prometheus, Cortex,
// VictoriaMetrics, etc.
// Spec Ref:
//
//	Prometheus remote write compatible senders MUST send stale markers when a time series will no longer be appended to.
func (prw *PrometheusRemoteWriter) Send(ctx context.Context, timeseries []TimeSeries, opts ...WriteOption) error {
	writeOptions := writeOptions{}

	writeOptions.applyFunctionalOptions(opts...)

	if err := writeOptions.validate(); err != nil {
		return fmt.Errorf("failed validating write options: %w", err)
	}

	protoTimeSeries, err := toProtoTimeSeries(timeseries)
	if err != nil {
		return fmt.Errorf("error converting time series to protobuf format: %w", err)
	}

	// Marshal proto and compress.
	protoReq := &prompb.WriteRequest{
		Timeseries: protoTimeSeries,
		// Sending metadata will be supported in a future release of prometheus.
		// The data types have been defined already, but the prometheus remote write server doesn't do anything with
		// the metadata sent in the request.
		// Example:
		//  Metadata: []prompb.MetricMetadata{
		//  	{
		//  		Type:             prompb.MetricMetadata_COUNTER,
		//  		Help:             "This is a certain kind of counter metric.",
		//  		MetricFamilyName: "my_metric_seconds_total",
		//  	},
		//  },
	}

	protoReqBytes, err := proto.Marshal(protoReq)
	if err != nil {
		return fmt.Errorf("failed marshaling remote write protobuf request: %w", err)
	}

	protoReqBytesCompressed := snappy.Encode(nil, protoReqBytes)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, prw.cfg.Endpoint, bytes.NewReader(protoReqBytesCompressed))
	if err != nil {
		return fmt.Errorf("failed to create request for remote write operation: %w", err)
	}

	httpReq.Header.Set("User-Agent", "prometheus-metrics-generator")
	httpReq.Header.Set("Content-Type", "application/x-protobuf")
	httpReq.Header.Set("Content-Encoding", "snappy")
	httpReq.Header.Set("X-Prometheus-Remote-Write-Version", "0.1.0")

	for headerKey, headerValues := range writeOptions.headers {
		for _, headerValue := range headerValues {
			httpReq.Header.Add(headerKey, headerValue)
		}
	}

	// Send http request.
	httpResp, err := prw.cfg.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to make request for remote write operation: %w", err)
	}
	defer httpResp.Body.Close()

	// The response body from the remote write receiver SHOULD be empty; clients MUST ignore the response body.
	// The response body is RESERVED for future use.
	//
	// Developer note: We read all the contents of the body anyway otherwise we can't properly reuse the TCP connection.
	//                 In case of non-successful status code, prometheus seems to be returning a simple string.
	responseBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response data for remote write operation: %w", err)
	}

	// We check if the response is in the 2xx range because the server might send a 200 or a 204 in case of no content.
	// Spec Ref:
	//  Prometheus remote Write compatible receivers MUST respond with a HTTP 2xx status code when the write is
	//  successful.
	//  They MUST respond with HTTP status code 5xx when the write fails and SHOULD be retried.
	//  They MUST respond with HTTP status code 4xx when the request is invalid, will never be able to succeed and
	//  should not be retried.
	if httpResp.StatusCode/100 != 2 {
		return fmt.Errorf("failed http request for remote write operation: got bad status code %d and body %q",
			httpResp.StatusCode,
			string(responseBody),
		)
	}

	return nil
}

type writeOptions struct {
	headers map[string][]string
}

// applyFunctionalOptions applies the set of WriteOption onto the writeOptions.
func (wo *writeOptions) applyFunctionalOptions(opts ...WriteOption) {
	for _, opt := range opts {
		opt(wo)
	}
}

// validate validates the writeOptions struct.
func (wo *writeOptions) validate() error {
	if _, ok := wo.headers["X-Prometheus-Remote-Write-Version"]; ok {
		return fmt.Errorf("failed validating write options: setting header %q not allowed", "X-Prometheus-Remote-Write-Version")
	}

	if _, ok := wo.headers["Content-Type"]; ok {
		return fmt.Errorf("failed validating write options: setting header %q not allowed", "Content-Type")
	}

	if _, ok := wo.headers["Content-Encoding"]; ok {
		return fmt.Errorf("failed validating write options: setting header %q not allowed", "Content-Encoding")
	}

	// if the User-Agent header is set, it cannot be empty
	if userAgentHeaderValues, ok := wo.headers["User-Agent"]; ok {
		for _, headerValue := range userAgentHeaderValues {
			if headerValue == "" {
				return fmt.Errorf("failed validating write options: user agent header value cannot be empty")
			}
		}
	}

	return nil
}

type WriteOption func(opts *writeOptions)

// WithWriteHeader adds an HTTP header to be used with the HTTP request the Send method() performs.
// Pass this functional option multiple times to set multiple headers.
func WithWriteHeader(key string, value string) WriteOption {
	return func(o *writeOptions) {
		if header, ok := o.headers[key]; ok {
			o.headers[key] = append(header, value)
			return
		}

		o.headers[key] = []string{value}
	}
}

// toProtoTimeSeries converts our []TimeSeries structs into protobuf structs, ready to be sent down the wire.
func toProtoTimeSeries(timeSeries []TimeSeries) ([]prompb.TimeSeries, error) {
	protoTimeSeries := make([]prompb.TimeSeries, len(timeSeries))

	for i, singleTimeSeries := range timeSeries {
		labels, err := convertLabels(singleTimeSeries.Labels)
		if err != nil {
			return nil, fmt.Errorf("error converting labels for time series: %w", err)
		}

		samples := make([]prompb.Sample, len(singleTimeSeries.Samples))
		for sampleIndex, sample := range singleTimeSeries.Samples {
			samples[sampleIndex] = prompb.Sample{
				// Timestamps MUST be int64 counted as milliseconds since the Unix epoch.
				Timestamp: sample.Time.UnixMilli(),
				// Values MUST be float64.
				Value: sample.Value,
			}
		}

		protoSingleTimeSeries := prompb.TimeSeries{
			Labels:  labels,
			Samples: samples,
		}

		protoTimeSeries[i] = protoSingleTimeSeries
	}

	return protoTimeSeries, nil
}

// convertLabels checks whether the labels are valid, formats and converts them according to the spec.
//
// Spec Ref:
// The complete set of labels MUST be sent with each sample. Whatsmore, the label set associated with samples:
//   - SHOULD contain a __name__ label.
//   - MUST NOT contain repeated label names.
//   - MUST have label names sorted in lexicographical order.
//   - MUST NOT contain any empty label names or values.
//
// Senders MUST only send valid metric names, label names, and label values:
//   - Metric names MUST adhere to the regex [a-zA-Z_:]([a-zA-Z0-9_:])*.
//   - Label names MUST adhere to the regex [a-zA-Z_]([a-zA-Z0-9_])*.
//   - Label values MAY be any sequence of UTF-8 characters .
//
// Label names beginning with "__" are RESERVED for system usage and SHOULD NOT be used.
//
// Labels - Every series MAY include a "job" and/or "instance" label, as these are typically added by service discovery
// in the Sender. These are not mandatory.
func convertLabels(labels []Label) ([]prompb.Label, error) {
	protoLabels := make([]prompb.Label, len(labels))

	// keep track of seen labels in order to flag labels that are duplicates.
	seenLabels := make(map[string]struct{})

	for labelIndex, label := range labels {
		if label.Name == "" {
			return nil, fmt.Errorf("label name must not be empty")
		}

		if _, ok := seenLabels[label.Name]; ok {
			return nil, fmt.Errorf("labels must not be repeated")
		}
		seenLabels[label.Name] = struct{}{}

		if label.Name == "__name__" {
			// check regex pattern for metric name
			if !metricNameRegex.MatchString(label.Name) {
				return nil, fmt.Errorf("metric name must comply with regex pattern specified in the remote write spec")
			}
		} else {
			// check regex pattern for label name.
			if !labelNameRegex.MatchString(label.Name) {
				return nil, fmt.Errorf("label name must comply with regex pattern specified in the remote write spec")
			}
		}

		protoLabels[labelIndex] = prompb.Label{
			Name:  label.Name,
			Value: label.Value,
		}
	}

	// Sort labels lexicographically.
	sort.Slice(protoLabels, func(i int, j int) bool {
		return protoLabels[i].Name < protoLabels[j].Name
	})

	return protoLabels, nil
}

// ----------

// For the chunks sender we should take this into account:
// Prometheus Remote Write compatible senders MUST send samples for any given series in timestamp order.
// Prometheus Remote Write compatible Senders MAY send multiple requests for different series in parallel.

// Basically specify a maximum of samples per send.
// And have "shards" where we can send multiple timeseries per shard since it's fine to send them in parallel.

// Sharding the current sharding scheme in Prometheus for remote write parallelisation is very much an implementation detail, and isn’t part of the spec. When senders do implement parallelisation they MUST preserve per-series sample ordering.
// How can we parallelise requests with the in-order constraint? Samples must be in-order for a given series. Remote write requests can be sent in parallel as long as they are for different series. In Prometheus, we shard the samples by their labels into separate queues, and then writes happen sequentially in each queue. This guarantees samples for the same series are delivered in order, but samples for different series are sent in parallel - and potentially "out of order" between different series.

// Prometheus Remote Write compatible senders MUST retry write requests on HTTP 5xx responses and MUST use a backoff
// algorithm to prevent overwhelming the server.
// They MUST NOT retry write requests on HTTP 2xx and 4xx responses other than 429.
// They MAY retry on HTTP 429 responses, which could result in senders "falling behind" if the server cannot keep up.
// This is done to ensure data is not lost when there are server side errors, and progress is made when there are
// client side errors.
