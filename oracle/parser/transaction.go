package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseCommitStmt parses a COMMIT statement.
//
// BNF: oracle/parser/bnf/COMMIT.bnf
//
//	COMMIT [ WORK ]
//	    [ [ COMMENT 'text' ]
//	      [ WRITE [ WAIT | NOWAIT ] [ IMMEDIATE | BATCH ] ]
//	    | FORCE { 'text' [, integer ] | CORRUPT_XID 'text' | CORRUPT_XID_ALL }
//	    ] ;
//
// Oracle 23ai takes FORCE only right after COMMIT [ WORK ] (ORA-00922 at
// FORCE after COMMENT or WRITE), WAIT and IMMEDIATE in either order, and
// after FORCE 'text' a number, a string or a word, with or without a comma
// before a string or a word but not before a number (ORA-00922).
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

	if p.cur.Type == kwFORCE {
		p.advance() // consume FORCE
		switch {
		case p.isKeywordStr("CORRUPT_XID"):
			p.advance()
			if p.cur.Type != tokSCONST {
				return nil, p.syntaxErrorAtCur()
			}
			stmt.Force = p.cur.Str
			p.advance()
		case p.isKeywordStr("CORRUPT_XID_ALL"):
			p.advance()
		case p.cur.Type == tokSCONST:
			stmt.Force = p.cur.Str
			p.advance()
			switch {
			case p.cur.Type == ',':
				p.advance()
				if p.cur.Type != tokSCONST && !p.isIdentLike() {
					return nil, p.syntaxErrorAtCur()
				}
				p.advance()
			case p.cur.Type == tokICONST:
				var err error
				stmt.ForceInteger, err = p.parseIntValue()
				if err != nil {
					return nil, err
				}
				stmt.HasForceInt = true
			case p.cur.Type == tokFCONST || p.cur.Type == tokSCONST || p.isIdentLike():
				p.advance()
			}
		default:
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

	// Optional COMMENT 'text'
	if p.cur.Type == kwCOMMENT {
		p.advance() // consume COMMENT
		if p.cur.Type != tokSCONST {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Comment = p.cur.Str
		p.advance()
	}

	// Optional WRITE [WAIT|NOWAIT] [IMMEDIATE|BATCH], in either order
	if p.isIdentLikeStr("WRITE") {
		p.advance() // consume WRITE
		for {
			if stmt.WriteWait == "" && p.isIdentLikeStr("WAIT") {
				stmt.WriteWait = "WAIT"
				p.advance()
			} else if stmt.WriteWait == "" && p.isIdentLikeStr("NOWAIT") {
				stmt.WriteWait = "NOWAIT"
				p.advance()
			} else if stmt.WriteBatch == "" && p.cur.Type == kwIMMEDIATE {
				stmt.WriteBatch = "IMMEDIATE"
				p.advance()
			} else if stmt.WriteBatch == "" && p.isIdentLikeStr("BATCH") {
				stmt.WriteBatch = "BATCH"
				p.advance()
			} else {
				break
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

	switch p.cur.Type {
	case kwTO:
		p.advance() // consume TO
		if p.cur.Type == kwSAVEPOINT {
			p.advance() // consume SAVEPOINT
		}
		var err error
		stmt.ToSavepoint, err = p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		if stmt.ToSavepoint == "" {
			return nil, p.syntaxErrorAtCur()
		}
	case kwFORCE:
		p.advance() // consume FORCE
		if p.cur.Type != tokSCONST {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Force = p.cur.Str
		p.advance()
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

// parseSetTransactionAtSet parses SET TRANSACTION ... with the current token
// at SET.
func (p *Parser) parseSetTransactionAtSet() (nodes.StmtNode, error) {
	start := p.pos()
	p.advance() // consume SET
	p.advance() // consume TRANSACTION
	stmt, err := p.parseSetTransactionStmt()
	if err != nil {
		return nil, err
	}
	stmt.(*nodes.SetTransactionStmt).Loc.Start = start
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

	// One clause, then NAME 'text'; or NAME 'text' alone. Oracle 23ai
	// takes no second clause (ORA-03049) and no clause after NAME
	// (ORA-00933).
	switch {
	case p.cur.Type == kwREAD:
		p.advance() // consume READ
		if p.cur.Type == kwONLY {
			stmt.ReadOnly = true
			p.advance()
		} else if p.isKeywordStr("WRITE") {
			stmt.ReadWrite = true
			p.advance()
		} else {
			return nil, p.syntaxErrorAtCur()
		}
	case p.cur.Type == kwISOLATION:
		p.advance() // consume ISOLATION
		if p.cur.Type != kwLEVEL {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance() // consume LEVEL
		// SERIALIZABLE or READ COMMITTED
		if p.isKeywordStr("SERIALIZABLE") {
			stmt.IsolLevel = "SERIALIZABLE"
			p.advance()
		} else if p.cur.Type == kwREAD {
			p.advance() // consume READ
			if err := p.expectWords("COMMITTED"); err != nil {
				return nil, err
			}
			stmt.IsolLevel = "READ COMMITTED"
		} else {
			return nil, p.syntaxErrorAtCur()
		}
	case p.isKeywordStr("USE"):
		p.advance() // consume USE
		if p.cur.Type != kwROLLBACK {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance() // consume ROLLBACK
		if err := p.expectWords("SEGMENT"); err != nil {
			return nil, err
		}
		var parseErr1115 error
		// segment_name
		stmt.UseRollbackSegment, parseErr1115 = p.parseIdentifier()
		if parseErr1115 != nil {
			return nil, parseErr1115
		}
		if stmt.UseRollbackSegment == "" {
			return nil, p.syntaxErrorAtCur()
		}
	case p.isKeywordStr("NAME"):
	default:
		return nil, p.syntaxErrorAtCur()
	}
	if p.isKeywordStr("NAME") {
		p.advance() // consume NAME
		if p.cur.Type != tokSCONST {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}
