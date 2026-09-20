package review

import (
	"strings"
	"unicode/utf8"
)

// Position converts a byte offset in sql into a 1-based line and a 1-based
// column counted in code points, the units of Bytebase's Position. Lines
// end at '\n'; a '\r' before it counts as a column like any other rune.
// An offset outside [0, len(sql)] is clamped, so len(sql) is the position
// just past the last rune. An offset inside a multi-byte rune counts that
// rune as reached.
func Position(sql string, offset int) (line, column int) {
	if offset < 0 {
		offset = 0
	}
	if offset > len(sql) {
		offset = len(sql)
	}
	before := sql[:offset]
	lineStart := strings.LastIndexByte(before, '\n') + 1
	line = strings.Count(before, "\n") + 1
	column = utf8.RuneCountInString(before[lineStart:]) + 1
	return line, column
}
