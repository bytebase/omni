// Package parser implements a recursive descent SQL parser for PostgreSQL.
//
// This parser reuses the lexer and keyword definitions from this package
// and produces AST nodes from the pg/ast package.
package parser

import (
	"fmt"

	nodes "github.com/bytebase/omni/pg/ast"
)

// Parser is a recursive descent parser for PostgreSQL SQL.
type Parser struct {
	lexer   *Lexer
	source  string
	cur     Token // current token
	prev    Token // previous token (for error reporting)
	nextBuf Token // buffered next token for 2-token lookahead
	hasNext bool  // whether nextBuf is valid

	// Completion mode fields.
	completing   bool          // true when collecting completion candidates
	cursorOff    int           // byte offset of the cursor in source
	candidates   *CandidateSet // collected candidates
	collecting   bool          // true once cursor position is reached
	collectDepth int           // recursion depth in collect mode
	maxCollect   int           // max exploration depth
	stmtStart    int           // byte offset of the current statement's first token
}

// Parse parses a SQL string into an AST list.
// Currently supports basic SELECT statements for expression testing.
// Full statement dispatch will be implemented in batch 34.
func Parse(sql string) (*nodes.List, error) {
	p := &Parser{
		lexer:  NewLexer(sql),
		source: sql,
	}
	p.advance()

	var stmts []nodes.Node
	for p.cur.Type != 0 {
		// Skip semicolons
		if p.cur.Type == ';' {
			p.advance()
			continue
		}
		if p.lexer.Err != nil {
			return nil, p.lexerError()
		}
		stmtStart := p.pos()
		p.stmtStart = stmtStart
		stmt, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		if stmt == nil {
			if p.cur.Type != 0 {
				return nil, p.syntaxErrorAtCur()
			}
			break
		}
		if p.lexer.Err != nil {
			return nil, p.lexerError()
		}
		raw := &nodes.RawStmt{
			Stmt: stmt,
			Loc:  nodes.Loc{Start: stmtStart, End: p.prev.End},
		}
		stmts = append(stmts, raw)
	}

	if len(stmts) == 0 {
		return &nodes.List{}, nil
	}
	return &nodes.List{Items: stmts}, nil
}

