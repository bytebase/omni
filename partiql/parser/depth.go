package parser

import "errors"

// maxNestingDepth bounds how deeply the parser may descend into nested
// constructs: expressions, queries, table references, and graph patterns.
// Without a bound, about 300KB of '(' overflowed the Go stack, and a stack
// overflow kills the process instead of returning an error.
const maxNestingDepth = 10000

const nestingErrorMessage = "statement is nested too deeply"

// enterNesting is called on entry to the productions that recurse, placed so
// that every recursion cycle in the parser passes through one. It fails when
// the statement nests deeper than maxNestingDepth; otherwise the caller
// defers leaveNesting.
func (p *Parser) enterNesting() error {
	if p.depth >= maxNestingDepth {
		return &ParseError{
			Message:  nestingErrorMessage,
			Position: p.cur.Loc.Start, End: p.cur.Loc.End,
		}
	}
	p.depth++
	return nil
}

func (p *Parser) leaveNesting() { p.depth-- }

// isNestingError reports whether err is the error enterNesting returns.
func isNestingError(err error) bool {
	var pe *ParseError
	return errors.As(err, &pe) && pe.Message == nestingErrorMessage
}
