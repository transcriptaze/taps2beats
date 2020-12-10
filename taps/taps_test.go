package taps

import (
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/twystd/ckmeans.1d.dp/go/ckmeans"
)

var taps = [][]float64{
	{4.570271991, 5.063594027, 5.603539973, 6.102690998, 6.642708943, 7.141796968, 7.710649857, 8.192470916},
	{4.506176116, 5.045971061, 5.591722996, 6.114172975, 6.619153989, 7.135788980, 7.693071891, 8.203885893},
	{4.529560070, 5.057670039, 5.591721996, 6.137423930, 6.630941966, 7.176683900, 7.698974880, 8.227207848},
	{4.529560070, 5.069284016, 5.603428973, 6.102591998, 6.613455000, 7.147644957, 7.699120880, 8.215609871},
	{4.517865093, 5.022782107, 5.580101018, 6.096715009, 6.654118921, 7.176371900, 7.681405914, 8.215537871},
	{5.133092891, 5.545395086, 6.067721066, 6.578564068, 7.130096991, 7.652464971, 8.134273030},
	{4.494581138, 5.040234073, 5.562732052, 6.079333043, 6.624973977, 7.141650968, 7.664070948, 8.198270905},
	{4.529408070, 5.040295073, 5.556940064, 6.131584941, 6.654145921, 7.193876866, 7.722112835, 8.244539814},
	{4.523631082, 5.046071061, 5.586102007, 6.090995020, 6.596029034, 7.130224991, 7.652501971, 8.180805939},
	{4.517979093, 5.046165061, 5.551068075, 6.073547054, 6.607636011, 7.165018923, 7.687334903, 8.238953825},
	{4.517911093, 5.069403016, 5.586174007, 6.108568986, 6.578649068, 7.147523957, 7.681606914, 8.262110780},
}

var bins = [][]float64{
	{4.570271991, 4.506176116, 4.529560070, 4.529560070, 4.517865093, 4.494581138, 4.529408070, 4.523631082, 4.517979093, 4.517911093},
	{5.063594027, 5.045971061, 5.057670039, 5.069284016, 5.022782107, 5.133092891, 5.040234073, 5.040295073, 5.046071061, 5.046165061, 5.069403016},
	{5.603539973, 5.591722996, 5.591721996, 5.603428973, 5.580101018, 5.545395086, 5.562732052, 5.556940064, 5.586102007, 5.551068075, 5.586174007},
	{6.102690998, 6.114172975, 6.137423930, 6.102591998, 6.096715009, 6.067721066, 6.079333043, 6.131584941, 6.090995020, 6.073547054, 6.108568986},
	{6.642708943, 6.619153989, 6.630941966, 6.613455000, 6.654118921, 6.578564068, 6.624973977, 6.654145921, 6.596029034, 6.607636011, 6.578649068},
	{7.141796968, 7.135788980, 7.176683900, 7.147644957, 7.176371900, 7.130096991, 7.141650968, 7.193876866, 7.130224991, 7.165018923, 7.147523957},
	{7.710649857, 7.693071891, 7.698974880, 7.699120880, 7.681405914, 7.652464971, 7.664070948, 7.722112835, 7.652501970, 7.687334903, 7.681606914},
	{8.192470916, 8.203885893, 8.227207848, 8.215609871, 8.215537871, 8.134273030, 8.198270905, 8.244539814, 8.180805939, 8.238953824, 8.262110780},
}

var clusters = []ckmeans.Cluster{
	{Center: 4.523694381, Variance: 0.003525498, Values: bins[0]},
	{Center: 5.057687493, Variance: 0.008223081, Values: bins[1]},
	{Center: 5.578084204, Variance: 0.004277370, Values: bins[2]},
	{Center: 6.100485910, Variance: 0.004944514, Values: bins[3]},
	{Center: 6.618216081, Variance: 0.007153066, Values: bins[4]},
	{Center: 7.153334490, Variance: 0.004573754, Values: bins[5]},
	{Center: 7.685755996, Variance: 0.005071400, Values: bins[6]},
	{Center: 8.210333335, Variance: 0.012172972, Values: bins[7]},
}

