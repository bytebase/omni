package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseSelectStmt parses a SELECT statement.
//
// BNF: oracle/parser/bnf/SELECT.bnf
//
//	select:
//	    [ with_clause ] subquery [ for_update_clause ] ;
//
//	subquery:
//	    query_block
//	    [ { UNION [ ALL ] | INTERSECT | MINUS } query_block ]...
//	    [ order_by_clause ]
//	    [ row_limiting_clause ]
//
//	query_block:
//	    SELECT [ hint ] [ DISTINCT | UNIQUE | ALL ]
//	    select_list
//	    [ FROM { table_reference | join_clause } [, { table_reference | join_clause } ]... ]
//	    [ where_clause ]
//	    [ hierarchical_query_clause ]
//	    [ group_by_clause ]
//	    [ model_clause ]
//	    [ window_clause ]
//	    [ qualify_clause ]
//
//	order_by_clause:
//	    ORDER [ SIBLINGS ] BY
//	    { expr | position | c_alias } [ ASC | DESC ] [ NULLS FIRST | NULLS LAST ]
//	    [, { expr | position | c_alias } [ ASC | DESC ] [ NULLS FIRST | NULLS LAST ] ]...
//
//	row_limiting_clause:
//	    [ OFFSET offset { ROW | ROWS } ]
//	    [ FETCH { FIRST | NEXT } [ { rowcount | percent PERCENT } ] { ROW | ROWS }
//	      { ONLY | WITH TIES } ]
//
//	for_update_clause:
//	    FOR UPDATE
//	    [ OF [ [ schema. ] { table | view } . ] column
//	      [, [ [ schema. ] { table | view } . ] column ]... ]
//	    [ NOWAIT | WAIT integer | SKIP LOCKED ]
//
//	window_clause:
//	    WINDOW window_name AS ( window_specification )
//	    [, window_name AS ( window_specification ) ]...
//
//	qualify_clause:
//	    QUALIFY condition
func (p *Parser) parseSelectStmt() (*nodes.SelectStmt, error) {
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()
	defer p.restorePLSQLCode(p.setPLSQLCode(false))
	// A query's own FROM is no DML statement's.
	defer func(saved bool) { p.dmlFromClause = saved }(p.dmlFromClause)
	p.dmlFromClause = false

	sel, more, err := p.parseQueryOperand()
	if err != nil || !more {
		return sel, err
	}
	return p.parseSetOperands(sel)
}

// parseQueryOperand parses a query up to a set operator after it: a query
// block, or a parenthesized subquery with what follows its parentheses,
// either after a WITH clause. more reports that a set operator ended it, its
// right operand still to parse (parseSetOperands).
func (p *Parser) parseQueryOperand() (sel *nodes.SelectStmt, more bool, err error) {
	defer p.restorePLSQLCode(p.setPLSQLCode(false))

	start := p.pos()
	sel = &nodes.SelectStmt{
		TargetList: &nodes.List{},
		Loc:        nodes.Loc{Start: start},
	}

	// WITH clause
	if p.cur.Type == kwWITH {
		var parseErr941 error
		sel.WithClause, parseErr941 = p.parseWithClause()
		if parseErr941 != nil {
			return nil, false, parseErr941
		}
	}

	if p.cur.Type == '(' {
		return p.parseParenthesizedQuery(sel)
	}
	if p.cur.Type != kwSELECT {
		// A WITH clause is followed by its query (ORA-00928).
		if sel.WithClause != nil {
			return nil, false, p.syntaxErrorAtCur()
		}
		sel.Loc.End = p.prev.End
		return sel, false, nil
	}
	p.advance() // consume SELECT
	return p.parseQueryBlock(sel)
}

// parseSetOperands parses the right operands of the set operation that
// first starts, its operator read: each operand is a query that the next
// operator may end, so the operands chain right-deep through Rarg, and
// every query of the chain ends where the last one does. The chain is read
// in a loop: it is long, not nested, and Oracle 23ai runs one of 10,001
// UNION ALL arms, so its operands do not count toward the nesting limit or
// grow the stack.
func (p *Parser) parseSetOperands(first *nodes.SelectStmt) (*nodes.SelectStmt, error) {
	open := []*nodes.SelectStmt{first}
	for sel := first; ; {
		next, more, err := p.parseQueryOperand()
		if err != nil {
			return nil, err
		}
		sel.Rarg = next
		if !more {
			break
		}
		open = append(open, next)
		sel = next
	}
	for _, sel := range open {
		sel.Loc.End = p.prev.End
	}
	// The set operation locks no rows: Oracle 23ai rejects a FOR UPDATE
	// after any operator with ORA-01786 at FOR, as for SELECT a FROM t
	// EXCEPT SELECT a FROM t FOR UPDATE. Only the last operand can hold one:
	// no operator follows a FOR UPDATE (parseQueryTail).
	if lock := trailingForUpdate(first.Rarg); lock != nil {
		return nil, p.syntaxErrorAtNode(lock)
	}
	return first, nil
}

// parseParenthesizedQuery parses a query whose first operand is a
// parenthesized subquery, the current token at '(', into sel, which holds
// any WITH clause before it, up to a set operator after it (more, as for
// parseQueryOperand).
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/SELECT.html
//
//	( subquery ) [ order_by_clause ] [ for_update_clause ] [ row_limiting_clause ]
//	| ( subquery ) { UNION [ ALL ] | INTERSECT [ ALL ] | MINUS [ ALL ] | EXCEPT [ ALL ] } subquery
//
// Oracle 23ai rejects a WITH clause inside the parentheses (ORA-32034), an
// ORDER BY that ends the parenthesized subquery with no row-limiting clause
// after it (ORA-00907, for ... ORDER BY 1 FOR UPDATE too), and a FOR UPDATE
// that ends it (ORA-00907 for (SELECT a FROM t FOR UPDATE) UNION ...),
// wherever the parenthesized subquery stands; an inline view in FROM is a table
// reference, which takes both. Without anything around the parentheses the
// query is the subquery itself; otherwise sel is a wrapper whose Larg is
// the subquery.
func (p *Parser) parseParenthesizedQuery(sel *nodes.SelectStmt) (*nodes.SelectStmt, bool, error) {
	open := p.pos()
	p.advance() // consume '('
	if p.cur.Type != kwSELECT && p.cur.Type != '(' {
		return nil, false, p.syntaxErrorAtCur()
	}
	inner, err := p.parseSelectStmt()
	if err != nil {
		return nil, false, err
	}
	if sort := trailingOrderBy(inner); sort != nil {
		return nil, false, p.syntaxErrorAtNode(sort)
	}
	if lock := trailingForUpdate(inner); lock != nil {
		return nil, false, p.syntaxErrorAtNode(lock)
	}
	if p.cur.Type != ')' {
		return nil, false, p.syntaxErrorAtCur()
	}
	p.advance() // consume ')'

	more, err := p.parseQueryTail(sel)
	if err != nil {
		return nil, false, err
	}
	if more {
		// A set operation: sel wraps the subquery, and ends where its
		// right operand does (parseSetOperands).
		sel.Larg = inner
		return sel, true, nil
	}
	// A set operation locks no rows, parenthesized whole too (ORA-01786
	// at FOR for (SELECT a FROM t UNION SELECT a FROM t) FOR UPDATE), nor
	// does a query unlockable refuses.
	if sel.ForUpdate != nil && (isSetQuery(inner) || unlockable(inner)) {
		return nil, false, p.syntaxErrorAtNode(sel.ForUpdate)
	}
	if sel.WithClause == nil && sel.OrderBy == nil && sel.ForUpdate == nil && sel.FetchFirst == nil && sel.Op == nodes.SETOP_NONE {
		// The query spans its parentheses, so a consumer that edits or
		// quotes it by position keeps them.
		inner.Loc = nodes.Loc{Start: open, End: p.prev.End}
		return inner, false, nil
	}
	sel.Larg = inner
	sel.Loc.End = p.prev.End
	return sel, false, nil
}

// parenQuery says how a parenthesized group reads as a query: whether the
// first token after its '(' and any '(' that follow is SELECT or WITH
// (leads), and whether a set operator, ORDER BY, OFFSET, or FETCH follows
// its ')' (continues).
type parenQuery struct {
	leads, continues bool
}

// decideParenQueries reads the tokens of lexer, a lookahead under the
// selection of conditional compilation branches the parser attempts, and
// reports, for each '(' whose group leads to a query, whether a query
// continues after the group. Matching parentheses with a stack and finding
// the first token after a run of '(' right to left visit each token a
// constant number of times.
func decideParenQueries(lexer *Lexer) map[int]parenQuery {
	var toks []Token
	for {
		tok := lexer.NextToken()
		if tok.Type == tokEOF || lexer.Err != nil {
			break
		}
		toks = append(toks, tok)
	}
	n := len(toks)

	// closing[i] is the index of the ')' that closes the '(' at i, or n.
	closing := make([]int, n)
	var open []int
	for i, tok := range toks {
		switch tok.Type {
		case '(':
			closing[i] = n
			open = append(open, i)
		case ')':
			if len(open) > 0 {
				closing[open[len(open)-1]] = i
				open = open[:len(open)-1]
			}
		}
	}

	// first[i] is the index of the first token from i on that is not '('.
	first := make([]int, n+1)
	first[n] = n
	for i := n - 1; i >= 0; i-- {
		first[i] = i
		if toks[i].Type == '(' {
			first[i] = first[i+1]
		}
	}

	queries := make(map[int]parenQuery)
	for i, tok := range toks {
		if tok.Type != '(' {
			continue
		}
		f := first[i+1]
		if f == n || (toks[f].Type != kwSELECT && toks[f].Type != kwWITH) {
			continue
		}
		q := parenQuery{leads: true}
		if c := closing[i]; c+1 < n {
			switch toks[c+1].Type {
			case kwUNION, kwINTERSECT, kwMINUS, kwEXCEPT, kwORDER, kwOFFSET, kwFETCH:
				q.continues = true
			case kwFOR:
				// FOR UPDATE, of a cursor's query: FOR r IN ((SELECT a
				// FROM t) FOR UPDATE) LOOP.
				q.continues = c+2 < n && toks[c+2].Type == kwUPDATE
			}
		}
		queries[tok.Loc] = q
	}
	return queries
}

// parenQueryAt returns how the group of the '(' at offset loc reads as a
// query.
func (p *Parser) parenQueryAt(loc int) parenQuery {
	if p.parenQueries == nil {
		p.parenQueries = decideParenQueries(p.selectionLexer(p.cur.Loc))
	}
	return p.parenQueries[loc]
}

// atQueryStart reports whether the current token starts a query where an
// expression or a table reference may stand as well: SELECT, WITH, or a
// parenthesized subquery that a set operator, ORDER BY, OFFSET, or FETCH
// follows, as in ((SELECT ...) UNION (SELECT ...)). A parenthesized
// subquery alone stays an expression, ((SELECT ...)) + 1, or a table
// reference, ((SELECT ...) v JOIN t ON ...).
func (p *Parser) atQueryStart() bool {
	switch p.cur.Type {
	case kwSELECT, kwWITH:
		return true
	case '(':
		return p.parenQueryAt(p.cur.Loc).continues
	}
	return false
}

// atQuery reports whether the current token starts a query where only a
// query may stand: SELECT, WITH, or a parenthesized subquery.
func (p *Parser) atQuery() bool {
	switch p.cur.Type {
	case kwSELECT, kwWITH:
		return true
	case '(':
		return p.parenQueryAt(p.cur.Loc).leads
	}
	return false
}

// isSetQuery reports whether query is a set operation, possibly wrapped in
// parentheses.
func isSetQuery(query *nodes.SelectStmt) bool {
	for ; query != nil; query = query.Larg {
		if query.Op != nodes.SETOP_NONE {
			return true
		}
	}
	return false
}

// unlockable reports whether FOR UPDATE cannot lock the rows of query, or
// of a parenthesized query it wraps. Oracle 23ai rejects it over DISTINCT,
// UNIQUE, GROUP BY, or a CURSOR expression in the select list (ORA-01786),
// and over any row limit but FETCH FIRST n
// ROWS ONLY: an OFFSET, a PERCENT, WITH TIES, or an ORDER BY with the row
// limit (ORA-02014). An aggregate's name is the catalog's to tell, and
// Oracle's to reject.
func unlockable(query *nodes.SelectStmt) bool {
	for ; query != nil; query = query.Larg {
		if query.Distinct || query.UniqueKw || query.GroupClause != nil && query.GroupClause.Len() > 0 || selectsCursor(query.TargetList) {
			return true
		}
		if fetch := query.FetchFirst; fetch != nil && (fetch.Offset != nil || fetch.Percent || fetch.WithTies || query.OrderBy != nil) {
			return true
		}
		if query.Op != nodes.SETOP_NONE {
			return false
		}
	}
	return false
}

