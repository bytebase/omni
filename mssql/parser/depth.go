package parser

// maxNestingDepth bounds how deeply the parser may descend into nested
// constructs: statements inside IF, WHILE, BEGIN...END and TRY...CATCH,
// expressions, queries, and the other productions that recurse. Without a
// bound, a few hundred kilobytes of '(' overflowed the Go stack, and a stack
// overflow kills the process instead of returning an error. SQL Server 2022
// itself rejects 1,016 nested parentheses, 505 nested BEGIN...END blocks and
// 169 nested scalar subqueries with Msg 191 ("Some part of your SQL statement
// is nested too deeply"), so nothing it accepts reaches this limit.
const maxNestingDepth = 10000

// enterNesting is called on entry to the productions that recurse, placed so
// that every recursion cycle in the parser passes through one. It fails when
// the statement nests deeper than maxNestingDepth; otherwise the caller
// defers leaveNesting.
//
// The error leaves Source unset so that Error() reports the message rather
// than the generic "syntax error at or near" text. ParseError carries no
// error number, so SQL Server's Msg 191 is not recorded on it.
func (p *Parser) enterNesting() error {
	if p.depth >= maxNestingDepth {
		return &ParseError{Message: "statement is nested too deeply", Position: p.cur.Loc}
	}
	p.depth++
	return nil
}

func (p *Parser) leaveNesting() { p.depth-- }
