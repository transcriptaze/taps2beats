// Provides a set of functions for estimating the beats of a piece of music from a set of 'taps'.
//
// The 'taps' would typically (but not exclusively) be generated by a person tapping on a keyboard
// in time to the music, but the functions can be used to estimate and linearize the beats from
// other beat detection algorithms.
//
// The beats are estimated by clustering the 'taps' into a set of optimal beats using the ckmeans.1d.dp
// algorithm, followed by optional quantization to adjust the estimated beats to a match a fixed BPM and
// (optionally) extrapolation to fill in missing beats.
package taps2beats

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/transcriptaze/taps2beats/taps2beats/ckmeans"
	"github.com/transcriptaze/taps2beats/taps2beats/regression"
)

// Contains the estimated BPM, offset of the first beat and the beats estimated from a set of 'taps'.
type Beats struct {
	BPM      uint          `json:"BPM"`
	Offset   time.Duration `json:"offset"`
	Beats    []Beat        `json:"beats"`
	Variance *float64      `json:"-"`
}

// Contains the estimated time of a single beat, the mean and variance of the 'taps' that were
// used to estimate the beat and a list of the 'taps' that were assigned to this beat.
type Beat struct {
	beat     int             `json:"-"`
	At       time.Duration   `json:"at"`
	Mean     time.Duration   `json:"mean"`
	Variance time.Duration   `json:"variance"`
	Taps     []time.Duration `json:"taps"`
}

// Used for marshaling and unmarshaling time as untyped seconds when marshaling and unmarshaling
// to/from JSON.
type instant time.Duration

// Marshals an instant as the equivalent seconds value, rounded to the nearest millisecond.
func (t instant) MarshalText() ([]byte, error) {
	s := fmt.Sprintf("%.3f", time.Duration(t).Seconds())

	return []byte(s), nil
}

// Marshals an instant as the equivalent seconds value, rounded to the nearest millisecond.
func (t instant) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("%.3f", time.Duration(t).Seconds())

	return []byte(s), nil
}

// Unmarshals an instant as seconds from a float64 value.
func (t *instant) UnmarshalJSON(s []byte) error {
	if t != nil {
		var v float64

		_, err := fmt.Sscanf(string(s), "%f", &v)
		if err != nil {
			return err
		}

		*t = instant(v * float64(time.Second))
	}

	return nil
}

const (
	MaxBPM         int = 200 // Maximum BPM that will be used when quantizing and interpolating beats
	MinSubdivision int = 8   // Minimum subdivision (eighths) that will be used when quantizing and interpolating beats
)

// Clusters the provided 'taps' into an optimal set of beats and estimates the average BPM and the offset of the
// first beats (on the assumption that the BPM is fixed).
//
// The supplied list of 'taps' may contain multiple loops, with each loop being a separate row in the list.
//
// The forgetting factor is used to discount earlier loops in favour of later loops, on the grounds that later 'taps'
// will probably be more accurate. A forgetting factor of 0.0 assumes all taps are equally accurate, while a value of
// 0.1 discounts each loop by 10% over the subsequent loop. A forgetting factor of -0.1 discounts each subsequent loop
// by 10% over the preceding loop.
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

	result := Beats{
		BPM:    BPM,
		Offset: offset,
		Beats:  beats,
	}

	if len(beats) > 0 {
		variance := 0.0
		for _, b := range beats {
			variance += b.Variance.Seconds()
		}

		variance = variance / float64(len(beats))

		result.Variance = &variance
	}

	return result
}

// Adjusts the times of the beats by performing a least squares reqression to fit the estimated beats to
// a straight line (on the assumption that the BPM is reasonably constant).
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

// Estimates beats that are not in the provided list by using least squares regression to fit the beats
// to a straight line (assumes the BPM is reasonably constant).
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

// Discards 'outlier' beats in a desperate attempt to obtain a better estimate of the beats and BPM.
// Outlier beats are identified as those beats which have 'fewer than expected' taps, where 'fewer
// than expected' is defined as less than a third of the median.
//
// It's an ad hoc estimation because interquartile range and other similar statistical outlier detection
// techniques seem to work better with human supervision, at least in this particular application.
//
// The clustering weights for the initial estimates are not retained in the Beats struct so the
// heuristics should only be used when the forgetting factor is zero i.e. all taps are equally
// weighted.
func (beats *Beats) Clean() (Beats, error) {
	// ... calculate the median taps per beat
	taps := []int{}
	for _, beat := range beats.Beats {
		if N := len(beat.Taps); N > 0 {
			taps = append(taps, N)
		}
	}

	sort.Ints(taps)

	N := len(taps)
	median := 1.0
	switch {
	case N == 0:
		median = 1.0

	case N%2 == 0:
		median = float64(taps[N/2-1]+taps[N/2+1]) / 2.0

	default:
		median = float64(taps[N/2])
	}

	// ... discard any beats with 'too few taps'
	fence := median / 3.0
	data := []float64{}
	for _, beat := range beats.Beats {
		if float64(len(beat.Taps)) >= fence {
			for _, t := range beat.Taps {
				data = append(data, t.Seconds())
			}
		}
	}

	weights := make([]float64, len(data))
	for i := 0; i < len(weights); i++ {
		weights[i] = 1.0
	}

	clusters := ckmeans.CKMeans1dDp(data, weights)

	cleaned := make([]Beat, len(clusters))
	for i, cluster := range clusters {
		cleaned[i] = makeBeat(cluster.Center, cluster)
	}

	BPM, offset := bpm(cleaned)

	sort.SliceStable(cleaned, func(i, j int) bool { return cleaned[i].At < cleaned[j].At })

	result := Beats{
		BPM:    BPM,
		Offset: offset,
		Beats:  cleaned,
	}

	if len(cleaned) > 0 {
		variance := 0.0
		for _, b := range cleaned {
			variance += b.Variance.Seconds()
		}

		variance = variance / float64(len(cleaned))

		result.Variance = &variance
	}

	return result, nil
}