var beats = []Beat{
	{At: Seconds(0.316)},
	{At: Seconds(0.842)},
	{At: Seconds(1.368)},
	{At: Seconds(1.894)},
	{At: Seconds(2.420)},
	{At: Seconds(2.947)},
	{At: Seconds(3.472)},
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
	{At: Seconds(9.260)},
	{At: Seconds(9.786)},
	{At: Seconds(10.312)},
}

var quantized = []Beat{
	{At: Seconds(0.316)},
	{At: Seconds(0.842)},
	{At: Seconds(1.368)},
	{At: Seconds(1.894)},
	{At: Seconds(2.420)},
	{At: Seconds(2.946)},
	{At: Seconds(3.472)},
	{At: Seconds(3.998)},
	{At: Seconds(4.525), Mean: Seconds(4.523694381), Variance: Seconds(0.003525498), Taps: seconds(bins[0]...)},
	{At: Seconds(5.051), Mean: Seconds(5.057687493), Variance: Seconds(0.008223081), Taps: seconds(bins[1]...)},
	{At: Seconds(5.577), Mean: Seconds(5.578084204), Variance: Seconds(0.004277370), Taps: seconds(bins[2]...)},
	{At: Seconds(6.103), Mean: Seconds(6.100485910), Variance: Seconds(0.004944514), Taps: seconds(bins[3]...)},
	{At: Seconds(6.629), Mean: Seconds(6.618216081), Variance: Seconds(0.007153066), Taps: seconds(bins[4]...)},
	{At: Seconds(7.155), Mean: Seconds(7.153334490), Variance: Seconds(0.004573754), Taps: seconds(bins[5]...)},
	{At: Seconds(7.681), Mean: Seconds(7.685755996), Variance: Seconds(0.005071400), Taps: seconds(bins[6]...)},
	{At: Seconds(8.207), Mean: Seconds(8.210333335), Variance: Seconds(0.012172972), Taps: seconds(bins[7]...)},
	{At: Seconds(8.733)},
	{At: Seconds(9.260)},
	{At: Seconds(9.786)},
	{At: Seconds(10.312)},
}

func TestTaps2Beats(t *testing.T) {
	expected := Beats{
		BPM:    114,
		Offset: 316 * time.Millisecond,
		Beats:  []Beat{beats[8], beats[9], beats[10], beats[11], beats[12], beats[13], beats[14], beats[15]},
	}

	t2b := T2B{
		Latency:    Default.Latency,
		Forgetting: Default.Forgetting,
	}

	beats := t2b.Taps2Beats(Floats2Seconds(taps))

	if beats.BPM != expected.BPM {
		t.Errorf("Incorrect BPM - expected:%v, got:%v", expected.BPM, beats.BPM)
	}

	if math.Abs(beats.Offset.Seconds()-expected.Offset.Seconds()) > 0.0011 {
		t.Errorf("Incorrect offset - expected:%v, got:%v", expected.Offset, beats.Offset)
	}

	compare(beats.Beats, expected.Beats, t)
}

