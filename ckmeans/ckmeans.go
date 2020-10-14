package ckmeans

import (
	"fmt"
	"sort"
)

func CKMeans(data []float64, k int) ([][]float64, error) {
	if k < 1 {
		return nil, fmt.Errorf("Invalid K: %v - cannot classify into 0 or less clusters", k)
	}

	if k > len(data) {
		return nil, fmt.Errorf("Invalid K: %v - cannot generate more classes than there are data values", k)
	}

	sort.Float64s(data)

	unique := 1
	p := data[0]
	for _, q := range data[1:] {
		if p != q {
			p = q
			unique += 1
		}
	}

	if unique == 1 {
		return [][]float64{data}, nil
	}

	N := len(data)

	S := make([][]float64, k)
	for i := 0; i < k; i++ {
		S[i] = make([]float64, N)
	}

	J := make([][]uint64, k)
	for i := 0; i < k; i++ {
		J[i] = make([]uint64, N)
	}

	initialise(data, S, J, k, N)

	clusters := [][]float64{}
	clusterR := N - 1
	for cluster := k - 1; cluster > -1; cluster-- {
		clusterL := int(J[cluster][clusterR])
		clusters = append(clusters, data[clusterL:clusterR+1])

		if cluster > 0 {
			clusterR = clusterL - 1
		}
	}

	for i, j := 0, len(clusters)-1; i < j; i, j = i+1, j-1 {
		clusters[i], clusters[j] = clusters[j], clusters[i]
	}

	return clusters, nil
}

func initialise(data []float64, S [][]float64, J [][]uint64, K, N int) {
	sumx := make([]float64, N)
	sumxsquared := make([]float64, N)

	// shifting values to centre around median apparently improves numerical stability
	shift := data[N/2]

	for i := 0; i < N; i++ {
		if i == 0 {
			sumx[0] = data[0] - shift
			sumxsquared[0] = (data[0] - shift) * (data[0] - shift)
		} else {
			sumx[i] = sumx[i-1] + data[i] - shift
			sumxsquared[i] = sumxsquared[i-1] + (data[i]-shift)*(data[i]-shift)
		}

		S[0][i] = ssq(0, i, sumx, sumxsquared)
		J[0][i] = 0
	}

	for k := 1; k < K; k++ {
		var imin int
		if k < K-1 {
			if k > 1 {
				imin = k
			} else {
				imin = 1
			}
		} else {
			imin = N - 1
		}

		fill_row(imin, N-1, k, S, J, sumx, sumxsquared, N)
	}
}

func fill_row(imin, imax, k int, S [][]float64, J [][]uint64, sumx, sumxsquared []float64, N int) {
	if imin > imax {
		return
	}

	var i int = (imin + imax) / 2

	S[k][i] = S[k-1][i-1]
	J[k][i] = uint64(i)

	jlow := k
	if imin > k {
		if jlow < int(J[k][imin-1]) {
			jlow = int(J[k][imin-1])
		}
	}

	if jlow < int(J[k-1][i]) {
		jlow = int(J[k-1][i])
	}

	jhigh := i - 1
	if imax < N-1 {
		if jhigh > int(J[k][imax+1]) {
			jhigh = int(J[k][imax+1])
		}
	}

	for j := jhigh; j > jlow-1; j-- {
		sji := ssq(j, i, sumx, sumxsquared)

		if sji+S[k-1][jlow-1] >= S[k][i] {
			break
		}

		// examine the lower bound of the cluster border
		// compute s(jlow, i)
		sjlowi := ssq(jlow, i, sumx, sumxsquared)
		SSQ_jlow := sjlowi + S[k-1][jlow-1]

		if SSQ_jlow < S[k][i] {
			S[k][i] = SSQ_jlow
			J[k][i] = uint64(jlow)
		}

		jlow += 1
		SSQ_j := sji + S[k-1][j-1]
		if SSQ_j < S[k][i] {
			S[k][i] = SSQ_j
			J[k][i] = uint64(j)
		}
	}

	fill_row(imin, i-1, k, S, J, sumx, sumxsquared, N)
	fill_row(i+1, imax, k, S, J, sumx, sumxsquared, N)
}

func ssq(j, i int, sumx, sumxsquared []float64) float64 {
	var sji float64

	if j > 0 {
		muji := (sumx[i] - sumx[j-1]) / float64(i-j+1)
		sji = sumxsquared[i] - sumxsquared[j-1] - float64(i-j+1)*muji*muji
	} else {
		sji = sumxsquared[i] - sumx[i]*sumx[i]/float64(i+1)
	}

	if sji > 0 {
		return sji
	}

	return 0
}
