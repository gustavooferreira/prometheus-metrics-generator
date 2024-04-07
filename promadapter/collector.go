package promadapter

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/gustavooferreira/prometheus-metrics-generator/metrics"
)

// Check at compile time whether Collector implements prometheus.Collector interface.
var _ prometheus.Collector = (*Collector)(nil)

// Collector implements the prometheus.Collector interface.
// It needs to be registered with a prometheus.Registerer in order for prometheus to be able to scrape metrics.
type Collector struct {
	// metricObservable is a list of metrics we should scrape.
	// Each metric may have multiple time series!
	// This field acts as a read-only variable, once set in the constructor, it's never changed.
	metricObservables []MetricObservable

	// mu protects the fields below
	mu sync.Mutex

	// firstIterationTime represents the time at which the very first iteration (scrape) happened.
	firstIterationTime time.Time

	// iterIndex keeps track of the current iteration.
	iterIndex int
}

// NewCollector returns a new collector to be registered with the prometheus.Registerer.
func NewCollector(metrics []MetricObservable) *Collector {
	return &Collector{
		metricObservables: metrics,
	}
}

// Describe is part of the implementation of the promentheus.Collector interface.
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, metricObservable := range c.metricObservables {
		for i := 0; i < metricObservable.TimeSeriesCount(); i++ {
			ch <- metricObservable.PromDesc()
		}
	}
}

// Collect runs the logic to collect the metrics.
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	now := time.Now().UTC()

	c.mu.Lock()
	// is this the first iteration?
	if c.iterIndex == 0 {
		c.firstIterationTime = now
	}

	scrapeInfo := metrics.ScrapeInfo{
		FirstIterationTime: c.firstIterationTime,
		IterationIndex:     c.iterIndex,
		IterationTime:      now,
	}

	// Make sure to increment the iterator index before leaving the function
	defer func() { c.iterIndex++ }()

	c.mu.Unlock()

	for _, metricObservable := range c.metricObservables {
		metricResults := metricObservable.Evaluate(scrapeInfo)

		for _, metricResult := range metricResults {
			var metricType prometheus.ValueType
			switch metricResult.Desc.MetricType {
			case MetricTypeCounter:
				metricType = prometheus.CounterValue
			case MetricTypeGauge:
				metricType = prometheus.GaugeValue
			default:
				// for now we ignore the error!
				continue
			}

			// Create array of label values in the same order the label names were specified!
			var labelValues []string

			for _, labelName := range metricObservable.Desc().LabelsNames {
				labelValues = append(labelValues, metricResult.LabelsValues[labelName])
			}

			metric, err := prometheus.NewConstMetric(
				metricResult.PromDesc,
				metricType,
				metricResult.Value,
				labelValues...,
			)
			if err != nil {
				// for now we ignore the error
				continue
			}

			ch <- metric
		}
	}
}
