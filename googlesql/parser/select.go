package parser

import (
	"github.com/bytebase/omni/googlesql/ast"
)

// qLoc returns the Loc of any node this node (parser-select) builds, plus any
// expression node (via the expressions node's nodeLoc). The expressions node's
// nodeLoc enumerates only expression types and returns NoLoc for the query /
// FROM node types defined here, so qLoc adds those cases and otherwise delegates.
// Used wherever a span is computed from a query_primary / FROM-source / clause
// node whose concrete type may be a query node rather than an expression.
func qLoc(n ast.Node) ast.Loc {
	switch v := n.(type) {
	case *ast.QueryStmt:
		return v.Loc
	case *ast.SelectStmt:
		return v.Loc
	case *ast.SetOperation:
		return v.Loc
	case *ast.TableExpr:
		return v.Loc
	case *ast.UnnestExpr:
		return v.Loc
	case *ast.JoinExpr:
		return v.Loc
	default:
		return nodeLoc(n)
	}
}

// This file is part of the `parser-select` DAG node. It implements GoogleSQL's
// query stack — the `query` / `query_without_pipe_operators` / `select` rules
// (GoogleSQLParser.g4 §2.13, §2.16), a hand-port of Google's open-source ZetaSQL
// reference and the grammar bytebase consumes today. The companion files are
// from_join.go (FROM / table sources / joins), set_ops.go (UNION/INTERSECT/
// EXCEPT), and cte.go (WITH / CTE).
//
// One omni parser serves both BigQuery and Spanner; this is the UNION grammar.
// Adjudication: the LIVE Cloud Spanner emulator oracle (oracle.md) decides
// accept/reject for shared + Spanner-only query forms. BigQuery-only forms
// (QUALIFY, WITH RECURSIVE, etc.) parse in the union grammar even though the
// Spanner emulator *feature-rejects* them — that feature-reject is NOT a grammar
// verdict (divergence ledger #9 QUALIFY, #11 WITH RECURSIVE/MERGE). The parser
// therefore accepts them; the differential test (select_oracle_test.go) only
// asserts the Spanner-authoritative forms.
//
// Grammar-faithful tree shape (a key ZetaSQL nuance): the trailing ORDER BY /
// LIMIT-OFFSET / FOR UPDATE attach to the whole query
// (query_without_pipe_operators), not to the inner SELECT. parseQuery builds an
// ast.QueryStmt that wraps the set-op/select Body and owns those trailing
// modifiers; parseSelectStmt builds the inner ast.SelectStmt with only
// SELECT…FROM…WHERE…GROUP BY…HAVING…QUALIFY…WINDOW.
//
// This node does NOT own PIVOT / UNPIVOT / TABLESAMPLE / MATCH_RECOGNIZE /
// AT SYSTEM TIME (table-source operators) — those belong to the
// parser-query-clauses node. UNNEST as a FROM source IS owned here.

// parseQuery parses a complete query (query / query_without_pipe_operators):
//
//	[WITH …] <query_primary_or_set_operation> [ORDER BY …] [LIMIT …] [FOR UPDATE]
//
// It is the entry point for the SELECT / WITH / `(` / FROM dispatch cases and
// for every subquery body (table_subquery, parenthesized_query, IN/EXISTS/ARRAY
// subqueries). The returned node is always an *ast.QueryStmt.
func (p *Parser) parseQuery() (*ast.QueryStmt, error) {
	start := p.cur.Loc.Start

	q := &ast.QueryStmt{Loc: ast.Loc{Start: start}}

	// Optional leading WITH clause (cte.go).
	if p.cur.Type == kwWITH {
		with, err := p.parseWithClause()
		if err != nil {
			return nil, err
		}
		q.With = with
	}

	// The query body: a query_primary or a left-recursive set-operation chain.
	body, err := p.parseQueryPrimaryOrSetOp()
	if err != nil {
		return nil, err
	}
	q.Body = body
	q.Loc.End = qLoc(body).End

	// Trailing query-level ORDER BY / LIMIT-OFFSET / FOR UPDATE
	// (query_without_pipe_operators). These bind to the whole query, not the
	// inner SELECT.
	if p.cur.Type == kwORDER {
		items, err := p.parseOrderByClause()
		if err != nil {
			return nil, err
		}
		q.OrderBy = items
		q.Loc.End = p.prev.Loc.End
	}
	if p.cur.Type == kwLIMIT {
		if err := p.parseQueryLimitOffset(q); err != nil {
			return nil, err
		}
	}
	// FOR UPDATE (Spanner row-locking; parses in the union grammar).
	if p.cur.Type == kwFOR && p.peekNext().Type == kwUPDATE {
		p.advance()           // FOR
		updTok := p.advance() // UPDATE
		q.ForUpdate = true
		q.Loc.End = updTok.Loc.End
	}

	return q, nil
}

