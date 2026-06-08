package parser

import (
	"github.com/bytebase/omni/googlesql/ast"
)

// ---------------------------------------------------------------------------
// BigQuery DDL — CREATE [SEARCH|VECTOR] INDEX / ALTER VECTOR INDEX REBUILD /
// DROP [SEARCH|VECTOR] INDEX (parser-ddl-bigquery node)
// ---------------------------------------------------------------------------
//
// Ports the SEARCH / VECTOR index_type variants of the legacy ANTLR
// create_index_statement (GoogleSQLParser.g4 §2.2/§2.6), the BigQuery
// `DROP index_type INDEX … ON …` alternative of drop_statement, and the
// BigQuery-documented `ALTER VECTOR INDEX … REBUILD`. The plain / UNIQUE /
// NULL_FILTERED CREATE INDEX is owned by the core parser-ddl node
// (create_index.go).
//
//	create_index_statement (search/vector):
//	  CREATE opt_or_replace? UNIQUE? opt_spanner_null_filtered? index_type INDEX
//	    opt_if_not_exists? path_expression on_path_expression as_alias?
//	    index_unnest_expression_list? index_order_by_and_options index_storing_list?
//	    opt_create_index_statement_suffix?
//	  index_type: SEARCH | VECTOR
//	  index_order_by_and_options: ( column_ordering_and_options_expr, … ) | index_all_columns
//	  index_all_columns: ( ALL COLUMNS opt_with_column_options? )
//	drop_statement (search/vector):
//	  DROP index_type INDEX opt_if_exists? path_expression on_path_expression?
//
// ORACLE NOTE — BigQuery-only at the union level (oracle.md). Probed 2026-06-05:
// a bare `CREATE VECTOR INDEX i ON t(c) OPTIONS(…)` ACCEPTS on the Spanner
// emulator (Spanner has a vector-index grammar), but `CREATE SEARCH INDEX i ON
// t(ALL COLUMNS)` REJECTS ("Expecting ')' but found 'ALL'") and `DROP {SEARCH|
// VECTOR} INDEX … ON …` / `ALTER VECTOR INDEX … REBUILD` all REJECT (Spanner's
// shapes differ). So the Spanner verdict is non-authoritative here; triangulated
// against the legacy .g4 + BigQuery truth1 (DDL-022/023/041/054); proven by the
// unit tests in bq_search_vector_index_test.go.
//
// DIVERGENCE (flagged) — `ALTER VECTOR INDEX … REBUILD` (DDL-041) is a documented
// BigQuery form that the legacy GoogleSQLParser.g4 does NOT parse (no REBUILD
// token anywhere; ALTER VECTOR INDEX rejects in the .g4 because VECTOR is neither
// a table_or_table_function nor a schema_object_kind nor a generic_entity_type).
// We implement the documented behavior (accept it) since the omni parser targets
// the BigQuery+Spanner union and this is strictly additive over one of the
// legacy grammar's known gaps. See the divergence ledger entry for this node.

