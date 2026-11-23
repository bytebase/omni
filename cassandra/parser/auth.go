package parser

import (
	"strings"

	"github.com/bytebase/omni/cassandra/ast"
)

func (p *Parser) parseGrant() (ast.StmtNode, error) {
	start := p.curLoc()
	p.advance() // GRANT

	// GRANT role TO grantee (Parser.g:1105) when a role name is followed by
	// TO; a permission is always followed by ON or ','.
	if (isIdentLike(p.cur.Type) || isStringLiteral(p.cur.Type)) && p.peekNext().Type == tokTO {
		roleName, err := p.parseRoleName()
		if err != nil {
			return nil, err
		}
		if err := p.expectKeyword(tokTO); err != nil {
			return nil, err
		}
		grantee, err := p.parseRoleName()
		if err != nil {
			return nil, err
		}
		return &ast.GrantRoleStmt{
			RoleName: roleName,
			Grantee:  grantee,
			Loc:      p.makeLoc(start),
		}, nil
	}

	priv, privs, err := p.parsePermissionOrAll()
	if err != nil {
		return nil, err
	}

	if err := p.expectKeyword(tokON); err != nil {
		return nil, err
	}

	resourceTok := p.cur
	resource, err := p.parseResource()
	if err != nil {
		return nil, err
	}
	if err := p.checkApplicablePermissions(privs, resource, resourceTok); err != nil {
		return nil, err
	}

	if err := p.expectKeyword(tokTO); err != nil {
		return nil, err
	}

	role, err := p.parseRoleName()
	if err != nil {
		return nil, err
	}

	return &ast.GrantStmt{
		Privilege:  priv,
		Privileges: privs,
		Resource:   resource,
		Role:       role,
		Loc:        p.makeLoc(start),
	}, nil
}

func (p *Parser) parseRevoke() (ast.StmtNode, error) {
	start := p.curLoc()
	p.advance() // REVOKE

	// REVOKE role FROM revokee (Parser.g:1116).
	if (isIdentLike(p.cur.Type) || isStringLiteral(p.cur.Type)) && p.peekNext().Type == tokFROM {
		roleName, err := p.parseRoleName()
		if err != nil {
			return nil, err
		}
		if err := p.expectKeyword(tokFROM); err != nil {
			return nil, err
		}
		revokee, err := p.parseRoleName()
		if err != nil {
			return nil, err
		}
		return &ast.RevokeRoleStmt{
			RoleName: roleName,
			Revokee:  revokee,
			Loc:      p.makeLoc(start),
		}, nil
	}

	priv, privs, err := p.parsePermissionOrAll()
	if err != nil {
		return nil, err
	}

	if err := p.expectKeyword(tokON); err != nil {
		return nil, err
	}

	resourceTok := p.cur
	resource, err := p.parseResource()
	if err != nil {
		return nil, err
	}
	if err := p.checkApplicablePermissions(privs, resource, resourceTok); err != nil {
		return nil, err
	}

	if err := p.expectKeyword(tokFROM); err != nil {
		return nil, err
	}

	role, err := p.parseRoleName()
	if err != nil {
		return nil, err
	}

	return &ast.RevokeStmt{
		Privilege:  priv,
		Privileges: privs,
		Resource:   resource,
		Role:       role,
		Loc:        p.makeLoc(start),
	}, nil
}

// parseList parses LIST ROLES (Parser.g:1363), LIST USERS (Parser.g:1274)
// and listPermissionsStatement (Parser.g:1124):
//
//	LIST permissionOrAll [ON resource] [OF roleName] [NORECURSIVE]
func (p *Parser) parseList() (ast.StmtNode, error) {
	start := p.curLoc()
	p.advance() // LIST

	if p.cur.Type == tokROLES {
		return p.parseListRoles(start)
	}
	if p.curIs("users") {
		p.advance()
		return &ast.ListUsersStmt{Loc: p.makeLoc(start)}, nil
	}

	priv, privs, err := p.parsePermissionOrAll()
	if err != nil {
		return nil, err
	}

	stmt := &ast.ListPermissionsStmt{Privilege: priv, Privileges: privs}

	if p.match(tokON) {
		resource, err := p.parseResource()
		if err != nil {
			return nil, err
		}
		stmt.Resource = resource
	}

	if p.match(tokOF) {
		role, err := p.parseRoleName()
		if err != nil {
			return nil, err
		}
		stmt.Role = role
	}

	if p.match(tokNORECURSIVE) {
		stmt.NoRecursive = true
	}

	stmt.Loc = p.makeLoc(start)
	return stmt, nil
}