// parseStmt dispatches to statement-specific parsers.
// Minimal implementation — SELECT/VALUES/TABLE/WITH/INSERT/UPDATE/DELETE/MERGE/CREATE TABLE/ALTER TABLE are supported.
// Full dispatch will be implemented in batch 34.
func (p *Parser) parseStmt() (nodes.Node, error) {
	if p.collectMode() {
		p.cachedCollect("parseStmt", func() {
			// In collection mode, add all statement-starting keywords as candidates.
			stmtTokens := []int{
				SELECT, VALUES, TABLE, WITH, INSERT, UPDATE, DELETE_P, MERGE,
				CREATE, COMMENT, SECURITY, ALTER, REFRESH,
				BEGIN_P, START, COMMIT, END_P, ABORT_P, SAVEPOINT, RELEASE,
				ROLLBACK, PREPARE, EXECUTE, DEALLOCATE,
				SET, SHOW, RESET, GRANT, REVOKE, DROP, TRUNCATE,
				LOCK_P, DECLARE, FETCH, MOVE, CLOSE,
				VACUUM, ANALYZE, ANALYSE, CLUSTER, REINDEX,
				COPY, IMPORT_P, EXPLAIN, DO, CHECKPOINT,
				DISCARD, LISTEN, UNLISTEN, NOTIFY, LOAD, CALL, REASSIGN,
			}
			for _, t := range stmtTokens {
				p.addTokenCandidate(t)
			}
		})
		return nil, errCollecting
	}
	switch p.cur.Type {
	case '(':
		// Parenthesized SELECT: (SELECT ...) UNION ... or (SELECT ...)
		return p.parseSelectNoParens()
	case SELECT, VALUES, TABLE:
		return p.parseSelectNoParens()
	case WITH:
		return p.parseWithStmt()
	case INSERT:
		return p.parseInsertStmt(nil)
	case UPDATE:
		return p.parseUpdateStmt(nil)
	case DELETE_P:
		return p.parseDeleteStmt(nil)
	case MERGE:
		return p.parseMergeStmt(nil)
	case CREATE:
		return p.parseCreateDispatch()
	case COMMENT:
		commentLoc := p.pos()
		p.advance() // consume COMMENT
		n, err := p.parseCommentStmt()
		if cs, ok := n.(*nodes.CommentStmt); ok && err == nil {
			cs.Loc = nodes.Loc{Start: commentLoc, End: p.prev.End}
		}
		return n, err
	case SECURITY:
		secLoc := p.pos()
		p.advance() // consume SECURITY
		n, err := p.parseSecLabelStmt()
		if sl, ok := n.(*nodes.SecLabelStmt); ok && err == nil {
			sl.Loc = nodes.Loc{Start: secLoc, End: p.prev.End}
		}
		return n, err
	case ALTER:
		alterLoc := p.pos() // capture ALTER position before consuming
		p.advance()         // consume ALTER
		if p.collectMode() {
			alterTokens := []int{
				DATABASE, ROLE, USER, SERVER, GROUP_P, POLICY,
				PUBLICATION, SUBSCRIPTION, STATISTICS, OPERATOR,
				SCHEMA, DEFAULT, FUNCTION, PROCEDURE, ROUTINE,
				TYPE_P, DOMAIN_P, COLLATION, CONVERSION_P, EXTENSION,
				AGGREGATE, TEXT_P, LANGUAGE, PROCEDURAL, LARGE_P,
				EVENT, SYSTEM_P, TABLESPACE, TRIGGER, RULE, TABLE,
			}
			for _, t := range alterTokens {
				p.addTokenCandidate(t)
			}
			return nil, errCollecting
		}
		switch p.cur.Type {
		case DATABASE:
			return p.parseAlterDatabaseDispatch(alterLoc)
		case ROLE:
			return p.parseAlterRoleStmt()
		case USER:
			// ALTER USER MAPPING or ALTER USER (role)
			if p.peekNext().Type == MAPPING {
				return p.parseAlterUserMappingStmt()
			}
			return p.parseAlterRoleStmt()
		case SERVER:
			return p.parseAlterForeignServerStmt()
		case GROUP_P:
			return p.parseAlterGroupStmt()
		case POLICY:
			return p.parseAlterPolicyStmt()
		case PUBLICATION:
			return p.parseAlterPublicationStmt()
		case SUBSCRIPTION:
			return p.parseAlterSubscriptionStmt()
		case STATISTICS:
			return p.parseAlterStatisticsStmt(alterLoc)
		case OPERATOR:
			return p.parseAlterOperatorStmt(alterLoc)
		case SCHEMA:
			return p.parseAlterSchemaOwner()
		case DEFAULT:
			return p.parseAlterDefaultPrivilegesStmt(alterLoc)
		case FUNCTION, PROCEDURE, ROUTINE:
			return p.parseAlterFunctionStmt()
		case TYPE_P:
			return p.parseAlterTypeStmt()
		case DOMAIN_P:
			return p.parseAlterDomainOwnerOrOther()
		case COLLATION:
			return p.parseAlterCollationStmt()
		case CONVERSION_P:
			return p.parseAlterConversionStmt()
		case EXTENSION:
			return p.parseAlterExtensionStmt()
		case AGGREGATE:
			return p.parseAlterAggregateStmt()
		case TEXT_P:
			return p.parseAlterTextSearchStmt()
		case LANGUAGE:
			return p.parseAlterLanguageStmt()
		case PROCEDURAL:
			return p.parseAlterLanguageStmt()
		case LARGE_P:
			return p.parseAlterLargeObjectStmt()
		case EVENT:
			return p.parseAlterEventTriggerOwner()
		case SYSTEM_P:
			return p.parseAlterSystemStmt()
		case TABLESPACE:
			return p.parseAlterTablespaceOwner()
		case TRIGGER:
			return p.parseAlterTriggerDependsOnExtension()
		case RULE:
			return p.parseAlterRuleStmt()
		default:
			return p.parseAlterTableStmt()
		}
	case REFRESH:
		p.advance() // consume REFRESH
		return p.parseRefreshMatViewStmt()
	case BEGIN_P, START, COMMIT, END_P, ABORT_P, SAVEPOINT, RELEASE:
		return p.parseTransactionStmt()
	case ROLLBACK:
		return p.parseTransactionStmt()
	case PREPARE:
		// PREPARE TRANSACTION is a transaction stmt; plain PREPARE is a prepared stmt.
		if p.peekNext().Type == TRANSACTION {
			return p.parseTransactionStmt()
		}
		return p.parsePrepareStmt()
	case EXECUTE:
		return p.parseExecuteStmt()
	case DEALLOCATE:
		return p.parseDeallocateStmt()
	case SET:
		p.advance() // consume SET
		// SET CONSTRAINTS is a different statement type.
		if p.cur.Type == CONSTRAINTS {
			return p.parseConstraintsSetStmt()
		}
		return p.parseVariableSetStmt()
	case SHOW:
		p.advance() // consume SHOW
		return p.parseVariableShowStmt()
	case RESET:
		p.advance() // consume RESET
		return p.parseVariableResetStmt()
	case GRANT:
		p.advance() // consume GRANT
		return p.parseGrantStmt()
	case REVOKE:
		p.advance() // consume REVOKE
		return p.parseRevokeStmt()
	case DROP:
		dropLoc := p.pos()
		p.advance() // consume DROP
		if p.collectMode() {
			dropTokens := []int{
				TABLE, VIEW, MATERIALIZED, INDEX, SEQUENCE,
				FUNCTION, PROCEDURE, ROUTINE, AGGREGATE,
				OPERATOR, TYPE_P, DOMAIN_P, COLLATION, CONVERSION_P,
				SCHEMA, DATABASE, ROLE, USER, GROUP_P,
				POLICY, TRIGGER, RULE, EXTENSION, EVENT,
				FOREIGN, SERVER, PUBLICATION, SUBSCRIPTION,
				TABLESPACE, TEXT_P, ACCESS, CAST, TRANSFORM,
				LANGUAGE, TRUSTED, PROCEDURAL, OWNED,
			}
			for _, t := range dropTokens {
				p.addTokenCandidate(t)
			}
			return nil, errCollecting
		}
		return p.parseDropStmt(dropLoc)
	case TRUNCATE:
		p.advance() // consume TRUNCATE
		return p.parseTruncateStmt()
	case LOCK_P:
		return p.parseLockStmt()
	case DECLARE:
		return p.parseDeclareCursorStmt()
	case FETCH, MOVE:
		return p.parseFetchStmt()
	case CLOSE:
		return p.parseClosePortalStmt()
	case VACUUM:
		return p.parseVacuumStmt()
	case ANALYZE, ANALYSE:
		return p.parseAnalyzeStmt()
	case CLUSTER:
		return p.parseClusterStmt()
	case REINDEX:
		return p.parseReindexStmt()
	case COPY:
		p.advance() // consume COPY
		return p.parseCopyStmt()
	case IMPORT_P:
		p.advance() // consume IMPORT
		return p.parseImportForeignSchemaStmt()
	case EXPLAIN:
		p.advance() // consume EXPLAIN
		return p.parseExplainStmt()
	case DO:
		p.advance() // consume DO
		return p.parseDoStmt()
	case CHECKPOINT:
		p.advance() // consume CHECKPOINT
		return p.parseCheckPointStmt()
	case DISCARD:
		p.advance() // consume DISCARD
		return p.parseDiscardStmt()
	case LISTEN:
		p.advance() // consume LISTEN
		return p.parseListenStmt()
	case UNLISTEN:
		p.advance() // consume UNLISTEN
		return p.parseUnlistenStmt()
	case NOTIFY:
		p.advance() // consume NOTIFY
		return p.parseNotifyStmt()
	case LOAD:
		p.advance() // consume LOAD
		return p.parseLoadStmt()
	case CALL:
		p.advance() // consume CALL
		return p.parseCallStmt()
	case REASSIGN:
		p.advance() // consume REASSIGN
		return p.parseReassignOwnedStmt()
	default:
		return nil, nil
	}
}

