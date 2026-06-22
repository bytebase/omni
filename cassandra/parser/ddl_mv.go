package parser

import (
	"github.com/bytebase/omni/cassandra/ast"
)

func (p *Parser) parseMVOptions(stmt *ast.CreateMVStmt) error {
	for {
		if p.cur.Type == tokCLUSTERING {
			co, err := p.parseClusteringOrderClause()
			if err != nil {
				return err
			}
			stmt.ClusteringOrders = co
		} else {
			opt, err := p.parseTableOption()
			if err != nil {
				return err
			}
			stmt.Options = append(stmt.Options, opt)
		}
		if !p.match(tokAND) {
			break
		}
	}
	return nil
}

func (p *Parser) isColumnNotNull() bool {
	if !isIdentLike(p.cur.Type) {
		return false
	}
	// We need to look ahead: ident IS NOT NULL
	// This is tricky with limited lookahead. Use peekNext for the IS keyword.
	next := p.peekNext()
	return next.Type == tokIS
}

func (p *Parser) isNextColumnNotNull(next Token) bool {
	// After AND, check if the next token starts ident IS NOT NULL
	return isIdentLike(next.Type)
}

func (p *Parser) parseAlterMV() (*ast.AlterMVStmt, error) {
	start := p.curLoc()
	p.advance() // VIEW

	name, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}

	stmt := &ast.AlterMVStmt{Name: name}

	if p.match(tokWITH) {
		for {
			opt, err := p.parseTableOption()
			if err != nil {
				return nil, err
			}
			stmt.Options = append(stmt.Options, opt)
			if !p.match(tokAND) {
				break
			}
		}
	}

	stmt.Loc = p.makeLoc(start)
	return stmt, nil
}

func (p *Parser) parseDropMV() (*ast.DropMVStmt, error) {
	start := p.curLoc()
	p.advance() // VIEW

	ifExists := p.parseIfExists()

	name, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}

	return &ast.DropMVStmt{
		IfExists: ifExists,
		Name:     name,
		Loc:      p.makeLoc(start),
	}, nil
}
