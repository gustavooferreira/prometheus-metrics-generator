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
	})

	t.Run("should generate a big number of scrapes given that there is no stop condition", func(t *testing.T) {
		scraper, err := metrics.NewScraper(
			metrics.ScraperConfig{
				StartTime:      time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
				ScrapeInterval: 15 * time.Second,
			},
		)
		require.NoError(t, err)

		// N is an arbitrary "big" number of scrapes to test that the scraper keeps generating scrapes.
		N := 10000
		scrapesCount := 0

		for iter := scraper.Iterator(); iter.HasNext(); {
			_ = iter.Next()

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

		for iter := scraper.Iterator(); iter.HasNext(); {
			scrapeInfo := iter.Next()
			scrapeInfoArr = append(scrapeInfoArr, scrapeInfo)
		}

		for _, r := range scrapeInfoArr {
			t.Logf("[%3d] Timestamp: %s\n",
				r.IterationCount,
				r.IterationTime,
			)
		}

		require.Equal(t, 5, len(scrapeInfoArr))
		assert.Equal(t, metrics.ScrapeInfo{
			FirstIterationTime: time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			IterationCount:     0,
			IterationTime:      time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
		}, scrapeInfoArr[0])
		assert.Equal(t, metrics.ScrapeInfo{
			FirstIterationTime: time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			IterationCount:     1,
			IterationTime:      time.Date(2023, 1, 1, 10, 30, 15, 0, time.UTC),
		}, scrapeInfoArr[1])
		assert.Equal(t, metrics.ScrapeInfo{
			FirstIterationTime: time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			IterationCount:     2,
			IterationTime:      time.Date(2023, 1, 1, 10, 30, 30, 0, time.UTC),
		}, scrapeInfoArr[2])
		assert.Equal(t, metrics.ScrapeInfo{
			FirstIterationTime: time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			IterationCount:     3,
			IterationTime:      time.Date(2023, 1, 1, 10, 30, 45, 0, time.UTC),
		}, scrapeInfoArr[3])
		assert.Equal(t, metrics.ScrapeInfo{
			FirstIterationTime: time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			IterationCount:     4,
			IterationTime:      time.Date(2023, 1, 1, 10, 31, 0, 0, time.UTC),
		}, scrapeInfoArr[4])
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
		for iter.HasNext() {
			scrapeInfo := iter.Next()
			scrapeInfoArr1 = append(scrapeInfoArr1, scrapeInfo)
		}

		iter.Reset()

		var scrapeInfoArr2 []metrics.ScrapeInfo
		for iter.HasNext() {
			scrapeInfo := iter.Next()
			scrapeInfoArr2 = append(scrapeInfoArr2, scrapeInfo)
		}

		assert.Equal(t, scrapeInfoArr1, scrapeInfoArr2)
	})
}
