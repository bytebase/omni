package parser

import (
	"github.com/bytebase/omni/googlesql/ast"
)

// ---------------------------------------------------------------------------
// BigQuery DDL — CREATE MATERIALIZED|APPROX VIEW + ALTER MATERIALIZED|APPROX VIEW
// (parser-ddl-bigquery node)
// ---------------------------------------------------------------------------
//
// Ports the MATERIALIZED and APPROX alternatives of the legacy ANTLR
// create_view_statement (GoogleSQLParser.g4 §2.2) and the BigQuery
// ALTER MATERIALIZED VIEW SET OPTIONS form (the schema_object_kind
// MATERIALIZED VIEW alternative of alter_statement). The plain VIEW alternative
// is owned by the core parser-ddl node (view.go).
//
//	create_view_statement (materialized):
//	  CREATE opt_or_replace? MATERIALIZED RECURSIVE? VIEW opt_if_not_exists?
//	    maybe_dashed_path column_with_options_list? opt_sql_security_clause?
//	    partition_by_clause_prefix_no_hint? cluster_by_clause_prefix_no_hint?
//	    opt_options_list? AS query_or_replica_source
//	create_view_statement (approx):
//	  CREATE opt_or_replace? APPROX RECURSIVE? VIEW opt_if_not_exists?
//	    maybe_dashed_path column_with_options_list? opt_sql_security_clause?
//	    opt_options_list? as_query
//	query_or_replica_source: query | REPLICA OF maybe_dashed_path
//
// ORACLE NOTE — BigQuery-only at the union level (oracle.md): the Spanner
// emulator rejects `CREATE MATERIALIZED VIEW` outright (probed 2026-06-05:
// "Encountered 'MATERIALIZED' while parsing: create_statement"). Triangulated
// against the legacy .g4 + BigQuery truth1 (DDL-011/012/040); proven by the unit
// tests in bq_materialized_view_test.go.