// parseCreateDispatch handles CREATE ... statements by peeking at what follows.
//
// The current token is CREATE. We peek at the next token to determine which
// CREATE sub-statement to parse.
func (p *Parser) parseCreateDispatch() (nodes.Node, error) {
	createLoc := p.pos() // byte offset of CREATE keyword
	// In collect mode, check if the next token is at/past cursor.
	// We need to peek ahead because CREATE has not been consumed yet.
	if p.completing && !p.collecting {
		next := p.peekNext()
		if next.Loc >= p.cursorOff || next.Type == 0 {
			p.collecting = true
		}
	}
	if p.collectMode() {
		createTokens := []int{
			OR, VIEW, RECURSIVE, MATERIALIZED, TABLE,
			TEMPORARY, TEMP, LOCAL, UNLOGGED,
			UNIQUE, INDEX, SEQUENCE, DOMAIN_P, TYPE_P,
			AGGREGATE, OPERATOR, TEXT_P, COLLATION, STATISTICS,
			FUNCTION, PROCEDURE, DATABASE, ROLE, USER, GROUP_P,
			POLICY, TRIGGER, CONSTRAINT, EVENT, FOREIGN, SERVER,
			LANGUAGE, TRUSTED, PROCEDURAL, GLOBAL,
			PUBLICATION, SUBSCRIPTION, RULE, EXTENSION, ACCESS,
			CAST, TRANSFORM, CONVERSION_P, DEFAULT, TABLESPACE, SCHEMA,
		}
		for _, t := range createTokens {
			p.addTokenCandidate(t)
		}
		return nil, errCollecting
	}
	next := p.peekNext()
	switch next.Type {
	case OR:
		// CREATE OR REPLACE VIEW/FUNCTION/PROCEDURE/TRIGGER ...
		createLoc := p.pos() // position of CREATE
		p.advance()          // consume CREATE
		p.advance()          // consume OR
		p.expect(REPLACE)
		switch p.cur.Type {
		case FUNCTION, PROCEDURE:
			return p.parseCreateFunctionStmt(createLoc, true)
		case TRIGGER, CONSTRAINT:
			return p.parseCreateTrigStmt(true)
		case TRUSTED, PROCEDURAL, LANGUAGE:
			return p.parseCreatePLangStmt(true)
		case RULE:
			return p.parseCreateRuleStmt(true, createLoc)
		case AGGREGATE:
			return p.parseDefineStmtAggregate(true, createLoc)
		case TRANSFORM:
			return p.parseCreateTransformStmt(true, createLoc)
		default:
			return p.parseViewStmt(true)
		}
	case VIEW:
		// CREATE VIEW ...
		p.advance() // consume CREATE
		return p.parseViewStmt(false)
	case RECURSIVE:
		// CREATE RECURSIVE VIEW ...
		p.advance() // consume CREATE
		return p.parseViewStmt(false)
	case MATERIALIZED:
		// CREATE MATERIALIZED VIEW ...
		p.advance() // consume CREATE
		relpersistence := byte(nodes.RELPERSISTENCE_PERMANENT)
		p.advance() // consume MATERIALIZED
		return p.parseCreateMatViewStmt(relpersistence)
	case TABLE:
		// CREATE TABLE ... (could be regular CREATE TABLE or CREATE TABLE AS)
		return p.parseCreateOrCTAS()
	case TEMPORARY, TEMP:
		// CREATE TEMP/TEMPORARY TABLE|VIEW ...
		return p.parseCreateTempDispatch()
	case LOCAL:
		// CREATE LOCAL TEMP/TEMPORARY TABLE|VIEW ...
		return p.parseCreateTempDispatch()
	case UNLOGGED:
		// CREATE UNLOGGED TABLE|MATERIALIZED VIEW ...
		return p.parseCreateUnloggedDispatch()
	case UNIQUE:
		// CREATE UNIQUE INDEX ...
		p.advance() // consume CREATE
		return p.parseIndexStmt()
	case INDEX:
		// CREATE INDEX ...
		p.advance() // consume CREATE
		return p.parseIndexStmt()
	case SEQUENCE:
		// CREATE SEQUENCE ...
		loc := p.pos()
		p.advance() // consume CREATE
		return p.parseCreateSeqStmt(loc, byte(nodes.RELPERSISTENCE_PERMANENT))
	case DOMAIN_P:
		// CREATE DOMAIN ...
		loc := p.pos()
		p.advance() // consume CREATE
		return p.parseCreateDomainStmt(loc)
	case TYPE_P:
		// CREATE TYPE ... (base, composite, enum, range, shell)
		p.advance() // consume CREATE
		return p.parseDefineStmtType(createLoc)
	case AGGREGATE:
		// CREATE AGGREGATE ...
		p.advance() // consume CREATE
		return p.parseDefineStmtAggregate(false, createLoc)
	case OPERATOR:
		// CREATE OPERATOR ... / CREATE OPERATOR CLASS ... / CREATE OPERATOR FAMILY ...
		p.advance() // consume CREATE
		return p.parseDefineStmtOperator(createLoc)
	case TEXT_P:
		// CREATE TEXT SEARCH ...
		p.advance() // consume CREATE
		return p.parseDefineStmtTextSearch(createLoc)
	case COLLATION:
		// CREATE COLLATION ...
		p.advance() // consume CREATE
		return p.parseDefineStmtCollation(createLoc)
	case STATISTICS:
		// CREATE STATISTICS ...
		p.advance() // consume CREATE
		return p.parseCreateStatsStmt(createLoc)
	case FUNCTION:
		// CREATE FUNCTION ...
		p.advance() // consume CREATE
		return p.parseCreateFunctionStmt(createLoc, false)
	case PROCEDURE:
		// CREATE PROCEDURE ...
		p.advance() // consume CREATE
		return p.parseCreateFunctionStmt(createLoc, false)
	case DATABASE:
		// CREATE DATABASE ...
		p.advance() // consume CREATE
		return p.parseCreatedbStmt(createLoc)
	case ROLE:
		// CREATE ROLE ...
		p.advance() // consume CREATE
		return p.parseCreateRoleStmt()
	case USER:
		// CREATE USER ... or CREATE USER MAPPING ...
		p.advance() // consume CREATE
		// Peek: if next after USER is MAPPING, it's CREATE USER MAPPING
		if p.peekNext().Type == MAPPING {
			return p.parseCreateUserMappingIfNotExistsStmt()
		}
		return p.parseCreateUserStmt()
	case GROUP_P:
		// CREATE GROUP ...
		p.advance() // consume CREATE
		return p.parseCreateGroupStmt()
	case POLICY:
		// CREATE POLICY ...
		p.advance() // consume CREATE
		return p.parseCreatePolicyStmt()
	case TRIGGER:
		// CREATE TRIGGER ...
		p.advance() // consume CREATE
		return p.parseCreateTrigStmt(false)
	case CONSTRAINT:
		// CREATE CONSTRAINT TRIGGER ...
		p.advance() // consume CREATE
		return p.parseCreateTrigStmt(false)
	case EVENT:
		// CREATE EVENT TRIGGER ...
		p.advance() // consume CREATE
		p.advance() // consume EVENT
		return p.parseCreateEventTrigStmt()
	case FOREIGN:
		// CREATE FOREIGN DATA WRAPPER or CREATE FOREIGN TABLE
		p.advance() // consume CREATE
		p.advance() // consume FOREIGN
		if p.cur.Type == DATA_P {
			// CREATE FOREIGN DATA WRAPPER
			return p.parseCreateFdwStmt()
		}
		// CREATE FOREIGN TABLE
		return p.parseCreateForeignTableStmt()
	case SERVER:
		// CREATE SERVER ...
		p.advance() // consume CREATE
		return p.parseCreateForeignServerStmt()
	case LANGUAGE:
		// CREATE LANGUAGE ...
		p.advance() // consume CREATE
		return p.parseCreatePLangStmt(false)
	case TRUSTED:
		// CREATE TRUSTED [PROCEDURAL] LANGUAGE ...
		p.advance() // consume CREATE
		return p.parseCreatePLangStmt(false)
	case PROCEDURAL:
		// CREATE PROCEDURAL LANGUAGE ...
		p.advance() // consume CREATE
		return p.parseCreatePLangStmt(false)
	case GLOBAL:
		// CREATE GLOBAL TEMP TABLE ... (same as CREATE TEMP)
		return p.parseCreateTempDispatch()
	case PUBLICATION:
		// CREATE PUBLICATION ...
		p.advance() // consume CREATE
		return p.parseCreatePublicationStmt()
	case SUBSCRIPTION:
		// CREATE SUBSCRIPTION ...
		p.advance() // consume CREATE
		return p.parseCreateSubscriptionStmt()
	case RULE:
		// CREATE RULE ...
		p.advance() // consume CREATE
		return p.parseCreateRuleStmt(false, createLoc)
	case EXTENSION:
		// CREATE EXTENSION ...
		p.advance() // consume CREATE
		return p.parseCreateExtensionStmt()
	case ACCESS:
		// CREATE ACCESS METHOD ...
		p.advance() // consume CREATE
		return p.parseCreateAmStmt()
	case CAST:
		// CREATE CAST ...
		p.advance() // consume CREATE
		return p.parseCreateCastStmt()
	case TRANSFORM:
		// CREATE TRANSFORM ...
		createLoc := p.pos() // position of CREATE
		p.advance()          // consume CREATE
		return p.parseCreateTransformStmt(false, createLoc)
	case CONVERSION_P:
		// CREATE CONVERSION ...
		p.advance() // consume CREATE
		return p.parseCreateConversionStmt(false, createLoc)
	case DEFAULT:
		// CREATE DEFAULT CONVERSION ...
		p.advance() // consume CREATE
		p.advance() // consume DEFAULT
		return p.parseCreateConversionStmt(true, createLoc)
	case TABLESPACE:
		// CREATE TABLESPACE ...
		p.advance() // consume CREATE
		return p.parseCreateTableSpaceStmt(createLoc)
	case SCHEMA:
		// CREATE SCHEMA ...
		p.advance() // consume CREATE
		return p.parseCreateSchemaStmt(createLoc)
	default:
		return nil, nil
	}
}

