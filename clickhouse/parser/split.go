package parser

// Segment is one segment of a ClickHouse SQL script delimited by top-level
// semicolons. Text is script[ByteStart:ByteEnd], preserving leading whitespace
// and comments, the terminating semicolon when present, and any same-line
// trailing line comment (`-- ...` or `# ...`) after the semicolon.
type Segment struct {
	Text      string
	ByteStart int
	ByteEnd   int
}

// Empty reports whether the segment holds no ClickHouse SQL statement (only
// whitespace, comments, and semicolons). Malformed spans such as unclosed
// strings, dollar-quoted strings, or block comments emit TokIllegal and are
// therefore non-empty.
func (s Segment) Empty() bool {
	l := NewLexer(s.Text)
	for {
		tok := l.NextToken()
		switch tok.Type {
		case TokEOF:
			return true
		case TokSemicolon:
			continue
		default:
			return false
		}
	}
}

// Split splits a ClickHouse SQL script into segments at top-level semicolons
// using the ClickHouse lexer. Semicolons inside single-quoted strings,
// dollar-quoted strings ($$...$$ / $tag$...$tag$), double-quoted or
// backtick-quoted identifiers, line comments (-- and #), nested block comments
// (/* ... */), and inline INSERT ... FORMAT <name> data payloads do not split
// statements.
//
// Byte spans are contiguous across [0, len(script)) with Text ==
// script[ByteStart:ByteEnd]. Empty segments (lone semicolons, whitespace, or
// comment-only chunks) are retained with Empty() == true so callers can
// account for every byte and semicolon. Lexical errors do not abort splitting;
// malformed spans stay in their segment as non-empty so Parse or callers can
// report them at their byte offset.
func Split(script string) []Segment {
	if len(script) == 0 {
		return nil
	}
	l := NewLexer(script)
	var buf []Token
	nextTok := func() Token {
		if len(buf) > 0 {
			t := buf[0]
			buf = buf[1:]
			return t
		}
		return l.NextToken()
	}
	peekTok := func() Token {
		if len(buf) == 0 {
			buf = append(buf, l.NextToken())
		}
		return buf[0]
	}

	var segments []Segment
	start := 0
	stmtStart := true
	for {
		tok := nextTok()
		if tok.Type == TokEOF {
			break
		}
		if tok.Type == TokSemicolon {
			segEnd := extendSameLineComment(script, tok.End)
			if segEnd > tok.End {
				buf = buf[:0]
				l.pos = segEnd
			}
			segments = append(segments, Segment{
				Text:      script[start:segEnd],
				ByteStart: start,
				ByteEnd:   segEnd,
			})
			start = segEnd
			stmtStart = true
			continue
		}
		if stmtStart {
			stmtStart = false
			if tok.Type == KwExplain && peekTok().Type == KwAst {
				nextTok() // consume AST
				if peekTok().Type == KwInsert {
					tok = nextTok() // consume INSERT
				}
			}
			if tok.Type == KwInsert && consumeInsertFormatHeader(nextTok, peekTok) {
				// Per ClickHouse's INSERT ... FORMAT <name> syntax (and
				// parseInsertStmt), all remaining bytes after FORMAT <name>
				// are raw inline payload data and must not be split on ';'.
				break
			}
		}
	}
	if start < len(script) {
		segments = append(segments, Segment{
			Text:      script[start:],
			ByteStart: start,
			ByteEnd:   len(script),
		})
	}
	return segments
}

// consumeInsertFormatHeader scans the SQL header of an INSERT statement (after
// KwInsert has already been consumed) and reports whether it reaches a
// `FORMAT <identifier>` inline-data boundary. Any top-level semicolon
// encountered before the boundary is left unconsumed in peekTok so Split can
// still emit the segment boundary.
func consumeInsertFormatHeader(nextTok, peekTok func() Token) bool {
	if peekTok().Type != KwInto {
		return false
	}
	nextTok() // INTO
	if peekTok().Type == KwTable {
		nextTok() // TABLE
	}

	// Target: FUNCTION fn(...) or [db.]table
	if peekTok().Type == KwFunction {
		nextTok() // FUNCTION
		if !isIdentifier(peekTok()) {
			return false
		}
		nextTok() // function name
		if peekTok().Type != TokLParen {
			return false
		}
		if !skipBalancedParens(nextTok, peekTok) {
			return false
		}
	} else {
		if !isIdentifier(peekTok()) {
			return false
		}
		nextTok() // db or table
		if peekTok().Type == TokDot {
			nextTok() // .
			if !isIdentifier(peekTok()) {
				return false
			}
			nextTok() // table
		}
	}

	// Optional column list: (col1, col2, ...), excluding parenthesized SELECT.
	if peekTok().Type == TokLParen {
		nextTok() // (
		switch peekTok().Type {
		case KwSelect, KwWith, KwFrom, TokLParen:
			return false
		}
		depth := 1
		for depth > 0 {
			t := peekTok()
			if t.Type == TokEOF || t.Type == TokSemicolon {
				return false
			}
			nextTok()
			switch t.Type {
			case TokLParen:
				depth++
			case TokRParen:
				depth--
			}
		}
	}

	// Optional SETTINGS clause before FORMAT / VALUES / SELECT.
	if peekTok().Type == KwSettings {
		nextTok() // SETTINGS
		depth := 0
		for {
			t := peekTok()
			if t.Type == TokEOF || t.Type == TokSemicolon {
				return false
			}
			if depth == 0 {
				if t.Type == KwFormat {
					nextTok() // FORMAT
					if peekTok().Type == TokEq {
						// Setting named `format = ...`, not the FORMAT clause.
						continue
					}
					return isIdentifier(peekTok())
				}
				if t.Type == KwValues || t.Type == KwSelect || t.Type == KwWith || t.Type == TokLParen {
					return false
				}
			}
			nextTok()
			switch t.Type {
			case TokLParen, TokLBracket, TokLBrace:
				depth++
			case TokRParen, TokRBracket, TokRBrace:
				if depth > 0 {
					depth--
				}
			}
		}
	}

	if peekTok().Type == KwFormat {
		nextTok() // FORMAT
		return isIdentifier(peekTok())
	}
	return false
}

func skipBalancedParens(nextTok, peekTok func() Token) bool {
	if peekTok().Type != TokLParen {
		return false
	}
	nextTok() // (
	depth := 1
	for depth > 0 {
		t := peekTok()
		if t.Type == TokEOF || t.Type == TokSemicolon {
			return false
		}
		nextTok()
		switch t.Type {
		case TokLParen:
			depth++
		case TokRParen:
			depth--
		}
	}
	return true
}

// extendSameLineComment advances past optional horizontal spaces/tabs and a
// single-line comment (`-- ...` or `# ...`) on the same line as a terminating
// semicolon, stopping before `\r`, `\n`, or EOF. If no such comment follows on
// the same line, it returns pos unchanged.
func extendSameLineComment(script string, pos int) int {
	i := pos
	for i < len(script) && (script[i] == ' ' || script[i] == '\t') {
		i++
	}
	if (i+1 < len(script) && script[i] == '-' && script[i+1] == '-') ||
		(i < len(script) && script[i] == '#') {
		for i < len(script) && script[i] != '\n' && script[i] != '\r' {
			i++
		}
		return i
	}
	return pos
}
