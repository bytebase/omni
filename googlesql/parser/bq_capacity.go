package parser

import (
	"github.com/bytebase/omni/googlesql/ast"
)

// ---------------------------------------------------------------------------
// BigQuery DDL — generic-entity CREATE / ALTER / DROP
// (parser-ddl-bigquery node)
// ---------------------------------------------------------------------------
//
// Ports the legacy ANTLR generic-entity mechanism (GoogleSQLParser.g4 §2.2/§2.3):
// create_entity_statement plus the generic-entity alternatives of alter_statement
// and drop_statement. This is how the legacy grammar parses extensible object
// kinds whose KEYWORD is a bare identifier — which is exactly how BigQuery
// CAPACITY / RESERVATION / ASSIGNMENT (DDL-024/025/026/053) are parsed: those
// words are NOT reserved tokens in the grammar, so they match
// `generic_entity_type: IDENTIFIER | PROJECT`. The node's display name calls out
// CAPACITY/RESERVATION/ASSIGNMENT specifically; implementing the generic-entity
// path covers them faithfully (the legacy grammar has no dedicated rule for any
// of them) plus any other extensible entity.
//
//	create_entity_statement:
//	  CREATE opt_or_replace? generic_entity_type opt_if_not_exists? path_expression
//	    opt_options_list? opt_generic_entity_body?
//	  opt_generic_entity_body: AS generic_entity_body
//	  generic_entity_body: json_literal | string_literal
//	  generic_entity_type: IDENTIFIER | PROJECT
//	alter_statement (generic entity):
//	  ALTER generic_entity_type opt_if_exists? path_expression alter_action_list
//	  | ALTER generic_entity_type opt_if_exists? alter_action_list
//	drop_statement (generic entity):
//	  DROP generic_entity_type opt_if_exists? path_expression
//
// ORACLE NOTE — BigQuery-only at the union level (oracle.md). Probed 2026-06-05:
// `CREATE CAPACITY \`p.l.c\` OPTIONS(…)` REJECTS on the Spanner emulator
// ("Encountered 'CAPACITY' while parsing: create_statement" — Spanner has no
// generic-entity mechanism). Triangulated against the legacy .g4 + BigQuery
// truth1 (DDL-024/025/026/053); proven by the unit tests in bq_capacity_test.go.

// isGenericEntityType reports whether tokenType can begin a generic_entity_type
// (`IDENTIFIER | PROJECT`). CAPACITY / RESERVATION / ASSIGNMENT lex as
// tokIdentifier (they are not reserved keywords), so they are admitted here.
func isGenericEntityType(tokenType int) bool {
	return tokenType == tokIdentifier || tokenType == kwPROJECT
}

