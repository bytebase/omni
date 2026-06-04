package parser

import (
	"github.com/bytebase/omni/trino/ast"
)

// This file is part of the `parser-utility` DAG node (with show.go and
// explain_call.go): it implements Trino's session / catalog-context statements —
// the legacy TrinoParser.g4 `statement` alternatives labelled use, setSession,
// resetSession, setSessionAuthorization, resetSessionAuthorization, setRole,
// setPath, and setTimeZone.
//
// As in show.go, the statement nodes are PARSER-PACKAGE types satisfying
// ast.Node via tags declared in trino/ast/nodetags.go.
//
// The legacy alternatives implemented here:
//
//	USE schema                                # use
//	USE catalog.schema                        # use
//	SET SESSION qualifiedName EQ expression   # setSession
//	RESET SESSION qualifiedName               # resetSession
//	SET SESSION AUTHORIZATION authorizationUser   # setSessionAuthorization
//	RESET SESSION AUTHORIZATION               # resetSessionAuthorization
//	SET ROLE (ALL | NONE | identifier) (IN catalog)?  # setRole
//	SET PATH pathSpecification                # setPath
//	SET TIME ZONE (LOCAL | expression)        # setTimeZone
//
//	authorizationUser : identifier | string_ ;
//	pathSpecification : pathElement (COMMA pathElement)* ;
//	pathElement       : identifier DOT identifier | identifier ;
//
// Adjudicated against the live Trino 481 oracle. The `SET` leading keyword
// fans out on the second keyword (SESSION / ROLE / PATH / TIME), and `SET
// SESSION` further splits on AUTHORIZATION vs a property name; `RESET SESSION`
// splits on AUTHORIZATION vs a property name. All of these forms are this
// node's: SET ROLE is the session-role activation (not a DCL grant — that is
// the parser-dcl-tcl node's grant_revoke.go), and SET PATH / SET TIME ZONE have
// no other home.

// ---------------------------------------------------------------------------
// USE
// ---------------------------------------------------------------------------

// UseStmt is USE [catalog.]schema. Catalog is nil for the schema-only form.
type UseStmt struct {
	// Catalog is the catalog identifier in `USE catalog.schema`; nil for the
	// `USE schema` form.
	Catalog *ast.Identifier
	// Schema is the schema identifier (always present).
	Schema *ast.Identifier
	Loc    ast.Loc
}

// Tag implements ast.Node.
func (s *UseStmt) Tag() ast.NodeTag { return ast.T_UseStmt }

// Span returns the source byte range.
func (s *UseStmt) Span() ast.Loc { return s.Loc }

var _ ast.Node = (*UseStmt)(nil)

// parseUseStmt parses USE schema | USE catalog.schema. USE is the current token.
func (p *Parser) parseUseStmt() (ast.Node, error) {
	useTok := p.advance() // consume USE
	first, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	s := &UseStmt{Schema: first, Loc: ast.Loc{Start: useTok.Loc.Start, End: first.Loc.End}}
	if _, ok := p.match(int('.')); ok {
		schema, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		s.Catalog = first
		s.Schema = schema
		s.Loc.End = schema.Loc.End
	}
	return s, nil
}

// ---------------------------------------------------------------------------
// SET dispatch
// ---------------------------------------------------------------------------

// parseSetStmt parses a statement whose leading keyword is SET, dispatching on
// the second keyword: SESSION (setSession / setSessionAuthorization), ROLE
// (setRole), PATH (setPath), TIME (setTimeZone). SET is the current token.
func (p *Parser) parseSetStmt() (ast.Node, error) {
	setTok := p.advance() // consume SET
	start := setTok.Loc.Start
	switch p.cur.Kind {
	case kwSESSION:
		return p.parseSetSession(start)
	case kwROLE:
		return p.parseSetRole(start)
	case kwPATH:
		return p.parseSetPath(start)
	case kwTIME:
		return p.parseSetTimeZone(start)
	default:
		return nil, p.syntaxErrorAtCur()
	}
}

