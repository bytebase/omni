package parser

import (
	"github.com/bytebase/omni/trino/ast"
)

// This file is part of the parser-dcl-tcl DAG node: it implements Trino's
// transaction-control (TCL) statements — START TRANSACTION, COMMIT, and
// ROLLBACK — as hand-written recursive-descent parsers over the token stream.
//
// The statement nodes are returned from parseStmt as ast.Node and so carry an
// ast.NodeTag (declared in trino/ast/nodetags.go); their concrete fields live
// here in package parser, matching the Trino convention for parser-built node
// types (Expr in expr.go, DataType in datatypes.go).
//
// Legacy ANTLR grammar (TrinoParser.g4) the implementation tracks:
//
//	statement
//	    : ...
//	    | START_ TRANSACTION_ (transactionMode (COMMA_ transactionMode)*)? # startTransaction
//	    | COMMIT_ WORK_?                                                   # commit
//	    | ROLLBACK_ WORK_?                                                 # rollback
//	    ;
//	transactionMode
//	    : ISOLATION_ LEVEL_ levelOfIsolation  # isolationLevel
//	    | READ_ accessMode = (ONLY_ | WRITE_) # transactionAccessMode
//	    ;
//	levelOfIsolation
//	    : READ_ UNCOMMITTED_ # readUncommitted
//	    | READ_ COMMITTED_   # readCommitted
//	    | REPEATABLE_ READ_  # repeatableRead
//	    | SERIALIZABLE_      # serializable
//	    ;
//
// Trino 481 docs (truth1) confirm the same surface, including combining
// multiple comma-separated modes in one statement, e.g.
//
//	START TRANSACTION ISOLATION LEVEL READ COMMITTED, READ ONLY;
//	START TRANSACTION READ WRITE, ISOLATION LEVEL SERIALIZABLE;
//
// The implementation is adjudicated against the live Trino 481 oracle.

// ---------------------------------------------------------------------------
// Transaction-mode AST
// ---------------------------------------------------------------------------

// TxnIsolationLevel enumerates the four SQL isolation levels Trino accepts in
// an ISOLATION LEVEL transaction mode.
type TxnIsolationLevel int

const (
	// IsolationReadUncommitted is READ UNCOMMITTED.
	IsolationReadUncommitted TxnIsolationLevel = iota
	// IsolationReadCommitted is READ COMMITTED.
	IsolationReadCommitted
	// IsolationRepeatableRead is REPEATABLE READ.
	IsolationRepeatableRead
	// IsolationSerializable is SERIALIZABLE.
	IsolationSerializable
)

// String renders the isolation level in its canonical SQL spelling.
func (l TxnIsolationLevel) String() string {
	switch l {
	case IsolationReadUncommitted:
		return "READ UNCOMMITTED"
	case IsolationReadCommitted:
		return "READ COMMITTED"
	case IsolationRepeatableRead:
		return "REPEATABLE READ"
	case IsolationSerializable:
		return "SERIALIZABLE"
	default:
		return "UNKNOWN"
	}
}

// TransactionMode is one mode of a START TRANSACTION statement: either an
// ISOLATION LEVEL clause or a READ ONLY / READ WRITE access mode. Exactly one
// of the two flavors is active, selected by IsAccessMode.
type TransactionMode struct {
	// IsAccessMode distinguishes the two transactionMode alternatives: false
	// for ISOLATION LEVEL (Isolation is meaningful), true for READ ONLY /
	// READ WRITE (ReadOnly is meaningful).
	IsAccessMode bool
	// Isolation is the isolation level when IsAccessMode is false.
	Isolation TxnIsolationLevel
	// ReadOnly is true for READ ONLY, false for READ WRITE, when
	// IsAccessMode is true.
	ReadOnly bool
	Loc      ast.Loc
}

// ---------------------------------------------------------------------------
// Statement AST
// ---------------------------------------------------------------------------

// StartTransactionStmt is a START TRANSACTION statement with zero or more
// comma-separated transaction modes.
type StartTransactionStmt struct {
	Modes []*TransactionMode
	Loc   ast.Loc
}

// Tag implements ast.Node.
func (n *StartTransactionStmt) Tag() ast.NodeTag { return ast.T_StartTransactionStmt }

// CommitStmt is a COMMIT [WORK] statement. Work records whether the optional,
// no-op WORK keyword was present (preserved for source-faithful deparse).
type CommitStmt struct {
	Work bool
	Loc  ast.Loc
}

// Tag implements ast.Node.
func (n *CommitStmt) Tag() ast.NodeTag { return ast.T_CommitStmt }

// RollbackStmt is a ROLLBACK [WORK] statement. Work records whether the
// optional, no-op WORK keyword was present.
type RollbackStmt struct {
	Work bool
	Loc  ast.Loc
}

// Tag implements ast.Node.
func (n *RollbackStmt) Tag() ast.NodeTag { return ast.T_RollbackStmt }

