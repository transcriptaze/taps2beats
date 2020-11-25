package taps

import (
	"fmt"
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

func Taps2Beats(taps [][]time.Duration, start, end time.Duration) []Beat {
	array := []float64{}
	for _, row := range taps {
		for _, t := range row {
			array = append(array, t.Seconds())
		}
	}

	clusters := ckmeans.CKMeans1dDp(array, nil)

	// TODO sort clusters by time (just in case)

	b, err := interpolate(clusters)
	if err != nil {
		beats := make([]Beat, len(clusters))
		for i, cluster := range clusters {
			taps := make([]time.Duration, len(cluster.Values))

			for i, v := range cluster.Values {
				taps[i] = Seconds(v)
			}

			beats[i] = Beat{
				At:       Seconds(cluster.Center),
				Mean:     Seconds(cluster.Center),
				Variance: Seconds(cluster.Variance),
				Taps:     taps,
			}
		}

		return beats
	}

	beats := extrapolate(clusters, b, start, end)

	// TODO sort beats by time

	return beats
}

func extrapolate(clusters []ckmeans.Cluster, b []int, start, end time.Duration) []Beat {
	if len(clusters) != len(b) {
		panic("Expected 'b' index array to match list of clusters")
	}

	beats := []Beat{}

	if len(clusters) < 2 {
		for _, cluster := range clusters {
			beats = append(beats, makeBeat(cluster.Center, cluster))
		}
	} else {
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
					taps[i] = Seconds(v)
				}

				beats = append(beats, Beat{
					At:       Seconds(tt),
					Mean:     Seconds(cluster.Center),
					Variance: Seconds(cluster.Variance),
					Taps:     taps,
				})
			}
		}
	}

	return beats
}

/* NOTE: assumes clusters are time sorted
 */
func interpolate(clusters []ckmeans.Cluster) ([]int, error) {
	N := len(clusters)
	beats := make([]int, N)

	for i := range clusters {
		beats[i] = i + 1
	}

	// ... trivial cases
	if N <= 2 {
		return beats, nil
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
		beats[j] = int(math.Round(y))
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

	// check beats are unique and monotonically increasinge
	for i := 1; i < len(beats); i++ {
		if beats[i] <= beats[i-1] {
			return nil, fmt.Errorf("Error interpolating beats: %v", beats)
		}
	}

	return beats, nil
}

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

func makeBeat(at float64, cluster ckmeans.Cluster) Beat {
	taps := make([]time.Duration, len(cluster.Values))

	for i, v := range cluster.Values {
		taps[i] = Seconds(v)
	}

	return Beat{
		At:       Seconds(at),
		Mean:     Seconds(cluster.Center),
		Variance: Seconds(cluster.Variance),
		Taps:     taps,
	}
}
