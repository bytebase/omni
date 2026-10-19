package parser

// maxNestingDepth bounds how deeply the parser may descend into nested
// constructs: expressions, queries, types, and the other productions that
// recurse. Without a bound, a few hundred kilobytes of '(' overflowed the Go
// stack, and a stack overflow kills the process instead of returning an
// error. PostgreSQL 17 itself rejects 10,000 nested parentheses ("memory
// exhausted") and 3,000 nested subqueries ("stack depth limit exceeded").
// It is a variable only so tests that nest deeper on purpose can raise it.
var maxNestingDepth = 10000

// enterNesting is called on entry to the productions that recurse, placed so
// that every recursion cycle in the parser passes through one. It fails when
// the statement nests deeper than maxNestingDepth; otherwise the caller
// defers leaveNesting.
func (p *Parser) enterNesting() error {
	if p.depth >= maxNestingDepth {
		return &ParseError{Code: "54001", Message: "statement is nested too deeply", Position: p.cur.Loc}
	}
	p.depth++
	return nil
}

func (p *Parser) leaveNesting() { p.depth-- }