// parseQueryLimitOffset parses the query-level `LIMIT count [OFFSET skip]`
// (limit_offset_clause) into q. LIMIT is the current token.
func (p *Parser) parseQueryLimitOffset(q *ast.QueryStmt) error {
	p.advance() // LIMIT
	lim, err := p.parseExpr()
	if err != nil {
		return err
	}
	q.Limit = lim
	q.Loc.End = nodeLoc(lim).End
	if p.cur.Type == kwOFFSET {
		p.advance()
		off, err := p.parseExpr()
		if err != nil {
			return err
		}
		q.Offset = off
		q.Loc.End = nodeLoc(off).End
	}
	return nil
}

// parseQueryPrimaryOrSetOp parses a query_primary_or_set_operation: a query
// primary, optionally followed by a left-associative chain of set operations
// (set_ops.go). The result is a *ast.SelectStmt, a *ast.QueryStmt (parenthesized
// primary), or an *ast.SetOperation.
func (p *Parser) parseQueryPrimaryOrSetOp() (ast.Node, error) {
	left, err := p.parseQueryPrimary()
	if err != nil {
		return nil, err
	}
	return p.parseSetOpChain(left)
}

// parseQueryPrimary parses a query_primary:
//
//	select | ( query ) [AS alias]
//
// A leading SELECT builds an ast.SelectStmt; a leading '(' parses a nested query
// and returns it as a parenthesized ast.QueryStmt (Parens=true). The optional
// `AS alias` after a parenthesized primary is accepted (grammar:
// `query_primary: ( parenthesized_query ) [AS alias]`) — the alias is consumed
// but, lacking an alias slot on QueryStmt, is not retained (it is only meaningful
// for set-op column naming, which bytebase's query-span resolves structurally).
func (p *Parser) parseQueryPrimary() (ast.Node, error) {
	switch p.cur.Type {
	case kwSELECT:
		return p.parseSelectStmt()
	case int('('):
		openTok := p.advance() // '('
		inner, err := p.parseQuery()
		if err != nil {
			return nil, err
		}
		closeTok, err := p.expect(int(')'))
		if err != nil {
			return nil, err
		}
		inner.Parens = true
		inner.Loc = ast.Loc{Start: openTok.Loc.Start, End: closeTok.Loc.End}
		// Optional `AS alias` (consumed, not retained — see doc comment).
		if p.cur.Type == kwAS {
			p.advance()
			if _, err := p.expectIdentifier(); err != nil {
				return nil, err
			}
		}
		return inner, nil
	case kwFROM:
		// FROM-first query (with_clause? from_clause): the legacy grammar
		// recognizes a leading FROM as a query head, then ZetaSQL rejects it
		// ("Syntax error: Unexpected FROM"). The Spanner emulator likewise
		// rejects `FROM t`. We surface the same reject by reporting at FROM
		// rather than building a tree.
		return nil, &ParseError{Loc: p.cur.Loc, Msg: "syntax error: unexpected FROM (a query must begin with SELECT)"}
	default:
		return nil, p.syntaxErrorAtCur()
	}
}

// ---------------------------------------------------------------------------
// SELECT statement
// ---------------------------------------------------------------------------

