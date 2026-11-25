package parser

import (
	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// INSERT statement parser
// ---------------------------------------------------------------------------

// parseInsertStmt is the top-level dispatcher for INSERT statements.
// It handles both single-table INSERT and multi-table INSERT ALL/FIRST.
//
//	INSERT [OVERWRITE] INTO ...            → parseInsertSingleStmt
//	INSERT [OVERWRITE] ALL ...             → parseInsertMultiStmt
//	INSERT [OVERWRITE] FIRST ...           → parseInsertMultiStmt
func (p *Parser) parseInsertStmt() (ast.Node, error) {
	insertTok, err := p.expect(kwINSERT)
	if err != nil {
		return nil, err
	}
	startLoc := insertTok.Loc.Start

	// Optional OVERWRITE
	overwrite := false
	if p.curIs(kwOVERWRITE) {
		p.advance()
		overwrite = true
	}

	// Dispatch: ALL or FIRST → multi-table INSERT
	if p.curIs(kwALL) || p.curIs(kwFIRST) {
		return p.parseInsertMultiBody(startLoc, overwrite)
	}

	// Default: single-table INSERT INTO ...
	return p.parseInsertSingleBody(startLoc, overwrite)
}

// parseInsertSingleBody parses the body of a single-table INSERT:
//
//	INTO table [(cols)] {VALUES (exprs)[, ...] | SELECT ...}
func (p *Parser) parseInsertSingleBody(startLoc int, overwrite bool) (*ast.InsertStmt, error) {
	if _, err := p.expect(kwINTO); err != nil {
		return nil, err
	}

	target, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}

	stmt := &ast.InsertStmt{
		Overwrite: overwrite,
		Target:    target,
		Loc:       ast.Loc{Start: startLoc},
	}

	// Optional column list: (col1, col2, ...)
	if p.curIs('(') && p.isColumnList() {
		p.advance() // consume '('
		cols, err := p.parseIdentList()
		if err != nil {
			return nil, err
		}
		stmt.Columns = cols
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	// VALUES or SELECT
	if p.curIs(kwVALUES) {
		values, err := p.parseValuesRows()
		if err != nil {
			return nil, err
		}
		stmt.Values = values
	} else {
		// SELECT / WITH SELECT
		query, err := p.parseInsertQuery()
		if err != nil {
			return nil, err
		}
		stmt.Select = query
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseInsertMultiBody parses the body of INSERT ALL / INSERT FIRST.
//
// Two forms:
//
//  1. Unconditional:
//     INSERT [OVERWRITE] ALL INTO t1 [(cols)] VALUES (exprs) [INTO t2 ...]
//     select_statement
//
//  2. Conditional (with WHEN):
//     INSERT [OVERWRITE] {ALL | FIRST}
//     WHEN cond THEN INTO t [(cols)] VALUES (exprs) [INTO t ...]
//     [WHEN ... THEN INTO ...]
//     [ELSE INTO t [(cols)] VALUES (exprs)]
//     select_statement
func (p *Parser) parseInsertMultiBody(startLoc int, overwrite bool) (*ast.InsertMultiStmt, error) {
	first := p.curIs(kwFIRST)
	p.advance() // consume ALL or FIRST

	stmt := &ast.InsertMultiStmt{
		Overwrite: overwrite,
		First:     first,
		Loc:       ast.Loc{Start: startLoc},
	}

	if p.curIs(kwWHEN) {
		// Conditional form: WHEN cond THEN INTO ... [ELSE INTO ...]
		for p.curIs(kwWHEN) {
			p.advance() // consume WHEN
			cond, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(kwTHEN); err != nil {
				return nil, err
			}
			// One or more INTO clauses under this WHEN
			for p.curIs(kwINTO) {
				branch, err := p.parseInsertMultiBranch(cond)
				if err != nil {
					return nil, err
				}
				stmt.Branches = append(stmt.Branches, branch)
			}
		}
		// Optional ELSE INTO ...
		if p.curIs(kwELSE) {
			p.advance() // consume ELSE
			for p.curIs(kwINTO) {
				branch, err := p.parseInsertMultiBranch(nil)
				if err != nil {
					return nil, err
				}
				stmt.Branches = append(stmt.Branches, branch)
			}
		}
	} else if p.curIs(kwINTO) {
		// Unconditional form: INTO t1 [(cols)] VALUES (exprs) [INTO t2 ...]
		for p.curIs(kwINTO) {
			branch, err := p.parseInsertMultiBranch(nil)
			if err != nil {
				return nil, err
			}
			stmt.Branches = append(stmt.Branches, branch)
		}
	}

	// Driving SELECT
	query, err := p.parseInsertQuery()
	if err != nil {
		return nil, err
	}
	stmt.Select = query
	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseInsertMultiBranch parses one INTO clause inside INSERT ALL/FIRST:
//
//	INTO table [(cols)] [VALUES (exprs)]
//
// The when parameter is the WHEN condition (nil for unconditional branches).
func (p *Parser) parseInsertMultiBranch(when ast.Node) (*ast.InsertMultiBranch, error) {
	intoTok, err := p.expect(kwINTO)
	if err != nil {
		return nil, err
	}

	target, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}

	branch := &ast.InsertMultiBranch{
		When:   when,
		Target: target,
		Loc:    ast.Loc{Start: intoTok.Loc.Start},
	}

	// Optional column list
	if p.curIs('(') && p.isColumnList() {
		p.advance() // consume '('
		cols, err := p.parseIdentList()
		if err != nil {
			return nil, err
		}
		branch.Columns = cols
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	// Optional VALUES (exprs)
	if p.curIs(kwVALUES) {
		p.advance() // consume VALUES
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		vals, err := p.parseExprList()
		if err != nil {
			return nil, err
		}
		branch.Values = vals
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	branch.Loc.End = p.prev.Loc.End
	return branch, nil
}

// parseValuesRows parses one or more VALUES rows:
//
//	VALUES (expr, ...) [, (expr, ...) ...]
func (p *Parser) parseValuesRows() ([][]ast.Node, error) {
	if _, err := p.expect(kwVALUES); err != nil {
		return nil, err
	}

	var rows [][]ast.Node

	row, err := p.parseValueRow()
	if err != nil {
		return nil, err
	}
	rows = append(rows, row)

	for p.curIs(',') {
		// Peek: is the next token '(' to start another row?
		if !p.nextIs('(') {
			break
		}
		p.advance() // consume ','
		row, err = p.parseValueRow()
		if err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}

	return rows, nil
}

// parseValueRow parses one parenthesized row: (expr, ...)
func (p *Parser) parseValueRow() ([]ast.Node, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	vals, err := p.parseExprList()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return vals, nil
}

// parseInsertQuery parses the driving SELECT for an INSERT statement.
// Supports both SELECT and WITH ... SELECT.
func (p *Parser) parseInsertQuery() (ast.Node, error) {
	if p.curIs(kwWITH) {
		return p.parseWithQueryExpr()
	}
	return p.parseQueryExpr()
}

// parseIdentList parses a comma-separated list of identifiers.
// Used for column lists in INSERT.
func (p *Parser) parseIdentList() ([]ast.Ident, error) {
	var idents []ast.Ident

	id, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	idents = append(idents, id)

	for p.curIs(',') {
		p.advance() // consume ','
		id, err = p.parseIdent()
		if err != nil {
			return nil, err
		}
		idents = append(idents, id)
	}

	return idents, nil
}

// isColumnList looks ahead to determine whether the current '(' starts a
// column list (identifiers separated by commas) as opposed to a subquery.
// It peeks at the next token: if it's a SELECT or WITH keyword, it's a
// subquery, not a column list.
func (p *Parser) isColumnList() bool {
	// If the first token inside '(' is SELECT or WITH, it's a subquery.
	return !p.nextIs(kwSELECT) && !p.nextIs(kwWITH)
}
