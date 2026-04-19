package sudokusolver

import (
	"reflect"
	"testing"
)

func TestSolveSudoku(t *testing.T) {
	tests := []struct {
		name  string
		board [][]byte
		want  interface{}
	}{
		{"example 1",
			[][]byte{[]byte{'5', '3', '.', '.', '7', '.', '.', '.', '.'}, []byte{'6', '.', '.', '1', '9', '5', '.', '.', '.'}, []byte{'.', '9', '8', '.', '.', '.', '.', '6', '.'}, []byte{'8', '.', '.', '.', '6', '.', '.', '.', '3'}, []byte{'4', '.', '.', '8', '.', '3', '.', '.', '1'}, []byte{'7', '.', '.', '.', '2', '.', '.', '.', '6'}, []byte{'.', '6', '.', '.', '.', '.', '2', '8', '.'}, []byte{'.', '.', '.', '4', '1', '9', '.', '.', '5'}, []byte{'.', '.', '.', '.', '8', '.', '.', '7', '9'}},
			nil /* TODO: fill expected */},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := solveSudoku(tt.board)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("solveSudoku(%v) = %v, want %v", tt.board, got, tt.want)
			}
		})
	}
}
