package regression

import (
	"math"
	"testing"
)

func TestOrdinaryLeastSquares(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5, 6, 7, 8}
	y := []float64{4.52369438, 5.05768749, 5.57808420, 6.10048591, 6.61821608, 7.15333449, 7.68575600, 8.21033334}

	m, b := OLS(x, y)

	if math.Abs(m-0.5261) > 0.0001 || math.Abs(b-3.9986) > 0.0001 {
		t.Errorf("Bad fit - expected:(%-.4f,%-.4f), got:(%-.4f,%-.4f)", 0.5261, 3.9986, m, b)
	}
}

func TestOrdinaryLeastSquaresWithNoData(t *testing.T) {
	t.Skip()
}

func TestOrdinaryLeastSquaresWithInsufficientData(t *testing.T) {
	t.Skip()
}

func TestOrdinaryLeastSquaresWithInvalidData(t *testing.T) {
	t.Skip()
}