// parseSelectStmt parses a single SELECT block (select):
//
//	SELECT [hint] [WITH <name> [OPTIONS(…)]] [ALL|DISTINCT] [AS STRUCT|VALUE|path]
//	  <select_list>
//	  [FROM …] [WHERE …] [GROUP BY …] [HAVING …] [QUALIFY …] [WINDOW …]
//
// SELECT is the current token. ORDER BY / LIMIT / FOR UPDATE are NOT consumed
// here — they belong to the enclosing query (parseQuery).
func (p *Parser) parseSelectStmt() (*ast.SelectStmt, error) {
	selectTok, err := p.expect(kwSELECT)
	if err != nil {
		return nil, err
	}
	stmt := &ast.SelectStmt{Loc: ast.Loc{Start: selectTok.Loc.Start}}

	// Optional SELECT-level hint @{...} / @N.
	if p.cur.Type == int('@') {
		if herr := p.skipHint(); herr != nil {
			return nil, herr
		}
	}

	// Optional `WITH <name> [OPTIONS(…)]` differential-privacy / aggregation-
	// threshold clause (opt_select_with). This is BEFORE ALL/DISTINCT.
	// Disambiguate from a column named via WITH(...) inline expression: the
	// select-with clause is `WITH <identifier>` (NOT `WITH (`).
	if p.cur.Type == kwWITH && p.peekNext().Type != int('(') {
		p.advance() // WITH
		nameTok, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}
		name, err := p.identifierText(nameTok)
		if err != nil {
			return nil, err
		}
		stmt.SelectWith = name
		if p.cur.Type == kwOPTIONS {
			p.advance()
			if err := p.skipOptionsList(); err != nil {
				return nil, err
			}
		}
	}

	// ALL | DISTINCT set quantifier.
	switch p.cur.Type {
	case kwALL:
		p.advance()
		stmt.All = true
	case kwDISTINCT:
		p.advance()
		stmt.Distinct = true
	}

	// AS STRUCT | AS VALUE | AS <path> (opt_select_as_clause).
	if p.cur.Type == kwAS {
		p.advance()
		switch p.cur.Type {
		case kwSTRUCT:
			p.advance()
			stmt.As = ast.SelectAsStruct
		case kwVALUE:
			p.advance()
			stmt.As = ast.SelectAsValue
		default:
			// AS <path> (proto/type name).
			path, err := p.parsePathExpr()
			if err != nil {
				return nil, err
			}
			stmt.As = ast.SelectAsTypeName
			stmt.AsTypeName = path
		}
	}

	// select_list (>= 1 item, trailing comma allowed). The grammar's empty-list-
	// before-FROM error alternative is surfaced by parseSelectList.
	items, err := p.parseSelectList()
	if err != nil {
		return nil, err
	}
	stmt.Items = items

	// FROM clause (from_join.go).
	if p.cur.Type == kwFROM {
		from, err := p.parseFromClause()
		if err != nil {
			return nil, err
		}
		stmt.From = from
	}

	// WHERE → GROUP BY → HAVING → QUALIFY → WINDOW (opt_clauses_following_from).
	// The grammar fixes this order; each clause is optional.
	if p.cur.Type == kwWHERE {
		p.advance()
		where, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Where = where
	}
	if p.cur.Type == kwGROUP {
		gb, err := p.parseGroupByClause()
		if err != nil {
			return nil, err
		}
		stmt.GroupBy = gb
	}
	if p.cur.Type == kwHAVING {
		p.advance()
		having, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Having = having
	}
	if p.cur.Type == kwQUALIFY {
		// QUALIFY: valid BigQuery, parses in the union grammar; Spanner
		// feature-rejects (divergence #9). We accept it here.
		p.advance()
		qualify, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Qualify = qualify
	}
	if p.cur.Type == kwWINDOW {
		wins, err := p.parseWindowClause()
		if err != nil {
			return nil, err
		}
		stmt.Window = wins
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// SELECT list
// ---------------------------------------------------------------------------

// parseSelectList parses a select_list: one or more select_list_items separated
// by commas, with a trailing comma allowed. At least one item is required; an
// empty list (`SELECT FROM …` / `SELECT` with nothing) is a syntax error
// (grammar: select_clause's empty-list-before-FROM error alt; oracle: Spanner
// rejects `SELECT FROM t` with "SELECT list must not be empty").
func (p *Parser) parseSelectList() ([]*ast.SelectItem, error) {
	if p.atSelectListEnd() {
		return nil, &ParseError{Loc: p.cur.Loc, Msg: "syntax error: SELECT list must not be empty"}
	}
	var items []*ast.SelectItem
	for {
		item, err := p.parseSelectItem()
		if err != nil {
			return nil, err
		}
		items = append(items, item)
		if p.cur.Type != int(',') {
			break
		}
		p.advance() // consume ','
		// Trailing comma: `SELECT a, FROM t` — the comma may be followed by a
		// clause keyword / EOF, which ends the list.
		if p.atSelectListEnd() {
			break
		}
	}
	return items, nil
}

// atSelectListEnd reports whether the current token ends the select list (a
// following clause keyword or end-of-statement), used to detect an empty list
// and to allow a trailing comma.
func (p *Parser) atSelectListEnd() bool {
	switch p.cur.Type {
	case kwFROM, kwWHERE, kwGROUP, kwHAVING, kwQUALIFY, kwWINDOW,
		kwORDER, kwLIMIT, kwUNION, kwINTERSECT, kwEXCEPT,
		int(')'), int(';'), tokEOF:
		return true
	}
	return false
}

// parseSelectItem parses one select_list_item:
//
//	expr [[AS] alias]      (select_column_expr)
//	*  [star_modifiers]    (select_column_star)
//	expr.* [star_modifiers] (select_column_dot_star)
func (p *Parser) parseSelectItem() (*ast.SelectItem, error) {
	// Bare star: `* [EXCEPT(...) [REPLACE(...)]]`.
	if p.cur.Type == int('*') {
		starTok := p.advance()
		item := &ast.SelectItem{Star: true, Loc: starTok.Loc}
		mods, end, err := p.tryParseStarModifiers(starTok.Loc.Start)
		if err != nil {
			return nil, err
		}
		if mods != nil {
			item.Modifiers = mods
			item.Loc.End = end
		}
		return item, nil
	}

	// Dot-star (select_column_dot_star): `<path>.* [star_modifiers]`. The
	// qualifier is a path_expression (e.g. `g.*`, `l.location.*`). The expression
	// parser would error on the trailing `.*` (it tries to read `*` as a field
	// name), so we detect a path-qualified dot-star up front via a throwaway scan
	// and parse the qualifier path ourselves (parsePathExpr stops before `.*`).
	if p.pathDotStarAhead() {
		qual, err := p.parsePathExpr()
		if err != nil {
			return nil, err
		}
		p.advance()            // '.'
		starTok := p.advance() // '*'
		item := &ast.SelectItem{Expr: qual, Star: true, Loc: ast.Loc{Start: qual.Loc.Start, End: starTok.Loc.End}}
		mods, end, err := p.tryParseStarModifiers(qual.Loc.Start)
		if err != nil {
			return nil, err
		}
		if mods != nil {
			item.Modifiers = mods
			item.Loc.End = end
		}
		return item, nil
	}

	// Expression. After it: `[AS] alias` or nothing. (A non-path `expr.*` —
	// e.g. `f(x).*` — is not a documented select-list form; GoogleSQL's
	// select_column_dot_star qualifier is a path. If a `.*` trails a non-path
	// expression the expression parser surfaces the syntax error, matching the
	// oracle.)
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	item := &ast.SelectItem{Expr: expr, Loc: nodeLoc(expr)}
	// Optional alias: `AS name` or a bare `name` (implicit alias). A bare alias
	// is only an identifier-start that is NOT a clause keyword (so `SELECT a FROM`
	// does not read FROM as an alias). Reserved keywords are not valid bare
	// aliases.
	if p.cur.Type == kwAS {
		p.advance()
		aliasTok, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}
		alias, err := p.identifierText(aliasTok)
		if err != nil {
			return nil, err
		}
		item.Alias = alias
		item.Loc.End = aliasTok.Loc.End
	} else if isIdentifierStart(p.cur.Type) && !p.atSelectListEnd() {
		// Implicit alias: `SELECT a b` → column a aliased b. isIdentifierStart
		// already excludes reserved keywords (FROM/WHERE/etc are reserved), so a
		// following clause keyword is never mistaken for an alias.
		aliasTok := p.advance()
		alias, err := p.identifierText(aliasTok)
		if err != nil {
			return nil, err
		}
		item.Alias = alias
		item.Loc.End = aliasTok.Loc.End
	}
	return item, nil
}

