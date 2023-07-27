package series

import "time"

// DataIterator defines the function type returned by all data functions in the datafuncs package.
// This is the primary way of creating series.
// The function is an iterator and therefore each time it's called returns the next value in the series.
// The field Exhausted in the ScrapeResult struct reports whether there is no more data to be returned by the iterator.
// Counter resets can be simulated by setting the Value field to zero.
// Missing scrapes can be simulated by setting the Missing field to true.
// TODO: Nope, let's define an interface insteas!
// TODO: Let's define a DataIteratorFunc just line the HandlerFunc!
type DataIterator2 func(scrapeInfo ScrapeInfo) ScrapeResult

type DataIterator interface {
	Iterate(scrapeInfo ScrapeInfo) ScrapeResult
}

// ScrapeInfo contains information about the scrape.
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
