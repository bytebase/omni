package parser

import (
	"strings"

	"github.com/bytebase/omni/trino/ast"
)

// This file is part of the parser-ddl DAG node. It implements Trino's CREATE /
// ALTER statements for schemas, views, and materialized views, plus REFRESH
// MATERIALIZED VIEW:
//
//	CREATE SCHEMA / ALTER SCHEMA (RENAME TO | SET AUTHORIZATION)
//	CREATE [OR REPLACE] VIEW … AS query / ALTER VIEW (RENAME | REFRESH | SET AUTHORIZATION)
//	CREATE [OR REPLACE] MATERIALIZED VIEW … AS query
//	ALTER MATERIALIZED VIEW (RENAME | SET PROPERTIES) / REFRESH MATERIALIZED VIEW
//
// DROP {SCHEMA,VIEW,MATERIALIZED VIEW} live in drop.go.
//
// Legacy ANTLR grammar (TrinoParser.g4) the implementation tracks:
//
//	| CREATE_ SCHEMA_ (IF_ NOT_ EXISTS_)? qualifiedName (AUTHORIZATION_ principal)? (WITH_ properties)? # createSchema
//	| ALTER_ SCHEMA_ qualifiedName RENAME_ TO_ identifier              # renameSchema
//	| ALTER_ SCHEMA_ qualifiedName SET_ AUTHORIZATION_ principal       # setSchemaAuthorization
//	| CREATE_ (OR_ REPLACE_)? MATERIALIZED_ VIEW_ (IF_ NOT_ EXISTS_)? qualifiedName
//	    (GRACE_ PERIOD_ interval)? (COMMENT_ string_)? (WITH_ properties)? AS_ rootQuery # createMaterializedView
//	| CREATE_ (OR_ REPLACE_)? VIEW_ qualifiedName (COMMENT_ string_)? (SECURITY_ (DEFINER_ | INVOKER_))? AS_ rootQuery # createView
//	| REFRESH_ MATERIALIZED_ VIEW_ qualifiedName                       # refreshMaterializedView
//	| ALTER_ MATERIALIZED_ VIEW_ (IF_ EXISTS_)? from RENAME_ TO_ to    # renameMaterializedView
//	| ALTER_ MATERIALIZED_ VIEW_ qualifiedName SET_ PROPERTIES_ propertyAssignments # setMaterializedViewProperties
//	| ALTER_ VIEW_ from RENAME_ TO_ to                                 # renameView
//	| ALTER_ VIEW_ from SET_ AUTHORIZATION_ principal                  # setViewAuthorization
//
// Adjudicated against the live Trino 481 oracle (which is ahead of the legacy
// grammar). Oracle-confirmed facts baked in:
//
//	D-MV1 (WHEN STALE + FIXED clause order). Trino 481 accepts `CREATE
//	   MATERIALIZED VIEW … WHEN STALE (INLINE | FAIL) …` (ddl.md
//	   create-materialized-view), absent from the legacy grammar.
//	   WHEN/STALE/INLINE/FAIL are NON-reserved (STALE/INLINE/FAIL are not even
//	   keyword tokens; they lex as plain identifiers), so STALE is matched by
//	   identifier text. The optional clauses appear in a FIXED order — GRACE
//	   PERIOD, then WHEN STALE, then COMMENT, then WITH, each at most once: 481
//	   REJECTS out-of-order clauses (`COMMENT 'c' GRACE PERIOD …`,
//	   `WITH (…) COMMENT …` are SYNTAX_ERRORs), so they are parsed as a strict
//	   sequence (verified against the oracle), not an order-tolerant set.
//	D-MV2 (object name ≤ 3 parts). Names are bounded to catalog.schema.object;
//	   the RENAME TO target is UNBOUNDED (over-bound failure is semantic — see
//	   D-AT5 in alter_table.go).
//	D-MV3 (SET AUTHORIZATION + IF EXISTS scoping). 481 accepts `ALTER MATERIALIZED
//	   VIEW name SET AUTHORIZATION principal` (ddl.md
//	   alter-materialized-view-set-authorization), absent from the legacy grammar
//	   — a docs-ahead extension. IF EXISTS is valid ONLY on the RENAME form (the
//	   legacy grammar's `(IF_ EXISTS_)?` is only on renameMaterializedView); 481
//	   rejects `ALTER MATERIALIZED VIEW IF EXISTS … SET …` as a SYNTAX_ERROR. Both
//	   facts are oracle-confirmed.
//	D-SC1 (CREATE SCHEMA name ≤ 2 parts). A schema is catalog.schema; `CREATE
//	   SCHEMA a.b.c` is a SYNTAX_ERROR. ALTER SCHEMA targets are likewise ≤ 2.
//	D-VW1 (CREATE VIEW has no IF NOT EXISTS). The grammar offers only `CREATE
//	   [OR REPLACE] VIEW` — there is no IF NOT EXISTS for a plain view, matching
//	   the docs (only OR REPLACE). MATERIALIZED VIEW does carry IF NOT EXISTS.

