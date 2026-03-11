package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseCommitStmt parses a COMMIT statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/COMMIT.html
//
//	COMMIT [WORK] [COMMENT 'text'] [FORCE 'text']
func (p *Parser) parseCommitStmt() nodes.StmtNode {
	start := p.pos()
	p.advance() // consume COMMIT

	stmt := &nodes.CommitStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Optional WORK
	if p.cur.Type == kwWORK {
		stmt.Work = true
		p.advance()
	}

	// Optional COMMENT 'text'
	if p.cur.Type == kwCOMMENT {
		p.advance() // consume COMMENT
		if p.cur.Type == tokSCONST {
			stmt.Comment = p.cur.Str
			p.advance()
		}
	}

	// Optional FORCE 'text'
	if p.cur.Type == kwFORCE {
		p.advance() // consume FORCE
		if p.cur.Type == tokSCONST {
			stmt.Force = p.cur.Str
			p.advance()
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseRollbackStmt parses a ROLLBACK statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/ROLLBACK.html
//
//	ROLLBACK [WORK] [TO [SAVEPOINT] savepoint_name] [FORCE 'text']
func (p *Parser) parseRollbackStmt() nodes.StmtNode {
	start := p.pos()
	p.advance() // consume ROLLBACK

	stmt := &nodes.RollbackStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Optional WORK
	if p.cur.Type == kwWORK {
		stmt.Work = true
		p.advance()
	}

	// Optional TO [SAVEPOINT] name
	if p.cur.Type == kwTO {
		p.advance() // consume TO
		if p.cur.Type == kwSAVEPOINT {
			p.advance() // consume SAVEPOINT
		}
		stmt.ToSavepoint = p.parseIdentifier()
	}

	// Optional FORCE 'text'
	if p.cur.Type == kwFORCE {
		p.advance() // consume FORCE
		if p.cur.Type == tokSCONST {
			stmt.Force = p.cur.Str
			p.advance()
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseSavepointStmt parses a SAVEPOINT statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/SAVEPOINT.html
//
//	SAVEPOINT name
func (p *Parser) parseSavepointStmt() nodes.StmtNode {
	start := p.pos()
	p.advance() // consume SAVEPOINT

	stmt := &nodes.SavepointStmt{
		Name: p.parseIdentifier(),
		Loc:  nodes.Loc{Start: start},
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseSetTransactionStmt parses a SET TRANSACTION statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/SET-TRANSACTION.html
//
//	SET TRANSACTION
//	    { READ ONLY | READ WRITE |
//	      ISOLATION LEVEL { SERIALIZABLE | READ COMMITTED } |
//	      USE ROLLBACK SEGMENT ... }
//	    [NAME 'text']
func (p *Parser) parseSetTransactionStmt() nodes.StmtNode {
	start := p.pos()
	// SET and TRANSACTION already consumed by the dispatcher
	// p.advance() for SET is done in parseStmt; TRANSACTION consumed there too

	stmt := &nodes.SetTransactionStmt{
		Loc: nodes.Loc{Start: start},
	}

	for {
		switch p.cur.Type {
		case kwREAD:
			p.advance() // consume READ
			if p.cur.Type == kwONLY {
				stmt.ReadOnly = true
				p.advance()
			} else if p.isIdentLike() && p.cur.Str == "WRITE" {
				stmt.ReadWrite = true
				p.advance()
			}
		case kwISOLATION:
			p.advance() // consume ISOLATION
			if p.cur.Type == kwLEVEL {
				p.advance() // consume LEVEL
			}
			// SERIALIZABLE or READ COMMITTED
			if p.isIdentLike() && p.cur.Str == "SERIALIZABLE" {
				stmt.IsolLevel = "SERIALIZABLE"
				p.advance()
			} else if p.cur.Type == kwREAD {
				p.advance() // consume READ
				// "COMMITTED" is not a keyword, it's an ident
				if p.isIdentLike() && p.cur.Str == "COMMITTED" {
					stmt.IsolLevel = "READ COMMITTED"
					p.advance()
				}
			}
		default:
			// Check for NAME 'text' (NAME is an identifier, not a keyword)
			if p.isIdentLike() && p.cur.Str == "NAME" {
				p.advance() // consume NAME
				if p.cur.Type == tokSCONST {
					stmt.Name = p.cur.Str
					p.advance()
				}
			} else {
				goto done
			}
		}
	}

done:
	stmt.Loc.End = p.pos()
	return stmt
}
