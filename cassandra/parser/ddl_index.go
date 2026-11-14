package parser

import (
	"github.com/bytebase/omni/cassandra/ast"
)

// parseCreateIndex parses createIndexStatement (Parser.g:848):
//
//	CREATE [CUSTOM] INDEX [IF NOT EXISTS] [name] ON t '(' [indexIdent (',' indexIdent)*] ')'
//	[USING 'class'] [WITH properties]
func (p *Parser) parseCreateIndex() (*ast.CreateIndexStmt, error) {
	start := p.curLoc()
	isCustom := false
	if p.cur.Type == tokCUSTOM {
		isCustom = true
		p.advance()
	}
	if err := p.expectKeyword(tokINDEX); err != nil {
		return nil, err
	}

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
	if p.cur.Type != tokRPAREN {
		for {
			target, err := p.parseIndexIdent()
			if err != nil {
				return nil, err
			}
			stmt.Columns = append(stmt.Columns, target)
			if !p.match(tokCOMMA) {
				break
			}
		}
		stmt.Column = stmt.Columns[0]
	}
	if _, err := p.expect(tokRPAREN); err != nil {
		return nil, err
	}

	if p.match(tokUSING) {
		cls, err := p.parseStringLit()
		if err != nil {
			return nil, err
		}
		stmt.UsingClass = cls
	}

	if p.match(tokWITH) {
		props, err := p.parseProperties()
		if err != nil {
			return nil, err
		}
		stmt.Properties = props
		for _, o := range props {
			if h, ok := o.Value.(*ast.OptionHash); ok && identKey(o.Name) == "options" {
				stmt.Options = h
			}
		}
	}

	stmt.Loc = p.makeLoc(start)
	return stmt, nil
}

// parseIndexIdent parses indexIdent (Parser.g:862): a column, or VALUES,
// KEYS, ENTRIES or FULL of one. KEYS and VALUES are also column names unless
// '(' follows.
func (p *Parser) parseIndexIdent() (ast.ExprNode, error) {
	switch p.cur.Type {
	case tokFULL, tokKEYS, tokVALUES, tokENTRIES:
		if p.peekNext().Type != tokLPAREN {
			break
		}
		fStart := p.curLoc()
		fname := tokIdent(p.cur)
		p.advance() // FULL | KEYS | VALUES | ENTRIES
		p.advance() // (
		col, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(tokRPAREN); err != nil {
			return nil, err
		}
		return &ast.FunctionCall{Name: fname, Args: []ast.ExprNode{col}, Loc: p.makeLoc(fStart)}, nil
	}
	return p.parseIdentifier()
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
