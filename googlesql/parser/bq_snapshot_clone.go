package parser

import (
	"github.com/bytebase/omni/googlesql/ast"
)

// ---------------------------------------------------------------------------
// BigQuery DDL — CREATE SNAPSHOT TABLE + DROP SNAPSHOT TABLE
// (parser-ddl-bigquery node)
// ---------------------------------------------------------------------------
//
// Ports the legacy ANTLR create_snapshot_statement (GoogleSQLParser.g4 §2.2) and
// the `DROP SNAPSHOT TABLE` alternative of drop_statement. (CREATE TABLE … CLONE
// / COPY / LIKE — DDL-006/007/009 — are part of the core CREATE TABLE grammar and
// are owned by the parser-ddl node in create_table.go; this file is only the
// SNAPSHOT object.)
//
//	create_snapshot_statement:
//	  CREATE opt_or_replace? SNAPSHOT (TABLE | schema_object_kind) opt_if_not_exists?
//	    maybe_dashed_path CLONE clone_data_source opt_options_list?
//	  clone_data_source: maybe_dashed_path opt_at_system_time? where_clause?
//	  opt_at_system_time: FOR SYSTEM_TIME AS OF expr   (also `FOR SYSTEM TIME AS OF`)
//	drop_statement: DROP SNAPSHOT TABLE opt_if_exists? maybe_dashed_path
//
// ORACLE NOTE — BigQuery-only at the union level (oracle.md). Probed 2026-06-05:
// `CREATE SNAPSHOT TABLE … CLONE …` REJECTS on the Spanner emulator
// ("Encountered 'SNAPSHOT' while parsing: create_statement"). Triangulated
// against the legacy .g4 + BigQuery truth1 (DDL-008/045); proven by the unit
// tests in bq_snapshot_clone_test.go.

// parseCreateSnapshot parses a CREATE SNAPSHOT TABLE statement. The shared CREATE
// prefix (OR REPLACE) has been consumed by parseCreateStmt; cur is at SNAPSHOT.
func (p *Parser) parseCreateSnapshot(create Token, orReplace bool) (ast.Node, error) {
	p.advance() // SNAPSHOT

	stmt := &ast.CreateSnapshotStmt{OrReplace: orReplace}
	stmt.Loc.Start = create.Loc.Start

	// SNAPSHOT (TABLE | schema_object_kind) — the documented form is SNAPSHOT TABLE;
	// the grammar also admits a schema_object_kind, but TABLE is the BigQuery form.
	if _, err := p.expect(kwTABLE); err != nil {
		return nil, err
	}

	ifNotExists, err := p.parseIfNotExists()
	if err != nil {
		return nil, err
	}
	stmt.IfNotExists = ifNotExists

	name, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// CLONE clone_data_source.
	if _, err := p.expect(kwCLONE); err != nil {
		return nil, err
	}
	src, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	stmt.CloneSource = src

	// opt_at_system_time? — FOR SYSTEM_TIME AS OF expr.
	if p.cur.Type == kwFOR {
		ts, err := p.parseForSystemTimeExpr()
		if err != nil {
			return nil, err
		}
		stmt.ForSystemTime = ts
	}

	// where_clause? — WHERE expr (part of clone_data_source).
	if p.cur.Type == kwWHERE {
		p.advance() // WHERE
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Where = expr
	}

	// opt_options_list?
	if p.cur.Type == kwOPTIONS {
		opts, err := p.parseOptionsList()
		if err != nil {
			return nil, err
		}
		stmt.Options = opts
	}

	stmt.Loc.End = p.prev.Loc.End
	p.fillSubqueries(stmt)
	return stmt, nil
}

// parseForSystemTimeExpr parses an opt_at_system_time clause and RETURNS its time
// expression (as opposed to from_join.go's parseForSystemTime, which attaches it
// to a TableExpr):
//
//	FOR SYSTEM_TIME AS OF expr
//	| FOR SYSTEM TIME AS OF expr
//
// cur is at FOR. The lexer emits SYSTEM_TIME as one token (kwSYSTEM_TIME) or
// SYSTEM + TIME as two; both spellings are accepted.
func (p *Parser) parseForSystemTimeExpr() (ast.Node, error) {
	p.advance() // FOR
	switch p.cur.Type {
	case kwSYSTEM_TIME:
		p.advance() // SYSTEM_TIME
	case kwSYSTEM:
		p.advance() // SYSTEM
		if _, err := p.expect(kwTIME); err != nil {
			return nil, err
		}
	default:
		return nil, p.syntaxErrorAtCur()
	}
	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwOF); err != nil {
		return nil, err
	}
	return p.parseExpr()
}

// parseBQDropSnapshot parses a DROP SNAPSHOT TABLE statement. The DROP keyword has
// been consumed by parseDropStmt; cur is at SNAPSHOT. Grammar:
// `DROP SNAPSHOT TABLE opt_if_exists? maybe_dashed_path` (DDL-045).
func (p *Parser) parseBQDropSnapshot(drop Token) (ast.Node, error) {
	p.advance() // SNAPSHOT
	if _, err := p.expect(kwTABLE); err != nil {
		return nil, err
	}
	return p.finishBQDropObject(drop, ast.BQDropSnapshotTable, false /*allowFunctionParams*/)
}
