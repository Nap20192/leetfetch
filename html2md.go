package main

import (
	"fmt"
	"html"
	"regexp"
	"strings"
)

var (
	rePre    = regexp.MustCompile(`(?si)<pre[^>]*>(.*?)</pre>`)
	reStrong = regexp.MustCompile(`(?si)<(?:strong|b)[^>]*>(.*?)</(?:strong|b)>`)
	reEm     = regexp.MustCompile(`(?si)<(?:em|i)[^>]*>(.*?)</(?:em|i)>`)
	reCode   = regexp.MustCompile(`(?s)<code[^>]*>(.*?)</code>`)
	reLI     = regexp.MustCompile(`(?si)<li[^>]*>(.*?)</li>`)
	reSup    = regexp.MustCompile(`(?si)<sup[^>]*>(.*?)</sup>`)
	reSub    = regexp.MustCompile(`(?si)<sub[^>]*>(.*?)</sub>`)
	reBR     = regexp.MustCompile(`<br\s*/?>`)
	reAnyTag = regexp.MustCompile(`<[^>]+>`)
	reSpaces = regexp.MustCompile(`[^\S\n]{2,}`)
	reBlanks = regexp.MustCompile(`\n{3,}`)
)

func stripTags(s string) string {
	return reAnyTag.ReplaceAllString(s, "")
}

// HTMLToMarkdown converts LeetCode HTML problem content to Markdown.
func HTMLToMarkdown(src string) string {
	s := strings.ReplaceAll(src, "\r\n", "\n")

	type preEntry struct{ placeholder, block string }
	var preBlocks []preEntry

	// Protect <pre> blocks before any other transformation.
	s = rePre.ReplaceAllStringFunc(s, func(m string) string {
		inner := rePre.FindStringSubmatch(m)[1]
		inner = stripTags(inner)
		inner = html.UnescapeString(inner)
		inner = strings.TrimSpace(inner)
		ph := fmt.Sprintf("\x00LFPRE%d\x00", len(preBlocks))
		preBlocks = append(preBlocks, preEntry{ph, "```\n" + inner + "\n```"})
		return ph
	})

	// Inline: <strong>/<b>
	s = reStrong.ReplaceAllStringFunc(s, func(m string) string {
		return "**" + stripTags(reStrong.FindStringSubmatch(m)[1]) + "**"
	})
	// Inline: <em>/<i>
	s = reEm.ReplaceAllStringFunc(s, func(m string) string {
		return "*" + stripTags(reEm.FindStringSubmatch(m)[1]) + "*"
	})
	// Inline: <code>
	// Do NOT call html.UnescapeString here. Entities like &lt; must stay encoded
	// until after the final stripTags pass; otherwise bare < would be mistaken
	// for an HTML tag start and stripped.
	s = reCode.ReplaceAllStringFunc(s, func(m string) string {
		inner := stripTags(reCode.FindStringSubmatch(m)[1])
		return "`" + inner + "`"
	})
	// Inline: <sup>/<sub>
	s = reSup.ReplaceAllStringFunc(s, func(m string) string {
		return "^" + stripTags(reSup.FindStringSubmatch(m)[1])
	})
	s = reSub.ReplaceAllStringFunc(s, func(m string) string {
		return "_" + stripTags(reSub.FindStringSubmatch(m)[1])
	})

	// List items — entities still encoded, no early unescape.
	s = reLI.ReplaceAllStringFunc(s, func(m string) string {
		inner := strings.TrimSpace(stripTags(reLI.FindStringSubmatch(m)[1]))
		return "\n- " + inner
	})

	// <br> → newline.
	s = reBR.ReplaceAllString(s, "\n")

	// Strip remaining block tags. At this point bare < / > do NOT appear in
	// processed content (they are still &lt; / &gt;), so the tag regex is safe.
	s = stripTags(s)

	// Decode all HTML entities in one final pass.
	s = html.UnescapeString(s)
	s = strings.ReplaceAll(s, "\u00a0", " ")
	s = strings.ReplaceAll(s, "\t", " ")

	// Restore pre blocks.
	for _, e := range preBlocks {
		s = strings.ReplaceAll(s, e.placeholder, "\n\n"+e.block+"\n\n")
	}

	s = reSpaces.ReplaceAllString(s, " ")
	s = reBlanks.ReplaceAllString(s, "\n\n")
	s = strings.TrimSpace(s)
	return s
}
