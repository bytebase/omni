package analysis

import (
	"unicode"
	"unicode/utf8"

	"github.com/bytebase/omni/tidb/parser"
)

// ByteRange is a half-open byte range [Start, End) in a script.
type ByteRange struct {
	Start int
	End   int
}

// StatementRanges returns the range of each top-level statement in a script,
// for editors that run the statement under the cursor. It splits lexically,
// so a statement with a syntax error keeps its range and so do the statements
// around it.
//
// A range starts at the statement's first non-whitespace byte and ends past
// its ';' when present. With a custom DELIMITER the range stops before the
// delimiter. A script with no statements returns nil.
func StatementRanges(statement string) []ByteRange {
	segments := parser.Split(statement)
	if len(segments) == 0 {
		return nil
	}

	ranges := make([]ByteRange, 0, len(segments))
	for _, seg := range segments {
		start := min(seg.ByteStart+leadingWhitespaceBytes(statement, seg.ByteStart, seg.ByteEnd), seg.ByteEnd)
		end := seg.ByteEnd + delimLenAt(statement, seg.ByteEnd)
		ranges = append(ranges, ByteRange{Start: start, End: end})
	}
	return ranges
}

// leadingWhitespaceBytes returns the number of leading Unicode whitespace
// bytes within sql[start:end].
func leadingWhitespaceBytes(sql string, start, end int) int {
	if start < 0 {
		start = 0
	}
	if end > len(sql) {
		end = len(sql)
	}
	i := start
	for i < end {
		r, size := utf8.DecodeRuneInString(sql[i:])
		if !unicode.IsSpace(r) {
			break
		}
		i += size
	}
	return i - start
}

// delimLenAt returns 1 when the default ';' delimiter sits at sql[byteEnd],
// and 0 otherwise, including for a custom DELIMITER.
func delimLenAt(sql string, byteEnd int) int {
	if byteEnd < 0 || byteEnd >= len(sql) {
		return 0
	}
	if sql[byteEnd] == ';' {
		return 1
	}
	return 0
}