func TestTaps2BeatsWithMissingBeat(t *testing.T) {
	expected := Beats{
		BPM:    114,
		Offset: 316 * time.Millisecond,
		Beats:  []Beat{beats[8], beats[9], beats[10], beats[11], beats[12], beats[14], beats[15]},
	}

	taps := [][]float64{
		{4.570271991, 5.063594027, 5.603539973, 6.102690998, 6.642708943, 7.710649857, 8.192470916},
		{4.506176116, 5.045971061, 5.591722996, 6.114172975, 6.619153989, 7.693071891, 8.203885893},
		{4.52956007, 5.057670039, 5.591721996, 6.13742393, 6.630941966, 7.69897488, 8.227207848},
		{4.52956007, 5.069284016, 5.603428973, 6.102591998, 6.613455, 7.69912088, 8.215609871},
		{4.517865093, 5.022782107, 5.580101018, 6.096715009, 6.654118921, 7.681405914, 8.215537871},
		{5.133092891, 5.545395086, 6.067721066, 6.578564068, 7.652464971, 8.13427303},
		{4.494581138, 5.040234073, 5.562732052, 6.079333043, 6.624973977, 7.664070948, 8.198270905},
		{4.52940807, 5.040295073, 5.556940064, 6.131584941, 6.654145921, 7.722112835, 8.244539814},
		{4.523631082, 5.046071061, 5.586102007, 6.09099502, 6.596029034, 7.652501971, 8.180805939},
		{4.517979093, 5.046165061, 5.551068075, 6.073547054, 6.607636011, 7.687334903, 8.238953825},
		{4.517911093, 5.069403016, 5.586174007, 6.108568986, 6.578649068, 7.681606914, 8.26211078},
	}

	t2b := T2B{
		Latency:    Default.Latency,
		Forgetting: Default.Forgetting,
	}

	beats := t2b.Taps2Beats(Floats2Seconds(taps))

	if beats.BPM != expected.BPM {
		t.Errorf("Incorrect BPM - expected:%v, got:%v", expected.BPM, beats.BPM)
	}

	if math.Abs(beats.Offset.Seconds()-expected.Offset.Seconds()) > 0.0011 {
		t.Errorf("Incorrect offset - expected:%v, got:%v", expected.Offset, beats.Offset)
	}

	compare(beats.Beats, expected.Beats, t)
}

func TestTaps2BeatsWithWeirdData(t *testing.T) {
	expected := Beats{
		BPM:    58,
		Offset: 978 * time.Millisecond,
		Beats: []Beat{
			{At: Seconds(1.000), Mean: Seconds(1), Variance: Seconds(0.1), Taps: seconds(1, 1.1, 1.2, 0.9, 0.8)},
			{At: Seconds(2.000), Mean: Seconds(2), Variance: Seconds(0.1), Taps: seconds(2, 2.1, 2.2, 1.9, 1.8)},
			{At: Seconds(50.000), Mean: Seconds(50), Variance: Seconds(0.100), Taps: seconds(50, 50.1, 50.2, 49.9, 49.8)},
		},
	}

	taps := [][]float64{
		{1.0, 1.1, 1.2, 0.9, 0.8},
		{2.0, 2.1, 2.2, 1.9, 1.8},
		{50.0, 50.1, 50.2, 49.9, 49.8},
	}

	t2b := T2B{
		Latency:    Default.Latency,
		Forgetting: Default.Forgetting,
	}

	beats := t2b.Taps2Beats(Floats2Seconds(taps))

	if beats.BPM != expected.BPM {
		t.Errorf("Incorrect BPM - expected:%v, got:%v", expected.BPM, beats.BPM)
	}

	if math.Abs(beats.Offset.Seconds()-expected.Offset.Seconds()) > 0.0011 {
		t.Errorf("Incorrect offset - expected:%v, got:%v", expected.Offset, beats.Offset)
	}

	compare(beats.Beats, expected.Beats, t)
}

