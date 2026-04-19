package main

import (
	"testing"
)

func TestLCTypeToGo(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"integer", "int"},
		{"integer[]", "[]int"},
		{"integer[][]", "[][]int"},
		{"long", "int64"},
		{"long[]", "[]int64"},
		{"string", "string"},
		{"string[]", "[]string"},
		{"string[][]", "[][]string"},
		{"boolean", "bool"},
		{"boolean[]", "[]bool"},
		{"character", "byte"},
		{"character[]", "[]byte"},
		{"double", "float64"},
		{"double[]", "[]float64"},
		{"ListNode", "*ListNode"},
		{"TreeNode", "*TreeNode"},
		{"unknownType", "interface{}"},
	}
	for _, tt := range tests {
		if got := LCTypeToGo(tt.input); got != tt.want {
			t.Errorf("LCTypeToGo(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestZeroLiteral(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"integer", "0"},
		{"long", "0"},
		{"double", "0"},
		{"string", `""`},
		{"boolean", "false"},
		{"integer[]", "nil"},
		{"string[]", "nil"},
		{"ListNode", "nil"},
	}
	for _, tt := range tests {
		if got := ZeroLiteral(tt.input); got != tt.want {
			t.Errorf("ZeroLiteral(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestParseLiteral(t *testing.T) {
	tests := []struct {
		raw    string
		lcType string
		want   string
		ok     bool
	}{
		// Primitives
		{"9", "integer", "9", true},
		{"-3", "integer", "-3", true},
		{"true", "boolean", "true", true},
		{"false", "boolean", "false", true},
		{"3.14", "double", "3.14", true},
		{`"hello"`, "string", `"hello"`, true},
		{`"a"`, "character", `'a'`, true},

		// Integer slices
		{"[2,7,11,15]", "integer[]", "[]int{2, 7, 11, 15}", true},
		{"[]", "integer[]", "[]int{}", true},
		{"[0]", "integer[]", "[]int{0}", true},

		// Integer matrices
		{"[[1,2],[3,4]]", "integer[][]", "[][]int{{1, 2}, {3, 4}}", true},
		{"[]", "integer[][]", "[][]int{}", true},

		// String slices
		{`["a","b","c"]`, "string[]", `[]string{"a", "b", "c"}`, true},
		{"[]", "string[]", "[]string{}", true},

		// Boolean slices
		{"[true,false,true]", "boolean[]", "[]bool{true, false, true}", true},

		// character[] from string
		{`"abc"`, "character[]", "[]byte{'a', 'b', 'c'}", true},

		// long
		{"1000000000000", "long", "1000000000000", true},
	}
	for _, tt := range tests {
		got, ok := ParseLiteral(tt.raw, tt.lcType)
		if got != tt.want || ok != tt.ok {
			t.Errorf("ParseLiteral(%q, %q)\n  got  (%q, %v)\n  want (%q, %v)",
				tt.raw, tt.lcType, got, ok, tt.want, tt.ok)
		}
	}
}
