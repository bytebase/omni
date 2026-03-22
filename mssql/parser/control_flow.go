// Package parser - control_flow.go implements T-SQL control flow statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseIfStmt parses an IF...ELSE statement.
//
// BNF: mssql/parser/bnf/if-else-transact-sql.bnf
//
//	IF boolean_expression
//	    { sql_statement | statement_block }
//	[ ELSE
//	    { sql_statement | statement_block } ]
func (p *Parser) parseIfStmt() (*nodes.IfStmt, error) {
	loc := p.pos()
	p.advance() // consume IF

	stmt := &nodes.IfStmt{
		Loc: nodes.Loc{Start: loc},
	}

	cond, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if cond == nil {
		return nil, &ParseError{
			Message:  "expected condition after IF",
			Position: p.cur.Loc,
		}
	}
	stmt.Condition = cond
	then, err := p.parseStmt()
	if err != nil {
		return nil, err
	}
	if then == nil {
		return nil, &ParseError{
			Message:  "expected statement after IF condition",
			Position: p.cur.Loc,
		}
	}
	stmt.Then = then

	// ELSE
	if _, ok := p.match(kwELSE); ok {
		elseStmt, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		stmt.Else = elseStmt
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseWhileStmt parses a WHILE loop.
//
// BNF: mssql/parser/bnf/while-transact-sql.bnf
//
//	WHILE boolean_expression
//	    { sql_statement | statement_block | BREAK | CONTINUE }
func (p *Parser) parseWhileStmt() (*nodes.WhileStmt, error) {
	loc := p.pos()
	p.advance() // consume WHILE

	stmt := &nodes.WhileStmt{
		Loc: nodes.Loc{Start: loc},
	}

	cond, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if cond == nil {
		return nil, &ParseError{
			Message:  "expected condition after WHILE",
			Position: p.cur.Loc,
		}
	}
	stmt.Condition = cond
	body, err := p.parseStmt()
	if err != nil {
		return nil, err
	}
	if body == nil {
		return nil, &ParseError{
			Message:  "expected statement after WHILE condition",
			Position: p.cur.Loc,
		}
	}
	stmt.Body = body

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseBeginStmt parses BEGIN...END, BEGIN TRY...END TRY BEGIN CATCH...END CATCH,
// or BEGIN TRAN/TRANSACTION.
//
// BNF: mssql/parser/bnf/begin-end-transact-sql.bnf
//
//	BEGIN [ ; ]
//	    { sql_statement | statement_block }
//	END [ ; ]
func (p *Parser) parseBeginStmt() (nodes.StmtNode, error) {
	loc := p.pos()

	// Check for BEGIN TRY
	next := p.peekNext()
	if next.Type == kwTRY {
		stmt, err := p.parseTryCatchStmt()
		return stmt, err
	}

	// Check for BEGIN DISTRIBUTED TRAN/TRANSACTION
	if next.Type == kwDISTRIBUTED {
		stmt, err := p.parseBeginDistributedTransStmt()
		return stmt, err
	}

	// Check for BEGIN TRAN/TRANSACTION
	if next.Type == kwTRAN || next.Type == kwTRANSACTION {
		stmt, err := p.parseBeginTransStmt()
		return stmt, err
	}

	// Check for BEGIN CONVERSATION TIMER (service broker)
	if next.Str != "" && matchesKeywordCI(next.Str, "CONVERSATION") {
		p.advance() // consume BEGIN
		p.advance() // consume CONVERSATION
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "TIMER") {
			p.advance() // consume TIMER
			stmt, err := p.parseBeginConversationTimerStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc = nodes.Loc{Start: loc}
			stmt.Loc.End = p.pos()
			return stmt, nil
		}
		// Not a valid statement; fall through to BEGIN...END block parsing
		// (BEGIN already consumed above, skip the p.advance() below)
		goto parseBlock
	}

	// Check for BEGIN DIALOG [CONVERSATION] (service broker)
	if next.Str != "" && matchesKeywordCI(next.Str, "DIALOG") {
		p.advance() // consume BEGIN
		p.advance() // consume DIALOG
		// optionally consume CONVERSATION
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "CONVERSATION") {
			p.advance()
		}
		stmt, err := p.parseBeginConversationStmt()
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = loc
		return stmt, nil
	}

	p.advance() // consume BEGIN

parseBlock:
	// Parse statements until END
	var stmts []nodes.Node
	for p.cur.Type != kwEND && p.cur.Type != tokEOF {
		if p.cur.Type == ';' {
			p.advance()
			continue
		}
		s, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		if s == nil {
			break
		}
		stmts = append(stmts, s)
	}
	p.match(kwEND)

	return &nodes.BeginEndStmt{
		Stmts: &nodes.List{Items: stmts},
		Loc:   nodes.Loc{Start: loc, End: p.pos()},
	}, nil
}

// parseTryCatchStmt parses BEGIN TRY...END TRY BEGIN CATCH...END CATCH.
//
// BNF: mssql/parser/bnf/try-catch-transact-sql.bnf
//
//	BEGIN TRY
//	    { sql_statement | statement_block }
//	END TRY
//	BEGIN CATCH
//	    [ { sql_statement | statement_block } ]
//	END CATCH
//	[ ; ]
func (p *Parser) parseTryCatchStmt() (*nodes.TryCatchStmt, error) {
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
		s, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
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
		s, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
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
	return stmt, nil
}

// parseReturnStmt parses a RETURN statement.
//
// BNF: mssql/parser/bnf/return-transact-sql.bnf
//
//	RETURN [ integer_expression ]
func (p *Parser) parseReturnStmt() (*nodes.ReturnStmt, error) {
	loc := p.pos()
	p.advance() // consume RETURN

	stmt := &nodes.ReturnStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Optional return expression
	if p.cur.Type != ';' && p.cur.Type != tokEOF && p.cur.Type != kwEND &&
		p.cur.Type != kwELSE && p.cur.Type != kwGO {
		stmt.Value, _ = p.parseExpr()
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseBreakStmt parses a BREAK statement.
//
// BNF: mssql/parser/bnf/break-transact-sql.bnf
//
//	BREAK
func (p *Parser) parseBreakStmt() (*nodes.BreakStmt, error) {
	loc := p.pos()
	p.advance() // consume BREAK
	return &nodes.BreakStmt{
		Loc: nodes.Loc{Start: loc, End: p.pos()},
	}, nil
}

// parseContinueStmt parses a CONTINUE statement.
//
// BNF: mssql/parser/bnf/continue-transact-sql.bnf
//
//	CONTINUE
func (p *Parser) parseContinueStmt() (*nodes.ContinueStmt, error) {
	loc := p.pos()
	p.advance() // consume CONTINUE
	return &nodes.ContinueStmt{
		Loc: nodes.Loc{Start: loc, End: p.pos()},
	}, nil
}

// parseGotoStmt parses a GOTO label statement.
//
// BNF: mssql/parser/bnf/goto-transact-sql.bnf
//
//	GOTO label
func (p *Parser) parseGotoStmt() (*nodes.GotoStmt, error) {
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
	}, nil
}

// parseWaitForStmt parses a WAITFOR statement.
//
// BNF: mssql/parser/bnf/waitfor-transact-sql.bnf
//
//	WAITFOR
//	{
//	    DELAY 'time_to_pass'
//	  | TIME 'time_to_execute'
//	  | [ ( receive_statement ) | ( get_conversation_group_statement ) ]
//	    [ , TIMEOUT timeout ]
//	}
func (p *Parser) parseWaitForStmt() (*nodes.WaitForStmt, error) {
	loc := p.pos()
	p.advance() // consume WAITFOR

	stmt := &nodes.WaitForStmt{
		Loc: nodes.Loc{Start: loc},
	}

	if p.cur.Type == kwDELAY {
		stmt.WaitType = "DELAY"
		p.advance()
		stmt.Value, _ = p.parseExpr()
	} else if p.cur.Type == kwTIME || (p.cur.Type == tokIDENT && strings.EqualFold(p.cur.Str, "time")) {
		stmt.WaitType = "TIME"
		p.advance()
		stmt.Value, _ = p.parseExpr()
	} else if p.cur.Type == '(' {
		// Parenthesized form: WAITFOR ( receive_statement | get_conversation_group_statement )
		p.advance() // consume '('
		innerStmt, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		stmt.InnerStmt = innerStmt
		_, _ = p.expect(')')
		// Optional: , TIMEOUT timeout
		if _, ok := p.match(','); ok {
			if p.isIdentLike() && strings.EqualFold(p.cur.Str, "TIMEOUT") {
				p.advance()
				stmt.Timeout, _ = p.parseExpr()
			}
		}
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}