// ---------------------------------------------------------------------------
// Authorization target (shared: SCHEMA / VIEW / TABLE / MATERIALIZED VIEW)
// ---------------------------------------------------------------------------

// parseAuthorizationPrincipal parses the principal after AUTHORIZATION /
// SET AUTHORIZATION (`identifier | USER identifier | ROLE identifier`). It is a
// thin wrapper over the shared parsePrincipal (grant_revoke.go) so every DDL
// authorization clause resolves the principal identically.
func (p *Parser) parseAuthorizationPrincipal() (*Principal, error) {
	return p.parsePrincipal()
}

// ---------------------------------------------------------------------------
// CREATE SCHEMA
// ---------------------------------------------------------------------------

// CreateSchemaStmt is `CREATE SCHEMA [IF NOT EXISTS] name
// [AUTHORIZATION principal] [WITH (props)]`.
type CreateSchemaStmt struct {
	IfNotExists   bool
	Name          *ast.QualifiedName // ≤ 2 parts (catalog.schema)
	Authorization *Principal         // nil when no AUTHORIZATION clause
	Properties    []*Property
	Loc           ast.Loc
}

func (n *CreateSchemaStmt) Tag() ast.NodeTag { return ast.T_Invalid }
func (n *CreateSchemaStmt) Span() ast.Loc    { return n.Loc }

