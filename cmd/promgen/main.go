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
		metrics.NewEndStrategyRemoveTimeSeries(),
		// metrics.NewEndStrategyLoop(),
	)

	metricName := "this_is_a_metric"

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

	// this check here doesn't make sense, but we will be moving all this logic to a function!
	if scraper.IsInfinite() {
		panic("can't have infinite scrapers with remote write")
	}

	if metric.HasInfiniteTimeSeries() {
		panic("can't have infinite time series with remote write")
	}

	iter := scraper.Iterator()
	for scrapeInfo, ok := iter.Next(); ok; scrapeInfo, ok = iter.Next() {
		metricResults := metric.Evaluate(scrapeInfo)

		// We have no more metrics to send
		if len(metricResults) == 0 {
			break
		}

		remoteWriterTimeSeries := promwrite.ConvertToRemoteWriterTimeSeries(metricName, metricResults)

		fmt.Printf("Time series: %+v\n", remoteWriterTimeSeries)

		err := prometheusRemoteWriter.Send(ctx, remoteWriterTimeSeries)
		if err != nil {
			panic(err)
		}
	}

	// TODO: Send Stale markers for metrics that didn't have them sent already!
	// For this we need to keep track of all the time series we've sent (the key should be all the labels serialized,
	// make sure the labels are sorted first)
	// Then we should remove them from the set as we see them being marked as stale.
	// Finally send the stale markers for the left overs.
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
