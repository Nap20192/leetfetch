package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// LCTypeToGo converts a LeetCode type string to its Go equivalent.
func LCTypeToGo(lcType string) string {
	switch lcType {
	case "integer":
		return "int"
	case "integer[]":
		return "[]int"
	case "integer[][]":
		return "[][]int"
	case "long":
		return "int64"
	case "long[]":
		return "[]int64"
	case "long[][]":
		return "[][]int64"
	case "string":
		return "string"
	case "string[]":
		return "[]string"
	case "string[][]":
		return "[][]string"
	case "character":
		return "byte"
	case "character[]":
		return "[]byte"
	case "character[][]":
		return "[][]byte"
	case "boolean":
		return "bool"
	case "boolean[]":
		return "[]bool"
	case "double":
		return "float64"
	case "double[]":
		return "[]float64"
	case "double[][]":
		return "[][]float64"
	case "ListNode":
		return "*ListNode"
	case "TreeNode":
		return "*TreeNode"
	default:
		return "interface{}"
	}
}

// ZeroLiteral returns the Go zero-value literal for a LeetCode type.
func ZeroLiteral(lcType string) string {
	switch LCTypeToGo(lcType) {
	case "int", "int64", "float64", "byte":
		return "0"
	case "string":
		return `""`
	case "bool":
		return "false"
	default:
		return "nil"
	}
}

// ParseLiteral converts a raw LeetCode example-testcase value to a Go source literal.
// Returns (literal, true) on success. On failure returns (zero_with_comment, false).
func ParseLiteral(raw, lcType string) (string, bool) {
	raw = strings.TrimSpace(raw)
	switch lcType {
	case "integer":
		var v int64
		if err := json.Unmarshal([]byte(raw), &v); err != nil {
			return fmt.Sprintf("0 /* TODO: could not parse %q */", raw), false
		}
		return fmt.Sprintf("%d", v), true

	case "long":
		var v int64
		if err := json.Unmarshal([]byte(raw), &v); err != nil {
			return fmt.Sprintf("0 /* TODO: could not parse %q */", raw), false
		}
		return fmt.Sprintf("%d", v), true

	case "double":
		var v float64
		if err := json.Unmarshal([]byte(raw), &v); err != nil {
			return fmt.Sprintf("0.0 /* TODO: could not parse %q */", raw), false
		}
		return fmt.Sprintf("%g", v), true

	case "boolean":
		var v bool
		if err := json.Unmarshal([]byte(raw), &v); err != nil {
			return fmt.Sprintf("false /* TODO: could not parse %q */", raw), false
		}
		if v {
			return "true", true
		}
		return "false", true

	case "string":
		var v string
		if err := json.Unmarshal([]byte(raw), &v); err != nil {
			return fmt.Sprintf("%q", raw), true
		}
		return fmt.Sprintf("%q", v), true

	case "character":
		var v string
		if err := json.Unmarshal([]byte(raw), &v); err != nil || len(v) == 0 {
			return fmt.Sprintf("'?' /* TODO: could not parse %q */", raw), false
		}
		return fmt.Sprintf("'%c'", rune(v[0])), true

	case "integer[]":
		var v []int64
		if err := json.Unmarshal([]byte(raw), &v); err != nil {
			return fmt.Sprintf("nil /* TODO: could not parse %q */", raw), false
		}
		return fmtIntSlice("[]int", v), true

	case "integer[][]":
		var v [][]int64
		if err := json.Unmarshal([]byte(raw), &v); err != nil {
			return fmt.Sprintf("nil /* TODO: could not parse %q */", raw), false
		}
		return fmtIntMatrix("[][]int", v), true

	case "long[]":
		var v []int64
		if err := json.Unmarshal([]byte(raw), &v); err != nil {
			return fmt.Sprintf("nil /* TODO: could not parse %q */", raw), false
		}
		return fmtIntSlice("[]int64", v), true

	case "long[][]":
		var v [][]int64
		if err := json.Unmarshal([]byte(raw), &v); err != nil {
			return fmt.Sprintf("nil /* TODO: could not parse %q */", raw), false
		}
		return fmtIntMatrix("[][]int64", v), true

	case "double[]":
		var v []float64
		if err := json.Unmarshal([]byte(raw), &v); err != nil {
			return fmt.Sprintf("nil /* TODO: could not parse %q */", raw), false
		}
		return fmtFloat64Slice(v), true

	case "string[]":
		var v []string
		if err := json.Unmarshal([]byte(raw), &v); err != nil {
			return fmt.Sprintf("nil /* TODO: could not parse %q */", raw), false
		}
		return fmtStringSlice(v), true

	case "string[][]":
		var v [][]string
		if err := json.Unmarshal([]byte(raw), &v); err != nil {
			return fmt.Sprintf("nil /* TODO: could not parse %q */", raw), false
		}
		return fmtStringMatrix(v), true

	case "boolean[]":
		var v []bool
		if err := json.Unmarshal([]byte(raw), &v); err != nil {
			return fmt.Sprintf("nil /* TODO: could not parse %q */", raw), false
		}
		return fmtBoolSlice(v), true

	case "character[]":
		// LeetCode may send "abc" or ["a","b","c"].
		var str string
		if err := json.Unmarshal([]byte(raw), &str); err == nil {
			return fmtByteSliceFromString(str), true
		}
		var arr []string
		if err := json.Unmarshal([]byte(raw), &arr); err == nil {
			return fmtByteSliceFromArray(arr), true
		}
		return fmt.Sprintf("nil /* TODO: could not parse %q */", raw), false

	case "character[][]":
		var v [][]string
		if err := json.Unmarshal([]byte(raw), &v); err != nil {
			return fmt.Sprintf("nil /* TODO: could not parse %q */", raw), false
		}
		rows := make([]string, len(v))
		for i, row := range v {
			rows[i] = fmtByteSliceFromArray(row)
		}
		if len(rows) == 0 {
			return "[][]byte{}", true
		}
		return "[][]byte{" + strings.Join(rows, ", ") + "}", true

	default:
		return fmt.Sprintf("nil /* TODO: unsupported type %q, fill manually */", lcType), false
	}
}

