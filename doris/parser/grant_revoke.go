package parser

import (
	"strings"

	"github.com/bytebase/omni/doris/ast"
)

// ---------------------------------------------------------------------------
// GRANT
// ---------------------------------------------------------------------------

// parseGrant parses:
//
//	GRANT priv_list ON [object_type] object TO {USER 'user'@'host' | ROLE 'role'} [WITH GRANT OPTION]
//	GRANT ROLE role_list TO 'user'@'host' | ROLE 'role'
//
// On entry, cur is the token AFTER kwGRANT has been consumed.
func (p *Parser) parseGrant(startLoc ast.Loc) (ast.Node, error) {
	stmt := &ast.GrantStmt{}
	endLoc := startLoc

	// Detect the two forms:
	//   GRANT ROLE role_list TO ...          (explicit role-grant form)
	//   GRANT 'role1', 'role2' TO ...        (implicit role-grant: string literal grantees)
	//   GRANT priv_list ON ...               (privilege-grant form)
	isRoleGrant := p.cur.Kind == kwROLE || p.cur.Kind == tokString

	if isRoleGrant {
		// Role-grant form: GRANT [ROLE] 'role1'[,'role2'...] TO user|role
		if p.cur.Kind == kwROLE {
			p.advance() // consume ROLE keyword
		}
		roles, loc, err := p.parseGranteeList()
		if err != nil {
			return nil, err
		}
		stmt.Roles = roles
		endLoc = loc

		if _, err := p.expect(kwTO); err != nil {
			return nil, err
		}
		grantees, granteeType, loc2, err := p.parseGranteeSpec()
		if err != nil {
			return nil, err
		}
		stmt.Grantees = grantees
		stmt.ToType = granteeType
		endLoc = loc2
	} else {
		// Privilege-grant form: GRANT priv_list ON [object_type] object TO ...
		privs, loc, err := p.parsePrivilegeList()
		if err != nil {
			return nil, err
		}
		stmt.Privileges = privs
		endLoc = loc

		if _, err := p.expect(kwON); err != nil {
			return nil, err
		}

		objType, obj, objLoc, err := p.parseObjectTypeAndName()
		if err != nil {
			return nil, err
		}
		stmt.ObjectType = objType
		stmt.Object = obj
		endLoc = objLoc

		if _, err := p.expect(kwTO); err != nil {
			return nil, err
		}

		grantees, granteeType, loc2, err := p.parseGranteeSpec()
		if err != nil {
			return nil, err
		}
		stmt.Grantees = grantees
		stmt.ToType = granteeType
		endLoc = loc2

		// Optional WITH GRANT OPTION
		if p.cur.Kind == kwWITH {
			p.advance()
			if p.cur.Kind == kwGRANT {
				p.advance()
				// OPTION is likely a non-reserved keyword parsed as identifier
				if p.cur.Kind == tokIdent && strings.ToUpper(p.cur.Str) == "OPTION" {
					endLoc = p.cur.Loc
					p.advance()
				}
			}
			stmt.WithGrantOption = true
		}
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// ---------------------------------------------------------------------------
// REVOKE
// ---------------------------------------------------------------------------

// parseRevoke parses:
//
//	REVOKE priv_list ON [object_type] object FROM {USER 'user'@'host' | ROLE 'role'} [, ...]
//	REVOKE ROLE role_list FROM 'user'@'host' | ROLE 'role'
//
// On entry, cur is the token AFTER kwREVOKE has been consumed.
func (p *Parser) parseRevoke(startLoc ast.Loc) (ast.Node, error) {
	stmt := &ast.RevokeStmt{}
	endLoc := startLoc

	isRoleRevoke := p.cur.Kind == kwROLE || p.cur.Kind == tokString

	if isRoleRevoke {
		// Role-revoke form: REVOKE [ROLE] 'role1'[,'role2'...] FROM user|role
		if p.cur.Kind == kwROLE {
			p.advance() // consume ROLE keyword
		}
		roles, loc, err := p.parseGranteeList()
		if err != nil {
			return nil, err
		}
		stmt.Roles = roles
		endLoc = loc

		if _, err := p.expect(kwFROM); err != nil {
			return nil, err
		}
		revokees, revokeeType, loc2, err := p.parseGranteeSpec()
		if err != nil {
			return nil, err
		}
		stmt.Revokees = revokees
		stmt.FromType = revokeeType
		endLoc = loc2
	} else {
		// Privilege-revoke form
		privs, loc, err := p.parsePrivilegeList()
		if err != nil {
			return nil, err
		}
		stmt.Privileges = privs
		endLoc = loc

		if _, err := p.expect(kwON); err != nil {
			return nil, err
		}

		objType, obj, objLoc, err := p.parseObjectTypeAndName()
		if err != nil {
			return nil, err
		}
		stmt.ObjectType = objType
		stmt.Object = obj
		endLoc = objLoc

		if _, err := p.expect(kwFROM); err != nil {
			return nil, err
		}

		revokees, revokeeType, loc2, err := p.parseGranteeSpec()
		if err != nil {
			return nil, err
		}
		stmt.Revokees = revokees
		stmt.FromType = revokeeType
		endLoc = loc2
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// parsePrivilegeList parses a comma-separated list of privilege names.
// Privilege names are identifiers like SELECT_PRIV, USAGE_PRIV, ALL, etc.
// Optional trailing PRIVILEGES keyword after ALL is consumed.
//
// Returns the list of privilege names (uppercased) and the end location.
func (p *Parser) parsePrivilegeList() ([]string, ast.Loc, error) {
	var privs []string
	endLoc := p.cur.Loc

	for {
		if !isPrivilegeToken(p.cur.Kind) {
			return nil, ast.NoLoc(), p.syntaxErrorAtCur()
		}
		priv := strings.ToUpper(p.cur.Str)
		endLoc = p.cur.Loc
		p.advance()

		// ALL [PRIVILEGES] — consume optional PRIVILEGES
		if strings.ToUpper(priv) == "ALL" && p.cur.Kind == kwPRIVILEGES {
			priv = "ALL PRIVILEGES"
			endLoc = p.cur.Loc
			p.advance()
		}

		privs = append(privs, priv)

		if p.cur.Kind != int(',') {
			break
		}
		p.advance() // consume ','
	}

	return privs, endLoc, nil
}

// isPrivilegeToken reports whether the current token can be a privilege name.
// Privileges are identifiers (SELECT_PRIV, USAGE_PRIV, LOAD_PRIV, …) or
// certain reserved keywords used as privilege names (ALL, SELECT, INSERT,
// UPDATE, DELETE, CREATE, DROP, ALTER, LOAD, ADMIN, SHOW, GRANT, USAGE).
func isPrivilegeToken(kind TokenKind) bool {
	switch kind {
	case kwALL, kwSELECT, kwINSERT, kwUPDATE, kwDELETE,
		kwCREATE, kwDROP, kwALTER, kwLOAD, kwADMIN, kwSHOW, kwGRANT:
		return true
	case tokIdent, tokQuotedIdent:
		return true
	default:
		// non-reserved keywords may be privilege names too
		if kind >= 700 && !IsReserved(kind) {
			return true
		}
		return false
	}
}

// parseObjectTypeAndName parses the optional object-type keyword(s) followed
// by the object name. Returns (objectType, objectName, endLoc, error).
//
// Object types:
//
//	TABLE (default when no keyword)
//	DATABASE
//	RESOURCE
//	CLUSTER
//	COMPUTE GROUP
//	STAGE
//	STORAGE VAULT
//	WORKLOAD GROUP
//
// Object names can be:
//   - *.* (three-part: catalog.db.table)
//   - db.* or db.table (two-part)
//   - * or name (single-part)
//   - A quoted string (for RESOURCE, STORAGE VAULT, COMPUTE GROUP, WORKLOAD GROUP names)
func (p *Parser) parseObjectTypeAndName() (string, *ast.ObjectName, ast.Loc, error) {
	objType := "TABLE"

	switch p.cur.Kind {
	case kwDATABASE:
		objType = "DATABASE"
		p.advance()
	case kwRESOURCE:
		objType = "RESOURCE"
		p.advance()
	case kwCLUSTER:
		objType = "CLUSTER"
		p.advance()
	case kwSTAGE:
		objType = "STAGE"
		p.advance()
	case kwCOMPUTE:
		// COMPUTE GROUP
		p.advance()
		if p.cur.Kind == kwGROUP {
			p.advance()
		}
		objType = "COMPUTE GROUP"
	case kwSTORAGE:
		// STORAGE VAULT
		p.advance()
		if p.cur.Kind == kwVAULT {
			p.advance()
		}
		objType = "STORAGE VAULT"
	case kwWORKLOAD:
		// WORKLOAD GROUP
		p.advance()
		if p.cur.Kind == kwGROUP {
			p.advance()
		}
		objType = "WORKLOAD GROUP"
	}

	// Parse object name: may be an identifier path (*.* / db.* / name) or
	// a quoted string (for non-table resources).
	obj, objLoc, err := p.parseGrantObject()
	if err != nil {
		return "", nil, ast.NoLoc(), err
	}

	return objType, obj, objLoc, nil
}

// parseGrantObject parses the object name in GRANT/REVOKE ON <object>.
// Handles the forms:
//   - *.*.* (three wildcards)
//   - *.* (two wildcards / catalog.db)
//   - * (single wildcard)
//   - db.table  / db.*
//   - 'name' (quoted string, for RESOURCE / WORKLOAD GROUP / COMPUTE GROUP / STORAGE VAULT)
func (p *Parser) parseGrantObject() (*ast.ObjectName, ast.Loc, error) {
	startLoc := p.cur.Loc

	// Quoted string object name (e.g., RESOURCE 'spark_resource')
	if p.cur.Kind == tokString {
		name := p.cur.Str
		loc := p.cur.Loc
		p.advance()
		obj := &ast.ObjectName{Parts: []string{name}, Loc: loc}
		return obj, loc, nil
	}

	// Wildcard or identifier path
	var parts []string
	var endLoc ast.Loc

	part, loc, err := p.parseGrantObjectPart()
	if err != nil {
		return nil, ast.NoLoc(), err
	}
	parts = append(parts, part)
	endLoc = loc

	for p.cur.Kind == int('.') {
		p.advance() // consume '.'
		next, nextLoc, err2 := p.parseGrantObjectPart()
		if err2 != nil {
			return nil, ast.NoLoc(), err2
		}
		parts = append(parts, next)
		endLoc = nextLoc
	}

	obj := &ast.ObjectName{
		Parts: parts,
		Loc:   ast.Loc{Start: startLoc.Start, End: endLoc.End},
	}
	return obj, endLoc, nil
}

// parseGrantObjectPart parses one segment of a grant object name:
// either '*' (wildcard) or an identifier/keyword.
func (p *Parser) parseGrantObjectPart() (string, ast.Loc, error) {
	if p.cur.Kind == int('*') {
		loc := p.cur.Loc
		p.advance()
		return "*", loc, nil
	}
	// identifier or keyword used as identifier
	return p.parseIdentifierQualified()
}

// parseGranteeList parses a comma-separated list of grantee identifiers
// (role names or user strings). Used for GRANT ROLE list and REVOKE ROLE list.
func (p *Parser) parseGranteeList() ([]string, ast.Loc, error) {
	var names []string
	endLoc := p.cur.Loc

	for {
		name, loc, err := p.parseIdentifierOrString()
		if err != nil {
			return nil, ast.NoLoc(), err
		}
		names = append(names, name)
		endLoc = loc

		if p.cur.Kind != int(',') {
			break
		}
		p.advance() // consume ','
	}

	return names, endLoc, nil
}

// parseGranteeSpec parses the grantee / revokee clause after TO / FROM.
// Handles:
//
//	USER 'user'@'host'
//	ROLE 'role'
//	'user'@'host'   (no keyword)
//
// Returns (grantees, granteeType, endLoc, error).
// granteeType is "USER", "ROLE", or "" when no keyword is present.
//
// Multiple grantees separated by commas are supported.
func (p *Parser) parseGranteeSpec() ([]string, string, ast.Loc, error) {
	var granteeType string
	endLoc := p.cur.Loc

	// Optional USER or ROLE keyword before first grantee
	switch p.cur.Kind {
	case kwUSER:
		granteeType = "USER"
		p.advance()
	case kwROLE:
		granteeType = "ROLE"
		p.advance()
	}

	var grantees []string

	for {
		raw, loc, err := p.parseRawGrantee()
		if err != nil {
			return nil, "", ast.NoLoc(), err
		}
		grantees = append(grantees, raw)
		endLoc = loc

		if p.cur.Kind != int(',') {
			break
		}
		p.advance() // consume ','

		// After the first grantee, subsequent items may have USER/ROLE prefix too.
		switch p.cur.Kind {
		case kwUSER:
			p.advance()
		case kwROLE:
			p.advance()
		}
	}

	return grantees, granteeType, endLoc, nil
}

// parseRawGrantee parses a single grantee as a raw string in 'user'@'host' form
// or just 'name'. Returns the concatenated raw string and end location.
func (p *Parser) parseRawGrantee() (string, ast.Loc, error) {
	name, nameLoc, err := p.parseIdentifierOrString()
	if err != nil {
		return "", ast.NoLoc(), err
	}
	endLoc := nameLoc

	// Optional @'host' suffix
	if p.cur.Kind == int('@') {
		p.advance() // consume '@'
		host, hostLoc, err2 := p.parseIdentifierOrString()
		if err2 != nil {
			return "", ast.NoLoc(), err2
		}
		name = name + "@" + host
		endLoc = hostLoc
	}

	return name, endLoc, nil
}
