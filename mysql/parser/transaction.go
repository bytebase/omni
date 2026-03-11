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
		if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "work") {
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

	stmt := &nodes.CommitStmt{Loc: nodes.Loc{Start: start}}

	// Optional WORK
	if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "work") {
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

	stmt := &nodes.RollbackStmt{Loc: nodes.Loc{Start: start}}

	// Optional WORK
	if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "work") {
		p.advance()
	}

	// TO [SAVEPOINT] identifier
	if _, ok := p.match(kwTO); ok {
		p.match(kwSAVEPOINT)
		name, _, err := p.parseIdentifier()
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

	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	return &nodes.SavepointStmt{
		Loc:  nodes.Loc{Start: start, End: p.pos()},
		Name: name,
	}, nil
}

// parseReleaseSavepointStmt parses RELEASE SAVEPOINT.
func (p *Parser) parseReleaseSavepointStmt() (*nodes.SavepointStmt, error) {
	start := p.pos()
	p.advance() // consume RELEASE
	p.match(kwSAVEPOINT)

	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	return &nodes.SavepointStmt{
		Loc:  nodes.Loc{Start: start, End: p.pos()},
		Name: name,
	}, nil
}

// parseLockTablesStmt parses a LOCK TABLES statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/lock-tables.html
//
//	LOCK TABLES tbl_name [[AS] alias] lock_type [, tbl_name [[AS] alias] lock_type] ...
//	lock_type: READ [LOCAL] | [LOW_PRIORITY] WRITE
func (p *Parser) parseLockTablesStmt() (*nodes.LockTablesStmt, error) {
	start := p.pos()
	p.advance() // consume TABLES

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
//	UNLOCK TABLES
func (p *Parser) parseUnlockTablesStmt() (*nodes.UnlockTablesStmt, error) {
	start := p.pos()
	p.advance() // consume TABLES

	return &nodes.UnlockTablesStmt{
		Loc: nodes.Loc{Start: start, End: p.pos()},
	}, nil
}
