package taps2beats

import (
	"math"
	"testing"
	"time"
)

func TestSub(t *testing.T) {
	v := Beats{
		BPM:    114,
		Offset: Seconds(0.3162627535),
		Beats: []Beat{
			{At: Seconds(3.999)},
			{At: Seconds(4.523694381), Mean: Seconds(4.523694381), Variance: Seconds(0.003525498), Taps: seconds(bins[0]...)},
			{At: Seconds(5.057687493), Mean: Seconds(5.057687493), Variance: Seconds(0.008223081), Taps: seconds(bins[1]...)},
			{At: Seconds(5.578084204), Mean: Seconds(5.578084204), Variance: Seconds(0.004277370), Taps: seconds(bins[2]...)},
			{At: Seconds(6.100485910), Mean: Seconds(6.100485910), Variance: Seconds(0.004944514), Taps: seconds(bins[3]...)},
			{At: Seconds(6.618216081), Mean: Seconds(6.618216081), Variance: Seconds(0.007153066), Taps: seconds(bins[4]...)},
			{At: Seconds(7.153334490), Mean: Seconds(7.153334490), Variance: Seconds(0.004573754), Taps: seconds(bins[5]...)},
			{At: Seconds(7.685755996), Mean: Seconds(7.685755996), Variance: Seconds(0.005071400), Taps: seconds(bins[6]...)},
			{At: Seconds(8.210333335), Mean: Seconds(8.210333335), Variance: Seconds(0.012172972), Taps: seconds(bins[7]...)},
			{At: Seconds(8.733)},
		},
	}

	expected := Beats{
		BPM:    114,
		Offset: Seconds(0.3162627535 - 0.037),
		Beats: []Beat{
			{At: Seconds(3.999 - 0.037)},
			{
				At:       Seconds(4.523694381 - 0.037),
				Mean:     Seconds(4.523694381 - 0.037),
				Variance: Seconds(0.003525498),
				Taps: seconds(
					4.570271991-0.037,
					4.506176116-0.037,
					4.529560070-0.037,
					4.529560070-0.037,
					4.517865093-0.037,
					4.494581138-0.037,
					4.529408070-0.037,
					4.523631082-0.037,
					4.517979093-0.037,
					4.517911093-0.037)},
			{
				At:       Seconds(5.057687493 - 0.037),
				Mean:     Seconds(5.057687493 - 0.037),
				Variance: Seconds(0.008223081),
				Taps: seconds(
					5.063594027-0.037,
					5.045971061-0.037,
					5.057670039-0.037,
					5.069284016-0.037,
					5.022782107-0.037,
					5.133092891-0.037,
					5.040234073-0.037,
					5.040295073-0.037,
					5.046071061-0.037,
					5.046165061-0.037,
					5.069403016-0.037)},
			{
				At:       Seconds(5.578084204 - 0.037),
				Mean:     Seconds(5.578084204 - 0.037),
				Variance: Seconds(0.004277370),
				Taps: seconds(
					5.603539973-0.037,
					5.591722996-0.037,
					5.591721996-0.037,
					5.603428973-0.037,
					5.580101018-0.037,
					5.545395086-0.037,
					5.562732052-0.037,
					5.556940064-0.037,
					5.586102007-0.037,
					5.551068075-0.037,
					5.586174007-0.037)},

			{
				At:       Seconds(6.100485910 - 0.037),
				Mean:     Seconds(6.100485910 - 0.037),
				Variance: Seconds(0.004944514),
				Taps: seconds(
					6.102690998-0.037,
					6.114172975-0.037,
					6.137423930-0.037,
					6.102591998-0.037,
					6.096715009-0.037,
					6.067721066-0.037,
					6.079333043-0.037,
					6.131584941-0.037,
					6.090995020-0.037,
					6.073547054-0.037,
					6.108568986-0.037)},

			{
				At:       Seconds(6.618216081 - 0.037),
				Mean:     Seconds(6.618216081 - 0.037),
				Variance: Seconds(0.007153066),
				Taps: seconds(
					6.642708943-0.037,
					6.619153989-0.037,
					6.630941966-0.037,
					6.613455000-0.037,
					6.654118921-0.037,
					6.578564068-0.037,
					6.624973977-0.037,
					6.654145921-0.037,
					6.596029034-0.037,
					6.607636011-0.037,
					6.578649068-0.037)},

			{
				At:       Seconds(7.153334490 - 0.037),
				Mean:     Seconds(7.153334490 - 0.037),
				Variance: Seconds(0.004573754),
				Taps: seconds(
					7.141796968-0.037,
					7.135788980-0.037,
					7.176683900-0.037,
					7.147644957-0.037,
					7.176371900-0.037,
					7.130096991-0.037,
					7.141650968-0.037,
					7.193876866-0.037,
					7.130224991-0.037,
					7.165018923-0.037,
					7.147523957-0.037)},

			{
				At:       Seconds(7.685755996 - 0.037),
				Mean:     Seconds(7.685755996 - 0.037),
				Variance: Seconds(0.005071400),
				Taps: seconds(
					7.710649857-0.037,
					7.693071891-0.037,
					7.698974880-0.037,
					7.699120880-0.037,
					7.681405914-0.037,
					7.652464971-0.037,
					7.664070948-0.037,
					7.722112835-0.037,
					7.652501970-0.037,
					7.687334903-0.037,
					7.681606914-0.037)},

			{
				At:       Seconds(8.210333335 - 0.037),
				Mean:     Seconds(8.210333335 - 0.037),
				Variance: Seconds(0.012172972),
				Taps: seconds(
					8.192470916-0.037,
					8.203885893-0.037,
					8.227207848-0.037,
					8.215609871-0.037,
					8.215537871-0.037,
					8.134273030-0.037,
					8.198270905-0.037,
					8.244539814-0.037,
					8.180805939-0.037,
					8.238953824-0.037,
					8.262110780-0.037)},
			{At: Seconds(8.733 - 0.037)},
		},
	}

	v.Sub(37 * time.Millisecond)

	if v.BPM != expected.BPM {
		t.Errorf("Incorrect BPM - expected:%v, got:%v", expected.BPM, v.BPM)
	}

	if v.Offset != expected.Offset {
		t.Errorf("Incorrect offset - expected:%v, got:%v", expected.Offset, v.Offset)
	}

	if len(v.Beats) != len(expected.Beats) {
		t.Errorf("Incorrect beats - expected:%v, got:%v", len(expected.Beats), len(v.Beats))
	} else {
		for i, x := range expected.Beats {
			b := v.Beats[i]

			if b.At != x.At {
				t.Errorf("Beat %d - incorrect 'at' - expected:%v, got:%v", i+1, x.At, b.At)
			}

			if b.Mean != x.Mean {
				t.Errorf("Beat %d - incorrect 'mean' - expected:%v, got:%v", i+1, x.Mean, b.Mean)
			}

			if b.Variance != x.Variance {
				t.Errorf("Beat %d - incorrect 'variance' - expected:%v, got:%v", i+1, x.Variance, b.Variance)
			}

			if len(b.Taps) != len(x.Taps) {
				t.Errorf("Incorrect taps - expected:%v, got:%v", len(x.Taps), len(b.Taps))
			} else {
				for j, tap := range x.Taps {
					if math.Abs(b.Taps[j].Seconds()-tap.Seconds()) > 0.00001 {
						t.Errorf("Beat %d - incorrect 'tap %d' - expected:%v, got:%v", i+1, j+1, tap, b.Taps[j])
					}
				}
			}
		}
	}
}
