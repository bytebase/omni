// Package parser - control_flow.go implements T-SQL control flow statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseIfStmt parses an IF...ELSE statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/if-else-transact-sql
//
//	IF condition stmt [ELSE stmt]
func (p *Parser) parseIfStmt() *nodes.IfStmt {
	loc := p.pos()
	p.advance() // consume IF

	stmt := &nodes.IfStmt{
		Loc: nodes.Loc{Start: loc},
	}

	stmt.Condition = p.parseExpr()
	stmt.Then = p.parseStmt()

	// ELSE
	if _, ok := p.match(kwELSE); ok {
		stmt.Else = p.parseStmt()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseWhileStmt parses a WHILE loop.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/while-transact-sql
//
//	WHILE condition stmt
func (p *Parser) parseWhileStmt() *nodes.WhileStmt {
	loc := p.pos()
	p.advance() // consume WHILE

	stmt := &nodes.WhileStmt{
		Loc: nodes.Loc{Start: loc},
	}

	stmt.Condition = p.parseExpr()
	stmt.Body = p.parseStmt()

	stmt.Loc.End = p.pos()
	return stmt
}

// parseBeginStmt parses BEGIN...END, BEGIN TRY...END TRY BEGIN CATCH...END CATCH,
// or BEGIN TRAN/TRANSACTION.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/begin-end-transact-sql
func (p *Parser) parseBeginStmt() nodes.StmtNode {
	loc := p.pos()

	// Check for BEGIN TRY
	next := p.peekNext()
	if next.Type == kwTRY {
		return p.parseTryCatchStmt()
	}

	// Check for BEGIN DISTRIBUTED TRAN/TRANSACTION
	if next.Type == kwDISTRIBUTED {
		return p.parseBeginDistributedTransStmt()
	}

	// Check for BEGIN TRAN/TRANSACTION
	if next.Type == kwTRAN || next.Type == kwTRANSACTION {
		return p.parseBeginTransStmt()
	}

	p.advance() // consume BEGIN

	// Parse statements until END
	var stmts []nodes.Node
	for p.cur.Type != kwEND && p.cur.Type != tokEOF {
		if p.cur.Type == ';' {
			p.advance()
			continue
		}
		s := p.parseStmt()
		if s == nil {
			break
		}
		stmts = append(stmts, s)
	}
	p.match(kwEND)

	return &nodes.BeginEndStmt{
		Stmts: &nodes.List{Items: stmts},
		Loc:   nodes.Loc{Start: loc, End: p.pos()},
	}
}

// parseTryCatchStmt parses BEGIN TRY...END TRY BEGIN CATCH...END CATCH.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/try-catch-transact-sql
func (p *Parser) parseTryCatchStmt() *nodes.TryCatchStmt {
	loc := p.pos()
	p.advance() // consume BEGIN
	p.advance() // consume TRY

	stmt := &nodes.TryCatchStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// TRY block
	var tryStmts []nodes.Node
	for p.cur.Type != kwEND && p.cur.Type != tokEOF {
		if p.cur.Type == ';' {
			p.advance()
			continue
		}
		s := p.parseStmt()
		if s == nil {
			break
		}
		tryStmts = append(tryStmts, s)
	}
	stmt.TryBlock = &nodes.List{Items: tryStmts}

	// END TRY
	p.match(kwEND)
	p.match(kwTRY)

	// BEGIN CATCH
	p.match(kwBEGIN)
	p.match(kwCATCH)

	// CATCH block
	var catchStmts []nodes.Node
	for p.cur.Type != kwEND && p.cur.Type != tokEOF {
		if p.cur.Type == ';' {
			p.advance()
			continue
		}
		s := p.parseStmt()
		if s == nil {
			break
		}
		catchStmts = append(catchStmts, s)
	}
	stmt.CatchBlock = &nodes.List{Items: catchStmts}

	// END CATCH
	p.match(kwEND)
	p.match(kwCATCH)

	stmt.Loc.End = p.pos()
	return stmt
}

// parseReturnStmt parses a RETURN statement.
//
//	RETURN [expr]
func (p *Parser) parseReturnStmt() *nodes.ReturnStmt {
	loc := p.pos()
	p.advance() // consume RETURN

	stmt := &nodes.ReturnStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Optional return expression
	if p.cur.Type != ';' && p.cur.Type != tokEOF && p.cur.Type != kwEND &&
		p.cur.Type != kwELSE && p.cur.Type != kwGO {
		stmt.Value = p.parseExpr()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseBreakStmt parses a BREAK statement.
func (p *Parser) parseBreakStmt() *nodes.BreakStmt {
	loc := p.pos()
	p.advance() // consume BREAK
	return &nodes.BreakStmt{
		Loc: nodes.Loc{Start: loc, End: p.pos()},
	}
}

// parseContinueStmt parses a CONTINUE statement.
func (p *Parser) parseContinueStmt() *nodes.ContinueStmt {
	loc := p.pos()
	p.advance() // consume CONTINUE
	return &nodes.ContinueStmt{
		Loc: nodes.Loc{Start: loc, End: p.pos()},
	}
}

// parseGotoStmt parses a GOTO label statement.
//
//	GOTO label
func (p *Parser) parseGotoStmt() *nodes.GotoStmt {
	loc := p.pos()
	p.advance() // consume GOTO

	label := ""
	if p.isIdentLike() {
		label = p.cur.Str
		p.advance()
	}

	return &nodes.GotoStmt{
		Label: label,
		Loc:   nodes.Loc{Start: loc, End: p.pos()},
	}
}

// parseWaitForStmt parses a WAITFOR statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/waitfor-transact-sql
//
//	WAITFOR DELAY|TIME expr
func (p *Parser) parseWaitForStmt() *nodes.WaitForStmt {
	loc := p.pos()
	p.advance() // consume WAITFOR

	stmt := &nodes.WaitForStmt{
		Loc: nodes.Loc{Start: loc},
	}

	if p.cur.Type == kwDELAY {
		stmt.WaitType = "DELAY"
		p.advance()
	} else if p.cur.Type == kwTIME || (p.cur.Type == tokIDENT && strings.EqualFold(p.cur.Str, "time")) {
		stmt.WaitType = "TIME"
		p.advance()
	}

	stmt.Value = p.parseExpr()

	stmt.Loc.End = p.pos()
	return stmt
}
