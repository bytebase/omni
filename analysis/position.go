package analysis

import (
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

// PositionAt returns the 1-based line and code-point column of byte offset in
// script. "\r", "\n" and "\r\n" each end a line. An offset outside the script
// is clamped, and one inside a code point stops before it.
func PositionAt(script string, offset int) Position {
	offset = max(0, min(offset, len(script)))
	line, column := 1, 0
	for i := 0; i < offset; {
		r, size := utf8.DecodeRuneInString(script[i:])
		if i+size > offset {
			break
		}
		switch r {
		case '\r':
			line++
			column = 0
		case '\n':
			if i == 0 || script[i-1] != '\r' {
				line++
				column = 0
			}
		default:
			column++
		}
		i += size
	}
	return Position{Line: line, Column: column + 1}
}

// Positions maps byte offsets in one script to positions, by the rules of
// PositionAt, without rescanning the script for each offset.
type Positions struct {
	script string
	// lines holds the byte offset each line starts at.
	lines []int
	// The last offset At counted to, so increasing offsets on one line are
	// counted once.
	lastOffset, lastLine, lastColumn int
}

// NewPositions indexes script's line starts.
func NewPositions(script string) *Positions {
	lines := []int{0}
	for i := 0; i < len(script); i++ {
		switch script[i] {
		case '\r':
			lines = append(lines, i+1)
		case '\n':
			if i == 0 || script[i-1] != '\r' {
				lines = append(lines, i+1)
			}
		default:
		}
	}
	return &Positions{script: script, lines: lines}
}

// At returns the position of offset, as PositionAt(script, offset) does.
func (p *Positions) At(offset int) Position {
	offset = max(0, min(offset, len(p.script)))
	line := sort.Search(len(p.lines), func(i int) bool { return p.lines[i] > offset }) - 1
	start, column := p.lines[line], 0
	if line == p.lastLine && offset >= p.lastOffset && p.lastOffset > start {
		start, column = p.lastOffset, p.lastColumn
	}
	i := start
	for i < offset {
		r, size := utf8.DecodeRuneInString(p.script[i:])
		if i+size > offset {
			break
		}
		if r != '\r' && r != '\n' {
			column++
		}
		i += size
	}
	p.lastOffset, p.lastLine, p.lastColumn = i, line, column
	return Position{Line: line + 1, Column: column + 1}
}

// Shift moves p, a position in a segment, to the script that holds the segment
// starting at base: lines move by base's, and a column on the segment's first
// line also moves by base's.
func (p Position) Shift(base Position) Position {
	if p.Line == 1 {
		p.Column += base.Column - 1
	}
	p.Line += base.Line - 1
	return p
}

// LSPRanges converts byte ranges in script to LSP ranges. A range outside the
// script, or one ending before it starts, is dropped.
func LSPRanges(script string, ranges []ByteRange) []LSPRange {
	if ranges == nil {
		return nil
	}
	cursor := &lspCursor{script: script}
	out := make([]LSPRange, 0, len(ranges))
	for _, r := range ranges {
		if r.Start < 0 || r.End > len(script) || r.Start > r.End {
			continue
		}
		out = append(out, LSPRange{Start: cursor.position(r.Start), End: cursor.position(r.End)})
	}
	return out
}

// lspCursor walks forward through a script, restarting for an earlier offset.
type lspCursor struct {
	script    string
	offset    int
	line      uint32
	character uint32
}

// position returns the LSP position of target. An offset inside a code point
// is the zero position.
func (c *lspCursor) position(target int) LSPPosition {
	if target < c.offset {
		c.offset, c.line, c.character = 0, 0, 0
	}
	for c.offset < target {
		r, size := utf8.DecodeRuneInString(c.script[c.offset:])
		if c.offset+size > target {
			return LSPPosition{}
		}
		switch {
		case r == '\n':
			c.line++
			c.character = 0
		case r > 0xFFFF:
			c.character += 2
		default:
			c.character++
		}
		c.offset += size
	}
	return LSPPosition{Line: c.line, Character: c.character}
}

// TrimStatement trims a statement's surrounding whitespace and its trailing
// semicolons, for an engine that wraps it in a larger statement.
func TrimStatement(statement string) string {
	return strings.TrimLeftFunc(strings.TrimRightFunc(statement, func(r rune) bool {
		return r == ';' || unicode.IsSpace(r)
	}), unicode.IsSpace)
}
