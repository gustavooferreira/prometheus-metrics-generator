package metrics

type EndStrategyType string

const (
	EndStrategyTypeLoop             EndStrategyType = "end_strategy_type-loop"
	EndStrategyTypeSendLastValue    EndStrategyType = "end_strategy_type-send_last_value"
	EndStrategyTypeSendCustomValue  EndStrategyType = "end_strategy_type-send_custom_value"
	EndStrategyTypeRemoveTimeSeries EndStrategyType = "end_strategy_type-remove_time_series"
)

// EndStrategy represents the end strategy.
// When the series finishes cycling through the data it reaches the end strategy.
// There are 4 possibilities when it comes to the end strategy.
// - Loop back to the beginning of the data
// - Send the last value sent forever
// - Send a custom value (can be a missing scrape or even zero) forever
// - Remove the time series
// The zero value is not useful. Use one of the helper functions prefixed with NewEndStrategy.
type EndStrategy struct {
	EndStrategyType EndStrategyType

	// CustomValue represents the value to be sent if StrategyType is SendCustomValueEndStrategy
	customValue ScrapeResult
}

func (es EndStrategy) CustomValue() ScrapeResult {
	return es.customValue
}

// NewEndStrategyLoop returns an end strategy that loops back to the beginning of data.
func NewEndStrategyLoop() EndStrategy {
	return EndStrategy{
		EndStrategyType: EndStrategyTypeLoop,
	}
}

// NewEndStrategySendLastValue returns an end strategy that sends the last value forever.
func NewEndStrategySendLastValue() EndStrategy {
	return EndStrategy{
		EndStrategyType: EndStrategyTypeSendLastValue,
	}
}

// NewEndStrategySendCustomValue returns an end strategy that sends a custom value forever.
func NewEndStrategySendCustomValue(customValue ScrapeResult) EndStrategy {
	return EndStrategy{
		EndStrategyType: EndStrategyTypeSendCustomValue,
		customValue:     customValue,
	}
}

// NewEndStrategyRemoveTimeSeries returns an end strategy that removes the time series.
func NewEndStrategyRemoveTimeSeries() EndStrategy {
	return EndStrategy{
		EndStrategyType: EndStrategyTypeRemoveTimeSeries,
	}
}
