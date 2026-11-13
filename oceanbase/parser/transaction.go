package parser

import (
	nodes "github.com/bytebase/omni/oceanbase/ast"
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
		return p.parseStartTransactionTail(start)
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

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseStartTransactionTail parses START TRANSACTION after START:
//
//	START TRANSACTION [WITH CONSISTENT SNAPSHOT [, access] | access [, WITH CONSISTENT SNAPSHOT]]
//
//	access: READ ONLY | READ WRITE
//
// Each characteristic is whole and named once, and TRANSACTION is required.
//
// Ref: begin_stmt, opt_with_consistent_snapshot and transaction_access_mode in
// sql_parser_mysql_mode.y (oceanbase v4.3.5_CE_BP6).
func (p *Parser) parseStartTransactionTail(start int) (*nodes.BeginStmt, error) {
	if _, err := p.expect(kwTRANSACTION); err != nil {
		return nil, err
	}

	// Completion: after START TRANSACTION, offer transaction characteristic keywords.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwWITH)
		p.addTokenCandidate(kwREAD)
		return nil, &ParseError{Message: "collecting"}
	}

	stmt := &nodes.BeginStmt{Loc: nodes.Loc{Start: start}}
	snapshot := func() error {
		for _, kw := range []int{kwWITH, kwCONSISTENT, kwSNAPSHOT} {
			if _, err := p.expect(kw); err != nil {
				return err
			}
		}
		stmt.WithConsistentSnapshot = true
		return nil
	}
	access := func() error {
		if _, err := p.expect(kwREAD); err != nil {
			return err
		}
		switch p.cur.Type {
		case kwONLY:
			stmt.ReadOnly = true
		case kwWRITE:
			stmt.ReadWrite = true
		default:
			return p.syntaxErrorAtCur()
		}
		p.advance()
		return nil
	}
	first, second := snapshot, access
	switch p.cur.Type {
	case kwREAD:
		first, second = access, snapshot
	case kwWITH:
	default:
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}
	if err := first(); err != nil {
		return nil, err
	}
	if p.cur.Type == ',' {
		p.advance()
		if err := second(); err != nil {
			return nil, err
		}
	}
	stmt.Loc.End = p.prev.End
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

	stmt.Loc.End = p.prev.End
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

		// ROLLBACK TO SAVEPOINT var_name, or ROLLBACK TO var_name, where
		// SAVEPOINT, an unreserved keyword, may be the name itself.
		if p.cur.Type == kwSAVEPOINT && p.peekNext().Type != tokEOF && p.peekNext().Type != ';' {
			p.advance()
		}
		name, _, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		stmt.Savepoint = name
		stmt.Loc.End = p.prev.End
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

	stmt.Loc.End = p.prev.End
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
		Loc:  nodes.Loc{Start: start, End: p.prev.End},
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
	if _, err := p.expect(kwSAVEPOINT); err != nil {
		return nil, err
	}

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
		Loc:     nodes.Loc{Start: start, End: p.prev.End},
		Name:    name,
		Release: true,
	}, nil
}

