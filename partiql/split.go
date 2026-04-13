// Package partiql provides a hand-written parser for the PartiQL query language
// (AWS DynamoDB, Azure Cosmos DB).
package partiql

// Segment represents a portion of PartiQL text delimited by top-level
// semicolons. Produced by Split.
type Segment struct {
	Text      string // the raw text of this segment (may include trailing ;)
	ByteStart int    // byte offset of start in original input
	ByteEnd   int    // byte offset of end (exclusive) in original input
}

// Empty returns true if the segment contains only whitespace, comments,
// and semicolons — no meaningful PartiQL content.
func (s Segment) Empty() bool {
	t := s.Text
	i := 0
	for i < len(t) {
		b := t[i]
		// Skip whitespace and semicolons.
		if b == ' ' || b == '\t' || b == '\n' || b == '\r' || b == ';' {
			i++
			continue
		}
		// Skip line comments (-- to end of line).
		if b == '-' && i+1 < len(t) && t[i+1] == '-' {
			i += 2
			for i < len(t) && t[i] != '\n' && t[i] != '\r' {
				i++
			}
			continue
		}
		// Skip block comments (/* ... */), supporting nesting.
		if b == '/' && i+1 < len(t) && t[i+1] == '*' {
			i += 2
			depth := 1
			for i+1 < len(t) && depth > 0 {
				if t[i] == '/' && t[i+1] == '*' {
					depth++
					i += 2
				} else if t[i] == '*' && t[i+1] == '/' {
					depth--
					i += 2
				} else {
					i++
				}
			}
			continue
		}
		return false
	}
	return true
}

// Split splits a PartiQL input into segments at top-level semicolons.
// It respects single-quoted strings (with ” escape), double-quoted
// identifiers (with "" escape), backtick-delimited Ion literals, line
// comments (--), and block comments (/* */ with nesting).
//
// Each segment's Text includes its trailing semicolon if present.
// The final segment may not have a trailing semicolon.
func Split(input string) []Segment {
	var segs []Segment
	n := len(input)
	start := 0

	i := 0
	for i <= n {
		if i == n {
			// End of input: flush remaining text.
			if i > start {
				segs = append(segs, Segment{
					Text:      input[start:],
					ByteStart: start,
					ByteEnd:   n,
				})
			}
			break
		}

		ch := input[i]

		switch {
		case ch == ';':
			// Include the semicolon in the current segment.
			segs = append(segs, Segment{
				Text:      input[start : i+1],
				ByteStart: start,
				ByteEnd:   i + 1,
			})
			start = i + 1
			i++

		case ch == '\'':
			// Single-quoted string: skip to closing quote. '' is escape.
			i++
			for i < n {
				if input[i] == '\'' {
					i++
					if i < n && input[i] == '\'' {
						i++ // escaped ''
						continue
					}
					break
				}
				i++
			}

		case ch == '"':
			// Double-quoted identifier: skip to closing quote. "" is escape.
			i++
			for i < n {
				if input[i] == '"' {
					i++
					if i < n && input[i] == '"' {
						i++ // escaped ""
						continue
					}
					break
				}
				i++
			}

		case ch == '`':
			// Ion literal: skip to closing backtick.
			i++
			for i < n && input[i] != '`' {
				i++
			}
			if i < n {
				i++ // consume closing `
			}

		case ch == '-' && i+1 < n && input[i+1] == '-':
			// Line comment: skip to end of line (\n or \r).
			i += 2
			for i < n && input[i] != '\n' && input[i] != '\r' {
				i++
			}

		case ch == '/' && i+1 < n && input[i+1] == '*':
			// Block comment: skip to */. Supports nesting.
			i += 2
			depth := 1
			for i+1 < n && depth > 0 {
				if input[i] == '/' && input[i+1] == '*' {
					depth++
					i += 2
				} else if input[i] == '*' && input[i+1] == '/' {
					depth--
					i += 2
				} else {
					i++
				}
			}

		default:
			i++
		}
	}

	return segs
}