// parseCreateSearchVectorIndex parses a CREATE SEARCH|VECTOR INDEX statement. The
// shared CREATE prefix (OR REPLACE / scope) has been consumed by parseCreateStmt;
// cur is at SEARCH or VECTOR. scope is not part of create_index_statement, so a
// non-empty scope is a syntax error.
func (p *Parser) parseCreateSearchVectorIndex(create Token, orReplace bool, scope string) (ast.Node, error) {
	if scope != "" {
		return nil, p.syntaxErrorAtCur()
	}

	isVector := p.cur.Type == kwVECTOR
	p.advance() // SEARCH | VECTOR

	stmt := &ast.SearchVectorIndexStmt{IsVector: isVector, OrReplace: orReplace}
	stmt.Loc.Start = create.Loc.Start

	if _, err := p.expect(kwINDEX); err != nil {
		return nil, err
	}

	ifNotExists, err := p.parseIfNotExists()
	if err != nil {
		return nil, err
	}
	stmt.IfNotExists = ifNotExists

	// index name (path_expression).
	name, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// on_path_expression: ON path_expression.
	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	table, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	stmt.Table = table

	// index_order_by_and_options: `( column_ordering_and_options_expr, … )` or
	// index_all_columns `( ALL COLUMNS opt_with_column_options? )`.
	if p.cur.Type != int('(') {
		return nil, p.syntaxErrorAtCur()
	}
	if p.peekNext().Type == kwALL {
		// index_all_columns: `( ALL COLUMNS opt_with_column_options? )`. The grammar
		// REQUIRES COLUMNS after ALL, so validate the `( ALL COLUMNS` prefix rather
		// than blindly skipping balanced parens (a bare `( ALL )` must reject). The
		// remaining body (the optional WITH COLUMN OPTIONS(...) and the closing ')')
		// is captured structurally — the per-column options are not load-bearing for
		// the query-span consumer.
		if err := p.parseIndexAllColumns(stmt); err != nil {
			return nil, err
		}
	} else {
		keys, err := p.parseKeyPartList(false /*index element*/)
		if err != nil {
			return nil, err
		}
		stmt.Keys = keys
	}

	// index_storing_list? — STORING ( expr, … ).
	if p.cur.Type == kwSTORING {
		p.advance() // STORING
		if _, err := p.expect(int('(')); err != nil {
			return nil, err
		}
		exprs, _, err := p.parseExprListThroughParen()
		if err != nil {
			return nil, err
		}
		stmt.Storing = exprs
	}

	// opt_create_index_statement_suffix?:
	//   partition_by_clause_prefix_no_hint opt_options_list?
	//   | opt_options_list? spanner_index_interleave_clause   (not valid for search/vector — left unhandled)
	//   | opt_options_list
	if p.cur.Type == kwPARTITION {
		exprs, err := p.parsePartitionByNoHint()
		if err != nil {
			return nil, err
		}
		stmt.PartitionBy = exprs
	}
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

// parseIndexAllColumns parses an index_all_columns body:
//
//	'(' ALL COLUMNS opt_with_column_options? ')'
//	opt_with_column_options: WITH COLUMN OPTIONS '(' column_ordering_and_options_expr (, …) ')'
//
// cur is at the opening '(' (the next token is ALL). It sets stmt.AllColumns and
// validates that COLUMNS follows ALL (a bare `( ALL )` rejects). The optional
// WITH COLUMN OPTIONS(...) parenthesized run is consumed structurally.
func (p *Parser) parseIndexAllColumns(stmt *ast.SearchVectorIndexStmt) error {
	p.advance() // '('
	if _, err := p.expect(kwALL); err != nil {
		return err
	}
	if _, err := p.expect(kwCOLUMNS); err != nil {
		return err
	}
	stmt.AllColumns = true

	// opt_with_column_options?: WITH COLUMN OPTIONS ( … ).
	if p.cur.Type == kwWITH {
		p.advance() // WITH
		if _, err := p.expect(kwCOLUMN); err != nil {
			return err
		}
		if _, err := p.expect(kwOPTIONS); err != nil {
			return err
		}
		// all_column_column_options is a parenthesized column_ordering_and_options_expr
		// list; consume it structurally (its per-column options are not load-bearing).
		if p.cur.Type != int('(') {
			return p.syntaxErrorAtCur()
		}
		if err := p.skipBalancedParens(); err != nil {
			return err
		}
	}

	if _, err := p.expect(int(')')); err != nil {
		return err
	}
	return nil
}

// parseBQAlterVectorIndex parses an `ALTER VECTOR INDEX [IF EXISTS] index_name ON
// table_name REBUILD [OPTIONS(…)]` statement (DDL-041). The ALTER keyword has been
// consumed by parseAlterStmt; cur is at VECTOR.
//
// This form is NOT in the legacy GoogleSQLParser.g4 (see the file header
// divergence note); it is implemented from the BigQuery docs as an additive union
// form.
func (p *Parser) parseBQAlterVectorIndex(alter Token) (ast.Node, error) {
	p.advance() // VECTOR
	if _, err := p.expect(kwINDEX); err != nil {
		return nil, err
	}

	stmt := &ast.BQAlterStmt{Object: ast.BQAlterVectorIndex}
	stmt.Loc.Start = alter.Loc.Start

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

	// ON table_name.
	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	table, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	stmt.OnTable = table

	// REBUILD (the keyword lexes as an identifier; match it case-insensitively by
	// source spelling via curIsWord).
	if !p.curIsWord("REBUILD") {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // REBUILD
	stmt.Rebuild = true

	// opt OPTIONS(…) (index_rebuild_option_list).
	if p.cur.Type == kwOPTIONS {
		opts, err := p.parseOptionsList()
		if err != nil {
			return nil, err
		}
		stmt.SetOptions = opts
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseBQDropSearchVectorIndex parses a `DROP {SEARCH|VECTOR} INDEX [IF EXISTS]
// path ON table` statement. The DROP keyword has been consumed by parseDropStmt;
// cur is at SEARCH or VECTOR. Grammar (drop_statement index_type alt):
// `DROP index_type INDEX opt_if_exists? path_expression on_path_expression?`
// (DDL-054). The on_path_expression is documented as required for BigQuery
// SEARCH/VECTOR index drops, but the grammar marks it optional; we accept both
// (an absent ON yields OnTable nil).
func (p *Parser) parseBQDropSearchVectorIndex(drop Token) (ast.Node, error) {
	kind := ast.BQDropSearchIndex
	if p.cur.Type == kwVECTOR {
		kind = ast.BQDropVectorIndex
	}
	p.advance() // SEARCH | VECTOR
	if _, err := p.expect(kwINDEX); err != nil {
		return nil, err
	}

	stmt := &ast.BQDropStmt{Object: kind}
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

	// on_path_expression? — ON path.
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
