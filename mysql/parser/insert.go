package parser

import (
	nodes "github.com/bytebase/omni/mysql/ast"
)

// parseInsertStmt parses an INSERT statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/insert.html
//
//	INSERT [LOW_PRIORITY | DELAYED | HIGH_PRIORITY] [IGNORE]
//	    [INTO] tbl_name
//	    [(col_name [, col_name] ...)]
//	    { {VALUES | VALUE} (value_list) [, (value_list)] ...
//	      | SET col_name=value [, col_name=value] ...
//	      | SELECT ... }
//	    [ON DUPLICATE KEY UPDATE col_name=value [, col_name=value] ...]
func (p *Parser) parseInsertStmt() (*nodes.InsertStmt, error) {
	return p.parseInsertOrReplace(false)
}

// parseReplaceStmt parses a REPLACE statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/replace.html
//
//	REPLACE [LOW_PRIORITY | DELAYED]
//	    [INTO] tbl_name
//	    [(col_name [, col_name] ...)]
//	    { {VALUES | VALUE} (value_list) [, (value_list)] ...
//	      | SET col_name=value [, col_name=value] ...
//	      | SELECT ... }
func (p *Parser) parseReplaceStmt() (*nodes.InsertStmt, error) {
	return p.parseInsertOrReplace(true)
}

