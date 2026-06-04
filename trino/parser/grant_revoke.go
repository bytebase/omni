package parser

import (
	"strings"

	"github.com/bytebase/omni/trino/ast"
)

// This file is the core of the parser-dcl-tcl DAG node: it implements Trino's
// data-control-language (DCL) statements — CREATE/DROP ROLE, GRANT/REVOKE of
// both roles and privileges, and DENY — as hand-written recursive-descent
// parsers over the token stream.
//
// The statement nodes are returned from parseStmt as ast.Node and carry an
// ast.NodeTag (trino/ast/nodetags.go); their concrete fields live here in
// package parser (Trino convention for parser-built node types — see expr.go,
// datatypes.go).
//
// Legacy ANTLR grammar (TrinoParser.g4) the implementation tracks:
//
//	statement
//	    : ...
//	    | CREATE_ ROLE_ name=identifier (WITH_ ADMIN_ grantor)? (IN_ catalog=identifier)?  # createRole
//	    | DROP_ ROLE_ name=identifier (IN_ catalog=identifier)?                            # dropRole
//	    | GRANT_ roles TO_ principal (COMMA_ principal)* (WITH_ ADMIN_ OPTION_)?
//	          (GRANTED_ BY_ grantor)? (IN_ catalog=identifier)?                            # grantRoles
//	    | REVOKE_ (ADMIN_ OPTION_ FOR_)? roles FROM_ principal (COMMA_ principal)*
//	          (GRANTED_ BY_ grantor)? (IN_ catalog=identifier)?                            # revokeRoles
//	    | GRANT_ (privilege (COMMA_ privilege)* | ALL_ PRIVILEGES_)
//	          ON_ (SCHEMA_ | TABLE_)? qualifiedName TO_ grantee=principal
//	          (WITH_ GRANT_ OPTION_)?                                                      # grant
//	    | DENY_ (privilege (COMMA_ privilege)* | ALL_ PRIVILEGES_)
//	          ON_ (SCHEMA_ | TABLE_)? qualifiedName TO_ grantee=principal                  # deny
//	    | REVOKE_ (GRANT_ OPTION_ FOR_)? (privilege (COMMA_ privilege)* | ALL_ PRIVILEGES_)
//	          ON_ (SCHEMA_ | TABLE_)? qualifiedName FROM_ grantee=principal                # revoke
//	    ;
//	privilege : CREATE_ | SELECT_ | DELETE_ | INSERT_ | UPDATE_ ;
//	principal : identifier | USER_ identifier | ROLE_ identifier ;
//	grantor   : principal | CURRENT_USER_ | CURRENT_ROLE_ ;
//	roles     : identifier (COMMA_ identifier)* ;
//
// Trino 481 (truth1 + the live oracle) diverges from this stale legacy grammar
// in three oracle-confirmed ways (recorded in the divergence ledger):
//
//	D1 (privilege vocabulary is broadened). A privilege is the five legacy
//	   reserved keywords CREATE/SELECT/DELETE/INSERT/UPDATE PLUS any non-reserved
//	   identifier: 481 accepts `GRANT mypriv ON test TO alice`, `GRANT "My Priv"
//	   ON ...`, and even `GRANT TO ON t ...` (TO is non-reserved). It is NOT "any
//	   keyword" — other reserved keywords (FROM, WHERE, TABLE, ...) are rejected
//	   as privileges. Validity beyond this is the access-control layer's concern.
//	D2 (ON BRANCH … IN target). 481 accepts a branch-qualified target
//	   `GRANT INSERT ON BRANCH audit IN orders TO alice` (and the REVOKE/DENY
//	   forms), absent from the legacy grammar — the docs-ahead "@branch"
//	   extension proven real by the oracle.
//	D3 (bare ALL). 481 accepts `GRANT ALL ON test TO alice` (ALL without the
//	   PRIVILEGES keyword) as an all-privileges grant; the legacy grammar
//	   requires `ALL PRIVILEGES`.
//	D5 (open entity-kind qualifier). The ON-target entity kind is an arbitrary
//	   single identifier in 481, not the legacy grammar's fixed `(SCHEMA |
//	   TABLE)?`: `ON VIEW v`, `ON a b` (qualifier "a", object "b") are accepted.
//	   At most one leading qualifier word is allowed (`ON a b c` is a
//	   SYNTAX_ERROR), and it must be a single non-dotted identifier (`ON a.b c`
//	   is a SYNTAX_ERROR).
//
// Disambiguation (the central difficulty): GRANT/REVOKE/DENY each have a
// role-grant and a privilege-grant alternative that share a leading
// comma-separated name list. They are told apart purely structurally by the
// keyword that terminates the list — `ON` ⇒ privilege grant; `TO` (GRANT/DENY)
// or `FROM` (REVOKE) ⇒ role grant. scanGrantFormIsPrivilege performs a bounded
// speculative scan to that terminator. The two forms then differ in their list
// element rules (privileges allow reserved keywords; roles require valid
// identifiers — `GRANT SELECT TO foo` is a SYNTAX_ERROR because reserved SELECT
// is not a legal role) and in their trailing options (WITH GRANT OPTION vs
// WITH ADMIN OPTION / GRANTED BY / IN catalog), all oracle-adjudicated.
//
// The implementation is adjudicated against the live Trino 481 oracle.

