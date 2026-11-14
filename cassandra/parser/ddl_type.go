package parser

import (
	"github.com/bytebase/omni/cassandra/ast"
)

func (p *Parser) parseCreateType() (*ast.CreateTypeStmt, error) {
	start := p.curLoc()
	p.advance() // TYPE

	ifNotExists, err := p.parseIfNotExists()
	if err != nil {
		return nil, err
	}

	name, err := p.parseUserTypeQualifiedName()
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(tokLPAREN); err != nil {
		return nil, err
	}

	stmt := &ast.CreateTypeStmt{
		IfNotExists: ifNotExists,
		Name:        name,
	}

	// typeColumns (',' typeColumns?)* (Parser.g:833): a trailing ',' is allowed.
	for first := true; ; first = false {
		if first || (p.cur.Type != tokCOMMA && p.cur.Type != tokRPAREN) {
			fStart := p.curLoc()
			col, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			dt, err := p.parseDataType()
			if err != nil {
				return nil, err
			}
			stmt.Fields = append(stmt.Fields, &ast.ColumnDef{
				Name: col,
				Type: dt,
				Loc:  p.makeLoc(fStart),
			})
		}
		if !p.match(tokCOMMA) {
			break
		}
	}

	if _, err := p.expect(tokRPAREN); err != nil {
		return nil, err
	}

	stmt.Loc = p.makeLoc(start)
	return stmt, nil
}

func (p *Parser) parseAlterType() (*ast.AlterTypeStmt, error) {
	start := p.curLoc()
	p.advance() // TYPE

	ifExists := p.parseIfExists()

	name, err := p.parseUserTypeQualifiedName()
	if err != nil {
		return nil, err
	}

	stmt := &ast.AlterTypeStmt{IfExists: ifExists, Name: name}

	switch p.cur.Type {
	case tokALTER:
		stmt.Op = ast.AlterTypeAlter
		p.advance()
		col, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		if err := p.expectKeyword(tokTYPE); err != nil {
			return nil, err
		}
		dt, err := p.parseDataType()
		if err != nil {
			return nil, err
		}
		stmt.AlterColumn = col
		stmt.AlterType = dt

	case tokADD:
		stmt.Op = ast.AlterTypeAdd
		p.advance()
		ine, err := p.parseIfNotExists()
		if err != nil {
			return nil, err
		}
		stmt.AddIfNotExists = ine
		// ADD takes one field (Parser.g:1012).
		fStart := p.curLoc()
		col, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		dt, err := p.parseDataType()
		if err != nil {
			return nil, err
		}
		stmt.AddFields = append(stmt.AddFields, &ast.ColumnDef{
			Name: col,
			Type: dt,
			Loc:  p.makeLoc(fStart),
		})

	case tokRENAME:
		stmt.Op = ast.AlterTypeRename
		p.advance()
		stmt.RenameIfExists = p.parseIfExists()
		for {
			rStart := p.curLoc()
			from, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			if err := p.expectKeyword(tokTO); err != nil {
				return nil, err
			}
			to, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			stmt.Renames = append(stmt.Renames, &ast.AlterTypeRenameItem{
				From: from,
				To:   to,
				Loc:  p.makeLoc(rStart),
			})
			if !p.match(tokAND) {
				break
			}
		}

	default:
		return nil, p.errorf("expected ALTER, ADD, or RENAME after ALTER TYPE, got %s", p.tokenDesc())
	}

	stmt.Loc = p.makeLoc(start)
	return stmt, nil
}

func (p *Parser) parseDropType() (*ast.DropTypeStmt, error) {
	start := p.curLoc()
	p.advance() // TYPE

	ifExists := p.parseIfExists()

	name, err := p.parseUserTypeQualifiedName()
	if err != nil {
		return nil, err
	}

	return &ast.DropTypeStmt{
		IfExists: ifExists,
		Name:     name,
		Loc:      p.makeLoc(start),
	}, nil
}

// parseUserTypeQualifiedName parses userTypeName (Parser.g:1489) as a
// QualifiedName.
func (p *Parser) parseUserTypeQualifiedName() (*ast.QualifiedName, error) {
	start := p.curLoc()
	ks, name, err := p.parseUserTypeName()
	if err != nil {
		return nil, err
	}
	parts := []*ast.Identifier{name}
	if ks != nil {
		parts = []*ast.Identifier{ks, name}
	}
	return &ast.QualifiedName{Parts: parts, Loc: p.makeLoc(start)}, nil
}
