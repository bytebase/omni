package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/pg/ast"
)

func (p *Parser) parseGrantStmt() (nodes.Node, error) {
	loc := p.prev.Loc // position of GRANT keyword
	switch p.cur.Type {
	case ALL:
		p.advance()
		if p.cur.Type == PRIVILEGES {
			p.advance()
		}
		if p.cur.Type == ON {
			p.advance()
			return p.finishGrantOnObject(loc, true, nil)
		}
		if p.cur.Type == ',' || p.cur.Type == TO {
			roles := &nodes.List{Items: []nodes.Node{&nodes.AccessPriv{PrivName: "all"}}}
			for p.cur.Type == ',' {
				p.advance()
				name, _ := p.parseColId()
				roles.Items = append(roles.Items, &nodes.AccessPriv{PrivName: name})
			}
			return p.finishGrantRole(loc, true, roles)
		}
		return nil, nil
	case SELECT, INSERT, UPDATE, DELETE_P, TRUNCATE, REFERENCES, TRIGGER,
		CREATE, TEMPORARY, TEMP, EXECUTE:
		privs, err := p.parsePrivilegeList()
		if err != nil {
			return nil, err
		}
		if p.cur.Type == ON {
			p.advance()
			return p.finishGrantOnObject(loc, false, privs)
		}
		return p.finishGrantRole(loc, true, privs)
	default:
		privs, err := p.parsePrivilegeList()
		if err != nil {
			return nil, err
		}
		if p.cur.Type == ON {
			p.advance()
			return p.finishGrantOnObject(loc, false, privs)
		}
		return p.finishGrantRole(loc, true, privs)
	}
}