// ---------------------------------------------------------------------------
// SET / RESET SESSION
// ---------------------------------------------------------------------------

// SetSessionStmt is SET SESSION name = expression. Name is a qualifiedName
// (system property `name` or catalog property `catalog.name`).
type SetSessionStmt struct {
	Name  *ast.QualifiedName
	Value Expr
	Loc   ast.Loc
}

// Tag implements ast.Node.
func (s *SetSessionStmt) Tag() ast.NodeTag { return ast.T_SetSessionStmt }

// Span returns the source byte range.
func (s *SetSessionStmt) Span() ast.Loc { return s.Loc }

var _ ast.Node = (*SetSessionStmt)(nil)

// ResetSessionStmt is RESET SESSION name.
type ResetSessionStmt struct {
	Name *ast.QualifiedName
	Loc  ast.Loc
}

// Tag implements ast.Node.
func (s *ResetSessionStmt) Tag() ast.NodeTag { return ast.T_ResetSessionStmt }

// Span returns the source byte range.
func (s *ResetSessionStmt) Span() ast.Loc { return s.Loc }

var _ ast.Node = (*ResetSessionStmt)(nil)

// SetSessionAuthorizationStmt is SET SESSION AUTHORIZATION user, where user is
// an identifier or a string literal (authorizationUser).
type SetSessionAuthorizationStmt struct {
	// User is the identifier form (`SET SESSION AUTHORIZATION alice`); nil when
	// the string form was used.
	User *ast.Identifier
	// UserString is the string-literal form (`SET SESSION AUTHORIZATION 'alice'`)
	// decoded value; HasUserString marks its presence.
	UserString    string
	HasUserString bool
	Loc           ast.Loc
}

// Tag implements ast.Node.
func (s *SetSessionAuthorizationStmt) Tag() ast.NodeTag { return ast.T_SetSessionAuthorizationStmt }

// Span returns the source byte range.
func (s *SetSessionAuthorizationStmt) Span() ast.Loc { return s.Loc }

var _ ast.Node = (*SetSessionAuthorizationStmt)(nil)

// ResetSessionAuthorizationStmt is RESET SESSION AUTHORIZATION.
type ResetSessionAuthorizationStmt struct {
	Loc ast.Loc
}

// Tag implements ast.Node.
func (s *ResetSessionAuthorizationStmt) Tag() ast.NodeTag {
	return ast.T_ResetSessionAuthorizationStmt
}

// Span returns the source byte range.
func (s *ResetSessionAuthorizationStmt) Span() ast.Loc { return s.Loc }

var _ ast.Node = (*ResetSessionAuthorizationStmt)(nil)

// parseSetSession parses SET SESSION { AUTHORIZATION user | name = expression }.
// SESSION is the current token; start is the SET keyword's start offset.
//
// AUTHORIZATION is a non-reserved keyword, so it can ALSO be the first part of a
// session-property name. Trino 481 accepts `SET SESSION AUTHORIZATION = 1` and
// `SET SESSION AUTHORIZATION.foo = 1` as property assignments (the property is
// named `authorization` / `authorization.foo`). So the AUTHORIZATION-user form
// fires only when AUTHORIZATION is NOT immediately followed by '.' or '=' — i.e.
// when it is genuinely the SET SESSION AUTHORIZATION keyword rather than a
// property name. Otherwise we fall through to the qualifiedName property path.
func (p *Parser) parseSetSession(start int) (ast.Node, error) {
	p.advance() // consume SESSION
	if p.cur.Kind == kwAUTHORIZATION && p.peekNext().Kind != int('.') && p.peekNext().Kind != int('=') {
		p.advance() // consume AUTHORIZATION
		return p.parseSetSessionAuthorization(start)
	}
	name, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(int('=')); err != nil {
		return nil, err
	}
	value, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	return &SetSessionStmt{
		Name:  name,
		Value: value,
		Loc:   ast.Loc{Start: start, End: value.Span().End},
	}, nil
}

