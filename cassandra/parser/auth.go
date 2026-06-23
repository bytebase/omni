package parser

import (
	"strings"

	"github.com/bytebase/omni/cassandra/ast"
)

func (p *Parser) parseGrant() (ast.StmtNode, error) {
	start := p.curLoc()
	p.advance() // GRANT

	// Check for GRANT role_name TO grantee (role grant).
	// If the current token is identifier-like and the next token is TO,
	// and the current token is NOT a privilege keyword that would be followed by ON,
	// then parse as role grant.
	if isIdentLike(p.cur.Type) && p.peekNext().Type == tokTO {
		roleName, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		if err := p.expectKeyword(tokTO); err != nil {
			return nil, err
		}
		grantee, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		return &ast.GrantRoleStmt{
			RoleName: roleName,
			Grantee:  grantee,
			Loc:      p.makeLoc(start),
		}, nil
	}

	priv, err := p.parsePrivilege()
	if err != nil {
		return nil, err
	}

	if err := p.expectKeyword(tokON); err != nil {
		return nil, err
	}

	resource, err := p.parseResource()
	if err != nil {
		return nil, err
	}

	if err := p.expectKeyword(tokTO); err != nil {
		return nil, err
	}

	role, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	return &ast.GrantStmt{
		Privilege: priv,
		Resource:  resource,
		Role:      role,
		Loc:       p.makeLoc(start),
	}, nil
}

func (p *Parser) parseRevoke() (ast.StmtNode, error) {
	start := p.curLoc()
	p.advance() // REVOKE

	// Check for REVOKE role_name FROM revokee (role revoke).
	// If the current token is identifier-like and the next token is FROM,
	// then parse as role revoke.
	if isIdentLike(p.cur.Type) && p.peekNext().Type == tokFROM {
		roleName, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		if err := p.expectKeyword(tokFROM); err != nil {
			return nil, err
		}
		revokee, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		return &ast.RevokeRoleStmt{
			RoleName: roleName,
			Revokee:  revokee,
			Loc:      p.makeLoc(start),
		}, nil
	}

	priv, err := p.parsePrivilege()
	if err != nil {
		return nil, err
	}

	if err := p.expectKeyword(tokON); err != nil {
		return nil, err
	}

	resource, err := p.parseResource()
	if err != nil {
		return nil, err
	}

	if err := p.expectKeyword(tokFROM); err != nil {
		return nil, err
	}

	role, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	return &ast.RevokeStmt{
		Privilege: priv,
		Resource:  resource,
		Role:      role,
		Loc:       p.makeLoc(start),
	}, nil
}

func (p *Parser) parseList() (ast.StmtNode, error) {
	start := p.curLoc()
	p.advance() // LIST

	if p.cur.Type == tokROLES {
		return p.parseListRoles(start)
	}

	// LIST permissions [ON resource] [OF role]
	priv, err := p.parsePrivilege()
	if err != nil {
		return nil, err
	}

	stmt := &ast.ListPermissionsStmt{Privilege: priv}

	if p.match(tokON) {
		resource, err := p.parseResource()
		if err != nil {
			return nil, err
		}
		stmt.Resource = resource
	}

	if p.match(tokOF) {
		role, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Role = role
	}

	stmt.Loc = p.makeLoc(start)
	return stmt, nil
}

