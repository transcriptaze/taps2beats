// Statistical regression functions for fitting beats to an idealized function.
package regression

// Ref. https://stackoverflow.com/questions/16422287/linear-regression-library-for-go-language

// Implementation of ordinary least squares to fit a 2-d set of data to a straight line.
// Returns the gradient and offset of the line.
//
// Panics if the x and y arrays do not contain at least 2 data points or if the
// x and y arrays are different lengths.
func OrdinaryLeastSquares(x, y []float64) (float64, float64) {
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