func (p *Parser) parseListRoles(start int) (*ast.ListRolesStmt, error) {
	p.advance() // ROLES

	stmt := &ast.ListRolesStmt{}

	if p.match(tokOF) {
		role, err := p.parseRoleName()
		if err != nil {
			return nil, err
		}
		stmt.Of = role
	}

	if p.cur.Type == tokNORECURSIVE {
		stmt.NoRecursive = true
		p.advance()
	}

	stmt.Loc = p.makeLoc(start)
	return stmt, nil
}

// parsePermissionOrAll parses permissionOrAll (Parser.g:1143):
//
//	ALL [PERMISSIONS] | permission [PERMISSION] (',' permission [PERMISSION])*
//
// It returns the permissions joined by ", " ("ALL" or "ALL PERMISSIONS" for
// the first form) and each permission.
func (p *Parser) parsePermissionOrAll() (string, []string, error) {
	if p.cur.Type == tokALL {
		p.advance()
		if p.match(tokPERMISSIONS) {
			return "ALL PERMISSIONS", []string{"ALL"}, nil
		}
		return "ALL", []string{"ALL"}, nil
	}
	var privs []string
	for {
		priv, err := p.parsePermission()
		if err != nil {
			return "", nil, err
		}
		privs = append(privs, priv)
		p.match(tokPERMISSION)
		if !p.match(tokCOMMA) {
			break
		}
	}
	return strings.Join(privs, ", "), privs, nil
}

// parsePermission parses permission (Parser.g:1138).
func (p *Parser) parsePermission() (string, error) {
	switch p.cur.Type {
	case tokCREATE, tokALTER, tokDROP, tokSELECT, tokMODIFY, tokAUTHORIZE, tokDESCRIBE, tokEXECUTE:
		priv := strings.ToUpper(p.cur.Str)
		p.advance()
		return priv, nil
	}
	if p.curIs("unmask") || p.curIs("select_masked") {
		priv := strings.ToUpper(p.cur.Str)
		p.advance()
		return priv, nil
	}
	return "", p.errorf("expected privilege keyword, got %s", p.tokenDesc())
}

