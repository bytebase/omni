package analysis

import (
	"strings"
	"unicode"

	"github.com/bytebase/omni/starrocks/parser"
)

// Range is a byte range in the script, start inclusive and end exclusive.
type Range struct {
	Start int
	End   int
}

// StatementRanges returns the byte range of each statement in a script.
//
// Empty (whitespace-only) segments are skipped, and a range includes the
// trailing ';' when present. A comment-only script — or comment-only text
// after the last statement — still gets a range, because an editor that counts
// statement positions must not lose ground to a buffer of comments.
func StatementRanges(script string) []Range {
	trimmed := strings.TrimRightFunc(script, unicode.IsSpace)

	segments := parser.Split(trimmed)
	if len(segments) == 0 {
		// The splitter drops comment-only buffers entirely. An editor still
		// wants a range covering the non-whitespace span.
		return commentOnlyRanges(trimmed)
	}

	var ranges []Range
	for _, segment := range segments {
		// Skip pure whitespace segments, but preserve comment-only segments.
		if segment.Empty() && !hasNonWhitespace(trimmed[segment.ByteStart:segment.ByteEnd]) {
			continue
		}
		// Skip only leading whitespace — leading comments are kept inside the
		// range, matching the prior ANTLR-helper behaviour, which includes
		// hidden-channel tokens but excludes pure whitespace.
		trimStart := segment.ByteStart
		for trimStart < segment.ByteEnd && isRangeSpace(trimmed[trimStart]) {
			trimStart++
		}
		end := segment.ByteEnd
		if end < len(trimmed) && trimmed[end] == ';' {
			end++
		}
		ranges = append(ranges, Range{Start: trimStart, End: end})
	}

	// Trailing comment-only content: the splitter filters out comment-only
	// segments, so any non-whitespace bytes between the last segment's end and
	// the end of the trimmed input would otherwise be lost. Redundant
	// delimiters open no statement, so a run of them alone adds no range.
	if len(segments) > 0 {
		i := segments[len(segments)-1].ByteEnd
		for i < len(trimmed) && (isRangeSpace(trimmed[i]) || trimmed[i] == ';') {
			i++
		}
		if i < len(trimmed) {
			ranges = append(ranges, Range{Start: i, End: len(trimmed)})
		}
	}

	return ranges
}

// commentOnlyRanges emits a single range covering the non-whitespace span of
// trimmed, or nil when nothing in it looks like a comment. Pure delimiters
// (a single ';') are not turned into ranges; the legacy helper dropped
// single-terminator inputs.
func commentOnlyRanges(trimmed string) []Range {
	if !containsCommentMarker(trimmed) {
		return nil
	}
	start := 0
	for start < len(trimmed) && isRangeSpace(trimmed[start]) {
		start++
	}
	if start >= len(trimmed) {
		return nil
	}
	return []Range{{Start: start, End: len(trimmed)}}
}

// containsCommentMarker reports whether s contains one of the comment
// introducers parser.Split's lexer skips (`--`, `//`, `#`, `/*`). A string
// literal could false-positive, but this runs only when the splitter found no
// segments at all.
func containsCommentMarker(s string) bool {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '#' {
			return true
		}
		if i+1 < len(s) {
			if c == '-' && s[i+1] == '-' {
				return true
			}
			if c == '/' && (s[i+1] == '*' || s[i+1] == '/') {
				return true
			}
		}
	}
	return false
}

// hasNonWhitespace reports whether s contains any non-whitespace character.
// The splitter treats a comment-only segment as empty since it holds no SQL
// token, but an editor still wants its range.
func hasNonWhitespace(s string) bool {
	for i := 0; i < len(s); i++ {
		if !isRangeSpace(s[i]) {
			return true
		}
	}
	return false
}

func isRangeSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}
