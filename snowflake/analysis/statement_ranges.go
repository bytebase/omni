package analysis

import (
	"strings"
	"unicode"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/snowflake/parser"
)

// ByteRange is a half-open byte range [Start, End) in a script.
type ByteRange = analysis.ByteRange

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

	var segs []parser.Segment
	for _, seg := range parser.Split(trimmed) {
		if !seg.Empty() {
			segs = append(segs, seg)
		}
	}
	if len(segs) == 0 {
		// No segment holds a statement. An editor still wants a range
		// covering a comment.
		return commentOnlyRanges(trimmed)
	}

	ranges := make([]ByteRange, 0, len(segs))
	for _, seg := range segs {
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

	// The loop above skips comment-only segments, so a comment after the last
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