// parseResource parses resource (Parser.g:1148): dataResource,
// roleResource, functionResource or jmxResource. ALL, ROLE, FUNCTION, MBEAN
// and MBEANS are unreserved keywords, so each is a table name when what
// follows it does not continue its own form (LIST ALL ON ROLE lists the
// permissions on a table named role).
func (p *Parser) parseResource() (*ast.Resource, error) {
	start := p.curLoc()
	next := p.peekNext()

	switch {
	case p.cur.Type == tokALL && (next.Type == tokKEYSPACES || next.Type == tokROLES || next.Type == tokMBEANS || next.Type == tokFUNCTIONS || isWord(next, "tables") || isWord(next, "columnfamilies")):
		p.advance() // ALL
		switch {
		case p.cur.Type == tokFUNCTIONS:
			p.advance()
			if p.match(tokIN) {
				ks, err := p.parseInKeyspace()
				if err != nil {
					return nil, err
				}
				return &ast.Resource{
					Type: "ALL FUNCTIONS IN KEYSPACE",
					Name: &ast.QualifiedName{Parts: []*ast.Identifier{ks}, Loc: ks.Loc},
					Loc:  p.makeLoc(start),
				}, nil
			}
			return &ast.Resource{Type: "ALL FUNCTIONS", Loc: p.makeLoc(start)}, nil
		case p.cur.Type == tokKEYSPACES:
			p.advance()
			return &ast.Resource{Type: "ALL KEYSPACES", Loc: p.makeLoc(start)}, nil
		case p.cur.Type == tokROLES:
			p.advance()
			return &ast.Resource{Type: "ALL ROLES", Loc: p.makeLoc(start)}, nil
		case p.cur.Type == tokMBEANS:
			p.advance()
			return &ast.Resource{Type: "ALL MBEANS", Loc: p.makeLoc(start)}, nil
		default: // TABLES IN KEYSPACE ks; K_TABLES is also COLUMNFAMILIES (Lexer.g:96)
			p.advance() // TABLES
			if err := p.expectKeyword(tokIN); err != nil {
				return nil, err
			}
			ks, err := p.parseInKeyspace()
			if err != nil {
				return nil, err
			}
			return &ast.Resource{
				Type: "ALL TABLES IN KEYSPACE",
				Name: &ast.QualifiedName{Parts: []*ast.Identifier{ks}, Loc: ks.Loc},
				Loc:  p.makeLoc(start),
			}, nil
		}

	case p.cur.Type == tokFUNCTION && (isCQLName(next) || isAllowedFunctionName(next)):
		// FUNCTION ks.fn(types): argument types are mandatory (Parser.g:1182)
		// and FunctionResource.functionFromCql unconditionally requires an
		// explicit keyspace qualifier during parsing (Parser.g:1191).
		p.advance()
		nameTok := p.cur
		name, err := p.parseFunctionQualifiedName()
		if err != nil {
			return nil, err
		}
		if len(name.Parts) < 2 {
			return nil, p.errorAt(nameTok, "function name in FUNCTION resource must be explicitly qualified by a keyspace")
		}
		res := &ast.Resource{Type: "FUNCTION", Name: name}
		if _, err := p.expect(tokLPAREN); err != nil {
			return nil, err
		}
		if p.cur.Type != tokRPAREN {
			argTypes, err := p.parseTypeList()
			if err != nil {
				return nil, err
			}
			res.ArgTypes = argTypes
		}
		if _, err := p.expect(tokRPAREN); err != nil {
			return nil, err
		}
		res.Loc = p.makeLoc(start)
		return res, nil

	case (p.cur.Type == tokMBEAN || p.cur.Type == tokMBEANS) && isStringLiteral(next.Type):
		// jmxResource (Parser.g:1162): the name must be a JMX object name,
		// and MBEANS takes only a pattern (canonicalizeObjectName, Parser.g:162).
		plural := p.cur.Type == tokMBEANS
		typ := "MBEAN"
		if plural {
			typ = "MBEANS"
		}
		p.advance()
		tok := p.cur
		if msg := checkJMXObjectName(tok.Str, plural); msg != "" {
			return nil, p.errorf("%s", msg)
		}
		p.advance()
		ident := &ast.Identifier{Name: tok.Str, Loc: ast.Loc{Start: tok.Loc, End: tok.End}}
		return &ast.Resource{
			Type: typ,
			Name: &ast.QualifiedName{Parts: []*ast.Identifier{ident}, Loc: ident.Loc},
			Loc:  p.makeLoc(start),
		}, nil

	case p.cur.Type == tokKEYSPACE:
		p.advance()
		ks, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		return &ast.Resource{
			Type: "KEYSPACE",
			Name: &ast.QualifiedName{Parts: []*ast.Identifier{ks}, Loc: ks.Loc},
			Loc:  p.makeLoc(start),
		}, nil

	// roleName (Parser.g:1519) is no reserved word, so GRANT ... ON ROLE TO
	// grants on a table named role.
	case p.cur.Type == tokROLE && (isCQLName(next) || isStringLiteral(next.Type)):
		p.advance()
		role, err := p.parseRoleName()
		if err != nil {
			return nil, err
		}
		return &ast.Resource{
			Type: "ROLE",
			Name: &ast.QualifiedName{Parts: []*ast.Identifier{role}, Loc: role.Loc},
			Loc:  p.makeLoc(start),
		}, nil

	case p.cur.Type == tokTABLE:
		p.advance()
		name, err := p.parseQualifiedName()
		if err != nil {
			return nil, err
		}
		return &ast.Resource{Type: "TABLE", Name: name, Loc: p.makeLoc(start)}, nil

	default:
		// Bare table reference: [keyspace.]table (TABLE keyword is optional)
		name, err := p.parseQualifiedName()
		if err != nil {
			return nil, err
		}
		return &ast.Resource{Type: "TABLE", Name: name, Loc: p.makeLoc(start)}, nil
	}
}

