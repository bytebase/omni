package parser

// Split splits a CQL input into segments at top-level semicolons, correctly
// handling string literals, quoted identifiers, code blocks and comments.
//
// A semicolon inside a batch does not end the statement: BEGIN … BATCH and
// APPLY BATCH frame one statement whose inner statements are separated by
// semicolons, and the server reads the batch as one request.
func Split(sql string) []Segment {
	var segments []Segment
	start := 0
	inBatch := false
	emit := func(end int) {
		seg := makeSegment(sql, start, end)
		if !seg.Empty {
			segments = append(segments, seg)
		}
	}
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
			text := sql[start:i]
			switch {
			case inBatch:
				// Inside a batch a semicolon separates the inner statements;
				// only APPLY BATCH closes it.
				if endsBatch(text) {
					emit(i)
					inBatch = false
					start = i + 1
				}
			case beginsBatch(text):
				inBatch = true
			default:
				emit(i)
				start = i + 1
			}
			i++
		default:
			i++
		}
	}
	// Trailing segment after last semicolon (or the whole input if no semicolons).
	if start < len(sql) {
		emit(len(sql))
	}
	return segments
}

// beginsBatch reports whether text's first tokens are BEGIN, optionally
// UNLOGGED, LOGGED or COUNTER, then BATCH.
func beginsBatch(text string) bool {
	lexer := NewLexer(text)
	first := lexer.Next()
	if first.Type != tokBEGIN {
		return false
	}
	next := lexer.Next()
	switch next.Type {
	case tokUNLOGGED, tokLOGGED, tokCOUNTER:
		return lexer.Next().Type == tokBATCH
	case tokBATCH:
		return true
	default:
		return false
	}
}

// endsBatch reports whether text's last two tokens are APPLY BATCH, with only
// whitespace and comments after them.
func endsBatch(text string) bool {
	lexer := NewLexer(text)
	var last, secondLast Token
	for {
		token := lexer.Next()
		if token.Type == tokEOF || lexer.Err != nil {
			break
		}
		secondLast, last = last, token
	}
	return secondLast.Type == tokAPPLY && last.Type == tokBATCH
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
