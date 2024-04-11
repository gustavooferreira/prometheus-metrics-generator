package promwrite

import "time"

// TimeSeries represents a time series that contains labels and a series of samples.
type TimeSeries struct {
	Labels  []Label
	Samples []Sample
}

// Label represents a label that can be attached to a time series.
type Label struct {
	Name  string
	Value string
}

// Sample represents a sample in a time series.
type Sample struct {
	Time  time.Time
	Value float64
}
