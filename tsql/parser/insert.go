// Package parser - insert.go implements T-SQL INSERT statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/tsql/ast"
)

// parseInsertStmt parses an INSERT statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/insert-transact-sql
//
//	INSERT [INTO] table [(col_list)]
//	    { VALUES (...) [, (...)] | SELECT ... | EXEC ... | DEFAULT VALUES }
//	    [OUTPUT clause]
func (p *Parser) parseInsertStmt() *nodes.InsertStmt {
	loc := p.pos()
	p.advance() // consume INSERT

	stmt := &nodes.InsertStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Optional TOP
	if p.cur.Type == kwTOP {
		stmt.Top = p.parseTopClause()
	}

	// Optional INTO
	p.match(kwINTO)

	// Table name
	stmt.Relation = p.parseTableRef()

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
		stmt.Source = p.parseSelectStmt()
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
		row := p.parseExprList()
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
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/queries/output-clause-transact-sql
//
//	OUTPUT output_columns [INTO table [(col_list)]]
func (p *Parser) parseOutputClause() *nodes.OutputClause {
	loc := p.pos()
	p.advance() // consume OUTPUT

	oc := &nodes.OutputClause{
		Loc: nodes.Loc{Start: loc},
	}

	// Parse output targets (comma-separated expressions like inserted.col, deleted.col)
	oc.Targets = p.parseTargetList()

	// Optional INTO
	if _, ok := p.match(kwINTO); ok {
		oc.IntoTable = p.parseTableRef()
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
