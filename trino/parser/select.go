package parser

import "github.com/bytebase/omni/trino/ast"

// This file is the `parser-select` DAG node (with relation.go, ctes.go,
// window.go and setops.go): it implements Trino's query layer — the `query`,
// `queryNoWith`, `queryTerm`, `queryPrimary`, and `querySpecification` grammar
// rules and their clause helpers (SELECT list, FROM, WHERE, GROUP BY, HAVING,
// the named-WINDOW clause, ORDER BY / OFFSET / LIMIT / FETCH, and set
// operations) — as a hand-written recursive-descent parser over the token
// stream, producing a Query node tree.
//
// Like datatypes.go's DataType and expr.go's Expr, the query nodes are
// PARSER-PACKAGE types (not ast.Node): the Trino ast tag set is closed to the
// ast-core nodes (File/Identifier/QualifiedName), so — matching the
// types/expressions precedent — query nodes live in package parser and carry
// their own Loc. They embed the Expr values produced by the expressions node;
// later DAG nodes (parser-ddl's CREATE-TABLE-AS / CREATE-VIEW, parser-dml's
// INSERT-from-query, analysis, deparse) embed these Query values.
//
// Legacy ANTLR grammar (TrinoParser.g4) the implementation tracks:
//
//	query             : with? queryNoWith ;
//	queryNoWith       : queryTerm (ORDER BY sortItem (, sortItem)*)?
//	                      (OFFSET rowCount (ROW|ROWS)?)?
//	                      (LIMIT limitRowCount | FETCH (FIRST|NEXT) rowCount? (ROW|ROWS) (ONLY|WITH TIES))? ;
//	queryTerm         : queryPrimary                                  # queryTermDefault
//	                  | left INTERSECT setQuantifier? right           # setOperation
//	                  | left (UNION|EXCEPT) setQuantifier? right      # setOperation ;
//	queryPrimary      : querySpecification          # queryPrimaryDefault
//	                  | TABLE qualifiedName         # table
//	                  | VALUES expression (, expression)*  # inlineTable
//	                  | ( queryNoWith )             # subquery ;
//	querySpecification: SELECT setQuantifier? selectItem (, selectItem)*
//	                      (FROM relation (, relation)*)?
//	                      (WHERE booleanExpression)?
//	                      (GROUP BY groupBy)?
//	                      (HAVING booleanExpression)?
//	                      (WINDOW windowDefinition (, windowDefinition)*)? ;
//	selectItem        : expression as_column_alias?            # selectSingle
//	                  | primaryExpression . *  (AS columnAliases)?  # selectAll
//	                  | *                                      # selectAll ;
//	groupBy           : setQuantifier? groupingElement (, groupingElement)* ;
//
// The implementation is adjudicated against the live Trino 481 oracle, not the
// literal legacy grammar. Oracle-confirmed facts baked in:
//
//	S1 (set-op precedence & associativity). INTERSECT binds tighter than
//	   UNION/EXCEPT; within a precedence level operators are left-associative
//	   (`a UNION b EXCEPT c` == `(a UNION b) EXCEPT c`, `a INTERSECT b UNION c`
//	   == `(a INTERSECT b) UNION c`). See setops.go.
//	S2 (ORDER BY / OFFSET / LIMIT / FETCH attach to the WHOLE queryTerm). They
//	   are part of queryNoWith, OUTSIDE queryTerm, so in `SELECT 1 UNION SELECT 2
//	   ORDER BY 1` the ORDER BY orders the union, not the right SELECT.
//	S3 (GROUP BY setQuantifier disambiguation). `ALL`/`DISTINCT` is the
//	   leading setQuantifier ONLY when a grouping element follows it; `GROUP BY
//	   ALL` (nothing after) and `GROUP BY ALL, x` (a comma after) read `ALL` as
//	   an ordinary grouping expression — ALL is a non-reserved keyword usable as
//	   a column reference. `GROUP BY ALL x` reads ALL as the quantifier. Probed
//	   against Trino 481. See parseGroupBy.
//
// Node-scope boundaries (recorded in the migration divergence ledger), all
// following the analysis §5 target policy "implement the legacy grammar scope
// correctly; treat docs-ahead-of-legacy items as explicit deferred P1
// extensions":
//
//	D1 (JSON_TABLE relation — divergence #3). JSON_TABLE is a relationPrimary
//	   valid in Trino 481 but ABSENT from the legacy grammar. It is a separate
//	   P1 extension; this node does not parse it (omni rejects, 481 accepts).
//	D2 (WITH SESSION — divergence #4). `WITH SESSION prop = v <query>` is a 481
//	   query prefix absent from the legacy grammar (where WITH is CTE-only). P1
//	   extension; not parsed here.
//	D3 (set-op CORRESPONDING — divergence #6). `UNION/INTERSECT/EXCEPT
//	   CORRESPONDING` is a 481 set-op modifier absent from the legacy grammar
//	   (where the modifier is setQuantifier = DISTINCT/ALL only). P1 extension;
//	   not parsed here. See setops.go.
//	D4 (MATCH_RECOGNIZE — parser-match-recognize node). sampledRelation's
//	   patternRecognition MATCH_RECOGNIZE subsystem is its own DAG node; here a
//	   relation is its aliasedRelation with the optional TABLESAMPLE only. See
//	   relation.go.

