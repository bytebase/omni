package parser

import (
	"github.com/bytebase/omni/cassandra/ast"
)

func (p *Parser) parseCreateKeyspace() (*ast.CreateKeyspaceStmt, error) {
	start := p.curLoc()
	p.advance() // KEYSPACE

	ifNotExists := p.parseIfNotExists()

	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	if err := p.expectKeyword(tokWITH); err != nil {
		return nil, err
	}

	stmt := &ast.CreateKeyspaceStmt{
		IfNotExists: ifNotExists,
		Name:        name,
	}

	if err := p.expectKeyword(tokREPLICATION); err != nil {
		return nil, err
	}
	if _, err := p.expect(tokEQ); err != nil {
		return nil, err
	}
	repl, err := p.parseOptionHash()
	if err != nil {
		return nil, err
	}
	stmt.Replication = repl

	if p.match(tokAND) {
		if err := p.expectKeyword(tokDURABLE_WRITES); err != nil {
			return nil, err
		}
		if _, err := p.expect(tokEQ); err != nil {
			return nil, err
		}
		boolVal, err := p.parseBoolLit()
		if err != nil {
			return nil, err
		}
		stmt.DurableWrites = boolVal
	}

	stmt.Loc = p.makeLoc(start)
	return stmt, nil
}

func (p *Parser) parseAlterKeyspace() (*ast.AlterKeyspaceStmt, error) {
	start := p.curLoc()
	p.advance() // KEYSPACE

	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	stmt := &ast.AlterKeyspaceStmt{Name: name}

	if err := p.expectKeyword(tokWITH); err != nil {
		return nil, err
	}
	if err := p.expectKeyword(tokREPLICATION); err != nil {
		return nil, err
	}
	if _, err := p.expect(tokEQ); err != nil {
		return nil, err
	}
	repl, err := p.parseOptionHash()
	if err != nil {
		return nil, err
	}
	stmt.Replication = repl

	if p.match(tokAND) {
		if err := p.expectKeyword(tokDURABLE_WRITES); err != nil {
			return nil, err
		}
		if _, err := p.expect(tokEQ); err != nil {
			return nil, err
		}
		boolVal, err := p.parseBoolLit()
		if err != nil {
			return nil, err
		}
		stmt.DurableWrites = boolVal
	}

	stmt.Loc = p.makeLoc(start)
	return stmt, nil
}

func (p *Parser) parseDropKeyspace() (*ast.DropKeyspaceStmt, error) {
	start := p.curLoc()
	p.advance() // KEYSPACE

	ifExists := p.parseIfExists()

	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	return &ast.DropKeyspaceStmt{
		IfExists: ifExists,
		Name:     name,
		Loc:      p.makeLoc(start),
	}, nil
}

func (p *Parser) parseBoolLit() (*ast.BoolLit, error) {
	tok := p.cur
	switch tok.Type {
	case tokTRUE:
		p.advance()
		return &ast.BoolLit{Val: true, Loc: ast.Loc{Start: tok.Loc, End: tok.End}}, nil
	case tokFALSE:
		p.advance()
		return &ast.BoolLit{Val: false, Loc: ast.Loc{Start: tok.Loc, End: tok.End}}, nil
	default:
		return nil, p.errorf("expected TRUE or FALSE, got %s", p.tokenDesc())
	}
}
