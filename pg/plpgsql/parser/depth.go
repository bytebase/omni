package parser

// maxNestingDepth bounds how deeply blocks and control statements (BEGIN,
// IF, LOOP, CASE, ...) may nest. Without a bound, deeply nested blocks
// overflowed the Go stack, and a stack overflow kills the process instead of
// returning an error.
const maxNestingDepth = 10000

// enterNesting is called on entry to every production that recurses. It
// fails when the body nests deeper than maxNestingDepth; otherwise the caller
// defers leaveNesting.
func (p *Parser) enterNesting() error {
	if p.depth >= maxNestingDepth {
		return &ParseError{Message: "statement is nested too deeply", Position: p.cur.Loc}
	}
	p.depth++
	return nil
}

func (p *Parser) leaveNesting() { p.depth-- }