// Compile-time assertions that the statement types satisfy ast.Node.
var (
	_ ast.Node = (*StartTransactionStmt)(nil)
	_ ast.Node = (*CommitStmt)(nil)
	_ ast.Node = (*RollbackStmt)(nil)
)

// ---------------------------------------------------------------------------
// Parsers
// ---------------------------------------------------------------------------

// parseStartTransactionStmt parses START TRANSACTION [mode [, ...]].
// On entry cur is the START keyword (not yet consumed).
func (p *Parser) parseStartTransactionStmt() (ast.Node, error) {
	startTok := p.advance() // consume START
	if _, err := p.expect(kwTRANSACTION); err != nil {
		return nil, err
	}

	stmt := &StartTransactionStmt{Loc: ast.Loc{Start: startTok.Loc.Start, End: p.prev.Loc.End}}

	// Optional mode list. A mode begins with ISOLATION or READ; anything else
	// terminates the statement (EOF / ';'). The grammar makes the whole list
	// optional, so a bare START TRANSACTION is valid.
	if p.cur.Kind == kwISOLATION || p.cur.Kind == kwREAD {
		for {
			mode, err := p.parseTransactionMode()
			if err != nil {
				return nil, err
			}
			stmt.Modes = append(stmt.Modes, mode)
			if _, ok := p.match(int(',')); !ok {
				break
			}
		}
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseTransactionMode parses one transactionMode:
//
//	ISOLATION LEVEL <levelOfIsolation>
//	READ ( ONLY | WRITE )
//
// On entry cur is ISOLATION or READ.
func (p *Parser) parseTransactionMode() (*TransactionMode, error) {
	startLoc := p.cur.Loc
	switch p.cur.Kind {
	case kwISOLATION:
		p.advance() // consume ISOLATION
		if _, err := p.expect(kwLEVEL); err != nil {
			return nil, err
		}
		level, err := p.parseIsolationLevel()
		if err != nil {
			return nil, err
		}
		return &TransactionMode{
			IsAccessMode: false,
			Isolation:    level,
			Loc:          ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End},
		}, nil

	case kwREAD:
		p.advance() // consume READ
		switch p.cur.Kind {
		case kwONLY:
			p.advance() // consume ONLY
			return &TransactionMode{
				IsAccessMode: true,
				ReadOnly:     true,
				Loc:          ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End},
			}, nil
		case kwWRITE:
			p.advance() // consume WRITE
			return &TransactionMode{
				IsAccessMode: true,
				ReadOnly:     false,
				Loc:          ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End},
			}, nil
		default:
			return nil, p.syntaxErrorAtCur()
		}

	default:
		return nil, p.syntaxErrorAtCur()
	}
}

// parseIsolationLevel parses a levelOfIsolation:
//
//	READ UNCOMMITTED | READ COMMITTED | REPEATABLE READ | SERIALIZABLE
//
// On entry cur is READ, REPEATABLE, or SERIALIZABLE (the LEVEL keyword has
// already been consumed).
func (p *Parser) parseIsolationLevel() (TxnIsolationLevel, error) {
	switch p.cur.Kind {
	case kwREAD:
		p.advance() // consume READ
		switch p.cur.Kind {
		case kwUNCOMMITTED:
			p.advance()
			return IsolationReadUncommitted, nil
		case kwCOMMITTED:
			p.advance()
			return IsolationReadCommitted, nil
		default:
			return 0, p.syntaxErrorAtCur()
		}
	case kwREPEATABLE:
		p.advance() // consume REPEATABLE
		if _, err := p.expect(kwREAD); err != nil {
			return 0, err
		}
		return IsolationRepeatableRead, nil
	case kwSERIALIZABLE:
		p.advance() // consume SERIALIZABLE
		return IsolationSerializable, nil
	default:
		return 0, p.syntaxErrorAtCur()
	}
}

// parseCommitStmt parses COMMIT [WORK]. On entry cur is the COMMIT keyword.
func (p *Parser) parseCommitStmt() (ast.Node, error) {
	commitTok := p.advance() // consume COMMIT
	stmt := &CommitStmt{Loc: ast.Loc{Start: commitTok.Loc.Start, End: commitTok.Loc.End}}
	if _, ok := p.match(kwWORK); ok {
		stmt.Work = true
		stmt.Loc.End = p.prev.Loc.End
	}
	return stmt, nil
}

// parseRollbackStmt parses ROLLBACK [WORK]. On entry cur is the ROLLBACK
// keyword.
//
// Note Trino's ROLLBACK statement is transaction rollback only; it has no
// SAVEPOINT / TO form (the grammar's rollback alternative is just
// `ROLLBACK_ WORK_?`).
func (p *Parser) parseRollbackStmt() (ast.Node, error) {
	rollbackTok := p.advance() // consume ROLLBACK
	stmt := &RollbackStmt{Loc: ast.Loc{Start: rollbackTok.Loc.Start, End: rollbackTok.Loc.End}}
	if _, ok := p.match(kwWORK); ok {
		stmt.Work = true
		stmt.Loc.End = p.prev.Loc.End
	}
	return stmt, nil
}
