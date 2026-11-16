package cosmosdb

// Segment is a portion of a Cosmos DB script delimited by a top-level
// semicolon. Split produces it.
//
// Text is the statement the engine runs: the segment without its terminating
// semicolon, since the query language has no statement terminator and the
// gateway rejects one. ByteStart and ByteEnd delimit the same text in the
// script.
type Segment struct {
	Text      string
	ByteStart int
	ByteEnd   int
}

// Empty reports a segment that holds no statement: only whitespace and
// comments. A lone semicolon leaves one behind, so a caller that accounts for
// every semicolon of a script can see it.
func (s Segment) Empty() bool {
	t := s.Text
	for i := 0; i < len(t); {
		switch {
		case t[i] == ' ' || t[i] == '\t' || t[i] == '\n' || t[i] == '\r' || t[i] == '\f' || t[i] == ';':
			i++
		case t[i] == '-' && i+1 < len(t) && t[i+1] == '-':
			for i < len(t) && t[i] != '\n' {
				i++
			}
		default:
			return false
		}
	}
	return true
}

// Split splits a script at its top-level semicolons. It is lexical and never
// fails, so the boundaries of text the query language does not parse are its
// punctuation and not its grammar. Semicolons inside a string literal and
// semicolons a comment holds do not split: the scan matches the lexer, where a
// string is delimited by ' or ", backslash escapes the next byte, and a comment
// runs from -- to the end of its line.
//
// The last segment is the text after the last semicolon, or nil when the script
// ends with one.
func Split(script string) []Segment {
	if script == "" {
		return nil
	}

	var segments []Segment
	start := 0
	for i := 0; i < len(script); {
		switch {
		case script[i] == ';':
			segments = append(segments, segmentAt(script, start, i))
			i++
			start = i
		case script[i] == '\'' || script[i] == '"':
			i = skipLiteral(script, i)
		case script[i] == '-' && i+1 < len(script) && script[i+1] == '-':
			for i < len(script) && script[i] != '\n' {
				i++
			}
		default:
			i++
		}
	}
	if start < len(script) {
		segments = append(segments, segmentAt(script, start, len(script)))
	}
	return segments
}

// segmentAt returns script[start:end] with the whitespace after the last token
// left out, so a segment's text is the statement and not the layout around it.
func segmentAt(script string, start, end int) Segment {
	for end > start && isSpace(script[end-1]) {
		end--
	}
	return Segment{Text: script[start:end], ByteStart: start, ByteEnd: end}
}

// skipLiteral returns the offset just past the string literal that starts at i,
// or len(script) when the literal does not close: a backslash escapes the byte
// after it, as in the lexer.
func skipLiteral(script string, i int) int {
	quote := script[i]
	for i++; i < len(script); i++ {
		if script[i] == '\\' {
			i++
			continue
		}
		if script[i] == quote {
			return i + 1
		}
	}
	return len(script)
}
