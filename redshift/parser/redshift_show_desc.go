package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/redshift/ast"
)

func (p *Parser) isRedshiftShowStart() bool {
	return p.redshiftWordEqual("databases") ||
		p.cur.Type == SCHEMAS ||
		p.cur.Type == TABLES ||
		p.cur.Type == TABLE ||
		(p.cur.Type == EXTERNAL && p.peekNext().Type == TABLE) ||
		p.cur.Type == COLUMNS ||
		p.redshiftWordEqual("grants") ||
		p.redshiftWordEqual("datashares") ||
		p.redshiftWordEqual("model") ||
		p.cur.Type == PROCEDURE ||
		p.cur.Type == VIEW
}

func (p *Parser) parseRedshiftShowStmt(loc int) (nodes.Node, error) {
	command := strings.ToLower(p.advance().Str)
	if command == "external" {
		if !p.consumeRedshiftWord("table") {
			return nil, p.syntaxErrorAtCur()
		}
		command = "external table"
	}
	options := p.parseRedshiftRawUtilityOptions()
	return &nodes.RedshiftShowStmt{
		Command: command,
		Options: options,
		Loc:     nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

func (p *Parser) parseRedshiftDescStmt(loc int) (nodes.Node, error) {
	objectType := strings.ToLower(p.advance().Str)
	if objectType == "identity" {
		if !p.consumeRedshiftWord("provider") {
			return nil, p.syntaxErrorAtCur()
		}
		objectType = "identity provider"
	}
	name, err := p.parseAnyName()
	if err != nil {
		return nil, err
	}
	if name == nil {
		return nil, p.syntaxErrorAtCur()
	}
	options := p.parseRedshiftRawUtilityOptions()
	return &nodes.RedshiftDescStmt{
		ObjectType: objectType,
		Name:       name,
		Options:    options,
		Loc:        nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

func (p *Parser) parseRedshiftRawUtilityOptions() *nodes.List {
	var items []nodes.Node
	for p.cur.Type != 0 && p.cur.Type != ';' {
		tok := p.advance()
		items = append(items, &nodes.String{Str: tok.Str})
	}
	if len(items) == 0 {
		return nil
	}
	return &nodes.List{Items: items}
}