func TestTaps2BeatsWithLatency(t *testing.T) {
	expected := Beats{
		BPM:    114,
		Offset: 279 * time.Millisecond,
		Beats: []Beat{
			{
				At:       4487 * time.Millisecond,
				Mean:     4487 * time.Millisecond,
				Variance: 3 * time.Millisecond,
				Taps: seconds(
					4.57027199-0.037,
					4.50617612-0.037,
					4.52956007-0.037,
					4.52956007-0.037,
					4.51786509-0.037,
					4.49458114-0.037,
					4.52940807-0.037,
					4.52363108-0.037,
					4.51797909-0.037,
					4.51791109-0.037)},
			{
				At:       5020 * time.Millisecond,
				Mean:     5020 * time.Millisecond,
				Variance: 8 * time.Millisecond,
				Taps: seconds(
					5.06359402-0.037,
					5.04597106-0.037,
					5.05767004-0.037,
					5.06928402-0.037,
					5.02278211-0.037,
					5.13309289-0.037,
					5.04023407-0.037,
					5.04029507-0.037,
					5.04607106-0.037,
					5.04616506-0.037,
					5.06940302-0.037)},
			{
				At:       5541 * time.Millisecond,
				Mean:     5541 * time.Millisecond,
				Variance: 4 * time.Millisecond,
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
				At:       6064 * time.Millisecond,
				Mean:     6064 * time.Millisecond,
				Variance: 4 * time.Millisecond,
				Taps: seconds(
					6.102690998-0.037,
					6.114172975-0.037,
					6.13742393-0.037,
					6.102591998-0.037,
					6.096715009-0.037,
					6.067721066-0.037,
					6.079333043-0.037,
					6.131584941-0.037,
					6.09099502-0.037,
					6.073547054-0.037,
					6.108568986-0.037)},
			{
				At:       6581 * time.Millisecond,
				Mean:     6581 * time.Millisecond,
				Variance: 7 * time.Millisecond,
				Taps: seconds(
					6.642708943-0.037,
					6.619153989-0.037,
					6.630941966-0.037,
					6.613455-0.037,
					6.654118921-0.037,
					6.578564068-0.037,
					6.624973977-0.037,
					6.654145921-0.037,
					6.596029034-0.037,
					6.607636011-0.037,
					6.578649068-0.037)},
			{
				At:       7116 * time.Millisecond,
				Mean:     7116 * time.Millisecond,
				Variance: 4 * time.Millisecond,
				Taps: seconds(
					7.14179697-0.037,
					7.13578898-0.037,
					7.17668390-0.037,
					7.14764496-0.037,
					7.17637190-0.037,
					7.13009699-0.037,
					7.14165097-0.037,
					7.19387687-0.037,
					7.13022499-0.037,
					7.16501892-0.037,
					7.14752397-0.037)},
			{
				At:       7648 * time.Millisecond,
				Mean:     7648 * time.Millisecond,
				Variance: 5 * time.Millisecond,
				Taps: seconds(
					7.710649857-0.037,
					7.693071891-0.037,
					7.69897488-0.037,
					7.69912088-0.037,
					7.681405914-0.037,
					7.652464971-0.037,
					7.664070948-0.037,
					7.722112835-0.037,
					7.652501970-0.037,
					7.687334903-0.037,
					7.681606914-0.037)},
			{
				At:       8173 * time.Millisecond,
				Mean:     8173 * time.Millisecond,
				Variance: 12 * time.Millisecond,
				Taps: seconds(
					8.192470916-0.037,
					8.203885893-0.037,
					8.227207848-0.037,
					8.215609871-0.037,
					8.215537871-0.037,
					8.13427303-0.037,
					8.198270905-0.037,
					8.244539814-0.037,
					8.180805939-0.037,
					8.238953824-0.037,
					8.26211078-0.037)},
		},
	}

	t2b := T2B{
		Latency:    37 * time.Millisecond,
		Forgetting: Default.Forgetting,
	}

	beats := t2b.Taps2Beats(Floats2Seconds(taps))

	if beats.BPM != expected.BPM {
		t.Errorf("Incorrect BPM - expected:%v, got:%v", expected.BPM, beats.BPM)
	}

	if math.Abs(beats.Offset.Seconds()-expected.Offset.Seconds()) > 0.0011 {
		t.Errorf("Incorrect BPM - expected:%v, got:%v", expected.Offset, beats.Offset)
	}

	compare(beats.Beats, expected.Beats, t)
}

