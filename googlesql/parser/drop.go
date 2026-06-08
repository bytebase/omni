package parser

import (
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
			return p.parseBQDropFunction(drop)
		}
		return p.parseDropObject(drop, ast.DropTable, false /*external*/)
	case kwVIEW:
		return p.parseDropObject(drop, ast.DropView, false)
	case kwINDEX:
		return p.parseDropIndex(drop)
	case kwSEARCH, kwVECTOR:
		// DROP {SEARCH|VECTOR} INDEX … ON table — BigQuery-only
		// (bq_search_vector_index.go).
		return p.parseBQDropSearchVectorIndex(drop)
	case kwSCHEMA:
		return p.parseDropObject(drop, ast.DropSchema, false)
	case kwDATABASE:
		return p.parseDropObject(drop, ast.DropDatabase, false)
	case kwEXTERNAL:
		// DROP EXTERNAL SCHEMA (BigQuery) — owned here as a SCHEMA-family drop;
		// DROP EXTERNAL TABLE [FUNCTION] is dialect-node territory.
		if p.peekNext().Type == kwSCHEMA {
			p.advance() // EXTERNAL
			return p.parseDropObject(drop, ast.DropSchema, true /*external*/)
		}
		return p.unsupported("DROP EXTERNAL")
	// --- BigQuery-only object drops (parser-ddl-bigquery node) ---
	case kwFUNCTION:
		return p.parseBQDropFunction(drop)
	case kwPROCEDURE:
		return p.parseBQDropProcedure(drop)
	case kwMATERIALIZED, kwAPPROX:
		return p.parseBQDropMaterializedView(drop)
	case kwSNAPSHOT:
		return p.parseBQDropSnapshot(drop)
	case kwROW:
		// DROP ROW ACCESS POLICY … ON table (bq_row_access_policy.go).
		return p.parseDropRowAccessPolicy(drop)
	case kwALL:
		// DROP ALL ROW ACCESS POLICIES ON table
		// (drop_all_row_access_policies_statement; bq_row_access_policy.go).
		return p.parseDropRowAccessPolicy(drop)
	// --- Spanner-only object DROPs (parser-ddl-spanner node) ---
	case kwSEQUENCE:
		// DROP SEQUENCE [IF EXISTS] name (spanner_sequence.go).
		return p.parseDropSequence(drop)
	default:
		// Spanner objects whose leading word lexes as a BARE IDENTIFIER (CHANGE
		// STREAM, LOCALITY GROUP) and ROLE — matched by spelling before the
		// generic-entity fallback (parser-ddl-spanner node).
		if p.tokIsWord(p.cur, "CHANGE") && p.tokIsWord(p.peekNext(), "STREAM") {
			return p.parseDropChangeStream(drop) // spanner_change_stream.go
		}
		if p.tokIsWord(p.cur, "LOCALITY") && p.tokIsWord(p.peekNext(), "GROUP") {
			return p.parseDropLocalityGroup(drop) // spanner_sequence.go
		}
		if p.curIsWord("ROLE") {
			// DROP ROLE name (spanner_schema_role.go). First-class — supersedes the
			// generic-entity drop the bare `ROLE` identifier would otherwise get.
			return p.parseDropRole(drop)
		}
		// DROP <generic-entity> (CAPACITY / RESERVATION / ASSIGNMENT — DDL-053; the
		// keyword lexes as an identifier) routes to the generic-entity drop
		// (bq_capacity.go). DROP MODEL / CONNECTION / CONSTANT / PRIVILEGE
		// RESTRICTION and the bare-identifier dispatch-coverage probe otherwise emit
		// the not-yet-supported diagnostic.
		if isGenericEntityType(p.cur.Type) {
			return p.parseBQDropEntity(drop)
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
	// (everything here EXCEPT a bare TABLE) carries it.
	if kind != ast.DropTable {
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
