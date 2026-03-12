// Package parser - transaction.go implements T-SQL transaction statement parsing.
package parser

import (
	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseBeginTransStmt parses a BEGIN TRAN/TRANSACTION statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/begin-transaction-transact-sql
//
//	BEGIN TRAN[SACTION] [name]
func (p *Parser) parseBeginTransStmt() *nodes.BeginTransStmt {
	loc := p.pos()
	p.advance() // consume BEGIN

	// TRAN or TRANSACTION
	p.match(kwTRAN)
	p.match(kwTRANSACTION)

	stmt := &nodes.BeginTransStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Optional transaction name
	if p.isIdentLike() || p.cur.Type == tokVARIABLE {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseBeginDistributedTransStmt parses a BEGIN DISTRIBUTED TRAN[SACTION] statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/begin-distributed-transaction-transact-sql
//
//	BEGIN DISTRIBUTED TRAN[SACTION] [tran_name | @tran_name_variable]
func (p *Parser) parseBeginDistributedTransStmt() *nodes.BeginDistributedTransStmt {
	loc := p.pos()
	p.advance() // consume BEGIN
	p.advance() // consume DISTRIBUTED

	// TRAN or TRANSACTION
	p.match(kwTRAN)
	p.match(kwTRANSACTION)

	stmt := &nodes.BeginDistributedTransStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Optional transaction name
	if p.isIdentLike() || p.cur.Type == tokVARIABLE {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseCommitStmt parses a COMMIT [TRAN/TRANSACTION] statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/commit-transaction-transact-sql
//
//	COMMIT [TRAN[SACTION]] [name]
func (p *Parser) parseCommitStmt() *nodes.CommitTransStmt {
	loc := p.pos()
	p.advance() // consume COMMIT

	// Optional TRAN/TRANSACTION
	p.match(kwTRAN)
	p.match(kwTRANSACTION)

	stmt := &nodes.CommitTransStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Optional name
	if p.isIdentLike() || p.cur.Type == tokVARIABLE {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseRollbackStmt parses a ROLLBACK [TRAN/TRANSACTION] statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/rollback-transaction-transact-sql
//
//	ROLLBACK [TRAN[SACTION]] [name|savepoint]
func (p *Parser) parseRollbackStmt() *nodes.RollbackTransStmt {
	loc := p.pos()
	p.advance() // consume ROLLBACK

	// Optional TRAN/TRANSACTION
	p.match(kwTRAN)
	p.match(kwTRANSACTION)

	stmt := &nodes.RollbackTransStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Optional name/savepoint
	if p.isIdentLike() || p.cur.Type == tokVARIABLE {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseSaveTransStmt parses a SAVE TRAN/TRANSACTION statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/save-transaction-transact-sql
//
//	SAVE TRAN[SACTION] name
func (p *Parser) parseSaveTransStmt() *nodes.SaveTransStmt {
	loc := p.pos()
	p.advance() // consume SAVE

	// TRAN or TRANSACTION
	p.match(kwTRAN)
	p.match(kwTRANSACTION)

	stmt := &nodes.SaveTransStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Name
	if p.isIdentLike() || p.cur.Type == tokVARIABLE {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt
}
