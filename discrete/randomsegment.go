package discrete

import (
	"fmt"
	"math/rand"

	"github.com/gustavooferreira/prometheus-metrics-generator/metrics"
)

// RandomSegmentDataGeneratorOptions contains the options for the RandomSegmentDataGenerator.
// The range for the random numbers is half-open [AmplitudeMin,AmplitudeMax[
type RandomSegmentDataGeneratorOptions struct {
	// AmplitudeMin represents the minimum value the data iterator will return.
	AmplitudeMin float64
	// AmplitudeMax represents the maximum value the data iterator will return (open-interval).
	AmplitudeMax float64

	// IterationCountLimit sets the number of iterations to be used by the segment.
	IterationCountLimit int
}

func (o *RandomSegmentDataGeneratorOptions) validate() error {
	if o.IterationCountLimit <= 0 {
		return fmt.Errorf("iteration count limit cannot be less than or equal to zero")
	}

	return nil
}

// Check at compile time whether RandomSegmentDataGenerator implements DataGenerator interface.
var _ DataGenerator = (*RandomSegmentDataGenerator)(nil)

// RandomSegmentDataGenerator returns a DataIterator representing a random sequence of samples.
// Note that it's an error to use negative values with counters. It's the user responsibility to make sure negative
// numbers only appear in gauges.
// The zero value is not useful.
type RandomSegmentDataGenerator struct {
	options RandomSegmentDataGeneratorOptions
}

// NewRandomDataGenerator returns a new instance of RandomSegmentDataGenerator.
func NewRandomDataGenerator(options RandomSegmentDataGeneratorOptions) (*RandomSegmentDataGenerator, error) {
	if err := options.validate(); err != nil {
		return &RandomSegmentDataGenerator{}, fmt.Errorf("error validating random segment data generator configuration: %w", err)
	}

	return &RandomSegmentDataGenerator{
		options: options,
	}, nil
}

func (dg *RandomSegmentDataGenerator) Iterator() metrics.DataIterator {
	return &RandomSegmentDataIterator{
		randomSegmentDataGenerator: *dg,
	}
}

func (dg *RandomSegmentDataGenerator) Describe() DataSpec {
	return DataNodeDataSpec{
		name: "Random",
	}
}

// Check at compile time whether RandomSegmentDataIterator implements metrics.DataIterator interface.
var _ metrics.DataIterator = (*RandomSegmentDataIterator)(nil)

type RandomSegmentDataIterator struct {
	// read-only access
	randomSegmentDataGenerator RandomSegmentDataGenerator

	// iterIndex keeps track of the current iteration.
	iterIndex int
}

// Evaluate fulfills the metrics.DataIterator interface.
// This function is responsible for returning the data points one at a time.
func (di *RandomSegmentDataIterator) Evaluate(scrapeInfo metrics.ScrapeInfo) metrics.ScrapeResult {
	// Have we reached the end?
	if di.iterIndex >= di.randomSegmentDataGenerator.options.IterationCountLimit {
		return metrics.ScrapeResult{Exhausted: true}
	}

	// Make sure to increment the iterator index before leaving the function
	defer func() { di.iterIndex++ }()

	randomRange := di.randomSegmentDataGenerator.options.AmplitudeMax - di.randomSegmentDataGenerator.options.AmplitudeMin
	randomValue := rand.Float64()*(randomRange) + di.randomSegmentDataGenerator.options.AmplitudeMin

	return metrics.ScrapeResult{Value: randomValue}
}
