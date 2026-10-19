package parser

// maxNestingDepth bounds how deeply the parser may descend into nested
// constructs: expressions, queries, FROM sources, types, option groups, and
// the statements that embed other statements. Without a bound, a few hundred
// kilobytes of '(' overflowed the Go stack, and a stack overflow kills the
// process instead of returning an error. Snowflake Scripting blocks keep
// their own, much lower ceiling (maxScriptDepth in scripting.go).
const maxNestingDepth = 10000

// enterNesting is called on entry to the productions that recurse, placed so
// that every recursion cycle in the parser passes through one (the scripting
// statement cycles pass through enterScript instead). It fails when the
// statement nests deeper than maxNestingDepth; otherwise the caller defers
// leaveNesting.
func (p *Parser) enterNesting() error {
	if p.depth >= maxNestingDepth {
		return &ParseError{Position: p.cur.Loc.Start, End: p.cur.Loc.End, Message: "statement is nested too deeply"}
	}
	p.depth++
	return nil
}

func (p *Parser) leaveNesting() { p.depth-- }
