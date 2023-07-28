package discrete_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	discrete2 "github.com/gustavooferreira/prometheus-metrics-generator/discrete"
)

func TestLoopDataIterator(t *testing.T) {
	t.Run("should not return any sample when count is zero", func(t *testing.T) {
		lsDataGenerator, err := discrete2.NewLinearSegment(discrete2.LinearSegmentOptions{
			AmplitudeStart:      10,
			AmplitudeEnd:        20,
			IterationCountLimit: 2,
		})
		require.NoError(t, err)

		dataGenerator := discrete2.Loop(lsDataGenerator, 0)

		results := helperScraper(t, dataGenerator.Iterator())

		require.Equal(t, 0, len(results))
	})

	t.Run("should not return any sample when count is negative", func(t *testing.T) {
		lsDataGenerator, err := discrete2.NewLinearSegment(discrete2.LinearSegmentOptions{
			AmplitudeStart:      10,
			AmplitudeEnd:        20,
			IterationCountLimit: 2,
		})
		require.NoError(t, err)

		dataGenerator := discrete2.Loop(lsDataGenerator, -5)

		results := helperScraper(t, dataGenerator.Iterator())

		require.Equal(t, 0, len(results))
	})

	t.Run("should produce valid results for the given data generator and count", func(t *testing.T) {
		lsDataGenerator, err := discrete2.NewLinearSegment(discrete2.LinearSegmentOptions{
			AmplitudeStart:      10,
			AmplitudeEnd:        20,
			IterationCountLimit: 2,
		})
		require.NoError(t, err)

		dataGenerator := discrete2.Loop(lsDataGenerator, 3)

		results := helperScraper(t, dataGenerator.Iterator())

		require.Equal(t, 6, len(results))
		assert.InDelta(t, 10, results[0].scrapeResult.Value, 0.001)
		assert.InDelta(t, 20, results[1].scrapeResult.Value, 0.001)
		assert.InDelta(t, 10, results[2].scrapeResult.Value, 0.001)
		assert.InDelta(t, 20, results[3].scrapeResult.Value, 0.001)
		assert.InDelta(t, 10, results[4].scrapeResult.Value, 0.001)
		assert.InDelta(t, 20, results[5].scrapeResult.Value, 0.001)
	})
}
