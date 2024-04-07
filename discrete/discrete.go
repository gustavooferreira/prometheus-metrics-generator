package discrete

import (
	"fmt"
	"strings"

	"github.com/gustavooferreira/prometheus-metrics-generator/metrics"
)

// DataGenerator generates data according to the generator.
// It's meant to be used by Counters and Gauges.
type DataGenerator interface {
	Iterator() metrics.DataIterator
	Describe() DataSpec
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

// Describe generates the tree of all nodes.
func Describe(rootDataSpec DataSpec) string {
	result := describe(rootDataSpec, 0, nil)
	return strings.Join(result, "\n")
}

func describe(dataSpec DataSpec, indent int, result []string) []string {
	prefix := strings.Repeat("  ", indent)

	switch dataSpecConcrete := dataSpec.(type) {
	case DataNodeDataSpec:
		result = append(result, fmt.Sprintf("%s%s", prefix, dataSpecConcrete.Name()))
		return result
	case JoinDataSpec:
		result = append(result, fmt.Sprintf("%s%s", prefix, dataSpecConcrete.Name()))
		for _, children := range dataSpecConcrete.Children {
			result = describe(children, indent+1, result)
		}
		return result
	case LoopDataSpec:
		result = append(result, fmt.Sprintf("%s%s [%d]", prefix, dataSpecConcrete.Name(), dataSpecConcrete.Count))
		result = describe(dataSpecConcrete.Func, indent+1, result)
		return result
	default:
		result = append(result, "||--- error")
		return result
	}
}
