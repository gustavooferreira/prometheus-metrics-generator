package datafuncs_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gustavooferreira/prometheus-metrics-generator/datafuncs"
	"gustavooferreira/prometheus-metrics-generator/series"
)

func TestCustomValuesDataIterator(t *testing.T) {
	t.Run("should produce valid values for given input array", func(t *testing.T) {
		scraper, err := series.NewScraper(
			series.ScraperConfig{
				StartTime:      time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
				ScrapeInterval: 15 * time.Second,
			},
			series.WithScraperIterationCountLimit(100), // It's good practice to set an upper bound in tests
		)
		require.NoError(t, err)

		lsDataIterator := datafuncs.CustomValuesDataIterator([]datafuncs.CustomValue{
			{Value: 1.0},
			{Value: 5.0},
			{Value: 9.0},
			{Value: 3.0},
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

		err = scraper.Scrape(lsDataIterator, scrapeHandler)
		require.NoError(t, err)

		for _, r := range results {
			t.Logf("[%3d] Timestamp: %s - Value: %.2f\n",
				r.scrapeInfo.IterationCount,
				r.scrapeInfo.IterationTime,
				r.scrapeResult.Value,
			)
		}

		require.Equal(t, 4, len(results))
		assert.Equal(t, 1.0, results[0].scrapeResult.Value)
		assert.Equal(t, 5.0, results[1].scrapeResult.Value)
		assert.Equal(t, 9.0, results[2].scrapeResult.Value)
		assert.Equal(t, 3.0, results[3].scrapeResult.Value)
	})

	t.Run("should produce valid values even though data is shifted in time", func(t *testing.T) {
		scraper, err := series.NewScraper(
			series.ScraperConfig{
				StartTime:      time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
				ScrapeInterval: 15 * time.Second,
			},
			series.WithScraperIterationCountLimit(100), // It's good practice to set an upper bound in tests
		)
		require.NoError(t, err)

		lsDataIterator := datafuncs.CustomValuesDataIterator([]datafuncs.CustomValue{
			{Value: 1.0},
			{Value: 5.0},
			{Value: 9.0},
			{Value: 3.0},
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

		// ----------------------------

		skipNTimes := 25
		skipCount := 0
		for iter := scraper.Iterator(); iter.HasNext(); {
			scrapeInfo := iter.Next()

			if skipCount < skipNTimes {
				skipCount++
				continue
			}

			scrapeResult := lsDataIterator(scrapeInfo)
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

		require.Equal(t, 4, len(results))

		assert.Equal(t, 25, results[0].scrapeInfo.IterationCount)
		assert.Equal(t, 1.0, results[0].scrapeResult.Value)
		assert.Equal(t,
			time.Date(2023, 1, 1, 10, 36, 15, 0, time.UTC),
			results[0].scrapeInfo.IterationTime,
		)

		assert.Equal(t, 26, results[1].scrapeInfo.IterationCount)
		assert.Equal(t, 5.0, results[1].scrapeResult.Value)
		assert.Equal(t,
			time.Date(2023, 1, 1, 10, 36, 30, 0, time.UTC),
			results[1].scrapeInfo.IterationTime,
		)

		assert.Equal(t, 27, results[2].scrapeInfo.IterationCount)
		assert.Equal(t, 9.0, results[2].scrapeResult.Value)
		assert.Equal(t,
			time.Date(2023, 1, 1, 10, 36, 45, 0, time.UTC),
			results[2].scrapeInfo.IterationTime,
		)

		assert.Equal(t, 28, results[3].scrapeInfo.IterationCount)
		assert.Equal(t, 3.0, results[3].scrapeResult.Value)
		assert.Equal(t,
			time.Date(2023, 1, 1, 10, 37, 0, 0, time.UTC),
			results[3].scrapeInfo.IterationTime,
		)
	})
}
