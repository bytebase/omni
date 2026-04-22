// Package parser - transaction.go implements T-SQL transaction statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseBeginTransStmt parses a BEGIN TRAN/TRANSACTION statement.
//
// BNF: mssql/parser/bnf/begin-transaction-transact-sql.bnf
//
//	BEGIN { TRAN | TRANSACTION }
//	    [ { transaction_name | @tran_name_variable }
//	      [ WITH MARK [ 'description' ] ]
//	    ]
//	[ ; ]
func (p *Parser) parseBeginTransStmt() (*nodes.BeginTransStmt, error) {
	loc := p.pos()
	p.advance() // consume BEGIN

	// TRAN or TRANSACTION
	p.match(kwTRAN)
	p.match(kwTRANSACTION)

	// Completion: after BEGIN TRANSACTION → identifier (optional name)
	if p.collectMode() {
		p.addRuleCandidate("identifier")
		return nil, errCollecting
	}

	stmt := &nodes.BeginTransStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Optional transaction name
	if p.isIdentLike() || p.cur.Type == tokVARIABLE {
		stmt.Name = p.cur.Str
		p.advance()

		// Optional WITH MARK ['description']
		if _, ok := p.match(kwWITH); ok {
			if p.matchIdentCI("MARK") {
				stmt.WithMark = true
				// Optional description string
				if p.cur.Type == tokSCONST {
					stmt.MarkDescription = p.cur.Str
					p.advance()
				}
			}
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseBeginDistributedTransStmt parses a BEGIN DISTRIBUTED TRAN[SACTION] statement.
//
// BNF: mssql/parser/bnf/begin-transaction-transact-sql.bnf
//
//	BEGIN DISTRIBUTED { TRAN | TRANSACTION }
//	    [ transaction_name | @tran_name_variable ]
//	[ ; ]
func (p *Parser) parseBeginDistributedTransStmt() (*nodes.BeginDistributedTransStmt, error) {
	loc := p.pos()
	p.advance() // consume BEGIN
	p.advance() // consume DISTRIBUTED

	// TRAN or TRANSACTION
	p.match(kwTRAN)
	p.match(kwTRANSACTION)

	stmt := &nodes.BeginDistributedTransStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Optional transaction name
	if p.isIdentLike() || p.cur.Type == tokVARIABLE {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseCommitStmt parses a COMMIT [TRAN/TRANSACTION] statement.
//
// BNF: mssql/parser/bnf/commit-transaction-transact-sql.bnf
//
//	COMMIT [ { TRAN | TRANSACTION }
//	    [ transaction_name | @tran_name_variable ] ]
//	    [ WITH ( DELAYED_DURABILITY = { OFF | ON } ) ]
//	[ ; ]
func (p *Parser) parseCommitStmt() (*nodes.CommitTransStmt, error) {
	loc := p.pos()
	p.advance() // consume COMMIT

	// Completion: after COMMIT → TRANSACTION, WORK
	if p.collectMode() {
		p.addTokenCandidate(kwTRANSACTION)
		return nil, errCollecting
	}

	// Optional WORK (ignored, just skip)
	p.matchIdentCI("WORK")

	// Optional TRAN/TRANSACTION
	p.match(kwTRAN)
	p.match(kwTRANSACTION)

	stmt := &nodes.CommitTransStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Optional name (but not WITH, which starts the delayed durability clause)
	if (p.isIdentLike() || p.cur.Type == tokVARIABLE) && p.cur.Type != kwWITH {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// Optional WITH ( DELAYED_DURABILITY = { OFF | ON } )
	if _, ok := p.match(kwWITH); ok {
		if _, ok2 := p.match('('); ok2 {
			if p.matchIdentCI("DELAYED_DURABILITY") {
				p.match('=')
				// COMMIT ... DELAYED_DURABILITY = { ON | OFF }
				if p.isValidOption(onOffOptionValues) {
					stmt.DelayedDurability = strings.ToUpper(p.cur.Str)
					p.advance()
				}
			}
			p.match(')')
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseRollbackStmt parses a ROLLBACK [TRAN/TRANSACTION] statement.
//
// BNF: mssql/parser/bnf/rollback-transaction-transact-sql.bnf
//
//	ROLLBACK { TRAN | TRANSACTION }
//	    [ transaction_name | @tran_name_variable
//	    | savepoint_name | @savepoint_variable ]
//	[ ; ]
func (p *Parser) parseRollbackStmt() (*nodes.RollbackTransStmt, error) {
	loc := p.pos()
	p.advance() // consume ROLLBACK

	// Completion: after ROLLBACK → TRANSACTION, WORK
	if p.collectMode() {
		p.addTokenCandidate(kwTRANSACTION)
		return nil, errCollecting
	}

	// Optional WORK (ignored, just skip)
	p.matchIdentCI("WORK")

	// Optional TRAN/TRANSACTION
	p.match(kwTRAN)
	p.match(kwTRANSACTION)

	stmt := &nodes.RollbackTransStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Optional name/savepoint
	if p.isIdentLike() || p.cur.Type == tokVARIABLE {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseSaveTransStmt parses a SAVE TRAN/TRANSACTION statement.
//
// BNF: mssql/parser/bnf/save-transaction-transact-sql.bnf
//
//	SAVE { TRAN | TRANSACTION } { savepoint_name | @savepoint_variable }
//	[ ; ]
func (p *Parser) parseSaveTransStmt() (*nodes.SaveTransStmt, error) {
	loc := p.pos()
	p.advance() // consume SAVE

	// TRAN or TRANSACTION
	p.match(kwTRAN)
	p.match(kwTRANSACTION)

	// Completion: after SAVE TRANSACTION → identifier
	if p.collectMode() {
		p.addRuleCandidate("identifier")
		return nil, errCollecting
	}

	stmt := &nodes.SaveTransStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Name
	if p.isIdentLike() || p.cur.Type == tokVARIABLE {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}
