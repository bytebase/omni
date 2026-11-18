package parser

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/bytebase/omni/starrocks/ast"
)

// ---------------------------------------------------------------------------
// SUBMIT TASK / DROP TASK (StarRocks.g4:625-642)
// ---------------------------------------------------------------------------

// parseSubmitTask parses submitTaskStatement. On entry SUBMIT has been
// consumed and cur is TASK.
//
//	SUBMIT TASK qualifiedName? taskClause*
//	    AS (createTableAsSelectStatement | insertStatement | dataCacheSelectStatement)
//
// taskClause is properties or taskScheduleDesc.
func (p *Parser) parseSubmitTask(startLoc ast.Loc) (ast.Node, error) {
	if _, err := p.expect(kwTASK); err != nil {
		return nil, err
	}

	stmt := &ast.SubmitTaskStmt{Loc: startLoc}

	// The name is optional. PROPERTIES and SCHEDULE are non-reserved
	// identifiers, so they are a name unless the token after them opens a
	// taskClause.
	if p.submitTaskNameAt() {
		name, err := p.parseMultipartIdentifier()
		if err != nil {
			return nil, err
		}
		if err := checkTaskName(name); err != nil {
			return nil, err
		}
		stmt.Name = name
	}

	for {
		switch p.cur.Kind {
		case kwPROPERTIES:
			props, err := p.parseProperties()
			if err != nil {
				return nil, err
			}
			for _, prop := range props {
				if strings.EqualFold(prop.Key, "warehouse") {
					return nil, &ParseError{
						Position: prop.Loc.Start,
						End:      prop.Loc.End,
						Message:  "task properties may not set warehouse",
					}
				}
			}
			stmt.Properties = append(stmt.Properties, props...)
		case kwSCHEDULE:
			schedule, err := p.parseTaskSchedule()
			if err != nil {
				return nil, err
			}
			stmt.Schedule = schedule
		default:
			// AS starts the submitted statement.
			goto clausesDone
		}
	}
clausesDone:

	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}

	var err error
	switch p.cur.Kind {
	case kwCREATE:
		// createTableAsSelectStatement (StarRocks.g4:488). parseCreateTable
		// reads it from the CREATE keyword on, so hold the result to that
		// rule: a query must follow AS, and the column list is names only.
		p.advance() // consume CREATE
		p.strictCTAS = true
		stmt.Statement, err = p.parseCreateTable()
		p.strictCTAS = false
		if err == nil {
			ctas, ok := stmt.Statement.(*ast.CreateTableStmt)
			if !ok || ctas.AsSelect == nil || ctas.External || ctas.Like != nil ||
				ctas.Engine != "" || len(ctas.Rollup) > 0 || len(ctas.BrokerProperties) > 0 ||
				len(ctas.Columns) > 0 || len(ctas.Constraints) > 0 {
				return nil, p.syntaxErrorAtCur()
			}
		}
	case kwINSERT:
		// insertStatement requires INTO or OVERWRITE. parseInsert also admits
		// a bare INSERT destination, which is a known top-level gap; the
		// grammar inside SUBMIT TASK has no such alternative.
		if next := p.peekNext().Kind; next != kwINTO && next != kwOVERWRITE {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Statement, err = p.parseInsert()
	case kwCACHE:
		stmt.Statement, err = p.parseDataCacheSelect()
	case kwEXPLAIN, kwDESC, kwDESCRIBE:
		// Only insertStatement, of the three submitted forms, takes an
		// explainDesc prefix.
		var explain ast.Node
		explain, err = p.parseExplain()
		if err == nil {
			if explained := explain.(*ast.ExplainStmt); explained.Query == nil {
				err = p.syntaxErrorAtCur()
			} else if _, ok := explained.Query.(*ast.InsertStmt); !ok {
				err = p.syntaxErrorAtCur()
			}
		}
		stmt.Statement = explain
	default:
		return nil, p.syntaxErrorAtCur()
	}
	if err != nil {
		return nil, err
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// checkTaskName enforces qualifiedNameToTaskName's shape: a task is named
// `task` or `database.task`, and a longer name is a parse error.
func checkTaskName(name *ast.ObjectName) error {
	if name == nil || len(name.Parts) <= 2 {
		return nil
	}
	return &ParseError{
		Position: name.Loc.Start,
		End:      name.Loc.End,
		Message:  "task name must be task or database.task",
	}
}

// submitTaskNameAt reports whether cur is the optional task name rather than
// the start of a taskClause. Both clause keywords are non-reserved, so
// `SUBMIT TASK properties AS ...` names the task properties while
// `SUBMIT TASK properties('a' = 'a') AS ...` starts a clause.
func (p *Parser) submitTaskNameAt() bool {
	if !isIdentifierToken(p.cur.Kind) {
		return false
	}
	switch p.cur.Kind {
	case kwPROPERTIES:
		return p.peekNext().Kind != int('(')
	case kwSCHEDULE:
		next := p.peekNext().Kind
		return next != kwSTART && next != kwEVERY
	}
	return true
}

// parseTaskSchedule parses taskScheduleDesc (StarRocks.g4:640):
//
//	SCHEDULE (START '(' string ')')? EVERY '(' INTERVAL integer unit ')'
//
// The grammar's taskInterval takes an expression, but StarRocks's AST builder
// takes an integer value (as alterMultiRangePartition documents for the same
// shape), so the value is an integer literal and the unit is one of DAY, HOUR,
// MINUTE or SECOND. On entry cur is SCHEDULE.
func (p *Parser) parseTaskSchedule() (*ast.TaskSchedule, error) {
	startLoc := p.cur.Loc
	p.advance() // consume SCHEDULE

	schedule := &ast.TaskSchedule{Loc: startLoc}

	if p.cur.Kind == kwSTART {
		p.advance() // consume START
		if _, err := p.expect(int('(')); err != nil {
			return nil, err
		}
		start, err := p.expect(tokString)
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(int(')')); err != nil {
			return nil, err
		}
		if !validTaskStart(start.Str) {
			return nil, &ParseError{
				Position: start.Loc.Start,
				End:      start.Loc.End,
				Message:  "invalid date format: " + start.Str,
			}
		}
		schedule.Start = start.Str
	}

	if _, err := p.expect(kwEVERY); err != nil {
		return nil, err
	}
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwINTERVAL); err != nil {
		return nil, err
	}
	// The grammar's value is an expression. StarRocks unwraps parentheses and
	// then requires an integer literal, so `(1)` is accepted and
	// `CURRENT_USER` is not.
	parens := 0
	for {
		switch p.cur.Kind {
		case int('('):
			p.advance() // consume '('
			parens++
		case int('+'):
			p.advance() // unary plus folds to the literal
		default:
			goto intervalValue
		}
	}
intervalValue:
	value, err := p.expect(tokInt)
	if err != nil {
		return nil, err
	}
	// The local lexer also folds 0x10 and 0b10 into tokInt, but the grammar
	// tokenizes those as identifiers and the engine needs an IntLiteral.
	if !isDecimalInteger(value.Str) {
		return nil, &ParseError{
			Position: value.Loc.Start,
			End:      value.Loc.End,
			Message:  "interval value must be a decimal integer",
		}
	}
	// StarRocks converts only values through MaxInt64 to an IntLiteral; a
	// larger one is a LargeIntLiteral and the task builder refuses it.
	if _, err := strconv.ParseInt(value.Str, 10, 64); err != nil {
		return nil, &ParseError{
			Position: value.Loc.Start,
			End:      value.Loc.End,
			Message:  "interval value is out of range",
		}
	}
	for range parens {
		if _, err := p.expect(int(')')); err != nil {
			return nil, err
		}
	}
	schedule.Interval = value.Str

	switch p.cur.Kind {
	case kwDAY, kwHOUR, kwMINUTE, kwSECOND:
		schedule.Unit = strings.ToUpper(p.advance().Str)
	default:
		return nil, p.syntaxErrorAtCur()
	}

	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}

	seconds, ok := taskIntervalSeconds(schedule.Interval, schedule.Unit)
	if !ok || seconds < 10 {
		return nil, &ParseError{
			Position: value.Loc.Start,
			End:      value.Loc.End,
			Message:  "task schedule interval must be at least 10 seconds",
		}
	}

	schedule.Loc.End = p.prev.Loc.End
	return schedule, nil
}

