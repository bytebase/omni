package parser

// maxNestingDepth bounds how deeply the parser may descend into nested
// constructs: expressions, queries, relations, types, row patterns, routine
// statements, and the statements that hold another statement. Without a
// bound, a few hundred kilobytes of '(' overflowed the Go stack, and a stack
// overflow kills the process instead of returning an error. Trino 482 itself
// rejects a few thousand nested parentheses and a few hundred nested
// subqueries ("statement is too large (stack overflow while parsing)"). It is
// a variable only so tests that nest deeper on purpose can raise it.
var maxNestingDepth = 10000

// nestingTooDeep is the message of the error a statement fails with when it
// nests deeper than maxNestingDepth.
const nestingTooDeep = "statement is nested too deeply"

// enterNesting is called on entry to the productions that recurse, placed so
// that every recursion cycle in the parser passes through one. It fails when
// the statement nests deeper than maxNestingDepth; otherwise the caller
// defers leaveNesting.
//
// The first failure is kept, and every later call fails with it: the parse
// continues only to unwind, and the statement fails with that error (see
// parseSingle) whatever reading a speculative parse falls back to. The
// readings the parser speculates on, and the failures it records by offset
// (parenQueryFail, parenRelationFail), assume that a parse depends only on
// where it starts, but a nesting failure depends on how deep the parse
// started too. Failing the whole statement keeps that assumption sound, and
// a fallback reading cannot descend to the limit again and again.
func (p *Parser) enterNesting() error {
	if p.nestingErr == nil && p.depth >= maxNestingDepth {
		p.nestingErr = &ParseError{Position: p.cur.Loc.Start, End: p.cur.Loc.End, Message: nestingTooDeep}
	}
	if p.nestingErr != nil {
		e := *p.nestingErr
		return &e
	}
	p.depth++
	return nil
}

func (p *Parser) leaveNesting() { p.depth-- }