// pathDotStarAhead reports whether the current position begins a path-qualified
// dot-star select item: `identifier (. identifier)* . *`. It scans a throwaway
// lexer from the current token WITHOUT mutating parser state (the parser's
// two-token lookahead cannot see past an arbitrary-length path). A leading token
// that is not an identifier-start, or a path not terminated by exactly `. *`,
// returns false (the item is an ordinary expression).
func (p *Parser) pathDotStarAhead() bool {
	if !isIdentifierStart(p.cur.Type) {
		return false
	}
	startIdx := absIndex(p, p.cur.Loc.Start)
	lx := NewLexerWithOffset(p.input[startIdx:], 0)
	// First path component.
	if !isIdentifierStart(lx.NextToken().Type) {
		return false
	}
	for {
		t := lx.NextToken()
		if t.Type != int('.') {
			return false
		}
		nx := lx.NextToken()
		if nx.Type == int('*') {
			return true // `… . *`
		}
		if !isAnyKeywordIdentifier(nx.Type) {
			return false // `. <non-name>` — not a path-dot-star
		}
		// `. identifier` — continue scanning the path.
	}
}

// tryParseStarModifiers parses an optional `EXCEPT ( col, … ) [REPLACE ( item, … )]`
// star_modifiers suffix following a `*` or `expr.*`. It returns (nil, 0, nil)
// when neither EXCEPT nor REPLACE follows. start anchors the returned Loc.
func (p *Parser) tryParseStarModifiers(start int) (*ast.StarModifiers, int, error) {
	if p.cur.Type != kwEXCEPT && p.cur.Type != kwREPLACE {
		return nil, 0, nil
	}
	mods := &ast.StarModifiers{Loc: ast.Loc{Start: start}}
	end := start

	// EXCEPT ( col, … )  (star_except_list).
	if p.cur.Type == kwEXCEPT {
		p.advance() // EXCEPT
		if _, err := p.expect(int('(')); err != nil {
			return nil, 0, err
		}
		for {
			colTok, err := p.expectIdentifier()
			if err != nil {
				return nil, 0, err
			}
			col, err := p.identifierText(colTok)
			if err != nil {
				return nil, 0, err
			}
			mods.Except = append(mods.Except, col)
			if p.cur.Type != int(',') {
				break
			}
			p.advance()
		}
		closeTok, err := p.expect(int(')'))
		if err != nil {
			return nil, 0, err
		}
		end = closeTok.Loc.End
	}

	// REPLACE ( expr AS col, … )  (star_replace_list).
	if p.cur.Type == kwREPLACE {
		p.advance() // REPLACE
		if _, err := p.expect(int('(')); err != nil {
			return nil, 0, err
		}
		for {
			valExpr, err := p.parseExpr()
			if err != nil {
				return nil, 0, err
			}
			if _, err := p.expect(kwAS); err != nil {
				return nil, 0, err
			}
			aliasTok, err := p.expectIdentifier()
			if err != nil {
				return nil, 0, err
			}
			alias, err := p.identifierText(aliasTok)
			if err != nil {
				return nil, 0, err
			}
			mods.Replace = append(mods.Replace, &ast.StructFieldExpr{
				Value: valExpr, Alias: alias,
				Loc: ast.Loc{Start: nodeLoc(valExpr).Start, End: aliasTok.Loc.End},
			})
			if p.cur.Type != int(',') {
				break
			}
			p.advance()
		}
		closeTok, err := p.expect(int(')'))
		if err != nil {
			return nil, 0, err
		}
		end = closeTok.Loc.End
	}

	mods.Loc.End = end
	return mods, end, nil
}

