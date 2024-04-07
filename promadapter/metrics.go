package promadapter

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/gustavooferreira/prometheus-metrics-generator/metrics"
)

// MetricObservable defines the interface metrics should implement.
type MetricObservable interface {
	Desc() Desc
	PromDescs() []*prometheus.Desc
	Evaluate(scrapeInfo metrics.ScrapeInfo) []MetricResult
}

// MetricTimeSeriesObservable is the interface implemented by any time series wanting to be scraped.
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
type Metric struct {
	// desc represents the descriptor that describes this metric
	desc Desc

	// timeSeries contains all the time series attached to this metric
	timeSeries []MetricTimeSeriesObservable

	// timeSeriesIterators contains the iterators for all time series
	timeSeriesIterators []metrics.DataIterator

	// timeSeriesDesc contains the prometheus.Desc for all time series
	timeSeriesDesc []*prometheus.Desc
}

func NewMetric(fqName string, help string, metricType MetricType, labelsNames []string) *Metric {
	return &Metric{
		desc: Desc{
			FQName:      fqName,
			Help:        help,
			MetricType:  metricType,
			LabelsNames: labelsNames,
		},
	}
}

// Attach attaches a time series (counter or gauge) to the metric.
func (m *Metric) Attach(metricTimeSeries MetricTimeSeriesObservable) error {
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
	return nil
}

// Prepare gets the iterators and descriptors for all time series.
// This function needs to be called before being passed to a collector.
func (m *Metric) Prepare() {
	for _, timeSeries := range m.timeSeries {
		m.timeSeriesIterators = append(m.timeSeriesIterators, timeSeries.Iterator())

		desc := prometheus.NewDesc(
			m.desc.FQName,
			m.desc.Help,
			m.desc.LabelsNames,
			nil,
		)

		m.timeSeriesDesc = append(m.timeSeriesDesc, desc)
	}
}

func (m *Metric) Desc() Desc {
	return m.desc
}

func (m *Metric) PromDescs() []*prometheus.Desc {
	return m.timeSeriesDesc
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
			PromDesc:     m.timeSeriesDesc[i],
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
