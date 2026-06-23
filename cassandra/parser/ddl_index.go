package parser

import (
	"github.com/bytebase/omni/cassandra/ast"
)

func (p *Parser) parseCreateIndex() (*ast.CreateIndexStmt, error) {
	start := p.curLoc()
	isCustom := false
	if p.cur.Type == tokCUSTOM {
		isCustom = true
		p.advance()
	}
	p.advance() // INDEX

	ifNotExists, err := p.parseIfNotExists()
	if err != nil {
		return nil, err
	}

	stmt := &ast.CreateIndexStmt{
		IsCustom:    isCustom,
		IfNotExists: ifNotExists,
	}

	// Optional index name before ON.
	if p.cur.Type != tokON {
		name, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.IndexName = name
	}

	if err := p.expectKeyword(tokON); err != nil {
		return nil, err
	}

	table, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	stmt.Table = table

	if _, err := p.expect(tokLPAREN); err != nil {
		return nil, err
	}

	// Index column: IDENT or FULL(IDENT), KEYS(IDENT), VALUES(IDENT), ENTRIES(IDENT)
	switch p.cur.Type {
	case tokFULL, tokKEYS, tokVALUES, tokENTRIES:
		fStart := p.curLoc()
		fname, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		fc, err := p.parseFunctionCallWithName(fname)
		if err != nil {
			return nil, err
		}
		fc.Loc.Start = fStart
		stmt.Column = fc
	default:
		col, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Column = col
	}

	if _, err := p.expect(tokRPAREN); err != nil {
		return nil, err
	}

	// Optional USING class
	if p.match(tokUSING) {
		switch p.cur.Type {
		case tokSTRING:
			stmt.UsingClass = &ast.StringLit{Val: p.cur.Str, Loc: ast.Loc{Start: p.cur.Loc, End: p.cur.End}}
			p.advance()
		case tokSAI:
			stmt.UsingClass = &ast.Identifier{Name: p.cur.Str, Loc: ast.Loc{Start: p.cur.Loc, End: p.cur.End}}
			p.advance()
		case tokSTORAGEATTACHEDINDEX:
			stmt.UsingClass = &ast.Identifier{Name: p.cur.Str, Loc: ast.Loc{Start: p.cur.Loc, End: p.cur.End}}
			p.advance()
		default:
			val, err := p.parseConstant()
			if err != nil {
				return nil, err
			}
			stmt.UsingClass = val
		}
	}

	// Optional WITH OPTIONS = { ... } or WITH { ... }
	if p.match(tokWITH) {
		if p.cur.Type == tokOPTIONS {
			p.advance()
			if _, err := p.expect(tokEQ); err != nil {
				return nil, err
			}
		}
		opts, err := p.parseOptionHash()
		if err != nil {
			return nil, err
		}
		stmt.Options = opts
	}

	stmt.Loc = p.makeLoc(start)
	return stmt, nil
}

func (p *Parser) parseDropIndex() (*ast.DropIndexStmt, error) {
	start := p.curLoc()
	p.advance() // INDEX

	ifExists := p.parseIfExists()

	name, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}

	return &ast.DropIndexStmt{
		IfExists: ifExists,
		Name:     name,
		Loc:      p.makeLoc(start),
	}, nil
}
