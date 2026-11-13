package parser

import (
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// parseUpdateStmt parses an UPDATE statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/update.html
//
// Single-table:
//
//	UPDATE [LOW_PRIORITY] [IGNORE] table_reference
//	    SET assignment_list
//	    [WHERE where_condition]
//	    [ORDER BY ...]
//	    [LIMIT row_count]
//
// Multi-table:
//
//	UPDATE [LOW_PRIORITY] [IGNORE] table_references
//	    SET assignment_list
//	    [WHERE where_condition]
func (p *Parser) parseUpdateStmt() (*nodes.UpdateStmt, error) {
	start := p.pos()
	p.advance() // consume UPDATE

	// Completion: after UPDATE keyword, offer table_ref candidates.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("table_ref")
		p.addRuleCandidate("database_ref")
		return nil, &ParseError{Message: "collecting"}
	}

	stmt := &nodes.UpdateStmt{Loc: nodes.Loc{Start: start}}

	// Optional modifiers
	if _, ok := p.match(kwLOW_PRIORITY); ok {
		stmt.LowPriority = true
	}
	if _, ok := p.match(kwIGNORE); ok {
		stmt.Ignore = true
	}

	// Parse table references
	tables, err := p.parseTableReferenceList()
	if err != nil {
		return nil, err
	}
	if err := refuseFromDuplicates(tables); err != nil {
		return nil, err
	}
	stmt.Tables = tables

	// FOR PORTION OF <period> FROM x TO y (application-time). Valid only on a
	// single bare base table (no join, no alias before the clause); for UPDATE
	// the alias, if any, follows the clause:
	//   UPDATE t FOR PORTION OF p FROM x TO y [AS] alias SET ...
	if p.cur.Type == kwFOR {
		// A FOR SYSTEM_TIME-bearing target cannot also take FOR PORTION OF
		// (MariaDB 11.8.8 rejects the combination); leave the trailing FOR to
		// the expect(kwSET) below, which reports 1064.
		if tref := singleBaseTable(stmt.Tables); tref != nil && tref.Alias == "" && tref.SystemTime == nil {
			fp, err := p.parseForPortionOf()
			if err != nil {
				return nil, err
			}
			stmt.ForPortionOf = fp
			if err := p.parseOptionalTableAlias(tref); err != nil {
				return nil, err
			}
		}
		// Otherwise the misplaced FOR is rejected by expect(kwSET) below (1064).
	}

	// SET clause (required)
	if _, err := p.expect(kwSET); err != nil {
		return nil, err
	}

	setList, err := p.parseAssignmentList()
	if err != nil {
		return nil, err
	}
	stmt.SetList = setList

	// Completion: after SET clause, offer WHERE/ORDER/LIMIT.
	p.checkCursor()
	if p.collectMode() {
		for _, tok := range []int{kwWHERE, kwORDER, kwLIMIT} {
			p.addTokenCandidate(tok)
		}
		return nil, &ParseError{Message: "collecting"}
	}

	// WHERE clause (optional)
	if _, ok := p.match(kwWHERE); ok {
		where, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if !p.collectMode() && where == nil {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Where = where
	}

	// ORDER BY clause (optional, single-table only)
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
	}

	// LIMIT clause (optional, single-table only)
	if _, ok := p.match(kwLIMIT); ok {
		limit, err := p.parseLimitClause()
		if err != nil {
			return nil, err
		}
		stmt.Limit = limit
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseForPortionOf parses the application-time clause
//
//	FOR PORTION OF <period> FROM <expr> TO <expr>
//
// on a single-table UPDATE / DELETE. FOR is the current token. The period name
// follows the identifier rule (non-reserved keywords allowed); SYSTEM_TIME is
// rejected (1064), matching MariaDB. PORTION is a non-reserved keyword.
func (p *Parser) parseForPortionOf() (*nodes.ForPortionOf, error) {
	start := p.pos()
	p.advance() // FOR
	if _, err := p.expect(kwPORTION); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwOF); err != nil {
		return nil, err
	}
	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	if strings.EqualFold(name, "SYSTEM_TIME") {
		return nil, p.syntaxErrorAtCur()
	}
	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}
	from, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if !p.collectMode() && from == nil {
		return nil, p.syntaxErrorAtCur()
	}
	if _, err := p.expect(kwTO); err != nil {
		return nil, err
	}
	to, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if !p.collectMode() && to == nil {
		return nil, p.syntaxErrorAtCur()
	}
	return &nodes.ForPortionOf{
		Loc:        nodes.Loc{Start: start, End: p.prev.End},
		PeriodName: name,
		From:       from,
		To:         to,
	}, nil
}

