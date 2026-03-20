// Package parser - rowset_functions.go implements T-SQL rowset function parsing.
// Rowset functions: OPENROWSET, OPENQUERY, OPENJSON, OPENDATASOURCE, OPENXML
// These appear as table sources in FROM clauses.
package parser

import (
	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseRowsetFunction parses OPENROWSET, OPENQUERY, OPENJSON, OPENDATASOURCE, OPENXML
// as table sources. These are treated as FuncCallExpr nodes with keyword names.
func (p *Parser) parseRowsetFunction() nodes.TableExpr {
	loc := p.pos()
	funcName := p.cur.Str
	p.advance() // consume the keyword

	// OPENDATASOURCE('provider', 'connstr').db.schema.table - special case
	if p.cur.Type == '(' {
		fc := &nodes.FuncCallExpr{
			Name: &nodes.TableRef{Object: funcName, Loc: nodes.Loc{Start: loc}},
			Loc:  nodes.Loc{Start: loc},
		}
		p.advance() // consume (

		if p.cur.Type != ')' {
			var args []nodes.Node
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				arg, _ := p.parseExpr()
				args = append(args, arg)
				if _, ok := p.match(','); !ok {
					break
				}
			}
			fc.Args = &nodes.List{Items: args}
		}
		_, _ = p.expect(')')

		// WITH clause for OPENJSON and OPENXML
		var withClause *nodes.List
		if p.cur.Type == kwWITH {
			next := p.peekNext()
			if next.Type == '(' {
				p.advance() // consume WITH
				p.advance() // consume (
				withClause = p.parseRowsetWithClause()
			}
		}

		alias := p.parseOptionalAlias()

		result := &nodes.AliasedTableRef{
			Table: fc,
			Alias: alias,
			Loc:   nodes.Loc{Start: loc},
		}
		_ = withClause // WITH clause columns are consumed but not stored in AST for now
		return result
	}

	// Shouldn't happen but handle gracefully
	return &nodes.AliasedTableRef{
		Table: &nodes.FuncCallExpr{
			Name: &nodes.TableRef{Object: funcName, Loc: nodes.Loc{Start: loc}},
			Loc:  nodes.Loc{Start: loc},
		},
		Loc: nodes.Loc{Start: loc},
	}
}

// parseRowsetWithClause parses the WITH (...) column definitions for OPENJSON/OPENXML.
// The opening '(' has already been consumed.
func (p *Parser) parseRowsetWithClause() *nodes.List {
	var cols []nodes.Node
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		// column_name data_type [path]
		colName, ok := p.parseIdentifier()
		if !ok {
			break
		}
		dt , _ := p.parseDataType()
		col := &nodes.ColumnDef{
			Name:     colName,
			DataType: dt,
		}
		// Optional path expression (string literal or XPath)
		if p.cur.Type != ',' && p.cur.Type != ')' {
			// Consume optional path/expression
			p.parseExpr()
		}
		cols = append(cols, col)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	_, _ = p.expect(')')
	return &nodes.List{Items: cols}
}