// parseCreateTempDispatch dispatches CREATE TEMP/TEMPORARY/LOCAL ... statements.
// The CREATE keyword has NOT been consumed yet.
//
// We need to look past OptTemp to see if it's TABLE or VIEW.
func (p *Parser) parseCreateTempDispatch() (nodes.Node, error) {
	loc := p.pos()
	p.advance() // consume CREATE
	relpersistence := p.parseOptTemp()

	if p.cur.Type == VIEW || p.cur.Type == RECURSIVE {
		// CREATE TEMP VIEW ... or CREATE TEMP RECURSIVE VIEW ...
		stmt, err := p.parseViewStmt(false)
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			stmt.View.Relpersistence = relpersistence
		}
		return stmt, nil
	}

	if p.cur.Type == SEQUENCE {
		// CREATE TEMP SEQUENCE ...
		return p.parseCreateSeqStmt(loc, relpersistence)
	}

	// CREATE TEMP TABLE ...
	p.expect(TABLE)

	ifNotExists := false
	if p.cur.Type == IF_P {
		p.advance()
		p.expect(NOT)
		p.expect(EXISTS)
		ifNotExists = true
	}

	names, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}

	// Determine if this is CTAS or regular CREATE TABLE.
	// CTAS: no '(' for column defs, just optional column list then AS.
	// Regular: has '(' for column defs, or PARTITION OF, or OF.
	if p.cur.Type == '(' {
		return p.parseCreateTableOrCTASAfterParen(names, relpersistence, ifNotExists)
	}

	// No '(' - check for AS (CTAS with no column list), PARTITION, OF
	if p.cur.Type == AS || p.cur.Type == USING {
		return p.finishCTAS(names, nil, relpersistence, ifNotExists)
	}

	// Must be PARTITION OF, OF, or error
	return p.finishCreateStmt(names, relpersistence, ifNotExists)
}

