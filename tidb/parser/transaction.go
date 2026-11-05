package parser

import (
	nodes "github.com/bytebase/omni/tidb/ast"
)

// TiDB's transaction statements are its own and not MySQL's, which this
// parser began as a copy of. TiDB has no WORK at all: BEGIN WORK, COMMIT WORK
// and ROLLBACK WORK are each 1064 on v8.5.5. BEGIN takes a transaction mode
// MySQL lacks, START TRANSACTION takes one characteristic rather than a list,
// and two of the characteristics are TiDB's own. Every verdict quoted in this
// file was measured against pingcap/tidb:v8.5.5 through PREPARE on 2026-11-05
// and is pinned by TestTransactionFormsMatchTiDB.
//
// Ref: BeginTransactionStmt, CommitStmt, RollbackStmt,
// CompletionTypeWithinTransaction, SavepointStmt and ReleaseSavepointStmt in
// TiDB v8.5.5 pkg/parser/parser.y (lines 3431, 3599, 9209, 9227, 5603, 5609).

// parseBeginStmt parses BEGIN, with the transaction mode TiDB lets it name.
//
//	BEGIN [PESSIMISTIC | OPTIMISTIC]
//
// The two modes are not keywords in this lexer, so they are read as bare
// words: v8.5.5 answers 1064 for BEGIN `PESSIMISTIC`.
func (p *Parser) parseBeginStmt() (*nodes.BeginStmt, error) {
	start := p.pos()
	p.advance() // consume BEGIN

	// Completion: after BEGIN. The two modes are bare words, which are no
	// token candidates.
	p.checkCursor()
	if p.collectMode() {
		return nil, &ParseError{Message: "collecting"}
	}

	stmt := &nodes.BeginStmt{Loc: nodes.Loc{Start: start}}
	switch {
	case p.isBareWord("PESSIMISTIC"):
		stmt.Mode = "PESSIMISTIC"
		p.advance()
	case p.isBareWord("OPTIMISTIC"):
		stmt.Mode = "OPTIMISTIC"
		p.advance()
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseStartTransactionStmt parses START TRANSACTION, with START consumed at
// start.
//
//	START TRANSACTION [
//	    READ WRITE
//	  | READ ONLY [AS OF TIMESTAMP expr]
//	  | WITH CONSISTENT SNAPSHOT
//	  | WITH CAUSAL CONSISTENCY ONLY
//	]
//
// One characteristic at most: upstream spells out each arm, so MySQL's
// comma-separated list is 1064 on v8.5.5 (`START TRANSACTION READ ONLY, WITH
// CONSISTENT SNAPSHOT`), and so is a characteristic cut short (`START
// TRANSACTION WITH CONSISTENT`) or a bare START. AS OF TIMESTAMP belongs to
// READ ONLY alone and takes any expression. CAUSAL and CONSISTENCY are bare
// words, as the BEGIN modes are.
func (p *Parser) parseStartTransactionStmt(start int) (*nodes.BeginStmt, error) {
	// Completion: after START, offer TRANSACTION.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwTRANSACTION)
		return nil, &ParseError{Message: "collecting"}
	}
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
	switch p.cur.Type {
	case kwREAD:
		p.advance()
		// Completion: after READ, offer ONLY and WRITE.
		p.checkCursor()
		if p.collectMode() {
			p.addTokenCandidate(kwONLY)
			p.addTokenCandidate(kwWRITE)
			return nil, &ParseError{Message: "collecting"}
		}
		switch p.cur.Type {
		case kwWRITE:
			p.advance()
			stmt.ReadWrite = true
		case kwONLY:
			p.advance()
			stmt.ReadOnly = true
			if _, ok := p.match(kwAS); ok {
				if _, err := p.expect(kwOF); err != nil {
					return nil, err
				}
				if _, err := p.expect(kwTIMESTAMP); err != nil {
					return nil, err
				}
				ts, err := p.parseExpr()
				if err != nil {
					return nil, err
				}
				if ts == nil {
					return nil, p.syntaxErrorAtCur()
				}
				stmt.AsOf = ts
			}
		default:
			return nil, p.syntaxErrorAtCur()
		}
	case kwWITH:
		p.advance()
		switch {
		case p.cur.Type == kwCONSISTENT:
			p.advance()
			if _, err := p.expect(kwSNAPSHOT); err != nil {
				return nil, err
			}
			stmt.WithConsistentSnapshot = true
		case p.isBareWord("CAUSAL"):
			p.advance()
			if !p.isBareWord("CONSISTENCY") {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
			if _, err := p.expect(kwONLY); err != nil {
				return nil, err
			}
			stmt.CausalConsistencyOnly = true
		default:
			return nil, p.syntaxErrorAtCur()
		}
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseCommitStmt parses a COMMIT statement.
//
//	COMMIT [AND [NO] CHAIN] [[NO] RELEASE]
func (p *Parser) parseCommitStmt() (*nodes.CommitStmt, error) {
	start := p.pos()
	p.advance() // consume COMMIT

	// Completion: after COMMIT, offer AND.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwAND)
		return nil, &ParseError{Message: "collecting"}
	}

	stmt := &nodes.CommitStmt{Loc: nodes.Loc{Start: start}}
	chain, release, err := p.parseCompletionType()
	if err != nil {
		return nil, err
	}
	stmt.Chain, stmt.Release = chain, release
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseRollbackStmt parses a ROLLBACK statement.
//
//	ROLLBACK [AND [NO] CHAIN] [[NO] RELEASE]
//	ROLLBACK TO [SAVEPOINT] identifier
//
// SAVEPOINT is an unreserved keyword in TiDB and so a savepoint name as well:
// in ROLLBACK TO SAVEPOINT it is the keyword when a name follows it and the
// name when nothing does. v8.5.5 parses `ROLLBACK TO SAVEPOINT` (a savepoint
// called savepoint) and `ROLLBACK TO SAVEPOINT savepoint`, and answers 1064
// for `ROLLBACK TO SAVEPOINT savepoint savepoint`. ROLLBACK TO takes no
// completion type after the name.
func (p *Parser) parseRollbackStmt() (*nodes.RollbackStmt, error) {
	start := p.pos()
	p.advance() // consume ROLLBACK

	// Completion: after ROLLBACK, offer TO.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwTO)
		return nil, &ParseError{Message: "collecting"}
	}

	stmt := &nodes.RollbackStmt{Loc: nodes.Loc{Start: start}}

	// TO [SAVEPOINT] identifier
	if _, ok := p.match(kwTO); ok {
		// Completion: after ROLLBACK TO, offer SAVEPOINT keyword.
		p.checkCursor()
		if p.collectMode() {
			p.addTokenCandidate(kwSAVEPOINT)
			return nil, &ParseError{Message: "collecting"}
		}

		if p.cur.Type == kwSAVEPOINT {
			if next := p.peekNext().Type; next != tokEOF && next != ';' {
				p.advance() // consume SAVEPOINT, the keyword
			}
		}
		name, _, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		stmt.Savepoint = name
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

	chain, release, err := p.parseCompletionType()
	if err != nil {
		return nil, err
	}
	stmt.Chain, stmt.Release = chain, release
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseCompletionType reads what may follow COMMIT or ROLLBACK, and reports
// whether the statement chains a new transaction and whether it releases the
// connection.
//
//	[AND [NO] CHAIN] [[NO] RELEASE]
//
// Upstream's CompletionTypeWithinTransaction spells out seven arms, which are
// every combination of the two optional parts but one: AND CHAIN RELEASE,
// which v8.5.5 answers 1064 for, as MySQL does. Each part is whole, so a
// dangling AND or NO is 1064 too.
func (p *Parser) parseCompletionType() (chain, release bool, err error) {
	if _, ok := p.match(kwAND); ok {
		_, no := p.match(kwNO)
		if _, err := p.expect(kwCHAIN); err != nil {
			return false, false, err
		}
		chain = !no
	}
	switch p.cur.Type {
	case kwNO:
		p.advance()
		if _, err := p.expect(kwRELEASE); err != nil {
			return false, false, err
		}
	case kwRELEASE:
		if chain {
			return false, false, p.syntaxErrorAtCur()
		}
		p.advance()
		release = true
	}
	return chain, release, nil
}

// parseSavepointStmt parses a SAVEPOINT statement.
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
//	RELEASE SAVEPOINT identifier
//
// SAVEPOINT is not optional: v8.5.5 answers 1064 for RELEASE s.
func (p *Parser) parseReleaseSavepointStmt() (*nodes.SavepointStmt, error) {
	start := p.pos()
	p.advance() // consume RELEASE

	// Completion: after RELEASE, offer SAVEPOINT.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwSAVEPOINT)
		return nil, &ParseError{Message: "collecting"}
	}
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
//	    tbl_name [[AS] alias] lock_type
//	    [, tbl_name [[AS] alias] lock_type] ...
//
//	lock_type: {
//	    READ [LOCAL]
//	  | [LOW_PRIORITY] WRITE
//	}
func (p *Parser) parseLockTablesStmt(start int) (*nodes.LockTablesStmt, error) {
	// LOCK names the kind of thing it locks before the list, and the keyword is
	// not optional: pingcap/tidb:v8.5.5 answers 1064 for LOCK t READ and for LOCK = t READ. The
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
// Ref: table_or_tables in MySQL 8.0 sql/sql_yacc.yy.
// https://dev.mysql.com/doc/refman/8.0/en/lock-tables.html
//
//	UNLOCK {TABLE | TABLES}
//
// The keyword is not optional, as it is not after LOCK: pingcap/tidb:v8.5.5
// answers 1064 for UNLOCK nosuch and for a bare UNLOCK STATS (read on
// 2026-11-04), and taking the next token without reading it let both through.
func (p *Parser) parseUnlockTablesStmt(start int) (*nodes.UnlockTablesStmt, error) {
	if _, ok := p.match(kwTABLES, kwTABLE); !ok {
		return nil, p.syntaxErrorAtCur()
	}

	return &nodes.UnlockTablesStmt{
		Loc: nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parseLockStatsStmt parses LOCK STATS, which TiDB has and MySQL has not.
//
//	LOCK STATS tbl_name [, tbl_name] ...
//	LOCK STATS tbl_name PARTITION [(] p_name [, p_name] ... [)]
//
// Ref: LockStatsStmt in pingcap/tidb pkg/parser/parser.y.
// https://docs.pingcap.com/tidb/stable/sql-statement-lock-stats/
//
// STATS is a token in TiDB's lexer and is none in this one, so it is read as a
// bare word: v8.5.5 answers 1064 for LOCK `STATS` t, where a backquoted name
// is an ordinary table to lock (read on 2026-11-04).
func (p *Parser) parseLockStatsStmt(start int) (*nodes.LockStatsStmt, error) {
	p.advance() // consume STATS
	tables, partitions, err := p.parseStatsLockTargets()
	if err != nil {
		return nil, err
	}
	return &nodes.LockStatsStmt{
		Loc:        nodes.Loc{Start: start, End: p.prev.End},
		Tables:     tables,
		Partitions: partitions,
	}, nil
}

// parseUnlockStatsStmt parses UNLOCK STATS, which takes what LOCK STATS takes.
//
//	UNLOCK STATS tbl_name [, tbl_name] ...
//	UNLOCK STATS tbl_name PARTITION [(] p_name [, p_name] ... [)]
//
// Ref: UnlockStatsStmt in pingcap/tidb pkg/parser/parser.y.
// https://docs.pingcap.com/tidb/stable/sql-statement-unlock-stats/
func (p *Parser) parseUnlockStatsStmt(start int) (*nodes.UnlockStatsStmt, error) {
	p.advance() // consume STATS
	tables, partitions, err := p.parseStatsLockTargets()
	if err != nil {
		return nil, err
	}
	return &nodes.UnlockStatsStmt{
		Loc:        nodes.Loc{Start: start, End: p.prev.End},
		Tables:     tables,
		Partitions: partitions,
	}, nil
}

// parseStatsLockTargets reads what LOCK STATS and UNLOCK STATS name after the
// STATS word: a list of tables, or one table and a list of its partitions.
//
// Neither list is optional and the partition arm takes one table, which is
// what separates the two arms. Read against pingcap/tidb:v8.5.5 on 2026-11-04:
// a bare LOCK STATS and a bare UNLOCK STATS are both 1064, so is
// `LOCK STATS t, t2 PARTITION (p0)`, and the partition list's parentheses are
// optional -- `PARTITION p0, p1` and `PARTITION (p0, p1)` both parse, while a
// doubled paren, an empty list and a trailing comma are each 1064.
func (p *Parser) parseStatsLockTargets() ([]*nodes.TableRef, []string, error) {
	var tables []*nodes.TableRef
	for {
		// Completion: each position in the list offers a table.
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("table_ref")
			return nil, nil, &ParseError{Message: "collecting"}
		}
		ref, err := p.parseTableRef()
		if err != nil {
			return nil, nil, err
		}
		tables = append(tables, ref)
		if _, ok := p.match(','); !ok {
			break
		}
	}

	if p.cur.Type != kwPARTITION {
		return tables, nil, nil
	}
	if len(tables) != 1 {
		return nil, nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume PARTITION
	_, parens := p.match('(')
	var partitions []string
	for {
		name, _, err := p.parseIdent()
		if err != nil {
			return nil, nil, err
		}
		partitions = append(partitions, name)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	if parens {
		if _, err := p.expect(')'); err != nil {
			return nil, nil, err
		}
	}
	return tables, partitions, nil
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

	stmt.Loc.End = p.prev.End
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

	stmt.Loc.End = p.prev.End
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