// applicablePermissions holds, for each resource a GRANT or REVOKE names,
// the permissions it takes. Parser.g:1086 and :1099 pass the permissions
// through filterPermissions (Parser.g:149), which keeps those the resource's
// applicablePermissions() holds and refuses the statement when none is left.
// Each set is what the pinned server answered for every permission on every
// resource. A FUNCTION resource has none here: its applicablePermissions()
// looks the function up, and fails on one that does not exist, so the
// server cannot be asked and omni takes every permission there.
var applicablePermissions = map[string]map[string]bool{
	"ALL KEYSPACES":             dataKeyspacePermissions,
	"KEYSPACE":                  dataKeyspacePermissions,
	"ALL TABLES IN KEYSPACE":    dataKeyspacePermissions,
	"TABLE":                     permissionSet("ALTER", "DROP", "SELECT", "MODIFY", "AUTHORIZE", "UNMASK", "SELECT_MASKED"),
	"ALL ROLES":                 permissionSet("CREATE", "ALTER", "DROP", "AUTHORIZE", "DESCRIBE"),
	"ROLE":                      permissionSet("ALTER", "DROP", "AUTHORIZE", "DESCRIBE"),
	"ALL FUNCTIONS":             functionKeyspacePermissions,
	"ALL FUNCTIONS IN KEYSPACE": functionKeyspacePermissions,
	"ALL MBEANS":                mbeanPermissions,
	"MBEAN":                     mbeanPermissions,
	"MBEANS":                    mbeanPermissions,
}

var (
	dataKeyspacePermissions     = permissionSet("CREATE", "ALTER", "DROP", "SELECT", "MODIFY", "AUTHORIZE", "UNMASK", "SELECT_MASKED")
	functionKeyspacePermissions = permissionSet("CREATE", "ALTER", "DROP", "AUTHORIZE", "EXECUTE")
	mbeanPermissions            = permissionSet("SELECT", "MODIFY", "AUTHORIZE", "DESCRIBE", "EXECUTE")
)

func permissionSet(perms ...string) map[string]bool {
	set := make(map[string]bool, len(perms))
	for _, perm := range perms {
		set[perm] = true
	}
	return set
}

// checkApplicablePermissions refuses a GRANT or REVOKE none of whose
// permissions the resource takes, as filterPermissions (Parser.g:149) does.
// ALL takes every permission, so it always leaves one.
func (p *Parser) checkApplicablePermissions(privs []string, res *ast.Resource, at Token) error {
	applicable, ok := applicablePermissions[res.Type]
	if !ok {
		return nil
	}
	for _, priv := range privs {
		if priv == "ALL" || applicable[priv] {
			return nil
		}
	}
	return p.errorAt(at, "Resource type %s does not support any of the requested permissions", res.Type)
}

// parseInKeyspace parses KEYSPACE ks after IN.
func (p *Parser) parseInKeyspace() (*ast.Identifier, error) {
	if err := p.expectKeyword(tokKEYSPACE); err != nil {
		return nil, err
	}
	return p.parseIdentifier()
}

func (p *Parser) parseCreateRole() (*ast.CreateRoleStmt, error) {
	start := p.curLoc()
	p.advance() // ROLE

	ifNotExists, err := p.parseIfNotExists()
	if err != nil {
		return nil, err
	}

	name, err := p.parseRoleName()
	if err != nil {
		return nil, err
	}

	stmt := &ast.CreateRoleStmt{
		IfNotExists: ifNotExists,
		Name:        name,
	}

	if p.match(tokWITH) {
		opts, err := p.parseRoleWithOptions()
		if err != nil {
			return nil, err
		}
		stmt.Options = opts
	}

	stmt.Loc = p.makeLoc(start)
	return stmt, nil
}