// ---------------------------------------------------------------------------
// Shared sub-grammar AST
// ---------------------------------------------------------------------------

// PrincipalKind classifies a principal: a bare name (defaults to a role/user
// per Trino semantics), an explicit USER, or an explicit ROLE.
type PrincipalKind int

const (
	// PrincipalUnspecified is a bare `identifier` principal (no USER/ROLE
	// keyword).
	PrincipalUnspecified PrincipalKind = iota
	// PrincipalUser is `USER identifier`.
	PrincipalUser
	// PrincipalRole is `ROLE identifier`.
	PrincipalRole
)

// Principal is a `principal` grammar node: an optionally USER/ROLE-qualified
// name. Name is nil only for the keyword grantor forms handled by Grantor.
type Principal struct {
	Kind PrincipalKind
	Name *ast.Identifier
	Loc  ast.Loc
}

// GrantorKind classifies a GRANTED BY grantor: a specified principal, the
// CURRENT_USER keyword, or the CURRENT_ROLE keyword.
type GrantorKind int

const (
	// GrantorPrincipal is a grantor given as a principal.
	GrantorPrincipal GrantorKind = iota
	// GrantorCurrentUser is the CURRENT_USER keyword grantor.
	GrantorCurrentUser
	// GrantorCurrentRole is the CURRENT_ROLE keyword grantor.
	GrantorCurrentRole
)

// Grantor is a `grantor` grammar node (the principal after GRANTED BY).
// Principal is set only when Kind is GrantorPrincipal.
type Grantor struct {
	Kind      GrantorKind
	Principal *Principal
	Loc       ast.Loc
}

// GrantTarget is the ON clause of a privilege GRANT/REVOKE/DENY. Trino 481's
// actual grammar (decoded from the live oracle — see the file-header D5 note) is
//
//	ON [ BRANCH branch_name IN ] [ entity_kind ] qualifiedName
//
// where:
//   - branch_name (the D2 docs-ahead extension) is a single identifier; it is
//     recognized only as `BRANCH <ident> IN …`, so a target literally named
//     "branch" (`ON branch`) is an ordinary object name.
//   - entity_kind is an OPTIONAL single leading qualifier word. The legacy
//     grammar hard-coded it to `(SCHEMA | TABLE)?`, but 481 accepts ANY single
//     non-reserved identifier here (`ON SCHEMA s`, `ON TABLE t`, `ON VIEW v`,
//     `ON a b` are all accepted) — so it is captured as a generic Identifier,
//     not a fixed enum (D5).
//   - qualifiedName is the object name (may be dotted: catalog.schema.table).
//
// Field meanings: Branch is the branch identifier (nil unless the BRANCH form).
// Qualifier is the optional entity-kind word (nil when the target is a bare
// object name). Object is the object qualifiedName. For the BRANCH form,
// Qualifier/Object describe the entity that follows IN.
type GrantTarget struct {
	Branch    *ast.Identifier    // nil unless ON BRANCH <branch> IN ...
	Qualifier *ast.Identifier    // nil unless an entity-kind word precedes the name
	Object    *ast.QualifiedName // the object name (may be dotted)
	Loc       ast.Loc
}

// IsSchema reports whether the target's entity-kind qualifier is the SCHEMA
// keyword (case-insensitive), the common "grant on a schema" form.
func (t *GrantTarget) IsSchema() bool {
	return t.Qualifier != nil && strings.EqualFold(t.Qualifier.Value, "SCHEMA")
}

// IsTable reports whether the target's entity-kind qualifier is the TABLE
// keyword (case-insensitive).
func (t *GrantTarget) IsTable() bool {
	return t.Qualifier != nil && strings.EqualFold(t.Qualifier.Value, "TABLE")
}