func fmtIntSlice(goType string, v []int64) string {
	if len(v) == 0 {
		return goType + "{}"
	}
	parts := make([]string, len(v))
	for i, x := range v {
		parts[i] = fmt.Sprintf("%d", x)
	}
	return goType + "{" + strings.Join(parts, ", ") + "}"
}

func fmtIntMatrix(goType string, v [][]int64) string {
	if len(v) == 0 {
		return goType + "{}"
	}
	rows := make([]string, len(v))
	for i, row := range v {
		parts := make([]string, len(row))
		for j, x := range row {
			parts[j] = fmt.Sprintf("%d", x)
		}
		rows[i] = "{" + strings.Join(parts, ", ") + "}"
	}
	return goType + "{" + strings.Join(rows, ", ") + "}"
}

func fmtFloat64Slice(v []float64) string {
	if len(v) == 0 {
		return "[]float64{}"
	}
	parts := make([]string, len(v))
	for i, x := range v {
		parts[i] = fmt.Sprintf("%g", x)
	}
	return "[]float64{" + strings.Join(parts, ", ") + "}"
}

func fmtStringSlice(v []string) string {
	if len(v) == 0 {
		return "[]string{}"
	}
	parts := make([]string, len(v))
	for i, x := range v {
		parts[i] = fmt.Sprintf("%q", x)
	}
	return "[]string{" + strings.Join(parts, ", ") + "}"
}

func fmtStringMatrix(v [][]string) string {
	if len(v) == 0 {
		return "[][]string{}"
	}
	rows := make([]string, len(v))
	for i, row := range v {
		parts := make([]string, len(row))
		for j, x := range row {
			parts[j] = fmt.Sprintf("%q", x)
		}
		rows[i] = "{" + strings.Join(parts, ", ") + "}"
	}
	return "[][]string{" + strings.Join(rows, ", ") + "}"
}

func fmtBoolSlice(v []bool) string {
	if len(v) == 0 {
		return "[]bool{}"
	}
	parts := make([]string, len(v))
	for i, x := range v {
		if x {
			parts[i] = "true"
		} else {
			parts[i] = "false"
		}
	}
	return "[]bool{" + strings.Join(parts, ", ") + "}"
}

func fmtByteSliceFromString(s string) string {
	if len(s) == 0 {
		return "[]byte{}"
	}
	parts := make([]string, len(s))
	for i, c := range s {
		parts[i] = fmt.Sprintf("'%c'", c)
	}
	return "[]byte{" + strings.Join(parts, ", ") + "}"
}

func fmtByteSliceFromArray(arr []string) string {
	if len(arr) == 0 {
		return "[]byte{}"
	}
	parts := make([]string, len(arr))
	for i, s := range arr {
		if len(s) > 0 {
			parts[i] = fmt.Sprintf("'%c'", rune(s[0]))
		} else {
			parts[i] = "0"
		}
	}
	return "[]byte{" + strings.Join(parts, ", ") + "}"
}
