package discrete

import (
	"github.com/gustavooferreira/prometheus-metrics-generator/metrics"
)

// Check at compile time whether LoopDataGenerator implements DataGenerator interface.
var _ DataGenerator = (*LoopDataGenerator)(nil)

type LoopDataGenerator struct {
	dataGenerator DataGenerator
	count         int
}

// NewLoopDataGenerator loops over the DataGenerator N times.
func NewLoopDataGenerator(dataGenerator DataGenerator, count int) *LoopDataGenerator {
	return &LoopDataGenerator{
		dataGenerator: dataGenerator,
		count:         count,
	}
}

func (ldg *LoopDataGenerator) Iterator() metrics.DataIterator {
	return &LoopDataIterator{
		loopDataGenerator: *ldg,
	}
}

func (ldg *LoopDataGenerator) Describe() DataSpec {
	return LoopDataSpec{
		Count: ldg.count,
		Func:  ldg.dataGenerator.Describe(),
	}
}

// Check at compile time whether LoopDataIterator implements metrics.DataIterator interface.
var _ metrics.DataIterator = (*LoopDataIterator)(nil)

type LoopDataIterator struct {
	loopDataGenerator LoopDataGenerator

	// these variables keep track of the current state of the iterator
	dataGeneratorLoopCount int
	dataIterator           metrics.DataIterator
}

// Evaluate fulfills the metrics.DataIterator interface.
// This function is responsible for returning the data points one at a time.
func (ldi *LoopDataIterator) Evaluate(scrapeInfo metrics.ScrapeInfo) metrics.ScrapeResult {
	for ; ldi.dataGeneratorLoopCount < ldi.loopDataGenerator.count; ldi.dataGeneratorLoopCount++ {
		if ldi.dataIterator == nil {
			ldi.dataIterator = ldi.loopDataGenerator.dataGenerator.Iterator()
		}

		result := ldi.dataIterator.Evaluate(scrapeInfo)
		if result.Exhausted {
			ldi.dataIterator = nil
			continue
		}

		return result
	}

	return metrics.ScrapeResult{Exhausted: true}
}

// Check at compile time whether LoopDataSpec implements DataSpec interface.
var _ DataSpec = (*LoopDataSpec)(nil)

// LoopDataSpec implements a generic DataSpec for the LoopDataGenerator container.
type LoopDataSpec struct {
	Count int
	Func  DataSpec
}

func (lds LoopDataSpec) DataGeneratorNodeType() DataGeneratorNodeType {
	return DataGeneratorNodeTypeLoop
}

func (lds LoopDataSpec) Name() string {
	return "Loop"
}
