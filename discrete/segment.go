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

func (o *LinearSegmentDataGeneratorOptions) validate() error {
	if o.IterationCountLimit <= 0 {
		return fmt.Errorf("iteration count limit cannot be less than or equal to zero")
	}

	return nil
}

// Check at compile time whether LinearSegmentDataGenerator implements DataGenerator interface.
var _ DataGenerator = (*LinearSegmentDataGenerator)(nil)

// LinearSegmentDataGenerator returns a DataGenerator representing a linear segment.
// A linear segment can be horizontal or have a positive or negative slope.
// Linear segments can be put together, with the help of the JoinDataGenerator to form more complex structures.
// Note that it's an error to use a linear segment containing negative values with counters. It's the user's
// responsibility to make sure negative numbers only appear in gauges.
// The LinearSegmentDataGenerator can only be used for metrics representing Counters or Gauges.
// The zero value is not useful.
type LinearSegmentDataGenerator struct {
	options LinearSegmentDataGeneratorOptions
	slope   float64
}

// NewLinearSegmentDataGenerator returns a new instance of LinearSegmentDataGenerator.
func NewLinearSegmentDataGenerator(options LinearSegmentDataGeneratorOptions) (*LinearSegmentDataGenerator, error) {
	if err := options.validate(); err != nil {
		return &LinearSegmentDataGenerator{}, fmt.Errorf("error validating linear segment data generator configuration: %w", err)
	}

	slope := 0.0
	if options.IterationCountLimit >= 2 {
		slope = (options.AmplitudeEnd - options.AmplitudeStart) / float64(options.IterationCountLimit-1)
	}

	return &LinearSegmentDataGenerator{
		options: options,
		slope:   slope,
	}, nil
}

func (ls *LinearSegmentDataGenerator) Iterator() metrics.DataIterator {
	return &LinearSegmentDataIterator{
		linearSegmentDataGenerator: *ls,
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

	// iterIndex keeps track of the current iteration.
	iterIndex int
}

// Evaluate fulfills the metrics.DataIterator interface.
// This function is responsible for returning the data points one at a time.
func (lsi *LinearSegmentDataIterator) Evaluate(scrapeInfo metrics.ScrapeInfo) metrics.ScrapeResult {
	// Have we reached the end?
	if lsi.iterIndex >= lsi.linearSegmentDataGenerator.options.IterationCountLimit {
		return metrics.ScrapeResult{Exhausted: true}
	}

	// Make sure to increment the iterator index before leaving the function
	defer func() { lsi.iterIndex++ }()

	// If we have a horizontal line, there is no need to do any computation
	if lsi.linearSegmentDataGenerator.options.AmplitudeStart == lsi.linearSegmentDataGenerator.options.AmplitudeEnd {
		return metrics.ScrapeResult{Value: lsi.linearSegmentDataGenerator.options.AmplitudeStart}
	}

	value := lsi.linearSegmentDataGenerator.options.AmplitudeStart + lsi.linearSegmentDataGenerator.slope*float64(lsi.iterIndex)
	return metrics.ScrapeResult{Value: value}
}
