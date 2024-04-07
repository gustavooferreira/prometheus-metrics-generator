package discrete_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gustavooferreira/prometheus-metrics-generator/discrete"
)

func TestCustomValuesDataIterator(t *testing.T) {
	t.Run("should not return any sample when values array is nil", func(t *testing.T) {
		dataGenerator := discrete.NewCustomValuesDataGenerator(nil)
		results := helperScraper(t, dataGenerator.Iterator())
		require.Equal(t, 0, len(results))
	})

	t.Run("should not return any sample when values array is empty", func(t *testing.T) {
		dataGenerator := discrete.NewCustomValuesDataGenerator([]discrete.CustomValueSample{})
		results := helperScraper(t, dataGenerator.Iterator())
		require.Equal(t, 0, len(results))
	})

	t.Run("should produce valid results for the given count", func(t *testing.T) {
		values := []discrete.CustomValueSample{
			{Value: 1},
			{Value: 2},
			{Value: 3},
			{Missing: true},
			{Value: 5},
		}
		dataGenerator := discrete.NewCustomValuesDataGenerator(values)
		results := helperScraper(t, dataGenerator.Iterator())

		require.Equal(t, 5, len(results))
		assert.InDelta(t, 1, results[0].scrapeResult.Value, 0.001)
		assert.False(t, results[0].scrapeResult.Missing)
		assert.InDelta(t, 2, results[1].scrapeResult.Value, 0.001)
		assert.False(t, results[1].scrapeResult.Missing)
		assert.InDelta(t, 3, results[2].scrapeResult.Value, 0.001)
		assert.False(t, results[2].scrapeResult.Missing)
		assert.True(t, results[3].scrapeResult.Missing)
		assert.InDelta(t, 5, results[4].scrapeResult.Value, 0.001)
		assert.False(t, results[4].scrapeResult.Missing)
	})
}
