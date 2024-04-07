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

// ScrapeHandler defines the function type to be used when calling the ScrapeDataIterator method of the Scraper.
// This function is specifically used with counter and gauge metrics.
// Return an error to stop the scrapping from proceeding any further. The ScrapeDataIterator method will wrap the
// returned error and return it.
// The field 'Exhausted' in the struct ScrapeResult will never be set, since if the scraper has exhausted the
// DataIterator it will automatically stop.
type ScrapeHandler func(scrapeInfo ScrapeInfo, scrapeResult ScrapeResult) error

// DataHistogramIterator defines the interface iterators returned by DataGenerators need to comply with.
// The DataIterator interface will be used regardless of whether the time series is being consumed by a prometheus
// collector or by the logic responsible for the prometheus remote write.
// Both continuous and discrete time series will return iterators that comply with this interface.
// The interface defines an iterator and therefore each time it's called, it returns the next value in the series.
// The field Exhausted in the ScrapeResult struct reports whether there is no more data to be returned by the iterator.
// Counter resets can be simulated by setting the Value field to zero.
// Missing scrapes can be simulated by setting the Missing field to true.
type DataHistogramIterator interface {
	Evaluate(scrapeInfo ScrapeInfo) ScrapeHistogramResult
}

// The DataHistogramIteratorFunc type is an adapter to allow the use of ordinary functions as DataHistogramIterator.
// If f is a function with the appropriate signature, DataHistogramIteratorFunc(f) is a DataHistogramIterator that
// calls f.
type DataHistogramIteratorFunc func(scrapeInfo ScrapeInfo) ScrapeHistogramResult

// Evaluate calls f(scrapeInfo).
func (f DataHistogramIteratorFunc) Evaluate(scrapeInfo ScrapeInfo) ScrapeHistogramResult {
	return f(scrapeInfo)
}

// ScrapeHistogramHandler defines the function type to be used when calling the ScrapeDataIterator method of the Scraper.
// This function is specifically used with histogram metrics.
// Return an error to stop the scrapping from proceeding any further. The ScrapeDataIterator method will wrap the
// returned error and return it.
// The field 'Exhausted' in the struct ScrapeResult will never be set, since if the scraper has exhausted the
// DataIterator it will automatically stop.
type ScrapeHistogramHandler func(scrapeInfo ScrapeInfo, scrapeHistogramResult ScrapeHistogramResult) error

// ScrapeInfo contains information about the scrape.
// Namely, information of when the scrape is happening.
type ScrapeInfo struct {
	// FirstIterationTime represents the time at which the very first iteration (scrape) happened.
	FirstIterationTime time.Time

	// IterationIndex specifies the index for this iteration.
	// An index of zero means this is the first iteration.
	IterationIndex int

	// IterationTime specifies the time of this iteration.
	IterationTime time.Time
}

// ScrapeResult contains the scrape outcome.
// Used for counters and gauges.
type ScrapeResult struct {
	// Value is the value of the sample.
	Value float64

	// Missing indicates whether the scrape failed to retrieve a sample.
	// Used to simulate failed scrapes.
	Missing bool

	// Exhausted indicates whether the data generator has no more data to return.
	Exhausted bool
}

// ScrapeHistogramResult contains the scrape outcome.
// Used for histograms.
type ScrapeHistogramResult struct {
	// Buckets represents the buckets of the histogram.
	Buckets []HistogramBucketScrape
	// Count is the number of occurances recorded by this histogram.
	Count float64
	// Sum is the total sum of all the values in the histogram.
	Sum float64

	// Missing indicates whether the scrape failed to retrieve a sample.
	// Used to simulate failed scrapes.
	Missing bool

	// Exhausted indicates whether the data generator has no more data to return.
	Exhausted bool
}

// HistogramBucketScrape represents a single scrape for a single histogram bucket.
type HistogramBucketScrape struct {
	// LE represents the less than or equal to threshold of the bucket
	LE float64

	// Value is the value of the sample.
	Value float64
}
