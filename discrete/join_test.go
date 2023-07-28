package discrete_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	discrete2 "github.com/gustavooferreira/prometheus-metrics-generator/discrete"
)

func TestJoinDataIterator(t *testing.T) {
	t.Run("should not return any sample when no data iterators are provided", func(t *testing.T) {
		dataGenerator := discrete2.Join([]discrete2.DataGenerator{})

		results := helperScraper(t, dataGenerator.Iterator())

		require.Equal(t, 0, len(results))
	})

	t.Run("should produce valid results for the given list of data iterators", func(t *testing.T) {
		var dataGeneratorsArray []discrete2.DataGenerator

		lsDataGenerator, err := discrete2.NewLinearSegment(discrete2.LinearSegmentOptions{
			AmplitudeStart:      10,
			AmplitudeEnd:        20,
			IterationCountLimit: 2,
		})
		require.NoError(t, err)
		dataGeneratorsArray = append(dataGeneratorsArray, lsDataGenerator)

		lsDataGenerator, err = discrete2.NewLinearSegment(discrete2.LinearSegmentOptions{
			AmplitudeStart:      40,
			AmplitudeEnd:        50,
			IterationCountLimit: 3,
		})
		require.NoError(t, err)
		dataGeneratorsArray = append(dataGeneratorsArray, lsDataGenerator)

		lsDataGenerator, err = discrete2.NewLinearSegment(discrete2.LinearSegmentOptions{
			AmplitudeStart:      70,
			AmplitudeEnd:        70,
			IterationCountLimit: 4,
		})
		require.NoError(t, err)
		dataGeneratorsArray = append(dataGeneratorsArray, lsDataGenerator)

		dataGenerator := discrete2.Join(dataGeneratorsArray)

		results := helperScraper(t, dataGenerator.Iterator())

		require.Equal(t, 9, len(results))
		assert.InDelta(t, 10, results[0].scrapeResult.Value, 0.001)
		assert.InDelta(t, 20, results[1].scrapeResult.Value, 0.001)
		assert.InDelta(t, 40, results[2].scrapeResult.Value, 0.001)
		assert.InDelta(t, 45, results[3].scrapeResult.Value, 0.001)
		assert.InDelta(t, 50, results[4].scrapeResult.Value, 0.001)
		assert.InDelta(t, 70, results[5].scrapeResult.Value, 0.001)
		assert.InDelta(t, 70, results[6].scrapeResult.Value, 0.001)
		assert.InDelta(t, 70, results[7].scrapeResult.Value, 0.001)
		assert.InDelta(t, 70, results[8].scrapeResult.Value, 0.001)
	})

	t.Run("should produce valid results when joining the same data generator multiple times", func(t *testing.T) {
		var dataGeneratorsArray []discrete2.DataGenerator

		lsDataGenerator, err := discrete2.NewLinearSegment(discrete2.LinearSegmentOptions{
			AmplitudeStart:      10,
			AmplitudeEnd:        20,
			IterationCountLimit: 2,
		})
		require.NoError(t, err)
		dataGeneratorsArray = append(dataGeneratorsArray, lsDataGenerator)
		dataGeneratorsArray = append(dataGeneratorsArray, lsDataGenerator)
		dataGeneratorsArray = append(dataGeneratorsArray, lsDataGenerator)

		dataGenerator := discrete2.Join(dataGeneratorsArray)

		results := helperScraper(t, dataGenerator.Iterator())

		require.Equal(t, 6, len(results))
		assert.InDelta(t, 10, results[0].scrapeResult.Value, 0.001)
		assert.InDelta(t, 20, results[1].scrapeResult.Value, 0.001)
		assert.InDelta(t, 10, results[2].scrapeResult.Value, 0.001)
		assert.InDelta(t, 20, results[3].scrapeResult.Value, 0.001)
		assert.InDelta(t, 10, results[4].scrapeResult.Value, 0.001)
		assert.InDelta(t, 20, results[5].scrapeResult.Value, 0.001)
	})

	t.Run("should produce valid results for the given list of data iterators, join other joins", func(t *testing.T) {
		lsDataGenerator1, err := discrete2.NewLinearSegment(discrete2.LinearSegmentOptions{
			AmplitudeStart:      10,
			AmplitudeEnd:        20,
			IterationCountLimit: 2,
		})
		require.NoError(t, err)

		lsDataGenerator2, err := discrete2.NewLinearSegment(discrete2.LinearSegmentOptions{
			AmplitudeStart:      40,
			AmplitudeEnd:        50,
			IterationCountLimit: 3,
		})
		require.NoError(t, err)

		join1 := []discrete2.DataGenerator{lsDataGenerator1, lsDataGenerator2}
		dataGeneratorJoin1 := discrete2.Join(join1)

		lsDataGenerator3, err := discrete2.NewLinearSegment(discrete2.LinearSegmentOptions{
			AmplitudeStart:      70,
			AmplitudeEnd:        70,
			IterationCountLimit: 4,
		})
		require.NoError(t, err)

		lsDataGenerator4, err := discrete2.NewLinearSegment(discrete2.LinearSegmentOptions{
			AmplitudeStart:      60,
			AmplitudeEnd:        30,
			IterationCountLimit: 4,
		})
		require.NoError(t, err)

		join2 := []discrete2.DataGenerator{lsDataGenerator3, lsDataGenerator4}
		dataGeneratorJoin2 := discrete2.Join(join2)

		// ------------

		greaterJoin := []discrete2.DataGenerator{dataGeneratorJoin1, dataGeneratorJoin2}
		dataGenerator := discrete2.Join(greaterJoin)

		results := helperScraper(t, dataGenerator.Iterator())

		require.Equal(t, 13, len(results))
		assert.InDelta(t, 10, results[0].scrapeResult.Value, 0.001)
		assert.InDelta(t, 20, results[1].scrapeResult.Value, 0.001)
		assert.InDelta(t, 40, results[2].scrapeResult.Value, 0.001)
		assert.InDelta(t, 45, results[3].scrapeResult.Value, 0.001)
		assert.InDelta(t, 50, results[4].scrapeResult.Value, 0.001)
		assert.InDelta(t, 70, results[5].scrapeResult.Value, 0.001)
		assert.InDelta(t, 70, results[6].scrapeResult.Value, 0.001)
		assert.InDelta(t, 70, results[7].scrapeResult.Value, 0.001)
		assert.InDelta(t, 70, results[8].scrapeResult.Value, 0.001)
		assert.InDelta(t, 60, results[9].scrapeResult.Value, 0.001)
		assert.InDelta(t, 50, results[10].scrapeResult.Value, 0.001)
		assert.InDelta(t, 40, results[11].scrapeResult.Value, 0.001)
		assert.InDelta(t, 30, results[12].scrapeResult.Value, 0.001)
	})
}
