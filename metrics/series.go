package metrics

import (
	"time"
)

type TimeSeriesIteratorState string

const (
	TimeSeriesIteratorStateRunning     TimeSeriesIteratorState = "time_series_iterator_state-running"
	TimeSeriesIteratorStateEndStrategy TimeSeriesIteratorState = "time_series_iterator_state-end_strategy"
)

// DataIterator defines the interface iterators returned by DataGenerators need to comply with.
// The DataIterator interface will be used regardless of whether the time series is being consumed by a prometheus
// collector or by the logic responsible for the prometheus remote write.
// Both continuous and discrete time series will return iterators that comply with this interface.
// The interface defines an iterator and therefore each time it's called, it returns the next value in the series.
// The field Exhausted in the ScrapeResult struct reports whether there is no more data to be returned by the iterator.
// Counter resets can be simulated by setting the Value field to zero.
// Missing scrapes can be simulated by setting the Missing field to true.
type DataIterator interface {
	Evaluate(scrapeInfo ScrapeInfo) ScrapeResult
}

// The DataIteratorFunc type is an adapter to allow the use of ordinary functions as DataIterator. If f is a function
// with the appropriate signature, DataIteratorFunc(f) is a DataIterator that calls f.
type DataIteratorFunc func(scrapeInfo ScrapeInfo) ScrapeResult

// Evaluate calls f(scrapeInfo).
func (f DataIteratorFunc) Evaluate(scrapeInfo ScrapeInfo) ScrapeResult {
	return f(scrapeInfo)
}

// ScrapeInfo contains information about the scrape.
// Namely, information of when the scrape is happening.
type ScrapeInfo struct {
	// FirstIterationTime represents the time at which the very first iteration (scrape) happened.
	FirstIterationTime time.Time

	// IterationCount specifies the count for this iteration.
	// A count of zero means this is the first iteration.
	IterationCount int

	// IterationTime specifies the time of this iteration.
	IterationTime time.Time
}

// ScrapeResult contains the scrape outcome.
type ScrapeResult struct {
	// Value is the value of the sample.
	Value float64

	// Missing indicates whether the scrape failed to retrieve a sample.
	// Used to simulate failed scrapes.
	Missing bool

	// Exhausted indicates whether the data generator has no more data to return.
	Exhausted bool
}
