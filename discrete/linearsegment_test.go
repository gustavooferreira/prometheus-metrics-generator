package discrete_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gustavooferreira/prometheus-metrics-generator/discrete"
	"github.com/gustavooferreira/prometheus-metrics-generator/metrics"
)

func TestLinearSegmentDataIterator(t *testing.T) {
	t.Run("should fail given that iteration count limit is set to zero", func(t *testing.T) {
		_, err := discrete.NewLinearSegmentDataGenerator(discrete.LinearSegmentDataGeneratorOptions{
			AmplitudeStart: 11,
			AmplitudeEnd:   20,
		})
		require.Error(t, err)
		expectedErrorMessage := "error validating linear segment data generator configuration: iteration count limit cannot be less than or equal to zero"
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("should produce valid results for the given iteration count limit", func(t *testing.T) {
		lsDataGenerator, err := discrete.NewLinearSegmentDataGenerator(discrete.LinearSegmentDataGeneratorOptions{
			AmplitudeStart:      10,
			AmplitudeEnd:        20,
			IterationCountLimit: 10,
		})
		require.NoError(t, err)

		results := helperScraper(t, lsDataGenerator.Iterator())

		require.Equal(t, 10, len(results))
		assert.InDelta(t, 10, results[0].scrapeResult.Value, 0.001)
		assert.InDelta(t, 20, results[9].scrapeResult.Value, 0.001)
	})

	t.Run("should produce valid results for an iteration count limit of 1", func(t *testing.T) {
		lsDataGenerator, err := discrete.NewLinearSegmentDataGenerator(discrete.LinearSegmentDataGeneratorOptions{
			AmplitudeStart:      10,
			AmplitudeEnd:        20,
			IterationCountLimit: 1,
		})
		require.NoError(t, err)

		results := helperScraper(t, lsDataGenerator.Iterator())

		require.Equal(t, 1, len(results))
		assert.InDelta(t, 10, results[0].scrapeResult.Value, 0.001)
	})

	t.Run("should produce valid results for an iteration count limit of 2", func(t *testing.T) {
		lsDataGenerator, err := discrete.NewLinearSegmentDataGenerator(discrete.LinearSegmentDataGeneratorOptions{
			AmplitudeStart:      10,
			AmplitudeEnd:        20,
			IterationCountLimit: 2,
		})
		require.NoError(t, err)

		results := helperScraper(t, lsDataGenerator.Iterator())

		require.Equal(t, 2, len(results))
		assert.InDelta(t, 10, results[0].scrapeResult.Value, 0.001)
		assert.InDelta(t, 20, results[1].scrapeResult.Value, 0.001)
	})

	t.Run("should produce valid results for a horizontal line", func(t *testing.T) {
		lsDataGenerator, err := discrete.NewLinearSegmentDataGenerator(discrete.LinearSegmentDataGeneratorOptions{
			AmplitudeStart:      50,
			AmplitudeEnd:        50,
			IterationCountLimit: 7,
		})
		require.NoError(t, err)

		results := helperScraper(t, lsDataGenerator.Iterator())

		require.Equal(t, 7, len(results))
		assert.InDelta(t, 50, results[0].scrapeResult.Value, 0.001)
		assert.InDelta(t, 50, results[6].scrapeResult.Value, 0.001)
	})

	t.Run("should produce valid results even though data is shifted in time", func(t *testing.T) {
		scraper, err := metrics.NewScraper(
			metrics.ScraperConfig{
				StartTime:      time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
				ScrapeInterval: 15 * time.Second,
			},
			metrics.WithScraperIterationCountLimit(100), // It's good practice to set an upper bound in tests
		)
		require.NoError(t, err)

		lsDataGenerator, err := discrete.NewLinearSegmentDataGenerator(discrete.LinearSegmentDataGeneratorOptions{
			AmplitudeStart:      20,
			AmplitudeEnd:        40,
			IterationCountLimit: 9,
		})
		require.NoError(t, err)

		type resultContainer struct {
			scrapeInfo   metrics.ScrapeInfo
			scrapeResult metrics.ScrapeResult
		}

		var results []resultContainer
		scrapeHandler := func(scrapeInfo metrics.ScrapeInfo, scrapeResult metrics.ScrapeResult) error {
			results = append(results, resultContainer{
				scrapeInfo:   scrapeInfo,
				scrapeResult: scrapeResult,
			})
			return nil
		}

		iterator := lsDataGenerator.Iterator()

		skipNTimes := 25
		skipCount := 0
		iter := scraper.Iterator()
		for scrapeInfo, ok := iter.Next(); ok; scrapeInfo, ok = iter.Next() {

			if skipCount < skipNTimes {
				skipCount++
				continue
			}

			scrapeResult := iterator.Evaluate(scrapeInfo)
			if scrapeResult.Exhausted {
				// exhausted time series samples
				break
			}

			err := scrapeHandler(scrapeInfo, scrapeResult)
			require.NoError(t, err)
		}

		for _, r := range results {
			t.Logf("[%3d] Timestamp: %s - Value: %.2f\n",
				r.scrapeInfo.IterationIndex,
				r.scrapeInfo.IterationTime,
				r.scrapeResult.Value,
			)
		}

		require.Equal(t, 9, len(results))
		assert.Equal(t, 25, results[0].scrapeInfo.IterationIndex)
		assert.Equal(t,
			time.Date(2023, 1, 1, 10, 36, 15, 0, time.UTC),
			results[0].scrapeInfo.IterationTime,
		)
		assert.InDelta(t, 20, results[0].scrapeResult.Value, 0.001)

		assert.Equal(t, 33, results[8].scrapeInfo.IterationIndex)
		assert.Equal(t,
			time.Date(2023, 1, 1, 10, 38, 15, 0, time.UTC),
			results[8].scrapeInfo.IterationTime,
		)
		assert.InDelta(t, 40, results[8].scrapeResult.Value, 0.001)
	})
}
