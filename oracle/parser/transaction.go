package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseCommitStmt parses a COMMIT statement.
//
// BNF: oracle/parser/bnf/COMMIT.bnf
//
//	COMMIT [ WORK ]
//	    [ COMMENT 'text' ]
//	    [ WRITE [ WAIT | NOWAIT ] [ IMMEDIATE | BATCH ] ]
//	    [ FORCE 'string' [, integer ] ] ;
func (p *Parser) parseCommitStmt() (nodes.StmtNode, error) {
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

	// Optional WRITE [WAIT|NOWAIT] [IMMEDIATE|BATCH]
	if p.isIdentLikeStr("WRITE") {
		p.advance() // consume WRITE
		if p.isIdentLikeStr("WAIT") {
			stmt.WriteWait = "WAIT"
			p.advance()
		} else if p.isIdentLikeStr("NOWAIT") {
			stmt.WriteWait = "NOWAIT"
			p.advance()
		}
		if p.cur.Type == kwIMMEDIATE {
			stmt.WriteBatch = "IMMEDIATE"
			p.advance()
		} else if p.isIdentLikeStr("BATCH") {
			stmt.WriteBatch = "BATCH"
			p.advance()
		}
	}

	// Optional FORCE 'text' [, integer]
	if p.cur.Type == kwFORCE {
		p.advance() // consume FORCE
		if p.cur.Type == tokSCONST {
			stmt.Force = p.cur.Str
			p.advance()
		}
		// optional , integer
		if p.cur.Type == ',' {
			p.advance()
			if p.cur.Type == tokICONST {
				var parseErr1113 error
				stmt.ForceInteger, parseErr1113 = p.parseIntValue()
				if parseErr1113 != nil {
					return nil, parseErr1113
				}
				stmt.HasForceInt = true
			}
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseRollbackStmt parses a ROLLBACK statement.
//
// BNF: oracle/parser/bnf/ROLLBACK.bnf
//
//	ROLLBACK [ WORK ]
//	    [ TO [ SAVEPOINT ] savepoint_name
//	    | FORCE 'string'
//	    ] ;
func (p *Parser) parseRollbackStmt() (nodes.StmtNode, error) {
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
		var parseErr1114 error
		stmt.ToSavepoint, parseErr1114 = p.parseIdentifier()
		if parseErr1114 !=

			// Optional FORCE 'text'
			nil {
			return nil, parseErr1114
		}
	}

	if p.cur.Type == kwFORCE {
		p.advance() // consume FORCE
		if p.cur.Type == tokSCONST {
			stmt.Force = p.cur.Str
			p.advance()
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseSavepointStmt parses a SAVEPOINT statement.
//
// BNF: oracle/parser/bnf/SAVEPOINT.bnf
//
//	SAVEPOINT savepoint_name ;
func (p *Parser) parseSavepointStmt() (nodes.StmtNode, error) {
	start := p.pos()
	p.advance()
	parseValue93, // consume SAVEPOINT
		parseErr94 := p.parseIdentifier()
	if parseErr94 != nil {
		return nil, parseErr94
	}
	stmt := &nodes.SavepointStmt{
		Name: parseValue93,
		Loc:  nodes.Loc{Start: start},
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseSetTransactionStmt parses a SET TRANSACTION statement.
//
// BNF: oracle/parser/bnf/SET-TRANSACTION.bnf
//
//	SET TRANSACTION
//	    { { READ ONLY | READ WRITE }
//	    | ISOLATION LEVEL { SERIALIZABLE | READ COMMITTED }
//	    | USE ROLLBACK SEGMENT rollback_segment
//	    | NAME 'string'
//	    }... ;
func (p *Parser) parseSetTransactionStmt() (nodes.StmtNode, error) {
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
			if p.isIdentLike() && p.cur.Str == "USE" {
				p.advance() // consume USE
				if p.cur.Type == kwROLLBACK {
					p.advance() // consume ROLLBACK
					// SEGMENT is not a keyword
					if p.isIdentLike() && p.cur.Str == "SEGMENT" {
						p.advance() // consume SEGMENT
					}
					var parseErr1115 error
					// segment_name
					stmt.UseRollbackSegment, parseErr1115 = p.parseIdentifier()
					if parseErr1115 != nil {
						return nil, parseErr1115
					}
				}
			} else if p.isIdentLike() && p.cur.Str == "NAME" {
				// Check for NAME 'text' (NAME is an identifier, not a keyword)
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
	stmt.Loc.End = p.prev.End
	return stmt, nil
}
