package regression

import (
	"math"
	"testing"
)

func TestOrdinaryLeastSquaresFit(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5, 6, 7, 8}
	y := []float64{4.52369438, 5.05768749, 5.57808420, 6.10048591, 6.61821608, 7.15333449, 7.68575600, 8.21033334}

	m, b, err := Fit(x, y)
	if err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	if math.Abs(m-0.5261) > 0.0001 || math.Abs(b-3.9986) > 0.0001 {
		t.Errorf("Bad fit - expected:(%-.4f,%-.4f), got:(%-.4f,%-.4f)", 0.5261, 3.9986, m, b)
	}
}

func TestOrdinaryLeastSquaresTrend(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5, 6, 7, 8}
	y := []float64{4.52369438, 5.05768749, 5.57808420, 6.10048591, 6.61821608, 7.15333449, 7.68575600, 8.21033334}
	p := []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	expected := []float64{3.9986, 4.5247, 5.0508, 5.5768, 6.1029, 6.6290, 7.1551, 7.6811, 8.2072, 8.7333}

	q, err := Trend(x, y, p)
	if err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	if len(q) != len(expected) {
		t.Errorf("Bad fit\n   expected:%v\n   got:     %v", expected, q)
	} else {
		for i := range expected {
			if math.Abs(q[i]-expected[i]) > 0.0001 {
				t.Errorf("Bad fit\n   expected:%v\n   got:     %v", expected, q)
				break
			}
		}
	}
}
