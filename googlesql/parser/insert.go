package parser

import (
	"github.com/bytebase/omni/googlesql/ast"
)

// This file is part of the `parser-dml` DAG node. It implements GoogleSQL's
// INSERT statement plus the DML helpers shared across INSERT / UPDATE / DELETE /
// MERGE (target-path parsing, ASSERT_ROWS_MODIFIED, THEN RETURN, the SET
// update_item_list, and expression_or_default). The grammar is a hand-port of
// Google's open-source ZetaSQL reference (GoogleSQLParser.g4 §2.7); one omni
// parser serves the BigQuery + Spanner union.
//
// INSERT grammar (insert_statement):
//
//	insert_statement_prefix column_list? insert_values_or_query            opt_assert? opt_returning?
//	insert_statement_prefix column_list? insert_values_list_or_table_clause on_conflict opt_assert? opt_returning?
//	insert_statement_prefix column_list? '(' query ')'                       on_conflict opt_assert? opt_returning?
//	insert_statement_prefix: INSERT opt_or_ignore_replace_update? INTO? <target> hint?
//	opt_or_ignore_replace_update: [OR] IGNORE | [OR] REPLACE | [OR] UPDATE
//	insert_values_or_query: insert_values_list | query
//	insert_values_list_or_table_clause: insert_values_list | TABLE table_clause_no_keyword
//
// The three alternatives collapse to a single unified parse here: after the
// prefix + optional column list we read the data source (VALUES list, a query,
// a `( query )`, or a TABLE clause), then the optional ON CONFLICT, then the
// trailing opt_assert / opt_returning. ON CONFLICT is grammatically restricted
// to the VALUES / TABLE / `( query )` sources (NOT a bare top-level query); we
// enforce that mapping by where ON CONFLICT is allowed.

