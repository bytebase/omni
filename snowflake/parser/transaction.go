package parser

import "github.com/bytebase/omni/snowflake/ast"

// ---------------------------------------------------------------------------
// TCL (transaction control) statements — T6.2
//
// Grammar (docs.snowflake.com + legacy SnowflakeParser.g4, which agree):
//
//	BEGIN [ { WORK | TRANSACTION } ] [ NAME <name> ]
//	START TRANSACTION [ NAME <name> ]
//	COMMIT [ WORK ]
//	ROLLBACK [ WORK ]
//
// START TRANSACTION is a documented synonym for BEGIN. WORK after COMMIT /
// ROLLBACK (and the WORK / TRANSACTION modifier after BEGIN) is a no-op kept
// for cross-database compatibility. Snowflake has no SAVEPOINT, RELEASE
// SAVEPOINT, or ROLLBACK TO — those keywords are absent from both the docs and
// the legacy grammar, so they are not accepted here.
//
// Note on trailing tokens: like the rest of omni's snowflake parser (e.g.
// parseDropStmt, parseTruncateStmt), these functions consume exactly the
// grammar they own and do not reject extra tokens that follow a complete
// statement — the F3 splitter is trusted to segment input, and trailing-token
// strictness is an engine-wide concern, not enforced per node.
// ---------------------------------------------------------------------------

// parseBeginStmt parses the transaction-opening statement in its BEGIN form:
//
//	BEGIN [ { WORK | TRANSACTION } ] [ NAME <name> ]
//
// The BEGIN keyword has NOT yet been consumed when this function is called.
func (p *Parser) parseBeginStmt() (ast.Node, error) {
	beginTok := p.advance() // consume BEGIN
	stmt := &ast.BeginStmt{
		Kind: ast.BeginBare,
		Loc:  ast.Loc{Start: beginTok.Loc.Start, End: beginTok.Loc.End},
	}

	// Optional { WORK | TRANSACTION } modifier (at most one).
	switch p.cur.Type {
	case kwWORK:
		p.advance()
		stmt.Kind = ast.BeginWork
	case kwTRANSACTION:
		p.advance()
		stmt.Kind = ast.BeginTransaction
	}

	// Optional NAME <name>.
	if err := p.parseOptionalTxnName(stmt); err != nil {
		return nil, err
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseStartTransactionStmt parses the transaction-opening statement in its
// START TRANSACTION form:
//
//	START TRANSACTION [ NAME <name> ]
//
// The START keyword has NOT yet been consumed when this function is called.
// TRANSACTION is required after START (START WORK and a bare START are not
// valid transaction openers).
func (p *Parser) parseStartTransactionStmt() (ast.Node, error) {
	startTok := p.advance() // consume START
	stmt := &ast.BeginStmt{
		Kind: ast.BeginStartTransaction,
		Loc:  ast.Loc{Start: startTok.Loc.Start, End: startTok.Loc.End},
	}

	if _, err := p.expect(kwTRANSACTION); err != nil {
		return nil, err
	}

	// Optional NAME <name>.
	if err := p.parseOptionalTxnName(stmt); err != nil {
		return nil, err
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseOptionalTxnName consumes an optional `NAME <name>` clause and stores the
// identifier on stmt. NAME is a non-reserved keyword (kwNAME); when present it
// MUST be followed by an identifier (a missing identifier is a parse error).
func (p *Parser) parseOptionalTxnName(stmt *ast.BeginStmt) error {
	if p.cur.Type != kwNAME {
		return nil
	}
	p.advance() // consume NAME
	name, err := p.parseIdent()
	if err != nil {
		return err
	}
	stmt.Name = name
	return nil
}

// parseCommitStmt parses `COMMIT [ WORK ]`. The COMMIT keyword has NOT yet been
// consumed when this function is called.
func (p *Parser) parseCommitStmt() (ast.Node, error) {
	commitTok := p.advance() // consume COMMIT
	stmt := &ast.CommitStmt{
		Loc: ast.Loc{Start: commitTok.Loc.Start, End: commitTok.Loc.End},
	}

	if p.cur.Type == kwWORK {
		p.advance()
		stmt.Work = true
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseRollbackStmt parses `ROLLBACK [ WORK ]`. The ROLLBACK keyword has NOT
// yet been consumed when this function is called.
func (p *Parser) parseRollbackStmt() (ast.Node, error) {
	rollbackTok := p.advance() // consume ROLLBACK
	stmt := &ast.RollbackStmt{
		Loc: ast.Loc{Start: rollbackTok.Loc.Start, End: rollbackTok.Loc.End},
	}

	if p.cur.Type == kwWORK {
		p.advance()
		stmt.Work = true
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}
