package regression

// Ref. https://stackoverflow.com/questions/16422287/linear-regression-library-for-go-language
func OLS(x, y []float64) (float64, float64) {
	if len(x) < 2 {
		panic("Insufficient data for a least squares fit")
	}

	if len(x) == 0 || len(x) != len(y) {
		panic("x and y data should be the same length")
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

	return m, b
}
