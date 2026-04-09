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

		// Single-table: optional [AS] alias and PARTITION
		if len(tables) == 1 {
			if tblRef, ok := tables[0].(*nodes.TableRef); ok {
				// Optional [AS] alias
				if _, ok := p.match(kwAS); ok {
					alias, _, err := p.parseIdent()
					if err != nil {
						return nil, err
					}
					tblRef.Alias = alias
					tblRef.Loc.End = p.pos()
				} else if p.cur.Type == tokIDENT {
					alias, _, _ := p.parseIdent()
					tblRef.Alias = alias
					tblRef.Loc.End = p.pos()
				}
				// Optional PARTITION (p0, p1, ...)
				if p.cur.Type == kwPARTITION {
					p.advance()
					parts, err := p.parseParenIdentList()
					if err != nil {
						return nil, err
					}
					tblRef.Partitions = parts
					tblRef.Loc.End = p.pos()
				}
			}
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
	// Completion: at the start of an assignment, offer columnref candidates.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		return nil, &ParseError{Message: "collecting"}
	}

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
	if val == nil {
		return nil, p.syntaxErrorAtCur()
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
			ref.Loc.End = p.pos()
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

	ref.Loc.End = p.pos()
	return ref, nil
}