// parseCreateEntity parses a CREATE <generic-entity> statement. The shared CREATE
// prefix (OR REPLACE) has been consumed by parseCreateStmt; cur is at the entity
// type token (a bare identifier or PROJECT). The grammar's create_entity_statement
// has no opt_create_scope (so a TEMP/etc. before the entity type would already
// have routed elsewhere).
func (p *Parser) parseCreateEntity(create Token, orReplace bool) (ast.Node, error) {
	entityTok := p.advance() // entity type
	entityType := p.tokenSource(entityTok)

	stmt := &ast.CreateEntityStmt{OrReplace: orReplace, EntityType: entityType}
	stmt.Loc.Start = create.Loc.Start

	ifNotExists, err := p.parseIfNotExists()
	if err != nil {
		return nil, err
	}
	stmt.IfNotExists = ifNotExists

	name, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// opt_options_list?
	if p.cur.Type == kwOPTIONS {
		opts, err := p.parseOptionsList()
		if err != nil {
			return nil, err
		}
		stmt.Options = opts
	}

	// opt_generic_entity_body? — AS generic_entity_body (json_literal | string_literal).
	if p.cur.Type == kwAS {
		p.advance() // AS
		body, err := p.parseGenericEntityBody()
		if err != nil {
			return nil, err
		}
		stmt.BodyText = body
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseGenericEntityBody parses a generic_entity_body (`json_literal |
// string_literal`) and returns its source text. json_literal is `JSON 'string'`;
// a bare string_literal is also accepted. Captured as text since the body is a
// free-form literal the query-span consumer never inspects.
func (p *Parser) parseGenericEntityBody() (string, error) {
	start := p.cur.Loc
	switch {
	case p.cur.Type == kwJSON:
		p.advance() // JSON
		strTok, err := p.expect(tokString)
		if err != nil {
			return "", err
		}
		return p.sourceSpan(start, strTok.Loc), nil
	case p.cur.Type == tokString:
		strTok := p.advance()
		return strTok.Str, nil
	default:
		return "", p.syntaxErrorAtCur()
	}
}

// sourceSpan returns the verbatim source text spanning from start.Start through
// end.End in the segment input, or "" if the span is out of range.
func (p *Parser) sourceSpan(start, end ast.Loc) string {
	s := absIndex(p, start.Start)
	e := absIndex(p, end.End)
	if s >= 0 && e <= len(p.input) && s < e {
		return p.input[s:e]
	}
	return ""
}

// parseBQAlterEntity parses an ALTER <generic-entity> statement. The ALTER keyword
// has been consumed by parseAlterStmt; cur is at the entity type token (a bare
// identifier or PROJECT). The grammar admits two shapes:
//
//	ALTER generic_entity_type opt_if_exists? path_expression alter_action_list
//	ALTER generic_entity_type opt_if_exists? alter_action_list           (no path)
//
// The alter_action of interest for these entities is SET OPTIONS(…) or
// SET AS generic_entity_body. We accept those; other alter_actions are out of the
// documented entity surface and reject as a syntax error.
func (p *Parser) parseBQAlterEntity(alter Token) (ast.Node, error) {
	entityTok := p.advance() // entity type
	entityType := p.tokenSource(entityTok)

	stmt := &ast.BQAlterStmt{Object: ast.BQAlterEntity, EntityType: entityType}
	stmt.Loc.Start = alter.Loc.Start

	ifExists, err := p.parseIfExists()
	if err != nil {
		return nil, err
	}
	stmt.IfExists = ifExists

	// Optional path_expression: present unless the next token starts the
	// alter_action_list directly (SET). A leading SET means the no-path shape.
	if p.cur.Type != kwSET {
		name, err := p.parseTablePath()
		if err != nil {
			return nil, err
		}
		stmt.Name = name
	}

	// alter_action_list — for entities, SET OPTIONS(…) | SET AS generic_entity_body.
	if _, err := p.expect(kwSET); err != nil {
		return nil, err
	}
	switch p.cur.Type {
	case kwOPTIONS:
		opts, err := p.parseOptionsList()
		if err != nil {
			return nil, err
		}
		stmt.SetOptions = opts
	case kwAS:
		p.advance() // AS
		body, err := p.parseGenericEntityBody()
		if err != nil {
			return nil, err
		}
		stmt.SetAsBody = body
	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseBQDropEntity parses a DROP <generic-entity> statement. The DROP keyword has
// been consumed by parseDropStmt; cur is at the entity type token (a bare
// identifier or PROJECT). Grammar: `DROP generic_entity_type opt_if_exists?
// path_expression` (DDL-053 DROP CAPACITY/RESERVATION/ASSIGNMENT).
func (p *Parser) parseBQDropEntity(drop Token) (ast.Node, error) {
	entityTok := p.advance() // entity type
	entityType := p.tokenSource(entityTok)

	stmt := &ast.BQDropStmt{Object: ast.BQDropEntity, EntityType: entityType}
	stmt.Loc.Start = drop.Loc.Start

	ifExists, err := p.parseIfExists()
	if err != nil {
		return nil, err
	}
	stmt.IfExists = ifExists

	name, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}
