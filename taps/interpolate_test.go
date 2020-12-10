package taps

import (
	"math"
	"testing"
	"time"
)

func TestInterpolate(t *testing.T) {
	expected := Beats{
		BPM:    114,
		Offset: 316 * time.Millisecond,
		Beats: []Beat{
			beats[0], beats[1], beats[2], beats[3], beats[4], beats[5], beats[6], beats[7],
			beats[8], beats[9], beats[10], beats[11], beats[12], beats[13], beats[14], beats[15],
			beats[16], beats[17], beats[18], beats[19],
		},
	}

	data := Beats{
		BPM:    114,
		Offset: 316 * time.Millisecond,
		Beats:  []Beat{beats[8], beats[9], beats[10], beats[11], beats[12], beats[13], beats[14], beats[15]},
	}

	t2b := T2B{
		Latency:    Default.Latency,
		Forgetting: Default.Forgetting,
	}

	beats, err := t2b.Interpolate(data, Seconds(0), Seconds(10.5))
	if err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	if beats.BPM != expected.BPM {
		t.Errorf("Incorrect BPM - expected:%v, got:%v", expected.BPM, beats.BPM)
	}

	if math.Abs(beats.Offset.Seconds()-expected.Offset.Seconds()) > 0.0011 {
		t.Errorf("Incorrect offset - expected:%v, got:%v", expected.Offset, beats.Offset)
	}

	compare(beats.Beats, expected.Beats, t)
}

func TestInterpolateWithNoData(t *testing.T) {
	data := Beats{
		BPM:    0,
		Offset: 0 * time.Millisecond,
		Beats:  []Beat{},
	}

	t2b := T2B{
		Latency:    Default.Latency,
		Forgetting: Default.Forgetting,
	}

	_, err := t2b.Interpolate(data, Seconds(0), Seconds(10.5))
	if err == nil {
		t.Fatalf("Expected error, got (%v)", err)
	}
}

func TestInterpolateWithInsufficientData(t *testing.T) {
	data := Beats{
		BPM:    0,
		Offset: 0 * time.Millisecond,
		Beats:  []Beat{beats[8]},
	}

	t2b := T2B{
		Latency:    Default.Latency,
		Forgetting: Default.Forgetting,
	}

	_, err := t2b.Interpolate(data, Seconds(0), Seconds(10.5))
	if err == nil {
		t.Fatalf("Expected error, got (%v)", err)
	}
}

func TestInterpolateWithMinimalInformation(t *testing.T) {
	expected := Beats{
		BPM:    114,
		Offset: 313 * time.Millisecond,
		Beats: []Beat{
			{At: 313 * time.Millisecond},
			{At: 839 * time.Millisecond},
			{At: 1365 * time.Millisecond},
			{At: 1892 * time.Millisecond},
			{At: 2418 * time.Millisecond},
			{At: 2945 * time.Millisecond},
			{At: 3471 * time.Millisecond},
			{At: 3997 * time.Millisecond},
			beats[8],
			{At: 5050 * time.Millisecond},
			{At: 5576 * time.Millisecond},
			{At: 6103 * time.Millisecond},
			{At: 6629 * time.Millisecond},
			{At: 7155 * time.Millisecond},
			{At: 7682 * time.Millisecond},
			{At: 8208 * time.Millisecond},
			{At: 8734 * time.Millisecond},
			{At: 9261 * time.Millisecond},
			{At: 9787 * time.Millisecond},
			{At: 10313 * time.Millisecond},
		},
	}

	data := Beats{
		BPM:    114,
		Offset: 0 * time.Millisecond,
		Beats:  []Beat{beats[8]},
	}

	t2b := T2B{
		Latency:    Default.Latency,
		Forgetting: Default.Forgetting,
	}

	beats, err := t2b.Interpolate(data, Seconds(0), Seconds(10.5))
	if err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	if beats.BPM != expected.BPM {
		t.Errorf("Incorrect BPM - expected:%v, got:%v", expected.BPM, beats.BPM)
	}

	if math.Abs(beats.Offset.Seconds()-expected.Offset.Seconds()) > 0.0011 {
		t.Errorf("Incorrect offset - expected:%v, got:%v", expected.Offset, beats.Offset)
	}

	compare(beats.Beats, expected.Beats, t)
}

