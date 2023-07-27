package discrete

import (
	"github.com/gustavooferreira/prometheus-metrics-generator/series"
)

// Check at compile time whether JoinDataGenerator implements DataGenerator interface.
var _ DataGenerator = (*JoinDataGenerator)(nil)

type JoinDataGenerator struct {
	dataGenerators []DataGenerator
}

// Join joins all DataGenerators, one after the next.
func Join(dataGenerators []DataGenerator) *JoinDataGenerator {
	return &JoinDataGenerator{
		dataGenerators: dataGenerators,
	}
}

func (jdg *JoinDataGenerator) Iterator() DataIterator {
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

// Check at compile time whether JoinDataIterator implements DataIterator interface.
var _ DataIterator = (*JoinDataIterator)(nil)

type JoinDataIterator struct {
	joinDataGenerator JoinDataGenerator

	// these variables keep track of the current state of the iterator
	dataGeneratorIndex int
	dataIterator       DataIterator
}

// Iterate fulfills the DataIterator interface.
// This function is responsible for returning the data points one at a time.
func (jdi *JoinDataIterator) Iterate(scrapeInfo series.ScrapeInfo) series.ScrapeResult {
	for ; jdi.dataGeneratorIndex < len(jdi.joinDataGenerator.dataGenerators); jdi.dataGeneratorIndex++ {
		if jdi.dataIterator == nil {
			jdi.dataIterator = jdi.joinDataGenerator.dataGenerators[jdi.dataGeneratorIndex].Iterator()
		}

		result := jdi.dataIterator.Iterate(scrapeInfo)
		if result.Exhausted {
			jdi.dataIterator = nil
			continue
		}

		return result
	}

	return series.ScrapeResult{Exhausted: true}
}

// Check at compile time whether JoinDataSpec implements DataSpec interface.
var _ DataSpec = (*JoinDataSpec)(nil)

// JoinDataSpec implements a generic DataSpec for the Join container.
type JoinDataSpec struct {
	Children []DataSpec
}

func (jds JoinDataSpec) DataGeneratorNodeType() DataGeneratorNodeType {
	return DataGeneratorNodeTypeJoin
}

func (jds JoinDataSpec) Name() string {
	return "Join"
}
