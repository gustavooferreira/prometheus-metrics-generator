package discrete

import (
	"github.com/gustavooferreira/prometheus-metrics-generator/series"
)

// Check at compile time whether LoopDataGenerator implements DataGenerator interface.
var _ DataGenerator = (*LoopDataGenerator)(nil)

type LoopDataGenerator struct {
	dataGenerator DataGenerator
	count         int
}

// Loop loops over the DataGenerator N times.
func Loop(dataGenerator DataGenerator, count int) *LoopDataGenerator {
	return &LoopDataGenerator{
		dataGenerator: dataGenerator,
		count:         count,
	}
}

func (ldg *LoopDataGenerator) Iterator() DataIterator {
	return &LoopDataIterator{
		loopDataGenerator: *ldg,
	}
}

func (ldg *LoopDataGenerator) Describe() DataSpec {
	return nil
}

// Check at compile time whether LoopDataIterator implements DataIterator interface.
var _ DataIterator = (*LoopDataIterator)(nil)

type LoopDataIterator struct {
	loopDataGenerator LoopDataGenerator

	// these variables keep track of the current state of the iterator
	dataGeneratorLoopCount int
	dataIterator           DataIterator
}

// Iterate fulfills the DataIterator interface.
// This function is responsible for returning the data points one at a time.
func (ldi *LoopDataIterator) Iterate(scrapeInfo series.ScrapeInfo) series.ScrapeResult {
	for ; ldi.dataGeneratorLoopCount < ldi.loopDataGenerator.count; ldi.dataGeneratorLoopCount++ {
		if ldi.dataIterator == nil {
			ldi.dataIterator = ldi.loopDataGenerator.dataGenerator.Iterator()
		}

		result := ldi.dataIterator.Iterate(scrapeInfo)
		if result.Exhausted {
			ldi.dataIterator = nil
			continue
		}

		return result
	}

	return series.ScrapeResult{Exhausted: true}
}

// Check at compile time whether LoopDataSpec implements DataSpec interface.
var _ DataSpec = (*LoopDataSpec)(nil)

// LoopDataSpec implements a generic DataSpec for the Loop container.
type LoopDataSpec struct {
	Count int
}

func (lds LoopDataSpec) DataGeneratorNodeType() DataGeneratorNodeType {
	return DataGeneratorNodeTypeLoop
}

func (lds LoopDataSpec) Name() string {
	return "Loop"
}
