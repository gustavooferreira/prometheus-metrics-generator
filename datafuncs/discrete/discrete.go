package discrete

import (
	"github.com/gustavooferreira/prometheus-metrics-generator/series"
)

type DataGenerator interface {
	Iterator() DataIterator
	Describe() DataSpec
}

type DataIterator interface {
	Iterate(series.ScrapeInfo) series.ScrapeResult
}

// DataSpec defines the data node type.
// It's necessary to type assert to the type returned by the DataGeneratorNodeType method.
type DataSpec interface {
	DataGeneratorNodeType() DataGeneratorNodeType
	Name() string
}

type DataGeneratorNodeType string

const (
	DataGeneratorNodeTypeData DataGeneratorNodeType = "data_generator_node_type-data"
	DataGeneratorNodeTypeJoin DataGeneratorNodeType = "data_generator_node_type-join"
	DataGeneratorNodeTypeLoop DataGeneratorNodeType = "data_generator_node_type-loop"
)

// Check at compile time whether DataNodeDataSpec implements DataSpec interface.
var _ DataSpec = (*DataNodeDataSpec)(nil)

// DataNodeDataSpec implements a generic DataSpec for data shapes.
type DataNodeDataSpec struct {
	name string
}

func (ds DataNodeDataSpec) DataGeneratorNodeType() DataGeneratorNodeType {
	return DataGeneratorNodeTypeData
}

func (ds DataNodeDataSpec) Name() string {
	return ds.name
}
