package parser

import (
	nodes "github.com/bytebase/omni/mariadb/ast"
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

	// Parse IGNORE
	if _, ok := p.match(kwIGNORE); ok {
		stmt.Ignore = true
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

	// Optional column list: (col1, col2, ...)
	//
	// Distinguish a column list from a parenthesized query SOURCE. A column
	// list always starts with an identifier (or ')'); a parenthesized query
	// source starts with SELECT, VALUES, WITH, or a nested '(' — none of which
	// can begin a column name. MariaDB's lexer decides the same way
	// (sql/sql_lex.cc, the left_paren case of MYSQLlex). So when the token
	// after '(' is one of those, this is INSERT ... (query) — possibly the left
	// operand of a set operation, e.g. (SELECT 1) UNION (SELECT 2) or
	// ((SELECT 1) UNION (SELECT 2)) — which the source switch below parses,
	// followed like any source by ON DUPLICATE KEY UPDATE and RETURNING.
	if p.cur.Type == '(' {
		switch p.peekNext().Type {
		case kwSELECT, kwVALUES, kwWITH, '(':
		default:
			cols, err := p.parseInsertColumnList()
			if err != nil {
				return nil, err
			}
			stmt.Columns = cols
		}
	}

	// Parse the data source: VALUES, SET, or a query.
	//
	// MariaDB parses VALUES rows as a table value constructor, the first
	// query primary of a query expression: insert_values (MariaDB 11.8.8
	// sql/sql_yacc.yy, line 13744) is create_select_query_expression. Rows
	// followed by a set operation or by ORDER BY, LIMIT, or locking are
	// therefore a query, as in INSERT INTO t VALUES (1) UNION VALUES (2), and
	// become the Select of the statement; rows alone stay its Values. Right
	// after the table or the column list MariaDB's lexer reads VALUE as VALUES
	// (sql/sql_lex.cc, BEFORE_OPT_LIST and AFTER_LIST), so VALUE starts the
	// rows too, but only these first ones. MySQL's ROW(...) is not a row there,
	// and MariaDB has no INSERT ... TABLE.
	switch p.cur.Type {
	case kwVALUES, kwVALUE:
		values, err := p.parseValuesConstructor()
		if err != nil {
			return nil, err
		}
		if p.atQueryExpressionTail() {
			sel, err := p.parseQueryAfterValues(values)
			if err != nil {
				return nil, err
			}
			stmt.Select = sel
		} else {
			stmt.Values = values.Rows
		}

	case kwSET:
		p.advance() // consume SET
		setList, err := p.parseAssignmentList()
		if err != nil {
			return nil, err
		}
		stmt.SetList = setList

	case kwSELECT, kwWITH, '(':
		// SELECT / WITH / parenthesized query source. A leading '(' reaches here
		// when an explicit column list preceded it, e.g.
		// INSERT INTO t (a) (SELECT 1) or INSERT INTO t (a) ((SELECT 1) UNION ...).
		sel, err := p.parseSelectStmt()
		if err != nil {
			return nil, err
		}
		stmt.Select = sel

	default:
		return nil, &ParseError{
			Message:  "expected VALUES, SET, or SELECT in INSERT statement",
			Position: p.cur.Loc,
		}
	}

	// MariaDB does not support the MySQL 8.0.19+ INSERT ... AS row_alias syntax
	// (it uses the VALUES(col) function for old-value references, not a row
	// alias) and rejects it at parse (1064). Reject explicitly rather than
	// leaving AS to be silently ignored or consumed by a later clause.
	if p.cur.Type == kwAS {
		return nil, &ParseError{
			Message:  "INSERT ... AS row alias is not supported in MariaDB",
			Position: p.cur.Loc,
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

	// RETURNING (MariaDB) — valid after every INSERT/REPLACE form.
	if p.cur.Type == kwRETURNING {
		ret, err := p.parseReturningClause()
		if err != nil {
			return nil, err
		}
		stmt.Returning = ret
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseInsertColumnList parses (col1, col2, ...) for INSERT.
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

	var cols []*nodes.ColumnRef
	for {
		// Completion: inside column list, offer columnref candidates.
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			return nil, &ParseError{Message: "collecting"}
		}

		ref, err := p.parseColumnRef()
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

// parseValuesRows parses the rows of INSERT ... VALUES and of a table value
// constructor: one or more parenthesized lists of values, each possibly
// empty, as in INSERT INTO t VALUES ().
//
// Ref: MariaDB 11.8.8 sql/sql_yacc.yy, values_list (line 13748) and
// no_braces (13788), a list of '(' opt_values ')'.
func (p *Parser) parseValuesRows() ([][]nodes.ExprNode, error) {
	var rows [][]nodes.ExprNode
	for {
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		var row []nodes.ExprNode
		if p.cur.Type != ')' {
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
		rows = append(rows, row)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return rows, nil
}

// parseValuesItem parses one value of a VALUES row: an expression or
// DEFAULT, which parseExpr reads. parseExpr also reads * and t.*, which only
// a select list takes.
func (p *Parser) parseValuesItem() (nodes.ExprNode, error) {
	e, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, p.syntaxErrorAtCur()
	}
	if star, ok := e.(*nodes.StarExpr); ok {
		return nil, &ParseError{Message: `syntax error at or near "*"`, Position: star.Loc.Start}
	}
	if ref, ok := e.(*nodes.ColumnRef); ok && ref.Star {
		return nil, &ParseError{Message: `syntax error at or near "*"`, Position: ref.Loc.End - 1}
	}
	return e, nil
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
	// Here VALUES is the VALUES(col) function, not a table value constructor
	// (opt_insert_update sets IN_UPDATE_ON_DUP_KEY, MariaDB 11.8.8
	// sql/sql_yacc.yy line 13888).
	saved := p.valuesIsFunction
	p.valuesIsFunction = true
	defer func() { p.valuesIsFunction = saved }()
	return p.parseAssignmentList()
}
