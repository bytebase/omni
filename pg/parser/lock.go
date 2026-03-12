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
func (p *Parser) parseLockStmt() nodes.Node {
	p.advance() // consume LOCK_P

	// opt_table
	if p.cur.Type == TABLE {
		p.advance()
	}

	// relation_expr_list
	relations := p.parseRelationExprList()

	// opt_lock
	mode := p.parseOptLock()

	// opt_nowait
	nowait := p.parseOptNowait()

	return &nodes.LockStmt{
		Relations: relations,
		Mode:      mode,
		Nowait:    nowait,
	}
}

// parseOptLock parses the optional lock mode clause.
//
//	opt_lock:
//	    IN_P lock_type MODE
//	    | /* EMPTY */           -> AccessExclusiveLock
func (p *Parser) parseOptLock() int {
	if p.cur.Type != IN_P {
		return nodes.AccessExclusiveLock
	}
	p.advance() // consume IN

	mode := p.parseLockType()

	p.expect(MODE)

	return mode
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
func (p *Parser) parseLockType() int {
	switch p.cur.Type {
	case ACCESS:
		p.advance() // consume ACCESS
		if p.cur.Type == SHARE {
			p.advance()
			return nodes.AccessShareLock
		}
		// ACCESS EXCLUSIVE
		p.expect(EXCLUSIVE)
		return nodes.AccessExclusiveLock

	case ROW:
		p.advance() // consume ROW
		if p.cur.Type == SHARE {
			p.advance()
			return nodes.RowShareLock
		}
		// ROW EXCLUSIVE
		p.expect(EXCLUSIVE)
		return nodes.RowExclusiveLock

	case SHARE:
		p.advance() // consume SHARE
		if p.cur.Type == UPDATE {
			p.advance() // consume UPDATE
			p.expect(EXCLUSIVE)
			return nodes.ShareUpdateExclusiveLock
		}
		if p.cur.Type == ROW {
			p.advance() // consume ROW
			p.expect(EXCLUSIVE)
			return nodes.ShareRowExclusiveLock
		}
		// plain SHARE
		return nodes.ShareLock

	case EXCLUSIVE:
		p.advance()
		return nodes.ExclusiveLock

	default:
		// Default to AccessExclusiveLock if unrecognized
		return nodes.AccessExclusiveLock
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
