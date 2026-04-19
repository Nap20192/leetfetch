package twosum

import (
	"reflect"
	"testing"
)

func TestTwoSum(t *testing.T) {
	tests := []struct {
		name   string
		nums   []int
		target int
		want   []int
	}{
		{"example 1",
			[]int{2, 7, 11, 15},
			9,
			nil /* TODO: fill expected */},
		{"example 2",
			[]int{3, 2, 4},
			6,
			nil /* TODO: fill expected */},
		{"example 3",
			[]int{3, 3},
			6,
			nil /* TODO: fill expected */},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := twoSum(tt.nums, tt.target)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("twoSum(%v, %v) = %v, want %v", tt.nums, tt.target, got, tt.want)
			}
		})
	}
}
