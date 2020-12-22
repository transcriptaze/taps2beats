package regression

import (
	"math"
)

// Deming regression for least squares fitting of 2-d data where both axes are uncertain.
//
// Ref. https://davegiles.blogspot.com/2014/11/orthogonal-regression-first-steps.html
// Ref. https://docs.scipy.org/doc/scipy-0.14.0/reference/odr.html
// Ref. https://en.wikipedia.org/wiki/Deming_regression#Orthogonal_regression
func orthogonal(x, y []float64) (float64, float64) {
	if len(x) < 2 {
		panic("Insufficient data for an orthogonal least squares fit")
	}

	if len(x) == 0 || len(x) != len(y) {
		panic("x and y data should be the same length")
	}

	N := len(x)
	sumx := 0.0
	sumy := 0.0

	for i := range x {
		sumx += x[i]
		sumy += y[i]
	}

	xmean := sumx / float64(N)
	ymean := sumy / float64(N)
	sumxx := 0.0
	sumyy := 0.0
	sumxy := 0.0

	for i := range x {
		sumxx += x[i]*x[i] - 2*x[i]*xmean + xmean*xmean
		sumyy += y[i]*y[i] - 2*y[i]*ymean + ymean*ymean
		sumxy += x[i]*y[i] - x[i]*ymean - y[i]*xmean + xmean*ymean
	}

	Sxx := sumxx / float64(N-1)
	Syy := sumyy / float64(N-1)
	Sxy := sumxy / float64(N-1)

	B1 := (Syy - Sxx + math.Sqrt((Syy-Sxx)*(Syy-Sxx)+4*Sxy*Sxy)) / (2 * Sxy)
	B0 := ymean - B1*xmean

	return B1, B0
}
