package series

import "time"

// Validate exports the private function validate().
func (sc *ScraperConfig) Validate() error {
	return sc.validate()
}

// ApplyFunctionalOptions exports the private function applyFunctionalOptions().
func (sc *ScraperConfig) ApplyFunctionalOptions(opts ...ScraperOption) {
	sc.applyFunctionalOptions(opts...)
}

// GetEndTime returns the unexported endTime field.
func (sc *ScraperConfig) GetEndTime() time.Time {
	return sc.endTime
}

// GetIterationCountLimit returns the unexported iterationCountLimit field.
func (sc *ScraperConfig) GetIterationCountLimit() int {
	return sc.iterationCountLimit
}
