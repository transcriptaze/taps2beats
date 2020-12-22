// Basic and uninformed port of the R Ckmeans.1d.dp_v4.3.3 implementation for the optimal
// clustering of 1-dimensional data.
//
// This implemenation is simplified specifically for use in this applications and only
// implements L2 dissimilarity and linear clustering. It also assumes the inputs are
// always weighted.
package ckmeans

import (
	"sort"
)

// A cluster resulting from the clustering algorithm.
type Cluster struct {
	Center   float64   // mean of the values assigned to this cluster from the data set
	Variance float64   // variance of the values assigned to this cluster from the data set
	Values   []float64 // the values from the data set that are assigned to this cluster
}

// ckmeans.1d.dp implementation using L2 dissimilarity and linear clustering. Panics
// if the weights array does not match the supplied data array.
func CKMeans1dDp(data, weights []float64) []Cluster {
	// validate inputs
	if data == nil || len(data) == 0 {
		return []Cluster{}
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

	k, clusters, centers, variance := ckmeans(x, w, kmin, kmax)
	index := make([]int, len(x))
	for i := range clusters {
		index[order[i]] = clusters[i]
	}

	clustered := make([]Cluster, k)

	for i := 0; i < k; i++ {
		clustered[i].Center = centers[i]
		clustered[i].Variance = variance[i]
	}

	for i, ix := range index {
		clustered[ix].Values = append(clustered[ix].Values, data[i])
	}

	return clustered
}

func ckmeans(x, w []float64, kmin, kmax int) (int, []int, []float64, []float64) {
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

	clusters := backtrackWeightedX(x, w, J)

	// ... calulate mean and variance
	centers := make([]float64, kopt)
	variance := make([]float64, kopt)
	count := make([]int, kopt)
	withinss := make([]float64, kopt)
	sum := make([]float64, kopt)
	sumw := make([]float64, kopt)

	for i := range x {
		ix := clusters[i]
		sum[ix] += x[i] * w[i]
		sumw[ix] += w[i]
	}

	for i := 0; i < kopt; i++ {
		centers[i] = sum[i] / sumw[i]
	}

	for i := range x {
		ix := clusters[i]
		withinss[ix] += w[i] * (x[i]*x[i] - 2*x[i]*centers[ix] + centers[ix]*centers[ix])
		count[ix] += 1
	}

	for i := 0; i < kopt; i++ {
		if count[i] > 1 {
			variance[i] = withinss[i] / float64(count[i]-1)
		} else {
			variance[i] = 0
		}
	}

	return kopt, clusters, centers, variance
}
