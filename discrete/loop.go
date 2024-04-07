package discrete

import (
	"github.com/gustavooferreira/prometheus-metrics-generator/metrics"
)

// Check at compile time whether LoopDataGenerator implements DataGenerator interface.
var _ DataGenerator = (*LoopDataGenerator)(nil)

// LoopDataGenerator loops over the DataGenerator N times.
type LoopDataGenerator struct {
	dataGenerator DataGenerator
	count         int
}

// NewLoopDataGenerator creates a new instance of LoopDataGenerator.
// A negative count will produce no samples.
func NewLoopDataGenerator(dataGenerator DataGenerator, count int) *LoopDataGenerator {
	return &LoopDataGenerator{
		dataGenerator: dataGenerator,
		count:         count,
	}
}

func (dg *LoopDataGenerator) Iterator() metrics.DataIterator {
	return &LoopDataIterator{
		loopDataGenerator: *dg,
	}
}

func (dg *LoopDataGenerator) Describe() DataSpec {
	return LoopDataSpec{
		Count: dg.count,
		Func:  dg.dataGenerator.Describe(),
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

func (ds LoopDataSpec) DataGeneratorNodeType() DataGeneratorNodeType {
	return DataGeneratorNodeTypeLoop
}

func (ds LoopDataSpec) Name() string {
	return "Loop"
}
