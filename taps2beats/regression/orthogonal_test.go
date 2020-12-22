package regression

import (
	"math"
	"testing"
)

// 'R'
//
// library(onls)
//
// x <- c(9.8, 9.7, 10.7, 10.9, 12.4, 12.5, 12.8, 12.8, 12.9, 13.3, 13.4, 13.5, 13.7, 14.9, 15.2, 15.5)
// y <- c(10.1, 11.4, 10.8, 11.3, 11.8, 12.1, 12.3, 13.6, 14.2, 14.4, 14.6, 15.3, 15.5, 15.8, 16.2, 16.5)
// DAT <- data.frame(x, y)
// mod7 <- onls(y ~ a + b * x, data = DAT, start = list(a = 2, b = 3))
//
// print(mod7) ## -1.909/1.208 as on webpage
// plot(mod7)
// ------
// Nonlinear orthogonal regression model
//   model: y ~ a + b * x
//    data: DAT
//      a      b
// -1.909  1.208
//  vertical residual sum-of-squares: 8.91
//  orthogonal residual sum-of-squares: 3.751
//  PASSED: 16 out of 16 fitted points are orthogonal.
//
// Number of iterations to convergence: 4
// Achieved convergence tolerance: 1.49e-08

func TestOrthogonalRegression(t *testing.T) {
	x := []float64{9.8, 9.7, 10.7, 10.9, 12.4, 12.5, 12.8, 12.8, 12.9, 13.3, 13.4, 13.5, 13.7, 14.9, 15.2, 15.5}
	y := []float64{10.1, 11.4, 10.8, 11.3, 11.8, 12.1, 12.3, 13.6, 14.2, 14.4, 14.6, 15.3, 15.5, 15.8, 16.2, 16.5}

	m, c := orthogonal(x, y)

	if math.Abs(m-1.2080) > 0.0001 || math.Abs(c - -1.9088) > 0.0001 {
		t.Errorf("Incorrect orthogonal regression - expected: %.4f,%.4f, got:%.4f,%.4f", 1.2080, -1.9088, m, c)
	}
}
