package parser

// maxNestingDepth bounds how deeply the parser may descend into nested
// constructs: expressions, collection literals, function arguments, INITCOND
// literals, and data types. Without a bound, about 880KB of '[' overflowed
// the Go stack, and a stack overflow kills the process instead of returning
// an error.
const maxNestingDepth = 10000

const nestingErrorMessage = "statement is nested too deeply"

// enterNesting is called on entry to the productions that recurse, placed so
// that every recursion cycle in the parser passes through one. It fails when
// the statement nests deeper than maxNestingDepth; otherwise the caller
// defers leaveNesting.
func (p *Parser) enterNesting() error {
	if p.depth >= maxNestingDepth {
		return p.errorf(nestingErrorMessage)
	}
	p.depth++
	return nil
}

func (p *Parser) leaveNesting() { p.depth-- }
