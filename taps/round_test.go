package taps

import (
	"testing"
	"time"
)

func TestRound(t *testing.T) {
	v := Beats{
		BPM:    114,
		Offset: Seconds(0.3162627535),
		Beats: []Beat{
			{At: Seconds(4.523694381), Mean: Seconds(4.523694381), Variance: Seconds(0.003525498), Taps: seconds(bins[0]...)},
			{At: Seconds(5.057687493), Mean: Seconds(5.057687493), Variance: Seconds(0.008223081), Taps: seconds(bins[1]...)},
			{At: Seconds(5.578084204), Mean: Seconds(5.578084204), Variance: Seconds(0.004277370), Taps: seconds(bins[2]...)},
			{At: Seconds(6.100485910), Mean: Seconds(6.100485910), Variance: Seconds(0.004944514), Taps: seconds(bins[3]...)},
			{At: Seconds(6.618216081), Mean: Seconds(6.618216081), Variance: Seconds(0.007153066), Taps: seconds(bins[4]...)},
			{At: Seconds(7.153334490), Mean: Seconds(7.153334490), Variance: Seconds(0.004573754), Taps: seconds(bins[5]...)},
			{At: Seconds(7.685755996), Mean: Seconds(7.685755996), Variance: Seconds(0.005071400), Taps: seconds(bins[6]...)},
			{At: Seconds(8.210333335), Mean: Seconds(8.210333335), Variance: Seconds(0.012172972), Taps: seconds(bins[7]...)},
		},
	}

	expected := Beats{
		BPM:    114,
		Offset: 316 * time.Millisecond,
		Beats: []Beat{
			{
				At:       4524 * time.Millisecond,
				Mean:     4524 * time.Millisecond,
				Variance: 4 * time.Millisecond,
				Taps: []time.Duration{
					4570 * time.Millisecond,
					4506 * time.Millisecond,
					4530 * time.Millisecond,
					4530 * time.Millisecond,
					4518 * time.Millisecond,
					4495 * time.Millisecond,
					4529 * time.Millisecond,
					4524 * time.Millisecond,
					4518 * time.Millisecond,
					4518 * time.Millisecond}},
			{
				At:       5058 * time.Millisecond,
				Mean:     5058 * time.Millisecond,
				Variance: 8 * time.Millisecond,
				Taps: []time.Duration{
					5064 * time.Millisecond,
					5046 * time.Millisecond,
					5058 * time.Millisecond,
					5069 * time.Millisecond,
					5023 * time.Millisecond,
					5133 * time.Millisecond,
					5040 * time.Millisecond,
					5040 * time.Millisecond,
					5046 * time.Millisecond,
					5046 * time.Millisecond,
					5069 * time.Millisecond}},
			{
				At:       5578 * time.Millisecond,
				Mean:     5578 * time.Millisecond,
				Variance: 4 * time.Millisecond,
				Taps: []time.Duration{
					5604 * time.Millisecond,
					5592 * time.Millisecond,
					5592 * time.Millisecond,
					5603 * time.Millisecond,
					5580 * time.Millisecond,
					5545 * time.Millisecond,
					5563 * time.Millisecond,
					5557 * time.Millisecond,
					5586 * time.Millisecond,
					5551 * time.Millisecond,
					5586 * time.Millisecond}},

			{
				At:       6100 * time.Millisecond,
				Mean:     6100 * time.Millisecond,
				Variance: 5 * time.Millisecond,
				Taps: []time.Duration{
					6103 * time.Millisecond,
					6114 * time.Millisecond,
					6137 * time.Millisecond,
					6103 * time.Millisecond,
					6097 * time.Millisecond,
					6068 * time.Millisecond,
					6079 * time.Millisecond,
					6132 * time.Millisecond,
					6091 * time.Millisecond,
					6074 * time.Millisecond,
					6109 * time.Millisecond}},
			{
				At:       6618 * time.Millisecond,
				Mean:     6618 * time.Millisecond,
				Variance: 7 * time.Millisecond,
				Taps: []time.Duration{
					6643 * time.Millisecond,
					6619 * time.Millisecond,
					6631 * time.Millisecond,
					6613 * time.Millisecond,
					6654 * time.Millisecond,
					6579 * time.Millisecond,
					6625 * time.Millisecond,
					6654 * time.Millisecond,
					6596 * time.Millisecond,
					6608 * time.Millisecond,
					6579 * time.Millisecond}},
			{
				At:       7153 * time.Millisecond,
				Mean:     7153 * time.Millisecond,
				Variance: 5 * time.Millisecond,
				Taps: []time.Duration{
					7142 * time.Millisecond,
					7136 * time.Millisecond,
					7177 * time.Millisecond,
					7148 * time.Millisecond,
					7176 * time.Millisecond,
					7130 * time.Millisecond,
					7142 * time.Millisecond,
					7194 * time.Millisecond,
					7130 * time.Millisecond,
					7165 * time.Millisecond,
					7148 * time.Millisecond}},
			{
				At:       7686 * time.Millisecond,
				Mean:     7686 * time.Millisecond,
				Variance: 5 * time.Millisecond,
				Taps: []time.Duration{
					7711 * time.Millisecond,
					7693 * time.Millisecond,
					7699 * time.Millisecond,
					7699 * time.Millisecond,
					7681 * time.Millisecond,
					7652 * time.Millisecond,
					7664 * time.Millisecond,
					7722 * time.Millisecond,
					7653 * time.Millisecond,
					7687 * time.Millisecond,
					7682 * time.Millisecond}},
			{
				At:       8210 * time.Millisecond,
				Mean:     8210 * time.Millisecond,
				Variance: 12 * time.Millisecond,
				Taps: []time.Duration{
					8192 * time.Millisecond,
					8204 * time.Millisecond,
					8227 * time.Millisecond,
					8216 * time.Millisecond,
					8216 * time.Millisecond,
					8134 * time.Millisecond,
					8198 * time.Millisecond,
					8245 * time.Millisecond,
					8181 * time.Millisecond,
					8239 * time.Millisecond,
					8262 * time.Millisecond}},
		},
	}

	v.Round(1 * time.Millisecond)

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
					if b.Taps[j] != tap {
						t.Errorf("Beat %d - incorrect 'tap %d' - expected:%v, got:%v", i+1, j+1, tap, b.Taps[j])
					}
				}
			}
		}
	}
}