// ---------------------------------------------------------------------------
// Statement AST
// ---------------------------------------------------------------------------

// CreateRoleStmt is CREATE ROLE <name> [WITH ADMIN <grantor>] [IN <catalog>].
type CreateRoleStmt struct {
	Name    *ast.Identifier
	Admin   *Grantor        // nil unless WITH ADMIN <grantor>
	Catalog *ast.Identifier // nil unless IN <catalog>
	Loc     ast.Loc
}

// Tag implements ast.Node.
func (n *CreateRoleStmt) Tag() ast.NodeTag { return ast.T_CreateRoleStmt }

// DropRoleStmt is DROP ROLE [IF EXISTS] <name> [IN <catalog>].
//
// IfExists is the D4 docs-ahead extension: 481 accepts DROP ROLE IF EXISTS,
// absent from the legacy grammar.
type DropRoleStmt struct {
	IfExists bool
	Name     *ast.Identifier
	Catalog  *ast.Identifier // nil unless IN <catalog>
	Loc      ast.Loc
}

// Tag implements ast.Node.
func (n *DropRoleStmt) Tag() ast.NodeTag { return ast.T_DropRoleStmt }

// GrantRolesStmt is
//
//	GRANT <roles> TO <principal> [, ...] [WITH ADMIN OPTION]
//	    [GRANTED BY <grantor>] [IN <catalog>].
type GrantRolesStmt struct {
	Roles       []*ast.Identifier
	Grantees    []*Principal
	AdminOption bool
	GrantedBy   *Grantor        // nil unless GRANTED BY <grantor>
	Catalog     *ast.Identifier // nil unless IN <catalog>
	Loc         ast.Loc
}

// Tag implements ast.Node.
func (n *GrantRolesStmt) Tag() ast.NodeTag { return ast.T_GrantRolesStmt }

// RevokeRolesStmt is
//
//	REVOKE [ADMIN OPTION FOR] <roles> FROM <principal> [, ...]
//	    [GRANTED BY <grantor>] [IN <catalog>].
type RevokeRolesStmt struct {
	AdminOptionFor bool
	Roles          []*ast.Identifier
	Grantees       []*Principal
	GrantedBy      *Grantor        // nil unless GRANTED BY <grantor>
	Catalog        *ast.Identifier // nil unless IN <catalog>
	Loc            ast.Loc
}

// Tag implements ast.Node.
func (n *RevokeRolesStmt) Tag() ast.NodeTag { return ast.T_RevokeRolesStmt }

// GrantPrivStmt is
//
//	GRANT ( <privileges> | ALL [PRIVILEGES] ) ON <target> TO <grantee>
//	    [WITH GRANT OPTION].
//
// AllPrivileges is true for the ALL [PRIVILEGES] form (Privileges is then nil).
type GrantPrivStmt struct {
	AllPrivileges bool
	Privileges    []*ast.Identifier
	On            *GrantTarget
	Grantee       *Principal
	GrantOption   bool
	Loc           ast.Loc
}

// Tag implements ast.Node.
func (n *GrantPrivStmt) Tag() ast.NodeTag { return ast.T_GrantPrivStmt }

// RevokePrivStmt is
//
//	REVOKE [GRANT OPTION FOR] ( <privileges> | ALL [PRIVILEGES] )
//	    ON <target> FROM <grantee>.
type RevokePrivStmt struct {
	GrantOptionFor bool
	AllPrivileges  bool
	Privileges     []*ast.Identifier
	On             *GrantTarget
	Grantee        *Principal
	Loc            ast.Loc
}

// Tag implements ast.Node.
func (n *RevokePrivStmt) Tag() ast.NodeTag { return ast.T_RevokePrivStmt }

// DenyStmt is
//
//	DENY ( <privileges> | ALL [PRIVILEGES] ) ON <target> TO <grantee>.
type DenyStmt struct {
	AllPrivileges bool
	Privileges    []*ast.Identifier
	On            *GrantTarget
	Grantee       *Principal
	Loc           ast.Loc
}

// Tag implements ast.Node.
func (n *DenyStmt) Tag() ast.NodeTag { return ast.T_DenyStmt }

// Compile-time assertions that the statement types satisfy ast.Node.
var (
	_ ast.Node = (*CreateRoleStmt)(nil)
	_ ast.Node = (*DropRoleStmt)(nil)
	_ ast.Node = (*GrantRolesStmt)(nil)
	_ ast.Node = (*RevokeRolesStmt)(nil)
	_ ast.Node = (*GrantPrivStmt)(nil)
	_ ast.Node = (*RevokePrivStmt)(nil)
	_ ast.Node = (*DenyStmt)(nil)
)

