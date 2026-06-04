package parser

import (
	"github.com/bytebase/omni/googlesql/ast"
)

// This file is part of the `parser-select` DAG node. It implements GoogleSQL's
// FROM clause, table sources, and joins (GoogleSQLParser.g4 §2.14 from_clause /
// table_primary / join / join_item, plus the UNNEST table source from §2.15),
// a hand-port of Google's open-source ZetaSQL reference. parseFromClause is
// called from parseSelectStmt (select.go).
//
// FROM-clause shape. GoogleSQL's from_clause_contents is
// `table_primary from_clause_contents_suffix*`, where a suffix is either
// `, table_primary` (comma cross join) or an explicit `[NATURAL] join_type?
// join_hint? JOIN hint? table_primary [on/using]`. Comma cross joins bind
// LOOSER than explicit JOINs (standard SQL precedence: `A, B JOIN C` =
// `A CROSS JOIN (B JOIN C)`), so we model FROM as a comma-separated LIST of
// items (SelectStmt.From []Node), and within each item fold the explicit JOINs
// into a left-deep *JoinExpr tree. This matches the omni snowflake/pg FROM
// representation and is ergonomic for the downstream query-span extractor.
//
// table_primary cases handled here: a (possibly dashed/slashed/dotted) table
// path, a table-valued function call, a parenthesized join `( join )`, a
// parenthesized subquery `( query )`, and UNNEST(...). PIVOT / UNPIVOT /
// TABLESAMPLE / MATCH_RECOGNIZE / AT SYSTEM TIME table-source operators are NOT
// owned by this node (parser-query-clauses); FOR SYSTEM_TIME AS OF is captured
// minimally on TableExpr since it is a path-source time-travel suffix in the
// shared grammar.

// parseFromClause parses a FROM clause (from_clause):
// `FROM <table_primary> <from_clause_contents_suffix>*`. FROM is the current
// token. It returns the comma-separated top-level FROM items (each a join tree).
func (p *Parser) parseFromClause() ([]ast.Node, error) {
	p.advance() // FROM

	// The legacy grammar has error alternatives for `@`/`?`/`@@` used as table
	// names (a parameter is not a table). Those naturally fall out as syntax
	// errors from parseTablePrimary's identifier/path requirement.
	var items []ast.Node
	for {
		item, err := p.parseJoinTree()
		if err != nil {
			return nil, err
		}
		items = append(items, item)
		if p.cur.Type != int(',') {
			break
		}
		p.advance() // consume ',' (comma cross join — a top-level list separator)
	}
	return items, nil
}

// parseJoinTree parses one comma-separated FROM item: a primary source followed
// by zero or more explicit JOIN suffixes, folded into a left-deep *JoinExpr.
func (p *Parser) parseJoinTree() (ast.Node, error) {
	left, err := p.parseTablePrimary()
	if err != nil {
		return nil, err
	}
	for {
		joined, ok, err := p.tryParseJoinSuffix(left)
		if err != nil {
			return nil, err
		}
		if !ok {
			return left, nil
		}
		left = joined
	}
}

// tryParseJoinSuffix parses one explicit join suffix on left if the current
// position starts one (a [NATURAL] [join_type] [join_hint] JOIN). It returns
// (joinExpr, true, nil) on success, (nil, false, nil) when no JOIN follows, or
// (nil, false, err) on a malformed join. A comma is NOT consumed here — the
// top-level comma is handled by parseFromClause (comma cross join binds looser).
func (p *Parser) tryParseJoinSuffix(left ast.Node) (ast.Node, bool, error) {
	natural, joinType, joinHint, ok, err := p.tryParseJoinKeywords()
	if err != nil {
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}

	je := &ast.JoinExpr{Type: joinType, Natural: natural, JoinHint: joinHint, Left: left}

	// Optional per-join `@{...}` hint between JOIN and the right source.
	if p.cur.Type == int('@') {
		if herr := p.skipHint(); herr != nil {
			return nil, false, herr
		}
	}

	right, err := p.parseTablePrimary()
	if err != nil {
		return nil, false, err
	}
	je.Right = right
	je.Loc = ast.Loc{Start: qLoc(left).Start, End: qLoc(right).End}

	// ON / USING criteria. NATURAL and CROSS joins take no criteria; the others
	// require exactly one (the grammar's on_or_using_clause_list is one-or-more,
	// but a single criterion is the universal case — repeated ON/USING is a rare
	// ZetaSQL extension we accept by consuming additional ones into the same
	// node, keeping the last for the structured fields).
	if !natural && joinType != ast.JoinCross {
		if err := p.parseJoinCriteria(je); err != nil {
			return nil, false, err
		}
	}
	if je.On != nil || je.Using != nil {
		je.Loc.End = p.prev.Loc.End
	}
	return je, true, nil
}

