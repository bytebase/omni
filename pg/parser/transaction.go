package parser

import (
	nodes "github.com/bytebase/omni/pg/ast"
)

// parseTransactionStmt parses a transaction control statement.
//
// TransactionStmt:
//
//	ABORT_P opt_transaction opt_transaction_chain
//	| BEGIN_P opt_transaction transaction_mode_list_or_empty
//	| START TRANSACTION transaction_mode_list_or_empty
//	| PREPARE TRANSACTION Sconst
//	| COMMIT PREPARED Sconst
//	| ROLLBACK PREPARED Sconst
//	| COMMIT opt_transaction opt_transaction_chain
//	| END_P opt_transaction opt_transaction_chain
//	| ROLLBACK opt_transaction opt_transaction_chain
//	| ROLLBACK opt_transaction TO [SAVEPOINT] ColId
//	| SAVEPOINT ColId
//	| RELEASE [SAVEPOINT] ColId
func (p *Parser) parseTransactionStmt() (nodes.Node, error) {
	loc := p.pos()
	n, err := p.parseTransactionStmtInner()
	if err != nil {
		return nil, err
	}
	if ts, ok := n.(*nodes.TransactionStmt); ok {
		ts.Loc = nodes.Loc{Start: loc, End: p.pos()}
	}
	return n, nil
}

func (p *Parser) parseTransactionStmtInner() (nodes.Node, error) {
	switch p.cur.Type {
	case ABORT_P:
		p.advance() // consume ABORT
		p.parseOptTransaction()
		chain, err := p.parseOptTransactionChain()
		if err != nil {
			return nil, err
		}
		return &nodes.TransactionStmt{
			Kind:  nodes.TRANS_STMT_ROLLBACK,
			Chain: chain,
		}, nil

	case BEGIN_P:
		p.advance() // consume BEGIN
		p.parseOptTransaction()
		options, err := p.parseTransactionModeListOrEmpty()
		if err != nil {
			return nil, err
		}
		return &nodes.TransactionStmt{
			Kind:    nodes.TRANS_STMT_BEGIN,
			Options: options,
		}, nil

	case START:
		p.advance() // consume START
		if _, err := p.expect(TRANSACTION); err != nil {
			return nil, err
		}
		options, err := p.parseTransactionModeListOrEmpty()
		if err != nil {
			return nil, err
		}
		return &nodes.TransactionStmt{
			Kind:    nodes.TRANS_STMT_START,
			Options: options,
		}, nil

	case PREPARE:
		// PREPARE TRANSACTION Sconst
		p.advance() // consume PREPARE
		if _, err := p.expect(TRANSACTION); err != nil {
			return nil, err
		}
		gid := p.cur.Str
		if _, err := p.expect(SCONST); err != nil {
			return nil, err
		}
		return &nodes.TransactionStmt{
			Kind: nodes.TRANS_STMT_PREPARE,
			Gid:  gid,
		}, nil

	case COMMIT:
		p.advance() // consume COMMIT
		// COMMIT PREPARED Sconst
		if p.cur.Type == PREPARED {
			p.advance() // consume PREPARED
			gid := p.cur.Str
			if _, err := p.expect(SCONST); err != nil {
				return nil, err
			}
			return &nodes.TransactionStmt{
				Kind: nodes.TRANS_STMT_COMMIT_PREPARED,
				Gid:  gid,
			}, nil
		}
		// COMMIT opt_transaction opt_transaction_chain
		p.parseOptTransaction()
		chain, err := p.parseOptTransactionChain()
		if err != nil {
			return nil, err
		}
		return &nodes.TransactionStmt{
			Kind:  nodes.TRANS_STMT_COMMIT,
			Chain: chain,
		}, nil

	case END_P:
		p.advance() // consume END
		p.parseOptTransaction()
		chain, err := p.parseOptTransactionChain()
		if err != nil {
			return nil, err
		}
		return &nodes.TransactionStmt{
			Kind:  nodes.TRANS_STMT_COMMIT,
			Chain: chain,
		}, nil

	case ROLLBACK:
		p.advance() // consume ROLLBACK
		// ROLLBACK PREPARED Sconst
		if p.cur.Type == PREPARED {
			p.advance() // consume PREPARED
			gid := p.cur.Str
			if _, err := p.expect(SCONST); err != nil {
				return nil, err
			}
			return &nodes.TransactionStmt{
				Kind: nodes.TRANS_STMT_ROLLBACK_PREPARED,
				Gid:  gid,
			}, nil
		}
		// ROLLBACK opt_transaction TO [SAVEPOINT] ColId
		// ROLLBACK opt_transaction opt_transaction_chain
		p.parseOptTransaction()
		if p.cur.Type == TO {
			p.advance() // consume TO
			if p.cur.Type == SAVEPOINT {
				p.advance() // consume SAVEPOINT
			}
			name, err := p.parseColId()
			if err != nil {
				return nil, err
			}
			return &nodes.TransactionStmt{
				Kind:      nodes.TRANS_STMT_ROLLBACK_TO,
				Savepoint: name,
			}, nil
		}
		chain, err := p.parseOptTransactionChain()
		if err != nil {
			return nil, err
		}
		return &nodes.TransactionStmt{
			Kind:  nodes.TRANS_STMT_ROLLBACK,
			Chain: chain,
		}, nil

	case SAVEPOINT:
		p.advance() // consume SAVEPOINT
		name, err := p.parseColId()
		if err != nil {
			return nil, err
		}
		return &nodes.TransactionStmt{
			Kind:      nodes.TRANS_STMT_SAVEPOINT,
			Savepoint: name,
		}, nil

	case RELEASE:
		p.advance() // consume RELEASE
		if p.cur.Type == SAVEPOINT {
			p.advance() // consume SAVEPOINT
		}
		name, err := p.parseColId()
		if err != nil {
			return nil, err
		}
		return &nodes.TransactionStmt{
			Kind:      nodes.TRANS_STMT_RELEASE,
			Savepoint: name,
		}, nil

	default:
		return nil, nil
	}
}

