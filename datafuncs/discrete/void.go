package discrete

import (
	"github.com/gustavooferreira/prometheus-metrics-generator/series"
)

// Check at compile time whether VoidDataGenerator implements DataGenerator interface.
var _ DataGenerator = (*VoidDataGenerator)(nil)

type VoidDataGenerator struct {
	count int
}

// Void voids over the DataGenerator N times.
func Void(count int) *VoidDataGenerator {
	return &VoidDataGenerator{
		count: count,
	}
}

func (vdg *VoidDataGenerator) Iterator() DataIterator {
	return &VoidDataIterator{
		voidDataGenerator: *vdg,
	}
}

func (vdg *VoidDataGenerator) Describe() DataSpec {
	return DataNodeDataSpec{
		name: "Void",
	}
}

// Check at compile time whether VoidDataIterator implements DataIterator interface.
var _ DataIterator = (*VoidDataIterator)(nil)

type VoidDataIterator struct {
	voidDataGenerator VoidDataGenerator

	// these variables keep track of the current state of the iterator
	dataGeneratorVoidCount int
}

// Iterate fulfills the DataIterator interface.
// This function is responsible for returning the data points one at a time.
func (vdi *VoidDataIterator) Iterate(scrapeInfo series.ScrapeInfo) series.ScrapeResult {
	for vdi.dataGeneratorVoidCount < vdi.voidDataGenerator.count {
		vdi.dataGeneratorVoidCount++
		return series.ScrapeResult{Missing: true}
	}

	return series.ScrapeResult{Exhausted: true}
}
