package main

import (
	"context"
	"time"

	"github.com/gustavooferreira/prometheus-metrics-generator/promwrite"
)

func main() {

	prometheusRemoteWriterConfig := promwrite.PrometheusRemoteWriterConfig{
		Endpoint: "http://localhost:9090/api/v1/write",
	}

	prometheusRemoteWriter, err := promwrite.NewPrometheusRemoteWriter(prometheusRemoteWriterConfig)
	if err != nil {
		panic(err)
	}

	timeseries := promwrite.TimeSeries{
		Labels: []promwrite.Label{
			{
				Name:  "__name__",
				Value: "my_metric_1",
			},
			{
				Name:  "label_1",
				Value: "value11",
			},
		},
		Samples: []promwrite.Sample{
			{
				Time:  time.Now(),
				Value: 100,
			},
		},
	}

	ctx := context.Background()

	err = prometheusRemoteWriter.Send(ctx, []promwrite.TimeSeries{timeseries})
	if err != nil {
		panic(err)
	}
}
