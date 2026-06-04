package parser

import (
	"github.com/bytebase/omni/trino/ast"
)

// This file is part of the parser-ddl DAG node. It implements Trino's CREATE
// CATALOG and DROP CATALOG statements.
//
// Legacy ANTLR grammar (TrinoParser.g4) the implementation tracks:
//
//	| CREATE_ CATALOG_ (IF_ NOT_ EXISTS_)? catalog = identifier USING_ connectorName = identifier
//	    (COMMENT_ string_)? (AUTHORIZATION_ principal)? (WITH_ properties)? # createCatalog
//	| DROP_ CATALOG_ (IF_ EXISTS_)? catalog = identifier (CASCADE_ | RESTRICT_)? # dropCatalog
//
// Adjudicated against the live Trino 481 oracle. Oracle-confirmed facts:
//
//	D-CAT1 (catalog name is a SINGLE identifier). Unlike tables/views, a CREATE
//	   CATALOG / DROP CATALOG target is one identifier, not a dotted name;
//	   `CREATE CATALOG a.b USING …` is a SYNTAX_ERROR.
//	D-CAT2 (USING is mandatory on CREATE). `CREATE CATALOG c` without
//	   `USING connector` is a SYNTAX_ERROR.
//	D-CAT3 (connector name is a SINGLE identifier).
//
// NOTE on the docs vs grammar: ddl.md (drop-catalog) says "No IF EXISTS clause
// is documented" and shows no CASCADE/RESTRICT, but the legacy grammar (and
// the oracle) accept `DROP CATALOG [IF EXISTS] c [CASCADE | RESTRICT]`. The
// oracle is authoritative, so the optional clauses are accepted; the
// differential corpus pins this.

// ---------------------------------------------------------------------------
// CREATE CATALOG
// ---------------------------------------------------------------------------

// CreateCatalogStmt is `CREATE CATALOG [IF NOT EXISTS] name USING connector
// [COMMENT str] [AUTHORIZATION principal] [WITH (props)]`.
type CreateCatalogStmt struct {
	IfNotExists   bool
	Name          *ast.Identifier
	Connector     *ast.Identifier
	Comment       *string
	Authorization *Principal
	Properties    []*Property
	Loc           ast.Loc
}

func (n *CreateCatalogStmt) Tag() ast.NodeTag { return ast.T_Invalid }
func (n *CreateCatalogStmt) Span() ast.Loc    { return n.Loc }

// parseCreateCatalogStmt parses CREATE CATALOG. On entry cur is CREATE (dispatch
// peeked CATALOG).
func (p *Parser) parseCreateCatalogStmt(startOffset int) (ast.Node, error) {
	p.advance() // consume CREATE
	p.advance() // consume CATALOG
	ifNotExists, err := p.parseOptionalIfNotExists()
	if err != nil {
		return nil, err
	}
	// D-CAT1: catalog name is a single identifier.
	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	// D-CAT2: USING connector is mandatory.
	if _, err := p.expect(kwUSING); err != nil {
		return nil, err
	}
	// D-CAT3: connector name is a single identifier.
	connector, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	stmt := &CreateCatalogStmt{
		IfNotExists: ifNotExists,
		Name:        name,
		Connector:   connector,
		Loc:         ast.Loc{Start: startOffset, End: connector.Loc.End},
	}

	// Optional COMMENT, AUTHORIZATION, WITH — in the grammar's fixed order.
	comment, err := p.parseOptionalComment()
	if err != nil {
		return nil, err
	}
	if comment != nil {
		stmt.Comment = comment
		stmt.Loc.End = p.prev.Loc.End
	}
	if _, ok := p.match(kwAUTHORIZATION); ok {
		princ, err := p.parseAuthorizationPrincipal()
		if err != nil {
			return nil, err
		}
		stmt.Authorization = princ
		stmt.Loc.End = princ.Loc.End
	}
	props, err := p.parseOptionalWithProperties()
	if err != nil {
		return nil, err
	}
	if props != nil {
		stmt.Properties = props
		stmt.Loc.End = p.prev.Loc.End
	}
	return stmt, nil
}

// ---------------------------------------------------------------------------
// DROP CATALOG
// ---------------------------------------------------------------------------

// DropCatalogStmt is `DROP CATALOG [IF EXISTS] name [CASCADE | RESTRICT]`.
type DropCatalogStmt struct {
	IfExists bool
	Name     *ast.Identifier
	Behavior DropBehavior
	Loc      ast.Loc
}

func (n *DropCatalogStmt) Tag() ast.NodeTag { return ast.T_Invalid }
func (n *DropCatalogStmt) Span() ast.Loc    { return n.Loc }

// parseDropCatalogStmt parses DROP CATALOG. On entry cur is DROP (dispatch
// peeked CATALOG).
func (p *Parser) parseDropCatalogStmt(startOffset int) (ast.Node, error) {
	p.advance() // consume DROP
	p.advance() // consume CATALOG
	ifExists, err := p.parseOptionalIfExists()
	if err != nil {
		return nil, err
	}
	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	stmt := &DropCatalogStmt{IfExists: ifExists, Name: name, Loc: ast.Loc{Start: startOffset, End: name.Loc.End}}
	if tok, ok := p.match(kwCASCADE, kwRESTRICT); ok {
		if tok.Kind == kwCASCADE {
			stmt.Behavior = DropBehaviorCascade
		} else {
			stmt.Behavior = DropBehaviorRestrict
		}
		stmt.Loc.End = tok.Loc.End
	}
	return stmt, nil
}

var (
	_ ast.Node = (*CreateCatalogStmt)(nil)
	_ ast.Node = (*DropCatalogStmt)(nil)
)
