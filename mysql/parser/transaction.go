package parser

import (
	nodes "github.com/bytebase/omni/mysql/ast"
)

// parseBeginStmt parses a BEGIN or START TRANSACTION statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/commit.html
//
//	START TRANSACTION [transaction_characteristic [, transaction_characteristic] ...]
//	BEGIN [WORK]
//
//	transaction_characteristic:
//	    WITH CONSISTENT SNAPSHOT
//	    | READ WRITE
//	    | READ ONLY
func (p *Parser) parseBeginStmt() (*nodes.BeginStmt, error) {
	start := p.pos()
	stmt := &nodes.BeginStmt{Loc: nodes.Loc{Start: start}}

	if p.cur.Type == kwSTART {
		p.advance() // consume START
		p.match(kwTRANSACTION)

		// Transaction characteristics
		for {
			if p.cur.Type == kwWITH {
				p.advance()
				p.match(kwCONSISTENT)
				p.match(kwSNAPSHOT)
				stmt.WithConsistentSnapshot = true
			} else if p.cur.Type == kwREAD {
				p.advance()
				if _, ok := p.match(kwONLY); ok {
					stmt.ReadOnly = true
				} else if _, ok := p.match(kwWRITE); ok {
					stmt.ReadWrite = true
				}
			} else {
				break
			}
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	} else {
		// BEGIN [WORK]
		p.advance() // consume BEGIN

		// Completion: after BEGIN, offer WORK keyword.
		p.checkCursor()
		if p.collectMode() {
			return nil, &ParseError{Message: "collecting"}
		}

		if p.cur.Type == kwWORK {
			p.advance()
		}
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseCommitStmt parses a COMMIT statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/commit.html
//
//	COMMIT [WORK] [AND [NO] CHAIN] [[NO] RELEASE]
func (p *Parser) parseCommitStmt() (*nodes.CommitStmt, error) {
	start := p.pos()
	p.advance() // consume COMMIT

	// Completion: after COMMIT, offer AND and WORK keywords.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwAND)
		return nil, &ParseError{Message: "collecting"}
	}

	stmt := &nodes.CommitStmt{Loc: nodes.Loc{Start: start}}

	// Optional WORK
	if p.cur.Type == kwWORK {
		p.advance()
	}

	// AND [NO] CHAIN
	if p.cur.Type == kwAND {
		p.advance()
		if p.cur.Type == kwNO {
			p.advance()
			p.match(kwCHAIN)
		} else {
			p.match(kwCHAIN)
			stmt.Chain = true
		}
	}

	// [NO] RELEASE
	if p.cur.Type == kwNO {
		p.advance()
		p.match(kwRELEASE)
	} else if _, ok := p.match(kwRELEASE); ok {
		stmt.Release = true
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseRollbackStmt parses a ROLLBACK statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/commit.html
//
//	ROLLBACK [WORK] [AND [NO] CHAIN] [[NO] RELEASE]
//	ROLLBACK [WORK] TO [SAVEPOINT] identifier
func (p *Parser) parseRollbackStmt() (*nodes.RollbackStmt, error) {
	start := p.pos()
	p.advance() // consume ROLLBACK

	// Completion: after ROLLBACK, offer TO and WORK keywords.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwTO)
		return nil, &ParseError{Message: "collecting"}
	}

	stmt := &nodes.RollbackStmt{Loc: nodes.Loc{Start: start}}

	// Optional WORK
	if p.cur.Type == kwWORK {
		p.advance()
	}

	// TO [SAVEPOINT] identifier
	if _, ok := p.match(kwTO); ok {
		// Completion: after ROLLBACK TO, offer SAVEPOINT keyword.
		p.checkCursor()
		if p.collectMode() {
			p.addTokenCandidate(kwSAVEPOINT)
			return nil, &ParseError{Message: "collecting"}
		}

		p.match(kwSAVEPOINT)
		name, _, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		stmt.Savepoint = name
		stmt.Loc.End = p.pos()
		return stmt, nil
	}

	// AND [NO] CHAIN
	if p.cur.Type == kwAND {
		p.advance()
		if p.cur.Type == kwNO {
			p.advance()
			p.match(kwCHAIN)
		} else {
			p.match(kwCHAIN)
			stmt.Chain = true
		}
	}

	// [NO] RELEASE
	if p.cur.Type == kwNO {
		p.advance()
		p.match(kwRELEASE)
	} else if _, ok := p.match(kwRELEASE); ok {
		stmt.Release = true
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseSavepointStmt parses a SAVEPOINT statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/savepoint.html
//
//	SAVEPOINT identifier
func (p *Parser) parseSavepointStmt() (*nodes.SavepointStmt, error) {
	start := p.pos()
	p.advance() // consume SAVEPOINT

	// Completion: after SAVEPOINT, identifier context (user defines a new savepoint name).
	p.checkCursor()
	if p.collectMode() {
		return nil, &ParseError{Message: "collecting"}
	}

	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}

	return &nodes.SavepointStmt{
		Loc:  nodes.Loc{Start: start, End: p.pos()},
		Name: name,
	}, nil
}

// parseReleaseSavepointStmt parses RELEASE SAVEPOINT.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/savepoint.html
//
//	RELEASE SAVEPOINT identifier
func (p *Parser) parseReleaseSavepointStmt() (*nodes.SavepointStmt, error) {
	start := p.pos()
	p.advance() // consume RELEASE
	p.match(kwSAVEPOINT)

	// Completion: after RELEASE SAVEPOINT, identifier context.
	p.checkCursor()
	if p.collectMode() {
		return nil, &ParseError{Message: "collecting"}
	}

	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}

	return &nodes.SavepointStmt{
		Loc:     nodes.Loc{Start: start, End: p.pos()},
		Name:    name,
		Release: true,
	}, nil
}

// parseLockTablesStmt parses a LOCK TABLES statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/lock-tables.html
//
//	LOCK {TABLE | TABLES}
//	    tbl_name [[AS] alias] lock_type
//	    [, tbl_name [[AS] alias] lock_type] ...
//
//	lock_type: {
//	    READ [LOCAL]
//	  | [LOW_PRIORITY] WRITE
//	}
func (p *Parser) parseLockTablesStmt() (*nodes.LockTablesStmt, error) {
	start := p.pos()
	p.advance() // consume TABLES

	// Completion: after LOCK TABLES, offer table_ref.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("table_ref")
		return nil, &ParseError{Message: "collecting"}
	}

	stmt := &nodes.LockTablesStmt{Loc: nodes.Loc{Start: start}}

	for {
		lockStart := p.pos()
		ref, err := p.parseTableRefWithAlias()
		if err != nil {
			return nil, err
		}

		lt := &nodes.LockTable{
			Loc:   nodes.Loc{Start: lockStart},
			Table: ref,
		}

		// Completion: after table name, offer READ/WRITE.
		p.checkCursor()
		if p.collectMode() {
			p.addTokenCandidate(kwREAD)
			p.addTokenCandidate(kwWRITE)
			return nil, &ParseError{Message: "collecting"}
		}

		// Lock type
		if _, ok := p.match(kwREAD); ok {
			if _, ok := p.match(kwLOCAL); ok {
				lt.LockType = "READ LOCAL"
			} else {
				lt.LockType = "READ"
			}
		} else if _, ok := p.match(kwLOW_PRIORITY); ok {
			p.match(kwWRITE)
			lt.LockType = "LOW_PRIORITY WRITE"
		} else if _, ok := p.match(kwWRITE); ok {
			lt.LockType = "WRITE"
		}

		lt.Loc.End = p.pos()
		stmt.Tables = append(stmt.Tables, lt)

		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseUnlockTablesStmt parses an UNLOCK TABLES statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/lock-tables.html
//
//	UNLOCK {TABLE | TABLES}
func (p *Parser) parseUnlockTablesStmt() (*nodes.UnlockTablesStmt, error) {
	start := p.pos()
	p.advance() // consume TABLES

	return &nodes.UnlockTablesStmt{
		Loc: nodes.Loc{Start: start, End: p.pos()},
	}, nil
}

// parseSetTransactionStmt parses SET [GLOBAL|SESSION] TRANSACTION.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/set-transaction.html
//
//	SET [GLOBAL | SESSION] TRANSACTION
//	    transaction_characteristic [, transaction_characteristic] ...
//
//	transaction_characteristic: {
//	    ISOLATION LEVEL level
//	  | access_mode
//	}
//
//	level: {
//	     REPEATABLE READ
//	   | READ COMMITTED
//	   | READ UNCOMMITTED
//	   | SERIALIZABLE
//	}
//
//	access_mode: {
//	     READ WRITE
//	   | READ ONLY
//	}
func (p *Parser) parseSetTransactionStmt(start int, scope string) (*nodes.SetTransactionStmt, error) {
	p.advance() // consume TRANSACTION

	stmt := &nodes.SetTransactionStmt{
		Loc:   nodes.Loc{Start: start},
		Scope: scope,
	}

	// Parse transaction characteristics
	for {
		if p.cur.Type == kwISOLATION {
			p.advance()
			p.match(kwLEVEL)
			switch p.cur.Type {
			case kwREPEATABLE:
				p.advance()
				p.match(kwREAD)
				stmt.IsolationLevel = "REPEATABLE READ"
			case kwREAD:
				p.advance()
				if _, ok := p.match(kwCOMMITTED); ok {
					stmt.IsolationLevel = "READ COMMITTED"
				} else if _, ok := p.match(kwUNCOMMITTED); ok {
					stmt.IsolationLevel = "READ UNCOMMITTED"
				}
			case kwSERIALIZABLE:
				p.advance()
				stmt.IsolationLevel = "SERIALIZABLE"
			}
		} else if p.cur.Type == kwREAD {
			p.advance()
			if _, ok := p.match(kwONLY); ok {
				stmt.AccessMode = "READ ONLY"
			} else if _, ok := p.match(kwWRITE); ok {
				stmt.AccessMode = "READ WRITE"
			}
		} else {
			break
		}

		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseXAStmt parses an XA statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/xa-statements.html
//
//	XA {START|BEGIN} xid [JOIN|RESUME]
//	XA END xid [SUSPEND [FOR MIGRATE]]
//	XA PREPARE xid
//	XA COMMIT xid [ONE PHASE]
//	XA ROLLBACK xid
//	XA RECOVER [CONVERT XID]
//
//	xid: gtrid [, bqual [, formatID]]
func (p *Parser) parseXAStmt() (*nodes.XAStmt, error) {
	start := p.pos()
	p.advance() // consume XA

	stmt := &nodes.XAStmt{Loc: nodes.Loc{Start: start}}

	switch p.cur.Type {
	case kwSTART, kwBEGIN:
		stmt.Type = nodes.XAStart
		p.advance()
		xid, err := p.parseXid()
		if err != nil {
			return nil, err
		}
		stmt.Xid = xid
		if _, ok := p.match(kwJOIN); ok {
			stmt.Join = true
		} else if _, ok := p.match(kwRESUME); ok {
			stmt.Resume = true
		}

	case kwEND:
		stmt.Type = nodes.XAEnd
		p.advance()
		xid, err := p.parseXid()
		if err != nil {
			return nil, err
		}
		stmt.Xid = xid
		if _, ok := p.match(kwSUSPEND); ok {
			stmt.Suspend = true
			if p.cur.Type == kwFOR {
				p.advance()
				p.match(kwMIGRATE)
				stmt.Migrate = true
			}
		}

	case kwPREPARE:
		stmt.Type = nodes.XAPrepare
		p.advance()
		xid, err := p.parseXid()
		if err != nil {
			return nil, err
		}
		stmt.Xid = xid

	case kwCOMMIT:
		stmt.Type = nodes.XACommit
		p.advance()
		xid, err := p.parseXid()
		if err != nil {
			return nil, err
		}
		stmt.Xid = xid
		if _, ok := p.match(kwONE); ok {
			p.match(kwPHASE)
			stmt.OnePhase = true
		}

	case kwROLLBACK:
		stmt.Type = nodes.XARollback
		p.advance()
		xid, err := p.parseXid()
		if err != nil {
			return nil, err
		}
		stmt.Xid = xid

	case kwRECOVER:
		stmt.Type = nodes.XARecover
		p.advance()
		if _, ok := p.match(kwCONVERT); ok {
			// CONVERT XID
			if p.cur.Type == kwXID {
				p.advance()
			}
			stmt.Convert = true
		}

	default:
		return nil, &ParseError{
			Message:  "expected START, END, PREPARE, COMMIT, ROLLBACK, or RECOVER after XA",
			Position: p.cur.Loc,
		}
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseXid parses an xid: gtrid [, bqual [, formatID]].
func (p *Parser) parseXid() ([]nodes.ExprNode, error) {
	var parts []nodes.ExprNode

	gtrid, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if gtrid == nil {
		return nil, p.syntaxErrorAtCur()
	}
	parts = append(parts, gtrid)

	if p.cur.Type == ',' {
		p.advance()
		bqual, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if bqual == nil {
			return nil, p.syntaxErrorAtCur()
		}
		parts = append(parts, bqual)

		if p.cur.Type == ',' {
			p.advance()
			formatID, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if formatID == nil {
				return nil, p.syntaxErrorAtCur()
			}
			parts = append(parts, formatID)
		}
	}

	return parts, nil
}
