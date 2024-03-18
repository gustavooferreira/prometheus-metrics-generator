package metrics

import "time"

// Validate exports the private method validate().
func (sc *ScraperConfig) Validate() error {
	return sc.validate()
}

// ApplyFunctionalOptions exports the private method applyFunctionalOptions().
func (sc *ScraperConfig) ApplyFunctionalOptions(opts ...ScraperOption) {
	sc.applyFunctionalOptions(opts...)
}

// GetEndTime returns the unexported field 'endTime'.
func (sc *ScraperConfig) GetEndTime() time.Time {
	return sc.endTime
}

// GetIterationCountLimit returns the unexported field 'iterationCountLimit'.
func (sc *ScraperConfig) GetIterationCountLimit() int {
	return sc.iterationCountLimit
}
