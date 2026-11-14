package parser

import (
	"strings"

	"github.com/bytebase/omni/doris/ast"
)

// GRANT and REVOKE are unsupportedGrantRevokeStatement (DorisParser.g4:532),
// which Doris also parses with its CUP parser's grant_stmt and revoke_stmt
// (DorisLegacyParser.g4:895, :939). The readers here take what both take.

// grantHead is a GRANT or REVOKE read up to its grantee.
type grantHead struct {
	roles      []string
	privileges []string
	columns    [][]string
	objectType string
	object     *ast.ObjectName
}

// parseGrant parses GRANT. On entry, cur is the token after GRANT.
func (p *Parser) parseGrant(startLoc ast.Loc) (ast.Node, error) {
	head, err := p.parseGrantHead(kwTO)
	if err != nil {
		return nil, err
	}
	kind, grantee, err := p.parseGrantee(head.roles == nil)
	if err != nil {
		return nil, err
	}
	return &ast.GrantStmt{
		Roles:            head.roles,
		Privileges:       head.privileges,
		PrivilegeColumns: head.columns,
		ObjectType:       head.objectType,
		Object:           head.object,
		ToType:           kind,
		Grantees:         []string{grantee},
		Loc:              ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End},
	}, nil
}

// parseRevoke parses REVOKE. On entry, cur is the token after REVOKE.
func (p *Parser) parseRevoke(startLoc ast.Loc) (ast.Node, error) {
	head, err := p.parseGrantHead(kwFROM)
	if err != nil {
		return nil, err
	}
	kind, revokee, err := p.parseGrantee(head.roles == nil)
	if err != nil {
		return nil, err
	}
	return &ast.RevokeStmt{
		Roles:            head.roles,
		Privileges:       head.privileges,
		PrivilegeColumns: head.columns,
		ObjectType:       head.objectType,
		Object:           head.object,
		FromType:         kind,
		Revokees:         []string{revokee},
		Loc:              ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End},
	}, nil
}

// parseGrantHead reads a GRANT or REVOKE through to (TO or FROM):
//
//	STRING_LITERAL (COMMA STRING_LITERAL)*                     (a role list)
//	privilegeList ON multipartIdentifierOrAsterisk
//	privilegeList ON (RESOURCE | CLUSTER | COMPUTE GROUP | STAGE | STORAGE VAULT
//	    | WORKLOAD GROUP) identifierOrTextOrAsterisk
func (p *Parser) parseGrantHead(to int) (*grantHead, error) {
	head := &grantHead{}
	if p.cur.Kind == tokString {
		for {
			if p.cur.Kind != tokString {
				return nil, p.syntaxErrorAtCur()
			}
			head.roles = append(head.roles, p.advance().Str)
			if p.cur.Kind != int(',') {
				break
			}
			p.advance() // consume ','
		}
		if _, err := p.expect(to); err != nil {
			return nil, err
		}
		return head, nil
	}

	// privilegeList (:552): privilege (COMMA privilege)*, where privilege is
	// identifier identifierList? | ALL, as the CUP parser's privilege_type.
	for {
		var name string
		if p.cur.Kind == kwALL {
			name = "ALL"
			p.advance()
		} else {
			n, _, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			name = strings.ToUpper(n)
		}
		var cols []string
		if name != "ALL" && p.cur.Kind == int('(') {
			p.advance() // consume '('
			names, err := p.parseIdentifierList()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(int(')')); err != nil {
				return nil, err
			}
			cols = names
		}
		head.privileges = append(head.privileges, name)
		head.columns = append(head.columns, cols)
		if p.cur.Kind != int(',') {
			break
		}
		p.advance() // consume ','
	}
	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}

	if typ, words := p.grantResourceType(to); typ != "" {
		for range words {
			p.advance()
		}
		head.objectType = typ
		// identifierOrTextOrAsterisk, as the CUP parser's resource_pattern
		// (ident_or_star | STRING_LITERAL).
		tok := p.cur
		switch {
		case tok.Kind == int('*') || tok.Kind == tokString:
			p.advance()
			name := tok.Str
			if tok.Kind == int('*') {
				name = "*"
			}
			head.object = &ast.ObjectName{Parts: []string{name}, Loc: tok.Loc}
		default:
			name, loc, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			head.object = &ast.ObjectName{Parts: []string{name}, Loc: loc}
		}
	} else {
		head.objectType = "TABLE"
		obj, err := p.parseGrantTablePattern()
		if err != nil {
			return nil, err
		}
		head.object = obj
	}
	if _, err := p.expect(to); err != nil {
		return nil, err
	}
	return head, nil
}

// grantResourceType reports the resource type the words at cur name, and
// how many words it takes, when they lead a resource pattern rather than a
// table named RESOURCE, CLUSTER or STAGE: a two-word type always does, a
// one-word type when to (TO or FROM) or '.' does not follow it.
func (p *Parser) grantResourceType(to int) (string, []int) {
	next := p.peekNext().Kind
	switch p.cur.Kind {
	case kwRESOURCE, kwCLUSTER, kwSTAGE:
		if next == to || next == int('.') {
			return "", nil
		}
		return strings.ToUpper(p.cur.Str), []int{p.cur.Kind}
	case kwCOMPUTE, kwWORKLOAD:
		if next == kwGROUP {
			return strings.ToUpper(p.cur.Str) + " GROUP", []int{p.cur.Kind, kwGROUP}
		}
	case kwSTORAGE:
		if next == kwVAULT {
			return "STORAGE VAULT", []int{kwSTORAGE, kwVAULT}
		}
	}
	return "", nil
}

// parseGrantTablePattern reads a table pattern as both grammars take it:
// DorisParser.g4's multipartIdentifierOrAsterisk (:1043) and the CUP
// parser's tbl_pattern (DorisLegacyParser.g4:914), so one to three parts,
// each an identifier or '*'.
func (p *Parser) parseGrantTablePattern() (*ast.ObjectName, error) {
	start := p.cur.Loc
	var parts []string
	for {
		if p.cur.Kind == int('*') {
			p.advance()
			parts = append(parts, "*")
		} else {
			name, _, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			parts = append(parts, name)
		}
		if len(parts) == 3 || p.cur.Kind != int('.') {
			break
		}
		p.advance() // consume '.'
	}
	return &ast.ObjectName{Parts: parts, Loc: ast.Loc{Start: start.Start, End: p.prev.Loc.End}}, nil
}

// parseGrantee reads the grantee after TO or FROM: userIdentify, or, after
// privileges, ROLE STRING_LITERAL. It returns ROLE or "" and the grantee as
// written.
//
// userIdentify (:1052) puts a host in parentheses where the CUP parser's
// user_identity (DorisLegacyParser.g4:571) puts it in brackets, so neither
// is taken: user ('@' host)?, each an identifier or a string.
func (p *Parser) parseGrantee(roleAllowed bool) (string, string, error) {
	if roleAllowed && p.cur.Kind == kwROLE && p.peekNext().Kind == tokString {
		p.advance() // consume ROLE
		return "ROLE", p.advance().Str, nil
	}
	user, _, err := p.parseIdentifierOrString()
	if err != nil {
		return "", "", err
	}
	if p.cur.Kind == int('@') {
		p.advance() // consume '@'
		host, _, err := p.parseIdentifierOrString()
		if err != nil {
			return "", "", err
		}
		user += "@" + host
	}
	return "", user, nil
}
