package taps

import (
	"reflect"
	"testing"
	"time"
)

func TestReindexWithNoBeats(t *testing.T) {
	expected := []Beat{}
	beats := []Beat{}

	if err := reindex(beats); err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	if len(beats) != len(expected) {
		t.Errorf("Invalid result\n   expected: %v beats\n   got:      %v beats", len(expected), len(beats))
	}
}

func TestReindexWithOneBeat(t *testing.T) {
	expected := []Beat{
		{beat: 1, At: Seconds(4.523694381), Mean: Seconds(4.523694381), Variance: Seconds(0.000391722), Taps: seconds(bins[0]...)},
	}

	beats := []Beat{
		{At: Seconds(4.523694381), Mean: Seconds(4.523694381), Variance: Seconds(0.000391722), Taps: seconds(bins[0]...)},
	}

	if err := reindex(beats); err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	if len(beats) != len(expected) {
		t.Errorf("Invalid result\n   expected: %v beats\n   got:      %v beats", len(expected), len(beats))
	}
}

func TestReindexWithTwoBeats(t *testing.T) {
	expected := []Beat{
		{beat: 1, At: Seconds(4.523694381), Mean: Seconds(4.523694381), Variance: Seconds(0.000391722), Taps: seconds(bins[0]...)},
		{beat: 2, At: Seconds(8.210333335), Mean: Seconds(8.210333335), Variance: Seconds(0.001217297), Taps: seconds(bins[7]...)},
	}

	beats := []Beat{
		{At: Seconds(4.523694381), Mean: Seconds(4.523694381), Variance: Seconds(0.000391722), Taps: seconds(bins[0]...)},
		{At: Seconds(8.210333335), Mean: Seconds(8.210333335), Variance: Seconds(0.001217297), Taps: seconds(bins[7]...)},
	}

	if err := reindex(beats); err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	for i, x := range expected {
		if beats[i].beat != x.beat {
			t.Errorf("Invalid interpolation [%d] - expected:%v, got:%v", i+1, x.beat, beats[i].beat)
		}
	}
}

func TestReindexWithThreeBeats(t *testing.T) {
	samples := [][]int{
		{1, 2, 3},
		{1, 2, 4}, {1, 3, 4},
		{1, 2, 5}, {1, 3, 5}, {1, 4, 5},
		{1, 2, 6}, {1, 3, 6}, {1, 4, 6}, {1, 5, 6},
		{1, 2, 7}, {1, 3, 7}, {1, 4, 7}, {1, 5, 7}, {1, 6, 7},
		{1, 2, 8}, {1, 3, 8}, {1, 4, 8}, {1, 5, 8}, {1, 6, 8}, {1, 7, 8},
	}

	expected := [][]int{
		{1, 2, 3},
		{1, 2, 4},
		{1, 3, 4},
		{1, 2, 5},
		{1, 2, 3}, // {1, 3, 5}
		{1, 4, 5},
		{1, 2, 6},
		{1, 3, 6},
		{1, 4, 6},
		{1, 5, 6},
		{1, 2, 7},
		{1, 2, 4}, // {1, 3, 7},
		{1, 2, 3}, // {1, 4, 7},
		{1, 3, 4}, // {1, 5, 7},
		{1, 6, 7},
		{1, 2, 8},
		{1, 3, 8},
		{1, 4, 8},
		{1, 5, 8},
		{1, 6, 8},
		{1, 7, 8},
	}

	for i, v := range samples {
		p := []Beat{}
		for _, ix := range v {
			p = append(p, clone(beats[8+ix-1]))
		}

		err := reindex(p)
		if err != nil {
			t.Fatalf("[%d] unexpected error (%v)", i+1, err)
		}

		for i, x := range expected[i] {
			if p[i].beat != x {
				t.Errorf("Invalid interpolation [%d] - expected:%v, got:%v", i+1, x, p[i].beat)
			}
		}
	}
}