func TestInterpolateWithMissingBeat(t *testing.T) {
	expected := Beats{
		BPM:    114,
		Offset: 316 * time.Millisecond,
		Beats: []Beat{
			beats[0], beats[1], beats[2], beats[3], beats[4], beats[5], beats[6], beats[7],
			beats[8],
			beats[9],
			{At: 5577 * time.Millisecond},
			beats[11],
			beats[12],
			beats[13],
			beats[14],
			beats[15],
			beats[16], beats[17], beats[18], beats[19],
		},
	}

	data := Beats{
		BPM:    114,
		Offset: 316 * time.Millisecond,
		Beats:  []Beat{beats[8], beats[9], beats[11], beats[12], beats[13], beats[14], beats[15]},
	}

	t2b := T2B{
		Latency:    Default.Latency,
		Forgetting: Default.Forgetting,
	}

	beats, err := t2b.Interpolate(data, Seconds(0), Seconds(10.5))
	if err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	if beats.BPM != expected.BPM {
		t.Errorf("Incorrect BPM - expected:%v, got:%v", expected.BPM, beats.BPM)
	}

	if math.Abs(beats.Offset.Seconds()-expected.Offset.Seconds()) > 0.0011 {
		t.Errorf("Incorrect offset - expected:%v, got:%v", expected.Offset, beats.Offset)
	}

	compare(beats.Beats, expected.Beats, t)
}

// func TestInterpolateWithTwoClusters(t *testing.T) {
// 	expected := []Beat{
// 		Beat{
// 			At:       Seconds(0.306),
// 			Mean:     Seconds(0),
// 			Variance: Seconds(0),
// 			Taps:     seconds(),
// 		},
// 		Beat{
// 			At:       Seconds(1.361),
// 			Mean:     Seconds(0),
// 			Variance: Seconds(0),
// 			Taps:     seconds(),
// 		},
// 		Beat{
// 			At:       Seconds(2.415),
// 			Mean:     Seconds(0),
// 			Variance: Seconds(0),
// 			Taps:     seconds(),
// 		},
// 		Beat{
// 			At:       Seconds(3.469),
// 			Mean:     Seconds(0),
// 			Variance: Seconds(0),
// 			Taps:     seconds(),
// 		},
// 		Beat{
// 			At:       Seconds(4.524),
// 			Mean:     Seconds(4.524),
// 			Variance: Seconds(0.003),
// 			Taps:     seconds(4.570, 4.506, 4.530, 4.530, 4.518, 4.495, 4.529, 4.524, 4.518, 4.518),
// 		},
// 		Beat{
// 			At:       Seconds(5.578),
// 			Mean:     Seconds(5.578),
// 			Variance: Seconds(0.004),
// 			Taps:     seconds(5.604, 5.592, 5.592, 5.603, 5.580, 5.545, 5.563, 5.557, 5.586, 5.551, 5.586),
// 		},
// 		Beat{
// 			At:       Seconds(6.632),
// 			Mean:     Seconds(0),
// 			Variance: Seconds(0),
// 			Taps:     seconds(),
// 		},
// 		Beat{
// 			At:       Seconds(7.687),
// 			Mean:     Seconds(0),
// 			Variance: Seconds(0),
// 			Taps:     seconds(),
// 		},
// 	}
//
// 	beats, BPM, offset := extrapolate(map[int]ckmeans.Cluster{
// 		1: clusters[0],
// 		2: clusters[2],
// 	}, Seconds(0), Seconds(8.5))
//
// 	if BPM != 57 {
// 		t.Errorf("Incorrect BPM - expected:%v, got:%v", 57, BPM)
// 	}
//
// 	if math.Abs(offset.Seconds()-0.306) > 0.0011 {
// 		t.Errorf("Incorrect BPM - expected:%v, got:%v", 306*time.Millisecond, offset)
// 	}
//
// 	compare(beats, expected, t)
// }
