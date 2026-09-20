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

// checkpointBytes bounds how many bytes a Position call decodes: the
// index records the rune count at a rune boundary at least this often.
const checkpointBytes = 512

// Text is an index over one SQL text for converting byte offsets to
// positions. Lines end at '\n'; a '\r' before it counts as a column like
// any other rune.
type Text struct {
	sql        string
	lineStarts []int // byte offset of each line's first byte; lineStarts[0] == 0
	// Rune-count checkpoints: the number of runes starting before byte
	// offset checkpointAt[i] is checkpointRunes[i]. Every offset is a rune
	// boundary, at most checkpointBytes apart, and checkpointAt[0] == 0.
	checkpointAt    []int
	checkpointRunes []int
}

// Index scans sql once and returns a Text whose Position calls each cost
// at most checkpointBytes of decoding, whatever the line length.
func Index(sql string) *Text {
	t := &Text{sql: sql, lineStarts: []int{0}}
	next := 0 // the next offset at or after which a checkpoint is recorded
	for pos, runes := 0, 0; pos < len(sql); runes++ {
		if pos >= next {
			t.checkpointAt = append(t.checkpointAt, pos)
			t.checkpointRunes = append(t.checkpointRunes, runes)
			next = pos + checkpointBytes
		}
		if sql[pos] == '\n' {
			t.lineStarts = append(t.lineStarts, pos+1)
			pos++
			continue
		}
		_, size := utf8.DecodeRuneInString(sql[pos:])
		pos += size
	}
	if len(t.checkpointAt) == 0 {
		t.checkpointAt, t.checkpointRunes = []int{0}, []int{0}
	}
	return t
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
	column = t.runesBefore(offset) - t.runesBefore(t.lineStarts[i]) + 1
	return line, column
}

// runesBefore returns the number of runes that start before the byte
// offset, so an offset inside a multi-byte rune counts that rune. It
// decodes from the nearest checkpoint at or before offset, which is a
// rune boundary, so partial leading bytes are never counted separately.
func (t *Text) runesBefore(offset int) int {
	k := sort.Search(len(t.checkpointAt), func(k int) bool { return t.checkpointAt[k] > offset }) - 1
	runes := t.checkpointRunes[k]
	for pos := t.checkpointAt[k]; pos < offset; runes++ {
		_, size := utf8.DecodeRuneInString(t.sql[pos:])
		pos += size
	}
	return runes
}