func TestReindexWithFourBeats(t *testing.T) {
	samples := [][]int{
		{1, 2, 3, 4}, {1, 2, 3, 5}, {1, 2, 3, 6}, {1, 2, 3, 7}, {1, 2, 3, 8},
		{1, 2, 4, 5}, {1, 2, 4, 6}, {1, 2, 4, 7}, {1, 2, 4, 8},
		{1, 2, 5, 6}, {1, 2, 5, 7}, {1, 2, 5, 8},
		{1, 2, 6, 7}, {1, 2, 6, 8},
		{1, 2, 7, 8},
		{1, 3, 4, 5}, {1, 3, 4, 6}, {1, 3, 4, 7}, {1, 3, 4, 8},
		{1, 3, 5, 6}, {1, 3, 5, 7}, {1, 3, 5, 8},
		{1, 3, 6, 7}, {1, 3, 6, 8},
		{1, 3, 7, 8},
		{1, 4, 5, 6}, {1, 4, 5, 7}, {1, 4, 5, 8},
		{1, 4, 6, 7}, {1, 4, 6, 8},
		{1, 4, 7, 8},
		{1, 5, 6, 7}, {1, 5, 6, 8},
		{1, 5, 7, 8},
		{1, 6, 7, 8},
	}

	expected := [][]int{
		{1, 2, 3, 4}, {1, 2, 3, 5}, {1, 2, 3, 6}, {1, 2, 3, 7}, {1, 2, 3, 8},
		{1, 2, 4, 5}, {1, 2, 4, 6}, {1, 2, 4, 7}, {1, 2, 4, 8},
		{1, 2, 5, 6}, {1, 2, 5, 7}, {1, 2, 5, 8},
		{1, 2, 6, 7}, {1, 2, 6, 8},
		{1, 2, 7, 8},
		{1, 3, 4, 5}, {1, 3, 4, 6}, {1, 3, 4, 7}, {1, 3, 4, 8},
		{1, 3, 5, 6} /* {1, 3, 5, 7}, */, {1, 2, 3, 4}, {1, 3, 5, 8},
		{1, 3, 6, 7}, {1, 3, 6, 8},
		{1, 3, 7, 8},
		{1, 4, 5, 6}, {1, 4, 5, 7}, {1, 4, 5, 8},
		{1, 4, 6, 7}, {1, 4, 6, 8},
		{1, 4, 7, 8},
		{1, 5, 6, 7}, {1, 5, 6, 8},
		{1, 5, 7, 8},
		{1, 6, 7, 8},
	}

	for i, v := range samples {
		p := []Beat{}
		for _, ix := range v {
			p = append(p, clone(beats[8+ix-1]))
		}

		err := reindex(p)
		if err != nil {
			t.Fatalf("[%d] unexpected error (%v)", i+1, err)
		}

		for i, x := range expected[i] {
			if p[i].beat != x {
				t.Errorf("Invalid interpolation [%d] - expected:%v, got:%v", i+1, x, p[i].beat)
			}
		}
	}
}

func TestReindexWithCombinations(t *testing.T) {
	exceptions := [][][]int{
		{{1, 3, 5}, {1, 2, 3}},
		{{1, 3, 7}, {1, 2, 4}},
		{{1, 4, 7}, {1, 2, 3}},
		{{1, 5, 7}, {1, 3, 4}},
		{{1, 3, 5, 7}, {1, 2, 3, 4}},
	}

	test := func(v []int) {
		p := []Beat{}
		for _, ix := range v {
			p = append(p, clone(beats[8+ix-1]))
		}

		err := reindex(p)
		if err != nil {
			t.Fatalf("[%v] unexpected error (%v)", v, err)
		}

		expected := make([]int, len(v))

		copy(expected, v)
		for _, x := range exceptions {
			if reflect.DeepEqual(v, x[0]) {
				copy(expected, x[1])
			}
		}

		for i, x := range expected {
			if p[i].beat != x {
				t.Errorf("Invalid interpolation [%d] - expected:%v, got:%v", i+1, x, p[i].beat)
			}
		}
	}

	K := []int{2, 3, 4, 5, 6, 7}
	suffix := []int{2, 3, 4, 5, 6, 7, 8}

	for _, k := range K {
		combinations(k, []int{1}, suffix, test)
	}
}

func TestReindexWithPathologicalData(t *testing.T) {
	beats := []Beat{
		Beat{At: Seconds(1.0)},
		Beat{At: Seconds(1.1)},
		Beat{At: Seconds(11.0)},
	}

	if err := reindex(beats); err == nil {
		t.Fatalf("Expected error, got %v", err)
	}
}

func combinations(k int, head, tail []int, f func([]int)) {
	if k > 0 {
		for i, v := range tail {
			h := append(head, v)
			combinations(k-1, h, tail[i+1:], f)
		}
		return
	}

	f(head)
}

func clone(beat Beat) Beat {
	b := Beat{
		beat:     beat.beat,
		At:       beat.At,
		Mean:     beat.Mean,
		Variance: beat.Variance,
		Taps:     make([]time.Duration, len(beat.Taps)),
	}

	copy(b.Taps, beat.Taps)

	return b
}