// singleBaseTable returns the table reference when tables is a single bare base
// table (not a join or comma-list), else nil.
func singleBaseTable(tables []nodes.TableExpr) *nodes.TableRef {
	if len(tables) != 1 {
		return nil
	}
	tref, _ := tables[0].(*nodes.TableRef)
	return tref
}

// parseOptionalTableAlias consumes an optional [AS] alias and records it on tref.
func (p *Parser) parseOptionalTableAlias(tref *nodes.TableRef) error {
	_, hasAS := p.match(kwAS)
	if !hasAS && !p.isIdentToken() {
		return nil
	}
	// WINDOW names a table anywhere but as its alias: `UPDATE t window SET a = 1` is a 1064.
	if p.cur.Type == kwWINDOW {
		return p.syntaxErrorAtCur()
	}
	alias, _, err := p.parseIdent()
	if err != nil {
		return err
	}
	tref.Alias = alias
	tref.Loc.End = p.prev.End
	return nil
}

// parseDeleteStmt parses a DELETE statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/delete.html
//
// Single-table:
//
//	DELETE [LOW_PRIORITY] [QUICK] [IGNORE] FROM tbl_name
//	    [WHERE where_condition]
//	    [ORDER BY ...]
//	    [LIMIT row_count]
//
// Multi-table (syntax 1):
//
//	DELETE [LOW_PRIORITY] [QUICK] [IGNORE]
//	    tbl_name[.*] [, tbl_name[.*]] ...
//	    FROM table_references
//	    [WHERE where_condition]
//
// Multi-table (syntax 2):
//
//	DELETE [LOW_PRIORITY] [QUICK] [IGNORE]
//	    FROM tbl_name[.*] [, tbl_name[.*]] ...
//	    USING table_references
//	    [WHERE where_condition]
func (p *Parser) parseDeleteStmt() (*nodes.DeleteStmt, error) {
	start := p.pos()
	p.advance() // consume DELETE

	stmt := &nodes.DeleteStmt{Loc: nodes.Loc{Start: start}}

	if p.cur.Type == kwHISTORY {
		return p.parseDeleteHistory(stmt)
	}

	// Optional modifiers
	if _, ok := p.match(kwLOW_PRIORITY); ok {
		stmt.LowPriority = true
	}
	if _, ok := p.match(kwQUICK); ok {
		stmt.Quick = true
	}
	if _, ok := p.match(kwIGNORE); ok {
		stmt.Ignore = true
	}

	if _, ok := p.match(kwFROM); ok {
		// Could be single-table: DELETE FROM tbl WHERE ...
		// or multi-table syntax 2: DELETE FROM t1, t2 USING ...
		tables, err := p.parseDeleteTableList()
		if err != nil {
			return nil, err
		}
		stmt.Tables = tables

		// Single-table: optional [AS | =] alias, index hints and PARTITION.
		// They belong to delete_single_table only: the USING form's targets are
		// a table_alias_ref_list, bare names.
		suffixed := false
		if len(tables) == 1 {
			if tblRef, ok := tables[0].(*nodes.TableRef); ok {
				at := p.pos()
				if err := p.parseDeleteSingleTableTail(tblRef); err != nil {
					return nil, err
				}
				suffixed = p.pos() != at
			}
		}
		if suffixed && p.cur.Type == kwUSING {
			return nil, p.syntaxErrorAtCur()
		}

		// Completion: after DELETE FROM table, offer WHERE/ORDER/LIMIT/USING.
		p.checkCursor()
		if p.collectMode() {
			for _, tok := range []int{kwWHERE, kwORDER, kwLIMIT, kwUSING} {
				p.addTokenCandidate(tok)
			}
			return nil, &ParseError{Message: "collecting"}
		}

		// Check for USING (multi-table syntax 2)
		if _, ok := p.match(kwUSING); ok {
			using, err := p.parseTableReferenceList()
			if err != nil {
				return nil, err
			}
			stmt.Using = using
		}
	} else {
		// Multi-table syntax 1: DELETE t1, t2 FROM table_references ...
		tables, err := p.parseDeleteTableList()
		if err != nil {
			return nil, err
		}
		stmt.Tables = tables

		if _, err := p.expect(kwFROM); err != nil {
			return nil, err
		}

		using, err := p.parseTableReferenceList()
		if err != nil {
			return nil, err
		}
		stmt.Using = using
	}

	if err := refuseFromDuplicates(stmt.Using); err != nil {
		return nil, err
	}
	if err := refuseDeleteTargets(stmt.Tables, stmt.Using); err != nil {
		return nil, err
	}

	// FOR PORTION OF <period> FROM x TO y (application-time, single-table only).
	if len(stmt.Tables) == 1 && len(stmt.Using) == 0 && p.cur.Type == kwFOR {
		fp, err := p.parseForPortionOf()
		if err != nil {
			return nil, err
		}
		stmt.ForPortionOf = fp
	}

	// WHERE clause (optional)
	if _, ok := p.match(kwWHERE); ok {
		where, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if !p.collectMode() && where == nil {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Where = where
	}

	// ORDER BY clause (optional, single-table only)
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
	}

	// LIMIT clause (optional, single-table only)
	if _, ok := p.match(kwLIMIT); ok {
		limit, err := p.parseLimitClause()
		if err != nil {
			return nil, err
		}
		stmt.Limit = limit
	}

	// RETURNING (MariaDB) — single-table DELETE only. Multi-table delete (either
	// syntax) has >1 table or USING; MariaDB rejects RETURNING there (1064), so
	// the clause is not reachable on those paths.
	if len(stmt.Tables) == 1 && stmt.Using == nil && p.cur.Type == kwRETURNING {
		ret, err := p.parseReturningClause()
		if err != nil {
			return nil, err
		}
		stmt.Returning = ret
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAssignmentList parses a comma-separated list of column=value assignments.
//
//	assignment_list:
//	    col_name = value [, col_name = value] ...
func (p *Parser) parseAssignmentList() ([]*nodes.Assignment, error) {
	var list []*nodes.Assignment

	for {
		asgn, err := p.parseAssignment()
		if err != nil {
			return nil, err
		}
		list = append(list, asgn)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	return list, nil
}

// parseAssignment parses a single col_name = value assignment.
func (p *Parser) parseAssignment() (*nodes.Assignment, error) {
	// Completion: at the start of an assignment, offer columnref candidates.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		return nil, &ParseError{Message: "collecting"}
	}

	start := p.pos()

	col, err := p.parseNamedColumn()
	if err != nil {
		return nil, err
	}

	if err := p.expectEqual(); err != nil {
		return nil, err
	}

	if ignore, ok := p.parseIgnoreValue(); ok {
		return &nodes.Assignment{
			Loc:    nodes.Loc{Start: start, End: p.prev.End},
			Column: col,
			Value:  ignore,
		}, nil
	}

	val, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, p.syntaxErrorAtCur()
	}

	return &nodes.Assignment{
		Loc:    nodes.Loc{Start: start, End: p.prev.End},
		Column: col,
		Value:  val,
	}, nil
}

// parseNamedColumn parses the column an assignment names. parseColumnRef
// also reads a user variable, which is no column (simple_ident_nospvar), and
// t.*, which names none; `@a` in backquotes is a column like any other.
func (p *Parser) parseNamedColumn() (*nodes.ColumnRef, error) {
	ref, err := p.parseColumnNotVariable()
	if err != nil {
		return nil, err
	}
	if ref.Star {
		return nil, &ParseError{Message: `syntax error at or near "*"`, Position: ref.Loc.End - 1}
	}
	return ref, nil
}

// parseColumnNotVariable parses a column reference that is not a user variable.
func (p *Parser) parseColumnNotVariable() (*nodes.ColumnRef, error) {
	if p.isUserVariable() {
		return nil, p.syntaxErrorAtCur()
	}
	return p.parseColumnRef()
}

// matchEqual consumes the = or := of an option, where the grammar makes it
// optional, and reports whether it found one. A position the grammar writes
// with a literal = takes no :=, such as a view's ALGORITHM, DEFINER, FORMAT, and
// every CHANGE REPLICATION option, so those read = alone.
//
// Ref: MariaDB 11.8 sql/sql_yacc.yy: opt_equal is empty or equal, and opt_with and
// opt_by read opt_equal too.
func (p *Parser) matchEqual() bool {
	if p.cur.Type == '=' || p.cur.Type == tokAssign {
		p.advance()
		return true
	}
	return false
}

// expectEqual consumes the = or := between a column and its value in an
// assignment list.
//
// Ref: MariaDB 11.8 sql/sql_yacc.yy: update_elem, insert_update_elem, ident_eq_value
// and load_data_set_elem all read a column, equal, and a value, and equal is
// '=' or SET_VAR.
func (p *Parser) expectEqual() error {
	if p.cur.Type == tokAssign {
		p.advance()
		return nil
	}
	_, err := p.expect('=')
	return err
}

// parseDeleteTableList parses a comma-separated list of table names
// for DELETE statements (with optional .* suffix).
func (p *Parser) parseDeleteTableList() ([]nodes.TableExpr, error) {
	var tables []nodes.TableExpr

	for {
		// Completion: at the start of each table name, offer table_ref candidates.
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("table_ref")
			p.addRuleCandidate("database_ref")
			return nil, &ParseError{Message: "collecting"}
		}

		ref, err := p.parseDeleteTableName()
		if err != nil {
			return nil, err
		}

		tables = append(tables, ref)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	return tables, nil
}

// parseDeleteTableName parses a single table name for DELETE,
// handling the optional .* suffix and schema.table.* form.
func (p *Parser) parseDeleteTableName() (*nodes.TableRef, error) {
	start := p.pos()
	dotted := p.peekNext().Type == '.'
	ref, err := p.parseUncheckedDeleteTableName()
	if err != nil {
		return nil, err
	}
	// name.* names no database, where `name`.t does, and an empty one is no Schema.
	qualified := ref.Schema != "" || (dotted && !strings.HasSuffix(p.inputText(start, p.prev.End), ".*"))
	return ref, refuseTableRefName(ref, qualified, start)
}

func (p *Parser) parseUncheckedDeleteTableName() (*nodes.TableRef, error) {
	start := p.pos()
	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}

	ref := &nodes.TableRef{
		Loc:  nodes.Loc{Start: start},
		Name: name,
	}

	// Check for dot: could be schema.table or name.*
	if p.cur.Type == '.' {
		next := p.peekNext()
		if next.Type == '*' {
			// name.* -- skip the .* suffix
			p.advance() // consume '.'
			p.advance() // consume '*'
			ref.Loc.End = p.prev.End
			return ref, nil
		}
		// schema.table or schema.table.*
		p.advance() // consume '.'
		// Completion: after "db.", offer table_ref.
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("table_ref")
			return nil, &ParseError{Message: "collecting"}
		}
		name2, _, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		ref.Schema = name
		ref.Name = name2

		// Check for .* after schema.table
		if p.cur.Type == '.' {
			next := p.peekNext()
			if next.Type == '*' {
				p.advance() // consume '.'
				p.advance() // consume '*'
			}
		}
	}

	ref.Loc.End = p.prev.End
	return ref, nil
}

