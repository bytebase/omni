package parser

import (
	"github.com/bytebase/omni/trino/ast"
)

// This file is part of the parser-ddl DAG node. It implements Trino's DROP
// statements for the schema-object family — TABLE, SCHEMA, VIEW, and
// MATERIALIZED VIEW. DROP CATALOG lives in catalog_ddl.go (it is dispatched
// from the CREATE/DROP CATALOG pair); DROP ROLE and DROP FUNCTION are owned by
// other nodes (grant_revoke.go / function_def.go) and are dispatched before
// this file is reached.
//
// Legacy ANTLR grammar (TrinoParser.g4) the implementation tracks:
//
//	statement
//	    : DROP_ TABLE_ (IF_ EXISTS_)? qualifiedName                  # dropTable
//	    | DROP_ SCHEMA_ (IF_ EXISTS_)? qualifiedName (CASCADE_ | RESTRICT_)? # dropSchema
//	    | DROP_ VIEW_ (IF_ EXISTS_)? qualifiedName                   # dropView
//	    | DROP_ MATERIALIZED_ VIEW_ (IF_ EXISTS_)? qualifiedName     # dropMaterializedView
//	    ;
//
// Adjudicated against the live Trino 481 oracle. Oracle-confirmed facts:
//
//	D-DR1 (object name ≤ catalog.schema.object). `DROP TABLE a.b.c.d` is a
//	   SYNTAX_ERROR for every object kind; the target is bounded to 3 parts.
//	D-DR2 (DROP SCHEMA CASCADE/RESTRICT). Only DROP SCHEMA carries the optional
//	   trailing CASCADE | RESTRICT; DROP TABLE/VIEW/MATERIALIZED VIEW do NOT
//	   (`DROP TABLE t CASCADE` is a SYNTAX_ERROR).

// DropObjectKind classifies which schema object a DROP targets.
type DropObjectKind int

const (
	// DropTable is DROP TABLE.
	DropTable DropObjectKind = iota
	// DropSchema is DROP SCHEMA.
	DropSchema
	// DropView is DROP VIEW.
	DropView
	// DropMaterializedView is DROP MATERIALIZED VIEW.
	DropMaterializedView
)

// DropBehavior is the trailing CASCADE / RESTRICT on DROP SCHEMA.
type DropBehavior int

const (
	// DropBehaviorDefault is no explicit CASCADE/RESTRICT (RESTRICT semantics).
	DropBehaviorDefault DropBehavior = iota
	// DropBehaviorCascade is CASCADE.
	DropBehaviorCascade
	// DropBehaviorRestrict is RESTRICT.
	DropBehaviorRestrict
)

// DropStmt is a `DROP { TABLE | SCHEMA | VIEW | MATERIALIZED VIEW }
// [IF EXISTS] name [CASCADE | RESTRICT]` statement. Behavior is only ever
// non-default for DROP SCHEMA (D-DR2).
type DropStmt struct {
	Object   DropObjectKind
	IfExists bool
	Name     *ast.QualifiedName
	Behavior DropBehavior
	Loc      ast.Loc
}

func (n *DropStmt) Tag() ast.NodeTag { return ast.T_Invalid }
func (n *DropStmt) Span() ast.Loc    { return n.Loc }

var _ ast.Node = (*DropStmt)(nil)

// parseDropStmt parses DROP TABLE/SCHEMA/VIEW/MATERIALIZED VIEW. On entry cur
// is the DROP keyword; startOffset is DROP's byte offset. DROP CATALOG, DROP
// ROLE and DROP FUNCTION are dispatched elsewhere before this is called, so the
// object keyword here is exactly one of TABLE / SCHEMA / VIEW / MATERIALIZED.
func (p *Parser) parseDropStmt(startOffset int) (ast.Node, error) {
	p.advance() // consume DROP

	var object DropObjectKind
	switch p.cur.Kind {
	case kwTABLE:
		p.advance()
		object = DropTable
	case kwSCHEMA:
		p.advance()
		object = DropSchema
	case kwVIEW:
		p.advance()
		object = DropView
	case kwMATERIALIZED:
		p.advance()
		if _, err := p.expect(kwVIEW); err != nil {
			return nil, err
		}
		object = DropMaterializedView
	default:
		// DROP with an object keyword this node does not own. The dispatch only
		// routes the four kinds above here, so reaching this is a syntax error.
		return nil, p.syntaxErrorAtCur()
	}

	ifExists, err := p.parseOptionalIfExists()
	if err != nil {
		return nil, err
	}

	// D-DR1: object name part-count limit. A SCHEMA is catalog.schema (≤ 2);
	// TABLE / VIEW / MATERIALIZED VIEW are catalog.schema.object (≤ 3). Trino
	// 481 enforces the per-object cap at parse time even though the legacy
	// grammar uses an unbounded qualifiedName everywhere.
	maxParts := 3
	context := "drop target"
	if object == DropSchema {
		maxParts = 2
		context = "schema name"
	}
	name, err := p.parseBoundedQualifiedName(maxParts, context)
	if err != nil {
		return nil, err
	}

	stmt := &DropStmt{
		Object:   object,
		IfExists: ifExists,
		Name:     name,
		Loc:      ast.Loc{Start: startOffset, End: name.Loc.End},
	}

	// D-DR2: only DROP SCHEMA accepts a trailing CASCADE | RESTRICT.
	if object == DropSchema {
		if tok, ok := p.match(kwCASCADE, kwRESTRICT); ok {
			if tok.Kind == kwCASCADE {
				stmt.Behavior = DropBehaviorCascade
			} else {
				stmt.Behavior = DropBehaviorRestrict
			}
			stmt.Loc.End = tok.Loc.End
		}
	}
	return stmt, nil
}
