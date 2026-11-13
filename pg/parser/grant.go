package parser

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/pg/ast"
)

// parseGrantStmt parses GRANT after its keyword.
//
//	GrantStmt: GRANT privileges ON privilege_target TO grantee_list
//	           opt_grant_grant_option opt_granted_by
//	GrantRoleStmt: GRANT privilege_list TO role_list opt_granted_by
//	             | GRANT privilege_list TO role_list WITH grant_role_opt_list opt_granted_by
func (p *Parser) parseGrantStmt() (nodes.Node, error) {
	loc := p.prev.Loc // position of GRANT keyword
	privs, all, err := p.parsePrivilegesOrAll()
	if err != nil {
		return nil, err
	}
	if all || p.cur.Type == ON {
		if _, err := p.expect(ON); err != nil {
			return nil, err
		}
		return p.finishGrantOnObject(loc, privs)
	}
	return p.finishGrantRole(loc, true, privs)
}

// parseRevokeStmt parses REVOKE after its keyword.
//
//	RevokeStmt: REVOKE privileges ON privilege_target FROM grantee_list opt_granted_by opt_drop_behavior
//	          | REVOKE GRANT OPTION FOR privileges ON privilege_target FROM grantee_list opt_granted_by opt_drop_behavior
//	RevokeRoleStmt: REVOKE privilege_list FROM role_list opt_granted_by opt_drop_behavior
//	              | REVOKE ColId OPTION FOR privilege_list FROM role_list opt_granted_by opt_drop_behavior
func (p *Parser) parseRevokeStmt() (nodes.Node, error) {
	loc := p.prev.Loc // position of REVOKE keyword
	if p.cur.Type == GRANT {
		p.advance()
		if _, err := p.expect(OPTION); err != nil {
			return nil, err
		}
		if _, err := p.expect(FOR); err != nil {
			return nil, err
		}
		privs, _, err := p.parsePrivilegesOrAll()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(ON); err != nil {
			return nil, err
		}
		return p.finishRevokeOnObject(loc, true, privs)
	}
	roleOptionName := "" // non-empty => "<name> OPTION FOR" form
	if p.isColId() && p.peekNext().Type == OPTION {
		roleOptionName = strings.ToLower(p.cur.Str)
		p.advance() // ColId
		p.advance() // OPTION
		if _, err := p.expect(FOR); err != nil {
			return nil, err
		}
		privs, err := p.parsePrivilegeList()
		if err != nil {
			return nil, err
		}
		return p.finishRevokeRole(loc, privs, roleOptionName)
	}
	privs, all, err := p.parsePrivilegesOrAll()
	if err != nil {
		return nil, err
	}
	if all || p.cur.Type == ON {
		if _, err := p.expect(ON); err != nil {
			return nil, err
		}
		return p.finishRevokeOnObject(loc, false, privs)
	}
	return p.finishRevokeRole(loc, privs, "")
}