// parseDeleteSingleTableTail reads what follows a single-table DELETE's table:
// opt_table_alias_clause, whose table_alias is AS, '=' or nothing,
// opt_key_definition (index hints) and opt_use_partition.
//
// Ref: delete_single_table and table_alias in MariaDB 11.8 sql/sql_yacc.yy.
func (p *Parser) parseDeleteSingleTableTail(tblRef *nodes.TableRef) error {
	written := false
	if _, ok := p.match(kwAS); ok {
		written = true
	} else if _, ok := p.match('='); ok {
		written = true
	}
	// table_alias reads an ident, which is any unreserved keyword as well as a
	// plain identifier, so `DELETE FROM t wait WHERE wait.a = 1` names an alias
	// `wait` (container-verified, as are the `locked` and `nowait` forms).
	if written || (p.isIdentToken() && p.cur.Type != kwWINDOW) {
		if p.cur.Type == kwWINDOW {
			return p.syntaxErrorAtCur()
		}
		alias, _, err := p.parseIdent()
		if err != nil {
			return err
		}
		tblRef.Alias = alias
		tblRef.Loc.End = p.prev.End
	}
	if (p.cur.Type == kwUSE || p.cur.Type == kwFORCE || p.cur.Type == kwIGNORE) &&
		(p.peekNext().Type == kwINDEX || p.peekNext().Type == kwKEY) {
		hints, err := p.parseIndexHints()
		if err != nil {
			return err
		}
		tblRef.IndexHints = hints
		tblRef.Loc.End = p.prev.End
	}
	if p.cur.Type == kwPARTITION {
		p.advance()
		parts, err := p.parseParenIdentList()
		if err != nil {
			return err
		}
		tblRef.Partitions = parts
		tblRef.Loc.End = p.prev.End
	}
	return nil
}

