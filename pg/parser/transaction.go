package parser

import (
	nodes "github.com/bytebase/omni/pg/ast"
)

// parseTransactionStmt parses a transaction control statement.
//
// TransactionStmt:
//   ABORT_P opt_transaction opt_transaction_chain
//   | BEGIN_P opt_transaction transaction_mode_list_or_empty
//   | START TRANSACTION transaction_mode_list_or_empty
//   | PREPARE TRANSACTION Sconst
//   | COMMIT PREPARED Sconst
//   | ROLLBACK PREPARED Sconst
//   | COMMIT opt_transaction opt_transaction_chain
//   | END_P opt_transaction opt_transaction_chain
//   | ROLLBACK opt_transaction opt_transaction_chain
//   | ROLLBACK opt_transaction TO [SAVEPOINT] ColId
//   | SAVEPOINT ColId
//   | RELEASE [SAVEPOINT] ColId
func (p *Parser) parseTransactionStmt() nodes.Node {
	switch p.cur.Type {
	case ABORT_P:
		p.advance() // consume ABORT
		p.parseOptTransaction()
		chain := p.parseOptTransactionChain()
		return &nodes.TransactionStmt{
			Kind:  nodes.TRANS_STMT_ROLLBACK,
			Chain: chain,
		}

	case BEGIN_P:
		p.advance() // consume BEGIN
		p.parseOptTransaction()
		options := p.parseTransactionModeListOrEmpty()
		return &nodes.TransactionStmt{
			Kind:    nodes.TRANS_STMT_BEGIN,
			Options: options,
		}

	case START:
		p.advance() // consume START
		p.expect(TRANSACTION)
		options := p.parseTransactionModeListOrEmpty()
		return &nodes.TransactionStmt{
			Kind:    nodes.TRANS_STMT_START,
			Options: options,
		}

	case PREPARE:
		// PREPARE TRANSACTION Sconst
		p.advance() // consume PREPARE
		p.expect(TRANSACTION)
		gid := p.cur.Str
		p.expect(SCONST)
		return &nodes.TransactionStmt{
			Kind: nodes.TRANS_STMT_PREPARE,
			Gid:  gid,
		}

	case COMMIT:
		p.advance() // consume COMMIT
		// COMMIT PREPARED Sconst
		if p.cur.Type == PREPARED {
			p.advance() // consume PREPARED
			gid := p.cur.Str
			p.expect(SCONST)
			return &nodes.TransactionStmt{
				Kind: nodes.TRANS_STMT_COMMIT_PREPARED,
				Gid:  gid,
			}
		}
		// COMMIT opt_transaction opt_transaction_chain
		p.parseOptTransaction()
		chain := p.parseOptTransactionChain()
		return &nodes.TransactionStmt{
			Kind:  nodes.TRANS_STMT_COMMIT,
			Chain: chain,
		}

	case END_P:
		p.advance() // consume END
		p.parseOptTransaction()
		chain := p.parseOptTransactionChain()
		return &nodes.TransactionStmt{
			Kind:  nodes.TRANS_STMT_COMMIT,
			Chain: chain,
		}

	case ROLLBACK:
		p.advance() // consume ROLLBACK
		// ROLLBACK PREPARED Sconst
		if p.cur.Type == PREPARED {
			p.advance() // consume PREPARED
			gid := p.cur.Str
			p.expect(SCONST)
			return &nodes.TransactionStmt{
				Kind: nodes.TRANS_STMT_ROLLBACK_PREPARED,
				Gid:  gid,
			}
		}
		// ROLLBACK opt_transaction TO [SAVEPOINT] ColId
		// ROLLBACK opt_transaction opt_transaction_chain
		p.parseOptTransaction()
		if p.cur.Type == TO {
			p.advance() // consume TO
			if p.cur.Type == SAVEPOINT {
				p.advance() // consume SAVEPOINT
			}
			name, _ := p.parseColId()
			return &nodes.TransactionStmt{
				Kind:      nodes.TRANS_STMT_ROLLBACK_TO,
				Savepoint: name,
			}
		}
		chain := p.parseOptTransactionChain()
		return &nodes.TransactionStmt{
			Kind:  nodes.TRANS_STMT_ROLLBACK,
			Chain: chain,
		}

	case SAVEPOINT:
		p.advance() // consume SAVEPOINT
		name, _ := p.parseColId()
		return &nodes.TransactionStmt{
			Kind:      nodes.TRANS_STMT_SAVEPOINT,
			Savepoint: name,
		}

	case RELEASE:
		p.advance() // consume RELEASE
		if p.cur.Type == SAVEPOINT {
			p.advance() // consume SAVEPOINT
		}
		name, _ := p.parseColId()
		return &nodes.TransactionStmt{
			Kind:      nodes.TRANS_STMT_RELEASE,
			Savepoint: name,
		}

	default:
		return nil
	}
}

