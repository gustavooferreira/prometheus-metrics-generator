package discrete_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/gustavooferreira/prometheus-metrics-generator/metrics"
)

// helperScraper computes the data function results given a DataIterator.
func helperScraper(t *testing.T, dataIterator metrics.DataIterator) []resultContainer {
	t.Helper()

	return helperScraperCustom(
		t,
		time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
		15*time.Second,
		100,
		dataIterator,
	)
}

func helperScraperCustom(t *testing.T, startTime time.Time, scrapeInterval time.Duration, scrapeCountLimit int, dataIterator metrics.DataIterator) []resultContainer {
	t.Helper()

	scraper, err := metrics.NewScraper(
		metrics.ScraperConfig{
			StartTime:      startTime,
			ScrapeInterval: scrapeInterval,
		},
		metrics.WithScraperIterationCountLimit(scrapeCountLimit), // It's good practice to set an upper bound in tests
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

	err = scraper.ScrapeDataIterator(dataIterator, scrapeHandler)
	require.NoError(t, err)

	for _, r := range results {
		t.Logf("[%3d] Timestamp: %s - Value: %6.2f - Missing: %t\n",
			r.scrapeInfo.IterationCount,
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
