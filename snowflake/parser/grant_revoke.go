package parser

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// DCL — GRANT / REVOKE (roles, privileges, ownership, shares) [T6.1]
// ---------------------------------------------------------------------------
//
// Snowflake's GRANT/REVOKE grammar is large and continuously extended with new
// privileges and object types. Rather than mirror the legacy ANTLR grammar's
// finite privilege/object enumerations — which the official documentation
// corpus already exceeds (e.g. ON NOTEBOOK, ON WORKSPACE, CREATE PROVISIONED
// THROUGHPUT, CREATE SNOWFLAKE.CORE.BUDGET) — these parse functions treat
// privileges and object types as free-form token runs. Structural keywords
// (ON, TO, FROM, IN, ALL, FUTURE, WITH GRANT OPTION, CURRENT GRANTS,
// CASCADE/RESTRICT, GRANT OPTION FOR) anchor the grammar; the spans between
// them are captured verbatim. Validating that a privilege is legal for an
// object is the catalog/semantic layer's job, not the parser's.

// ---------------------------------------------------------------------------
// GRANT dispatch
// ---------------------------------------------------------------------------

// parseGrantStmt parses a GRANT statement in any of its shapes:
//
//	GRANT { <privileges> | ALL [PRIVILEGES] } ON <target> TO <grantee> [WITH GRANT OPTION]
//	GRANT [DATABASE] ROLE <name> TO <grantee>
//	GRANT OWNERSHIP ON <target> TO <grantee> [{REVOKE|COPY} CURRENT GRANTS]
//
// The GRANT keyword has NOT yet been consumed when this function is called.
func (p *Parser) parseGrantStmt() (ast.Node, error) {
	grantTok := p.advance() // consume GRANT
	start := grantTok.Loc

	switch {
	case p.cur.Type == kwOWNERSHIP:
		return p.parseGrantOwnership(start)
	case p.cur.Type == kwROLE:
		return p.parseGrantRole(start, ast.GrantedAccountRole)
	case p.cur.Type == kwDATABASE && p.peekNext().Type == kwROLE:
		return p.parseGrantRole(start, ast.GrantedDatabaseRole)
	case p.cur.Type == kwAPPLICATION && p.peekNext().Type == kwROLE:
		return p.parseGrantRole(start, ast.GrantedApplicationRole)
	default:
		return p.parseGrantPrivileges(start)
	}
}

