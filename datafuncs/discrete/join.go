package discrete

import (
	"github.com/gustavooferreira/prometheus-metrics-generator/series"
)

// JoinDataIterator joins all series.DataIterators one after the next.
func JoinDataIterator(dataIterators ...series.DataIterator) series.DataIterator {
	return func(scrapeInfo series.ScrapeInfo) series.ScrapeResult {
		for _, dataIterator := range dataIterators {
			// if the returned sample is not exhausted, return that one
			scrapeResult := dataIterator(scrapeInfo)
			if scrapeResult.Exhausted {
				continue
			}

			return scrapeResult
		}

		return series.ScrapeResult{Exhausted: true}
	}
}
