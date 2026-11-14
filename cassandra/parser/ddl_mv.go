package parser

import (
	"github.com/bytebase/omni/cassandra/ast"
)

// parseMVOptions parses viewProperty (AND viewProperty)* (Parser.g:904): a
// property or CLUSTERING ORDER BY (...). COMPACT STORAGE, which the grammar
// lists, its action always refuses.
func (p *Parser) parseMVOptions(stmt *ast.CreateMVStmt) error {
	props := map[string]bool{}
	orders := map[string]bool{}
	for {
		if p.cur.Type == tokCLUSTERING && p.peekNext().Type == tokORDER {
			co, err := p.parseClusteringOrderClause(orders)
			if err != nil {
				return err
			}
			stmt.ClusteringOrders = append(stmt.ClusteringOrders, co...)
		} else {
			opt, err := p.parseProperty(props)
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

// parseAlterMV parses alterMaterializedViewStatement (Parser.g:994):
// ALTER MATERIALIZED VIEW [IF EXISTS] v WITH properties.
func (p *Parser) parseAlterMV() (*ast.AlterMVStmt, error) {
	start := p.curLoc()

	ifExists := p.parseIfExists()

	name, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}

	stmt := &ast.AlterMVStmt{IfExists: ifExists, Name: name}

	if err := p.expectKeyword(tokWITH); err != nil {
		return nil, err
	}
	opts, err := p.parseProperties()
	if err != nil {
		return nil, err
	}
	stmt.Options = opts

	stmt.Loc = p.makeLoc(start)
	return stmt, nil
}

func (p *Parser) parseDropMV() (*ast.DropMVStmt, error) {
	start := p.curLoc()

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
