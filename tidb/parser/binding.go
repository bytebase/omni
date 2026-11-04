package parser

import (
	nodes "github.com/bytebase/omni/tidb/ast"
)

// TiDB's SQL plan bindings: CREATE BINDING, DROP BINDING, and SET BINDING,
// each of which carries one or two complete statements.
//
// Ref: TiDB's SQL Plan Management,
// https://docs.pingcap.com/tidb/v8.5/sql-plan-management/. Every form and
// every rejection below is measured against v8.5.5 in
// binding_oracle_test.go.

// parseBindingScope reads the GLOBAL or SESSION a binding statement may
// carry, leaving any other token where it is.
func (p *Parser) parseBindingScope() nodes.BindingScope {
	switch p.cur.Type {
	case kwGLOBAL:
		p.advance()
		return nodes.BindingScopeGlobal
	case kwSESSION:
		p.advance()
		return nodes.BindingScopeSession
	}
	return nodes.BindingScopeUnset
}

// parseCreateBindingStmt parses CREATE [GLOBAL | SESSION] BINDING. CREATE
// and the scope word are consumed; p.cur is BINDING.
//
//	CREATE [GLOBAL | SESSION] BINDING FOR stmt USING stmt
//	CREATE [GLOBAL | SESSION] BINDING USING stmt
//	CREATE [GLOBAL | SESSION] BINDING FROM HISTORY USING PLAN DIGEST digest [, digest] ...
func (p *Parser) parseCreateBindingStmt(start int, scope nodes.BindingScope) (*nodes.CreateBindingStmt, error) {
	p.advance() // consume BINDING

	stmt := &nodes.CreateBindingStmt{Loc: nodes.Loc{Start: start}, Scope: scope}

	p.checkCursor()
	if p.collectMode() {
		for _, t := range []int{kwFOR, kwUSING, kwFROM} {
			p.addTokenCandidate(t)
		}
		return nil, &ParseError{Message: "collecting"}
	}

	switch p.cur.Type {
	case kwFROM:
		// CREATE BINDING FROM HISTORY USING PLAN DIGEST ... binds plans
		// already in the statement history, named by digest.
		p.advance() // consume FROM
		if _, err := p.expect(kwHISTORY); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwUSING); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwPLAN); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwDIGEST); err != nil {
			return nil, err
		}
		digests, err := p.parseBindingDigestList()
		if err != nil {
			return nil, err
		}
		stmt.PlanDigests = digests

	case kwUSING:
		// CREATE BINDING USING stmt, where the statement carries its own
		// hints and the shape it binds is read off the same text.
		p.advance() // consume USING
		hinted, err := p.parseBindableStmt()
		if err != nil {
			return nil, err
		}
		stmt.HintedStmt = hinted

	case kwFOR:
		p.advance() // consume FOR
		origin, err := p.parseBindableStmt()
		if err != nil {
			return nil, err
		}
		stmt.OriginStmt = origin
		if _, err := p.expect(kwUSING); err != nil {
			return nil, err
		}
		hinted, err := p.parseBindableStmt()
		if err != nil {
			return nil, err
		}
		stmt.HintedStmt = hinted

	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseDropBindingStmt parses DROP [GLOBAL | SESSION] BINDING. DROP and the
