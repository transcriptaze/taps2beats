package ckmeans

import (
	"reflect"
	"testing"
)

type test struct {
	data     []float64
	k        int
	expected [][]float64
}

func TestCKMeansWithInvalidK(t *testing.T) {
	if _, err := CKMeans([]float64{1, 2, 3, 4, 5, 6}, 0); err == nil {
		t.Errorf("Expected error with K=0, got:%v", err)
	}

	if _, err := CKMeans([]float64{1, 2, 3, 4, 5, 6}, 7); err == nil {
		t.Errorf("Expected error with K > len(data), got:%v", err)
	}
}

func TestCKMeansWithOneUniqueValue(t *testing.T) {
	expected := [][]float64{{3, 3, 3, 3, 3, 3}}

	f, err := CKMeans([]float64{3, 3, 3, 3, 3, 3}, 4)
	if err != nil {
		t.Errorf("Unexpected error (%v)", err)
	}

	if !reflect.DeepEqual(f, expected) {
		t.Errorf("Incorrect classification - expected:%v, got:%v", expected, f)
	}
}

func TestCKMeans(t *testing.T) {
	tests := []test{
		test{[]float64{1}, 1, [][]float64{{1}}},
		test{[]float64{0, 3, 4}, 2, [][]float64{{0}, {3, 4}}},
		test{[]float64{-3, 0, 4}, 2, [][]float64{{-3, 0}, {4}}},
		test{[]float64{1, 1, 1, 1}, 1, [][]float64{{1, 1, 1, 1}}},
		test{[]float64{1, 2, 3}, 3, [][]float64{{1}, {2}, {3}}},
		test{[]float64{1, 2, 2, 3}, 3, [][]float64{{1}, {2, 2}, {3}}},
		test{[]float64{1, 2, 2, 3, 3}, 3, [][]float64{{1}, {2, 2}, {3, 3}}},
		test{[]float64{1, 2, 3, 2, 3}, 3, [][]float64{{1}, {2, 2}, {3, 3}}},
		test{[]float64{3, 2, 3, 2, 1}, 3, [][]float64{{1}, {2, 2}, {3, 3}}},
		test{[]float64{3, 2, 3, 5, 2, 1}, 3, [][]float64{{1, 2, 2}, {3, 3}, {5}}},
		test{[]float64{0, 1, 2, 100, 101, 103}, 2, [][]float64{{0, 1, 2}, {100, 101, 103}}},
		test{[]float64{0, 1, 2, 50, 100, 101, 103}, 3, [][]float64{{0, 1, 2}, {50}, {100, 101, 103}}},
		test{[]float64{-1, 2, -1, 2, 4, 5, 6, -1, 2, -1}, 3, [][]float64{{-1, -1, -1, -1}, {2, 2, 2}, {4, 5, 6}}},
	}

	for _, v := range tests {
		f, err := CKMeans(v.data, v.k)
		if err != nil {
			t.Errorf("Unexpected error (%v)", err)
		}

		if !reflect.DeepEqual(f, v.expected) {
			t.Errorf("Incorrect classification - expected:%v, got:%v", v.expected, f)
		}
	}
}
