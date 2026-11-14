package analysis

import (
	"strings"
	"unicode"

	"github.com/bytebase/omni/snowflake/parser"
)

// ByteRange is a half-open byte range [Start, End) in a script.
type ByteRange struct {
	Start int
	End   int
}

// StatementRanges returns the range of each statement in a script, for
// editors that count statement positions. Trailing whitespace of the script
// is ignored.
//
// A range starts at the statement's first non-whitespace byte, so a leading
// comment stays inside it, and ends past its ';' when present. Comments after
// the last statement get a range of their own, and so does a script holding
// only comments. A script holding only ';' gets none.
func StatementRanges(statement string) []ByteRange {
	trimmed := strings.TrimRightFunc(statement, unicode.IsSpace)

	segs := parser.Split(trimmed)
	if len(segs) == 0 {
		return commentOnlyRanges(trimmed)
	}

	ranges := make([]ByteRange, 0, len(segs))
	for _, seg := range segs {
		if !hasNonWhitespace(trimmed[seg.ByteStart:seg.ByteEnd]) {
			continue
		}
		start := seg.ByteStart
		for start < seg.ByteEnd {
			c := trimmed[start]
			if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
				start++
				continue
			}
			break
		}
		end := seg.ByteEnd
		if end < len(trimmed) && trimmed[end] == ';' {
			end++
		}
		ranges = append(ranges, ByteRange{Start: start, End: end})
	}

	// parser.Split drops comment-only segments, so a comment after the last
	// statement would otherwise have no range.
	last := segs[len(segs)-1]
	tailStart := last.ByteEnd
	if tailStart < len(trimmed) && trimmed[tailStart] == ';' {
		tailStart++
	}
	for tailStart < len(trimmed) {
		c := trimmed[tailStart]
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			tailStart++
			continue
		}
		break
	}
	if tailStart < len(trimmed) {
		ranges = append(ranges, ByteRange{Start: tailStart, End: len(trimmed)})
	}

	return ranges
}

// commentOnlyRanges returns one range over the non-whitespace span of a
// script with no statements, or nil when it holds no comment.
func commentOnlyRanges(trimmed string) []ByteRange {
	if !containsCommentMarker(trimmed) {
		return nil
	}
	start := 0
	for start < len(trimmed) {
		c := trimmed[start]
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			start++
			continue
		}
		break
	}
	if start >= len(trimmed) {
		return nil
	}
	return []ByteRange{{Start: start, End: len(trimmed)}}
}

// containsCommentMarker reports whether s contains `--`, `//` or `/*`. It is
// only called on scripts with no statements, so no string literal can hold one.
func containsCommentMarker(s string) bool {
	for i := 0; i < len(s); i++ {
		c := s[i]
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

func hasNonWhitespace(s string) bool {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c != ' ' && c != '\t' && c != '\n' && c != '\r' {
			return true
		}
	}
	return false
}