// parseCreateUnloggedDispatch dispatches CREATE UNLOGGED ... statements.
// The CREATE keyword has NOT been consumed yet.
func (p *Parser) parseCreateUnloggedDispatch() (nodes.Node, error) {
	p.advance() // consume CREATE
	p.advance() // consume UNLOGGED
	relpersistence := byte(nodes.RELPERSISTENCE_UNLOGGED)

	if p.cur.Type == MATERIALIZED {
		p.advance() // consume MATERIALIZED
		return p.parseCreateMatViewStmt(relpersistence)
	}

	if p.cur.Type == VIEW || p.cur.Type == RECURSIVE {
		// CREATE UNLOGGED VIEW ... or CREATE UNLOGGED RECURSIVE VIEW ...
		stmt, err := p.parseViewStmt(false)
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			stmt.View.Relpersistence = relpersistence
		}
		return stmt, nil
	}

	// CREATE UNLOGGED TABLE ...
	p.expect(TABLE)

	ifNotExists := false
	if p.cur.Type == IF_P {
		p.advance()
		p.expect(NOT)
		p.expect(EXISTS)
		ifNotExists = true
	}

	names, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}

	if p.cur.Type == '(' {
		return p.parseCreateTableOrCTASAfterParen(names, relpersistence, ifNotExists)
	}

	if p.cur.Type == AS || p.cur.Type == USING {
		return p.finishCTAS(names, nil, relpersistence, ifNotExists)
	}

	return p.finishCreateStmt(names, relpersistence, ifNotExists)
}

// parseCreateOrCTAS parses CREATE TABLE ... which could be either a regular
// CREATE TABLE or a CREATE TABLE AS (CTAS) statement.
// The CREATE keyword has NOT been consumed yet.
func (p *Parser) parseCreateOrCTAS() (nodes.Node, error) {
	p.advance() // consume CREATE
	relpersistence := byte(nodes.RELPERSISTENCE_PERMANENT)
	p.expect(TABLE)

	ifNotExists := false
	if p.cur.Type == IF_P {
		p.advance()
		p.expect(NOT)
		p.expect(EXISTS)
		ifNotExists = true
	}

	names, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}

	// Determine if this is CTAS or regular CREATE TABLE.
	if p.cur.Type == '(' {
		return p.parseCreateTableOrCTASAfterParen(names, relpersistence, ifNotExists)
	}

	// No '(' - check for AS or USING (CTAS with no column list)
	if p.cur.Type == AS || p.cur.Type == USING {
		return p.finishCTAS(names, nil, relpersistence, ifNotExists)
	}

	// PARTITION OF, OF, etc.
	return p.finishCreateStmt(names, relpersistence, ifNotExists)
}

