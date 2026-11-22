package parser

// Segment is one segment of a HiveQL script delimited by top-level semicolons.
// Text is script[ByteStart:ByteEnd], preserving leading whitespace/comments and
// including the terminating semicolon when present.
type Segment struct {
	Text      string
	ByteStart int
	ByteEnd   int
}

// Empty reports whether the segment holds no HiveQL statement (only whitespace,
// comments, and semicolons). Malformed spans such as unclosed block comments
// emit TokInvalid and are therefore non-empty.
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

// Split splits a HiveQL script into segments at top-level semicolons using the
// Hive 4.0.1 lexer. Semicolons inside string literals, backtick-quoted
// identifiers, line comments, and nested block comments do not split
// statements. Empty segments (lone semicolons, whitespace, or comment-only
// chunks) are retained and report Empty() == true so callers can account for
// every semicolon. Lexical errors do not abort splitting; malformed spans stay
// in their segment so Parse can report them at their byte offset.
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
			segments = append(segments, Segment{
				Text:      script[start:tok.Loc.End],
				ByteStart: start,
				ByteEnd:   tok.Loc.End,
			})
			start = tok.Loc.End
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
