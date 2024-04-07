package discrete

import (
	"github.com/gustavooferreira/prometheus-metrics-generator/metrics"
)

// Check at compile time whether JoinDataGenerator implements DataGenerator interface.
var _ DataGenerator = (*JoinDataGenerator)(nil)

// JoinDataGenerator joins several segments together to form a bigger and more complex segments.
// It joins all DataGenerators one after the next.
type JoinDataGenerator struct {
	dataGenerators []DataGenerator
}

// NewJoinDataGenerator creates a new instance of JoinDataGenerator.
func NewJoinDataGenerator(dataGenerators []DataGenerator) *JoinDataGenerator {
	return &JoinDataGenerator{
		dataGenerators: dataGenerators,
	}
}

func (dg *JoinDataGenerator) Iterator() metrics.DataIterator {
	return &JoinDataIterator{
		joinDataGenerator: *dg,
	}
}

func (dg *JoinDataGenerator) Describe() DataSpec {
	var dataSpecs []DataSpec

	for _, dataGenerator := range dg.dataGenerators {
		dataSpecs = append(dataSpecs, dataGenerator.Describe())
	}

	return JoinDataSpec{
		Children: dataSpecs,
	}
}

// Check at compile time whether JoinDataIterator implements metrics.DataIterator interface.
var _ metrics.DataIterator = (*JoinDataIterator)(nil)

type JoinDataIterator struct {
	joinDataGenerator JoinDataGenerator

	// these variables keep track of the current state of the iterator
	dataGeneratorIndex int
	dataIterator       metrics.DataIterator
}

// Evaluate fulfills the metrics.DataIterator interface.
// This function is responsible for returning the data points one at a time.
func (di *JoinDataIterator) Evaluate(scrapeInfo metrics.ScrapeInfo) metrics.ScrapeResult {
	for ; di.dataGeneratorIndex < len(di.joinDataGenerator.dataGenerators); di.dataGeneratorIndex++ {
		if di.dataIterator == nil {
			di.dataIterator = di.joinDataGenerator.dataGenerators[di.dataGeneratorIndex].Iterator()
		}

		result := di.dataIterator.Evaluate(scrapeInfo)
		if result.Exhausted {
			di.dataIterator = nil
			continue
		}

		return result
	}

	return metrics.ScrapeResult{Exhausted: true}
}

// Check at compile time whether JoinDataSpec implements DataSpec interface.
var _ DataSpec = (*JoinDataSpec)(nil)

// JoinDataSpec implements a generic DataSpec for the JoinDataGenerator container.
type JoinDataSpec struct {
	Children []DataSpec
}

func (ds JoinDataSpec) DataGeneratorNodeType() DataGeneratorNodeType {
	return DataGeneratorNodeTypeJoin
}

func (ds JoinDataSpec) Name() string {
	return "Join"
}
