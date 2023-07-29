package discrete_test

import (
	"fmt"

	"github.com/gustavooferreira/prometheus-metrics-generator/discrete"
)

func ExampleDataSpec() {
	lsDataGenerator1, err := discrete.NewLinearSegment(discrete.LinearSegmentOptions{
		AmplitudeStart:      50,
		AmplitudeEnd:        70,
		IterationCountLimit: 5,
	})
	if err != nil {
		panic(err)
	}

	lsDataGenerator2, err := discrete.NewLinearSegment(discrete.LinearSegmentOptions{
		AmplitudeStart:      50,
		AmplitudeEnd:        70,
		IterationCountLimit: 5,
	})
	if err != nil {
		panic(err)
	}

	lsDataGenerator2a := discrete.Loop(lsDataGenerator2, 3)

	dataGenerator := discrete.Join([]discrete.DataGenerator{lsDataGenerator1, lsDataGenerator2a})

	rootDataSpec := dataGenerator.Describe()
	result := discrete.Describe(rootDataSpec)
	fmt.Println(result)
	// Output:
	// Join
	//   Linear Segment
	//   Loop [3]
	//     Linear Segment
}
