package promwrite_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/gustavooferreira/prometheus-metrics-generator/promwrite"
)

// TODO: This is an integration test!
func TestPrometheusRemoteWriter(t *testing.T) {
	t.Run("testing writer", func(t *testing.T) {
		// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		// defer cancel()
		ctx := context.Background()

		cfg := promwrite.PrometheusRemoteWriterConfig{
			Endpoint: "http://localhost:9090/api/v1/write",
		}

		remoteWriter, err := promwrite.NewPrometheusRemoteWriter(cfg)
		require.NoError(t, err)

		timeseries := []promwrite.TimeSeries{
			{
				Labels: []promwrite.Label{
					{
						Name:  "__name__",
						Value: "gf_test_metric_6_total",
					},
					{
						Name:  "gus_label",
						Value: "gus_val",
					},
				},
				Samples: []promwrite.Sample{
					{
						Time:  time.Now().UTC(),
						Value: 1000,
					},
				},
			},
		}

		err = remoteWriter.Send(ctx, timeseries)
		require.NoError(t, err)
	})
}
