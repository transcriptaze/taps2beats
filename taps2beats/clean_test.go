package taps2beats

import (
	"math"
	"testing"
	"time"
)

func TestCleanWithNothingToDo(t *testing.T) {
	expected := Beats{
		BPM:    114,
		Offset: 316 * time.Millisecond,
		Beats:  []Beat{beats[8], beats[9], beats[10], beats[11], beats[12], beats[13], beats[14], beats[15]},
	}

	beats := Beats{
		BPM:    114,
		Offset: 316 * time.Millisecond,
		Beats: []Beat{
			{
				At:       Seconds(4.523694381),
				Mean:     Seconds(4.523694381),
				Variance: Seconds(0.000391722),
				Taps:     seconds(4.570271991, 4.506176116, 4.529560070, 4.529560070, 4.517865093, 4.494581138, 4.529408070, 4.523631082, 4.517979093, 4.517911093)},
			{
				At:       Seconds(5.057687493),
				Mean:     Seconds(5.057687493),
				Variance: Seconds(0.000822308),
				Taps:     seconds(5.063594027, 5.045971061, 5.057670039, 5.069284016, 5.022782107, 5.133092891, 5.040234073, 5.040295073, 5.046071061, 5.046165061, 5.069403016)},
			{
				At:       Seconds(5.578084204),
				Mean:     Seconds(5.578084204),
				Variance: Seconds(0.000427737),
				Taps:     seconds(5.603539973, 5.591722996, 5.591721996, 5.603428973, 5.580101018, 5.545395086, 5.562732052, 5.556940064, 5.586102007, 5.551068075, 5.586174007)},
			{
				At:       Seconds(6.100485910),
				Mean:     Seconds(6.100485910),
				Variance: Seconds(0.000494451),
				Taps:     seconds(6.102690998, 6.114172975, 6.137423930, 6.102591998, 6.096715009, 6.067721066, 6.079333043, 6.131584941, 6.090995020, 6.073547054, 6.108568986)},
			{
				At:       Seconds(6.618216081),
				Mean:     Seconds(6.618216081),
				Variance: Seconds(0.000715306),
				Taps:     seconds(6.642708943, 6.619153989, 6.630941966, 6.613455000, 6.654118921, 6.578564068, 6.624973977, 6.654145921, 6.596029034, 6.607636011, 6.578649068)},
			{
				At:       Seconds(7.153334490),
				Mean:     Seconds(7.153334490),
				Variance: Seconds(0.000457375),
				Taps:     seconds(7.141796968, 7.135788980, 7.176683900, 7.147644957, 7.176371900, 7.130096991, 7.141650968, 7.193876866, 7.130224991, 7.165018923, 7.147523957)},
			{
				At:       Seconds(7.685755996),
				Mean:     Seconds(7.685755996),
				Variance: Seconds(0.000507140),
				Taps:     seconds(7.710649857, 7.693071891, 7.698974880, 7.699120880, 7.681405914, 7.652464971, 7.664070948, 7.722112835, 7.652501970, 7.687334903, 7.681606914)},
			{
				At:       Seconds(8.210333335),
				Mean:     Seconds(8.210333335),
				Variance: Seconds(0.001217297),
				Taps:     seconds(8.192470916, 8.203885893, 8.227207848, 8.215609871, 8.215537871, 8.134273030, 8.198270905, 8.244539814, 8.180805939, 8.238953824, 8.262110780)},
		},
	}

	cleaned, err := beats.Clean()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if cleaned.BPM != expected.BPM {
		t.Errorf("Incorrect BPM - expected:%v, got:%v", expected.BPM, cleaned.BPM)
	}

	if math.Abs(cleaned.Offset.Seconds()-expected.Offset.Seconds()) > 0.0011 {
		t.Errorf("Incorrect offset - expected:%v, got:%v", expected.Offset, cleaned.Offset)
	}

	compare(cleaned.Beats, expected.Beats, t)
}

func TestCleanWithSingleTapBeat(t *testing.T) {
	expected := Beats{
		BPM:    114,
		Offset: 316 * time.Millisecond,
		Beats: []Beat{
			beats[8],
			beats[9],
			beats[10],
			beats[12],
			beats[13],
			beats[14],
			beats[15]},
	}

	beats := Beats{
		BPM:    114,
		Offset: 316 * time.Millisecond,
		Beats: []Beat{
			{
				At:       Seconds(4.523694381),
				Mean:     Seconds(4.523694381),
				Variance: Seconds(0.000391722),
				Taps:     seconds(4.570271991, 4.506176116, 4.529560070, 4.529560070, 4.517865093, 4.494581138, 4.529408070, 4.523631082, 4.517979093, 4.517911093)},
			{
				At:       Seconds(5.057687493),
				Mean:     Seconds(5.057687493),
				Variance: Seconds(0.000822308),
				Taps:     seconds(5.063594027, 5.045971061, 5.057670039, 5.069284016, 5.022782107, 5.133092891, 5.040234073, 5.040295073, 5.046071061, 5.046165061, 5.069403016)},
			{
				At:       Seconds(5.578084204),
				Mean:     Seconds(5.578084204),
				Variance: Seconds(0.000427737),
				Taps:     seconds(5.603539973, 5.591722996, 5.591721996, 5.603428973, 5.580101018, 5.545395086, 5.562732052, 5.556940064, 5.586102007, 5.551068075, 5.586174007)},
			{
				At:       Seconds(6.100485910),
				Mean:     Seconds(6.100485910),
				Variance: Seconds(0.000494451),
				Taps:     seconds(6.100485910),
			},
			{
				At:       Seconds(6.618216081),
				Mean:     Seconds(6.618216081),
				Variance: Seconds(0.000715306),
				Taps:     seconds(6.642708943, 6.619153989, 6.630941966, 6.613455000, 6.654118921, 6.578564068, 6.624973977, 6.654145921, 6.596029034, 6.607636011, 6.578649068)},
			{
				At:       Seconds(7.153334490),
				Mean:     Seconds(7.153334490),
				Variance: Seconds(0.000457375),
				Taps:     seconds(7.141796968, 7.135788980, 7.176683900, 7.147644957, 7.176371900, 7.130096991, 7.141650968, 7.193876866, 7.130224991, 7.165018923, 7.147523957)},
			{
				At:       Seconds(7.685755996),
				Mean:     Seconds(7.685755996),
				Variance: Seconds(0.000507140),
				Taps:     seconds(7.710649857, 7.693071891, 7.698974880, 7.699120880, 7.681405914, 7.652464971, 7.664070948, 7.722112835, 7.652501970, 7.687334903, 7.681606914)},
			{
				At:       Seconds(8.210333335),
				Mean:     Seconds(8.210333335),
				Variance: Seconds(0.001217297),
				Taps:     seconds(8.192470916, 8.203885893, 8.227207848, 8.215609871, 8.215537871, 8.134273030, 8.198270905, 8.244539814, 8.180805939, 8.238953824, 8.262110780)},
		},
	}

	cleaned, err := beats.Clean()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if cleaned.BPM != expected.BPM {
		t.Errorf("Incorrect BPM - expected:%v, got:%v", expected.BPM, cleaned.BPM)
	}

	if math.Abs(cleaned.Offset.Seconds()-expected.Offset.Seconds()) > 0.0011 {
		t.Errorf("Incorrect offset - expected:%v, got:%v", expected.Offset, cleaned.Offset)
	}

	compare(cleaned.Beats, expected.Beats, t)

	t.Errorf("OOOPS %v\n", cleaned)
}
