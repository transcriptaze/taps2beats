package taps

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/twystd/taps2beats/ckmeans"
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

const (
	MaxBPM         int = 200
	MinSubdivision int = 8
)

func Taps2Beats(taps [][]time.Duration, forgetting float64) Beats {
	data := []float64{}
	for _, row := range taps {
		for _, t := range row {
			data = append(data, t.Seconds())
		}
	}

	w := weights(taps, forgetting)
	clusters := ckmeans.CKMeans1dDp(data, w)
	beats, BPM, offset := bpm(clusters)

	sort.SliceStable(beats, func(i, j int) bool { return beats[i].At < beats[j].At })

	return Beats{
		BPM:    BPM,
		Offset: offset,
		Beats:  beats,
	}
}

func (beats *Beats) Quantize() error {
	if beats != nil {
		if len(beats.Beats) < 1 {
			beats.BPM = 0
			beats.Offset = 0 * time.Millisecond

			return nil
		}

		if len(beats.Beats) < 2 {
			beats.BPM = 0
			beats.Offset = beats.Beats[0].At

			return nil
		}

		index, err := reindex(beats.Beats)
		if err != nil {
			return err
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

		beats.BPM = uint(math.Round(60.0 / m))
		beats.Offset = Seconds(t0)
		beats.Beats = quantized
	}

	return nil
}

func (beats *Beats) Interpolate(start, end time.Duration) error {
	if beats != nil {
		if len(beats.Beats) == 0 {
			return fmt.Errorf("Insufficient data")
		}

		if len(beats.Beats) == 1 && beats.BPM == 0 {
			return fmt.Errorf("Insufficient data")
		}

		// TODO simplify
		if len(beats.Beats) == 1 {
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

			beats.BPM = uint(math.Round(60.0 / m))
			beats.Offset = Seconds(t0)
			beats.Beats = interpolated

			return nil
		}

		index, err := reindex(beats.Beats)
		if err != nil {
			return err
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

		beats.BPM = uint(math.Round(60.0 / m))
		beats.Offset = Seconds(t0)
		beats.Beats = interpolated

		return nil
	}

	return nil
}

func (beats *Beats) Round(precision time.Duration) {
	if beats != nil {
		beats.Offset = beats.Offset.Round(precision)

		for i, b := range beats.Beats {
			beats.Beats[i].At = b.At.Round(precision)
			beats.Beats[i].Mean = b.Mean.Round(precision)
			beats.Beats[i].Variance = b.Variance.Round(precision)
			for j, tap := range b.Taps {
				beats.Beats[i].Taps[j] = tap.Round(precision)
			}
		}
	}
}

func (beats *Beats) Sub(dt time.Duration) {
	if beats != nil {
		beats.Offset -= dt
		for i, b := range beats.Beats {
			beats.Beats[i].At = b.At - dt

			if len(b.Taps) > 0 {
				beats.Beats[i].Mean = b.Mean - dt
				for j, t := range b.Taps {
					beats.Beats[i].Taps[j] = t - dt
				}
			}
		}
	}
}

func weights(taps [][]time.Duration, forgetting float64) []float64 {
	N := 0
	for _, row := range taps {
		N += len(row)
	}

	array := make([]float64, N)
	switch {
	case forgetting == 0.0:
		for i := range array {
			array[i] = 1.0
		}

	case forgetting > 0.0:
		ix := len(array) - 1
		w := 1.0
		f := 1.0 - forgetting
		for _, row := range taps {
			for range row {
				array[ix] = w
				ix--
			}

			w = w * f
		}

	case forgetting < 0.0:
		ix := 0
		w := 1.0
		f := 1.0 + forgetting
		for _, row := range taps {
			for range row {
				array[ix] = w
				ix++
			}

			w = w * f
		}
	}

	return array
}

func bpm(clusters []ckmeans.Cluster) ([]Beat, uint, time.Duration) {
	sort.SliceStable(clusters, func(i, j int) bool { return clusters[i].Center < clusters[j].Center })

	beats := make([]Beat, len(clusters))
	for i, cluster := range clusters {
		beats[i] = makeBeat(cluster.Center, cluster)
	}

	if len(beats) < 2 {
		return beats, 0, 0
	}

	index, err := reindex(beats)
	if err != nil {
		return beats, 0, 0
	}

	x := []float64{}
	t := []float64{}
	for i, c := range index {
		x = append(x, float64(i))
		t = append(t, c.At.Seconds())
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

func reindex(beats []Beat) (map[int]Beat, error) {
	sort.SliceStable(beats, func(i, j int) bool { return beats[i].At < beats[j].At })

	at := make([]float64, len(beats))
	for i, b := range beats {
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
			m[ix] = beats[i]
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
				m[ix] = beats[i]
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
