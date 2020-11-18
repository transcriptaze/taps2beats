package taps

import (
	"math"
	"reflect"
	"testing"

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

func TestTaps2Beats(t *testing.T) {
	v := floats2seconds(bins)
	expected := []Beat{
		{at: seconds(4.524686644), mean: seconds(4.523694381), variance: seconds(0.003525498), taps: v[0]},
		{at: seconds(5.050761599), mean: seconds(5.057687493), variance: seconds(0.008223081), taps: v[1]},
		{at: seconds(5.576836554), mean: seconds(5.578084204), variance: seconds(0.004277370), taps: v[2]},
		{at: seconds(6.102911509), mean: seconds(6.100485910), variance: seconds(0.004944514), taps: v[3]},
		{at: seconds(6.628986464), mean: seconds(6.618216081), variance: seconds(0.007153066), taps: v[4]},
		{at: seconds(7.155061419), mean: seconds(7.153334490), variance: seconds(0.004573754), taps: v[5]},
		{at: seconds(7.681136374), mean: seconds(7.685755996), variance: seconds(0.005071400), taps: v[6]},
		{at: seconds(8.207211329), mean: seconds(8.210333335), variance: seconds(0.012172972), taps: v[7]},
	}

	beats, err := taps2beats(floats2seconds(taps), seconds(0.0), seconds(8.5))
	if err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	if len(beats) != len(expected) {
		t.Errorf("Invalid result\n   expected: %v\n   got:      %v", expected, beats)
	} else {
		for i, v := range expected {
			if !reflect.DeepEqual(v, beats[i]) {
				if v.at != beats[i].at {
					t.Errorf("Invalid beat %d 'at' - expected:%v, got:%v", i+1, v.at, beats[i].at)
				}

				if v.mean != beats[i].mean {
					t.Errorf("Invalid beat %d 'mean' - expected:%v, got:%v", i+1, v.mean, beats[i].mean)
				}

				if v.variance != beats[i].variance {
					t.Errorf("Invalid beat %d 'variance' - expected:%v, got:%v", i+1, v.variance, beats[i].variance)
				}

				if !reflect.DeepEqual(v.taps, beats[i].taps) {
					for j := range v.taps {
						if math.Abs(beats[i].taps[j].Seconds()-v.taps[j].Seconds()) > 0.0001 {
							t.Errorf("Invalid beat %d\n   expected: %v\n   got:      %v", i+1, v, beats[i])
							break
						}
					}
				}
			}
		}
	}
}

func TestExtrapolate(t *testing.T) {
	clusters := []ckmeans.Cluster{
		{Center: 4.524686644},
		{Center: 5.050761599},
		{Center: 5.576836554},
		{Center: 6.102911509},
		{Center: 6.628986464},
		{Center: 7.155061419},
		{Center: 7.681136374},
		{Center: 8.207211329},
	}

	expected := []Beat{
		{at: seconds(4.524686644)},
		{at: seconds(5.050761599)},
		{at: seconds(5.576836554)},
		{at: seconds(6.102911509)},
		{at: seconds(6.628986464)},
		{at: seconds(7.155061419)},
		{at: seconds(7.681136374)},
		{at: seconds(8.207211329)},
	}

	beats, err := extrapolate(clusters, seconds(4.5), seconds(8.5))
	if err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	if len(beats) != len(expected) {
		t.Errorf("Invalid result\n   expected: %v\n   got:      %v", expected, beats)
	} else {
		for i, v := range expected {
			if !reflect.DeepEqual(v, beats[i]) {
				if math.Abs(beats[i].at.Seconds()-v.at.Seconds()) > 0.0001 {
					t.Errorf("Invalid beat %d 'at' - expected:%v, got:%v", i+1, v.at.Seconds(), beats[i].at.Seconds())
				}

				//				if v.mean != beats[i].mean {
				//					t.Errorf("Invalid beat %d 'mean' - expected:%v, got:%v", i+1, v.mean, beats[i].mean)
				//				}
				//
				//				if v.variance != beats[i].variance {
				//					t.Errorf("Invalid beat %d 'variance' - expected:%v, got:%v", i+1, v.variance, beats[i].variance)
				//				}
				//
				//				if !reflect.DeepEqual(v.taps, beats[i].taps) {
				//					for j := range v.taps {
				//						if math.Abs(beats[i].taps[j].Seconds()-v.taps[j].Seconds()) > 0.0001 {
				//							t.Errorf("Invalid beat %d\n   expected: %v\n   got:      %v", i+1, v, beats[i])
				//							break
				//						}
				//					}
				//				}
			}
		}
	}
}
