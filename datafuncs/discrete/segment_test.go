package discrete_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gustavooferreira/prometheus-metrics-generator/datafuncs/discrete"
	"github.com/gustavooferreira/prometheus-metrics-generator/series"
)

func TestLinearSegmentDataIterator(t *testing.T) {
	t.Run("should fail given that iteration count limit is set to zero", func(t *testing.T) {
		_, err := discrete.NewLinearSegmentDataIterator(discrete.LinearSegmentDataIteratorOptions{
			AmplitudeStart: 11,
			AmplitudeEnd:   20,
		})
		require.Error(t, err)
		expectedErrorMessage := "iteration count limit cannot be less than or equal to zero"
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("should produce valid results for the given iteration count limit", func(t *testing.T) {
		scraper, err := series.NewScraper(
			series.ScraperConfig{
				StartTime:      time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
				ScrapeInterval: 15 * time.Second,
			},
			series.WithScraperIterationCountLimit(100), // It's good practice to set an upper bound in tests
		)
		require.NoError(t, err)

		lsDataIterator, err := discrete.NewLinearSegmentDataIterator(discrete.LinearSegmentDataIteratorOptions{
			AmplitudeStart:      10,
			AmplitudeEnd:        20,
			IterationCountLimit: 10,
		})
		require.NoError(t, err)

		type resultContainer struct {
			scrapeInfo   series.ScrapeInfo
			scrapeResult series.ScrapeResult
		}

		var results []resultContainer
		scrapeHandler := func(scrapeInfo series.ScrapeInfo, scrapeResult series.ScrapeResult) error {
			results = append(results, resultContainer{
				scrapeInfo:   scrapeInfo,
				scrapeResult: scrapeResult,
			})
			return nil
		}

		err = scraper.Scrape(lsDataIterator.Evaluate, scrapeHandler)
		require.NoError(t, err)

		for _, r := range results {
			t.Logf("[%3d] Timestamp: %s - Value: %.2f\n",
				r.scrapeInfo.IterationCount,
				r.scrapeInfo.IterationTime,
				r.scrapeResult.Value,
			)
		}

		require.Equal(t, 10, len(results))
		assert.InDelta(t, 10, results[0].scrapeResult.Value, 0.001)
		assert.InDelta(t, 20, results[9].scrapeResult.Value, 0.001)
	})

	t.Run("should produce valid results for a horizontal line", func(t *testing.T) {
		scraper, err := series.NewScraper(
			series.ScraperConfig{
				StartTime:      time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
				ScrapeInterval: 15 * time.Second,
			},
			series.WithScraperIterationCountLimit(100), // It's good practice to set an upper bound in tests
		)
		require.NoError(t, err)

		lsDataIterator, err := discrete.NewLinearSegmentDataIterator(discrete.LinearSegmentDataIteratorOptions{
			AmplitudeStart:      50,
			AmplitudeEnd:        50,
			IterationCountLimit: 7,
		})
		require.NoError(t, err)

		type resultContainer struct {
			scrapeInfo   series.ScrapeInfo
			scrapeResult series.ScrapeResult
		}

		var results []resultContainer
		scrapeHandler := func(scrapeInfo series.ScrapeInfo, scrapeResult series.ScrapeResult) error {
			results = append(results, resultContainer{
				scrapeInfo:   scrapeInfo,
				scrapeResult: scrapeResult,
			})
			return nil
		}

		err = scraper.Scrape(lsDataIterator.Evaluate, scrapeHandler)
		require.NoError(t, err)

		for _, r := range results {
			t.Logf("[%3d] Timestamp: %s - Value: %.2f\n",
				r.scrapeInfo.IterationCount,
				r.scrapeInfo.IterationTime,
				r.scrapeResult.Value,
			)
		}

		require.Equal(t, 7, len(results))
		assert.InDelta(t, 50, results[0].scrapeResult.Value, 0.001)
		assert.InDelta(t, 50, results[6].scrapeResult.Value, 0.001)
	})

	t.Run("should produce valid results even though data is shifted in time", func(t *testing.T) {
		scraper, err := series.NewScraper(
			series.ScraperConfig{
				StartTime:      time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
				ScrapeInterval: 15 * time.Second,
			},
			series.WithScraperIterationCountLimit(100), // It's good practice to set an upper bound in tests
		)
		require.NoError(t, err)

		lsDataIterator, err := discrete.NewLinearSegmentDataIterator(discrete.LinearSegmentDataIteratorOptions{
			AmplitudeStart:      20,
			AmplitudeEnd:        40,
			IterationCountLimit: 9,
		})
		require.NoError(t, err)

		type resultContainer struct {
			scrapeInfo   series.ScrapeInfo
			scrapeResult series.ScrapeResult
		}

		var results []resultContainer
		scrapeHandler := func(scrapeInfo series.ScrapeInfo, scrapeResult series.ScrapeResult) error {
			results = append(results, resultContainer{
				scrapeInfo:   scrapeInfo,
				scrapeResult: scrapeResult,
			})
			return nil
		}

		skipNTimes := 25
		skipCount := 0
		for iter := scraper.Iterator(); iter.HasNext(); {
			scrapeInfo := iter.Next()

			if skipCount < skipNTimes {
				skipCount++
				continue
			}

			scrapeResult := lsDataIterator.Evaluate(scrapeInfo)
			if scrapeResult.Exhausted {
				// exhausted time series samples
				break
			}

			err := scrapeHandler(scrapeInfo, scrapeResult)
			require.NoError(t, err)
		}

		for _, r := range results {
			t.Logf("[%3d] Timestamp: %s - Value: %.2f\n",
				r.scrapeInfo.IterationCount,
				r.scrapeInfo.IterationTime,
				r.scrapeResult.Value,
			)
		}

		require.Equal(t, 9, len(results))
		assert.Equal(t, 25, results[0].scrapeInfo.IterationCount)
		assert.Equal(t,
			time.Date(2023, 1, 1, 10, 36, 15, 0, time.UTC),
			results[0].scrapeInfo.IterationTime,
		)
		assert.InDelta(t, 20, results[0].scrapeResult.Value, 0.001)

		assert.Equal(t, 33, results[8].scrapeInfo.IterationCount)
		assert.Equal(t,
			time.Date(2023, 1, 1, 10, 38, 15, 0, time.UTC),
			results[8].scrapeInfo.IterationTime,
		)
		assert.InDelta(t, 40, results[8].scrapeResult.Value, 0.001)
	})
}
