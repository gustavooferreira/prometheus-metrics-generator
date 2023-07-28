// Package continuous provides continuous data functions.
// It simulates the passing of time and computes what the value of the sample should be when a scrape happens.
// All functions provide data for a given time window. The user can specify whether the interval is open, half-open or
// closed.
//
// Intervals:
// (a,b) == ]a,b[ == { a < x < b } --> Open Interval
// [a,b) == [a,b[ == { a <= x < b } --> Half-open Interval
// (a,b] == ]a,b] == { a < x <= b } --> Half-open Interval
// [a,b] == ]a,b] == { a <= x <= b } --> Closed Interval
//
// For more information on how intervals work, please read: https://en.wikipedia.org/wiki/Interval_(mathematics)
package continuous
