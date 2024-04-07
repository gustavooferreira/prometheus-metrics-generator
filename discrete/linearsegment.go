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
// Note that it's an error to use negative values with counters. It's the user responsibility to make sure negative
// numbers only appear in gauges.
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

func (dg *LinearSegmentDataGenerator) Iterator() metrics.DataIterator {
	return &LinearSegmentDataIterator{
		linearSegmentDataGenerator: *dg,
	}
}

func (dg *LinearSegmentDataGenerator) Describe() DataSpec {
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
func (di *LinearSegmentDataIterator) Evaluate(scrapeInfo metrics.ScrapeInfo) metrics.ScrapeResult {
	// Have we reached the end?
	if di.iterIndex >= di.linearSegmentDataGenerator.options.IterationCountLimit {
		return metrics.ScrapeResult{Exhausted: true}
	}

	// Make sure to increment the iterator index before leaving the function
	defer func() { di.iterIndex++ }()

	// If we have a horizontal line, there is no need to do any computation
	if di.linearSegmentDataGenerator.options.AmplitudeStart == di.linearSegmentDataGenerator.options.AmplitudeEnd {
		return metrics.ScrapeResult{Value: di.linearSegmentDataGenerator.options.AmplitudeStart}
	}

	value := di.linearSegmentDataGenerator.options.AmplitudeStart + di.linearSegmentDataGenerator.slope*float64(di.iterIndex)
	return metrics.ScrapeResult{Value: value}
}
