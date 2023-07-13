package continuous

import (
	"time"

	"github.com/gustavooferreira/prometheus-metrics-generator/series"
)

// DataIterator defines the interface a continuous data iterator needs to comply with.
type DataIterator interface {
	Evaluate(scrapeInfo ScrapeInfo) series.ScrapeResult
	Duration() time.Duration
}

// ScrapeInfo contains information about the scrape for a continuous data iterator.
type ScrapeInfo struct {
	// FirstIterationTime represents the time at which the very first iteration (scrape) happened.
	FirstIterationTime time.Time

	// IterationCount specifies the count for this iteration.
	// A count of zero means this is the first iteration.
	IterationCount int

	// IterationTime specifies the time of this iteration.
	IterationTime time.Time

	// FunctionStartTime represents the time at which the continuous function started.
	FunctionStartTime time.Time
}