// ---------------------------------------------------------------------------
// Query node hierarchy (parser-package; not ast.Node — see file header)
// ---------------------------------------------------------------------------

// QueryNode is the interface implemented by every Trino query-layer node
// produced by this node's parsers. Span returns the source byte range; concrete
// fields are reached by a Go type switch (mirroring Expr / DataType). It is the
// `query` / `queryTerm` / `queryPrimary` result type.
type QueryNode interface {
	// Span returns the source byte range covered by the query node.
	Span() ast.Loc
	// queryNode is a marker preventing unrelated types from satisfying QueryNode.
	queryNode()
}

// Query is the top-level `query` (with? queryNoWith): an optional WITH clause
// wrapping a body. It is what parseStmt returns for a query statement and what
// a subquery / CTE / set-op operand reduces to.
//
// With is nil when there is no leading WITH. Body is the queryNoWith — a
// QuerySpec, TableQuery, ValuesQuery, set-operation, or a parenthesized query —
// already decorated with this query's ORDER BY / OFFSET / LIMIT / FETCH (which
// belong to queryNoWith, outside the queryTerm; see S2). OrderBy/Offset/Limit/
// Fetch on this node carry the queryNoWith-level modifiers.
type Query struct {
	With    *WithClause // nil when no WITH
	Body    QueryNode   // the queryTerm (QuerySpec / set-op / parenthesized / TABLE / VALUES)
	OrderBy []SortItem  // queryNoWith ORDER BY, nil when absent
	Offset  *OffsetClause
	Limit   *LimitClause
	Fetch   *FetchClause
	Loc     ast.Loc
}

func (n *Query) Span() ast.Loc { return n.Loc }
func (*Query) queryNode()      {}

// OffsetClause is `OFFSET count [ROW|ROWS]` (the OFFSET part of queryNoWith).
// Count is an INTEGER_VALUE_ or a `?` parameter (the rowCount rule). RowKeyword
// records the optional ROW/ROWS spelling for deparse fidelity.
type OffsetClause struct {
	Count      Expr   // integer literal or parameter
	RowKeyword string // "", "ROW", or "ROWS"
	Loc        ast.Loc
}

// LimitClause is `LIMIT (ALL | count)` (the LIMIT part of queryNoWith). All
// marks `LIMIT ALL`; Count is the rowCount otherwise (nil when All is true).
type LimitClause struct {
	All   bool
	Count Expr // integer literal or parameter; nil for LIMIT ALL
	Loc   ast.Loc
}

// FetchClause is `FETCH (FIRST|NEXT) [count] (ROW|ROWS) (ONLY | WITH TIES)`
// (the FETCH part of queryNoWith). Count is nil for the count-less form
// (`FETCH FIRST ROW ONLY`). WithTies marks the WITH TIES spelling (vs ONLY).
type FetchClause struct {
	Keyword    string // "FIRST" or "NEXT"
	Count      Expr   // rowCount, nil when absent
	RowKeyword string // "ROW" or "ROWS"
	WithTies   bool   // true for WITH TIES, false for ONLY
	Loc        ast.Loc
}

// QuerySpec is a `querySpecification` (queryPrimaryDefault): a SELECT block.
// SetQuantifier is "", "DISTINCT", or "ALL". Items is the non-empty select
// list. From is the comma-separated relation list (nil when no FROM). Where /
// Having carry the optional predicates. GroupBy is nil when absent. Windows is
// the named-WINDOW clause (nil when absent).
type QuerySpec struct {
	SetQuantifier string // "", "DISTINCT", or "ALL"
	Items         []SelectItem
	From          []Relation // nil when no FROM clause
	Where         Expr       // nil when absent
	GroupBy       *GroupBy   // nil when absent
	Having        Expr       // nil when absent
	Windows       []WindowDefinition
	Loc           ast.Loc
}

func (n *QuerySpec) Span() ast.Loc { return n.Loc }
func (*QuerySpec) queryNode()      {}

// TableQuery is the `TABLE qualifiedName` query primary (the `table` rule): a
// shorthand for `SELECT * FROM qualifiedName`.
type TableQuery struct {
	Name *ast.QualifiedName
	Loc  ast.Loc
}

func (n *TableQuery) Span() ast.Loc { return n.Loc }
func (*TableQuery) queryNode()      {}

// ValuesQuery is the `VALUES expression (, expression)*` query primary (the
// `inlineTable` rule): an inline row set. Each row is an expression (commonly a
// rowConstructor for multi-column rows, e.g. `VALUES (1, 'a'), (2, 'b')`).
type ValuesQuery struct {
	Rows []Expr
	Loc  ast.Loc
}

func (n *ValuesQuery) Span() ast.Loc { return n.Loc }
func (*ValuesQuery) queryNode()      {}

// ParenQuery is the `( queryNoWith )` query primary (the `subquery` rule): a
// parenthesized query expression. Inner is the wrapped queryNoWith (a Query
// without a WITH clause but possibly with its own ORDER BY / LIMIT etc.).
type ParenQuery struct {
	Inner *Query
	Loc   ast.Loc
}

func (n *ParenQuery) Span() ast.Loc { return n.Loc }
func (*ParenQuery) queryNode()      {}

