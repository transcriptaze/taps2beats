package ckmeans

import (
	"math"
	"reflect"
	"testing"
)

func TestCKMeansWithNilData(t *testing.T) {
	expected := []Cluster{}
	clusters := CKMeans1dDp(nil, nil)

	if !reflect.DeepEqual(clusters, expected) {
		t.Errorf("Expected %v, return, got:%v", expected, clusters)
	}
}

func TestCKMeansWithEmptyData(t *testing.T) {
	expected := []Cluster{}
	clusters := CKMeans1dDp([]float64{}, nil)

	if !reflect.DeepEqual(clusters, expected) {
		t.Errorf("Expected %v, return, got:%v", expected, clusters)
	}
}

func TestCKMeansWithInvalidWeights(t *testing.T) {
	defer func() { recover() }()

	CKMeans1dDp([]float64{-1, 2, -1, 2, 4, 5, 6, -1, 2, -1}, []float64{0.4, 1.3})

	t.Errorf("Expected panic if weights array did not match data array")
}

func TestCKMeans(t *testing.T) {
	x := []float64{-0.9, 1.0, 1.1, 1.9, 2.0, 2.1}
	w := []float64{3, 1, 2, 2, 1, 1}
	expected := []Cluster{
		Cluster{Center: -0.9, Variance: 0.0, Values: []float64{-0.9}},
		Cluster{Center: 1.06666667, Variance: 0.00666667, Values: []float64{1.0, 1.1}},
		Cluster{Center: 1.975, Variance: 0.0275 / 2, Values: []float64{1.9, 2.0, 2.1}},
	}

	clusters := CKMeans1dDp(x, w)

	if !reflect.DeepEqual(clusters, expected) {
		for i := range expected {
			if math.Abs(clusters[i].Center-expected[i].Center) > 0.000001 {
				t.Errorf("(cluster %d) invalid 'center' - expected:%v, got:%v", i, expected[i].Center, clusters[i].Center)
			}

			if math.Abs(clusters[i].Variance-expected[i].Variance) > 0.000001 {
				t.Errorf("(cluster %d) invalid 'variance' - expected:%v, got:%v", i, expected[i].Variance, clusters[i].Variance)
			}

			if len(clusters[i].Values) != len(expected[i].Values) {
				t.Errorf("(cluster %d) invalid 'values' - expected:%v, got:%v", i, expected[i].Values, clusters[i].Values)
			} else {
				for j := range expected[i].Values {
					if math.Abs(clusters[i].Values[j]-expected[i].Values[j]) > 0.000001 {
						t.Errorf("(cluster %d) invalid 'values' - expected:%v, got:%v", i, expected[i].Values, clusters[i].Values)
					}
				}
			}
		}
	}
}

