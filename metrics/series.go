package metrics

type TimeSeriesType string

const (
	TimeSeriesTypeCounter TimeSeriesType = "time_series_type-counter"
	TimeSeriesTypeGauge   TimeSeriesType = "time_series_type-gauge"
)

type TimeSeriesIteratorState string

const (
	TimeSeriesIteratorStateRunning     TimeSeriesIteratorState = "time_series_iterator_state-running"
	TimeSeriesIteratorStateEndStrategy TimeSeriesIteratorState = "time_series_iterator_state-end_strategy"
)

// TimeSeriesInfo contains information about the time series.
type TimeSeriesInfo struct {
	// Type specifies the type of series (counter, gauge, etc)
	Type TimeSeriesType
	// Name specifies the series name
	Name string
	// Labels represents the labels part of the time series
	Labels map[string]string
}

// TimeSeries defines the interface time series must implement.
type TimeSeries interface {
	Iterator() DataIterator
	Info() TimeSeriesInfo
}
