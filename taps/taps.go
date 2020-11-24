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
	return taps2beats(floats2seconds(taps), seconds(0), seconds(8.5)), nil
}

func taps2beats(taps [][]time.Duration, start, end time.Duration) []Beat {
	array := []float64{}
	for _, row := range taps {
		for _, t := range row {
			array = append(array, t.Seconds())
		}
	}

	clusters := ckmeans.CKMeans1dDp(array, nil)

	return extrapolate(clusters, start, end)
}

func extrapolate(clusters []ckmeans.Cluster, start, end time.Duration) []Beat {
	beats := []Beat{}

	if len(clusters) < 2 {
		// TODO
	} else {

		b := interpolate(clusters)
		index := map[int]ckmeans.Cluster{}
		x := make([]float64, len(clusters))
		t := make([]float64, len(clusters))

		for i, c := range clusters {
			x[i] = float64(b[i])
			t[i] = c.Center
			index[b[i]] = c
		}

		m, c := regression.OLS(x, t)

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
	}

	return beats
}

// TODO assumes clusters are time sorted
func interpolate(clusters []ckmeans.Cluster) []int {
	N := len(clusters)
	beats := make([]int, N)

	for i := range clusters {
		beats[i] = i + 1
	}

	// ... trivial cases
	if N <= 2 {
		return beats
	}

	// ... 3+ intervals

	x0 := clusters[0].Center
	xn := clusters[N-1].Center
	y0 := 1.0

	// ... b0

	yn := float64(N)
	m := (yn - y0) / (xn - x0)
	c := yn - m*xn

	sumsq := 0.0
	for j := 0; j < N; j++ {
		x := clusters[j].Center
		y := m*x + c
		beatf := math.Round(y)
		sumsq += y*y - 2*y*beatf + beatf*beatf
	}

	variance := sumsq / float64(N-1)

	// TODO gradient descent (?)
	for bn := N + 1; variance > 0.001; bn++ {
		yn := float64(bn)
		m := (yn - y0) / (xn - x0)
		c := yn - m*xn

		sumsq := 0.0
		for j := 0; j < N; j++ {
			x := clusters[j].Center
			y := m*x + c
			beatf := math.Round(y)
			sumsq += y*y - 2*y*beatf + beatf*beatf
		}

		v := sumsq / float64(N-1)

		if v < variance {
			for j := 0; j < N; j++ {
				x := clusters[j].Center
				y := m*x + c
				beats[j] = int(math.Round(y))
			}
			variance = v
		}
	}

	return beats
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
