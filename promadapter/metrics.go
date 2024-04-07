package promadapter

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/gustavooferreira/prometheus-metrics-generator/metrics"
)

// MetricObservable defines the interface metrics should implement.
// This only applies to Counter and Gauge metrics.
type MetricObservable interface {
	Desc() Desc
	PromDesc() *prometheus.Desc
	Evaluate(scrapeInfo metrics.ScrapeInfo) []MetricResult
	TimeSeriesCount() int
}

// MetricTimeSeriesObservable is the interface implemented by any time series wanting to be scraped.
// This is only valid for Counter and Gauge metrics.
type MetricTimeSeriesObservable interface {
	Iterator() metrics.DataIterator
	Labels() map[string]string
}

// Desc represents the description of the metric.
type Desc struct {
	// FQName represents the name of the metric.
	FQName string

	// Help represent the Help string of the metric.
	Help string

	// MetricType represents the type of the metric (counter or gauge).
	MetricType MetricType

	// LabelsNames contains the names of the labels to be use by the time series attached to this metric
	LabelsNames []string
}

type MetricType string

const (
	MetricTypeCounter   MetricType = "time_series_type-counter"
	MetricTypeGauge     MetricType = "time_series_type-gauge"
	MetricTypeHistogram MetricType = "time_series_type-histogram"
)

// Check at compile time whether Metric implements MetricObservable interface.
var _ MetricObservable = (*Metric)(nil)

// Metric represents a metric.
// It's only meant to be used by metrics that are Counters or Gauges.
// The zero value is not useful. Use the NewMetric function instead.
type Metric struct {
	// desc represents the descriptor that describes this metric.
	desc Desc

	// promDesc contains the prometheus.Desc for the metric.
	promDesc *prometheus.Desc

	// timeSeries contains all the time series attached to this metric.
	timeSeries []MetricTimeSeriesObservable

	// timeSeriesIterators contains the iterators for all time series.
	timeSeriesIterators []metrics.DataIterator

	// timeSeriesCount is the number of time series contained in this metric.
	timeSeriesCount int
}

// NewMetric creates a new instance of Metric.
// It's only meant to be used by metrics that are Counters or Gauges.
func NewMetric(fqName string, help string, metricType MetricType, labelsNames []string) *Metric {
	desc := Desc{
		FQName:      fqName,
		Help:        help,
		MetricType:  metricType,
		LabelsNames: labelsNames,
	}
	promDesc := prometheus.NewDesc(fqName, help, labelsNames, nil)

	return &Metric{
		desc:     desc,
		promDesc: promDesc,
	}
}

// AddTimeSeries adds a time series (counter or gauge) to the metric.
func (m *Metric) AddTimeSeries(metricTimeSeries MetricTimeSeriesObservable) error {
	// validate labels match
	labelsNamesMap := make(map[string]struct{})
	for _, labelName := range m.desc.LabelsNames {
		labelsNamesMap[labelName] = struct{}{}
	}

	metricsLabels := metricTimeSeries.Labels()

	for k := range metricsLabels {
		// time series includes an unexpected label
		if _, ok := labelsNamesMap[k]; !ok {
			return fmt.Errorf("label mismatch: unexpected label in time series")
		}
		delete(labelsNamesMap, k)
	}

	// time series doesn't include an expected label
	if len(labelsNamesMap) != 0 {
		return fmt.Errorf("label mismatch: missing expected label in time series")
	}

	m.timeSeries = append(m.timeSeries, metricTimeSeries)
	m.timeSeriesIterators = append(m.timeSeriesIterators, metricTimeSeries.Iterator())
	m.timeSeriesCount++

	return nil
}

func (m *Metric) Desc() Desc {
	return m.desc
}

func (m *Metric) PromDesc() *prometheus.Desc {
	return m.promDesc
}

func (m *Metric) TimeSeriesCount() int {
	return m.timeSeriesCount
}

// Evaluate returns the computed samples as well as the labels for the time series.
// It returns an array as the Metric may have multiple time series attached.
// If the sample for a given time series is missing or the time series itself has been exhausted, then the result
// won't be included in the returned array.
func (m *Metric) Evaluate(scrapeInfo metrics.ScrapeInfo) []MetricResult {
	var results []MetricResult

	// loop over iterators, get result and decide what to do
	for i, timeSeriesIterator := range m.timeSeriesIterators {
		scrapeResult := timeSeriesIterator.Evaluate(scrapeInfo)

		// We do not send a given time series result if the sample has been flagged as missing or the time series has
		// exhausted.
		if scrapeResult.Exhausted || scrapeResult.Missing {
			continue
		}

		result := MetricResult{
			Desc:         m.desc,
			PromDesc:     m.promDesc,
			LabelsValues: m.timeSeries[i].Labels(),
			Value:        scrapeResult.Value,
		}

		results = append(results, result)
	}

	return results
}

// MetricResult represents the result of a Metric.
// It's only meant to be used by a Counter or a Gauge.
type MetricResult struct {
	Desc     Desc
	PromDesc *prometheus.Desc

	LabelsValues map[string]string

	// Value represents the value of the sample
	Value float64
}
