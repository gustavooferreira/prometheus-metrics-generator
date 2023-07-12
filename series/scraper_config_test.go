package series_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"gustavooferreira/prometheus-metrics-generator/series"
)

func TestScraperConfig(t *testing.T) {
	t.Run("should fail validation check when provided with a negative scrape interval", func(t *testing.T) {
		scraperConfig := series.ScraperConfig{
			StartTime:      time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			ScrapeInterval: -1 * time.Second,
		}

		err := scraperConfig.Validate()
		require.Error(t, err)
		expectedErrorMessage := "scrape interval cannot be less or equal to zero"
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("should fail validation check when provided with an end time before start time", func(t *testing.T) {
		scraperConfig := series.ScraperConfig{
			StartTime:      time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			ScrapeInterval: 1 * time.Second,
		}

		scraperConfig.ApplyFunctionalOptions(
			series.WithScraperEndTime(time.Date(2023, 1, 1, 9, 0, 0, 0, time.UTC)),
		)

		err := scraperConfig.Validate()
		require.Error(t, err)
		expectedErrorMessage := "end time cannot be before start time"
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("should fail validation check when provided with a negative iteration count limit", func(t *testing.T) {
		scraperConfig := series.ScraperConfig{
			StartTime:      time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			ScrapeInterval: 1 * time.Second,
		}

		scraperConfig.ApplyFunctionalOptions(
			series.WithScraperIterationCountLimit(-10),
		)

		err := scraperConfig.Validate()
		require.Error(t, err)
		expectedErrorMessage := "iteration count limit cannot be less than zero"
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("should pass validation check when provided sane values", func(t *testing.T) {
		scraperConfig := series.ScraperConfig{
			StartTime:      time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			ScrapeInterval: 15 * time.Second,
		}

		err := scraperConfig.Validate()
		require.NoError(t, err)
	})

	t.Run("should pass validation check when provided sane values including functional options", func(t *testing.T) {
		scraperConfig := series.ScraperConfig{
			StartTime:      time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			ScrapeInterval: 15 * time.Second,
		}

		scraperConfig.ApplyFunctionalOptions(
			series.WithScraperEndTime(time.Date(2023, 1, 1, 11, 0, 0, 0, time.UTC)),
			series.WithScraperIterationCountLimit(10),
		)

		err := scraperConfig.Validate()
		require.NoError(t, err)

		assert.Equal(t,
			time.Date(2023, 1, 1, 11, 0, 0, 0, time.UTC),
			scraperConfig.GetEndTime())
		assert.Equal(t, 10, scraperConfig.GetIterationCountLimit())
	})
}
