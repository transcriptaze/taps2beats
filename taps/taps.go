package taps

import (
	"math"
	"time"

	"github.com/twystd/ckmeans.1d.dp/go/ckmeans"
	"github.com/twystd/taps2beats/regression"
)

type Beat struct {
	At       time.Duration
	Mean     time.Duration
	Variance time.Duration
	Taps     []time.Duration
}

func Taps2Beats(taps [][]float64) ([]Beat, error) {
	return taps2beats(floats2seconds(taps), seconds(0), seconds(8.5))
}

func taps2beats(taps [][]time.Duration, start, end time.Duration) ([]Beat, error) {
	//	beats := []Beat{}

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

	return extrapolate(clusters, start, end)
}

func extrapolate(clusters []ckmeans.Cluster, start, end time.Duration) ([]Beat, error) {
	b := interpolate(clusters)
	index := map[int]ckmeans.Cluster{}
	x := make([]float64, len(clusters))
	t := make([]float64, len(clusters))

	for i, c := range clusters {
		x[i] = float64(b[i])
		t[i] = c.Center
		index[b[i]] = c
	}

	m, c, err := regression.Fit(x, t)
	if err != nil {
		return nil, err
	}

	beats := []Beat{}
	bmin := int(math.Floor((start.Seconds() - c) / m))
	bmax := int(math.Ceil((end.Seconds() - c) / m))

	for bb := bmin; bb <= bmax; bb++ {
		tt := float64(bb)*m + c
		if tt >= start.Seconds() && tt <= end.Seconds() {
			cluster := index[bb]
			taps := make([]time.Duration, len(cluster.Values))

			for i, v := range cluster.Values {
				taps[i] = seconds(v)
			}

			beats = append(beats, Beat{
				At:       seconds(tt),
				Mean:     seconds(cluster.Center),
				Variance: seconds(cluster.Variance),
				Taps:     taps,
			})
		}
	}

	return beats, nil
}

func interpolate(clusters []ckmeans.Cluster) []int {
	return []int{1, 2, 3, 4, 5, 6, 7, 8}
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
