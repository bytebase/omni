package parser

// Segment is one top-level HiveQL statement extracted from a script. Text is
// script[ByteStart:ByteEnd], preserving leading whitespace and comments since
// the previous segment and including the terminating semicolon when present.
type Segment struct {
	Text      string
	ByteStart int
	ByteEnd   int
}

// Empty reports whether the segment holds no HiveQL tokens (only whitespace
// and comments).
func (s Segment) Empty() bool {
	return NewLexer(s.Text).NextToken().Type == TokEOF
}

// Split splits a HiveQL script into non-empty statements using the Hive 4.0.1
// lexer. Semicolons inside string literals, backtick-quoted identifiers, line
// comments, and nested block comments do not split statements. Lexical errors
// do not abort splitting; malformed spans stay in their segment so Parse can
// report them at their byte offset.
func Split(script string) []Segment {
	l := NewLexer(script)
	var segments []Segment
	prevEnd := 0
	hasToken := false
	lastEnd := 0
	for {
		tok := l.NextToken()
		if tok.Type == TokEOF {
			break
		}
		if tok.Type == TokSemicolon {
			if hasToken {
				segments = append(segments, Segment{
					Text:      script[prevEnd:tok.Loc.End],
					ByteStart: prevEnd,
					ByteEnd:   tok.Loc.End,
				})
				hasToken = false
			}
			prevEnd = tok.Loc.End
			continue
		}
		hasToken = true
		lastEnd = tok.Loc.End
	}
	if hasToken {
		segments = append(segments, Segment{
			Text:      script[prevEnd:lastEnd],
			ByteStart: prevEnd,
			ByteEnd:   lastEnd,
		})
	}
	return segments
}
