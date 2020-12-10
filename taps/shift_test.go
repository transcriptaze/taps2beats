package taps

import (
	"testing"
	"time"
)

func TestShift(t *testing.T) {
	expected := Beats{
		BPM:    123,
		Offset: 0 * time.Millisecond,
		Beats: []Beat{
			{
				At:       0 * time.Millisecond,
				Mean:     -1 * time.Millisecond,
				Variance: 3 * time.Millisecond,
				Taps:     seconds(0.045, -0.018, 0.004, 0.004, -0.007, -0.030, 0.004, -0.001, -0.007, -0.007)},
			{
				At:       526 * time.Millisecond,
				Mean:     532 * time.Millisecond,
				Variance: 8 * time.Millisecond,
				Taps:     seconds(0.538, 0.520, 0.532, 0.544, 0.497, 0.608, 0.515, 0.515, 0.521, 0.521, 0.544)},
			{
				At:       1052 * time.Millisecond,
				Mean:     1053 * time.Millisecond,
				Variance: 4 * time.Millisecond,
				Taps:     seconds(1.078, 1.067, 1.067, 1.078, 1.055, 1.020, 1.038, 1.032, 1.061, 1.026, 1.061)},
			{
				At:       1578 * time.Millisecond,
				Mean:     1576 * time.Millisecond,
				Variance: 4 * time.Millisecond,
				Taps:     seconds(1.5777, 1.589, 1.612, 1.578, 1.572, 1.543, 1.554, 1.607, 1.566, 1.549, 1.584)},
			{
				At:       2104 * time.Millisecond,
				Mean:     2093 * time.Millisecond,
				Variance: 7 * time.Millisecond,
				Taps:     seconds(2.118, 2.094, 2.106, 2.088, 2.129, 2.054, 2.100, 2.129, 2.071, 2.083, 2.054)},
			{
				At:       2630 * time.Millisecond,
				Mean:     2628 * time.Millisecond,
				Variance: 4 * time.Millisecond,
				Taps:     seconds(2.617, 2.611, 2.652, 2.623, 2.651, 2.605, 2.617, 2.669, 2.605, 2.640, 2.623)},
			{
				At:       3156 * time.Millisecond,
				Mean:     3160 * time.Millisecond,
				Variance: 5 * time.Millisecond,
				Taps:     seconds(3.186, 3.168, 3.174, 3.174, 3.156, 3.127, 3.139, 3.197, 3.128, 3.162, 3.157)},
			{
				At:       3682 * time.Millisecond,
				Mean:     3685 * time.Millisecond,
				Variance: 12 * time.Millisecond,
				Taps:     seconds(3.667, 3.679, 3.702, 3.690, 3.691, 3.609, 3.673, 3.720, 3.656, 3.714, 3.737)},
		},
	}

	beats := Beats{
		BPM:    123,
		Offset: 117 * time.Millisecond,
		Beats: []Beat{
			quantized[8], quantized[10], quantized[9], quantized[11], quantized[13], quantized[12], quantized[15], quantized[14],
		},
	}

	t2b := T2B{
		Latency:    Default.Latency,
		Forgetting: Default.Forgetting,
	}

	beats = t2b.Shift(beats)

	if beats.BPM != expected.BPM {
		t.Errorf("BPM unexpectedly modified - expected:%v, got:%v", expected.BPM, beats.BPM)
	}

	if beats.Offset != expected.Offset {
		t.Errorf("Incorrect offset - expected:%v, got:%v", expected.Offset, beats.Offset)
	}

	compare(beats.Beats, expected.Beats, t)
}

func TestShiftWithNoData(t *testing.T) {
	expected := Beats{
		BPM:    123,
		Offset: 0 * time.Millisecond,
		Beats:  []Beat{},
	}

	beats := Beats{
		BPM:    123,
		Offset: 117 * time.Millisecond,
		Beats:  []Beat{},
	}

	t2b := T2B{
		Latency:    Default.Latency,
		Forgetting: Default.Forgetting,
	}

	beats = t2b.Shift(beats)

	if beats.BPM != expected.BPM {
		t.Errorf("BPM unexpectedly modified - expected:%v, got:%v", expected.BPM, beats.BPM)
	}

	if beats.Offset != expected.Offset {
		t.Errorf("Incorrect offset - expected:%v, got:%v", expected.Offset, beats.Offset)
	}

	compare(beats.Beats, expected.Beats, t)
}

