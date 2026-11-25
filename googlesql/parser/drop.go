package parser

import (
	"strings"

	"github.com/bytebase/omni/googlesql/ast"
)

// ---------------------------------------------------------------------------
// Core DDL — DROP (parser-ddl node)
// ---------------------------------------------------------------------------
//
// Ports the table-like-object alternatives of the legacy ANTLR drop_statement
// (GoogleSQLParser.g4 §2.3):
//
//	DROP index_type INDEX opt_if_exists? path_expression on_path_expression?
//	DROP table_or_table_function opt_if_exists? maybe_dashed_path_expression …
//	DROP SNAPSHOT TABLE …
//	DROP schema_object_kind opt_if_exists? path_expression opt_function_parameters? opt_drop_mode?
//
// where schema_object_kind ⊇ { INDEX, VIEW, SCHEMA, DATABASE, TABLE-via-the-
// table alt, … }. This node owns DROP for TABLE / VIEW / INDEX / SCHEMA /
// DATABASE (incl. the BigQuery `DROP {SEARCH|VECTOR} INDEX … ON table`, the
// `DROP EXTERNAL SCHEMA`, and `DROP SCHEMA … RESTRICT|CASCADE` forms). Object
// kinds owned elsewhere (FUNCTION / TABLE FUNCTION / PROCEDURE / MODEL /
// CONNECTION / CONSTANT / MATERIALIZED|APPROX VIEW / SNAPSHOT TABLE / ROW ACCESS
// POLICY / PRIVILEGE RESTRICTION / generic entity) route to the unsupported stub.
//
// Verified against the live Spanner emulator (drop_oracle_test.go): DROP TABLE /
// VIEW / INDEX / SCHEMA [IF EXISTS] all accept.

// parseDropStmt parses a DROP statement. The DROP keyword has NOT yet been
// consumed (parseStmt peeks it). It dispatches on the object keyword(s).
func (p *Parser) parseDropStmt() (ast.Node, error) {
	drop := p.advance() // consume DROP

	switch p.cur.Type {
	case kwTABLE:
		// DROP TABLE — but DROP TABLE FUNCTION is BigQuery-only
		// (bq_function_procedure.go).
		if p.peekNext().Type == kwFUNCTION {
			return nil, p.syntaxErrorAtCur()
		}
		return p.parseDropObject(drop, ast.DropTable, false /*external*/)
	case kwVIEW:
		return p.parseDropObject(drop, ast.DropView, false)
	case kwINDEX:
		return p.parseDropIndex(drop)
	case kwSEARCH, kwVECTOR:
		// DROP {SEARCH|VECTOR} INDEX … ON table — BigQuery-only ON table
		// (bq_search_vector_index.go); Spanner allows bare DROP {SEARCH|VECTOR} INDEX name.
		return p.parseBQDropSearchVectorIndex(drop)
	case kwSCHEMA:
		return p.parseDropObject(drop, ast.DropSchema, false)
	case kwDATABASE:
		return nil, p.syntaxErrorAtCur()
	case kwEXTERNAL:
		return nil, p.syntaxErrorAtCur()
	case kwMODEL:
		return p.parseDropObject(drop, ast.DropModel, false)
	case kwCONNECTION:
		return p.unsupported("DROP CONNECTION")
	case kwCONSTANT:
		return p.unsupported("DROP CONSTANT")
	case kwPRIVILEGE:
		return p.unsupported("DROP PRIVILEGE")
	// --- BigQuery-only object drops (parser-ddl-bigquery node) ---
	case kwFUNCTION:
		return p.parseBQDropFunction(drop)
	case kwPROCEDURE:
		return nil, p.syntaxErrorAtCur()
	case kwMATERIALIZED, kwAPPROX:
		return nil, p.syntaxErrorAtCur()
	case kwSNAPSHOT:
		return nil, p.syntaxErrorAtCur()
	case kwROW:
		return nil, p.syntaxErrorAtCur()
	case kwALL:
		return nil, p.syntaxErrorAtCur()
	// --- Spanner-only object DROPs (parser-ddl-spanner node) ---
	case kwSEQUENCE:
		// DROP SEQUENCE [IF EXISTS] name (spanner_sequence.go).
		return p.parseDropSequence(drop)
	case kwPROPERTY:
		// DROP PROPERTY GRAPH [IF EXISTS] name [RESTRICT|CASCADE] (schema_object_kind).
		if p.peekNext().Type == kwGRAPH {
			p.advance() // consume PROPERTY; parseDropObject consumes GRAPH
			return p.parseDropObject(drop, ast.DropPropertyGraph, false)
		}
		return p.unsupported("DROP")
	default:
		// Spanner objects whose leading word lexes as a BARE IDENTIFIER (CHANGE
		// STREAM, LOCALITY GROUP, PLACEMENT) and ROLE — matched by spelling before
		// the generic-entity fallback (parser-ddl-spanner node).
		if p.tokIsWord(p.cur, "CHANGE") && p.tokIsWord(p.peekNext(), "STREAM") {
			return p.parseDropChangeStream(drop) // spanner_change_stream.go
		}
		if p.tokIsWord(p.cur, "LOCALITY") && p.tokIsWord(p.peekNext(), "GROUP") {
			return p.parseDropLocalityGroup(drop) // spanner_sequence.go
		}
		if p.tokIsWord(p.cur, "PLACEMENT") {
			return p.parseDropPlacement(drop) // spanner_placement.go
		}
		if p.curIsWord("ROLE") {
			// DROP ROLE name (spanner_schema_role.go). First-class — supersedes the
			// generic-entity drop the bare `ROLE` identifier would otherwise get.
			return p.parseDropRole(drop)
		}
		return p.unsupported("DROP")
	}
}