// SelectItem is one entry of a SELECT list. Kind distinguishes the three
// selectItem alternatives. For SelectSingle, Expr is the projected expression
// and Alias the optional `[AS] column_alias`. For SelectAllFrom (`expr.*`),
// Expr is the row-valued primary and Aliases the optional `AS (a, b, …)`. For
// SelectAll (bare `*`), all fields are zero.
type SelectItem struct {
	Kind    SelectItemKind
	Expr    Expr              // the projected/row expression (nil for bare SelectAll)
	Alias   *ast.Identifier   // SelectSingle [AS] alias, nil when absent
	Aliases []*ast.Identifier // SelectAllFrom AS (a, b, …), nil when absent
	Loc     ast.Loc
}

// SelectItemKind classifies a SelectItem.
type SelectItemKind int

const (
	// SelectSingle is `expression [[AS] column_alias]` (selectSingle).
	SelectSingle SelectItemKind = iota
	// SelectAllFrom is `primaryExpression . * [AS columnAliases]` — all columns
	// of a row-valued expression (the first selectAll alternative).
	SelectAllFrom
	// SelectAll is the bare `*` — all columns (the second selectAll alternative).
	SelectAll
)

// GroupBy is the `GROUP BY groupBy` clause body (the groupBy rule):
// `setQuantifier? groupingElement (, groupingElement)*`. SetQuantifier is "",
// "DISTINCT", or "ALL" (the leading quantifier, disambiguated per S3). Elements
// is the non-empty list of grouping elements.
type GroupBy struct {
	SetQuantifier string // "", "DISTINCT", or "ALL"
	Elements      []GroupingElement
	Loc           ast.Loc
}

// GroupingElementKind classifies a GroupingElement.
type GroupingElementKind int

const (
	// GroupingExpr is a single grouping set `( e, … )` or a bare expression
	// (the singleGroupingSet rule via groupingSet). Exprs holds the expressions;
	// for the bare-expression form it is a single element.
	GroupingExprKind GroupingElementKind = iota
	// GroupingRollup is `ROLLUP ( e, … )` (the rollup rule); the list may be empty.
	GroupingRollup
	// GroupingCube is `CUBE ( e, … )` (the cube rule); the list may be empty.
	GroupingCube
	// GroupingSets is `GROUPING SETS ( groupingSet, … )` (the multipleGroupingSets
	// rule). Sets holds each grouping set's expression list.
	GroupingSets
)

// GroupingElement is one element of a GROUP BY (the groupingElement rule):
// a single grouping set / bare expression, a ROLLUP, a CUBE, or a GROUPING SETS.
type GroupingElement struct {
	Kind  GroupingElementKind
	Exprs []Expr   // for GroupingExprKind / GroupingRollup / GroupingCube
	Sets  [][]Expr // for GroupingSets — each inner slice is one grouping set
	Loc   ast.Loc
}

// ---------------------------------------------------------------------------
// Statement-dispatch wiring
// ---------------------------------------------------------------------------

// parseQueryStmt parses a top-level query statement and returns it as an
// ast.Node (the dispatch contract of parseStmt). The leading token is SELECT,
// WITH, TABLE, VALUES, or '('. The WITH-FUNCTION inline-routine prefix
// (`WITH FUNCTION …`) is the parser-routines node; here a leading WITH is a CTE
// list only. parseQueryStmt wraps the parsed *Query in a QueryStmt so the File
// statement list holds an ast.Node.
func (p *Parser) parseQueryStmt() (ast.Node, error) {
	start := p.cur.Loc.Start
	q, err := p.parseQuery()
	if err != nil {
		return nil, err
	}
	return &QueryStmt{Query: q, Loc: ast.Loc{Start: start, End: q.Loc.End}}, nil
}

// QueryStmt is the ast.Node wrapper for a top-level query statement (the
// statementDefault / rootQuery alternative). It adapts a parser-package *Query
// to the ast.Node interface parseStmt returns, mirroring how the foundation's
// File holds ast.Node statements. Query is the parsed query body.
type QueryStmt struct {
	Query *Query
	Loc   ast.Loc
}

// Tag implements ast.Node.
func (n *QueryStmt) Tag() ast.NodeTag { return ast.T_Invalid }

// Compile-time assertion that *QueryStmt satisfies ast.Node.
var _ ast.Node = (*QueryStmt)(nil)

// ---------------------------------------------------------------------------
// query / queryNoWith
// ---------------------------------------------------------------------------

// parseQuery parses the `query` rule (`with? queryNoWith`). A leading WITH (not
// WITH FUNCTION) introduces a CTE list parsed by parseWithClause (ctes.go); the
// body is then a queryNoWith.
func (p *Parser) parseQuery() (*Query, error) {
	var with *WithClause
	start := p.cur.Loc.Start
	if p.cur.Kind == kwWITH {
		// WITH FUNCTION (the rootQuery inline-routine prefix, parser-routines
		// node) is valid ONLY at the statement level — parseStmt routes it to
		// parseWithFunctionQuery there. Reaching parseQuery with a leading
		// WITH FUNCTION means it appeared in a subquery / nested-query position,
		// which Trino 481 rejects with a SYNTAX_ERROR; a CTE WITH is
		// `WITH [RECURSIVE] namedQuery …`. Report the syntax error rather than
		// mis-parsing it as a CTE.
		if p.peekNext().Kind == kwFUNCTION {
			return nil, p.syntaxErrorAtCur()
		}
		w, err := p.parseWithClause()
		if err != nil {
			return nil, err
		}
		with = w
	}
	return p.parseQueryNoWith(with, start)
}

