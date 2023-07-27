package discrete

import (
	"fmt"

	"github.com/gustavooferreira/prometheus-metrics-generator/series"
)

// LinearSegmentOptions contains the options for the LinearSegment.
type LinearSegmentOptions struct {
	// AmplitudeStart represents the initial value for the segment.
	AmplitudeStart float64

	// AmplitudeEnd represents the end value for the segment.
	AmplitudeEnd float64

	// IterationCountLimit sets the number of iterations to be used by the segment.
	IterationCountLimit int
}

// Check at compile time whether LinearSegment implements DataGenerator interface.
var _ DataGenerator = (*LinearSegment)(nil)

// LinearSegment returns a DataGenerator representing a linear segment.
// A linear segment can be horizontal or have a positive or negative slope.
// Linear segments can be put together, with the help of the Join DataGenerator to form more complex structures.
// Note that it's an error to use a linear segment containing negative values with counters. It's the user's
// responsibility to make sure negative numbers only appear in gauges.
// The zero value is not useful.
type LinearSegment struct {
	options LinearSegmentOptions
}

// NewLinearSegment returns a new instance of LinearSegment.
func NewLinearSegment(options LinearSegmentOptions) (*LinearSegment, error) {
	if options.IterationCountLimit <= 0 {
		return &LinearSegment{}, fmt.Errorf("iteration count limit cannot be less than or equal to zero")
	}

	return &LinearSegment{
		options: options,
	}, nil
}

func (ls *LinearSegment) Iterator() DataIterator {
	slope := 0.0
	if ls.options.IterationCountLimit >= 2 {
		slope = (ls.options.AmplitudeEnd - ls.options.AmplitudeStart) / float64(ls.options.IterationCountLimit-1)
	}

	return &LinearSegmentIterator{
		linearSegment: *ls,
		slope:         slope,
	}
}

func (ls *LinearSegment) Describe() DataSpec {
	return DataNodeDataSpec{
		name: "Linear Segment",
	}
}

// Check at compile time whether LinearSegmentIterator implements DataIterator interface.
var _ DataIterator = (*LinearSegmentIterator)(nil)

type LinearSegmentIterator struct {
	// read-only access
	linearSegment LinearSegment
	slope         float64

	// These 2 variables keep track of the first scrape when running the iterator.
	// This allows us to keep track of how many iterations we've been running for.
	// All calculations are performed relative to the first detected scrape.
	firstScrapeHappened bool
	firstIterationCount int
}

// Iterate fulfills the DataIterator interface.
// This function is responsible for returning the data points one at a time.
func (lsi *LinearSegmentIterator) Iterate(scrapeInfo series.ScrapeInfo) series.ScrapeResult {
	// Is this the first scrape?
	if !lsi.firstScrapeHappened {
		lsi.firstScrapeHappened = true
		lsi.firstIterationCount = scrapeInfo.IterationCount
	}

	// Normalize
	currentIterationCount := scrapeInfo.IterationCount - lsi.firstIterationCount

	// Have we reached the end?
	if currentIterationCount >= lsi.linearSegment.options.IterationCountLimit {
		return series.ScrapeResult{Exhausted: true}
	}

	// If we have a horizontal line, there is no need to do any computation
	if lsi.linearSegment.options.AmplitudeStart == lsi.linearSegment.options.AmplitudeEnd {
		return series.ScrapeResult{Value: lsi.linearSegment.options.AmplitudeStart}
	}

	value := lsi.linearSegment.options.AmplitudeStart + lsi.slope*float64(currentIterationCount)
	return series.ScrapeResult{Value: value}
}
