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
// backtick-quoted identifiers, line comments (-- and #), and nested block
// comments (/* ... */) do not split statements.
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
	var segments []Segment
	start := 0
	for {
		tok := l.NextToken()
		if tok.Type == TokEOF {
			break
		}
		if tok.Type == TokSemicolon {
			segEnd := extendSameLineComment(script, tok.End)
			if segEnd > l.pos {
				l.pos = segEnd
			}
			segments = append(segments, Segment{
				Text:      script[start:segEnd],
				ByteStart: start,
				ByteEnd:   segEnd,
			})
			start = segEnd
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
