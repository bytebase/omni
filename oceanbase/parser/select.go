package parser

import (
	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// parseWithClause parses a WITH clause (Common Table Expressions).
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/with.html
//
//	with_clause:
//	    WITH [RECURSIVE]
//	        cte_name [(col_name [, col_name] ...)] AS (subquery)
//	        [, cte_name [(col_name [, col_name] ...)] AS (subquery)] ...
func (p *Parser) parseWithClause() ([]*nodes.CommonTableExpr, error) {
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()
	// Record CTE position for completion context (before consuming WITH).
	if p.completing {
		p.addCTEPosition(p.pos())
	}
	p.advance() // consume WITH

	// Completion: after WITH keyword, offer RECURSIVE and identifier context for CTE name.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwRECURSIVE)
		return nil, &ParseError{Message: "collecting"}
	}

	recursive := false
	if _, ok := p.match(kwRECURSIVE); ok {
		recursive = true
	}

	var ctes []*nodes.CommonTableExpr
	for {
		start := p.pos()
		name, _, err := p.parseIdent()
		if err != nil {
			return nil, err
		}

		cte := &nodes.CommonTableExpr{
			Loc:       nodes.Loc{Start: start},
			Name:      name,
			Recursive: recursive,
		}

		// Optional column list: (col1, col2, ...)
		if p.cur.Type == '(' {
			cols, err := p.parseParenIdentList()
			if err != nil {
				return nil, err
			}
			cte.Columns = cols
		}

		// AS (subquery)
		if _, err := p.expect(kwAS); err != nil {
			return nil, err
		}
		if _, err := p.expect('('); err != nil {
			return nil, err
		}

		// Completion: inside CTE AS (, offer SELECT keyword.
		p.checkCursor()
		if p.collectMode() {
			p.addTokenCandidate(kwSELECT)
			return nil, &ParseError{Message: "collecting"}
		}

		sel, err := p.parseSelectStmt()
		if err != nil {
			return nil, err
		}
		cte.Select = sel

		if _, err := p.expect(')'); err != nil {
			return nil, err
		}

		cte.Loc.End = p.prev.End
		ctes = append(ctes, cte)

		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	return ctes, nil
}

// parseWithStmt parses a statement starting with WITH: WITH ... SELECT, WITH ... UPDATE, WITH ... DELETE.
func (p *Parser) parseWithStmt() (nodes.Node, error) {
	start := p.pos()
	ctes, err := p.parseWithClause()
	if err != nil {
		return nil, err
	}

	switch p.cur.Type {
	case kwSELECT, kwVALUES, '(':
		// A SELECT leaf, a VALUES primary or a parenthesized query body:
		//   WITH ... SELECT ...
		//   WITH ... VALUES ROW(...) ...
		//   WITH ... (SELECT ...) [set-op] [ORDER BY / LIMIT]
		return p.parseSelectClauseAndTrailing(start, ctes)

	case kwINSERT, kwREPLACE:
		isReplace := p.cur.Type == kwREPLACE
		stmt, err := p.parseInsertOrReplace(isReplace)
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = start
		return stmt, nil

	case kwUPDATE:
		stmt, err := p.parseUpdateStmt()
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = start
		return stmt, nil

	case kwDELETE:
		stmt, err := p.parseDeleteStmt()
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = start
		return stmt, nil

	default:
		return nil, &ParseError{
			Message:  "expected SELECT, INSERT, UPDATE, DELETE, or '(' after WITH clause",
			Position: p.cur.Loc,
		}
	}
}

// parseSelectStmt parses a complete query expression as a statement: an
// optional WITH clause, a query-expression body (a SELECT leaf, a parenthesized
// query, or a UNION/INTERSECT/EXCEPT chain), and any trailing ORDER BY / LIMIT.
// It is the entry point every query-bearing context uses (statement level,
// subqueries, INSERT/CTAS/VIEW sources, routine bodies).
func (p *Parser) parseSelectStmt() (*nodes.SelectStmt, error) {
	return p.parseSelectNoParens()
}

// parseSelectNoParens parses a query expression that is not itself wrapped in
// parentheses (PostgreSQL's select_no_parens): [WITH] select_clause
// [ORDER BY] [LIMIT/OFFSET]. The trailing ORDER BY / LIMIT bind at this OUTER
// level, which is what gives a parenthesized inner query and an outer clause
// distinct scopes.
func (p *Parser) parseSelectNoParens() (*nodes.SelectStmt, error) {
	start := p.pos()

	// Optional WITH clause (applies to the whole query expression).
	var ctes []*nodes.CommonTableExpr
	if p.cur.Type == kwWITH {
		var err error
		ctes, err = p.parseWithClause()
		if err != nil {
			return nil, err
		}
	}

	return p.parseSelectClauseAndTrailing(start, ctes)
}

// parseSelectClauseAndTrailing parses a select_clause followed by the trailing
// ORDER BY / LIMIT that applies to a parenthesized query or set-operation
// result. ctes (already parsed by the caller) and start (the position of the
// leading WITH / '(' / SELECT) are attached to the resulting node.
func (p *Parser) parseSelectClauseAndTrailing(start int, ctes []*nodes.CommonTableExpr) (*nodes.SelectStmt, error) {
	body, err := p.parseSelectClause(0)
	if err != nil {
		return nil, err
	}
	if body == nil { // collecting completion candidates
		return nil, nil
	}
	if len(ctes) > 0 {
		body.CTEs = ctes
	}
	body.Loc.Start = start

	// A trailing ORDER BY / LIMIT after a parenthesized query or a set operation
	// (i.e. after a ')') applies to the whole result. A bare SELECT leaf has
	// already consumed its own ORDER BY / LIMIT / locking greedily, so it is
	// skipped here (a second ORDER BY on a leaf is a duplicate the leaf rejects).
	if body.ParenSource != nil || body.SetOp != nodes.SetOpNone {
		if err := p.parseOuterOrderByLimit(body); err != nil {
			return nil, err
		}
	}

	body.Loc.End = p.prev.End
	return body, nil
}

// parseOuterOrderByLimit parses trailing ORDER BY, LIMIT/OFFSET, and locking
// clauses applying to a parenthesized query or set-operation result, onto stmt's
// own fields.
//
// OceanBase takes a locking clause after a set operation, and after a
// parenthesized query only once an ORDER BY or LIMIT follows the ')'. In
// OceanBase CE v4.3.5_CE_BP1 src/sql/parser/sql_parser_mysql_mode.y,
// select_no_parens (line 11016) is select_clause, select_clause_set or
// select_clause_set_with_order_and_limit, each followed by opt_lock_type, and
// a parenthesized query is a select_clause only as
// select_with_parens_with_order_and_limit (line 11259), a select_with_parens
// with an order_by or a limit_clause after it. A lone select_with_parens is a
// select_stmt (line 10969) by itself, with no opt_lock_type, so "(SELECT 1)
// FOR UPDATE" and "((SELECT 1)) FOR UPDATE" are 1064 where MySQL 8.0 takes
// them, while "(SELECT 1) ORDER BY 1 FOR UPDATE" and "(SELECT 1) UNION
// (SELECT 2) FOR UPDATE" parse.
func (p *Parser) parseOuterOrderByLimit(stmt *nodes.SelectStmt) error {
	// A set operation's last operand has read its clauses already, and they have
	// moved here, so a second ORDER BY, LIMIT or locking clause has no place: "SELECT 1
	// UNION SELECT 2 FOR UPDATE FOR UPDATE" is a 1064 (measured).
	if stmt.SetOp != nodes.SetOpNone && (len(stmt.OrderBy) > 0 || stmt.Limit != nil || stmt.ForUpdate != nil || stmt.Into != nil) {
		return nil
	}
	if p.cur.Type == kwORDER {
		p.advance()
		if _, err := p.expect(kwBY); err != nil {
			return err
		}
		orderBy, err := p.parseOrderByList()
		if err != nil {
			return err
		}
		stmt.OrderBy = orderBy
		if p.cur.Type == kwORDER {
			return &ParseError{Message: "duplicate ORDER BY clause", Position: p.cur.Loc}
		}
	}
	if _, ok := p.match(kwLIMIT); ok {
		limit, err := p.parseLimitClause()
		if err != nil {
			return err
		}
		stmt.Limit = limit
		if p.cur.Type == kwLIMIT {
			return &ParseError{Message: "duplicate LIMIT clause", Position: p.cur.Loc}
		}
	}
	if stmt.SetOp == nodes.SetOpNone && len(stmt.OrderBy) == 0 && stmt.Limit == nil {
		// A lone parenthesized query, which takes no locking clause.
		return nil
	}
	return p.parseLockingClause(stmt)
}

// parseLockingClause parses the opt_lock_type of select_no_parens onto stmt, if
// there is one.
func (p *Parser) parseLockingClause(stmt *nodes.SelectStmt) error {
	switch p.cur.Type {
	case kwFOR:
		fu, err := p.parseForUpdateClause()
		if err != nil {
			return err
		}
		stmt.ForUpdate = fu
	case kwLOCK:
		fu, err := p.parseLockInShareMode()
		if err != nil {
			return err
		}
		stmt.ForUpdate = fu
	}
	return nil
}

// operandEndsInClauses reports whether a query expression that a set operator follows
// ends in an ORDER BY, LIMIT or locking clause outside parentheses, which OceanBase
// refuses there: SELECT 1 ORDER BY 1 UNION SELECT 2 and VALUES ROW(1) FOR UPDATE UNION
// VALUES ROW(2) are each a 1064 (measured), though the grammar's
// select_clause_set_left_body reads as allowing the first. Those of a set operation's
// last operand have moved to the set operation by now, and a parenthesized query has
// read none of its own yet.
func operandEndsInClauses(sel *nodes.SelectStmt) bool {
	return len(sel.OrderBy) > 0 || sel.Limit != nil || sel.ForUpdate != nil
}

// parseQueryStatement parses a query expression where the grammar reads
// select_stmt: as a statement, an EXPLAIN target, a view or CREATE TABLE body, or an
// INSERT source. There select_into adds an INTO after select_no_parens, which a lone
// parenthesized query is not, and with_select has none of its own: a SELECT reads its
// own INTO as a clause, so what comes here is an INTO after VALUES, a set operation
// or a closing ')'.
func (p *Parser) parseQueryStatement() (*nodes.SelectStmt, error) {
	sel, err := p.parseSelectStmt()
	if err != nil || sel == nil || p.cur.Type != kwINTO {
		return sel, err
	}
	if len(sel.CTEs) > 0 || hasInto(sel) {
		return sel, nil
	}
	if sel.ParenSource != nil && len(sel.OrderBy) == 0 && sel.Limit == nil {
		return sel, nil
	}
	p.advance() // consume INTO
	into, err := p.parseIntoClause()
	if err != nil {
		return nil, err
	}
	sel.Into = into
	sel.Loc.End = p.prev.End
	return sel, nil
}

// hasInto reports whether a query expression already holds an INTO, on itself or on
// the last operand or body it ends in.
func hasInto(sel *nodes.SelectStmt) bool {
	for sel != nil {
		if sel.Into != nil {
			return true
		}
		switch {
		case sel.SetOp != nodes.SetOpNone:
			sel = sel.Right
		case sel.ParenSource != nil:
			sel = sel.ParenSource
		default:
			return false
		}
	}
	return false
}

// parseSelectClause parses a select_clause: a query primary optionally combined
// with others via UNION / INTERSECT / EXCEPT, using precedence climbing
// (INTERSECT binds tighter than UNION/EXCEPT). minPrec is the minimum operator
// precedence this call will consume; the top-level no_parens call passes 0.
func (p *Parser) parseSelectClause(minPrec int) (*nodes.SelectStmt, error) {
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()
	left, err := p.parseSelectClausePrimary()
	if err != nil {
		return nil, err
	}
	if left == nil { // collecting completion candidates
		return nil, nil
	}

	for {
		prec := setOpPrecedence(p.cur.Type)
		if prec == 0 || prec < minPrec {
			return left, nil
		}
		if operandEndsInClauses(left) {
			return nil, p.syntaxErrorAtCur()
		}

		var op nodes.SetOperation
		switch p.cur.Type {
		case kwUNION:
			op = nodes.SetOpUnion
		case kwEXCEPT:
			op = nodes.SetOpExcept
		case kwINTERSECT:
			op = nodes.SetOpIntersect
		}
		p.advance() // consume UNION/INTERSECT/EXCEPT

		// Completion: after the set-op keyword, offer ALL, DISTINCT, SELECT.
		p.checkCursor()
		if p.collectMode() {
			p.addTokenCandidate(kwALL)
			p.addTokenCandidate(kwDISTINCT)
			p.addTokenCandidate(kwSELECT)
			return nil, &ParseError{Message: "collecting"}
		}

		all := false
		if _, ok := p.match(kwALL); ok {
			all = true
		} else {
			p.match(kwDISTINCT)
		}

		// Completion: after ALL/DISTINCT, offer SELECT.
		p.checkCursor()
		if p.collectMode() {
			p.addTokenCandidate(kwSELECT)
			return nil, &ParseError{Message: "collecting"}
		}

		// Right operand binds higher-precedence operators first; prec+1 keeps
		// same-precedence chains left-associative.
		right, err := p.parseSelectClause(prec + 1)
		if err != nil {
			return nil, err
		}

		// Hoist the clauses after an UNPARENTHESIZED last operand onto the
		// set-operation node, which is what they apply to, so they stay in the order
		// they were written: an ORDER BY, LIMIT, locking clause and INTO are read
		// together by select_no_parens. A parenthesized operand is a ParenSource
		// wrapper whose own clauses are the OUTER scope, so its inner clauses stay in
		// place — e.g. "SELECT 1 UNION (SELECT 2 LIMIT 1)" keeps LIMIT 1 on operand 2. An
		// operand that is itself a set operation, as the right of a UNION can be an
		// INTERSECT, has already gathered its clauses at its own root.
		var orderBy []*nodes.OrderByItem
		var limit *nodes.Limit
		var forUpdate *nodes.ForUpdate
		var into *nodes.IntoClause
		if right.ParenSource == nil {
			orderBy, right.OrderBy = right.OrderBy, nil
			limit, right.Limit = right.Limit, nil
			forUpdate, right.ForUpdate = right.ForUpdate, nil
			into, right.Into = right.Into, nil
		}

		left = &nodes.SelectStmt{
			Loc:       nodes.Loc{Start: left.Loc.Start, End: right.Loc.End},
			SetOp:     op,
			SetAll:    all,
			Left:      left,
			Right:     right,
			OrderBy:   orderBy,
			Limit:     limit,
			ForUpdate: forUpdate,
			Into:      into,
		}
	}
}

// parseSelectClausePrimary parses a single query primary: a parenthesized query
// (select_with_parens), a simple SELECT leaf, or VALUES.
//
// OceanBase has no TABLE query. Its select_clause (OceanBase CE
// v4.3.5_CE_BP1 src/sql/parser/sql_parser_mysql_mode.y:11088) and
// select_clause_set_body (11162) take a SELECT, a parenthesized query and
// table_values_clause (14920), which is VALUES values_row_list, and nothing
// else in the grammar reads TABLE as a query, so "TABLE t" is 1064 wherever a
// query goes: as a statement, a subquery, a CTE, a set-operation operand, or
// the source of INSERT, CREATE TABLE or CREATE VIEW.
func (p *Parser) parseSelectClausePrimary() (*nodes.SelectStmt, error) {
	if p.cur.Type == '(' {
		return p.parseSelectWithParens()
	}
	if p.cur.Type == kwSELECT {
		return p.parseSimpleSelect()
	}
	// A VALUES primary becomes a SelectStmt that holds the source and is the
	// only owner of the clauses after it: the ORDER BY and LIMIT of
	// table_values_clause_with_order_by_and_limit (line 14931), and the
	// opt_lock_type select_no_parens puts after any select_clause (line 11016).
	// Had the source kept them too, every walk would visit them twice. An INTO
	// is select_into's, so only a query standing as a statement takes one (see
	// parseQueryStatement).
	if p.cur.Type == kwVALUES {
		values, err := p.parseValuesConstructor()
		if err != nil {
			return nil, err
		}
		stmt := &nodes.SelectStmt{Loc: values.Loc, ValuesSource: values}
		if stmt.OrderBy, stmt.Limit, err = p.parseSourceClauses(); err != nil {
			return nil, err
		}
		if err := p.parseLockingClause(stmt); err != nil {
			return nil, err
		}
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}
	if p.cur.Type == tokEOF {
		return nil, p.syntaxErrorAtCur()
	}
	return nil, &ParseError{
		Message:  "expected SELECT, VALUES, or '('",
		Position: p.cur.Loc,
	}
}

func isQueryPrimaryStartToken(tokenType int) bool {
	switch tokenType {
	case kwSELECT, kwWITH, kwVALUES:
		return true
	default:
		return false
	}
}

func (p *Parser) isQueryExpressionStart() bool {
	if isQueryPrimaryStartToken(p.cur.Type) {
		return true
	}
	if p.cur.Type != '(' {
		return false
	}
	// A leading run of '(' may open either a parenthesized query expression
	// — `((SELECT ...))`, `((SELECT 1) UNION (SELECT 2))` — or a parenthesized
	// scalar expression — `((a = b))`, the redundant-paren form MySQL emits for
	// a view's join ON condition (SHOW CREATE VIEW renders `on(((a) and (b)))`).
	// A fixed one-token lookahead cannot tell them apart; scan past the run of
	// '(' to the first non-'(' token — a query primary (SELECT/WITH/VALUES)
	// means a query expression, anything else means a scalar expression.
	return isQueryPrimaryStartToken(p.firstTokenPastParens().Type)
}

// scanState captures the full lexical scan position so a speculative lookahead
// can rewind without consuming input. The Lexer is a flat value (string + ints),
// so a shallow copy snapshots its position; the buffered next token and the
// current/previous tokens complete the parser's scan cursor. The completion
// `collecting` latch is position-derived (it arms when the scan reaches the
// cursor offset), so rewinding the scan rewinds it too — a speculative parse
// that crossed the cursor must not leave collect mode armed for the reparse.
type scanState struct {
	lexer      Lexer
	cur        Token
	prev       Token
	nextBuf    Token
	hasNext    bool
	collecting bool
	// Lengths of the CandidateSet's position slices at save time. Unlike the
	// set-valued token/rule candidates, these are plain appends, so a rewound
	// speculative parse would duplicate every offset it recorded; restore
	// truncates back to the snapshot.
	ctePosLen   int
	aliasPosLen int
}

func (p *Parser) saveScan() scanState {
	s := scanState{lexer: *p.lexer, cur: p.cur, prev: p.prev, nextBuf: p.nextBuf, hasNext: p.hasNext, collecting: p.collecting}
	if p.candidates != nil {
		s.ctePosLen = len(p.candidates.CTEPositions)
		s.aliasPosLen = len(p.candidates.SelectAliasPositions)
	}
	return s
}

func (p *Parser) restoreScan(s scanState) {
	*p.lexer = s.lexer
	p.cur = s.cur
	p.prev = s.prev
	p.nextBuf = s.nextBuf
	p.hasNext = s.hasNext
	p.collecting = s.collecting
	if p.candidates != nil {
		if len(p.candidates.CTEPositions) > s.ctePosLen {
			p.candidates.CTEPositions = p.candidates.CTEPositions[:s.ctePosLen]
		}
		if len(p.candidates.SelectAliasPositions) > s.aliasPosLen {
			p.candidates.SelectAliasPositions = p.candidates.SelectAliasPositions[:s.aliasPosLen]
		}
	}
}

// firstTokenPastParens returns the first token at or after the current position
// that is not an opening parenthesis, without consuming input. It is the
// disambiguator for constructs that begin with a run of '(' whose meaning turns
// on the innermost head:
//   - a FROM-clause table factor — a (possibly nested) derived-table subquery
//     `((SELECT ...))` vs. a parenthesized join group `((a JOIN b) JOIN c)`;
//   - an expression '(' — a parenthesized query expression `((SELECT ...))` vs.
//     a parenthesized scalar expression `((a = b))`.
//
// In each case the grammar distinguishes by whether the innermost head is a
// query primary (SELECT/WITH/VALUES) or not; a fixed one-token lookahead
// on '(' cannot, so we scan past the parens speculatively and rewind.
//
// Every '(' of a run shares the run's answer, and the parser asks again at each
// nesting level as it descends, so the last scan is cached: a '(' starting
// between the run's first '(' and the token past it is in that run. Without
// the cache, n nested parens cost O(n²) token scans.
func (p *Parser) firstTokenPastParens() Token {
	if p.cur.Type == '(' && p.cur.Loc >= p.parenRunStart && p.cur.Loc < p.parenRunPast.Loc {
		return p.parenRunPast
	}
	saved := p.saveScan()
	defer p.restoreScan(saved)
	start := p.cur.Loc
	for p.cur.Type == '(' {
		p.advance()
	}
	p.parenRunStart, p.parenRunPast = start, p.cur
	return p.cur
}

// parseSelectWithParens parses a parenthesized query expression (PostgreSQL's
// select_with_parens): '(' select_no_parens ')'. It returns a WRAPPER node
// whose ParenSource is the inner query; the wrapper's own ORDER BY / LIMIT hold
// any OUTER trailing clauses attached by the enclosing no_parens. Unlike
// PostgreSQL (which unwraps the inner query and would overwrite its clause),
// MySQL accepts and must preserve BOTH scopes — e.g. "(SELECT 1 LIMIT 5) LIMIT
// 2" — so the parens materialize as a real node. The unconditional delegation to
// parseSelectNoParens left-factors the grammar so a nested "((...))" and a
// parenthesized set-op operand are parsed by the same path.
//
// The result depends only on where the parse starts and how deeply it is
// nested, so a failure is recorded (see parenQueryFail) and returned again
// without reparsing. Completion mode records its failures with the effects a
// reparse would repeat (see parenQueryReplay).
//
// A record made at nesting depth d stands for a parse that passed every
// enterNesting it reached from d; started no deeper, a reparse passes them
// all again and fails the same way, so the record is replayed only then.
// Started deeper, a reparse may reach the nesting limit first, so it parses
// again. A record never holds a nesting error, which depends on the depth the
// parse started at.
func (p *Parser) parseSelectWithParens() (*nodes.SelectStmt, error) {
	start := p.pos()
	if r, ok := p.parenQueryFail[start]; ok && p.depth <= r.depth {
		return nil, r.err
	}
	if p.completing {
		return p.parseSelectWithParensCompleting(start)
	}
	depth := p.depth
	sel, err := p.parseSelectWithParensAt(start)
	if err != nil && !isNestingError(err) {
		if p.parenQueryFail == nil {
			p.parenQueryFail = make(map[int]parenQueryFailRecord)
		}
		p.parenQueryFail[start] = parenQueryFailRecord{err: err, depth: depth}
	}
	return sel, err
}

// parenQueryFailRecord is a parseSelectWithParens failure outside completion
// mode and the nesting depth the failed parse started at.
type parenQueryFailRecord struct {
	err   error
	depth int
}

// parenQueryKey identifies a completion-mode parse of a parenthesized query
// expression: the offset of its '(' and whether collection was armed when it
// began, the one piece of parser state besides the input that steers it.
type parenQueryKey struct {
	start      int
	collecting bool
}

// parenQueryFailure is a completion-mode parseSelectWithParens failure and
// what reparsing from the same place would do again: stop the scan at the
// same token, and record the same CTE and select-alias positions. The token
// and rule candidates it collected are in their sets already, which a reparse
// would leave unchanged.
type parenQueryFailure struct {
	err      error
	end      scanState
	ctePos   []int
	aliasPos []int
	depth    int // the nesting depth the failed parse started at
}

// parseSelectWithParensCompleting is parseSelectWithParens in completion
// mode. Returning a failure again without its effects would lose candidates,
// so a replay restores them instead of reparsing, which keeps the speculation
// at every level of `((((SELECT 1) = 0) = 0) = 0)` linear: reparsing took 8 s
// at 10,000 levels. As outside completion mode, a record is replayed only at
// its own nesting depth or shallower, and a failure at the nesting limit is
// not recorded.
func (p *Parser) parseSelectWithParensCompleting(start int) (*nodes.SelectStmt, error) {
	key := parenQueryKey{start: start, collecting: p.collecting}
	if f, ok := p.parenQueryReplay[key]; ok && p.depth <= f.depth {
		e := f.end
		*p.lexer = e.lexer
		p.cur, p.prev, p.nextBuf, p.hasNext, p.collecting = e.cur, e.prev, e.nextBuf, e.hasNext, e.collecting
		if p.candidates != nil {
			p.candidates.CTEPositions = append(p.candidates.CTEPositions, f.ctePos...)
			p.candidates.SelectAliasPositions = append(p.candidates.SelectAliasPositions, f.aliasPos...)
		}
		return nil, f.err
	}
	depth := p.depth
	before := p.saveScan()
	sel, err := p.parseSelectWithParensAt(start)
	if err == nil {
		return sel, nil
	}
	if isNestingError(err) {
		return nil, err
	}
	f := &parenQueryFailure{err: err, end: p.saveScan(), depth: depth}
	if cs := p.candidates; cs != nil {
		if len(cs.CTEPositions) > before.ctePosLen {
			f.ctePos = append([]int(nil), cs.CTEPositions[before.ctePosLen:]...)
		}
		if len(cs.SelectAliasPositions) > before.aliasPosLen {
			f.aliasPos = append([]int(nil), cs.SelectAliasPositions[before.aliasPosLen:]...)
		}
	}
	if p.parenQueryReplay == nil {
		p.parenQueryReplay = make(map[parenQueryKey]*parenQueryFailure)
	}
	p.parenQueryReplay[key] = f
	return nil, err
}

func (p *Parser) parseSelectWithParensAt(start int) (*nodes.SelectStmt, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	inner, err := p.parseSelectNoParens()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.SelectStmt{
		Loc:         nodes.Loc{Start: start, End: p.prev.End},
		ParenSource: inner,
	}, nil
}

// parseSimpleSelect parses a single SELECT query primary (no leading WITH, no
// trailing set operation). It greedily consumes its own ORDER BY / LIMIT /
// locking; the enclosing select_clause hoists an unparenthesized last operand's
// ORDER BY / LIMIT onto the set-op node.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/select.html
//
//	SELECT
//	    [ALL | DISTINCT | DISTINCTROW] [HIGH_PRIORITY] [STRAIGHT_JOIN]
//	    [SQL_CALC_FOUND_ROWS] select_expr [, select_expr] ...
//	    [FROM table_references] [WHERE where_condition]
//	    [GROUP BY ... [WITH ROLLUP]] [HAVING ...] [WINDOW ...]
//	    [ORDER BY ...] [LIMIT {[offset,] row_count | row_count OFFSET offset}]
//	    [FOR {UPDATE | SHARE} ...] [INTO OUTFILE 'file_name']
func (p *Parser) parseSimpleSelect() (*nodes.SelectStmt, error) {
	start := p.pos()

	p.advance() // consume SELECT

	// Completion: after SELECT keyword, offer select-expression candidates.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwDISTINCT)
		p.addTokenCandidate(kwALL)
		p.addRuleCandidate("columnref")
		p.addRuleCandidate("func_name")
		return nil, &ParseError{Message: "collecting"}
	}

	stmt := &nodes.SelectStmt{Loc: nodes.Loc{Start: start}}

	// Parse SELECT options
	for {
		switch p.cur.Type {
		case kwALL:
			stmt.DistinctKind = nodes.DistinctAll
			p.advance()
			continue
		case kwDISTINCT, kwDISTINCTROW:
			stmt.DistinctKind = nodes.DistinctOn
			p.advance()
			continue
		case kwHIGH_PRIORITY:
			stmt.HighPriority = true
			p.advance()
			continue
		case kwSTRAIGHT_JOIN:
			stmt.StraightJoin = true
			p.advance()
			continue
		case kwSQL_CALC_FOUND_ROWS:
			stmt.CalcFoundRows = true
			p.advance()
			continue
		case kwSQL_SMALL_RESULT:
			stmt.SmallResult = true
			p.advance()
			continue
		case kwSQL_BIG_RESULT:
			stmt.BigResult = true
			p.advance()
			continue
		case kwSQL_BUFFER_RESULT:
			stmt.BufferResult = true
			p.advance()
			continue
		case kwSQL_NO_CACHE:
			stmt.NoCache = true
			p.advance()
			continue
		}
		break
	}

	// Parse select expression list
	targets, err := p.parseSelectExprList()
	if err != nil {
		return nil, err
	}
	stmt.TargetList = targets

	// Completion: after select target list, offer clause keywords.
	p.checkCursor()
	if p.collectMode() {
		for _, tok := range []int{kwINTO, kwFROM, kwWHERE, kwGROUP, kwHAVING, kwORDER, kwLIMIT, kwFOR} {
			p.addTokenCandidate(tok)
		}
		return nil, &ParseError{Message: "collecting"}
	}

	// INTO clause (position 1: after select_expr, before FROM)
	if p.cur.Type == kwINTO {
		p.advance()
		into, err := p.parseIntoClause()
		if err != nil {
			return nil, err
		}
		stmt.Into = into
	}

	// FROM clause
	if _, ok := p.match(kwFROM); ok {
		from, err := p.parseTableReferenceList()
		if err != nil {
			return nil, err
		}
		stmt.From = from
	}

	// WHERE clause — MySQL allows WHERE without FROM (e.g., SELECT 1 WHERE 1=1)
	if p.cur.Type == kwWHERE {
		p.advance()
		where, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if !p.collectMode() && where == nil {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Where = where

		// Reject duplicate WHERE
		if p.cur.Type == kwWHERE {
			return nil, &ParseError{
				Message:  "duplicate WHERE clause",
				Position: p.cur.Loc,
			}
		}
	}

	// GROUP BY clause — MySQL allows GROUP BY without FROM
	if p.cur.Type == kwGROUP {
		p.advance()
		if _, err := p.expect(kwBY); err != nil {
			return nil, err
		}
		groupBy, err := p.parseExprList()
		if err != nil {
			return nil, err
		}
		stmt.GroupBy = groupBy

		// Completion: after GROUP BY list, offer follow-set keywords.
		p.checkCursor()
		if p.collectMode() {
			p.addTokenCandidate(kwHAVING)
			p.addTokenCandidate(kwORDER)
			p.addTokenCandidate(kwLIMIT)
			p.addTokenCandidate(kwWITH)
			return nil, &ParseError{Message: "collecting"}
		}

		// WITH ROLLUP
		if p.cur.Type == kwWITH {
			p.advance()
			if _, ok := p.match(kwROLLUP); ok {
				stmt.WithRollup = true
			}
		}

		// Reject duplicate GROUP BY
		if p.cur.Type == kwGROUP {
			return nil, &ParseError{
				Message:  "duplicate GROUP BY clause",
				Position: p.cur.Loc,
			}
		}
	}

	// HAVING clause — MySQL allows HAVING without FROM
	if p.cur.Type == kwHAVING {
		p.advance()
		having, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if !p.collectMode() && having == nil {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Having = having
	}

	// WINDOW clause: WINDOW window_name AS (window_spec) [, ...]
	if p.cur.Type == kwWINDOW {
		p.advance()
		defs, err := p.parseNamedWindowList()
		if err != nil {
			return nil, err
		}
		stmt.WindowClause = defs
	}

	// ORDER BY clause
	if p.cur.Type == kwORDER {
		p.advance()
		if _, err := p.expect(kwBY); err != nil {
			return nil, err
		}
		orderBy, err := p.parseOrderByList()
		if err != nil {
			return nil, err
		}
		stmt.OrderBy = orderBy

		// WITH ROLLUP (MySQL 8.0.12+)
		if p.cur.Type == kwWITH {
			next := p.peekNext()
			if next.Type == kwROLLUP {
				p.advance() // consume WITH
				p.advance() // consume ROLLUP
				stmt.OrderByWithRollup = true
			}
		}

		// Reject duplicate ORDER BY
		if p.cur.Type == kwORDER {
			return nil, &ParseError{
				Message:  "duplicate ORDER BY clause",
				Position: p.cur.Loc,
			}
		}
	}

	// LIMIT clause
	if _, ok := p.match(kwLIMIT); ok {
		limit, err := p.parseLimitClause()
		if err != nil {
			return nil, err
		}
		stmt.Limit = limit
	}

	// FOR UPDATE / FOR SHARE / LOCK IN SHARE MODE
	if p.cur.Type == kwFOR {
		fu, err := p.parseForUpdateClause()
		if err != nil {
			return nil, err
		}
		stmt.ForUpdate = fu
	} else if p.cur.Type == kwLOCK {
		fu, err := p.parseLockInShareMode()
		if err != nil {
			return nil, err
		}
		stmt.ForUpdate = fu
	}

	// INTO clause (position 2: last, after the locking clause).
	//
	// OceanBase takes INTO only here, after every other clause, so
	// "SELECT * FROM t INTO OUTFILE 'f' ORDER BY a",
	// "... INTO OUTFILE 'f' LIMIT 1" and "... INTO OUTFILE 'f' FOR UPDATE"
	// are syntax errors (container-verified on
	// oceanbase/oceanbase-ce:4.3.5-lts). An INTO before FROM (position 1
	// above) is a separate rule.
	if p.cur.Type == kwINTO && stmt.Into == nil {
		p.advance()
		into, err := p.parseIntoClause()
		if err != nil {
			return nil, err
		}
		stmt.Into = into
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseSelectExprList parses a comma-separated list of select expressions.
func (p *Parser) parseSelectExprList() ([]nodes.ExprNode, error) {
	var list []nodes.ExprNode

	for {
		// Completion: at the start of each select expression, offer column/function candidates.
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addRuleCandidate("func_name")
			return nil, &ParseError{Message: "collecting"}
		}

		// Soft-fail: if the current token is a clause keyword (e.g. FROM in
		// "SELECT FROM t"), return what we have instead of parsing it as an
		// identifier. MySQL keywords are non-reserved and would otherwise be
		// consumed as column names, preventing completion on subsequent clauses.
		if p.completing && p.isClauseKeyword() {
			break
		}

		target, err := p.parseSelectExpr()
		if err != nil {
			return nil, err
		}
		if target == nil {
			break // soft-fail: no expression found, return what we have
		}
		list = append(list, target)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	return list, nil
}

// parseSelectExpr parses a single select expression (expr [AS alias]).
func (p *Parser) parseSelectExpr() (nodes.ExprNode, error) {
	start := p.pos()

	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if expr == nil {
		return nil, nil // soft-fail: no expression found (e.g. SELECT FROM)
	}

	// Check for AS alias or implicit alias
	var alias string
	if _, ok := p.match(kwAS); ok {
		aliasLoc := p.pos()
		name, _, err := p.parseIdentOrText()
		if err != nil {
			return nil, err
		}
		alias = name
		if p.completing {
			p.addSelectAliasPosition(aliasLoc)
		}
	} else if (p.isIdentToken() || p.cur.Type == tokSCONST) && !p.isSelectTerminator() {
		// Implicit alias (identifier or string literal without AS), but not if it's a keyword that starts the next clause
		aliasLoc := p.pos()
		alias, _, err = p.parseIdentOrText()
		if err != nil {
			return nil, err
		}
		if p.completing {
			p.addSelectAliasPosition(aliasLoc)
		}
	}

	if alias != "" {
		return &nodes.ResTarget{
			Loc:  nodes.Loc{Start: start, End: p.prev.End},
			Name: alias,
			Val:  expr,
		}, nil
	}

	return expr, nil
}

// parseDerivedColumnList parses an optional parenthesized column alias list
// for derived tables: (col1, col2, ...).
// MySQL 8.0.14+ syntax: (SELECT ...) AS alias(c1, c2).
func (p *Parser) parseDerivedColumnList() ([]string, error) {
	p.advance() // consume '('
	var cols []string
	for {
		name, _, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		cols = append(cols, name)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return cols, nil
}

// isClauseKeyword returns true if the current token is a core SQL clause keyword.
// In MySQL, these keywords are non-reserved (can be used as identifiers), so the
// bottom-level expression parser won't soft-fail on them. This check is needed in
// completion mode to prevent clause keywords from being consumed as column names
// when the SELECT expression list is empty (e.g. "SELECT FROM t").
func (p *Parser) isClauseKeyword() bool {
	switch p.cur.Type {
	case kwFROM, kwWHERE, kwGROUP, kwHAVING, kwORDER, kwLIMIT, kwFOR, kwINTO,
		kwUNION, kwINTERSECT, kwEXCEPT:
		return true
	}
	return false
}

// isSelectTerminator returns true if the current token starts a clause that terminates
// the select expression list. Used to avoid consuming clause keywords as aliases.
func (p *Parser) isSelectTerminator() bool {
	switch p.cur.Type {
	case kwFROM, kwWHERE, kwGROUP, kwHAVING, kwORDER, kwLIMIT, kwFOR, kwINTO,
		kwUNION, kwINTERSECT, kwEXCEPT, kwON, kwUSING, kwJOIN, kwINNER, kwLEFT,
		kwRIGHT, kwCROSS, kwNATURAL, kwWINDOW, kwWITH, kwLOCK, ';', tokEOF:
		return true
	}
	return false
}

// parseTableReferenceList parses a comma-separated list of table references (with joins).
func (p *Parser) parseTableReferenceList() ([]nodes.TableExpr, error) {
	var refs []nodes.TableExpr

	for {
		// Completion: at the start of each table reference, offer table_ref candidates.
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("table_ref")
			p.addRuleCandidate("database_ref")
			return nil, &ParseError{Message: "collecting"}
		}

		ref, err := p.parseTableReference()
		if err != nil {
			return nil, err
		}
		refs = append(refs, ref)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	return refs, nil
}

func foldTableReferenceList(refs []nodes.TableExpr) nodes.TableExpr {
	if len(refs) == 0 {
		return nil
	}
	left := refs[0]
	for i := 1; i < len(refs); i++ {
		left = &nodes.JoinClause{
			Loc:   nodes.Loc{Start: tableExprLoc(left).Start, End: tableExprLoc(refs[i]).End},
			Type:  nodes.JoinInner,
			Left:  left,
			Right: refs[i],
		}
	}
	return left
}

// parseTableReference parses a table reference (table_factor with optional joins).
func (p *Parser) parseTableReference() (nodes.TableExpr, error) {
	left, err := p.parseTableFactor()
	if err != nil {
		return nil, err
	}

	// Completion: after a table factor, offer join keywords and post-FROM keywords.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwWHERE)
		p.addTokenCandidate(kwJOIN)
		p.addTokenCandidate(kwLEFT)
		p.addTokenCandidate(kwRIGHT)
		p.addTokenCandidate(kwCROSS)
		p.addTokenCandidate(kwNATURAL)
		p.addTokenCandidate(kwINNER)
		p.addTokenCandidate(kwSTRAIGHT_JOIN)
		p.addTokenCandidate(kwORDER)
		p.addTokenCandidate(kwGROUP)
		p.addTokenCandidate(kwHAVING)
		p.addTokenCandidate(kwLIMIT)
		p.addTokenCandidate(kwUNION)
		p.addTokenCandidate(kwFOR)
		return nil, &ParseError{Message: "collecting"}
	}

	// Parse joins
	for {
		jt, ok, err := p.matchJoinType()
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}

		// Completion: after JOIN keyword, offer table_ref candidates for the right-side table.
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("table_ref")
			p.addRuleCandidate("database_ref")
			return nil, &ParseError{Message: "collecting"}
		}

		right, err := p.parseTableFactor()
		if err != nil {
			return nil, err
		}

		join := &nodes.JoinClause{
			Loc:   nodes.Loc{Start: tableExprLoc(left).Start},
			Type:  jt,
			Left:  left,
			Right: right,
		}

		// ON or USING condition
		condStart := p.pos()
		if _, ok := p.match(kwON); ok {
			// Completion: after ON keyword, offer columnref candidates for join condition.
			p.checkCursor()
			if p.collectMode() {
				p.addRuleCandidate("columnref")
				return nil, &ParseError{Message: "collecting"}
			}
			condStart = p.pos()
			cond, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if cond == nil {
				return nil, p.syntaxErrorAtCur()
			}
			join.Condition = &nodes.OnCondition{
				Loc:  nodes.Loc{Start: condStart, End: p.prev.End},
				Expr: cond,
			}
		} else if _, ok := p.match(kwUSING); ok {
			// For USING, we manually consume '(' so we can add a completion checkpoint inside.
			if _, err := p.expect('('); err != nil {
				return nil, err
			}
			// Completion: after USING (, offer columnref candidates for shared columns.
			p.checkCursor()
			if p.collectMode() {
				p.addRuleCandidate("columnref")
				return nil, &ParseError{Message: "collecting"}
			}
			var cols []string
			for {
				name, _, err := p.parseIdent()
				if err != nil {
					return nil, err
				}
				cols = append(cols, name)
				if p.cur.Type != ',' {
					break
				}
				p.advance()
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
			join.Condition = &nodes.UsingCondition{Loc: nodes.Loc{Start: condStart, End: p.prev.End}, Columns: cols}
		}

		join.Loc.End = p.prev.End
		left = join
	}

	return left, nil
}

// matchJoinType checks for a join keyword combination and returns the join type.
// Returns an error if a join modifier keyword (INNER, LEFT, etc.) is found but
// the required JOIN keyword is missing.
func (p *Parser) matchJoinType() (nodes.JoinType, bool, error) {
	switch p.cur.Type {
	case kwJOIN:
		p.advance()
		return nodes.JoinInner, true, nil
	case kwINNER:
		p.advance()
		if _, err := p.expect(kwJOIN); err != nil {
			return 0, false, err
		}
		return nodes.JoinInner, true, nil
	case kwLEFT:
		p.advance()
		p.match(kwOUTER)
		if _, err := p.expect(kwJOIN); err != nil {
			return 0, false, err
		}
		return nodes.JoinLeft, true, nil
	case kwRIGHT:
		p.advance()
		p.match(kwOUTER)
		if _, err := p.expect(kwJOIN); err != nil {
			return 0, false, err
		}
		return nodes.JoinRight, true, nil
	case kwCROSS:
		p.advance()
		if _, err := p.expect(kwJOIN); err != nil {
			return 0, false, err
		}
		return nodes.JoinCross, true, nil
	case kwNATURAL:
		p.advance()
		// NATURAL [LEFT|RIGHT] [OUTER] JOIN
		if _, ok := p.match(kwLEFT); ok {
			p.match(kwOUTER)
			if _, err := p.expect(kwJOIN); err != nil {
				return 0, false, err
			}
			return nodes.JoinNaturalLeft, true, nil
		}
		if _, ok := p.match(kwRIGHT); ok {
			p.match(kwOUTER)
			if _, err := p.expect(kwJOIN); err != nil {
				return 0, false, err
			}
			return nodes.JoinNaturalRight, true, nil
		}
		if _, err := p.expect(kwJOIN); err != nil {
			return 0, false, err
		}
		return nodes.JoinNatural, true, nil
	case kwSTRAIGHT_JOIN:
		p.advance()
		return nodes.JoinStraight, true, nil
	}
	return 0, false, nil
}

// parenOpensQuery reports whether the '(' at p.cur opens a parenthesized
// query expression, as opposed to a join group or a parenthesized expression,
// without consuming input. It does when its content begins with a query
// primary, or begins with a block that opens a query and is followed by a
// token that continues a query expression: the ')', a set operator, ORDER BY,
// LIMIT, or a locking clause. Any other token past that inner block makes it
// the first operand of something else: a derived table in a join group,
// `((SELECT 1) d JOIN t ON ...)`, as MySQL writes a view over a derived table
// joined to a table.
//
// MySQL 8.0.46 sql/sql_yacc.yy: a derived_table (line 12072) is a
// query_expression_parens (9965) with its alias, and joined_table_parens
// (12067) and table_reference_list_parens (12045) may begin with one; inside
// query_expression_parens, a parenthesized operand is followed by UNION,
// EXCEPT, INTERSECT (query_expression_body, 9930), ORDER BY, LIMIT, a
// locking clause (query_expression_with_opt_locking_clauses, 13299), or the
// ')'.
//
// A block's verdict depends only on the run of '(' it starts and the token
// after each of their blocks, so the verdicts of a run are recorded together,
// and a scan records the token after every block it passes: asking again at
// each level as the parser descends costs nothing.
func (p *Parser) parenOpensQuery() bool {
	if p.cur.Type != '(' {
		return false
	}
	if v, ok := p.parenQuery[p.cur.Loc]; ok {
		return v
	}
	if !p.isQueryExpressionStart() {
		return false
	}
	saved := p.saveScan()
	defer p.restoreScan(saved)
	var run []int
	for p.cur.Type == '(' {
		run = append(run, p.cur.Loc)
		p.advance()
	}
	if len(run) == 1 {
		return true
	}
	if _, ok := p.parenAfter[run[0]]; !ok {
		p.restoreScan(saved)
		p.scanParenBlock()
	}
	if p.parenQuery == nil {
		p.parenQuery = make(map[int]bool)
	}
	query := true
	p.parenQuery[run[len(run)-1]] = query
	for i := len(run) - 2; i >= 0; i-- {
		query = query && continuesParenQuery(p.parenAfter[run[i+1]])
		p.parenQuery[run[i]] = query
	}
	return query
}

// continuesParenQuery reports whether a token of type tokenType, following a
// parenthesized query that opens a block, keeps the block a query expression.
// The end of input does, so an unclosed block reads as it did before the
// inner block was looked at.
func continuesParenQuery(tokenType int) bool {
	switch tokenType {
	case ')', kwUNION, kwEXCEPT, kwINTERSECT, kwORDER, kwLIMIT, kwFOR, kwLOCK, tokEOF:
		return true
	default:
		return false
	}
}

// scanParenBlock consumes the block the '(' at p.cur opens, through its
// matching ')', and records in parenAfter, for every '(' in it, the type of
// the token after its matching ')' (tokEOF for one never closed).
func (p *Parser) scanParenBlock() {
	if p.parenAfter == nil {
		p.parenAfter = make(map[int]int)
	}
	var open []int
	for {
		switch p.cur.Type {
		case '(':
			open = append(open, p.cur.Loc)
		case ')':
			loc := open[len(open)-1]
			open = open[:len(open)-1]
			p.advance()
			p.parenAfter[loc] = p.cur.Type
			if len(open) == 0 {
				return
			}
			continue
		case tokEOF:
			for _, loc := range open {
				p.parenAfter[loc] = tokEOF
			}
			return
		}
		p.advance()
	}
}

// parseTableFactor parses a table factor: table_ref, subquery, LATERAL subquery, or parenthesized table_references.
func (p *Parser) parseTableFactor() (nodes.TableExpr, error) {
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()
	// LATERAL (SELECT ...) — LATERAL derived table
	if p.cur.Type == kwLATERAL && p.peekNext().Type == '(' {
		latStart := p.pos()
		p.advance() // consume LATERAL

		if _, err := p.expect('('); err != nil {
			return nil, err
		}

		if p.cur.Type != kwSELECT {
			return nil, &ParseError{
				Message:  "expected SELECT after LATERAL (",
				Position: p.cur.Loc,
			}
		}

		sel, err := p.parseSelectStmt()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}

		sub := &nodes.SubqueryExpr{
			Loc:     nodes.Loc{Start: latStart},
			Select:  sel,
			Lateral: true,
		}

		// Optional alias: [AS] alias
		if _, ok := p.match(kwAS); ok {
			alias, _, err := p.parseIdentOrText()
			if err != nil {
				return nil, err
			}
			sub.Alias = alias
		} else if (p.isIdentToken() || p.cur.Type == tokSCONST) && !p.isSelectTerminator() {
			alias, _, err := p.parseIdentOrText()
			if err != nil {
				return nil, err
			}
			sub.Alias = alias
		}

		// Optional derived column alias list: AS t(c1, c2)
		if sub.Alias != "" && p.cur.Type == '(' {
			cols, err := p.parseDerivedColumnList()
			if err != nil {
				return nil, err
			}
			sub.Columns = cols
		}

		sub.Loc.End = p.prev.End
		return sub, nil
	}

	if p.cur.Type == '(' {
		startPos := p.pos()

		// Disambiguate the parenthesized table factor. A leading run of '(' may
		// open either a derived-table subquery — `(SELECT ...)`, `((SELECT ...))`,
		// `((SELECT ...) UNION (...))` — or a parenthesized join group using
		// MySQL nested-join semantics — `(a JOIN b ON ...)`,
		// `((a LEFT JOIN b ON ...) JOIN c ON ...)` (the canonical SHOW CREATE VIEW
		// form for multi-table joins), whose first operand may itself be a
		// derived table — `((SELECT ...) d JOIN b ON ...)`. parenOpensQuery
		// tells them apart.
		derived := p.parenOpensQuery()
		p.advance()
		if derived {
			sel, err := p.parseSelectStmt()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}

			sub := &nodes.SubqueryExpr{
				Loc:    nodes.Loc{Start: startPos},
				Select: sel,
			}

			// Optional alias: [AS] alias
			if _, ok := p.match(kwAS); ok {
				alias, _, err := p.parseIdentOrText()
				if err != nil {
					return nil, err
				}
				sub.Alias = alias
			} else if (p.isIdentToken() || p.cur.Type == tokSCONST) && !p.isSelectTerminator() {
				alias, _, err := p.parseIdentOrText()
				if err != nil {
					return nil, err
				}
				sub.Alias = alias
			}

			// Optional derived column alias list: AS t(c1, c2)
			if sub.Alias != "" && p.cur.Type == '(' {
				cols, err := p.parseDerivedColumnList()
				if err != nil {
					return nil, err
				}
				sub.Columns = cols
			}

			sub.Loc.End = p.prev.End
			return sub, nil
		}

		// Parenthesized table references use MySQL nested-join semantics.
		refs, err := p.parseTableReferenceList()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		// OceanBase delta. MySQL parenthesizes a single table, a join, or a
		// list of two or more table references, but not a lone derived table
		// or table function (8.0.46 sql_yacc.yy table_factor, line 12034), so
		// mysql/parser rejects `((SELECT 1) d)` at the ')'. OceanBase accepts
		// a lone derived table, LATERAL derived table, or JSON_TABLE in
		// parentheses: 5.7.25-OceanBase_CE-v4.3.5.6 runs
		// `SELECT * FROM ((SELECT 1) d)`, `SELECT * FROM t0, (LATERAL
		// (SELECT 1) d)` and `SELECT * FROM (JSON_TABLE(...) AS jt)`
		// (loneParenthesizedFactors, TestLoneParenthesizedFactorsOceanBase).
		// The parentheses add nothing, so the factor is returned as written.
		return foldTableReferenceList(refs), nil
	}

	// JSON_TABLE(expr, path COLUMNS (...)) AS alias
	if p.cur.Type == kwJSON_TABLE && p.peekNext().Type == '(' {
		return p.parseJsonTable()
	}

	// DUAL — MySQL reserved keyword used as a dummy table in FROM clause.
	// After registration as a reserved keyword, parseIdent rejects it,
	// so we handle it explicitly here.
	if p.cur.Type == kwDUAL {
		tok := p.advance()
		return &nodes.TableRef{
			Loc:  nodes.Loc{Start: tok.Loc, End: p.prev.End},
			Name: "dual",
		}, nil
	}

	// Regular table reference with alias
	return p.parseTableRefWithAlias()
}

// parseJsonTable parses JSON_TABLE(expr, path COLUMNS (column_list)) [AS] alias.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/json-table-functions.html
//
//	JSON_TABLE(
//	    expr,
//	    path COLUMNS (column_list)
//	) [AS] alias
//
//	column_list:
//	    column[, column][, ...]
//
//	column:
//	    name FOR ORDINALITY
//	  | name type PATH string_path [on_empty] [on_error]
//	  | name type EXISTS PATH string_path
//	  | NESTED [PATH] path COLUMNS (column_list)
//
//	on_empty:
//	    {NULL | DEFAULT json_string | ERROR} ON EMPTY
//
//	on_error:
//	    {NULL | DEFAULT json_string | ERROR} ON ERROR
func (p *Parser) parseJsonTable() (nodes.TableExpr, error) {
	start := p.pos()
	p.advance() // consume JSON_TABLE

	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	// Parse the JSON expression
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if expr == nil {
		return nil, p.syntaxErrorAtCur()
	}

	if _, err := p.expect(','); err != nil {
		return nil, err
	}

	// Parse the path expression
	path, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if path == nil {
		return nil, p.syntaxErrorAtCur()
	}

	// COLUMNS
	if p.cur.Type != kwCOLUMNS {
		return nil, &ParseError{Message: "expected COLUMNS in JSON_TABLE", Position: p.cur.Loc}
	}
	p.advance()

	cols, err := p.parseJsonTableColumns()
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	// [AS] alias (required for JSON_TABLE)
	p.match(kwAS)
	alias, _, aErr := p.parseIdentOrText()
	if aErr != nil {
		return nil, aErr
	}

	return &nodes.JsonTableExpr{
		Loc:     nodes.Loc{Start: start, End: p.prev.End},
		Expr:    expr,
		Path:    path,
		Columns: cols,
		Alias:   alias,
	}, nil
}

// parseJsonTableColumns parses (column_list) for JSON_TABLE.
func (p *Parser) parseJsonTableColumns() ([]*nodes.JsonTableColumn, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	var cols []*nodes.JsonTableColumn
	for {
		col, err := p.parseJsonTableColumn()
		if err != nil {
			return nil, err
		}
		cols = append(cols, col)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return cols, nil
}

// parseJsonTableColumn parses a single JSON_TABLE column definition.
func (p *Parser) parseJsonTableColumn() (*nodes.JsonTableColumn, error) {
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()
	start := p.pos()

	// NESTED [PATH] path COLUMNS (column_list)
	if p.cur.Type == kwNESTED {
		p.advance()
		p.match(kwPATH) // optional PATH keyword
		if p.cur.Type != tokSCONST {
			return nil, &ParseError{Message: "expected path string in NESTED", Position: p.cur.Loc}
		}
		nestedPath := p.cur.Str
		p.advance()

		if p.cur.Type != kwCOLUMNS {
			return nil, &ParseError{Message: "expected COLUMNS after NESTED PATH", Position: p.cur.Loc}
		}
		p.advance()
		nestedCols, err := p.parseJsonTableColumns()
		if err != nil {
			return nil, err
		}
		return &nodes.JsonTableColumn{
			Loc:        nodes.Loc{Start: start, End: p.prev.End},
			Nested:     true,
			NestedPath: nestedPath,
			NestedCols: nestedCols,
		}, nil
	}

	// name FOR ORDINALITY | name type PATH path | name type EXISTS PATH path
	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}

	col := &nodes.JsonTableColumn{
		Loc:  nodes.Loc{Start: start},
		Name: name,
	}

	// FOR ORDINALITY
	if p.cur.Type == kwFOR {
		p.advance()
		p.match(kwORDINALITY)
		col.Ordinality = true
		col.Loc.End = p.prev.End
		return col, nil
	}

	// type
	typeName, err := p.parseDataType()
	if err != nil {
		return nil, err
	}
	col.TypeName = typeName

	// EXISTS PATH
	if p.cur.Type == kwEXISTS_KW {
		p.advance()
		col.Exists = true
	}

	// PATH
	if p.cur.Type == kwPATH {
		p.advance()
		if p.cur.Type != tokSCONST {
			return nil, &ParseError{Message: "expected path string", Position: p.cur.Loc}
		}
		col.Path = p.cur.Str
		p.advance()
	}

	// [on_empty] [on_error]
	// {NULL | DEFAULT json_string | ERROR} ON {EMPTY | ERROR}
	for p.cur.Type == kwNULL || p.cur.Type == kwDEFAULT || p.cur.Type == kwERROR_KW {
		option := ""
		switch p.cur.Type {
		case kwNULL:
			option = "NULL"
			p.advance()
		case kwDEFAULT:
			p.advance()
			if p.cur.Type == tokSCONST {
				option = "DEFAULT " + p.cur.Str
				p.advance()
			}
		case kwERROR_KW:
			option = "ERROR"
			p.advance()
		}

		if p.cur.Type != kwON {
			break
		}
		p.advance() // consume ON

		if p.cur.Type == kwEMPTY {
			col.EmptyOption = option
			p.advance()
		} else if p.cur.Type == kwERROR_KW {
			col.ErrorOption = option
			p.advance()
		}
	}

	col.Loc.End = p.prev.End
	return col, nil
}

// parseOrderByList parses ORDER BY items.
func (p *Parser) parseOrderByList() ([]*nodes.OrderByItem, error) {
	var items []*nodes.OrderByItem

	for {
		start := p.pos()
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if expr == nil {
			return nil, p.syntaxErrorAtCur()
		}

		item := &nodes.OrderByItem{
			Loc:  nodes.Loc{Start: start},
			Expr: expr,
		}

		// Completion: after ORDER BY expression, offer ASC/DESC/LIMIT and comma.
		p.checkCursor()
		if p.collectMode() {
			p.addTokenCandidate(kwASC)
			p.addTokenCandidate(kwDESC)
			p.addTokenCandidate(kwLIMIT)
			p.addTokenCandidate(',')
			return nil, &ParseError{Message: "collecting"}
		}

		if _, ok := p.match(kwDESC); ok {
			item.Direction = nodes.OrderDirectionDesc
			item.Desc = true
		} else if _, ok := p.match(kwASC); ok {
			item.Direction = nodes.OrderDirectionAsc
		} else {
			item.Direction = nodes.OrderDirectionDefault
		}

		item.Loc.End = p.prev.End
		items = append(items, item)

		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	return items, nil
}

// parseLimitValue reads one LIMIT or OFFSET value. It is no expression but an
// integer literal (limit_expr in OceanBase's grammar, measured on 4.3.5.6), so `LIMIT 1 + 1`, `LIMIT (1)`, `LIMIT -1`, `LIMIT 1 OR
// 2`, `LIMIT 1 || 2`, `LIMIT @a`, `LIMIT a`, `LIMIT 0x10` and `LIMIT 1.5` are all
// 1064.
func (p *Parser) parseLimitValue() (nodes.ExprNode, error) {
	if p.cur.Type == tokICONST {
		return p.parsePrimaryExpr()
	}
	return nil, p.syntaxErrorAtCur()
}

// parseLimitClause parses LIMIT [offset,] count or LIMIT count OFFSET offset.
func (p *Parser) parseLimitClause() (*nodes.Limit, error) {
	start := p.pos()

	count, err := p.parseLimitValue()
	if err != nil {
		return nil, err
	}
	if count == nil {
		return nil, p.syntaxErrorAtCur()
	}

	limit := &nodes.Limit{
		Loc:   nodes.Loc{Start: start},
		Count: count,
	}

	// LIMIT offset, count
	if p.cur.Type == ',' {
		p.advance()
		count2, err := p.parseLimitValue()
		if err != nil {
			return nil, err
		}
		if count2 == nil {
			return nil, p.syntaxErrorAtCur()
		}
		// In MySQL, LIMIT offset, count — first is offset, second is count
		limit.Offset = limit.Count
		limit.Count = count2
	}

	// LIMIT count OFFSET offset
	if _, ok := p.match(kwOFFSET); ok {
		offset, err := p.parseLimitValue()
		if err != nil {
			return nil, err
		}
		if offset == nil {
			return nil, p.syntaxErrorAtCur()
		}
		limit.Offset = offset
	}

	limit.Loc.End = p.prev.End
	return limit, nil
}

// parseForUpdateClause parses FOR UPDATE. OceanBase has neither FOR SHARE nor FOR
// UPDATE OF: opt_for_update is FOR UPDATE and a wait option, and both of the others
// are a 1064 (measured).
func (p *Parser) parseForUpdateClause() (*nodes.ForUpdate, error) {
	start := p.pos()
	p.advance() // consume FOR

	// Completion: after FOR keyword, offer UPDATE.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwUPDATE)
		return nil, &ParseError{Message: "collecting"}
	}

	fu := &nodes.ForUpdate{Loc: nodes.Loc{Start: start}}

	if _, err := p.expect(kwUPDATE); err != nil {
		return nil, err
	}

	// Completion: after FOR UPDATE, offer NOWAIT, SKIP, and WAIT.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwNOWAIT)
		p.addTokenCandidate(kwSKIP)
		p.addTokenCandidate(kwWAIT)
		return nil, &ParseError{Message: "collecting"}
	}

	// At most one of NOWAIT, WAIT n, SKIP LOCKED. OceanBase 4.3 spells each as
	// one alternative, so SKIP needs its LOCKED and "NOWAIT WAIT 5",
	// "WAIT 5 NOWAIT" and "WAIT 1 WAIT 2" are syntax errors. WAIT takes a
	// numeric literal, not an expression: WAIT 5, WAIT 5.5 and WAIT 1e3 are
	// accepted, WAIT @v, WAIT ?, WAIT (5), WAIT -5 and WAIT 0x10 are not
	// (MariaDB takes the hex form, OceanBase does not). The literal is kept as
	// written rather than converted, so nothing is invented about how
	// OceanBase rounds or clamps it. All container-verified on
	// oceanbase/oceanbase-ce:4.3.5-lts.
	switch {
	case p.cur.Type == kwNOWAIT:
		p.advance()
		fu.NoWait = true
	case p.cur.Type == kwWAIT:
		p.advance()
		if p.cur.Type != tokICONST && p.cur.Type != tokFCONST {
			return nil, &ParseError{
				Message:  "expected a number of seconds after WAIT",
				Position: p.cur.Loc,
			}
		}
		fu.WaitSeconds = p.cur.Str
		fu.HasWait = true
		p.advance()
	case p.cur.Type == kwSKIP:
		p.advance()
		if _, err := p.expect(kwLOCKED); err != nil {
			return nil, err
		}
		fu.SkipLocked = true
	}

	fu.Loc.End = p.prev.End
	return fu, nil
}