// parseQueryNoWith parses `queryNoWith` (a queryTerm followed by the optional
// ORDER BY / OFFSET / LIMIT / FETCH tail) and attaches the already-parsed
// optional WITH clause. start is the byte offset of the whole query (the WITH
// keyword or the queryTerm's first token).
//
// The ORDER BY / OFFSET / LIMIT / FETCH attach to the WHOLE queryTerm (S2), so
// they are parsed here, after the (possibly set-operation) term — not inside a
// querySpecification.
func (p *Parser) parseQueryNoWith(with *WithClause, start int) (*Query, error) {
	body, err := p.parseQueryTerm()
	if err != nil {
		return nil, err
	}

	q := &Query{
		With: with,
		Body: body,
		Loc:  ast.Loc{Start: start, End: body.Span().End},
	}

	if p.cur.Kind == kwORDER {
		items, err := p.parseOrderByItems()
		if err != nil {
			return nil, err
		}
		q.OrderBy = items
		if len(items) > 0 {
			q.Loc.End = items[len(items)-1].Loc.End
		}
	}

	if p.cur.Kind == kwOFFSET {
		off, err := p.parseOffset()
		if err != nil {
			return nil, err
		}
		q.Offset = off
		q.Loc.End = off.Loc.End
	}

	switch p.cur.Kind {
	case kwLIMIT:
		lim, err := p.parseLimit()
		if err != nil {
			return nil, err
		}
		q.Limit = lim
		q.Loc.End = lim.Loc.End
	case kwFETCH:
		fetch, err := p.parseFetch()
		if err != nil {
			return nil, err
		}
		q.Fetch = fetch
		q.Loc.End = fetch.Loc.End
	}

	return q, nil
}

// parseOffset parses `OFFSET rowCount (ROW|ROWS)?` (OFFSET is current). The
// count is an INTEGER_VALUE_ or a `?` parameter (the rowCount rule).
func (p *Parser) parseOffset() (*OffsetClause, error) {
	offTok := p.advance() // consume OFFSET
	count, err := p.parseRowCount()
	if err != nil {
		return nil, err
	}
	oc := &OffsetClause{Count: count, Loc: ast.Loc{Start: offTok.Loc.Start, End: count.Span().End}}
	if tok, ok := p.match(kwROW, kwROWS); ok {
		oc.RowKeyword = tok.Str
		oc.Loc.End = tok.Loc.End
	}
	return oc, nil
}

// parseLimit parses `LIMIT (ALL | rowCount)` (LIMIT is current).
func (p *Parser) parseLimit() (*LimitClause, error) {
	limTok := p.advance() // consume LIMIT
	if p.cur.Kind == kwALL {
		allTok := p.advance()
		return &LimitClause{All: true, Loc: ast.Loc{Start: limTok.Loc.Start, End: allTok.Loc.End}}, nil
	}
	count, err := p.parseRowCount()
	if err != nil {
		return nil, err
	}
	return &LimitClause{Count: count, Loc: ast.Loc{Start: limTok.Loc.Start, End: count.Span().End}}, nil
}

// parseFetch parses `FETCH (FIRST|NEXT) rowCount? (ROW|ROWS) (ONLY | WITH TIES)`
// (FETCH is current).
func (p *Parser) parseFetch() (*FetchClause, error) {
	fetchTok := p.advance() // consume FETCH
	kwTok, ok := p.match(kwFIRST, kwNEXT)
	if !ok {
		return nil, p.exprErrorAt("expected FIRST or NEXT after FETCH")
	}
	fc := &FetchClause{Keyword: kwTok.Str, Loc: ast.Loc{Start: fetchTok.Loc.Start}}

	// Optional row count before ROW/ROWS. Present only when the next token starts
	// a rowCount (an integer or '?'); ROW/ROWS follows directly otherwise.
	if p.cur.Kind == tokInteger || p.cur.Kind == tokQuestion || p.cur.Kind == int('?') {
		count, err := p.parseRowCount()
		if err != nil {
			return nil, err
		}
		fc.Count = count
	}

	rowTok, ok := p.match(kwROW, kwROWS)
	if !ok {
		return nil, p.exprErrorAt("expected ROW or ROWS in FETCH clause")
	}
	fc.RowKeyword = rowTok.Str

	switch p.cur.Kind {
	case kwONLY:
		onlyTok := p.advance()
		fc.Loc.End = onlyTok.Loc.End
	case kwWITH:
		p.advance() // consume WITH
		tiesTok, err := p.expect(kwTIES)
		if err != nil {
			return nil, err
		}
		fc.WithTies = true
		fc.Loc.End = tiesTok.Loc.End
	default:
		return nil, p.exprErrorAt("expected ONLY or WITH TIES in FETCH clause")
	}
	return fc, nil
}

// parseRowCount parses the `rowCount` rule: an INTEGER_VALUE_ or a `?`
// parameter. Used by OFFSET / LIMIT / FETCH.
func (p *Parser) parseRowCount() (Expr, error) {
	switch p.cur.Kind {
	case tokInteger:
		tok := p.advance()
		return &Literal{Kind: LiteralInteger, Value: tok.Str, Loc: tok.Loc}, nil
	case tokQuestion, int('?'):
		tok := p.advance()
		return &Parameter{Loc: tok.Loc}, nil
	default:
		return nil, p.exprErrorAt("expected an integer or '?' row count")
	}
}

