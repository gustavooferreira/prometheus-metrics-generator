package datafuncs

import (
	"fmt"
	"gustavooferreira/prometheus-metrics-generator/series"
	"math/rand"
	"time"
)

// Generate random noise
// Provide amplitude range
// Generate int range but then return it as a float

func RandomDataIterator(options RandomDataIteratorOptions) (series.DataIterator, error) {
	// validation
	if options.LengthDuration == 0 && options.LengthIterationCount == 0 {
		return nil, fmt.Errorf("stop condition needs to be provided, either set the length duration or " +
			"length iteration count")
	}

	if options.LengthDuration != 0 && options.LengthIterationCount != 0 {
		return nil, fmt.Errorf("only one stop condition should be provided")
	}

	if options.LengthDuration < 0 {
		return nil, fmt.Errorf("length duration cannot be negative")
	}

	if options.LengthIterationCount < 0 {
		return nil, fmt.Errorf("length iteration count cannot be negative")
	}

	if options.LengthDurationExclusive && options.LengthDuration == 0 {
		return nil, fmt.Errorf("length duration exclusive option applies to length duration option only, but " +
			"length duration is not set")
	}

	// These 3 closure variables keep track of the first scrape when running the iterator.
	// This allows us to keep track of how many iterations or how long (timewise) we've been running for.
	// All calculations are performed relative to the first detected scrape.
	var firstScrapeHappened bool
	var firstIterationCount int
	var firstScrapeTime time.Time

	// We might never return a single sample if LengthDuration is less than the time it took to scrape for the first
	// time.
	return func(scrapeInfo series.ScrapeInfo) series.ScrapeResult {
		// Is this the first scrape?
		if !firstScrapeHappened {
			firstScrapeHappened = true
			firstIterationCount = scrapeInfo.IterationCount
			firstScrapeTime = scrapeInfo.IterationTime
		}

		// Normalize
		currentIterationCount := scrapeInfo.IterationCount - firstIterationCount
		currentElapsedTime := scrapeInfo.IterationTime.Sub(firstScrapeTime)

		// Have we reached the end?
		if options.LengthIterationCount != 0 && currentIterationCount >= options.LengthIterationCount {
			return series.ScrapeResult{Exhausted: true}
		} else if options.LengthDuration != 0 {
			if options.LengthDurationExclusive {
				if currentElapsedTime >= options.LengthDuration {
					return series.ScrapeResult{Exhausted: true}
				}
			} else {
				if currentElapsedTime > options.LengthDuration {
					return series.ScrapeResult{Exhausted: true}
				}
			}
		}

		randomValue := rand.Float64()*(options.AmplitudeMax-options.AmplitudeMin) + options.AmplitudeMin
		return series.ScrapeResult{Value: randomValue}
	}, nil
}

// RandomDataIteratorOptions contains the options for the RandomDataIterator.
// Either the LengthIterationCount or the LengthDuration fields must be set.
// It's an error to not set one and only one of those fields.
type RandomDataIteratorOptions struct {
	// AmplitudeMax represents the maximum value the data iterator will return.
	AmplitudeMax float64
	// AmplitudeMin represents the minimum value the data iterator will return.
	AmplitudeMin float64

	// LengthIterationCount sets the number of iterations to be used by the segment.
	LengthIterationCount int

	// LengthDuration sets the max duration between the first and the last value.
	LengthDuration time.Duration

	// LengthDurationExclusive sets the LengthDuration to be exclusive.
	// This means that the range is inclusive on the end time.
	// It's an error to set this field without using LengthDuration.
	LengthDurationExclusive bool
}