// parseCreateSchemaStmt parses CREATE SCHEMA. On entry cur is CREATE.
func (p *Parser) parseCreateSchemaStmt(startOffset int) (ast.Node, error) {
	p.advance() // consume CREATE
	if _, err := p.expect(kwSCHEMA); err != nil {
		return nil, err
	}
	ifNotExists, err := p.parseOptionalIfNotExists()
	if err != nil {
		return nil, err
	}
	// D-SC1: a schema name is catalog.schema (≤ 2 parts).
	name, err := p.parseBoundedQualifiedName(2, "schema name")
	if err != nil {
		return nil, err
	}
	stmt := &CreateSchemaStmt{IfNotExists: ifNotExists, Name: name, Loc: ast.Loc{Start: startOffset, End: name.Loc.End}}

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
// ALTER SCHEMA
// ---------------------------------------------------------------------------

// AlterSchemaKind classifies an ALTER SCHEMA sub-form.
type AlterSchemaKind int

const (
	// AlterSchemaRename is ALTER SCHEMA name RENAME TO new_name.
	AlterSchemaRename AlterSchemaKind = iota
	// AlterSchemaSetAuthorization is ALTER SCHEMA name SET AUTHORIZATION principal.
	AlterSchemaSetAuthorization
)

// AlterSchemaStmt is an ALTER SCHEMA statement. NewName is set for the rename
// form; Authorization for the set-authorization form.
type AlterSchemaStmt struct {
	Kind          AlterSchemaKind
	Name          *ast.QualifiedName
	NewName       *ast.Identifier // rename target (RENAME form)
	Authorization *Principal      // SET AUTHORIZATION form
	Loc           ast.Loc
}

func (n *AlterSchemaStmt) Tag() ast.NodeTag { return ast.T_Invalid }
func (n *AlterSchemaStmt) Span() ast.Loc    { return n.Loc }

// parseAlterSchemaStmt parses ALTER SCHEMA. On entry cur is ALTER (the dispatch
// has peeked SCHEMA as the second keyword).
func (p *Parser) parseAlterSchemaStmt(startOffset int) (ast.Node, error) {
	p.advance() // consume ALTER
	p.advance() // consume SCHEMA
	name, err := p.parseBoundedQualifiedName(2, "schema name")
	if err != nil {
		return nil, err
	}
	stmt := &AlterSchemaStmt{Name: name, Loc: ast.Loc{Start: startOffset}}

	switch p.cur.Kind {
	case kwRENAME:
		p.advance() // consume RENAME
		if _, err := p.expect(kwTO); err != nil {
			return nil, err
		}
		newName, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Kind = AlterSchemaRename
		stmt.NewName = newName
		stmt.Loc.End = newName.Loc.End
	case kwSET:
		p.advance() // consume SET
		if _, err := p.expect(kwAUTHORIZATION); err != nil {
			return nil, err
		}
		princ, err := p.parseAuthorizationPrincipal()
		if err != nil {
			return nil, err
		}
		stmt.Kind = AlterSchemaSetAuthorization
		stmt.Authorization = princ
		stmt.Loc.End = princ.Loc.End
	default:
		return nil, p.syntaxErrorAtCur()
	}
	return stmt, nil
}

// ---------------------------------------------------------------------------
// CREATE VIEW
// ---------------------------------------------------------------------------

// ViewSecurity is the SECURITY mode of a CREATE VIEW.
type ViewSecurity int

const (
	// ViewSecurityNone means no SECURITY clause was given.
	ViewSecurityNone ViewSecurity = iota
	// ViewSecurityDefiner is SECURITY DEFINER.
	ViewSecurityDefiner
	// ViewSecurityInvoker is SECURITY INVOKER.
	ViewSecurityInvoker
)

// CreateViewStmt is `CREATE [OR REPLACE] VIEW name [COMMENT str]
// [SECURITY {DEFINER | INVOKER}] AS query`.
type CreateViewStmt struct {
	OrReplace bool
	Name      *ast.QualifiedName
	Comment   *string
	Security  ViewSecurity
	Query     *Query
	Loc       ast.Loc
}

func (n *CreateViewStmt) Tag() ast.NodeTag { return ast.T_Invalid }
func (n *CreateViewStmt) Span() ast.Loc    { return n.Loc }

// parseCreateViewStmt parses CREATE [OR REPLACE] VIEW. On entry cur is CREATE.
func (p *Parser) parseCreateViewStmt(startOffset int) (ast.Node, error) {
	p.advance() // consume CREATE
	orReplace, err := p.parseOptionalOrReplace()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(kwVIEW); err != nil {
		return nil, err
	}
	// D-VW1: no IF NOT EXISTS for a plain view; the name follows directly.
	name, err := p.parseBoundedQualifiedName(3, "view name")
	if err != nil {
		return nil, err
	}
	stmt := &CreateViewStmt{OrReplace: orReplace, Name: name, Loc: ast.Loc{Start: startOffset}}

	comment, err := p.parseOptionalComment()
	if err != nil {
		return nil, err
	}
	stmt.Comment = comment

	if _, ok := p.match(kwSECURITY); ok {
		secTok, ok := p.match(kwDEFINER, kwINVOKER)
		if !ok {
			return nil, p.syntaxErrorAtCur()
		}
		if secTok.Kind == kwDEFINER {
			stmt.Security = ViewSecurityDefiner
		} else {
			stmt.Security = ViewSecurityInvoker
		}
	}

	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}
	query, err := p.parseQuery()
	if err != nil {
		return nil, err
	}
	stmt.Query = query
	stmt.Loc.End = query.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// ALTER VIEW
// ---------------------------------------------------------------------------

// AlterViewKind classifies an ALTER VIEW sub-form.
type AlterViewKind int

const (
	// AlterViewRename is ALTER VIEW name RENAME TO new_name.
	AlterViewRename AlterViewKind = iota
	// AlterViewSetAuthorization is ALTER VIEW name SET AUTHORIZATION principal.
	AlterViewSetAuthorization
	// AlterViewRefresh is ALTER VIEW name REFRESH (docs-ahead extension, D-AV1).
	AlterViewRefresh
)

// AlterViewStmt is an ALTER VIEW statement.
type AlterViewStmt struct {
	Kind          AlterViewKind
	Name          *ast.QualifiedName
	NewName       *ast.QualifiedName // RENAME target
	Authorization *Principal         // SET AUTHORIZATION form
	Loc           ast.Loc
}

func (n *AlterViewStmt) Tag() ast.NodeTag { return ast.T_Invalid }
func (n *AlterViewStmt) Span() ast.Loc    { return n.Loc }

// parseAlterViewStmt parses ALTER VIEW. On entry cur is ALTER (dispatch peeked
// VIEW). The docs (alter-view.html) add `ALTER VIEW name REFRESH` over the
// legacy grammar (D-AV1); the oracle adjudicates whether 481 accepts it.
func (p *Parser) parseAlterViewStmt(startOffset int) (ast.Node, error) {
	p.advance() // consume ALTER
	p.advance() // consume VIEW
	name, err := p.parseBoundedQualifiedName(3, "view name")
	if err != nil {
		return nil, err
	}
	stmt := &AlterViewStmt{Name: name, Loc: ast.Loc{Start: startOffset}}

	switch p.cur.Kind {
	case kwRENAME:
		p.advance()
		if _, err := p.expect(kwTO); err != nil {
			return nil, err
		}
		// The RENAME TO target is UNBOUNDED (oracle-confirmed, ledger DDL4): a
		// 4+-part target fails semantically (TABLE_NOT_FOUND), not at parse.
		newName, err := p.parseQualifiedName()
		if err != nil {
			return nil, err
		}
		stmt.Kind = AlterViewRename
		stmt.NewName = newName
		stmt.Loc.End = newName.Loc.End
	case kwSET:
		p.advance()
		if _, err := p.expect(kwAUTHORIZATION); err != nil {
			return nil, err
		}
		princ, err := p.parseAuthorizationPrincipal()
		if err != nil {
			return nil, err
		}
		stmt.Kind = AlterViewSetAuthorization
		stmt.Authorization = princ
		stmt.Loc.End = princ.Loc.End
	case kwREFRESH:
		refreshTok := p.advance()
		stmt.Kind = AlterViewRefresh
		stmt.Loc.End = refreshTok.Loc.End
	default:
		return nil, p.syntaxErrorAtCur()
	}
	return stmt, nil
}

// ---------------------------------------------------------------------------
// CREATE MATERIALIZED VIEW
// ---------------------------------------------------------------------------

// MaterializedViewStaleMode is the WHEN STALE mode (D-MV1 docs-ahead extension).
type MaterializedViewStaleMode int

const (
	// MVStaleNone means no WHEN STALE clause was given.
	MVStaleNone MaterializedViewStaleMode = iota
	// MVStaleInline is WHEN STALE INLINE.
	MVStaleInline
	// MVStaleFail is WHEN STALE FAIL.
	MVStaleFail
)

// CreateMaterializedViewStmt is `CREATE [OR REPLACE] MATERIALIZED VIEW
// [IF NOT EXISTS] name [GRACE PERIOD interval] [WHEN STALE (INLINE|FAIL)]
// [COMMENT str] [WITH (props)] AS query`.
type CreateMaterializedViewStmt struct {
	OrReplace   bool
	IfNotExists bool
	Name        *ast.QualifiedName
	GracePeriod Expr // INTERVAL literal; nil when absent
	StaleMode   MaterializedViewStaleMode
	Comment     *string
	Properties  []*Property
	Query       *Query
	Loc         ast.Loc
}

func (n *CreateMaterializedViewStmt) Tag() ast.NodeTag { return ast.T_Invalid }
func (n *CreateMaterializedViewStmt) Span() ast.Loc    { return n.Loc }

// parseCreateMaterializedViewStmt parses CREATE [OR REPLACE] MATERIALIZED VIEW.
// On entry cur is CREATE. The optional clauses (GRACE PERIOD, WHEN STALE,
// COMMENT, WITH) are parsed in an order-tolerant loop, each at most once, then
// the mandatory `AS query` (D-MV1).
func (p *Parser) parseCreateMaterializedViewStmt(startOffset int) (ast.Node, error) {
	p.advance() // consume CREATE
	orReplace, err := p.parseOptionalOrReplace()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(kwMATERIALIZED); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwVIEW); err != nil {
		return nil, err
	}
	ifNotExists, err := p.parseOptionalIfNotExists()
	if err != nil {
		return nil, err
	}
	name, err := p.parseBoundedQualifiedName(3, "materialized view name")
	if err != nil {
		return nil, err
	}
	stmt := &CreateMaterializedViewStmt{
		OrReplace:   orReplace,
		IfNotExists: ifNotExists,
		Name:        name,
		Loc:         ast.Loc{Start: startOffset},
	}

	// D-MV1: the optional clauses appear in a FIXED order — GRACE PERIOD, then
	// WHEN STALE, then COMMENT, then WITH — each at most once. Trino 481 REJECTS
	// out-of-order clauses (`COMMENT 'c' GRACE PERIOD …`, `WITH (…) COMMENT …`,
	// `WHEN STALE … GRACE PERIOD …` are SYNTAX_ERRORs), so they are consumed as a
	// strict sequence, not an order-tolerant set — verified against the oracle.
	if p.cur.Kind == kwGRACE {
		p.advance() // consume GRACE
		if _, err := p.expect(kwPERIOD); err != nil {
			return nil, err
		}
		if p.cur.Kind != kwINTERVAL {
			return nil, p.exprErrorAt("expected INTERVAL after GRACE PERIOD")
		}
		gp, err := p.parseIntervalLiteral()
		if err != nil {
			return nil, err
		}
		stmt.GracePeriod = gp
	}
	if p.cur.Kind == kwWHEN {
		// WHEN STALE (INLINE | FAIL). STALE/INLINE/FAIL are not keyword tokens;
		// match them by identifier text.
		p.advance() // consume WHEN
		if p.cur.Kind != tokIdent || !strings.EqualFold(p.cur.Str, "STALE") {
			return nil, &ParseError{Loc: p.cur.Loc, Msg: "expected STALE after WHEN"}
		}
		p.advance() // consume STALE
		switch {
		case p.cur.Kind == tokIdent && strings.EqualFold(p.cur.Str, "INLINE"):
			stmt.StaleMode = MVStaleInline
		case p.cur.Kind == tokIdent && strings.EqualFold(p.cur.Str, "FAIL"):
			stmt.StaleMode = MVStaleFail
		default:
			return nil, &ParseError{Loc: p.cur.Loc, Msg: "expected INLINE or FAIL after WHEN STALE"}
		}
		p.advance() // consume INLINE / FAIL
	}
	comment, err := p.parseOptionalComment()
	if err != nil {
		return nil, err
	}
	stmt.Comment = comment
	props, err := p.parseOptionalWithProperties()
	if err != nil {
		return nil, err
	}
	stmt.Properties = props

	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}
	query, err := p.parseQuery()
	if err != nil {
		return nil, err
	}
	stmt.Query = query
	stmt.Loc.End = query.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// ALTER MATERIALIZED VIEW
