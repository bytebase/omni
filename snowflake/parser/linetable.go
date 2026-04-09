package parser

// LineTable provides O(log n) byte-offset → (line, column) conversion
// for a single source string. Construct once per input, then call
// Position repeatedly.
//
// Handles LF, CRLF, and bare-CR line endings. Each of the following is
// treated as exactly ONE line break:
//   - "\n" (Unix / LF)
//   - "\r\n" (Windows / CRLF)
//   - "\r" (classic Mac / bare CR)
//
// Lines and columns are 1-based and measured in bytes (not runes).
// Position is a byte count, so a single multi-byte UTF-8 character
// occupies multiple column positions.
type LineTable struct {
	// lineStarts[i] is the byte offset of the start of line (i+1).
	// lineStarts[0] is always 0 (line 1 starts at the beginning of input).
	lineStarts []int
}

// NewLineTable builds a LineTable for the given input. O(n) time and
// O(lines) space.
func NewLineTable(input string) *LineTable {
	// Line 1 always starts at byte 0.
	starts := []int{0}
	for i := 0; i < len(input); i++ {
		switch input[i] {
		case '\n':
			// LF line break. Next line starts at i+1.
			starts = append(starts, i+1)
		case '\r':
			if i+1 < len(input) && input[i+1] == '\n' {
				// CRLF. Skip the \n and start the next line at i+2.
				starts = append(starts, i+2)
				i++ // consume the \n so we don't double-count it
			} else {
				// Bare CR. Next line starts at i+1.
				starts = append(starts, i+1)
			}
		}
	}
	return &LineTable{lineStarts: starts}
}

// Position returns the 1-based (line, column) for the given byte offset.
// If byteOffset is negative, returns (1, 1). If byteOffset exceeds the
// length of the input, the line is the last line of input and the column
// continues past the end of that line (e.g. for input "abc", Position(100)
// returns (1, 101)) — the column count is monotonic and unbounded beyond
// EOF. Callers that want an EOF-clamped position should min() the offset
// against len(input) before calling.
//
// Column is measured in bytes: a byte offset at the start of a line
// returns column 1; the byte after returns column 2; and so on.
func (lt *LineTable) Position(byteOffset int) (line, col int) {
	if byteOffset < 0 {
		return 1, 1
	}
	if len(lt.lineStarts) == 0 {
		return 1, 1
	}
	// Binary search for the largest index i such that lineStarts[i] <= byteOffset.
	lo, hi := 0, len(lt.lineStarts)-1
	for lo < hi {
		mid := (lo + hi + 1) / 2
		if lt.lineStarts[mid] <= byteOffset {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	line = lo + 1                            // 1-based line number
	col = byteOffset - lt.lineStarts[lo] + 1 // 1-based column
	return line, col
}
