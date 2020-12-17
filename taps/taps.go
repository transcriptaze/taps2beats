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
	BPM    uint          `json:"BPM"`
	Offset time.Duration `json:"offset"`
	Beats  []Beat        `json:"beats"`
}

type Beat struct {
	beat     int             `json:"-"`
	At       time.Duration   `json:"at"`
	Mean     time.Duration   `json:"mean"`
	Variance time.Duration   `json:"variance"`
	Taps     []time.Duration `json:"taps"`
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

	clusters := ckmeans.CKMeans1dDp(data, weights(taps, forgetting))

	beats := make([]Beat, len(clusters))
	for i, cluster := range clusters {
		beats[i] = makeBeat(cluster.Center, cluster)
	}

	BPM, offset := bpm(beats)

	sort.SliceStable(beats, func(i, j int) bool { return beats[i].At < beats[j].At })

	return Beats{
		BPM:    BPM,
		Offset: offset,
		Beats:  beats,
	}
}

func (beats *Beats) Quantize() error {
	switch {
	case beats == nil:
		return nil

	case len(beats.Beats) < 1:
		beats.BPM = 0
		beats.Offset = 0 * time.Millisecond
		return nil

	case len(beats.Beats) < 2:
		beats.BPM = 0
		beats.Offset = beats.Beats[0].At
		return nil

	default:
		m, c, err := fit(beats.Beats)
		if err != nil {
			return err
		}

		quantized := []Beat{}
		for _, b := range beats.Beats {
			quantized = append(quantized, Beat{
				At:       Seconds(float64(b.beat)*m + c),
				Mean:     b.Mean,
				Variance: b.Variance,
				Taps:     b.Taps,
			})
		}

		beats.BPM, beats.Offset = bpm(quantized)
		beats.Beats = quantized

		return nil
	}
}

func (beats *Beats) Interpolate(start, end time.Duration) error {
	switch {
	case beats == nil:
		return nil

	case len(beats.Beats) == 0:
		return fmt.Errorf("Insufficient data")

	case len(beats.Beats) == 1 && beats.BPM == 0:
		return fmt.Errorf("Insufficient data")

	case len(beats.Beats) == 1:
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

		beats.BPM, beats.Offset = bpm(interpolated)
		beats.Beats = interpolated

		return nil

	default:
		m, c, err := fit(beats.Beats)
		if err != nil {
			return err
		}

		index := map[int]Beat{}
		for _, b := range beats.Beats {
			index[b.beat] = b
		}

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

		beats.BPM, beats.Offset = bpm(interpolated)
		beats.Beats = interpolated

		return nil
	}
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

func bpm(beats []Beat) (uint, time.Duration) {
	if len(beats) < 2 {
		return 0, 0
	}

	m, c, err := fit(beats)
	if err != nil {
		return 0, 0
	}

	bpm := uint(math.Round(60.0 / m))

	b0 := int(math.Floor(-c / m))
	t0 := float64(b0)*m + c
	for t0 < 0.0 {
		b0++
		t0 = float64(b0)*m + c
	}

	offset := Seconds(t0)

	return bpm, offset
}

func fit(beats []Beat) (float64, float64, error) {
	if len(beats) < 2 {
		panic("Insufficient data")
	}

	err := reindex(beats)
	if err != nil {
		return 0, 0, err
	}

	x := []float64{}
	t := []float64{}
	for _, b := range beats {
		x = append(x, float64(b.beat))
		t = append(t, b.At.Seconds())
	}

	m, c := regression.OLS(x, t)

	return m, c, nil
}

func reindex(beats []Beat) error {
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
		for i := range beats {
			beats[i].beat = index[i]
		}

		return nil
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
			for i := range beats {
				beats[i].beat = index[i]
			}

			return nil
		}
	}

	return fmt.Errorf("Error mapping taps to beats: %v", beats)
}

func makeBeat(at float64, cluster ckmeans.Cluster) Beat {
	taps := make([]time.Duration, len(cluster.Values))

	for i, v := range cluster.Values {
		taps[i] = Seconds(v)
	}

	//	sort.SliceStable(taps, func(i, j int) bool { return taps[i] < taps[j] })

	return Beat{
		At:       Seconds(at),
		Mean:     Seconds(cluster.Center),
		Variance: Seconds(cluster.Variance),
		Taps:     taps,
	}
}