func (p *Parser) parseListRoles(start int) (*ast.ListRolesStmt, error) {
	p.advance() // ROLES

	stmt := &ast.ListRolesStmt{}

	if p.match(tokOF) {
		role, err := p.parseIdentifier()
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

func (p *Parser) parsePrivilege() (string, error) {
	tok := p.cur
	switch tok.Type {
	case tokALL:
		p.advance()
		if p.cur.Type == tokPERMISSIONS {
			p.advance()
			return "ALL PERMISSIONS", nil
		}
		return "ALL", nil
	case tokALTER:
		p.advance()
		return "ALTER", nil
	case tokAUTHORIZE:
		p.advance()
		return "AUTHORIZE", nil
	case tokDESCRIBE:
		p.advance()
		return "DESCRIBE", nil
	case tokEXECUTE:
		p.advance()
		return "EXECUTE", nil
	case tokCREATE:
		p.advance()
		return "CREATE", nil
	case tokDROP:
		p.advance()
		return "DROP", nil
	case tokMODIFY:
		p.advance()
		return "MODIFY", nil
	case tokSELECT:
		p.advance()
		return "SELECT", nil
	default:
		return "", p.errorf("expected privilege keyword, got %s", p.tokenDesc())
	}
}

func (p *Parser) parseResource() (*ast.Resource, error) {
	start := p.curLoc()

	switch p.cur.Type {
	case tokALL:
		p.advance()
		switch p.cur.Type {
		case tokFUNCTIONS:
			p.advance()
			if p.match(tokIN) {
				if err := p.expectKeyword(tokKEYSPACE); err != nil {
					return nil, err
				}
				ks, err := p.parseIdentifier()
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
		case tokKEYSPACES:
			p.advance()
			return &ast.Resource{Type: "ALL KEYSPACES", Loc: p.makeLoc(start)}, nil
		case tokROLES:
			p.advance()
			return &ast.Resource{Type: "ALL ROLES", Loc: p.makeLoc(start)}, nil
		case tokMBEANS:
			p.advance()
			return &ast.Resource{Type: "ALL MBEANS", Loc: p.makeLoc(start)}, nil
		default:
			return nil, p.errorf("expected FUNCTIONS, KEYSPACES, ROLES, or MBEANS after ALL, got %s", p.tokenDesc())
		}

	case tokFUNCTION:
		p.advance()
		name, err := p.parseQualifiedName()
		if err != nil {
			return nil, err
		}
		res := &ast.Resource{Type: "FUNCTION", Name: name, Loc: p.makeLoc(start)}
		if p.cur.Type == tokLPAREN {
			p.advance()
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
		}
		return res, nil

	case tokMBEAN:
		p.advance()
		val, err := p.parseConstant()
		if err != nil {
			return nil, err
		}
		ident := &ast.Identifier{Name: val.(*ast.StringLit).Val, Loc: val.GetLoc()}
		return &ast.Resource{
			Type: "MBEAN",
			Name: &ast.QualifiedName{Parts: []*ast.Identifier{ident}, Loc: ident.Loc},
			Loc:  p.makeLoc(start),
		}, nil

	case tokMBEANS:
		p.advance()
		val, err := p.parseConstant()
		if err != nil {
			return nil, err
		}
		ident := &ast.Identifier{Name: val.(*ast.StringLit).Val, Loc: val.GetLoc()}
		return &ast.Resource{
			Type: "MBEANS",
			Name: &ast.QualifiedName{Parts: []*ast.Identifier{ident}, Loc: ident.Loc},
			Loc:  p.makeLoc(start),
		}, nil

	case tokKEYSPACE:
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

	case tokROLE:
		p.advance()
		role, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		return &ast.Resource{
			Type: "ROLE",
			Name: &ast.QualifiedName{Parts: []*ast.Identifier{role}, Loc: role.Loc},
			Loc:  p.makeLoc(start),
		}, nil

	case tokTABLE:
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

func (p *Parser) parseCreateRole() (*ast.CreateRoleStmt, error) {
	start := p.curLoc()
	p.advance() // ROLE

	ifNotExists, err := p.parseIfNotExists()
	if err != nil {
		return nil, err
	}

	name, err := p.parseIdentifier()
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

	name, err := p.parseIdentifier()
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

	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	return &ast.DropRoleStmt{
		IfExists: ifExists,
		Name:     name,
		Loc:      p.makeLoc(start),
	}, nil
}

func (p *Parser) parseRoleWithOptions() ([]*ast.RoleOption, error) {
	var opts []*ast.RoleOption
	for {
		opt, err := p.parseRoleOption()
		if err != nil {
			return nil, err
		}
		opts = append(opts, opt)
		if !p.match(tokAND) {
			break
		}
	}
	return opts, nil
}

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
		val, err := p.parseConstant()
		if err != nil {
			return nil, err
		}
		return &ast.RoleOption{Key: "HASHED PASSWORD", Value: val, Loc: p.makeLoc(start)}, nil

	case tokPASSWORD:
		p.advance()
		if _, err := p.expect(tokEQ); err != nil {
			return nil, err
		}
		val, err := p.parseConstant()
		if err != nil {
			return nil, err
		}
		return &ast.RoleOption{Key: key, Value: val, Loc: p.makeLoc(start)}, nil

	case tokACCESS:
		p.advance() // ACCESS
		if err := p.expectKeyword(tokTO); err != nil {
			return nil, err
		}
		if p.cur.Type == tokALL {
			p.advance()
			if err := p.expectKeyword(tokDATACENTERS); err != nil {
				return nil, err
			}
			return &ast.RoleOption{Key: "ACCESS TO ALL DATACENTERS", Loc: p.makeLoc(start)}, nil
		}
		if err := p.expectKeyword(tokDATACENTERS); err != nil {
			return nil, err
		}
		// Parse set of datacenter names: { 'dc1', 'dc2' }
		val, err := p.parseCollectionLiteral()
		if err != nil {
			return nil, err
		}
		return &ast.RoleOption{Key: "ACCESS TO DATACENTERS", Value: val, Loc: p.makeLoc(start)}, nil

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
		return nil, p.errorf("expected PASSWORD, LOGIN, SUPERUSER, or OPTIONS, got %s", p.tokenDesc())
	}
}

func (p *Parser) parseCreateUser() (*ast.CreateUserStmt, error) {
	start := p.curLoc()
	p.advance() // USER

	ifNotExists, err := p.parseIfNotExists()
	if err != nil {
		return nil, err
	}

	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	stmt := &ast.CreateUserStmt{
		IfNotExists: ifNotExists,
		Name:        name,
	}

	if err := p.expectKeyword(tokWITH); err != nil {
		return nil, err
	}
	if err := p.expectKeyword(tokPASSWORD); err != nil {
		return nil, err
	}
	pwd, err := p.parseConstant()
	if err != nil {
		return nil, err
	}
	stmt.Password = pwd

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

func (p *Parser) parseAlterUser() (*ast.AlterUserStmt, error) {
	start := p.curLoc()
	p.advance() // USER

	ifExists := p.parseIfExists()

	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	stmt := &ast.AlterUserStmt{IfExists: ifExists, Name: name}

	if err := p.expectKeyword(tokWITH); err != nil {
		return nil, err
	}
	if err := p.expectKeyword(tokPASSWORD); err != nil {
		return nil, err
	}
	pwd, err := p.parseConstant()
	if err != nil {
		return nil, err
	}
	stmt.Password = pwd

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

func (p *Parser) parseDropUser() (*ast.DropUserStmt, error) {
	start := p.curLoc()
	p.advance() // USER

	ifExists := p.parseIfExists()

	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	return &ast.DropUserStmt{
		IfExists: ifExists,
		Name:     name,
		Loc:      p.makeLoc(start),
	}, nil
}
