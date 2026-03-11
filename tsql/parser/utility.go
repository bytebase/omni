// Package parser - utility.go implements T-SQL utility statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/tsql/ast"
)

// parseUseStmt parses a USE database statement.
//
//	USE database
func (p *Parser) parseUseStmt() *nodes.UseStmt {
	loc := p.pos()
	p.advance() // consume USE

	stmt := &nodes.UseStmt{
		Loc: nodes.Loc{Start: loc},
	}

	if p.isIdentLike() {
		stmt.Database = p.cur.Str
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parsePrintStmt parses a PRINT statement.
//
//	PRINT expr
func (p *Parser) parsePrintStmt() *nodes.PrintStmt {
	loc := p.pos()
	p.advance() // consume PRINT

	stmt := &nodes.PrintStmt{
		Loc: nodes.Loc{Start: loc},
	}

	stmt.Expr = p.parseExpr()

	stmt.Loc.End = p.pos()
	return stmt
}

// parseRaiseErrorStmt parses a RAISERROR statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/raiserror-transact-sql
//
//	RAISERROR (msg, severity, state [, args]) [WITH options]
func (p *Parser) parseRaiseErrorStmt() *nodes.RaiseErrorStmt {
	loc := p.pos()
	p.advance() // consume RAISERROR

	stmt := &nodes.RaiseErrorStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// RAISERROR can use parens or not: RAISERROR('msg', 16, 1)
	if _, err := p.expect('('); err == nil {
		// Message
		stmt.Message = p.parseExpr()

		// Severity
		if _, ok := p.match(','); ok {
			stmt.Severity = p.parseExpr()
		}

		// State
		if _, ok := p.match(','); ok {
			stmt.State = p.parseExpr()
		}

		// Optional args
		var args []nodes.Node
		for {
			if _, ok := p.match(','); !ok {
				break
			}
			arg := p.parseExpr()
			args = append(args, arg)
		}
		if len(args) > 0 {
			stmt.Args = &nodes.List{Items: args}
		}

		_, _ = p.expect(')')
	}

	// WITH options (LOG, NOWAIT, SETERROR)
	if p.cur.Type == kwWITH {
		p.advance()
		var opts []nodes.Node
		for {
			if p.isIdentLike() || p.cur.Type == kwNOWAIT {
				opts = append(opts, &nodes.String{Str: strings.ToUpper(p.cur.Str)})
				p.advance()
			} else {
				break
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
		if len(opts) > 0 {
			stmt.Options = &nodes.List{Items: opts}
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseThrowStmt parses a THROW statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/throw-transact-sql
//
//	THROW [number, message, state]
func (p *Parser) parseThrowStmt() *nodes.ThrowStmt {
	loc := p.pos()
	p.advance() // consume THROW

	stmt := &nodes.ThrowStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// THROW without arguments = rethrow
	if p.cur.Type == ';' || p.cur.Type == tokEOF || p.cur.Type == kwEND ||
		p.isStatementStart() {
		stmt.Loc.End = p.pos()
		return stmt
	}

	// Error number
	stmt.ErrorNumber = p.parseExpr()

	// Message
	if _, ok := p.match(','); ok {
		stmt.Message = p.parseExpr()
	}

	// State
	if _, ok := p.match(','); ok {
		stmt.State = p.parseExpr()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseTruncateStmt parses a TRUNCATE TABLE statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/truncate-table-transact-sql
//
//	TRUNCATE TABLE name
func (p *Parser) parseTruncateStmt() *nodes.TruncateStmt {
	loc := p.pos()
	p.advance() // consume TRUNCATE

	// TABLE
	p.match(kwTABLE)

	stmt := &nodes.TruncateStmt{
		Loc: nodes.Loc{Start: loc},
	}

	stmt.Table = p.parseTableRef()

	stmt.Loc.End = p.pos()
	return stmt
}
