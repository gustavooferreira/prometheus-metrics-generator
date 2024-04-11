package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gustavooferreira/prometheus-metrics-generator/discrete"
	"github.com/gustavooferreira/prometheus-metrics-generator/metrics"
	"github.com/gustavooferreira/prometheus-metrics-generator/promadapter"
	"github.com/gustavooferreira/prometheus-metrics-generator/promwrite"
)

func main() {
	RemoteWrite()
}

func RemoteWrite() {
	ctx := context.Background()

	fmt.Println("Running prometheus remote write")

	// create scraper
	scraper, err := metrics.NewScraper(
		metrics.ScraperConfig{
			StartTime:      time.Date(2024, 4, 11, 0, 0, 0, 0, time.UTC),
			ScrapeInterval: 15 * time.Second,
		},
		metrics.WithScraperIterationCountLimit(100),
	)
	if err != nil {
		panic(err)
	}

	// create the data generators and add them to a time series.

	dataGenerator, err := discrete.NewLinearSegmentDataGenerator(
		discrete.LinearSegmentDataGeneratorOptions{
			AmplitudeStart:      50,
			AmplitudeEnd:        70,
			IterationCountLimit: 5,
		},
	)
	if err != nil {
		panic(err)
	}

	timeSeries := discrete.NewMetricTimeSeries(
		map[string]string{"label1": "value1"},
		dataGenerator,
		// metrics.NewEndStrategyRemoveTimeSeries(),
		metrics.NewEndStrategyLoop(),
	)

	metricName := "this_is_a_metric_2"

	metric := promadapter.NewMetric(metricName, "my metric help", promadapter.MetricTypeGauge, []string{"label1"})

	err = metric.AddTimeSeries(timeSeries)
	if err != nil {
		panic(err)
	}

	prometheusRemoteWriterConfig := promwrite.PrometheusRemoteWriterConfig{
		Endpoint: "http://localhost:9090/api/v1/write",
	}

	prometheusRemoteWriter, err := promwrite.NewPrometheusRemoteWriter(prometheusRemoteWriterConfig)
	if err != nil {
		panic(err)
	}

	err = promwrite.GenerateAndImportMetrics(ctx, prometheusRemoteWriter, scraper, []promadapter.MetricObservable{metric})
	if err != nil {
		panic(err)
	}
}

func ExposeServer() {
	fmt.Println("Running prometheus metrics generator server")

	dataGenerator, err := discrete.NewLinearSegmentDataGenerator(
		discrete.LinearSegmentDataGeneratorOptions{
			AmplitudeStart:      50,
			AmplitudeEnd:        70,
			IterationCountLimit: 5,
		},
	)
	if err != nil {
		panic(err)
	}

	timeSeries := discrete.NewMetricTimeSeries(
		map[string]string{"label1": "value1"},
		dataGenerator,
		metrics.NewEndStrategyLoop(),
	)

	metric := promadapter.NewMetric("my_metric", "my metric help", promadapter.MetricTypeGauge, []string{"label1"})
	err = metric.AddTimeSeries(timeSeries)
	if err != nil {
		panic(err)
	}

	collector := promadapter.NewCollector([]promadapter.MetricObservable{metric})

	reg := prometheus.NewPedanticRegistry()

	err = reg.Register(collector)
	if err != nil {
		panic(err)
	}

	//http.Handle("/metrics", promhttp.Handler())
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	err = http.ListenAndServe(":2112", nil)
	if err != nil {
		fmt.Printf("error while listening on metrics handler: %s", err)
	}
}
