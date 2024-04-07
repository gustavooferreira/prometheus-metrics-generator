package main

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gustavooferreira/prometheus-metrics-generator/discrete"
	"github.com/gustavooferreira/prometheus-metrics-generator/metrics"
	"github.com/gustavooferreira/prometheus-metrics-generator/promadapter"
)

// This is a CLI tool that reads metrics/timeseries from a yaml file.
// It can serve them with a http server exposing the metrics or write to prometheus using remote write.
func main() {
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
	err = metric.Attach(timeSeries)
	if err != nil {
		panic(err)
	}

	metric.Prepare()

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