func (p *Parser) finishGrantOnObject(loc int, privs *nodes.List) (nodes.Node, error) {
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
	grantOption, err := p.parseOptGrantGrantOption()
	if err != nil {
		return nil, err
	}
	grantor, err := p.parseOptGrantedBy()
	if err != nil {
		return nil, err
	}
	return &nodes.GrantStmt{
		IsGrant: true, Targtype: targtype, Objtype: objtype,
		Objects: objects, Privileges: privs, Grantees: grantees,
		GrantOption: grantOption, Grantor: grantor,
		Loc: nodes.Loc{Start: loc, End: p.prev.End},
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
	grantor, err := p.parseOptGrantedBy()
	if err != nil {
		return nil, err
	}
	behavior := p.parseOptDropBehavior()
	return &nodes.GrantStmt{
		IsGrant: false, Targtype: targtype, Objtype: objtype,
		Objects: objects, Privileges: privs, Grantees: grantees,
		GrantOption: grantOptionFor, Grantor: grantor, Behavior: nodes.DropBehavior(behavior),
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

// parseGrantTarget parses privilege_target.
//
// Its object type keywords are unreserved, so each can also start the
// qualified_name_list of tables: GRANT SELECT ON sequence TO a names a
// table called sequence. gram.y reads the keyword as a name when the token
// after it can only follow a name.
func (p *Parser) parseGrantTarget() (nodes.GrantTargetType, nodes.ObjectType, *nodes.List, error) {
	if !p.collectMode() && p.cur.Type != IDENT && p.isColId() {
		switch p.peekNext().Type {
		case TO, FROM, ',', '.':
			list, err := p.parseGrantObjectAnyNameList()
			if err != nil {
				return 0, 0, nil, err
			}
			return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_TABLE, list, nil
		}
	}
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
		if _, err := p.expect(OBJECT_P); err != nil {
			return 0, 0, nil, err
		}
		list, err := p.parseNumericOnlyList()
		if err != nil {
			return 0, 0, nil, err
		}
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_LARGEOBJECT, list, nil
	case PARAMETER:
		p.advance()
		list, err := p.parseParameterNameList()
		if err != nil {
			return 0, 0, nil, err
		}
		return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_PARAMETER_ACL, list, nil
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
			if _, err := p.expect(WRAPPER); err != nil {
				return 0, 0, nil, err
			}
			names, err := p.parseNameList()
			if err != nil {
				return 0, 0, nil, err
			}
			return nodes.ACL_TARGET_OBJECT, nodes.OBJECT_FDW, makeNameListAsAnyNameList(names), nil
		}
		if _, err := p.expect(SERVER); err != nil {
			return 0, 0, nil, err
		}
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
			if p.collectMode() {
				p.addTokenCandidate(TABLES)
				p.addTokenCandidate(SEQUENCES)
				p.addTokenCandidate(FUNCTIONS)
				p.addTokenCandidate(PROCEDURES)
				p.addTokenCandidate(ROUTINES)
				return 0, 0, nil, errCollecting
			}
			return 0, 0, nil, p.syntaxErrorAtCur()
		}
		if _, err := p.expect(IN_P); err != nil {
			return 0, 0, nil, err
		}
		if _, err := p.expect(SCHEMA); err != nil {
			return 0, 0, nil, err
		}
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

// parseParameterNameList parses the configuration parameters of GRANT ON
// PARAMETER, each given as one string, as gram.y joins it.
//
//	parameter_name_list: parameter_name | parameter_name_list ',' parameter_name
//	parameter_name: ColId | parameter_name '.' ColId
func (p *Parser) parseParameterNameList() (*nodes.List, error) {
	result := &nodes.List{}
	for {
		name, err := p.parseColId()
		if err != nil {
			return nil, err
		}
		for p.cur.Type == '.' {
			p.advance()
			part, err := p.parseColId()
			if err != nil {
				return nil, err
			}
			name += "." + part
		}
		result.Items = append(result.Items, &nodes.String{Str: name})
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

// parsePrivileges parses privileges, ALL alone giving a nil list.
func (p *Parser) parsePrivileges() (*nodes.List, error) {
	privs, _, err := p.parsePrivilegesOrAll()
	return privs, err
}

// parsePrivilegesOrAll parses privileges and reports whether it began with ALL.
//
//	privileges: privilege_list
//	          | ALL | ALL PRIVILEGES
//	          | ALL '(' columnList ')' | ALL PRIVILEGES '(' columnList ')'
func (p *Parser) parsePrivilegesOrAll() (*nodes.List, bool, error) {
	if p.cur.Type != ALL {
		privs, err := p.parsePrivilegeList()
		return privs, false, err
	}
	loc := p.pos()
	p.advance()
	if p.cur.Type == PRIVILEGES {
		p.advance()
	}
	if p.cur.Type != '(' {
		return nil, true, nil
	}
	p.advance()
	cols, err := p.parseColumnList()
	if err != nil {
		return nil, true, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, true, err
	}
	all := &nodes.AccessPriv{Cols: cols, Loc: nodes.Loc{Start: loc, End: p.prev.End}}
	return &nodes.List{Items: []nodes.Node{all}}, true, nil
}

// parsePrivilegeList parses privilege_list: privilege | privilege_list ',' privilege
func (p *Parser) parsePrivilegeList() (*nodes.List, error) {
	result := &nodes.List{}
	for {
		priv, err := p.parsePrivilege()
		if err != nil {
			return nil, err
		}
		result.Items = append(result.Items, priv)
		if p.cur.Type != ',' {
			return result, nil
		}
		p.advance()
	}
}

// parsePrivilege parses privilege.
//
//	privilege: SELECT opt_column_list | REFERENCES opt_column_list
//	         | CREATE opt_column_list | ALTER SYSTEM_P | ColId opt_column_list
func (p *Parser) parsePrivilege() (*nodes.AccessPriv, error) {
	loc := p.pos()
	var privName string
	switch p.cur.Type {
	case SELECT, REFERENCES, CREATE:
		privName = strings.ToLower(p.cur.Str)
		p.advance()
	case ALTER:
		// ALTER is unreserved: alone it is a ColId, the name of a role.
		if p.peekNext().Type != SYSTEM_P {
			name, err := p.parseColId()
			if err != nil {
				return nil, err
			}
			privName = name
			break
		}
		p.advance()
		p.advance()
		return &nodes.AccessPriv{PrivName: "alter system", Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
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

func (p *Parser) parseOptGrantGrantOption() (bool, error) {
	if p.cur.Type != WITH || p.peekNext().Type != GRANT {
		return false, nil
	}
	p.advance()
	p.advance()
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

// parseGrantRoleOptList parses the options after WITH in GrantRoleStmt.
//
//	grant_role_opt_list: grant_role_opt | grant_role_opt_list ',' grant_role_opt
//	grant_role_opt: ColLabel grant_role_opt_value
//	grant_role_opt_value: OPTION | TRUE_P | FALSE_P
//
// Any label parses; the server names the unknown one when it runs.
func (p *Parser) parseGrantRoleOptList() (*nodes.List, error) {
	if p.cur.Type != WITH {
		return nil, nil
	}
	p.advance() // consume WITH
	result := &nodes.List{}
	for {
		name, err := p.parseColLabel()
		if err != nil {
			return nil, err
		}
		var val bool
		switch p.cur.Type {
		case OPTION, TRUE_P:
			val = true
		case FALSE_P:
			val = false
		default:
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		result.Items = append(result.Items, makeDefElem(name, &nodes.Boolean{Boolval: val}))
		if p.cur.Type != ',' {
			return result, nil
		}
		p.advance()
	}
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
			var err error
			elem, err = p.parseAlterOptRoleElem()
			if err != nil {
				return nil, err
			}
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

func (p *Parser) parseAlterOptRoleElem() (*nodes.DefElem, error) {
	expectSconst := func() (string, error) {
		if p.cur.Type != SCONST {
			return "", p.syntaxErrorAtCur()
		}
		return p.advance().Str, nil
	}
	if p.collectMode() {
		// Offer the role option names at the cursor.
		return p.parseRoleOptionIdent(), nil
	}
	start := p.pos()
	switch p.cur.Type {
	case PASSWORD:
		p.advance()
		if p.cur.Type == NULL_P {
			p.advance()
			return makeDefElem("password", nil), nil
		}
		pw, err := expectSconst()
		if err != nil {
			return nil, err
		}
		return makeDefElem("password", &nodes.String{Str: pw}), nil
	case ENCRYPTED, UNENCRYPTED:
		unencrypted := p.cur.Type == UNENCRYPTED
		p.advance()
		if _, err := p.expect(PASSWORD); err != nil {
			return nil, err
		}
		pw, err := expectSconst()
		if err != nil {
			return nil, err
		}
		if unencrypted {
			return nil, &ParseError{Message: "UNENCRYPTED PASSWORD is no longer supported", Position: start}
		}
		return makeDefElem("password", &nodes.String{Str: pw}), nil
	case INHERIT:
		p.advance()
		return makeDefElem("inherit", &nodes.Boolean{Boolval: true}), nil
	case CONNECTION:
		p.advance()
		if _, err := p.expect(LIMIT); err != nil {
			return nil, err
		}
		val, err := p.parseNumericOnlyRequired()
		if err != nil {
			return nil, err
		}
		limit, ok := val.(*nodes.Integer)
		if !ok {
			// SignedIconst: an integer.
			return nil, &ParseError{Message: "syntax error", Position: p.prev.Loc}
		}
		return makeDefElem("connectionlimit", limit), nil
	case VALID:
		p.advance()
		if _, err := p.expect(UNTIL); err != nil {
			return nil, err
		}
		until, err := expectSconst()
		if err != nil {
			return nil, err
		}
		return makeDefElem("validUntil", &nodes.String{Str: until}), nil
	case USER:
		// Supported but not documented for roles, for use by ALTER GROUP.
		p.advance()
		roles, err := p.parseRoleList()
		if err != nil {
			return nil, err
		}
		return makeDefElem("rolemembers", roles), nil
	case IDENT:
		elem := p.parseRoleOptionIdent()
		if elem == nil {
			return nil, &ParseError{Message: fmt.Sprintf("unrecognized role option \"%s\"", p.cur.Str), Position: start}
		}
		return elem, nil
	}
	return nil, nil
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
		if p.cur.Type != ICONST {
			return nil, p.syntaxErrorAtCur()
		}
		val := p.advance().Ival
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
		return p.parseAlterOptRoleElem()
	}
}

func (p *Parser) parseAlterRoleStmt() (nodes.Node, error) {
	loc := p.prev.Loc // position of ALTER keyword
	p.advance()       // consume ROLE or USER
	// ALTER ROLE ALL SET/RESET/IN DATABASE ...
	if p.cur.Type == ALL {
		p.advance()
		return p.parseAlterRoleSetStmtSuffix(loc, nil)
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
	if p.cur.Type == SET || p.cur.Type == RESET || p.cur.Type == IN_P {
		return p.parseAlterRoleSetStmtSuffix(loc, role)
	}
	p.parseGrantOptWith()
	options, err := p.parseOptRoleList(false)
	if err != nil {
		return nil, err
	}
	return &nodes.AlterRoleStmt{Role: role, Options: options, Action: 1,
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
	return &nodes.DropRoleStmt{Roles: roles, MissingOk: missingOk,
		Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
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
