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
	{4.506176116, 5.045971061, 5.591722996, 6.114172975, 6.619153989, 7.13578898, 7.693071891, 8.203885893},
	{4.52956007, 5.057670039, 5.591721996, 6.13742393, 6.630941966, 7.1766839, 7.69897488, 8.227207848},
	{4.52956007, 5.069284016, 5.603428973, 6.102591998, 6.613455, 7.147644957, 7.69912088, 8.215609871},
	{4.517865093, 5.022782107, 5.580101018, 6.096715009, 6.654118921, 7.1763719, 7.681405914, 8.215537871},
	{5.133092891, 5.545395086, 6.067721066, 6.578564068, 7.130096991, 7.652464971, 8.13427303},
	{4.494581138, 5.040234073, 5.562732052, 6.079333043, 6.624973977, 7.141650968, 7.664070948, 8.198270905},
	{4.52940807, 5.040295073, 5.556940064, 6.131584941, 6.654145921, 7.193876866, 7.722112835, 8.244539814},
	{4.523631082, 5.046071061, 5.586102007, 6.09099502, 6.596029034, 7.130224991, 7.652501971, 8.180805939},
	{4.517979093, 5.046165061, 5.551068075, 6.073547054, 6.607636011, 7.165018923, 7.687334903, 8.238953825},
	{4.517911093, 5.069403016, 5.586174007, 6.108568986, 6.578649068, 7.147523957, 7.681606914, 8.26211078},
}

var bins = [][]float64{
	{4.570271991, 4.506176116, 4.52956007, 4.52956007, 4.517865093, 4.494581138, 4.52940807, 4.523631082, 4.517979093, 4.517911093},
	{5.063594027, 5.045971061, 5.057670039, 5.069284016, 5.022782107, 5.133092891, 5.040234073, 5.040295073, 5.046071061, 5.046165061, 5.069403016},
	{5.603539973, 5.591722996, 5.591721996, 5.603428973, 5.580101018, 5.545395086, 5.562732052, 5.556940064, 5.586102007, 5.551068075, 5.586174007},
	{6.102690998, 6.114172975, 6.13742393, 6.102591998, 6.096715009, 6.067721066, 6.079333043, 6.131584941, 6.09099502, 6.073547054, 6.108568986},
	{6.642708943, 6.619153989, 6.630941966, 6.613455, 6.654118921, 6.578564068, 6.624973977, 6.654145921, 6.596029034, 6.607636011, 6.578649068},
	{7.141796968, 7.13578898, 7.1766839, 7.147644957, 7.1763719, 7.130096991, 7.141650968, 7.193876866, 7.130224991, 7.165018923, 7.147523957},
	{7.710649857, 7.693071891, 7.69897488, 7.69912088, 7.681405914, 7.652464971, 7.664070948, 7.722112835, 7.652501970, 7.687334903, 7.681606914},
	{8.192470916, 8.203885893, 8.227207848, 8.215609871, 8.215537871, 8.13427303, 8.198270905, 8.244539814, 8.180805939, 8.238953824, 8.26211078},
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
	{At: 316 * time.Millisecond},
	{At: 842 * time.Millisecond},
	{At: 1368 * time.Millisecond},
	{At: 1894 * time.Millisecond},
	{At: 2420 * time.Millisecond},
	{At: 2947 * time.Millisecond},
	{At: 3472 * time.Millisecond},
	{At: 3999 * time.Millisecond},
	{At: Seconds(4.523694381), Mean: 4524 * time.Millisecond, Variance: 3 * time.Millisecond, Taps: Floats2Seconds(bins)[0]},
	{At: Seconds(5.057687493), Mean: 5057 * time.Millisecond, Variance: 8 * time.Millisecond, Taps: Floats2Seconds(bins)[1]},
	{At: Seconds(5.578084204), Mean: 5578 * time.Millisecond, Variance: 4 * time.Millisecond, Taps: Floats2Seconds(bins)[2]},
	{At: Seconds(6.100485910), Mean: 6101 * time.Millisecond, Variance: 4 * time.Millisecond, Taps: Floats2Seconds(bins)[3]},
	{At: Seconds(6.618216081), Mean: 6618 * time.Millisecond, Variance: 7 * time.Millisecond, Taps: Floats2Seconds(bins)[4]},
	{At: Seconds(7.153334490), Mean: 7153 * time.Millisecond, Variance: 4 * time.Millisecond, Taps: Floats2Seconds(bins)[5]},
	{At: Seconds(7.685755996), Mean: 7685 * time.Millisecond, Variance: 5 * time.Millisecond, Taps: Floats2Seconds(bins)[6]},
	{At: Seconds(8.210333335), Mean: 8210 * time.Millisecond, Variance: 12 * time.Millisecond, Taps: Floats2Seconds(bins)[7]},
	{At: 8733 * time.Millisecond},
	{At: 9260 * time.Millisecond},
	{At: 9786 * time.Millisecond},
	{At: 10312 * time.Millisecond},
}

