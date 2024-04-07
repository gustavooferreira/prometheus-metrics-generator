// Package discrete provides discrete data functions.
// This means the user specifies what the data function returns on every single scrape.
// Be mindful of the caveat that if there is a missing scrape, the value of the sample will be returned on the next
// scrape instead which might lead to distortions.
// Example:
// If the user specifies that it wants a linear segment to go from the value 10 to 30 in 3 iterations, if the second
// scrape fails for whatever reason, then the time series will become: [10, missing, 20, 30], instead of the expected
// [10, 20, 30]. This means it would take a whole scrape interval to reach to the same point.
package discrete
