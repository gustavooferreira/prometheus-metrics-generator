package discrete_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gustavooferreira/prometheus-metrics-generator/discrete"
	"github.com/gustavooferreira/prometheus-metrics-generator/metrics"
)

func TestCounterTimeSeries(t *testing.T) {
	t.Run("should return valid time series information", func(t *testing.T) {
		dataGenerator := discrete.NewVoidDataGenerator(5)

		timeSeries := discrete.NewMetricTimeSeries(
			"some-series",
			map[string]string{"key": "value"},
			dataGenerator,
			metrics.NewEndStrategyLoop())

		assert.Equal(t, metrics.MetricTypeCounter, timeSeries.Info().Type)
		assert.Equal(t, "some-series", timeSeries.Info().Name)
		assert.Equal(t, map[string]string{"key": "value"}, timeSeries.Info().Labels)
	})

	t.Run("should return a looped series given the loop end strategy", func(t *testing.T) {
		dataGenerator, err := discrete.NewLinearSegmentDataGenerator(discrete.LinearSegmentDataGeneratorOptions{
			AmplitudeStart:      50,
			AmplitudeEnd:        70,
			IterationCountLimit: 3,
		})
		require.NoError(t, err)

		timeSeries := discrete.NewMetricTimeSeries(
			"some-series",
			map[string]string{"key": "value"},
			dataGenerator,
			metrics.NewEndStrategyLoop())

		results := helperScraperCustom(
			t,
			time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			15*time.Second,
			9,
			timeSeries.Iterator(),
		)

		require.Equal(t, 9, len(results))
		assert.InDelta(t, 50, results[0].scrapeResult.Value, 0.001)
		assert.InDelta(t, 60, results[1].scrapeResult.Value, 0.001)
		assert.InDelta(t, 70, results[2].scrapeResult.Value, 0.001)
		assert.InDelta(t, 50, results[3].scrapeResult.Value, 0.001)
		assert.InDelta(t, 60, results[4].scrapeResult.Value, 0.001)
		assert.InDelta(t, 70, results[5].scrapeResult.Value, 0.001)
		assert.InDelta(t, 50, results[6].scrapeResult.Value, 0.001)
		assert.InDelta(t, 60, results[7].scrapeResult.Value, 0.001)
		assert.InDelta(t, 70, results[8].scrapeResult.Value, 0.001)
	})

	t.Run("should return a single run of the iterator followed by last value given the end strategy", func(t *testing.T) {
		dataGenerator, err := discrete.NewLinearSegmentDataGenerator(discrete.LinearSegmentDataGeneratorOptions{
			AmplitudeStart:      50,
			AmplitudeEnd:        70,
			IterationCountLimit: 3,
		})
		require.NoError(t, err)

		timeSeries := discrete.NewMetricTimeSeries(
			"some-series",
			map[string]string{"key": "value"},
			dataGenerator,
			metrics.NewEndStrategySendLastValue())

		results := helperScraperCustom(
			t,
			time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			15*time.Second,
			9,
			timeSeries.Iterator(),
		)

		require.Equal(t, 9, len(results))
		assert.InDelta(t, 50, results[0].scrapeResult.Value, 0.001)
		assert.InDelta(t, 60, results[1].scrapeResult.Value, 0.001)
		assert.InDelta(t, 70, results[2].scrapeResult.Value, 0.001)
		assert.InDelta(t, 70, results[3].scrapeResult.Value, 0.001)
		assert.InDelta(t, 70, results[4].scrapeResult.Value, 0.001)
		assert.InDelta(t, 70, results[5].scrapeResult.Value, 0.001)
		assert.InDelta(t, 70, results[6].scrapeResult.Value, 0.001)
		assert.InDelta(t, 70, results[7].scrapeResult.Value, 0.001)
		assert.InDelta(t, 70, results[8].scrapeResult.Value, 0.001)
	})

	t.Run("should return a single run of the iterator followed by custom value given the end strategy", func(t *testing.T) {
		dataGenerator, err := discrete.NewLinearSegmentDataGenerator(discrete.LinearSegmentDataGeneratorOptions{
			AmplitudeStart:      50,
			AmplitudeEnd:        70,
			IterationCountLimit: 3,
		})
		require.NoError(t, err)

		timeSeries := discrete.NewMetricTimeSeries(
			"some-series",
			map[string]string{"key": "value"},
			dataGenerator,
			metrics.NewEndStrategySendCustomValue(metrics.ScrapeResult{Value: 123}))

		results := helperScraperCustom(
			t,
			time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			15*time.Second,
			9,
			timeSeries.Iterator(),
		)

		require.Equal(t, 9, len(results))
		assert.InDelta(t, 50, results[0].scrapeResult.Value, 0.001)
		assert.InDelta(t, 60, results[1].scrapeResult.Value, 0.001)
		assert.InDelta(t, 70, results[2].scrapeResult.Value, 0.001)
		assert.InDelta(t, 123, results[3].scrapeResult.Value, 0.001)
		assert.InDelta(t, 123, results[4].scrapeResult.Value, 0.001)
		assert.InDelta(t, 123, results[5].scrapeResult.Value, 0.001)
		assert.InDelta(t, 123, results[6].scrapeResult.Value, 0.001)
		assert.InDelta(t, 123, results[7].scrapeResult.Value, 0.001)
		assert.InDelta(t, 123, results[8].scrapeResult.Value, 0.001)
	})

	t.Run("should return a single run of the iterator followed the removal of the time series given the end strategy", func(t *testing.T) {
		dataGenerator, err := discrete.NewLinearSegmentDataGenerator(discrete.LinearSegmentDataGeneratorOptions{
			AmplitudeStart:      50,
			AmplitudeEnd:        70,
			IterationCountLimit: 3,
		})
		require.NoError(t, err)

		timeSeries := discrete.NewMetricTimeSeries(
			"some-series",
			map[string]string{"key": "value"},
			dataGenerator,
			metrics.NewEndStrategyRemoveTimeSeries())

		results := helperScraperCustom(
			t,
			time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			15*time.Second,
			9,
			timeSeries.Iterator(),
		)

		require.Equal(t, 3, len(results))
		assert.InDelta(t, 50, results[0].scrapeResult.Value, 0.001)
		assert.InDelta(t, 60, results[1].scrapeResult.Value, 0.001)
		assert.InDelta(t, 70, results[2].scrapeResult.Value, 0.001)
	})
}
