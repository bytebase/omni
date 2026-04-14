package parser

import (
	"github.com/bytebase/omni/doris/ast"
)

// parseCreateView parses a CREATE [OR REPLACE] VIEW statement.
// On entry, CREATE has already been consumed (p.prev is CREATE).
// If OR REPLACE were seen, they are already consumed; cur is VIEW.
//
// Syntax:
//
//	CREATE [OR REPLACE] VIEW [IF NOT EXISTS] view_name
//	    [(col1 [COMMENT 'text'], ...)]
//	    [COMMENT 'view comment']
//	    AS query
func (p *Parser) parseCreateView(startLoc ast.Loc, orReplace bool) (ast.Node, error) {
	// Consume VIEW keyword (cur is VIEW on entry)
	if _, err := p.expect(kwVIEW); err != nil {
		return nil, err
	}

	stmt := &ast.CreateViewStmt{
		OrReplace: orReplace,
	}

	// Optional IF NOT EXISTS
	if p.cur.Kind == kwIF {
		p.advance() // consume IF
		if _, err := p.expect(kwNOT); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfNotExists = true
	}

	// View name
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Optional column list: (col1 [COMMENT 'text'], ...)
	if p.cur.Kind == int('(') {
		cols, err := p.parseViewColumns()
		if err != nil {
			return nil, err
		}
		stmt.Columns = cols
	}

	// Optional COMMENT 'text'
	if p.cur.Kind == kwCOMMENT {
		p.advance() // consume COMMENT
		if p.cur.Kind != tokString {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Comment = p.cur.Str
		p.advance()
	}

	// AS keyword
	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}

	// Query — parse the SELECT statement
	query, err := p.parseSelectStmt()
	if err != nil {
		return nil, err
	}
	stmt.Query = query

	stmt.Loc = startLoc.Merge(ast.NodeLoc(query))
	return stmt, nil
}

// parseAlterView parses an ALTER VIEW statement.
// On entry, ALTER has already been consumed; cur is VIEW.
//
// Syntax:
//
//	ALTER VIEW view_name [(col1 [COMMENT 'text'], ...)] AS query
func (p *Parser) parseAlterView() (ast.Node, error) {
	startLoc := p.prev.Loc // loc of ALTER token

	// Consume VIEW keyword
	if _, err := p.expect(kwVIEW); err != nil {
		return nil, err
	}

	stmt := &ast.AlterViewStmt{}

	// View name
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Optional column list
	if p.cur.Kind == int('(') {
		cols, err := p.parseViewColumns()
		if err != nil {
			return nil, err
		}
		stmt.Columns = cols
	}

	// AS keyword
	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}

	// Query
	query, err := p.parseSelectStmt()
	if err != nil {
		return nil, err
	}
	stmt.Query = query

	stmt.Loc = startLoc.Merge(ast.NodeLoc(query))
	return stmt, nil
}

// parseDropView parses a DROP VIEW statement.
// On entry, DROP has already been consumed; cur is VIEW.
//
// Syntax:
//
//	DROP VIEW [IF EXISTS] view_name
func (p *Parser) parseDropView(startLoc ast.Loc) (ast.Node, error) {
	// Consume VIEW keyword
	if _, err := p.expect(kwVIEW); err != nil {
		return nil, err
	}

	stmt := &ast.DropViewStmt{}

	// Optional IF EXISTS
	if p.cur.Kind == kwIF {
		p.advance() // consume IF
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfExists = true
	}

	// View name
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	stmt.Loc = startLoc.Merge(ast.NodeLoc(name))
	return stmt, nil
}

// parseViewColumns parses the optional column list in CREATE/ALTER VIEW:
//
//	( col_name [COMMENT 'text'] [, col_name [COMMENT 'text']] ... )
//
// cur must be '(' on entry.
func (p *Parser) parseViewColumns() ([]*ast.ViewColumn, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}

	var cols []*ast.ViewColumn

	for p.cur.Kind != int(')') && p.cur.Kind != tokEOF {
		colStart := p.cur.Loc

		// Column name — accept any valid identifier (including non-reserved keywords)
		colName, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}

		col := &ast.ViewColumn{
			Name: colName,
			Loc:  colStart,
		}

		// Optional COMMENT 'text'
		if p.cur.Kind == kwCOMMENT {
			p.advance() // consume COMMENT
			if p.cur.Kind != tokString {
				return nil, p.syntaxErrorAtCur()
			}
			col.Comment = p.cur.Str
			col.Loc.End = p.cur.Loc.End
			p.advance()
		} else {
			col.Loc.End = colStart.End
		}

		cols = append(cols, col)

		// Optional comma
		if p.cur.Kind == int(',') {
			p.advance()
		}
	}

	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}

	return cols, nil
}