// ---------------------------------------------------------------------------
// GROUP BY
// ---------------------------------------------------------------------------

// parseGroupByClause parses a GROUP BY clause (group_by_clause):
//
//	GROUP [hint] [AND ORDER] BY { ALL | <grouping_item>, … }
//
// GROUP is the current token. The grouping items are plain expressions, the
// empty grouping `()`, or ROLLUP / CUBE / GROUPING SETS.
func (p *Parser) parseGroupByClause() (*ast.GroupByClause, error) {
	groupTok := p.advance() // GROUP
	gb := &ast.GroupByClause{Loc: groupTok.Loc}
	if p.cur.Type == int('@') {
		if herr := p.skipHint(); herr != nil {
			return nil, herr
		}
	}
	if p.cur.Type == kwAND {
		p.advance()
		if _, err := p.expect(kwORDER); err != nil {
			return nil, err
		}
		gb.AndOrder = true
	}
	if _, err := p.expect(kwBY); err != nil {
		return nil, err
	}

	// GROUP BY ALL.
	if p.cur.Type == kwALL {
		allTok := p.advance()
		gb.Kind = ast.GroupByAll
		gb.Loc.End = allTok.Loc.End
		return gb, nil
	}

	gb.Kind = ast.GroupByItems
	for {
		item, err := p.parseGroupingItem()
		if err != nil {
			return nil, err
		}
		gb.Items = append(gb.Items, item)
		if p.cur.Type != int(',') {
			break
		}
		p.advance()
	}
	gb.Loc.End = p.prev.Loc.End
	return gb, nil
}

