package continuous_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gustavooferreira/prometheus-metrics-generator/datafuncs/continuous"
)

func TestLinearSegmentDataIterator(t *testing.T) {
	t.Run("should fail given that duration length is negative", func(t *testing.T) {
		_, err := continuous.NewLinearSegmentDataIterator(continuous.LinearSegmentDataIteratorOptions{
			AmplitudeStart: 11,
			AmplitudeEnd:   20,
			DurationLength: -1 * time.Second,
		})
		require.Error(t, err)
		expectedErrorMessage := "duration length cannot be less than or equal to zero"
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("should produce valid results for the given duration length, closed interval", func(t *testing.T) {
		lsDataIterator, err := continuous.NewLinearSegmentDataIterator(continuous.LinearSegmentDataIteratorOptions{
			AmplitudeStart: 10,
			AmplitudeEnd:   20,
			DurationLength: time.Minute,
		})
		require.NoError(t, err)

		functionStartTime := time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC)

		results := helperScraper(t, lsDataIterator, functionStartTime)

		require.Equal(t, 5, len(results))
		assert.InDelta(t, 10, results[0].scrapeResult.Value, 0.001)
		assert.InDelta(t, 20, results[4].scrapeResult.Value, 0.001)
	})

	t.Run("should produce horizontal line for the given duration length, closed interval", func(t *testing.T) {
		lsDataIterator, err := continuous.NewLinearSegmentDataIterator(continuous.LinearSegmentDataIteratorOptions{
			AmplitudeStart: 50,
			AmplitudeEnd:   50,
			DurationLength: 2 * time.Minute,
		})
		require.NoError(t, err)

		functionStartTime := time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC)

		results := helperScraper(t, lsDataIterator, functionStartTime)

		require.Equal(t, 9, len(results))
		assert.InDelta(t, 50, results[0].scrapeResult.Value, 0.001)
		assert.InDelta(t, 50, results[8].scrapeResult.Value, 0.001)
	})

	t.Run("should produce segment with negative slope for the given duration length, closed interval", func(t *testing.T) {
		lsDataIterator, err := continuous.NewLinearSegmentDataIterator(continuous.LinearSegmentDataIteratorOptions{
			AmplitudeStart: 50,
			AmplitudeEnd:   20,
			DurationLength: time.Minute,
		})
		require.NoError(t, err)

		functionStartTime := time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC)

		results := helperScraper(t, lsDataIterator, functionStartTime)

		require.Equal(t, 5, len(results))
		assert.InDelta(t, 50, results[0].scrapeResult.Value, 0.001)
		assert.InDelta(t, 20, results[4].scrapeResult.Value, 0.001)
	})

	t.Run("should produce segment with negative values for the given duration length, closed interval", func(t *testing.T) {
		lsDataIterator, err := continuous.NewLinearSegmentDataIterator(continuous.LinearSegmentDataIteratorOptions{
			AmplitudeStart: -20,
			AmplitudeEnd:   20,
			DurationLength: time.Minute,
		})
		require.NoError(t, err)

		functionStartTime := time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC)

		results := helperScraper(t, lsDataIterator, functionStartTime)

		require.Equal(t, 5, len(results))
		assert.InDelta(t, -20, results[0].scrapeResult.Value, 0.001)
		assert.InDelta(t, 20, results[4].scrapeResult.Value, 0.001)
	})

	t.Run("should produce segment with negative slope and negative values for the given duration length, closed interval", func(t *testing.T) {
		lsDataIterator, err := continuous.NewLinearSegmentDataIterator(continuous.LinearSegmentDataIteratorOptions{
			AmplitudeStart: 20,
			AmplitudeEnd:   -20,
			DurationLength: time.Minute,
		})
		require.NoError(t, err)

		functionStartTime := time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC)

		results := helperScraper(t, lsDataIterator, functionStartTime)

		require.Equal(t, 5, len(results))
		assert.InDelta(t, 20, results[0].scrapeResult.Value, 0.001)
		assert.InDelta(t, -20, results[4].scrapeResult.Value, 0.001)
	})

	t.Run("should produce valid results for the given duration length, open interval", func(t *testing.T) {
		lsDataIterator, err := continuous.NewLinearSegmentDataIterator(continuous.LinearSegmentDataIteratorOptions{
			AmplitudeStart:         10,
			AmplitudeEnd:           20,
			DurationLength:         time.Minute,
			IntervalLeftBoundOpen:  true,
			IntervalRightBoundOpen: true,
		})
		require.NoError(t, err)

		functionStartTime := time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC)

		results := helperScraper(t, lsDataIterator, functionStartTime)

		require.Equal(t, 4, len(results))

		assert.Equal(t,
			time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			results[0].scrapeInfo.IterationTime,
		)
		assert.InDelta(t, 0, results[0].scrapeResult.Value, 0.001)
		assert.True(t, results[0].scrapeResult.Missing)

		assert.Equal(t,
			time.Date(2023, 1, 1, 10, 30, 45, 0, time.UTC),
			results[3].scrapeInfo.IterationTime,
		)
		assert.InDelta(t, 17.5, results[3].scrapeResult.Value, 0.001)
	})

	t.Run("should produce missing value followed by exhausted given duration length and function start time", func(t *testing.T) {
		lsDataIterator, err := continuous.NewLinearSegmentDataIterator(continuous.LinearSegmentDataIteratorOptions{
			AmplitudeStart: 5,
			AmplitudeEnd:   10,
			DurationLength: time.Second,
		})
		require.NoError(t, err)

		functionStartTime := time.Date(2023, 1, 1, 10, 30, 10, 0, time.UTC)

		results := helperScraper(t, lsDataIterator, functionStartTime)

		require.Equal(t, 1, len(results))
		assert.True(t, results[0].scrapeResult.Missing)
		assert.InDelta(t, 0, results[0].scrapeResult.Value, 0.001)
	})

	t.Run("should produce valid value once given duration length and function start time", func(t *testing.T) {
		lsDataIterator, err := continuous.NewLinearSegmentDataIterator(continuous.LinearSegmentDataIteratorOptions{
			AmplitudeStart: 10,
			AmplitudeEnd:   20,
			DurationLength: time.Minute,
		})
		require.NoError(t, err)

		functionStartTime := time.Date(2023, 1, 1, 10, 29, 0, 0, time.UTC)

		results := helperScraper(t, lsDataIterator, functionStartTime)

		require.Equal(t, 1, len(results))
		assert.InDelta(t, 20, results[0].scrapeResult.Value, 0.001)
	})

	t.Run("should produce valid results from the middle of the function for the given duration length and function start time", func(t *testing.T) {
		lsDataIterator, err := continuous.NewLinearSegmentDataIterator(continuous.LinearSegmentDataIteratorOptions{
			AmplitudeStart: 20,
			AmplitudeEnd:   60,
			DurationLength: time.Minute,
		})
		require.NoError(t, err)

		functionStartTime := time.Date(2023, 1, 1, 10, 29, 30, 0, time.UTC)

		results := helperScraper(t, lsDataIterator, functionStartTime)

		require.Equal(t, 3, len(results))
		assert.InDelta(t, 40, results[0].scrapeResult.Value, 0.001)
		assert.InDelta(t, 60, results[2].scrapeResult.Value, 0.001)
	})

	t.Run("should return no value given duration length and function start time", func(t *testing.T) {
		lsDataIterator, err := continuous.NewLinearSegmentDataIterator(continuous.LinearSegmentDataIteratorOptions{
			AmplitudeStart: 10,
			AmplitudeEnd:   20,
			DurationLength: time.Minute,
		})
		require.NoError(t, err)

		functionStartTime := time.Date(2023, 1, 1, 10, 25, 0, 0, time.UTC)

		results := helperScraper(t, lsDataIterator, functionStartTime)

		require.Equal(t, 0, len(results))
	})
}
