package parser

import nodes "github.com/bytebase/omni/pg/ast"

// parseLockStmt parses a LOCK TABLE statement.
//
// LockStmt:
//
//	LOCK_P opt_table relation_expr_list opt_lock opt_nowait
//
// opt_table:
//
//	TABLE
//	| /* EMPTY */
//
// opt_lock:
//
//	IN_P lock_type MODE
//	| /* EMPTY */           -> AccessExclusiveLock
//
// lock_type:
//
//	ACCESS SHARE
//	| ROW SHARE
//	| ROW EXCLUSIVE
//	| SHARE UPDATE EXCLUSIVE
//	| SHARE
//	| SHARE ROW EXCLUSIVE
//	| EXCLUSIVE
//	| ACCESS EXCLUSIVE
//
// opt_nowait:
//
//	NOWAIT  -> true
//	| /* EMPTY */ -> false
func (p *Parser) parseLockStmt() (nodes.Node, error) {
	loc := p.pos()
	p.advance() // consume LOCK_P

	// opt_table
	if p.cur.Type == TABLE {
		p.advance()
	}

	// relation_expr_list
	relations := p.parseRelationExprList()

	// opt_lock
	mode, err := p.parseOptLock()
	if err != nil {
		return nil, err
	}

	// opt_nowait
	nowait := p.parseOptNowait()

	return &nodes.LockStmt{
		Relations: relations,
		Mode:      mode,
		Nowait:    nowait,
		Loc:       nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

// parseOptLock parses the optional lock mode clause.
//
//	opt_lock:
//	    IN_P lock_type MODE
//	    | /* EMPTY */           -> AccessExclusiveLock
func (p *Parser) parseOptLock() (int, error) {
	if p.cur.Type != IN_P {
		return nodes.AccessExclusiveLock, nil
	}
	p.advance() // consume IN

	mode, err := p.parseLockType()
	if err != nil {
		return 0, err
	}

	if _, err := p.expect(MODE); err != nil {
		return 0, err
	}

	return mode, nil
}

// parseLockType parses the lock type keyword(s).
//
//	lock_type:
//	    ACCESS SHARE
//	    | ROW SHARE
//	    | ROW EXCLUSIVE
//	    | SHARE UPDATE EXCLUSIVE
//	    | SHARE
//	    | SHARE ROW EXCLUSIVE
//	    | EXCLUSIVE
//	    | ACCESS EXCLUSIVE
func (p *Parser) parseLockType() (int, error) {
	switch p.cur.Type {
	case ACCESS:
		p.advance() // consume ACCESS
		if p.cur.Type == SHARE {
			p.advance()
			return nodes.AccessShareLock, nil
		}
		// ACCESS EXCLUSIVE
		if _, err := p.expect(EXCLUSIVE); err != nil {
			return 0, err
		}
		return nodes.AccessExclusiveLock, nil

	case ROW:
		p.advance() // consume ROW
		if p.cur.Type == SHARE {
			p.advance()
			return nodes.RowShareLock, nil
		}
		// ROW EXCLUSIVE
		if _, err := p.expect(EXCLUSIVE); err != nil {
			return 0, err
		}
		return nodes.RowExclusiveLock, nil

	case SHARE:
		p.advance() // consume SHARE
		if p.cur.Type == UPDATE {
			p.advance() // consume UPDATE
			if _, err := p.expect(EXCLUSIVE); err != nil {
				return 0, err
			}
			return nodes.ShareUpdateExclusiveLock, nil
		}
		if p.cur.Type == ROW {
			p.advance() // consume ROW
			if _, err := p.expect(EXCLUSIVE); err != nil {
				return 0, err
			}
			return nodes.ShareRowExclusiveLock, nil
		}
		// plain SHARE
		return nodes.ShareLock, nil

	case EXCLUSIVE:
		p.advance()
		return nodes.ExclusiveLock, nil

	default:
		// Default to AccessExclusiveLock if unrecognized
		return nodes.AccessExclusiveLock, nil
	}
}

// parseOptNowait parses the optional NOWAIT keyword.
//
//	opt_nowait:
//	    NOWAIT      -> true
//	    | /* EMPTY */ -> false
func (p *Parser) parseOptNowait() bool {
	if p.cur.Type == NOWAIT {
		p.advance()
		return true
	}
	return false
}
