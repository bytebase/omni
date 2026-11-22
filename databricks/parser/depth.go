package parser

// maxNestingDepth bounds how deeply the parser may descend into nested
// constructs: expressions, subqueries, relations, grouping elements, data
// types, compound script blocks, and unary/NOT chains.
const maxNestingDepth = 10000

const nestingErrorMessage = "maximum nesting depth exceeded"

// enterNesting is called on entry to recursive productions so that every
// recursion cycle passes through one check. It fails when the statement nests
// deeper than maxNestingDepth; otherwise the caller defers leaveNesting.
func (p *Parser) enterNesting() error {
	if p.depth >= maxNestingDepth {
		return p.errorf(nestingErrorMessage)
	}
	p.depth++
	return nil
}

func (p *Parser) leaveNesting() {
	p.depth--
}