// parseLockInShareMode parses LOCK IN SHARE MODE (legacy syntax).
//
//	LOCK IN SHARE MODE
func (p *Parser) parseLockInShareMode() (*nodes.ForUpdate, error) {
	start := p.pos()
	p.advance() // consume LOCK

	if _, err := p.expect(kwIN); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwSHARE); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwMODE); err != nil {
		return nil, err
	}

	return &nodes.ForUpdate{
		Loc:             nodes.Loc{Start: start, End: p.prev.End},
		Share:           true,
		LockInShareMode: true,
	}, nil
}

// parseOutfileS3Options parses the trailing Aurora MySQL export options of
// INTO OUTFILE S3, in documented order:
//
//	[MANIFEST {ON | OFF}]
//	[OVERWRITE {ON | OFF}]
//	[ENCRYPTION {ON | OFF | SSE_S3 | SSE_KMS ['cmk_id']}]
func (p *Parser) parseOutfileS3Options(into *nodes.IntoClause) error {
	if p.cur.Type == kwMANIFEST {
		p.advance()
		v, err := p.parseOnOff()
		if err != nil {
			return err
		}
		into.Manifest = v
	}
	if p.cur.Type == kwOVERWRITE {
		p.advance()
		v, err := p.parseOnOff()
		if err != nil {
			return err
		}
		into.Overwrite = v
	}
	if p.cur.Type == kwENCRYPTION {
		p.advance()
		switch p.cur.Type {
		case kwON:
			into.Encryption = "ON"
			p.advance()
		case kwOFF:
			into.Encryption = "OFF"
			p.advance()
		case kwSSE_S3:
			into.Encryption = "SSE_S3"
			p.advance()
		case kwSSE_KMS:
			into.Encryption = "SSE_KMS"
			p.advance()
			if p.cur.Type == tokSCONST {
				into.EncryptionKey = p.cur.Str
				p.advance()
			}
		default:
			return p.syntaxErrorAtCur()
		}
	}
	return nil
}

