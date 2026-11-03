package review

import (
	"strings"
	"unicode/utf8"
)

// Bytebase renders a finding's message as Markdown, so text taken from the
// SQL must not be read as Markdown: an underscore pair would turn into
// emphasis and a name like shop.app into a link. Code, CodeBlock, and
// Escape are the only ways a message includes such text.

// codeLimit is the longest text, in characters, Code shows before it cuts
// the text with an ellipsis.
const codeLimit = 120

// Code returns s as a Markdown code span. Line breaks and tabs are written
// as \n, \r, and \t, so the span stays on one line, and text longer than
// codeLimit characters is cut with an ellipsis.
func Code(s string) string {
	s = strings.NewReplacer("\n", `\n`, "\r", `\r`, "\t", `\t`).Replace(s)
	if utf8.RuneCountInString(s) > codeLimit {
		s = string([]rune(s)[:codeLimit-1]) + "…"
	}
	fence := strings.Repeat("`", longestRun(s, '`')+1)
	// CommonMark strips one space from each end of a span that starts and
	// ends with one, and a span's text cannot touch its fence's backticks,
	// so pad either case with a space the renderer strips again.
	if s == "" || strings.HasPrefix(s, "`") || strings.HasSuffix(s, "`") ||
		(strings.HasPrefix(s, " ") && strings.HasSuffix(s, " ") && strings.Trim(s, " ") != "") {
		s = " " + s + " "
	}
	return fence + s + fence
}

// CodeBlock returns s as a fenced Markdown code block, s kept verbatim.
func CodeBlock(s string) string {
	fence := strings.Repeat("`", max(3, longestRun(s, '`')+1))
	return fence + "\n" + strings.TrimSuffix(s, "\n") + "\n" + fence
}

// Escape returns s with the characters that start inline Markdown
// backslash-escaped, for prose such as a parser's message, used inside a
// sentence. It does not keep a dotted name from being linked, and it does
// not escape what only means something at the start of a line; SQL text
// belongs in Code.
func Escape(s string) string {
	var b strings.Builder
	for _, r := range s {
		if strings.ContainsRune("\\`*_[]<>~|", r) {
			b.WriteByte('\\')
		}
		b.WriteRune(r)
	}
	return b.String()
}

// longestRun is the length of the longest run of c in s.
func longestRun(s string, c byte) int {
	longest, run := 0, 0
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			run++
			longest = max(longest, run)
		} else {
			run = 0
		}
	}
	return longest
}

// Names writes each name as a code span and joins them as English does:
// "`a`", "`a` and `b`", "`a`, `b`, and `c`".
func Names(names ...string) string {
	spans := make([]string, len(names))
	for i, n := range names {
		spans[i] = Code(n)
	}
	return Join(spans...)
}

// Join joins phrases already written as Markdown the way Names joins
// names.
func Join(items ...string) string {
	switch len(items) {
	case 0:
		return ""
	case 1:
		return items[0]
	case 2:
		return items[0] + " and " + items[1]
	}
	return strings.Join(items[:len(items)-1], ", ") + ", and " + items[len(items)-1]
}

// Plural returns kind as it reads before n names: unchanged when n is 1,
// else with its last word in the plural, as "tables", "indexes",
// "policies", or "operator classes".
func Plural(kind string, n int) string {
	if n == 1 || kind == "" {
		return kind
	}
	head, last := "", kind
	if i := strings.LastIndexByte(kind, ' '); i >= 0 {
		head, last = kind[:i+1], kind[i+1:]
	}
	return head + pluralWord(last)
}

// pluralWord is the plural of one English word as kinds use them.
func pluralWord(w string) string {
	switch w {
	case "statistics", "series":
		return w
	case "index":
		return "indexes"
	}
	switch {
	case strings.HasSuffix(w, "y") && len(w) > 1 && !strings.ContainsRune("aeiou", rune(w[len(w)-2])):
		return w[:len(w)-1] + "ies"
	case strings.HasSuffix(w, "s"), strings.HasSuffix(w, "x"), strings.HasSuffix(w, "ch"), strings.HasSuffix(w, "sh"):
		return w + "es"
	}
	return w + "s"
}
