package parser

import (
	"github.com/bytebase/omni/googlesql/ast"
)

// This file is part of the `parser-utility` DAG node. It implements GoogleSQL's
// transaction and batch statements (GoogleSQLParser.g4 §2.9), a hand-port of
// Google's open-source ZetaSQL reference and the grammar bytebase consumes:
//
//	begin_statement:            begin_transaction_keywords transaction_mode_list?
//	begin_transaction_keywords: START TRANSACTION | BEGIN TRANSACTION?
//	commit_statement:           COMMIT TRANSACTION?
//	rollback_statement:         ROLLBACK TRANSACTION?
//	transaction_mode_list:      transaction_mode (',' transaction_mode)*
//	transaction_mode:           READ ONLY | READ WRITE
//	                            | ISOLATION LEVEL id | ISOLATION LEVEL id id
//	start_batch_statement:      START BATCH identifier?
//	run_batch_statement:        RUN BATCH
//	abort_batch_statement:      ABORT BATCH
//
// CORRECTNESS (correctness-protocol.md). The live Cloud Spanner emulator parses
// every one of these (accepts the leading form, then feature-rejects "Statement
// not supported: BeginStatement / CommitStatement / …"), so the leading-form
// ACCEPT is oracle-authoritative (utility_oracle_test.go). The emulator's
// recognizer, however, swallows arbitrary trailing tokens after these keywords
// (it accepts `COMMIT WORK`, `START BATCH a b`, `COMMIT garbage`), so the precise
// trailing grammar is NON-authoritative on Spanner and follows the ZetaSQL .g4:
// `COMMIT TRANSACTION?` rejects `COMMIT WORK`; `START BATCH id?` rejects a second
// word. (Divergence ledger: Spanner's shallow recognizer over-accepts the tails.)

// parseBeginStmt parses begin_statement OR a START-led statement (START
// TRANSACTION / START BATCH). The leading keyword (BEGIN or START) is the
// current token; the disambiguation between transaction and batch is done here.
//
// IMPORTANT (top-level BEGIN). A bare top-level `BEGIN` is begin_statement (a TCL
// transaction); a `BEGIN … END` is a procedural block owned by parser-scripting.
// parseStmt's BEGIN arm disambiguates with the SAME predicate the splitter uses
// (isTCLBeginFollower): a TCL follower (';' / EOF / TRANSACTION / READ /
// ISOLATION) means a transaction, anything else means a block. So this function
// is only reached for the transaction form.
func (p *Parser) parseBeginStmt() (ast.Node, error) {
	begin := p.advance() // BEGIN

	stmt := &ast.TransactionStmt{Kind: ast.TransactionBegin, Loc: begin.Loc}
	// Optional TRANSACTION keyword.
	if _, ok := p.match(kwTRANSACTION); ok {
		stmt.Transaction = true
	}
	stmt.Loc.End = p.prev.Loc.End

	// Optional transaction_mode_list.
	if p.atTransactionModeStart() {
		modes, err := p.parseTransactionModeList()
		if err != nil {
			return nil, err
		}
		stmt.Modes = modes
		stmt.Loc.End = modes[len(modes)-1].Loc.End
	}
	return stmt, nil
}

// parseStartStmt parses a START-led statement: `START TRANSACTION [modes]`
// (begin_statement via begin_transaction_keywords) or `START BATCH [id]`
// (start_batch_statement). START is the current token.
func (p *Parser) parseStartStmt() (ast.Node, error) {
	start := p.advance() // START

	switch p.cur.Type {
	case kwTRANSACTION:
		p.advance() // TRANSACTION
		stmt := &ast.TransactionStmt{Kind: ast.TransactionStart, Transaction: true, Loc: start.Loc}
		stmt.Loc.End = p.prev.Loc.End
		if p.atTransactionModeStart() {
			modes, err := p.parseTransactionModeList()
			if err != nil {
				return nil, err
			}
			stmt.Modes = modes
			stmt.Loc.End = modes[len(modes)-1].Loc.End
		}
		return stmt, nil
	case kwBATCH:
		return p.parseStartBatch(start)
	default:
		return nil, p.syntaxErrorAtCur()
	}
}

// parseStartBatch parses `START BATCH identifier?`. BATCH is the current token;
// start is the consumed START token (for the statement Loc).
func (p *Parser) parseStartBatch(start Token) (ast.Node, error) {
	p.advance() // BATCH
	stmt := &ast.BatchStmt{Kind: ast.BatchStart, Loc: start.Loc}
	stmt.Loc.End = p.prev.Loc.End

	// Optional batch-type identifier (e.g. `ddl`). The grammar allows exactly one;
	// a second word is junk and left for parseSingle's EOF check to reject.
	if isIdentifierStart(p.cur.Type) {
		tok := p.advance()
		name, err := p.identifierText(tok)
		if err != nil {
			return nil, err
		}
		stmt.Name = name
		stmt.Loc.End = tok.Loc.End
	}
	return stmt, nil
}

