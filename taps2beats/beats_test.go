package taps2beats

import (
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/twystd/taps2beats/taps2beats/ckmeans"
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
	{At: Seconds(4.523694381), Mean: Seconds(4.523694381), Variance: Seconds(0.000391722), Taps: seconds(bins[0]...)},
	{At: Seconds(5.057687493), Mean: Seconds(5.057687493), Variance: Seconds(0.000822308), Taps: seconds(bins[1]...)},
	{At: Seconds(5.578084204), Mean: Seconds(5.578084204), Variance: Seconds(0.000427737), Taps: seconds(bins[2]...)},
	{At: Seconds(6.100485910), Mean: Seconds(6.100485910), Variance: Seconds(0.000494451), Taps: seconds(bins[3]...)},
	{At: Seconds(6.618216081), Mean: Seconds(6.618216081), Variance: Seconds(0.000715306), Taps: seconds(bins[4]...)},
	{At: Seconds(7.153334490), Mean: Seconds(7.153334490), Variance: Seconds(0.000457375), Taps: seconds(bins[5]...)},
	{At: Seconds(7.685755996), Mean: Seconds(7.685755996), Variance: Seconds(0.000507140), Taps: seconds(bins[6]...)},
	{At: Seconds(8.210333335), Mean: Seconds(8.210333335), Variance: Seconds(0.001217297), Taps: seconds(bins[7]...)},
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
	{At: Seconds(4.525), Mean: Seconds(4.523694381), Variance: Seconds(0.000391722), Taps: seconds(bins[0]...)},
	{At: Seconds(5.051), Mean: Seconds(5.057687493), Variance: Seconds(0.000822308), Taps: seconds(bins[1]...)},
	{At: Seconds(5.577), Mean: Seconds(5.578084204), Variance: Seconds(0.000427737), Taps: seconds(bins[2]...)},
	{At: Seconds(6.103), Mean: Seconds(6.100485910), Variance: Seconds(0.000494451), Taps: seconds(bins[3]...)},
	{At: Seconds(6.629), Mean: Seconds(6.618216081), Variance: Seconds(0.000715306), Taps: seconds(bins[4]...)},
	{At: Seconds(7.155), Mean: Seconds(7.153334490), Variance: Seconds(0.000457375), Taps: seconds(bins[5]...)},
	{At: Seconds(7.681), Mean: Seconds(7.685755996), Variance: Seconds(0.000507140), Taps: seconds(bins[6]...)},
	{At: Seconds(8.207), Mean: Seconds(8.210333335), Variance: Seconds(0.001217297), Taps: seconds(bins[7]...)},
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

	beats := Taps2Beats(Floats2Seconds(taps), 0.0)

	if beats.BPM != expected.BPM {
		t.Errorf("Incorrect BPM - expected:%v, got:%v", expected.BPM, beats.BPM)
	}

	if math.Abs(beats.Offset.Seconds()-expected.Offset.Seconds()) > 0.0011 {
		t.Errorf("Incorrect offset - expected:%v, got:%v", expected.Offset, beats.Offset)
	}

	compare(beats.Beats, expected.Beats, t)
}

// func TestTaps2BeatsJS(t *testing.T) {
// 	expected := Beats{
// 		BPM:    114,
// 		Offset: 316 * time.Millisecond,
// 		Beats:  []Beat{beats[8], beats[9], beats[10], beats[11], beats[12], beats[13], beats[14], beats[15]},
// 	}
//
// 	taps := [][]float64{
// 		{
// 			1.6808709809265137, 2.680558198364258, 4.384705977111817, 5.07266901335144, 5.744292851226807,
// 			7.21743113923645, 7.864426975204468, 8.521181040054321, 12.617829851226807, 13.296893062942505,
// 			13.985408954223633, 15.345035160217286, 16.017740091552735, 16.705076910354613, 18.137180977111818,
// 			18.775760009536743, 19.46510804577637},
// 	}
//
// 	fmt.Printf("%v\n", Floats2Seconds(taps))
// 	beats := Taps2Beats(Floats2Seconds(taps), 0.0)
//
// 	if beats.BPM != expected.BPM {
// 		t.Errorf("Incorrect BPM - expected:%v, got:%v", expected.BPM, beats.BPM)
// 	}
//
// 	if math.Abs(beats.Offset.Seconds()-expected.Offset.Seconds()) > 0.0011 {
// 		t.Errorf("Incorrect offset - expected:%v, got:%v", expected.Offset, beats.Offset)
// 	}
//
// 	compare(beats.Beats, expected.Beats, t)
// }

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

	beats := Taps2Beats(Floats2Seconds(taps), 0.0)

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
			{At: Seconds(1.000), Mean: Seconds(1), Variance: Seconds(0.025), Taps: seconds(1, 1.1, 1.2, 0.9, 0.8)},
			{At: Seconds(2.000), Mean: Seconds(2), Variance: Seconds(0.025), Taps: seconds(2, 2.1, 2.2, 1.9, 1.8)},
			{At: Seconds(50.000), Mean: Seconds(50), Variance: Seconds(0.025), Taps: seconds(50, 50.1, 50.2, 49.9, 49.8)},
		},
	}

	taps := [][]float64{
		{1.0, 1.1, 1.2, 0.9, 0.8},
		{2.0, 2.1, 2.2, 1.9, 1.8},
		{50.0, 50.1, 50.2, 49.9, 49.8},
	}

	beats := Taps2Beats(Floats2Seconds(taps), 0.0)

	if beats.BPM != expected.BPM {
		t.Errorf("Incorrect BPM - expected:%v, got:%v", expected.BPM, beats.BPM)
	}

	if math.Abs(beats.Offset.Seconds()-expected.Offset.Seconds()) > 0.0011 {
		t.Errorf("Incorrect offset - expected:%v, got:%v", expected.Offset, beats.Offset)
	}

	compare(beats.Beats, expected.Beats, t)
}

