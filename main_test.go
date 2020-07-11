package main

import (
	"testing"
)

func TestParallelSum(t *testing.T) {
	var slice1 = []int{0, 1, 2}
	var slice2 = []int{2, -1, 7}
	var slice3 = []int{1, 2, 3}
	var slice4 = []int{2, -1, -7}
	result := ParallelSum(slice1, slice2, slice3, slice4)
	if result[0] != 5 || result[1] != 1 || result[2] != 5 || len(result) != 3 {
		t.Fail()
	}
}
