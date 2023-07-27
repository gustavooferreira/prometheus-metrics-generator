package discrete_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/gustavooferreira/prometheus-metrics-generator/series"
)

// helperScraper computes the data function results given a DataIterator.
func helperScraper(t *testing.T, dataIterator series.DataIterator) []resultContainer {
	t.Helper()

	scraper, err := series.NewScraper(
		series.ScraperConfig{
			StartTime:      time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			ScrapeInterval: 15 * time.Second,
		},
		series.WithScraperIterationCountLimit(100), // It's good practice to set an upper bound in tests
	)
	require.NoError(t, err)

	var results []resultContainer
	scrapeHandler := func(scrapeInfo series.ScrapeInfo, scrapeResult series.ScrapeResult) error {
		results = append(results, resultContainer{
			scrapeInfo:   scrapeInfo,
			scrapeResult: scrapeResult,
		})
		return nil
	}

	err = scraper.Scrape(dataIterator, scrapeHandler)
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
	scrapeInfo   series.ScrapeInfo
	scrapeResult series.ScrapeResult
}