// ---------------------------------------------------------------------------
// querySpecification
// ---------------------------------------------------------------------------

// parseQuerySpecification parses a `querySpecification` (SELECT block). SELECT is
// the current token. The clause order is fixed: SELECT [quantifier] items
// [FROM …] [WHERE …] [GROUP BY …] [HAVING …] [WINDOW …].
func (p *Parser) parseQuerySpecification() (*QuerySpec, error) {
	selTok := p.advance() // consume SELECT
	spec := &QuerySpec{Loc: ast.Loc{Start: selTok.Loc.Start}}

	// Optional setQuantifier (DISTINCT | ALL). Here ALL/DISTINCT before the select
	// list is unambiguously the quantifier (a select item must follow either way),
	// matching the grammar `SELECT setQuantifier? selectItem …`.
	if tok, ok := p.match(kwDISTINCT, kwALL); ok {
		spec.SetQuantifier = tok.Str
	}

	items, err := p.parseSelectItems()
	if err != nil {
		return nil, err
	}
	spec.Items = items
	spec.Loc.End = items[len(items)-1].Loc.End

	if p.cur.Kind == kwFROM {
		p.advance() // consume FROM
		rels, err := p.parseRelationList()
		if err != nil {
			return nil, err
		}
		spec.From = rels
		spec.Loc.End = rels[len(rels)-1].Span().End
	}

	if p.cur.Kind == kwWHERE {
		p.advance() // consume WHERE
		where, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		spec.Where = where
		spec.Loc.End = where.Span().End
	}

	if p.cur.Kind == kwGROUP {
		gb, err := p.parseGroupBy()
		if err != nil {
			return nil, err
		}
		spec.GroupBy = gb
		spec.Loc.End = gb.Loc.End
	}

	if p.cur.Kind == kwHAVING {
		p.advance() // consume HAVING
		having, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		spec.Having = having
		spec.Loc.End = having.Span().End
	}

	if p.cur.Kind == kwWINDOW {
		defs, err := p.parseWindowClause()
		if err != nil {
			return nil, err
		}
		spec.Windows = defs
		if len(defs) > 0 {
			spec.Loc.End = defs[len(defs)-1].Loc.End
		}
	}

	return spec, nil
}

// parseSelectItems parses `selectItem (, selectItem)*` — the non-empty SELECT
// list.
func (p *Parser) parseSelectItems() ([]SelectItem, error) {
	first, err := p.parseSelectItem()
	if err != nil {
		return nil, err
	}
	items := []SelectItem{first}
	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		next, err := p.parseSelectItem()
		if err != nil {
			return nil, err
		}
		items = append(items, next)
	}
	return items, nil
}

// parseSelectItem parses one `selectItem`, disambiguating the three alternatives:
//
//   - bare `*` → SelectAll.
//   - `primaryExpression . *` → SelectAllFrom (`t.*`, `a.b.*`, `(row).*`), with
//     an optional `AS (col, …)`.
//   - `expression [[AS] column_alias]` → SelectSingle.
//
// The `expr.*` form is the selectAll alternative `primaryExpression DOT_
// ASTERISK_`. Two readings reach it:
//   - a dotted NAME ending in `.*` (`t.*`, `cat.sch.tbl.*`). The expression
//     parser's qualifiedName reader (parseIdentifierOrFunction → parseQualifiedName)
//     greedily commits to `.identifier` and so would fail on the trailing `.*`;
//     parseSelectAllName therefore speculatively reads the dotted name and, when
//     it ends in `.*`, builds the row reference itself.
//   - a non-name primary followed by `.*` (`(row).*`), where the expression
//     parser stops the dotted chain before `.*` (parsePostfix leaves the `.`),
//     so the trailing `. *` is consumed after parseExpr returns.
func (p *Parser) parseSelectItem() (SelectItem, error) {
	// Bare `*` — all columns.
	if p.cur.Kind == int('*') {
		starTok := p.advance()
		return SelectItem{Kind: SelectAll, Loc: starTok.Loc}, nil
	}

	// `name . *` — a dotted name whose final step is `.*`. Detected speculatively
	// because the expression parser's qualifiedName reader cannot stop before a
	// trailing `.*`.
	if isIdentifierStart(p.cur.Kind) {
		if item, ok, err := p.tryParseSelectAllName(); err != nil {
			return SelectItem{}, err
		} else if ok {
			return item, nil
		}
	}

	expr, err := p.parseExpr()
	if err != nil {
		return SelectItem{}, err
	}

	// `(row) . *` — a non-name primary followed by `.*`.
	if p.cur.Kind == int('.') && p.peekNext().Kind == int('*') {
		p.advance()            // consume '.'
		starTok := p.advance() // consume '*'
		item := SelectItem{
			Kind: SelectAllFrom,
			Expr: expr,
			Loc:  ast.Loc{Start: expr.Span().Start, End: starTok.Loc.End},
		}
		if err := p.parseSelectAllAliases(&item); err != nil {
			return SelectItem{}, err
		}
		return item, nil
	}

	// `expression [[AS] column_alias]` — a projected expression with an optional
	// alias.
	item := SelectItem{Kind: SelectSingle, Expr: expr, Loc: expr.Span()}
	if alias, ok, err := p.parseOptionalColumnAlias(); err != nil {
		return SelectItem{}, err
	} else if ok {
		item.Alias = alias
		item.Loc.End = alias.Loc.End
	}
	return item, nil
}