// parseCreateTableOrCTASAfterParen handles the ambiguous case where we see '('
// after the table name. This could be either:
//   - CREATE TABLE t (col1 int, ...) -- regular CREATE TABLE with column defs
//   - CREATE TABLE t (col1, col2) AS SELECT ... -- CTAS with column list
//
// In CTAS, the column list is just names (no types). In CREATE TABLE, we have
// full table element definitions. We can distinguish them by looking at what
// follows the first identifier inside the parens.
func (p *Parser) parseCreateTableOrCTASAfterParen(names *nodes.List, relpersistence byte, ifNotExists bool) (nodes.Node, error) {
	// Save state: we're at '('
	// We need to look ahead to determine if this is CTAS or CREATE TABLE.
	// Strategy: parse as CREATE TABLE (which expects '(' TableElementList ')').
	// The caller already parsed up to the '('.
	// For a regular CREATE TABLE, we expect '(' then column defs.
	// For CTAS, we expect '(' then column names ')' then AS.
	//
	// We use a simple heuristic: parse inside parens. If the first token
	// after a name is ',' or ')' (no type), it's likely a CTAS column list.
	// If it's a type name or constraint keyword, it's a CREATE TABLE.
	// Also special cases: LIKE, CHECK, CONSTRAINT, PRIMARY, UNIQUE, FOREIGN, EXCLUDE
	// are table element starters.

	// Peek inside the parens to determine the type.
	// After '(' we see the first token. If it's a table constraint keyword, it's CREATE TABLE.
	// If it's a column name followed by ',' or ')', it's CTAS.
	// If it's a column name followed by a type, it's CREATE TABLE.
	// If it's empty '()', it's CREATE TABLE.

	// We need to just parse it as CREATE TABLE. The existing parseCreateStmt
	// code handles everything after names are parsed.
	rv := makeRangeVarFromNames(names)
	rv.Relpersistence = relpersistence

	stmt := &nodes.CreateStmt{
		Relation:    rv,
		IfNotExists: ifNotExists,
	}

	// '(' OptTableElementList ')'
	p.expect('(')

	// Check for immediate ')' (empty table)
	if p.cur.Type == ')' {
		p.advance()
		stmt.InhRelations = p.parseOptInherit()
		stmt.Partspec = p.parseOptPartitionSpec()
		stmt.AccessMethod = p.parseOptAccessMethod()
		stmt.Options = p.parseOptWith()
		stmt.OnCommit = p.parseOnCommitOption()
		stmt.Tablespacename = p.parseOptTableSpace()
		return stmt, nil
	}

	// Check if this is a CTAS column list by looking at what the first identifier
	// is followed by. In CTAS, column names are bare ColIds separated by commas.
	// In CREATE TABLE, column names are followed by type names.
	// Special keywords like LIKE, CHECK, CONSTRAINT, PRIMARY, UNIQUE, FOREIGN, EXCLUDE
	// indicate CREATE TABLE.
	if p.isCreateTableElement() {
		// This is a regular CREATE TABLE.
		tableElts, err := p.parseOptTableElementList()
		if err != nil {
			return nil, err
		}
		stmt.TableElts = tableElts
		p.expect(')')
		stmt.InhRelations = p.parseOptInherit()
		stmt.Partspec = p.parseOptPartitionSpec()
		stmt.AccessMethod = p.parseOptAccessMethod()
		stmt.Options = p.parseOptWith()
		stmt.OnCommit = p.parseOnCommitOption()
		stmt.Tablespacename = p.parseOptTableSpace()
		return stmt, nil
	}

	// Might be CTAS with column list, but we need to verify.
	// If the first token is a ColId, we check the next token.
	// Actually, table columns also start with ColId (the column name).
	// The difference is what follows:
	//   - CTAS: ColId followed by ',' or ')'
	//   - CREATE TABLE: ColId followed by type name
	//
	// Since a type name starts with a ColId too (like 'integer', 'text', etc.),
	// we actually need more lookahead. But we can use the following heuristic:
	// If after the first ColId we see ',' or ')', it's CTAS.
	// Otherwise it's CREATE TABLE.

	if p.isColId() {
		// Save the current token (the ColId) for backtracking.
		savedPrev := p.prev
		firstName := p.advance() // consume ColId; firstName = the ColId token

		if p.cur.Type == ',' || p.cur.Type == ')' {
			// This is a CTAS column list. Parse the rest.
			colNames := []nodes.Node{&nodes.String{Str: firstName.Str}}
			for p.cur.Type == ',' {
				p.advance()
				name, err := p.parseName()
				if err != nil {
					break
				}
				colNames = append(colNames, &nodes.String{Str: name})
			}
			p.expect(')')
			colList := &nodes.List{Items: colNames}
			return p.finishCTAS(names, colList, relpersistence, ifNotExists)
		}

		// It's a CREATE TABLE. Restore: push current token into lookahead
		// buffer and restore the ColId as the current token.
		p.nextBuf = p.cur
		p.hasNext = true
		p.cur = firstName
		p.prev = savedPrev

		tableElts, err := p.parseOptTableElementList()
		if err != nil {
			return nil, err
		}
		stmt.TableElts = tableElts
		p.expect(')')
		stmt.InhRelations = p.parseOptInherit()
		stmt.Partspec = p.parseOptPartitionSpec()
		stmt.AccessMethod = p.parseOptAccessMethod()
		stmt.Options = p.parseOptWith()
		stmt.OnCommit = p.parseOnCommitOption()
		stmt.Tablespacename = p.parseOptTableSpace()
		return stmt, nil
	}

	// Fallback: parse as CREATE TABLE
	tableElts, err := p.parseOptTableElementList()
	if err != nil {
		return nil, err
	}
	stmt.TableElts = tableElts
	p.expect(')')
	stmt.InhRelations = p.parseOptInherit()
	stmt.Partspec = p.parseOptPartitionSpec()
	stmt.AccessMethod = p.parseOptAccessMethod()
	stmt.Options = p.parseOptWith()
	stmt.OnCommit = p.parseOnCommitOption()
	stmt.Tablespacename = p.parseOptTableSpace()
	return stmt, nil
}