// ---------------------------------------------------------------------------

// AlterMaterializedViewKind classifies an ALTER MATERIALIZED VIEW sub-form.
type AlterMaterializedViewKind int

const (
	// AlterMVRename is ALTER MATERIALIZED VIEW [IF EXISTS] name RENAME TO new.
	AlterMVRename AlterMaterializedViewKind = iota
	// AlterMVSetProperties is ALTER MATERIALIZED VIEW name SET PROPERTIES …
	AlterMVSetProperties
	// AlterMVSetAuthorization is ALTER MATERIALIZED VIEW name SET AUTHORIZATION
	// principal (D-MV3 docs-ahead extension — absent from the legacy grammar but
	// documented and oracle-accepted in 481).
	AlterMVSetAuthorization
)

// AlterMaterializedViewStmt is an ALTER MATERIALIZED VIEW statement.
type AlterMaterializedViewStmt struct {
	Kind          AlterMaterializedViewKind
	IfExists      bool // only valid on the RENAME form (D-MV3)
	Name          *ast.QualifiedName
	NewName       *ast.QualifiedName // RENAME target
	Properties    []*Property        // SET PROPERTIES form
	Authorization *Principal         // SET AUTHORIZATION form
	Loc           ast.Loc
}

func (n *AlterMaterializedViewStmt) Tag() ast.NodeTag { return ast.T_Invalid }
func (n *AlterMaterializedViewStmt) Span() ast.Loc    { return n.Loc }

