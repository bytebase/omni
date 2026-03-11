// Package parser - create_index.go implements T-SQL CREATE INDEX statement parsing.
package parser

import (
	nodes "github.com/bytebase/omni/tsql/ast"
)

// parseCreateIndexStmt parses a CREATE INDEX statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-index-transact-sql
//
//	CREATE [UNIQUE] [CLUSTERED|NONCLUSTERED] INDEX name ON table (cols)
//	    [INCLUDE (cols)] [WHERE expr] [WITH (options)]
func (p *Parser) parseCreateIndexStmt(unique bool) *nodes.CreateIndexStmt {
	loc := p.pos()

	stmt := &nodes.CreateIndexStmt{
		Unique: unique,
		Loc:    nodes.Loc{Start: loc},
	}

	// CLUSTERED / NONCLUSTERED
	if p.cur.Type == kwCLUSTERED {
		p.advance()
		v := true
		stmt.Clustered = &v
	} else if p.cur.Type == kwNONCLUSTERED {
		p.advance()
		v := false
		stmt.Clustered = &v
	}

	// COLUMNSTORE (optional)
	if p.cur.Type == kwCOLUMNSTORE {
		p.advance()
		stmt.Columnstore = true
	}

	// INDEX keyword
	p.match(kwINDEX)

	// Index name
	name, _ := p.parseIdentifier()
	stmt.Name = name

	// ON table
	if _, ok := p.match(kwON); ok {
		stmt.Table = p.parseTableRef()
	}

	// Column list
	if p.cur.Type == '(' {
		stmt.Columns = p.parseIndexColumnList()
	}

	// INCLUDE
	if p.cur.Type == kwINCLUDE {
		p.advance()
		if p.cur.Type == '(' {
			stmt.IncludeCols = p.parseParenIdentList()
		}
	}

	// WHERE (filtered index)
	if _, ok := p.match(kwWHERE); ok {
		stmt.WhereClause = p.parseExpr()
	}

	// WITH (options)
	if p.cur.Type == kwWITH {
		p.advance()
		if p.cur.Type == '(' {
			stmt.Options = p.parseOptionList()
		}
	}

	// ON filegroup
	if _, ok := p.match(kwON); ok {
		if p.isIdentLike() {
			stmt.OnFileGroup = p.cur.Str
			p.advance()
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseIndexColumnList parses (col [ASC|DESC], ...).
func (p *Parser) parseIndexColumnList() *nodes.List {
	p.advance() // consume (
	var items []nodes.Node
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		loc := p.pos()
		name, ok := p.parseIdentifier()
		if !ok {
			break
		}
		dir := nodes.SortDefault
		if _, ok := p.match(kwASC); ok {
			dir = nodes.SortAsc
		} else if _, ok := p.match(kwDESC); ok {
			dir = nodes.SortDesc
		}
		items = append(items, &nodes.IndexColumn{
			Name:    name,
			SortDir: dir,
			Loc:     nodes.Loc{Start: loc},
		})
		if _, ok := p.match(','); !ok {
			break
		}
	}
	_, _ = p.expect(')')
	return &nodes.List{Items: items}
}

// parseOptionList parses (option = value, ...) used in WITH clauses.
func (p *Parser) parseOptionList() *nodes.List {
	p.advance() // consume (
	var items []nodes.Node
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		expr := p.parseExpr()
		if expr != nil {
			items = append(items, expr)
		}
		if _, ok := p.match(','); !ok {
			break
		}
	}
	_, _ = p.expect(')')
	return &nodes.List{Items: items}
}