// parseGroupingItem parses one grouping_item:
//
//	() | expr [AS alias] [ASC|DESC] | ROLLUP(…) | CUBE(…) | GROUPING SETS(…)
func (p *Parser) parseGroupingItem() (*ast.GroupingItem, error) {
	switch p.cur.Type {
	case int('('):
		// Empty grouping `()` OR a parenthesized expression used as a grouping
		// item. GoogleSQL's grouping_item allows a bare `()`; a `(expr …)` is a
		// normal parenthesized expression handled by the default case.
		if p.peekNext().Type == int(')') {
			openTok := p.advance()  // '('
			closeTok := p.advance() // ')'
			return &ast.GroupingItem{Kind: ast.GroupingEmpty, Loc: ast.Loc{Start: openTok.Loc.Start, End: closeTok.Loc.End}}, nil
		}
	case kwROLLUP:
		return p.parseRollupCube(ast.GroupingRollup)
	case kwCUBE:
		return p.parseRollupCube(ast.GroupingCube)
	case kwGROUPING:
		// GROUPING SETS ( <set>, … ). Note GROUPING(x) is also a function call;
		// here GROUPING followed by SETS is the grouping-sets form.
		if p.peekNext().Type == kwSETS {
			return p.parseGroupingSets()
		}
	}

	// Plain `expr [AS alias] [ASC|DESC]`.
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	gi := &ast.GroupingItem{Kind: ast.GroupingExpr, Expr: expr, Loc: nodeLoc(expr)}
	if p.cur.Type == kwAS {
		p.advance()
		aliasTok, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}
		alias, err := p.identifierText(aliasTok)
		if err != nil {
			return nil, err
		}
		gi.Alias = alias
		gi.Loc.End = aliasTok.Loc.End
	}
	if p.cur.Type == kwASC || p.cur.Type == kwDESC {
		dirTok := p.advance()
		gi.Loc.End = dirTok.Loc.End
	}
	return gi, nil
}

// parseRollupCube parses `ROLLUP ( expr, … )` or `CUBE ( expr, … )` (rollup_list
// / cube_list). The ROLLUP/CUBE keyword is the current token.
func (p *Parser) parseRollupCube(kind ast.GroupingKind) (*ast.GroupingItem, error) {
	kwTok := p.advance() // ROLLUP / CUBE
	gi := &ast.GroupingItem{Kind: kind, Loc: kwTok.Loc}
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	items, err := p.parseExprCommaList()
	if err != nil {
		return nil, err
	}
	gi.Items = items
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	gi.Loc.End = closeTok.Loc.End
	return gi, nil
}

