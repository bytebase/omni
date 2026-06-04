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
		// DROP TABLE — but DROP TABLE FUNCTION is dialect-node territory.
		if p.peekNext().Type == kwFUNCTION {
			return p.unsupported("DROP TABLE FUNCTION")
		}
		return p.parseDropObject(drop, ast.DropTable, false /*external*/)
	case kwVIEW:
		return p.parseDropObject(drop, ast.DropView, false)
	case kwINDEX:
		return p.parseDropIndex(drop)
	case kwSEARCH, kwVECTOR:
		// DROP {SEARCH|VECTOR} INDEX … ON table — dialect-specific.
		return p.unsupported("DROP " + p.cur.Str + " INDEX")
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
	default:
		// DROP FUNCTION / PROCEDURE / MODEL / SNAPSHOT TABLE / ROW ACCESS POLICY /
		// PRIVILEGE RESTRICTION / generic entity / a bare-identifier object — not
		// owned here.
		return p.unsupported("DROP")
	}
}

// parseDropObject parses `DROP <object-keyword> opt_if_exists? path_expression
// opt_drop_mode?` where the object keyword has already been consumed-by-peek but
// NOT advanced. It consumes the object keyword, the IF EXISTS, the path, and a
// trailing RESTRICT|CASCADE (opt_drop_mode — DROP SCHEMA CASCADE/RESTRICT).
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

	// opt_drop_mode?  — RESTRICT | CASCADE (DROP SCHEMA).
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

// parseDropIndex parses `DROP INDEX opt_if_exists? path_expression
// on_path_expression?`. The INDEX keyword is the current token. The optional
// `ON <table>` is the BigQuery search/vector-index form; for a plain Spanner
// DROP INDEX it is absent.
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

	// on_path_expression?  — ON <table>.
	if p.cur.Type == kwON {
		p.advance() // ON
		onTable, err := p.parseTablePath()
		if err != nil {
			return nil, err
		}
		stmt.OnTable = onTable
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}
