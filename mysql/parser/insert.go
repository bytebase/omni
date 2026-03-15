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

	// Parse IGNORE
	if _, ok := p.match(kwIGNORE); ok {
		stmt.Ignore = true
	}

	// Optional INTO
	p.match(kwINTO)

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

	// Optional column list: (col1, col2, ...)
	if p.cur.Type == '(' {
		// Peek ahead to distinguish column list from VALUES (...) or subquery SELECT
		// If the next token after '(' is SELECT, this is INSERT ... (SELECT ...)
		next := p.peekNext()
		if next.Type == kwSELECT {
			// This is INSERT ... (SELECT ...)
			p.advance() // consume '('
			sel, err := p.parseSelectStmt()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
			stmt.Select = sel
			stmt.Loc.End = p.pos()
			return stmt, nil
		}
		// Otherwise, parse column list
		cols, err := p.parseInsertColumnList()
		if err != nil {
			return nil, err
		}
		stmt.Columns = cols
	}

	// Parse the data source: VALUES, SET, or SELECT
	switch p.cur.Type {
	case kwVALUES:
		p.advance() // consume VALUES
		rows, err := p.parseValuesRows()
		if err != nil {
			return nil, err
		}
		stmt.Values = rows

	case kwSET:
		p.advance() // consume SET
		setList, err := p.parseAssignmentList()
		if err != nil {
			return nil, err
		}
		stmt.SetList = setList

	case kwSELECT, kwWITH:
		sel, err := p.parseSelectStmt()
		if err != nil {
			return nil, err
		}
		stmt.Select = sel

	case kwTABLE:
		// TABLE table_name (MySQL 8.0.19+)
		tblStmt, err := p.parseTableStmt()
		if err != nil {
			return nil, err
		}
		stmt.TableSource = tblStmt

	default:
		// Also accept VALUE (MySQL alias for VALUES)
		if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "value") {
			p.advance()
			rows, err := p.parseValuesRows()
			if err != nil {
				return nil, err
			}
			stmt.Values = rows
		} else {
			return nil, &ParseError{
				Message:  "expected VALUES, SET, or SELECT in INSERT statement",
				Position: p.cur.Loc,
			}
		}
	}

	// AS row_alias [(col_alias, ...)] (MySQL 8.0.19+)
	if p.cur.Type == kwAS {
		p.advance()
		alias, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.RowAlias = alias
		if p.cur.Type == '(' {
			p.advance()
			var colAliases []string
			for {
				name, _, err := p.parseIdentifier()
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

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseInsertColumnList parses (col1, col2, ...) for INSERT.
func (p *Parser) parseInsertColumnList() ([]*nodes.ColumnRef, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	var cols []*nodes.ColumnRef
	for {
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

// parseValuesRows parses one or more value rows: (val_list), (val_list), ...
func (p *Parser) parseValuesRows() ([][]nodes.ExprNode, error) {
	var rows [][]nodes.ExprNode
	for {
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		row, err := p.parseExprList()
		if err != nil {
			return nil, err
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