// selectsCursor reports whether a CURSOR expression stands in list, a
// select list, outside the subqueries it holds.
func selectsCursor(list *nodes.List) bool {
	found := false
	if list != nil {
		inspect(list, func(n nodes.Node) bool {
			switch n.(type) {
			case *nodes.CursorExpr:
				found = true
			case *nodes.SelectStmt:
				return false
			}
			return !found
		})
	}
	return found
}

// trailingForUpdate returns the FOR UPDATE clause that ends query, its own
// or that of the last operand of its set operation, or nil.
func trailingForUpdate(query *nodes.SelectStmt) nodes.Node {
	for query != nil {
		if query.ForUpdate != nil {
			return query.ForUpdate
		}
		if query.Op == nodes.SETOP_NONE {
			return nil
		}
		query = query.Rarg
	}
	return nil
}

// trailingOrderBy returns the first item of the ORDER BY that ends query, if
// one does: the query's own, or that of the last operand of its set
// operation.
func trailingOrderBy(query *nodes.SelectStmt) nodes.Node {
	for query != nil {
		if query.OrderBy != nil && query.OrderBy.Len() > 0 {
			// One a row-limiting clause follows is that clause's: Oracle
			// 23ai runs (SELECT ... ORDER BY d FETCH FIRST 10 ROWS ONLY)
			// UNION ALL ..., with OFFSET too.
			if query.FetchFirst != nil {
				return nil
			}
			return query.OrderBy.Items[0]
		}
		if query.Op == nodes.SETOP_NONE {
			return nil
		}
		query = query.Rarg
	}
	return nil
}

// parseSelectAfterKeyword parses the rest of a query into sel, from the
// token after the SELECT keyword: hints, the select list, the clauses that
// follow, and any set operation after them. CREATE JAVA ... USING calls it
// directly for the legacy subquery that omits the SELECT keyword.
func (p *Parser) parseSelectAfterKeyword(sel *nodes.SelectStmt) (*nodes.SelectStmt, error) {
	sel, more, err := p.parseQueryBlock(sel)
	if err != nil || !more {
		return sel, err
	}
	return p.parseSetOperands(sel)
}

// parseQueryBlock parses the rest of a query block into sel, from the token
// after the SELECT keyword, up to a set operator after it (more, as for
// parseQueryOperand).
func (p *Parser) parseQueryBlock(sel *nodes.SelectStmt) (*nodes.SelectStmt, bool, error) {
	if res, err := p.parseQueryBlockClauses(sel); err != nil {
		return res, false, err
	}
	more, err := p.parseQueryTail(sel)
	if err != nil {
		return nil, false, err
	}
	if !more {
		sel.Loc.End = p.prev.End
	}
	return sel, more, nil
}

