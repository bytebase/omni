package parser

import (
	"github.com/bytebase/omni/cassandra/ast"
)

// parseCreateKeyspace parses createKeyspaceStatement (Parser.g:763):
// CREATE KEYSPACE [IF NOT EXISTS] ks WITH properties.
func (p *Parser) parseCreateKeyspace() (*ast.CreateKeyspaceStmt, error) {
	start := p.curLoc()
	p.advance() // KEYSPACE

	ifNotExists, err := p.parseIfNotExists()
	if err != nil {
		return nil, err
	}

	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	if err := p.expectKeyword(tokWITH); err != nil {
		return nil, err
	}
	opts, err := p.parseProperties()
	if err != nil {
		return nil, err
	}

	stmt := &ast.CreateKeyspaceStmt{
		IfNotExists: ifNotExists,
		Name:        name,
		Options:     opts,
	}
	stmt.Replication, stmt.DurableWrites = keyspaceOptions(opts)
	stmt.Loc = p.makeLoc(start)
	return stmt, nil
}

// parseAlterKeyspace parses alterKeyspaceStatement (Parser.g:939):
// ALTER KEYSPACE [IF EXISTS] ks WITH properties.
func (p *Parser) parseAlterKeyspace() (*ast.AlterKeyspaceStmt, error) {
	start := p.curLoc()
	p.advance() // KEYSPACE

	ifExists := p.parseIfExists()

	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	if err := p.expectKeyword(tokWITH); err != nil {
		return nil, err
	}
	opts, err := p.parseProperties()
	if err != nil {
		return nil, err
	}

	stmt := &ast.AlterKeyspaceStmt{IfExists: ifExists, Name: name, Options: opts}
	stmt.Replication, stmt.DurableWrites = keyspaceOptions(opts)
	stmt.Loc = p.makeLoc(start)
	return stmt, nil
}

// keyspaceOptions picks the replication map and the durable_writes boolean
// out of a keyspace's properties, where they are given in those forms.
func keyspaceOptions(opts []*ast.TableOption) (*ast.OptionHash, *ast.BoolLit) {
	var repl *ast.OptionHash
	var durable *ast.BoolLit
	for _, o := range opts {
		switch identKey(o.Name) {
		case "replication":
			if h, ok := o.Value.(*ast.OptionHash); ok {
				repl = h
			}
		case "durable_writes":
			if b, ok := o.Value.(*ast.BoolLit); ok {
				durable = b
			}
		}
	}
	return repl, durable
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