// isCreateTableElement returns true if the current token starts a table element
// that is NOT a simple column name (i.e., a table constraint or LIKE).
func (p *Parser) isCreateTableElement() bool {
	switch p.cur.Type {
	case LIKE, CHECK, CONSTRAINT, PRIMARY, UNIQUE, FOREIGN, EXCLUDE:
		return true
	}
	return false
}

// finishCTAS completes parsing a CTAS statement after names and optional column
// list have been parsed. Handles OptAccessMethod, OptWith, OnCommitOption,
// OptTableSpace, AS, SelectStmt, and opt_with_data.
func (p *Parser) finishCTAS(names *nodes.List, colNames *nodes.List, relpersistence byte, ifNotExists bool) (*nodes.CreateTableAsStmt, error) {
	rv := makeRangeVarFromAnyName(names)
	rv.Relpersistence = relpersistence

	accessMethod := p.parseOptAccessMethod()
	options := p.parseOptWith()
	onCommit := p.parseOnCommitOption()
	tableSpace := p.parseOptTableSpace()

	// AS
	p.expect(AS)

	var query nodes.Node
	if p.cur.Type == EXECUTE {
		p.advance()
		name, err := p.parseName()
		if err != nil {
			return nil, err
		}
		params, err := p.parseExecuteParamClause()
		if err != nil {
			return nil, err
		}
		query = &nodes.ExecuteStmt{
			Name:   name,
			Params: params,
		}
	} else {
		var err error
		query, err = p.parseSelectNoParens()
		if err != nil {
			return nil, err
		}
	}

	withData := p.parseOptWithData()

	into := &nodes.IntoClause{
		Rel:            rv,
		ColNames:       colNames,
		AccessMethod:   accessMethod,
		Options:        options,
		OnCommit:       onCommit,
		TableSpaceName: tableSpace,
		SkipData:       !withData,
		Loc:            rv.Loc,
	}

	return &nodes.CreateTableAsStmt{
		Query:       query,
		Into:        into,
		Objtype:     nodes.OBJECT_TABLE,
		IfNotExists: ifNotExists,
	}, nil
}

// finishCreateStmt completes parsing a regular CREATE TABLE after the name
// has been parsed. Handles PARTITION OF, OF, etc.
func (p *Parser) finishCreateStmt(names *nodes.List, relpersistence byte, ifNotExists bool) (*nodes.CreateStmt, error) {
	rv := makeRangeVarFromNames(names)
	rv.Relpersistence = relpersistence

	stmt := &nodes.CreateStmt{
		Relation:    rv,
		IfNotExists: ifNotExists,
	}

	switch p.cur.Type {
	case PARTITION:
		p.advance()
		p.expect(OF)
		return p.parseCreateStmtPartitionOf(stmt, relpersistence)
	case OF:
		p.advance()
		return p.parseCreateStmtOf(stmt)
	}

	// Should have '(' but we already handled that case
	p.expect('(')
	tableElts, err := p.parseOptTableElementList()
	if err != nil {
		return nil, err
	}
	stmt.TableElts = tableElts
	p.expect(')')
	stmt.InhRelations = p.parseOptInherit()
	stmt.Partspec = p.parseOptPartitionSpec()
	stmt.AccessMethod = p.parseOptAccessMethod()
	stmt.Options = p.parseOptWith()
	stmt.OnCommit = p.parseOnCommitOption()
	stmt.Tablespacename = p.parseOptTableSpace()
	return stmt, nil
}

// parseWithStmt parses a WITH clause followed by SELECT, INSERT, UPDATE, DELETE, or MERGE.
func (p *Parser) parseWithStmt() (nodes.Node, error) {
	withClause, err := p.parseWithClause()
	if err != nil {
		return nil, err
	}
	switch p.cur.Type {
	case INSERT:
		return p.parseInsertStmt(withClause)
	case UPDATE:
		return p.parseUpdateStmt(withClause)
	case DELETE_P:
		return p.parseDeleteStmt(withClause)
	case MERGE:
		return p.parseMergeStmt(withClause)
	default:
		// SELECT
		stmt, err := p.parseSelectClause(setOpPrecNone)
		if err != nil {
			return nil, err
		}
		if stmt == nil {
			return nil, nil
		}
		if p.cur.Type == ORDER {
			p.advance()
			p.expect(BY)
			stmt.SortClause, err = p.parseSortByList()
			if err != nil {
				return nil, err
			}
		}
		p.parseSelectOptions(stmt)
		stmt.WithClause = withClause
		return stmt, nil
	}
}