// scope word are consumed; p.cur is BINDING.
//
//	DROP [GLOBAL | SESSION] BINDING FOR stmt [USING stmt]
//	DROP [GLOBAL | SESSION] BINDING FOR SQL DIGEST digest [, digest] ...
func (p *Parser) parseDropBindingStmt(start int, scope nodes.BindingScope) (*nodes.DropBindingStmt, error) {
	p.advance() // consume BINDING

	stmt := &nodes.DropBindingStmt{Loc: nodes.Loc{Start: start}, Scope: scope}

	if _, err := p.expect(kwFOR); err != nil {
		return nil, err
	}

	if p.cur.Type == kwSQL {
		p.advance() // consume SQL
		if _, err := p.expect(kwDIGEST); err != nil {
			return nil, err
		}
		digests, err := p.parseBindingDigestList()
		if err != nil {
			return nil, err
		}
		stmt.SQLDigests = digests
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

	origin, err := p.parseBindableStmt()
	if err != nil {
		return nil, err
	}
	stmt.OriginStmt = origin

	// USING is optional here: DROP BINDING FOR stmt drops every binding of
	// the shape, and DROP BINDING FOR stmt USING stmt drops the one plan.
	if _, ok := p.match(kwUSING); ok {
		hinted, err := p.parseBindableStmt()
		if err != nil {
			return nil, err
		}
		stmt.HintedStmt = hinted
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseSetBindingStmt parses SET BINDING {ENABLED | DISABLED} FOR ... . SET
// is consumed; p.cur is BINDING and the token after it is ENABLED or
// DISABLED, which is what tells this statement from the variable assignment
// `SET binding = 1`, which v8.5.5 also accepts.
//
//	SET BINDING {ENABLED | DISABLED} FOR stmt [USING stmt]
//	SET BINDING {ENABLED | DISABLED} FOR SQL DIGEST digest
func (p *Parser) parseSetBindingStmt(start int) (*nodes.SetBindingStmt, error) {
	p.advance() // consume BINDING

	stmt := &nodes.SetBindingStmt{Loc: nodes.Loc{Start: start}}
	switch p.cur.Type {
	case kwENABLED:
		stmt.Enabled = true
	case kwDISABLED:
	default:
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume ENABLED or DISABLED

	if _, err := p.expect(kwFOR); err != nil {
		return nil, err
	}

	if p.cur.Type == kwSQL {
		p.advance() // consume SQL
		if _, err := p.expect(kwDIGEST); err != nil {
			return nil, err
		}
		// One digest, and a string literal only: v8.5.5 answers 1064 for
		// `SET BINDING ENABLED FOR SQL DIGEST 'a', 'b'` and for `... SQL
		// DIGEST @a`, both of which DROP BINDING FOR SQL DIGEST takes.
		if p.cur.Type != tokSCONST {
			return nil, p.syntaxErrorAtCur()
		}
		digest, err := p.parseBindingDigest()
		if err != nil {
			return nil, err
		}
		stmt.SQLDigest = digest
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

	origin, err := p.parseBindableStmt()
	if err != nil {
		return nil, err
	}
	stmt.OriginStmt = origin

	if _, ok := p.match(kwUSING); ok {
		hinted, err := p.parseBindableStmt()
		if err != nil {
			return nil, err
		}
		stmt.HintedStmt = hinted
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseBindingDigestList parses the comma-separated digest list of CREATE
// BINDING FROM HISTORY USING PLAN DIGEST and DROP BINDING FOR SQL DIGEST.
// It holds at least one digest; a trailing comma is a syntax error.
func (p *Parser) parseBindingDigestList() ([]nodes.ExprNode, error) {
	var digests []nodes.ExprNode
	for {
		digest, err := p.parseBindingDigest()
		if err != nil {
			return nil, err
		}
		digests = append(digests, digest)
		if p.cur.Type != ',' {
			return digests, nil
		}
		p.advance() // consume ','
	}
}

// parseBindingDigest parses one digest: a string literal or a user
// variable. Nothing else is one -- v8.5.5 answers 1064 for a number, a
// bare name, a placeholder, a system variable, and a function call.
//
// A user variable goes through parseVariableRef, so a digest reads the
// same variables every other position does. The lexer scans one as a
// single tokIDENT whose Str keeps the '@'; a system variable is tokAt2,
// which the test below leaves to the error. v8.5.5 also takes a quoted
// user variable here (`@'a'`, `@"a"`, and a bare `@`), which this parser
// reads nowhere -- corpus line 688, `SET @' ' = 1`, is the same gap.
func (p *Parser) parseBindingDigest() (nodes.ExprNode, error) {
	switch {
	case p.cur.Type == tokSCONST:
		tok := p.advance()
		return &nodes.StringLit{Loc: nodes.Loc{Start: tok.Loc, End: tok.End}, Value: tok.Str}, nil
	case p.cur.Type == tokIDENT && len(p.cur.Str) > 1 && p.cur.Str[0] == '@':
		ref, err := p.parseVariableRef()
		if err != nil {
			return nil, err
		}
		return ref, nil
	}
	return nil, p.syntaxErrorAtCur()
}

// parseBindableStmt parses the statement in a binding's FOR or USING
// position. TiDB's bindable statements are a query -- a SELECT, a set
// operation, a parenthesized query, TABLE, or VALUES, with or without a
// WITH clause -- an UPDATE, with or without a WITH clause, an INSERT or a
// REPLACE, and the DELETE forms that take no USING clause. Nothing else is
// one: SHOW, EXPLAIN, SET, LOAD DATA, ANALYZE, and every DDL statement
// answer 1064 here on v8.5.5.
func (p *Parser) parseBindableStmt() (nodes.StmtNode, error) {
	p.checkCursor()
	if p.collectMode() {
		for _, t := range []int{kwSELECT, kwWITH, kwTABLE, kwVALUES, kwINSERT, kwREPLACE, kwUPDATE, kwDELETE} {
			p.addTokenCandidate(t)
		}
		return nil, &ParseError{Message: "collecting"}
	}

	switch p.cur.Type {
	case kwWITH:
		return p.parseBindableWithStmt()
	case kwSELECT, '(':
		stmt, err := p.parseSelectStmt()
		if err != nil {
			return nil, err
		}
		return stmt, nil
	case kwTABLE:
		stmt, err := p.parseTableStmt()
		if err != nil {
			return nil, err
		}
		return stmt, nil
	case kwVALUES:
		stmt, err := p.parseValuesStmt()
		if err != nil {
			return nil, err
		}
		return stmt, nil
	case kwINSERT:
		stmt, err := p.parseInsertStmt()
		if err != nil {
			return nil, err
		}
		return stmt, nil
	case kwREPLACE:
		stmt, err := p.parseReplaceStmt()
		if err != nil {
			return nil, err
		}
		return stmt, nil
	case kwUPDATE:
		stmt, err := p.parseUpdateStmt()
		if err != nil {
			return nil, err
		}
		return stmt, nil
	case kwDELETE:
		stmt, err := p.parseDeleteStmtNoUsing()
		if err != nil {
			return nil, err
		}
		return stmt, nil
	default:
		return nil, p.syntaxErrorAtCur()
	}
}

// parseBindableWithStmt parses a bindable statement that opens with a WITH
// clause. Only a query and an UPDATE take one here: v8.5.5 answers 1064 for
// `CREATE BINDING USING WITH c AS (SELECT 1) DELETE FROM t WHERE a = 1`
// though it accepts that DELETE on its own, and it has no WITH before
// INSERT or REPLACE at all.
//
// TABLE and VALUES after a WITH clause are left out with the rest of their
// gap: v8.5.5 accepts `WITH cte AS (TABLE t) TABLE cte`, which this parser
// rejects everywhere, binding or not, and closing that belongs to the WITH
// family rather than here.
func (p *Parser) parseBindableWithStmt() (nodes.StmtNode, error) {
	start := p.pos()
	ctes, err := p.parseWithClause()
	if err != nil {
		return nil, err
	}

	switch p.cur.Type {
	case kwSELECT, '(':
		stmt, err := p.parseSelectClauseAndTrailing(start, ctes)
		if err != nil {
			return nil, err
		}
		return stmt, nil
	case kwUPDATE:
		// As in parseWithStmt, an UPDATE node keeps no CTE list; the WITH
		// clause shows only in the statement's start offset.
		stmt, err := p.parseUpdateStmt()
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = start
		return stmt, nil
	default:
		return nil, p.syntaxErrorAtCur()
	}
}
