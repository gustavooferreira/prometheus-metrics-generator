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

	// iterationCount specifies the count for this iteration.
	// A count of zero means this is the first iteration.
	iterationCount int
}

// NewCollector returns a new collector to be registered with the prometheus.Registerer.
// The metrics passed to the collector need to have been prepared already, meaning the Prepare() method on the Metric
// struct must be called before the metrics can be passed to the Collector.
// Failing to do so will result in no metrics being returned by the collector.
func NewCollector(metrics []MetricObservable) *Collector {
	return &Collector{
		metricObservables: metrics,
	}
}

// Describe is part of the implementation of the promentheus.Collector interface.
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, metricObservable := range c.metricObservables {
		for _, promDesc := range metricObservable.PromDescs() {
			ch <- promDesc
		}
	}
}

// Collect runs the logic to collect the metrics.
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	now := time.Now().UTC()

	c.mu.Lock()
	// is this the first iteration?
	if c.iterationCount == 0 {
		c.firstIterationTime = now
	}

	scrapeInfo := metrics.ScrapeInfo{
		FirstIterationTime: c.firstIterationTime,
		IterationIndex:     c.iterationCount,
		IterationTime:      now,
	}

	// increment iteration count here as we won't use it anymore in this function
	c.iterationCount++
	c.mu.Unlock()

	for _, metricObservable := range c.metricObservables {
		metricsResults := metricObservable.Evaluate(scrapeInfo)

		for _, metricResult := range metricsResults {
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
