package taps2beats

import (
	"time"
)

// Utility function to convert an array of 'taps' in seconds to
// the array of 'taps' as Durations used by Taps2Beats.
func Floats2Seconds(floats [][]float64) [][]time.Duration {
	l := [][]time.Duration{}

	for _, f := range floats {
		p := []time.Duration{}
		for _, g := range f {
			p = append(p, Seconds(g))
		}
		l = append(l, p)
	}

	return l
}

// Utility function to convert a float 'seconds' value to the
// equivalent Duration.
func Seconds(g float64) time.Duration {
	return time.Duration(g * float64(time.Second))
}