func (p *Parser) parseAlterRole() (*ast.AlterRoleStmt, error) {
	start := p.curLoc()
	p.advance() // ROLE

	ifExists := p.parseIfExists()

	name, err := p.parseRoleName()
	if err != nil {
		return nil, err
	}

	stmt := &ast.AlterRoleStmt{IfExists: ifExists, Name: name}

	if p.match(tokWITH) {
		opts, err := p.parseRoleWithOptions()
		if err != nil {
			return nil, err
		}
		stmt.Options = opts
	}

	stmt.Loc = p.makeLoc(start)
	return stmt, nil
}

func (p *Parser) parseDropRole() (*ast.DropRoleStmt, error) {
	start := p.curLoc()
	p.advance() // ROLE

	ifExists := p.parseIfExists()

	name, err := p.parseRoleName()
	if err != nil {
		return nil, err
	}

	return &ast.DropRoleStmt{
		IfExists: ifExists,
		Name:     name,
		Loc:      p.makeLoc(start),
	}, nil
}

// roleOptionGroup names the option a role option sets, for the duplicate
// check of parseRoleWithOptions: the ACCESS forms set the datacenter or CIDR
// permissions rather than an option.
var roleOptionGroup = map[string]string{
	"ACCESS TO ALL DATACENTERS": "dc all",
	"ACCESS TO DATACENTERS":     "dc list",
	"ACCESS FROM ALL CIDRS":     "cidr all",
	"ACCESS FROM CIDRS":         "cidr list",
}

// roleOptionConflict is the ACCESS form each ACCESS form may not be combined with.
var roleOptionConflict = map[string]string{
	"dc all":    "dc list",
	"dc list":   "dc all",
	"cidr all":  "cidr list",
	"cidr list": "cidr all",
}

// parseRoleWithOptions parses roleOptions (Parser.g:1374) and refuses what
// the actions refuse while parsing: an option given twice
// (RoleOptions.setOption, "Multiple definition for property"), ACCESS TO ALL
// DATACENTERS beside a datacenter list or ACCESS FROM ALL CIDRS beside a
// CIDR list (the DCPermissions and CIDRPermissions builders' preconditions),
// and PASSWORD beside HASHED PASSWORD (Parser.g:1312, Parser.g:1343).
func (p *Parser) parseRoleWithOptions() ([]*ast.RoleOption, error) {
	var opts []*ast.RoleOption
	seen := map[string]bool{}
	for {
		optStart := p.cur
		opt, err := p.parseRoleOption()
		if err != nil {
			return nil, err
		}
		if group, ok := roleOptionGroup[opt.Key]; ok {
			if seen[roleOptionConflict[group]] {
				return nil, p.errorAt(optStart, "%s cannot be combined with the other ACCESS %s form", opt.Key, strings.Fields(opt.Key)[1])
			}
			seen[group] = true
		} else {
			if seen[opt.Key] {
				return nil, p.errorAt(optStart, "Multiple definition for property '%s'", strings.ReplaceAll(opt.Key, " ", "_"))
			}
			seen[opt.Key] = true
		}
		if seen["PASSWORD"] && seen["HASHED PASSWORD"] {
			return nil, p.errorAt(optStart, "Options 'password' and 'hashed password' are mutually exclusive")
		}
		opts = append(opts, opt)
		if !p.match(tokAND) {
			break
		}
	}
	return opts, nil
}

