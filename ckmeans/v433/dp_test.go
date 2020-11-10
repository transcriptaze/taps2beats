package ckmeans

import (
	"math"
	"reflect"
	"testing"
)

func TestFillDPMatrix(t *testing.T) {
	x := []float64{-0.9, 1.0, 1.1, 1.9, 2.0, 2.1}
	w := []float64{3.0, 1.0, 2.0, 2.0, 1.0, 1.0}

	N := len(x)
	S := make([][]float64, 6)
	J := make([][]int, 6)

	for i := range S {
		S[i] = make([]float64, N)
		J[i] = make([]int, N)
	}

	fill_dp_matrix(x, w, S, J)

	s := [][]float64{
		{0.0, 2.7075, 5.80833, 10.75875, 12.6570, 14.42400},
		{0.0, 0.0000, 0.00667, 0.840000, 1.14000, 1.448570},
		{0.0, 0.0000, 0.00000, 0.006670, 0.01333, 0.034167},
		{0.0, 0.0000, 0.00000, 0.000000, 0.00667, 0.011667},
		{0.0, 0.0000, 0.00000, 0.000000, 0.00000, 0.005000},
		{0.0, 0.0000, 0.00000, 0.000000, 0.00000, 0.000000},
	}

	j := [][]int{
		{0, 0, 0, 0, 0, 0},
		{0, 1, 1, 1, 1, 1},
		{0, 0, 2, 3, 3, 3},
		{0, 0, 0, 3, 4, 4},
		{0, 0, 0, 0, 4, 4},
		{0, 0, 0, 0, 0, 5},
	}

loop:
	for i := range S {
		for j := range S[i] {
			if delta := math.Abs(S[i][j] - s[i][j]); delta > 0.01 {
				t.Errorf("Incorrect S matrix - expected: %v, got: %v", s, S)
				break loop
			}
		}
	}

	if !reflect.DeepEqual(J, j) {
		t.Errorf("Incorrect J matrix - expected: %v, got: %v", j, J)
	}
}
