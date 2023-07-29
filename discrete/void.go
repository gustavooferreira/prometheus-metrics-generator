package discrete

import (
	"github.com/gustavooferreira/prometheus-metrics-generator/metrics"
)

// Check at compile time whether VoidDataGenerator implements DataGenerator interface.
var _ DataGenerator = (*VoidDataGenerator)(nil)

// VoidDataGenerator returns a DataGenerator representing a period of missing scrapes.
// This function is useful to simulate a chunk of time when a given metric goes missing.
type VoidDataGenerator struct {
	count int
}

// NewVoidDataGenerator returns a DataGenerator representing a void segment.
func NewVoidDataGenerator(count int) *VoidDataGenerator {
	return &VoidDataGenerator{
		count: count,
	}
}

func (vdg *VoidDataGenerator) Iterator() metrics.DataIterator {
	return &VoidDataIterator{
		voidDataGenerator: *vdg,
	}
}

func (vdg *VoidDataGenerator) Describe() DataSpec {
	return DataNodeDataSpec{
		name: "Void",
	}
}

// Check at compile time whether VoidDataIterator implements metrics.DataIterator interface.
var _ metrics.DataIterator = (*VoidDataIterator)(nil)

type VoidDataIterator struct {
	voidDataGenerator VoidDataGenerator

	// iterCount represents the cycle the iterator is in
	iterCount int
}

// Evaluate fulfills the metrics.DataIterator interface.
// This function is responsible for returning the data points one at a time.
func (vdi *VoidDataIterator) Evaluate(scrapeInfo metrics.ScrapeInfo) metrics.ScrapeResult {
	for vdi.iterCount < vdi.voidDataGenerator.count {
		vdi.iterCount++
		return metrics.ScrapeResult{Missing: true}
	}

	return metrics.ScrapeResult{Exhausted: true}
}
