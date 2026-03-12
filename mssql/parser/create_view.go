// Package parser - create_view.go implements T-SQL CREATE VIEW statement parsing.
package parser

import (
	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseCreateViewStmt parses a CREATE [OR ALTER] VIEW statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-view-transact-sql
//
//	CREATE [ OR ALTER ] VIEW [ schema_name . ] view_name [ ( column [ ,...n ] ) ]
//	[ WITH <view_attribute> [ ,...n ] ]
//	AS select_statement
//	[ WITH CHECK OPTION ]
//
//	<view_attribute> ::=
//	{
//	    [ ENCRYPTION ]
//	    [ SCHEMABINDING ]
//	    [ VIEW_METADATA ]
//	}
func (p *Parser) parseCreateViewStmt(orAlter bool) *nodes.CreateViewStmt {
	loc := p.pos()

	stmt := &nodes.CreateViewStmt{
		OrAlter: orAlter,
		Loc:     nodes.Loc{Start: loc},
	}

	// View name
	stmt.Name = p.parseTableRef()

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
		stmt.Columns = &nodes.List{Items: cols}
	}

	// WITH <view_attribute> [,...n]
	if p.cur.Type == kwWITH {
		next := p.peekNext()
		if p.isRoutineOption(next) {
			p.advance() // consume WITH
			stmt.Options = p.parseRoutineOptionList()
			// Set SchemaBinding flag for backward compat
			if stmt.Options != nil {
				for _, item := range stmt.Options.Items {
					if s, ok := item.(*nodes.String); ok && s.Str == "SCHEMABINDING" {
						stmt.SchemaBinding = true
					}
				}
			}
		}
	}

	// AS
	p.match(kwAS)

	// SELECT query
	stmt.Query = p.parseSelectStmt()

	// WITH CHECK OPTION
	if p.cur.Type == kwWITH {
		next := p.peekNext()
		if next.Type == kwCHECK {
			p.advance() // WITH
			p.advance() // CHECK
			// OPTION
			if p.cur.Type == kwOPTION {
				p.advance()
			}
			stmt.WithCheck = true
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}