// parseInsertOrReplace is the shared implementation for INSERT and REPLACE.
func (p *Parser) parseInsertOrReplace(isReplace bool) (*nodes.InsertStmt, error) {
	start := p.pos()
	p.advance() // consume INSERT or REPLACE

	stmt := &nodes.InsertStmt{
		Loc:       nodes.Loc{Start: start},
		IsReplace: isReplace,
	}

	// Parse priority modifiers
	switch p.cur.Type {
	case kwLOW_PRIORITY:
		stmt.Priority = nodes.InsertPriorityLow
		p.advance()
	case kwDELAYED:
		stmt.Priority = nodes.InsertPriorityDelayed
		p.advance()
	case kwHIGH_PRIORITY:
		if !isReplace {
			stmt.Priority = nodes.InsertPriorityHigh
			p.advance()
		}
	}

	// Parse IGNORE. REPLACE takes none: replace_stmt has no opt_ignore.
	if !isReplace {
		if _, ok := p.match(kwIGNORE); ok {
			stmt.Ignore = true
		}
	}

	// Optional INTO
	p.match(kwINTO)

	// Completion: after INSERT [INTO], offer table_ref candidates.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("table_ref")
		p.addRuleCandidate("database_ref")
		return nil, &ParseError{Message: "collecting"}
	}

	// Parse table name
	tbl, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	stmt.Table = tbl

	// Optional PARTITION (p0, p1, ...)
	if p.cur.Type == kwPARTITION {
		p.advance()
		parts, err := p.parseParenIdentList()
		if err != nil {
			return nil, err
		}
		stmt.Partitions = parts
	}

	// Completion: after table name, offer VALUES, SET, SELECT, PARTITION keywords.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwVALUES)
		p.addTokenCandidate(kwSET)
		p.addTokenCandidate(kwSELECT)
		p.addTokenCandidate(kwPARTITION)
		return nil, &ParseError{Message: "collecting"}
	}

	// Optional column list: (col1, col2, ...), possibly empty.
	if p.cur.Type == '(' {
		// Distinguish a column list from a parenthesized query SOURCE. A column
		// list starts with a name or ')'; a parenthesized query source starts
		// with SELECT, TABLE, VALUES, WITH, or a nested '(' — none of which can begin
		// a column name. Such a source is left for the switch below, which parses
		// the query and any trailing set-op / ORDER BY / LIMIT, and is followed
		// like any source by ON DUPLICATE KEY UPDATE.
		switch p.peekNext().Type {
		case kwSELECT, kwTABLE, kwVALUES, kwWITH, '(':
		default:
			cols, err := p.parseInsertColumnList()
			if err != nil {
				return nil, err
			}
			stmt.Columns = cols
		}
	}

	// Parse the data source. MySQL 8.0 sql/sql_yacc.yy: insert_from_constructor,
	// insert_query_expression and the SET alternative of insert_stmt.
	//
	// Rows in parentheses follow VALUE or VALUES and may be empty. ROW(...) rows
	// follow only VALUES, and never mix with the others: VALUES ROW(...) is a
	// query, which takes a set operation or ORDER BY after it and no row alias.
	// A bare one still lands in Values, so that a reader finds the rows of a
	// plain INSERT in one place.
	aliasable := false
	switch p.cur.Type {
	case kwVALUE, kwVALUES:
		if p.cur.Type == kwVALUES && p.peekNext().Type == kwROW {
			if err := p.parseInsertValuesQuery(stmt); err != nil {
				return nil, err
			}
			break
		}
		p.advance() // consume VALUE or VALUES
		rows, err := p.parseValuesRows()
		if err != nil {
			return nil, err
		}
		stmt.Values = rows
		aliasable = true

	case kwSET:
		// SET names its own columns, so it follows the table directly.
		if stmt.Columns != nil {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance() // consume SET
		setList, err := p.parseAssignmentList()
		if err != nil {
			return nil, err
		}
		stmt.SetList = setList
		aliasable = true

	case kwSELECT, kwWITH, kwTABLE, '(':
		// SELECT / WITH / TABLE / parenthesized query source. A leading '('
		// reaches here with or without a column list before it, e.g.
		// INSERT INTO t (a) (SELECT 1) or INSERT INTO t ((SELECT 1) UNION ...).
		// TABLE is a query primary, so "INSERT INTO t TABLE t2 ORDER BY a" and
		// "INSERT INTO t TABLE t2 UNION TABLE t3" are the query expressions
		// they read as (container-verified).
		sel, err := p.parseSelectStmt()
		if err != nil {
			return nil, err
		}
		if err := rejectInsertSourceInto(sel); err != nil {
			return nil, err
		}
		stmt.Select = sel

	default:
		return nil, &ParseError{
			Message:  "expected VALUES, SET, or SELECT in INSERT statement",
			Position: p.cur.Loc,
		}
	}

	// AS row_alias [(col_alias, ...)] (MySQL 8.0.19+), only after rows in
	// parentheses or SET (opt_values_reference), and never on REPLACE.
	if aliasable && !isReplace && p.cur.Type == kwAS {
		p.advance()
		alias, _, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		stmt.RowAlias = alias
		if p.cur.Type == '(' {
			p.advance()
			var colAliases []string
			for {
				name, _, err := p.parseIdent()
				if err != nil {
					return nil, err
				}
				colAliases = append(colAliases, name)
				if p.cur.Type != ',' {
					break
				}
				p.advance()
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
			stmt.ColAliases = colAliases
		}
	}

	// ON DUPLICATE KEY UPDATE (only for INSERT, not REPLACE)
	if !isReplace && p.cur.Type == kwON {
		onDup, err := p.parseOnDuplicateKeyUpdate()
		if err != nil {
			return nil, err
		}
		stmt.OnDuplicateKey = onDup
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseInsertColumnList parses (col1, col2, ...) for INSERT. The list may be
// empty, and then comes back non-nil and empty, so that a reader can tell
// INSERT INTO t () VALUES () from INSERT INTO t VALUES ().
func (p *Parser) parseInsertColumnList() ([]*nodes.ColumnRef, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	// Completion: inside column list (including empty parens), offer columnref.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		return nil, &ParseError{Message: "collecting"}
	}

	if p.cur.Type == ')' {
		p.advance()
		return []*nodes.ColumnRef{}, nil
	}

	var cols []*nodes.ColumnRef
	for {
		// Completion: inside column list, offer columnref candidates.
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			return nil, &ParseError{Message: "collecting"}
		}

		ref, err := p.parseNamedColumn()
		if err != nil {
			return nil, err
		}
		cols = append(cols, ref)
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

// parseInsertValuesQuery parses a VALUES ROW(...) source, the current token
// being VALUES. It is a query expression, and goes in Select when anything
// follows the rows that makes it more than rows.
func (p *Parser) parseInsertValuesQuery(stmt *nodes.InsertStmt) error {
	sel, err := p.parseSelectStmt()
	if err != nil {
		return err
	}
	if err := rejectInsertSourceInto(sel); err != nil {
		return err
	}
	if sel != nil && sel.ValuesSource != nil && len(sel.OrderBy) == 0 && sel.Limit == nil && sel.ForUpdate == nil && sel.Into == nil {
		stmt.Values = sel.ValuesSource.Rows
		return nil
	}
	stmt.Select = sel
	return nil
}

// rejectInsertSourceInto refuses an INTO on a TABLE or VALUES primary used as
// an INSERT source, alone or as the last operand of a set operation. The query
// primary reads one as it does in a statement of its own, and MySQL 8.0 refuses it
// here (measured: 1064).
func rejectInsertSourceInto(sel *nodes.SelectStmt) error {
	for ; sel != nil; sel = sel.Right {
		if sel.Into != nil && (sel.TableSource != nil || sel.ValuesSource != nil) {
			return &ParseError{Message: `syntax error at or near "INTO"`, Position: sel.Into.Loc.Start}
		}
	}
	return nil
}

// parseValuesRows parses the rows of INSERT ... {VALUE | VALUES}: one or more
// parenthesized lists of values. A row may be empty, as in
// INSERT INTO t VALUES (), (). Rows written ROW(...) are a query, which
// parseValuesConstructor reads.
//
// Ref: MySQL 8.0 sql/sql_yacc.yy, insert_values is value_or_values values_list,
// and row_value is '(' opt_values ')'.
func (p *Parser) parseValuesRows() ([][]nodes.ExprNode, error) {
	var rows [][]nodes.ExprNode
	for {
		row, err := p.parseRowValue()
		if err != nil {
			return nil, err
		}
		rows = append(rows, row)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return rows, nil
}

// parseRowValue parses one parenthesized list of values, possibly empty.
func (p *Parser) parseRowValue() ([]nodes.ExprNode, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	var row []nodes.ExprNode
	// Completion: inside an empty row, parseExpr still offers its candidates.
	p.checkCursor()
	if p.cur.Type != ')' || p.collectMode() {
		for {
			e, err := p.parseValuesItem()
			if err != nil {
				return nil, err
			}
			row = append(row, e)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return row, nil
}

// parseValuesItem parses one value of a row: an expression or DEFAULT, which
// parseExpr reads.
func (p *Parser) parseValuesItem() (nodes.ExprNode, error) {
	return p.parseExprNotStar()
}

// parseExprNotStar parses an expression, refusing the * and t.* that parseExpr also
// reads, which only a select list takes.
func (p *Parser) parseExprNotStar() (nodes.ExprNode, error) {
	e, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, p.syntaxErrorAtCur()
	}
	if err := checkNotStar(e); err != nil {
		return nil, err
	}
	return e, nil
}

// checkNotStar refuses the stars parseExpr reads, which only a select list takes.
func checkNotStar(e nodes.ExprNode) error {
	if star, ok := e.(*nodes.StarExpr); ok {
		return &ParseError{Message: `syntax error at or near "*"`, Position: star.Loc.Start}
	}
	if ref, ok := e.(*nodes.ColumnRef); ok && ref.Star {
		return &ParseError{Message: `syntax error at or near "*"`, Position: ref.Loc.End - 1}
	}
	return nil
}

// parseOnDuplicateKeyUpdate parses ON DUPLICATE KEY UPDATE col=val, col=val, ...
func (p *Parser) parseOnDuplicateKeyUpdate() ([]*nodes.Assignment, error) {
	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwDUPLICATE); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwKEY); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwUPDATE); err != nil {
		return nil, err
	}
	return p.parseAssignmentList()
}
