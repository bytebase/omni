package parser

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/redshift/ast"
)

func (p *Parser) parseGrantStmt() (nodes.Node, error) {
	loc := p.prev.Loc // position of GRANT keyword
	if p.redshiftKeywordIs("assumerole") {
		return p.parseRedshiftAssumeRoleStmt(loc, true)
	}
	if p.cur.Type == ROLE {
		p.advance()
		roles, err := p.parseRedshiftGrantRoleNameList()
		if err != nil {
			return nil, err
		}
		return p.finishGrantRole(loc, true, roles)
	}
	if p.cur.Type == CREATE && p.redshiftTokIsWord(p.peekNext(), "model") {
		privLoc := nodes.Loc{Start: p.pos(), End: p.cur.End}
		p.advance() // CREATE
		privLoc.End = p.cur.End
		p.advance() // MODEL
		return p.finishRedshiftModelPrivilege(loc, privLoc, true, false)
	}
	if p.isRedshiftConnectStart() {
		return p.parseRedshiftConnectStmt(loc, true)
	}
	if p.isRedshiftDebugStart() {
		return p.parseRedshiftDebugStmt(loc, true)
	}
	if p.isRedshiftScopedPermissionStart() {
		privs, _, err := p.parseRedshiftGrantPrivileges()
		if err != nil {
			return nil, err
		}
		return p.finishRedshiftScopedPermission(loc, true, false, privs)
	}
	if roles, ok, err := p.parseRedshiftSystemPrivileges(TO); ok || err != nil {
		if err != nil {
			return nil, err
		}
		return p.finishGrantRole(loc, true, roles)
	}
	privs, bareAll, err := p.parseRedshiftGrantPrivileges()
	if err != nil {
		return nil, err
	}
	switch {
	case p.cur.Type == FOR:
		return p.finishRedshiftScopedPermission(loc, true, false, privs)
	case p.cur.Type == ON:
		p.advance()
		return p.finishGrantOnObject(loc, privs)
	case bareAll:
		// A bare ALL with no ON clause is the system privilege r_GRANT writes as
		// { ALL [ PRIVILEGES ] } TO ROLE role_name [, ...].
		roles := &nodes.List{Items: []nodes.Node{&nodes.AccessPriv{PrivName: "all"}}}
		for p.cur.Type == ',' {
			p.advance()
			name, err := p.parseColId()
			if err != nil {
				return nil, err
			}
			roles.Items = append(roles.Items, &nodes.AccessPriv{PrivName: name})
		}
		if p.cur.Type != TO {
			return nil, p.syntaxErrorAtCur()
		}
		return p.finishGrantRole(loc, true, roles)
	default:
		return p.finishGrantRole(loc, true, privs)
	}
}

func (p *Parser) parseRevokeStmt() (nodes.Node, error) {
	loc := p.prev.Loc // position of REVOKE keyword
	if p.redshiftKeywordIs("assumerole") {
		return p.parseRedshiftAssumeRoleStmt(loc, false)
	}
	grantOptionFor := false
	grantOptionWithoutFor := false
	grantOptionWithFor := false
	roleOptionName := "" // non-empty => "<name> OPTION FOR" form
	if p.cur.Type == GRANT {
		// r_REVOKE's scoped forms write 'GRANT_OPTION' as one token, and the
		// notes beside them write GRANT OPTION: the one-token spelling is the
		// docs' notation, not SQL, and stays rejected.
		p.advance()
		if _, err := p.expect(OPTION); err != nil {
			return nil, err
		}
		grantOptionFor = true
		if p.cur.Type == FOR {
			p.advance()
			grantOptionWithFor = true
		} else {
			grantOptionWithoutFor = true
		}
	} else if p.isColId() && p.peekNext().Type == OPTION {
		// PG's grammar allows any ColId before OPTION FOR in REVOKE role:
		//   REVOKE ColId OPTION FOR privilege_list FROM role_list
		// ADMIN/INHERIT/SET are the meaningful names; PG parses any ColId
		// and defers validation to semantic time.
		roleOptionName = strings.ToLower(p.cur.Str)
		p.advance() // ColId
		if _, err := p.expect(OPTION); err != nil {
			return nil, err
		}
		if _, err := p.expect(FOR); err != nil {
			return nil, err
		}
	}
	if grantOptionWithoutFor {
		// The docs write GRANT OPTION without FOR only before the scoped forms,
		// whose privileges end at FOR.
		privs, _, err := p.parseRedshiftGrantPrivileges()
		if err != nil {
			return nil, err
		}
		if p.cur.Type != FOR {
			return nil, p.syntaxErrorAtCur()
		}
		return p.finishRedshiftScopedPermission(loc, false, grantOptionFor, privs)
	}
	if p.cur.Type == ROLE {
		roles, err := p.parseRedshiftGrantRoleNameList()
		if err != nil {
			return nil, err
		}
		return p.finishRevokeRole(loc, roles, roleOptionName)
	}
	if p.cur.Type == CREATE && p.redshiftTokIsWord(p.peekNext(), "model") {
		privLoc := nodes.Loc{Start: p.pos(), End: p.cur.End}
		p.advance() // CREATE
		privLoc.End = p.cur.End
		p.advance() // MODEL
		return p.finishRedshiftModelPrivilege(loc, privLoc, false, grantOptionFor)
	}
	if p.isRedshiftDebugStart() && !grantOptionFor && roleOptionName == "" {
		// r_REVOKE's DEBUG alternative starts at DEBUG: no GRANT OPTION FOR and
		// no ADMIN OPTION FOR prefix.
		return p.parseRedshiftDebugStmt(loc, false)
	}
	if p.isRedshiftScopedPermissionStart() {
		// r_REVOKE's scoped alternatives take `GRANT OPTION` only, without FOR,
		// and no other option prefix.
		if roleOptionName != "" || grantOptionWithFor {
			return nil, p.syntaxErrorAtCur()
		}
		privs, _, err := p.parseRedshiftGrantPrivileges()
		if err != nil {
			return nil, err
		}
		return p.finishRedshiftScopedPermission(loc, false, grantOptionFor, privs)
	}
	if roles, ok, err := p.parseRedshiftSystemPrivileges(FROM); ok || err != nil {
		if err != nil {
			return nil, err
		}
		return p.finishRevokeRole(loc, roles, roleOptionName)
	}
	privs, bareAll, err := p.parseRedshiftGrantPrivileges()
	if err != nil {
		return nil, err
	}
	switch {
	case p.cur.Type == FOR:
		if roleOptionName != "" || grantOptionWithFor {
			return nil, p.syntaxErrorAtCur()
		}
		return p.finishRedshiftScopedPermission(loc, false, grantOptionFor, privs)
	case p.cur.Type == ON:
		p.advance()
		return p.finishRevokeOnObject(loc, grantOptionFor, roleOptionName != "", privs)
	case bareAll:
		// A bare ALL with no ON clause is the system privilege r_REVOKE writes as
		// { ALL [ PRIVILEGES ] } FROM ROLE role_name [, ...].
		roles := &nodes.List{Items: []nodes.Node{&nodes.AccessPriv{PrivName: "all"}}}
		for p.cur.Type == ',' {
			p.advance()
			name, err := p.parseColId()
			if err != nil {
				return nil, err
			}
			roles.Items = append(roles.Items, &nodes.AccessPriv{PrivName: name})
		}
		return p.finishRevokeRole(loc, roles, roleOptionName)
	default:
		return p.finishRevokeRole(loc, privs, roleOptionName)
	}
}