// parseAlterMaterializedViewStmt parses ALTER MATERIALIZED VIEW. On entry cur
// is ALTER (dispatch peeked MATERIALIZED).
func (p *Parser) parseAlterMaterializedViewStmt(startOffset int) (ast.Node, error) {
	p.advance() // consume ALTER
	p.advance() // consume MATERIALIZED
	if _, err := p.expect(kwVIEW); err != nil {
		return nil, err
	}
	// D-MV3: IF EXISTS is valid ONLY on the RENAME form (legacy grammar carries
	// `(IF_ EXISTS_)?` only on renameMaterializedView, not on the SET forms). It
	// appears before the name, so it is parsed here and the SET branch rejects
	// it. Trino 481 rejects `ALTER MATERIALIZED VIEW IF EXISTS … SET …` as a
	// SYNTAX_ERROR.
	ifExists, err := p.parseOptionalIfExists()
	if err != nil {
		return nil, err
	}
	name, err := p.parseBoundedQualifiedName(3, "materialized view name")
	if err != nil {
		return nil, err
	}
	stmt := &AlterMaterializedViewStmt{IfExists: ifExists, Name: name, Loc: ast.Loc{Start: startOffset}}

	switch p.cur.Kind {
	case kwRENAME:
		p.advance()
		if _, err := p.expect(kwTO); err != nil {
			return nil, err
		}
		// The RENAME TO target is UNBOUNDED (oracle-confirmed, ledger DDL4): a
		// 4+-part target fails semantically (TABLE_NOT_FOUND), not at parse.
		newName, err := p.parseQualifiedName()
		if err != nil {
			return nil, err
		}
		stmt.Kind = AlterMVRename
		stmt.NewName = newName
		stmt.Loc.End = newName.Loc.End
	case kwSET:
		// D-MV3: IF EXISTS is not allowed on the SET forms.
		if ifExists {
			return nil, &ParseError{Loc: name.Loc, Msg: "IF EXISTS is not allowed on ALTER MATERIALIZED VIEW … SET"}
		}
		p.advance() // consume SET
		switch p.cur.Kind {
		case kwPROPERTIES:
			p.advance() // consume PROPERTIES
			props, err := p.parsePropertyAssignments()
			if err != nil {
				return nil, err
			}
			stmt.Kind = AlterMVSetProperties
			stmt.Properties = props
			stmt.Loc.End = p.prev.Loc.End
		case kwAUTHORIZATION:
			// D-MV3: SET AUTHORIZATION (docs-ahead; absent from legacy grammar,
			// accepted by 481).
			p.advance() // consume AUTHORIZATION
			princ, err := p.parseAuthorizationPrincipal()
			if err != nil {
				return nil, err
			}
			stmt.Kind = AlterMVSetAuthorization
			stmt.Authorization = princ
			stmt.Loc.End = princ.Loc.End
		default:
			return nil, p.syntaxErrorAtCur()
		}
	default:
		return nil, p.syntaxErrorAtCur()
	}
	return stmt, nil
}

