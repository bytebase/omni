package parser

import (
	"github.com/bytebase/omni/googlesql/ast"
)

// ---------------------------------------------------------------------------
// Spanner DDL — ROLE + role-based GRANT/REVOKE (parser-ddl-spanner node)
// ---------------------------------------------------------------------------
//
// Cloud Spanner fine-grained access control (truth1 DDL-032..037):
//
//	CREATE ROLE role_name
//	DROP   ROLE role_name
//	GRANT  privilege [, ...] ON { TABLE | VIEW | CHANGE STREAM | TABLE FUNCTION } name [, ...]
//	         TO ROLE role_name [, ...]
//	GRANT  ROLE role_name [, ...] TO ROLE role_name [, ...]
//	REVOKE … FROM ROLE role_name [, ...]
//	REVOKE ROLE role_name [, ...] FROM ROLE role_name [, ...]
//
// DIALECT NOTE — the legacy ANTLR grammar's grant_statement/revoke_statement
// accept ONLY string-literal/parameter grantees (string_literal_or_parameter) and
// have NO `ROLE` grantee, NO `GRANT ROLE` head, and NO CREATE/DROP ROLE rule, so
// the legacy parser over-rejects the entire Spanner role surface. The omni parser
// — whose authoritative oracle for Spanner DDL is the live emulator (oracle.md) —
// accepts both dialects: the existing parser-dcl code (grant_revoke.go) still
// handles the ZetaSQL string-grantee form, and the hooks here add the Spanner role
// forms to that same union path. Verified ACCEPT/REJECT in
// spanner_ddl_oracle_test.go; the emulator REJECTS a string grantee (`Expecting
// 'ROLE'`) and ACCEPTS `TO ROLE name [, name]` / `GRANT ROLE …`.
//
// `ROLE` lexes as a bare (non-reserved) identifier, so it is matched by spelling.

// parseCreateRole parses a CREATE ROLE body after the shared CREATE prefix. cur is
// at `ROLE` (dispatcher confirmed). create_role takes no OR REPLACE / scope / IF
// NOT EXISTS in the Spanner grammar.
//
// The role name is a SINGLE identifier, NOT a dotted path: the emulator REJECTS
// `CREATE ROLE a.b` ("Expecting 'EOF' but found '.'"). A backtick-quoted name
// (`\`dotted.name\``) is one identifier token and is accepted. (DROP ROLE, by
// contrast, accepts a dotted path — see parseDropRole.)
func (p *Parser) parseCreateRole(create Token, orReplace bool, scope string) (ast.Node, error) {
	if orReplace || scope != "" {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // ROLE

	stmt := &ast.CreateRoleStmt{}
	stmt.Loc.Start = create.Loc.Start

	name, err := p.parseRoleName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseDropRole parses a DROP ROLE body. The DROP keyword has been consumed by
// parseDropStmt; cur is at `ROLE` (dispatcher confirmed). No IF EXISTS in the
// Spanner grammar.
//
// Unlike CREATE ROLE, the emulator ACCEPTS a dotted path here (`DROP ROLE a.b.c`
// accepts), so parseTablePath (which folds dotted parts) is correct.
func (p *Parser) parseDropRole(drop Token) (ast.Node, error) {
	p.advance() // ROLE

	stmt := &ast.DropRoleStmt{}
	stmt.Loc.Start = drop.Loc.Start

	name, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseRoleName parses a single role-name identifier as a one-part *PathExpr. A
// role name in CREATE ROLE / GRANT / REVOKE is a lone identifier (Spanner rejects
// a dotted role name); a backtick-quoted name lands in a single part.
func (p *Parser) parseRoleName() (*ast.PathExpr, error) {
	tok, err := p.expectIdentifier()
	if err != nil {
		return nil, err
	}
	name, err := p.identifierText(tok)
	if err != nil {
		return nil, err
	}
	return &ast.PathExpr{Parts: []string{name}, Loc: tok.Loc}, nil
}

// curIsRoleKeyword reports whether cur is the bare-identifier `ROLE` word that
// introduces a Spanner role grantee list / role-grant subject. (`ROLE` is not a
// reserved token; the lexer emits it as tokIdentifier.)
func (p *Parser) curIsRoleKeyword() bool {
	return p.curIsWord("ROLE")
}

// atRoleToRoleGrant reports whether the head is a Spanner role-to-role
// grant/revoke — `GRANT|REVOKE ROLE <role_name> …` — as opposed to a legacy
// privilege grant whose first privilege happens to be named `ROLE`
// (`GRANT ROLE ON …`, `GRANT ROLE, SELECT ON …`, `GRANT ROLE(cols) ON …`).
//
// The disambiguator: `ROLE` here lexes as a generic identifier, so it is BOTH a
// valid privilege_name AND the role-list introducer. They are told apart by the
// token AFTER `ROLE`: a role-to-role grant is `ROLE <role_name> …` where
// role_name is an identifier; a privilege grant is `ROLE { ON | , | ( }`. So the
// role-to-role form is exactly "ROLE followed by an identifier-start token". This
// matches the live emulator, which REJECTS `GRANT ROLE ON …` (privilege ROLE is
// not a Spanner privilege) yet the BigQuery/ZetaSQL UNION must still accept the
// legacy `GRANT ROLE ON foo TO 'x'` (privilege_name: identifier), so the
// role-to-role branch must NOT swallow it.
func (p *Parser) atRoleToRoleGrant() bool {
	return p.curIsRoleKeyword() && isIdentifierStart(p.peekNext().Type)
}

// parseRoleNameList parses `ROLE role_name (',' role_name)*` — a Spanner role list
// introduced by a single `ROLE` keyword. cur must be at `ROLE`. Used for both the
// GRANT/REVOKE grantee list (`TO ROLE r [, r]`) and the role-grant subject
// (`GRANT ROLE r [, r] …`). Each role_name is a SINGLE identifier — the emulator
// REJECTS a dotted role name (`TO ROLE r1, r2.x` and `GRANT ROLE r1, r2.x …` both
// fail at the '.'). Returns the role names as one-part NamePaths.
func (p *Parser) parseRoleNameList() ([]ast.NamePath, error) {
	p.advance() // ROLE
	var roles []ast.NamePath
	for {
		tok, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}
		name, err := p.identifierText(tok)
		if err != nil {
			return nil, err
		}
		roles = append(roles, ast.NamePath{Parts: []string{name}, Loc: tok.Loc})
		if _, ok := p.match(int(',')); ok {
			continue
		}
		break
	}
	return roles, nil
}

// roleNamesToGrantees converts a Spanner role list into GranteeRole grantees (the
// TO/FROM recipients of a GRANT/REVOKE).
func roleNamesToGrantees(roles []ast.NamePath) []*ast.Grantee {
	out := make([]*ast.Grantee, 0, len(roles))
	for _, r := range roles {
		out = append(out, &ast.Grantee{Kind: ast.GranteeRole, Value: r.String(), Loc: r.Loc})
	}
	return out
}
