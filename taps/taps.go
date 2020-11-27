package taps

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/twystd/ckmeans.1d.dp/go/ckmeans"
	"github.com/twystd/taps2beats/regression"
)

type Beats struct {
	BPM    *uint
	Offset *time.Duration
	Beats  []Beat
}

type Beat struct {
	At       time.Duration
	Mean     time.Duration
	Variance time.Duration
	Taps     []time.Duration
}

type T2B struct {
	Precision  time.Duration
	Latency    time.Duration
	Forgetting float64
}

const (
	MaxBPM         int = 200
	MinSubdivision int = 8
)

var Default = T2B{
	Precision:  1 * time.Millisecond,
	Latency:    0 * time.Millisecond,
	Forgetting: 0.0,
}

func (t2b *T2B) Taps2Beats(taps [][]time.Duration, start, end time.Duration) Beats {
	// ... cluster taps into beats
	data := []float64{}
	weights := []float64{}
	w := 1.0
	for _, row := range taps {
		for _, t := range row {
			data = append(data, t.Seconds())
			weights = append(weights, w)
		}

		w = w * (1.0 - t2b.Forgetting)
	}

	clusters := ckmeans.CKMeans1dDp(data, weights)

	sort.SliceStable(clusters, func(i, j int) bool { return clusters[i].Center < clusters[j].Center })

	// ... estimate beats
	var beats = []Beat{}
	var BPM *uint

	if b, err := interpolate(clusters); err != nil {
		for _, cluster := range clusters {
			beats = append(beats, makeBeat(cluster.Center, cluster))
		}
	} else {
		m := map[int]ckmeans.Cluster{}
		for i, c := range clusters {
			m[b[i]] = c
		}

		beats, BPM = extrapolate(m, start, end)
	}

	// ... compensate for latency
	for i, b := range beats {
		beats[i].At = (b.At - t2b.Latency)

		if len(b.Taps) > 0 {
			beats[i].Mean = (b.Mean - t2b.Latency)
		}
	}

	// ... round to precision
	for i, b := range beats {
		beats[i].At = b.At.Round(t2b.Precision)
		beats[i].Mean = b.Mean.Round(t2b.Precision)
		beats[i].Variance = b.Variance.Round(t2b.Precision)

		for j, t := range b.Taps {
			beats[i].Taps[j] = t.Round(t2b.Precision)
		}
	}

	sort.SliceStable(beats, func(i, j int) bool { return beats[i].At < beats[j].At })

	return Beats{
		BPM:   BPM,
		Beats: beats,
	}
}

func extrapolate(clusters map[int]ckmeans.Cluster, start, end time.Duration) ([]Beat, *uint) {
	beats := []Beat{}

	if len(clusters) < 2 {
		for _, cluster := range clusters {
			beats = append(beats, makeBeat(cluster.Center, cluster))
		}

		return beats, nil
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

	bpm := uint(math.Round(60.0 / m))

	return beats, &bpm
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