// ---------------------------------------------------------------------------
// REFRESH MATERIALIZED VIEW
// ---------------------------------------------------------------------------

// RefreshMaterializedViewStmt is `REFRESH MATERIALIZED VIEW name`.
type RefreshMaterializedViewStmt struct {
	Name *ast.QualifiedName
	Loc  ast.Loc
}

func (n *RefreshMaterializedViewStmt) Tag() ast.NodeTag { return ast.T_Invalid }
func (n *RefreshMaterializedViewStmt) Span() ast.Loc    { return n.Loc }

// parseRefreshMaterializedViewStmt parses REFRESH MATERIALIZED VIEW name. On
// entry cur is REFRESH.
func (p *Parser) parseRefreshMaterializedViewStmt(startOffset int) (ast.Node, error) {
	p.advance() // consume REFRESH
	if _, err := p.expect(kwMATERIALIZED); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwVIEW); err != nil {
		return nil, err
	}
	name, err := p.parseBoundedQualifiedName(3, "materialized view name")
	if err != nil {
		return nil, err
	}
	return &RefreshMaterializedViewStmt{Name: name, Loc: ast.Loc{Start: startOffset, End: name.Loc.End}}, nil
}

// ---------------------------------------------------------------------------
// Shared helpers
// ---------------------------------------------------------------------------

// parseOptionalOrReplace consumes an `OR REPLACE` clause when present.
func (p *Parser) parseOptionalOrReplace() (bool, error) {
	if _, ok := p.match(kwOR); !ok {
		return false, nil
	}
	if _, err := p.expect(kwREPLACE); err != nil {
		return false, err
	}
	return true, nil
}

// parsePropertyAssignments parses a bare `property (, property)*` list (no
// surrounding parentheses) — the SET PROPERTIES form's propertyAssignments rule.
func (p *Parser) parsePropertyAssignments() ([]*Property, error) {
	first, err := p.parseProperty()
	if err != nil {
		return nil, err
	}
	props := []*Property{first}
	for {
		if _, ok := p.match(int(',')); !ok {
			break
		}
		next, err := p.parseProperty()
		if err != nil {
			return nil, err
		}
		props = append(props, next)
	}
	return props, nil
}

var (
	_ ast.Node = (*CreateSchemaStmt)(nil)
	_ ast.Node = (*AlterSchemaStmt)(nil)
	_ ast.Node = (*CreateViewStmt)(nil)
	_ ast.Node = (*AlterViewStmt)(nil)
	_ ast.Node = (*CreateMaterializedViewStmt)(nil)
	_ ast.Node = (*AlterMaterializedViewStmt)(nil)
	_ ast.Node = (*RefreshMaterializedViewStmt)(nil)
)