// parseGroupingSets parses `GROUPING SETS ( <grouping_set>, … )`
// (grouping_set_list). GROUPING is the current token (SETS confirmed by caller).
// Each set is `()`, an expression, or a nested ROLLUP/CUBE — their inner
// expressions are captured flat into Items for walking (the grouping-set
// structure is not consumed by bytebase, so this preserves the expressions for
// span/coverage without a bespoke nested node).
func (p *Parser) parseGroupingSets() (*ast.GroupingItem, error) {
	groupingTok := p.advance() // GROUPING
	p.advance()                // SETS
	gi := &ast.GroupingItem{Kind: ast.GroupingSets, Loc: groupingTok.Loc}
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	for {
		if err := p.parseGroupingSetInto(gi); err != nil {
			return nil, err
		}
		if p.cur.Type != int(',') {
			break
		}
		p.advance()
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	gi.Loc.End = closeTok.Loc.End
	return gi, nil
}

// parseGroupingSetInto parses one grouping_set (`() | expr | ROLLUP(…) | CUBE(…)`)
// and appends its expression(s) to gi.Items. A parenthesized multi-expression
// set `(a, b)` contributes each expression.
func (p *Parser) parseGroupingSetInto(gi *ast.GroupingItem) error {
	switch p.cur.Type {
	case kwROLLUP, kwCUBE:
		p.advance() // ROLLUP / CUBE
		if _, err := p.expect(int('(')); err != nil {
			return err
		}
		items, err := p.parseExprCommaList()
		if err != nil {
			return err
		}
		gi.Items = append(gi.Items, items...)
		if _, err := p.expect(int(')')); err != nil {
			return err
		}
		return nil
	case int('('):
		// `()` empty set, or `(a, b, …)` a multi-column set.
		p.advance() // '('
		if p.cur.Type == int(')') {
			p.advance() // ')'
			return nil
		}
		items, err := p.parseExprCommaList()
		if err != nil {
			return err
		}
		gi.Items = append(gi.Items, items...)
		if _, err := p.expect(int(')')); err != nil {
			return err
		}
		return nil
	default:
		expr, err := p.parseExpr()
		if err != nil {
			return err
		}
		gi.Items = append(gi.Items, expr)
		return nil
	}
}

// ---------------------------------------------------------------------------
// WINDOW clause
// ---------------------------------------------------------------------------

// parseWindowClause parses `WINDOW name AS <spec> (, name AS <spec>)*`
// (window_clause). WINDOW is the current token. Each definition's spec reuses
// the expressions node's parseInlineWindowSpec / named-reference machinery via
// parseWindowSpecBody.
func (p *Parser) parseWindowClause() ([]*ast.WindowDef, error) {
	p.advance() // WINDOW
	var defs []*ast.WindowDef
	for {
		nameTok, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}
		name, err := p.identifierText(nameTok)
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(kwAS); err != nil {
			return nil, err
		}
		spec, err := p.parseWindowSpecBody(nameTok.Loc.Start)
		if err != nil {
			return nil, err
		}
		defs = append(defs, &ast.WindowDef{Name: name, Spec: spec, Loc: ast.Loc{Start: nameTok.Loc.Start, End: spec.Loc.End}})
		if p.cur.Type != int(',') {
			break
		}
		p.advance()
	}
	return defs, nil
}

// parseWindowSpecBody parses a window_specification in a WINDOW definition: a
// parenthesized inline spec `( [base] [PARTITION BY] [ORDER BY] [frame] )`, or a
// bare base-window name reference. startOff anchors Loc. It reuses the
// expressions node's parseInlineWindowSpec for the parenthesized form so the
// PARTITION BY / ORDER BY / frame parsing stays single-sourced.
func (p *Parser) parseWindowSpecBody(startOff int) (*ast.WindowSpec, error) {
	if p.cur.Type == int('(') {
		return p.parseInlineWindowSpec(startOff)
	}
	// Bare named reference.
	if !isIdentifierStart(p.cur.Type) {
		return nil, p.syntaxErrorAtCur()
	}
	nameTok := p.advance()
	name, err := p.identifierText(nameTok)
	if err != nil {
		return nil, err
	}
	return &ast.WindowSpec{Name: name, Loc: ast.Loc{Start: startOff, End: nameTok.Loc.End}}, nil
}
