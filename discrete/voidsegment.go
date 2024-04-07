package discrete

import (
	"github.com/gustavooferreira/prometheus-metrics-generator/metrics"
)

// Check at compile time whether VoidDataGenerator implements DataGenerator interface.
var _ DataGenerator = (*VoidSegmentDataGenerator)(nil)

// VoidSegmentDataGenerator returns a DataGenerator representing a period of missing scrapes.
// This function is useful to simulate a chunk of time when a given metric goes missing.
type VoidSegmentDataGenerator struct {
	count int
}

// NewVoidSegmentDataGenerator returns a DataGenerator representing a void segment.
func NewVoidSegmentDataGenerator(count int) *VoidSegmentDataGenerator {
	return &VoidSegmentDataGenerator{
		count: count,
	}
}

func (dg *VoidSegmentDataGenerator) Iterator() metrics.DataIterator {
	return &VoidSegmentDataIterator{
		voidSegmentDataGenerator: *dg,
	}
}

func (dg *VoidSegmentDataGenerator) Describe() DataSpec {
	return DataNodeDataSpec{
		name: "Void",
	}
}

// Check at compile time whether VoidDataIterator implements metrics.DataIterator interface.
var _ metrics.DataIterator = (*VoidSegmentDataIterator)(nil)

type VoidSegmentDataIterator struct {
	// read-only access
	voidSegmentDataGenerator VoidSegmentDataGenerator

	// iterIndex keeps track of the current iteration.
	iterIndex int
}

// Evaluate fulfills the metrics.DataIterator interface.
// This function is responsible for returning the data points one at a time.
func (di *VoidSegmentDataIterator) Evaluate(scrapeInfo metrics.ScrapeInfo) metrics.ScrapeResult {
	// Have we reached the end?
	if di.iterIndex >= di.voidSegmentDataGenerator.count {
		return metrics.ScrapeResult{Exhausted: true}
	}

	di.iterIndex++
	return metrics.ScrapeResult{Missing: true}
}
