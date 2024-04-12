package continuous

import (
	"time"

	"github.com/gustavooferreira/prometheus-metrics-generator/metrics"
)

// startTime represents the time at which the function started being evaluated.
// A given iterator might not return any sample if too long has passed since the startTime.
type DataGenerator interface {
	Iterator(startTime time.Time) DataIterator
	// Describe() DataSpec
}

// DataIterator defines the interface a continuous data iterator needs to comply with.
type DataIterator interface {
	Evaluate(scrapeInfo ScrapeInfo) metrics.ScrapeResult
	Duration() time.Duration
}

// ScrapeInfo contains information about the scrape for a continuous data iterator.
type ScrapeInfo struct {
	// FirstIterationTime represents the time at which the very first iteration (scrape) happened.
	FirstIterationTime time.Time

	// IterationIndex specifies the index for this iteration.
	// An index of zero means this is the first iteration.
	IterationIndex int

	// IterationTime specifies the time of this iteration.
	IterationTime time.Time

	// FunctionStartTime represents the time at which the continuous function started.
	FunctionStartTime time.Time
}
