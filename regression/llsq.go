package regression

import (
	"fmt"
)

// Ref. https://stackoverflow.com/questions/16422287/linear-regression-library-for-go-language
func Fit(x, y []float64) (float64, float64, error) {
	if len(x) < 2 {
		return 0, 0, fmt.Errorf("Insufficient data for a least squares fit")
	}

	if len(x) == 0 || len(x) != len(y) {
		return 0, 0, fmt.Errorf("x and y data should be the same length")
	}

	sumX := 0.0
	sumY := 0.0
	sumXX := 0.0
	sumXY := 0.0

	for i := range x {
		sumX += x[i]
		sumY += y[i]
		sumXX += x[i] * x[i]
		sumXY += x[i] * y[i]
	}

	p := float64(len(x))
	m := (p*sumXY - sumX*sumY) / (p*sumXX - sumX*sumX)
	b := (sumY / p) - (m * sumX / p)

	return m, b, nil
}

func Trend(x, y, p []float64) ([]float64, error) {
	m, b, err := Fit(x, y)
	if err != nil {
		return nil, err
	}

	q := make([]float64, len(p))

	for i, v := range p {
		q[i] = m*v + b
	}

	return q, nil
}