// advance consumes the current token and moves to the next one.
// It maps lexer-internal token types (lex_*) to parser token constants.
// It also performs lookahead-based token reclassification (NOT -> NOT_LA).
func (p *Parser) advance() Token {
	p.prev = p.cur
	if p.hasNext {
		p.cur = p.nextBuf
		p.hasNext = false
	} else {
		tok := p.lexer.NextToken()
		tok.Type = mapTokenType(tok.Type)
		p.cur = tok
	}
	// Lookahead-based token reclassification (mirrors PostgreSQL's lexer lookahead).
	switch p.cur.Type {
	case NOT:
		// NOT -> NOT_LA when followed by BETWEEN, IN, LIKE, ILIKE, SIMILAR
		next := p.peekNext()
		switch next.Type {
		case BETWEEN, IN_P, LIKE, ILIKE, SIMILAR:
			p.cur.Type = NOT_LA
		}
	case WITH:
		// WITH -> WITH_LA when followed by TIME or ORDINALITY
		next := p.peekNext()
		switch next.Type {
		case TIME, ORDINALITY:
			p.cur.Type = WITH_LA
		}
	case WITHOUT:
		// WITHOUT -> WITHOUT_LA when followed by TIME
		if p.peekNext().Type == TIME {
			p.cur.Type = WITHOUT_LA
		}
	case NULLS_P:
		// NULLS -> NULLS_LA when followed by FIRST or LAST
		next := p.peekNext()
		switch next.Type {
		case FIRST_P, LAST_P:
			p.cur.Type = NULLS_LA
		}
	}
	if p.completing && !p.collecting {
		p.checkCursor()
	}
	return p.prev
}

// peekNext returns the next token after cur without consuming it.
// Used for 2-token lookahead.
func (p *Parser) peekNext() Token {
	if !p.hasNext {
		tok := p.lexer.NextToken()
		tok.Type = mapTokenType(tok.Type)
		p.nextBuf = tok
		p.hasNext = true
	}
	return p.nextBuf
}

// mapTokenType maps lexer-internal token types (lex_* constants starting at
// nonKeywordTokenBase=800) to parser token constants (from tokens.go).
// Single-char tokens (0-255) and keyword tokens are passed through unchanged.
func mapTokenType(typ int) int {
	if typ == 0 || (typ > 0 && typ < 256) {
		return typ
	}
	if typ >= nonKeywordTokenBase && typ < nonKeywordTokenBase+100 {
		offset := typ - nonKeywordTokenBase
		switch offset {
		case 0: // lex_ICONST
			return ICONST
		case 1: // lex_FCONST
			return FCONST
		case 2: // lex_SCONST
			return SCONST
		case 3: // lex_BCONST
			return BCONST
		case 4: // lex_XCONST
			return XCONST
		case 5: // lex_USCONST
			return SCONST
		case 6: // lex_IDENT
			return IDENT
		case 7: // lex_UIDENT
			return IDENT
		case 8: // lex_TYPECAST
			return TYPECAST
		case 9: // lex_DOT_DOT
			return DOT_DOT
		case 10: // lex_COLON_EQUALS
			return COLON_EQUALS
		case 11: // lex_EQUALS_GREATER
			return EQUALS_GREATER
		case 12: // lex_LESS_EQUALS
			return LESS_EQUALS
		case 13: // lex_GREATER_EQUALS
			return GREATER_EQUALS
		case 14: // lex_NOT_EQUALS
			return NOT_EQUALS
		case 15: // lex_PARAM
			return PARAM
		case 16: // lex_Op
			return Op
		}
		return 0
	}
	// Keywords and other tokens pass through directly.
	return typ
}

// peek returns the current token without consuming it.
func (p *Parser) peek() Token {
	return p.cur
}

// match checks if the current token type matches any of the given types.
// If it matches, the token is consumed and returned with ok=true.
func (p *Parser) match(types ...int) (Token, bool) {
	if p.collectMode() {
		for _, t := range types {
			p.addTokenCandidate(t)
		}
		return Token{}, false
	}
	for _, t := range types {
		if p.cur.Type == t {
			return p.advance(), true
		}
	}
	return Token{}, false
}

// expect consumes the current token if it matches the expected type.
// Returns an error if the token does not match.
func (p *Parser) expect(tokenType int) (Token, error) {
	if p.collectMode() {
		p.addTokenCandidate(tokenType)
		return Token{}, errCollecting
	}
	if p.cur.Type == tokenType {
		return p.advance(), nil
	}
	return Token{}, p.syntaxErrorAtCur()
}

// pos returns the byte position of the current token.
func (p *Parser) pos() int {
	return p.cur.Loc
}

// ParseError represents a parse error with position information.
type ParseError struct {
	Severity string
	Code     string
	Message  string
	Position int
}

func (e *ParseError) Error() string {
	sev := e.Severity
	if sev == "" {
		sev = "ERROR"
	}
	code := e.Code
	if code == "" {
		code = "42601"
	}
	return fmt.Sprintf("%s: %s (SQLSTATE %s)", sev, e.Message, code)
}

func (p *Parser) syntaxErrorAtCur() *ParseError {
	return p.syntaxErrorAtTok(p.cur)
}

func (p *Parser) syntaxErrorAtTok(tok Token) *ParseError {
	text := p.tokenText(tok)
	var msg string
	if text == "" {
		msg = "syntax error at end of input"
	} else {
		msg = fmt.Sprintf("syntax error at or near \"%s\"", text)
	}
	return &ParseError{Message: msg, Position: tok.Loc}
}

func (p *Parser) tokenText(tok Token) string {
	if tok.Type == 0 {
		return ""
	}
	if tok.Loc >= 0 && tok.End > tok.Loc && tok.End <= len(p.source) {
		return p.source[tok.Loc:tok.End]
	}
	if tok.Str != "" {
		return tok.Str
	}
	if tok.Type > 0 && tok.Type < 256 {
		return string(rune(tok.Type))
	}
	return ""
}

func (p *Parser) lexerError() *ParseError {
	text := p.tokenText(p.cur)
	var msg string
	if text != "" {
		msg = fmt.Sprintf("%s at or near \"%s\"", p.lexer.Err.Error(), text)
	} else {
		msg = p.lexer.Err.Error()
	}
	return &ParseError{Message: msg, Position: p.cur.Loc}
}