// tryParseJoinKeywords parses the join-type keyword prefix
// `[NATURAL] [INNER | CROSS | FULL [OUTER] | LEFT [OUTER] | RIGHT [OUTER]]
// [HASH|LOOKUP] JOIN`. It returns (natural, joinType, joinHint, ok, err). When
// the current position is not a join it returns ok=false having consumed
// nothing. Once a join_type / NATURAL prefix is committed, a missing trailing
// JOIN keyword is a syntax error (returned via err). A bare JOIN (no type) is an
// INNER join.
//
// Set operations only appear after a complete query_primary, never mid-FROM, so
// a leading LEFT/FULL/RIGHT here is always a join (not a corresponding-outer
// set-op prefix).
func (p *Parser) tryParseJoinKeywords() (natural bool, joinType ast.JoinType, joinHint string, ok bool, err error) {
	switch p.cur.Type {
	case kwNATURAL, kwCROSS, kwFULL, kwLEFT, kwRIGHT, kwINNER, kwJOIN, kwHASH, kwLOOKUP:
		// committed to a join below
	default:
		return false, 0, "", false, nil
	}

	// NATURAL prefix.
	if p.cur.Type == kwNATURAL {
		p.advance()
		natural = true
	}

	// join_type.
	jt := ast.JoinInner
	switch p.cur.Type {
	case kwCROSS:
		p.advance()
		jt = ast.JoinCross
	case kwFULL:
		p.advance()
		jt = ast.JoinFull
		if p.cur.Type == kwOUTER {
			p.advance()
		}
	case kwLEFT:
		p.advance()
		jt = ast.JoinLeft
		if p.cur.Type == kwOUTER {
			p.advance()
		}
	case kwRIGHT:
		p.advance()
		jt = ast.JoinRight
		if p.cur.Type == kwOUTER {
			p.advance()
		}
	case kwINNER:
		p.advance()
		jt = ast.JoinInner
	}

	// Optional join_hint HASH | LOOKUP before JOIN (join_hint). (Spanner's
	// `HASH JOIN` form.)
	if p.cur.Type == kwHASH || p.cur.Type == kwLOOKUP {
		joinHint = TokenName(p.cur.Type)
		p.advance()
	}

	// JOIN keyword (required after any join prefix, including a bare JOIN).
	if p.cur.Type != kwJOIN {
		return false, 0, "", false, p.syntaxErrorAtCur()
	}
	p.advance() // JOIN
	return natural, jt, joinHint, true, nil
}

// parseJoinCriteria parses the ON / USING clause(s) of a join (on_or_using_
// clause_list). The current token is just past the right source. It consumes one
// or more ON/USING clauses; the structured fields hold the criteria (multiple
// criteria — a ZetaSQL extension — are all consumed, with On/Using reflecting
// the last, which suffices for parse parity and query-span).
func (p *Parser) parseJoinCriteria(je *ast.JoinExpr) error {
	for {
		switch p.cur.Type {
		case kwON:
			p.advance()
			expr, err := p.parseExpr()
			if err != nil {
				return err
			}
			je.On = expr
			je.Using = nil
		case kwUSING:
			cols, err := p.parseUsingClause()
			if err != nil {
				return err
			}
			je.Using = cols
			je.On = nil
		default:
			return nil
		}
		// A second ON/USING immediately following is the repeated-criteria form.
		if p.cur.Type != kwON && p.cur.Type != kwUSING {
			return nil
		}
	}
}

// parseUsingClause parses `USING ( column (, column)* )` (using_clause). USING is
// the current token. Returns the column names.
func (p *Parser) parseUsingClause() ([]string, error) {
	p.advance() // USING
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	var cols []string
	for {
		colTok, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}
		col, err := p.identifierText(colTok)
		if err != nil {
			return nil, err
		}
		cols = append(cols, col)
		if p.cur.Type != int(',') {
			break
		}
		p.advance()
	}
	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}
	return cols, nil
}

// ---------------------------------------------------------------------------
// Table primaries
// ---------------------------------------------------------------------------

