package parser

import (
	"github.com/bytebase/omni/googlesql/ast"
)

// ---------------------------------------------------------------------------
// Core DDL — CREATE INDEX (parser-ddl node)
// ---------------------------------------------------------------------------
//
// Ports the legacy ANTLR create_index_statement (GoogleSQLParser.g4 §2.2/§2.6):
//
//	CREATE opt_or_replace? UNIQUE? opt_spanner_null_filtered? index_type? INDEX
//	  opt_if_not_exists? path_expression on_path_expression as_alias?
//	  index_unnest_expression_list? index_order_by_and_options index_storing_list?
//	  opt_create_index_statement_suffix?
//
// This node owns the plain / UNIQUE / NULL_FILTERED forms. The SEARCH / VECTOR
// index_type variants are dialect-specific (parser-ddl-{bigquery,spanner}); a
// SEARCH/VECTOR token after the qualifiers routes to the unsupported stub.
//
// Verified against the live Spanner emulator in create_index_oracle_test.go:
// plain/UNIQUE/NULL_FILTERED indexes, ASC/DESC key parts, STORING, the Spanner
// `, INTERLEAVE IN <table>` suffix, and the partial-index WHERE filter all accept.

// parseCreateIndex parses a CREATE INDEX body after the shared CREATE prefix (OR
// REPLACE / scope) has been consumed. cur is at one of UNIQUE / NULL_FILTERED /
// INDEX (the qualifiers and the INDEX keyword are consumed here).
func (p *Parser) parseCreateIndex(create Token, orReplace bool, scope string) (ast.Node, error) {
	// scope (TEMP/etc.) is not part of create_index_statement; if one was consumed
	// before the INDEX path, it is not a valid index statement. Treat a non-empty
	// scope as a syntax error at the current token (the qualifier/INDEX keyword).
	if scope != "" {
		return nil, p.syntaxErrorAtCur()
	}

	stmt := &ast.CreateIndexStmt{OrReplace: orReplace}
	stmt.Loc.Start = create.Loc.Start

	// UNIQUE? opt_spanner_null_filtered? — in that grammar order, but tolerate
	// either order defensively (both are simple flags).
	for {
		switch p.cur.Type {
		case kwUNIQUE:
			p.advance()
			stmt.Unique = true
			continue
		case kwNULL_FILTERED:
			p.advance()
			stmt.NullFiltered = true
			continue
		}
		break
	}

	// index_type? — SEARCH | VECTOR. These are dialect-specific; route to stub.
	if p.cur.Type == kwSEARCH || p.cur.Type == kwVECTOR {
		return p.unsupported("CREATE " + p.cur.Str + " INDEX")
	}

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

	// as_alias? index_unnest_expression_list? — the table alias + UNNEST list are
	// rare BigQuery index forms over array columns. They are not in the
	// Spanner-authoritative set and add no query-span value here; recognized only
	// to the extent the key list below disambiguates. We do NOT consume an alias /
	// UNNEST list (the common case has neither and the key list '(' follows the
	// table directly); if present they would surface as a syntax error at the key
	// list, which is acceptable for this core node (the dialect index nodes own
	// the exotic array-index forms).

	// index_order_by_and_options: `( column_ordering_and_options_expr, … )` or
	// `( ALL COLUMNS … )`.
	if p.cur.Type != int('(') {
		return nil, p.syntaxErrorAtCur()
	}
	if p.peekNext().Type == kwALL {
		return nil, &ParseError{
			Position: p.peekNext().Loc.Start, End: p.peekNext().Loc.End,
			Message: "syntax error: ALL COLUMNS is not supported in Spanner indexes",
		}
	} else {
		keys, err := p.parseKeyPartListWithOptions(false /*index element*/, true /*allowSpannerParenExpr*/)
		if err != nil {
			return nil, err
		}
		stmt.Keys = keys
	}

	// index_storing_list?  — STORING ( expr, … ). In Spanner, STORING takes column identifiers.
	if p.cur.Type == kwSTORING {
		p.advance() // STORING
		exprs, err := p.parseStoringIdentList()
		if err != nil {
			return nil, err
		}
		stmt.Storing = exprs
	}

	// opt_create_index_statement_suffix?:
	//   partition_by_clause_prefix_no_hint opt_options_list?
	//   | opt_options_list? spanner_index_interleave_clause
	//   | opt_options_list
	if err := p.parseCreateIndexSuffix(stmt); err != nil {
		return nil, err
	}

	// NULL_FILTERED and `, INTERLEAVE IN` are Spanner-only index forms (BigQuery
	// rejects both); require Spanner-valid index keys (column identifiers or
	// parenthesized expressions, never bare literals such as `NULL` or
	// `ALL COLUMNS`) so hybrids rejected by both engines are not accepted.
	if stmt.NullFiltered || stmt.Interleave != nil {
		if stmt.AllColumns {
			return nil, &ParseError{
				Position: stmt.Loc.Start, End: p.prev.Loc.End,
				Message: "syntax error: ALL COLUMNS is not supported in Spanner indexes",
			}
		}
		for _, k := range stmt.Keys {
			if !isSpannerIndexKey(k) {
				return nil, &ParseError{
					Position: k.Loc.Start, End: k.Loc.End,
					Message: "syntax error: Spanner index key must be a column name or parenthesized expression",
				}
			}
		}
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// isSpannerIndexKey reports whether k is a column identifier or parenthesized
// expression, as required by Spanner ordinary CREATE INDEX keys.
func isSpannerIndexKey(k *ast.KeyPart) bool {
	if k.Name != "" {
		return true
	}
	_, ok := k.Expr.(*ast.ParenExpr)
	return ok
}

// parseStoringIdentList parses a Spanner index STORING column identifier list
// `( identifier (, identifier)* )` preserving source locations on each Identifier.
func (p *Parser) parseStoringIdentList() ([]ast.Node, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	var exprs []ast.Node
	for {
		tok, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}
		name, err := p.identifierText(tok)
		if err != nil {
			return nil, err
		}
		exprs = append(exprs, &ast.Identifier{Name: name, Loc: tok.Loc})
		if _, ok := p.match(int(',')); ok {
			continue
		}
		break
	}
	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}
	return exprs, nil
}

