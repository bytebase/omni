// Package parser - insert.go implements T-SQL INSERT statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseInsertStmt parses an INSERT statement.
//
// BNF: mssql/parser/bnf/insert-transact-sql.bnf
//
//	[ WITH <common_table_expression> [ ,...n ] ]
//	INSERT
//	{
//	        [ TOP ( expression ) [ PERCENT ] ]
//	        [ INTO ]
//	        { <object> | rowset_function_limited
//	          [ WITH ( <Table_Hint_Limited> [ ...n ] ) ]
//	        }
//	    {
//	        [ ( column_list ) ]
//	        [ <OUTPUT Clause> ]
//	        { VALUES ( { DEFAULT | NULL | expression } [ ,...n ] ) [ ,...n ]
//	        | derived_table
//	        | execute_statement
//	        | <dml_table_source>
//	        | DEFAULT VALUES
//	        }
//	    }
//	}
//	[;]
func (p *Parser) parseInsertStmt() (*nodes.InsertStmt, error) {
	loc := p.pos()
	p.advance() // consume INSERT

	stmt := &nodes.InsertStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Optional TOP
	if p.cur.Type == kwTOP {
		top, err := p.parseTopClause()
		if err != nil {
			return nil, err
		}
		stmt.Top = top
	}

	// Optional INTO
	p.match(kwINTO)

	// Table name or @table_variable
	rel, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	stmt.Relation = rel

	// Optional WITH ( <Table_Hint_Limited> ) on target
	if p.cur.Type == kwWITH && p.peekNext().Type == '(' {
		hints, err := p.parseTableHints()
		if err != nil {
			return nil, err
		}
		stmt.Relation.Hints = hints
	}

	// Optional column list
	if p.cur.Type == '(' {
		p.advance()
		var cols []nodes.Node
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			colName, ok := p.parseIdentifier()
			if !ok {
				break
			}
			cols = append(cols, &nodes.String{Str: colName})
			if _, ok := p.match(','); !ok {
				break
			}
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		stmt.Cols = &nodes.List{Items: cols}
	}

	// OUTPUT clause (before source)
	if p.cur.Type == kwOUTPUT {
		oc, err := p.parseOutputClause()
		if err != nil {
			return nil, err
		}
		stmt.OutputClause = oc
	}

	// Source: VALUES, SELECT, EXEC, DEFAULT VALUES
	switch {
	case p.cur.Type == kwVALUES:
		vc, err := p.parseValuesClause()
		if err != nil {
			return nil, err
		}
		stmt.Source = vc
	case p.cur.Type == kwSELECT || p.cur.Type == kwWITH:
		sel, err := p.parseSelectStmt()
		if err != nil {
			return nil, err
		}
		stmt.Source = sel
	case p.cur.Type == kwEXEC || p.cur.Type == kwEXECUTE:
		execStmt, _ := p.parseExecStmt()
		stmt.Source = execStmt
	case p.cur.Type == kwDEFAULT:
		// DEFAULT VALUES
		defLoc := p.pos()
		p.advance()
		if p.cur.Type == kwVALUES {
			p.advance()
		}
		stmt.Source = &nodes.Literal{
			Type: nodes.LitDefault,
			Loc:  nodes.Loc{Start: defLoc},
		}
	}

	// Optional OPTION clause
	if p.cur.Type == kwOPTION {
		oc, err := p.parseOptionClause()
		if err != nil {
			return nil, err
		}
		stmt.OptionClause = oc
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseValuesClause parses VALUES (...), (...), ...
//
//	VALUES '(' expr_list ')' { ',' '(' expr_list ')' }
func (p *Parser) parseValuesClause() (*nodes.ValuesClause, error) {
	loc := p.pos()
	p.advance() // consume VALUES

	vc := &nodes.ValuesClause{
		Loc: nodes.Loc{Start: loc},
	}

	var rows []nodes.Node
	for {
		if _, err := p.expect('('); err != nil {
			break
		}
		row, err := p.parseExprList()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		rows = append(rows, row)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	vc.Rows = &nodes.List{Items: rows}
	vc.Loc.End = p.pos()
	return vc, nil
}

// parseOutputClause parses an OUTPUT clause.
//
// BNF: OUTPUT { dml_select_list } [ INTO { @table_variable | output_table } [ ( column_list ) ] ]
//
//	OUTPUT output_columns [INTO table [(col_list)]]
func (p *Parser) parseOutputClause() (*nodes.OutputClause, error) {
	loc := p.pos()
	p.advance() // consume OUTPUT

	oc := &nodes.OutputClause{
		Loc: nodes.Loc{Start: loc},
	}

	// Parse output targets (comma-separated expressions like inserted.col, deleted.col, $action)
	targets, err := p.parseTargetList()
	if err != nil {
		return nil, err
	}
	oc.Targets = targets

	// Optional INTO @table_variable | output_table
	if _, ok := p.match(kwINTO); ok {
		if p.cur.Type == tokVARIABLE {
			// @table_variable
			oc.IntoTable = &nodes.TableRef{
				Object: p.cur.Str,
				Loc:    nodes.Loc{Start: p.pos(), End: p.pos()},
			}
			p.advance()
		} else {
			tbl, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			oc.IntoTable = tbl
		}
		// Optional column list
		if p.cur.Type == '(' {
			p.advance()
			var cols []nodes.Node
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				colName, ok := p.parseIdentifier()
				if !ok {
					break
				}
				cols = append(cols, &nodes.String{Str: colName})
				if _, ok := p.match(','); !ok {
					break
				}
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
			oc.IntoCols = &nodes.List{Items: cols}
		}
	}

	oc.Loc.End = p.pos()
	return oc, nil
}

// isOutputKeyword checks if the current token is the OUTPUT keyword.
// This is used to detect OUTPUT clause in INSERT, UPDATE, DELETE.
func (p *Parser) isOutputKeyword() bool {
	return p.cur.Type == kwOUTPUT && !strings.EqualFold(p.cur.Str, "out")
}
