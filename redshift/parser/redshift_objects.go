package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/redshift/ast"
)

func (p *Parser) redshiftCreateObjectType() string {
	switch {
	case p.redshiftWordEqual("datashare"):
		return "datashare"
	case p.cur.Type == EXTERNAL:
		next := p.peekNext()
		switch next.Type {
		case SCHEMA:
			return "external schema"
		case TABLE:
			return "external table"
		case VIEW:
			return "external view"
		case FUNCTION:
			return "external function"
		}
		if strings.EqualFold(next.Str, "model") {
			return "external model"
		}
		if strings.EqualFold(next.Str, "protected") {
			return "external protected view"
		}
	case p.redshiftWordEqual("masking") && strings.EqualFold(p.peekNext().Str, "policy"):
		return "masking policy"
	case p.redshiftWordEqual("rls") && strings.EqualFold(p.peekNext().Str, "policy"):
		return "rls policy"
	case p.cur.Type == IDENTITY_P && strings.EqualFold(p.peekNext().Str, "provider"):
		return "identity provider"
	case p.redshiftWordEqual("model"):
		return "model"
	case p.redshiftWordEqual("library"):
		return "library"
	}
	return ""
}

func isRedshiftCreateObjectLead(tok Token) bool {
	return strings.EqualFold(tok.Str, "datashare") ||
		tok.Type == EXTERNAL ||
		strings.EqualFold(tok.Str, "masking") ||
		strings.EqualFold(tok.Str, "rls") ||
		tok.Type == IDENTITY_P ||
		strings.EqualFold(tok.Str, "model") ||
		strings.EqualFold(tok.Str, "library")
}

func (p *Parser) redshiftAlterObjectType() string {
	switch {
	case p.redshiftWordEqual("datashare"):
		return "datashare"
	case p.cur.Type == EXTERNAL:
		next := p.peekNext()
		if next.Type == SCHEMA {
			return "external schema"
		}
		if next.Type == VIEW {
			return "external view"
		}
	case p.redshiftWordEqual("masking") && strings.EqualFold(p.peekNext().Str, "policy"):
		return "masking policy"
	case p.redshiftWordEqual("rls") && strings.EqualFold(p.peekNext().Str, "policy"):
		return "rls policy"
	case p.cur.Type == IDENTITY_P && strings.EqualFold(p.peekNext().Str, "provider"):
		return "identity provider"
	}
	return ""
}

func (p *Parser) redshiftDropObjectType() string {
	switch {
	case p.redshiftWordEqual("datashare"):
		return "datashare"
	case p.cur.Type == EXTERNAL && p.peekNext().Type == VIEW:
		return "external view"
	case p.cur.Type == IDENTITY_P && strings.EqualFold(p.peekNext().Str, "provider"):
		return "identity provider"
	case p.redshiftWordEqual("library"):
		return "library"
	case p.redshiftWordEqual("masking") && strings.EqualFold(p.peekNext().Str, "policy"):
		return "masking policy"
	case p.redshiftWordEqual("rls") && strings.EqualFold(p.peekNext().Str, "policy"):
		return "rls policy"
	case p.redshiftWordEqual("model"):
		return "model"
	}
	return ""
}

func (p *Parser) parseRedshiftObjectStmt(command, objectType string, loc int) (nodes.Node, error) {
	p.consumeRedshiftObjectType(objectType)

	if command == "drop" {
		p.parseOptIfExists()
	}

	name, err := p.parseOptionalRedshiftObjectName()
	if err != nil {
		return nil, err
	}
	options := p.parseRedshiftRawUtilityOptions()
	return &nodes.RedshiftObjectStmt{
		Command:    command,
		ObjectType: objectType,
		Name:       name,
		Options:    options,
		Loc:        nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

func (p *Parser) parseRedshiftAttachDetachStmt(command string, loc int) (nodes.Node, error) {
	objectType := ""
	if p.redshiftWordEqual("masking") && strings.EqualFold(p.peekNext().Str, "policy") {
		objectType = "masking policy"
	} else if p.redshiftWordEqual("rls") && strings.EqualFold(p.peekNext().Str, "policy") {
		objectType = "rls policy"
	} else {
		return nil, p.syntaxErrorAtCur()
	}
	p.consumeRedshiftObjectType(objectType)
	name, err := p.parseOptionalRedshiftObjectName()
	if err != nil {
		return nil, err
	}
	options := p.parseRedshiftRawUtilityOptions()
	return &nodes.RedshiftObjectStmt{
		Command:    command,
		ObjectType: objectType,
		Name:       name,
		Options:    options,
		Loc:        nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

func (p *Parser) parseRedshiftCancelStmt(loc int) (nodes.Node, error) {
	options := p.parseRedshiftRawUtilityOptions()
	if options == nil {
		return nil, p.syntaxErrorAtCur()
	}
	return &nodes.RedshiftObjectStmt{
		Command:    "cancel",
		ObjectType: "query",
		Options:    options,
		Loc:        nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

func (p *Parser) consumeRedshiftObjectType(objectType string) {
	for _, part := range strings.Split(objectType, " ") {
		if !p.consumeRedshiftWord(part) {
			p.advance()
		}
	}
}

func (p *Parser) parseOptionalRedshiftObjectName() (*nodes.List, error) {
	if p.cur.Type == 0 || p.cur.Type == ';' {
		return nil, nil
	}
	return p.parseAnyName()
}
