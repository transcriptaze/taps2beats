package ckmeans

func fill_dp_matrix(x, w []float64, S [][]float64, J [][]int) {
	K := len(S)
	N := len(S[0])

	sum_x := make([]float64, N)
	sum_x_sq := make([]float64, N)
	sum_w := make([]float64, len(w))
	sum_w_sq := make([]float64, len(w))

	// 	//jseq := []int{}

	shift := x[N/2] // median. used to shift the values of x to improve numerical stability

	sum_x[0] = w[0] * (x[0] - shift)
	sum_x_sq[0] = w[0] * (x[0] - shift) * (x[0] - shift)
	sum_w[0] = w[0]
	sum_w_sq[0] = w[0] * w[0]

	S[0][0] = 0
	J[0][0] = 0

	for i := 1; i < N; i++ {
		sum_x[i] = sum_x[i-1] + w[i]*(x[i]-shift)
		sum_x_sq[i] = sum_x_sq[i-1] + w[i]*(x[i]-shift)*(x[i]-shift)
		sum_w[i] = sum_w[i-1] + w[i]
		sum_w_sq[i] = sum_w_sq[i-1] + w[i]*w[i]

		// NOTE: using same dissimilarity as SMAWK - original algorithm potentially (but not really) allowed for alternative criterion here
		//       i.e. not convinced embedding criterion in SMAWK is all that correct
		S[0][i] = dissimilarity(0, i, sum_x, sum_x_sq, sum_w, sum_w_sq)
		J[0][i] = 0
	}

	for q := 1; q < K; q++ {
		var imin int
		if q < K-1 {
			imin = 1
			if q > imin {
				imin = q
			}
		} else {
			imin = N - 1
		}

		fill_row_q_SMAWK(imin, N-1, q, S, J, sum_x, sum_x_sq, sum_w, sum_w_sq)
	}
}

func backtrackWeighted(x, y []float64, J [][]int, counts []int, weights []float64, K int) {
	N := len(J[0])
	cluster_right := N - 1

	for k := K - 1; k >= 0; k-- {
		cluster_left := J[k][cluster_right]
		counts[k] = cluster_right - cluster_left + 1

		weights[k] = 0
		for i := cluster_left; i <= cluster_right; i++ {
			weights[k] += y[i]
		}

		if k > 0 {
			cluster_right = cluster_left - 1
		}
	}
}

func backtrackWeightedX(x, y []float64, J [][]int) []int {
	K := len(J)
	N := len(J[0])
	clusters := make([]int, N)

	cluster_right := N - 1

	for k := K - 1; k >= 0; k-- {
		cluster_left := J[k][cluster_right]

		for i := cluster_left; i <= cluster_right; i++ {
			clusters[i] = k
		}

		if k > 0 {
			cluster_right = cluster_left - 1
		}
	}

	return clusters
}
