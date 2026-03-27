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
func (p *Parser) parseCreateViewStmt(orAlter bool) (*nodes.CreateViewStmt, error) {
	loc := p.pos()

	stmt := &nodes.CreateViewStmt{
		OrAlter: orAlter,
		Loc:     nodes.Loc{Start: loc},
	}

	// View name
	name, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

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
		stmt.Columns = &nodes.List{Items: cols}
	}

	// WITH <view_attribute> [,...n]
	if p.cur.Type == kwWITH {
		next := p.peekNext()
		if p.isRoutineOption(next) {
			p.advance() // consume WITH
			opts, err := p.parseRoutineOptionList()
			if err != nil {
				return nil, err
			}
			stmt.Options = opts
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
	query, err := p.parseSelectStmt()
	if err != nil {
		return nil, err
	}
	stmt.Query = query

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

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
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
func (p *Parser) parseCreateMaterializedViewStmt() (*nodes.CreateMaterializedViewStmt, error) {
	loc := p.pos()

	stmt := &nodes.CreateMaterializedViewStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// View name
	name, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// WITH ( <distribution_option> [, FOR_APPEND] )
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		if _, err := p.expect('('); err != nil {
			return nil, err
		}

		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			ident, ok := p.parseIdentifier()
			if !ok {
				break
			}
			switch {
			case strings.EqualFold(ident, "DISTRIBUTION"):
				if _, err := p.expect('='); err != nil {
					return nil, err
				}
				distType, ok := p.parseIdentifier()
				if !ok {
					break
				}
				if strings.EqualFold(distType, "HASH") {
					stmt.Distribution = "HASH"
					if _, err := p.expect('('); err != nil {
						return nil, err
					}
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
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	// AS
	p.match(kwAS)

	// SELECT query
	query, err := p.parseSelectStmt()
	if err != nil {
		return nil, err
	}
	stmt.Query = query

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterMaterializedViewStmt parses an ALTER MATERIALIZED VIEW statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-materialized-view-transact-sql
//
//	ALTER MATERIALIZED VIEW [ schema_name. ] view_name
//	{
//	    REBUILD | DISABLE
//	}
func (p *Parser) parseAlterMaterializedViewStmt() (*nodes.AlterMaterializedViewStmt, error) {
	loc := p.pos()

	stmt := &nodes.AlterMaterializedViewStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// View name
	name, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// REBUILD | DISABLE
	action, ok := p.parseIdentifier()
	if ok {
		stmt.Action = action
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}
