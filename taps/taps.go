package taps

import (
	"math"
	"time"

	"github.com/twystd/ckmeans.1d.dp/go/ckmeans"
	"github.com/twystd/taps2beats/regression"
)

type Beat struct {
	at       time.Duration
	mean     time.Duration
	variance time.Duration
	taps     []time.Duration
}

func Taps2Beats(taps []float64) ([]Beat, error) {
	return taps2beats(floats2seconds([][]float64{taps}), seconds(0), seconds(8.5))
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

	extrapolate(clusters, start, end)

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

func extrapolate(clusters []ckmeans.Cluster, start, end time.Duration) ([]Beat, error) {
	b := make([]float64, len(clusters))
	t := make([]float64, len(clusters))

	for i, c := range clusters {
		b[i] = float64(i + 1)
		t[i] = c.Center
	}

	m, c, err := regression.Fit(b, t)
	if err != nil {
		return nil, err
	}

	bmin := int(math.Floor((start.Seconds() - c) / m))
	bmax := int(math.Ceil((end.Seconds() - c) / m))

	beats := []Beat{}
	for bb := bmin; bb <= bmax; bb++ {
		tt := float64(bb)*m + c
		if tt >= start.Seconds() && tt <= end.Seconds() {
			beats = append(beats, Beat{at: seconds(tt)})
		}
	}

	return beats, nil
}

func floats2seconds(floats [][]float64) [][]time.Duration {
	l := [][]time.Duration{}

	for _, f := range floats {
		p := []time.Duration{}
		for _, g := range f {
			p = append(p, seconds(g))
		}
		l = append(l, p)
	}

	return l
}

func seconds(g float64) time.Duration {
	return time.Duration(g * float64(time.Second))
}