// parseTablePrimary parses a table_primary (the non-suffix part):
//
//	( join | query )         — parenthesized join or subquery
//	UNNEST ( … ) [alias] …   — unnest table source
//	<path> [( args )] …      — table path, or a table-valued function call
//
// followed by an optional `[AS] alias` and (for path sources) `WITH OFFSET` /
// `FOR SYSTEM_TIME AS OF`. PIVOT/UNPIVOT/TABLESAMPLE/MATCH_RECOGNIZE suffixes are
// left for the parser-query-clauses node and are not consumed here.
func (p *Parser) parseTablePrimary() (ast.Node, error) {
	switch p.cur.Type {
	case int('('):
		return p.parseParenTableSource()
	case kwUNNEST:
		return p.parseUnnestTableSource()
	default:
		return p.parsePathTableSource()
	}
}

// parseParenTableSource parses a parenthesized table source: either a
// parenthesized join `( join )` or a parenthesized subquery `( query )`
// (table_subquery). The current token is '('. The disambiguator is atQueryStart:
// a SELECT/WITH inside is a subquery; otherwise it is a parenthesized join.
func (p *Parser) parseParenTableSource() (ast.Node, error) {
	openTok := p.advance() // '('

	if p.atQueryStart() {
		// Parenthesized subquery `( query )`.
		query, err := p.parseQuery()
		if err != nil {
			return nil, err
		}
		closeTok, err := p.expect(int(')'))
		if err != nil {
			return nil, err
		}
		te := &ast.TableExpr{Subquery: query, Loc: ast.Loc{Start: openTok.Loc.Start, End: closeTok.Loc.End}}
		if err := p.parseTableSuffixes(te); err != nil {
			return nil, err
		}
		return te, nil
	}

	// Parenthesized join `( join )`: parse a full join tree, then ')'.
	inner, err := p.parseJoinTree()
	if err != nil {
		return nil, err
	}
	// A parenthesized join may itself be comma-separated inside the parens
	// (`( A, B )` is a valid table_primary join input). Fold top-level commas
	// here into comma cross joins so the parenthesized group is a single Node.
	for p.cur.Type == int(',') {
		p.advance()
		right, err := p.parseJoinTree()
		if err != nil {
			return nil, err
		}
		inner = &ast.JoinExpr{Type: ast.JoinComma, Left: inner, Right: right, Loc: ast.Loc{Start: qLoc(inner).Start, End: qLoc(right).End}}
	}
	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}
	// A parenthesized join can itself be the left of a further join / carry an
	// alias; return it directly (the grammar's table_primary: ( join ) has no
	// trailing alias of its own). Further JOIN suffixes fold via the caller.
	return inner, nil
}

// parseUnnestTableSource parses an UNNEST table source: `UNNEST ( array … )
// [[AS] alias] [WITH OFFSET [[AS] name]]` (unnest_expression as a
// table_path_expression_base). UNNEST is the current token. The UNNEST(...) call
// itself is parsed by the expressions node's parseUnnestExpression (a *FuncCall).
func (p *Parser) parseUnnestTableSource() (ast.Node, error) {
	start := p.cur.Loc.Start
	arr, err := p.parseUnnestExpression()
	if err != nil {
		return nil, err
	}
	ue := &ast.UnnestExpr{Array: arr, Loc: ast.Loc{Start: start, End: nodeLoc(arr).End}}

	// Optional alias and WITH OFFSET. Reuse the shared suffix parser via a
	// TableExpr-shaped staging then copy, but UnnestExpr has its own fields; do
	// it inline.
	if alias, ok := p.tryParseTableAlias(); ok {
		ue.Alias = alias
		ue.Loc.End = p.prev.Loc.End
	}
	if p.cur.Type == kwWITH && p.peekNext().Type == kwOFFSET {
		p.advance()           // WITH
		offTok := p.advance() // OFFSET
		ue.WithOffset = true
		ue.Loc.End = offTok.Loc.End
		if alias, ok := p.tryParseTableAlias(); ok {
			ue.WithOffsetAlias = alias
			ue.Loc.End = p.prev.Loc.End
		}
	}
	return ue, nil
}

