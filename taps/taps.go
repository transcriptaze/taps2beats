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

type T2B struct {
	Precision time.Duration
	Latency   time.Duration
}

const (
	MaxBPM         int = 200
	MinSubdivision int = 8
)

var Default = T2B{
	Precision: time.Millisecond,
	Latency:   0,
}

func (t2b *T2B) Taps2Beats(taps [][]time.Duration, start, end time.Duration) []Beat {
	// ... cluster taps into beats
	array := []float64{}
	for _, row := range taps {
		for _, t := range row {
			array = append(array, t.Round(t2b.Precision).Seconds())
		}
	}

	clusters := ckmeans.CKMeans1dDp(array, nil)

	// TODO sort clusters by time (just in case)

	// ... fill in gaps
	b, err := t2b.interpolate(clusters)
	if err != nil {
		beats := []Beat{}
		for _, cluster := range clusters {
			beats = append(beats, t2b.makeBeat(cluster.Center, cluster))
		}

		return beats
	}

	m := map[int]ckmeans.Cluster{}

	for i, c := range clusters {
		m[b[i]] = c
	}

	beats := t2b.extrapolate(m, start, end)

	// TODO sort beats by time

	// ... compensate for latency
	for i, b := range beats {
		beats[i].At = (b.At - t2b.Latency).Round(t2b.Precision)

		if len(b.Taps) > 0 {
			beats[i].Mean = (b.Mean - t2b.Latency).Round(t2b.Precision)
		}
	}

	return beats
}

func (t2b *T2B) extrapolate(clusters map[int]ckmeans.Cluster, start, end time.Duration) []Beat {
	beats := []Beat{}

	if len(clusters) < 2 {
		for _, cluster := range clusters {
			beats = append(beats, t2b.makeBeat(cluster.Center, cluster))
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
			beats = append(beats, t2b.makeBeat(tt, clusters[bb]))
		}
	}

	return beats
}

/* NOTE: assumes clusters are time sorted
 */
func (t2b *T2B) interpolate(clusters []ckmeans.Cluster) ([]int, error) {
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

	dt := t2b.Seconds(xn - x0).Minutes()
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

func (t2b *T2B) Floats2Seconds(floats [][]float64) [][]time.Duration {
	l := [][]time.Duration{}

	for _, f := range floats {
		p := []time.Duration{}
		for _, g := range f {
			p = append(p, t2b.Seconds(g))
		}
		l = append(l, p)
	}

	return l
}

func (t2b *T2B) Seconds(g float64) time.Duration {
	return time.Duration(g * float64(time.Second)).Round(t2b.Precision)
}

func (t2b *T2B) makeBeat(at float64, cluster ckmeans.Cluster) Beat {
	taps := make([]time.Duration, len(cluster.Values))

	for i, v := range cluster.Values {
		taps[i] = t2b.Seconds(v)
	}

	return Beat{
		At:       t2b.Seconds(at),
		Mean:     t2b.Seconds(cluster.Center),
		Variance: t2b.Seconds(cluster.Variance),
		Taps:     taps,
	}
}
