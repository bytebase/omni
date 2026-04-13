package parser

// Segment represents one top-level SQL statement extracted from a source string.
//
// Text is the raw substring of the input from ByteStart (inclusive) to
// ByteEnd (exclusive). The trailing `;` delimiter (if present) is NOT part
// of Text or ByteEnd.
//
// NOTE: This is a minimal stub to support the parser framework (F4). It will
// be replaced when the split branch (F3) merges.
type Segment struct {
	Text      string // the raw text of the statement (no trailing semicolon)
	ByteStart int    // inclusive start byte offset in the original source
	ByteEnd   int    // exclusive end byte offset
}

// Split extracts top-level SQL statements from input by splitting on
// top-level semicolons.
//
// NOTE: This is a minimal stub. It handles the common case of
// semicolon-separated statements but does not handle strings, comments, or
// nesting. It will be replaced by the full F3 Split implementation.
func Split(input string) []Segment {
	if len(input) == 0 {
		return nil
	}

	var segments []Segment
	stmtStart := 0

	l := NewLexer(input)
	for {
		tok := l.NextToken()
		if tok.Kind == tokEOF {
			break
		}
		if tok.Kind == int(';') {
			seg := Segment{
				Text:      input[stmtStart:tok.Loc.Start],
				ByteStart: stmtStart,
				ByteEnd:   tok.Loc.Start,
			}
			if !segmentEmpty(seg) {
				segments = append(segments, seg)
			}
			stmtStart = tok.Loc.End
		}
	}

	// Trailing segment (no trailing semicolon).
	if stmtStart < len(input) {
		seg := Segment{
			Text:      input[stmtStart:],
			ByteStart: stmtStart,
			ByteEnd:   len(input),
		}
		if !segmentEmpty(seg) {
			segments = append(segments, seg)
		}
	}

	return segments
}

// segmentEmpty reports whether a segment contains no meaningful SQL content
// (only whitespace and comments).
func segmentEmpty(s Segment) bool {
	return NewLexer(s.Text).NextToken().Kind == tokEOF
}
