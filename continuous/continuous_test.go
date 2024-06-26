package continuous_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/gustavooferreira/prometheus-metrics-generator/continuous"
	"github.com/gustavooferreira/prometheus-metrics-generator/metrics"
)

// helperScraper computes the data function results given a DataIterator and a functionStartTime.
func helperScraper(t *testing.T, dataIterator continuous.DataIterator, functionStartTime time.Time) []resultContainer {
	t.Helper()

	scraper, err := metrics.NewScraper(
		metrics.ScraperConfig{
			StartTime:      time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			ScrapeInterval: 15 * time.Second,
		},
		metrics.WithScraperIterationCountLimit(100), // It's good practice to set an upper bound in tests
	)
	require.NoError(t, err)

	var results []resultContainer
	scrapeHandler := func(scrapeInfo metrics.ScrapeInfo, scrapeResult metrics.ScrapeResult) error {
		results = append(results, resultContainer{
			scrapeInfo:   scrapeInfo,
			scrapeResult: scrapeResult,
		})
		return nil
	}

	iter := scraper.Iterator()
	for scrapeInfo, ok := iter.Next(); ok; scrapeInfo, ok = iter.Next() {

		continuousScrapeInfo := continuous.ScrapeInfo{
			FirstIterationTime: scrapeInfo.FirstIterationTime,
			IterationIndex:     scrapeInfo.IterationIndex,
			IterationTime:      scrapeInfo.IterationTime,
			FunctionStartTime:  functionStartTime,
		}

		scrapeResult := dataIterator.Evaluate(continuousScrapeInfo)
		if scrapeResult.Exhausted {
			// exhausted time series samples
			break
		}

		err := scrapeHandler(scrapeInfo, scrapeResult)
		require.NoError(t, err)
	}

	for _, r := range results {
		t.Logf("[%3d] Timestamp: %s - Value: %6.2f - Missing: %t\n",
			r.scrapeInfo.IterationIndex,
			r.scrapeInfo.IterationTime,
			r.scrapeResult.Value,
			r.scrapeResult.Missing,
		)
	}

	return results
}

type resultContainer struct {
	scrapeInfo   metrics.ScrapeInfo
	scrapeResult metrics.ScrapeResult
}