func TestShiftWithExtrapolatedData(t *testing.T) {
	expected := Beats{
		BPM:    123,
		Offset: 0 * time.Millisecond,
		Beats: []Beat{
			{At: 0 * time.Millisecond},
			{At: 526 * time.Millisecond},
			{
				At:       1053 * time.Millisecond,
				Mean:     1052 * time.Millisecond,
				Variance: 3 * time.Millisecond,
				Taps:     seconds(1.098, 1.034, 1.058, 1.058, 1.046, 1.023, 1.057, 1.052, 1.046, 1.046)},
			{
				At:       1579 * time.Millisecond,
				Mean:     1585 * time.Millisecond,
				Variance: 8 * time.Millisecond,
				Taps:     seconds(1.592, 1.574, 1.586, 1.597, 1.551, 1.661, 1.568, 1.568, 1.574, 1.574, 1.597)},
			{
				At:       2105 * time.Millisecond,
				Mean:     2106 * time.Millisecond,
				Variance: 4 * time.Millisecond,
				Taps:     seconds(2.131, 2.120, 2.120, 2.131, 2.108, 2.073, 2.091, 2.085, 2.114, 2.079, 2.1141)},
			{
				At:       2631 * time.Millisecond,
				Mean:     2629 * time.Millisecond,
				Variance: 4 * time.Millisecond,
				Taps:     seconds(2.631, 2.642, 2.665, 2.631, 2.625, 2.596, 2.607, 2.660, 2.619, 2.602, 2.637)},
			{
				At:       3157 * time.Millisecond,
				Mean:     3146 * time.Millisecond,
				Variance: 7 * time.Millisecond,
				Taps:     seconds(3.171, 3.147, 3.159, 3.141, 3.182, 3.107, 3.153, 3.182, 3.124, 3.136, 3.1066)},
			{
				At:       3683 * time.Millisecond,
				Mean:     3681 * time.Millisecond,
				Variance: 4 * time.Millisecond,
				Taps:     seconds(3.670, 3.664, 3.705, 3.676, 3.704, 3.658, 3.670, 3.722, 3.658, 3.693, 3.676)},
			{
				At:       4209 * time.Millisecond,
				Mean:     4213 * time.Millisecond,
				Variance: 5 * time.Millisecond,
				Taps:     seconds(4.239, 4.221, 4.227, 4.227, 4.209, 4.180, 4.192, 4.250, 4.181, 4.215, 4.210)},
			{
				At:       4735 * time.Millisecond,
				Mean:     4738 * time.Millisecond,
				Variance: 12 * time.Millisecond,
				Taps:     seconds(4.7205, 4.732, 4.755, 4.744, 4.744, 4.662, 4.726, 4.773, 4.709, 4.767, 4.790)},
			{At: 5261 * time.Millisecond},
			{At: 5788 * time.Millisecond},
		},
	}

	beats := Beats{
		BPM:    123,
		Offset: 117 * time.Millisecond,
		Beats: []Beat{
			quantized[6], quantized[7],
			quantized[8], quantized[10], quantized[9], quantized[11], quantized[13], quantized[12], quantized[15], quantized[14],
			quantized[16], quantized[17],
		},
	}

	t2b := T2B{
		Latency:    Default.Latency,
		Forgetting: Default.Forgetting,
	}

	beats = t2b.Shift(beats)

	if beats.BPM != expected.BPM {
		t.Errorf("BPM unexpectedly modified - expected:%v, got:%v", expected.BPM, beats.BPM)
	}

	if beats.Offset != expected.Offset {
		t.Errorf("Incorrect offset - expected:%v, got:%v", expected.Offset, beats.Offset)
	}

	compare(beats.Beats, expected.Beats, t)
}
