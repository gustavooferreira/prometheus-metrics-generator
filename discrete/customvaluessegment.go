package discrete

import (
	"github.com/gustavooferreira/prometheus-metrics-generator/metrics"
)

// Check at compile time whether CustomValuesSegmentDataGenerator implements DataGenerator interface.
var _ DataGenerator = (*CustomValuesSegmentDataGenerator)(nil)

// CustomValuesSegmentDataGenerator returns a DataGenerator containing the array of values passed in.
// Each value is returned in sequence on each scrape.
// Note that it's an error to use negative values with counters. It's the user responsibility to make sure negative
// numbers only appear in gauges.
// The zero value is not useful.
type CustomValuesSegmentDataGenerator struct {
	values []CustomValueSample
}

// NewCustomValuesDataGenerator returns an instance of CustomValuesDataGenerator.
func NewCustomValuesDataGenerator(values []CustomValueSample) *CustomValuesSegmentDataGenerator {
	return &CustomValuesSegmentDataGenerator{
		values: values,
	}
}

func (dg *CustomValuesSegmentDataGenerator) Iterator() metrics.DataIterator {
	return &CustomValuesSegmentDataIterator{
		customValuesSegmentDataGenerator: *dg,
	}
}

func (dg *CustomValuesSegmentDataGenerator) Describe() DataSpec {
	return DataNodeDataSpec{
		name: "Custom Values",
	}
}

// Check at compile time whether CustomValuesSegmentDataIterator implements DataIterator interface.
var _ metrics.DataIterator = (*CustomValuesSegmentDataIterator)(nil)

type CustomValuesSegmentDataIterator struct {
	customValuesSegmentDataGenerator CustomValuesSegmentDataGenerator

	// iterIndex keeps track of the current iteration.
	iterIndex int
}

// Evaluate fulfills the metrics.DataIterator interface.
// This function is responsible for returning the data points one at a time.
func (di *CustomValuesSegmentDataIterator) Evaluate(scrapeInfo metrics.ScrapeInfo) metrics.ScrapeResult {
	// Have we reached the end?
	if di.iterIndex >= len(di.customValuesSegmentDataGenerator.values) {
		return metrics.ScrapeResult{Exhausted: true}
	}

	// Make sure to increment the iterator index before leaving the function
	defer func() { di.iterIndex++ }()

	result := di.customValuesSegmentDataGenerator.values[di.iterIndex]

	return metrics.ScrapeResult{
		Value:   result.Value,
		Missing: result.Missing,
	}
}

// CustomValueSample contains the scrape value to be returned by CustomValuesSegmentDataIterator.
type CustomValueSample struct {
	// Value is the value of the sample.
	Value float64

	// Missing indicates whether the scrape failed to retrieve a sample.
	// Used to simulate failed scrapes.
	Missing bool
}
