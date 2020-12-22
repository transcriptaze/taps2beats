package taps2beats

import (
	"time"
)

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

func Seconds(g float64) time.Duration {
	return time.Duration(g * float64(time.Second))
}