// Reduces the precision of the beats to the specified time value.
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

// Subtracts a time value from all timeline values in the Beats e.g. to compensate
// for a known latency or to adjust the beats so that the first beat falls on 0.
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

// Custom JSON marshaler for the Beats struct that represents the internal times as (float) seconds.
func (beats Beats) MarshalJSON() ([]byte, error) {
	type beat struct {
		At       instant   `json:"at"`
		Mean     instant   `json:"mean"`
		Variance instant   `json:"variance"`
		Taps     []instant `json:"taps"`
	}

	b := struct {
		BPM    uint    `json:"BPM"`
		Offset instant `json:"offset"`
		Beats  []beat  `json:"beats"`
	}{
		BPM:    beats.BPM,
		Offset: instant(beats.Offset),
		Beats:  make([]beat, len(beats.Beats)),
	}

	for i, bb := range beats.Beats {
		b.Beats[i] = beat{
			At:       instant(bb.At),
			Mean:     instant(bb.Mean),
			Variance: instant(bb.Variance),
			Taps:     make([]instant, len(bb.Taps)),
		}

		for j, t := range bb.Taps {
			b.Beats[i].Taps[j] = instant(t)
		}
	}

	return json.Marshal(b)
}

// Custom JSON unmarshaler for the Beats struct that unmarshals times stored as (float) seconds.
func (beats *Beats) UnmarshalJSON(bytes []byte) error {
	if beats != nil {
		type beat struct {
			At       instant   `json:"at"`
			Mean     instant   `json:"mean"`
			Variance instant   `json:"variance"`
			Taps     []instant `json:"taps"`
		}

		b := struct {
			BPM    uint    `json:"BPM"`
			Offset instant `json:"offset"`
			Beats  []beat  `json:"beats"`
		}{}

		if err := json.Unmarshal(bytes, &b); err != nil {
			return err
		}

		beats.BPM = b.BPM
		beats.Offset = time.Duration(b.Offset)
		beats.Beats = make([]Beat, len(b.Beats))

		for i, bb := range b.Beats {
			beats.Beats[i] = Beat{
				At:       time.Duration(bb.At),
				Mean:     time.Duration(bb.Mean),
				Variance: time.Duration(bb.Variance),
				Taps:     make([]time.Duration, len(bb.Taps)),
			}

			for j, tap := range bb.Taps {
				beats.Beats[i].Taps[j] = time.Duration(tap)
			}
		}
	}

	return nil
}

// Implementation of the Stringer interface.
func (beats Beats) String() string {
	var width = 0

	for _, beat := range beats.Beats {
		if s := fmt.Sprintf("%v", beat.At); len(s) > width {
			width = len(s)
		}

		if s := fmt.Sprintf("%v", beat.Mean); len(s) > width {
			width = len(s)
		}

		if s := fmt.Sprintf("%v", beat.Variance); len(s) > width {
			width = len(s)
		}

		for _, t := range beat.Taps {
			if s := fmt.Sprintf("%v", t); len(s) > width {
				width = len(s)
			}
		}
	}

	var b bytes.Buffer

	fmt.Fprintf(&b, "BPM:    %d\n", beats.BPM)
	fmt.Fprintf(&b, "Offset: %v\n", beats.Offset)
	fmt.Fprintln(&b)
	for i, beat := range beats.Beats {
		s := ""
		s += fmt.Sprintf("%-3d", i+1)
		s += fmt.Sprintf(" %-[1]*s", width, beat.At)

		if len(beat.Taps) > 0 {
			s += fmt.Sprintf(" %-[1]*s", width, beat.Mean)
			s += fmt.Sprintf(" %-[1]*s", width, beat.Variance)
			for _, t := range beat.Taps {
				s += fmt.Sprintf(" %-[1]*s", width, t)
			}
		}

		fmt.Fprintf(&b, "%s\n", strings.TrimSpace(s))
	}

	return string(b.Bytes())
}

// Generates the weights array for a set of taps. The returned weights use 1.0 as a base value,
// with the weights of the taps in ieach discounted row being multiplied by (1.0 - forgetting).
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

// Estimate the BPM and offset of the first beats by applying least squares reqression to a set of beats.
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

// Performs a least squares reqression on a set of beats and returns the gradient
// and offset of the calculated line.
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

	m, c := regression.OrdinaryLeastSquares(x, t)

	return m, c, nil
}

// Performs a 'best guess' as to which beat number corresponds to which estimated beat,by fitting the
// beats to a line, adjusting the beats and calculating the variance between the beats and adjusted
// beats. Returns when the average variance is sufficiently low.
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
	variance := math.MaxFloat64

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

		sigmasq := sumsq / float64(N-1)

		if sigmasq < variance {
			for i := range beats {
				beats[i].beat = index[i]
			}

			variance = sigmasq
		}

		if variance < 0.001 {
			break
		}
	}

	if variance > 0.05 {
		return fmt.Errorf("Error mapping taps to beats: %v", beats)
	}

	return nil
}

// Converts a result from the ckmeans.1d.dp algorithm to a Beat.
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
