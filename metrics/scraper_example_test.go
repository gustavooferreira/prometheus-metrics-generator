package metrics_test

import (
	"fmt"
	"time"

	"github.com/gustavooferreira/prometheus-metrics-generator/discrete"
	"github.com/gustavooferreira/prometheus-metrics-generator/metrics"
)

func ExampleScraper() {
	scraper, err := metrics.NewScraper(
		metrics.ScraperConfig{
			StartTime:      time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			ScrapeInterval: 15 * time.Second,
		},
		metrics.WithScraperIterationCountLimit(100), // It's good practice to set an upper bound in tests
	)
	if err != nil {
		panic(err)
	}

	lsDataGenerator, err := discrete.NewLinearSegmentDataGenerator(discrete.LinearSegmentDataGeneratorOptions{
		AmplitudeStart:      50,
		AmplitudeEnd:        70,
		IterationCountLimit: 5,
	})
	if err != nil {
		panic(err)
	}

	scrapeHandler := func(scrapeInfo metrics.ScrapeInfo, scrapeResult metrics.ScrapeResult) error {
		fmt.Printf("[%3d] Timestamp: %s - Value: %.2f\n",
			scrapeInfo.IterationCount,
			scrapeInfo.IterationTime,
			scrapeResult.Value,
		)
		return nil
	}

	err = scraper.ScrapeDataIterator(lsDataGenerator.Iterator(), scrapeHandler)
	if err != nil {
		panic(err)
	}
	// Output:
	// [  0] Timestamp: 2023-01-01 10:30:00 +0000 UTC - Value: 50.00
	// [  1] Timestamp: 2023-01-01 10:30:15 +0000 UTC - Value: 55.00
	// [  2] Timestamp: 2023-01-01 10:30:30 +0000 UTC - Value: 60.00
	// [  3] Timestamp: 2023-01-01 10:30:45 +0000 UTC - Value: 65.00
	// [  4] Timestamp: 2023-01-01 10:31:00 +0000 UTC - Value: 70.00
}