func TestTaps2BeatsWithForgetting(t *testing.T) {
	expected := Beats{
		BPM:    114,
		Offset: 309 * time.Millisecond,
		Beats: []Beat{
			{At: 4521 * time.Millisecond, Mean: 4521 * time.Millisecond, Variance: 1535 * time.Microsecond, Taps: seconds(bins[0]...)},
			{At: 5057 * time.Millisecond, Mean: 5057 * time.Millisecond, Variance: 5000 * time.Microsecond, Taps: seconds(bins[1]...)},
			{At: 5575 * time.Millisecond, Mean: 5575 * time.Millisecond, Variance: 2507 * time.Microsecond, Taps: seconds(bins[2]...)},
			{At: 6099 * time.Millisecond, Mean: 6099 * time.Millisecond, Variance: 4000 * time.Microsecond, Taps: seconds(bins[3]...)},
			{At: 6614 * time.Millisecond, Mean: 6614 * time.Millisecond, Variance: 4000 * time.Microsecond, Taps: seconds(bins[4]...)},
			{At: 7153 * time.Millisecond, Mean: 7153 * time.Millisecond, Variance: 2892 * time.Microsecond, Taps: seconds(bins[5]...)},
			{At: 7683 * time.Millisecond, Mean: 7683 * time.Millisecond, Variance: 3000 * time.Microsecond, Taps: seconds(bins[6]...)},
			{At: 8215 * time.Millisecond, Mean: 8215 * time.Millisecond, Variance: 9000 * time.Microsecond, Taps: seconds(bins[7]...)},
		},
	}

	t2b := T2B{
		Latency:    Default.Latency,
		Forgetting: 0.1,
	}

	beats := t2b.Taps2Beats(Floats2Seconds(taps))

	if beats.BPM != expected.BPM {
		t.Errorf("Incorrect BPM - expected:%v, got:%v", expected.BPM, beats.BPM)
	}

	if math.Abs(beats.Offset.Seconds()-expected.Offset.Seconds()) > 0.0011 {
		t.Errorf("Incorrect offset - expected:%v, got:%v", expected.Offset, beats.Offset)
	}

	compare(beats.Beats, expected.Beats, t)
}

func TestReindexForDegenerateCases(t *testing.T) {
	tests := [][]Beat{
		{},
		{beats[8]},
		{beats[8], beats[15]},
	}

	expected := []map[int]Beat{
		{},
		{1: beats[8]},
		{1: beats[8], 2: beats[15]},
	}

	for i, v := range tests {
		beats, err := reindex(v)
		if err != nil {
			t.Fatalf("[%d] unexpected error (%v)", i+1, err)
		}

		if !reflect.DeepEqual(beats, expected[i]) {
			t.Errorf("[%d] invalid interpolation\n   expected: %v\n   got:      %v", i+1, expected[i], beats)
		}
	}
}

func TestReindexForThreeBeats(t *testing.T) {
	samples := [][]int{
		{1, 2, 3},
		{1, 2, 4}, {1, 3, 4},
		{1, 2, 5}, {1, 3, 5}, {1, 4, 5},
		{1, 2, 6}, {1, 3, 6}, {1, 4, 6}, {1, 5, 6},
		{1, 2, 7}, {1, 3, 7}, {1, 4, 7}, {1, 5, 7}, {1, 6, 7},
		{1, 2, 8}, {1, 3, 8}, {1, 4, 8}, {1, 5, 8}, {1, 6, 8}, {1, 7, 8},
	}

	expected := []map[int]Beat{
		{1: beats[8], 2: beats[9], 3: beats[10]},
		{1: beats[8], 2: beats[9], 4: beats[11]},
		{1: beats[8], 3: beats[10], 4: beats[11]},
		{1: beats[8], 2: beats[9], 5: beats[12]},
		{1: beats[8], 2: beats[10], 3: beats[12]}, // {1, 3, 5}
		{1: beats[8], 4: beats[11], 5: beats[12]},
		{1: beats[8], 2: beats[9], 6: beats[13]},
		{1: beats[8], 3: beats[10], 6: beats[13]},
		{1: beats[8], 4: beats[11], 6: beats[13]},
		{1: beats[8], 5: beats[12], 6: beats[13]},
		{1: beats[8], 2: beats[9], 7: beats[14]},
		{1: beats[8], 2: beats[10], 4: beats[14]}, // {1, 3, 7},
		{1: beats[8], 2: beats[11], 3: beats[14]}, // {1, 4, 7},
		{1: beats[8], 3: beats[12], 4: beats[14]}, // {1, 5, 7},
		{1: beats[8], 6: beats[13], 7: beats[14]},
		{1: beats[8], 2: beats[9], 8: beats[15]},
		{1: beats[8], 3: beats[10], 8: beats[15]},
		{1: beats[8], 4: beats[11], 8: beats[15]},
		{1: beats[8], 5: beats[12], 8: beats[15]},
		{1: beats[8], 6: beats[13], 8: beats[15]},
		{1: beats[8], 7: beats[14], 8: beats[15]},
	}
	for i, v := range samples {
		c := []Beat{}
		for _, ix := range v {
			c = append(c, beats[8+ix-1])
		}

		beats, err := reindex(c)
		if err != nil {
			t.Fatalf("[%d] unexpected error (%v)", i+1, err)
		}

		if !reflect.DeepEqual(beats, expected[i]) {
			t.Errorf("Invalid interpolation [%d] - expected:%v, got:%v", i+1, expected[i], beats)
		}
	}
}