// parseCommitStmt parses `COMMIT TRANSACTION?`. COMMIT is the current token.
func (p *Parser) parseCommitStmt() (ast.Node, error) {
	commit := p.advance() // COMMIT
	stmt := &ast.TransactionStmt{Kind: ast.TransactionCommit, Loc: commit.Loc}
	if _, ok := p.match(kwTRANSACTION); ok {
		stmt.Transaction = true
	}
	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseRollbackStmt parses `ROLLBACK TRANSACTION?`. ROLLBACK is the current
// token. (There is no SAVEPOINT / ROLLBACK TO in this grammar — a trailing TO is
// left for the EOF check to reject.)
func (p *Parser) parseRollbackStmt() (ast.Node, error) {
	rollback := p.advance() // ROLLBACK
	stmt := &ast.TransactionStmt{Kind: ast.TransactionRollback, Loc: rollback.Loc}
	if _, ok := p.match(kwTRANSACTION); ok {
		stmt.Transaction = true
	}
	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseRunBatchStmt parses `RUN BATCH`. RUN is the current token.
func (p *Parser) parseRunBatchStmt() (ast.Node, error) {
	run := p.advance() // RUN
	if _, err := p.expect(kwBATCH); err != nil {
		return nil, err
	}
	return &ast.BatchStmt{Kind: ast.BatchRun, Loc: ast.Loc{Start: run.Loc.Start, End: p.prev.Loc.End}}, nil
}

// parseAbortBatchStmt parses `ABORT BATCH`. ABORT is the current token.
func (p *Parser) parseAbortBatchStmt() (ast.Node, error) {
	abort := p.advance() // ABORT
	if _, err := p.expect(kwBATCH); err != nil {
		return nil, err
	}
	return &ast.BatchStmt{Kind: ast.BatchAbort, Loc: ast.Loc{Start: abort.Loc.Start, End: p.prev.Loc.End}}, nil
}

// atTransactionModeStart reports whether the current token begins a
// transaction_mode (READ or ISOLATION) — the gate for the optional
// transaction_mode_list after BEGIN / START TRANSACTION.
func (p *Parser) atTransactionModeStart() bool {
	return p.cur.Type == kwREAD || p.cur.Type == kwISOLATION
}

// parseTransactionModeList parses transaction_mode_list: transaction_mode (','
// transaction_mode)*. The first token is a mode start (READ / ISOLATION). A
// trailing comma is a syntax error (the loop requires a mode after each ',').
func (p *Parser) parseTransactionModeList() ([]*ast.TransactionMode, error) {
	var modes []*ast.TransactionMode
	for {
		mode, err := p.parseTransactionMode()
		if err != nil {
			return nil, err
		}
		modes = append(modes, mode)
		if _, ok := p.match(int(',')); ok {
			continue
		}
		break
	}
	return modes, nil
}

// parseTransactionMode parses one transaction_mode:
//
//	READ ONLY | READ WRITE | ISOLATION LEVEL id | ISOLATION LEVEL id id
func (p *Parser) parseTransactionMode() (*ast.TransactionMode, error) {
	switch p.cur.Type {
	case kwREAD:
		read := p.advance() // READ
		switch p.cur.Type {
		case kwONLY:
			only := p.advance()
			return &ast.TransactionMode{Kind: ast.TransactionModeReadOnly, Loc: ast.Loc{Start: read.Loc.Start, End: only.Loc.End}}, nil
		case kwWRITE:
			write := p.advance()
			return &ast.TransactionMode{Kind: ast.TransactionModeReadWrite, Loc: ast.Loc{Start: read.Loc.Start, End: write.Loc.End}}, nil
		default:
			return nil, p.syntaxErrorAtCur()
		}
	case kwISOLATION:
		iso := p.advance() // ISOLATION
		if _, err := p.expect(kwLEVEL); err != nil {
			return nil, err
		}
		mode := &ast.TransactionMode{Kind: ast.TransactionModeIsolationLevel, Loc: iso.Loc}
		// First isolation-level word (required).
		first, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}
		firstName, err := p.identifierText(first)
		if err != nil {
			return nil, err
		}
		mode.Levels = append(mode.Levels, firstName)
		mode.Loc.End = first.Loc.End
		// Optional second isolation-level word (e.g. `READ COMMITTED`). It is taken
		// ONLY when it does not start the next transaction_mode after a ','/EOF — a
		// bare identifier here is always part of THIS isolation level (the grammar's
		// `ISOLATION LEVEL id id` alternative), since a transaction_mode never begins
		// with a bare identifier (only READ / ISOLATION).
		if isIdentifierStart(p.cur.Type) {
			second := p.advance()
			secondName, err := p.identifierText(second)
			if err != nil {
				return nil, err
			}
			mode.Levels = append(mode.Levels, secondName)
			mode.Loc.End = second.Loc.End
		}
		return mode, nil
	default:
		return nil, p.syntaxErrorAtCur()
	}
}
