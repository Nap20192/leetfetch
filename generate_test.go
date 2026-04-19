package main

import (
	"reflect"
	"testing"
)

func TestParseTestCases(t *testing.T) {
	tests := []struct {
		name      string
		examples  string
		numParams int
		want      [][]string
	}{
		{
			name:      "two-sum style: 2 params, 3 examples",
			examples:  "[2,7,11,15]\n9\n[3,2,4]\n6\n[3,3]\n6",
			numParams: 2,
			want: [][]string{
				{"[2,7,11,15]", "9"},
				{"[3,2,4]", "6"},
				{"[3,3]", "6"},
			},
		},
		{
			name:      "single param",
			examples:  "1\n2\n3",
			numParams: 1,
			want:      [][]string{{"1"}, {"2"}, {"3"}},
		},
		{
			name:      "trailing newline",
			examples:  "1\n2\n",
			numParams: 1,
			want:      [][]string{{"1"}, {"2"}},
		},
		{
			name:      "empty input",
			examples:  "",
			numParams: 1,
			want:      nil,
		},
		{
			name:      "zero numParams",
			examples:  "1\n2",
			numParams: 0,
			want:      nil,
		},
		{
			name:      "windows line endings",
			examples:  "[1,2]\r\n3\r\n[4,5]\r\n6",
			numParams: 2,
			want: [][]string{
				{"[1,2]", "3"},
				{"[4,5]", "6"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseTestCases(tt.examples, tt.numParams)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseTestCases()\n  got  %v\n  want %v", got, tt.want)
			}
		})
	}
}

func TestSlugToPackage(t *testing.T) {
	tests := []struct {
		slug string
		want string
	}{
		{"two-sum", "twosum"},
		{"add-two-numbers", "addtwonumbers"},
		{"lru-cache", "lrucache"},
		{"longest-substring-without-repeating-characters", "longestsubstringwithoutrepeatingcharacters"},
		{"3sum", "p3sum"},
		{"", "solution"},
	}
	for _, tt := range tests {
		if got := slugToPackage(tt.slug); got != tt.want {
			t.Errorf("slugToPackage(%q) = %q, want %q", tt.slug, got, tt.want)
		}
	}
}

func TestPaddedID(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"1", "0001"},
		{"42", "0042"},
		{"1337", "1337"},
		{"100", "0100"},
		{"10000", "10000"},
	}
	for _, tt := range tests {
		if got := paddedID(tt.input); got != tt.want {
			t.Errorf("paddedID(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestExtractSlug(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"two-sum", "two-sum"},
		{"https://leetcode.com/problems/two-sum/", "two-sum"},
		{"https://leetcode.com/problems/two-sum/description/", "two-sum"},
		{"https://leetcode.com/problems/two-sum", "two-sum"},
		{"https://leetcode.com/problems/longest-palindromic-substring/", "longest-palindromic-substring"},
	}
	for _, tt := range tests {
		if got := extractSlug(tt.input); got != tt.want {
			t.Errorf("extractSlug(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestCapitalizeFirst(t *testing.T) {
	tests := []struct{ input, want string }{
		{"twoSum", "TwoSum"},
		{"isValid", "IsValid"},
		{"lruCache", "LruCache"},
		{"", ""},
	}
	for _, tt := range tests {
		if got := capitalizeFirst(tt.input); got != tt.want {
			t.Errorf("capitalizeFirst(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
