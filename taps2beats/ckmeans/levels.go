package ckmeans

import (
	"math"
)

func select_levels_weighted(x, y []float64, J [][]int, Kmin, Kmax int, bic []float64) int {
	N := len(x)

	if Kmin > Kmax || N < 2 {
		if Kmax < Kmin {
			return Kmax
		}
		return Kmin
	}

	Kopt := Kmin
	maxBIC := 0.0

	lambda := make([]float64, Kmax)
	mu := make([]float64, Kmax)
	sigma2 := make([]float64, Kmax)
	coeff := make([]float64, Kmax)
	counts := make([]int, Kmax)
	weights := make([]float64, Kmax)

	for K := Kmin; K <= Kmax; K++ {
		backtrackWeighted(x, y, J, counts, weights, K)

		totalweight := 0.0
		for k := 0; k < K; k++ {
			totalweight += weights[k]
		}

		indexLeft := 0

		for k := 0; k < K; k++ { // Estimate GMM parameters first
			lambda[k] = weights[k] / totalweight

			indexRight := indexLeft + counts[k] - 1

			mu[k], sigma2[k] = shiftedDataVarianceWeighted(x, y, weights[k], indexLeft, indexRight)

			if sigma2[k] == 0 || counts[k] == 1 {
				var dmin float64

				if indexLeft > 0 && indexRight < N-1 {
					dmin = x[indexLeft] - x[indexLeft-1]
					if dmin > (x[indexRight+1] - x[indexRight]) {
						dmin = x[indexRight+1] - x[indexRight]
					}
				} else if indexLeft > 0 {
					dmin = x[indexLeft] - x[indexLeft-1]
				} else {
					dmin = x[indexRight+1] - x[indexRight]
				}

				if sigma2[k] == 0 {
					sigma2[k] = dmin * dmin / 4.0 / 9.0
				}

				if counts[k] == 1 {
					sigma2[k] = dmin * dmin
				}
			}

			coeff[k] = lambda[k] / math.Sqrt(2.0*math.Pi*sigma2[k])
			indexLeft = indexRight + 1
		}

		loglikelihood := 0.0

		for i := 0; i < N; i++ {
			L := 0.0
			for k := 0; k < K; k++ {
				L += coeff[k] * math.Exp(-(x[i]-mu[k])*(x[i]-mu[k])/(2.0*sigma2[k]))
			}

			loglikelihood += y[i] * math.Log(L)
		}

		// Compute the Bayesian information criterion
		bicx := 2*loglikelihood - float64(3*K-1)*math.Log(totalweight) //(K*3-1)

		if K == Kmin {
			maxBIC = bicx
			Kopt = Kmin
		} else {
			if bicx > maxBIC {
				maxBIC = bicx
				Kopt = K
			}
		}

		bic[K-Kmin] = bicx
	}

	return Kopt
}

func shiftedDataVarianceWeighted(x, y []float64, total_weight float64, left, right int) (float64, float64) {
	sum := 0.0
	sumsq := 0.0

	mean := 0.0
	variance := 0.0

	n := right - left + 1

	if right >= left {

		median := x[(left+right)/2]

		for i := left; i <= right; i++ {
			sum += (x[i] - median) * y[i]
			sumsq += (x[i] - median) * (x[i] - median) * y[i]
		}
		mean = sum/total_weight + median

		if n > 1 {
			variance = (sumsq - sum*sum/total_weight) / (total_weight - 1)
		}
	}
	return mean, variance
}
