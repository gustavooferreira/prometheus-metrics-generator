package metrics_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gustavooferreira/prometheus-metrics-generator/metrics"
)

func TestScraper(t *testing.T) {
	t.Run("should return a scraper that runs forever", func(t *testing.T) {
		scraper, err := metrics.NewScraper(
			metrics.ScraperConfig{
				StartTime:      time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
				ScrapeInterval: 15 * time.Second,
			},
		)
		require.NoError(t, err)
		assert.True(t, scraper.IsInfinite())

		// should generate a big number of scrapes given that there is no stop condition

		// N is an arbitrary "big" number of scrapes to test that the scraper keeps generating scrapes.
		N := 10000
		scrapesCount := 0

		iter := scraper.Iterator()
		for _, ok := iter.Next(); ok; _, ok = iter.Next() {
			if scrapesCount == N {
				break
			}

			scrapesCount++
		}

		assert.Equal(t, N, scrapesCount)
	})

	t.Run("should produce the right amount of scrapes given the presence of endTime", func(t *testing.T) {
		scraper, err := metrics.NewScraper(
			metrics.ScraperConfig{
				StartTime:      time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
				ScrapeInterval: 15 * time.Second,
			},
			metrics.WithScraperEndTime(time.Date(2023, 1, 1, 10, 31, 0, 0, time.UTC)),
		)
		require.NoError(t, err)

		var scrapeInfoArr []metrics.ScrapeInfo

		iter := scraper.Iterator()
		for scrapeInfo, ok := iter.Next(); ok; scrapeInfo, ok = iter.Next() {
			scrapeInfoArr = append(scrapeInfoArr, scrapeInfo)
		}

		for _, r := range scrapeInfoArr {
			t.Logf("[%3d] Timestamp: %s",
				r.IterationIndex,
				r.IterationTime,
			)
		}

		require.Equal(t, 5, len(scrapeInfoArr))
		assert.Equal(t, metrics.ScrapeInfo{
			FirstIterationTime: time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			IterationIndex:     0,
			IterationTime:      time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
		}, scrapeInfoArr[0])
		assert.Equal(t, metrics.ScrapeInfo{
			FirstIterationTime: time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			IterationIndex:     1,
			IterationTime:      time.Date(2023, 1, 1, 10, 30, 15, 0, time.UTC),
		}, scrapeInfoArr[1])
		assert.Equal(t, metrics.ScrapeInfo{
			FirstIterationTime: time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			IterationIndex:     2,
			IterationTime:      time.Date(2023, 1, 1, 10, 30, 30, 0, time.UTC),
		}, scrapeInfoArr[2])
		assert.Equal(t, metrics.ScrapeInfo{
			FirstIterationTime: time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			IterationIndex:     3,
			IterationTime:      time.Date(2023, 1, 1, 10, 30, 45, 0, time.UTC),
		}, scrapeInfoArr[3])
		assert.Equal(t, metrics.ScrapeInfo{
			FirstIterationTime: time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			IterationIndex:     4,
			IterationTime:      time.Date(2023, 1, 1, 10, 31, 0, 0, time.UTC),
		}, scrapeInfoArr[4])
	})

	t.Run("should produce one scrape given that the end time is the same as the start time", func(t *testing.T) {
		scraper, err := metrics.NewScraper(
			metrics.ScraperConfig{
				StartTime:      time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
				ScrapeInterval: 15 * time.Second,
			},
			metrics.WithScraperEndTime(time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC)),
		)
		require.NoError(t, err)

		var scrapeInfoArr []metrics.ScrapeInfo

		iter := scraper.Iterator()
		for scrapeInfo, ok := iter.Next(); ok; scrapeInfo, ok = iter.Next() {
			scrapeInfoArr = append(scrapeInfoArr, scrapeInfo)
		}

		require.Equal(t, 1, len(scrapeInfoArr))
		assert.Equal(t, metrics.ScrapeInfo{
			FirstIterationTime: time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			IterationIndex:     0,
			IterationTime:      time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
		}, scrapeInfoArr[0])
	})

	t.Run("should produce the right amount of scrapes given the presence of iterationCountLimit", func(t *testing.T) {
		scraper, err := metrics.NewScraper(
			metrics.ScraperConfig{
				StartTime:      time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
				ScrapeInterval: 15 * time.Second,
			},
			metrics.WithScraperIterationCountLimit(4),
		)
		require.NoError(t, err)

		var scrapeInfoArr []metrics.ScrapeInfo

		iter := scraper.Iterator()
		for scrapeInfo, ok := iter.Next(); ok; scrapeInfo, ok = iter.Next() {
			scrapeInfoArr = append(scrapeInfoArr, scrapeInfo)
		}

		for _, r := range scrapeInfoArr {
			t.Logf("[%3d] Timestamp: %s",
				r.IterationIndex,
				r.IterationTime,
			)
		}

		require.Equal(t, 4, len(scrapeInfoArr))
		assert.Equal(t, metrics.ScrapeInfo{
			FirstIterationTime: time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			IterationIndex:     0,
			IterationTime:      time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
		}, scrapeInfoArr[0])
		assert.Equal(t, metrics.ScrapeInfo{
			FirstIterationTime: time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			IterationIndex:     1,
			IterationTime:      time.Date(2023, 1, 1, 10, 30, 15, 0, time.UTC),
		}, scrapeInfoArr[1])
		assert.Equal(t, metrics.ScrapeInfo{
			FirstIterationTime: time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			IterationIndex:     2,
			IterationTime:      time.Date(2023, 1, 1, 10, 30, 30, 0, time.UTC),
		}, scrapeInfoArr[2])
		assert.Equal(t, metrics.ScrapeInfo{
			FirstIterationTime: time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			IterationIndex:     3,
			IterationTime:      time.Date(2023, 1, 1, 10, 30, 45, 0, time.UTC),
		}, scrapeInfoArr[3])
	})

	t.Run("should be able to go over the iterator twice given that we reset the iterator", func(t *testing.T) {
		scraper, err := metrics.NewScraper(
			metrics.ScraperConfig{
				StartTime:      time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
				ScrapeInterval: 15 * time.Second,
			},
			metrics.WithScraperEndTime(time.Date(2023, 1, 1, 10, 31, 0, 0, time.UTC)),
		)
		require.NoError(t, err)

		iter := scraper.Iterator()

		var scrapeInfoArr1 []metrics.ScrapeInfo
		for scrapeInfo, ok := iter.Next(); ok; scrapeInfo, ok = iter.Next() {
			scrapeInfoArr1 = append(scrapeInfoArr1, scrapeInfo)
		}

		iter.Reset()

		var scrapeInfoArr2 []metrics.ScrapeInfo
		for scrapeInfo, ok := iter.Next(); ok; scrapeInfo, ok = iter.Next() {
			scrapeInfoArr2 = append(scrapeInfoArr2, scrapeInfo)
		}

		assert.Equal(t, scrapeInfoArr1, scrapeInfoArr2)
		require.Equal(t, 5, len(scrapeInfoArr1))
	})

	t.Run("the method HasNext should always report correctly whether we have another scrape available when provided with iteration count limit", func(t *testing.T) {
		scraper, err := metrics.NewScraper(
			metrics.ScraperConfig{
				StartTime:      time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
				ScrapeInterval: 15 * time.Second,
			},
			metrics.WithScraperIterationCountLimit(4),
		)
		require.NoError(t, err)

		var scrapeInfoArr []metrics.ScrapeInfo

		iter := scraper.Iterator()

		for iter.HasNext() {
			scrapeInfo, ok := iter.Next()
			require.True(t, ok)

			scrapeInfoArr = append(scrapeInfoArr, scrapeInfo)
		}

		for _, r := range scrapeInfoArr {
			t.Logf("[%3d] Timestamp: %s",
				r.IterationIndex,
				r.IterationTime,
			)
		}

		require.Equal(t, 4, len(scrapeInfoArr))
		assert.Equal(t, metrics.ScrapeInfo{
			FirstIterationTime: time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			IterationIndex:     0,
			IterationTime:      time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
		}, scrapeInfoArr[0])
		assert.Equal(t, metrics.ScrapeInfo{
			FirstIterationTime: time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			IterationIndex:     1,
			IterationTime:      time.Date(2023, 1, 1, 10, 30, 15, 0, time.UTC),
		}, scrapeInfoArr[1])
		assert.Equal(t, metrics.ScrapeInfo{
			FirstIterationTime: time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			IterationIndex:     2,
			IterationTime:      time.Date(2023, 1, 1, 10, 30, 30, 0, time.UTC),
		}, scrapeInfoArr[2])
		assert.Equal(t, metrics.ScrapeInfo{
			FirstIterationTime: time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			IterationIndex:     3,
			IterationTime:      time.Date(2023, 1, 1, 10, 30, 45, 0, time.UTC),
		}, scrapeInfoArr[3])
	})

	t.Run("the method HasNext should always report correctly whether we have another scrape available when provided with end time", func(t *testing.T) {
		scraper, err := metrics.NewScraper(
			metrics.ScraperConfig{
				StartTime:      time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
				ScrapeInterval: 15 * time.Second,
			},
			metrics.WithScraperEndTime(time.Date(2023, 1, 1, 10, 31, 0, 0, time.UTC)),
		)
		require.NoError(t, err)

		var scrapeInfoArr []metrics.ScrapeInfo

		iter := scraper.Iterator()

		for iter.HasNext() {
			scrapeInfo, ok := iter.Next()
			require.True(t, ok)

			scrapeInfoArr = append(scrapeInfoArr, scrapeInfo)
		}

		for _, r := range scrapeInfoArr {
			t.Logf("[%3d] Timestamp: %s",
				r.IterationIndex,
				r.IterationTime,
			)
		}

		require.Equal(t, 5, len(scrapeInfoArr))
		assert.Equal(t, metrics.ScrapeInfo{
			FirstIterationTime: time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			IterationIndex:     0,
			IterationTime:      time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
		}, scrapeInfoArr[0])
		assert.Equal(t, metrics.ScrapeInfo{
			FirstIterationTime: time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			IterationIndex:     1,
			IterationTime:      time.Date(2023, 1, 1, 10, 30, 15, 0, time.UTC),
		}, scrapeInfoArr[1])
		assert.Equal(t, metrics.ScrapeInfo{
			FirstIterationTime: time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			IterationIndex:     2,
			IterationTime:      time.Date(2023, 1, 1, 10, 30, 30, 0, time.UTC),
		}, scrapeInfoArr[2])
		assert.Equal(t, metrics.ScrapeInfo{
			FirstIterationTime: time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			IterationIndex:     3,
			IterationTime:      time.Date(2023, 1, 1, 10, 30, 45, 0, time.UTC),
		}, scrapeInfoArr[3])
		assert.Equal(t, metrics.ScrapeInfo{
			FirstIterationTime: time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			IterationIndex:     4,
			IterationTime:      time.Date(2023, 1, 1, 10, 31, 0, 0, time.UTC),
		}, scrapeInfoArr[4])
	})

	t.Run("should correctly iterator through all the possible values using the ScrapeDataIterator method", func(t *testing.T) {
		scraper, err := metrics.NewScraper(
			metrics.ScraperConfig{
				StartTime:      time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
				ScrapeInterval: 15 * time.Second,
			},
			metrics.WithScraperIterationCountLimit(4),
		)
		require.NoError(t, err)

		valueCount := 0.0
		dataIteratorFunc := func(scrapeInfo metrics.ScrapeInfo) metrics.ScrapeResult {
			valueCount++

			return metrics.ScrapeResult{
				Value:     valueCount,
				Missing:   false,
				Exhausted: false,
			}
		}

		var scrapeInfoArr []metrics.ScrapeInfo
		scrapeHandler := func(scrapeInfo metrics.ScrapeInfo, scrapeResult metrics.ScrapeResult) error {
			scrapeInfoArr = append(scrapeInfoArr, scrapeInfo)
			return nil
		}

		err = scraper.ScrapeDataIterator(metrics.DataIteratorFunc(dataIteratorFunc), scrapeHandler)
		require.NoError(t, err)

		for _, r := range scrapeInfoArr {
			t.Logf("[%3d] Timestamp: %s",
				r.IterationIndex,
				r.IterationTime,
			)
		}

		require.Equal(t, 4, len(scrapeInfoArr))
		assert.Equal(t, metrics.ScrapeInfo{
			FirstIterationTime: time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			IterationIndex:     0,
			IterationTime:      time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
		}, scrapeInfoArr[0])
		assert.Equal(t, metrics.ScrapeInfo{
			FirstIterationTime: time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			IterationIndex:     1,
			IterationTime:      time.Date(2023, 1, 1, 10, 30, 15, 0, time.UTC),
		}, scrapeInfoArr[1])
		assert.Equal(t, metrics.ScrapeInfo{
			FirstIterationTime: time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			IterationIndex:     2,
			IterationTime:      time.Date(2023, 1, 1, 10, 30, 30, 0, time.UTC),
		}, scrapeInfoArr[2])
		assert.Equal(t, metrics.ScrapeInfo{
			FirstIterationTime: time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			IterationIndex:     3,
			IterationTime:      time.Date(2023, 1, 1, 10, 30, 45, 0, time.UTC),
		}, scrapeInfoArr[3])
	})
}
