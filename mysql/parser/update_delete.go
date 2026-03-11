package parser

import (
	nodes "github.com/bytebase/omni/mysql/ast"
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
	stmt.Tables = tables

	// SET clause (required)
	if _, err := p.expect(kwSET); err != nil {
		return nil, err
	}

	setList, err := p.parseAssignmentList()
	if err != nil {
		return nil, err
	}
	stmt.SetList = setList

	// WHERE clause (optional)
	if _, ok := p.match(kwWHERE); ok {
		where, err := p.parseExpr()
		if err != nil {
			return nil, err
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

	stmt.Loc.End = p.pos()
	return stmt, nil
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

	// WHERE clause (optional)
	if _, ok := p.match(kwWHERE); ok {
		where, err := p.parseExpr()
		if err != nil {
			return nil, err
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

	stmt.Loc.End = p.pos()
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
	start := p.pos()

	col, err := p.parseColumnRef()
	if err != nil {
		return nil, err
	}

	if _, err := p.expect('='); err != nil {
		return nil, err
	}

	val, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	return &nodes.Assignment{
		Loc:    nodes.Loc{Start: start, End: p.pos()},
		Column: col,
		Value:  val,
	}, nil
}

// parseDeleteTableList parses a comma-separated list of table names
// for DELETE statements (with optional .* suffix).
func (p *Parser) parseDeleteTableList() ([]nodes.TableExpr, error) {
	var tables []nodes.TableExpr

	for {
		ref, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}

		// Skip optional .* suffix (multi-table delete syntax)
		if p.cur.Type == '.' {
			next := p.peekNext()
			if next.Type == '*' {
				p.advance() // consume '.'
				p.advance() // consume '*'
			}
		}

		tables = append(tables, ref)
		if p.cur.Type != ',' {
			break
		}
		// Peek ahead to see if the next identifier is followed by FROM or USING,
		// which would mean this comma is NOT part of the table list.
		p.advance()
	}

	return tables, nil
}