// parseSetSessionAuthorization parses the tail of SET SESSION AUTHORIZATION:
// an identifier or a string literal. AUTHORIZATION has been consumed.
func (p *Parser) parseSetSessionAuthorization(start int) (ast.Node, error) {
	if p.cur.Kind == tokString || p.cur.Kind == tokUnicodeString {
		strTok := p.advance()
		return &SetSessionAuthorizationStmt{
			UserString:    strTok.Str,
			HasUserString: true,
			Loc:           ast.Loc{Start: start, End: strTok.Loc.End},
		}, nil
	}
	ident, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	return &SetSessionAuthorizationStmt{
		User: ident,
		Loc:  ast.Loc{Start: start, End: ident.Loc.End},
	}, nil
}

// parseResetStmt parses RESET SESSION { AUTHORIZATION | name }. RESET is the
// current token.
func (p *Parser) parseResetStmt() (ast.Node, error) {
	resetTok := p.advance() // consume RESET
	start := resetTok.Loc.Start
	if _, err := p.expect(kwSESSION); err != nil {
		return nil, err
	}
	// As in parseSetSession, AUTHORIZATION is the RESET SESSION AUTHORIZATION
	// keyword only when not followed by '.' (a dotted property name like
	// `authorization.foo`); Trino 481 accepts `RESET SESSION AUTHORIZATION.foo`
	// as a property reset. (RESET takes no '=', so only '.' disambiguates.)
	if p.cur.Kind == kwAUTHORIZATION && p.peekNext().Kind != int('.') {
		authTok := p.advance() // consume AUTHORIZATION
		return &ResetSessionAuthorizationStmt{
			Loc: ast.Loc{Start: start, End: authTok.Loc.End},
		}, nil
	}
	name, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	return &ResetSessionStmt{
		Name: name,
		Loc:  ast.Loc{Start: start, End: name.Loc.End},
	}, nil
}

// ---------------------------------------------------------------------------
// SET ROLE
// ---------------------------------------------------------------------------

// RoleSpec classifies the role argument of SET ROLE.
type RoleSpec int

const (
	// RoleNamed is SET ROLE <identifier>.
	RoleNamed RoleSpec = iota
	// RoleAll is SET ROLE ALL.
	RoleAll
	// RoleNone is SET ROLE NONE.
	RoleNone
)

// SetRoleStmt is SET ROLE (role | ALL | NONE) [IN catalog].
type SetRoleStmt struct {
	Spec RoleSpec
	// Role is the role identifier when Spec == RoleNamed; nil otherwise.
	Role *ast.Identifier
	// Catalog is the `IN catalog` scope identifier; nil when absent.
	Catalog *ast.Identifier
	Loc     ast.Loc
}

// Tag implements ast.Node.
func (s *SetRoleStmt) Tag() ast.NodeTag { return ast.T_SetRoleStmt }

// Span returns the source byte range.
func (s *SetRoleStmt) Span() ast.Loc { return s.Loc }

var _ ast.Node = (*SetRoleStmt)(nil)

// parseSetRole parses the tail of SET ROLE: (ALL | NONE | identifier)
// (IN catalog)?. ROLE is the current token; start is the SET start offset.
func (p *Parser) parseSetRole(start int) (ast.Node, error) {
	roleTok := p.advance() // consume ROLE
	s := &SetRoleStmt{Loc: ast.Loc{Start: start, End: roleTok.Loc.End}}
	switch p.cur.Kind {
	case kwALL:
		end := p.advance()
		s.Spec = RoleAll
		s.Loc.End = end.Loc.End
	case kwNONE:
		end := p.advance()
		s.Spec = RoleNone
		s.Loc.End = end.Loc.End
	default:
		ident, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		s.Spec = RoleNamed
		s.Role = ident
		s.Loc.End = ident.Loc.End
	}
	if _, ok := p.match(kwIN); ok {
		catalog, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		s.Catalog = catalog
		s.Loc.End = catalog.Loc.End
	}
	return s, nil
}