var quantized = []Beat{
	{At: 316 * time.Millisecond},
	{At: 842 * time.Millisecond},
	{At: 1368 * time.Millisecond},
	{At: 1894 * time.Millisecond},
	{At: 2420 * time.Millisecond},
	{At: 2946 * time.Millisecond},
	{At: 3472 * time.Millisecond},
	{At: 3998 * time.Millisecond},
	{At: 4525 * time.Millisecond, Mean: 4524 * time.Millisecond, Variance: 3 * time.Millisecond, Taps: Floats2Seconds(bins)[0]},
	{At: 5051 * time.Millisecond, Mean: 5057 * time.Millisecond, Variance: 8 * time.Millisecond, Taps: Floats2Seconds(bins)[1]},
	{At: 5577 * time.Millisecond, Mean: 5578 * time.Millisecond, Variance: 4 * time.Millisecond, Taps: Floats2Seconds(bins)[2]},
	{At: 6103 * time.Millisecond, Mean: 6101 * time.Millisecond, Variance: 4 * time.Millisecond, Taps: Floats2Seconds(bins)[3]},
	{At: 6629 * time.Millisecond, Mean: 6618 * time.Millisecond, Variance: 7 * time.Millisecond, Taps: Floats2Seconds(bins)[4]},
	{At: 7155 * time.Millisecond, Mean: 7153 * time.Millisecond, Variance: 4 * time.Millisecond, Taps: Floats2Seconds(bins)[5]},
	{At: 7681 * time.Millisecond, Mean: 7685 * time.Millisecond, Variance: 5 * time.Millisecond, Taps: Floats2Seconds(bins)[6]},
	{At: 8207 * time.Millisecond, Mean: 8210 * time.Millisecond, Variance: 12 * time.Millisecond, Taps: Floats2Seconds(bins)[7]},
	{At: 8733 * time.Millisecond},
	{At: 9260 * time.Millisecond},
	{At: 9786 * time.Millisecond},
	{At: 10312 * time.Millisecond},
}