func TestTaps2BeatsWithForgetting(t *testing.T) {
	expected := Beats{
		BPM:    114,
		Offset: 309 * time.Millisecond,
		Beats: []Beat{
			{At: 4521 * time.Millisecond, Mean: 4521 * time.Millisecond, Variance: 170 * time.Microsecond, Taps: seconds(bins[0]...)},
			{At: 5057 * time.Millisecond, Mean: 5057 * time.Millisecond, Variance: 492 * time.Microsecond, Taps: seconds(bins[1]...)},
			{At: 5575 * time.Millisecond, Mean: 5575 * time.Millisecond, Variance: 251 * time.Microsecond, Taps: seconds(bins[2]...)},
			{At: 6099 * time.Millisecond, Mean: 6099 * time.Millisecond, Variance: 307 * time.Microsecond, Taps: seconds(bins[3]...)},
			{At: 6614 * time.Millisecond, Mean: 6614 * time.Millisecond, Variance: 482 * time.Microsecond, Taps: seconds(bins[4]...)},
			{At: 7153 * time.Millisecond, Mean: 7153 * time.Millisecond, Variance: 289 * time.Microsecond, Taps: seconds(bins[5]...)},
			{At: 7683 * time.Millisecond, Mean: 7683 * time.Millisecond, Variance: 321 * time.Microsecond, Taps: seconds(bins[6]...)},
			{At: 8215 * time.Millisecond, Mean: 8215 * time.Millisecond, Variance: 864 * time.Microsecond, Taps: seconds(bins[7]...)},
		},
	}

	beats := Taps2Beats(Floats2Seconds(taps), 0.1)

	if beats.BPM != expected.BPM {
		t.Errorf("Incorrect BPM - expected:%v, got:%v", expected.BPM, beats.BPM)
	}

	if math.Abs(beats.Offset.Seconds()-expected.Offset.Seconds()) > 0.0011 {
		t.Errorf("Incorrect offset - expected:%v, got:%v", expected.Offset, beats.Offset)
	}

	compare(beats.Beats, expected.Beats, t)
}

func TestTaps2BeatsWithUnorderedTaps(t *testing.T) {
	taps := [][]float64{
		{4.570271991, 5.603539973, 7.141796968, 5.063594027, 6.102690998, 6.642708943, 7.710649857, 8.192470916},
		{4.506176116, 5.591722996, 7.135788980, 6.114172975, 5.045971061, 6.619153989, 7.693071891, 8.203885893},
		{4.529560070, 5.591721996, 7.176683900, 6.137423930, 6.630941966, 5.057670039, 7.698974880, 8.227207848},
		{4.529560070, 5.603428973, 5.069284016, 6.102591998, 6.613455000, 7.147644957, 7.699120880, 8.215609871},
		{4.517865093, 5.580101018, 7.681405914, 6.096715009, 6.654118921, 7.176371900, 5.022782107, 8.215537871},
		{5.133092891, 6.067721066, 8.134273030, 6.578564068, 7.130096991, 7.652464971, 5.545395086},
		{4.494581138, 5.562732052, 7.664070948, 6.079333043, 6.624973977, 7.141650968, 5.040234073, 8.198270905},
		{4.529408070, 5.556940064, 6.131584941, 7.722112835, 6.654145921, 7.193876866, 5.040295073, 8.244539814},
		{4.523631082, 5.586102007, 6.090995020, 7.652501971, 6.596029034, 7.130224991, 5.046071061, 8.180805939},
		{4.517979093, 5.551068075, 7.687334903, 6.073547054, 6.607636011, 7.165018923, 5.046165061, 8.238953825},
		{4.517911093, 5.586174007, 7.681606914, 6.108568986, 6.578649068, 7.147523957, 5.069403016, 8.262110780},
	}
	expected := Beats{
		BPM:    114,
		Offset: 316 * time.Millisecond,
		Beats:  []Beat{beats[8], beats[9], beats[10], beats[11], beats[12], beats[13], beats[14], beats[15]},
	}

	beats := Taps2Beats(Floats2Seconds(taps), 0.0)

	if beats.BPM != expected.BPM {
		t.Errorf("Incorrect BPM - expected:%v, got:%v", expected.BPM, beats.BPM)
	}

	if math.Abs(beats.Offset.Seconds()-expected.Offset.Seconds()) > 0.0011 {
		t.Errorf("Incorrect offset - expected:%v, got:%v", expected.Offset, beats.Offset)
	}

	compare(beats.Beats, expected.Beats, t)
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