// parsePathTableSource parses a path-based table source: a path expression (a
// table name, possibly schema/project-qualified or a correlated array-field
// path, including dashed/slashed BigQuery paths) optionally followed by a `(
// args )` to form a table-valued function call. The current token begins the
// path. Trailing alias / WITH OFFSET / FOR SYSTEM_TIME suffixes are parsed by
// parseTableSuffixes.
func (p *Parser) parsePathTableSource() (ast.Node, error) {
	path, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}

	// Table-valued function call: `path ( args )`.
	if p.cur.Type == int('(') {
		fc, err := p.parseFuncCallSuffix(path)
		if err != nil {
			return nil, err
		}
		te := &ast.TableExpr{Func: fc.(*ast.FuncCall), Loc: nodeLoc(fc)}
		if err := p.parseTableSuffixes(te); err != nil {
			return nil, err
		}
		return te, nil
	}

	te := &ast.TableExpr{Path: path, Loc: path.Loc}
	if err := p.parseTableSuffixes(te); err != nil {
		return nil, err
	}
	return te, nil
}

// parseTablePath parses a table path (maybe_dashed_path_expression /
// table_path_expression_base): a dotted identifier path whose components may be
// joined by '-' (dashed BigQuery paths, e.g. `project-id.dataset.table`) or '/'
// (slashed paths). It builds an ast.PathExpr whose Parts hold the normalized
// components. The dashed/slashed segments are folded into a single component
// spelling so the path round-trips for query-span name resolution.
func (p *Parser) parseTablePath() (*ast.PathExpr, error) {
	if !isIdentifierStart(p.cur.Type) {
		return nil, p.syntaxErrorAtCur()
	}
	first := p.advance()
	part0, err := p.identifierText(first)
	if err != nil {
		return nil, err
	}
	// A dashed first component: `a-b-c` (BigQuery project ids). The lexer emits
	// '-' as a separate token; fold consecutive `- <word/int>` into the part.
	part0 = p.consumeDashedTail(part0)

	path := &ast.PathExpr{Parts: []string{part0}, Loc: first.Loc}
	path.Loc.End = p.prev.Loc.End

	for p.cur.Type == int('.') {
		if !isAnyKeywordIdentifier(p.peekNext().Type) {
			break
		}
		p.advance() // '.'
		partTok := p.advance()
		part, err := p.identifierText(partTok)
		if err != nil {
			return nil, err
		}
		part = p.consumeDashedTail(part)
		path.Parts = append(path.Parts, part)
		path.Loc.End = p.prev.Loc.End
	}
	return path, nil
}

// consumeDashedTail folds a trailing run of `- <word/integer>` into a dashed
// path component (BigQuery dashed paths like `my-project` or `region-us`). It
// appends each `-segment` to base. A '-' not followed by a word/integer is left
// for the caller (it is not part of the path). Each segment uses its SOURCE
// spelling (sliced from the input), not TokenName — so a keyword segment such as
// `project` keeps its lower-case source form rather than the upper-cased keyword
// name, matching how identifiers preserve source case.
func (p *Parser) consumeDashedTail(base string) string {
	out := base
	for p.cur.Type == int('-') {
		next := p.peekNext()
		if !isAnyKeywordIdentifier(next.Type) && next.Type != tokInteger {
			break
		}
		p.advance() // '-'
		segTok := p.advance()
		out += "-" + p.tokenSource(segTok)
	}
	return out
}

// tokenSource returns the source spelling of a path-segment token. For a
// tokIdentifier or tokInteger it uses Token.Str/the raw slice; for a keyword
// token (whose Str is empty) it slices the original input over the token's Loc so
// the source case is preserved (e.g. `project`, not `PROJECT`).
func (p *Parser) tokenSource(tok Token) string {
	if tok.Str != "" {
		return tok.Str
	}
	s := absIndex(p, tok.Loc.Start)
	e := absIndex(p, tok.Loc.End)
	if s < e && e <= len(p.input) {
		return p.input[s:e]
	}
	return TokenName(tok.Type)
}

