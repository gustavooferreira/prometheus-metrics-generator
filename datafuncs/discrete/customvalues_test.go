package discrete_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gustavooferreira/prometheus-metrics-generator/datafuncs/discrete"
)

func TestCustomValuesDataIterator(t *testing.T) {
	t.Run("should not return any sample when values array is nil", func(t *testing.T) {
		dataGenerator := discrete.NewCustomValuesDataGenerator(nil)
		results := helperScraper(t, dataGenerator.Iterator())
		require.Equal(t, 0, len(results))
	})

	t.Run("should not return any sample when values array is empty", func(t *testing.T) {
		dataGenerator := discrete.NewCustomValuesDataGenerator([]discrete.CustomValue{})
		results := helperScraper(t, dataGenerator.Iterator())
		require.Equal(t, 0, len(results))
	})

	t.Run("should produce valid results for the given count", func(t *testing.T) {
		values := []discrete.CustomValue{
			{Value: 1},
			{Value: 2},
			{Value: 3},
		}
		dataGenerator := discrete.NewCustomValuesDataGenerator(values)
		results := helperScraper(t, dataGenerator.Iterator())

		require.Equal(t, 3, len(results))
		assert.InDelta(t, 1, results[0].scrapeResult.Value, 0.001)
		assert.InDelta(t, 2, results[1].scrapeResult.Value, 0.001)
		assert.InDelta(t, 3, results[2].scrapeResult.Value, 0.001)
	})
}