// parseIndexKeyElement parses one column_ordering_and_options_expr: `expression
// collate_clause? asc_or_desc? null_order? opt_options_list?`. The common case
// is a bare column reference; the expression is kept in Expr (and Name set when
// it is a single identifier) so the query-span walker can descend.
// In Spanner, ordinary CREATE INDEX allows `{ column_name | ( expression ) }`
// (allowSpannerParenExpr=true), while SEARCH/VECTOR indexes require column identifiers.
func (p *Parser) parseIndexKeyElement(allowSpannerParenExpr bool) (*ast.KeyPart, error) {
	start := p.cur.Loc
	var part *ast.KeyPart
	if allowSpannerParenExpr && p.cur.Type == int('(') {
		openTok := p.advance() // '('
		inner, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		closeTok, err := p.expect(int(')'))
		if err != nil {
			return nil, err
		}
		part = &ast.KeyPart{
			Expr: &ast.ParenExpr{Expr: inner, Loc: ast.Loc{Start: openTok.Loc.Start, End: closeTok.Loc.End}},
			Loc:  ast.Loc{Start: start.Start},
		}
	} else {
		// Spanner index keys are column names (or parenthesized expressions for ordinary indexes).
		tok, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}
		name, err := p.identifierText(tok)
		if err != nil {
			return nil, err
		}
		part = &ast.KeyPart{
			Name: name,
			Expr: &ast.Identifier{Name: name, Loc: tok.Loc},
			Loc:  ast.Loc{Start: start.Start},
		}
	}

	// collate_clause?
	if p.cur.Type == kwCOLLATE {
		coll, _, err := p.parseCollateClause()
		if err != nil {
			return nil, err
		}
		part.Collate = coll
	}
	part.Direction = p.matchAscDesc()
	part.NullOrder = p.matchNullOrder()

	// opt_options_list?
	if p.cur.Type == kwOPTIONS {
		opts, err := p.parseDDLOptionsList()
		if err != nil {
			return nil, err
		}
		part.Options = opts
	}
	part.Loc.End = p.prev.Loc.End
	return part, nil
}

// parseCreateIndexSuffix parses the opt_create_index_statement_suffix into stmt:
// a PARTITION BY (+ optional OPTIONS), or an optional OPTIONS followed by a
// Spanner `, INTERLEAVE IN <table>`, or a bare OPTIONS.
func (p *Parser) parseCreateIndexSuffix(stmt *ast.CreateIndexStmt) error {
	// WHERE filter (Spanner partial index). Not in the legacy
	// opt_create_index_statement_suffix rule order shown in antlr_rules, but the
	// Spanner emulator accepts `CREATE INDEX … (cols) WHERE expr` and the truth1
	// corpus documents it (DDL-019). Recognized here, before the option/interleave
	// suffix, so the documented Spanner partial index parses.
	if p.cur.Type == kwWHERE {
		p.advance() // WHERE
		expr, err := p.parseExpr()
		if err != nil {
			return err
		}
		stmt.Where = expr
	}

	switch p.cur.Type {
	case kwPARTITION:
		return p.syntaxErrorAtCur()
	case kwOPTIONS:
		opts, err := p.parseDDLOptionsList()
		if err != nil {
			return err
		}
		stmt.Options = opts
		// opt_options_list? spanner_index_interleave_clause — an interleave may
		// follow the options.
		if p.cur.Type == int(',') && p.peekNext().Type == kwINTERLEAVE {
			if err := p.parseIndexInterleave(stmt); err != nil {
				return err
			}
		}
	case int(','):
		// spanner_index_interleave_clause with no preceding OPTIONS.
		if p.peekNext().Type == kwINTERLEAVE {
			if err := p.parseIndexInterleave(stmt); err != nil {
				return err
			}
		}
	}
	return nil
}

// parseIndexInterleave parses a spanner_index_interleave_clause: `, INTERLEAVE IN
// maybe_dashed_path_expression`. cur is at the ','.
func (p *Parser) parseIndexInterleave(stmt *ast.CreateIndexStmt) error {
	p.advance() // ','
	p.advance() // INTERLEAVE
	if _, err := p.expect(kwIN); err != nil {
		return err
	}
	parent, err := p.parseTablePath()
	if err != nil {
		return err
	}
	stmt.Interleave = parent
	return nil
}
