package parser

import "testing"

// TestParenQueryRecordDepth: parseSelectWithParens records a failed parse by
// the offset of its '(' so that speculation does not repeat it. A parse
// started deeper can reach the nesting limit where the recorded one did not,
// so a record is replayed only at its own depth or shallower; and a nesting
// error depends on the depth the parse started at, so it is never recorded.
// Both held only outside completion mode before, and neither inside it.
func TestParenQueryRecordDepth(t *testing.T) {
	// The inner query parses, then `= 0` fails the outer one: an ordinary
	// error from depth 0, a nesting error from one level below the limit.
	const sql = "((SELECT 1) = 0)"
	for _, completing := range []bool{false, true} {
		newParser := func() *Parser {
			p := &Parser{lexer: NewLexer(sql)}
			if completing {
				p.completing = true
				p.cursorOff = len(sql) + 1
				p.candidates = newCandidateSet()
				p.maxCollect = 100
			}
			p.advance()
			return p
		}
		parseAt := func(p *Parser, depth int) error {
			start := p.saveScan()
			p.depth = depth
			_, err := p.parseSelectWithParens()
			p.restoreScan(start)
			p.depth = 0
			return err
		}
		deep := maxNestingDepth - 1

		p := newParser()
		if err := parseAt(p, 0); err == nil || isNestingError(err) {
			t.Fatalf("completing=%v: from depth 0: err = %v, want a syntax error", completing, err)
		}
		if err := parseAt(p, deep); !isNestingError(err) {
			t.Errorf("completing=%v: after a failure recorded at depth 0, from depth %d: err = %v, want the nesting error", completing, deep, err)
		}

		p = newParser()
		if err := parseAt(p, deep); !isNestingError(err) {
			t.Fatalf("completing=%v: from depth %d: err = %v, want the nesting error", completing, deep, err)
		}
		if err := parseAt(p, 0); err == nil || isNestingError(err) {
			t.Errorf("completing=%v: after the nesting error, from depth 0: err = %v, want a syntax error", completing, err)
		}
	}
}
