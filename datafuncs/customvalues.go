package datafuncs

import (
	"gustavooferreira/prometheus-metrics-generator/series"
)

// CustomValuesDataIterator returns a DataIterator containing the array of values passed in.
// Each value is returned in sequence on each scrape.
// It won't take time into consideration, simply return the next value when the scrape is performed.
// This means that if the scrape fails, a given value will be returned on the next scrape.
func CustomValuesDataIterator(values []CustomValue) series.DataIterator {
	// Keeps track of where we are in the array
	var iterationCount int

	return func(scrapeInfo series.ScrapeInfo) series.ScrapeResult {
		if iterationCount >= len(values) {
			return series.ScrapeResult{Exhausted: true}
		}

		iterationValue := values[iterationCount]
		iterationCount++
		return series.ScrapeResult{Value: iterationValue.Value, Missing: iterationValue.Missing}
	}
}

// CustomValue contains the scrape value to be returned by CustomValuesDataIterator.
type CustomValue struct {
	// Value is the value of the sample.
	Value float64

	// Missing indicates whether the scrape failed to retrieve a sample.
	// Used to simulate failed scrapes.
	Missing bool
}