// parseTableSuffixes parses the optional trailing suffixes of a table source:
// `[AS] alias`, `WITH OFFSET [[AS] name]`, and `FOR SYSTEM_TIME AS OF expr`. It
// fills the corresponding TableExpr fields. (PIVOT/UNPIVOT/TABLESAMPLE/
// MATCH_RECOGNIZE are NOT handled — they belong to parser-query-clauses.)
func (p *Parser) parseTableSuffixes(te *ast.TableExpr) error {
	// Per-source hint @{...} (table_hint_expr) may precede the alias.
	if p.cur.Type == int('@') {
		if herr := p.skipHint(); herr != nil {
			return herr
		}
	}

	// FOR SYSTEM_TIME AS OF expr (opt_at_system_time). Recognized here because it
	// is a shared path-source time-travel suffix; the more exotic AT SYSTEM TIME
	// variants are parser-query-clauses' concern.
	if p.cur.Type == kwFOR && (p.peekNext().Type == kwSYSTEM_TIME || p.peekNext().Type == kwSYSTEM) {
		if err := p.parseForSystemTime(te); err != nil {
			return err
		}
	}

	if alias, ok := p.tryParseTableAlias(); ok {
		te.Alias = alias
		te.Loc.End = p.prev.Loc.End
	}

	if p.cur.Type == kwWITH && p.peekNext().Type == kwOFFSET {
		p.advance()           // WITH
		offTok := p.advance() // OFFSET
		te.WithOffset = true
		te.Loc.End = offTok.Loc.End
		if alias, ok := p.tryParseTableAlias(); ok {
			te.WithOffsetAlias = alias
			te.Loc.End = p.prev.Loc.End
		}
	}

	// FOR SYSTEM_TIME may also appear after the alias in some doc forms; accept
	// it post-alias too.
	if te.SystemTime == nil && p.cur.Type == kwFOR && (p.peekNext().Type == kwSYSTEM_TIME || p.peekNext().Type == kwSYSTEM) {
		if err := p.parseForSystemTime(te); err != nil {
			return err
		}
	}
	return nil
}

// parseForSystemTime parses `FOR SYSTEM_TIME AS OF expr` (and the `FOR SYSTEM
// TIME` two-word spelling). FOR is the current token.
func (p *Parser) parseForSystemTime(te *ast.TableExpr) error {
	p.advance() // FOR
	// SYSTEM_TIME (single token) or SYSTEM TIME (two tokens).
	if p.cur.Type == kwSYSTEM_TIME {
		p.advance()
	} else {
		if _, err := p.expect(kwSYSTEM); err != nil {
			return err
		}
		// optional TIME word
		if p.cur.Type == kwTIME {
			p.advance()
		}
	}
	if _, err := p.expect(kwAS); err != nil {
		return err
	}
	if _, err := p.expect(kwOF); err != nil {
		return err
	}
	expr, err := p.parseExpr()
	if err != nil {
		return err
	}
	te.SystemTime = expr
	te.Loc.End = nodeLoc(expr).End
	return nil
}

// tryParseTableAlias parses an optional table alias `[AS] identifier`
// (as_alias). It returns (alias, true) when an alias is present. A bare alias is
// only consumed when the current token is an identifier-start that does NOT
// begin a following clause/join keyword (so `FROM a JOIN b` does not read JOIN
// as an alias of a, and `FROM a WHERE …` does not read WHERE as an alias).
// isIdentifierStart already excludes reserved keywords (JOIN/WHERE/ON/etc are
// reserved), so a bare non-reserved word is a valid implicit alias.
func (p *Parser) tryParseTableAlias() (string, bool) {
	if p.cur.Type == kwAS {
		p.advance()
		aliasTok, err := p.expectIdentifier()
		if err != nil {
			// AS with no following identifier: leave it; the caller's next expect
			// will surface a clear error. Return no-alias so we don't lose the AS.
			return "", false
		}
		alias, err := p.identifierText(aliasTok)
		if err != nil {
			return "", false
		}
		return alias, true
	}
	// Implicit alias: a bare identifier that is not a clause/join continuation.
	if isIdentifierStart(p.cur.Type) && !p.atTableAliasStop() {
		aliasTok := p.advance()
		alias, err := p.identifierText(aliasTok)
		if err != nil {
			return "", false
		}
		return alias, true
	}
	return "", false
}

// atTableAliasStop reports whether the current (identifier-start) token must NOT
// be consumed as an implicit table alias because it begins a following
// clause/join/table-operator. These are the non-reserved keywords that, in
// table-source position, introduce a suffix rather than serve as an alias:
// WITH (OFFSET), and the table-operator keywords PIVOT/UNPIVOT/TABLESAMPLE/
// MATCH_RECOGNIZE (owned by parser-query-clauses, but must not be eaten as an
// alias here). Reserved clause keywords (FROM/WHERE/JOIN/…) are already excluded
// by isIdentifierStart.
func (p *Parser) atTableAliasStop() bool {
	switch p.cur.Type {
	case kwWITH, kwPIVOT, kwUNPIVOT, kwTABLESAMPLE, kwMATCH_RECOGNIZE:
		return true
	}
	return false
}
