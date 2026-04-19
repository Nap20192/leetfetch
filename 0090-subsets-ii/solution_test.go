package subsetsii

import (
	"reflect"
	"testing"
)

func TestSubsetsWithDup(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want interface{}
	}{
		{"example 1",
			[]int{1, 2, 2},
			nil /* TODO: fill expected */},
		{"example 2",
			[]int{0},
			nil /* TODO: fill expected */},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := subsetsWithDup(tt.nums)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("subsetsWithDup(%v) = %v, want %v", tt.nums, got, tt.want)
			}
		})
	}
}
