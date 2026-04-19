package main

import (
	"encoding/json"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

// MetaData is the parsed metaData JSON string from LeetCode.
type MetaData struct {
	Name   string      `json:"name"`
	Params []ParamMeta `json:"params"`
	Return ReturnMeta  `json:"return"`
}

// ParamMeta describes one function parameter.
type ParamMeta struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// ReturnMeta describes the function return type.
type ReturnMeta struct {
	Type string `json:"type"`
}

// slugToPackage converts a LeetCode slug to a valid Go package name.
// "two-sum" → "twosum", "3sum" → "threesum" prefix is not applied — digits are kept.
func slugToPackage(slug string) string {
	var sb strings.Builder
	for _, r := range slug {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			sb.WriteRune(unicode.ToLower(r))
		}
	}
	s := sb.String()
	if s == "" {
		return "solution"
	}
	// Package name must not start with a digit.
	if unicode.IsDigit(rune(s[0])) {
		s = "p" + s
	}
	return s
}

// capitalizeFirst uppercases the first rune of s.
func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// paddedID formats a question frontend ID as a 4-digit zero-padded string.
func paddedID(id string) string {
	var n int
	fmt.Sscanf(id, "%d", &n)
	return fmt.Sprintf("%04d", n)
}

// dirName returns the output directory name for a question.
func dirName(q *Question) string {
	return paddedID(q.QuestionFrontendID) + "-" + q.TitleSlug
}

func snippetForLang(q *Question, langSlug string) string {
	for _, s := range q.CodeSnippets {
		if s.LangSlug == langSlug {
			return s.Code
		}
	}
	return ""
}

func parseMetaData(s string) (*MetaData, error) {
	var md MetaData
	if err := json.Unmarshal([]byte(s), &md); err != nil {
		return nil, err
	}
	return &md, nil
}

// parseTestCases groups exampleTestcases lines into chunks of numParams each.
func parseTestCases(examples string, numParams int) [][]string {
	if numParams <= 0 {
		return nil
	}
	examples = strings.ReplaceAll(examples, "\r\n", "\n")
	examples = strings.ReplaceAll(examples, "\r", "\n")
	trimmed := strings.TrimSpace(examples)
	if trimmed == "" {
		return nil
	}
	lines := strings.Split(trimmed, "\n")
	var cases [][]string
	for i := 0; i+numParams <= len(lines); i += numParams {
		group := make([]string, numParams)
		copy(group, lines[i:i+numParams])
		cases = append(cases, group)
	}
	return cases
}