// parseGrantRole parses GRANT [DATABASE | APPLICATION] ROLE <name> TO <grantee>.
// start is the Loc of the GRANT token. roleKind selects the granted-role flavor;
// on entry cur is ROLE (account) or the DATABASE/APPLICATION qualifier.
func (p *Parser) parseGrantRole(start ast.Loc, roleKind ast.GrantedRoleKind) (ast.Node, error) {
	if roleKind != ast.GrantedAccountRole {
		p.advance() // consume DATABASE / APPLICATION qualifier
	}
	if _, err := p.expect(kwROLE); err != nil {
		return nil, err
	}

	stmt := &ast.GrantStmt{
		Kind:     ast.GrantRole,
		RoleKind: roleKind,
		Loc:      ast.Loc{Start: start.Start},
	}

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Role = name

	if _, err := p.expect(kwTO); err != nil {
		return nil, err
	}
	grantee, err := p.parseGrantee()
	if err != nil {
		return nil, err
	}
	stmt.Grantee = grantee

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseGrantOwnership parses GRANT OWNERSHIP ON <target> TO <grantee>
// [{REVOKE|COPY} CURRENT GRANTS]. On entry cur is OWNERSHIP.
func (p *Parser) parseGrantOwnership(start ast.Loc) (ast.Node, error) {
	p.advance() // consume OWNERSHIP

	stmt := &ast.GrantStmt{
		Kind: ast.GrantOwnership,
		Loc:  ast.Loc{Start: start.Start},
	}

	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	target, err := p.parseGrantTarget()
	if err != nil {
		return nil, err
	}
	stmt.On = target

	if _, err := p.expect(kwTO); err != nil {
		return nil, err
	}
	grantee, err := p.parseGrantee()
	if err != nil {
		return nil, err
	}
	stmt.Grantee = grantee

	// Optional { REVOKE | COPY } CURRENT GRANTS.
	switch p.cur.Type {
	case kwREVOKE:
		if p.peekNext().Type == kwCURRENT {
			p.advance() // consume REVOKE
			if err := p.expectCurrentGrants(); err != nil {
				return nil, err
			}
			stmt.CurrentGrants = ast.CurrentGrantsRevoke
		}
	case kwCOPY:
		if p.peekNext().Type == kwCURRENT {
			p.advance() // consume COPY
			if err := p.expectCurrentGrants(); err != nil {
				return nil, err
			}
			stmt.CurrentGrants = ast.CurrentGrantsCopy
		}
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseGrantPrivileges parses
// GRANT { <privileges> | ALL [PRIVILEGES] } ON <target> TO <grantee>
// [WITH GRANT OPTION]. On entry cur is the first privilege token (or ALL).
func (p *Parser) parseGrantPrivileges(start ast.Loc) (ast.Node, error) {
	stmt := &ast.GrantStmt{
		Kind: ast.GrantPrivileges,
		Loc:  ast.Loc{Start: start.Start},
	}

	all, privs, err := p.parsePrivilegeList()
	if err != nil {
		return nil, err
	}
	stmt.AllPrivileges = all
	stmt.Privileges = privs

	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	target, err := p.parseGrantTarget()
	if err != nil {
		return nil, err
	}
	stmt.On = target

	if _, err := p.expect(kwTO); err != nil {
		return nil, err
	}
	grantee, err := p.parseGrantee()
	if err != nil {
		return nil, err
	}
	stmt.Grantee = grantee

	// Optional WITH GRANT OPTION.
	if p.cur.Type == kwWITH && p.peekNext().Type == kwGRANT {
		p.advance() // consume WITH
		p.advance() // consume GRANT
		if _, err := p.expect(kwOPTION); err != nil {
			return nil, err
		}
		stmt.GrantOption = true
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// REVOKE dispatch
// ---------------------------------------------------------------------------

// parseRevokeStmt parses a REVOKE statement in either of its shapes:
//
//	REVOKE [GRANT OPTION FOR] { <privileges> | ALL [PRIVILEGES] } ON <target> FROM <grantee> [CASCADE|RESTRICT]
//	REVOKE [DATABASE] ROLE <name> FROM <grantee>
//
// The REVOKE keyword has NOT yet been consumed when this function is called.
func (p *Parser) parseRevokeStmt() (ast.Node, error) {
	revokeTok := p.advance() // consume REVOKE
	start := revokeTok.Loc

	// REVOKE [DATABASE | APPLICATION] ROLE ... — GRANT OPTION FOR applies only
	// to privilege revocation, so a leading ROLE/DATABASE ROLE/APPLICATION ROLE
	// is unambiguously a role revoke.
	if p.cur.Type == kwROLE {
		return p.parseRevokeRole(start, ast.GrantedAccountRole)
	}
	if p.cur.Type == kwDATABASE && p.peekNext().Type == kwROLE {
		return p.parseRevokeRole(start, ast.GrantedDatabaseRole)
	}
	if p.cur.Type == kwAPPLICATION && p.peekNext().Type == kwROLE {
		return p.parseRevokeRole(start, ast.GrantedApplicationRole)
	}
	return p.parseRevokePrivileges(start)
}

// parseRevokeRole parses REVOKE [DATABASE | APPLICATION] ROLE <name> FROM <grantee>.
// On entry cur is ROLE (account) or the DATABASE/APPLICATION qualifier.
func (p *Parser) parseRevokeRole(start ast.Loc, roleKind ast.GrantedRoleKind) (ast.Node, error) {
	if roleKind != ast.GrantedAccountRole {
		p.advance() // consume DATABASE / APPLICATION qualifier
	}
	if _, err := p.expect(kwROLE); err != nil {
		return nil, err
	}

	stmt := &ast.RevokeStmt{
		Kind:     ast.RevokeRole,
		RoleKind: roleKind,
		Loc:      ast.Loc{Start: start.Start},
	}

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Role = name

	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}
	grantee, err := p.parseGrantee()
	if err != nil {
		return nil, err
	}
	stmt.Grantee = grantee

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseRevokePrivileges parses
// REVOKE [GRANT OPTION FOR] { <privileges> | ALL [PRIVILEGES] } ON <target>
// FROM <grantee> [CASCADE|RESTRICT]. On entry cur is GRANT (for GRANT OPTION
// FOR) or the first privilege token (or ALL).
func (p *Parser) parseRevokePrivileges(start ast.Loc) (ast.Node, error) {
	stmt := &ast.RevokeStmt{
		Kind: ast.RevokePrivileges,
		Loc:  ast.Loc{Start: start.Start},
	}

	// Optional GRANT OPTION FOR.
	if p.cur.Type == kwGRANT && p.peekNext().Type == kwOPTION {
		p.advance() // consume GRANT
		p.advance() // consume OPTION
		if _, err := p.expect(kwFOR); err != nil {
			return nil, err
		}
		stmt.GrantOptionFor = true
	}

	all, privs, err := p.parsePrivilegeList()
	if err != nil {
		return nil, err
	}
	stmt.AllPrivileges = all
	stmt.Privileges = privs

	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	target, err := p.parseGrantTarget()
	if err != nil {
		return nil, err
	}
	stmt.On = target

	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}
	grantee, err := p.parseGrantee()
	if err != nil {
		return nil, err
	}
	stmt.Grantee = grantee

	// Optional CASCADE | RESTRICT.
	switch p.cur.Type {
	case kwCASCADE:
		p.advance()
		stmt.Cascade = true
	case kwRESTRICT:
		p.advance()
		stmt.Restrict = true
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// Shared sub-grammar: privilege list, ON target, grantee
// ---------------------------------------------------------------------------

// parsePrivilegeList parses { <privileges> | ALL [PRIVILEGES] }, the span
// immediately before ON. Returns (allPrivileges, privileges, err); exactly one
// of the two result paths applies.
//
//	ALL [PRIVILEGES]          -> (true, nil, nil)
//	priv [, priv ...]         -> (false, [...], nil)
//
// A single privilege is one or more consecutive name tokens (keywords and/or
// identifiers, including dotted class names like SNOWFLAKE.CORE.BUDGET),
// captured verbatim and uppercased. Privileges are separated by commas and the
// whole list is terminated by ON.
func (p *Parser) parsePrivilegeList() (bool, []*ast.Privilege, error) {
	// ALL [PRIVILEGES] — but only treat a leading ALL as "all privileges" when
	// it is immediately followed by PRIVILEGES or ON (i.e. ALL is the entire
	// privilege spec). A privilege literally beginning with ALL is not a
	// documented Snowflake form, so this is unambiguous.
	if p.cur.Type == kwALL {
		next := p.peekNext()
		if next.Type == kwPRIVILEGES || next.Type == kwON {
			p.advance() // consume ALL
			if p.cur.Type == kwPRIVILEGES {
				p.advance() // consume PRIVILEGES
			}
			return true, nil, nil
		}
	}

	var privs []*ast.Privilege
	for {
		priv, err := p.parsePrivilege()
		if err != nil {
			return false, nil, err
		}
		privs = append(privs, priv)

		if p.cur.Type == ',' {
			p.advance() // consume ','
			continue
		}
		break
	}
	return false, privs, nil
}

// parsePrivilege parses one privilege: a run of one or more name tokens
// (keywords or identifiers), joining multi-word privileges (e.g.
// CREATE MATERIALIZED VIEW) with single spaces and dotted class names (e.g.
// SNOWFLAKE.CORE.BUDGET) with dots. The run stops at ON, a comma, or any
// non-name token. Requires at least one name token.
func (p *Parser) parsePrivilege() (*ast.Privilege, error) {
	if !p.isPrivilegeWord(p.cur.Type) {
		return nil, p.syntaxErrorAtCur()
	}

	startLoc := p.cur.Loc
	var b strings.Builder
	endLoc := p.cur.Loc

	first := p.advance()
	b.WriteString(strings.ToUpper(first.Str))
	endLoc = first.Loc

	for {
		// Dotted continuation: name . name (class privilege).
		if p.cur.Type == '.' && p.isPrivilegeWord(p.peekNext().Type) {
			p.advance() // consume '.'
			part := p.advance()
			b.WriteByte('.')
			b.WriteString(strings.ToUpper(part.Str))
			endLoc = part.Loc
			continue
		}
		// Space-separated continuation: stop at ON / comma / structural tokens.
		if p.cur.Type == kwON || p.cur.Type == ',' {
			break
		}
		if !p.isPrivilegeWord(p.cur.Type) {
			break
		}
		part := p.advance()
		b.WriteByte(' ')
		b.WriteString(strings.ToUpper(part.Str))
		endLoc = part.Loc
	}

	return &ast.Privilege{
		Name: b.String(),
		Loc:  ast.Loc{Start: startLoc.Start, End: endLoc.End},
	}, nil
}

// isPrivilegeWord reports whether a token type may appear inside a privilege
// name. Identifiers and any keyword qualify; structural punctuation and EOF do
// not. ON is excluded because it terminates the privilege list. The lexer emits
// most privilege words (SELECT, USAGE, CREATE, ...) as keywords, but some
// (e.g. PROVISIONED, THROUGHPUT, BUDGET) arrive as identifiers.
func (p *Parser) isPrivilegeWord(tokType int) bool {
	if tokType == kwON {
		return false
	}
	if tokType == tokIdent || tokType == tokQuotedIdent {
		return true
	}
	// Any keyword token (kw* constants are >= 700) is a valid privilege word.
	return tokType >= 700
}

// parseGrantTarget parses the ON clause body (the ON keyword has already been
// consumed):
//
//	ACCOUNT
//	<object_type> <name> [ ( signature ) ]
//	ALL <object_type_plural> IN { DATABASE <db> | SCHEMA <schema> }
//	FUTURE <object_type_plural> IN { DATABASE <db> | SCHEMA <schema> }
func (p *Parser) parseGrantTarget() (*ast.GrantTarget, error) {
	startLoc := p.cur.Loc

	switch p.cur.Type {
	case kwACCOUNT:
		p.advance() // consume ACCOUNT
		return &ast.GrantTarget{
			Kind: ast.GrantTargetAccount,
			Loc:  ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End},
		}, nil

	case kwALL:
		p.advance() // consume ALL
		return p.parseGrantTargetContainer(ast.GrantTargetAllIn, startLoc)

	case kwFUTURE:
		p.advance() // consume FUTURE
		return p.parseGrantTargetContainer(ast.GrantTargetFutureIn, startLoc)

	default:
		// <object_type> <name> [ ( signature ) ]
		return p.parseGrantTargetObject(startLoc)
	}
}

// parseGrantTargetObject parses the <object_type> <name> [ ( signature ) ] form
// of an ON clause. startLoc is the Loc of the first object-type token.
//
// The object type is open-ended and may be multiple words (TABLE, MATERIALIZED
// VIEW, CORTEX SEARCH SERVICE, STORAGE LIFECYCLE POLICY, ...). The boundary
// between the type and the object name is purely structural: the words of the
// ON target are a sequence of space-separated "name units" (each unit a dotted
// identifier path such as mydb.public.t), and the LAST unit is the object name
// while every preceding unit forms the object type. This needs no fixed type
// vocabulary, so new Snowflake object types parse without code changes.
func (p *Parser) parseGrantTargetObject(startLoc ast.Loc) (*ast.GrantTarget, error) {
	// Collect the type words (everything before the final name unit). We read
	// one name unit, then keep reading while another name unit follows, shifting
	// the previous unit into the type. The final unit is the object name.
	var typeWords []string
	name, err := p.parseGrantNameUnit()
	if err != nil {
		return nil, err
	}

	for p.startsGrantNameUnit() {
		// The unit we just parsed is actually part of the type; absorb it and
		// read the next unit as the new (tentative) name.
		typeWords = append(typeWords, objectNameToType(name))
		name, err = p.parseGrantNameUnit()
		if err != nil {
			return nil, err
		}
	}

	if len(typeWords) == 0 {
		// A bare single unit (e.g. "ON foo") has no object type — invalid.
		return nil, &ParseError{
			Loc: startLoc,
			Msg: "expected object type before object name in GRANT/REVOKE ON clause",
		}
	}

	target := &ast.GrantTarget{
		Kind:       ast.GrantTargetObject,
		ObjectType: strings.Join(typeWords, " "),
		Name:       name,
		Loc:        ast.Loc{Start: startLoc.Start},
	}

	// Optional FUNCTION/PROCEDURE argument-type signature.
	if p.cur.Type == '(' {
		sig, err := p.parseGrantSignature()
		if err != nil {
			return nil, err
		}
		target.Signature = sig
	}

	target.Loc.End = p.prev.Loc.End
	return target, nil
}

// startsGrantNameUnit reports whether the current token can begin another
// space-separated name unit inside an ON target — i.e. it is a name word and
// not a clause terminator. The object form's run ends at TO/FROM (the grantee
// clause), '(' (a signature), comma, or EOF.
func (p *Parser) startsGrantNameUnit() bool {
	switch p.cur.Type {
	case kwTO, kwFROM, '(', ',', tokEOF:
		return false
	}
	return p.isObjectTypeWord(p.cur.Type)
}

// parseGrantNameUnit parses one space-separated name unit: a dotted identifier
// path (one or more name parts joined by dots). Mirrors parseObjectName but is
// used to walk the ON target word-by-word. The type-vs-name split is decided by
// the caller; this only consumes a single dotted unit.
func (p *Parser) parseGrantNameUnit() (*ast.ObjectName, error) {
	if !p.isObjectTypeWord(p.cur.Type) {
		return nil, p.syntaxErrorAtCur()
	}
	return p.parseObjectName()
}

// objectNameToType renders a single-word object-type unit (parsed as an
// ObjectName) back to its uppercased source word. Object-type words are always
// single, unquoted identifiers/keywords, so the unit's Name part carries the
// text. A dotted unit here would be malformed input; we still render it
// faithfully so the error surfaces downstream rather than silently dropping it.
func objectNameToType(n *ast.ObjectName) string {
	return strings.ToUpper(n.String())
}

// parseGrantTargetContainer parses the tail of an ALL/FUTURE target:
//
//	<object_type_plural> IN { DATABASE <db> | SCHEMA <schema> }
//
// kind is GrantTargetAllIn or GrantTargetFutureIn. startLoc is the Loc of the
// ALL/FUTURE keyword. On entry ALL/FUTURE has already been consumed.
func (p *Parser) parseGrantTargetContainer(kind ast.GrantTargetKind, startLoc ast.Loc) (*ast.GrantTarget, error) {
	plural, err := p.parsePluralObjectType()
	if err != nil {
		return nil, err
	}
	target := &ast.GrantTarget{
		Kind:             kind,
		ObjectTypePlural: plural,
		Loc:              ast.Loc{Start: startLoc.Start},
	}

	if _, err := p.expect(kwIN); err != nil {
		return nil, err
	}
	switch p.cur.Type {
	case kwDATABASE:
		p.advance() // consume DATABASE
		target.Container = ast.GrantContainerDatabase
	case kwSCHEMA:
		p.advance() // consume SCHEMA
		target.Container = ast.GrantContainerSchema
	default:
		return nil, p.syntaxErrorAtCur()
	}

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	target.ContainerName = name

	target.Loc.End = p.prev.Loc.End
	return target, nil
}

// parsePluralObjectType parses the plural object-type words of an
// ALL/FUTURE ... IN form (e.g. TABLES, EXTERNAL TABLES, MATERIALIZED VIEWS).
// It consumes name words until the IN keyword, so it is open-ended over
// Snowflake's growing object-type set. The loop also stops at the grantee
// keywords (TO/FROM) so that malformed input missing IN fails at a meaningful
// position rather than swallowing the grantee clause. Requires at least one
// word.
func (p *Parser) parsePluralObjectType() (string, error) {
	if !p.isObjectTypeWord(p.cur.Type) || p.cur.Type == kwIN {
		return "", p.syntaxErrorAtCur()
	}

	var words []string
	for p.isObjectTypeWord(p.cur.Type) {
		switch p.cur.Type {
		case kwIN, kwTO, kwFROM:
			return strings.Join(words, " "), nil
		}
		tok := p.advance()
		words = append(words, strings.ToUpper(tok.Str))
	}
	return strings.Join(words, " "), nil
}

// isObjectTypeWord reports whether tokType can start/continue an object-type or
// object-name word. Identifiers (for object types not in the lexer keyword set,
// e.g. NOTEBOOK, WORKSPACE) and keywords both qualify; punctuation and EOF do
// not.
func (p *Parser) isObjectTypeWord(tokType int) bool {
	if tokType == tokIdent || tokType == tokQuotedIdent {
		return true
	}
	return tokType >= 700
}

// parseGrantSignature parses a FUNCTION/PROCEDURE argument-type signature
// following the object name: ( [ data_type [, data_type ...] ] ). Each argument
// is a full Snowflake data type, so it reuses parseDataType and therefore
// handles parameterized and multi-word types such as NUMBER(38,0),
// VARCHAR(100), DOUBLE PRECISION, and ARRAY. Returns a non-nil (possibly empty)
// slice so callers can distinguish "()" from "no signature". On entry cur is
// '('.
func (p *Parser) parseGrantSignature() ([]*ast.TypeName, error) {
	p.advance() // consume '('

	sig := []*ast.TypeName{}
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		dt, err := p.parseDataType()
		if err != nil {
			return nil, err
		}
		sig = append(sig, dt)
		if p.cur.Type == ',' {
			p.advance() // consume ','
			continue
		}
		break
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return sig, nil
}

// parseGrantee parses the recipient after TO (GRANT) or the subject after FROM
// (REVOKE). The TO/FROM keyword has already been consumed.
//
//	[ROLE] <name>            (bare name defaults to ROLE)
//	DATABASE ROLE <name>
//	USER <name>
//	SHARE <name>
//	APPLICATION <name>
//	APPLICATION ROLE <name>
func (p *Parser) parseGrantee() (*ast.Grantee, error) {
	startLoc := p.cur.Loc
	var kind ast.GranteeKind

	switch p.cur.Type {
	case kwROLE:
		p.advance() // consume ROLE
		kind = ast.GranteeRole
	case kwDATABASE:
		p.advance() // consume DATABASE
		if _, err := p.expect(kwROLE); err != nil {
			return nil, err
		}
		kind = ast.GranteeDatabaseRole
	case kwUSER:
		p.advance() // consume USER
		kind = ast.GranteeUser
	case kwSHARE:
		p.advance() // consume SHARE
		kind = ast.GranteeShare
	case kwAPPLICATION:
		p.advance() // consume APPLICATION
		if p.cur.Type == kwROLE {
			p.advance() // consume ROLE
			kind = ast.GranteeApplicationRole
		} else {
			kind = ast.GranteeApplication
		}
	default:
		// Bare role name: TO <role_name> with the ROLE keyword omitted.
		// Only a name-like token is acceptable here.
		if !p.isGranteeNameStart(p.cur.Type) {
			return nil, p.syntaxErrorAtCur()
		}
		kind = ast.GranteeRole
	}

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}

	return &ast.Grantee{
		Kind: kind,
		Name: name,
		Loc:  ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End},
	}, nil
}

// isGranteeNameStart reports whether tokType can begin a bare grantee role
// name. Identifiers and keywords qualify (so a system role like PUBLIC, which
// the lexer emits as a keyword, is accepted), but structural punctuation and
// EOF do not — guarding against "GRANT ... TO" with nothing following.
func (p *Parser) isGranteeNameStart(tokType int) bool {
	if tokType == tokIdent || tokType == tokQuotedIdent {
		return true
	}
	return tokType >= 700
}

// expectCurrentGrants consumes the CURRENT GRANTS keywords that follow REVOKE
// or COPY in a GRANT OWNERSHIP trailing clause. The REVOKE/COPY keyword has
// already been consumed.
func (p *Parser) expectCurrentGrants() error {
	if _, err := p.expect(kwCURRENT); err != nil {
		return err
	}
	if _, err := p.expect(kwGRANTS); err != nil {
		return err
	}
	return nil
}