// parseOptTransaction consumes the optional WORK or TRANSACTION keyword.
//
// opt_transaction:
//
//	WORK
//	| TRANSACTION
//	| /* EMPTY */
func (p *Parser) parseOptTransaction() {
	if p.cur.Type == WORK || p.cur.Type == TRANSACTION {
		p.advance()
	}
}

// parseOptTransactionChain parses the optional AND [NO] CHAIN clause.
//
// opt_transaction_chain:
//
//	AND CHAIN        -> true
//	| AND NO CHAIN   -> false
//	| /* EMPTY */    -> false
func (p *Parser) parseOptTransactionChain() (bool, error) {
	if p.cur.Type == AND {
		p.advance() // consume AND
		if p.cur.Type == NO {
			p.advance() // consume NO
			if _, err := p.expect(CHAIN); err != nil {
				return false, err
			}
			return false, nil
		}
		if _, err := p.expect(CHAIN); err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

// parseTransactionModeListOrEmpty parses an optional list of transaction mode items.
//
// transaction_mode_list_or_empty:
//
//	transaction_mode_list
//	| /* EMPTY */
func (p *Parser) parseTransactionModeListOrEmpty() (*nodes.List, error) {
	item, err := p.parseTransactionModeItem()
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, nil
	}
	items := []nodes.Node{item}
	for {
		// Items can be separated by commas or just juxtaposed.
		if p.cur.Type == ',' {
			p.advance()
		}
		next, err := p.parseTransactionModeItem()
		if err != nil {
			return nil, err
		}
		if next == nil {
			break
		}
		items = append(items, next)
	}
	return &nodes.List{Items: items}, nil
}

// parseTransactionModeItem parses a single transaction mode item.
//
// transaction_mode_item:
//
//	ISOLATION LEVEL iso_level
//	| READ ONLY
//	| READ WRITE
//	| DEFERRABLE
//	| NOT DEFERRABLE
func (p *Parser) parseTransactionModeItem() (nodes.Node, error) {
	switch p.cur.Type {
	case ISOLATION:
		p.advance() // consume ISOLATION
		if _, err := p.expect(LEVEL); err != nil {
			return nil, err
		}
		level, err := p.parseIsoLevel()
		if err != nil {
			return nil, err
		}
		return makeDefElem("transaction_isolation", &nodes.A_Const{Val: &nodes.String{Str: level}}), nil

	case READ:
		p.advance() // consume READ
		if p.cur.Type == ONLY {
			p.advance() // consume ONLY
			return makeDefElem("transaction_read_only", makeIntConst(1)), nil
		}
		// READ WRITE
		if _, err := p.expect(WRITE); err != nil {
			return nil, err
		}
		return makeDefElem("transaction_read_only", makeIntConst(0)), nil

	case DEFERRABLE:
		p.advance() // consume DEFERRABLE
		return makeDefElem("transaction_deferrable", makeIntConst(1)), nil

	case NOT:
		// NOT DEFERRABLE — but we need to check that it's actually DEFERRABLE after NOT.
		// Note: NOT may have been reclassified to NOT_LA, but in statement context
		// (not expression context), it stays NOT.
		next := p.peekNext()
		if next.Type == DEFERRABLE {
			p.advance() // consume NOT
			p.advance() // consume DEFERRABLE
			return makeDefElem("transaction_deferrable", makeIntConst(0)), nil
		}
		return nil, nil

	case NOT_LA:
		// NOT_LA DEFERRABLE
		next := p.peekNext()
		if next.Type == DEFERRABLE {
			p.advance() // consume NOT_LA
			p.advance() // consume DEFERRABLE
			return makeDefElem("transaction_deferrable", makeIntConst(0)), nil
		}
		return nil, nil

	default:
		return nil, nil
	}
}

// parseIsoLevel parses the isolation level name.
//
// iso_level:
//
//	READ UNCOMMITTED   -> "read uncommitted"
//	| READ COMMITTED   -> "read committed"
//	| REPEATABLE READ  -> "repeatable read"
//	| SERIALIZABLE     -> "serializable"
func (p *Parser) parseIsoLevel() (string, error) {
	switch p.cur.Type {
	case READ:
		p.advance() // consume READ
		if p.cur.Type == UNCOMMITTED {
			p.advance()
			return "read uncommitted", nil
		}
		// READ COMMITTED
		if _, err := p.expect(COMMITTED); err != nil {
			return "", err
		}
		return "read committed", nil
	case REPEATABLE:
		p.advance() // consume REPEATABLE
		if _, err := p.expect(READ); err != nil {
			return "", err
		}
		return "repeatable read", nil
	case SERIALIZABLE:
		p.advance()
		return "serializable", nil
	default:
		return "", nil
	}
}