// parseLockTablesStmt parses a LOCK TABLES statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/lock-tables.html
//
//	LOCK {TABLE | TABLES}
//	    tbl_name [AS alias] lock_type
//	    [, tbl_name [AS alias] lock_type] ...
//
// OceanBase's lock_table takes an alias only after AS, and a lock_type always
// (lock_tables_stmt, lock_table and lock_type in sql_parser_mysql_mode.y,
// oceanbase v4.3.5_CE_BP6).
//
//	lock_type: {
//	    READ [LOCAL]
//	  | [LOW_PRIORITY] WRITE
//	}
func (p *Parser) parseLockTablesStmt(start int) (*nodes.LockTablesStmt, error) {
	// LOCK names the kind of thing it locks before the list, and the keyword is
	// not optional: 5.7.25-OceanBase_CE-v4.3.5.6 answers 1064 for LOCK t READ and for LOCK = t READ. The
	// statement is reached with LOCK consumed, so taking the next token without
	// reading it left a run-on -- DROP INDEX i ON t LOCK = SHARED -- parsing as
	// a DROP INDEX and a LOCK TABLES of a table called shared.
	//
	// Ref: table_or_tables in MySQL 8.0 sql/sql_yacc.yy.
	if _, ok := p.match(kwTABLES, kwTABLE); !ok {
		return nil, p.syntaxErrorAtCur()
	}

	// Completion: after LOCK TABLES, offer table_ref.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("table_ref")
		return nil, &ParseError{Message: "collecting"}
	}

	stmt := &nodes.LockTablesStmt{Loc: nodes.Loc{Start: start}}

	for {
		lockStart := p.pos()
		// relation_factor [AS relation_name]: the alias takes AS.
		ref, err := p.parseRelationFactor()
		if err != nil {
			return nil, err
		}
		if _, ok := p.match(kwAS); ok {
			alias, _, err := p.parseIdent()
			if err != nil {
				return nil, err
			}
			ref.Alias = alias
			ref.Loc.End = p.prev.End
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
			if _, err := p.expect(kwWRITE); err != nil {
				return nil, err
			}
			lt.LockType = "LOW_PRIORITY WRITE"
		} else if _, ok := p.match(kwWRITE); ok {
			lt.LockType = "WRITE"
		} else {
			// lock_type is required.
			return nil, p.syntaxErrorAtCur()
		}

		lt.Loc.End = p.prev.End
		stmt.Tables = append(stmt.Tables, lt)

		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseUnlockTablesStmt parses an UNLOCK TABLES statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/lock-tables.html
//
//	UNLOCK {TABLE | TABLES}
//
// The keyword is not optional, as it is not after LOCK: 5.7.25-OceanBase_CE-v4.3.5.6 answers 1064 for
// UNLOCK nosuch, and taking the next token without reading it let it through.
//
// Ref: table_or_tables in MySQL 8.0 sql/sql_yacc.yy.
func (p *Parser) parseUnlockTablesStmt(start int) (*nodes.UnlockTablesStmt, error) {
	if _, ok := p.match(kwTABLES, kwTABLE); !ok {
		return nil, p.syntaxErrorAtCur()
	}

	return &nodes.UnlockTablesStmt{
		Loc: nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parseSetTransactionStmt parses SET [scope] TRANSACTION, with the current token on
// TRANSACTION.
//
// Ref: set_transaction_stmt and transaction_characteristics in OceanBase CE v4.3.5_CE_BP6
// src/sql/parser/sql_parser_mysql_mode.y (lines 23975-24051).
//
//	SET [GLOBAL | SESSION | LOCAL] TRANSACTION
//	    {access_mode [, ISOLATION LEVEL level] | ISOLATION LEVEL level [, access_mode]}
//
//	access_mode: READ WRITE | READ ONLY
//
//	level: REPEATABLE READ | READ COMMITTED | READ UNCOMMITTED | SERIALIZABLE
//
// The list may not be empty, each characteristic has to be whole, and it names each kind
// at most once, so `SET TRANSACTION READ ONLY, READ WRITE` is a 1064. TiDB's TransactionChars
// differs: it takes any number of either.
func (p *Parser) parseSetTransactionStmt(start int, scope string) (*nodes.SetTransactionStmt, error) {
	p.advance() // consume TRANSACTION

	stmt := &nodes.SetTransactionStmt{
		Loc:   nodes.Loc{Start: start},
		Scope: scope,
	}
	if err := p.parseTransactionChar(stmt); err != nil {
		return nil, err
	}
	if _, ok := p.match(','); ok {
		if err := p.parseTransactionChar(stmt); err != nil {
			return nil, err
		}
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseTransactionChar parses one characteristic of SET TRANSACTION into stmt, taking
// only a kind stmt does not hold yet.
func (p *Parser) parseTransactionChar(stmt *nodes.SetTransactionStmt) error {
	wantLevel, wantAccess := stmt.IsolationLevel == "", stmt.AccessMode == ""

	// Completion: where a characteristic starts.
	p.checkCursor()
	if p.collectMode() {
		if wantLevel {
			p.addTokenCandidate(kwISOLATION)
		}
		if wantAccess {
			p.addTokenCandidate(kwREAD)
		}
		return &ParseError{Message: "collecting"}
	}

	switch {
	case wantLevel && p.cur.Type == kwISOLATION:
		p.advance()
		if _, err := p.expect(kwLEVEL); err != nil {
			return err
		}
		level, err := p.parseIsolationLevel()
		if err != nil {
			return err
		}
		stmt.IsolationLevel = level
	case wantAccess && p.cur.Type == kwREAD:
		p.advance()
		// Completion: after READ, offer ONLY and WRITE.
		p.checkCursor()
		if p.collectMode() {
			p.addTokenCandidate(kwONLY)
			p.addTokenCandidate(kwWRITE)
			return &ParseError{Message: "collecting"}
		}
		switch p.cur.Type {
		case kwWRITE:
			p.advance()
			stmt.AccessMode = "READ WRITE"
		case kwONLY:
			p.advance()
			stmt.AccessMode = "READ ONLY"
		default:
			return p.syntaxErrorAtCur()
		}
	default:
		return p.syntaxErrorAtCur()
	}
	return nil
}

// parseIsolationLevel parses the level after ISOLATION LEVEL.
func (p *Parser) parseIsolationLevel() (string, error) {
	// Completion: after ISOLATION LEVEL, offer the words a level starts with.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwREPEATABLE)
		p.addTokenCandidate(kwREAD)
		p.addTokenCandidate(kwSERIALIZABLE)
		return "", &ParseError{Message: "collecting"}
	}

	switch p.cur.Type {
	case kwREPEATABLE:
		p.advance()
		if _, err := p.expect(kwREAD); err != nil {
			return "", err
		}
		return "REPEATABLE READ", nil
	case kwREAD:
		p.advance()
		// Completion: after READ, offer COMMITTED and UNCOMMITTED.
		p.checkCursor()
		if p.collectMode() {
			p.addTokenCandidate(kwCOMMITTED)
			p.addTokenCandidate(kwUNCOMMITTED)
			return "", &ParseError{Message: "collecting"}
		}
		switch p.cur.Type {
		case kwCOMMITTED:
			p.advance()
			return "READ COMMITTED", nil
		case kwUNCOMMITTED:
			p.advance()
			return "READ UNCOMMITTED", nil
		}
		return "", p.syntaxErrorAtCur()
	case kwSERIALIZABLE:
		p.advance()
		return "SERIALIZABLE", nil
	}
	return "", p.syntaxErrorAtCur()
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
				if _, err := p.expect(kwMIGRATE); err != nil {
					return nil, err
				}
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
			if _, err := p.expect(kwPHASE); err != nil {
				return nil, err
			}
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
			if _, err := p.expect(kwXID); err != nil {
				return nil, err
			}
			stmt.Convert = true
		}

	default:
		return nil, &ParseError{
			Message:  "expected START, END, PREPARE, COMMIT, ROLLBACK, or RECOVER after XA",
			Position: p.cur.Loc,
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseXid parses xa_xid:
//
//	text_string [, text_string [, ulong_num]]
//
// text_string is a string, hex or bit literal, ulong_num an integer or a hex
// literal (sql_parser_mysql_mode.y, oceanbase v4.3.5_CE_BP6). Strings side by
// side are one literal, as OceanBase's lexer reads them.
func (p *Parser) parseXid() ([]nodes.ExprNode, error) {
	var parts []nodes.ExprNode
	literal := func(ok func(int) bool) error {
		if !ok(p.cur.Type) {
			return p.syntaxErrorAtCur()
		}
		e, err := p.parsePrimaryExpr()
		if err != nil {
			return err
		}
		parts = append(parts, e)
		return nil
	}
	textString := func(t int) bool { return t == tokSCONST || t == tokXCONST || t == tokBCONST }
	ulongNum := func(t int) bool { return t == tokICONST || t == tokXCONST }
	if err := literal(textString); err != nil {
		return nil, err
	}
	if p.cur.Type == ',' {
		p.advance()
		if err := literal(textString); err != nil {
			return nil, err
		}
		if p.cur.Type == ',' {
			p.advance()
			if err := literal(ulongNum); err != nil {
				return nil, err
			}
		}
	}
	return parts, nil
}
