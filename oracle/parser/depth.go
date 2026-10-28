package parser

import "errors"

// maxNestingDepth bounds how deeply the parser may descend into nested
// constructs: expressions, queries, table references, PL/SQL blocks and
// declarations, and the other productions that recurse. Without a bound,
// about 700,000 nested parentheses overflowed the Go stack, and a stack
// overflow kills the process instead of returning an error. Oracle 23ai
// itself rejects 8,001 nested parentheses or function calls (ORA-03095),
// 256 nested subqueries in an expression (ORA-01778), and 99 nested PL/SQL
// blocks (PLS-00103), but takes 10,001 nested derived tables, which this
// limit rejects.
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

// nestingErrorMessage is the message of the error enterNesting fails with.
const nestingErrorMessage = "statement is nested too deeply"

// nestingError returns err as the error enterNesting fails with, or nil.
func nestingError(err error) *ParseError {
	var pe *ParseError
	if errors.As(err, &pe) && pe.Message == nestingErrorMessage {
		return pe
	}
	return nil
}