func TestTaps2Beats(t *testing.T) {
	expected := Beats{
		BPM:    114,
		Offset: 316 * time.Millisecond,
		Beats:  []Beat{beats[8], beats[9], beats[10], beats[11], beats[12], beats[13], beats[14], beats[15]},
	}

	t2b := T2B{
		Precision:  Default.Precision,
		Latency:    Default.Latency,
		Forgetting: Default.Forgetting,
	}

	beats := t2b.Taps2Beats(Floats2Seconds(taps), 0, Seconds(10.5))

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
		Precision:  Default.Precision,
		Latency:    Default.Latency,
		Forgetting: Default.Forgetting,
	}

	beats := t2b.Taps2Beats(Floats2Seconds(taps), 0.0, Seconds(10.5))

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
		Precision:  Default.Precision,
		Latency:    Default.Latency,
		Forgetting: Default.Forgetting,
	}

	beats := t2b.Taps2Beats(Floats2Seconds(taps), 0, Seconds(60.0))

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
		Offset: (316 - 37) * time.Millisecond,
		Beats: []Beat{
			{At: (4524 - 37) * time.Millisecond, Mean: (4524 - 37) * time.Millisecond, Variance: 3 * time.Millisecond, Taps: Floats2Seconds(bins)[0]},
			{At: (5057 - 37) * time.Millisecond, Mean: (5057 - 37) * time.Millisecond, Variance: 8 * time.Millisecond, Taps: Floats2Seconds(bins)[1]},
			{At: (5578 - 37) * time.Millisecond, Mean: (5578 - 37) * time.Millisecond, Variance: 4 * time.Millisecond, Taps: Floats2Seconds(bins)[2]},
			{At: (6101 - 37) * time.Millisecond, Mean: (6101 - 37) * time.Millisecond, Variance: 4 * time.Millisecond, Taps: Floats2Seconds(bins)[3]},
			{At: (6618 - 37) * time.Millisecond, Mean: (6618 - 37) * time.Millisecond, Variance: 7 * time.Millisecond, Taps: Floats2Seconds(bins)[4]},
			{At: (7153 - 37) * time.Millisecond, Mean: (7153 - 37) * time.Millisecond, Variance: 4 * time.Millisecond, Taps: Floats2Seconds(bins)[5]},
			{At: (7685 - 37) * time.Millisecond, Mean: (7685 - 37) * time.Millisecond, Variance: 5 * time.Millisecond, Taps: Floats2Seconds(bins)[6]},
			{At: (8210 - 37) * time.Millisecond, Mean: (8210 - 37) * time.Millisecond, Variance: 12 * time.Millisecond, Taps: Floats2Seconds(bins)[7]},
		},
	}

	t2b := T2B{
		Precision:  Default.Precision,
		Latency:    37 * time.Millisecond,
		Forgetting: Default.Forgetting,
	}

	beats := t2b.Taps2Beats(Floats2Seconds(taps), 0, Seconds(10.5))

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
		Offset: 323 * time.Millisecond,
		Beats: []Beat{
			{At: 4528 * time.Millisecond, Mean: 4527 * time.Millisecond, Variance: 3 * time.Millisecond, Taps: Floats2Seconds(bins)[0]},
			{At: 5057 * time.Millisecond, Mean: 5057 * time.Millisecond, Variance: 5 * time.Millisecond, Taps: Floats2Seconds(bins)[1]},
			{At: 5582 * time.Millisecond, Mean: 5582 * time.Millisecond, Variance: 4 * time.Millisecond, Taps: Floats2Seconds(bins)[2]},
			{At: 6103 * time.Millisecond, Mean: 6103 * time.Millisecond, Variance: 4 * time.Millisecond, Taps: Floats2Seconds(bins)[3]},
			{At: 6622 * time.Millisecond, Mean: 6622 * time.Millisecond, Variance: 4 * time.Millisecond, Taps: Floats2Seconds(bins)[4]},
			{At: 7153 * time.Millisecond, Mean: 7153 * time.Millisecond, Variance: 4 * time.Millisecond, Taps: Floats2Seconds(bins)[5]},
			{At: 7689 * time.Millisecond, Mean: 7689 * time.Millisecond, Variance: 3 * time.Millisecond, Taps: Floats2Seconds(bins)[6]},
			{At: 8207 * time.Millisecond, Mean: 8207 * time.Millisecond, Variance: 6 * time.Millisecond, Taps: Floats2Seconds(bins)[7]},
		},
	}

	t2b := T2B{
		Precision:  Default.Precision,
		Latency:    Default.Latency,
		Forgetting: 0.1,
	}

	beats := t2b.Taps2Beats(Floats2Seconds(taps), 0, Seconds(10.5))

	if beats.BPM != expected.BPM {
		t.Errorf("Incorrect BPM - expected:%v, got:%v", expected.BPM, beats.BPM)
	}

	if math.Abs(beats.Offset.Seconds()-expected.Offset.Seconds()) > 0.0011 {
		t.Errorf("Incorrect offset - expected:%v, got:%v", expected.Offset, beats.Offset)
	}

	compare(beats.Beats, expected.Beats, t)
}

func TestInterpolateForDegenerateCases(t *testing.T) {
	tests := [][]ckmeans.Cluster{
		{},
		{clusters[0]},
		{clusters[0], clusters[7]},
	}

	expected := [][]int{
		{},
		{1},
		{1, 2},
	}

	for i, v := range tests {
		beats, err := interpolate(v)
		if err != nil {
			t.Fatalf("[%d] unexpected error (%v)", i+1, err)
		}

		if !reflect.DeepEqual(beats, expected[i]) {
			t.Errorf("[%d] invalid interpolation\n   expected: %v\n   got:      %v", i+1, expected[i], beats)
		}
	}
}

func TestInterpolateForThreeBeats(t *testing.T) {
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
		c := []ckmeans.Cluster{}
		for _, ix := range v {
			c = append(c, clusters[ix-1])
		}

		beats, err := interpolate(c)
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
		c := []ckmeans.Cluster{}
		for _, ix := range v {
			c = append(c, clusters[ix-1])
		}

		beats, err := interpolate(c)
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
		c := make([]ckmeans.Cluster, len(v))
		for i, ix := range v {
			c[i] = clusters[ix-1]
		}

		beats, err := interpolate(c)
		if err != nil {
			t.Fatalf("[%v] unexpected error (%v)", v, err)
		}

		expected := v
		for _, x := range exceptions {
			if reflect.DeepEqual(v, x[0]) {
				expected = x[1]
			}
		}

		if !reflect.DeepEqual(beats, expected) {
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
	clusters := []ckmeans.Cluster{
		{Center: 1.0, Variance: 0.001, Values: []float64{}},
		{Center: 1.1, Variance: 0.001, Values: []float64{}},
		{Center: 11.0, Variance: 0.001, Values: []float64{}},
	}

	_, err := interpolate(clusters)
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
