package main

import (
	"fmt"
	"testing"
)

func twoSum(nums []int, target int) []int {
	return []int{}
}

func TestTwoSum1(t *testing.T) {
	got := twoSum([]int{2, 7, 11, 15}, 9)
	want := []int{0, 1}
	if !areEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestTwoSum{{ID}}(t *testing.T) {
	got := twoSum({{INPUT1}}, {{INPUT2}})
	want := {{OUTPUT}}
	if !areEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestTwoSum(t *testing.T) {
	var tests = []struct {
		arr    []int
		target int
		want   []int
	}{
		{
			[]int{2, 7, 11, 15},
			9,
			[]int{0, 1},
		},
		{
			[]int{1, 2, 3, 4, 4},
			8,
			[]int{3, 4},
		},
		{
			[]int{-3, 4, 3, 90},
			0,
			[]int{0, 2},
		},
		{
			[]int{-1, -2, -3, -4, -5},
			-8,
			[]int{2, 4},
		},
		{
			[]int{230, 863, 916, 585, 981, 404, 316, 785, 88, 12},
			542,
			[]int{4, 9},
		},
		{
			[]int{10, 20, 30, 40, 50},
			60,
			[]int{1, 4},
		},
		{
			[]int{1, 2},
			3,
			[]int{0, 1},
		},
		{
			[]int{5, 5, 15},
			10,
			[]int{0, 1},
		},
		{
			[]int{1, 3, 3, 3, 4, 4, 4, 5, 6},
			8,
			[]int{3, 7},
		},
		{
			[]int{0, 4, 3, 0},
			0,
			[]int{0, 3},
		},
	}

	for _, test := range tests {
		testName := fmt.Sprintf("arr = %v, target = %d", test.arr, test.target)
		t.Run(testName, func(t *testing.T) {
			got := twoSum(test.arr, test.target)
			if !areEqual(got, test.want) {
				t.Errorf("got %v, want %v", got, test.want)
			}
		})
	}
}

func areEqual(got []int, want []int) bool {
	if len(got) != 2 {
		return false
	}

	if got[0] == want[0] && got[1] == want[1] {
		return true
	}
	if got[0] == want[1] && got[1] == want[0] {
		return true
	}

	return false
}
