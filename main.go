package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	var outDir string
	var force bool
	var lang string

	flag.StringVar(&outDir, "o", ".", "output directory")
	flag.BoolVar(&force, "f", false, "overwrite existing directory")
	flag.StringVar(&lang, "lang", "golang", "language slug for code snippet")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: leetfetch [flags] <slug|url>")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Examples:")
		fmt.Fprintln(os.Stderr, "  leetfetch two-sum")
		fmt.Fprintln(os.Stderr, "  leetfetch https://leetcode.com/problems/two-sum/")
		fmt.Fprintln(os.Stderr, "  leetfetch https://leetcode.com/problems/two-sum/description/")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Flags:")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	slug := extractSlug(flag.Arg(0))
	if slug == "" {
		fmt.Fprintln(os.Stderr, "error: could not extract slug from argument")
		os.Exit(1)
	}

	q, err := fetchQuestion(slug)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to fetch: %v\n", err)
		os.Exit(1)
	}
	if q == nil {
		fmt.Fprintf(os.Stderr, "question not found: %s\n", slug)
		os.Exit(1)
	}

	if err := generate(q, outDir, force, lang); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// extractSlug extracts the problem slug from a URL or returns the slug as-is.
func extractSlug(arg string) string {
	arg = strings.TrimRight(arg, "/")
	if idx := strings.Index(arg, "/problems/"); idx >= 0 {
		rest := arg[idx+len("/problems/"):]
		if i := strings.Index(rest, "/"); i >= 0 {
			rest = rest[:i]
		}
		return rest
	}
	return arg
}
