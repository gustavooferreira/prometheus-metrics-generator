package discrete

import (
	"github.com/gustavooferreira/prometheus-metrics-generator/metrics"
	"github.com/gustavooferreira/prometheus-metrics-generator/promadapter"
)

// Check at compile time whether MetricTimeSeries implements metrics.MetricTimeSeriesObservable interface.
var _ promadapter.MetricTimeSeriesObservable = (*MetricTimeSeries)(nil)

// MetricTimeSeries represents a metric time series (counter or gauge).
// When the time series iterator gets to the end of the DataGenerator provided it will evaluate the metrics.EndStrategy
// to decide on what to do next.
// The zero value of MetricTimeSeries is not useful. Use NewMetricTimeSeries function.
type MetricTimeSeries struct {
	labels map[string]string

	dataGenerator DataGenerator
	endStrategy   metrics.EndStrategy
}

// NewMetricTimeSeries creates a new instance of MetricTimeSeries.
func NewMetricTimeSeries(labels map[string]string, data DataGenerator, endStrategy metrics.EndStrategy) *MetricTimeSeries {
	return &MetricTimeSeries{
		labels:        labels,
		dataGenerator: data,
		endStrategy:   endStrategy,
	}
}

// Iterator returns a time series iterator that can be used to iterate over the data.
func (ts *MetricTimeSeries) Iterator() metrics.DataIterator {
	return &MetricTimeSeriesDataIterator{
		timeseries: *ts,
		state:      metrics.TimeSeriesIteratorStateRunning,
	}
}

// Labels returns the labels associated with the time series.
func (ts *MetricTimeSeries) Labels() map[string]string {
	return ts.labels
}

// IsInfinite reports whether this time series is infinite.
// In other words, whether this time series will never stop generating samples.
func (ts *MetricTimeSeries) IsInfinite() bool {
	return ts.endStrategy.EndStrategyType != metrics.EndStrategyTypeRemoveTimeSeries
}

// Check at compile time whether MetricTimeSeriesDataIterator implements metrics.DataIterator interface.
var _ metrics.DataIterator = (*MetricTimeSeriesDataIterator)(nil)

type MetricTimeSeriesDataIterator struct {
	timeseries MetricTimeSeries

	// currentDataIterator contains the DataIterator for the current run (in case we use a loop over strategy)
	currentDataIterator metrics.DataIterator

	// Reports whether we are evaluating data or we are in the end strategy stage
	state metrics.TimeSeriesIteratorState

	// loopCount keeps track of how many times we've looped over the DataGenerator
	// For example, if loopCount is 1, it means we've cycled through the iterator once and may be in the middle of
	// cycling through the iterator for a second time.
	loopCount int

	// lastValue represents the last value returned by the DataIterator
	lastValue metrics.ScrapeResult
}

// Evaluate fulfills the metrics.DataIterator interface.
// This function is responsible for returning the data points one at a time.
func (di *MetricTimeSeriesDataIterator) Evaluate(scrapeInfo metrics.ScrapeInfo) metrics.ScrapeResult {
	// Need the loop as when we reach the end of the iterator, regardless of what the end strategy is, we need to
	// evaluate the logic again, after setting the iterator state.
	for {
		if di.state == metrics.TimeSeriesIteratorStateEndStrategy {
			switch di.timeseries.endStrategy.EndStrategyType {
			case metrics.EndStrategyTypeLoop:
				di.currentDataIterator = nil
				di.loopCount++
				di.state = metrics.TimeSeriesIteratorStateRunning
				continue // it would be safe to let the logic run through as well
			case metrics.EndStrategyTypeSendLastValue:
				return di.lastValue
			case metrics.EndStrategyTypeSendCustomValue:
				return di.timeseries.endStrategy.CustomValue()
			case metrics.EndStrategyTypeRemoveTimeSeries:
				return metrics.ScrapeResult{Exhausted: true}
			default:
				// if the end strategy hasn't been set somehow, default to removing time series
				return metrics.ScrapeResult{Exhausted: true}
			}
		}

		// if we don't have an iterator, get one
		if di.currentDataIterator == nil {
			di.currentDataIterator = di.timeseries.dataGenerator.Iterator()
		}

		result := di.currentDataIterator.Evaluate(scrapeInfo)

		// We reached the end of the iterator
		if result.Exhausted {
			di.state = metrics.TimeSeriesIteratorStateEndStrategy
			continue
		}

		di.lastValue = result
		return result
	}
}
