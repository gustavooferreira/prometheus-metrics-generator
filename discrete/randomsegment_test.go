package discrete_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gustavooferreira/prometheus-metrics-generator/discrete"
)

func TestRandomDataIterator(t *testing.T) {
	t.Run("should fail given that iteration count limit is set to zero", func(t *testing.T) {
		_, err := discrete.NewRandomDataGenerator(discrete.RandomSegmentDataGeneratorOptions{
			AmplitudeMin: 11,
			AmplitudeMax: 20,
		})
		require.Error(t, err)
		expectedErrorMessage := "error validating random segment data generator configuration: iteration count limit cannot be less than or equal to zero"
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("should produce valid results for the given count", func(t *testing.T) {
		dataGenerator, err := discrete.NewRandomDataGenerator(discrete.RandomSegmentDataGeneratorOptions{
			AmplitudeMin:        11,
			AmplitudeMax:        20,
			IterationCountLimit: 10,
		})
		require.NoError(t, err)

		results := helperScraper(t, dataGenerator.Iterator())

		require.Equal(t, 10, len(results))

		// assert all values are within amplitude range provided
		for _, result := range results {
			assert.GreaterOrEqual(t, 20.0, result.scrapeResult.Value)
			assert.LessOrEqual(t, 11.0, result.scrapeResult.Value)
		}
	})
}
