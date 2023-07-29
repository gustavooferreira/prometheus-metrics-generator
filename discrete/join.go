package discrete

import (
	"github.com/gustavooferreira/prometheus-metrics-generator/metrics"
)

// Check at compile time whether JoinDataGenerator implements DataGenerator interface.
var _ DataGenerator = (*JoinDataGenerator)(nil)

type JoinDataGenerator struct {
	dataGenerators []DataGenerator
}

// NewJoinDataGenerator joins all DataGenerators, one after the next.
func NewJoinDataGenerator(dataGenerators []DataGenerator) *JoinDataGenerator {
	return &JoinDataGenerator{
		dataGenerators: dataGenerators,
	}
}

func (jdg *JoinDataGenerator) Iterator() metrics.DataIterator {
	return &JoinDataIterator{
		joinDataGenerator: *jdg,
	}
}

func (jdg *JoinDataGenerator) Describe() DataSpec {
	var dataSpecs []DataSpec

	for _, dataGenerator := range jdg.dataGenerators {
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
func (jdi *JoinDataIterator) Evaluate(scrapeInfo metrics.ScrapeInfo) metrics.ScrapeResult {
	for ; jdi.dataGeneratorIndex < len(jdi.joinDataGenerator.dataGenerators); jdi.dataGeneratorIndex++ {
		if jdi.dataIterator == nil {
			jdi.dataIterator = jdi.joinDataGenerator.dataGenerators[jdi.dataGeneratorIndex].Iterator()
		}

		result := jdi.dataIterator.Evaluate(scrapeInfo)
		if result.Exhausted {
			jdi.dataIterator = nil
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

func (jds JoinDataSpec) DataGeneratorNodeType() DataGeneratorNodeType {
	return DataGeneratorNodeTypeJoin
}

func (jds JoinDataSpec) Name() string {
	return "Join"
}
