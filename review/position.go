package review

import (
	"sort"
	"unicode/utf8"
)

// Position converts a byte offset in sql into a 1-based line and a 1-based
// column counted in code points, the units of Bytebase's Position. It is
// Index(sql).Position(offset); a caller converting many offsets of one
// text builds the Index once.
func Position(sql string, offset int) (line, column int) {
	return Index(sql).Position(offset)
}

// Text is an index over one SQL text for converting byte offsets to
// positions. Lines end at '\n'; a '\r' before it counts as a column like
// any other rune.
type Text struct {
	sql        string
	lineStarts []int // byte offset of each line's first byte; lineStarts[0] == 0
}

// Index scans sql once and returns a Text whose Position calls each cost
// the length of one line.
func Index(sql string) *Text {
	starts := []int{0}
	for i := 0; i < len(sql); i++ {
		if sql[i] == '\n' {
			starts = append(starts, i+1)
		}
	}
	return &Text{sql: sql, lineStarts: starts}
}

// Position converts a byte offset into a 1-based line and a 1-based column
// counted in code points. An offset outside [0, len(sql)] is clamped, so
// len(sql) is the position just past the last rune. An offset inside a
// multi-byte rune counts that rune as reached.
func (t *Text) Position(offset int) (line, column int) {
	if offset < 0 {
		offset = 0
	}
	if offset > len(t.sql) {
		offset = len(t.sql)
	}
	// The line is the last start at or before offset.
	i := sort.Search(len(t.lineStarts), func(i int) bool { return t.lineStarts[i] > offset }) - 1
	line = i + 1
	// Decode rune by rune from the line start rather than counting runes
	// in a slice, so an offset inside a multi-byte rune counts that rune
	// once instead of counting each of its leading bytes as a rune.
	column = 1
	for pos := t.lineStarts[i]; pos < offset; column++ {
		_, size := utf8.DecodeRuneInString(t.sql[pos:])
		pos += size
	}
	return line, column
}
