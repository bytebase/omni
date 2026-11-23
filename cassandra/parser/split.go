package parser

// Split splits a CQL input into segments at top-level semicolons, correctly
// handling string literals, quoted identifiers, code blocks and comments.
// Segments do not include the semicolon. A segment holding no statement (a
// lone semicolon, or only whitespace and comments before one or at the end of
// the input) is returned too, with Empty set, so a caller can account for
// every semicolon; callers that want statements skip it.
//
// A semicolon inside a batch does not end the statement: BEGIN … BATCH and
// APPLY BATCH frame one statement whose inner statements are separated by
// semicolons, and the server reads the batch as one request.
func Split(sql string) []Segment {
	var segments []Segment
	start := 0
	// pieceStart is where the statement since the last semicolon begins; the
	// batch checks read only it, so a long batch is not re-lexed at every inner
	// semicolon.
	pieceStart := 0
	inBatch := false
	emit := func(end int) {
		segments = append(segments, makeSegment(sql, start, end))
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
		case (ch == '-' || ch == '/') && i+1 < len(sql) && sql[i+1] == ch:
			// "--" and "//" both start a line comment, as in the lexer.
			i = skipLineComment(sql, i)
		case ch == '/' && i+1 < len(sql) && sql[i+1] == '*':
			i = skipBlockComment(sql, i)
		case ch == ';':
			piece := sql[pieceStart:i]
			switch {
			case inBatch:
				// Inside a batch a semicolon separates the inner statements;
				// only APPLY BATCH closes it.
				if endsBatch(piece) {
					emit(i)
					inBatch = false
					start = i + 1
				}
			case beginsBatch(sql[start:i]):
				// A batch that closes at once — BEGIN BATCH APPLY BATCH — is
				// one statement too.
				if endsBatch(piece) {
					emit(i)
					start = i + 1
				} else {
					inBatch = true
				}
			default:
				emit(i)
				start = i + 1
			}
			pieceStart = i + 1
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

// Segment represents a single SQL segment from splitting. Text is the
// segment with its surrounding whitespace trimmed, and ByteStart and ByteEnd
// are its range in the input. Empty reports that it holds no statement: it is
// only whitespace and comments, or nothing at all.
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
	return Segment{
		Text:      text,
		ByteStart: start,
		ByteEnd:   end,
		Empty:     holdsNoToken(text),
	}
}

// holdsNoToken reports whether text is only whitespace and comments, read as
// the lexer reads them. An unterminated block comment is a lexer error, not
// an empty segment: the parser's verdict on it is what the caller should read.
func holdsNoToken(text string) bool {
	lexer := NewLexer(text)
	lexer.skipWhitespaceAndComments()
	return lexer.Err == nil && lexer.pos == len(text)
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
	i += 2 // skip -- or //
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
