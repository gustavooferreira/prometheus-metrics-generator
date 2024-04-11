package promwrite

import "context"

// PrometheusRemoteWriterBuffer does lots of things.
// Create a new one.
// Function that closes the buffer, will make sure to send all the remaining stale markers.
// Accumulate samples and send after X amount.
// It retries in the even of errors.
// Parallelizes the sending of metrics using shards, etc.

// TODO: Send Stale markers for metrics that didn't have them sent already!
// For this we need to keep track of all the time series we've sent (the key should be all the labels serialized,
// make sure the labels are sorted first)
// Then we should remove them from the set as we see them being marked as stale.
// Finally send the stale markers for the left overs.

type PrometheusRemoteWriterBuffer struct {
	prometheusRemoteWriter *PrometheusRemoteWriter
}

// NewPrometheusRemoteWriterBuffer creates a new instance of PrometheusRemoteWriterBuffer.
func NewPrometheusRemoteWriterBuffer(prometheusRemoteWriter *PrometheusRemoteWriter) *PrometheusRemoteWriterBuffer {
	return &PrometheusRemoteWriterBuffer{
		prometheusRemoteWriter: prometheusRemoteWriter,
	}
}

func (prwb *PrometheusRemoteWriterBuffer) Send(ctx context.Context, timeseries []TimeSeries) error {

	return nil
}

// Flush flushes any remaining data in the buffers as well as making sure to send stale markers for all time series
// that have not been observed sending stale markers.
func (prwb *PrometheusRemoteWriterBuffer) Flush() error {

	return nil
}

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
