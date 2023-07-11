package series

import (
	"fmt"
	"time"
)

// ScraperConfig represents the Scraper config.
// The end time and iteration count limit represent stop conditions for the scraper.
// If both contain the zero value, then the scraper will generate scrapes forever.
// If both are set to a given value, then the scraper will stop when it hits the first stop condition.
type ScraperConfig struct {
	// StartTime defines the initial timestamp the scraper will use.
	StartTime time.Time
	// ScrapeInterval represents the scrape interval.
	ScrapeInterval time.Duration

	// -------------------------------------------------
	// Unexported fields are set via a functional option
	// -------------------------------------------------

	// endTime represents the time at which the scraper should stop generating timestamps.
	// endTime is inclusive, meaning, the scraper will also generate a scrape for the endTime.
	endTime time.Time

	// iterationCountLimit specifies how many iteration cycles the scraper should go through before stopping.
	iterationCountLimit int
}

// validate validates the configuration.
func (sc *ScraperConfig) validate() error {
	if sc.ScrapeInterval <= time.Duration(0) {
		return fmt.Errorf("scrape interval cannot be less or equal to zero")
	}

	// if end time is set, make sure it comes after start time.
	if !sc.endTime.IsZero() {
		if sc.endTime.Before(sc.StartTime) {
			return fmt.Errorf("end time cannot be before start time")
		}
	}

	if sc.iterationCountLimit < 0 {
		return fmt.Errorf("iteration count limit cannot be less than zero")
	}

	return nil
}

// applyFunctionalOptions applies the set of ScraperOption onto the ScraperConfig.
func (sc *ScraperConfig) applyFunctionalOptions(opts ...ScraperOption) {
	for _, opt := range opts {
		opt(sc)
	}
}

// Functional Options -----------------

type ScraperOption func(sc *ScraperConfig)

// WithScraperEndTime defines an end time for the Scraper.
func WithScraperEndTime(endTime time.Time) ScraperOption {
	return func(sc *ScraperConfig) {
		sc.endTime = endTime
	}
}

// WithScraperIterationCountLimit defines a max number of iterations for the scraper.
func WithScraperIterationCountLimit(n int) ScraperOption {
	return func(sc *ScraperConfig) {
		sc.iterationCountLimit = n
	}
}
