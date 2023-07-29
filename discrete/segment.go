package discrete

import (
	"fmt"

	"github.com/gustavooferreira/prometheus-metrics-generator/metrics"
)

// LinearSegmentDataGeneratorOptions contains the options for the LinearSegmentDataGenerator.
type LinearSegmentDataGeneratorOptions struct {
	// AmplitudeStart represents the initial value for the segment.
	AmplitudeStart float64

	// AmplitudeEnd represents the end value for the segment.
	AmplitudeEnd float64

	// IterationCountLimit sets the number of iterations to be used by the segment.
	IterationCountLimit int
}

// Check at compile time whether LinearSegmentDataGenerator implements DataGenerator interface.
var _ DataGenerator = (*LinearSegmentDataGenerator)(nil)

// LinearSegmentDataGenerator returns a DataGenerator representing a linear segment.
// A linear segment can be horizontal or have a positive or negative slope.
// Linear segments can be put together, with the help of the NewJoinDataGenerator DataGenerator to form more complex structures.
// Note that it's an error to use a linear segment containing negative values with counters. It's the user's
// responsibility to make sure negative numbers only appear in gauges.
// The zero value is not useful.
type LinearSegmentDataGenerator struct {
	options LinearSegmentDataGeneratorOptions
}

// NewLinearSegmentDataGenerator returns a new instance of LinearSegmentDataGenerator.
func NewLinearSegmentDataGenerator(options LinearSegmentDataGeneratorOptions) (*LinearSegmentDataGenerator, error) {
	if options.IterationCountLimit <= 0 {
		return &LinearSegmentDataGenerator{}, fmt.Errorf("iteration count limit cannot be less than or equal to zero")
	}

	return &LinearSegmentDataGenerator{
		options: options,
	}, nil
}

func (ls *LinearSegmentDataGenerator) Iterator() metrics.DataIterator {
	slope := 0.0
	if ls.options.IterationCountLimit >= 2 {
		slope = (ls.options.AmplitudeEnd - ls.options.AmplitudeStart) / float64(ls.options.IterationCountLimit-1)
	}

	return &LinearSegmentDataIterator{
		linearSegmentDataGenerator: *ls,
		slope:                      slope,
	}
}

func (ls *LinearSegmentDataGenerator) Describe() DataSpec {
	return DataNodeDataSpec{
		name: "Linear Segment",
	}
}

// Check at compile time whether LinearSegmentDataIterator implements metrics.DataIterator interface.
var _ metrics.DataIterator = (*LinearSegmentDataIterator)(nil)

type LinearSegmentDataIterator struct {
	// read-only access
	linearSegmentDataGenerator LinearSegmentDataGenerator
	slope                      float64

	// iterCount represents the cycle the iterator is in
	iterCount int
}

// Evaluate fulfills the metrics.DataIterator interface.
// This function is responsible for returning the data points one at a time.
func (lsi *LinearSegmentDataIterator) Evaluate(scrapeInfo metrics.ScrapeInfo) metrics.ScrapeResult {
	// Have we reached the end?
	if lsi.iterCount >= lsi.linearSegmentDataGenerator.options.IterationCountLimit {
		return metrics.ScrapeResult{Exhausted: true}
	}

	// Make sure to increment the iterator counter before leaving the function
	defer func() { lsi.iterCount++ }()

	// If we have a horizontal line, there is no need to do any computation
	if lsi.linearSegmentDataGenerator.options.AmplitudeStart == lsi.linearSegmentDataGenerator.options.AmplitudeEnd {
		return metrics.ScrapeResult{Value: lsi.linearSegmentDataGenerator.options.AmplitudeStart}
	}

	value := lsi.linearSegmentDataGenerator.options.AmplitudeStart + lsi.slope*float64(lsi.iterCount)
	return metrics.ScrapeResult{Value: value}
}
