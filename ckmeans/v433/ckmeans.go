package ckmeans

import (
	"sort"
)

func CKMeans(data, weights []float64) [][]float64 {
	// validate inputs
	if data == nil || len(data) == 0 {
		return [][]float64{}
	}

	if weights != nil && len(weights) != len(data) {
		panic("Invalid weights")
	}

	// sort and order data
	x := make([]float64, len(data))
	w := make([]float64, len(data))
	order := make([]int, len(data))

	for i := range order {
		order[i] = i
	}

	sort.SliceStable(order, func(i, j int) bool { return data[order[i]] < data[order[j]] })

	if weights == nil {
		for i := range data {
			x[i] = data[order[i]]
			w[i] = 1.0
		}
	} else {
		for i := range data {
			x[i] = data[order[i]]
			w[i] = weights[order[i]]
		}
	}

	// calculate range of K
	// TODO: should this include weights??
	kmin := 1
	kmax := 1

	p := x[0]
	for _, q := range x[1:] {
		if q != p {
			kmax++
			q = p
		}

	}

	k, clusters := ckmeans(x, w, kmin, kmax)
	index := make([]int, len(x))
	for i := range clusters {
		index[order[i]] = clusters[i]
	}

	clustered := make([][]float64, k)
	for i, ix := range index {
		clustered[ix] = append(clustered[ix], data[i])
	}

	return clustered
}

func ckmeans(x, w []float64, kmin, kmax int) (int, []int) {
	N := len(x)
	S := make([][]float64, kmax)
	J := make([][]int, kmax)

	for i := range S {
		S[i] = make([]float64, N)
		J[i] = make([]int, N)
	}

	fill_dp_matrix(x, w, S, J)

	bic := make([]float64, kmax)
	kopt := select_levels_weighted(x, w, J, kmin, kmax, bic)

	if kopt < kmax {
		J = J[0:kopt]
	}

	cluster_sorted := make([]int, N)
	centers := make([]float64, N)
	withinss := make([]float64, N)
	size := make([]float64, kmax)
	backtrackWeightedX(x, w, J, cluster_sorted, centers, withinss, size)

	return kopt, cluster_sorted
}