// parseOnOff consumes ON or OFF and returns its canonical spelling.
func (p *Parser) parseOnOff() (string, error) {
	switch p.cur.Type {
	case kwON:
		p.advance()
		return "ON", nil
	case kwOFF:
		p.advance()
		return "OFF", nil
	default:
		return "", p.syntaxErrorAtCur()
	}
}

// parseIntoClause parses INTO OUTFILE / DUMPFILE / var_list.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/select.html
//
//	into_option: {
//	    INTO OUTFILE 'file_name'
//	        [CHARACTER SET charset_name]
//	        [{FIELDS | COLUMNS}
//	            [TERMINATED BY 'string']
//	            [[OPTIONALLY] ENCLOSED BY 'char']
//	            [ESCAPED BY 'char']
//	        ]
//	        [LINES
//	            [STARTING BY 'string']
//	            [TERMINATED BY 'string']
//	        ]
//	  | INTO DUMPFILE 'file_name'
//	  | INTO var_name [, var_name] ...
//	}
//
// Aurora MySQL extension — export straight to S3. Aurora registers as the
// MYSQL engine, so the extension is accepted unconditionally; stock MySQL 8.0
// rejects it with ER_PARSE_ERROR (1064). The Aurora-only options (FORMAT,
// HEADER, MANIFEST, OVERWRITE, ENCRYPTION) are accepted only after OUTFILE S3.
//
// Ref: https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/AuroraMySQL.Integrating.SaveIntoS3.html
//
//	INTO OUTFILE S3 's3_uri'
//	    [CHARACTER SET charset_name]
//	    [FORMAT {CSV | TEXT} [HEADER]]
//	    [{FIELDS | COLUMNS} ...]
//	    [LINES ...]
//	    [MANIFEST {ON | OFF}]
//	    [OVERWRITE {ON | OFF}]
//	    [ENCRYPTION {ON | OFF | SSE_S3 | SSE_KMS ['cmk_id']}]
func (p *Parser) parseIntoClause() (*nodes.IntoClause, error) {
	start := p.pos()
	into := &nodes.IntoClause{Loc: nodes.Loc{Start: start}}

	// Completion: right after INTO, offer the file forms.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwOUTFILE)
		p.addTokenCandidate(kwDUMPFILE)
		return nil, &ParseError{Message: "collecting"}
	}

	switch {
	case p.cur.Type == kwOUTFILE:
		p.advance()
		// Completion: INTO OUTFILE | → a file literal or the Aurora S3 form.
		p.checkCursor()
		if p.collectMode() {
			p.addTokenCandidate(kwS3)
			return nil, &ParseError{Message: "collecting"}
		}
		// Aurora MySQL: INTO OUTFILE S3 's3_uri'
		if p.cur.Type == kwS3 {
			p.advance()
			into.OutfileS3 = true
		}
		outfile, err := p.parseIntoString()
		if err != nil {
			return nil, err
		}
		into.Outfile, into.HasOutfile = outfile, true

		// [CHARACTER SET charset_name]
		if p.cur.Type == kwCHARACTER {
			p.advance()
			if _, err := p.expect(kwSET); err != nil {
				return nil, err
			}
			charset, csErr := p.parseCharsetName()
			if csErr != nil {
				return nil, csErr
			}
			into.Charset = charset
		} else if p.cur.Type == kwCHARSET {
			p.advance()
			charset, csErr := p.parseCharsetName()
			if csErr != nil {
				return nil, csErr
			}
			into.Charset = charset
		}

		// Aurora MySQL: [FORMAT {CSV | TEXT} [HEADER]]
		if into.OutfileS3 && p.cur.Type == kwFORMAT {
			p.advance()
			switch p.cur.Type {
			case kwCSV:
				into.Format = "CSV"
			case kwTEXT:
				into.Format = "TEXT"
			default:
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
			if p.cur.Type == kwHEADER {
				into.Header = true
				p.advance()
			}
		}

		if err := p.parseIntoFields(into); err != nil {
			return nil, err
		}
		if err := p.parseIntoLines(into); err != nil {
			return nil, err
		}

		if into.OutfileS3 {
			if err := p.parseOutfileS3Options(into); err != nil {
				return nil, err
			}
		}

	// DUMPFILE is no reserved word, so with no file after it it is a variable's
	// name, which the default arm reads.
	case p.cur.Type == kwDUMPFILE && p.peekNext().Type == tokSCONST:
		p.advance()
		dumpfile, err := p.parseIntoString()
		if err != nil {
			return nil, err
		}
		into.Dumpfile, into.HasDumpfile = dumpfile, true

	default:
		// INTO var_name [, var_name ...]: into_var_list, each a @user_var or a bare
		// identifier (a stored-program variable). A system variable is no target,
		// so INTO @@autocommit is a 1064, and neither a missing target nor a comma
		// with nothing after it parses.
		for {
			if p.isVariableRef() {
				v, err := p.parseVariableRef()
				if err != nil {
					return nil, err
				}
				if v.System {
					return nil, &ParseError{Message: `syntax error at or near "@@"`, Position: v.Loc.Start}
				}
				into.Vars = append(into.Vars, v)
			} else if p.cur.Type == tokIDENT || (p.cur.Type >= 700 && isLvalueKeyword(p.cur.Type)) {
				// Bare identifier: SELECT ... INTO local_var (stored program)
				nameStart := p.pos()
				name, _, err := p.parseLvalueIdent()
				if err != nil {
					return nil, err
				}
				into.Vars = append(into.Vars, &nodes.VariableRef{
					Loc:  nodes.Loc{Start: nameStart, End: p.prev.End},
					Name: name,
					Bare: true,
				})
			} else {
				return nil, p.syntaxErrorAtCur()
			}
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	}

	into.Loc.End = p.prev.End
	return into, nil
}

// parseIntoString parses a string of INTO OUTFILE: one literal, or several that the
// server reads as one, as in INTO OUTFILE 'a' 'b' (measured).
func (p *Parser) parseIntoString() (string, error) {
	if p.cur.Type != tokSCONST {
		return "", p.syntaxErrorAtCur()
	}
	val := p.cur.Str
	p.advance()
	for p.cur.Type == tokSCONST {
		val += p.cur.Str
		p.advance()
	}
	return val, nil
}

// parseIntoByString parses the BY and the string of a FIELDS or LINES option. BY is
// required: FIELDS TERMINATED ',' is a 1064 (measured).
func (p *Parser) parseIntoByString() (string, error) {
	if _, err := p.expect(kwBY); err != nil {
		return "", err
	}
	return p.parseIntoString()
}

// parseIntoFields parses [{FIELDS | COLUMNS} field_term ...]. OceanBase's
// field_term_list takes at least one option, in any order, each as often as written,
// and the last one written wins, so each is read where it stands.
func (p *Parser) parseIntoFields(into *nodes.IntoClause) error {
	if p.cur.Type != kwFIELDS && p.cur.Type != kwCOLUMNS {
		return nil
	}
	p.advance()
	into.HasFieldsClause = true
	for count := 0; ; count++ {
		var err error
		switch p.cur.Type {
		case kwTERMINATED:
			p.advance()
			into.FieldsTerminatedBy, err = p.parseIntoByString()
			into.HasFieldsTerminatedBy = true
		case kwOPTIONALLY:
			p.advance()
			if _, err := p.expect(kwENCLOSED); err != nil {
				return err
			}
			into.FieldsEnclosedBy, err = p.parseIntoByString()
			into.HasFieldsEnclosedBy, into.FieldsOptionalEncl = true, true
		case kwENCLOSED:
			p.advance()
			into.FieldsEnclosedBy, err = p.parseIntoByString()
			into.HasFieldsEnclosedBy, into.FieldsOptionalEncl = true, false
		case kwESCAPED:
			p.advance()
			into.FieldsEscapedBy, err = p.parseIntoByString()
			into.HasFieldsEscapedBy = true
		default:
			if count == 0 {
				return p.syntaxErrorAtCur()
			}
			return nil
		}
		if err != nil {
			return err
		}
	}
}

// parseIntoLines parses [LINES line_term ...], which takes at least one option, in any
// order, as parseIntoFields does.
func (p *Parser) parseIntoLines(into *nodes.IntoClause) error {
	if p.cur.Type != kwLINES {
		return nil
	}
	p.advance()
	into.HasLinesClause = true
	for count := 0; ; count++ {
		var err error
		switch p.cur.Type {
		case kwSTARTING:
			p.advance()
			into.LinesStartingBy, err = p.parseIntoByString()
			into.HasLinesStartingBy = true
		case kwTERMINATED:
			p.advance()
			into.LinesTerminatedBy, err = p.parseIntoByString()
			into.HasLinesTerminatedBy = true
		default:
			if count == 0 {
				return p.syntaxErrorAtCur()
			}
			return nil
		}
		if err != nil {
			return err
		}
	}
}

// setOpPrecedence returns the precedence for a set operation.
// INTERSECT has higher precedence than UNION/EXCEPT in MySQL 8.0.31+.
func setOpPrecedence(tokType int) int {
	switch tokType {
	case kwINTERSECT:
		return 2
	case kwUNION, kwEXCEPT:
		return 1
	}
	return 0
}

// parseParenIdentList parses a parenthesized comma-separated list of identifiers.
func (p *Parser) parseParenIdentList() ([]string, error) {
	return p.parseParenIdentListWithCompletion("")
}

// parseParenIdentListWithCompletion is like parseParenIdentList but emits the
// given completion rule (e.g. "columnref") at each identifier position. Pass ""
// to disable completion candidate emission.
func (p *Parser) parseParenIdentListWithCompletion(completionRule string) ([]string, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	var names []string
	for {
		// Completion: at each identifier position (including empty parens).
		if completionRule != "" {
			p.checkCursor()
			if p.collectMode() {
				p.addRuleCandidate(completionRule)
				return nil, &ParseError{Message: "collecting"}
			}
		}

		name, _, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		names = append(names, name)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return names, nil
}

// tableExprLoc returns the location of a TableExpr.
func tableExprLoc(te nodes.TableExpr) nodes.Loc {
	switch t := te.(type) {
	case *nodes.TableRef:
		return t.Loc
	case *nodes.JoinClause:
		return t.Loc
	case *nodes.SubqueryExpr:
		return t.Loc
	case *nodes.JsonTableExpr:
		return t.Loc
	}
	return nodes.NoLoc()
}

// parseNamedWindowList parses WINDOW window_name AS (window_spec) [, ...].
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/window-functions-named-windows.html
//
//	window_clause:
//	    WINDOW window_name AS (window_spec) [, window_name AS (window_spec)] ...
func (p *Parser) parseNamedWindowList() ([]*nodes.WindowDef, error) {
	var defs []*nodes.WindowDef
	for {
		start := p.pos()
		name, _, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(kwAS); err != nil {
			return nil, err
		}
		if _, err := p.expect('('); err != nil {
			return nil, err
		}

		wd := &nodes.WindowDef{Loc: nodes.Loc{Start: start}, Name: name}

		// Optional reference to existing window
		if p.isIdentToken() && p.cur.Type != kwPARTITION && p.cur.Type != kwORDER &&
			p.cur.Type != kwROWS && p.cur.Type != kwRANGE && p.cur.Type != kwGROUPS {
			var refErr error
			wd.RefName, _, refErr = p.parseIdent()
			if refErr != nil {
				return nil, refErr
			}
		}

		if p.cur.Type == kwPARTITION {
			p.advance()
			if _, err := p.expect(kwBY); err != nil {
				return nil, err
			}
			exprs, err := p.parseExprList()
			if err != nil {
				return nil, err
			}
			wd.PartitionBy = exprs
		}

		if p.cur.Type == kwORDER {
			p.advance()
			if _, err := p.expect(kwBY); err != nil {
				return nil, err
			}
			orderBy, err := p.parseOrderByList()
			if err != nil {
				return nil, err
			}
			wd.OrderBy = orderBy
		}

		if p.cur.Type == kwROWS || p.cur.Type == kwRANGE || p.cur.Type == kwGROUPS {
			frame, err := p.parseFrameClause()
			if err != nil {
				return nil, err
			}
			wd.Frame = frame
		}

		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		wd.Loc.End = p.prev.End
		defs = append(defs, wd)

		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return defs, nil
}

// parseSourceClauses parses the ORDER BY and LIMIT that may follow a VALUES
// query primary.
func (p *Parser) parseSourceClauses() ([]*nodes.OrderByItem, *nodes.Limit, error) {
	var orderBy []*nodes.OrderByItem
	if p.cur.Type == kwORDER {
		p.advance()
		if _, err := p.expect(kwBY); err != nil {
			return nil, nil, err
		}
		list, err := p.parseOrderByList()
		if err != nil {
			return nil, nil, err
		}
		orderBy = list
	}

	var limit *nodes.Limit
	if _, ok := p.match(kwLIMIT); ok {
		l, err := p.parseLimitClause()
		if err != nil {
			return nil, nil, err
		}
		limit = l
	}
	return orderBy, limit, nil
}

// parseValuesConstructor parses VALUES ROW(...) [, ROW(...)] ..., MySQL's
// table_value_constructor, into a ValuesStmt that holds only the rows.
func (p *Parser) parseValuesConstructor() (*nodes.ValuesStmt, error) {
	start := p.pos()
	p.advance() // consume VALUES

	stmt := &nodes.ValuesStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Parse ROW(...) list
	for {
		if _, err := p.expect(kwROW); err != nil {
			return nil, err
		}
		row, err := p.parseRowValue()
		if err != nil {
			return nil, err
		}
		stmt.Rows = append(stmt.Rows, row)

		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseSubqueryExpr parses a subquery expression: SELECT ... or a
// parenthesized query expression. The caller owns the surrounding subquery
// parentheses used by IN/EXISTS/scalar-subquery syntax.
//
// Redundant bare paren-wrapper layers are folded into the SubqueryExpr, whose
// rendering owns the single canonical pair. This matches the engine: MySQL
// 8.0.32 and 5.7.25 both store `((SELECT ...))` (any extra depth, in any
// expression position — select item, binary operand, IN, EXISTS, quantified
// comparison) as exactly `(select ...)` in SHOW CREATE VIEW bodies, so keeping
// the extra layers would phantom-diff the user form against the stored form.
// A wrapper holding outer clauses (WITH / ORDER BY / LIMIT / locking)
// keeps its parens — they scope those clauses. The engine does NOT keep that
// pair: it merges wrapper clauses into the leaf and drops the parens
// (`(((SELECT 1) LIMIT 1) = 1)` is stored as `((select 1 limit 1) = 1)`;
// `((SELECT 1 LIMIT 2) LIMIT 1)` even discards the inner LIMIT). Keeping the
// wrapper is still safe — readbacks never contain the wrapper form — but
// user-form SDL written with clause-bearing wrappers phantom-diffs against
// the stored form (a perpetual no-op view replace), the same residual class
// as the unexpanded row-compare and set-op ARM paren forms.
func (p *Parser) parseSubqueryExpr() (*nodes.SubqueryExpr, error) {
	start := p.pos()
	sel, err := p.parseSelectStmt()
	if err != nil {
		return nil, err
	}
	for sel != nil && sel.ParenSource != nil && len(sel.CTEs) == 0 &&
		len(sel.OrderBy) == 0 && sel.Limit == nil && sel.ForUpdate == nil {
		sel = sel.ParenSource
	}
	return &nodes.SubqueryExpr{
		Loc:    nodes.Loc{Start: start, End: p.prev.End},
		Select: sel,
	}, nil
}

// parseIndexHints parses an optional list of index hints following a table reference.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/index-hints.html
//
//	index_hint_list:
//	    index_hint [index_hint] ...
//
//	index_hint:
//	    USE {INDEX|KEY}
//	      [FOR {JOIN|ORDER BY|GROUP BY}] ([index_list])
//	  | {IGNORE|FORCE} {INDEX|KEY}
//	      [FOR {JOIN|ORDER BY|GROUP BY}] (index_list)
//
//	index_list:
//	    index_name [, index_name] ...
func (p *Parser) parseIndexHints() ([]*nodes.IndexHint, error) {
	var hints []*nodes.IndexHint
	for p.cur.Type == kwUSE || p.cur.Type == kwFORCE || p.cur.Type == kwIGNORE {
		hint, err := p.parseIndexHint()
		if err != nil {
			return hints, err
		}
		hints = append(hints, hint)
	}
	return hints, nil
}

// parseIndexHint parses a single index hint.
//
//	index_hint:
//	    USE {INDEX|KEY}
//	      [FOR {JOIN|ORDER BY|GROUP BY}] ([index_list])
//	  | {IGNORE|FORCE} {INDEX|KEY}
//	      [FOR {JOIN|ORDER BY|GROUP BY}] (index_list)
func (p *Parser) parseIndexHint() (*nodes.IndexHint, error) {
	start := p.pos()
	hint := &nodes.IndexHint{
		Loc: nodes.Loc{Start: start},
	}

	// Parse hint type: USE | FORCE | IGNORE
	switch p.cur.Type {
	case kwUSE:
		hint.Type = nodes.HintUse
	case kwFORCE:
		hint.Type = nodes.HintForce
	case kwIGNORE:
		hint.Type = nodes.HintIgnore
	default:
		return nil, &ParseError{Message: "expected USE, FORCE, or IGNORE", Position: p.cur.Loc}
	}
	p.advance()

	// Parse INDEX | KEY
	if p.cur.Type != kwINDEX && p.cur.Type != kwKEY {
		return nil, &ParseError{Message: "expected INDEX or KEY", Position: p.cur.Loc}
	}
	p.advance()

	// Optional FOR {JOIN | ORDER BY | GROUP BY}
	hint.Scope = nodes.HintScopeAll
	if p.cur.Type == kwFOR {
		p.advance()
		switch p.cur.Type {
		case kwJOIN:
			hint.Scope = nodes.HintScopeJoin
			p.advance()
		case kwORDER:
			p.advance()
			if _, err := p.expect(kwBY); err != nil {
				return nil, err
			}
			hint.Scope = nodes.HintScopeOrderBy
		case kwGROUP:
			p.advance()
			if _, err := p.expect(kwBY); err != nil {
				return nil, err
			}
			hint.Scope = nodes.HintScopeGroupBy
		default:
			return nil, &ParseError{Message: "expected JOIN, ORDER BY, or GROUP BY", Position: p.cur.Loc}
		}
	}

	// Parse (index_list) — required parens, list may be empty for USE
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	// Completion: inside index hint parentheses, offer index_ref candidates.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("index_ref")
		return nil, &ParseError{Message: "collecting"}
	}

	if p.cur.Type != ')' {
		// Parse comma-separated index names
		for {
			name, _, err := p.parseIdent()
			if err != nil {
				return nil, err
			}
			hint.Indexes = append(hint.Indexes, name)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	hint.Loc.End = p.prev.End
	return hint, nil
}