// parseRoleOption parses roleOption (Parser.g:1378).
func (p *Parser) parseRoleOption() (*ast.RoleOption, error) {
	start := p.curLoc()
	tok := p.cur
	key := strings.ToUpper(tok.Str)

	switch tok.Type {
	case tokHASHED:
		p.advance() // HASHED
		if err := p.expectKeyword(tokPASSWORD); err != nil {
			return nil, err
		}
		if _, err := p.expect(tokEQ); err != nil {
			return nil, err
		}
		val, err := p.parseStringLit()
		if err != nil {
			return nil, err
		}
		return &ast.RoleOption{Key: "HASHED PASSWORD", Value: val, Loc: p.makeLoc(start)}, nil

	case tokPASSWORD:
		p.advance()
		if _, err := p.expect(tokEQ); err != nil {
			return nil, err
		}
		val, err := p.parseStringLit()
		if err != nil {
			return nil, err
		}
		return &ast.RoleOption{Key: key, Value: val, Loc: p.makeLoc(start)}, nil

	case tokACCESS:
		p.advance() // ACCESS
		// ACCESS TO [ALL] DATACENTERS or ACCESS FROM [ALL] CIDRS.
		var noun, toFrom string
		switch {
		case p.match(tokTO):
			noun, toFrom = "DATACENTERS", "TO"
		case p.match(tokFROM):
			noun, toFrom = "CIDRS", "FROM"
		default:
			return nil, p.errorf("expected TO or FROM after ACCESS, got %s", p.tokenDesc())
		}
		all := p.match(tokALL)
		if !p.curIs(strings.ToLower(noun)) {
			return nil, p.errorf("expected %s, got %s", noun, p.tokenDesc())
		}
		p.advance()
		if all {
			return &ast.RoleOption{Key: "ACCESS " + toFrom + " ALL " + noun, Loc: p.makeLoc(start)}, nil
		}
		val, err := p.parseStringSet()
		if err != nil {
			return nil, err
		}
		return &ast.RoleOption{Key: "ACCESS " + toFrom + " " + noun, Value: val, Loc: p.makeLoc(start)}, nil

	case tokLOGIN, tokSUPERUSER:
		p.advance()
		if _, err := p.expect(tokEQ); err != nil {
			return nil, err
		}
		val, err := p.parseBoolLit()
		if err != nil {
			return nil, err
		}
		return &ast.RoleOption{Key: key, Value: val, Loc: p.makeLoc(start)}, nil

	case tokOPTIONS:
		p.advance()
		if _, err := p.expect(tokEQ); err != nil {
			return nil, err
		}
		val, err := p.parseOptionHash()
		if err != nil {
			return nil, err
		}
		return &ast.RoleOption{Key: key, Value: val, Loc: p.makeLoc(start)}, nil

	default:
		return nil, p.errorf("expected PASSWORD, HASHED PASSWORD, LOGIN, SUPERUSER, OPTIONS or ACCESS, got %s", p.tokenDesc())
	}
}

// parseStringLit parses a STRING_LITERAL, quoted or dollar-quoted
// (isStringLiteral).
func (p *Parser) parseStringLit() (*ast.StringLit, error) {
	tok := p.cur
	if !isStringLiteral(tok.Type) {
		return nil, p.errorf("expected string, got %s", p.tokenDesc())
	}
	p.advance()
	return &ast.StringLit{Val: tok.Str, Loc: ast.Loc{Start: tok.Loc, End: tok.End}}, nil
}

// parseStringSet parses '{' STRING_LITERAL (',' STRING_LITERAL)* '}', the
// datacenter and CIDR lists of ACCESS (Parser.g:1385, Parser.g:1387).
func (p *Parser) parseStringSet() (*ast.SetLit, error) {
	start := p.curLoc()
	if _, err := p.expect(tokLBRACE); err != nil {
		return nil, err
	}
	set := &ast.SetLit{}
	for {
		s, err := p.parseStringLit()
		if err != nil {
			return nil, err
		}
		set.Elements = append(set.Elements, s)
		if !p.match(tokCOMMA) {
			break
		}
	}
	if _, err := p.expect(tokRBRACE); err != nil {
		return nil, err
	}
	set.Loc = p.makeLoc(start)
	return set, nil
}

// parseCreateUser parses createUserStatement (Parser.g:1197).
func (p *Parser) parseCreateUser() (*ast.CreateUserStmt, error) {
	start := p.curLoc()
	p.advance() // USER

	ifNotExists, err := p.parseIfNotExists()
	if err != nil {
		return nil, err
	}

	name, err := p.parseUsername("user name")
	if err != nil {
		return nil, err
	}

	stmt := &ast.CreateUserStmt{
		IfNotExists: ifNotExists,
		Name:        name,
	}

	if p.cur.Type == tokWITH {
		p.advance()
		if err := p.parseUserPasswordClause(stmt, nil); err != nil {
			return nil, err
		}
	}

	if p.cur.Type == tokSUPERUSER {
		v := true
		stmt.Superuser = &v
		p.advance()
	} else if p.cur.Type == tokNOSUPERUSER {
		v := false
		stmt.Superuser = &v
		p.advance()
	}

	stmt.Loc = p.makeLoc(start)
	return stmt, nil
}