func TestInterpolateForFourBeats(t *testing.T) {
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

	expected := []map[int]Beat{
		{1: beats[8], 2: beats[9], 3: beats[10], 4: beats[11]},
		{1: beats[8], 2: beats[9], 3: beats[10], 5: beats[12]},
		{1: beats[8], 2: beats[9], 3: beats[10], 6: beats[13]},
		{1: beats[8], 2: beats[9], 3: beats[10], 7: beats[14]},
		{1: beats[8], 2: beats[9], 3: beats[10], 8: beats[15]},

		{1: beats[8], 2: beats[9], 4: beats[11], 5: beats[12]},
		{1: beats[8], 2: beats[9], 4: beats[11], 6: beats[13]},
		{1: beats[8], 2: beats[9], 4: beats[11], 7: beats[14]},
		{1: beats[8], 2: beats[9], 4: beats[11], 8: beats[15]},

		{1: beats[8], 2: beats[9], 5: beats[12], 6: beats[13]},
		{1: beats[8], 2: beats[9], 5: beats[12], 7: beats[14]},
		{1: beats[8], 2: beats[9], 5: beats[12], 8: beats[15]},

		{1: beats[8], 2: beats[9], 6: beats[13], 7: beats[14]},
		{1: beats[8], 2: beats[9], 6: beats[13], 8: beats[15]},

		{1: beats[8], 2: beats[9], 7: beats[14], 8: beats[15]},

		{1: beats[8], 3: beats[10], 4: beats[11], 5: beats[12]},
		{1: beats[8], 3: beats[10], 4: beats[11], 6: beats[13]},
		{1: beats[8], 3: beats[10], 4: beats[11], 7: beats[14]},
		{1: beats[8], 3: beats[10], 4: beats[11], 8: beats[15]},

		{1: beats[8], 3: beats[10], 5: beats[12], 6: beats[13]}, /* {1, 3, 5, 7}, */
		{1: beats[8], 2: beats[10], 3: beats[12], 4: beats[14]},
		{1: beats[8], 3: beats[10], 5: beats[12], 8: beats[15]},

		{1: beats[8], 3: beats[10], 6: beats[13], 7: beats[14]},
		{1: beats[8], 3: beats[10], 6: beats[13], 8: beats[15]},

		{1: beats[8], 3: beats[10], 7: beats[14], 8: beats[15]},

		{1: beats[8], 4: beats[11], 5: beats[12], 6: beats[13]},
		{1: beats[8], 4: beats[11], 5: beats[12], 7: beats[14]},
		{1: beats[8], 4: beats[11], 5: beats[12], 8: beats[15]},

		{1: beats[8], 4: beats[11], 6: beats[13], 7: beats[14]},
		{1: beats[8], 4: beats[11], 6: beats[13], 8: beats[15]},

		{1: beats[8], 4: beats[11], 7: beats[14], 8: beats[15]},

		{1: beats[8], 5: beats[12], 6: beats[13], 7: beats[14]},
		{1: beats[8], 5: beats[12], 6: beats[13], 8: beats[15]},

		{1: beats[8], 5: beats[12], 7: beats[14], 8: beats[15]},

		{1: beats[8], 6: beats[13], 7: beats[14], 8: beats[15]},
	}

	for i, v := range samples {
		c := []Beat{}
		for _, ix := range v {
			c = append(c, beats[8+ix-1])
		}

		beats, err := reindex(c)
		if err != nil {
			t.Fatalf("[%d] unexpected error (%v)", i+1, err)
		}

		if !reflect.DeepEqual(beats, expected[i]) {
			t.Errorf("Invalid interpolation [%d] - expected:%v, got:%v", i+1, expected[i], beats)
		}
	}
}

