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
//	  CREATE opt_or_replace? SNAPSHOT (TABLE | SCHEMA) opt_if_not_exists?
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
