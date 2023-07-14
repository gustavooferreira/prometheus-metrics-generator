package discrete_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gustavooferreira/prometheus-metrics-generator/datafuncs/discrete"
	"github.com/gustavooferreira/prometheus-metrics-generator/series"
)

func TestJoinDataIterator(t *testing.T) {
	t.Run("should not return any sample when no data iterators are provided", func(t *testing.T) {
		dataIterator := discrete.JoinDataIterator()

		results := helperScraper(t, dataIterator)

		require.Equal(t, 0, len(results))
	})

	t.Run("should produce valid results for the given list of data iterators", func(t *testing.T) {
		var dataIteratorsArray []series.DataIterator

		lsDataIterator, err := discrete.NewLinearSegmentDataIterator(discrete.LinearSegmentDataIteratorOptions{
			AmplitudeStart:      10,
			AmplitudeEnd:        20,
			IterationCountLimit: 2,
		})
		require.NoError(t, err)
		dataIteratorsArray = append(dataIteratorsArray, lsDataIterator.Evaluate)

		lsDataIterator, err = discrete.NewLinearSegmentDataIterator(discrete.LinearSegmentDataIteratorOptions{
			AmplitudeStart:      40,
			AmplitudeEnd:        50,
			IterationCountLimit: 3,
		})
		require.NoError(t, err)
		dataIteratorsArray = append(dataIteratorsArray, lsDataIterator.Evaluate)

		lsDataIterator, err = discrete.NewLinearSegmentDataIterator(discrete.LinearSegmentDataIteratorOptions{
			AmplitudeStart:      70,
			AmplitudeEnd:        70,
			IterationCountLimit: 4,
		})
		require.NoError(t, err)
		dataIteratorsArray = append(dataIteratorsArray, lsDataIterator.Evaluate)

		// ------------

		dataIterator := discrete.JoinDataIterator(dataIteratorsArray...)

		results := helperScraper(t, dataIterator)

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

	t.Run("should produce valid results for the given list of data iterators, join other joins", func(t *testing.T) {
		lsDataIterator1, err := discrete.NewLinearSegmentDataIterator(discrete.LinearSegmentDataIteratorOptions{
			AmplitudeStart:      10,
			AmplitudeEnd:        20,
			IterationCountLimit: 2,
		})
		require.NoError(t, err)

		lsDataIterator2, err := discrete.NewLinearSegmentDataIterator(discrete.LinearSegmentDataIteratorOptions{
			AmplitudeStart:      40,
			AmplitudeEnd:        50,
			IterationCountLimit: 3,
		})
		require.NoError(t, err)

		dataIteratorJoin1 := discrete.JoinDataIterator(lsDataIterator1.Evaluate, lsDataIterator2.Evaluate)

		lsDataIterator3, err := discrete.NewLinearSegmentDataIterator(discrete.LinearSegmentDataIteratorOptions{
			AmplitudeStart:      70,
			AmplitudeEnd:        70,
			IterationCountLimit: 4,
		})
		require.NoError(t, err)

		lsDataIterator4, err := discrete.NewLinearSegmentDataIterator(discrete.LinearSegmentDataIteratorOptions{
			AmplitudeStart:      60,
			AmplitudeEnd:        30,
			IterationCountLimit: 4,
		})
		require.NoError(t, err)

		dataIteratorJoin2 := discrete.JoinDataIterator(lsDataIterator3.Evaluate, lsDataIterator4.Evaluate)

		// ------------

		dataIterator := discrete.JoinDataIterator(dataIteratorJoin1, dataIteratorJoin2)

		results := helperScraper(t, dataIterator)

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