func (p *Parser) parseRevokeStmt() (nodes.Node, error) {
	loc := p.prev.Loc // position of REVOKE keyword
	grantOptionFor := false
	roleOptionName := "" // non-empty => "<name> OPTION FOR" form
	if p.cur.Type == GRANT {
		p.advance()
		if _, err := p.expect(OPTION); err != nil {
			return nil, err
		}
		if _, err := p.expect(FOR); err != nil {
			return nil, err
		}
		grantOptionFor = true
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
	switch p.cur.Type {
	case ALL:
		p.advance()
		if p.cur.Type == PRIVILEGES {
			p.advance()
		}
		if p.cur.Type == ON {
			p.advance()
			return p.finishRevokeOnObject(loc, grantOptionFor, nil)
		}
		if p.cur.Type == ',' || p.cur.Type == FROM {
			roles := &nodes.List{Items: []nodes.Node{&nodes.AccessPriv{PrivName: "all"}}}
			for p.cur.Type == ',' {
				p.advance()
				name, _ := p.parseColId()
				roles.Items = append(roles.Items, &nodes.AccessPriv{PrivName: name})
			}
			return p.finishRevokeRole(loc, roles, roleOptionName)
		}
		return nil, nil
	case SELECT, INSERT, UPDATE, DELETE_P, TRUNCATE, REFERENCES, TRIGGER,
		CREATE, TEMPORARY, TEMP, EXECUTE:
		privs, err := p.parsePrivilegeList()
		if err != nil {
			return nil, err
		}
		if p.cur.Type == ON {
			p.advance()
			return p.finishRevokeOnObject(loc, grantOptionFor, privs)
		}
		return p.finishRevokeRole(loc, privs, roleOptionName)
	default:
		privs, err := p.parsePrivilegeList()
		if err != nil {
			return nil, err
		}
		if p.cur.Type == ON {
			p.advance()
			return p.finishRevokeOnObject(loc, grantOptionFor, privs)
		}
		return p.finishRevokeRole(loc, privs, roleOptionName)
	}
}

func (p *Parser) finishGrantOnObject(loc int, allPrivs bool, privs *nodes.List) (nodes.Node, error) {
	var privileges *nodes.List
	if !allPrivs {
		privileges = privs
	}
	targtype, objtype, objects, err := p.parseGrantTarget()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(TO); err != nil {
		return nil, err
	}
	grantees, err := p.parseGranteeList()
	if err != nil {
		return nil, err
	}
	grantOption := p.parseOptGrantGrantOption()
	return &nodes.GrantStmt{
		IsGrant: true, Targtype: targtype, Objtype: objtype,
		Objects: objects, Privileges: privileges, Grantees: grantees,
		GrantOption: grantOption,
		Loc:         nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

func (p *Parser) finishRevokeOnObject(loc int, grantOptionFor bool, privs *nodes.List) (nodes.Node, error) {
	targtype, objtype, objects, err := p.parseGrantTarget()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(FROM); err != nil {
		return nil, err
	}
	grantees, err := p.parseGranteeList()
	if err != nil {
		return nil, err
	}
	behavior := p.parseOptDropBehavior()
	return &nodes.GrantStmt{
		IsGrant: false, Targtype: targtype, Objtype: objtype,
		Objects: objects, Privileges: privs, Grantees: grantees,
		GrantOption: grantOptionFor, Behavior: nodes.DropBehavior(behavior),
		Loc: nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

func (p *Parser) finishGrantRole(loc int, isGrant bool, roles *nodes.List) (nodes.Node, error) {
	if _, err := p.expect(TO); err != nil {
		return nil, err
	}
	grantees, err := p.parseRoleList()
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
	grantees, err := p.parseRoleList()
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

func (p *Parser) parseGrantTarget() (nodes.GrantTargetType, nodes.ObjectType, *nodes.List, error) {
	switch p.cur.Type {
	case TABLE:
		p.advance()
		list, err := p.parseGrantObjectAnyNameList()
		if err != nil {
			return 0, 0, nil, err
		}
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_TABLE, list, nil
	case SEQUENCE:
		p.advance()
		list, err := p.parseGrantObjectAnyNameList()
		if err != nil {
			return 0, 0, nil, err
		}
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_SEQUENCE, list, nil
	case FUNCTION:
		p.advance()
		list, err := p.parseGrantFunctionWithArgtypesList()
		if err != nil {
			return 0, 0, nil, err
		}
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_FUNCTION, list, nil
	case PROCEDURE:
		p.advance()
		list, err := p.parseGrantFunctionWithArgtypesList()
		if err != nil {
			return 0, 0, nil, err
		}
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_PROCEDURE, list, nil
	case ROUTINE:
		p.advance()
		list, err := p.parseGrantFunctionWithArgtypesList()
		if err != nil {
			return 0, 0, nil, err
		}
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_ROUTINE, list, nil
	case DATABASE:
		p.advance()
		names, err := p.parseNameList()
		if err != nil {
			return 0, 0, nil, err
		}
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_DATABASE, makeNameListAsAnyNameList(names), nil
	case DOMAIN_P:
		p.advance()
		objects, err := p.parseAnyNameList()
		if err != nil {
			return 0, 0, nil, err
		}
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_DOMAIN, objects, nil
	case LANGUAGE:
		p.advance()
		names, err := p.parseNameList()
		if err != nil {
			return 0, 0, nil, err
		}
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_LANGUAGE, makeNameListAsAnyNameList(names), nil
	case LARGE_P:
		p.advance()
		p.expect(OBJECT_P)
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_LARGEOBJECT, p.parseNumericOnlyList(), nil
	case SCHEMA:
		p.advance()
		names, err := p.parseNameList()
		if err != nil {
			return 0, 0, nil, err
		}
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_SCHEMA, makeNameListAsAnyNameList(names), nil
	case TABLESPACE:
		p.advance()
		names, err := p.parseNameList()
		if err != nil {
			return 0, 0, nil, err
		}
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_TABLESPACE, makeNameListAsAnyNameList(names), nil
	case TYPE_P:
		p.advance()
		objects, err := p.parseAnyNameList()
		if err != nil {
			return 0, 0, nil, err
		}
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_TYPE, objects, nil
	case FOREIGN:
		p.advance()
		if p.cur.Type == DATA_P {
			p.advance()
			p.expect(WRAPPER)
			names, err := p.parseNameList()
			if err != nil {
				return 0, 0, nil, err
			}
			return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_FDW, makeNameListAsAnyNameList(names), nil
		}
		p.expect(SERVER)
		names, err := p.parseNameList()
		if err != nil {
			return 0, 0, nil, err
		}
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_FOREIGN_SERVER, makeNameListAsAnyNameList(names), nil
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
			objType = nodes.OBJECT_TABLE
		}
		p.expect(IN_P)
		p.expect(SCHEMA)
		if p.collectMode() {
			p.addRuleCandidate("qualified_name")
			return nodes.ACL_TARGET_ALL_IN_SCHEMA, objType, nil, nil
		}
		names, err := p.parseNameList()
		if err != nil {
			return 0, 0, nil, err
		}
		return nodes.ACL_TARGET_ALL_IN_SCHEMA, objType, makeNameListAsAnyNameList(names), nil
	default:
		list, err := p.parseGrantObjectAnyNameList()
		if err != nil {
			return 0, 0, nil, err
		}
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_TABLE, list, nil
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

func (p *Parser) parseNumericOnlyList() *nodes.List {
	val := p.parseNumericOnly()
	if val == nil {
		return nil
	}
	result := &nodes.List{Items: []nodes.Node{val}}
	for p.cur.Type == ',' {
		p.advance()
		val = p.parseNumericOnly()
		if val == nil {
			break
		}
		result.Items = append(result.Items, val)
	}
	return result
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
	arg := p.parseFuncArg()
	// A func_arg always has a type; without one there is no argument, and
	// a typed nil must not reach the caller's nil check.
	if arg == nil || arg.ArgType == nil {
		return nil, nil
	}
	return arg.ArgType, nil
}

func (p *Parser) parsePrivileges() (*nodes.List, error) {
	if p.cur.Type == ALL {
		p.advance()
		if p.cur.Type == PRIVILEGES {
			p.advance()
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
	for p.cur.Type == ',' {
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
			return nil, nil
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

func (p *Parser) parseGrantee() (*nodes.RoleSpec, error) {
	if p.cur.Type == GROUP_P {
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

func (p *Parser) parseOptGrantGrantOption() bool {
	if p.cur.Type == WITH {
		next := p.peekNext()
		if next.Type == GRANT {
			p.advance()
			p.advance()
			p.expect(OPTION)
			return true
		}
	}
	return false
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
	name, err := p.parseRoleId()
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
	p.parseGrantOptWith()
	options, err := p.parseOptRoleList(true)
	if err != nil {
		return nil, err
	}
	return &nodes.CreateRoleStmt{StmtType: nodes.ROLESTMT_GROUP, Role: name, Options: options,
		Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
}

func (p *Parser) parseRoleId() (string, error) {
	return p.parseNonReservedWord()
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
		val := p.parseSignedIconst()
		return makeDefElem("connectionlimit", &nodes.Integer{Ival: val})
	case VALID:
		p.advance()
		p.expect(UNTIL)
		until := p.cur.Str
		p.advance()
		return makeDefElem("validUntil", &nodes.String{Str: until})
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
		} else {
			p.expect(ROLE)
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
	p.advance()       // consume ROLE or USER
	// ALTER ROLE ALL SET/RESET/IN DATABASE ...
	if p.cur.Type == ALL {
		p.advance()
		if p.cur.Type == SET || p.cur.Type == RESET {
			n := p.parseAlterRoleSetStmtSuffix(loc, nil, "")
			return n, nil
		}
		if p.cur.Type == IN_P {
			p.advance()
			if _, err := p.expect(DATABASE); err != nil {
				return nil, err
			}
			dbname, _ := p.parseName()
			n := p.parseAlterRoleSetStmtSuffix(loc, nil, dbname)
			return n, nil
		}
		return nil, nil
	}
	role, err := p.parseRoleSpec()
	if err != nil {
		return nil, err
	}
	if p.cur.Type == RENAME {
		p.advance()
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		newname, _ := p.parseName()
		return &nodes.RenameStmt{
			RenameType: nodes.OBJECT_ROLE,
			Subname:    role.Rolename,
			Newname:    newname,
		}, nil
	}
	if p.cur.Type == SET || p.cur.Type == RESET {
		n := p.parseAlterRoleSetStmtSuffix(loc, role, "")
		return n, nil
	}
	if p.cur.Type == IN_P {
		p.advance()
		if _, err := p.expect(DATABASE); err != nil {
			return nil, err
		}
		dbname, _ := p.parseName()
		n := p.parseAlterRoleSetStmtSuffix(loc, role, dbname)
		return n, nil
	}
	p.parseGrantOptWith()
	options, err := p.parseOptRoleList(false)
	if err != nil {
		return nil, err
	}
	return &nodes.AlterRoleStmt{Role: role, Options: options, Action: 1,
		Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
}

func (p *Parser) parseAlterRoleSetStmtSuffix(loc int, role *nodes.RoleSpec, dbname string) nodes.Node {
	var setstmt *nodes.VariableSetStmt
	if p.cur.Type == SET {
		p.advance()
		result, _ := p.parseVariableSetStmt()
		if vs, ok := result.(*nodes.VariableSetStmt); ok {
			setstmt = vs
		}
	} else if p.cur.Type == RESET {
		p.advance()
		result, _ := p.parseVariableResetStmt()
		if vs, ok := result.(*nodes.VariableSetStmt); ok {
			setstmt = vs
		}
	}
	return &nodes.AlterRoleSetStmt{Role: role, Database: dbname, Setstmt: setstmt,
		Loc: nodes.Loc{Start: loc, End: p.prev.End}}
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
		newname, _ := p.parseName()
		return &nodes.RenameStmt{
			RenameType: nodes.OBJECT_ROLE,
			Subname:    role.Rolename,
			Newname:    newname,
		}, nil
	}
	return nil, nil
}

func (p *Parser) parseDropRoleStmt() (nodes.Node, error) {
	loc := p.prev.Loc // position of DROP keyword
	p.advance()       // consume ROLE/USER/GROUP
	missingOk := p.parseOptIfExists()
	roles, err := p.parseRoleList()
	if err != nil {
		return nil, err
	}
	return &nodes.DropRoleStmt{Roles: roles, MissingOk: missingOk,
		Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
}

func (p *Parser) parseCreatePolicyStmt() (nodes.Node, error) {
	loc := p.prev.Loc // position of CREATE keyword
	p.advance()       // consume POLICY
	policyName, _ := p.parseName()
	if _, err := p.expect(ON); err != nil {
		return nil, err
	}
	nameStart := p.pos()
	names, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	table := makeRangeVarFromNames(names, nodes.Loc{Start: nameStart, End: p.prev.End})
	permissive := true
	if p.cur.Type == AS {
		p.advance()
		if p.isColId() && p.cur.Str == "permissive" {
			p.advance()
		} else if p.isColId() && p.cur.Str == "restrictive" {
			p.advance()
			permissive = false
		}
	}
	cmdName := "all"
	if p.cur.Type == FOR {
		p.advance()
		cmdName = p.parseRowSecurityCmd()
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
	var qual nodes.Node
	if p.cur.Type == USING {
		p.advance()
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		qual, _ = p.parseAExpr(0)
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}
	var withCheck nodes.Node
	if p.cur.Type == WITH {
		next := p.peekNext()
		if next.Type == CHECK {
			p.advance()
			p.advance()
			if _, err := p.expect('('); err != nil {
				return nil, err
			}
			withCheck, _ = p.parseAExpr(0)
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
		}
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
	policyName, _ := p.parseName()
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
	if p.cur.Type == RENAME {
		p.advance()
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		newname, _ := p.parseName()
		return &nodes.RenameStmt{
			RenameType: nodes.OBJECT_POLICY,
			Relation:   table,
			Subname:    policyName,
			Newname:    newname,
			Loc:        nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	}
	var roles *nodes.List
	var qual nodes.Node
	var withCheck nodes.Node
	if p.cur.Type == TO {
		p.advance()
		var err error
		roles, err = p.parseRoleList()
		if err != nil {
			return nil, err
		}
	}
	if p.cur.Type == USING {
		p.advance()
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		qual, _ = p.parseAExpr(0)
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}
	if p.cur.Type == WITH {
		next := p.peekNext()
		if next.Type == CHECK {
			p.advance()
			p.advance()
			if _, err := p.expect('('); err != nil {
				return nil, err
			}
			withCheck, _ = p.parseAExpr(0)
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
		}
	}
	return &nodes.AlterPolicyStmt{
		PolicyName: policyName, Table: table, Roles: roles, Qual: qual, WithCheck: withCheck,
		Loc: nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

func (p *Parser) parseRowSecurityCmd() string {
	switch p.cur.Type {
	case ALL:
		p.advance()
		return "all"
	case SELECT:
		p.advance()
		return "select"
	case INSERT:
		p.advance()
		return "insert"
	case UPDATE:
		p.advance()
		return "update"
	case DELETE_P:
		p.advance()
		return "delete"
	default:
		return "all"
	}
}
