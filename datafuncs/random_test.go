package datafuncs_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gustavooferreira/prometheus-metrics-generator/datafuncs"
	"gustavooferreira/prometheus-metrics-generator/series"
)

func TestRandomDataIterator(t *testing.T) {
	t.Run("should fail given that neither the LengthDuration nor LengthIterationCount were set", func(t *testing.T) {
		_, err := datafuncs.RandomDataIterator(datafuncs.RandomDataIteratorOptions{
			AmplitudeMin: 11,
			AmplitudeMax: 20,
		})
		require.Error(t, err)
		expectedErrorMessage := "stop condition needs to be provided, either set the length duration or length " +
			"iteration count"
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("should fail given that both the LengthDuration and LengthIterationCount options were set", func(t *testing.T) {
		_, err := datafuncs.RandomDataIterator(datafuncs.RandomDataIteratorOptions{
			AmplitudeMin:         11,
			AmplitudeMax:         20,
			LengthDuration:       time.Second,
			LengthIterationCount: 10,
		})
		require.Error(t, err)
		expectedErrorMessage := "only one stop condition should be provided"
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("should fail given that LengthDuration is negative", func(t *testing.T) {
		_, err := datafuncs.RandomDataIterator(datafuncs.RandomDataIteratorOptions{
			AmplitudeMin:   11,
			AmplitudeMax:   20,
			LengthDuration: -1 * time.Second,
		})
		require.Error(t, err)
		expectedErrorMessage := "length duration cannot be negative"
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("should fail given that LengthIterationCount is negative", func(t *testing.T) {
		_, err := datafuncs.RandomDataIterator(datafuncs.RandomDataIteratorOptions{
			AmplitudeMin:         11,
			AmplitudeMax:         20,
			LengthIterationCount: -10,
		})
		require.Error(t, err)
		expectedErrorMessage := "length iteration count cannot be negative"
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("should fail given that LengthDurationExclusive was set without setting LengthDuration", func(t *testing.T) {
		_, err := datafuncs.RandomDataIterator(datafuncs.RandomDataIteratorOptions{
			AmplitudeMin:            11,
			AmplitudeMax:            20,
			LengthIterationCount:    10,
			LengthDurationExclusive: true,
		})
		require.Error(t, err)
		expectedErrorMessage := "length duration exclusive option applies to length duration option only, but length " +
			"duration is not set"
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("should produce valid values for given iteration count", func(t *testing.T) {
		scraper, err := series.NewScraper(
			series.ScraperConfig{
				StartTime:      time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
				ScrapeInterval: 15 * time.Second,
			},
			series.WithScraperIterationCountLimit(100), // It's good practice to set an upper bound in tests
		)
		require.NoError(t, err)

		lsDataIterator, err := datafuncs.RandomDataIterator(datafuncs.RandomDataIteratorOptions{
			AmplitudeMin:         11,
			AmplitudeMax:         20,
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

		// assert all values are within amplitude range provided
		for _, result := range results {
			assert.GreaterOrEqual(t, 20.0, result.scrapeResult.Value)
			assert.LessOrEqual(t, 11.0, result.scrapeResult.Value)
		}
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

		lsDataIterator, err := datafuncs.RandomDataIterator(datafuncs.RandomDataIteratorOptions{
			AmplitudeMin:   11,
			AmplitudeMax:   20,
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

		require.Equal(t, 5, len(results))

		// assert all values are within amplitude range provided
		for _, result := range results {
			assert.GreaterOrEqual(t, 20.0, result.scrapeResult.Value)
			assert.LessOrEqual(t, 11.0, result.scrapeResult.Value)
		}
	})

	t.Run("should produce valid values for given length duration with exclusive option set", func(t *testing.T) {
		scraper, err := series.NewScraper(
			series.ScraperConfig{
				StartTime:      time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
				ScrapeInterval: 15 * time.Second,
			},
			series.WithScraperIterationCountLimit(100), // It's good practice to set an upper bound in tests
		)
		require.NoError(t, err)

		lsDataIterator, err := datafuncs.RandomDataIterator(datafuncs.RandomDataIteratorOptions{
			AmplitudeMin:            11,
			AmplitudeMax:            20,
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
			t.Logf("[%3d] Timestamp: %s - Value: %.2f\n",
				r.scrapeInfo.IterationCount,
				r.scrapeInfo.IterationTime,
				r.scrapeResult.Value,
			)
		}

		require.Equal(t, 4, len(results))

		// assert all values are within amplitude range provided
		for _, result := range results {
			assert.GreaterOrEqual(t, 20.0, result.scrapeResult.Value)
			assert.LessOrEqual(t, 11.0, result.scrapeResult.Value)
		}
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

		lsDataIterator, err := datafuncs.RandomDataIterator(datafuncs.RandomDataIteratorOptions{
			AmplitudeMin:   20,
			AmplitudeMax:   40,
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

		require.Equal(t, 5, len(results))
		assert.Equal(t, 25, results[0].scrapeInfo.IterationCount)
		assert.Equal(t,
			time.Date(2023, 1, 1, 10, 36, 15, 0, time.UTC),
			results[0].scrapeInfo.IterationTime,
		)

		assert.Equal(t, 29, results[4].scrapeInfo.IterationCount)
		assert.Equal(t,
			time.Date(2023, 1, 1, 10, 37, 15, 0, time.UTC),
			results[4].scrapeInfo.IterationTime,
		)

		// assert all values are within amplitude range provided
		for _, result := range results {
			assert.GreaterOrEqual(t, 40.0, result.scrapeResult.Value)
			assert.LessOrEqual(t, 20.0, result.scrapeResult.Value)
		}

	})
}
