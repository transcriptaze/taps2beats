package taps

import (
	"time"

	"github.com/twystd/beats/regression"
	"github.com/twystd/ckmeans.1d.dp/go/ckmeans"
)

type Beat struct {
	at       time.Duration
	mean     time.Duration
	variance time.Duration
	taps     []time.Duration
}

func taps2beats(taps [][]time.Duration, start, end time.Duration) ([]Beat, error) {
	beats := []Beat{}

	array := []float64{}
	for _, row := range taps {
		for _, t := range row {
			array = append(array, t.Seconds())
		}
	}

	clusters := ckmeans.CKMeans1dDp(array, nil)

	b := make([]float64, len(clusters))
	t := make([]float64, len(clusters))
	for i, c := range clusters {
		b[i] = float64(i + 1)
		t[i] = c.Center
	}

	at, err := regression.Trend(b, t, b)
	if err != nil {
		return nil, err
	}

	for i, c := range clusters {
		t := make([]time.Duration, len(c.Values))

		for i, v := range c.Values {
			t[i] = seconds(v)
		}

		beats = append(beats, Beat{
			at:       seconds(at[i]),
			mean:     seconds(c.Center),
			variance: seconds(c.Variance),
			taps:     t,
		})
	}

	return beats, nil
}

func seconds(g float64) time.Duration {
	return time.Duration(g * float64(time.Second))
}
