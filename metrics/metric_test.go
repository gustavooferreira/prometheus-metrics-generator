package metrics_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gustavooferreira/prometheus-metrics-generator/discrete"
	"github.com/gustavooferreira/prometheus-metrics-generator/metrics"
)

func TestMetric(t *testing.T) {
	t.Run("should return valid metric descriptor", func(t *testing.T) {
		metric := metrics.NewMetric(
			"some-metric",
			"some-help-description",
			metrics.MetricTypeGauge,
			[]string{"label1", "label2"},
		)

		desc := metric.Desc()

		assert.Equal(t, "some-metric", desc.FQName)
		assert.Equal(t, "some-help-description", desc.Help)
		assert.Equal(t, metrics.MetricTypeGauge, desc.MetricType)
		assert.Equal(t, []string{"label1", "label2"}, desc.LabelsNames)
	})

	t.Run("should fail to attach a time series which includes an unexpected label name", func(t *testing.T) {
		dataGenerator := discrete.NewVoidDataGenerator(5)

		timeSeries := discrete.NewMetricTimeSeries(
			map[string]string{"label_extra": "value"},
			dataGenerator,
			metrics.NewEndStrategyRemoveTimeSeries(),
		)

		metric := metrics.NewMetric(
			"some-metric",
			"some-help-description",
			metrics.MetricTypeCounter,
			[]string{"label1", "label2"},
		)

		err := metric.Attach(timeSeries)
		require.Error(t, err)
		assert.Equal(t, "label mismatch: unexpected label in time series", err.Error())

	})

	t.Run("should fail to attach a time series which does not include an expected label name", func(t *testing.T) {
		dataGenerator := discrete.NewVoidDataGenerator(5)

		timeSeries := discrete.NewMetricTimeSeries(
			map[string]string{"label1": "value1"},
			dataGenerator,
			metrics.NewEndStrategyRemoveTimeSeries(),
		)

		metric := metrics.NewMetric(
			"some-metric",
			"some-help-description",
			metrics.MetricTypeCounter,
			[]string{"label1", "label2"},
		)

		err := metric.Attach(timeSeries)
		require.Error(t, err)
		assert.Equal(t, "label mismatch: missing expected label in time series", err.Error())
	})

	t.Run("should never return any results given that the samples are market as being missed", func(t *testing.T) {
		scraper, err := metrics.NewScraper(
			metrics.ScraperConfig{
				StartTime:      time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
				ScrapeInterval: 15 * time.Second,
			},
			metrics.WithScraperIterationCountLimit(5), // It's good practice to set an upper bound in tests
		)
		require.NoError(t, err)

		dataGenerator1 := discrete.NewCustomValuesDataGenerator([]discrete.CustomValue{{Value: 10}})

		timeSeries1 := discrete.NewMetricTimeSeries(
			map[string]string{"label1": "valueA"},
			dataGenerator1,
			metrics.NewEndStrategyRemoveTimeSeries(),
		)

		dataGenerator2 := discrete.NewCustomValuesDataGenerator([]discrete.CustomValue{{Value: 1}, {Value: 2}})

		timeSeries2 := discrete.NewMetricTimeSeries(
			map[string]string{"label1": "valueB"},
			dataGenerator2,
			metrics.NewEndStrategyRemoveTimeSeries(),
		)

		metric := metrics.NewMetric(
			"some-metric",
			"some-help-description",
			metrics.MetricTypeCounter,
			[]string{"label1"},
		)

		err = metric.Attach(timeSeries1)
		require.NoError(t, err)
		err = metric.Attach(timeSeries2)
		require.NoError(t, err)

		metric.Prepare()

		// Evaluate
		var results []resultContainer
		for iter := scraper.Iterator(); iter.HasNext(); {
			scrapeInfo := iter.Next()

			metricResults := metric.Evaluate(scrapeInfo)
			results = append(results, resultContainer{
				scrapeInfo:    scrapeInfo,
				metricResults: metricResults,
			})
		}

		for _, result := range results {
			if len(result.metricResults) == 0 {
				t.Logf("[%3d] Timestamp: %s - No time series\n",
					result.scrapeInfo.IterationCount,
					result.scrapeInfo.IterationTime,
				)
			}

			for _, metricResult := range result.metricResults {
				t.Logf("[%3d] Timestamp: %s - Value: %6.2f\n",
					result.scrapeInfo.IterationCount,
					result.scrapeInfo.IterationTime,
					metricResult.Value,
				)
			}
		}

		// assert results
		require.Equal(t, 5, len(results))
		assert.InDelta(t, 10.0, results[0].metricResults[0].Value, 0.001)
		assert.InDelta(t, 1.0, results[0].metricResults[1].Value, 0.001)
		assert.InDelta(t, 2.0, results[1].metricResults[0].Value, 0.001)
		assert.Equal(t, 0, len(results[2].metricResults))
		assert.Equal(t, 0, len(results[3].metricResults))
		assert.Equal(t, 0, len(results[4].metricResults))
	})
}

type resultContainer struct {
	scrapeInfo    metrics.ScrapeInfo
	metricResults []metrics.MetricResult
}
