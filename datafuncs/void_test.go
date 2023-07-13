package datafuncs_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gustavooferreira/prometheus-metrics-generator/datafuncs"
	"github.com/gustavooferreira/prometheus-metrics-generator/series"
)

func TestVoidDataIterator(t *testing.T) {
	t.Run("should fail given that neither the LengthDuration nor LengthIterationCount were set", func(t *testing.T) {
		_, err := datafuncs.VoidDataIterator(datafuncs.VoidDataIteratorOptions{})
		require.Error(t, err)
		expectedErrorMessage := "stop condition needs to be provided, either set the length duration or length " +
			"iteration count"
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("should fail given that both the LengthDuration and LengthIterationCount options were set", func(t *testing.T) {
		_, err := datafuncs.VoidDataIterator(datafuncs.VoidDataIteratorOptions{
			LengthDuration:       time.Second,
			LengthIterationCount: 10,
		})
		require.Error(t, err)
		expectedErrorMessage := "only one stop condition should be provided"
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("should fail given that LengthDuration is negative", func(t *testing.T) {
		_, err := datafuncs.VoidDataIterator(datafuncs.VoidDataIteratorOptions{
			LengthDuration: -1 * time.Second,
		})
		require.Error(t, err)
		expectedErrorMessage := "length duration cannot be negative"
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("should fail given that LengthIterationCount is negative", func(t *testing.T) {
		_, err := datafuncs.VoidDataIterator(datafuncs.VoidDataIteratorOptions{
			LengthIterationCount: -10,
		})
		require.Error(t, err)
		expectedErrorMessage := "length iteration count cannot be negative"
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("should fail given that LengthDurationExclusive was set without setting LengthDuration", func(t *testing.T) {
		_, err := datafuncs.VoidDataIterator(datafuncs.VoidDataIteratorOptions{
			LengthIterationCount:    10,
			LengthDurationExclusive: true,
		})
		require.Error(t, err)
		expectedErrorMessage := "length duration exclusive option applies to length duration option only, but length " +
			"duration is not set"
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("should produce valid result for given iteration count", func(t *testing.T) {
		scraper, err := series.NewScraper(
			series.ScraperConfig{
				StartTime:      time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
				ScrapeInterval: 15 * time.Second,
			},
			series.WithScraperIterationCountLimit(100), // It's good practice to set an upper bound in tests
		)
		require.NoError(t, err)

		lsDataIterator, err := datafuncs.VoidDataIterator(datafuncs.VoidDataIteratorOptions{
			LengthIterationCount: 5,
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
			t.Logf("[%3d] Timestamp: %s - Missing: %t\n",
				r.scrapeInfo.IterationCount,
				r.scrapeInfo.IterationTime,
				r.scrapeResult.Missing,
			)
		}

		require.Equal(t, 5, len(results))
		assert.True(t, results[0].scrapeResult.Missing)
		assert.True(t, results[4].scrapeResult.Missing)
	})

	t.Run("should produce valid results for given length duration", func(t *testing.T) {
		scraper, err := series.NewScraper(
			series.ScraperConfig{
				StartTime:      time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
				ScrapeInterval: 15 * time.Second,
			},
			series.WithScraperIterationCountLimit(100), // It's good practice to set an upper bound in tests
		)
		require.NoError(t, err)

		lsDataIterator, err := datafuncs.VoidDataIterator(datafuncs.VoidDataIteratorOptions{
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
			t.Logf("[%3d] Timestamp: %s - Missing: %t\n",
				r.scrapeInfo.IterationCount,
				r.scrapeInfo.IterationTime,
				r.scrapeResult.Missing,
			)
		}

		require.Equal(t, 5, len(results))
		assert.True(t, results[0].scrapeResult.Missing)
		assert.True(t, results[4].scrapeResult.Missing)
	})

	t.Run("should produce valid results for given length duration with exclusive option set", func(t *testing.T) {
		scraper, err := series.NewScraper(
			series.ScraperConfig{
				StartTime:      time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
				ScrapeInterval: 15 * time.Second,
			},
			series.WithScraperIterationCountLimit(100), // It's good practice to set an upper bound in tests
		)
		require.NoError(t, err)

		lsDataIterator, err := datafuncs.VoidDataIterator(datafuncs.VoidDataIteratorOptions{
			LengthDuration:          time.Minute,
			LengthDurationExclusive: true,
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
			t.Logf("[%3d] Timestamp: %s - Missing: %t\n",
				r.scrapeInfo.IterationCount,
				r.scrapeInfo.IterationTime,
				r.scrapeResult.Missing,
			)
		}

		require.Equal(t, 4, len(results))
		assert.True(t, results[0].scrapeResult.Missing)
		assert.True(t, results[3].scrapeResult.Missing)
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

		lsDataIterator, err := datafuncs.VoidDataIterator(datafuncs.VoidDataIteratorOptions{
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

		// ----------------------------

		skipNTimes := 30
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
			t.Logf("[%3d] Timestamp: %s - Missing: %t\n",
				r.scrapeInfo.IterationCount,
				r.scrapeInfo.IterationTime,
				r.scrapeResult.Missing,
			)
		}

		require.Equal(t, 5, len(results))
		assert.Equal(t, 30, results[0].scrapeInfo.IterationCount)
		assert.Equal(t,
			time.Date(2023, 1, 1, 10, 37, 30, 0, time.UTC),
			results[0].scrapeInfo.IterationTime,
		)
		assert.True(t, results[0].scrapeResult.Missing)

		assert.Equal(t, 34, results[4].scrapeInfo.IterationCount)
		assert.Equal(t,
			time.Date(2023, 1, 1, 10, 38, 30, 0, time.UTC),
			results[4].scrapeInfo.IterationTime,
		)
		assert.True(t, results[4].scrapeResult.Missing)
	})

	t.Run("should return a single data point given the duration length", func(t *testing.T) {
		scraper, err := series.NewScraper(
			series.ScraperConfig{
				StartTime:      time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
				ScrapeInterval: 15 * time.Second,
			},
			series.WithScraperIterationCountLimit(100), // It's good practice to set an upper bound in tests
		)
		require.NoError(t, err)

		lsDataIterator, err := datafuncs.VoidDataIterator(datafuncs.VoidDataIteratorOptions{
			LengthDuration: 10 * time.Second,
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
			t.Logf("[%3d] Timestamp: %s - Missing: %t\n",
				r.scrapeInfo.IterationCount,
				r.scrapeInfo.IterationTime,
				r.scrapeResult.Missing,
			)
		}

		require.Equal(t, 5, len(results))
		assert.True(t, results[0].scrapeResult.Missing)
		assert.True(t, results[4].scrapeResult.Missing)
	})
}