func (p *Parser) finishGrantOnObject(loc int, privs *nodes.List) (nodes.Node, error) {
	targtype, objtype, objects, kind, err := p.parseGrantTarget(false)
	if err != nil {
		return nil, err
	}
	if err := p.checkRedshiftGrantPrivileges(kind, objtype, privs); err != nil {
		return nil, err
	}
	if _, err := p.expect(TO); err != nil {
		return nil, err
	}
	var grantees *nodes.List
	var grantOption bool
	if kind == redshiftGrantDatashare && !redshiftDatashareUsage(privs) {
		// r_GRANT writes the producer-side grantee as
		// ID ( WITH GRANT OPTION )? | ROLE role_name | GROUP group_name | PUBLIC.
		grantees, grantOption, err = p.parseRedshiftOptionableGrantees(true)
	} else {
		grantees, _, grantOption, err = p.parseRedshiftObjectGranteeList(kind, objtype, objects, privs)
	}
	if err != nil {
		return nil, err
	}
	return &nodes.GrantStmt{
		IsGrant: true, Targtype: targtype, Objtype: objtype,
		Objects: objects, Privileges: privs, Grantees: grantees,
		GrantOption: grantOption,
		Loc:         nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

func (p *Parser) finishRevokeOnObject(loc int, grantOptionFor, withRoleOption bool, privs *nodes.List) (nodes.Node, error) {
	targtype, objtype, objects, kind, err := p.parseGrantTarget(true)
	if err != nil {
		return nil, err
	}
	if err := p.checkRedshiftGrantPrivileges(kind, objtype, privs); err != nil {
		return nil, err
	}
	if _, err := p.expect(FROM); err != nil {
		return nil, err
	}
	var grantees *nodes.List
	var procedureDatashare bool
	if kind == redshiftGrantDatashare && !redshiftDatashareUsage(privs) {
		// r_REVOKE writes the producer-side grantee as
		// ID ( WITH GRANT OPTION )? | ROLE role_name | GROUP group_name | PUBLIC:
		// the option names the grant option it takes back, as GRANT OPTION FOR
		// does.
		var withGrantOption bool
		grantees, withGrantOption, err = p.parseRedshiftOptionableGrantees(true)
		grantOptionFor = grantOptionFor || withGrantOption
	} else {
		// The procedure-to-datashare form starts at REVOKE: no option prefix.
		grantees, procedureDatashare, err = p.parseRedshiftGranteeList(kind, objtype, objects, privs, false, !grantOptionFor && !withRoleOption)
	}
	if err != nil {
		return nil, err
	}
	var behavior int
	if objtype != nodes.OBJECT_TEMPLATE && kind == redshiftGrantObject && !procedureDatashare {
		// A template's revoke page, a procedure's datashare page and an external
		// table's or schema's end after the grantees: no CASCADE, no RESTRICT.
		behavior = p.parseOptDropBehavior()
	}
	return &nodes.GrantStmt{
		IsGrant: false, Targtype: targtype, Objtype: objtype,
		Objects: objects, Privileges: privs, Grantees: grantees,
		GrantOption: grantOptionFor, Behavior: nodes.DropBehavior(behavior),
		Loc: nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

// finishRedshiftModelPrivilege finishes GRANT or REVOKE of CREATE MODEL, which
// r_GRANT and r_REVOKE write without an ON clause:
//
//	GRANT CREATE MODEL TO { username [ WITH GRANT OPTION ] | ROLE ... } [, ...]
//	REVOKE [ GRANT OPTION FOR ] CREATE MODEL FROM ... [ RESTRICT ]
//
// The permission is kept in a GrantRoleStmt, the node this fork keeps a Redshift
// permission that names no catalog object in, so Catalog.Exec records nothing
// for it.
//
// https://docs.aws.amazon.com/redshift/latest/dg/r_GRANT.html
func (p *Parser) finishRedshiftModelPrivilege(loc int, privLoc nodes.Loc, isGrant, grantOptionFor bool) (nodes.Node, error) {
	stmt := &nodes.GrantRoleStmt{
		GrantedRoles: &nodes.List{Items: []nodes.Node{&nodes.AccessPriv{PrivName: "create_model", Loc: privLoc}}},
		IsGrant:      isGrant,
		Loc:          nodes.Loc{Start: loc},
	}
	if isGrant {
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
	} else if _, err := p.expect(FROM); err != nil {
		return nil, err
	}
	grantees, read, err := p.parseRedshiftOptionableGrantees(isGrant)
	if err != nil {
		return nil, err
	}
	stmt.GranteeRoles = grantees
	grantOption := grantOptionFor || read
	if !isGrant {
		stmt.Behavior = nodes.DropBehavior(p.parseOptDropBehavior())
	}
	if grantOption {
		stmt.Opt = &nodes.List{Items: []nodes.Node{makeDefElem("grant_option", &nodes.Boolean{Boolval: true})}}
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// finishRedshiftFederatedPermission finishes CONNECT and DEBUG, whose grantee is
// an AWS IAM Identity Center identity:
//
//	GRANT CONNECT [ ON WORKGROUP ] TO { <prefix>:<username> | ROLE <prefix>:<role> | PUBLIC }
//	GRANT DEBUG ON DATABASE <db> TO { <prefix>:<name> | ROLE <prefix>:<group> | ADMINS OF ACCOUNT '<id>' }
//
// Neither names a catalog object, so both are kept in a GrantRoleStmt, with the
// page's own clause in Opt.
//
// https://docs.aws.amazon.com/redshift/latest/dg/r_GRANT.html
func (p *Parser) finishRedshiftFederatedPermission(loc int, isGrant bool, privName string, privLoc nodes.Loc, opts *nodes.List, allowPublic, allowAdmins bool) (nodes.Node, error) {
	stmt := &nodes.GrantRoleStmt{
		GrantedRoles: &nodes.List{Items: []nodes.Node{&nodes.AccessPriv{PrivName: privName, Loc: privLoc}}},
		IsGrant:      isGrant,
		Opt:          opts,
		Loc:          nodes.Loc{Start: loc},
	}
	if isGrant {
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
	} else if _, err := p.expect(FROM); err != nil {
		return nil, err
	}
	grantee, err := p.parseRedshiftIdentityGrantee(allowPublic, allowAdmins)
	if err != nil {
		return nil, err
	}
	stmt.GranteeRoles = &nodes.List{Items: []nodes.Node{grantee}}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseRedshiftConnectStmt parses the CONNECT alternative of r_GRANT. ON
// WORKGROUP names no workgroup: it is recorded in Opt.
func (p *Parser) parseRedshiftConnectStmt(loc int, isGrant bool) (nodes.Node, error) {
	privLoc := nodes.Loc{Start: p.pos(), End: p.cur.End}
	p.advance() // consume CONNECT
	var opts *nodes.List
	if p.cur.Type == ON {
		p.advance()
		if !p.redshiftKeywordIs("workgroup") {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		opts = &nodes.List{Items: []nodes.Node{makeDefElem("on_workgroup", &nodes.Boolean{Boolval: true})}}
	}
	return p.finishRedshiftFederatedPermission(loc, isGrant, "connect", privLoc, opts, true, false)
}

// parseRedshiftDebugStmt parses the DEBUG alternative of r_GRANT and r_REVOKE.
func (p *Parser) parseRedshiftDebugStmt(loc int, isGrant bool) (nodes.Node, error) {
	privLoc := nodes.Loc{Start: p.pos(), End: p.cur.End}
	p.advance() // consume DEBUG
	if _, err := p.expect(ON); err != nil {
		return nil, err
	}
	if _, err := p.expect(DATABASE); err != nil {
		return nil, err
	}
	// DEBUG names one federated database, not a qualified path.
	name, err := p.parseColId()
	if err != nil {
		return nil, err
	}
	opts := &nodes.List{Items: []nodes.Node{makeDefElem("on_database", &nodes.String{Str: name})}}
	return p.finishRedshiftFederatedPermission(loc, isGrant, "debug", privLoc, opts, false, true)
}

// parseRedshiftIdentityGrantee parses one CONNECT or DEBUG grantee:
//
//	<prefix>:<name> | ROLE <prefix>:<name> | PUBLIC | ADMINS OF ACCOUNT '<id>'
//
// A quoted identity is one delimited identifier, as AWS writes it
// ("IAM:sales_auditor"), and ROLE is unreserved, so `ROLE : a` reads ROLE as the
// prefix rather than as the keyword.
func (p *Parser) parseRedshiftIdentityGrantee(allowPublic, allowAdmins bool) (*nodes.RoleSpec, error) {
	loc := p.pos()
	if allowAdmins && p.isRedshiftAdminsOfAccountStart() {
		p.advance() // consume ADMINS
		p.advance() // consume OF
		if !p.redshiftKeywordIs("account") {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		id, err := p.parseFileName()
		if err != nil {
			return nil, err
		}
		return &nodes.RoleSpec{
			Roletype: int(nodes.ROLESPEC_ADMINS_OF_ACCOUNT),
			Rolename: "admins_of_account:" + id,
			Loc:      nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	}
	if allowPublic && p.isRedshiftPublicGranteeStart() {
		p.advance()
		return &nodes.RoleSpec{Roletype: int(nodes.ROLESPEC_PUBLIC), Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
	}
	isRole := false
	if p.isRedshiftRoleIdentityStart() {
		p.advance()
		isRole = true
	}
	roletype := int(nodes.ROLESPEC_CSTRING)
	if isRole {
		roletype = int(nodes.ROLESPEC_IDC_ROLE)
	}
	// A quoted identity carries its own ":" inside the delimited identifier, and
	// must hold one <prefix>:<name> as the unquoted spelling does.
	if p.isDelimitedIdentToken(p.cur) {
		if !redshiftIdentityName(p.cur.Str) {
			return nil, p.syntaxErrorAtCur()
		}
		tok := p.advance()
		return &nodes.RoleSpec{
			Roletype: roletype,
			Rolename: tok.Str,
			Loc:      nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	}
	prefix, err := p.parseColId()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(':'); err != nil {
		return nil, err
	}
	name, err := p.parseColId()
	if err != nil {
		return nil, err
	}
	return &nodes.RoleSpec{
		Roletype: roletype,
		Rolename: prefix + ":" + name,
		Loc:      nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

// redshiftIdentityName reports whether s is one <prefix>:<name> pair: the AWS
// IAM Identity Center identity a quoted CONNECT or DEBUG grantee spells as one
// delimited identifier.
func redshiftIdentityName(s string) bool {
	prefix, name, ok := strings.Cut(s, ":")
	return ok && prefix != "" && name != "" && !strings.Contains(name, ":")
}

func (p *Parser) finishGrantRole(loc int, isGrant bool, roles *nodes.List) (nodes.Node, error) {
	if _, err := p.expect(TO); err != nil {
		return nil, err
	}
	grantees, err := p.parseGranteeList()
	if err != nil {
		return nil, err
	}
	opts, err := p.parseGrantRoleOptList()
	if err != nil {
		return nil, err
	}
	grantedBy, err := p.parseOptGrantedBy()
	if err != nil {
		return nil, err
	}
	return &nodes.GrantRoleStmt{
		GrantedRoles: roles, GranteeRoles: grantees, IsGrant: isGrant, Opt: opts,
		Grantor: grantedBy,
		Loc:     nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

func (p *Parser) finishRevokeRole(loc int, roles *nodes.List, roleOptionName string) (nodes.Node, error) {
	if _, err := p.expect(FROM); err != nil {
		return nil, err
	}
	grantees, err := p.parseGranteeList()
	if err != nil {
		return nil, err
	}
	grantedBy, err := p.parseOptGrantedBy()
	if err != nil {
		return nil, err
	}
	behavior := p.parseOptDropBehavior()
	var opts *nodes.List
	if roleOptionName != "" {
		opts = &nodes.List{Items: []nodes.Node{makeDefElem(roleOptionName, &nodes.Boolean{Boolval: false})}}
	}
	return &nodes.GrantRoleStmt{
		GrantedRoles: roles, GranteeRoles: grantees, IsGrant: false,
		Opt: opts, Grantor: grantedBy, Behavior: nodes.DropBehavior(behavior),
		Loc: nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

func (p *Parser) parseGrantTarget(isRevoke bool) (nodes.GrantTargetType, nodes.ObjectType, *nodes.List, redshiftGrantKind, error) {
	// GRANT ... ON DATASHARE datashare_name names one datashare. DATASHARE is
	// unreserved, so a grant that names no datashare after it
	// (`ON DATASHARE TO a`) stays on the plain object path below, which reads
	// DATASHARE as the object's name.
	if p.isRedshiftDatashareNameStart() {
		p.advance()
		nameStart := p.pos()
		name, err := p.parseAnyName()
		if err != nil {
			return 0, 0, nil, redshiftGrantObject, err
		}
		list := &nodes.List{Items: []nodes.Node{makeRangeVarFromAnyName(name, nodes.Loc{Start: nameStart, End: p.prev.End})}}
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_DATASHARE, list, redshiftGrantDatashare, nil
	}
	if p.redshiftKeywordIs("model") {
		p.advance()
		list, err := p.parseGrantObjectAnyNameList()
		if err != nil {
			return 0, 0, nil, redshiftGrantObject, err
		}
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_MODEL, list, redshiftGrantObject, nil
	}
	// GRANT ... ON TEMPLATE [ <database>. ] [ <schema>. ] <template> [ , ... ]
	// https://docs.aws.amazon.com/redshift/latest/dg/r_GRANT.html
	//
	// TEMPLATE is unreserved, so a grant that names no template after it
	// (`ON TEMPLATE TO a`) stays on the plain object path below, which reads
	// TEMPLATE as the object's name.
	if p.cur.Type == TEMPLATE && isColIdToken(p.peekNext()) {
		p.advance()
		list, err := p.parseRedshiftTemplateNameList(!isRevoke)
		if err != nil {
			return 0, 0, nil, redshiftGrantObject, err
		}
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_TEMPLATE, list, redshiftGrantObject, nil
	}
	switch p.cur.Type {
	case COPY:
		p.advance()
		if !p.consumeRedshiftWord("job") {
			return 0, 0, nil, redshiftGrantObject, p.syntaxErrorAtCur()
		}
		list, err := p.parseGrantObjectAnyNameList()
		if err != nil {
			return 0, 0, nil, redshiftGrantObject, err
		}
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_COPY_JOB, list, redshiftGrantObject, nil
	case EXTERNAL:
		p.advance()
		switch p.cur.Type {
		case SCHEMA:
			p.advance()
			names, err := p.parseNameList()
			if err != nil {
				return 0, 0, nil, redshiftGrantObject, err
			}
			return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_EXTERNAL_SCHEMA, makeNameListAsAnyNameList(names), redshiftGrantExternalSchema, nil
		case TABLE:
			p.advance()
			list, err := p.parseGrantObjectAnyNameList()
			if err != nil {
				return 0, 0, nil, redshiftGrantObject, err
			}
			return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_FOREIGN_TABLE, list, redshiftGrantExternalTable, nil
		default:
			// ON EXTERNAL <name> is not a form either page writes; it keeps the
			// page's object type on the plain rule lists.
			list, err := p.parseGrantObjectAnyNameList()
			if err != nil {
				return 0, 0, nil, redshiftGrantObject, err
			}
			return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_FOREIGN_TABLE, list, redshiftGrantObject, nil
		}
	case TABLE:
		p.advance()
		list, err := p.parseGrantObjectAnyNameList()
		if err != nil {
			return 0, 0, nil, redshiftGrantObject, err
		}
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_TABLE, list, redshiftGrantObject, nil
	case SEQUENCE:
		p.advance()
		list, err := p.parseGrantObjectAnyNameList()
		if err != nil {
			return 0, 0, nil, redshiftGrantObject, err
		}
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_SEQUENCE, list, redshiftGrantObject, nil
	case FUNCTION:
		p.advance()
		list, err := p.parseGrantFunctionWithArgtypesList()
		if err != nil {
			return 0, 0, nil, redshiftGrantObject, err
		}
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_FUNCTION, list, redshiftGrantObject, nil
	case PROCEDURE:
		p.advance()
		list, err := p.parseGrantFunctionWithArgtypesList()
		if err != nil {
			return 0, 0, nil, redshiftGrantObject, err
		}
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_PROCEDURE, list, redshiftGrantObject, nil
	case ROUTINE:
		p.advance()
		list, err := p.parseGrantFunctionWithArgtypesList()
		if err != nil {
			return 0, 0, nil, redshiftGrantObject, err
		}
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_ROUTINE, list, redshiftGrantObject, nil
	case DATABASE:
		p.advance()
		names, err := p.parseNameList()
		if err != nil {
			return 0, 0, nil, redshiftGrantObject, err
		}
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_DATABASE, makeNameListAsAnyNameList(names), redshiftGrantObject, nil
	case DOMAIN_P:
		p.advance()
		objects, err := p.parseAnyNameList()
		if err != nil {
			return 0, 0, nil, redshiftGrantObject, err
		}
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_DOMAIN, objects, redshiftGrantObject, nil
	case LANGUAGE:
		p.advance()
		names, err := p.parseNameList()
		if err != nil {
			return 0, 0, nil, redshiftGrantObject, err
		}
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_LANGUAGE, makeNameListAsAnyNameList(names), redshiftGrantObject, nil
	case LARGE_P:
		p.advance()
		if _, err := p.expect(OBJECT_P); err != nil {
			return 0, 0, nil, redshiftGrantObject, err
		}
		list, err := p.parseNumericOnlyList()
		if err != nil {
			return 0, 0, nil, redshiftGrantObject, err
		}
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_LARGEOBJECT, list, redshiftGrantObject, nil
	case SCHEMA:
		p.advance()
		names, err := p.parseNameList()
		if err != nil {
			return 0, 0, nil, redshiftGrantObject, err
		}
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_SCHEMA, makeNameListAsAnyNameList(names), redshiftGrantObject, nil
	case TABLESPACE:
		p.advance()
		names, err := p.parseNameList()
		if err != nil {
			return 0, 0, nil, redshiftGrantObject, err
		}
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_TABLESPACE, makeNameListAsAnyNameList(names), redshiftGrantObject, nil
	case TYPE_P:
		p.advance()
		objects, err := p.parseAnyNameList()
		if err != nil {
			return 0, 0, nil, redshiftGrantObject, err
		}
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_TYPE, objects, redshiftGrantObject, nil
	case FOREIGN:
		p.advance()
		if p.cur.Type == DATA_P {
			p.advance()
			if _, err := p.expect(WRAPPER); err != nil {
				return 0, 0, nil, redshiftGrantObject, err
			}
			names, err := p.parseNameList()
			if err != nil {
				return 0, 0, nil, redshiftGrantObject, err
			}
			return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_FDW, makeNameListAsAnyNameList(names), redshiftGrantObject, nil
		}
		if _, err := p.expect(SERVER); err != nil {
			return 0, 0, nil, redshiftGrantObject, err
		}
		names, err := p.parseNameList()
		if err != nil {
			return 0, 0, nil, redshiftGrantObject, err
		}
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_FOREIGN_SERVER, makeNameListAsAnyNameList(names), redshiftGrantObject, nil
	case ALL:
		p.advance()
		var objType nodes.ObjectType
		switch p.cur.Type {
		case TABLES:
			p.advance()
			objType = nodes.OBJECT_TABLE
		case SEQUENCES:
			p.advance()
			objType = nodes.OBJECT_SEQUENCE
		case FUNCTIONS:
			p.advance()
			objType = nodes.OBJECT_FUNCTION
		case PROCEDURES:
			p.advance()
			objType = nodes.OBJECT_PROCEDURE
		case ROUTINES:
			p.advance()
			objType = nodes.OBJECT_ROUTINE
		default:
			if p.collectMode() {
				p.addTokenCandidate(TABLES)
				p.addTokenCandidate(SEQUENCES)
				p.addTokenCandidate(FUNCTIONS)
				p.addTokenCandidate(PROCEDURES)
				p.addTokenCandidate(ROUTINES)
				return 0, 0, nil, redshiftGrantObject, errCollecting
			}
			return 0, 0, nil, redshiftGrantObject, p.syntaxErrorAtCur()
		}
		if _, err := p.expect(IN_P); err != nil {
			return 0, 0, nil, redshiftGrantObject, err
		}
		if _, err := p.expect(SCHEMA); err != nil {
			return 0, 0, nil, redshiftGrantObject, err
		}
		if p.collectMode() {
			p.addRuleCandidate("qualified_name")
			return nodes.ACL_TARGET_ALL_IN_SCHEMA, objType, nil, redshiftGrantObject, nil
		}
		names, err := p.parseNameList()
		if err != nil {
			return 0, 0, nil, redshiftGrantObject, err
		}
		return nodes.ACL_TARGET_ALL_IN_SCHEMA, objType, makeNameListAsAnyNameList(names), redshiftGrantObject, nil
	default:
		list, err := p.parseGrantObjectAnyNameList()
		if err != nil {
			return 0, 0, nil, redshiftGrantObject, err
		}
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_TABLE, list, redshiftGrantObject, nil
	}
}

func (p *Parser) parseRedshiftGrantRoleNameList() (*nodes.List, error) {
	if p.cur.Type == ROLE {
		p.advance()
	}
	name, err := p.parseColId()
	if err != nil {
		return nil, err
	}
	result := &nodes.List{Items: []nodes.Node{&nodes.AccessPriv{PrivName: name}}}
	for p.cur.Type == ',' {
		p.advance()
		if p.cur.Type == ROLE {
			p.advance()
		}
		name, err = p.parseColId()
		if err != nil {
			return nil, err
		}
		result.Items = append(result.Items, &nodes.AccessPriv{PrivName: name})
	}
	return result, nil
}

// objLocOf is the source range of a parsed object, for a diagnosis.
func objLocOf(n nodes.Node) int {
	if loc, ok := n.(interface{ GetLoc() nodes.Loc }); ok {
		return loc.GetLoc().Start
	}
	if rv, ok := n.(*nodes.RangeVar); ok {
		return rv.Loc.Start
	}
	return -1
}

// parseRedshiftObjectGranteeList parses an object grant's grantees, each with the
// WITH GRANT OPTION the page puts after a username:
//
//	{ username [ WITH GRANT OPTION ] | ROLE role_name | GROUP group_name | PUBLIC } [, ...]
//
// It also reports whether any grantee carried the option, and whether the
// grantees were the procedure-to-datashare form, which ends the statement.
func (p *Parser) parseRedshiftObjectGranteeList(kind redshiftGrantKind, objtype nodes.ObjectType, objects, privs *nodes.List) (*nodes.List, bool, bool, error) {
	if kind == redshiftGrantExternalTable && redshiftPrivilegesNameColumns(privs) && len(objects.Items) != 1 {
		// The column-level alternative names one external table.
		return nil, false, false, &ParseError{Position: objLocOf(objects.Items[1]), Message: "syntax error: an external table's column grant names one table"}
	}
	grantees, procedureDatashare, err := p.parseRedshiftGranteeList(kind, objtype, objects, privs, true, true)
	if err != nil || procedureDatashare {
		return grantees, procedureDatashare, false, err
	}
	grantOption := false
	for {
		if p.cur.Type == WITH {
			p.advance()
			if _, err := p.expect(GRANT); err != nil {
				return nil, false, false, err
			}
			if _, err := p.expect(OPTION); err != nil {
				return nil, false, false, err
			}
			// An external table's or schema's page writes one option after the
			// whole IAM role list, and it covers every role; elsewhere the option
			// follows the username it belongs to.
			all := kind == redshiftGrantExternalTable || kind == redshiftGrantExternalSchema
			if rs, ok := grantees.Items[len(grantees.Items)-1].(*nodes.RoleSpec); ok && !all {
				rs.WithGrantOption = true
			}
			if all {
				for _, item := range grantees.Items {
					if rs, ok := item.(*nodes.RoleSpec); ok {
						rs.WithGrantOption = true
					}
				}
			}
			grantOption = true
		}
		if p.cur.Type != ',' {
			return grantees, false, grantOption, nil
		}
		p.advance()
		grantee, err := p.parseGrantee()
		if err != nil {
			return nil, false, false, err
		}
		grantees.Items = append(grantees.Items, grantee)
	}
}

// redshiftGrantKind says which alternative of r_GRANT or r_REVOKE names the
// target, where the AST's Objtype alone cannot: an EXTERNAL SCHEMA is an
// OBJECT_SCHEMA and both of a DATASHARE's alternatives are OBJECT_DATASHARE.
// Each alternative lists its own privileges and grantees.
type redshiftGrantKind int

const (
	redshiftGrantObject redshiftGrantKind = iota // an object the alternative names by its own keywords
	redshiftGrantExternalTable
	redshiftGrantExternalSchema
	redshiftGrantDatashare
)

// parseRedshiftGrantPrivileges parses the privileges of a GRANT or REVOKE:
//
//	{ privilege [, ...] } | ALL [ PRIVILEGES ] [ ( column [, ...] ) ]
//
// A bare ALL [PRIVILEGES] is nil, PostgreSQL's spelling of every privilege,
// and bareAll says so: a GRANT of it with no ON clause is a system privilege.
func (p *Parser) parseRedshiftGrantPrivileges() (privs *nodes.List, bareAll bool, err error) {
	if p.cur.Type != ALL {
		privs, err = p.parsePrivilegeList()
		return privs, false, err
	}
	allLoc := p.pos()
	p.advance()
	if p.cur.Type == PRIVILEGES {
		p.advance()
	}
	cols, err := p.parseOptColumnList()
	if err != nil {
		return nil, false, err
	}
	if cols == nil {
		return nil, true, nil
	}
	return &nodes.List{Items: []nodes.Node{&nodes.AccessPriv{Cols: cols, Loc: nodes.Loc{Start: allLoc, End: p.prev.End}}}}, false, nil
}

// parseRedshiftGranteeList parses the grantees an alternative of r_GRANT or
// r_REVOKE lists for its target. Most objects take the generic list; the ones
// with a list of their own take it here.
func (p *Parser) parseRedshiftGranteeList(kind redshiftGrantKind, objtype nodes.ObjectType, objects, privs *nodes.List, isGrant, allowDatashare bool) (*nodes.List, bool, error) {
	switch kind {
	case redshiftGrantExternalTable:
		// The column-list alternative takes IAM roles only; PUBLIC belongs to the
		// one that lists privileges on the table.
		list, err := p.parseRedshiftIAMRoleGrantees(!redshiftPrivilegesNameColumns(privs))
		return list, false, err
	case redshiftGrantExternalSchema:
		list, err := p.parseRedshiftIAMRoleGrantees(false)
		return list, false, err
	case redshiftGrantDatashare:
		if redshiftDatashareUsage(privs) {
			list, err := p.parseRedshiftDatashareUsageGrantees(isGrant)
			return list, false, err
		}
		list, err := p.parseGranteeList()
		return list, false, err
	}
	if objtype == nodes.OBJECT_TEMPLATE {
		list, err := p.parseRedshiftPlainGranteeList()
		return list, false, err
	}
	if allowDatashare && p.isRedshiftDatashareNameStart() && redshiftProcedureDatashareTarget(objtype, objects, privs) {
		list, err := p.parseRedshiftDatashareGrantee()
		return list, true, err
	}
	list, err := p.parseGranteeList()
	return list, false, err
}

// parseRedshiftOptionableGrantees parses the grantees r_GRANT and r_REVOKE write
// as
//
//	ID [ WITH GRANT OPTION ] | ROLE role_name | GROUP group_name | PUBLIC
//
// one per comma, for a datashare's producer-side grantees and for CREATE MODEL.
// The pages write WITH GRANT OPTION after a username only, so a ROLE, GROUP or
// PUBLIC grantee takes none, and the statement's AST carries one option, as
// PostgreSQL's does: the option ends the list.
func (p *Parser) parseRedshiftOptionableGrantees(allowOption bool) (*nodes.List, bool, error) {
	result := &nodes.List{}
	withGrantOption := false
	for {
		loc := p.pos()
		var grantee *nodes.RoleSpec
		username := false
		switch {
		case p.isRedshiftRoleGranteeStart():
			p.advance()
			name, err := p.parseColId()
			if err != nil {
				return nil, false, err
			}
			grantee = &nodes.RoleSpec{Roletype: int(nodes.ROLESPEC_CSTRING), Rolename: name}
		case p.isRedshiftGroupGranteeStart():
			p.advance()
			name, err := p.parseColId()
			if err != nil {
				return nil, false, err
			}
			grantee = &nodes.RoleSpec{Roletype: int(nodes.ROLESPEC_CSTRING), Rolename: name}
		case p.redshiftKeywordIs("public"):
			p.advance()
			grantee = &nodes.RoleSpec{Roletype: int(nodes.ROLESPEC_PUBLIC)}
		default:
			name, err := p.parseColId()
			if err != nil {
				return nil, false, err
			}
			grantee = &nodes.RoleSpec{Roletype: int(nodes.ROLESPEC_CSTRING), Rolename: name}
			username = true
		}
		// The grantee ends at its own last token: the option is the statement's,
		// not part of the name.
		grantee.Loc = nodes.Loc{Start: loc, End: p.prev.End}
		result.Items = append(result.Items, grantee)
		if allowOption && username && p.cur.Type == WITH {
			p.advance()
			if _, err := p.expect(GRANT); err != nil {
				return nil, false, err
			}
			if _, err := p.expect(OPTION); err != nil {
				return nil, false, err
			}
			grantee.WithGrantOption = true
			withGrantOption = true
		}
		if p.cur.Type != ',' {
			return result, withGrantOption, nil
		}
		p.advance()
	}
}

// redshiftProcedureDatashareTarget reports whether r_GRANT's and r_REVOKE's
// procedure-to-datashare alternative is the one that fits: one procedure with
// its argument list, and EXECUTE. The grantee itself is
// isRedshiftDatashareNameStart's.
func redshiftProcedureDatashareTarget(objtype nodes.ObjectType, objects, privs *nodes.List) bool {
	if objtype != nodes.OBJECT_PROCEDURE || !redshiftPrivilegeIs(privs, "execute") {
		return false
	}
	if objects == nil || len(objects.Items) != 1 {
		return false
	}
	owa, ok := objects.Items[0].(*nodes.ObjectWithArgs)
	return ok && owa.Objargs != nil
}

// parseRedshiftIAMRoleGrantees parses the grantees of r_GRANT's and r_REVOKE's
// EXTERNAL TABLE and EXTERNAL SCHEMA alternatives:
//
//	{ IAM_ROLE iam_role } [, ...] | PUBLIC
//
// where only an external table takes PUBLIC. The docs write iam_role as a
// value, which their own examples spell as a quoted ARN, so both read.
func (p *Parser) parseRedshiftIAMRoleGrantees(allowPublic bool) (*nodes.List, error) {
	if allowPublic && p.redshiftKeywordIs("public") {
		loc := p.pos()
		p.advance()
		return &nodes.List{Items: []nodes.Node{&nodes.RoleSpec{Roletype: int(nodes.ROLESPEC_PUBLIC), Loc: nodes.Loc{Start: loc, End: p.prev.End}}}}, nil
	}
	result := &nodes.List{}
	for {
		loc := p.pos()
		if !p.redshiftKeywordIs("iam_role") {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		var value string
		roletype := nodes.ROLESPEC_IAM_ROLE
		switch {
		case p.cur.Type == DEFAULT:
			// IAM_ROLE DEFAULT is the keyword, not the string 'default'.
			tok := p.advance()
			value = tok.Str
			roletype = nodes.ROLESPEC_IAM_ROLE_DEFAULT
		case p.cur.Type == SCONST:
			value = p.advance().Str
		case p.isColId():
			value = p.advance().Str
		default:
			return nil, p.syntaxErrorAtCur()
		}
		result.Items = append(result.Items, &nodes.RoleSpec{
			Roletype: int(roletype),
			Rolename: "iam_role:" + value,
			Loc:      nodes.Loc{Start: loc, End: p.prev.End},
		})
		if p.cur.Type != ',' {
			return result, nil
		}
		p.advance() // consume ','
	}
}

// parseRedshiftDatashareGrantee parses TO DATASHARE datashare_name.
func (p *Parser) parseRedshiftDatashareGrantee() (*nodes.List, error) {
	loc := p.pos()
	p.advance() // consume DATASHARE
	name, err := p.parseColId()
	if err != nil {
		return nil, err
	}
	return &nodes.List{Items: []nodes.Node{&nodes.RoleSpec{
		Roletype: int(nodes.ROLESPEC_DATASHARE),
		Rolename: "datashare:" + name,
		Loc:      nodes.Loc{Start: loc, End: p.prev.End},
	}}}, nil
}

// parseRedshiftDatashareUsageGrantees parses the consumers r_GRANT and r_REVOKE
// name for USAGE ON DATASHARE:
//
//	GRANT:  NAMESPACE '<guid>' | ACCOUNT '<number>' [ VIA DATA CATALOG ]
//	REVOKE: NAMESPACE '<guid>' [, ...] | ACCOUNT '<number>' [ VIA DATA CATALOG ] [, ...]
//
// A GRANT names one consumer; a REVOKE names a list of them.
func (p *Parser) parseRedshiftDatashareUsageGrantees(isGrant bool) (*nodes.List, error) {
	loc := p.pos()
	prefix := ""
	switch {
	case p.redshiftKeywordIs("namespace"):
		prefix = "namespace"
	case p.redshiftKeywordIs("account"):
		prefix = "account"
	default:
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	item, err := p.parseRedshiftDatashareConsumer(prefix, loc)
	if err != nil {
		return nil, err
	}
	items := []nodes.Node{item}
	if isGrant {
		return &nodes.List{Items: items}, nil
	}
	for p.cur.Type == ',' {
		p.advance()
		itemLoc := p.pos()
		item, err := p.parseRedshiftDatashareConsumer(prefix, itemLoc)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return &nodes.List{Items: items}, nil
}

// parseRedshiftDatashareConsumer parses one consumer a datashare's USAGE grantee
// names: NAMESPACE '<guid>' or ACCOUNT '<number>' [ VIA DATA CATALOG ]. The
// IAM_ROLE, namespace and account prefixes keep the qualifier, which is per
// account, so a consumer can tell a Data Catalog account from a direct one.
func (p *Parser) parseRedshiftDatashareConsumer(prefix string, loc int) (*nodes.RoleSpec, error) {
	id, err := p.parseFileName()
	if err != nil {
		return nil, err
	}
	name := prefix + ":" + id
	roletype := nodes.ROLESPEC_NAMESPACE
	if prefix == "account" {
		roletype = nodes.ROLESPEC_ACCOUNT
		via, err := p.parseRedshiftViaDataCatalog()
		if err != nil {
			return nil, err
		}
		if via {
			roletype = nodes.ROLESPEC_ACCOUNT_VIA_DATA_CATALOG
			name = "account_via_data_catalog:" + id
		}
	}
	return &nodes.RoleSpec{
		Roletype: int(roletype),
		Rolename: name,
		Loc:      nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

// parseRedshiftViaDataCatalog parses the optional VIA DATA CATALOG that a
// datashare's ACCOUNT grantee takes, and reports whether it read it. The docs
// write DATA CATALOG as two words; the one-word DATA_CATALOG this fork has
// always read still reads.
func (p *Parser) parseRedshiftViaDataCatalog() (bool, error) {
	if !p.redshiftKeywordIs("via") {
		return false, nil
	}
	p.advance()
	if p.redshiftKeywordIs("data_catalog") {
		p.advance()
		return true, nil
	}
	if !p.redshiftKeywordIs("data") {
		return false, p.syntaxErrorAtCur()
	}
	p.advance()
	if !p.redshiftKeywordIs("catalog") {
		return false, p.syntaxErrorAtCur()
	}
	p.advance()
	return true, nil
}

// redshiftPrivilegesNameColumns reports whether any privilege names a column,
// which is how r_GRANT writes the external table alternative that takes an IAM
// role and nothing else.
func redshiftPrivilegesNameColumns(privs *nodes.List) bool {
	if privs == nil {
		return false
	}
	for _, item := range privs.Items {
		if priv, ok := item.(*nodes.AccessPriv); ok && priv.Cols != nil {
			return true
		}
	}
	return false
}

// redshiftDatashareUsage reports whether privs is the USAGE privilege, the one
// r_GRANT and r_REVOKE pair with consumer grantees.
func redshiftDatashareUsage(privs *nodes.List) bool {
	return redshiftPrivilegeIs(privs, "usage")
}

// redshiftPrivilegeIs reports whether privs holds exactly the privilege name.
func redshiftPrivilegeIs(privs *nodes.List, name string) bool {
	if privs == nil || len(privs.Items) != 1 {
		return false
	}
	priv, ok := privs.Items[0].(*nodes.AccessPriv)
	return ok && strings.EqualFold(priv.PrivName, name) && priv.Cols == nil
}

// checkRedshiftGrantPrivileges holds a grant or revoke on one of Redshift's own
// object types to the privileges its alternative of r_GRANT or r_REVOKE lists,
// and to the column list that alternative writes. A nil list is ALL PRIVILEGES.
// https://docs.aws.amazon.com/redshift/latest/dg/r_GRANT.html
func (p *Parser) checkRedshiftGrantPrivileges(kind redshiftGrantKind, objtype nodes.ObjectType, privs *nodes.List) error {
	if kind == redshiftGrantDatashare {
		// { ALTER | SHARE } for a producer, USAGE for a consumer: no list and no
		// column list either way.
		if privs == nil || len(privs.Items) != 1 {
			return p.redshiftPrivilegeError(nil, "a datashare takes ALTER, SHARE or USAGE")
		}
		name := redshiftPrivilegeName(privs.Items[0])
		if name != "alter" && name != "share" && name != "usage" {
			return p.redshiftPrivilegeError(privs.Items[0], name+" is not a datashare privilege")
		}
		return p.checkRedshiftPrivilegeShape(privs, false)
	}
	if objtype == nodes.OBJECT_TEMPLATE {
		return p.checkRedshiftPrivilegeNames(privs, "template", false, map[string]bool{"alter": true, "drop": true, "usage": true})
	}
	if objtype == nodes.OBJECT_COPY_JOB {
		return p.checkRedshiftPrivilegeNames(privs, "copy job", false, map[string]bool{"alter": true, "drop": true})
	}
	if objtype == nodes.OBJECT_MODEL {
		return p.checkRedshiftPrivilegeNames(privs, "model", false, map[string]bool{"execute": true})
	}
	switch kind {
	case redshiftGrantExternalTable:
		// The column-level alternative is the whole statement: { SELECT | ALL
		// [ PRIVILEGES ] } ( column_list ), one privilege and one list. The
		// privilege-list alternative names no column.
		if redshiftPrivilegesNameColumns(privs) {
			if len(privs.Items) != 1 {
				return p.redshiftPrivilegeError(privs.Items[1], "an external table's column grant is one privilege with one column list")
			}
			if name := redshiftPrivilegeName(privs.Items[0]); name != "" && name != "select" {
				return p.redshiftPrivilegeError(privs.Items[0], name+" takes no column list on an external table")
			}
			return p.checkRedshiftPrivilegeShape(privs, true)
		}
		return p.checkRedshiftPrivilegeNames(privs, "external table", false,
			map[string]bool{"select": true, "alter": true, "drop": true, "delete": true, "insert": true})
	case redshiftGrantExternalSchema:
		return p.checkRedshiftPrivilegeNames(privs, "external schema", false,
			map[string]bool{"create": true, "alter": true, "drop": true})
	}
	return nil
}

// checkRedshiftPrivilegeNames checks a privilege list against the names one
// object type takes.
func (p *Parser) checkRedshiftPrivilegeNames(privs *nodes.List, object string, columnsAllowed bool, allowed map[string]bool) error {
	if privs == nil {
		return nil
	}
	for _, item := range privs.Items {
		name := redshiftPrivilegeName(item)
		if name != "" && !allowed[name] {
			return p.redshiftPrivilegeError(item, name+" is not valid for "+object)
		}
	}
	return p.checkRedshiftPrivilegeShape(privs, columnsAllowed)
}

// checkRedshiftPrivilegeShape checks the spelling of each privilege: a
// privilege is a keyword, never a delimited identifier, and only the column
// lists r_GRANT writes may follow one. columnsAllowed admits the column list
// the external table and column-level table alternatives write.
func (p *Parser) checkRedshiftPrivilegeShape(privs *nodes.List, columnsAllowed bool) error {
	if privs == nil {
		return nil
	}
	for _, item := range privs.Items {
		priv, ok := item.(*nodes.AccessPriv)
		if !ok {
			continue
		}
		if p.isQuotedAt(priv.Loc) {
			return &ParseError{Position: priv.Loc.Start, Message: "syntax error: a privilege is a keyword, not a quoted name"}
		}
		if priv.Cols != nil {
			if !columnsAllowed {
				return &ParseError{Position: priv.Loc.Start, Message: "syntax error: a privilege here names no column"}
			}
			name := strings.ToLower(priv.PrivName)
			if name != "" && name != "select" {
				return &ParseError{Position: priv.Loc.Start, Message: "syntax error: " + priv.PrivName + " takes no column list here"}
			}
		}
	}
	return nil
}

// redshiftPrivilegeError reports a privilege the alternative does not list at
// the privilege's own token.
func (p *Parser) redshiftPrivilegeError(item nodes.Node, message string) error {
	pos := p.pos()
	if priv, ok := item.(*nodes.AccessPriv); ok && priv.Loc.Start >= 0 {
		pos = priv.Loc.Start
	}
	return &ParseError{Position: pos, Message: "syntax error: " + message}
}

// redshiftPrivilegeName is the lower-case name of a parsed privilege, empty for
// ALL PRIVILEGES.
func redshiftPrivilegeName(item nodes.Node) string {
	priv, ok := item.(*nodes.AccessPriv)
	if !ok {
		return ""
	}
	return strings.ToLower(priv.PrivName)
}

// parseRedshiftPlainGranteeList parses the grantees r_GRANT's TEMPLATE and
// ASSUMEROLE forms take: a username, ROLE role_name, GROUP group_name or
// PUBLIC. The generic list also takes an IAM role, a namespace, an account and
// an RLS policy as grantees of their own, which neither page does.
func (p *Parser) parseRedshiftPlainGranteeList() (*nodes.List, error) {
	result := &nodes.List{}
	for {
		grantee, err := p.parseRedshiftPlainGrantee()
		if err != nil {
			return nil, err
		}
		result.Items = append(result.Items, grantee)
		if p.cur.Type != ',' {
			return result, nil
		}
		p.advance() // consume ','
	}
}

func (p *Parser) parseRedshiftPlainGrantee() (nodes.Node, error) {
	loc := p.pos()
	// ROLE role_name or GROUP group_name; both words are unreserved, so either
	// names a grantee of that spelling where no name follows.
	if (p.cur.Type == ROLE || p.cur.Type == GROUP_P) && isColIdToken(p.peekNext()) {
		p.advance()
		name, err := p.parseColId()
		if err != nil {
			return nil, err
		}
		return &nodes.RoleSpec{Roletype: int(nodes.ROLESPEC_CSTRING), Rolename: name, Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
	}
	// PUBLIC is not a keyword; an unquoted word of that spelling is the
	// pseudo-role, and a quoted one is a user called public.
	if p.cur.Type == IDENT && p.cur.End-p.cur.Loc == len("public") && p.cur.Str == "public" {
		p.advance()
		return &nodes.RoleSpec{Roletype: int(nodes.ROLESPEC_PUBLIC), Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
	}
	name, err := p.parseColId()
	if err != nil {
		return nil, err
	}
	return &nodes.RoleSpec{Roletype: int(nodes.ROLESPEC_CSTRING), Rolename: name, Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
}

// isQuotedAt reports whether the token at loc is a delimited identifier.
func (p *Parser) isQuotedAt(loc nodes.Loc) bool {
	return loc.Start >= 0 && loc.Start < loc.End && loc.End <= len(p.source) && p.source[loc.Start] == '"'
}

// parseRedshiftTemplateNameList parses the template names a GRANT or REVOKE
// names (r_GRANT). GRANT writes `( <database>. )? ( <schema>. )? <template>`;
// REVOKE writes a one-part <template>. Either way the names after a comma are
// one part each, and a name has at most three parts.
func (p *Parser) parseRedshiftTemplateNameList(qualified bool) (*nodes.List, error) {
	result := &nodes.List{}
	for {
		nameStart := p.pos()
		parts := &nodes.List{}
		for {
			id, err := p.parseColId()
			if err != nil {
				return nil, err
			}
			parts.Items = append(parts.Items, &nodes.String{Str: id})
			if p.cur.Type != '.' || len(result.Items) > 0 || !qualified {
				break
			}
			if len(parts.Items) == 3 {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance() // consume '.'
		}
		result.Items = append(result.Items, makeRangeVarFromAnyName(parts, nodes.Loc{Start: nameStart, End: p.prev.End}))
		if p.cur.Type != ',' {
			return result, nil
		}
		p.advance() // consume ','
	}
}

// parseGrantObjectAnyNameList parses the qualified_name_list a GRANT or
// REVOKE names its tables or sequences by. A comma promises a name.
func (p *Parser) parseGrantObjectAnyNameList() (*nodes.List, error) {
	nameStart := p.pos()
	name, err := p.parseAnyName()
	if err != nil {
		return nil, err
	}
	result := &nodes.List{Items: []nodes.Node{makeRangeVarFromAnyName(name, nodes.Loc{Start: nameStart, End: p.prev.End})}}
	for p.cur.Type == ',' {
		p.advance()
		nameStart = p.pos()
		name, err = p.parseAnyName()
		if err != nil {
			return nil, err
		}
		result.Items = append(result.Items, makeRangeVarFromAnyName(name, nodes.Loc{Start: nameStart, End: p.prev.End}))
	}
	return result, nil
}

func (p *Parser) parseRedshiftSystemPrivileges(stopToken int) (*nodes.List, bool, error) {
	if !p.redshiftSystemPrivilegeStart() {
		return nil, false, nil
	}
	if next := p.peekNext(); next.Type == ON || next.Type == ',' {
		return nil, false, nil
	}

	loc := p.pos()
	var parts []string
	for p.cur.Type != 0 && p.cur.Type != ';' && p.cur.Type != stopToken {
		if p.cur.Type == ON {
			return nil, false, nil
		}
		if p.cur.Type == ',' {
			return nil, false, p.syntaxErrorAtCur()
		}
		parts = append(parts, strings.ToLower(p.cur.Str))
		p.advance()
	}
	if p.cur.Type != stopToken {
		return nil, false, p.syntaxErrorAtCur()
	}
	if len(parts) == 0 {
		return nil, false, nil
	}
	name := strings.Join(parts, "_")
	return &nodes.List{Items: []nodes.Node{&nodes.AccessPriv{PrivName: name, Loc: nodes.Loc{Start: loc, End: p.prev.End}}}}, true, nil
}

func (p *Parser) redshiftSystemPrivilegeStart() bool {
	switch p.cur.Type {
	case CREATE, DROP, ALTER, TRUNCATE, VACUUM, ANALYZE, ANALYSE:
		return true
	}
	return p.redshiftWordEqual("access") ||
		p.redshiftWordEqual("cancel") ||
		p.redshiftWordEqual("ignore") ||
		p.redshiftWordEqual("explain")
}

// finishRedshiftScopedPermission finishes a scoped permission, whose privileges
// the caller has already parsed:
//
//	{ privilege [, ...] | ALL [ PRIVILEGES ] } FOR { SCHEMAS | TABLES | ... } IN ...
//
// https://docs.aws.amazon.com/redshift/latest/dg/r_GRANT.html
func (p *Parser) finishRedshiftScopedPermission(loc int, isGrant bool, grantOption bool, privs *nodes.List) (nodes.Node, error) {
	if _, err := p.expect(FOR); err != nil {
		return nil, err
	}
	objtype, objects, err := p.parseRedshiftScopedPermissionTarget()
	if err != nil {
		return nil, err
	}
	if err := p.checkRedshiftScopedPrivileges(objtype, privs); err != nil {
		return nil, err
	}
	if isGrant {
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
	} else if _, err := p.expect(FROM); err != nil {
		return nil, err
	}
	// r_GRANT writes username ( WITH GRANT OPTION )? per grantee; a REVOKE's
	// scoped grantees take none. The statement's flag is their union.
	grantees, read, err := p.parseRedshiftScopedGrantees(isGrant)
	if err != nil {
		return nil, err
	}
	grantOption = grantOption || read
	// A scoped permission ends after its grantees: no CASCADE, no RESTRICT.
	return &nodes.GrantStmt{
		IsGrant:     isGrant,
		Targtype:    nodes.ACL_TARGET_SCOPED,
		Objtype:     objtype,
		Objects:     objects,
		Privileges:  privs,
		Grantees:    grantees,
		GrantOption: grantOption,
		Loc:         nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

// parseRedshiftScopedGrantees parses a scoped permission's grantees, which
// r_GRANT and r_REVOKE write as
// { username [ WITH GRANT OPTION ] | ROLE role_name } [, ...].
func (p *Parser) parseRedshiftScopedGrantees(allowOption bool) (*nodes.List, bool, error) {
	result := &nodes.List{}
	withGrantOption := false
	for {
		loc := p.pos()
		// A scoped grantee is a username or a role, never the PUBLIC pseudo-role.
		if p.redshiftKeywordIs("public") {
			return nil, false, p.syntaxErrorAtCur()
		}
		// ROLE is unreserved, so ROLE with no name after it is the username
		// "role".
		isRole := p.isRedshiftRoleGranteeStart()
		if isRole {
			p.advance()
		}
		name, err := p.parseColId()
		if err != nil {
			return nil, false, err
		}
		grantee := &nodes.RoleSpec{
			Roletype: int(nodes.ROLESPEC_CSTRING),
			Rolename: name,
			Loc:      nodes.Loc{Start: loc, End: p.prev.End},
		}
		result.Items = append(result.Items, grantee)
		if allowOption && !isRole && p.cur.Type == WITH {
			p.advance()
			if _, err := p.expect(GRANT); err != nil {
				return nil, false, err
			}
			if _, err := p.expect(OPTION); err != nil {
				return nil, false, err
			}
			grantee.WithGrantOption = true
			withGrantOption = true
		}
		if p.cur.Type != ',' {
			return result, withGrantOption, nil
		}
		p.advance()
	}
}

// checkRedshiftScopedPrivileges holds a scoped permission to the privileges its
// alternative of r_GRANT or r_REVOKE lists.
func (p *Parser) checkRedshiftScopedPrivileges(objtype nodes.ObjectType, privs *nodes.List) error {
	switch objtype {
	case nodes.OBJECT_SCHEMA:
		return p.checkRedshiftPrivilegeNames(privs, "scoped schema", false,
			map[string]bool{"create": true, "usage": true, "alter": true, "drop": true})
	case nodes.OBJECT_TABLE:
		return p.checkRedshiftPrivilegeNames(privs, "scoped table", false,
			map[string]bool{"select": true, "insert": true, "update": true, "delete": true,
				"drop": true, "alter": true, "truncate": true, "references": true})
	case nodes.OBJECT_FUNCTION, nodes.OBJECT_PROCEDURE:
		return p.checkRedshiftPrivilegeNames(privs, "scoped function", false,
			map[string]bool{"execute": true})
	case nodes.OBJECT_LANGUAGE:
		return p.checkRedshiftPrivilegeNames(privs, "scoped language", false,
			map[string]bool{"usage": true})
	case nodes.OBJECT_COPY_JOB:
		return p.checkRedshiftPrivilegeNames(privs, "scoped copy job", false,
			map[string]bool{"create": true, "alter": true, "drop": true})
	case nodes.OBJECT_TEMPLATE:
		return p.checkRedshiftPrivilegeNames(privs, "scoped template", false,
			map[string]bool{"alter": true, "drop": true, "usage": true})
	}
	return nil
}

func (p *Parser) parseRedshiftScopedPermissionTarget() (nodes.ObjectType, *nodes.List, error) {
	objtype := nodes.OBJECT_TABLE
	switch {
	case p.cur.Type == SCHEMAS:
		p.advance()
		objtype = nodes.OBJECT_SCHEMA
	case p.cur.Type == TABLES:
		p.advance()
		objtype = nodes.OBJECT_TABLE
	case p.cur.Type == FUNCTIONS:
		p.advance()
		objtype = nodes.OBJECT_FUNCTION
	case p.cur.Type == PROCEDURES:
		p.advance()
		objtype = nodes.OBJECT_PROCEDURE
	case p.redshiftKeywordIs("languages"):
		p.advance()
		objtype = nodes.OBJECT_LANGUAGE
	case p.redshiftKeywordIs("templates"):
		p.advance()
		objtype = nodes.OBJECT_TEMPLATE
	case p.cur.Type == COPY:
		p.advance()
		if !p.redshiftKeywordIs("jobs") {
			return 0, nil, p.syntaxErrorAtCur()
		}
		p.advance()
		objtype = nodes.OBJECT_COPY_JOB
	default:
		return 0, nil, p.syntaxErrorAtCur()
	}
	if _, err := p.expect(IN_P); err != nil {
		return 0, nil, err
	}
	var items []nodes.Node
	schemaScope := false
	switch p.cur.Type {
	case SCHEMA:
		p.advance()
		name, err := p.parseColId()
		if err != nil {
			return 0, nil, err
		}
		items = append(items, makeDefElem("schema", &nodes.String{Str: name}))
		schemaScope = true
	case DATABASE:
		p.advance()
		name, err := p.parseColId()
		if err != nil {
			return 0, nil, err
		}
		items = append(items, makeDefElem("database", &nodes.String{Str: name}))
	default:
		return 0, nil, p.syntaxErrorAtCur()
	}
	// { SCHEMA name [ DATABASE name ] | DATABASE name }: a database scope is
	// already the whole scope.
	if schemaScope && p.cur.Type == DATABASE {
		p.advance()
		name, err := p.parseColId()
		if err != nil {
			return 0, nil, err
		}
		items = append(items, makeDefElem("database", &nodes.String{Str: name}))
	}
	return objtype, &nodes.List{Items: items}, nil
}

func (p *Parser) parseRedshiftAssumeRoleStmt(loc int, isGrant bool) (nodes.Node, error) {
	p.advance() // consume ASSUMEROLE
	// ASSUMEROLE's roles are the ON clause's objects:
	//   GRANT ASSUMEROLE ON { '<role>' [, ...] | DEFAULT | ALL } TO ...
	// https://docs.aws.amazon.com/redshift/latest/dg/r_GRANT.html
	// ON is read where the page writes it; the shorter spelling without it is
	// read too, of the forms this fork already answered. A role named 'on' is a
	// string, so it never stands for the keyword.
	p.acceptRedshiftKeyword("on")
	roles, err := p.parseRedshiftAssumeRoleList()
	if err != nil {
		return nil, err
	}
	if isGrant {
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
	} else if _, err := p.expect(FROM); err != nil {
		return nil, err
	}
	grantees, err := p.parseRedshiftPlainGranteeList()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(FOR); err != nil {
		return nil, err
	}
	var opts *nodes.List
	if isGrant {
		// GRANT's FOR clause takes a list of actions; REVOKE's takes one.
		opts, err = p.parseRedshiftAssumeRoleActionList()
		if err != nil {
			return nil, err
		}
	} else {
		action, err := p.parseRedshiftAssumeRoleAction()
		if err != nil {
			return nil, err
		}
		opts = &nodes.List{Items: []nodes.Node{makeDefElem("for", &nodes.String{Str: action})}}
	}
	return &nodes.GrantRoleStmt{
		GrantedRoles: roles,
		GranteeRoles: grantees,
		IsGrant:      isGrant,
		Opt:          opts,
		Loc:          nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

func (p *Parser) parseRedshiftAssumeRoleList() (*nodes.List, error) {
	role, err := p.parseRedshiftAssumeRoleItem()
	if err != nil {
		return nil, err
	}
	items := []nodes.Node{role}
	// DEFAULT and ALL stand alone; only a string list takes commas.
	if p.cur.Type == ',' && (role.PrivName == "default" || role.PrivName == "all") {
		return nil, p.syntaxErrorAtCur()
	}
	for p.cur.Type == ',' {
		p.advance()
		role, err = p.parseRedshiftAssumeRoleItem()
		if err != nil {
			return nil, err
		}
		items = append(items, role)
	}
	return &nodes.List{Items: items}, nil
}

func (p *Parser) parseRedshiftAssumeRoleItem() (*nodes.AccessPriv, error) {
	loc := p.pos()
	switch p.cur.Type {
	case SCONST:
		tok := p.advance()
		return &nodes.AccessPriv{PrivName: tok.Str, Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
	case DEFAULT:
		p.advance()
		return &nodes.AccessPriv{PrivName: "default", Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
	case ALL:
		p.advance()
		return &nodes.AccessPriv{PrivName: "all", Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
	default:
		return nil, p.syntaxErrorAtCur()
	}
}

// parseRedshiftAssumeRoleActionList parses ASSUMEROLE's action list:
//
//	{ ALL | COPY | UNLOAD | EXTERNAL FUNCTION | CREATE MODEL } [ , ... ]
func (p *Parser) parseRedshiftAssumeRoleActionList() (*nodes.List, error) {
	items := &nodes.List{}
	for {
		action, err := p.parseRedshiftAssumeRoleAction()
		if err != nil {
			return nil, err
		}
		items.Items = append(items.Items, makeDefElem("for", &nodes.String{Str: action}))
		if p.cur.Type != ',' {
			return items, nil
		}
		p.advance() // consume ','
	}
}

// parseRedshiftAssumeRoleAction reads one action of ASSUMEROLE's FOR clause:
// one of ALL, COPY, UNLOAD, EXTERNAL FUNCTION and CREATE MODEL, each a whole
// alternative. A quoted spelling is a name, never one of them.
func (p *Parser) parseRedshiftAssumeRoleAction() (string, error) {
	switch {
	case p.redshiftKeywordIs("all"):
		p.advance()
		return "all", nil
	case p.redshiftKeywordIs("copy"):
		p.advance()
		return "copy", nil
	case p.redshiftKeywordIs("unload"):
		p.advance()
		return "unload", nil
	case p.redshiftKeywordIs("external"):
		p.advance()
		if !p.acceptRedshiftKeyword("function") {
			return "", p.syntaxErrorAtCur()
		}
		return "external_function", nil
	case p.redshiftKeywordIs("create"):
		p.advance()
		if !p.acceptRedshiftKeyword("model") {
			return "", p.syntaxErrorAtCur()
		}
		return "create_model", nil
	default:
		return "", p.syntaxErrorAtCur()
	}
}

// redshiftKeywordIs reports whether cur is the unquoted word word: a keyword
// token, or an unquoted identifier spelled exactly word. A delimited identifier
// ("all") is a name, never the grammar's keyword.
func (p *Parser) redshiftKeywordIs(word string) bool {
	switch {
	case p.cur.Type == IDENT:
		return p.cur.End-p.cur.Loc == len(word) && p.redshiftWordEqual(word)
	case p.cur.Type >= ABORT_P:
		return p.redshiftWordEqual(word)
	default:
		return false
	}
}

// acceptRedshiftKeyword consumes the unquoted word word, and reports whether it
// did.
func (p *Parser) acceptRedshiftKeyword(word string) bool {
	if !p.redshiftKeywordIs(word) {
		return false
	}
	p.advance()
	return true
}

func (p *Parser) parseNumericOnlyList() (*nodes.List, error) {
	result := &nodes.List{}
	for {
		val, err := p.parseNumericOnlyRequired()
		if err != nil {
			return nil, err
		}
		result.Items = append(result.Items, val)
		if p.cur.Type != ',' {
			return result, nil
		}
		p.advance()
	}
}

func (p *Parser) parseGrantFunctionWithArgtypesList() (*nodes.List, error) {
	fn, err := p.parseGrantFunctionWithArgtypes()
	if err != nil {
		return nil, err
	}
	if fn == nil {
		return nil, nil
	}
	result := &nodes.List{Items: []nodes.Node{fn}}
	for p.cur.Type == ',' {
		p.advance()
		fn, err = p.parseGrantFunctionWithArgtypes()
		if err != nil {
			return nil, err
		}
		if fn == nil {
			return nil, p.syntaxErrorAtCur()
		}
		result.Items = append(result.Items, fn)
	}
	return result, nil
}

func (p *Parser) parseGrantFunctionWithArgtypes() (*nodes.ObjectWithArgs, error) {
	funcName, err := p.parseFuncName()
	if err != nil {
		return nil, err
	}
	owa := &nodes.ObjectWithArgs{Objname: funcName}
	if p.cur.Type == '(' {
		p.advance()
		if p.cur.Type == ')' {
			p.advance()
			owa.Objargs = &nodes.List{}
		} else {
			args, err := p.parseGrantFuncArgsList()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
			owa.Objargs = args
		}
	} else {
		owa.ArgsUnspecified = true
	}
	return owa, nil
}

func (p *Parser) parseGrantFuncArgsList() (*nodes.List, error) {
	arg, err := p.parseGrantFuncArg()
	if err != nil {
		return nil, err
	}
	if arg == nil {
		return nil, nil
	}
	result := &nodes.List{Items: []nodes.Node{arg}}
	for p.cur.Type == ',' {
		p.advance()
		arg, err = p.parseGrantFuncArg()
		if err != nil {
			return nil, err
		}
		if arg == nil {
			if p.collectMode() {
				break
			}
			return nil, p.syntaxErrorAtCur()
		}
		result.Items = append(result.Items, arg)
	}
	return result, nil
}

func (p *Parser) parseGrantFuncArg() (nodes.Node, error) {
	arg, err := p.parseFuncArg()
	if err != nil {
		return nil, err
	}
	return arg.ArgType, nil
}

func (p *Parser) parsePrivileges() (*nodes.List, error) {
	if p.cur.Type == ALL {
		loc := p.pos()
		p.advance()
		if p.cur.Type == PRIVILEGES {
			p.advance()
		}
		cols, err := p.parseOptColumnList()
		if err != nil {
			return nil, err
		}
		if cols != nil {
			return &nodes.List{Items: []nodes.Node{&nodes.AccessPriv{Cols: cols, Loc: nodes.Loc{Start: loc, End: p.prev.End}}}}, nil
		}
		return nil, nil
	}
	return p.parsePrivilegeList()
}

func (p *Parser) parsePrivilegeList() (*nodes.List, error) {
	priv, err := p.parsePrivilege()
	if err != nil {
		return nil, err
	}
	if priv == nil {
		return nil, nil
	}
	result := &nodes.List{Items: []nodes.Node{priv}}
	columnListOnly := false
	for p.cur.Type == ',' {
		if p.peekNext().Type == '(' && redshiftColumnListPrivilege(priv) && (len(result.Items) == 1 || columnListOnly) {
			// r_GRANT's column-level alternative repeats the column list after
			// the privilege:
			//   { SELECT | UPDATE } ( column [, ...] ) [ , ( column [, ...] ) ... ]
			// and the alternative ends there: no other privilege follows it.
			p.advance() // consume ','
			listStart := p.pos()
			p.advance() // consume '('
			cols, err := p.parseColumnList()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
			priv = &nodes.AccessPriv{PrivName: priv.PrivName, Cols: cols, Loc: nodes.Loc{Start: listStart, End: p.prev.End}}
			result.Items = append(result.Items, priv)
			columnListOnly = true
			continue
		}
		if columnListOnly {
			break
		}
		p.advance()
		priv, err = p.parsePrivilege()
		if err != nil {
			return nil, err
		}
		// privilege_list: privilege | privilege_list ',' privilege
		if priv == nil {
			if p.collectMode() {
				break
			}
			return nil, p.syntaxErrorAtCur()
		}
		result.Items = append(result.Items, priv)
	}
	return result, nil
}

// redshiftColumnListPrivilege reports whether priv is one of the privileges
// r_GRANT's column-level alternative repeats a column list for: SELECT and
// UPDATE, which the alternative writes with a list.
func redshiftColumnListPrivilege(priv *nodes.AccessPriv) bool {
	if priv == nil || priv.Cols == nil {
		return false
	}
	switch strings.ToLower(priv.PrivName) {
	case "select", "update":
		return true
	}
	return false
}

func (p *Parser) parsePrivilege() (*nodes.AccessPriv, error) {
	loc := p.pos()
	var privName string
	switch p.cur.Type {
	case SELECT:
		privName = "select"
		p.advance()
	case INSERT:
		privName = "insert"
		p.advance()
	case UPDATE:
		privName = "update"
		p.advance()
	case DELETE_P:
		privName = "delete"
		p.advance()
	case TRUNCATE:
		privName = "truncate"
		p.advance()
	case REFERENCES:
		privName = "references"
		p.advance()
	case TRIGGER:
		privName = "trigger"
		p.advance()
	case CREATE:
		privName = "create"
		p.advance()
	case TEMPORARY:
		privName = "temporary"
		p.advance()
	case TEMP:
		privName = "temp"
		p.advance()
	case EXECUTE:
		privName = "execute"
		p.advance()
	default:
		name, err := p.parseColId()
		if err != nil {
			return nil, err
		}
		privName = name
	}
	cols, err := p.parseOptColumnList()
	if err != nil {
		return nil, err
	}
	return &nodes.AccessPriv{PrivName: privName, Cols: cols,
		Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
}

func (p *Parser) parseGranteeList() (*nodes.List, error) {
	g, err := p.parseGrantee()
	if err != nil {
		return nil, err
	}
	result := &nodes.List{Items: []nodes.Node{g}}
	for p.cur.Type == ',' {
		p.advance()
		g, err = p.parseGrantee()
		if err != nil {
			return nil, err
		}
		result.Items = append(result.Items, g)
	}
	return result, nil
}

// granteeEnds reports whether a token ends a grantee: after it, a ROLE
// before is the grantee's name, not ROLE role_name. ROLE is not among
// Redshift's reserved words, so a user may be named role
// (https://docs.aws.amazon.com/redshift/latest/dg/r_pg_keywords.html), and
// GRANT and REVOKE take TO { username | ROLE role_name | ... }
// (https://docs.aws.amazon.com/redshift/latest/dg/r_GRANT.html). CASCADE
// and RESTRICT are unreserved too, so ROLE cascade names a role.
func granteeEnds(t int) bool {
	switch t {
	case 0, ';', ',', WITH:
		return true
	}
	return false
}

func (p *Parser) parseGrantee() (*nodes.RoleSpec, error) {
	if p.redshiftWordEqual("rls") {
		loc := p.pos()
		p.advance()
		if p.cur.Type == POLICY {
			p.advance()
		}
		name, err := p.parseColId()
		if err != nil {
			return nil, err
		}
		return &nodes.RoleSpec{
			Roletype: int(nodes.ROLESPEC_RLS_POLICY),
			Rolename: "rls_policy:" + name,
			Loc:      nodes.Loc{Start: loc, End: p.pos()},
		}, nil
	}
	if p.redshiftWordEqual("namespace") || p.redshiftWordEqual("account") {
		loc := p.pos()
		prefix := strings.ToLower(p.cur.Str)
		p.advance()
		if p.cur.Type != SCONST {
			return p.parseRoleSpec()
		}
		tok := p.advance()
		roletype := nodes.ROLESPEC_NAMESPACE
		name := prefix + ":" + tok.Str
		if prefix == "account" {
			roletype = nodes.ROLESPEC_ACCOUNT
			if p.consumeRedshiftWord("via") {
				p.consumeRedshiftWord("data_catalog")
			}
		}
		return &nodes.RoleSpec{
			Roletype: int(roletype),
			Rolename: name,
			Loc:      nodes.Loc{Start: loc, End: p.pos()},
		}, nil
	}
	if p.consumeRedshiftWord("iam_role") {
		loc := p.prev.Loc
		value := ""
		if p.cur.Type == SCONST || p.cur.Type == DEFAULT {
			tok := p.advance()
			value = tok.Str
			if tok.Type == DEFAULT {
				value = "default"
			}
		}
		return &nodes.RoleSpec{
			Roletype: int(nodes.ROLESPEC_IAM_ROLE),
			Rolename: "iam_role:" + value,
			Loc:      nodes.Loc{Start: loc, End: p.pos()},
		}, nil
	}
	if p.cur.Type == GROUP_P || p.cur.Type == ROLE && !granteeEnds(p.peekNext().Type) {
		p.advance()
	}
	loc := p.pos()
	// PUBLIC is not a reserved keyword; it appears as IDENT
	if p.isColId() && p.cur.Str == "public" {
		p.advance()
		return &nodes.RoleSpec{Roletype: int(nodes.ROLESPEC_PUBLIC), Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
	}
	return p.parseRoleSpec()
}

// parseOptGrantGrantOption parses opt_grant_grant_option. WITH starts
// nothing else after the grantees, so it is read as WITH GRANT OPTION.
//
//	opt_grant_grant_option: WITH GRANT OPTION | /* EMPTY */
func (p *Parser) parseOptGrantGrantOption() (bool, error) {
	if p.cur.Type != WITH {
		return false, nil
	}
	p.advance()
	if _, err := p.expect(GRANT); err != nil {
		return false, err
	}
	if _, err := p.expect(OPTION); err != nil {
		return false, err
	}
	return true, nil
}

func (p *Parser) parseOptGrantedBy() (*nodes.RoleSpec, error) {
	if p.cur.Type == GRANTED {
		p.advance()
		if _, err := p.expect(BY); err != nil {
			return nil, err
		}
		return p.parseRoleSpec()
	}
	return nil, nil
}

func (p *Parser) parseGrantRoleOptList() (*nodes.List, error) {
	// PG grammar: opt_granted_by's sibling here is
	//   WITH grant_role_opt_list
	//   where grant_role_opt_list is a comma-separated list of grant_role_opt.
	// A single WITH introduces the whole list; entries are comma-separated.
	// Once WITH (or a ',') has been consumed, a valid grant_role_opt is
	// required — silently accepting a trailing comma or an unrecognized
	// option would make the parser permissive where PG rejects.
	if p.cur.Type != WITH {
		return nil, nil
	}
	p.advance() // consume WITH
	first := p.parseGrantRoleOptValue()
	if first == nil {
		return nil, p.syntaxErrorAtCur()
	}
	items := []nodes.Node{first}
	for p.cur.Type == ',' {
		p.advance()
		next := p.parseGrantRoleOptValue()
		if next == nil {
			return nil, p.syntaxErrorAtCur()
		}
		items = append(items, next)
	}
	return &nodes.List{Items: items}, nil
}

func (p *Parser) parseGrantRoleOptValue() *nodes.DefElem {
	var name string
	switch p.cur.Type {
	case ADMIN:
		name = "admin"
		p.advance()
	case INHERIT:
		name = "inherit"
		p.advance()
	case SET:
		name = "set"
		p.advance()
	default:
		return nil
	}
	var val nodes.Node
	switch p.cur.Type {
	case OPTION:
		p.advance()
		val = &nodes.Boolean{Boolval: true}
	case TRUE_P:
		p.advance()
		val = &nodes.Boolean{Boolval: true}
	case FALSE_P:
		p.advance()
		val = &nodes.Boolean{Boolval: false}
	default:
		val = &nodes.Boolean{Boolval: true}
	}
	return makeDefElem(name, val)
}

func (p *Parser) parseCreateRoleStmt() (nodes.Node, error) {
	loc := p.prev.Loc // position of CREATE keyword
	p.advance()       // consume ROLE
	name, err := p.parseRoleId()
	if err != nil {
		return nil, err
	}
	p.parseGrantOptWith()
	options, err := p.parseOptRoleList(true)
	if err != nil {
		return nil, err
	}
	return &nodes.CreateRoleStmt{StmtType: nodes.ROLESTMT_ROLE, Role: name, Options: options,
		Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
}

func (p *Parser) parseCreateUserStmt() (nodes.Node, error) {
	loc := p.prev.Loc // position of CREATE keyword
	p.advance()       // consume USER
	name, err := p.parseRedshiftUserId()
	if err != nil {
		return nil, err
	}
	p.parseGrantOptWith()
	options, err := p.parseOptRoleList(true)
	if err != nil {
		return nil, err
	}
	return &nodes.CreateRoleStmt{StmtType: nodes.ROLESTMT_USER, Role: name, Options: options,
		Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
}

func (p *Parser) parseCreateGroupStmt() (nodes.Node, error) {
	loc := p.prev.Loc // position of CREATE keyword
	p.advance()       // consume GROUP
	name, err := p.parseRoleId()
	if err != nil {
		return nil, err
	}
	// r_CREATE_GROUP writes the WITH clause once per comma-separated group of
	// users, and each WITH takes its USER list as optional:
	//	CREATE GROUP name [ WITH [ USER name [, ...] ] [, WITH [ USER ... ] ] ]
	//
	// The commas belong to that WITH block, so a statement with no WITH has
	// none of them.
	with := p.cur.Type == WITH
	p.parseGrantOptWith()
	savedInCreateGroup := p.inCreateGroup
	p.inCreateGroup = with
	options, err := p.parseOptRoleList(true)
	p.inCreateGroup = savedInCreateGroup
	if err != nil {
		return nil, err
	}
	for with && p.cur.Type == ',' {
		p.advance()
		if p.collectMode() {
			// The clause repeats after the comma, and WITH comes first.
			p.addTokenCandidate(WITH)
			return options, nil
		}
		if !p.acceptRedshiftKeyword("with") {
			return nil, p.syntaxErrorAtCur()
		}
		more, err := p.parseOptRoleList(true)
		if err != nil {
			return nil, err
		}
		options = appendNodeLists(options, more)
	}
	return &nodes.CreateRoleStmt{StmtType: nodes.ROLESTMT_GROUP, Role: name, Options: options,
		Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
}

func (p *Parser) parseRoleId() (string, error) {
	return p.parseNonReservedWord()
}

func (p *Parser) parseRedshiftUserId() (string, error) {
	name, err := p.parseRoleId()
	if err != nil {
		return "", err
	}
	for p.cur.Type == ':' {
		p.advance()
		part, err := p.parseColId()
		if err != nil {
			return "", err
		}
		name += ":" + part
	}
	return name, nil
}

func (p *Parser) parseRedshiftUserSpec() (*nodes.RoleSpec, error) {
	loc := p.pos()
	name, err := p.parseRedshiftUserId()
	if err != nil {
		return nil, err
	}
	return &nodes.RoleSpec{
		Roletype: int(nodes.ROLESPEC_CSTRING),
		Rolename: name,
		Loc:      nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

func (p *Parser) parseGrantOptWith() {
	if p.cur.Type == WITH {
		p.advance()
	}
}

func (p *Parser) parseOptRoleList(isCreate bool) (*nodes.List, error) {
	var items []nodes.Node
	for {
		var elem *nodes.DefElem
		if isCreate {
			var err error
			elem, err = p.parseCreateOptRoleElem()
			if err != nil {
				return nil, err
			}
		} else {
			elem = p.parseAlterOptRoleElem()
		}
		if elem == nil {
			break
		}
		items = append(items, elem)
	}
	if len(items) == 0 {
		return nil, nil
	}
	return &nodes.List{Items: items}, nil
}

func (p *Parser) parseAlterOptRoleElem() *nodes.DefElem {
	switch p.cur.Type {
	case PASSWORD:
		p.advance()
		if p.cur.Type == DISABLE_P {
			p.advance()
			return makeDefElem("password_disabled", &nodes.Boolean{Boolval: true})
		}
		if p.cur.Type == NULL_P {
			p.advance()
			return makeDefElem("password", nil)
		}
		pw := p.cur.Str
		p.advance()
		return makeDefElem("password", &nodes.String{Str: pw})
	case ENCRYPTED:
		p.advance()
		p.expect(PASSWORD)
		pw := p.cur.Str
		p.advance()
		return makeDefElem("password", &nodes.String{Str: pw})
	case UNENCRYPTED:
		p.advance()
		p.expect(PASSWORD)
		pw := p.cur.Str
		p.advance()
		return makeDefElem("password", &nodes.String{Str: pw})
	case INHERIT:
		p.advance()
		return makeDefElem("inherit", &nodes.Boolean{Boolval: true})
	case CONNECTION:
		p.advance()
		p.expect(LIMIT)
		if p.redshiftWordEqual("unlimited") {
			p.advance()
			return makeDefElem("connectionlimit", &nodes.String{Str: "unlimited"})
		}
		val := p.parseSignedIconst()
		return makeDefElem("connectionlimit", &nodes.Integer{Ival: val})
	case VALID:
		p.advance()
		p.expect(UNTIL)
		until := p.cur.Str
		p.advance()
		return makeDefElem("validUntil", &nodes.String{Str: until})
	case SESSION:
		p.advance()
		if !p.redshiftWordEqual("timeout") {
			return nil
		}
		p.advance()
		timeout := p.parseSignedIconst()
		return makeDefElem("session_timeout", &nodes.Integer{Ival: timeout})
	case RESET:
		p.advance()
		if p.cur.Type == SESSION {
			p.advance()
			if !p.redshiftWordEqual("timeout") {
				return nil
			}
			p.advance()
			return makeDefElem("reset_session_timeout", &nodes.Boolean{Boolval: true})
		}
		resetStmt, _ := p.parseVariableResetStmt()
		return makeDefElem("reset", resetStmt)
	case SET:
		p.advance()
		setStmt, _ := p.parseVariableSetStmt()
		return makeDefElem("set", setStmt)
	default:
		return p.parseRoleOptionIdent()
	}
}

func (p *Parser) parseRoleOptionIdent() *nodes.DefElem {
	if !p.isColId() {
		return nil
	}
	ident := p.cur.Str
	switch ident {
	case "superuser":
		p.advance()
		return makeDefElem("superuser", &nodes.Boolean{Boolval: true})
	case "nosuperuser":
		p.advance()
		return makeDefElem("superuser", &nodes.Boolean{Boolval: false})
	case "createrole":
		p.advance()
		return makeDefElem("createrole", &nodes.Boolean{Boolval: true})
	case "nocreaterole":
		p.advance()
		return makeDefElem("createrole", &nodes.Boolean{Boolval: false})
	case "createdb":
		p.advance()
		return makeDefElem("createdb", &nodes.Boolean{Boolval: true})
	case "nocreatedb":
		p.advance()
		return makeDefElem("createdb", &nodes.Boolean{Boolval: false})
	case "createuser":
		p.advance()
		return makeDefElem("superuser", &nodes.Boolean{Boolval: true})
	case "nocreateuser":
		p.advance()
		return makeDefElem("superuser", &nodes.Boolean{Boolval: false})
	case "login":
		p.advance()
		return makeDefElem("canlogin", &nodes.Boolean{Boolval: true})
	case "nologin":
		p.advance()
		return makeDefElem("canlogin", &nodes.Boolean{Boolval: false})
	case "replication":
		p.advance()
		return makeDefElem("isreplication", &nodes.Boolean{Boolval: true})
	case "noreplication":
		p.advance()
		return makeDefElem("isreplication", &nodes.Boolean{Boolval: false})
	case "bypassrls":
		p.advance()
		return makeDefElem("bypassrls", &nodes.Boolean{Boolval: true})
	case "nobypassrls":
		p.advance()
		return makeDefElem("bypassrls", &nodes.Boolean{Boolval: false})
	case "noinherit":
		p.advance()
		return makeDefElem("inherit", &nodes.Boolean{Boolval: false})
	case "externalid":
		p.advance()
		if p.cur.Type == TO {
			p.advance()
		}
		externalID, _ := p.parseColLabel()
		return makeDefElem("externalid", &nodes.String{Str: externalID})
	case "syslog":
		p.advance()
		p.expect(ACCESS)
		access, _ := p.parseColLabel()
		return makeDefElem("syslog_access", &nodes.String{Str: access})

	default:
		return nil
	}
}

func (p *Parser) parseCreateOptRoleElem() (*nodes.DefElem, error) {
	switch p.cur.Type {
	case SYSID:
		p.advance()
		val := p.cur.Ival
		p.advance()
		return makeDefElem("sysid", &nodes.Integer{Ival: val}), nil
	case ADMIN:
		p.advance()
		roles, err := p.parseRoleList()
		if err != nil {
			return nil, err
		}
		return makeDefElem("adminmembers", roles), nil
	case ROLE, USER:
		p.advance()
		roles, err := p.parseRoleList()
		if err != nil {
			return nil, err
		}
		return makeDefElem("rolemembers", roles), nil
	case IN_P:
		p.advance()
		if p.cur.Type == GROUP_P {
			p.advance()
		} else if _, err := p.expect(ROLE); err != nil {
			return nil, err
		}
		roles, err := p.parseRoleList()
		if err != nil {
			return nil, err
		}
		return makeDefElem("addroleto", roles), nil
	default:
		return p.parseAlterOptRoleElem(), nil
	}
}

func (p *Parser) parseAlterRoleStmt() (nodes.Node, error) {
	loc := p.prev.Loc // position of ALTER keyword
	isUser := p.cur.Type == USER
	p.advance() // consume ROLE or USER
	// ALTER ROLE ALL SET/RESET/IN DATABASE ...
	if p.cur.Type == ALL {
		p.advance()
		return p.parseAlterRoleSetStmtSuffix(loc, nil)
	}
	var role *nodes.RoleSpec
	var err error
	if isUser {
		role, err = p.parseRedshiftUserSpec()
	} else {
		role, err = p.parseRoleSpec()
	}
	if err != nil {
		return nil, err
	}
	p.parseGrantOptWith()
	if p.cur.Type == RENAME || p.cur.Type == OWNER {
		return p.parseRedshiftAlterRoleClauses(role, loc, isUser)
	}
	if (!isUser && (p.cur.Type == SET || p.cur.Type == RESET)) || p.cur.Type == IN_P {
		return p.parseAlterRoleSetStmtSuffix(loc, role)
	}
	if isUser {
		options, err := p.parseRedshiftAlterUserOptionList()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterRoleStmt{Role: role, Options: options, Action: 1,
			Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
	}
	options, err := p.parseOptRoleList(false)
	if err != nil {
		return nil, err
	}
	return &nodes.AlterRoleStmt{Role: role, Options: options, Action: 1,
		Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
}

// parseRedshiftAlterUserOptionList parses r_ALTER_USER's
// alterUser_option (',' alterUser_option)*; PostgreSQL's ALTER ROLE takes its
// options without commas between them, so the separators are ALTER USER's
// alone. RENAME TO is one of those options and takes part in the list.
func (p *Parser) parseRedshiftAlterUserOptionList() (*nodes.List, error) {
	var options *nodes.List
	for {
		started := false
		if p.cur.Type == RENAME {
			rename, err := p.parseRedshiftAlterUserRename()
			if err != nil {
				return nil, err
			}
			options = appendToList(options, rename)
			started = true
		} else {
			more, err := p.parseOptRoleList(false)
			if err != nil {
				return nil, err
			}
			if more != nil {
				options = appendNodeLists(options, more)
				started = true
			}
		}
		if !started {
			if options == nil {
				// alterUser_option (',' alterUser_option)*: a comma follows an
				// option, and this list holds none.
				return nil, nil
			}
			return nil, p.syntaxErrorAtCur()
		}
		if p.cur.Type != ',' {
			return options, nil
		}
		p.advance()
	}
}

// parseRedshiftAlterUserRename parses r_ALTER_USER's RENAME TO name option.
func (p *Parser) parseRedshiftAlterUserRename() (*nodes.DefElem, error) {
	p.advance() // consume RENAME
	if _, err := p.expect(TO); err != nil {
		return nil, err
	}
	name, err := p.parseName()
	if err != nil {
		return nil, err
	}
	return makeDefElem("rename_to", &nodes.String{Str: name}), nil
}

// parseRedshiftAlterRoleClauses parses r_ALTER_ROLE's RENAME TO and OWNER TO
// clauses, which follow the role name and repeat with commas between them,
// and the EXTERNALID TO name that may end the statement:
//
//	ALTER ROLE role_name [ WITH ] { RENAME TO name | OWNER TO role } [, ...]
//	    [ EXTERNALID TO name ]
//
// One clause on its own keeps the node PostgreSQL gives it, RenameStmt or
// AlterOwnerStmt, because that is what a consumer of ALTER ROLE RENAME TO
// reads. A repeat, or an EXTERNALID tail, is an AlterRoleStmt's option list
// instead: rename_to and owner_to hold the clauses in the order written, and
// externalid the tail. r_ALTER_USER writes the same EXTERNALID clause as one
// of its comma-separated options.
//
// Grammar reference: alterRole in RedshiftDocs.g4.
func (p *Parser) parseRedshiftAlterRoleClauses(role *nodes.RoleSpec, loc int, isUser bool) (nodes.Node, error) {
	var opts []nodes.Node
	clauses := 0
	for {
		switch p.cur.Type {
		case RENAME:
			p.advance()
			if _, err := p.expect(TO); err != nil {
				return nil, err
			}
			newname, err := p.parseName()
			if err != nil {
				return nil, err
			}
			opts = append(opts, makeDefElem("rename_to", &nodes.String{Str: newname}))
		case OWNER:
			p.advance()
			if _, err := p.expect(TO); err != nil {
				return nil, err
			}
			roleSpec, err := p.parseRoleSpec()
			if err != nil {
				return nil, err
			}
			opts = append(opts, makeDefElem("owner_to", roleSpec))
		default:
			return nil, p.syntaxErrorAtCur()
		}
		clauses++
		if p.cur.Type != ',' {
			break
		}
		p.advance()
		if p.cur.Type == RENAME || (!isUser && p.cur.Type == OWNER) {
			continue
		}
		if !isUser {
			// ALTER ROLE's comma takes another RENAME TO or OWNER TO clause.
			return nil, p.syntaxErrorAtCur()
		}
		// ALTER USER's comma goes on with its next option, RENAME TO and
		// CREATEDB among them.
		more, err := p.parseRedshiftAlterUserOptionList()
		if err != nil {
			return nil, err
		}
		if more == nil {
			return nil, p.syntaxErrorAtCur()
		}
		opts = append(opts, more.Items...)
		clauses = 0 // the option list holds more than the one clause
		break
	}

	hasExternalID := false
	if !isUser && p.redshiftKeywordIs("externalid") {
		p.advance()
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		externalID, err := p.parseColLabel()
		if err != nil {
			return nil, err
		}
		opts = append(opts, makeDefElem("externalid", &nodes.String{Str: externalID}))
		hasExternalID = true
	}

	if clauses == 1 && !hasExternalID {
		first := opts[0].(*nodes.DefElem)
		switch first.Defname {
		case "rename_to":
			return &nodes.RenameStmt{
				RenameType: nodes.OBJECT_ROLE,
				Subname:    role.Rolename,
				Newname:    first.Arg.(*nodes.String).Str,
				Loc:        nodes.Loc{Start: loc, End: p.prev.End},
			}, nil
		case "owner_to":
			return &nodes.AlterOwnerStmt{
				ObjectType: nodes.OBJECT_ROLE,
				Object:     &nodes.String{Str: role.Rolename},
				Newowner:   first.Arg.(*nodes.RoleSpec),
				Loc:        nodes.Loc{Start: loc, End: p.prev.End},
			}, nil
		}
	}
	return &nodes.AlterRoleStmt{Role: role, Options: &nodes.List{Items: opts}, Action: 1,
		Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
}

// parseAlterRoleSetStmtSuffix parses AlterRoleSetStmt after the role.
//
//	AlterRoleSetStmt: ALTER ROLE RoleSpec opt_in_database SetResetClause
//	    | ALTER ROLE ALL opt_in_database SetResetClause
//	opt_in_database: /* EMPTY */ | IN_P DATABASE name
//	SetResetClause: SET set_rest | VariableResetStmt
func (p *Parser) parseAlterRoleSetStmtSuffix(loc int, role *nodes.RoleSpec) (nodes.Node, error) {
	var dbname string
	if p.cur.Type == IN_P {
		p.advance()
		if _, err := p.expect(DATABASE); err != nil {
			return nil, err
		}
		var err error
		if dbname, err = p.parseName(); err != nil {
			return nil, err
		}
	}
	setstmt, err := p.parseSetResetClause()
	if err != nil {
		return nil, err
	}
	return &nodes.AlterRoleSetStmt{Role: role, Database: dbname, Setstmt: setstmt,
		Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
}

func (p *Parser) parseAlterGroupStmt() (nodes.Node, error) {
	loc := p.prev.Loc // position of ALTER keyword
	p.advance()       // consume GROUP
	role, err := p.parseRoleSpec()
	if err != nil {
		return nil, err
	}
	if p.cur.Type == ADD_P {
		p.advance()
		if _, err := p.expect(USER); err != nil {
			return nil, err
		}
		roles, err := p.parseRoleList()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterRoleStmt{
			Role: role, Options: &nodes.List{Items: []nodes.Node{makeDefElem("rolemembers", roles)}}, Action: 1,
			Loc: nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	} else if p.cur.Type == DROP {
		p.advance()
		if _, err := p.expect(USER); err != nil {
			return nil, err
		}
		roles, err := p.parseRoleList()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterRoleStmt{
			Role: role, Options: &nodes.List{Items: []nodes.Node{makeDefElem("rolemembers", roles)}}, Action: -1,
			Loc: nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	} else if p.cur.Type == RENAME {
		p.advance()
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		newname, err := p.parseRoleId()
		if err != nil {
			return nil, err
		}
		return &nodes.RenameStmt{
			RenameType: nodes.OBJECT_ROLE,
			Subname:    role.Rolename,
			Newname:    newname,
		}, nil
	}
	if p.collectMode() {
		for _, t := range []int{ADD_P, DROP, RENAME} {
			p.addTokenCandidate(t)
		}
		return nil, errCollecting
	}
	return nil, p.syntaxErrorAtCur()
}

func (p *Parser) parseDropRoleStmt() (nodes.Node, error) {
	loc := p.prev.Loc // position of DROP keyword
	p.advance()       // consume ROLE/USER/GROUP
	missingOk := p.parseOptIfExists()
	roles, err := p.parseRoleList()
	if err != nil {
		return nil, err
	}
	options := p.parseDropRoleOptions()
	return &nodes.DropRoleStmt{Roles: roles, Options: options, MissingOk: missingOk,
		Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
}

func (p *Parser) parseDropRoleOptions() *nodes.List {
	switch p.cur.Type {
	case FORCE:
		p.advance()
		return &nodes.List{Items: []nodes.Node{makeDefElem("force", nil)}}
	case RESTRICT:
		p.advance()
		return &nodes.List{Items: []nodes.Node{makeDefElem("restrict", nil)}}
	default:
		return nil
	}
}

func (p *Parser) parseCreatePolicyStmt() (nodes.Node, error) {
	loc := p.prev.Loc // position of CREATE keyword
	p.advance()       // consume POLICY
	policyName, err := p.parseName()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(ON); err != nil {
		return nil, err
	}
	nameStart := p.pos()
	names, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	table := makeRangeVarFromNames(names, nodes.Loc{Start: nameStart, End: p.prev.End})
	// RowSecurityDefaultPermissive: AS IDENT, whose action takes only
	// permissive and restrictive.
	permissive := true
	if p.cur.Type == AS {
		p.advance()
		if p.cur.Type != IDENT {
			if _, err := p.expect(IDENT); err != nil {
				return nil, err
			}
		}
		switch p.cur.Str {
		case "permissive":
		case "restrictive":
			permissive = false
		default:
			return nil, &ParseError{Message: fmt.Sprintf("unrecognized row security option \"%s\"", p.cur.Str), Position: p.cur.Loc}
		}
		p.advance()
	}
	cmdName := "all"
	if p.cur.Type == FOR {
		p.advance()
		if cmdName, err = p.parseRowSecurityCmd(); err != nil {
			return nil, err
		}
	}
	var roles *nodes.List
	if p.cur.Type == TO {
		p.advance()
		var err error
		roles, err = p.parseRoleList()
		if err != nil {
			return nil, err
		}
	}
	qual, withCheck, err := p.parseRowSecurityExprs()
	if err != nil {
		return nil, err
	}
	return &nodes.CreatePolicyStmt{
		PolicyName: policyName, Table: table, CmdName: cmdName,
		Permissive: permissive, Roles: roles, Qual: qual, WithCheck: withCheck,
		Loc: nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

func (p *Parser) parseAlterPolicyStmt() (nodes.Node, error) {
	loc := p.prev.Loc // position of ALTER keyword
	p.advance()       // consume POLICY
	// gram.y has IF EXISTS only in ALTER POLICY ... RENAME TO.
	missingOk := false
	if p.cur.Type == IF_P {
		p.advance()
		if _, err := p.expect(EXISTS); err != nil {
			return nil, err
		}
		missingOk = true
	}
	policyName, err := p.parseName()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(ON); err != nil {
		return nil, err
	}
	nameStart := p.pos()
	names, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	table := makeRangeVarFromNames(names, nodes.Loc{Start: nameStart, End: p.prev.End})
	// ALTER POLICY name ON qualified_name RENAME TO name
	if p.cur.Type == RENAME || missingOk {
		if _, err := p.expect(RENAME); err != nil {
			return nil, err
		}
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		newname, err := p.parseName()
		if err != nil {
			return nil, err
		}
		return &nodes.RenameStmt{
			RenameType: nodes.OBJECT_POLICY,
			Relation:   table,
			Subname:    policyName,
			Newname:    newname,
			MissingOk:  missingOk,
			Loc:        nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	}
	var roles *nodes.List
	if p.cur.Type == TO {
		p.advance()
		var err error
		roles, err = p.parseRoleList()
		if err != nil {
			return nil, err
		}
	}
	qual, withCheck, err := p.parseRowSecurityExprs()
	if err != nil {
		return nil, err
	}
	return &nodes.AlterPolicyStmt{
		PolicyName: policyName, Table: table, Roles: roles, Qual: qual, WithCheck: withCheck,
		Loc: nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

// parseRowSecurityExprs parses the RowSecurityOptionalExpr and
// RowSecurityOptionalWithCheck that end a policy statement.
//
//	RowSecurityOptionalExpr: USING '(' a_expr ')' | /* EMPTY */
//	RowSecurityOptionalWithCheck: WITH CHECK '(' a_expr ')' | /* EMPTY */
func (p *Parser) parseRowSecurityExprs() (qual, withCheck nodes.Node, err error) {
	if p.cur.Type == USING {
		p.advance()
		if qual, err = p.parseParenRequiredAExpr(); err != nil {
			return nil, nil, err
		}
	}
	if p.cur.Type == WITH {
		p.advance()
		if _, err = p.expect(CHECK); err != nil {
			return nil, nil, err
		}
		if withCheck, err = p.parseParenRequiredAExpr(); err != nil {
			return nil, nil, err
		}
	}
	return qual, withCheck, nil
}

// parseParenRequiredAExpr parses '(' a_expr ')'.
func (p *Parser) parseParenRequiredAExpr() (nodes.Node, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	expr, err := p.parseRequiredAExpr(0)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return expr, nil
}

// parseRowSecurityCmd parses row_security_cmd: ALL | SELECT | INSERT |
// UPDATE | DELETE_P.
func (p *Parser) parseRowSecurityCmd() (string, error) {
	var cmd string
	switch p.cur.Type {
	case ALL:
		cmd = "all"
	case SELECT:
		cmd = "select"
	case INSERT:
		cmd = "insert"
	case UPDATE:
		cmd = "update"
	case DELETE_P:
		cmd = "delete"
	default:
		if p.collectMode() {
			for _, t := range []int{ALL, SELECT, INSERT, UPDATE, DELETE_P} {
				p.addTokenCandidate(t)
			}
			return "", errCollecting
		}
		return "", p.syntaxErrorAtCur()
	}
	p.advance()
	return cmd, nil
}