func TestTapCombinations(t *testing.T) {
	exceptions := [][][]int{
		{{1, 3, 5}, {1, 2, 3}},
		{{1, 3, 7}, {1, 2, 4}},
		{{1, 4, 7}, {1, 2, 3}},
		{{1, 5, 7}, {1, 3, 4}},
		{{1, 3, 5, 7}, {1, 2, 3, 4}},
	}

	test := func(v []int) {
		c := make([]Beat, len(v))
		for i, ix := range v {
			c[i] = beats[8+ix-1]
		}

		expected := map[int]Beat{}
		for _, ix := range v {
			expected[ix] = beats[8+ix-1]
		}

		index, err := reindex(c)
		if err != nil {
			t.Fatalf("[%v] unexpected error (%v)", v, err)
		}

		for _, x := range exceptions {
			if reflect.DeepEqual(v, x[0]) {
				expected = map[int]Beat{}
				for i := range x[0] {
					expected[x[1][i]] = beats[8+x[0][i]-1]
				}
			}
		}

		if !reflect.DeepEqual(index, expected) {
			t.Errorf("[%v] invalid interpolation - expected:%v, got:%v", v, expected, beats)
		}
	}

	K := []int{2, 3, 4, 5, 6, 7}
	suffix := []int{2, 3, 4, 5, 6, 7, 8}

	for _, k := range K {
		combinations(k, []int{1}, suffix, test)
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

func TestInterpolateWithPathologicalData(t *testing.T) {
	clusters := []Beat{
		Beat{At: Seconds(1.0)},
		Beat{At: Seconds(1.1)},
		Beat{At: Seconds(11.0)},
	}

	_, err := reindex(clusters)
	if err == nil {
		t.Fatalf("Expected error, got %v", err)
	}
}

func seconds(floats ...float64) []time.Duration {
	l := []time.Duration{}

	for _, f := range floats {
		l = append(l, Seconds(f))
	}

	return l
}

func compare(beats, expected []Beat, t *testing.T) {
	if len(beats) != len(expected) {
		t.Errorf("Invalid result\n   expected: %v beats\n   got:      %v beats", len(expected), len(beats))
		return
	}

	for i, v := range expected {
		if !reflect.DeepEqual(v, beats[i]) {
			if math.Abs(beats[i].At.Seconds()-v.At.Seconds()) >= 0.0011 {
				t.Errorf("Invalid beat %d 'at' - expected:%v, got:%v (delta:%.4f)", i+1, v.At, beats[i].At, math.Abs(beats[i].At.Seconds()-v.At.Seconds()))
			}

			if math.Abs(beats[i].Mean.Seconds()-v.Mean.Seconds()) >= 0.0011 {
				t.Errorf("Invalid beat %d 'mean' - expected:%v, got:%v (delta:%.4f)", i+1, v.Mean, beats[i].Mean, math.Abs(beats[i].Mean.Seconds()-v.Mean.Seconds()))
			}

			if math.Abs(beats[i].Variance.Seconds()-v.Variance.Seconds()) >= 0.0011 {
				t.Errorf("Invalid beat %d 'variance' - expected:%v, got:%v (delta:%.4f)", i+1, v.Variance, beats[i].Variance, math.Abs(beats[i].Mean.Seconds()-v.Mean.Seconds()))
			}

			if !reflect.DeepEqual(v.Taps, beats[i].Taps) {
				if len(beats[i].Taps) != len(v.Taps) {
					t.Errorf("Invalid beat %d 'taps'\n   expected: %v\n   got:      %v", i+1, v.Taps, beats[i].Taps)
					continue
				}

				for j := range v.Taps {
					if math.Abs(beats[i].Taps[j].Seconds()-v.Taps[j].Seconds()) >= 0.0011 {
						t.Errorf("Invalid beat %d 'taps'\n   expected: %v\n   got:      %v (delta:%.4f)", i+1, v.Taps, beats[i].Taps, math.Abs(beats[i].Taps[j].Seconds()-v.Taps[j].Seconds()))
						break
					}
				}
			}
		}
	}
}