func TestCKMeansTaps(t *testing.T) {
	x := []float64{
		4.570271991, 5.063594027, 5.603539973, 6.102690998, 6.642708943, 7.141796968, 7.710649857, 8.192470916,
		4.506176116, 5.045971061, 5.591722996, 6.114172975, 6.619153989, 7.135788980, 7.693071891, 8.203885893,
		4.529560070, 5.057670039, 5.591721996, 6.137423930, 6.630941966, 7.176683900, 7.698974880, 8.227207848,
		4.529560070, 5.069284016, 5.603428973, 6.102591998, 6.613455000, 7.147644957, 7.699120880, 8.215609871,
		4.517865093, 5.022782107, 5.580101018, 6.096715009, 6.654118921, 7.176371900, 7.681405914, 8.215537871,
		5.133092891, 5.545395086, 6.067721066, 6.578564068, 7.130096991, 7.652464971, 8.134273030,
		4.494581138, 5.040234073, 5.562732052, 6.079333043, 6.624973977, 7.141650968, 7.664070948, 8.198270905,
		4.529408070, 5.040295073, 5.556940064, 6.131584941, 6.654145921, 7.193876866, 7.722112835, 8.244539814,
		4.523631082, 5.046071061, 5.586102007, 6.090995020, 6.596029034, 7.130224991, 7.652501971, 8.180805939,
		4.517979093, 5.046165061, 5.551068075, 6.073547054, 6.607636011, 7.165018923, 7.687334903, 8.238953825,
		4.517911093, 5.069403016, 5.586174007, 6.108568986, 6.578649068, 7.147523957, 7.681606914, 8.262110780,
	}

	expected := []Cluster{
		Cluster{
			Center:   4.52369438,
			Variance: 0.00039172,
			Values:   []float64{4.5702720, 4.5061761, 4.5295601, 4.5295601, 4.5178651, 4.4945811, 4.5294081, 4.5236311, 4.5179791, 4.5179111}},
		Cluster{
			Center:   5.05768749,
			Variance: 0.00082230,
			Values:   []float64{5.0635940, 5.0459710, 5.0576700, 5.0692840, 5.0227821, 5.1330929, 5.0402341, 5.0402951, 5.0460710, 5.0461651, 5.0694030}},
		Cluster{
			Center:   5.57808420,
			Variance: 0.00042773,
			Values:   []float64{5.6035400, 5.5917230, 5.5917220, 5.6034290, 5.5801010, 5.5453951, 5.5627321, 5.5569401, 5.5861020, 5.5510681, 5.5861740}},
		Cluster{
			Center:   6.10048591,
			Variance: 0.00049445,
			Values:   []float64{6.1026910, 6.1141730, 6.1374239, 6.1025920, 6.0967150, 6.0677211, 6.0793330, 6.1315849, 6.0909950, 6.0735471, 6.1085690}},
		Cluster{
			Center:   6.61821608,
			Variance: 0.00071530,
			Values:   []float64{6.6427089, 6.6191540, 6.6309420, 6.6134550, 6.6541189, 6.5785641, 6.6249740, 6.6541459, 6.5960290, 6.6076360, 6.5786491}},
		Cluster{
			Center:   7.15333449,
			Variance: 0.00045737,
			Values:   []float64{7.1417970, 7.1357890, 7.1766839, 7.1476450, 7.1763719, 7.1300970, 7.1416510, 7.1938769, 7.1302250, 7.1650189, 7.1475240}},
		Cluster{
			Center:   7.68575599,
			Variance: 0.00050714,
			Values:   []float64{7.7106499, 7.6930719, 7.6989749, 7.6991209, 7.6814059, 7.6524650, 7.6640709, 7.7221128, 7.6525020, 7.6873349, 7.6816069}},
		Cluster{
			Center:   8.21033333,
			Variance: 0.00121729,
			Values:   []float64{8.1924709, 8.2038859, 8.2272078, 8.2156099, 8.2155379, 8.1342730, 8.1982709, 8.2445398, 8.1808060, 8.2389538, 8.2621108}},
	}

	clusters := CKMeans1dDp(x, nil)

	if !reflect.DeepEqual(clusters, expected) {
		for i := range expected {
			if math.Abs(clusters[i].Center-expected[i].Center) > 0.000001 {
				t.Errorf("(cluster %d) invalid 'center' - expected:%v, got:%v", i, expected[i].Center, clusters[i].Center)
			}

			if math.Abs(clusters[i].Variance-expected[i].Variance) > 0.000001 {
				t.Errorf("(cluster %d) invalid 'variance' - expected:%v, got:%v", i, expected[i].Variance, clusters[i].Variance)
			}

			if len(clusters[i].Values) != len(expected[i].Values) {
				t.Errorf("(cluster %d) invalid 'values' - expected:%v, got:%v", i, expected[i].Values, clusters[i].Values)
			} else {
				for j := range expected[i].Values {
					if math.Abs(clusters[i].Values[j]-expected[i].Values[j]) > 0.000001 {
						t.Errorf("(cluster %d) invalid 'values' - expected:%v, got:%v", i, expected[i].Values, clusters[i].Values)
					}
				}
			}
		}
	}
}