// taskIntervalSeconds converts a task interval to seconds, reporting false on
// overflow. StarRocks's default task_min_schedule_interval_s is 10 seconds.
func taskIntervalSeconds(interval, unit string) (int64, bool) {
	value, err := strconv.ParseInt(interval, 10, 64)
	if err != nil {
		return 0, false
	}
	factor := int64(1)
	switch unit {
	case "SECOND":
	case "MINUTE":
		factor = 60
	case "HOUR":
		factor = 3600
	case "DAY":
		factor = 86400
	default:
		return 0, false
	}
	if value > math.MaxInt64/factor {
		// Java's TimeUnit conversion saturates a positive overflow to
		// Long.MAX_VALUE, so the engine accepts these intervals.
		return math.MaxInt64, true
	}
	return value * factor, true
}

// isDecimalInteger reports whether s is a bare decimal integer, the only
// spelling StarRocks's interval conversion accepts.
func isDecimalInteger(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

// validTaskStart reports whether s is a timestamp DateUtils.probeFormat
// accepts: yyyyMMdd, yyyy-MM-dd, yyyyMMddHHmmss, yyyy-MM-dd HH:mm:ss, or the
// same with microseconds.
func validTaskStart(s string) bool {
	var layouts []string
	switch len(s) {
	case 8:
		layouts = []string{"20060102", "06-01-02"}
	case 10:
		layouts = []string{"2006-01-02"}
	case 14:
		layouts = []string{"20060102150405", "06-01-02 15:04", "06-01-02T15:04"}
	case 16:
		layouts = []string{"2006-01-02 15:04", "2006-01-02T15:04"}
	case 17:
		layouts = []string{"06-01-02 15:04:05", "06-01-02T15:04:05"}
	case 19:
		layouts = []string{"2006-01-02 15:04:05", "2006-01-02T15:04:05"}
	case 26:
		layouts = []string{"2006-01-02 15:04:05.000000", "2006-01-02T15:04:05.000000"}
	default:
		return false
	}
	for _, layout := range layouts {
		if _, err := time.Parse(layout, s); err == nil {
			return true
		}
	}
	return false
}

// parseDropTask parses dropTaskStatement. On entry DROP has been consumed and
// cur is TASK.
//
//	DROP TASK qualifiedName FORCE?
func (p *Parser) parseDropTask(startLoc ast.Loc) (ast.Node, error) {
	if _, err := p.expect(kwTASK); err != nil {
		return nil, err
	}

	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	if err := checkTaskName(name); err != nil {
		return nil, err
	}

	stmt := &ast.DropTaskStmt{Name: name, Loc: startLoc}
	if p.cur.Kind == kwFORCE {
		p.advance() // consume FORCE
		stmt.Force = true
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// ALTER LOAD (StarRocks.g4:1499)
// ---------------------------------------------------------------------------

// parseAlterLoad parses alterLoadStatement. On entry ALTER has been consumed
// and cur is LOAD.
//
//	ALTER LOAD FOR (db=qualifiedName '.')? name=identifier jobProperties?
func (p *Parser) parseAlterLoad(startLoc ast.Loc) (ast.Node, error) {
	if _, err := p.expect(kwLOAD); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwFOR); err != nil {
		return nil, err
	}

	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}

	stmt := &ast.AlterLoadStmt{Name: name, Loc: startLoc}
	if p.cur.Kind == kwPROPERTIES {
		props, err := p.parseProperties()
		if err != nil {
			return nil, err
		}
		stmt.Properties = props
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// CACHE SELECT (StarRocks.g4:1867)
// ---------------------------------------------------------------------------

// parseDataCacheSelect parses dataCacheSelectStatement. On entry CACHE has
// been consumed and cur is SELECT.
//
//	CACHE SELECT selectItem (',' selectItem)* FROM qualifiedName
//	    (WHERE expression)? properties?
func (p *Parser) parseDataCacheSelect() (ast.Node, error) {
	startLoc := p.cur.Loc
	p.advance() // consume CACHE

	if _, err := p.expect(kwSELECT); err != nil {
		return nil, err
	}

	items, err := p.parseSelectList()
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}
	from, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	if len(from.Parts) > 3 {
		return nil, &ParseError{
			Position: from.Loc.Start,
			End:      from.Loc.End,
			Message:  "table name must be table, database.table or catalog.database.table",
		}
	}

	stmt := &ast.DataCacheSelectStmt{Items: items, From: from, Loc: startLoc}
	if p.cur.Kind == kwWHERE {
		p.advance() // consume WHERE
		if stmt.Where, err = p.parseExpr(); err != nil {
			return nil, err
		}
	}
	if p.cur.Kind == kwPROPERTIES {
		if stmt.Properties, err = p.parseProperties(); err != nil {
			return nil, err
		}
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// TRACE (StarRocks.g4:2584)
// ---------------------------------------------------------------------------

// parseOptimizerTrace parses optimizerTrace. On entry TRACE has been consumed.
//
//	TRACE (ALL | LOGS | TIMES | VALUES | REASON) identifier? query
//
// The trace identifier is optional and can be any identifier, so an identifier
// before SELECT/WITH is the trace's name and not the query.
func (p *Parser) parseOptimizerTrace(startLoc ast.Loc) (ast.Node, error) {
	stmt := &ast.OptimizerTraceStmt{Loc: startLoc}

	switch p.cur.Kind {
	case kwALL:
		stmt.Mode = "ALL"
	case kwLOGS:
		stmt.Mode = "LOGS"
	case kwTIMES:
		stmt.Mode = "TIMES"
	case kwVALUES:
		stmt.Mode = "VALUES"
	case kwREASON:
		stmt.Mode = "REASON"
	default:
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()

	if isIdentifierToken(p.cur.Kind) {
		identifier, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Identifier = identifier
	}

	query, err := p.parseTraceQuery()
	if err != nil {
		return nil, err
	}
	if p.cur.Kind == kwINTO {
		if queryHasInto(query) {
			return nil, p.syntaxErrorAtCur()
		}
		into, err := p.parseIntoOutfileClause()
		if err != nil {
			return nil, err
		}
		attachInto(query, into)
	}
	stmt.Query = query

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// attachWith sets the WITH clause on the query node that owns the whole
// queryNoWith, so a set-operation tail stays inside its scope.
func attachWith(node ast.Node, with *ast.WithClause) {
	switch n := node.(type) {
	case *ast.ParenSelect:
		// The outer clause belongs to the group, above the inner query's own
		// WITH, so the two lexical scopes stay separate.
		n.With = with
		n.Loc.Start = with.Loc.Start
	case *ast.SetOpStmt:
		// The clause scopes over both arms, not just the leftmost SELECT.
		n.With = with
		n.Loc.Start = with.Loc.Start
	case *ast.SelectStmt:
		n.With = with
		n.Loc.Start = with.Loc.Start
	}
}

// queryHasInto reports whether the query already carries an INTO OUTFILE
// clause, including one on an arm or inside a parenthesized group.
func queryHasInto(node ast.Node) bool {
	stack := []ast.Node{node}
	for len(stack) > 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		switch q := n.(type) {
		case *ast.SelectStmt:
			if q.Into != nil {
				return true
			}
		case *ast.ParenSelect:
			if q.Into != nil {
				return true
			}
			stack = append(stack, q.Sel)
		case *ast.SetOpStmt:
			if q.Into != nil {
				return true
			}
			stack = append(stack, q.Left, q.Right)
		}
	}
	return false
}

// attachInto stores an INTO OUTFILE clause on the query node that owns the
// whole queryNoWith, so a child's Loc never extends past its parent's.
func attachInto(node ast.Node, into *ast.IntoOutfileClause) {
	switch n := node.(type) {
	case *ast.ParenSelect:
		n.Into = into
		n.Loc.End = into.Loc.End
	case *ast.SetOpStmt:
		n.Into = into
		n.Loc.End = into.Loc.End
	case *ast.SelectStmt:
		n.Into = into
		n.Loc.End = into.Loc.End
	}
}

// takeRightmostSelectInto removes and returns the INTO OUTFILE clause of the
// rightmost unparenthesized SELECT. A parenthesized right operand keeps its
// clause inside a subquery, which has no outfile grammar, so it is not taken.
func takeRightmostSelectInto(node ast.Node) (*ast.IntoOutfileClause, int) {
	for {
		switch n := node.(type) {
		case *ast.SetOpStmt:
			node = n.Right
		case *ast.SelectStmt:
			into := n.Into
			n.Into = nil
			if into == nil {
				return nil, 0
			}
			n.Loc.End = selectEndBeforeInto(n)
			return into, n.Loc.End
		default:
			return nil, 0
		}
	}
}

// selectEndBeforeInto returns the end of the last token in a SELECT other than
// its INTO OUTFILE clause.
func selectEndBeforeInto(n *ast.SelectStmt) int {
	end := n.Loc.Start
	note := func(node ast.Node) {
		if loc := ast.NodeLoc(node); loc.End > end {
			end = loc.End
		}
	}
	if n.With != nil {
		note(n.With)
	}
	for _, item := range n.Items {
		note(item)
	}
	for _, from := range n.From {
		note(from)
	}
	for _, node := range []ast.Node{n.Where, n.Having, n.Qualify, n.Limit, n.Offset} {
		if node != nil {
			note(node)
		}
	}
	for _, g := range n.GroupBy {
		note(g)
	}
	for _, o := range n.OrderBy {
		note(o)
	}
	return end
}

// settleTraceSetOpInto moves a genuine query-level OUTFILE from the rightmost
// arm onto the set operation and rejects one anywhere else.
func (p *Parser) settleTraceSetOpInto(query ast.Node) error {
	setOp, ok := query.(*ast.SetOpStmt)
	if !ok {
		return nil
	}
	into, armEnd := takeRightmostSelectInto(setOp)
	if queryHasInto(setOp) {
		return p.syntaxErrorAtCur()
	}
	if into != nil {
		// The arm and the right spine below the root no longer own the clause.
		for node := setOp.Right; ; {
			next, ok := node.(*ast.SetOpStmt)
			if !ok {
				break
			}
			next.Loc.End = armEnd
			node = next.Right
		}
		setOp.Into = into
		if into.Loc.End > setOp.Loc.End {
			setOp.Loc.End = into.Loc.End
		}
	}
	return nil
}

// settleQueryInto validates an INTO OUTFILE already consumed inside a parsed
// query relation. A SELECT may own the query-level suffix; a set operation
// takes one from its rightmost arm; a parenthesized subquery may not carry one.
func (p *Parser) settleQueryInto(query ast.Node) error {
	switch q := query.(type) {
	case *ast.SetOpStmt:
		return p.settleTraceSetOpInto(q)
	case *ast.ParenSelect:
		if queryHasInto(q) {
			return p.syntaxErrorAtCur()
		}
	}
	return nil
}

// startsSetOp reports whether kind starts a set-operation tail.
func startsSetOp(kind TokenKind) bool {
	switch kind {
	case kwUNION, kwINTERSECT, kwEXCEPT, kwMINUS:
		return true
	}
	return false
}

// parseTraceQuery parses the queryRelation optimizerTrace wraps: a SELECT,
// WITH ... query, or parenthesized query, with any set operations attached.
// EXPLAIN takes any statement, but TRACE takes a query only.
func (p *Parser) parseTraceQuery() (ast.Node, error) {
	switch p.cur.Kind {
	case kwSELECT:
		selectStmt, err := p.parseSelectStmt()
		if err != nil {
			return nil, err
		}
		if selectStmt.Into != nil && startsSetOp(p.cur.Kind) {
			return nil, p.syntaxErrorAtCur()
		}
		query, err := p.parseSetOpTail(selectStmt)
		if err != nil {
			return nil, err
		}
		if err := p.settleTraceSetOpInto(query); err != nil {
			return nil, err
		}
		return query, nil
	case kwWITH:
		with, err := p.parseWithClause()
		if err != nil {
			return nil, err
		}
		switch p.cur.Kind {
		case kwSELECT:
			selectStmt, err := p.parseSelectStmt()
			if err != nil {
				return nil, err
			}
			if selectStmt.Into != nil && startsSetOp(p.cur.Kind) {
				return nil, p.syntaxErrorAtCur()
			}
			query, err := p.parseSetOpTail(selectStmt)
			if err != nil {
				return nil, err
			}
			if err := p.settleTraceSetOpInto(query); err != nil {
				return nil, err
			}
			attachWith(query, with)
			return query, nil
		case int('('):
			paren, err := p.parseParenSelect()
			if err != nil {
				return nil, err
			}
			query, err := p.parseSetOpTail(paren)
			if err != nil {
				return nil, err
			}
			// An OUTFILE inside the parentheses belongs to a subquery, which
			// has no outfile clause; only the query-level suffix is allowed.
			if queryHasInto(query) {
				return nil, p.syntaxErrorAtCur()
			}
			attachWith(query, with)
			return query, nil
		}
		return nil, p.syntaxErrorAtCur()
	case int('('):
		paren, err := p.parseParenSelect()
		if err != nil {
			return nil, err
		}
		query, err := p.parseSetOpTail(paren)
		if err != nil {
			return nil, err
		}
		// OUTFILE inside the parentheses is a subquery suffix, not the
		// query-level one queryStatement allows.
		if queryHasInto(query) {
			return nil, p.syntaxErrorAtCur()
		}
		return query, nil
	}
	return nil, p.syntaxErrorAtCur()
}