// parseDeleteHistory reads `DELETE HISTORY delete_single_table
// [BEFORE SYSTEM_TIME history_point]`, p.cur being HISTORY. history_point is
// TIMESTAMP 'string', TIMESTAMP(...), or [TRANSACTION | TIMESTAMP] bit_expr.
//
// Ref: delete_part2, opt_delete_system_time and history_point in MariaDB 11.8
// sql/sql_yacc.yy.
func (p *Parser) parseDeleteHistory(stmt *nodes.DeleteStmt) (*nodes.DeleteStmt, error) {
	p.advance() // consume HISTORY
	stmt.History = true
	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}
	ref, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	stmt.Tables = []nodes.TableExpr{ref}
	if err := p.parseDeleteSingleTableTail(ref); err != nil {
		return nil, err
	}
	if p.cur.Type == kwBEFORE {
		p.advance()
		// SYSTEM_TIME is a keyword this lexer does not know, written bare.
		if (p.cur.Type != tokIDENT && p.cur.Type < 700) || !strings.EqualFold(p.inputText(p.cur.Loc, p.cur.End), "SYSTEM_TIME") {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		point, err := p.parseHistoryPoint(stmt)
		if err != nil {
			return nil, err
		}
		stmt.HistoryBefore = point
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseHistoryPoint reads history_point: TIMESTAMP 'datetime', whose action
// refuses a string that is no datetime; TIMESTAMP(expr [, expr]); or
// [TRANSACTION | TIMESTAMP] bit_expr, the unit going to stmt.HistoryUnit.
//
// Ref: history_point and function_call_keyword_timestamp in MariaDB 11.8
// sql/sql_yacc.yy.
func (p *Parser) parseHistoryPoint(stmt *nodes.DeleteStmt) (nodes.ExprNode, error) {
	switch {
	case p.cur.Type == kwTIMESTAMP && p.peekNext().Type == tokSCONST:
		start := p.pos()
		p.advance() // consume TIMESTAMP
		if !historyDatetimeValid(p.cur.Str) {
			return nil, &ParseError{Message: "Incorrect DATETIME value: '" + p.cur.Str + "'", Position: p.cur.Loc}
		}
		val := p.advance().Str
		return &nodes.TemporalLit{Loc: nodes.Loc{Start: start, End: p.prev.End}, Type: "TIMESTAMP", Value: val}, nil
	case p.cur.Type == kwTIMESTAMP && p.peekNext().Type == '(':
		start := p.pos()
		p.advance() // consume TIMESTAMP
		p.advance() // consume '('
		fc := &nodes.FuncCallExpr{Loc: nodes.Loc{Start: start}, Name: "TIMESTAMP", HasParens: true}
		for {
			arg, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if arg == nil {
				return nil, p.syntaxErrorAtCur()
			}
			fc.Args = append(fc.Args, arg)
			if len(fc.Args) == 2 || p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		fc.Loc.End = p.prev.End
		return fc, nil
	case p.cur.Type == kwTIMESTAMP || p.cur.Type == kwTRANSACTION:
		stmt.HistoryUnit = strings.ToUpper(p.advance().Str)
	}
	point, err := p.parseBitExpr()
	if err != nil {
		return nil, err
	}
	if point == nil {
		return nil, p.syntaxErrorAtCur()
	}
	return point, nil
}

// historyDatetimeValid reports whether s reads as a DATETIME: digits alone
// (YYMMDD, YYYYMMDD, YYMMDDhhmmss or YYYYMMDDhhmmss), or a date of three numbers
// split by punctuation, then optionally a space or T and a time of two or three
// numbers split by punctuation and a fraction, each field in range. The
// server's reader takes a few rarer spellings too; this one refuses only text
// that is none of these.
func historyDatetimeValid(s string) bool {
	s = strings.TrimSpace(s)
	inRange := func(y, mo, d, h, mi, sec int) bool {
		return y <= 9999 && mo >= 1 && mo <= 12 && d >= 1 && d <= 31 && h <= 23 && mi <= 59 && sec <= 59
	}
	if _, err := strconv.ParseUint(s, 10, 64); err == nil {
		f := func(a, b int) int { v, _ := strconv.Atoi(s[a:b]); return v }
		switch len(s) {
		case 6:
			return inRange(f(0, 2), f(2, 4), f(4, 6), 0, 0, 0)
		case 8:
			return inRange(f(0, 4), f(4, 6), f(6, 8), 0, 0, 0)
		case 12:
			return inRange(f(0, 2), f(2, 4), f(4, 6), f(6, 8), f(8, 10), f(10, 12))
		case 14:
			return inRange(f(0, 4), f(4, 6), f(6, 8), f(8, 10), f(10, 12), f(12, 14))
		}
		return false
	}
	fields := func(part string) ([]int, bool) {
		var out []int
		cur := -1
		for i := 0; i < len(part); i++ {
			c := part[i]
			switch {
			case c >= '0' && c <= '9':
				if cur < 0 {
					cur = 0
				}
				cur = cur*10 + int(c-'0')
				if cur > 9999 {
					return nil, false
				}
			case cur >= 0 && (c == '-' || c == ':' || c == '/' || c == '.' || c == '_'):
				out = append(out, cur)
				cur = -1
			default:
				return nil, false
			}
		}
		if cur < 0 {
			return nil, false
		}
		return append(out, cur), true
	}
	date, clock := s, ""
	if i := strings.IndexAny(s, " T"); i >= 0 {
		date, clock = s[:i], strings.TrimSpace(s[i+1:])
	}
	d, ok := fields(date)
	if !ok || len(d) != 3 {
		return false
	}
	h, mi, sec := 0, 0, 0
	if clock != "" {
		if i := strings.IndexByte(clock, '.'); i >= 0 {
			if _, err := strconv.ParseUint(clock[i+1:], 10, 64); err != nil {
				return false
			}
			clock = clock[:i]
		}
		t, ok := fields(clock)
		if !ok || len(t) < 2 || len(t) > 3 {
			return false
		}
		h, mi = t[0], t[1]
		if len(t) == 3 {
			sec = t[2]
		}
	}
	return inRange(d[0], d[1], d[2], h, mi, sec)
}
