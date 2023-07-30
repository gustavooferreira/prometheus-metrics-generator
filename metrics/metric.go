package metrics

import (
	"fmt"
)

// MetricObservable defines the interface metrics should implement.
type MetricObservable interface {
	Desc() Desc
	Evaluate(scrapeInfo ScrapeInfo) []MetricResult
}

// MetricTimeSeriesObservable is the interface implemented by any time series wanting to be scraped.
type MetricTimeSeriesObservable interface {
	Iterator() DataIterator
	Labels() map[string]string
}

type Desc struct {
	// FQName represents the name of the metric
	FQName string

	// Help represent the Help string of the metric
	Help string

	// MetricType represents the type of the metric
	MetricType MetricType

	// LabelsNames contains the names of the labels to be use by the time series attached to this metric
	LabelsNames []string
}

type MetricType string

const (
	MetricTypeCounter MetricType = "time_series_type-counter"
	MetricTypeGauge   MetricType = "time_series_type-gauge"
)

type Metric struct {
	// desc represents the descriptor that describes this metric
	desc Desc

	// timeSeries contains all the time series attached to this metric
	timeSeries []MetricTimeSeriesObservable

	// timeSeriesIterators contains the iterators for all time series
	timeSeriesIterators []DataIterator
}

func NewMetric(fqName string, help string, metricType MetricType, labelsNames []string) Metric {
	return Metric{
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

	for k, _ := range metricsLabels {
		if _, ok := labelsNamesMap[k]; !ok {
			return fmt.Errorf("label mismatch")
		}
		delete(labelsNamesMap, k)
	}

	if len(labelsNamesMap) != 0 {
		return fmt.Errorf("label mismatch")
	}

	m.timeSeries = append(m.timeSeries, metricTimeSeries)
	return nil
}

// Prepare gets the iterators for all time series.
// This function needs to be called before we can start getting values for each time series.
func (m *Metric) Prepare() {
	for _, timeSeries := range m.timeSeries {
		m.timeSeriesIterators = append(m.timeSeriesIterators, timeSeries.Iterator())
	}
}

func (m *Metric) Desc() Desc {
	return m.desc
}

// Evaluate returns the scrapeResult as well as the labels for the time series.
// It's an array because we may have more than one time series attached to this metric.
// If the scrapeResult has Exhausted turned on, don't include the result for that time series.
func (m *Metric) Evaluate(scrapeInfo ScrapeInfo) []MetricResult {
	// loop over iterators, get result and decide what to do
	var results []MetricResult

	for i, timeSeriesIterator := range m.timeSeriesIterators {
		scrapeResult := timeSeriesIterator.Evaluate(scrapeInfo)

		// We do not send metrics if the sample has been flagged as missing or the time series has exhausted
		if scrapeResult.Exhausted || scrapeResult.Missing {
			continue
		}

		result := MetricResult{
			Desc:         m.desc,
			LabelsValues: m.timeSeries[i].Labels(),
			Result:       scrapeResult.Value,
		}

		results = append(results, result)
	}

	return results
}

type MetricResult struct {
	Desc Desc

	LabelsValues map[string]string

	// NOTE: We don't need to include the missing or exhausted fields because the Evaluate() method already decides
	// to not return a value if one of these conditions is met.
	Result float64
}
