package promwrite

import "github.com/gustavooferreira/prometheus-metrics-generator/promadapter"

// ConvertToRemoteWriterTimeSeries takes a slice of metric results and creates the corresponding slice of time series
// in the format the PrometheusRemoteWriter expects.
func ConvertToRemoteWriterTimeSeries(metricName string, metricResults []promadapter.MetricResult) []TimeSeries {
	var remoteWriterTimeSeries []TimeSeries

	for _, metricResult := range metricResults {
		labels := make([]Label, 0, len(metricResult.LabelsSet))

		labels = append(labels, Label{
			Name:  "__name__",
			Value: metricName,
		})

		for labelName, labelValue := range metricResult.LabelsSet {
			labels = append(labels, Label{
				Name:  labelName,
				Value: labelValue,
			})
		}

		sampleValue := metricResult.Value

		if metricResult.StaleMarker {
			sampleValue = staleMarker
		}

		remoteWriterSingleTimeSeries := TimeSeries{
			Labels: labels,
			Samples: []Sample{{
				Time:  metricResult.Timestamp,
				Value: sampleValue,
			}},
		}

		remoteWriterTimeSeries = append(remoteWriterTimeSeries, remoteWriterSingleTimeSeries)
	}

	return remoteWriterTimeSeries
}
