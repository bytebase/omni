package parser

import "errors"

// maxNestingDepth bounds how deeply the parser may descend into nested
// constructs: expressions, queries, table references, stored-program blocks,
// and the other productions that recurse. Without a bound, a few hundred
// kilobytes of '(' overflowed the Go stack, and a stack overflow kills the
// process instead of returning an error. TiDB 8.5 sets no nesting limit of
// its own: it accepts 50,000 nested parentheses or function calls and 20,000
// nested subqueries.
// It is a variable only so tests that nest deeper on purpose can raise it.
var maxNestingDepth = 10000

// enterNesting is called on entry to the productions that recurse, placed so
// that every recursion cycle in the parser passes through one. It fails when
// the statement nests deeper than maxNestingDepth; otherwise the caller
// defers leaveNesting.
func (p *Parser) enterNesting() error {
	if p.depth >= maxNestingDepth {
		return &ParseError{Message: nestingErrorMessage, Position: p.cur.Loc}
	}
	p.depth++
	return nil
}

func (p *Parser) leaveNesting() { p.depth-- }

const nestingErrorMessage = "statement is nested too deeply"

// isNestingError reports whether err is enterNesting's failure.
func isNestingError(err error) bool {
	var pe *ParseError
	return errors.As(err, &pe) && pe.Message == nestingErrorMessage
}