func TestRoundTo10ms(t *testing.T) {
	v := Beats{
		BPM:    114,
		Offset: Seconds(0.3162627535),
		Beats: []Beat{
			{At: Seconds(4.523694381), Mean: Seconds(4.523694381), Variance: Seconds(0.003525498), Taps: seconds(bins[0]...)},
			{At: Seconds(5.057687493), Mean: Seconds(5.057687493), Variance: Seconds(0.008223081), Taps: seconds(bins[1]...)},
			{At: Seconds(5.578084204), Mean: Seconds(5.578084204), Variance: Seconds(0.004277370), Taps: seconds(bins[2]...)},
			{At: Seconds(6.100485910), Mean: Seconds(6.100485910), Variance: Seconds(0.004944514), Taps: seconds(bins[3]...)},
			{At: Seconds(6.618216081), Mean: Seconds(6.618216081), Variance: Seconds(0.007153066), Taps: seconds(bins[4]...)},
			{At: Seconds(7.153334490), Mean: Seconds(7.153334490), Variance: Seconds(0.004573754), Taps: seconds(bins[5]...)},
			{At: Seconds(7.685755996), Mean: Seconds(7.685755996), Variance: Seconds(0.005071400), Taps: seconds(bins[6]...)},
			{At: Seconds(8.210333335), Mean: Seconds(8.210333335), Variance: Seconds(0.012172972), Taps: seconds(bins[7]...)},
		},
	}

	expected := Beats{
		BPM:    114,
		Offset: 320 * time.Millisecond,
		Beats: []Beat{
			{
				At:       4520 * time.Millisecond,
				Mean:     4520 * time.Millisecond,
				Variance: 0 * time.Millisecond,
				Taps: []time.Duration{
					4570 * time.Millisecond,
					4510 * time.Millisecond,
					4530 * time.Millisecond,
					4530 * time.Millisecond,
					4520 * time.Millisecond,
					4490 * time.Millisecond,
					4530 * time.Millisecond,
					4520 * time.Millisecond,
					4520 * time.Millisecond,
					4520 * time.Millisecond}},
			{
				At:       5060 * time.Millisecond,
				Mean:     5060 * time.Millisecond,
				Variance: 10 * time.Millisecond,
				Taps: []time.Duration{
					5060 * time.Millisecond,
					5050 * time.Millisecond,
					5060 * time.Millisecond,
					5070 * time.Millisecond,
					5020 * time.Millisecond,
					5130 * time.Millisecond,
					5040 * time.Millisecond,
					5040 * time.Millisecond,
					5050 * time.Millisecond,
					5050 * time.Millisecond,
					5070 * time.Millisecond}},
			{
				At:       5580 * time.Millisecond,
				Mean:     5580 * time.Millisecond,
				Variance: 0 * time.Millisecond,
				Taps: []time.Duration{
					5600 * time.Millisecond,
					5590 * time.Millisecond,
					5590 * time.Millisecond,
					5600 * time.Millisecond,
					5580 * time.Millisecond,
					5550 * time.Millisecond,
					5560 * time.Millisecond,
					5560 * time.Millisecond,
					5590 * time.Millisecond,
					5550 * time.Millisecond,
					5590 * time.Millisecond}},

			{
				At:       6100 * time.Millisecond,
				Mean:     6100 * time.Millisecond,
				Variance: 0 * time.Millisecond,
				Taps: []time.Duration{
					6100 * time.Millisecond,
					6110 * time.Millisecond,
					6140 * time.Millisecond,
					6100 * time.Millisecond,
					6100 * time.Millisecond,
					6070 * time.Millisecond,
					6080 * time.Millisecond,
					6130 * time.Millisecond,
					6090 * time.Millisecond,
					6070 * time.Millisecond,
					6110 * time.Millisecond}},
			{
				At:       6620 * time.Millisecond,
				Mean:     6620 * time.Millisecond,
				Variance: 10 * time.Millisecond,
				Taps: []time.Duration{
					6640 * time.Millisecond,
					6620 * time.Millisecond,
					6630 * time.Millisecond,
					6610 * time.Millisecond,
					6650 * time.Millisecond,
					6580 * time.Millisecond,
					6620 * time.Millisecond,
					6650 * time.Millisecond,
					6600 * time.Millisecond,
					6610 * time.Millisecond,
					6580 * time.Millisecond}},
			{
				At:       7150 * time.Millisecond,
				Mean:     7150 * time.Millisecond,
				Variance: 0 * time.Millisecond,
				Taps: []time.Duration{
					7140 * time.Millisecond,
					7140 * time.Millisecond,
					7180 * time.Millisecond,
					7150 * time.Millisecond,
					7180 * time.Millisecond,
					7130 * time.Millisecond,
					7140 * time.Millisecond,
					7190 * time.Millisecond,
					7130 * time.Millisecond,
					7170 * time.Millisecond,
					7150 * time.Millisecond}},
			{
				At:       7690 * time.Millisecond,
				Mean:     7690 * time.Millisecond,
				Variance: 10 * time.Millisecond,
				Taps: []time.Duration{
					7710 * time.Millisecond,
					7690 * time.Millisecond,
					7700 * time.Millisecond,
					7700 * time.Millisecond,
					7680 * time.Millisecond,
					7650 * time.Millisecond,
					7660 * time.Millisecond,
					7720 * time.Millisecond,
					7650 * time.Millisecond,
					7690 * time.Millisecond,
					7680 * time.Millisecond}},
			{
				At:       8210 * time.Millisecond,
				Mean:     8210 * time.Millisecond,
				Variance: 10 * time.Millisecond,
				Taps: []time.Duration{
					8190 * time.Millisecond,
					8200 * time.Millisecond,
					8230 * time.Millisecond,
					8220 * time.Millisecond,
					8220 * time.Millisecond,
					8130 * time.Millisecond,
					8200 * time.Millisecond,
					8240 * time.Millisecond,
					8180 * time.Millisecond,
					8240 * time.Millisecond,
					8260 * time.Millisecond}},
		},
	}

	v.Round(10 * time.Millisecond)

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
					if b.Taps[j] != tap {
						t.Errorf("Beat %d - incorrect 'tap %d' - expected:%v, got:%v", i+1, j+1, tap, b.Taps[j])
					}
				}
			}
		}
	}
}