// parseQueryBlockClauses parses the clauses of a query block into sel, from
// the token after the SELECT keyword: hints, the select list, and the
// clauses that follow, up to its ORDER BY.
func (p *Parser) parseQueryBlockClauses(sel *nodes.SelectStmt) (*nodes.SelectStmt, error) {
	// Hints
	if p.cur.Type == tokHINT {
		tok := p.cur
		sel.Hints = &nodes.List{}
		sel.Hints.Items = append(sel.Hints.Items, &nodes.Hint{
			Text: tok.Str,
			Loc:  nodes.Loc{Start: tok.Loc, End: tok.End},
		})
		p.advance()
	}

	// ALL | DISTINCT | UNIQUE
	switch p.cur.Type {
	case kwALL:
		sel.All = true
		p.advance()
	case kwDISTINCT:
		sel.Distinct = true
		p.advance()
	case kwUNIQUE:
		sel.UniqueKw = true
		sel.Distinct = true
		p.advance()
	}
	var parseErr942 error

	// Select list
	sel.TargetList, parseErr942 = p.parseSelectList()
	if parseErr942 != nil {
		return nil, parseErr942
	}
	if sel.TargetList.Len() == 0 {

		return sel, p.syntaxErrorAtCur(

		// INTO (PL/SQL SELECT ... INTO variable_list FROM ...)
		)
	}

	// BULK COLLECT INTO (PL/SQL SELECT ... BULK COLLECT INTO collection_list FROM ...)
	bulk, err := p.parseOptionalBulkCollect()
	if err != nil {
		return nil, err
	}
	sel.BulkCollect = bulk

	if p.cur.Type == kwINTO {
		p.advance()
		var // consume INTO
		parseErr943 error
		sel.IntoVars, parseErr943 = p.parseExprList()
		if parseErr943 !=

			// FROM
			nil {
			return nil, parseErr943
		}
		if err := p.checkTargets(sel.IntoVars); err != nil {
			return nil, err
		}
	}

	if p.cur.Type == kwFROM {
		p.advance()
		var parseErr944 error
		sel.FromClause, parseErr944 = p.parseFromClause()
		if parseErr944 !=

			// PIVOT / UNPIVOT (parsed after FROM, before WHERE)
			nil {
			return nil, parseErr944
		}
	}

	// PIVOT / UNPIVOT are parsed inside the table_reference chain
	// (parsePivotSuffix); nothing to do at the select level.

	if p.cur.Type == kwWHERE {
		p.advance()
		var parseErr947 error
		sel.WhereClause, parseErr947 = p.parseExpr()
		if parseErr947 != nil {
			return nil, parseErr947
		}
		if sel.WhereClause == nil {

			return sel, p.syntaxErrorAtCur(

			// START WITH / CONNECT BY (either order)
			)
		}
	}

	if p.cur.Type == kwSTART || p.cur.Type == kwCONNECT {
		var parseErr948 error
		sel.Hierarchical, parseErr948 = p.parseHierarchicalClause()
		if parseErr948 !=

			// GROUP BY
			nil {
			return nil, parseErr948
		}
	}

	if p.cur.Type == kwGROUP {
		p.advance()
		if p.cur.Type != kwBY {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		var parseErr949 error
		sel.GroupClause, parseErr949 = p.parseGroupByList()
		if parseErr949 != nil {
			return nil, parseErr949
		}
		if sel.GroupClause.Len() == 0 {

			return sel, p.syntaxErrorAtCur(

			// HAVING
			)
		}
	}

	if p.cur.Type == kwHAVING {
		p.advance()
		var parseErr950 error
		sel.HavingClause, parseErr950 = p.parseExpr()
		if parseErr950 != nil {
			return nil, parseErr950
		}
		if sel.HavingClause == nil {

			return sel, p.syntaxErrorAtCur(

			// MODEL clause
			)
		}
	}

	if p.cur.Type == kwGROUP && (sel.GroupClause == nil || sel.GroupClause.Len() == 0) {
		p.advance()
		if p.cur.Type != kwBY {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		var parseErr951 error
		sel.GroupClause, parseErr951 = p.parseGroupByList()
		if parseErr951 != nil {
			return nil, parseErr951
		}
		if sel.GroupClause.Len() == 0 {
			return sel, p.syntaxErrorAtCur()
		}
	}

	// A MODEL clause reads the rows of a FROM clause (ORA-00923 at MODEL
	// without one).
	if p.cur.Type == kwMODEL && (sel.FromClause == nil || sel.FromClause.Len() == 0) {
		return nil, p.syntaxErrorAtCur()
	}
	if p.cur.Type == kwMODEL {
		var parseErr951 error
		sel.ModelClause, parseErr951 = p.parseModelClause()
		if parseErr951 !=

			// WINDOW clause
			nil {
			return nil, parseErr951
		}
	}

	if p.isIdentLikeStr("WINDOW") {
		var parseErr952 error
		sel.WindowDefs, parseErr952 = p.parseWindowClause()
		if parseErr952 !=

			// QUALIFY clause
			nil {
			return nil, parseErr952
		}
	}

	if p.isIdentLikeStr("QUALIFY") {
		p.advance()
		var parseErr953 error
		sel.QualifyClause, parseErr953 = p.parseExpr()
		if parseErr953 != nil {
			return nil, parseErr953
		}
		if sel.QualifyClause == nil {
			return nil, p.syntaxErrorAtCur()
		}
	}
	return sel, nil
}

// parseQueryTail parses what may follow a query block or a parenthesized
// subquery into sel: the ORDER BY, row limiting, and FOR UPDATE clauses, in
// that order (ORA-03049 on Oracle 23ai for FOR UPDATE FETCH FIRST ...),
// then a set operator. more reports that one was read, a query standing
// after it for parseSetOperands to parse.
func (p *Parser) parseQueryTail(sel *nodes.SelectStmt) (more bool, err error) {
	if p.cur.Type == kwORDER {
		p.advance()
		if p.isIdentLikeStr("SIBLINGS") {
			// SIBLINGS orders a hierarchical query's rows (ORA-30929).
			if sel.Hierarchical == nil {
				return false, p.syntaxErrorAtCur()
			}
			sel.SiblingsOrder = true
			p.advance()
		}
		if p.cur.Type != kwBY {
			return false, p.syntaxErrorAtCur()
		}
		p.advance()
		sel.OrderBy, err = p.parseOrderByList()
		if err != nil {
			return false, err
		}
		if sel.OrderBy.Len() == 0 {
			return false, p.syntaxErrorAtCur()
		}
	}

	if p.cur.Type == kwOFFSET || p.cur.Type == kwFETCH {
		sel.FetchFirst, err = p.parseFetchFirstClause()
		if err != nil {
			return false, err
		}
	}

	if p.cur.Type == kwFOR {
		sel.ForUpdate, err = p.parseForUpdateClause()
		if err != nil {
			return false, err
		}
		if unlockable(sel) {
			return false, p.syntaxErrorAtNode(sel.ForUpdate)
		}
	}

	// Set operations: UNION, INTERSECT, MINUS, EXCEPT, each with ALL in
	// 23ai. None follows an ORDER BY, FOR UPDATE, or row limiting clause
	// (ORA-03048 on Oracle 23ai).
	switch p.cur.Type {
	case kwUNION, kwINTERSECT, kwMINUS, kwEXCEPT:
		if sel.OrderBy != nil || sel.ForUpdate != nil || sel.FetchFirst != nil {
			return false, p.syntaxErrorAtCur()
		}
	}
	switch p.cur.Type {
	case kwUNION:
		sel.Op = nodes.SETOP_UNION
	case kwINTERSECT:
		sel.Op = nodes.SETOP_INTERSECT
	case kwMINUS:
		sel.Op = nodes.SETOP_MINUS
	case kwEXCEPT:
		sel.Op = nodes.SETOP_EXCEPT
	default:
		return false, nil
	}
	p.advance()
	if p.cur.Type == kwALL {
		sel.SetAll = true
		p.advance()
	}
	// A query follows the operator (ORA-00928 on Oracle 23ai for UNION ),
	// UNION (1), or UNION at the end), without a WITH clause of its own
	// (ORA-32034 for UNION WITH w AS (...) SELECT ...).
	if p.cur.Type == kwWITH || !p.atQuery() {
		return false, p.syntaxErrorAtCur()
	}
	return true, nil
}

// parseSelectList parses the select list (target expressions).
func (p *Parser) parseSelectList() (*nodes.List, error) {
	list := &nodes.List{}

	for {
		if p.isSelectListTerminator() {
			break
		}
		// A bare * is the whole select list: Oracle 23ai fails at the
		// ',' after it (ORA-00923) and at a * after a ',' (ORA-00936).
		if p.cur.Type == '*' && list.Len() > 0 {
			return nil, p.syntaxErrorAtCur()
		}
		rt, parseErr960 := p.parseResTarget()
		if parseErr960 != nil {
			return nil, parseErr960
		}
		// Every item is an expression (ORA-00936 at the ',' of 1,,2).
		if rt == nil {
			return nil, p.syntaxErrorAtCur()
		}
		list.Items = append(list.Items, rt)
		if p.cur.Type != ',' {
			break
		}
		if _, star := rt.Expr.(*nodes.Star); star {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance() // consume ','
	}

	return list, nil
}

// isSelectListTerminator reports whether the current token ends the select
// list where an item would start. FETCH and OFFSET are absent: an item can be
// a column of that name, and the row_limiting_clause never follows an empty
// list.
func (p *Parser) isSelectListTerminator() bool {
	switch p.cur.Type {
	case tokEOF, kwFROM, kwWHERE, kwGROUP, kwHAVING, kwORDER, kwUNION,
		kwINTERSECT, kwMINUS, kwEXCEPT, kwFOR, kwCONNECT, kwSTART:
		return true
	default:
		return false
	}
}

// parseResTarget parses a single target expression (with optional alias).
func (p *Parser) parseResTarget() (*nodes.ResTarget, error) {
	start := p.pos()
	if p.cur.Type == '*' {
		return &nodes.ResTarget{Expr: p.takeSelectStar(), Loc: nodes.Loc{Start: start, End: p.prev.End}}, nil
	}
	expr, parseErr961 := p.parseExpr()
	if parseErr961 != nil {
		return nil, parseErr961
	}
	if expr == nil {
		return nil, nil
	}

	rt := &nodes.ResTarget{
		Expr: expr,
		Loc:  nodes.Loc{Start: start},
	}

	// * and t.* take no alias (ORA-00923 at it).
	if p.isWildcard(expr) {
		rt.Loc.End = p.prev.End
		return rt, nil
	}

	// Optional alias: AS name or just name
	if p.cur.Type == kwAS {
		p.advance()
		var parseErr962 error
		rt.Name, parseErr962 = p.parseIdentifier()
		if parseErr962 != nil {
			return nil, parseErr962
		}
	} else if p.cur.Type == tokBIND {
		// Oracle takes a bind variable after an expression as its alias
		// (compat.tsv's suite:lexer row SELECT :a :b FROM dual).
		rt.Name = p.cur.Str
		rt.AliasBind = true
		p.advance()
	} else if p.isAliasCandidate() {
		var parseErr963 error
		rt.Name, parseErr963 = p.parseIdentifier()
		if parseErr963 != nil {
			return nil, parseErr963
		}
	}

	rt.Loc.End = p.prev.End
	return rt, nil
}

// isAliasCandidate returns true if the current token can be an implicit alias.
// We exclude keywords that start clauses to avoid consuming FROM, WHERE, etc. as aliases.
func (p *Parser) isAliasCandidate() bool {
	if p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT {
		return true
	}
	// BULK is a valid alias, but BULK COLLECT starts the PL/SQL
	// bulk_collect_into_clause after the select list.
	if p.isBulkCollect() {
		return false
	}
	// Disallow clause-starting keywords as implicit aliases
	switch p.cur.Type {
	case kwFROM, kwWHERE, kwGROUP, kwHAVING, kwORDER, kwUNION, kwINTERSECT,
		kwMINUS, kwEXCEPT, kwFOR, kwCONNECT, kwSTART, kwFETCH, kwOFFSET, kwON,
		kwLEFT, kwRIGHT, kwINNER, kwOUTER, kwCROSS, kwFULL, kwNATURAL,
		kwJOIN, kwWHEN, kwTHEN, kwELSE, kwEND, kwAND, kwOR, kwNOT,
		kwIS, kwIN, kwBETWEEN, kwLIKE, kwLIKEC, kwLIKE2, kwLIKE4,
		kwINTO, kwVALUES, kwSET, kwRETURNING, kwPIVOT, kwUNPIVOT,
		kwMODEL, kwWITH, kwKEEP, kwOVER, kwUSING:
		return false
	}
	return p.cur.Type >= 2000 &&
		!isOracleSQLReservedKeyword(p.cur) &&
		!isOracleClauseStarterKeyword(p.cur.Type)
}

// isBulkCollect reports whether the current tokens are BULK COLLECT, the
// start of a PL/SQL bulk_collect_into_clause.
func (p *Parser) isBulkCollect() bool {
	return p.cur.Type == kwBULK && p.peekNext().Type == kwCOLLECT
}

// parseExprList parses a comma-separated list of expressions.
func (p *Parser) parseExprList() (*nodes.List, error) {
	list := &nodes.List{}
	for {
		expr, parseErr964 := p.parseExpr()
		if parseErr964 != nil {
			return nil, parseErr964
		}
		if expr == nil {
			return nil, p.syntaxErrorAtCur()

		}
		list.Items = append(list.Items, expr)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
		if p.cur.Type == tokEOF {
			return nil, p.syntaxErrorAtCur()

		}
	}
	return list, nil
}

// parseGroupByList parses a comma-separated GROUP BY list, handling
// GROUPING SETS, CUBE, and ROLLUP extensions.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/SELECT.html
//
//	group_by_clause ::=
//	    GROUP BY { expr | rollup_cube | grouping_sets } [, ...]
//	rollup_cube ::= { ROLLUP | CUBE } ( expr [, ...] )
//	grouping_sets ::= GROUPING SETS ( { rollup_cube | expr } [, ...] )
func (p *Parser) parseGroupByList() (*nodes.List, error) {
	list := &nodes.List{}
	for {
		item, parseErr965 := p.parseGroupByItem()
		if parseErr965 != nil {
			return nil, parseErr965
		}
		if item == nil {
			return nil, p.syntaxErrorAtCur()

		}
		list.Items = append(list.Items, item)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
		if p.cur.Type == tokEOF {
			return nil, p.syntaxErrorAtCur()

		}
	}
	return list, nil
}

// parseGroupByItem parses a single GROUP BY item: an expression, ROLLUP(...),
// CUBE(...), or GROUPING SETS(...).
//
//	rollup_cube_clause:
//	    { ROLLUP | CUBE } ( grouping_expression_list )
//
//	grouping_sets_clause:
//	    GROUPING SETS
//	    ( { rollup_cube_clause | grouping_expression_list }
//	      [, { rollup_cube_clause | grouping_expression_list } ]... )
//
// An empty ( ) is a grouping set, and a whole GROUP BY list, but not an
// element of ROLLUP or CUBE (ORA-00936 at its ')'). GROUPING SETS do not
// nest (ORA-00907 at the inner SETS), and the GROUPING function groups
// nothing (ORA-00934).
func (p *Parser) parseGroupByItem() (nodes.Node, error) {
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()
	start := p.pos()

	switch {
	case (p.cur.Type == kwROLLUP || p.cur.Type == kwCUBE) && p.peekNext().Type == '(':
		rollup := p.cur.Type == kwROLLUP
		p.advance() // consume ROLLUP or CUBE
		args, err := p.parseSelectGroupingList()
		if err != nil {
			return nil, err
		}
		if rollup {
			return &nodes.RollupClause{Args: args, Loc: nodes.Loc{Start: start, End: p.prev.End}}, nil
		}
		return &nodes.CubeClause{Args: args, Loc: nodes.Loc{Start: start, End: p.prev.End}}, nil

	case p.cur.Type == kwGROUPING && p.peekNext().Type == kwSETS:
		p.advance() // consume GROUPING
		p.advance() // consume SETS
		gs := &nodes.GroupingSetsClause{Loc: nodes.Loc{Start: start}}
		if p.cur.Type != '(' {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		gs.Sets = &nodes.List{}
		for {
			var item nodes.Node
			switch {
			case p.cur.Type == kwGROUPING && p.peekNext().Type == kwSETS:
				p.advance()
				return nil, p.syntaxErrorAtCur()
			case p.cur.Type == '(' && p.peekNext().Type == ')':
				parenStart := p.pos()
				p.advance()
				p.advance()
				item = &nodes.ParenExpr{Loc: nodes.Loc{Start: parenStart, End: p.prev.End}}
			default:
				var err error
				if item, err = p.parseGroupByItem(); err != nil {
					return nil, err
				}
				if item == nil {
					return nil, p.syntaxErrorAtCur()
				}
			}
			gs.Sets.Items = append(gs.Sets.Items, item)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if p.cur.Type != ')' {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		gs.Loc.End = p.prev.End
		return gs, nil

	case p.cur.Type == kwGROUPING && p.peekNext().Type == '(':
		return nil, p.syntaxErrorAtCur()
	}
	return p.parseExpr()
}

// parseSelectGroupingList parses the parenthesized grouping expression list
// of ROLLUP or CUBE: expressions, each maybe a parenthesized list, never an
// empty one.
func (p *Parser) parseSelectGroupingList() (*nodes.List, error) {
	if p.cur.Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	args := &nodes.List{}
	for {
		if p.cur.Type == '(' && p.peekNext().Type == ')' {
			p.advance()
			return nil, p.syntaxErrorAtCur()
		}
		e, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if e == nil {
			return nil, p.syntaxErrorAtCur()
		}
		args.Items = append(args.Items, e)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	return args, nil
}

// parseFromClause parses a FROM clause (comma-separated table references).
func (p *Parser) parseFromClause() (*nodes.List, error) {
	list := &nodes.List{}

	for {
		tref, parseErr969 := p.parseTableRef()
		if parseErr969 != nil {
			return nil, parseErr969

			// Check for JOINs
		}
		if tref == nil {
			break
		}
		var parseErr970 error

		tref, parseErr970 = p.parseJoinContinuation(tref)
		if parseErr970 != nil {
			return nil, parseErr970
		}

		list.Items = append(list.Items, tref)

		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	return list, nil
}

// parseTableRef parses a single table reference.
//
// BNF: oracle/parser/bnf/SELECT.bnf
//
//	table_reference:
//	    { query_table_expression | ( join_clause ) | inline_external_table } [ t_alias ]
//
//	query_table_expression:
//	    { [ schema. ] { table | view | materialized_view | analytic_view | hierarchy }
//	      [ partition_extension_clause ] [ @ dblink ]
//	    | ( subquery [ subquery_restriction_clause ] )
//	    | table_collection_expression
//	    | inline_analytic_view
//	    }
//	    [ flashback_query_clause ] [ sample_clause ]
//	    [ pivot_clause | unpivot_clause ] [ row_pattern_clause ]
//	    [ containers_clause ] [ shards_clause ]
func (p *Parser) parseTableRef() (nodes.TableExpr, error) {
	start := p.pos()

	// LATERAL ( subquery )
	if p.cur.Type == kwLATERAL {
		return p.parseLateralRef(start)
	}

	// XMLTABLE(...)
	if p.cur.Type == kwXMLTABLE {
		ref, err := p.parseXmlTableRef(start)
		if err != nil {
			return nil, err
		}
		return p.parsePivotSuffix(ref)
	}

	// JSON_TABLE(...)
	if p.cur.Type == kwJSON_TABLE {
		ref, err := p.parseJsonTableRef(start)
		if err != nil {
			return nil, err
		}
		return p.parsePivotSuffix(ref)
	}

	// TABLE(collection_expression) — table collection expression
	if p.cur.Type == kwTABLE && p.peekNext().Type == '(' {
		ref, err := p.parseTableCollectionExpr(start)
		if err != nil {
			return nil, err
		}
		return p.parsePivotSuffix(ref)
	}

	// EXTERNAL ((columns) TYPE ... DEFAULT DIRECTORY ... LOCATION ...)
	if p.isIdentLikeStr("EXTERNAL") && p.peekNext().Type == '(' {
		return p.parseInlineExternalTable(start)
	}

	// ANALYTIC VIEW ( sub_av_clause ) [ alias ]
	if p.isInlineAnalyticViewAhead() {
		return p.parseInlineAnalyticViewRef(start)
	}

	// CONTAINERS(table) or SHARDS(table)
	if p.isIdentLikeStr("CONTAINERS") && p.peekNext().Type == '(' {
		ref, err := p.parseContainersOrShards(start, false)
		if err != nil {
			return nil, err
		}
		return p.parsePivotSuffix(ref)
	}
	if p.isIdentLikeStr("SHARDS") && p.peekNext().Type == '(' {
		ref, err := p.parseContainersOrShards(start, true)
		if err != nil {
			return nil, err
		}
		return p.parsePivotSuffix(ref)
	}

	// Subquery: ( SELECT ... )
	if p.cur.Type == '(' {
		next := p.peekNext()
		if next.Type == kwSELECT || next.Type == kwWITH || next.Type == '(' && p.parenQueryAt(next.Loc).continues {
			return p.parseSubqueryRef(start)
		}
		return p.parseParenthesizedTableRef(start)
	}

	// MATCH_RECOGNIZE as standalone (rare, usually post-table)
	if p.isIdentLikeStr("MATCH_RECOGNIZE") {
		ref, err := p.parseMatchRecognize(start)
		if err != nil {
			return nil, err
		}
		return p.parsePivotSuffix(ref)
	}

	// Table name
	if !p.isIdentLike() {
		return nil, nil
	}
	if err := p.syntaxErrorIfReservedIdentifier(); err != nil {
		return nil, err
	}

	name, parseErr971 := p.parseReservedCheckedObjectName()
	if parseErr971 != nil {
		return nil, parseErr971
	}
	tr := &nodes.TableRef{
		Name: name,
		Loc:  nodes.Loc{Start: start},
	}

	// A table function called without TABLE (Oracle 23ai): name ( args ).
	if p.cur.Type == '(' && name.DBLink == "" {
		return p.parseSelectTableFunctionRef(name, start)
	}

	// Partition extension clause: PARTITION (name) | PARTITION FOR (key) |
	// SUBPARTITION ..., not the PARTITION BY of a partitioned outer join.
	if (p.cur.Type == kwPARTITION || p.cur.Type == kwSUBPARTITION) && p.peekNext().Type != kwBY {
		var parseErr972 error
		tr.PartitionExt, parseErr972 = p.parsePartitionExtClause(true)
		if parseErr972 !=

			// @ dblink
			nil {
			return nil, parseErr972
		}
	}

	if p.cur.Type == '@' {
		p.advance()
		var parseErr973 error
		tr.Dblink, parseErr973 = p.parseIdentifier()
		if parseErr973 !=

			// Flashback query: VERSIONS BETWEEN ... AND ... or AS OF SCN/TIMESTAMP
			nil {
			return nil, parseErr973
		}
	}

	if p.cur.Type == kwVERSIONS || (p.cur.Type == kwAS && p.peekNext().Type == kwOF) {
		var parseErr974 error
		tr.Flashback, parseErr974 = p.parseFlashbackClause()
		if parseErr974 !=

			// SAMPLE [BLOCK] (percent) [SEED (value)]
			nil {
			return nil, parseErr974
		}
	}

	// SAMPLE opens the clause before ( or BLOCK; otherwise Oracle 23ai reads
	// it as the table's alias (FROM t SAMPLE parses, FROM t SAMPLE 1 is
	// ORA-03047 at the 1).
	// A DML statement's FROM samples no table, nor takes SAMPLE as an alias
	// (ORA-03049 at SAMPLE).
	if p.cur.Type == kwSAMPLE && p.dmlFromClause {
		return nil, p.syntaxErrorAtCur()
	}
	if p.cur.Type == kwSAMPLE && !p.isSelectSampleClauseAhead() {
		alias, err := p.parseAlias()
		if err != nil {
			return nil, err
		}
		tr.Alias = alias
		tr.Loc.End = p.prev.End
		return p.parsePivotSuffix(tr)
	}
	if p.cur.Type == kwSAMPLE {
		var parseErr975 error
		tr.Sample, parseErr975 = p.parseSampleClause()
		if parseErr975 !=

			// MATCH_RECOGNIZE after table reference
			nil {
			return nil, parseErr975
		}
	}

	if p.pivotClauseAhead() || p.unpivotClauseAhead() {
		tr.Loc.End = p.prev.End
		return p.parsePivotSuffix(tr)
	}

	// Nor does it match patterns: Oracle 23ai fails in its own code on
	// MATCH_RECOGNIZE there, so the byte its judge recorded is the engine's
	// crash position and not a token this grammar cannot take (ParseRange).
	if p.isIdentLikeStr("MATCH_RECOGNIZE") && p.dmlFromClause {
		p.engineCrash = p.cur.Loc
		return nil, p.syntaxErrorAtCur()
	}
	if p.isIdentLikeStr("MATCH_RECOGNIZE") {
		// The clause wraps the table it reads, so it starts where the table
		// does, and the table ends before it.
		tr.Loc.End = p.prev.End
		mrRef, parseErr976 := p.parseMatchRecognize(tr.Loc.Start)
		if parseErr976 != nil {
			return nil, parseErr976
		}
		if mrClause, ok := mrRef.(*nodes.MatchRecognizeClause); ok {
			mrClause.Source = tr
			return p.parsePivotSuffix(mrRef)
		}
	}

	// Optional alias
	if p.cur.Type == kwAS {
		// Only consume AS as alias intro if next is not OF (flashback already handled)
		p.advance()
		var parseErr977 error
		tr.Alias, parseErr977 = p.parseAlias()
		if parseErr977 != nil {
			return nil, parseErr977
		}
	} else if p.isTableAliasCandidate() {
		var parseErr978 error
		tr.Alias, parseErr978 = p.parseAlias()
		if parseErr978 != nil {
			return nil, parseErr978
		}
	}

	tr.Loc.End = p.prev.End
	// Oracle also accepts the pivot after an alias (t alias PIVOT (...)),
	// even though the documented grammar orders the alias last.
	return p.parsePivotSuffix(tr)
}

// isTableAliasCandidate checks if current token can be a table alias.
// Excludes soft keywords that start clauses in SELECT statements.
func (p *Parser) isTableAliasCandidate() bool {
	if p.cur.Type == tokQIDENT {
		return true
	}
	// LOG ERRORS opens a DML statement's error_logging_clause.
	if p.cur.Type == kwLOG && p.peekNext().Type == kwERRORS {
		return false
	}
	// PIVOT and UNPIVOT name an alias where no clause of theirs follows:
	// FROM t PIVOT parses on Oracle 23ai.
	if p.cur.Type == kwPIVOT && !p.pivotClauseAhead() || p.cur.Type == kwUNPIVOT && !p.unpivotClauseAhead() {
		return true
	}
	// NATURAL names an alias where no join follows: FROM t NATURAL parses on
	// Oracle 23ai, and so does DELETE a FROM t NATURAL, where the natural
	// join is only the one that JOIN, INNER, LEFT, RIGHT or FULL continues.
	if p.cur.Type == kwNATURAL && !p.naturalJoinAhead() {
		return true
	}
	if p.cur.Type == tokIDENT {
		switch strings.ToUpper(p.cur.Str) {
		case "WINDOW", "QUALIFY", "MATCH_RECOGNIZE", "CONTAINERS", "SHARDS",
			"SIBLINGS", "APPLY", "SEARCH", "CYCLE", "VERSIONS", "PERIOD",
			"XML":
			return false
		}
		return true
	}
	switch p.cur.Type {
	case kwFROM, kwWHERE, kwGROUP, kwHAVING, kwORDER, kwUNION, kwINTERSECT,
		kwMINUS, kwEXCEPT, kwFOR, kwCONNECT, kwSTART, kwFETCH, kwOFFSET, kwON,
		kwLEFT, kwRIGHT, kwINNER, kwOUTER, kwCROSS, kwFULL, kwNATURAL,
		kwJOIN, kwWHEN, kwTHEN, kwELSE, kwEND, kwAND, kwOR, kwNOT,
		kwIS, kwIN, kwBETWEEN, kwLIKE, kwLIKEC, kwLIKE2, kwLIKE4,
		kwINTO, kwVALUES, kwSET, kwRETURNING, kwPIVOT, kwUNPIVOT,
		kwMODEL, kwWITH, kwKEEP, kwOVER, kwUSING,
		kwPARTITION, kwSUBPARTITION, kwSAMPLE, kwVERSIONS:
		return false
	}
	return p.cur.Type >= 2000 &&
		!isOracleSQLReservedKeyword(p.cur) &&
		!isOracleClauseStarterKeyword(p.cur.Type)
}

// parseSubqueryRef parses a subquery in FROM: ( SELECT ... ) alias.
func (p *Parser) parseSubqueryRef(start int) (nodes.TableExpr, error) {
	p.advance() // consume '('

	subSel, parseErr979 := p.parseSelectStmt()
	if parseErr979 != nil {
		return nil, parseErr979
	}

	ref := &nodes.SubqueryRef{
		Subquery: subSel,
		Loc:      nodes.Loc{Start: start},
	}
	// ( subquery [ subquery_restriction_clause ] )
	if ref.Restriction, parseErr979 = p.parseDMLSubqueryRestriction(subSel); parseErr979 != nil {
		return nil, parseErr979
	}

	if p.cur.Type == ')' {
		p.advance()
	}

	if p.pivotClauseAhead() || p.unpivotClauseAhead() {
		ref.Loc.End = p.prev.End
		return p.parsePivotSuffix(ref)
	}

	// Optional alias
	if p.cur.Type == kwAS {
		p.advance()
		var parseErr980 error
		ref.Alias, parseErr980 = p.parseAlias()
		if parseErr980 != nil {
			return nil, parseErr980
		}
	} else if p.isTableAliasCandidate() {
		var parseErr981 error
		ref.Alias, parseErr981 = p.parseAlias()
		if parseErr981 != nil {
			return nil, parseErr981
		}
	}

	ref.Loc.End = p.prev.End
	return p.parsePivotSuffix(ref)
}

func (p *Parser) parseParenthesizedTableRef(start int) (nodes.TableExpr, error) {
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()
	p.advance() // consume '('

	tref, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	if tref == nil {
		return nil, p.syntaxErrorAtCur()
	}
	tref, err = p.parseJoinContinuation(tref)
	if err != nil {
		return nil, err
	}
	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()

	if p.pivotClauseAhead() || p.unpivotClauseAhead() {
		return p.parsePivotSuffix(tref)
	}

	// An inline view inside redundant parentheses takes its alias after
	// them: FROM ((SELECT ...)) x and FROM (((q) UNION (q))) x compile on
	// Oracle 23ai.
	if sub, ok := tref.(*nodes.SubqueryRef); ok && sub.Alias == nil {
		if p.cur.Type == kwAS {
			p.advance()
			alias, err := p.parseAlias()
			if err != nil {
				return nil, err
			}
			sub.Alias = alias
		} else if p.isTableAliasCandidate() {
			alias, err := p.parseAlias()
			if err != nil {
				return nil, err
			}
			sub.Alias = alias
		}
	}
	if aliasable, ok := tref.(*nodes.TableRef); ok {
		if p.cur.Type == kwAS {
			p.advance()
			alias, err := p.parseAlias()
			if err != nil {
				return nil, err
			}
			aliasable.Alias = alias
		} else if p.isTableAliasCandidate() {
			alias, err := p.parseAlias()
			if err != nil {
				return nil, err
			}
			aliasable.Alias = alias
		}
	}

	switch n := tref.(type) {
	case *nodes.JoinClause:
		n.Loc.Start = start
		n.Loc.End = p.prev.End
	case *nodes.TableRef:
		n.Loc.Start = start
		n.Loc.End = p.prev.End
	case *nodes.SubqueryRef:
		n.Loc.Start = start
		n.Loc.End = p.prev.End
	case *nodes.InlineAnalyticView:
		n.Loc.Start = start
		n.Loc.End = p.prev.End
	}
	return p.parsePivotSuffix(tref)
}

// parseLateralRef parses a LATERAL inline view.
//
//	LATERAL ( subquery ) [ alias ]
func (p *Parser) parseLateralRef(start int) (nodes.TableExpr, error) {
	p.advance() // consume LATERAL

	if p.cur.Type != '(' {
		return nil, nil
	}
	p.advance() // consume '('

	subSel, parseErr982 := p.parseSelectStmt()
	if parseErr982 != nil {
		return nil, parseErr982
	}

	if p.cur.Type == ')' {
		p.advance()
	}

	ref := &nodes.LateralRef{
		Subquery: subSel,
		Loc:      nodes.Loc{Start: start},
	}

	if p.cur.Type == kwAS {
		p.advance()
		var parseErr983 error
		ref.Alias, parseErr983 = p.parseAlias()
		if parseErr983 != nil {
			return nil, parseErr983
		}
	} else if p.isTableAliasCandidate() {
		var parseErr984 error
		ref.Alias, parseErr984 = p.parseAlias()
		if parseErr984 != nil {
			return nil, parseErr984
		}
	}

	ref.Loc.End = p.prev.End
	return ref, nil
}

// parseXmlTableRef parses an XMLTABLE expression in FROM.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/XMLTABLE.html
//
//	XMLTABLE ( xpath_string PASSING xml_expr COLUMNS column_def [, ...] )
func (p *Parser) parseXmlTableRef(start int) (nodes.TableExpr, error) {
	p.advance() // consume XMLTABLE

	ref := &nodes.XmlTableRef{Loc: nodes.Loc{Start: start}}

	if p.cur.Type != '(' {
		return ref, nil
	}
	p.advance()
	var // consume '('
	parseErr985 error

	// XPath expression (usually a string literal)
	ref.XPath, parseErr985 = p.parseExpr()
	if parseErr985 !=

		// PASSING xml_expr
		nil {
		return nil, parseErr985
	}

	if p.cur.Type == kwPASSING {
		p.advance()
		var parseErr986 error
		ref.Passing, parseErr986 = p.parseExpr()
		if parseErr986 !=

			// COLUMNS column_def [, ...]
			nil {
			return nil, parseErr986
		}
	}

	if p.cur.Type == kwCOLUMNS {
		p.advance()
		var parseErr987 error
		ref.Columns, parseErr987 = p.parseXmlTableColumns()
		if parseErr987 != nil {
			return nil, parseErr987
		}
	}

	if p.cur.Type == ')' {
		p.advance()
	}

	// Optional alias
	if p.cur.Type == kwAS {
		p.advance()
		var parseErr988 error
		ref.Alias, parseErr988 = p.parseAlias()
		if parseErr988 != nil {
			return nil, parseErr988
		}
	} else if p.isTableAliasCandidate() {
		var parseErr989 error
		ref.Alias, parseErr989 = p.parseAlias()
		if parseErr989 != nil {
			return nil, parseErr989
		}
	}

	ref.Loc.End = p.prev.End
	return ref, nil
}

// parseXmlTableColumns parses column definitions in XMLTABLE.
func (p *Parser) parseXmlTableColumns() (*nodes.List, error) {
	list := &nodes.List{}
	for {
		col, parseErr990 := p.parseXmlTableColumn()
		if parseErr990 != nil {
			return nil, parseErr990
		}
		if col == nil {
			break
		}
		list.Items = append(list.Items, col)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return list, nil
}

// parseXmlTableColumn parses a single XMLTABLE column definition.
//
//	name { datatype [PATH path] [DEFAULT default] | FOR ORDINALITY }
func (p *Parser) parseXmlTableColumn() (*nodes.XmlTableColumn, error) {
	if !p.isIdentLike() {
		return nil, nil
	}
	start := p.pos()
	col := &nodes.XmlTableColumn{Loc: nodes.Loc{Start: start}}
	var parseErr991 error
	col.Name, parseErr991 = p.parseIdentifier()
	if parseErr991 !=

		// FOR ORDINALITY
		nil {
		return nil, parseErr991
	}

	if p.cur.Type == kwFOR && p.peekNext().Type == kwORDINALITY {
		p.advance() // consume FOR
		p.advance() // consume ORDINALITY
		col.ForOrdinality = true
		col.Loc.End = p.prev.End
		return col, nil
	}
	var parseErr992 error

	// Data type
	col.TypeName, parseErr992 = p.parseTypeName()
	if parseErr992 !=

		// PATH
		nil {
		return nil, parseErr992
	}

	if p.cur.Type == kwPATH {
		p.advance()
		var parseErr993 error
		col.Path, parseErr993 = p.parseExpr()
		if parseErr993 !=

			// DEFAULT
			nil {
			return nil, parseErr993
		}
	}

	if p.cur.Type == kwDEFAULT {
		p.advance()
		var parseErr994 error
		col.Default, parseErr994 = p.parseExpr()
		if parseErr994 != nil {
			return nil, parseErr994
		}
	}

	col.Loc.End = p.prev.End
	return col, nil
}

// parseJsonTableRef parses a JSON_TABLE expression in FROM.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/JSON_TABLE.html
//
//	JSON_TABLE ( expr, path_string COLUMNS ( column_def [, ...] ) )
func (p *Parser) parseJsonTableRef(start int) (nodes.TableExpr, error) {
	p.advance() // consume JSON_TABLE

	ref := &nodes.JsonTableRef{Loc: nodes.Loc{Start: start}}

	if p.cur.Type != '(' {
		return ref, nil
	}
	p.advance()
	var // consume '('
	parseErr995 error

	// JSON expression
	ref.Expr, parseErr995 = p.parseExpr()
	if parseErr995 !=

		// Comma separator
		nil {
		return nil, parseErr995
	}

	if p.cur.Type == ',' {
		p.advance()
		var parseErr996 error

		// Path expression (string literal)
		ref.Path, parseErr996 = p.parseExpr()
		if parseErr996 != nil {
			return nil, parseErr996
		}
	} else if p.cur.Type != kwCOLUMNS {
		var parseErr996 error
		ref.Path, parseErr996 = p.parseExpr()
		if parseErr996 != nil {
			return nil, parseErr996
		}
	}

	if p.cur.Type == kwCOLUMNS {
		p.advance()
		if p.cur.Type == '(' {
			p.advance()
			var parseErr997 error
			ref.Columns, parseErr997 = p.parseJsonTableColumns()
			if parseErr997 != nil {
				return nil, parseErr997
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}
	}

	if p.cur.Type == ')' {
		p.advance()
	}

	// Optional alias
	if p.cur.Type == kwAS {
		p.advance()
		var parseErr998 error
		ref.Alias, parseErr998 = p.parseAlias()
		if parseErr998 != nil {
			return nil, parseErr998
		}
	} else if p.isTableAliasCandidate() {
		var parseErr999 error
		ref.Alias, parseErr999 = p.parseAlias()
		if parseErr999 != nil {
			return nil, parseErr999
		}
	}

	ref.Loc.End = p.prev.End
	return ref, nil
}

// parseJsonTableColumns parses column definitions in JSON_TABLE.
func (p *Parser) parseJsonTableColumns() (*nodes.List, error) {
	list := &nodes.List{}
	for {
		col, parseErr1000 := p.parseJsonTableColumn()
		if parseErr1000 != nil {
			return nil, parseErr1000
		}
		if col == nil {
			break
		}
		list.Items = append(list.Items, col)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return list, nil
}

// parseJsonTableColumn parses a single JSON_TABLE column definition.
//
//	name datatype [PATH path] | name FOR ORDINALITY | NESTED [PATH] path COLUMNS ( ... )
func (p *Parser) parseJsonTableColumn() (*nodes.JsonTableColumn, error) {
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()
	start := p.pos()

	// NESTED [PATH] path COLUMNS (...)
	if p.cur.Type == kwNESTED {
		p.advance() // consume NESTED
		col := &nodes.JsonTableColumn{Loc: nodes.Loc{Start: start}}
		if p.cur.Type == kwPATH {
			p.advance() // consume PATH
		}
		nested := &nodes.JsonTableRef{Loc: nodes.Loc{Start: start}}
		var parseErr1001 error
		nested.Path, parseErr1001 = p.parseExpr()
		if parseErr1001 != nil {
			return nil, parseErr1001
		}
		if p.cur.Type == kwCOLUMNS {
			p.advance()
			if p.cur.Type == '(' {
				p.advance()
				var parseErr1002 error
				nested.Columns, parseErr1002 = p.parseJsonTableColumns()
				if parseErr1002 != nil {
					return nil, parseErr1002
				}
				if p.cur.Type == ')' {
					p.advance()
				}
			}
		}
		nested.Loc.End = p.prev.End
		col.Nested = nested
		col.Loc.End = p.prev.End
		return col, nil
	}

	if !p.isIdentLike() {
		return nil, nil
	}

	col := &nodes.JsonTableColumn{Loc: nodes.Loc{Start: start}}
	var parseErr1003 error
	col.Name, parseErr1003 = p.parseIdentifier()
	if parseErr1003 !=

		// FOR ORDINALITY
		nil {
		return nil, parseErr1003
	}

	if p.cur.Type == kwFOR && p.peekNext().Type == kwORDINALITY {
		p.advance() // consume FOR
		p.advance() // consume ORDINALITY
		col.ForOrdinality = true
		col.Loc.End = p.prev.End
		return col, nil
	}
	var parseErr1004 error

	// Data type
	col.TypeName, parseErr1004 = p.parseTypeName()
	if parseErr1004 !=

		// EXISTS
		nil {
		return nil, parseErr1004
	}

	if parseErr1004 = p.parseJsonTableValueColumnModifiers(); parseErr1004 != nil {
		return nil, parseErr1004
	}

	if p.cur.Type == kwEXISTS {
		col.Exists = true
		p.advance()
	}

	// PATH
	if p.cur.Type == kwPATH {
		p.advance()
		var parseErr1005 error
		col.Path, parseErr1005 = p.parseExpr()
		if parseErr1005 != nil {
			return nil, parseErr1005
		}
	}

	col.Loc.End = p.prev.End
	return col, nil
}

func (p *Parser) parseJsonTableValueColumnModifiers() error {
	if p.cur.Type == kwFORMAT {
		p.advance()
		if p.cur.Type == kwJSON {
			p.advance()
		} else {
			return p.syntaxErrorAtCur()
		}
	}

	switch {
	case p.cur.Type == kwWITH:
		p.advance()
		if p.isIdentLikeStr("CONDITIONAL") || p.isIdentLikeStr("UNCONDITIONAL") {
			p.advance()
		}
		if p.isIdentLikeStr("ARRAY") {
			p.advance()
		}
		if p.isIdentLikeStr("WRAPPER") {
			p.advance()
		}
	case p.isIdentLikeStr("WITHOUT"):
		p.advance()
		if p.isIdentLikeStr("ARRAY") {
			p.advance()
		}
		if p.isIdentLikeStr("WRAPPER") {
			p.advance()
		}
	}
	return nil
}

// parseHierarchicalClause parses START WITH / CONNECT BY clauses (either order).
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/Hierarchical-Queries.html
//
//	[ START WITH condition ] CONNECT BY [ NOCYCLE ] condition
//	CONNECT BY [ NOCYCLE ] condition [ START WITH condition ]
func (p *Parser) parseHierarchicalClause() (*nodes.HierarchicalClause, error) {
	start := p.pos()
	hc := &nodes.HierarchicalClause{
		Loc: nodes.Loc{Start: start},
	}
	startWith := func() error {
		p.advance() // START
		if p.cur.Type != kwWITH {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		var err error
		if hc.StartWith, err = p.parseExpr(); err != nil {
			return err
		}
		if hc.StartWith == nil {
			return p.syntaxErrorAtCur()
		}
		return nil
	}

	// START WITH first, then CONNECT BY
	if p.cur.Type == kwSTART {
		if err := startWith(); err != nil {
			return nil, err
		}
	}

	if p.cur.Type == kwCONNECT {
		p.advance() // CONNECT
		if p.cur.Type != kwBY {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		if p.isIdentLikeStr("NOCYCLE") {
			hc.IsNocycle = true
			p.advance()
		}
		var err error
		if hc.ConnectBy, err = p.parseExpr(); err != nil {
			return nil, err
		}
		if hc.ConnectBy == nil {
			return nil, p.syntaxErrorAtCur()
		}
	}

	// START WITH may come after CONNECT BY
	if hc.StartWith == nil && p.cur.Type == kwSTART {
		if err := startWith(); err != nil {
			return nil, err
		}
	}

	hc.Loc.End = p.prev.End
	return hc, nil
}

// parseOrderByList parses a comma-separated list of ORDER BY items.
func (p *Parser) parseOrderByList() (*nodes.List, error) {
	list := &nodes.List{}

	for {
		sb, parseErr1012 := p.parseSortBy()
		if parseErr1012 != nil {
			return nil, parseErr1012
		}
		if sb == nil {
			break
		}
		list.Items = append(list.Items, sb)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	return list, nil
}

// parseSortBy parses a single ORDER BY item.
//
//	sort_key [ ASC | DESC ] [ NULLS { FIRST | LAST } ]
func (p *Parser) parseSortBy() (*nodes.SortBy, error) {
	start := p.pos()
	expr, parseErr1013 := p.parseExpr()
	if parseErr1013 != nil {
		return nil, parseErr1013
	}
	if expr == nil {
		return nil, nil
	}

	sb := &nodes.SortBy{
		Expr: expr,
		Loc:  nodes.Loc{Start: start},
	}

	// ASC | DESC
	switch p.cur.Type {
	case kwASC:
		sb.Dir = nodes.SORTBY_ASC
		p.advance()
	case kwDESC:
		sb.Dir = nodes.SORTBY_DESC
		p.advance()
	}

	// NULLS FIRST | NULLS LAST
	if p.cur.Type == kwNULLS {
		p.advance()
		switch p.cur.Type {
		case kwFIRST:
			sb.NullOrder = nodes.SORTBY_NULLS_FIRST
			p.advance()
		case kwLAST:
			sb.NullOrder = nodes.SORTBY_NULLS_LAST
			p.advance()
		default:
			return nil, p.syntaxErrorAtCur()
		}
	}

	sb.Loc.End = p.prev.End
	return sb, nil
}

// parseForUpdateClause parses FOR UPDATE [OF ...] [NOWAIT | WAIT n | SKIP LOCKED].
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/SELECT.html#GUID-CFA006CA-6FF1-4972-821E-6996142A51C6
func (p *Parser) parseForUpdateClause() (*nodes.ForUpdateClause, error) {
	start := p.pos()
	p.advance() // consume FOR

	// FOR takes UPDATE (ORA-02000).
	if p.cur.Type != kwUPDATE {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()

	fu := &nodes.ForUpdateClause{
		Loc: nodes.Loc{Start: start},
	}

	// OF table_list
	if p.cur.Type == kwOF {
		p.advance()
		var parseErr1014 error
		fu.Tables, parseErr1014 = p.parseExprList()
		if parseErr1014 !=

			// NOWAIT | WAIT n | SKIP LOCKED
			nil {
			return nil, parseErr1014
		}
	}

	switch p.cur.Type {
	case kwNOWAIT:
		fu.NoWait = true
		p.advance()
	case kwWAIT:
		p.advance()
		var parseErr1015 error
		fu.Wait, parseErr1015 = p.parseExpr()
		if parseErr1015 != nil {
			return nil, parseErr1015
		}
		if fu.Wait == nil {
			return nil, p.syntaxErrorAtCur()
		}
	case kwSKIP:
		p.advance()
		if p.cur.Type != kwLOCKED {
			return nil, p.syntaxErrorAtCur()
		}
		fu.SkipLocked = true
		p.advance()
	}

	fu.Loc.End = p.prev.End
	return fu, nil
}

// parseFetchFirstClause parses OFFSET/FETCH FIRST clause.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/SELECT.html
//
//	[ OFFSET n { ROW | ROWS } ]
//	FETCH { FIRST | NEXT } [ n [ PERCENT ] ] { ROW | ROWS } { ONLY | WITH TIES }
func (p *Parser) parseFetchFirstClause() (*nodes.FetchFirstClause, error) {
	start := p.pos()
	fc := &nodes.FetchFirstClause{
		Loc: nodes.Loc{Start: start},
	}

	// OFFSET n ROWS
	if p.cur.Type == kwOFFSET {
		p.advance()
		if p.cur.Type == kwROW || p.cur.Type == kwROWS {
			return nil, p.syntaxErrorAtCur()
		}
		var parseErr1016 error
		fc.Offset, parseErr1016 = p.parseExpr()
		if parseErr1016 !=
			// ROW | ROWS
			nil {
			return nil, parseErr1016
		}
		if fc.Offset == nil {
			return nil, p.syntaxErrorAtCur()
		}

		// ROW | ROWS (ORA-02000)
		if p.cur.Type != kwROW && p.cur.Type != kwROWS {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
	}

	// FETCH FIRST|NEXT
	if p.cur.Type == kwFETCH {
		p.advance()
		// FIRST | NEXT
		if p.cur.Type != kwFIRST && p.cur.Type != kwNEXT {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()

		// count expression
		if p.cur.Type != kwROW && p.cur.Type != kwROWS {
			var parseErr1017 error
			fc.Count, parseErr1017 = p.parseExpr()
			if parseErr1017 !=

				// PERCENT
				nil {
				return nil, parseErr1017
			}
			if fc.Count == nil {
				return nil, p.syntaxErrorAtCur()
			}
		}

		if p.cur.Type == kwPERCENT {
			fc.Percent = true
			p.advance()
		}

		// ROW | ROWS
		if p.cur.Type != kwROW && p.cur.Type != kwROWS {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()

		// ONLY | WITH TIES
		if p.cur.Type == kwONLY {
			p.advance()
		} else if p.cur.Type == kwWITH {
			p.advance()
			if p.cur.Type != kwTIES {
				return nil, p.syntaxErrorAtCur()
			}
			fc.WithTies = true
			p.advance()
		} else {
			return nil, p.syntaxErrorAtCur()
		}
	}

	fc.Loc.End = p.prev.End
	return fc, nil
}

// parseWithClause parses a WITH clause (common table expressions).
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/SELECT.html
//
//	WITH [ RECURSIVE ] cte_name [ ( col1, col2, ... ) ] AS ( subquery ) [, ...]
func (p *Parser) parseWithClause() (*nodes.WithClause, error) {
	start := p.pos()
	p.advance() // consume WITH

	wc := &nodes.WithClause{
		CTEs: &nodes.List{},
		Loc:  nodes.Loc{Start: start},
	}

	if p.cur.Type == kwRECURSIVE {
		wc.Recursive = true
		p.advance()
	}

	// plsql_declarations: FUNCTION and PROCEDURE definitions, then the
	// query or more elements.
	if p.cur.Type == kwFUNCTION || p.cur.Type == kwPROCEDURE {
		var err error
		if wc.PLSQLDeclarations, err = p.parseSelectWithDeclarations(); err != nil {
			return nil, err
		}
		if p.cur.Type == kwSELECT || p.cur.Type == '(' {
			wc.Loc.End = p.prev.End
			return wc, nil
		}
	}

	for {
		cte, parseErr1018 := p.parseCTE()
		if parseErr1018 != nil {
			return nil, parseErr1018
		}
		// WITH names a query first (ORA-00903), and again after a ','.
		if cte == nil {
			return nil, p.syntaxErrorAtCur()
		}
		wc.CTEs.Items = append(wc.CTEs.Items, cte)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	wc.Loc.End = p.prev.End
	return wc, nil
}

// parseCTE parses a single common table expression.
func (p *Parser) parseCTE() (*nodes.CTE, error) {
	if !p.isIdentLike() {
		return nil, nil
	}

	start := p.pos()
	name, parseErr1019 := p.parseIdentifier()
	if parseErr1019 != nil {
		return nil, parseErr1019
	}

	cte := &nodes.CTE{
		Name: name,
		Loc:  nodes.Loc{Start: start},
	}

	// subav_factoring_clause: a derived analytic view,
	// subav_name ANALYTIC VIEW AS ( sub_av_clause ), instead of a query.
	// It takes no column list and no SEARCH or CYCLE.
	if p.isKeywordStr("ANALYTIC") {
		p.advance() // consume ANALYTIC
		if p.cur.Type != kwVIEW {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance() // consume VIEW
		if p.cur.Type != kwAS {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance() // consume AS
		if p.cur.Type != '(' {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance() // consume (
		var err error
		if cte.SubAV, err = p.parseSubAVClause(); err != nil {
			return nil, err
		}
		if err := p.anExpect(')'); err != nil {
			return nil, err
		}
		cte.Loc.End = p.prev.End
		return cte, nil
	}

	// Optional column list, each name once (ORA-32049 at its ')').
	if p.cur.Type == '(' {
		p.advance()
		cte.Columns = &nodes.List{}
		seen := map[string]bool{}
		dup := false
		for {
			if !p.isIdentLike() {
				return nil, p.syntaxErrorAtCur()
			}
			col, parseErr1020 := p.parseIdentifier()
			if parseErr1020 != nil {
				return nil, parseErr1020
			}
			dup = dup || seen[col]
			seen[col] = true
			cte.Columns.Items = append(cte.Columns.Items, &nodes.String{Str: col})
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if p.cur.Type != ')' || dup {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
	}

	// AS ( subquery ) (ORA-02000 without AS)
	if p.cur.Type != kwAS {
		return nil, p.syntaxErrorAtCur()
	}
	asTok := p.cur
	p.advance()
	if p.cur.Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	if !p.atQuery() {
		return nil, p.syntaxErrorAtCur()
	}
	var parseErr1021 error
	cte.Query, parseErr1021 = p.parseSelectStmt()
	if parseErr1021 != nil {
		return nil, parseErr1021
	}
	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()

	// SEARCH and CYCLE belong to a recursive WITH element, which names its
	// columns (ORA-32488 at its AS).
	if (p.isIdentLikeStr("SEARCH") || p.cur.Type == kwCYCLE) && cte.Columns == nil {
		return nil, p.syntaxErrorAtTok(asTok)
	}

	// SEARCH { BREADTH FIRST | DEPTH FIRST } BY col [, ...] SET ordering_column
	if p.isIdentLikeStr("SEARCH") {
		var parseErr1022 error
		cte.Search, parseErr1022 = p.parseCTESearchClause()
		if parseErr1022 !=

			// CYCLE col [, ...] SET cycle_mark TO value DEFAULT no_value
			nil {
			return nil, parseErr1022
		}
	}

	if p.cur.Type == kwCYCLE {
		var parseErr1023 error
		cte.Cycle, parseErr1023 = p.parseCTECycleClause()
		if parseErr1023 != nil {
			return nil, parseErr1023
		}
	}

	cte.Loc.End = p.prev.End
	return cte, nil
}

// parseCTESearchClause parses a SEARCH clause on a CTE.
//
// BNF: oracle/parser/bnf/SELECT.bnf
//
//	search_clause:
//	    SEARCH { BREADTH FIRST | DEPTH FIRST }
//	    BY c_alias [ ASC | DESC ] [ NULLS FIRST | NULLS LAST ]
//	       [, c_alias [ ASC | DESC ] [ NULLS FIRST | NULLS LAST ] ]...
//	    SET ordering_column
func (p *Parser) parseCTESearchClause() (*nodes.CTESearchClause, error) {
	start := p.pos()
	p.advance() // consume SEARCH

	sc := &nodes.CTESearchClause{Loc: nodes.Loc{Start: start}}

	// BREADTH FIRST | DEPTH FIRST
	if p.isIdentLikeStr("BREADTH") {
		sc.BreadthFirst = true
		p.advance()
	} else if p.isIdentLikeStr("DEPTH") {
		p.advance()
	}
	if p.cur.Type == kwFIRST {
		p.advance()
	}

	// BY col [ASC|DESC] [NULLS FIRST|LAST] [, ...]
	if p.cur.Type == kwBY {
		p.advance()
	}
	var parseErr1024 error
	sc.Columns, parseErr1024 = p.parseOrderByList()
	if parseErr1024 !=

		// SET ordering_column
		nil {
		return nil, parseErr1024
	}

	if p.cur.Type == kwSET {
		p.advance()
		var parseErr1025 error
		sc.SetColumn, parseErr1025 = p.parseIdentifier()
		if parseErr1025 != nil {
			return nil, parseErr1025
		}
	}

	sc.Loc.End = p.prev.End
	return sc, nil
}

// parseCTECycleClause parses a CYCLE clause on a CTE.
//
// BNF: oracle/parser/bnf/SELECT.bnf
//
//	cycle_clause:
//	    CYCLE c_alias [, c_alias ]...
//	    SET cycle_mark_c_alias TO cycle_value
//	    DEFAULT no_cycle_value
func (p *Parser) parseCTECycleClause() (*nodes.CTECycleClause, error) {
	start := p.pos()
	p.advance() // consume CYCLE

	cc := &nodes.CTECycleClause{Loc: nodes.Loc{Start: start}}

	// Column list
	cc.Columns = &nodes.List{}
	for {
		col, parseErr1026 := p.parseIdentifier()
		if parseErr1026 != nil {
			return nil, parseErr1026
		}
		if col == "" {
			break
		}
		cc.Columns.Items = append(cc.Columns.Items, &nodes.String{Str: col})
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	// SET cycle_mark_alias
	if p.cur.Type == kwSET {
		p.advance()
		var parseErr1027 error
		cc.SetColumn, parseErr1027 = p.parseIdentifier()
		if parseErr1027 !=

			// TO cycle_value
			nil {
			return nil, parseErr1027
		}
	}

	if p.cur.Type == kwTO {
		p.advance()
		var parseErr1028 error
		cc.CycleValue, parseErr1028 = p.parseExpr()
		if parseErr1028 !=

			// DEFAULT no_cycle_value
			nil {
			return nil, parseErr1028
		}
	}

	if p.cur.Type == kwDEFAULT {
		p.advance()
		var parseErr1029 error
		cc.NoCycleValue, parseErr1029 = p.parseExpr()
		if parseErr1029 != nil {
			return nil, parseErr1029
		}
	}

	cc.Loc.End = p.prev.End
	return cc, nil
}

// pivotClauseAhead reports whether the current position starts a real
// pivot_clause: the PIVOT keyword followed by '(' or XML. PIVOT and UNPIVOT
// are non-reserved in Oracle, so a trailing PIVOT not followed by its body is
// an alias position, not a clause.
func (p *Parser) pivotClauseAhead() bool {
	if p.cur.Type != kwPIVOT {
		return false
	}
	nxt := p.peekNext()
	return nxt.Type == '(' || (nxt.Type == tokIDENT && strings.EqualFold(nxt.Str, "XML"))
}

// unpivotClauseAhead reports whether the current position starts a real
// unpivot_clause: UNPIVOT followed by '(' or INCLUDE/EXCLUDE NULLS.
func (p *Parser) unpivotClauseAhead() bool {
	if p.cur.Type != kwUNPIVOT {
		return false
	}
	nxt := p.peekNext()
	return nxt.Type == '(' || nxt.Type == kwINCLUDE ||
		(nxt.Type == tokIDENT && strings.EqualFold(nxt.Str, "EXCLUDE"))
}

// parsePivotSuffix attaches any PIVOT / UNPIVOT clauses — each with its
// optional trailing t_alias — to a just-parsed table reference, per the
// table_reference grammar: query_table_expression
// [pivot_clause | unpivot_clause] [t_alias]. The clause becomes the from-item
// itself (both implement TableExpr), so join continuations and the alias
// participate in the FROM chain like any other table reference. Oracle also
// accepts chained clauses (t PIVOT (...) a PIVOT (...) b), so this loops.
// Each clause starts where the source it wraps does.
// A bare identifier alias is accepted; AS stays rejected (Oracle t_alias
// takes no AS keyword). No-op when no clause follows.
func (p *Parser) parsePivotSuffix(source nodes.TableExpr) (nodes.TableExpr, error) {
	for {
		switch {
		case p.pivotClauseAhead():
			pc, err := p.parsePivotClause()
			if err != nil {
				return nil, err
			}
			pc.Source = source
			pc.Loc.Start = exprLocStart(source, pc.Loc.Start)
			if p.isTableAliasCandidate() {
				pc.Alias, err = p.parseAlias()
				if err != nil {
					return nil, err
				}
				pc.Loc.End = p.prev.End
			}
			source = pc
		case p.unpivotClauseAhead():
			uc, err := p.parseUnpivotClause()
			if err != nil {
				return nil, err
			}
			uc.Source = source
			uc.Loc.Start = exprLocStart(source, uc.Loc.Start)
			if p.isTableAliasCandidate() {
				uc.Alias, err = p.parseAlias()
				if err != nil {
					return nil, err
				}
				uc.Loc.End = p.prev.End
			}
			source = uc
		default:
			return source, nil
		}
	}
}

func (p *Parser) parseIdentifierListInParens() (*nodes.List, error) {
	list := &nodes.List{}
	if p.cur.Type != '(' {
		return list, nil
	}
	p.advance()
	for {
		if !p.isIdentLike() {
			return nil, p.syntaxErrorAtCur()
		}
		name, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		list.Items = append(list.Items, &nodes.String{Str: name})
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	return list, nil
}

// parseModelClause parses an Oracle MODEL clause.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/SELECT.html
//
//	model_clause ::=
//	    MODEL
//	      [ cell_reference_options ]
//	      [ return_rows_clause ]
//	      [ reference_model ]...
//	      main_model
//
//	cell_reference_options ::=
//	    { IGNORE NAV | KEEP NAV }
//	    { UNIQUE DIMENSION | UNIQUE SINGLE REFERENCE }
//
//	return_rows_clause ::=
//	    RETURN { UPDATED | ALL } ROWS
//
//	reference_model ::=
//	    REFERENCE reference_model_name ON ( subquery )
//	        model_column_clauses
//	        [ cell_reference_options ]
//
//	main_model ::=
//	    [ MAIN main_model_name ]
//	        model_column_clauses
//	        [ cell_reference_options ]
//	        model_rules_clause
//
//	model_column_clauses ::=
//	    [ PARTITION BY ( expr [ [ AS ] c_alias ] [, ...] ) ]
//	    DIMENSION BY ( expr [ [ AS ] c_alias ] [, ...] )
//	    MEASURES ( expr [ [ AS ] c_alias ] [, ...] )
//
//	model_rules_clause ::=
//	    [ RULES ]
//	    [ { UPDATE | UPSERT [ ALL ] } ]
//	    [ { AUTOMATIC | SEQUENTIAL } ORDER ]
//	    [ model_iterate_clause ]
//	    ( cell_assignment [, ...] )
//
//	model_iterate_clause ::=
//	    ITERATE ( number ) [ UNTIL ( condition ) ]
//
//	cell_assignment ::=
//	    measure_column [ dimension_subscripts ] = expr
//
//	single_column_for_loop ::=
//	    FOR dimension_column
//	      { IN ( { literal [, ...] | subquery } )
//	      | [ LIKE pattern ] FROM literal TO literal { INCREMENT | DECREMENT } literal
//	      }
//
//	multi_column_for_loop ::=
//	    FOR ( dimension_column [, ...] ) IN
//	      ( ( literal [, ...] ) [, ( literal [, ...] ) ]...
//	      | subquery
//	      )
func (p *Parser) parseModelClause() (*nodes.ModelClause, error) {
	start := p.pos()
	p.advance() // consume MODEL

	mc := &nodes.ModelClause{
		Loc: nodes.Loc{Start: start},
	}
	var parseErr1043 error

	// cell_reference_options (optional, before RETURN/REFERENCE/MAIN/PARTITION/DIMENSION)
	mc.CellRefOptions, parseErr1043 = p.parseModelCellRefOptions()
	if parseErr1043 !=

		// RETURN { UPDATED | ALL } ROWS
		nil {
		return nil, parseErr1043
	}

	if p.cur.Type == kwRETURN {
		p.advance() // consume RETURN
		if p.cur.Type == kwUPDATED {
			mc.ReturnRows = "UPDATED"
			p.advance()
		} else if p.cur.Type == kwALL {
			mc.ReturnRows = "ALL"
			p.advance()
		}
		if p.cur.Type == kwROWS {
			p.advance()
		}
	}

	// REFERENCE models (zero or more)
	for p.cur.Type == kwREFERENCE {
		ref, parseErr1044 := p.parseModelRefModel()
		if parseErr1044 != nil {
			return nil, parseErr1044
		}
		mc.RefModels = append(mc.RefModels, ref)
	}
	var parseErr1045 error

	// main_model
	mc.MainModel, parseErr1045 = p.parseModelMainModel()
	if parseErr1045 != nil {
		return nil, parseErr1045

		// parseModelCellRefOptions parses optional cell_reference_options.
	}

	mc.Loc.End = p.prev.End
	return mc, nil
}

func (p *Parser) parseModelCellRefOptions() (*nodes.ModelCellRefOptions, error) {
	opts := &nodes.ModelCellRefOptions{Loc: nodes.Loc{Start: p.pos()}}
	found := false

	// IGNORE NAV | KEEP NAV
	if p.isIdentLikeStr("IGNORE") {
		p.advance()
		if p.cur.Type == kwNAV {
			p.advance()
		}
		opts.IgnoreNav = true
		found = true
	} else if p.isIdentLikeStr("KEEP") {
		p.advance()
		if p.cur.Type == kwNAV {
			p.advance()
		}
		opts.KeepNav = true
		found = true
	}

	// UNIQUE DIMENSION | UNIQUE SINGLE REFERENCE
	if p.cur.Type == kwUNIQUE {
		p.advance()
		if p.cur.Type == kwDIMENSION {
			opts.UniqueDimension = true
			p.advance()
			found = true
		} else if p.isIdentLikeStr("SINGLE") {
			p.advance()
			if p.cur.Type == kwREFERENCE {
				p.advance()
			}
			opts.UniqueSingleRef = true
			found = true
		}
	}

	if !found {
		return nil, nil
	}

	opts.Loc.End = p.prev.End
	return opts, nil
}

// parseModelRefModel parses a REFERENCE model.
func (p *Parser) parseModelRefModel() (*nodes.ModelRefModel, error) {
	start := p.pos()
	p.advance() // consume REFERENCE

	ref := &nodes.ModelRefModel{
		Loc: nodes.Loc{Start: start},
	}
	var parseErr1046 error

	// reference_model_name
	ref.Name, parseErr1046 = p.parseIdentifier()
	if parseErr1046 !=

		// ON ( subquery )
		nil {
		return nil, parseErr1046
	}

	if p.cur.Type == kwON {
		p.advance()
		if p.cur.Type == '(' {
			p.advance()
			var parseErr1047 error
			ref.Subquery, parseErr1047 = p.parseSelectStmt()
			if parseErr1047 != nil {
				return nil, parseErr1047
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}
	}
	var parseErr1048 error

	// model_column_clauses
	ref.ColumnClauses, parseErr1048 = p.parseModelColumnClauses()
	if parseErr1048 !=

		// optional cell_reference_options
		nil {
		return nil, parseErr1048
	}
	var parseErr1049 error

	ref.CellRefOptions, parseErr1049 = p.parseModelCellRefOptions()
	if parseErr1049 != nil {
		return nil, parseErr1049
	}

	ref.Loc.End = p.prev.End
	return ref, nil
}

// parseModelMainModel parses the main_model.
func (p *Parser) parseModelMainModel() (*nodes.ModelMainModel, error) {
	start := p.pos()
	mm := &nodes.ModelMainModel{
		Loc: nodes.Loc{Start: start},
	}

	// [ MAIN main_model_name ]
	if p.cur.Type == kwMAIN {
		p.advance()
		var parseErr1050 error
		mm.Name, parseErr1050 = p.parseIdentifier()
		if parseErr1050 !=

			// model_column_clauses
			nil {
			return nil, parseErr1050
		}
	}
	var parseErr1051 error

	mm.ColumnClauses, parseErr1051 = p.parseModelColumnClauses()
	if parseErr1051 !=

		// optional cell_reference_options
		nil {
		return nil, parseErr1051
	}
	var parseErr1052 error

	mm.CellRefOptions, parseErr1052 = p.parseModelCellRefOptions()
	if parseErr1052 !=

		// model_rules_clause
		nil {
		return nil, parseErr1052
	}
	var parseErr1053 error

	mm.RulesClause, parseErr1053 = p.parseModelRulesClause()
	if parseErr1053 != nil {
		return nil, parseErr1053

		// parseModelColumnClauses parses PARTITION BY / DIMENSION BY / MEASURES.
	}

	mm.Loc.End = p.prev.End
	return mm, nil
}

func (p *Parser) parseModelColumnClauses() (*nodes.ModelColumnClauses, error) {
	start := p.pos()
	cc := &nodes.ModelColumnClauses{
		Loc: nodes.Loc{Start: start},
	}

	// [ PARTITION BY ( expr [AS alias], ... ) ]
	if p.cur.Type == kwPARTITION {
		p.advance() // PARTITION
		if p.cur.Type == kwBY {
			p.advance() // BY
		}
		if p.cur.Type == '(' {
			p.advance()
			var parseErr1054 error
			cc.PartitionBy, parseErr1054 = p.parseModelColumnList()
			if parseErr1054 != nil {
				return nil, parseErr1054
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}
	}

	// DIMENSION BY ( expr [AS alias], ... )
	if p.cur.Type == kwDIMENSION {
		p.advance() // DIMENSION
		if p.cur.Type == kwBY {
			p.advance() // BY
		}
		if p.cur.Type == '(' {
			p.advance()
			var parseErr1055 error
			cc.DimensionBy, parseErr1055 = p.parseModelColumnList()
			if parseErr1055 != nil {
				return nil, parseErr1055
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}
	}

	// MEASURES ( expr [AS alias], ... )
	if p.cur.Type == kwMEASURES {
		p.advance() // MEASURES
		if p.cur.Type == '(' {
			p.advance()
			var parseErr1056 error
			cc.Measures, parseErr1056 = p.parseModelColumnList()
			if parseErr1056 != nil {
				return nil, parseErr1056
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}
	}

	cc.Loc.End = p.prev.End
	return cc, nil
}

// parseModelColumnList parses a comma-separated list of expr [ [AS] alias ] for MODEL clauses.
func (p *Parser) parseModelColumnList() (*nodes.List, error) {
	list := &nodes.List{}
	for {
		rt, parseErr1057 := p.parseResTarget()
		if parseErr1057 != nil {
			return nil, parseErr1057
		}
		if rt == nil {
			break
		}
		list.Items = append(list.Items, rt)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return list, nil
}

// parseModelRulesClause parses the model_rules_clause.
func (p *Parser) parseModelRulesClause() (*nodes.ModelRulesClause, error) {
	start := p.pos()
	rc := &nodes.ModelRulesClause{
		Loc: nodes.Loc{Start: start},
	}

	// [ RULES ]
	if p.cur.Type == kwRULES {
		p.advance()
	}

	// [ UPDATE | UPSERT [ALL] ]
	if p.cur.Type == kwUPDATE {
		rc.UpdateMode = "UPDATE"
		p.advance()
	} else if p.cur.Type == kwUPSERT {
		p.advance()
		if p.cur.Type == kwALL {
			rc.UpdateMode = "UPSERT ALL"
			p.advance()
		} else {
			rc.UpdateMode = "UPSERT"
		}
	}

	// [ AUTOMATIC | SEQUENTIAL ] ORDER
	if p.cur.Type == kwAUTOMATIC {
		rc.OrderMode = "AUTOMATIC"
		p.advance()
		if p.cur.Type == kwORDER {
			p.advance()
		}
	} else if p.cur.Type == kwSEQUENTIAL {
		rc.OrderMode = "SEQUENTIAL"
		p.advance()
		if p.cur.Type == kwORDER {
			p.advance()
		}
	}

	// [ ITERATE ( number ) [ UNTIL ( condition ) ] ]
	if p.cur.Type == kwITERATE {
		p.advance()
		if p.cur.Type == '(' {
			p.advance()
			var parseErr1058 error
			rc.Iterate, parseErr1058 = p.parseExpr()
			if parseErr1058 != nil {
				return nil, parseErr1058
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}
		// [ UNTIL ( condition ) ]
		if p.cur.Type == kwUNTIL {
			p.advance()
			if p.cur.Type == '(' {
				p.advance()
				var parseErr1059 error
				rc.Until, parseErr1059 = p.parseExpr()
				if parseErr1059 != nil {
					return nil, parseErr1059
				}
				if p.cur.Type == ')' {
					p.advance()
				}
			}
		}
	}

	// ( cell_assignment [, ...] )
	if p.cur.Type == '(' {
		p.advance()
		rc.Rules = &nodes.List{}
		for {
			rule, parseErr1060 := p.parseModelRule()
			if parseErr1060 != nil {
				return nil, parseErr1060
			}
			if rule == nil {
				break
			}
			rc.Rules.Items = append(rc.Rules.Items, rule)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if p.cur.Type == ')' {
			p.advance()
		}
	}

	rc.Loc.End = p.prev.End
	return rc, nil
}

// parseModelRule parses a single cell_assignment: measure_column[dim1, dim2] = expr.
func (p *Parser) parseModelRule() (*nodes.ModelRule, error) {
	start := p.pos()

	// Parse the left side: measure_column [ dim_subscripts ]
	// The left side is an identifier (possibly qualified) followed by [ subscripts ]
	cellRef, cell, parseErr1061 := p.parseModelCellRef()
	if parseErr1061 != nil {
		return nil, parseErr1061
	}
	if cellRef == nil {
		return nil, nil
	}

	rule := &nodes.ModelRule{
		CellRef: cellRef,
		Cell:    cell,
		Loc:     nodes.Loc{Start: start},
	}

	// = expr
	if p.cur.Type == '=' {
		p.advance()
		var parseErr1062 error
		rule.Expr, parseErr1062 = p.parseExpr()
		if parseErr1062 != nil {
			return nil, parseErr1062
		}
	}

	rule.Loc.End = p.prev.End
	return rule, nil
}

// parseModelCellRef parses a model cell reference: measure[dim1, dim2, ...].
// Dimension subscripts can be expressions, ANY, or FOR loops (single/multi column).
func (p *Parser) parseModelCellRef() (nodes.ExprNode, *nodes.ModelCellReference, error) {
	start := p.pos()

	// Parse the measure column name as a column reference
	if !p.isIdentLike() {
		return nil, nil, nil
	}
	name, parseErr1063 := p.parseIdentifier()
	if parseErr1063 != nil {
		return nil, nil, parseErr1063
	}
	col := &nodes.ColumnRef{
		Column: name,
		Loc:    nodes.Loc{Start: start, End: p.prev.End},
	}

	// [ dim_subscripts ]
	if p.cur.Type != '[' {
		return col, &nodes.ModelCellReference{
			Measure:    name,
			Dimensions: &nodes.List{},
			Loc:        col.Loc,
		}, nil
	}
	p.advance() // consume '['

	// Parse subscript expressions (comma-separated)
	args := &nodes.List{}
	for {
		if p.cur.Type == ']' {
			break
		}
		expr, parseErr1064 := p.parseExpr()
		if parseErr1064 != nil {
			return nil, nil, parseErr1064
		}
		if expr == nil {
			break
		}
		args.Items = append(args.Items, expr)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	if p.cur.Type == ']' {
		p.advance()
	}

	// Represent as a FuncCallExpr with the column name and subscript args
	fn := &nodes.FuncCallExpr{
		FuncName: &nodes.ObjectName{Name: name, Loc: nodes.Loc{Start: start, End: p.prev.End}},
		Args:     args,
		Loc:      nodes.Loc{Start: start, End: p.prev.End},
	}
	return fn, &nodes.ModelCellReference{
		Measure:    name,
		Dimensions: args,
		Loc:        fn.Loc,
	}, nil
}

// parseWindowClause parses a WINDOW clause.
//
// BNF: oracle/parser/bnf/SELECT.bnf
//
//	window_clause:
//	    WINDOW window_name AS ( window_specification )
//	    [, window_name AS ( window_specification ) ]...
func (p *Parser) parseWindowClause() ([]*nodes.WindowDef, error) {
	p.advance() // consume WINDOW (soft keyword)

	// Each window is named once (ORA-32782), and a window it builds on is
	// one named before it (ORA-32784). A name used twice is refused at the
	// second, once the list is read: the rule is the engine's, after its
	// parse, so a statement that breaks earlier in the list reports that
	// first, as Oracle's parser does.
	var defs []*nodes.WindowDef
	var dupErr *ParseError
	// ordered records each window named so far, and whether it orders its
	// rows, itself or through the window it builds on.
	ordered := map[string]bool{}
	for {
		start := p.pos()
		nameTok := p.cur
		if !p.isIdentLike() {
			return nil, p.syntaxErrorAtCur()
		}
		if err := p.syntaxErrorIfReservedIdentifier(); err != nil {
			return nil, err
		}
		name, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		if _, dup := ordered[name]; dup && dupErr == nil {
			dupErr = p.syntaxErrorAtTok(nameTok)
		}
		wd := &nodes.WindowDef{
			Name: name,
			Loc:  nodes.Loc{Start: start},
		}
		if p.cur.Type != kwAS {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		if p.cur.Type != '(' {
			return nil, p.syntaxErrorAtCur()
		}
		if wd.Spec, err = p.parseSelectWindowSpec(nameTok, ordered); err != nil {
			return nil, err
		}
		ordered[name] = wd.Spec.OrderBy != nil || wd.Spec.WindowName != "" && ordered[wd.Spec.WindowName]

		wd.Loc.End = p.prev.End
		defs = append(defs, wd)

		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if dupErr != nil {
		return nil, dupErr
	}
	return defs, nil
}

// parsePartitionExtClause parses a partition extension clause.
//
// BNF: oracle/parser/bnf/SELECT.bnf
//
//	partition_extension_clause:
//	    { PARTITION ( partition_name )
//	    | PARTITION FOR ( partition_key_value [, partition_key_value ]... )
//	    | SUBPARTITION ( subpartition_name )
//	    | SUBPARTITION FOR ( subpartition_key_value [, subpartition_key_value ]... )
//	    }
//
// Each form takes its parentheses and what they hold (ORA-00906, ORA-00936).
// Oracle 23ai reports a missing name at the '(' in a query, UPDATE and
// DELETE (emptyAtOpen), and at the ')' in INSERT (ORA-00928).
func (p *Parser) parsePartitionExtClause(emptyAtOpen bool) (*nodes.PartitionExtClause, error) {
	start := p.pos()
	pe := &nodes.PartitionExtClause{Loc: nodes.Loc{Start: start}}

	if p.cur.Type == kwSUBPARTITION {
		pe.IsSubpartition = true
	}
	p.advance() // consume PARTITION or SUBPARTITION

	// FOR keyword
	if p.cur.Type == kwFOR {
		pe.IsFor = true
		p.advance()
	}

	if p.cur.Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}
	open := p.cur
	p.advance()
	switch {
	case pe.IsFor:
		var err error
		if pe.Keys, err = p.parseExprList(); err != nil {
			return nil, err
		}
		if pe.Keys == nil || pe.Keys.Len() == 0 {
			return nil, p.syntaxErrorAtCur()
		}
	case p.cur.Type == ')' && emptyAtOpen:
		return nil, p.syntaxErrorAtTok(open)
	case p.cur.Type == tokICONST:
		// Oracle 23ai parses a number as the name and looks it up later.
		pe.Name = p.cur.Str
		p.advance()
	default:
		if !p.isIdentLike() {
			return nil, p.syntaxErrorAtCur()
		}
		var err error
		if pe.Name, err = p.parseIdentifier(); err != nil {
			return nil, err
		}
	}
	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()

	pe.Loc.End = p.prev.End
	return pe, nil
}

// parseTableCollectionExpr parses TABLE(collection_expression) [(+)].
//
// BNF: oracle/parser/bnf/SELECT.bnf
//
//	table_collection_expression:
//	    TABLE ( collection_expression ) [ ( + ) ]
func (p *Parser) parseTableCollectionExpr(start int) (nodes.TableExpr, error) {
	p.advance() // consume TABLE
	tc := &nodes.TableCollectionExpr{Loc: nodes.Loc{Start: start}}

	if p.cur.Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	// The collection expression may be a subquery, as it is where a DML
	// statement changes a nested table.
	if p.atQueryStart() {
		subStart := p.pos()
		sub, err := p.parseSelectStmt()
		if err != nil {
			return nil, err
		}
		tc.Expr = &nodes.SubqueryExpr{Subquery: sub, Loc: nodes.Loc{Start: subStart, End: p.prev.End}}
	} else {
		var err error
		if tc.Expr, err = p.parseExpr(); err != nil {
			return nil, err
		}
		if tc.Expr == nil {
			return nil, p.syntaxErrorAtCur()
		}
	}
	if fn, ok := tc.Expr.(*nodes.FuncCallExpr); ok {
		tc.FunctionCall = fn
	} else {
		tc.CollectionExpr = tc.Expr
	}
	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()

	// Optional (+)
	if p.cur.Type == '(' && p.peekNext().Type == '+' {
		p.advance() // (
		p.advance() // +
		if p.cur.Type != ')' {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		tc.OuterJoin = true
	}

	// Optional alias
	if p.cur.Type == kwAS {
		p.advance()
		var parseErr1081 error
		tc.Alias, parseErr1081 = p.parseAlias()
		if parseErr1081 != nil {
			return nil, parseErr1081
		}
	} else if p.isTableAliasCandidate() {
		var parseErr1082 error
		tc.Alias, parseErr1082 = p.parseAlias()
		if parseErr1082 != nil {
			return nil, parseErr1082
		}
	}
	if p.cur.Type == '(' {
		cols, parseErr1083 := p.parseIdentifierListInParens()
		if parseErr1083 != nil {
			return nil, parseErr1083
		}
		tc.ColumnAliases = cols
		if tc.Alias != nil {
			tc.Alias.Cols = cols
		}
	}

	tc.Loc.End = p.prev.End
	return tc, nil
}

// parseContainersOrShards parses CONTAINERS(table) or SHARDS(table).
//
// BNF: oracle/parser/bnf/SELECT.bnf
//
//	containers_clause:
//	    CONTAINERS ( [ schema. ] { table | view } )
//	shards_clause:
//	    SHARDS ( [ schema. ] { table | view } )
func (p *Parser) parseContainersOrShards(start int, isShards bool) (nodes.TableExpr, error) {
	p.advance() // consume CONTAINERS or SHARDS
	ce := &nodes.ContainersExpr{
		IsShards: isShards,
		Loc:      nodes.Loc{Start: start},
	}

	if p.cur.Type == '(' {
		p.advance()
		var parseErr1083 error
		ce.Name, parseErr1083 = p.parseObjectName()
		if parseErr1083 != nil {
			return nil, parseErr1083
		}
		if p.cur.Type == ')' {
			p.advance()
		}
	}

	// Optional alias
	if p.cur.Type == kwAS {
		p.advance()
		var parseErr1084 error
		ce.Alias, parseErr1084 = p.parseAlias()
		if parseErr1084 != nil {
			return nil, parseErr1084
		}
	} else if p.isTableAliasCandidate() {
		var parseErr1085 error
		ce.Alias, parseErr1085 = p.parseAlias()
		if parseErr1085 != nil {
			return nil, parseErr1085
		}
	}

	ce.Loc.End = p.prev.End
	return ce, nil
}