// tryParseSelectAllName speculatively parses a dotted name that ends in `.*`
// (the `primaryExpression DOT_ ASTERISK_` selectAll form, name reading). The
// current token is an identifier-start. It reads `identifier (. identifier)*`
// and, only when the very next step is `. *`, builds a SelectAllFrom whose Expr
// is the column-reference chain for the leading parts. When the name is not
// followed by `. *` (it is an ordinary column reference, function call, etc.) it
// rewinds and returns ok=false so the caller parses the item as an expression.
func (p *Parser) tryParseSelectAllName() (SelectItem, bool, error) {
	cp := p.checkpoint()

	first := identFromToken(p.advance())
	parts := []*ast.Identifier{first}
	for p.cur.Kind == int('.') {
		// A `. *` ends the name as a row reference; stop and build SelectAllFrom.
		if p.peekNext().Kind == int('*') {
			p.advance()            // consume '.'
			starTok := p.advance() // consume '*'
			name := &ast.QualifiedName{Parts: parts, Loc: ast.Loc{Start: parts[0].Loc.Start, End: parts[len(parts)-1].Loc.End}}
			item := SelectItem{
				Kind: SelectAllFrom,
				Expr: p.columnRefFromName(name),
				Loc:  ast.Loc{Start: parts[0].Loc.Start, End: starTok.Loc.End},
			}
			if err := p.parseSelectAllAliases(&item); err != nil {
				return SelectItem{}, false, err
			}
			return item, true, nil
		}
		// A `. identifier` continues the dotted name. If a non-identifier follows
		// the dot (and it is not the `*` handled above), this is not a clean dotted
		// name (it may be `m['k']`, a call, etc.) — rewind to the expression path.
		if !isIdentifierStart(p.peekNext().Kind) {
			p.restore(cp)
			return SelectItem{}, false, nil
		}
		p.advance() // consume '.'
		parts = append(parts, identFromToken(p.advance()))
	}

	// The name did not end in `. *` (no trailing `.*`); it is an ordinary
	// expression (column reference, function call, typeConstructor, …). Rewind.
	p.restore(cp)
	return SelectItem{}, false, nil
}

// parseSelectAllAliases parses the optional `AS ( col, … )` column-alias list of
// a `primaryExpression . *` select item and stores it on item.
func (p *Parser) parseSelectAllAliases(item *SelectItem) error {
	if p.cur.Kind != kwAS {
		return nil
	}
	p.advance() // consume AS
	aliases, end, err := p.parseColumnAliases()
	if err != nil {
		return err
	}
	item.Aliases = aliases
	item.Loc.End = end
	return nil
}

// parseOptionalColumnAlias parses an optional `[AS] column_alias` after a select
// expression (the as_column_alias rule). It returns ok=false (no alias) when the
// next token is neither AS nor an identifier-start. A bare identifier following
// an expression is an alias; an explicit AS makes the following identifier
// mandatory.
//
// As in a relation alias, a bare non-reserved keyword that introduces a
// following clause (WINDOW/LIMIT/OFFSET/FETCH in clause-introducing position —
// atRelationClauseStart) is NOT consumed as a column alias, so `SELECT 2 LIMIT 3`
// reads `2` then the LIMIT clause (not `2 AS limit` with a dangling `3`), while
// the bare `SELECT 2 limit` still aliases `2` as "limit". Oracle-probed against
// Trino 481. An explicit AS still forces even a clause keyword to be the alias.
func (p *Parser) parseOptionalColumnAlias() (*ast.Identifier, bool, error) {
	if p.cur.Kind == kwAS {
		p.advance() // consume AS — an identifier MUST follow (even a clause keyword)
		alias, err := p.parseIdentifier()
		if err != nil {
			return nil, false, err
		}
		return alias, true, nil
	}
	if isIdentifierStart(p.cur.Kind) && !p.atRelationClauseStart() {
		return identFromToken(p.advance()), true, nil
	}
	return nil, false, nil
}

// parseColumnAliases parses `( identifier (, identifier)* )` (the columnAliases
// rule), the opening '(' the current token. Returns the identifiers and the
// closing ')' end offset. Used by the `expr.* AS (…)` select item, relation
// aliases, and CTE column lists.
func (p *Parser) parseColumnAliases() ([]*ast.Identifier, int, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, 0, err
	}
	first, err := p.parseIdentifier()
	if err != nil {
		return nil, 0, err
	}
	cols := []*ast.Identifier{first}
	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		next, err := p.parseIdentifier()
		if err != nil {
			return nil, 0, err
		}
		cols = append(cols, next)
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, 0, err
	}
	return cols, closeTok.Loc.End, nil
}

// ---------------------------------------------------------------------------
// GROUP BY
// ---------------------------------------------------------------------------

