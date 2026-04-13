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
			for i < len(t) && t[i] != '\n' {
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
			// Ion literal: scan with awareness of Ion-mode strings and
			// comments so that backticks inside them don't prematurely
			// close the literal. Matches PartiQLLexer.g4's ION mode
			// (lines 408-428): Ion short strings ("..."), Ion symbols
			// ('...'), Ion long strings ('''...'''), and Ion inline
			// comments (/*...*/) can all contain backticks.
			i++ // skip opening `
			scanIon(input, &i)

		case ch == '-' && i+1 < n && input[i+1] == '-':
			// Line comment: skip to end of line.
			i += 2
			for i < n && input[i] != '\n' {
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

// scanIon advances *pos past an Ion literal body (after the opening
// backtick has been consumed). It handles Ion short strings ("..."),
// Ion symbols ('...'), Ion long strings (”'...”'), and Ion inline
// comments (/*...*/) — all of which can contain backticks without
// terminating the literal. Matches PartiQLLexer.g4's ION mode.
func scanIon(input string, pos *int) {
	n := len(input)
	i := *pos
	for i < n {
		switch input[i] {
		case '`':
			i++ // closing backtick
			*pos = i
			return
		case '"':
			// Ion short string: "..." with backslash escapes.
			i++
			for i < n && input[i] != '"' {
				if input[i] == '\\' && i+1 < n {
					i++ // skip escaped char
				}
				i++
			}
			if i < n {
				i++ // closing "
			}
		case '\'':
			// Check for triple-quote Ion long string '''...'''.
			if i+2 < n && input[i+1] == '\'' && input[i+2] == '\'' {
				i += 3
				for i+2 < n {
					if input[i] == '\'' && input[i+1] == '\'' && input[i+2] == '\'' {
						i += 3
						break
					}
					i++
				}
			} else {
				// Ion symbol: '...' with backslash escapes.
				i++
				for i < n && input[i] != '\'' {
					if input[i] == '\\' && i+1 < n {
						i++
					}
					i++
				}
				if i < n {
					i++ // closing '
				}
			}
		case '/':
			// Ion inline comment: /*...*/ (non-nesting per Ion spec).
			if i+1 < n && input[i+1] == '*' {
				i += 2
				for i+1 < n {
					if input[i] == '*' && input[i+1] == '/' {
						i += 2
						break
					}
					i++
				}
			} else {
				i++
			}
		default:
			i++
		}
	}
	*pos = i
}
