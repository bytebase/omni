package parser

import "github.com/bytebase/omni/snowflake/ast"

// ---------------------------------------------------------------------------
// INDEX (a secondary index on a hybrid table), read as an ast.ObjectStmt
// ---------------------------------------------------------------------------

func init() {
	registerObject(&objectForm{
		words:     []string{"INDEX"},
		create:    parseCreateIndex,
		orReplace: true,
		nameParts: 1, // the page writes <index_name>; the table qualifies it
	})
}

// parseCreateIndex reads what follows the name of
//
//	CREATE [ OR REPLACE ] INDEX [ IF NOT EXISTS ] <index_name>
//	  ON <table_name> ( <col_name> [ , <col_name> , ... ] )
//	  [ INCLUDE ( <col_name> [ , <col_name> , ... ] ) ]
//
// https://docs.snowflake.com/en/sql-reference/sql/create-index
//
// The ON clause holds the table as its name and the columns as one nameless
// option whose Group lists them; INCLUDE holds its columns the same way.
func parseCreateIndex(p *Parser, s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	if _, err := p.expect(kwON); err != nil {
		return err
	}
	table, err := p.parseObjectName()
	if err != nil {
		return err
	}
	cols, err := p.parseObjectColumnGroup()
	if err != nil {
		return err
	}
	c := p.addObjectKeywordsClause(s, "ON", start)
	c.Names = []*ast.ObjectName{table}
	c.Options = []*ast.CopyOption{cols}
	if p.cur.Type == kwINCLUDE {
		start := p.cur.Loc.Start
		p.advance() // consume INCLUDE
		cols, err := p.parseObjectColumnGroup()
		if err != nil {
			return err
		}
		p.addObjectKeywordsClause(s, "INCLUDE", start).Options = []*ast.CopyOption{cols}
	}
	return nil
}

// parseObjectColumnGroup reads ( <col_name> [ , ... ] ) into a nameless
// option whose Group holds one bare entry per column, named as written (a
// quoted name keeps its quotes). cur is '('.
func (p *Parser) parseObjectColumnGroup() (*ast.CopyOption, error) {
	open, err := p.expect('(')
	if err != nil {
		return nil, err
	}
	opt := &ast.CopyOption{Loc: ast.Loc{Start: open.Loc.Start}}
	for {
		col, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		entry := &ast.CopyOption{Name: col.Normalize(), Bare: true, Loc: col.Loc}
		if col.Quoted {
			entry.Name = col.String()
		}
		opt.Group = append(opt.Group, entry)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	opt.Loc.End = p.prev.Loc.End
	return opt, nil
}

// parseDropIndexStmt parses the tail of
//
//	DROP INDEX [ IF EXISTS ] <table_name>.<index_name>
//
// after INDEX (https://docs.snowflake.com/en/sql-reference/sql/drop-index).
// The table may be qualified, so the name has two to four parts: the last is
// the index, held in Name, and the rest the table, held in Table.
func (p *Parser) parseDropIndexStmt(start ast.Loc) (ast.Node, error) {
	stmt := &ast.DropStmt{Kind: ast.DropObject, Object: "INDEX", Loc: ast.Loc{Start: start.Start}}
	if err := p.parseIfExists(&stmt.IfExists); err != nil {
		return nil, err
	}
	var parts []ast.Ident
	for {
		part, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		parts = append(parts, part)
		if p.cur.Type != '.' || len(parts) == 4 {
			break
		}
		p.advance() // consume '.'
	}
	if len(parts) < 2 {
		return nil, p.syntaxErrorAtCur()
	}
	index := parts[len(parts)-1]
	tableParts := parts[:len(parts)-1]
	table := &ast.ObjectName{
		Name: tableParts[len(tableParts)-1],
		Loc:  ast.Loc{Start: tableParts[0].Loc.Start, End: tableParts[len(tableParts)-1].Loc.End},
	}
	switch len(tableParts) {
	case 3:
		table.Database, table.Schema = tableParts[0], tableParts[1]
	case 2:
		table.Schema = tableParts[0]
	}
	stmt.Table = table
	stmt.Name = &ast.ObjectName{Name: index, Loc: index.Loc}
	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}
