package discrete_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gustavooferreira/prometheus-metrics-generator/discrete"
)

func TestVoidDataIterator(t *testing.T) {
	t.Run("should not return any sample when count is zero", func(t *testing.T) {
		dataGenerator := discrete.Void(0)
		results := helperScraper(t, dataGenerator.Iterator())
		require.Equal(t, 0, len(results))
	})

	t.Run("should not return any sample when count is negative", func(t *testing.T) {
		dataGenerator := discrete.Void(-5)
		results := helperScraper(t, dataGenerator.Iterator())
		require.Equal(t, 0, len(results))
	})

	t.Run("should produce valid results for the given count", func(t *testing.T) {
		dataGenerator := discrete.Void(5)
		results := helperScraper(t, dataGenerator.Iterator())

		require.Equal(t, 5, len(results))
		assert.True(t, results[0].scrapeResult.Missing)
		assert.True(t, results[1].scrapeResult.Missing)
		assert.True(t, results[2].scrapeResult.Missing)
		assert.True(t, results[3].scrapeResult.Missing)
		assert.True(t, results[4].scrapeResult.Missing)
	})
}
