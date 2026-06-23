package parser

// Split splits a CQL input into segments at top-level semicolons,
// correctly handling string literals, quoted identifiers, code blocks,
// and comments.
func Split(sql string) []Segment {
	var segments []Segment
	start := 0
	i := 0
	for i < len(sql) {
		ch := sql[i]
		switch {
		case ch == '\'':
			i = skipString(sql, i)
		case ch == '"':
			i = skipQuotedIdent(sql, i)
		case ch == '$' && i+1 < len(sql) && sql[i+1] == '$':
			i = skipCodeBlock(sql, i)
		case ch == '-' && i+1 < len(sql) && sql[i+1] == '-':
			i = skipLineComment(sql, i)
		case ch == '/' && i+1 < len(sql) && sql[i+1] == '*':
			i = skipBlockComment(sql, i)
		case ch == ';':
			seg := makeSegment(sql, start, i)
			if !seg.Empty {
				segments = append(segments, seg)
			}
			i++
			start = i
		default:
			i++
		}
	}
	// Trailing segment after last semicolon (or the whole input if no semicolons).
	if start < len(sql) {
		seg := makeSegment(sql, start, len(sql))
		if !seg.Empty {
			segments = append(segments, seg)
		}
	}
	return segments
}

// Segment represents a single SQL segment from splitting.
type Segment struct {
	Text      string
	ByteStart int
	ByteEnd   int
	Empty     bool
}

func makeSegment(sql string, start, end int) Segment {
	for start < end && isWhitespace(sql[start]) {
		start++
	}
	for end > start && isWhitespace(sql[end-1]) {
		end--
	}
	text := sql[start:end]
	empty := len(text) == 0
	return Segment{
		Text:      text,
		ByteStart: start,
		ByteEnd:   end,
		Empty:     empty,
	}
}

func isWhitespace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f'
}

func skipString(sql string, i int) int {
	i++ // skip opening '
	for i < len(sql) {
		if sql[i] == '\'' {
			i++
			if i < len(sql) && sql[i] == '\'' {
				i++ // escaped quote ''
				continue
			}
			return i
		}
		i++
	}
	return i
}

func skipQuotedIdent(sql string, i int) int {
	i++ // skip opening "
	for i < len(sql) {
		if sql[i] == '"' {
			i++
			if i < len(sql) && sql[i] == '"' {
				i++ // escaped quote ""
				continue
			}
			return i
		}
		i++
	}
	return i
}

func skipCodeBlock(sql string, i int) int {
	i += 2 // skip opening $$
	for i+1 < len(sql) {
		if sql[i] == '$' && sql[i+1] == '$' {
			return i + 2
		}
		i++
	}
	return len(sql)
}

func skipLineComment(sql string, i int) int {
	i += 2 // skip --
	for i < len(sql) && sql[i] != '\n' {
		i++
	}
	return i
}

func skipBlockComment(sql string, i int) int {
	i += 2 // skip /*
	for i+1 < len(sql) {
		if sql[i] == '*' && sql[i+1] == '/' {
			return i + 2
		}
		i++
	}
	return len(sql)
}
