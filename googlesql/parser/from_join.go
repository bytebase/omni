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
// by zero or more explicit JOIN suffixes. Left-to-right joins with immediate
// ON/USING clauses fold left-deep; consecutive condition-needing joins followed
// by consecutive ON/USING clauses (e.g.
// `a JOIN b JOIN c USING (id) ON a.id = b.id`) reassociate right-to-left as
// each ON/USING closes the most recent open join, producing
// `a JOIN (b JOIN c USING (id)) ON a.id = b.id`.
func (p *Parser) parseJoinTree() (ast.Node, error) {
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()
	cur, err := p.parseTablePrimary()
	if err != nil {
		return nil, err
	}
	var open []*ast.JoinExpr
	for {
		next, ok, err := p.tryParseJoinSuffix(cur, &open)
		if err != nil {
			return nil, err
		}
		if !ok {
			for i := len(open) - 1; i >= 0; i-- {
				j := open[i]
				j.Right = cur
				j.Loc = ast.Loc{Start: qLoc(j.Left).Start, End: qLoc(cur).End}
				cur = j
			}
			return cur, nil
		}
		cur = next
	}
}

// tryParseJoinSuffix parses one explicit join suffix on cur if the current
// position starts one (a [NATURAL] [join_type] [join_hint] JOIN). It returns
// (nextCur, true, nil) on success, (nil, false, nil) when no JOIN follows, or
// (nil, false, err) on a malformed join. A comma is NOT consumed here — the
// top-level comma is handled by parseFromClause (comma cross join binds looser).
func (p *Parser) tryParseJoinSuffix(cur ast.Node, open *[]*ast.JoinExpr) (ast.Node, bool, error) {
	natural, joinType, joinHint, ok, err := p.tryParseJoinKeywords()
	if err != nil {
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}

	je := &ast.JoinExpr{Type: joinType, Natural: natural, JoinHint: joinHint, Left: cur}

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

	// NATURAL and CROSS joins take no criteria and attach directly to cur.
	if natural || joinType == ast.JoinCross {
		je.Right = right
		je.Loc = ast.Loc{Start: qLoc(cur).Start, End: qLoc(right).End}
		return je, true, nil
	}

	// Condition-needing join: push je onto open with cur as its left operand,
	// and let right become the active right-hand tree. Each immediately
	// following ON/USING closes the top open join, reassociating consecutive
	// joins right-to-left (e.g. `a JOIN b JOIN c USING (id) ON a.id = b.id`
	// -> `a JOIN (b JOIN c USING (id)) ON a.id = b.id`). In BigQuery, excess
	// ON/USING clauses beyond the open joins are left unconsumed so they fail.
	*open = append(*open, je)
	cur = right
	for p.cur.Type == kwON || p.cur.Type == kwUSING {
		var target *ast.JoinExpr
		if len(*open) > 0 {
			top := len(*open) - 1
			target = (*open)[top]
			*open = (*open)[:top]
			target.Right = cur
		} else {
			if p.isBigQuery() {
				break
			}
			j, ok := cur.(*ast.JoinExpr)
			if !ok {
				break
			}
			target = j
		}
		switch p.cur.Type {
		case kwON:
			p.advance()
			expr, err := p.parseExpr()
			if err != nil {
				return nil, false, err
			}
			target.On = expr
			target.Using = nil
		case kwUSING:
			cols, err := p.parseUsingClause()
			if err != nil {
				return nil, false, err
			}
			target.Using = cols
			target.On = nil
		}
		target.Loc = ast.Loc{Start: qLoc(target.Left).Start, End: p.prev.Loc.End}
		cur = target
	}
	return cur, true, nil
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
// followed by an optional `[AS] alias`, (for path/subquery/TVF sources) a
// PIVOT/UNPIVOT operator, (for path sources) `WITH OFFSET` / `FOR SYSTEM_TIME AS
// OF`, and finally an optional TABLESAMPLE that wraps the whole source.
// PIVOT/UNPIVOT/TABLESAMPLE are parsed by the parser-query-clauses node
// (pivot_unpivot.go / tablesample.go); MATCH_RECOGNIZE is not yet implemented.
func (p *Parser) parseTablePrimary() (ast.Node, error) {
	var (
		src ast.Node
		err error
	)
	switch {
	case p.cur.Type == int('('):
		src, err = p.parseParenTableSource()
	case p.cur.Type == kwUNNEST:
		src, err = p.parseUnnestTableSource()
	case p.curIsUnquotedWord("LATERAL"):
		src, err = p.parseLateralTableSource()
	case p.curIsUnquotedWord("GRAPH_TABLE") && p.peekNext().Type == int('('):
		src, err = p.parseGraphTableQuery()
	default:
		src, err = p.parsePathTableSource()
	}
	if err != nil {
		return nil, err
	}

	// TABLESAMPLE and MATCH_RECOGNIZE wrap a complete table_primary.
	for {
		if p.atTableSample() {
			sc, serr := p.parseTableSample()
			if serr != nil {
				return nil, serr
			}
			if !attachSample(src, sc) {
				return nil, p.syntaxErrorAtCur()
			}
			continue
		}
		if p.isBigQuery() && p.cur.Type == kwMATCH_RECOGNIZE {
			mr, merr := p.parseMatchRecognizeClause()
			if merr != nil {
				return nil, merr
			}
			if !attachMatchRecognize(src, mr) {
				return nil, p.syntaxErrorAtCur()
			}
			continue
		}
		break
	}
	return src, nil
}

// parseLateralTableSource parses `LATERAL table_subquery` or
// `LATERAL tvf [as_alias] [pivot_or_unpivot_clause]`. cur is at LATERAL.
func (p *Parser) parseLateralTableSource() (ast.Node, error) {
	latTok := p.advance() // LATERAL
	if p.cur.Type == int('(') {
		p.advance() // '('
		if !p.atQueryStart() {
			return nil, p.syntaxErrorAtCur()
		}
		query, err := p.parseQuery()
		if err != nil {
			return nil, err
		}
		closeTok, err := p.expect(int(')'))
		if err != nil {
			return nil, err
		}
		te := &ast.TableExpr{
			Subquery: query,
			Lateral:  true,
			Loc:      ast.Loc{Start: latTok.Loc.Start, End: closeTok.Loc.End},
		}
		if err := p.parseTableAliasOnly(te); err != nil {
			return nil, err
		}
		return te, nil
	}
	// LATERAL tvf: requires '(' after path.
	path, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	if p.cur.Type != int('(') {
		return nil, p.syntaxErrorAtCur()
	}
	fc, err := p.parseFuncCallSuffix(path)
	if err != nil {
		return nil, err
	}
	te := &ast.TableExpr{
		Func:    fc.(*ast.FuncCall),
		Lateral: true,
		Loc:     ast.Loc{Start: latTok.Loc.Start, End: nodeLoc(fc).End},
	}
	if err := p.parseTableAliasOnly(te); err != nil {
		return nil, err
	}
	return te, nil
}

// parseGraphTableQuery parses `GRAPH_TABLE '(' path_expression
// (graph_match_operator graph_shape_clause | graph_operation_block) ')' as_alias?`.
func (p *Parser) parseGraphTableQuery() (ast.Node, error) {
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()
	gtTok := p.advance() // GRAPH_TABLE
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	graphPath, err := p.parsePathExpr()
	if err != nil {
		return nil, err
	}
	var ops []ast.Node
	var cols []*ast.SelectItem
	if p.cur.Type == kwMATCH {
		// Could be `graph_match_operator graph_shape_clause` (no RETURN) or
		// `graph_operation_block` (ending in RETURN).
		matchOp, err := p.parseGraphMatchOp(false)
		if err != nil {
			return nil, err
		}
		if p.cur.Type == kwCOLUMNS || p.cur.Type == int(')') {
			ops = []ast.Node{matchOp}
			if p.cur.Type == kwCOLUMNS {
				p.advance() // COLUMNS
				if _, err := p.expect(int('(')); err != nil {
					return nil, err
				}
				cols, err = p.parseSelectList()
				if err != nil {
					return nil, err
				}
				if _, err := p.expect(int(')')); err != nil {
					return nil, err
				}
			}
		} else {
			linOps := []ast.Node{matchOp}
			for p.atGraphLinearOperator() {
				op, err := p.parseGraphLinearOperator()
				if err != nil {
					return nil, err
				}
				linOps = append(linOps, op)
			}
			firstLin, err := p.finishGraphLinearQuery(linOps)
			if err != nil {
				return nil, err
			}
			firstBlock, err := p.parseGraphSetOpChain(firstLin)
			if err != nil {
				return nil, err
			}
			ops = append(ops, firstBlock)
			for p.cur.Type == kwNEXT {
				p.advance()
				nextBlock, err := p.parseGraphCompositeQueryBlock()
				if err != nil {
					return nil, err
				}
				ops = append(ops, nextBlock)
			}
		}
	} else {
		var err error
		ops, _, err = p.parseGraphOperationBlock()
		if err != nil {
			return nil, err
		}
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	gt := &ast.GraphTableQuery{
		Graph:   graphPath,
		Ops:     ops,
		Columns: cols,
		Loc:     ast.Loc{Start: gtTok.Loc.Start, End: closeTok.Loc.End},
	}
	te := &ast.TableExpr{GraphTable: gt, Loc: gt.Loc}
	if alias, ok, err := p.tryParseTableAlias(); err != nil {
		return nil, err
	} else if ok {
		gt.Alias = alias
		te.Alias = alias
		gt.Loc.End = p.prev.Loc.End
		te.Loc.End = p.prev.Loc.End
	}
	return te, nil
}

// attachSample stores a parsed TABLESAMPLE clause on a table source, updating
// its end Loc. It supports the TableExpr / UnnestExpr source shapes; for any
// other node (e.g. a parenthesized join) it returns false so the caller can
// surface a syntax error.
func attachSample(src ast.Node, sc *ast.SampleClause) bool {
	switch s := src.(type) {
	case *ast.TableExpr:
		s.Sample = sc
		s.Loc.End = sc.Loc.End
	case *ast.UnnestExpr:
		s.Sample = sc
		s.Loc.End = sc.Loc.End
	default:
		return false
	}
	return true
}

// attachMatchRecognize stores a parsed MATCH_RECOGNIZE clause on a table
// source (TableExpr, UnnestExpr, or parenthesized JoinExpr), updating its end
// Loc. Returns false if the source shape is unsupported or already carries a
// MATCH_RECOGNIZE clause so the caller never silently drops it.
func attachMatchRecognize(src ast.Node, mr *ast.MatchRecognizeClause) bool {
	switch s := src.(type) {
	case *ast.TableExpr:
		if s.MatchRecognize != nil {
			return false
		}
		s.MatchRecognize = mr
		s.Loc.End = mr.Loc.End
	case *ast.UnnestExpr:
		if s.MatchRecognize != nil {
			return false
		}
		s.MatchRecognize = mr
		s.Loc.End = mr.Loc.End
	case *ast.JoinExpr:
		if s.MatchRecognize != nil {
			return false
		}
		s.MatchRecognize = mr
		s.Loc.End = mr.Loc.End
	default:
		return false
	}
	return true
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
		// A subquery source takes only an `[AS] alias` — NOT WITH OFFSET /
		// FOR SYSTEM_TIME (oracle: `SELECT * FROM (SELECT 1) FOR SYSTEM_TIME …`
		// rejects).
		if err := p.parseTableAliasOnly(te); err != nil {
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
	// table_primary: '(' join ')' requires an actual JOIN inside the parens in
	// BigQuery (`FROM (a a)` without a JOIN is a syntax error).
	if p.isBigQuery() {
		if _, ok := inner.(*ast.JoinExpr); !ok {
			return nil, p.syntaxErrorAtCur()
		}
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

	// Optional `[AS] alias`, then a PIVOT/UNPIVOT operator (parser-query-clauses),
	// then `WITH OFFSET [[AS] name]` (UNNEST supports the offset companion
	// column). UnnestExpr has its own fields, so it is done inline rather than via
	// the path-source suffix parser. PIVOT/UNPIVOT/TABLESAMPLE on an array scan
	// are always SEMANTIC errors but PARSE in the union grammar (retained for
	// parity); TABLESAMPLE wraps the source in parseTablePrimary.
	if alias, ok, err := p.tryParseTableAlias(); err != nil {
		return nil, err
	} else if ok {
		ue.Alias = alias
		ue.Loc.End = p.prev.Loc.End
	}
	if err := p.parsePivotOrUnpivot(ue); err != nil {
		return nil, err
	}
	if p.cur.Type == kwWITH && p.peekNext().Type == kwOFFSET {
		p.advance()           // WITH
		offTok := p.advance() // OFFSET
		ue.WithOffset = true
		ue.Loc.End = offTok.Loc.End
		if alias, ok, err := p.tryParseTableAlias(); err != nil {
			return nil, err
		} else if ok {
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

	// Table-valued function call: `path ( args )`. A TVF source takes only an
	// `[AS] alias` — NOT WITH OFFSET / FOR SYSTEM_TIME (those are path-source
	// suffixes; oracle: `SELECT * FROM f() WITH OFFSET` rejects).
	if p.cur.Type == int('(') {
		fc, err := p.parseFuncCallSuffix(path)
		if err != nil {
			return nil, err
		}
		te := &ast.TableExpr{Func: fc.(*ast.FuncCall), Loc: nodeLoc(fc)}
		if err := p.parseTableAliasOnly(te); err != nil {
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
	if !isIdentifierStart(p.cur.Type) || isExtraReservedIdent(p.cur) {
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
		// GoogleSQL dashed identifiers forbid whitespace around '-'.
		if p.prev.Loc.End != p.cur.Loc.Start || p.cur.Loc.End != next.Loc.Start {
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

// parseTableAliasOnly parses the optional `[AS] alias` suffix of a table source
// plus an optional PIVOT/UNPIVOT operator (parser-query-clauses), filling
// te.Alias / te.Pivot / te.Unpivot. It is used for subquery and TVF sources,
// which (unlike path / UNNEST sources) take NO WITH OFFSET or FOR SYSTEM_TIME
// suffix (oracle: those reject on subqueries/TVFs) but DO accept PIVOT/UNPIVOT
// (grammar: table_subquery / tvf_with_suffixes carry opt_pivot_or_unpivot_
// clause_and_alias; oracle: `(SELECT …) PIVOT(…)` parses). A per-source `@{...}`
// hint may precede the alias.
func (p *Parser) parseTableAliasOnly(te *ast.TableExpr) error {
	if p.cur.Type == int('@') {
		if herr := p.skipHint(); herr != nil {
			return herr
		}
	}
	alias, ok, err := p.tryParseTableAlias()
	if err != nil {
		return err
	}
	if ok {
		te.Alias = alias
		te.Loc.End = p.prev.Loc.End
	}
	// PIVOT / UNPIVOT operator (after the optional alias).
	return p.parsePivotOrUnpivot(te)
}

// parseTableSuffixes parses the optional trailing suffixes of a PATH table
// source: `[AS] alias`, a PIVOT/UNPIVOT operator, `WITH OFFSET [[AS] name]`, and
// `FOR SYSTEM_TIME AS OF expr`. It fills the corresponding TableExpr fields.
// (PIVOT/UNPIVOT are owned by parser-query-clauses via parsePivotOrUnpivot;
// TABLESAMPLE wraps the whole source in parseTablePrimary. Subquery and TVF
// sources use parseTableAliasOnly instead, since the offset / time-travel
// suffixes are path-source-only.) The grammar's ordering is
// `opt_pivot_or_unpivot_clause_and_alias? opt_with_offset_and_alias?
// opt_at_system_time?`, i.e. alias+pivot bind before WITH OFFSET / FOR SYSTEM_
// TIME.
func (p *Parser) parseTableSuffixes(te *ast.TableExpr) error {
	// Per-source hint @{...} (table_hint_expr) may precede the alias.
	if p.cur.Type == int('@') {
		if herr := p.skipHint(); herr != nil {
			return herr
		}
	}

	// FOR SYSTEM_TIME AS OF expr (opt_at_system_time). In GoogleSQL's grammar
	// (`table_path_expression`), `as_alias` precedes `at_system_time`; BigQuery
	// rejects a table alias after `FOR SYSTEM_TIME AS OF ...`.
	if !p.isBigQuery() && p.atForSystemTime() {
		if err := p.parseForSystemTime(te); err != nil {
			return err
		}
	}

	if alias, ok, err := p.tryParseTableAlias(); err != nil {
		return err
	} else if ok {
		te.Alias = alias
		te.Loc.End = p.prev.Loc.End
	}

	// PIVOT / UNPIVOT operator (after the optional alias, before WITH OFFSET).
	if err := p.parsePivotOrUnpivot(te); err != nil {
		return err
	}

	if p.cur.Type == kwWITH && p.peekNext().Type == kwOFFSET {
		p.advance()           // WITH
		offTok := p.advance() // OFFSET
		te.WithOffset = true
		te.Loc.End = offTok.Loc.End
		if alias, ok, err := p.tryParseTableAlias(); err != nil {
			return err
		} else if ok {
			te.WithOffsetAlias = alias
			te.Loc.End = p.prev.Loc.End
		}
	}

	// FOR SYSTEM_TIME may also appear after the alias in some doc forms; accept
	// it post-alias too.
	if te.SystemTime == nil && p.atForSystemTime() {
		if err := p.parseForSystemTime(te); err != nil {
			return err
		}
	}
	return nil
}

// atForSystemTime reports whether the current position begins a time-travel
// suffix: `FOR SYSTEM_TIME …` (single token) or `FOR SYSTEM TIME …` (two-word).
// It must NOT match `FOR UPDATE` (the query-level row-lock, handled by
// parseQuery), so it requires SYSTEM/SYSTEM_TIME after FOR.
func (p *Parser) atForSystemTime() bool {
	return p.cur.Type == kwFOR && (p.peekNext().Type == kwSYSTEM_TIME || p.peekNext().Type == kwSYSTEM)
}

// parseForSystemTime parses `FOR SYSTEM_TIME AS OF expr` (and the `FOR SYSTEM
// TIME` two-word spelling). FOR is the current token. The TIME word is REQUIRED
// in the two-word form (oracle: `FOR SYSTEM AS OF …` rejects with "Expected
// keyword TIME").
func (p *Parser) parseForSystemTime(te *ast.TableExpr) error {
	p.advance() // FOR
	// SYSTEM_TIME (single token) or SYSTEM TIME (two tokens, TIME required).
	if p.cur.Type == kwSYSTEM_TIME {
		p.advance()
	} else {
		if _, err := p.expect(kwSYSTEM); err != nil {
			return err
		}
		if _, err := p.expect(kwTIME); err != nil {
			return err
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
// (as_alias). It returns (alias, true, nil) when an alias is present,
// ("", false, nil) when none, and ("", false, err) when an explicit AS is given
// without a following identifier — that is a hard syntax error (oracle:
// `SELECT * FROM t AS` rejects), NOT an absent alias: the explicit AS commits to
// an alias, so a missing name must be diagnosed rather than silently dropped. A
// bare (AS-less) alias is consumed only when the current token is an
// identifier-start that does NOT begin a following clause/join keyword (so
// `FROM a JOIN b` does not read JOIN as an alias of a, and `FROM a WHERE …` does
// not read WHERE as an alias). isIdentifierStart already excludes reserved
// keywords (JOIN/WHERE/ON/etc are reserved), so a bare non-reserved word is a
// valid implicit alias.
func (p *Parser) tryParseTableAlias() (string, bool, error) {
	if p.cur.Type == kwAS {
		p.advance()
		aliasTok, err := p.expectIdentifier()
		if err != nil {
			return "", false, err
		}
		alias, err := p.identifierText(aliasTok)
		if err != nil {
			return "", false, err
		}
		return alias, true, nil
	}
	// Implicit alias: a bare identifier that is not a clause/join continuation.
	// In Spanner, QUALIFY is not a clause keyword and is accepted as a bare
	// table alias (e.g. `TABLE a() QUALIFY`).
	if ((isIdentifierStart(p.cur.Type) && !isExtraReservedIdent(p.cur)) || (p.isSpanner() && p.cur.Type == kwQUALIFY)) && !p.atTableAliasStop() {
		aliasTok := p.advance()
		alias, err := p.identifierText(aliasTok)
		if err != nil {
			return "", false, err
		}
		return alias, true, nil
	}
	return "", false, nil
}

// atTableAliasStop reports whether the current (identifier-start) token must NOT
// be consumed as an implicit table alias because it begins a following
// clause/table-operator suffix rather than serving as an alias.
//
//   - WITH always stops: it begins `WITH OFFSET` (the only WITH that follows a
//     table source); a table is never implicitly aliased `WITH`.
//   - PIVOT / UNPIVOT are NON-reserved, so a bare `FROM t pivot` / `FROM t
//     unpivot` is a valid implicit alias (oracle-confirmed accept). They stop
//     alias parsing ONLY when they actually START the operator: `PIVOT (`,
//     `UNPIVOT (`, or `UNPIVOT {INCLUDE|EXCLUDE}`. Anything else (a bare keyword
//     at end of source, before a comma/JOIN/clause) is the implicit alias.
//
// TABLESAMPLE / MATCH_RECOGNIZE are RESERVED keywords (keywordReserved), so they
// are never identifier-start and never reach this predicate as alias candidates;
// they are consumed as operators by parseTablePrimary. Reserved clause keywords
// (FROM/WHERE/JOIN/…) are likewise already excluded by isIdentifierStart.
func (p *Parser) atTableAliasStop() bool {
	switch p.cur.Type {
	case kwWITH:
		return true
	case kwPIVOT:
		// Operator iff `PIVOT (`; otherwise a bare implicit alias named "pivot".
		return p.peekNext().Type == int('(')
	case kwUNPIVOT:
		// Operator iff `UNPIVOT (` or `UNPIVOT INCLUDE|EXCLUDE`; otherwise alias.
		switch p.peekNext().Type {
		case int('('), kwINCLUDE, kwEXCLUDE:
			return true
		}
		return false
	}
	return false
}
