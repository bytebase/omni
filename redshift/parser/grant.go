package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/redshift/ast"
)

func (p *Parser) parseGrantStmt() (nodes.Node, error) {
	loc := p.prev.Loc // position of GRANT keyword
	if p.redshiftWordEqual("assumerole") {
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
	if p.cur.Type == CREATE && strings.EqualFold(p.peekNext().Str, "model") {
		p.advance()
		p.advance()
		roles := &nodes.List{Items: []nodes.Node{&nodes.AccessPriv{PrivName: "create_model"}}}
		return p.finishGrantRole(loc, true, roles)
	}
	if p.redshiftScopedPermissionStart() {
		return p.parseRedshiftScopedPermissionStmt(loc, true, false)
	}
	if roles, ok, err := p.parseRedshiftSystemPrivileges(TO); ok || err != nil {
		if err != nil {
			return nil, err
		}
		return p.finishGrantRole(loc, true, roles)
	}
	switch p.cur.Type {
	case ALL:
		allLoc := p.pos()
		p.advance()
		if p.cur.Type == PRIVILEGES {
			p.advance()
		}
		cols, err := p.parseOptColumnList()
		if err != nil {
			return nil, err
		}
		var privileges *nodes.List
		if cols != nil {
			privileges = &nodes.List{Items: []nodes.Node{&nodes.AccessPriv{Cols: cols, Loc: nodes.Loc{Start: allLoc, End: p.prev.End}}}}
		}
		if p.cur.Type == ON {
			p.advance()
			return p.finishGrantOnObject(loc, privileges == nil, privileges)
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
	if p.redshiftWordEqual("assumerole") {
		return p.parseRedshiftAssumeRoleStmt(loc, false)
	}
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
	if p.cur.Type == ROLE {
		roles, err := p.parseRedshiftGrantRoleNameList()
		if err != nil {
			return nil, err
		}
		return p.finishRevokeRole(loc, roles, roleOptionName)
	}
	if p.cur.Type == CREATE && strings.EqualFold(p.peekNext().Str, "model") {
		p.advance()
		p.advance()
		roles := &nodes.List{Items: []nodes.Node{&nodes.AccessPriv{PrivName: "create_model"}}}
		return p.finishRevokeRole(loc, roles, roleOptionName)
	}
	if p.redshiftScopedPermissionStart() {
		return p.parseRedshiftScopedPermissionStmt(loc, false, grantOptionFor)
	}
	if roles, ok, err := p.parseRedshiftSystemPrivileges(FROM); ok || err != nil {
		if err != nil {
			return nil, err
		}
		return p.finishRevokeRole(loc, roles, roleOptionName)
	}
	switch p.cur.Type {
	case ALL:
		allLoc := p.pos()
		p.advance()
		if p.cur.Type == PRIVILEGES {
			p.advance()
		}
		cols, err := p.parseOptColumnList()
		if err != nil {
			return nil, err
		}
		var privileges *nodes.List
		if cols != nil {
			privileges = &nodes.List{Items: []nodes.Node{&nodes.AccessPriv{Cols: cols, Loc: nodes.Loc{Start: allLoc, End: p.prev.End}}}}
		}
		if p.cur.Type == ON {
			p.advance()
			return p.finishRevokeOnObject(loc, grantOptionFor, privileges)
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
	grantees := p.parseGranteeList()
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
	grantees := p.parseGranteeList()
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
	grantees := p.parseGranteeList()
	opts, err := p.parseGrantRoleOptList()
	if err != nil {
		return nil, err
	}
	grantedBy := p.parseOptGrantedBy()
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
	grantees := p.parseGranteeList()
	grantedBy := p.parseOptGrantedBy()
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
	if p.redshiftWordEqual("datashare") {
		p.advance()
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_DATASHARE, p.parseGrantObjectAnyNameList(), nil
	}
	if p.redshiftWordEqual("model") {
		p.advance()
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_MODEL, p.parseGrantObjectAnyNameList(), nil
	}
	switch p.cur.Type {
	case COPY:
		p.advance()
		if !p.consumeRedshiftWord("job") {
			return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_COPY_JOB, nil, p.syntaxErrorAtCur()
		}
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_COPY_JOB, p.parseGrantObjectAnyNameList(), nil
	case EXTERNAL:
		p.advance()
		switch p.cur.Type {
		case SCHEMA:
			p.advance()
			names, _ := p.parseNameList()
			return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_SCHEMA, makeNameListAsAnyNameList(names), nil
		case TABLE:
			p.advance()
			return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_FOREIGN_TABLE, p.parseGrantObjectAnyNameList(), nil
		default:
			return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_FOREIGN_TABLE, p.parseGrantObjectAnyNameList(), nil
		}
	case TABLE:
		p.advance()
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_TABLE, p.parseGrantObjectAnyNameList(), nil
	case SEQUENCE:
		p.advance()
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_SEQUENCE, p.parseGrantObjectAnyNameList(), nil
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
		names, _ := p.parseNameList()
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_DATABASE, makeNameListAsAnyNameList(names), nil
	case DOMAIN_P:
		p.advance()
		objects, _ := p.parseAnyNameList()
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_DOMAIN, objects, nil
	case LANGUAGE:
		p.advance()
		names, _ := p.parseNameList()
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_LANGUAGE, makeNameListAsAnyNameList(names), nil
	case LARGE_P:
		p.advance()
		p.expect(OBJECT_P)
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_LARGEOBJECT, p.parseNumericOnlyList(), nil
	case SCHEMA:
		p.advance()
		names, _ := p.parseNameList()
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_SCHEMA, makeNameListAsAnyNameList(names), nil
	case TABLESPACE:
		p.advance()
		names, _ := p.parseNameList()
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_TABLESPACE, makeNameListAsAnyNameList(names), nil
	case TYPE_P:
		p.advance()
		objects, _ := p.parseAnyNameList()
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_TYPE, objects, nil
	case FOREIGN:
		p.advance()
		if p.cur.Type == DATA_P {
			p.advance()
			p.expect(WRAPPER)
			names, _ := p.parseNameList()
			return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_FDW, makeNameListAsAnyNameList(names), nil
		}
		p.expect(SERVER)
		names, _ := p.parseNameList()
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
		names, _ := p.parseNameList()
		return nodes.ACL_TARGET_ALL_IN_SCHEMA, objType, makeNameListAsAnyNameList(names), nil
	default:
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_TABLE, p.parseGrantObjectAnyNameList(), nil
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

func (p *Parser) parseGrantObjectAnyNameList() *nodes.List {
	name, err := p.parseAnyName()
	if err != nil {
		return nil
	}
	result := &nodes.List{Items: []nodes.Node{makeRangeVarFromAnyName(name)}}
	for p.cur.Type == ',' {
		p.advance()
		name, err = p.parseAnyName()
		if err != nil {
			break
		}
		result.Items = append(result.Items, makeRangeVarFromAnyName(name))
	}
	return result
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

func (p *Parser) redshiftScopedPermissionStart() bool {
	return p.peekNext().Type == FOR
}

func (p *Parser) parseRedshiftScopedPermissionStmt(loc int, isGrant bool, grantOption bool) (nodes.Node, error) {
	privs, err := p.parsePrivileges()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(FOR); err != nil {
		return nil, err
	}
	objtype, objects, err := p.parseRedshiftScopedPermissionTarget()
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
	grantees := p.parseGranteeList()
	behavior := nodes.DropBehavior(0)
	if !isGrant {
		behavior = nodes.DropBehavior(p.parseOptDropBehavior())
	}
	return &nodes.GrantStmt{
		IsGrant:     isGrant,
		Targtype:    nodes.ACL_TARGET_SCOPED,
		Objtype:     objtype,
		Objects:     objects,
		Privileges:  privs,
		Grantees:    grantees,
		GrantOption: grantOption,
		Behavior:    behavior,
		Loc:         nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
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
	case p.redshiftWordEqual("languages"):
		p.advance()
		objtype = nodes.OBJECT_LANGUAGE
	case p.cur.Type == COPY:
		p.advance()
		if !p.consumeRedshiftWord("jobs") {
			return 0, nil, p.syntaxErrorAtCur()
		}
		objtype = nodes.OBJECT_COPY_JOB
	default:
		return 0, nil, p.syntaxErrorAtCur()
	}
	if _, err := p.expect(IN_P); err != nil {
		return 0, nil, err
	}
	var items []nodes.Node
	switch p.cur.Type {
	case SCHEMA:
		p.advance()
		name, err := p.parseColId()
		if err != nil {
			return 0, nil, err
		}
		items = append(items, makeDefElem("schema", &nodes.String{Str: name}))
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
	if p.cur.Type == DATABASE {
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
	grantees := p.parseGranteeList()
	if _, err := p.expect(FOR); err != nil {
		return nil, err
	}
	action := p.parseRedshiftAssumeRoleAction()
	opts := &nodes.List{Items: []nodes.Node{makeDefElem("for", &nodes.String{Str: action})}}
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

func (p *Parser) parseRedshiftAssumeRoleAction() string {
	var parts []string
	for p.cur.Type != 0 && p.cur.Type != ';' {
		parts = append(parts, strings.ToLower(p.cur.Str))
		p.advance()
	}
	return strings.Join(parts, "_")
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
			break
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
			break
		}
		result.Items = append(result.Items, arg)
	}
	return result, nil
}

func (p *Parser) parseGrantFuncArg() (nodes.Node, error) {
	arg := p.parseFuncArg()
	if arg == nil {
		return nil, nil
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
	for p.cur.Type == ',' {
		p.advance()
		priv, err = p.parsePrivilege()
		if err != nil {
			return nil, err
		}
		if priv == nil {
			break
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

func (p *Parser) parseGranteeList() *nodes.List {
	g := p.parseGrantee()
	if g == nil {
		return nil
	}
	result := &nodes.List{Items: []nodes.Node{g}}
	for p.cur.Type == ',' {
		p.advance()
		g = p.parseGrantee()
		if g == nil {
			break
		}
		result.Items = append(result.Items, g)
	}
	return result
}

func (p *Parser) parseGrantee() *nodes.RoleSpec {
	if p.redshiftWordEqual("rls") {
		loc := p.pos()
		p.advance()
		if p.cur.Type == POLICY {
			p.advance()
		}
		name, err := p.parseColId()
		if err != nil {
			return nil
		}
		return &nodes.RoleSpec{
			Roletype: int(nodes.ROLESPEC_CSTRING),
			Rolename: "rls_policy:" + name,
			Loc:      nodes.Loc{Start: loc, End: p.pos()},
		}
	}
	if p.redshiftWordEqual("namespace") || p.redshiftWordEqual("account") {
		loc := p.pos()
		prefix := strings.ToLower(p.cur.Str)
		p.advance()
		if p.cur.Type != SCONST {
			return p.parseRoleSpec()
		}
		tok := p.advance()
		if prefix == "account" && p.consumeRedshiftWord("via") {
			p.consumeRedshiftWord("data_catalog")
		}
		return &nodes.RoleSpec{
			Roletype: int(nodes.ROLESPEC_CSTRING),
			Rolename: prefix + ":" + tok.Str,
			Loc:      nodes.Loc{Start: loc, End: p.pos()},
		}
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
			Roletype: int(nodes.ROLESPEC_CSTRING),
			Rolename: "iam_role:" + value,
			Loc:      nodes.Loc{Start: loc, End: p.pos()},
		}
	}
	if p.cur.Type == GROUP_P || p.cur.Type == ROLE {
		p.advance()
	}
	loc := p.pos()
	// PUBLIC is not a reserved keyword; it appears as IDENT
	if p.isColId() && p.cur.Str == "public" {
		p.advance()
		return &nodes.RoleSpec{Roletype: int(nodes.ROLESPEC_PUBLIC), Loc: nodes.Loc{Start: loc, End: p.pos()}}
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

func (p *Parser) parseOptGrantedBy() *nodes.RoleSpec {
	if p.cur.Type == GRANTED {
		p.advance()
		p.expect(BY)
		return p.parseRoleSpec()
	}
	return nil
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
	name := p.parseRoleId()
	p.parseGrantOptWith()
	options := p.parseOptRoleList(true)
	return &nodes.CreateRoleStmt{StmtType: nodes.ROLESTMT_ROLE, Role: name, Options: options,
		Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
}

func (p *Parser) parseCreateUserStmt() (nodes.Node, error) {
	loc := p.prev.Loc // position of CREATE keyword
	p.advance()       // consume USER
	name := p.parseRedshiftUserId()
	p.parseGrantOptWith()
	options := p.parseOptRoleList(true)
	return &nodes.CreateRoleStmt{StmtType: nodes.ROLESTMT_USER, Role: name, Options: options,
		Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
}

func (p *Parser) parseCreateGroupStmt() (nodes.Node, error) {
	loc := p.prev.Loc // position of CREATE keyword
	p.advance()       // consume GROUP
	name := p.parseRoleId()
	p.parseGrantOptWith()
	options := p.parseOptRoleList(true)
	return &nodes.CreateRoleStmt{StmtType: nodes.ROLESTMT_GROUP, Role: name, Options: options,
		Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
}

func (p *Parser) parseRoleId() string {
	name, _ := p.parseColId()
	return name
}

func (p *Parser) parseRedshiftUserId() string {
	name := p.parseRoleId()
	for p.cur.Type == ':' {
		p.advance()
		part, _ := p.parseColId()
		name += ":" + part
	}
	return name
}

func (p *Parser) parseRedshiftUserSpec() *nodes.RoleSpec {
	loc := p.pos()
	name := p.parseRedshiftUserId()
	return &nodes.RoleSpec{
		Roletype: int(nodes.ROLESPEC_CSTRING),
		Rolename: name,
		Loc:      nodes.Loc{Start: loc, End: p.pos()},
	}
}

func (p *Parser) parseGrantOptWith() {
	if p.cur.Type == WITH {
		p.advance()
	}
}

func (p *Parser) parseOptRoleList(isCreate bool) *nodes.List {
	var items []nodes.Node
	for {
		var elem *nodes.DefElem
		if isCreate {
			elem = p.parseCreateOptRoleElem()
		} else {
			elem = p.parseAlterOptRoleElem()
		}
		if elem == nil {
			break
		}
		items = append(items, elem)
	}
	if len(items) == 0 {
		return nil
	}
	return &nodes.List{Items: items}
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

func (p *Parser) parseCreateOptRoleElem() *nodes.DefElem {
	switch p.cur.Type {
	case SYSID:
		p.advance()
		val := p.cur.Ival
		p.advance()
		return makeDefElem("sysid", &nodes.Integer{Ival: val})
	case ADMIN:
		p.advance()
		roles := p.parseRoleList()
		return makeDefElem("adminmembers", roles)
	case ROLE, USER:
		p.advance()
		roles := p.parseRoleList()
		return makeDefElem("rolemembers", roles)
	case IN_P:
		p.advance()
		if p.cur.Type == GROUP_P {
			p.advance()
		} else {
			p.expect(ROLE)
		}
		roles := p.parseRoleList()
		return makeDefElem("addroleto", roles)
	default:
		return p.parseAlterOptRoleElem()
	}
}

func (p *Parser) parseAlterRoleStmt() (nodes.Node, error) {
	loc := p.prev.Loc // position of ALTER keyword
	isUser := p.cur.Type == USER
	p.advance() // consume ROLE or USER
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
	var role *nodes.RoleSpec
	if isUser {
		role = p.parseRedshiftUserSpec()
	} else {
		role = p.parseRoleSpec()
	}
	p.parseGrantOptWith()
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
	if p.cur.Type == OWNER {
		p.advance()
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		roleSpec := p.parseRoleSpec()
		return &nodes.AlterOwnerStmt{
			ObjectType: nodes.OBJECT_ROLE,
			Object:     &nodes.String{Str: role.Rolename},
			Newowner:   roleSpec,
			Loc:        nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	}
	if !isUser && (p.cur.Type == SET || p.cur.Type == RESET) {
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
	options := p.parseOptRoleList(false)
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
	role := p.parseRoleSpec()
	if p.cur.Type == ADD_P {
		p.advance()
		if _, err := p.expect(USER); err != nil {
			return nil, err
		}
		roles := p.parseRoleList()
		return &nodes.AlterRoleStmt{
			Role: role, Options: &nodes.List{Items: []nodes.Node{makeDefElem("rolemembers", roles)}}, Action: 1,
			Loc: nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	} else if p.cur.Type == DROP {
		p.advance()
		if _, err := p.expect(USER); err != nil {
			return nil, err
		}
		roles := p.parseRoleList()
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
	roles := p.parseRoleList()
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
	policyName, _ := p.parseName()
	if _, err := p.expect(ON); err != nil {
		return nil, err
	}
	names, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	table := makeRangeVarFromNames(names)
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
		roles = p.parseRoleList()
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
	names, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	table := makeRangeVarFromNames(names)
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
		roles = p.parseRoleList()
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
