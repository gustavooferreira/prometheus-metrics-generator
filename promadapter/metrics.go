package promadapter

import (
	"fmt"
	"time"

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
	IsInfinite() bool
}

// Desc represents the description of the metric.
type Desc struct {
	// FQName represents the name of the metric (also known as Metric Family).
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

	// timeSeriesStaleMarkers contains the state for stale markers for all time series.
	timeSeriesStaleMarkers []bool

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
	m.timeSeriesStaleMarkers = append(m.timeSeriesStaleMarkers, false)
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

// HasInfiniteTimeSeries checks whether any of the time series in this metric family is infinite.
func (m *Metric) HasInfiniteTimeSeries() bool {
	for _, singleTimeseries := range m.timeSeries {
		if singleTimeseries.IsInfinite() {
			return true
		}
	}

	return false
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

		// We do not send a given time series result if the sample has been flagged as missing.
		if scrapeResult.Missing {
			continue
		}

		// We do not send a given time series result if the time series has exhausted and the stale marker has already
		// been sent.
		if scrapeResult.Exhausted {
			if m.timeSeriesStaleMarkers[i] {
				continue
			}

			m.timeSeriesStaleMarkers[i] = true
		}

		result := MetricResult{
			Desc:        m.desc,
			PromDesc:    m.promDesc,
			LabelsSet:   m.timeSeries[i].Labels(),
			Timestamp:   scrapeInfo.IterationTime,
			Value:       scrapeResult.Value,
			StaleMarker: m.timeSeriesStaleMarkers[i],
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

	// LabelsSet is the set of labels associated with this sample.
	LabelsSet map[string]string

	// Timestamp represents the timestamp of the sample.
	Timestamp time.Time

	// Value represents the value of the sample.
	Value float64

	// StaleMarker represents whether this time series has come to an end.
	// Spec Ref:
	//	Prometheus remote write compatible senders MUST send stale markers when a time series will no longer be appended
	//  to.
	StaleMarker bool
}
