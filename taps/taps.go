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

const (
	MaxBPM         int = 200
	MinSubdivision int = 8
)

var Precision = time.Millisecond

func Taps2Beats(taps [][]time.Duration, start, end time.Duration) []Beat {
	array := []float64{}
	for _, row := range taps {
		for _, t := range row {
			array = append(array, t.Round(Precision).Seconds())
		}
	}

	clusters := ckmeans.CKMeans1dDp(array, nil)

	// TODO sort clusters by time (just in case)

	b, err := interpolate(clusters)
	if err != nil {
		beats := []Beat{}
		for _, cluster := range clusters {
			beats = append(beats, makeBeat(cluster.Center, cluster))
		}

		return beats
	}

	m := map[int]ckmeans.Cluster{}

	for i, c := range clusters {
		m[b[i]] = c
	}

	beats := extrapolate(m, start, end)

	// TODO sort beats by time

	return beats
}

func extrapolate(clusters map[int]ckmeans.Cluster, start, end time.Duration) []Beat {
	beats := []Beat{}

	if len(clusters) < 2 {
		for _, cluster := range clusters {
			beats = append(beats, makeBeat(cluster.Center, cluster))
		}

		return beats
	}

	x := []float64{}
	t := []float64{}
	for ix, c := range clusters {
		x = append(x, float64(ix))
		t = append(t, c.Center)
	}

	m, c := regression.OLS(x, t)

	bmin := int(math.Floor((start.Seconds() - c) / m))
	bmax := int(math.Ceil((end.Seconds() - c) / m))

	for bb := bmin; bb <= bmax; bb++ {
		tt := float64(bb)*m + c
		if tt >= start.Seconds() && tt <= end.Seconds() {
			beats = append(beats, makeBeat(tt, clusters[bb]))
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

	dt := Seconds(xn - x0).Minutes()
	bmax := int(math.Ceil(dt * float64(MaxBPM*MinSubdivision/4)))

loop:
	for i := N; i <= bmax; i++ {
		yn := float64(i)
		m := (yn - y0) / (xn - x0)
		c := yn - m*xn

		x := clusters[0].Center
		y := m*x + c
		b0 := math.Round(y)
		beats[0] = int(b0)
		sumsq := y*y - 2*y*b0 + b0*b0

		for j := 1; j < N; j++ {
			x := clusters[j].Center
			y := m*x + c
			bn := math.Round(y)

			beats[j] = int(bn)
			if beats[j] <= beats[j-1] {
				continue loop
			}

			sumsq += y*y - 2*y*bn + bn*bn
		}

		variance := sumsq / float64(N-1)

		if variance < 0.001 {
			return beats, nil
		}
	}

	return nil, fmt.Errorf("Error interpolating beats: %v", beats)
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
	return time.Duration(g * float64(time.Second)).Round(Precision)
}

func makeBeat(at float64, cluster ckmeans.Cluster) Beat {
	taps := make([]time.Duration, len(cluster.Values))

	for i, v := range cluster.Values {
		taps[i] = Seconds(v)
	}

	return Beat{
		At:       Seconds(at).Round(Precision),
		Mean:     Seconds(cluster.Center).Round(Precision),
		Variance: Seconds(cluster.Variance).Round(Precision),
		Taps:     taps,
	}
}
