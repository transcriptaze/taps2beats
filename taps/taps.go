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
	BPM    uint
	Offset time.Duration
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

func (t2b *T2B) Taps2Beats(taps [][]time.Duration) Beats {
	// ... cluster taps into beats
	data := []float64{}
	for _, row := range taps {
		for _, t := range row {
			data = append(data, t.Seconds())
		}
	}

	weights := make([]float64, len(data))
	switch {
	case t2b.Forgetting == 0.0:
		for i := range weights {
			weights[i] = 1.0
		}

	case t2b.Forgetting > 0.0:
		ix := len(weights) - 1
		w := 1.0
		f := 1.0 - t2b.Forgetting
		for _, row := range taps {
			for range row {
				weights[ix] = w
				ix--
			}

			w = w * f
		}

	case t2b.Forgetting < 0.0:
		ix := 0
		w := 1.0
		f := 1.0 + t2b.Forgetting
		for _, row := range taps {
			for range row {
				weights[ix] = w
				ix++
			}

			w = w * f
		}
	}

	clusters := ckmeans.CKMeans1dDp(data, weights)

	// ... estimate BPM and offset
	beats, BPM, offset := bpm(clusters)

	// ... compensate for latency
	offset -= t2b.Latency
	for i, b := range beats {
		beats[i].At = (b.At - t2b.Latency)
		beats[i].Mean = (b.Mean - t2b.Latency)
		for j, t := range b.Taps {
			beats[i].Taps[j] = t - t2b.Latency
		}
	}

	// ... round to precision

	offset = offset.Round(t2b.Precision)

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
		BPM:    BPM,
		Offset: offset,
		Beats:  beats,
	}
}

func (t2b *T2B) Quantize(beats Beats) (Beats, error) {
	if len(beats.Beats) < 2 {
		quantized := Beats{
			BPM:    0,
			Offset: 0 * time.Millisecond,
			Beats:  make([]Beat, len(beats.Beats)),
		}

		copy(quantized.Beats, beats.Beats)

		if len(quantized.Beats) > 0 {
			quantized.Offset = quantized.Beats[0].At
		}

		return quantized, nil
	}

	index, err := remap(beats)
	if err != nil {
		return beats, err
	}

	x := []float64{}
	t := []float64{}
	for ix, b := range index {
		x = append(x, float64(ix))
		t = append(t, b.At.Seconds())
	}

	m, c := regression.OLS(x, t)

	quantized := []Beat{}
	for ix, b := range index {
		quantized = append(quantized, Beat{
			At:       Seconds(float64(ix)*m + c),
			Mean:     b.Mean,
			Variance: b.Variance,
			Taps:     b.Taps,
		})
	}

	sort.SliceStable(quantized, func(i, j int) bool { return quantized[i].At < quantized[j].At })

	b0 := int(math.Floor(-c / m))
	t0 := float64(b0)*m + c
	for t0 < 0.0 {
		b0++
		t0 = float64(b0)*m + c
	}

	return Beats{
		BPM:    uint(math.Round(60.0 / m)),
		Offset: Seconds(t0),
		Beats:  quantized,
	}, nil
}

func (t2b *T2B) Interpolate(beats Beats, start, end time.Duration) (Beats, error) {
	if len(beats.Beats) == 0 {
		return beats, fmt.Errorf("Insufficient data")
	}

	if len(beats.Beats) == 1 {
		if beats.BPM == 0 {
			return beats, fmt.Errorf("Insufficient data")
		}

		m := 60.0 / float64(beats.BPM)
		c := beats.Beats[0].At.Seconds() - m
		bmin := int(math.Floor((start.Seconds() - c) / m))
		bmax := int(math.Ceil((end.Seconds() - c) / m))

		interpolated := []Beat{}
		for b := bmin; b <= bmax; b++ {
			tt := float64(b)*m + c
			if tt >= start.Seconds() && tt <= end.Seconds() {
				if b == 1 {
					interpolated = append(interpolated, beats.Beats[0])
				} else {
					interpolated = append(interpolated, Beat{At: Seconds(tt)})
				}
			}
		}

		b0 := int(math.Floor(-c / m))
		t0 := float64(b0)*m + c
		for t0 < 0.0 {
			b0++
			t0 = float64(b0)*m + c
		}

		return Beats{
			BPM:    uint(math.Round(60.0 / m)),
			Offset: Seconds(t0),
			Beats:  interpolated,
		}, nil
	}

	index, err := remap(beats)
	if err != nil {
		return beats, err
	}

	x := []float64{}
	t := []float64{}
	for ix, b := range index {
		x = append(x, float64(ix))
		t = append(t, b.At.Seconds())
	}

	m, c := regression.OLS(x, t)
	bmin := int(math.Floor((start.Seconds() - c) / m))
	bmax := int(math.Ceil((end.Seconds() - c) / m))

	interpolated := []Beat{}
	for b := bmin; b <= bmax; b++ {
		tt := float64(b)*m + c
		if tt >= start.Seconds() && tt <= end.Seconds() {
			if beat, ok := index[b]; ok {
				interpolated = append(interpolated, beat)
			} else {
				interpolated = append(interpolated, Beat{At: Seconds(tt)})
			}
		}
	}

	b0 := int(math.Floor(-c / m))
	t0 := float64(b0)*m + c
	for t0 < 0.0 {
		b0++
		t0 = float64(b0)*m + c
	}

	return Beats{
		BPM:    uint(math.Round(60.0 / m)),
		Offset: Seconds(t0),
		Beats:  interpolated,
	}, nil
}

