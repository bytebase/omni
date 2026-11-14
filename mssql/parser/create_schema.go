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
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()
	// Completion: after CREATE SCHEMA → identifier (new schema name)
	if p.collectMode() {
		p.addRuleCandidate("identifier")
		return nil, errCollecting
	}

	stmt := &nodes.CreateSchemaStmt{}

	// ( identifier authorizationOpt | authorization ) (createSchemaStatement,
	// TSql160.g:12878): a name, an owner, or both.
	if p.cur.Type == kwAUTHORIZATION {
		auth, err := p.schemaParseAuthorization()
		if err != nil {
			return nil, err
		}
		stmt.Authorization = auth
	} else {
		name, ok := p.parseIdentifier()
		if !ok {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Name = name
		if p.cur.Type == kwAUTHORIZATION {
			auth, err := p.schemaParseAuthorization()
			if err != nil {
				return nil, err
			}
			stmt.Authorization = auth
		}
	}

	// createSchemaElement (TSql160.g:12920): CREATE TABLE, CREATE VIEW,
	// GRANT, DENY or REVOKE, each followed by optional semicolons.
	var elements []nodes.Node
	for {
		var elem nodes.StmtNode
		var err error
		switch p.cur.Type {
		case kwCREATE:
			if next := p.peekNext().Type; next != kwTABLE && next != kwVIEW {
				return nil, p.syntaxErrorAtCur()
			}
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

// schemaParseAuthorization parses AUTHORIZATION owner_name, with
// AUTHORIZATION the current token.
func (p *Parser) schemaParseAuthorization() (string, error) {
	p.advance() // AUTHORIZATION
	owner, ok := p.parseIdentifier()
	if !ok {
		return "", p.syntaxErrorAtCur()
	}
	return owner, nil
}

// ddlIsWord reports whether the current token is the word word, written in
// upper case: its keyword token, or an identifier spelled so that is not
// bracketed or quoted. SQL Server reads an option word in its own text only.
func (p *Parser) ddlIsWord(word string) bool {
	if kw := lookupKeyword(word); kw != tokIDENT {
		return p.cur.Type == kw
	}
	return p.cur.Type == tokIDENT && !p.isDelimited(p.cur) && strings.ToUpper(p.cur.Str) == word
}

// ddlExpectWord consumes the word word or fails.
func (p *Parser) ddlExpectWord(word string) error {
	if !p.ddlIsWord(word) {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
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

	// alterSchemaStatement (TSql160.g:12650): every part is required, and
	// the securable has at most two parts (Msg 117).
	name, ok := p.parseIdentifier()
	if !ok {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Name = name

	if err := p.ddlExpectWord("TRANSFER"); err != nil {
		return nil, err
	}

	// Optional entity_type ::
	switch next := p.peekNext(); {
	case (p.cur.Type == kwOBJECT || p.cur.Type == kwTYPE) && next.Type == tokCOLONCOLON:
		stmt.TransferType = "OBJECT"
		if p.cur.Type == kwTYPE {
			stmt.TransferType = "TYPE"
		}
		p.advance()
		p.advance()
	case p.cur.Type == kwXML && next.Type == kwSCHEMA:
		p.advance() // XML
		p.advance() // SCHEMA
		if err := p.ddlExpectWord("COLLECTION"); err != nil {
			return nil, err
		}
		if _, err := p.expect(tokCOLONCOLON); err != nil {
			return nil, err
		}
		stmt.TransferType = "XML SCHEMA COLLECTION"
	}

	// securable_name: [ schema_name . ] name
	id, ok := p.parseIdentifier()
	if !ok {
		return nil, p.syntaxErrorAtCur()
	}
	securableParts := []string{id}
	if _, ok := p.match('.'); ok {
		part, ok := p.parseIdentifier()
		if !ok {
			return nil, p.syntaxErrorAtCur()
		}
		securableParts = append(securableParts, part)
	}
	stmt.TransferEntity = strings.Join(securableParts, ".")

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}