// parseGroupBy parses `GROUP BY groupBy` (GROUP is current). The body is
// `setQuantifier? groupingElement (, groupingElement)*`.
//
// setQuantifier disambiguation (S3): `ALL`/`DISTINCT` is the leading quantifier
// ONLY when a grouping element follows it on the same level. Because ALL is a
// non-reserved keyword usable as a column reference (and DISTINCT, though
// reserved, only appears here as the quantifier), the quantifier is taken only
// when the token after ALL/DISTINCT begins another grouping element (i.e. is not
// ',' and not a clause/segment terminator). Probed against Trino 481:
// `GROUP BY ALL x` → ALL is the quantifier; `GROUP BY ALL` and `GROUP BY ALL, x`
// → ALL is a grouping expression.
func (p *Parser) parseGroupBy() (*GroupBy, error) {
	groupTok := p.advance() // consume GROUP
	if _, err := p.expect(kwBY); err != nil {
		return nil, err
	}
	gb := &GroupBy{Loc: ast.Loc{Start: groupTok.Loc.Start}}

	if (p.cur.Kind == kwALL || p.cur.Kind == kwDISTINCT) && startsGroupingElement(p.peekNext().Kind) {
		gb.SetQuantifier = p.advance().Str
	}

	first, err := p.parseGroupingElement()
	if err != nil {
		return nil, err
	}
	gb.Elements = []GroupingElement{first}
	gb.Loc.End = first.Loc.End
	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		next, err := p.parseGroupingElement()
		if err != nil {
			return nil, err
		}
		gb.Elements = append(gb.Elements, next)
		gb.Loc.End = next.Loc.End
	}
	return gb, nil
}

// startsGroupingElement reports whether kind can begin a groupingElement: an
// expression-start, a '(' (parenthesized grouping set), or one of ROLLUP / CUBE
// / GROUPING (the keyword grouping forms). Used to disambiguate the GROUP BY
// setQuantifier (S3) — ALL/DISTINCT is the quantifier only when a grouping
// element follows.
func startsGroupingElement(kind TokenKind) bool {
	switch kind {
	case kwROLLUP, kwCUBE, kwGROUPING, int('('):
		return true
	default:
		return startsExpression(kind)
	}
}

// startsExpression reports whether kind can begin an `expression`
// (booleanExpression → … → primaryExpression). It enumerates the leading-token
// vocabulary of parsePrimaryAtom (expr.go) plus the prefix operators NOT / + / -
// that begin a boolean/value expression. Used where a token must be classified
// as "starts an expression" without committing to a parse (e.g. the GROUP BY
// setQuantifier disambiguation, S3).
//
// It is deliberately permissive on the keyword-led atoms (CASE, CAST, EXISTS,
// the special built-ins, the JSON functions, ARRAY/ROW constructors) and on any
// identifier-start (an unquoted/quoted name or a non-reserved keyword usable as
// a column reference). A token it returns true for may still fail to parse as a
// complete expression — the caller commits and surfaces any error normally.
func startsExpression(kind TokenKind) bool {
	switch kind {
	case kwNOT, // logicalNot prefix
		int('+'), int('-'), // arithmeticUnary prefix
		kwNULL, kwTRUE, kwFALSE,
		tokString, tokUnicodeString, tokInteger, tokDecimal, tokDouble, tokBinaryLiteral,
		tokQuestion, int('?'),
		kwINTERVAL, int('('),
		kwARRAY, kwROW, kwCASE, kwCAST, kwTRY_CAST, kwEXISTS,
		kwEXTRACT, kwTRIM, kwNORMALIZE, kwLISTAGG, kwGROUPING,
		kwJSON_EXISTS, kwJSON_VALUE, kwJSON_QUERY, kwJSON_OBJECT, kwJSON_ARRAY,
		kwRUNNING, kwFINAL,
		kwCURRENT_DATE, kwCURRENT_USER, kwCURRENT_CATALOG, kwCURRENT_SCHEMA, kwCURRENT_PATH,
		kwCURRENT_TIME, kwCURRENT_TIMESTAMP, kwLOCALTIME, kwLOCALTIMESTAMP,
		kwDOUBLE:
		return true
	default:
		return isIdentifierStart(kind)
	}
}

// parseGroupingElement parses one `groupingElement`: ROLLUP(…), CUBE(…),
// GROUPING SETS(…), a parenthesized grouping set `( e, … )`, or a bare
// expression.
func (p *Parser) parseGroupingElement() (GroupingElement, error) {
	switch p.cur.Kind {
	case kwROLLUP:
		return p.parseRollupCube(GroupingRollup)
	case kwCUBE:
		return p.parseRollupCube(GroupingCube)
	case kwGROUPING:
		// GROUPING SETS ( … ). A bare GROUPING(...) is the groupingOperation
		// expression (function.go); GROUPING SETS is the grouping-element form,
		// recognised by the SETS keyword following GROUPING.
		if p.peekNext().Kind == kwSETS {
			return p.parseGroupingSets()
		}
		// GROUPING not followed by SETS is a groupingOperation expression used as
		// a single grouping expression (e.g. GROUP BY GROUPING(x)).
		return p.parseGroupingSetOrExpr()
	default:
		return p.parseGroupingSetOrExpr()
	}
}

// parseRollupCube parses `ROLLUP ( e, … )` or `CUBE ( e, … )` (kind selects
// which; the keyword is current). The expression list may be empty.
func (p *Parser) parseRollupCube(kind GroupingElementKind) (GroupingElement, error) {
	kwTok := p.advance() // consume ROLLUP / CUBE
	if _, err := p.expect(int('(')); err != nil {
		return GroupingElement{}, err
	}
	exprs, closeTok, err := p.parseBracketedExprList(int(')'))
	if err != nil {
		return GroupingElement{}, err
	}
	return GroupingElement{
		Kind:  kind,
		Exprs: exprs,
		Loc:   ast.Loc{Start: kwTok.Loc.Start, End: closeTok.Loc.End},
	}, nil
}

