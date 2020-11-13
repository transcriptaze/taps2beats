package taps

import (
	"time"

	"github.com/twystd/ckmeans.1d.dp/go/ckmeans"
)

type Beat struct {
	at time.Duration
}

func taps2beats(taps [][]time.Duration) []Beat {
	beats := []Beat{}

	array := []float64{}
	for _, row := range taps {
		for _, t := range row {
			array = append(array, t.Seconds())
		}
	}

	clusters := ckmeans.CKMeans1dDp(array, nil)

	for _, c := range clusters {
		beats = append(beats, Beat{
			at: time.Duration(c.Center * float64(time.Second)),
		})
	}

	return beats
}