// parseInsertStmt parses an INSERT statement. INSERT is the current token.
func (p *Parser) parseInsertStmt() (ast.Node, error) {
	start := p.cur.Loc.Start
	p.advance() // INSERT

	stmt := &ast.InsertStmt{Loc: ast.Loc{Start: start}}

	// opt_or_ignore_replace_update: [OR] IGNORE | [OR] REPLACE | [OR] UPDATE.
	// The bare forms (no OR) are accepted by the grammar and the emulator.
	stmt.OrAction = p.parseInsertOrAction()

	// opt_into.
	if _, ok := p.match(kwINTO); ok {
		stmt.Into = true
	}

	// Target: maybe_dashed_generalized_path_expression.
	target, err := p.parseDMLTargetPath()
	if err != nil {
		return nil, err
	}
	stmt.Target = target

	// Per-target `@{…}` hint (insert_statement_prefix's trailing hint).
	if p.cur.Type == int('@') {
		if herr := p.skipHint(); herr != nil {
			return nil, herr
		}
	}

	// Optional explicit column_list. A '(' here opens either a column_list
	// (`( col, … )`) or the parenthesized-query source (`( query )` / a nested
	// `(( … ) UNION …)`). They are distinguished by the token after '(': a query
	// source opens with SELECT / WITH / FROM / GRAPH or a further '(' (a
	// parenthesized query primary), whereas a column list opens with a column
	// name (an identifier). atQueryStart only peeks one token, so it misses the
	// nested `((SELECT …))` case — handle the leading '(' explicitly here.
	if p.cur.Type == int('(') && !p.insertSourceParenFollows() {
		cols, err := p.parseColumnList()
		if err != nil {
			return nil, err
		}
		stmt.Columns = cols
	}

	// The data source: VALUES list | query | ( query ) | TABLE clause.
	// onConflictAllowed tracks whether the source admits a trailing ON CONFLICT
	// (the grammar's alternatives 2 and 3); a bare top-level query (alt 1's
	// `query` that is NOT parenthesized) does NOT.
	onConflictAllowed, err := p.parseInsertSource(stmt)
	if err != nil {
		return nil, err
	}

	// on_conflict_clause (only after the VALUES / TABLE / ( query ) sources).
	if p.cur.Type == kwON {
		if !onConflictAllowed {
			return nil, p.syntaxErrorAtCur()
		}
		oc, err := p.parseOnConflict()
		if err != nil {
			return nil, err
		}
		stmt.OnConflict = oc
	}

	// opt_assert_rows_modified.
	if p.cur.Type == kwASSERT_ROWS_MODIFIED {
		ar, err := p.parseAssertRowsModified()
		if err != nil {
			return nil, err
		}
		stmt.AssertRows = ar
	}

	// opt_returning_clause.
	if p.cur.Type == kwTHEN {
		ret, err := p.parseReturning()
		if err != nil {
			return nil, err
		}
		stmt.Returning = ret
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseInsertOrAction parses opt_or_ignore_replace_update and returns the
// canonical action spelling ("" when absent). All six forms are valid:
//
//	OR IGNORE | IGNORE | OR REPLACE | REPLACE | OR UPDATE | UPDATE
func (p *Parser) parseInsertOrAction() string {
	if p.cur.Type == kwOR {
		// OR (IGNORE | REPLACE | UPDATE)
		switch p.peekNext().Type {
		case kwIGNORE:
			p.advance()
			p.advance()
			return "OR IGNORE"
		case kwREPLACE:
			p.advance()
			p.advance()
			return "OR REPLACE"
		case kwUPDATE:
			p.advance()
			p.advance()
			return "OR UPDATE"
		}
		// `OR` not followed by IGNORE/REPLACE/UPDATE is not an or-action; leave it
		// for the (failing) target parse to report.
		return ""
	}
	// Bare IGNORE | REPLACE | UPDATE (no OR).
	switch p.cur.Type {
	case kwIGNORE:
		p.advance()
		return "IGNORE"
	case kwREPLACE:
		p.advance()
		return "REPLACE"
	case kwUPDATE:
		p.advance()
		return "UPDATE"
	}
	return ""
}

// insertSourceParenFollows reports whether the current '(' opens a parenthesized
// query SOURCE rather than a column_list. The discriminator is the token right
// after '(': a query source begins with the RESERVED query heads SELECT / WITH /
// FROM, or a nested '(' (a parenthesized query primary, e.g.
// `((SELECT 1) UNION ALL (SELECT 2))`); a column_list begins with a column-name
// identifier. Only RESERVED keywords are query-source signals — a NON-reserved
// keyword such as GRAPH is a valid column name (`INSERT INTO t (graph) …`), so it
// must fall through to the column-list path. The current token must be '('.
func (p *Parser) insertSourceParenFollows() bool {
	switch p.peekNext().Type {
	case kwSELECT, kwWITH, kwFROM, int('('):
		return true
	}
	return false
}

// parseInsertSource reads the INSERT data source into stmt and reports whether a
// trailing ON CONFLICT is grammatically allowed after it. The source is one of:
//
//	VALUES (row) [, (row)]…              -> stmt.Rows           (on-conflict OK)
//	TABLE <path/tvf> [WHERE …]           -> stmt.TableClause    (on-conflict OK)
//	( query )                            -> stmt.Query+Parens   (on-conflict OK)
//	query                                -> stmt.Query          (on-conflict NOT allowed)
func (p *Parser) parseInsertSource(stmt *ast.InsertStmt) (onConflictAllowed bool, err error) {
	switch {
	case p.cur.Type == kwVALUES:
		rows, err := p.parseInsertValuesList()
		if err != nil {
			return false, err
		}
		stmt.Rows = rows
		return true, nil

	case p.cur.Type == kwTABLE:
		tc, err := p.parseInsertTableClause()
		if err != nil {
			return false, err
		}
		stmt.TableClause = tc
		return true, nil

	case p.cur.Type == int('('):
		// A '('-led source is grammatically ambiguous between:
		//   alt-3  `( query )`                       — ON CONFLICT eligible
		//   alt-1  query = ( primary ) <continuation> — a bare query, NOT eligible
		// Both accept without ON CONFLICT (oracle); they differ only in whether a
		// trailing ON CONFLICT is legal. The discriminator is whether ANY
		// query-level continuation follows the closing ')': a set operation
		// (UNION/INTERSECT/EXCEPT) OR a trailing ORDER BY / LIMIT / FOR UPDATE. If
		// one does, the parens were a query_primary of a larger BARE query (alt-1,
		// no ON CONFLICT); if nothing query-continuing follows, the parens wrap
		// the whole source (alt-3, ON CONFLICT OK). This is exactly what the
		// emulator enforces:
		//   `(SELECT 1) ON CONFLICT …`            accepts (alt-3)
		//   `(SELECT 1) LIMIT 5`                  accepts (alt-1, bare query)
		//   `(SELECT 1) LIMIT 5 ON CONFLICT …`    REJECTS (alt-1 + ON CONFLICT)
		//   `SELECT 1 ON CONFLICT …`              REJECTS (alt-1, unparenthesized)
		openTok := p.advance() // '('
		inner, err := p.parseQuery()
		if err != nil {
			return false, err
		}
		closeTok, err := p.expect(int(')'))
		if err != nil {
			return false, err
		}
		inner.Parens = true
		inner.Loc = ast.Loc{Start: openTok.Loc.Start, End: closeTok.Loc.End}

		// Continue a set-op chain, if any, treating the parenthesized query as the
		// left primary. parseSetOpChain returns inner unchanged when no set op
		// follows.
		body, err := p.parseSetOpChain(inner)
		if err != nil {
			return false, err
		}
		setOpFollowed := body != ast.Node(inner)

		// A query-level trailing ORDER BY / LIMIT / FOR UPDATE also marks the
		// source as a bare query (alt-1), not a `( query )`. wrapSetOpQueryTail
		// consumes them and reports whether any were present.
		q, tailFollowed, err := p.wrapSetOpQueryTail(body)
		if err != nil {
			return false, err
		}

		if !setOpFollowed && !tailFollowed {
			// Nothing query-continuing followed the ')': a true `( query )` source
			// — ON CONFLICT eligible.
			stmt.Query = inner
			stmt.QueryParens = true
			return true, nil
		}
		// A set op and/or a query-level tail followed: a bare query (alt-1). NOT
		// ON CONFLICT eligible.
		stmt.Query = q
		return false, nil

	default:
		// A bare top-level query (SELECT / WITH). NOT on-conflict-eligible.
		q, err := p.parseQuery()
		if err != nil {
			return false, err
		}
		stmt.Query = q
		return false, nil
	}
}

// wrapSetOpQueryTail wraps an already-parsed query body in a *QueryStmt and
// consumes the query-level trailing ORDER BY / LIMIT[/OFFSET] / FOR UPDATE that
// bind to the whole query (query_without_pipe_operators). It reports whether any
// such tail was present (tailFollowed). Used for the INSERT alt-1
// `( primary ) [UNION …] [ORDER BY …] [LIMIT …]` source, where the parenthesized
// primary (and any set-op chain) was already pulled off before this call.
//
// body is always wrapped in a FRESH outer *QueryStmt — a parenthesized inner
// query (`( SELECT 1 )`) becomes the Body of the outer query so an outer-level
// ORDER BY / LIMIT stays distinct from any the inner query carried inside its
// own parens (`( SELECT 1 ORDER BY x ) ORDER BY y` keeps both). The caller only
// uses the returned QueryStmt on the bare-query (alt-1) path; on the alt-3
// `( query )` path it returns the inner QueryStmt directly and ignores this.
func (p *Parser) wrapSetOpQueryTail(body ast.Node) (q *ast.QueryStmt, tailFollowed bool, err error) {
	q = &ast.QueryStmt{Body: body, Loc: ast.Loc{Start: qLoc(body).Start, End: qLoc(body).End}}
	if p.cur.Type == kwORDER {
		items, err := p.parseOrderByClause()
		if err != nil {
			return nil, false, err
		}
		q.OrderBy = items
		q.Loc.End = p.prev.Loc.End
		tailFollowed = true
	}
	if p.cur.Type == kwLIMIT {
		if err := p.parseQueryLimitOffset(q); err != nil {
			return nil, false, err
		}
		tailFollowed = true
	}
	if p.cur.Type == kwFOR && p.peekNext().Type == kwUPDATE {
		p.advance()           // FOR
		updTok := p.advance() // UPDATE
		q.ForUpdate = true
		q.Loc.End = updTok.Loc.End
		tailFollowed = true
	}
	return q, tailFollowed, nil
}

// parseInsertValuesList parses insert_values_list:
// `VALUES insert_values_row (, insert_values_row)*`. VALUES is the current token.
func (p *Parser) parseInsertValuesList() ([]*ast.InsertRow, error) {
	p.advance() // VALUES
	var rows []*ast.InsertRow
	for {
		row, err := p.parseInsertValuesRow()
		if err != nil {
			return nil, err
		}
		rows = append(rows, row)
		if _, ok := p.match(int(',')); !ok {
			break
		}
	}
	return rows, nil
}

// parseInsertValuesRow parses insert_values_row: a parenthesized
// expression_or_default list `( e_or_d (, e_or_d)* )`. The current token is '('.
func (p *Parser) parseInsertValuesRow() (*ast.InsertRow, error) {
	openTok, err := p.expect(int('('))
	if err != nil {
		return nil, err
	}
	row := &ast.InsertRow{Loc: openTok.Loc}
	for {
		val, err := p.parseExprOrDefault()
		if err != nil {
			return nil, err
		}
		row.Values = append(row.Values, val)
		if _, ok := p.match(int(',')); !ok {
			break
		}
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	row.Loc.End = closeTok.Loc.End
	return row, nil
}

// parseInsertTableClause parses table_clause_unreversed → table_clause_no_keyword:
// `TABLE (path | tvf) [WHERE expr]`. TABLE is the current token.
func (p *Parser) parseInsertTableClause() (*ast.InsertTable, error) {
	tableTok := p.advance() // TABLE
	tc := &ast.InsertTable{Loc: tableTok.Loc}

	path, err := p.parsePathExpr()
	if err != nil {
		return nil, err
	}
	if p.cur.Type == int('(') {
		// path ( args ) — a tvf_with_suffixes table source.
		fc, err := p.parseFuncCallSuffix(path)
		if err != nil {
			return nil, err
		}
		tc.Func = fc.(*ast.FuncCall)
	} else {
		tc.Path = path
	}
	tc.Loc.End = p.prev.Loc.End

	if p.cur.Type == kwWHERE {
		p.advance() // WHERE
		where, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		tc.Where = where
		tc.Loc.End = nodeLoc(where).End
	}
	return tc, nil
}

// parseOnConflict parses on_conflict_clause:
//
//	ON CONFLICT [conflict_target] DO NOTHING
//	ON CONFLICT [conflict_target] DO UPDATE SET update_item_list [WHERE expr]
//
// conflict_target: column_list | ON UNIQUE CONSTRAINT identifier
//
// ON is the current token.
func (p *Parser) parseOnConflict() (*ast.OnConflict, error) {
	onTok := p.advance() // ON
	if _, err := p.expect(kwCONFLICT); err != nil {
		return nil, err
	}
	oc := &ast.OnConflict{Loc: onTok.Loc}

	// opt_conflict_target: column_list | ON UNIQUE CONSTRAINT identifier.
	switch {
	case p.cur.Type == int('('):
		cols, err := p.parseColumnList()
		if err != nil {
			return nil, err
		}
		oc.Columns = cols
	case p.cur.Type == kwON:
		p.advance() // ON
		if _, err := p.expect(kwUNIQUE); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwCONSTRAINT); err != nil {
			return nil, err
		}
		nameTok, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}
		name, err := p.identifierText(nameTok)
		if err != nil {
			return nil, err
		}
		oc.ConstraintName = name
	}

	if _, err := p.expect(kwDO); err != nil {
		return nil, err
	}
	switch p.cur.Type {
	case kwNOTHING:
		p.advance()
		oc.DoNothing = true
	case kwUPDATE:
		p.advance() // UPDATE
		if _, err := p.expect(kwSET); err != nil {
			return nil, err
		}
		items, err := p.parseUpdateItemList()
		if err != nil {
			return nil, err
		}
		oc.SetItems = items
		if p.cur.Type == kwWHERE {
			p.advance() // WHERE
			where, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			oc.Where = where
		}
	default:
		return nil, p.syntaxErrorAtCur()
	}
	oc.Loc.End = p.prev.Loc.End
	return oc, nil
}

// ---------------------------------------------------------------------------
// Shared DML helpers
// ---------------------------------------------------------------------------

// parseDMLTargetPath parses maybe_dashed_generalized_path_expression — the
// DML target path used by INSERT / UPDATE / DELETE. It is a dotted, optionally
// dashed (BigQuery `my-project.ds.tbl`) path, optionally extended by generalized
// accessors: `.field`, `.(extension)`, `[index]` (for nested / array DML
// targets, e.g. `UPDATE t.a[0].b SET …`). The base dotted/dashed path is held in
// the returned PathExpr.Parts; trailing generalized accessors are consumed (they
// fold into the path's source span) so the target round-trips and its base table
// name is recoverable for query-span resolution.
func (p *Parser) parseDMLTargetPath() (*ast.PathExpr, error) {
	// Base dotted/dashed path (reuses the FROM-source path parser, which folds
	// dashed `a-b-c` and dotted `a.b.c` components — the dashed BigQuery form is
	// a target-valid path per the .g4's dashed_path_expression alternative).
	path, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	// Trailing generalized accessors (generalized_path_expression's recursive
	// `.field` / `.(ext)` / `[index]`). These extend a nested/array target; the
	// base table name stays in path.Parts[0..]. parseTablePath already folded the
	// leading dotted/dashed run, so this loop only runs once a generalized
	// accessor (`[…]` or `.(…)`) has interrupted that run — after which a further
	// `.field` must be consumed here too (the run does not resume in
	// parseTablePath). We fold each accessor into path.Parts so the full target
	// span and base name round-trip for query-span resolution.
	for {
		switch {
		case p.cur.Type == int('.') && p.peekNext().Type == int('('):
			// `.(extension)` — generalized_extension_path.
			p.advance() // '.'
			p.advance() // '('
			ext, err := p.parsePathExpr()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(int(')')); err != nil {
				return nil, err
			}
			path.Parts = append(path.Parts, "("+ext.String()+")")
			path.Loc.End = p.prev.Loc.End
		case p.cur.Type == int('.') && isAnyKeywordIdentifier(p.peekNext().Type):
			// `.field` — a plain name part following a generalized accessor.
			p.advance() // '.'
			partTok := p.advance()
			part, err := p.identifierText(partTok)
			if err != nil {
				return nil, err
			}
			path.Parts = append(path.Parts, part)
			path.Loc.End = partTok.Loc.End
		case p.cur.Type == int('['):
			// `[index]` — array-element access on the target.
			p.advance() // '['
			if _, err := p.parseExpr(); err != nil {
				return nil, err
			}
			closeTok, err := p.expect(int(']'))
			if err != nil {
				return nil, err
			}
			path.Parts = append(path.Parts, "[...]")
			path.Loc.End = closeTok.Loc.End
		default:
			return path, nil
		}
	}
}

// parseExprOrDefault parses expression_or_default: an expression or the bare
// DEFAULT keyword. DEFAULT yields a *DefaultExpr marker.
func (p *Parser) parseExprOrDefault() (ast.Node, error) {
	if p.cur.Type == kwDEFAULT {
		tok := p.advance()
		return &ast.DefaultExpr{Loc: tok.Loc}, nil
	}
	return p.parseExpr()
}

// parseAssertRowsModified parses opt_assert_rows_modified:
// `ASSERT_ROWS_MODIFIED possibly_cast_int_literal_or_parameter`. The operand is
// NARROW — an integer literal, a query parameter (`@p` / `?`), a system variable
// (`@@v`), or a `CAST(... AS type)` of one — NOT an arbitrary expression. Using
// the full expression grammar here would wrongly accept `ASSERT_ROWS_MODIFIED
// 1 + 1`, a string, or a subquery (oracle: all three syntax-reject).
// ASSERT_ROWS_MODIFIED is the current token.
func (p *Parser) parseAssertRowsModified() (ast.Node, error) {
	p.advance() // ASSERT_ROWS_MODIFIED
	return p.parseIntLiteralOrParameterOrCast()
}

// parseIntLiteralOrParameterOrCast parses possibly_cast_int_literal_or_parameter:
//
//	integer_literal | parameter_expression | system_variable_expression
//	| CAST ( int_literal_or_parameter AS type [FORMAT …] )
//
// It accepts the restricted operand set the grammar allows for
// ASSERT_ROWS_MODIFIED / TABLESAMPLE counts, rejecting general expressions. The
// CAST arm reuses parseCastExpr (whose inner is a full expression — a CAST of a
// non-literal like `CAST(1+1 AS INT64)` is a rare over-accept relative to the
// grammar's int_literal_or_parameter inner, but every documented/corpus form
// casts a literal or parameter; the load-bearing fix is rejecting a bare
// non-literal operand, which this does).
func (p *Parser) parseIntLiteralOrParameterOrCast() (ast.Node, error) {
	switch p.cur.Type {
	case tokInteger:
		tok := p.advance()
		return &ast.Literal{Kind: ast.LitInt, Value: tok.Str, Loc: tok.Loc}, nil
	case int('@'):
		return p.parseNamedParameter()
	case tokAtAt:
		return p.parseSystemVariable()
	case int('?'):
		tok := p.advance()
		return &ast.Parameter{Positional: true, Loc: tok.Loc}, nil
	case kwCAST:
		return p.parseCastExpr(false)
	default:
		return nil, p.syntaxErrorAtCur()
	}
}

// parseReturning parses opt_returning_clause:
//
//	THEN RETURN [WITH ACTION [AS action_alias]] select_list
//
// THEN is the current token.
func (p *Parser) parseReturning() (*ast.Returning, error) {
	thenTok := p.advance() // THEN
	if _, err := p.expect(kwRETURN); err != nil {
		return nil, err
	}
	ret := &ast.Returning{Loc: thenTok.Loc}

	// `WITH ACTION [AS id]` is the action modifier. It must be distinguished from
	// a select_list item that is an inline WITH(...) expression (`THEN RETURN
	// WITH(x AS 1, x)`): only `WITH ACTION` (WITH followed by the ACTION keyword)
	// is the modifier; `WITH (` is a select-list expression, left for
	// parseSelectList (oracle: `THEN RETURN WITH(x AS 1, x)` accepts).
	if p.cur.Type == kwWITH && p.peekNext().Type == kwACTION {
		p.advance() // WITH
		p.advance() // ACTION
		ret.WithAction = true
		if p.cur.Type == kwAS {
			p.advance() // AS
			aliasTok, err := p.expectIdentifier()
			if err != nil {
				return nil, err
			}
			alias, err := p.identifierText(aliasTok)
			if err != nil {
				return nil, err
			}
			ret.ActionAlias = alias
		}
	}

	items, err := p.parseSelectList()
	if err != nil {
		return nil, err
	}
	ret.Items = items
	ret.Loc.End = p.prev.Loc.End
	return ret, nil
}

// parseUpdateItemList parses update_item_list: `update_item (, update_item)*`.
// The current token begins the first item.
func (p *Parser) parseUpdateItemList() ([]*ast.UpdateItem, error) {
	var items []*ast.UpdateItem
	for {
		item, err := p.parseUpdateItem()
		if err != nil {
			return nil, err
		}
		items = append(items, item)
		if _, ok := p.match(int(',')); !ok {
			break
		}
	}
	return items, nil
}

// parseUpdateItem parses update_item: a `generalized_path = expression_or_default`
// assignment OR a nested_dml_statement `( insert | update | delete )`.
func (p *Parser) parseUpdateItem() (*ast.UpdateItem, error) {
	start := p.cur.Loc
	// nested_dml_statement: ( dml_statement ).
	if p.cur.Type == int('(') {
		p.advance() // '('
		nested, err := p.parseNestedDML()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(int(')')); err != nil {
			return nil, err
		}
		return &ast.UpdateItem{Nested: nested, Loc: ast.Loc{Start: start.Start, End: p.prev.Loc.End}}, nil
	}

	// update_set_value: generalized_path_expression = expression_or_default.
	lhs, err := p.parseGeneralizedPath()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(int('=')); err != nil {
		return nil, err
	}
	rhs, err := p.parseExprOrDefault()
	if err != nil {
		return nil, err
	}
	return &ast.UpdateItem{Path: lhs, Value: rhs, Loc: ast.Loc{Start: start.Start, End: p.prev.Loc.End}}, nil
}

// parseNestedDML parses the body of a nested_dml_statement (the inner
// insert/update/delete of an `( dml )` update item). The current token is the
// inner statement's leading keyword.
func (p *Parser) parseNestedDML() (ast.Node, error) {
	switch p.cur.Type {
	case kwINSERT:
		return p.parseInsertStmt()
	case kwUPDATE:
		return p.parseUpdateStmt()
	case kwDELETE:
		return p.parseDeleteStmt()
	default:
		return nil, p.syntaxErrorAtCur()
	}
}
