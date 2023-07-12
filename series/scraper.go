package series

import (
	"fmt"
	"time"
)

// Scraper generates scrape that are then passed into a TimeSeries Iterator to generate the time series values.
// This can be used to generate metrics that can then be pushed to prometheus using the remote writer.
// It can also be useful for testing purposes.
// The generated scrapes are precise and do not include any jitter.
// The zero value of Scraper is not useful. Use NewScraper to create a new instance.
type Scraper struct {
	cfg ScraperConfig

	// infiniteGenerator indicates whether the scraper runs forever or not.
	infiniteGenerator bool
}

// NewScraper returns a new instance of Scraper.
func NewScraper(cfg ScraperConfig, opts ...ScraperOption) (Scraper, error) {
	cfg.applyFunctionalOptions(opts...)

	if err := cfg.validate(); err != nil {
		return Scraper{}, fmt.Errorf("error validating scraper configuration: %w", err)
	}

	scraper := Scraper{
		cfg: cfg,
	}

	if cfg.endTime.IsZero() && cfg.iterationCountLimit == 0 {
		scraper.infiniteGenerator = true
	}

	return scraper, nil
}

// IsInfinite reports whether the scraper will run forever or not.
func (s *Scraper) IsInfinite() bool {
	return s.infiniteGenerator
}

// Iterator returns an iterator that can be iterated over to exhaustion.
// Any number of iterators can be retrieved from a single scraper.
func (s *Scraper) Iterator() ScraperIterator {
	return ScraperIterator{
		scraper: *s,
	}
}

// Scrape scrapes the DataIterator according to the settings of the Scraper.
// This function can be used as an alternative to creating an iterator and manually iterate over the scrapes.
// For each generated scrape, this function will call the ScrapeHandler provided.
// This function terminates when the DataIterator has no more data left, or there are no more scrapes to be generated,
// or the ScrapeHandler returns an error.
func (s *Scraper) Scrape(dataIterator DataIterator, scrapeHandler ScrapeHandler) error {
	for iter := s.Iterator(); iter.HasNext(); {
		scrapeInfo := iter.Next()

		scrapeResult := dataIterator(scrapeInfo)
		if scrapeResult.Exhausted {
			// exhausted time series
			return nil
		}

		err := scrapeHandler(scrapeInfo, scrapeResult)
		if err != nil {
			return fmt.Errorf("failed while calling scrape handler: %s", err)
		}
	}

	// exhausted scraper
	return nil
}

// ScrapeHandler defines the function type to be used when calling the Scrape method of the Scraper.
// Return an error to stop the scrapping from proceeding any further. The Scrape method will wrap the returned error and
// return it.
// The field 'Exhausted' in the struct ScrapeResult will never be set, since if the scraper has exhausted the
// DataIterator it will automatically stop.
type ScrapeHandler func(scrapeInfo ScrapeInfo, scrapeResult ScrapeResult) error

// ScraperIterator iterates over the scraper.
type ScraperIterator struct {
	// scraper is a copy of the scraper this iterator was created from.
	scraper Scraper
	// currentIteration keeps track of the current iteration.
	currentIteration int
	// lastTimeStamp contains the timestamp of the last scrape. Used for calculations purposes.
	lastTimeStamp time.Time
}

// HasNext reports whether there are more scrapes to be generated or whether the iterator has been exhausted.
func (si *ScraperIterator) HasNext() bool {
	// Check if we are past the iteration count limit
	if si.scraper.cfg.iterationCountLimit > 0 && si.currentIteration >= si.scraper.cfg.iterationCountLimit {
		return false
	}

	var nextScrapeTime time.Time

	// compute next scrape time
	if si.currentIteration == 0 {
		nextScrapeTime = si.scraper.cfg.StartTime
	} else {
		nextScrapeTime = si.lastTimeStamp.Add(si.scraper.cfg.ScrapeInterval)
	}

	// Check if we are past the endTime
	if !si.scraper.cfg.endTime.IsZero() && nextScrapeTime.After(si.scraper.cfg.endTime) {
		return false
	}

	return true
}

// Next returns the next generated scrape.
func (si *ScraperIterator) Next() ScrapeInfo {
	var nextScrapeTime time.Time

	// compute next scrape time
	if si.currentIteration == 0 {
		nextScrapeTime = si.scraper.cfg.StartTime
	} else {
		nextScrapeTime = si.lastTimeStamp.Add(si.scraper.cfg.ScrapeInterval)
	}

	scrapeInfo := ScrapeInfo{
		FirstIterationTime: si.scraper.cfg.StartTime,
		IterationCount:     si.currentIteration,
		IterationTime:      nextScrapeTime,
	}

	si.currentIteration++
	si.lastTimeStamp = nextScrapeTime

	return scrapeInfo
}

// Reset resets the iterator, meaning the iterator will start from the beginning.
// This function allows for the possibility of reusing an iterator after it's been used.
func (si *ScraperIterator) Reset() {
	si.currentIteration = 0
}