// parseDropObject parses `DROP <object-keyword> opt_if_exists? path_expression
// opt_drop_mode?` where the object keyword has already been consumed-by-peek but
// NOT advanced. It consumes the object keyword, the IF EXISTS, the path, and a
// trailing RESTRICT|CASCADE (opt_drop_mode).
//
// opt_drop_mode is attached in the .g4 ONLY to the `schema_object_kind`
// alternative (VIEW / SCHEMA / DATABASE / INDEX / EXTERNAL SCHEMA / …). A bare
// DROP TABLE is the separate `table_or_table_function` alternative, which carries
// NO opt_drop_mode — so `DROP TABLE T CASCADE|RESTRICT` is a syntax error (the
// live Spanner oracle rejects it; real BigQuery DROP TABLE has no CASCADE/RESTRICT
// either). We therefore consume opt_drop_mode only when kind != DropTable; for a
// TABLE, a trailing CASCADE/RESTRICT is left unconsumed and parseSingle's EOF
// check reports it as the syntax error it is.
func (p *Parser) parseDropObject(drop Token, kind ast.DropObjectKind, external bool) (ast.Node, error) {
	p.advance() // consume the object keyword (TABLE / VIEW / SCHEMA / DATABASE)

	stmt := &ast.DropStmt{Object: kind, External: external}
	stmt.Loc.Start = drop.Loc.Start

	ifExists, err := p.parseIfExists()
	if err != nil {
		return nil, err
	}
	stmt.IfExists = ifExists

	name, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// opt_drop_mode?  — RESTRICT | CASCADE. Only the schema_object_kind alternative
	// (excluding bare TABLE and BigQuery CONNECTION) carries it.
	if kind != ast.DropTable && kind != ast.DropConnection {
		switch p.cur.Type {
		case kwRESTRICT:
			p.advance()
			stmt.DropMode = "RESTRICT"
		case kwCASCADE:
			p.advance()
			stmt.DropMode = "CASCADE"
		}
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseDropIndex parses a plain `DROP INDEX opt_if_exists? path_expression`. The
// INDEX keyword is the current token.
//
// A plain DROP INDEX is the `schema_object_kind` (INDEX) alternative, which has
// NO on_path_expression: the trailing `ON <table>` belongs ONLY to the
// `DROP index_type INDEX … on_path_expression?` alternative, where index_type is
// SEARCH | VECTOR (already routed to the dialect stub by parseDropStmt before we
// get here). So `DROP INDEX idx ON T` is a syntax error (the live Spanner oracle
// rejects it); a bare ON is left unconsumed and parseSingle's EOF check reports
// it. The opt_drop_mode (`DROP INDEX idx CASCADE|RESTRICT`) is part of
// schema_object_kind and is consumed below.
func (p *Parser) parseDropIndex(drop Token) (ast.Node, error) {
	p.advance() // consume INDEX

	stmt := &ast.DropStmt{Object: ast.DropIndex}
	stmt.Loc.Start = drop.Loc.Start

	ifExists, err := p.parseIfExists()
	if err != nil {
		return nil, err
	}
	stmt.IfExists = ifExists

	name, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// opt_drop_mode?  — RESTRICT | CASCADE (schema_object_kind INDEX).
	switch p.cur.Type {
	case kwRESTRICT:
		p.advance()
		stmt.DropMode = "RESTRICT"
	case kwCASCADE:
		p.advance()
		stmt.DropMode = "CASCADE"
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// UNDROP (parser-ddl node)
// ---------------------------------------------------------------------------
//
// Grammar (googlesql.tm:2934):
//
//	undrop_statement: UNDROP schema_object_kind opt_if_not_exists path_expression
//	                  at_system_time? opt_options_list
//
// The rule's own action accepts ONE object kind — SCHEMA — and reports
// "UNDROP <kind> is not supported" for every other kind, so only `UNDROP SCHEMA`
// is a statement the engine parses to completion. The parser follows that: any
// other kind is rejected at the kind keyword, as ZetaSQL does.
//
// This is why no UNDROP row of testdata/compat.tsv closes: every row the
// generator wrote is `UNDROP CONNECTION …` (the first alternative of
// schema_object_kind), which the grammar's action rejects, while the emulator
// answers "Statement not supported: UndropStmt" from dispatch and never parses
// the body at all — it accepts `UNDROP CONNECTION` with no name too.

// parseUndropStmt parses an undrop_statement. UNDROP is the current token.
func (p *Parser) parseUndropStmt() (ast.Node, error) {
	undrop := p.advance() // UNDROP

	// schema_object_kind — only SCHEMA survives the rule's action.
	if !isIdentifierStart(p.cur.Type) {
		return nil, p.syntaxErrorAtCur()
	}
	if p.cur.Type != kwSCHEMA {
		// The rule's action: "UNDROP <kind> is not supported".
		return nil, &ParseError{
			Position: p.cur.Loc.Start,
			End:      p.cur.Loc.End,
			Message:  "syntax error: UNDROP " + strings.ToUpper(p.tokenSource(p.cur)) + " is not supported",
		}
	}
	kindTok := p.advance() // SCHEMA
	stmt := &ast.UndropStmt{ObjectKind: "SCHEMA", Loc: kindTok.Loc}
	stmt.Loc.Start = undrop.Loc.Start

	// opt_if_not_exists.
	if p.cur.Type == kwIF {
		p.advance() // IF
		if _, err := p.expect(kwNOT); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfNotExists = true
	}

	// path_expression.
	name, err := p.parsePathExpr()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// at_system_time? — FOR SYSTEM_TIME AS OF expr.
	if p.cur.Type == kwFOR {
		ts, err := p.parseForSystemTimeExpr()
		if err != nil {
			return nil, err
		}
		stmt.ForSystemTime = ts
	}

	// opt_options_list.
	if p.cur.Type == kwOPTIONS {
		opts, err := p.parseOptionsList()
		if err != nil {
			return nil, err
		}
		stmt.Options = opts
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}
