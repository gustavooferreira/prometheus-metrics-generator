package discrete_test

import (
	"fmt"

	"github.com/gustavooferreira/prometheus-metrics-generator/discrete"
)

func ExampleDataSpec() {
	lsDataGenerator1, err := discrete.NewLinearSegmentDataGenerator(discrete.LinearSegmentDataGeneratorOptions{
		AmplitudeStart:      50,
		AmplitudeEnd:        70,
		IterationCountLimit: 5,
	})
	if err != nil {
		panic(err)
	}

	lsDataGenerator2, err := discrete.NewLinearSegmentDataGenerator(discrete.LinearSegmentDataGeneratorOptions{
		AmplitudeStart:      50,
		AmplitudeEnd:        70,
		IterationCountLimit: 5,
	})
	if err != nil {
		panic(err)
	}

	lsDataGenerator2Loop := discrete.NewLoopDataGenerator(lsDataGenerator2, 3)

	dataGenerator := discrete.NewJoinDataGenerator([]discrete.DataGenerator{lsDataGenerator1, lsDataGenerator2Loop})

	rootDataSpec := dataGenerator.Describe()
	result := discrete.Describe(rootDataSpec)
	fmt.Println(result)
	// Output:
	// Join
	//   Linear Segment
	//   Loop [3]
	//     Linear Segment
}
