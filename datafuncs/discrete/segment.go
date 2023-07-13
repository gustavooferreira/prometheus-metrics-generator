package discrete

import (
	"fmt"

	"github.com/gustavooferreira/prometheus-metrics-generator/series"
)

// LinearSegmentDataIterator returns a DataIterator representing a linear segment.
// A linear segment can be horizontal or have a positive or negative slope.
// Linear segments can be put together, with the help of the Join function, to form more complex structures.
// Note that it's an error to use a linear segment containing negative values with counters. It's the user
// responsibility to make sure negative numbers only appear in gauges.
// The data series generated by this function must be finite. A perpetual series can be created by setting the
// EndStrategy in the series.Series struct to loop over forever.
// The zero value is not useful.
type LinearSegmentDataIterator struct {
	options LinearSegmentDataIteratorOptions

	// These 2 variables keep track of the first scrape when running the iterator.
	// This allows us to keep track of how many iterations we've been running for.
	// All calculations are performed relative to the first detected scrape.
	firstScrapeHappened bool
	firstIterationCount int
}

// NewLinearSegmentDataIterator returns a new instance of LinearSegmentDataIterator.
func NewLinearSegmentDataIterator(options LinearSegmentDataIteratorOptions) (LinearSegmentDataIterator, error) {
	if options.IterationCountLimit <= 0 {
		return LinearSegmentDataIterator{}, fmt.Errorf("iteration count limit cannot be less than or equal to zero")
	}

	return LinearSegmentDataIterator{
		options: options,
	}, nil
}

// Evaluate fulfills the DataIterator function type.
// This function is responsible for returning the data points one at a time.
func (ls *LinearSegmentDataIterator) Evaluate(scrapeInfo series.ScrapeInfo) series.ScrapeResult {
	// Is this the first scrape?
	if !ls.firstScrapeHappened {
		ls.firstScrapeHappened = true
		ls.firstIterationCount = scrapeInfo.IterationCount
	}

	// Normalize
	currentIterationCount := scrapeInfo.IterationCount - ls.firstIterationCount

	// Have we reached the end?
	if currentIterationCount >= ls.options.IterationCountLimit {
		return series.ScrapeResult{Exhausted: true}
	}

	// If we have a horizontal line, there is no need to do any computation
	if ls.options.AmplitudeStart == ls.options.AmplitudeEnd {
		return series.ScrapeResult{Value: ls.options.AmplitudeStart}
	}

	slope := (ls.options.AmplitudeEnd - ls.options.AmplitudeStart) / float64(ls.options.IterationCountLimit-1)
	value := ls.options.AmplitudeStart + slope*float64(currentIterationCount)
	return series.ScrapeResult{Value: value}
}

// LinearSegmentDataIteratorOptions contains the options for the LinearSegmentDataIterator.
type LinearSegmentDataIteratorOptions struct {
	// AmplitudeStart represents the initial value for the segment.
	AmplitudeStart float64

	// AmplitudeEnd represents the end value for the segment.
	AmplitudeEnd float64

	// IterationCountLimit sets the number of iterations to be used by the segment.
	IterationCountLimit int
}