// parseGroupingSets parses `GROUPING SETS ( groupingSet (, groupingSet)* )`
// (GROUPING is current, SETS confirmed by the caller). Each grouping set is a
// parenthesized expression list `( e, … )` or a bare expression.
func (p *Parser) parseGroupingSets() (GroupingElement, error) {
	groupingTok := p.advance() // consume GROUPING
	if _, err := p.expect(kwSETS); err != nil {
		return GroupingElement{}, err
	}
	if _, err := p.expect(int('(')); err != nil {
		return GroupingElement{}, err
	}
	first, err := p.parseGroupingSetExprs()
	if err != nil {
		return GroupingElement{}, err
	}
	sets := [][]Expr{first}
	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		next, err := p.parseGroupingSetExprs()
		if err != nil {
			return GroupingElement{}, err
		}
		sets = append(sets, next)
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return GroupingElement{}, err
	}
	return GroupingElement{
		Kind: GroupingSets,
		Sets: sets,
		Loc:  ast.Loc{Start: groupingTok.Loc.Start, End: closeTok.Loc.End},
	}, nil
}

// parseGroupingSetExprs parses one grouping set inside GROUPING SETS: either a
// parenthesized expression list `( e, … )` (possibly empty, the `()` empty set)
// or a single bare expression. Returns the expression list (nil/empty for `()`).
func (p *Parser) parseGroupingSetExprs() ([]Expr, error) {
	if p.cur.Kind == int('(') {
		p.advance() // consume '('
		exprs, _, err := p.parseBracketedExprList(int(')'))
		return exprs, err
	}
	e, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	return []Expr{e}, nil
}

// parseGroupingSetOrExpr parses a single groupingSet (the singleGroupingSet
// alternative): a parenthesized grouping set `( (expression (, …)*)? )` — which
// may be empty (the `()` empty grouping set) or multi-column — or a single bare
// expression.
//
// A leading '(' is genuinely ambiguous (oracle-confirmed in Trino 481): it can
// open a grouping set (`GROUP BY (a, b)`, `GROUP BY ()`) OR a parenthesized
// sub-expression of a larger grouping expression (`GROUP BY (a + b) * 2`,
// `GROUP BY (a) IS NULL`, `GROUP BY (a, b)[1]` — all parse, the latter as a row
// subscript). ANTLR resolves this by adaptive prediction; here a checkpoint does
// it: read `( expr-list? )`, then decide by what follows the `)`:
//   - empty `()` → the empty grouping set (nothing else can start with `()`).
//   - the `)` followed by an expression-continuation token (an operator,
//     subscript, dereference, predicate, …) → the parentheses were a
//     sub-expression; rewind and parse the whole grouping element as one
//     expression.
//   - otherwise (the `)` is followed by ',' or a clause terminator) → a
//     parenthesized grouping set; keep its column list.
//
// A leading non-'(' token is always a single grouping expression.
func (p *Parser) parseGroupingSetOrExpr() (GroupingElement, error) {
	if p.cur.Kind == int('(') {
		cp := p.checkpoint()
		openTok := p.advance() // consume '('
		exprs, closeTok, err := p.parseBracketedExprList(int(')'))
		if err == nil && (len(exprs) == 0 || !continuesExpression(p.cur.Kind)) {
			return GroupingElement{
				Kind:  GroupingExprKind,
				Exprs: exprs,
				Loc:   ast.Loc{Start: openTok.Loc.Start, End: closeTok.Loc.End},
			}, nil
		}
		// The parentheses opened a larger expression (or the bracketed list failed
		// to parse, e.g. a single expression with internal structure). Rewind and
		// parse the whole grouping element as one expression.
		p.restore(cp)
	}
	e, err := p.parseExpr()
	if err != nil {
		return GroupingElement{}, err
	}
	return GroupingElement{
		Kind:  GroupingExprKind,
		Exprs: []Expr{e},
		Loc:   e.Span(),
	}, nil
}

// continuesExpression reports whether kind, appearing right after a closing ')',
// continues the surrounding expression — i.e. it is an infix/postfix operator,
// subscript, dereference, AT-TIME-ZONE, comparison, or boolean/predicate keyword.
// Used by parseGroupingSetOrExpr to tell a parenthesized grouping set from a
// parenthesized sub-expression. The conservative default (false) treats the ')'
// as ending a grouping set, which is correct at a grouping-element boundary
// (',' or a clause keyword).
func continuesExpression(kind TokenKind) bool {
	switch kind {
	case int('+'), int('-'), int('*'), int('/'), int('%'), // arithmetic
		tokConcat,                                                       // ||
		int('['),                                                        // subscript
		int('.'),                                                        // dereference
		kwAT,                                                            // AT TIME ZONE
		int('='), tokNotEq, int('<'), tokLessEq, int('>'), tokGreaterEq, // comparison
		kwAND, kwOR, // boolean
		kwNOT, kwBETWEEN, kwIN, kwLIKE, kwIS: // predicate
		return true
	default:
		return false
	}
}