// ---------------------------------------------------------------------------
// SET PATH
// ---------------------------------------------------------------------------

// PathElement is one element of a SET PATH specification: a bare schema
// (`schema`) or a qualified `catalog.schema`. Catalog is nil for the bare form.
type PathElement struct {
	Catalog *ast.Identifier
	Schema  *ast.Identifier
	Loc     ast.Loc
}

// SetPathStmt is SET PATH pathElement (, pathElement)*.
type SetPathStmt struct {
	Elements []PathElement
	Loc      ast.Loc
}

// Tag implements ast.Node.
func (s *SetPathStmt) Tag() ast.NodeTag { return ast.T_SetPathStmt }

// Span returns the source byte range.
func (s *SetPathStmt) Span() ast.Loc { return s.Loc }

var _ ast.Node = (*SetPathStmt)(nil)

// parsePathElement parses `identifier (DOT identifier)?` (pathElement).
func (p *Parser) parsePathElement() (PathElement, error) {
	first, err := p.parseIdentifier()
	if err != nil {
		return PathElement{}, err
	}
	el := PathElement{Schema: first, Loc: first.Loc}
	if _, ok := p.match(int('.')); ok {
		schema, err := p.parseIdentifier()
		if err != nil {
			return PathElement{}, err
		}
		el.Catalog = first
		el.Schema = schema
		el.Loc.End = schema.Loc.End
	}
	return el, nil
}

// parseSetPath parses the tail of SET PATH: pathElement (, pathElement)*. PATH
// is the current token; start is the SET start offset.
func (p *Parser) parseSetPath(start int) (ast.Node, error) {
	p.advance() // consume PATH
	first, err := p.parsePathElement()
	if err != nil {
		return nil, err
	}
	s := &SetPathStmt{Elements: []PathElement{first}, Loc: ast.Loc{Start: start, End: first.Loc.End}}
	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		next, err := p.parsePathElement()
		if err != nil {
			return nil, err
		}
		s.Elements = append(s.Elements, next)
		s.Loc.End = next.Loc.End
	}
	return s, nil
}

// ---------------------------------------------------------------------------
// SET TIME ZONE
// ---------------------------------------------------------------------------

// SetTimeZoneStmt is SET TIME ZONE (LOCAL | expression). Local is true for the
// LOCAL form; Value holds the expression otherwise.
type SetTimeZoneStmt struct {
	Local bool
	// Value is the time-zone expression (string / interval / function call);
	// nil when Local is true.
	Value Expr
	Loc   ast.Loc
}

// Tag implements ast.Node.
func (s *SetTimeZoneStmt) Tag() ast.NodeTag { return ast.T_SetTimeZoneStmt }

// Span returns the source byte range.
func (s *SetTimeZoneStmt) Span() ast.Loc { return s.Loc }

var _ ast.Node = (*SetTimeZoneStmt)(nil)

// parseSetTimeZone parses the tail of SET TIME ZONE: LOCAL | expression. TIME is
// the current token; start is the SET start offset.
func (p *Parser) parseSetTimeZone(start int) (ast.Node, error) {
	p.advance() // consume TIME
	if _, err := p.expect(kwZONE); err != nil {
		return nil, err
	}
	if localTok, ok := p.match(kwLOCAL); ok {
		return &SetTimeZoneStmt{
			Local: true,
			Loc:   ast.Loc{Start: start, End: localTok.Loc.End},
		}, nil
	}
	value, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	return &SetTimeZoneStmt{
		Value: value,
		Loc:   ast.Loc{Start: start, End: value.Span().End},
	}, nil
}