// generate creates the output directory with README.md, solution.go, solution_test.go.
func generate(q *Question, outDir string, force bool, langSlug string) error {
	if q.IsPaidOnly {
		return fmt.Errorf("this is a premium question, not accessible without subscription")
	}

	name := dirName(q)
	dir := filepath.Join(outDir, name)

	if _, err := os.Stat(dir); err == nil {
		if !force {
			return fmt.Errorf("directory %s already exists (use -f to overwrite)", name)
		}
		if err := os.RemoveAll(dir); err != nil {
			return fmt.Errorf("remove existing directory: %w", err)
		}
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	pkgName := slugToPackage(q.TitleSlug)

	readme := genREADME(q)
	solution := genSolution(q, pkgName, langSlug)
	tests := genTests(q, pkgName)

	gomod := []byte("module " + pkgName + "\n\ngo 1.21\n")

	type file struct {
		name    string
		content []byte
	}
	files := []file{
		{"go.mod", gomod},
		{"README.md", []byte(readme)},
		{"solution.go", solution},
		{"solution_test.go", tests},
	}

	for _, f := range files {
		if err := os.WriteFile(filepath.Join(dir, f.name), f.content, 0o644); err != nil {
			return fmt.Errorf("write %s: %w", f.name, err)
		}
	}

	fmt.Printf("Created %s/\n", name)
	for _, f := range []string{"README.md", "solution.go", "solution_test.go"} {
		fmt.Printf("  %s\n", f)
	}
	return nil
}

func genREADME(q *Question) string {
	md := HTMLToMarkdown(q.Content)
	var sb strings.Builder
	fmt.Fprintf(&sb, "# %s. %s\n\n", q.QuestionFrontendID, q.Title)
	fmt.Fprintf(&sb, "Difficulty: %s\n", q.Difficulty)
	fmt.Fprintf(&sb, "Link: https://leetcode.com/problems/%s/\n\n", q.TitleSlug)
	sb.WriteString(md)
	sb.WriteString("\n")
	return sb.String()
}

func genSolution(q *Question, pkgName, langSlug string) []byte {
	snippet := snippetForLang(q, langSlug)
	if snippet == "" {
		snippet = fmt.Sprintf("// TODO: no %s snippet available\nfunc solution() {}\n", langSlug)
	}

	var needsListNode, needsTreeNode bool
	if q.MetaData != "" {
		if md, err := parseMetaData(q.MetaData); err == nil && md != nil {
			for _, p := range md.Params {
				if strings.Contains(p.Type, "ListNode") {
					needsListNode = true
				}
				if strings.Contains(p.Type, "TreeNode") {
					needsTreeNode = true
				}
			}
			if strings.Contains(md.Return.Type, "ListNode") {
				needsListNode = true
			}
			if strings.Contains(md.Return.Type, "TreeNode") {
				needsTreeNode = true
			}
		}
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "package %s\n\n", pkgName)
	fmt.Fprintf(&sb, "// https://leetcode.com/problems/%s/\n", q.TitleSlug)
	fmt.Fprintf(&sb, "// Difficulty: %s\n\n", q.Difficulty)

	if needsListNode {
		sb.WriteString("// TODO: define ListNode (or import it if shared across problems)\n")
		sb.WriteString("// type ListNode struct { Val int; Next *ListNode }\n\n")
	}
	if needsTreeNode {
		sb.WriteString("// TODO: define TreeNode (or import it if shared across problems)\n")
		sb.WriteString("// type TreeNode struct { Val int; Left *TreeNode; Right *TreeNode }\n\n")
	}

	sb.WriteString(snippet)
	sb.WriteString("\n")

	formatted, err := format.Source([]byte(sb.String()))
	if err != nil {
		return []byte(sb.String())
	}
	return formatted
}

func genTests(q *Question, pkgName string) []byte {
	if q.MetaData == "" {
		return fallbackTest(pkgName, q.TitleSlug)
	}
	md, err := parseMetaData(q.MetaData)
	if err != nil || md == nil || md.Name == "" {
		return fallbackTest(pkgName, q.TitleSlug)
	}

	cases := parseTestCases(q.ExampleTestcases, len(md.Params))
	retGoType := LCTypeToGo(md.Return.Type)

	var sb strings.Builder
	fmt.Fprintf(&sb, "package %s\n\n", pkgName)
	sb.WriteString("import (\n\t\"reflect\"\n\t\"testing\"\n)\n\n")

	testName := "Test" + capitalizeFirst(md.Name)
	fmt.Fprintf(&sb, "func %s(t *testing.T) {\n", testName)

	// Struct definition.
	sb.WriteString("\ttests := []struct {\n")
	sb.WriteString("\t\tname string\n")
	for _, p := range md.Params {
		fmt.Fprintf(&sb, "\t\t%s %s\n", p.Name, LCTypeToGo(p.Type))
	}
	fmt.Fprintf(&sb, "\t\twant %s\n", retGoType)
	sb.WriteString("\t}{\n")

	// Test cases.
	for i, tc := range cases {
		fmt.Fprintf(&sb, "\t\t{\"example %d\",\n", i+1)
		for j, p := range md.Params {
			lit, _ := ParseLiteral(tc[j], p.Type)
			fmt.Fprintf(&sb, "\t\t\t%s,\n", lit)
		}
		fmt.Fprintf(&sb, "\t\t\t%s},\n", wantZero(md.Return.Type))
	}

	sb.WriteString("\t}\n\n")

	// Test loop.
	sb.WriteString("\tfor _, tt := range tests {\n")
	sb.WriteString("\t\tt.Run(tt.name, func(t *testing.T) {\n")

	argNames := make([]string, len(md.Params))
	for i, p := range md.Params {
		argNames[i] = "tt." + p.Name
	}
	fmt.Fprintf(&sb, "\t\t\tgot := %s(%s)\n", md.Name, strings.Join(argNames, ", "))
	sb.WriteString("\t\t\tif !reflect.DeepEqual(got, tt.want) {\n")

	fmtVerbs := make([]string, len(md.Params))
	for i := range md.Params {
		fmtVerbs[i] = "%v"
	}
	fmtStr := md.Name + "(" + strings.Join(fmtVerbs, ", ") + ") = %v, want %v"
	fmt.Fprintf(&sb, "\t\t\t\tt.Errorf(%q, %s, got, tt.want)\n",
		fmtStr, strings.Join(argNames, ", "))

	sb.WriteString("\t\t\t}\n")
	sb.WriteString("\t\t})\n")
	sb.WriteString("\t}\n")
	sb.WriteString("}\n")

	formatted, err := format.Source([]byte(sb.String()))
	if err != nil {
		// Emit compilable stub so the user can see what went wrong.
		stub := fmt.Sprintf("package %s\n\nimport \"testing\"\n\n// TODO: auto-generation failed, implement manually\n// Error: %v\nfunc %s(t *testing.T) { t.Skip(\"not implemented\") }\n",
			pkgName, err, testName)
		if src, e2 := format.Source([]byte(stub)); e2 == nil {
			return src
		}
		return []byte(stub)
	}
	return formatted
}

// wantZero returns the zero-value literal with a TODO comment for the want field.
func wantZero(lcType string) string {
	switch LCTypeToGo(lcType) {
	case "int", "int64", "float64", "byte":
		return "0 /* TODO: fill expected */"
	case "string":
		return `"" /* TODO: fill expected */`
	case "bool":
		return "false /* TODO: fill expected */"
	default:
		return "nil /* TODO: fill expected */"
	}
}

func fallbackTest(pkgName, slug string) []byte {
	src := fmt.Sprintf("package %s\n\nimport \"testing\"\n\n// TODO: implement tests for %s\nfunc TestSolution(t *testing.T) {\n\tt.Skip(\"not implemented\")\n}\n",
		pkgName, slug)
	if formatted, err := format.Source([]byte(src)); err == nil {
		return formatted
	}
	return []byte(src)
}
