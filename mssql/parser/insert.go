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
func (p *Parser) parseInsertStmt() *nodes.InsertStmt {
	loc := p.pos()
	p.advance() // consume INSERT

	stmt := &nodes.InsertStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Optional TOP
	if p.cur.Type == kwTOP {
		stmt.Top, _ = p.parseTopClause()
	}

	// Optional INTO
	p.match(kwINTO)

	// Table name or @table_variable
	stmt.Relation, _ = p.parseTableRef()

	// Optional WITH ( <Table_Hint_Limited> ) on target
	if p.cur.Type == kwWITH && p.peekNext().Type == '(' {
		stmt.Relation.Hints, _ = p.parseTableHints()
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
		_, _ = p.expect(')')
		stmt.Cols = &nodes.List{Items: cols}
	}

	// OUTPUT clause (before source)
	if p.cur.Type == kwOUTPUT {
		stmt.OutputClause = p.parseOutputClause()
	}

	// Source: VALUES, SELECT, EXEC, DEFAULT VALUES
	switch {
	case p.cur.Type == kwVALUES:
		stmt.Source = p.parseValuesClause()
	case p.cur.Type == kwSELECT || p.cur.Type == kwWITH:
		stmt.Source, _ = p.parseSelectStmt()
	case p.cur.Type == kwEXEC || p.cur.Type == kwEXECUTE:
		stmt.Source = p.parseExecStmt()
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
		stmt.OptionClause, _ = p.parseOptionClause()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseValuesClause parses VALUES (...), (...), ...
//
//	VALUES '(' expr_list ')' { ',' '(' expr_list ')' }
func (p *Parser) parseValuesClause() *nodes.ValuesClause {
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
		row, _ := p.parseExprList()
		_, _ = p.expect(')')
		rows = append(rows, row)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	vc.Rows = &nodes.List{Items: rows}
	vc.Loc.End = p.pos()
	return vc
}

// parseOutputClause parses an OUTPUT clause.
//
// BNF: OUTPUT { dml_select_list } [ INTO { @table_variable | output_table } [ ( column_list ) ] ]
//
//	OUTPUT output_columns [INTO table [(col_list)]]
func (p *Parser) parseOutputClause() *nodes.OutputClause {
	loc := p.pos()
	p.advance() // consume OUTPUT

	oc := &nodes.OutputClause{
		Loc: nodes.Loc{Start: loc},
	}

	// Parse output targets (comma-separated expressions like inserted.col, deleted.col, $action)
	oc.Targets, _ = p.parseTargetList()

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
			oc.IntoTable , _ = p.parseTableRef()
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
			_, _ = p.expect(')')
			oc.IntoCols = &nodes.List{Items: cols}
		}
	}

	oc.Loc.End = p.pos()
	return oc
}

// isOutputKeyword checks if the current token is the OUTPUT keyword.
// This is used to detect OUTPUT clause in INSERT, UPDATE, DELETE.
func (p *Parser) isOutputKeyword() bool {
	return p.cur.Type == kwOUTPUT && !strings.EqualFold(p.cur.Str, "out")
}
