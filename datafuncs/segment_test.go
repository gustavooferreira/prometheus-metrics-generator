package datafuncs_test

import (
	"gustavooferreira/prometheus-metrics-generator/datafuncs"
	"gustavooferreira/prometheus-metrics-generator/series"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLinearSegmentDataIterator(t *testing.T) {
	t.Run("should produce valid values for given iteration count", func(t *testing.T) {
		scraper, err := series.NewScraper(
			series.ScraperConfig{
				StartTime:      time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
				ScrapeInterval: 15 * time.Second,
			},
			series.WithScraperIterationCountLimit(100), // It's good practice to set an upper bound in tests
		)
		require.NoError(t, err)

		lsDataIterator, err := datafuncs.LinearSegmentDataIterator(datafuncs.LinearSegmentDataIteratorOptions{
			AmplitudeStart:       11,
			AmplitudeEnd:         20,
			LengthIterationCount: 10,
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

		require.Equal(t, 10, len(results))
		assert.InDelta(t, 11, results[0].scrapeResult.Value, 0.001)
		assert.InDelta(t, 20, results[9].scrapeResult.Value, 0.001)
	})

	t.Run("should produce valid values for given length duration", func(t *testing.T) {
		scraper, err := series.NewScraper(
			series.ScraperConfig{
				StartTime:      time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
				ScrapeInterval: 15 * time.Second,
			},
			series.WithScraperIterationCountLimit(100), // It's good practice to set an upper bound in tests
		)
		require.NoError(t, err)

		lsDataIterator, err := datafuncs.LinearSegmentDataIterator(datafuncs.LinearSegmentDataIteratorOptions{
			AmplitudeStart: 11,
			AmplitudeEnd:   20,
			LengthDuration: time.Minute,
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
		assert.InDelta(t, 11, results[0].scrapeResult.Value, 0.001)
		assert.InDelta(t, 17.75, results[3].scrapeResult.Value, 0.001)
	})

	t.Run("should produce valid values for a horizontal line", func(t *testing.T) {
		scraper, err := series.NewScraper(
			series.ScraperConfig{
				StartTime:      time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
				ScrapeInterval: 15 * time.Second,
			},
			series.WithScraperIterationCountLimit(100), // It's good practice to set an upper bound in tests
		)
		require.NoError(t, err)

		lsDataIterator, err := datafuncs.LinearSegmentDataIterator(datafuncs.LinearSegmentDataIteratorOptions{
			AmplitudeStart: 50,
			AmplitudeEnd:   50,
			LengthDuration: 2 * time.Minute,
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

		require.Equal(t, 8, len(results))
		assert.InDelta(t, 50, results[0].scrapeResult.Value, 0.001)
		assert.InDelta(t, 50, results[7].scrapeResult.Value, 0.001)
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

		lsDataIterator, err := datafuncs.LinearSegmentDataIterator(datafuncs.LinearSegmentDataIteratorOptions{
			AmplitudeStart: 20,
			AmplitudeEnd:   40,
			LengthDuration: time.Minute,
		})
		require.NoError(t, err)

		type resultContainer struct {
			scrapeInfo series.ScrapeInfo
			value      float64
		}

		var results []resultContainer
		scrapeHandler := func(scrapeInfo series.ScrapeInfo, scrapeResult series.ScrapeResult) error {
			results = append(results, resultContainer{
				scrapeInfo: scrapeInfo,
				value:      scrapeResult.Value,
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
				r.value,
			)
		}

		require.Equal(t, 4, len(results))
		assert.InDelta(t, 20, results[0].value, 0.001)
		assert.Equal(t, 25, results[0].scrapeInfo.IterationCount)
		assert.InDelta(t, 35, results[3].value, 0.001)
		assert.Equal(t, 28, results[3].scrapeInfo.IterationCount)
	})
}
