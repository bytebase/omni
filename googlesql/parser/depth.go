package parser

// maxNestingDepth bounds how deeply the parser may descend into nested
// constructs: expressions, queries, join groups, types, script blocks, and
// the other productions that recurse. Without a bound, a few hundred
// kilobytes of '(' overflowed the Go stack, and a stack overflow kills the
// process instead of returning an error. The Spanner emulator (ZetaSQL)
// itself gives up on 3,000 nested subqueries and on 10,000 nested unary
// minuses ("Out of stack space due to deeply nested query expression").
// It is a variable only so tests that nest deeper on purpose can raise it.
var maxNestingDepth = 10000

// nestingErrorMessage is the message of the error enterNesting returns.
const nestingErrorMessage = "statement is nested too deeply"

// enterNesting is called on entry to the productions that recurse, placed so
// that every recursion cycle in the parser passes through one. It fails when
// the statement nests deeper than maxNestingDepth; otherwise the caller
// defers leaveNesting.
func (p *Parser) enterNesting() error {
	if p.depth >= maxNestingDepth {
		return &ParseError{Message: nestingErrorMessage, Position: p.cur.Loc.Start, End: p.cur.Loc.End}
	}
	p.depth++
	return nil
}

func (p *Parser) leaveNesting() { p.depth-- }

// isNestingError reports whether err is the error enterNesting returns.
func isNestingError(err error) bool {
	pe, ok := err.(*ParseError)
	return ok && pe.Message == nestingErrorMessage
}
