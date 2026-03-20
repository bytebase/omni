// Package parser - create_view.go implements T-SQL CREATE VIEW statement parsing.
package parser

import (
	"strings"

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
	stmt.Name , _ = p.parseTableRef()

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
	stmt.Query, _ = p.parseSelectStmt()

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

// parseCreateMaterializedViewStmt parses a CREATE MATERIALIZED VIEW AS SELECT statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-materialized-view-as-select-transact-sql
//
//	CREATE MATERIALIZED VIEW [ schema_name. ] materialized_view_name
//	    WITH (
//	      <distribution_option>
//	      [, FOR_APPEND ]
//	    )
//	    AS <select_statement>
//
//	<distribution_option> ::=
//	    {
//	        DISTRIBUTION = HASH ( distribution_column_name [, ...n] )
//	      | DISTRIBUTION = ROUND_ROBIN
//	    }
func (p *Parser) parseCreateMaterializedViewStmt() *nodes.CreateMaterializedViewStmt {
	loc := p.pos()

	stmt := &nodes.CreateMaterializedViewStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// View name
	stmt.Name , _ = p.parseTableRef()

	// WITH ( <distribution_option> [, FOR_APPEND] )
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		p.expect('(')

		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			ident, ok := p.parseIdentifier()
			if !ok {
				break
			}
			switch {
			case strings.EqualFold(ident, "DISTRIBUTION"):
				p.expect('=')
				distType, ok := p.parseIdentifier()
				if !ok {
					break
				}
				if strings.EqualFold(distType, "HASH") {
					stmt.Distribution = "HASH"
					p.expect('(')
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
					p.expect(')')
					stmt.HashColumns = &nodes.List{Items: cols}
				} else if strings.EqualFold(distType, "ROUND_ROBIN") {
					stmt.Distribution = "ROUND_ROBIN"
				}
			case strings.EqualFold(ident, "FOR_APPEND"):
				stmt.ForAppend = true
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
		p.expect(')')
	}

	// AS
	p.match(kwAS)

	// SELECT query
	stmt.Query, _ = p.parseSelectStmt()

	stmt.Loc.End = p.pos()
	return stmt
}

// parseAlterMaterializedViewStmt parses an ALTER MATERIALIZED VIEW statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-materialized-view-transact-sql
//
//	ALTER MATERIALIZED VIEW [ schema_name. ] view_name
//	{
//	    REBUILD | DISABLE
//	}
func (p *Parser) parseAlterMaterializedViewStmt() *nodes.AlterMaterializedViewStmt {
	loc := p.pos()

	stmt := &nodes.AlterMaterializedViewStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// View name
	stmt.Name , _ = p.parseTableRef()

	// REBUILD | DISABLE
	action, ok := p.parseIdentifier()
	if ok {
		stmt.Action = action
	}

	stmt.Loc.End = p.pos()
	return stmt
}
