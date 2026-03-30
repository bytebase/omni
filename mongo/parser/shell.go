package parser

import "github.com/bytebase/omni/mongo/ast"

// parseShowCommand parses "show dbs", "show databases", "show collections", etc.
func (p *Parser) parseShowCommand() (ast.Node, error) {
	showTok := p.advance() // consume 'show'
	start := showTok.Loc

	var target string
	switch p.cur.Type {
	case kwDbs:
		target = "databases"
		p.advance()
	case kwDatabases:
		target = "databases"
		p.advance()
	case kwCollections:
		target = "collections"
		p.advance()
	case kwTables:
		target = "tables"
		p.advance()
	case kwProfile:
		target = "profile"
		p.advance()
	case kwUsers:
		target = "users"
		p.advance()
	case kwRoles:
		target = "roles"
		p.advance()
	case kwLog:
		target = "log"
		p.advance()
	case kwLogs:
		target = "logs"
		p.advance()
	case kwStartupWarnings:
		target = "startupWarnings"
		p.advance()
	default:
		return nil, p.errorAtCur("expected show target (dbs, databases, collections, etc.)")
	}

	return &ast.ShowCommand{
		Target: target,
		Loc:    ast.Loc{Start: start, End: p.prev.End},
	}, nil
}