// parseUserPasswordClause parses userPassword (Parser.g:1399), [HASHED]
// PASSWORD string, after WITH.
func (p *Parser) parseUserPasswordClause(create *ast.CreateUserStmt, alter *ast.AlterUserStmt) error {
	hashed := false
	if p.cur.Type == tokHASHED {
		hashed = true
		p.advance()
	}
	if err := p.expectKeyword(tokPASSWORD); err != nil {
		return err
	}
	pwd, err := p.parseStringLit()
	if err != nil {
		return err
	}
	if create != nil {
		create.Password = pwd
		create.Hashed = hashed
	}
	if alter != nil {
		alter.Password = pwd
		alter.Hashed = hashed
	}
	return nil
}

// parseAlterUser parses alterUserStatement (Parser.g:1219).
func (p *Parser) parseAlterUser() (*ast.AlterUserStmt, error) {
	start := p.curLoc()
	p.advance() // USER

	ifExists := p.parseIfExists()

	name, err := p.parseUsername("user name")
	if err != nil {
		return nil, err
	}

	stmt := &ast.AlterUserStmt{IfExists: ifExists, Name: name}

	if p.cur.Type == tokWITH {
		p.advance()
		if err := p.parseUserPasswordClause(nil, stmt); err != nil {
			return nil, err
		}
	}

	if p.cur.Type == tokSUPERUSER {
		v := true
		stmt.Superuser = &v
		p.advance()
	} else if p.cur.Type == tokNOSUPERUSER {
		v := false
		stmt.Superuser = &v
		p.advance()
	}

	stmt.Loc = p.makeLoc(start)
	return stmt, nil
}

// parseDropUser parses dropUserStatement (Parser.g:1241).
func (p *Parser) parseDropUser() (*ast.DropUserStmt, error) {
	start := p.curLoc()
	p.advance() // USER

	ifExists := p.parseIfExists()

	name, err := p.parseUsername("user name")
	if err != nil {
		return nil, err
	}

	return &ast.DropUserStmt{
		IfExists: ifExists,
		Name:     name,
		Loc:      p.makeLoc(start),
	}, nil
}

// parseAddIdentity parses addIdentityStatement (Parser.g:1251):
// ADD IDENTITY [IF NOT EXISTS] identity TO ROLE identity.
func (p *Parser) parseAddIdentity() (*ast.AddIdentityStmt, error) {
	start := p.curLoc()
	p.advance() // ADD
	if !p.curIs("identity") {
		return nil, p.errorf("expected IDENTITY after ADD, got %s", p.tokenDesc())
	}
	p.advance()

	ifNotExists, err := p.parseIfNotExists()
	if err != nil {
		return nil, err
	}
	identity, err := p.parseUsername("identity")
	if err != nil {
		return nil, err
	}
	if err := p.expectKeyword(tokTO); err != nil {
		return nil, err
	}
	if err := p.expectKeyword(tokROLE); err != nil {
		return nil, err
	}
	role, err := p.parseUsername("role")
	if err != nil {
		return nil, err
	}
	return &ast.AddIdentityStmt{
		IfNotExists: ifNotExists,
		Identity:    identity,
		Role:        role,
		Loc:         p.makeLoc(start),
	}, nil
}

// parseDropIdentity parses dropIdentityStatement (Parser.g:1263):
// DROP IDENTITY [IF EXISTS] identity, with DROP already consumed.
func (p *Parser) parseDropIdentity(start int) (*ast.DropIdentityStmt, error) {
	p.advance() // IDENTITY
	ifExists := p.parseIfExists()
	identity, err := p.parseUsername("identity")
	if err != nil {
		return nil, err
	}
	return &ast.DropIdentityStmt{
		IfExists: ifExists,
		Identity: identity,
		Loc:      p.makeLoc(start),
	}, nil
}