// ---------------------------------------------------------------------------
// CREATE / DROP ROLE
// ---------------------------------------------------------------------------

// parseCreateRoleStmt parses
// CREATE ROLE <name> [WITH ADMIN <grantor>] [IN <catalog>].
//
// On entry the CREATE and ROLE keywords have BOTH already been consumed by the
// dispatcher (CREATE is a multi-object keyword); start is the CREATE keyword's
// byte offset.
func (p *Parser) parseCreateRoleStmt(start int) (ast.Node, error) {
	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	stmt := &CreateRoleStmt{Name: name, Loc: ast.Loc{Start: start, End: p.prev.Loc.End}}

	// Optional WITH ADMIN <grantor>.
	if p.cur.Kind == kwWITH && p.peekNext().Kind == kwADMIN {
		p.advance() // consume WITH
		p.advance() // consume ADMIN
		grantor, err := p.parseGrantor()
		if err != nil {
			return nil, err
		}
		stmt.Admin = grantor
	}

	// Optional IN <catalog>.
	if catalog, err := p.parseOptionalInCatalog(); err != nil {
		return nil, err
	} else {
		stmt.Catalog = catalog
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseDropRoleStmt parses DROP ROLE [IF EXISTS] <name> [IN <catalog>].
//
// On entry the DROP and ROLE keywords have both already been consumed; start is
// the DROP keyword's byte offset.
func (p *Parser) parseDropRoleStmt(start int) (ast.Node, error) {
	stmt := &DropRoleStmt{Loc: ast.Loc{Start: start}}

	// Optional IF EXISTS (D4 docs-ahead extension).
	if p.cur.Kind == kwIF && p.peekNext().Kind == kwEXISTS {
		p.advance() // consume IF
		p.advance() // consume EXISTS
		stmt.IfExists = true
	}

	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	if catalog, err := p.parseOptionalInCatalog(); err != nil {
		return nil, err
	} else {
		stmt.Catalog = catalog
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// GRANT / REVOKE / DENY dispatch
// ---------------------------------------------------------------------------

// parseGrantStmt parses a GRANT statement, dispatching between the role-grant
// and privilege-grant forms by the structural lookahead in
// scanGrantFormIsPrivilege. On entry cur is the GRANT keyword.
func (p *Parser) parseGrantStmt() (ast.Node, error) {
	grantTok := p.advance() // consume GRANT
	start := grantTok.Loc.Start
	if p.scanGrantFormIsPrivilege(false) {
		return p.parseGrantPrivileges(start)
	}
	return p.parseGrantRoles(start)
}

// parseRevokeStmt parses a REVOKE statement. A leading ADMIN OPTION FOR forces
// the role form; a leading GRANT OPTION FOR forces the privilege form;
// otherwise the ON-vs-FROM lookahead decides. On entry cur is the REVOKE
// keyword.
func (p *Parser) parseRevokeStmt() (ast.Node, error) {
	revokeTok := p.advance() // consume REVOKE
	start := revokeTok.Loc.Start

	if p.cur.Kind == kwADMIN && p.peekNext().Kind == kwOPTION {
		return p.parseRevokeRoles(start)
	}
	if p.cur.Kind == kwGRANT && p.peekNext().Kind == kwOPTION {
		return p.parseRevokePrivileges(start)
	}
	if p.scanGrantFormIsPrivilege(true) {
		return p.parseRevokePrivileges(start)
	}
	return p.parseRevokeRoles(start)
}

// parseDenyStmt parses a DENY statement. DENY has only the privilege form
// (there is no role-deny), so it always parses as a privilege grant target.
// On entry cur is the DENY keyword.
func (p *Parser) parseDenyStmt() (ast.Node, error) {
	denyTok := p.advance() // consume DENY
	start := denyTok.Loc.Start

	stmt := &DenyStmt{Loc: ast.Loc{Start: start}}
	all, privs, err := p.parsePrivilegeSpec()
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
	grantee, err := p.parsePrincipal()
	if err != nil {
		return nil, err
	}
	stmt.Grantee = grantee

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// scanGrantFormIsPrivilege performs a bounded speculative scan from the current
// token (the first token after GRANT/REVOKE, or after a GRANT/ADMIN OPTION FOR
// prefix) to decide whether the statement is the privilege form (terminated by
// ON) or the role form (terminated by TO for GRANT/DENY, FROM for REVOKE). It
// rewinds via a checkpoint so the caller's cursor is unchanged.
//
// The scan is STRUCTURE-AWARE rather than a naive search for the first ON/TO/
// FROM token: a privilege or role NAME may itself be the non-reserved word TO
// or FROM (`GRANT TO ON t TO alice` grants a privilege literally named "TO").
// So it walks the leading `item (',' item)*` list — each item one token, except
// the two-token ALL PRIVILEGES — and only inspects the token that sits at a
// list boundary (right after an item, where a comma or the terminator must
// appear). ON at a boundary ⇒ privilege; TO/FROM at a boundary ⇒ role. A
// boundary token that is neither (malformed, or EOF) defaults to the privilege
// form so the concrete privilege parser emits a precise "expected ON" error.
func (p *Parser) scanGrantFormIsPrivilege(isRevoke bool) bool {
	cp := p.checkpoint()
	defer p.restore(cp)

	for {
		// Consume one list item. The first token after the keyword is always an
		// item (even if it is the non-reserved word TO/FROM/ON), never the
		// terminator. ALL PRIVILEGES is a single two-token item.
		if p.cur.Kind == tokEOF {
			return true
		}
		if p.cur.Kind == kwALL && p.peekNext().Kind == kwPRIVILEGES {
			p.advance() // ALL
			p.advance() // PRIVILEGES
		} else {
			p.advance() // the single-token item
		}

		// Now at a list boundary: a comma continues the list; ON/TO/FROM
		// terminates it; anything else is malformed (default to privilege).
		switch p.cur.Kind {
		case int(','):
			p.advance() // consume comma, read the next item
			continue
		case kwON:
			return true
		case kwTO:
			return false
		case kwFROM:
			if isRevoke {
				return false
			}
			return true
		default:
			return true
		}
	}
}

// ---------------------------------------------------------------------------
// GRANT / REVOKE roles
// ---------------------------------------------------------------------------

// parseGrantRoles parses the role-grant form (the GRANT keyword consumed):
//
//	<roles> TO <principal> [, ...] [WITH ADMIN OPTION]
//	    [GRANTED BY <grantor>] [IN <catalog>].
func (p *Parser) parseGrantRoles(start int) (ast.Node, error) {
	stmt := &GrantRolesStmt{Loc: ast.Loc{Start: start}}

	roles, err := p.parseRoleList()
	if err != nil {
		return nil, err
	}
	stmt.Roles = roles

	if _, err := p.expect(kwTO); err != nil {
		return nil, err
	}
	grantees, err := p.parsePrincipalList()
	if err != nil {
		return nil, err
	}
	stmt.Grantees = grantees

	// Optional WITH ADMIN OPTION.
	if p.cur.Kind == kwWITH && p.peekNext().Kind == kwADMIN {
		p.advance() // consume WITH
		p.advance() // consume ADMIN
		if _, err := p.expect(kwOPTION); err != nil {
			return nil, err
		}
		stmt.AdminOption = true
	}

	// Optional GRANTED BY <grantor>.
	if grantedBy, err := p.parseOptionalGrantedBy(); err != nil {
		return nil, err
	} else {
		stmt.GrantedBy = grantedBy
	}

	// Optional IN <catalog>.
	if catalog, err := p.parseOptionalInCatalog(); err != nil {
		return nil, err
	} else {
		stmt.Catalog = catalog
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseRevokeRoles parses the role-revoke form (the REVOKE keyword consumed):
//
//	[ADMIN OPTION FOR] <roles> FROM <principal> [, ...]
//	    [GRANTED BY <grantor>] [IN <catalog>].
func (p *Parser) parseRevokeRoles(start int) (ast.Node, error) {
	stmt := &RevokeRolesStmt{Loc: ast.Loc{Start: start}}

	// Optional ADMIN OPTION FOR.
	if p.cur.Kind == kwADMIN && p.peekNext().Kind == kwOPTION {
		p.advance() // consume ADMIN
		p.advance() // consume OPTION
		if _, err := p.expect(kwFOR); err != nil {
			return nil, err
		}
		stmt.AdminOptionFor = true
	}

	roles, err := p.parseRoleList()
	if err != nil {
		return nil, err
	}
	stmt.Roles = roles

	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}
	grantees, err := p.parsePrincipalList()
	if err != nil {
		return nil, err
	}
	stmt.Grantees = grantees

	if grantedBy, err := p.parseOptionalGrantedBy(); err != nil {
		return nil, err
	} else {
		stmt.GrantedBy = grantedBy
	}
	if catalog, err := p.parseOptionalInCatalog(); err != nil {
		return nil, err
	} else {
		stmt.Catalog = catalog
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseRoleList parses `roles := identifier (COMMA identifier)*`. Each role is
// a plain identifier; a reserved keyword is rejected here (so `GRANT SELECT TO
// foo` fails, matching Trino).
func (p *Parser) parseRoleList() ([]*ast.Identifier, error) {
	first, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	roles := []*ast.Identifier{first}
	for {
		if _, ok := p.match(int(',')); !ok {
			break
		}
		next, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		roles = append(roles, next)
	}
	return roles, nil
}

// ---------------------------------------------------------------------------
// GRANT / REVOKE privileges
// ---------------------------------------------------------------------------

// parseGrantPrivileges parses the privilege-grant form (the GRANT keyword
// consumed):
//
//	( <privileges> | ALL [PRIVILEGES] ) ON <target> TO <grantee>
//	    [WITH GRANT OPTION].
func (p *Parser) parseGrantPrivileges(start int) (ast.Node, error) {
	stmt := &GrantPrivStmt{Loc: ast.Loc{Start: start}}

	all, privs, err := p.parsePrivilegeSpec()
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
	grantee, err := p.parsePrincipal()
	if err != nil {
		return nil, err
	}
	stmt.Grantee = grantee

	// Optional WITH GRANT OPTION.
	if p.cur.Kind == kwWITH && p.peekNext().Kind == kwGRANT {
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

// parseRevokePrivileges parses the privilege-revoke form (the REVOKE keyword
// consumed):
//
//	[GRANT OPTION FOR] ( <privileges> | ALL [PRIVILEGES] )
//	    ON <target> FROM <grantee>.
func (p *Parser) parseRevokePrivileges(start int) (ast.Node, error) {
	stmt := &RevokePrivStmt{Loc: ast.Loc{Start: start}}

	// Optional GRANT OPTION FOR.
	if p.cur.Kind == kwGRANT && p.peekNext().Kind == kwOPTION {
		p.advance() // consume GRANT
		p.advance() // consume OPTION
		if _, err := p.expect(kwFOR); err != nil {
			return nil, err
		}
		stmt.GrantOptionFor = true
	}

	all, privs, err := p.parsePrivilegeSpec()
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
	grantee, err := p.parsePrincipal()
	if err != nil {
		return nil, err
	}
	stmt.Grantee = grantee

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parsePrivilegeSpec parses the privilege specification before ON:
//
//	ALL [PRIVILEGES]            -> (true, nil, nil)
//	privilege [, privilege ...] -> (false, [...], nil)
//
// Per oracle divergence D3, a bare ALL immediately followed by ON is the
// all-privileges form. A privilege is an arbitrary single identifier-or-keyword
// (D1) — privilege validity is not a grammar concern — so reserved keywords
// like SELECT/DELETE/UPDATE and ordinary identifiers both parse.
func (p *Parser) parsePrivilegeSpec() (bool, []*ast.Identifier, error) {
	// ALL [PRIVILEGES] — but only when ALL is the entire spec, i.e. it is
	// followed by PRIVILEGES or directly by ON. (A privilege/role literally
	// named ALL followed by a comma is not a documented Trino privilege form;
	// in the role grant `GRANT ALL TO foo`, ALL is parsed as a role elsewhere.)
	if p.cur.Kind == kwALL {
		next := p.peekNext().Kind
		if next == kwPRIVILEGES {
			p.advance() // consume ALL
			p.advance() // consume PRIVILEGES
			return true, nil, nil
		}
		if next == kwON {
			p.advance() // consume ALL
			return true, nil, nil
		}
	}

	var privs []*ast.Identifier
	for {
		priv, err := p.parsePrivilege()
		if err != nil {
			return false, nil, err
		}
		privs = append(privs, priv)
		if _, ok := p.match(int(',')); !ok {
			break
		}
	}
	return false, privs, nil
}

// parsePrivilege parses a single privilege, captured as an ast.Identifier (a
// quoted token keeps its quoting metadata). A privilege is one token; ON, a
// comma, or a non-privilege token terminates the list.
func (p *Parser) parsePrivilege() (*ast.Identifier, error) {
	if !isPrivilegeWord(p.cur.Kind) {
		return nil, p.syntaxErrorAtCur()
	}
	return identFromToken(p.advance()), nil
}

// isPrivilegeWord reports whether a token kind can be a privilege name. Per the
// oracle-pinned D1 vocabulary, a privilege is:
//   - an unquoted/quoted identifier (incl. any NON-reserved keyword, since the
//     lexer emits those as identifiers via isIdentifierStart) — e.g. mypriv,
//     "My Priv", and the non-reserved word TO; OR
//   - one of the five legacy reserved privilege keywords SELECT / INSERT /
//     UPDATE / DELETE / CREATE.
//
// It is NOT "any keyword": other reserved keywords (FROM, WHERE, TABLE, DROP, …)
// are rejected as privileges (`GRANT FROM ON t …` is a Trino SYNTAX_ERROR),
// which is why this cannot simply test kind >= 700.
func isPrivilegeWord(kind TokenKind) bool {
	if isIdentifierStart(kind) {
		return true
	}
	switch kind {
	case kwSELECT, kwINSERT, kwUPDATE, kwDELETE, kwCREATE:
		return true
	default:
		return false
	}
}

// parseGrantTarget parses the ON clause body (the ON keyword already consumed):
//
//	[ BRANCH branch_name IN ] [ entity_kind ] qualifiedName
//
// matching Trino 481 exactly (see the GrantTarget doc). The BRANCH prefix is the
// D2 docs-ahead extension and the open entity_kind qualifier is D5; both are
// oracle-confirmed. The clause ends at the grantee keyword (TO for GRANT/DENY,
// FROM for REVOKE) or a '(' — those terminate the trailing entity.
func (p *Parser) parseGrantTarget() (*GrantTarget, error) {
	startLoc := p.cur.Loc
	target := &GrantTarget{Loc: ast.Loc{Start: startLoc.Start}}

	// Optional BRANCH <branch> IN — recognized only as the three-token shape
	// `BRANCH <ident> IN`, so a target named "branch" (`ON branch`, or the
	// entity-kind word in `ON branch x`) stays an ordinary name.
	if p.branchClauseFollows() {
		p.advance() // consume BRANCH (an unquoted identifier)
		branch, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(kwIN); err != nil {
			return nil, err
		}
		target.Branch = branch
	}

	// [ entity_kind ] qualifiedName.
	qualifier, object, err := p.parseGrantEntity()
	if err != nil {
		return nil, err
	}
	target.Qualifier = qualifier
	target.Object = object

	target.Loc.End = p.prev.Loc.End
	return target, nil
}

// parseGrantEntity parses an ON-target entity: an OPTIONAL single leading
// entity-kind qualifier word followed by the object qualifiedName. It is used
// both for the plain ON target and for the IN-target of the BRANCH form.
//
// The entity-kind qualifier (D5) is `TABLE | SCHEMA | <non-reserved identifier>`
// — the two legacy reserved keywords plus any identifier; other reserved
// keywords (SELECT, FROM, …) are rejected. It is present iff a SECOND name unit
// follows the first before the terminator (TO/FROM/'('/EOF). A dotted qualifier
// (`ON a.b c`) is rejected — the qualifier is a single identifier.
func (p *Parser) parseGrantEntity() (*ast.Identifier, *ast.QualifiedName, error) {
	// The reserved keywords TABLE/SCHEMA can only be the qualifier (they are not
	// valid identifiers, so parseQualifiedName would reject them as an object).
	if p.cur.Kind == kwTABLE || p.cur.Kind == kwSCHEMA {
		qualifier := identFromToken(p.advance())
		object, err := p.parseQualifiedName()
		if err != nil {
			return nil, nil, err
		}
		return qualifier, object, nil
	}

	// Otherwise the leading unit is an identifier — either the object (one unit)
	// or the qualifier (a second unit follows).
	first, err := p.parseQualifiedName()
	if err != nil {
		return nil, nil, err
	}
	if p.startsGrantEntityUnit() {
		if len(first.Parts) != 1 {
			return nil, nil, &ParseError{Loc: first.Loc, Msg: "entity kind in GRANT/REVOKE ON clause must be a single identifier"}
		}
		object, err := p.parseQualifiedName()
		if err != nil {
			return nil, nil, err
		}
		return first.Parts[0], object, nil
	}
	return nil, first, nil
}

// branchClauseFollows reports whether the current position begins a
// `BRANCH <ident> IN` clause: the current token is the unquoted identifier
// "BRANCH" (case-insensitive), followed by a single identifier, followed by the
// IN keyword. It uses a bounded speculative scan via a checkpoint so the
// caller's cursor is unchanged. The three-token shape is required so that a
// plain `ON branch <object>` (entity-kind word "branch") or `ON branch` (object
// named branch) is NOT mistaken for the branch form.
func (p *Parser) branchClauseFollows() bool {
	if p.cur.Kind != tokIdent || !strings.EqualFold(p.cur.Str, "BRANCH") {
		return false
	}
	cp := p.checkpoint()
	defer p.restore(cp)
	p.advance() // consume BRANCH
	if !isIdentifierStart(p.cur.Kind) {
		return false
	}
	p.advance() // consume the branch name
	return p.cur.Kind == kwIN
}

// startsGrantEntityUnit reports whether the current token can begin another
// space-separated name unit inside an ON target — used to decide whether the
// already-parsed qualifiedName was the entity-kind qualifier (another unit
// follows) or the object (the clause ends). The entity run ends at the grantee
// keyword (TO/FROM), a '(' (a future signature), or EOF.
func (p *Parser) startsGrantEntityUnit() bool {
	switch p.cur.Kind {
	case kwTO, kwFROM, int('('), tokEOF:
		return false
	}
	return isIdentifierStart(p.cur.Kind)
}

// ---------------------------------------------------------------------------
// principal / grantor
// ---------------------------------------------------------------------------

// parsePrincipalList parses `principal (COMMA principal)*` (the grantee list of
// a role grant/revoke).
func (p *Parser) parsePrincipalList() ([]*Principal, error) {
	first, err := p.parsePrincipal()
	if err != nil {
		return nil, err
	}
	principals := []*Principal{first}
	for {
		if _, ok := p.match(int(',')); !ok {
			break
		}
		next, err := p.parsePrincipal()
		if err != nil {
			return nil, err
		}
		principals = append(principals, next)
	}
	return principals, nil
}

// parsePrincipal parses a `principal`:
//
//	identifier | USER identifier | ROLE identifier
//
// USER and ROLE are NON-RESERVED keywords, so they act as the USER/ROLE
// qualifier ONLY when directly followed by an identifier (the principal name).
// A bare `USER` / `ROLE` not followed by a name is itself the (unspecified)
// principal name — Trino accepts `GRANT r TO USER` and `GRANT r TO ROLE, bar`.
// A bare identifier principal is PrincipalUnspecified.
func (p *Parser) parsePrincipal() (*Principal, error) {
	startLoc := p.cur.Loc
	switch p.cur.Kind {
	case kwUSER:
		if isIdentifierStart(p.peekNext().Kind) {
			p.advance() // consume USER (qualifier)
			name, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			return &Principal{Kind: PrincipalUser, Name: name, Loc: ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End}}, nil
		}
	case kwROLE:
		if isIdentifierStart(p.peekNext().Kind) {
			p.advance() // consume ROLE (qualifier)
			name, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			return &Principal{Kind: PrincipalRole, Name: name, Loc: ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End}}, nil
		}
	}
	// Bare name (includes a bare USER/ROLE not used as a qualifier).
	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	return &Principal{Kind: PrincipalUnspecified, Name: name, Loc: ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End}}, nil
}

// parseGrantor parses a `grantor`:
//
//	principal | CURRENT_USER | CURRENT_ROLE
//
// USER/ROLE-qualified and bare-identifier grantors go through parsePrincipal.
func (p *Parser) parseGrantor() (*Grantor, error) {
	startLoc := p.cur.Loc
	switch p.cur.Kind {
	case kwCURRENT_USER:
		p.advance()
		return &Grantor{Kind: GrantorCurrentUser, Loc: ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End}}, nil
	case kwCURRENT_ROLE:
		p.advance()
		return &Grantor{Kind: GrantorCurrentRole, Loc: ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End}}, nil
	default:
		principal, err := p.parsePrincipal()
		if err != nil {
			return nil, err
		}
		return &Grantor{Kind: GrantorPrincipal, Principal: principal, Loc: ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End}}, nil
	}
}

// parseOptionalGrantedBy parses an optional `GRANTED BY <grantor>` clause,
// returning nil when absent.
func (p *Parser) parseOptionalGrantedBy() (*Grantor, error) {
	if p.cur.Kind != kwGRANTED {
		return nil, nil
	}
	p.advance() // consume GRANTED
	if _, err := p.expect(kwBY); err != nil {
		return nil, err
	}
	return p.parseGrantor()
}

// parseOptionalInCatalog parses an optional `IN <catalog>` clause, returning nil
// when absent. The catalog is a plain identifier.
func (p *Parser) parseOptionalInCatalog() (*ast.Identifier, error) {
	if p.cur.Kind != kwIN {
		return nil, nil
	}
	p.advance() // consume IN
	return p.parseIdentifier()
}