// parseCreateMaterializedView parses a CREATE MATERIALIZED|APPROX [RECURSIVE]
// VIEW statement. The shared CREATE prefix (OR REPLACE) has been consumed by
// parseCreateStmt; cur is at MATERIALIZED or APPROX. scope is unused for these
// forms (the grammar has no opt_create_scope before MATERIALIZED/APPROX), so a
// non-empty scope is a syntax error.
func (p *Parser) parseCreateMaterializedView(create Token, orReplace bool, scope string) (ast.Node, error) {
	if scope != "" {
		// MATERIALIZED/APPROX VIEW take no TEMP/TEMPORARY/PUBLIC/PRIVATE scope.
		return nil, p.syntaxErrorAtCur()
	}

	var kind ast.ViewKind
	switch p.cur.Type {
	case kwMATERIALIZED:
		kind = ast.ViewMaterialized
	case kwAPPROX:
		kind = ast.ViewApprox
	default:
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // MATERIALIZED | APPROX

	stmt := &ast.CreateMaterializedViewStmt{Kind: kind, OrReplace: orReplace}
	stmt.Loc.Start = create.Loc.Start

	// RECURSIVE?
	if p.cur.Type == kwRECURSIVE {
		p.advance()
		stmt.Recursive = true
	}

	if _, err := p.expect(kwVIEW); err != nil {
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

	// column_with_options_list? — `( identifier opt_options_list? (, …)* )`.
	if p.cur.Type == int('(') {
		cols, err := p.parseViewColumnList()
		if err != nil {
			return nil, err
		}
		stmt.Columns = cols
	}

	// opt_sql_security_clause? — SQL SECURITY {INVOKER|DEFINER}.
	if p.cur.Type == kwSQL {
		sec, err := p.parseSQLSecurityClause()
		if err != nil {
			return nil, err
		}
		stmt.SQLSecurity = sec
	}

	// partition_by_clause_prefix_no_hint? / cluster_by_clause_prefix_no_hint? —
	// MATERIALIZED only in the grammar. APPROX has neither; a PARTITION/CLUSTER after
	// an APPROX view falls through to the AS check and is reported there.
	if kind == ast.ViewMaterialized {
		if p.cur.Type == kwPARTITION {
			exprs, err := p.parsePartitionByNoHint()
			if err != nil {
				return nil, err
			}
			stmt.PartitionBy = exprs
		}
		if p.cur.Type == kwCLUSTER {
			exprs, err := p.parseClusterByNoHint()
			if err != nil {
				return nil, err
			}
			stmt.ClusterBy = exprs
		}
	}

	// opt_options_list?
	if p.cur.Type == kwOPTIONS {
		opts, err := p.parseOptionsList()
		if err != nil {
			return nil, err
		}
		stmt.Options = opts
	}

	// AS query_or_replica_source (materialized) | as_query (approx). Required.
	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}
	// query_or_replica_source: query | REPLICA OF maybe_dashed_path. The REPLICA OF
	// alternative is materialized-only; for APPROX a leading REPLICA would be a bad
	// query head and reject in parseQuery.
	if kind == ast.ViewMaterialized && p.cur.Type == kwREPLICA {
		p.advance() // REPLICA
		if _, err := p.expect(kwOF); err != nil {
			return nil, err
		}
		src, err := p.parseTablePath()
		if err != nil {
			return nil, err
		}
		stmt.ReplicaOf = src
		// Trailing OPTIONS(materialized_view_replica_option_list) after AS REPLICA OF
		// (BigQuery DDL-012). The legacy GoogleSQLParser.g4 puts opt_options_list only
		// BEFORE `AS`, so it rejects this trailing form; the BigQuery docs document it,
		// and accepting it is additive (Spanner rejects the whole REPLICA-OF view
		// regardless). Recognized here so a documented replica-OF view with options
		// does not draw a false Diagnose error.
		if p.cur.Type == kwOPTIONS {
			opts, err := p.parseOptionsList()
			if err != nil {
				return nil, err
			}
			stmt.Options = opts
		}
	} else {
		q, err := p.parseQuery()
		if err != nil {
			return nil, err
		}
		stmt.AsQuery = q
	}

	stmt.Loc.End = p.prev.Loc.End
	p.fillSubqueries(stmt)
	return stmt, nil
}

// parseBQAlterMaterializedView parses an ALTER MATERIALIZED|APPROX VIEW statement.
// The ALTER keyword has been consumed by parseAlterStmt; cur is at MATERIALIZED or
// APPROX. The only documented action is SET OPTIONS (DDL-040), but the legacy
// schema_object_kind alter path admits any alter_action; we accept SET OPTIONS
// (the BigQuery-documented form) and reject others as a syntax error.
func (p *Parser) parseBQAlterMaterializedView(alter Token) (ast.Node, error) {
	var kind ast.BQAlterObjectKind
	switch p.cur.Type {
	case kwMATERIALIZED:
		kind = ast.BQAlterMaterializedView
	case kwAPPROX:
		kind = ast.BQAlterApproxView
	default:
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // MATERIALIZED | APPROX
	if _, err := p.expect(kwVIEW); err != nil {
		return nil, err
	}

	stmt := &ast.BQAlterStmt{Object: kind}
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

	// SET OPTIONS options_list.
	if _, err := p.expect(kwSET); err != nil {
		return nil, err
	}
	if p.cur.Type != kwOPTIONS {
		return nil, p.syntaxErrorAtCur()
	}
	opts, err := p.parseOptionsList() // consumes OPTIONS and the parenthesized body
	if err != nil {
		return nil, err
	}
	stmt.SetOptions = opts

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseBQDropMaterializedView parses a DROP MATERIALIZED|APPROX VIEW statement.
// The DROP keyword has been consumed by parseDropStmt; cur is at MATERIALIZED or
// APPROX. Grammar: `DROP schema_object_kind opt_if_exists? path opt_drop_mode?`
// where schema_object_kind ⊇ {MATERIALIZED VIEW, APPROX VIEW}. A documented DROP
// MATERIALIZED VIEW (DDL-048) takes no RESTRICT/CASCADE, so the trailing
// opt_drop_mode is left for the EOF check to reject if present.
func (p *Parser) parseBQDropMaterializedView(drop Token) (ast.Node, error) {
	kind := ast.BQDropMaterializedView
	if p.cur.Type == kwAPPROX {
		kind = ast.BQDropApproxView
	}
	p.advance() // MATERIALIZED | APPROX
	if _, err := p.expect(kwVIEW); err != nil {
		return nil, err
	}
	return p.finishBQDropObject(drop, kind, false /*allowFunctionParams*/)
}