// parseOptTransaction consumes the optional WORK or TRANSACTION keyword.
//
// opt_transaction:
//   WORK
//   | TRANSACTION
//   | /* EMPTY */
func (p *Parser) parseOptTransaction() {
	if p.cur.Type == WORK || p.cur.Type == TRANSACTION {
		p.advance()
	}
}

// parseOptTransactionChain parses the optional AND [NO] CHAIN clause.
//
// opt_transaction_chain:
//   AND CHAIN        -> true
//   | AND NO CHAIN   -> false
//   | /* EMPTY */    -> false
func (p *Parser) parseOptTransactionChain() bool {
	if p.cur.Type == AND {
		p.advance() // consume AND
		if p.cur.Type == NO {
			p.advance() // consume NO
			p.expect(CHAIN)
			return false
		}
		p.expect(CHAIN)
		return true
	}
	return false
}

// parseTransactionModeListOrEmpty parses an optional list of transaction mode items.
//
// transaction_mode_list_or_empty:
//   transaction_mode_list
//   | /* EMPTY */
func (p *Parser) parseTransactionModeListOrEmpty() *nodes.List {
	item := p.parseTransactionModeItem()
	if item == nil {
		return nil
	}
	items := []nodes.Node{item}
	for {
		// Items can be separated by commas or just juxtaposed.
		if p.cur.Type == ',' {
			p.advance()
		}
		next := p.parseTransactionModeItem()
		if next == nil {
			break
		}
		items = append(items, next)
	}
	return &nodes.List{Items: items}
}

// parseTransactionModeItem parses a single transaction mode item.
//
// transaction_mode_item:
//   ISOLATION LEVEL iso_level
//   | READ ONLY
//   | READ WRITE
//   | DEFERRABLE
//   | NOT DEFERRABLE
func (p *Parser) parseTransactionModeItem() nodes.Node {
	switch p.cur.Type {
	case ISOLATION:
		p.advance() // consume ISOLATION
		p.expect(LEVEL)
		level := p.parseIsoLevel()
		return makeDefElem("transaction_isolation", &nodes.A_Const{Val: &nodes.String{Str: level}})

	case READ:
		p.advance() // consume READ
		if p.cur.Type == ONLY {
			p.advance() // consume ONLY
			return makeDefElem("transaction_read_only", makeIntConst(1))
		}
		// READ WRITE
		p.expect(WRITE)
		return makeDefElem("transaction_read_only", makeIntConst(0))

	case DEFERRABLE:
		p.advance() // consume DEFERRABLE
		return makeDefElem("transaction_deferrable", makeIntConst(1))

	case NOT:
		// NOT DEFERRABLE — but we need to check that it's actually DEFERRABLE after NOT.
		// Note: NOT may have been reclassified to NOT_LA, but in statement context
		// (not expression context), it stays NOT.
		next := p.peekNext()
		if next.Type == DEFERRABLE {
			p.advance() // consume NOT
			p.advance() // consume DEFERRABLE
			return makeDefElem("transaction_deferrable", makeIntConst(0))
		}
		return nil

	case NOT_LA:
		// NOT_LA DEFERRABLE
		next := p.peekNext()
		if next.Type == DEFERRABLE {
			p.advance() // consume NOT_LA
			p.advance() // consume DEFERRABLE
			return makeDefElem("transaction_deferrable", makeIntConst(0))
		}
		return nil

	default:
		return nil
	}
}

// parseIsoLevel parses the isolation level name.
//
// iso_level:
//   READ UNCOMMITTED   -> "read uncommitted"
//   | READ COMMITTED   -> "read committed"
//   | REPEATABLE READ  -> "repeatable read"
//   | SERIALIZABLE     -> "serializable"
func (p *Parser) parseIsoLevel() string {
	switch p.cur.Type {
	case READ:
		p.advance() // consume READ
		if p.cur.Type == UNCOMMITTED {
			p.advance()
			return "read uncommitted"
		}
		// READ COMMITTED
		p.expect(COMMITTED)
		return "read committed"
	case REPEATABLE:
		p.advance() // consume REPEATABLE
		p.expect(READ)
		return "repeatable read"
	case SERIALIZABLE:
		p.advance()
		return "serializable"
	default:
		return ""
	}
}
