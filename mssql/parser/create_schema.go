// Package parser - create_schema.go implements CREATE/ALTER SCHEMA parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseCreateSchemaStmt parses a CREATE SCHEMA statement.
//
// BNF: mssql/parser/bnf/create-schema-transact-sql.bnf
//
//	CREATE SCHEMA schema_name_clause [ <schema_element> [ ...n ] ]
//
//	<schema_name_clause> ::=
//	    { schema_name
//	    | AUTHORIZATION owner_name
//	    | schema_name AUTHORIZATION owner_name }
//
//	<schema_element> ::=
//	    { table_definition | view_definition | grant_statement |
//	      revoke_statement | deny_statement }
func (p *Parser) parseCreateSchemaStmt() (*nodes.CreateSchemaStmt, error) {
	// Completion: after CREATE SCHEMA → identifier (new schema name)
	if p.collectMode() {
		p.addRuleCandidate("identifier")
		return nil, errCollecting
	}

	stmt := &nodes.CreateSchemaStmt{}

	if p.cur.Type == kwAUTHORIZATION {
		p.advance()
		auth, _ := p.parseIdentifier()
		stmt.Authorization = auth
	} else {
		name, _ := p.parseIdentifier()
		stmt.Name = name
		if p.cur.Type == kwAUTHORIZATION {
			p.advance()
			auth, _ := p.parseIdentifier()
			stmt.Authorization = auth
		}
	}

	var elements []nodes.Node
	for {
		var elem nodes.StmtNode
		var err error
		switch p.cur.Type {
		case kwCREATE:
			elem, err = p.parseCreateStmt()
		case kwGRANT:
			elem, err = p.parseGrantStmt()
		case kwREVOKE:
			elem, err = p.parseRevokeStmt()
		case kwDENY:
			elem, err = p.parseDenyStmt()
		}
		if err != nil {
			return nil, err
		}
		if elem == nil {
			break
		}
		elements = append(elements, elem)
		for p.cur.Type == ';' {
			p.advance()
		}
	}
	if len(elements) > 0 {
		stmt.Elements = &nodes.List{Items: elements}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterSchemaStmt parses an ALTER SCHEMA statement.
//
// BNF: mssql/parser/bnf/alter-schema-transact-sql.bnf
//
//	ALTER SCHEMA schema_name TRANSFER [ <entity_type> :: ] securable_name
//
//	<entity_type> ::= { Object | Type | XML Schema Collection }
func (p *Parser) parseAlterSchemaStmt() (*nodes.AlterSchemaStmt, error) {
	stmt := &nodes.AlterSchemaStmt{}

	name, _ := p.parseIdentifier()
	stmt.Name = name

	p.matchIdentCI("TRANSFER")

	// Optional entity_type ::
	if p.isIdentLike() {
		next := p.peekNext()
		if next.Type == tokCOLONCOLON {
			switch p.cur.Type {
			case kwOBJECT:
				p.advance()
				p.advance()
				stmt.TransferType = "OBJECT"
			case kwTYPE:
				p.advance()
				p.advance()
				stmt.TransferType = "TYPE"
			}
		} else if p.cur.Type == kwXML && next.Type == kwSCHEMA {
			p.advance() // XML
			p.advance() // SCHEMA
			p.matchIdentCI("COLLECTION")
			p.match(tokCOLONCOLON)
			stmt.TransferType = "XML SCHEMA COLLECTION"
		}
	}

	// Parse securable_name (possibly schema-qualified)
	var securableParts []string
	id, ok := p.parseIdentifier()
	if ok {
		securableParts = append(securableParts, id)
		for p.cur.Type == '.' {
			p.advance()
			part, ok := p.parseIdentifier()
			if !ok {
				break
			}
			securableParts = append(securableParts, part)
		}
	}
	if len(securableParts) > 0 {
		stmt.TransferEntity = strings.Join(securableParts, ".")
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}