func (t2b *T2B) Shift(beats Beats) Beats {
	shifted := Beats{
		BPM:    beats.BPM,
		Offset: 0,
		Beats:  make([]Beat, len(beats.Beats)),
	}

	if len(beats.Beats) > 0 {
		sort.SliceStable(beats.Beats, func(i, j int) bool { return beats.Beats[i].At < beats.Beats[j].At })

		shift := beats.Beats[0].At
		for i, b := range beats.Beats {
			shifted.Beats[i] = Beat{
				At:       b.At - shift,
				Mean:     b.Mean,
				Variance: b.Variance,
				Taps:     make([]time.Duration, len(b.Taps)),
			}

			if len(b.Taps) > 0 {
				shifted.Beats[i].Mean = b.Mean - shift
				for j, t := range b.Taps {
					shifted.Beats[i].Taps[j] = t - shift
				}
			}
		}
	}

	return shifted
}

func bpm(clusters []ckmeans.Cluster) ([]Beat, uint, time.Duration) {
	beats := []Beat{}

	for _, cluster := range clusters {
		beats = append(beats, makeBeat(cluster.Center, cluster))
	}

	if len(beats) < 2 {
		return beats, 0, 0
	}

	interpolated, err := interpolate(clusters)
	if err != nil {
		return beats, 0, 0
	}

	index := map[int]ckmeans.Cluster{}
	for i, c := range clusters {
		index[interpolated[i]] = c
	}

	x := []float64{}
	t := []float64{}
	for i, c := range index {
		x = append(x, float64(i))
		t = append(t, c.Center)
	}

	m, c := regression.OLS(x, t)
	bpm := uint(math.Round(60.0 / m))

	b0 := int(math.Floor(-c / m))
	t0 := float64(b0)*m + c
	for t0 < 0.0 {
		b0++
		t0 = float64(b0)*m + c
	}

	offset := Seconds(t0)

	return beats, bpm, offset
}

/* NOTE: assumes clusters are time sorted
 */
func interpolate(clusters []ckmeans.Cluster) ([]int, error) {
	sort.SliceStable(clusters, func(i, j int) bool { return clusters[i].Center < clusters[j].Center })

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

func remap(beats Beats) (map[int]Beat, error) {
	sort.SliceStable(beats.Beats, func(i, j int) bool { return beats.Beats[i].At < beats.Beats[j].At })

	at := make([]float64, len(beats.Beats))
	for i, b := range beats.Beats {
		at[i] = b.At.Seconds()
	}

	N := len(at)
	index := make([]int, N)

	for i := range index {
		index[i] = i + 1
	}

	// ... trivial cases
	if N <= 2 {
		m := map[int]Beat{}

		for i, ix := range index {
			m[ix] = beats.Beats[i]
		}

		return m, nil
	}

	// ... 3+ intervals

	x0 := at[0]
	xn := at[N-1]
	y0 := 1.0

	dt := Seconds(xn - x0).Minutes()
	bmax := int(math.Ceil(dt * float64(MaxBPM*MinSubdivision/4)))

loop:
	for i := N; i <= bmax; i++ {
		yn := float64(i)
		m := (yn - y0) / (xn - x0)
		c := yn - m*xn

		x := at[0]
		y := m*x + c
		b0 := math.Round(y)
		index[0] = int(b0)
		sumsq := y*y - 2*y*b0 + b0*b0

		for j := 1; j < N; j++ {
			x := at[j]
			y := m*x + c
			bn := math.Round(y)

			index[j] = int(bn)
			if index[j] <= index[j-1] {
				continue loop
			}

			sumsq += y*y - 2*y*bn + bn*bn
		}

		variance := sumsq / float64(N-1)

		if variance < 0.001 {
			m := map[int]Beat{}

			for i, ix := range index {
				m[ix] = beats.Beats[i]
			}

			return m, nil
		}
	}

	return nil, fmt.Errorf("Error mapping taps to beats: %v", beats)
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
