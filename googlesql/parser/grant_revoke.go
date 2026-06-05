package parser

import (
	"github.com/bytebase/omni/googlesql/ast"
)

// ---------------------------------------------------------------------------
// DCL — GRANT / REVOKE (parser-dcl node)
// ---------------------------------------------------------------------------
//
// This file ports the legacy ANTLR grant_statement / revoke_statement family
// (GoogleSQLParser.g4 §2.8) — itself a hand-port of Google's open-source ZetaSQL
// reference, and the grammar bytebase consumes today. The full grammar:
//
//	grant_statement:  GRANT  privileges ON (identifier identifier?)? path_expression TO   grantee_list
//	revoke_statement: REVOKE privileges ON (identifier identifier?)? path_expression FROM grantee_list
//	privileges:       ALL PRIVILEGES? | privilege_list
//	privilege_list:   privilege (',' privilege)*
//	privilege:        privilege_name path_expression_list_with_parens?
//	privilege_name:   identifier | SELECT
//	path_expression_list_with_parens: '(' path_expression_list ')'
//	path_expression_list: path_expression (',' path_expression)*
//	path_expression:  identifier ('.' identifier)*
//	grantee_list:     string_literal_or_parameter (',' string_literal_or_parameter)*
//	string_literal_or_parameter: string_literal | parameter_expression | system_variable_expression
//	parameter_expression: named_parameter_expression ('@' identifier) | '?'
//	system_variable_expression: '@@' path_expression
//
// ORACLE NOTE — DCL is a DIALECT-DIVERGENT zone. The Spanner emulator (the
// differential harness) speaks a DIFFERENT GRANT dialect — `GRANT priv ON TABLE
// x TO ROLE r`, role-name grantees, no string-literal grantees — and rejects
// EVERY legacy ZetaSQL form here, while accepting `TO ROLE`/`GRANT ROLE` forms
// the ZetaSQL grammar rejects. The two grammars are disjoint on grantees, so the
// emulator is NON-AUTHORITATIVE for DCL. The accept/reject oracle is the legacy
// GoogleSQLParser.g4 plus the canonical ZetaSQL corpus
// (parser/googlesql/examples/.../grant_and_revoke.sql); this code follows that,
// not the emulator. See the parser-dcl divergence ledger entry.

// parseGrantStmt parses a GRANT statement. The GRANT keyword has NOT yet been
// consumed when this is called (dispatch peeks it).
//
// It serves BOTH dialects (oracle.md, "one grammar for both"): the legacy ZetaSQL
// privilege grant with string/parameter grantees AND the Spanner forms — a
// privilege grant `… TO ROLE role [, ...]` and the role-to-role grant `GRANT ROLE
// role [, ...] TO ROLE role [, ...]`. The Spanner role surface is added by the
// parser-ddl-spanner node (spanner_schema_role.go) and is authoritative against
// the live emulator; the grantee list (parseGranteeList) absorbs `ROLE …` as the
// recipient.
func (p *Parser) parseGrantStmt() (ast.Node, error) {
	grant := p.advance() // consume GRANT

	// Spanner role-to-role grant: `GRANT ROLE role [, ...] TO ROLE role [, ...]`.
	// Taken ONLY when the head is unambiguously a role list (`ROLE <name> …`), so a
	// legacy privilege grant whose first privilege is named `ROLE` (`GRANT ROLE ON
	// foo TO 'x'`) still falls through to the privilege path below. The target must
	// itself be a `ROLE` list (the emulator REJECTS `GRANT ROLE a TO 'x'`).
	if p.atRoleToRoleGrant() {
		roles, err := p.parseRoleNameList()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(kwTO); err != nil {
			return nil, err
		}
		grantees, err := p.parseRoleGranteeList()
		if err != nil {
			return nil, err
		}
		return &ast.GrantStmt{
			Roles:    roles,
			Grantees: grantees,
			Loc:      ast.Loc{Start: grant.Loc.Start, End: p.prev.Loc.End},
		}, nil
	}

	allPrivs, privs, objType, paths, err := p.parseGrantRevokeHead()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(kwTO); err != nil {
		return nil, err
	}
	grantees, err := p.parseGranteeList()
	if err != nil {
		return nil, err
	}
	if err := p.checkObjectListGrantees(paths, grantees); err != nil {
		return nil, err
	}

	return &ast.GrantStmt{
		AllPrivileges: allPrivs,
		Privileges:    privs,
		ObjectType:    objType,
		Path:          paths[0],
		Paths:         paths,
		Grantees:      grantees,
		Loc:           ast.Loc{Start: grant.Loc.Start, End: p.prev.Loc.End},
	}, nil
}

// parseRevokeStmt parses a REVOKE statement. The REVOKE keyword has NOT yet been
// consumed when this is called. Like GRANT it serves both dialects, including the
// Spanner role-from-role form `REVOKE ROLE role [, ...] FROM ROLE role [, ...]`
// and `… FROM ROLE role [, ...]` grantees.
func (p *Parser) parseRevokeStmt() (ast.Node, error) {
	revoke := p.advance() // consume REVOKE

	// Spanner role-from-role revoke: `REVOKE ROLE role [, ...] FROM ROLE …` (same
	// disambiguation + strict ROLE target as the GRANT branch above).
	if p.atRoleToRoleGrant() {
		roles, err := p.parseRoleNameList()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(kwFROM); err != nil {
			return nil, err
		}
		grantees, err := p.parseRoleGranteeList()
		if err != nil {
			return nil, err
		}
		return &ast.RevokeStmt{
			Roles:    roles,
			Grantees: grantees,
			Loc:      ast.Loc{Start: revoke.Loc.Start, End: p.prev.Loc.End},
		}, nil
	}

	allPrivs, privs, objType, paths, err := p.parseGrantRevokeHead()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}
	grantees, err := p.parseGranteeList()
	if err != nil {
		return nil, err
	}
	if err := p.checkObjectListGrantees(paths, grantees); err != nil {
		return nil, err
	}

	return &ast.RevokeStmt{
		AllPrivileges: allPrivs,
		Privileges:    privs,
		ObjectType:    objType,
		Path:          paths[0],
		Paths:         paths,
		Grantees:      grantees,
		Loc:           ast.Loc{Start: revoke.Loc.Start, End: p.prev.Loc.End},
	}, nil
}

// checkObjectListGrantees enforces the union constraint that a comma-separated
// object list (>1 object) is a SPANNER-only form, which requires ROLE grantees.
// The legacy ZetaSQL grant names exactly one object, so a multi-object grant with
// a non-ROLE (string/parameter) grantee belongs to NEITHER dialect — the emulator
// REJECTS `GRANT SELECT ON TABLE t1, t2 TO 'user'`. A single-object grant is
// unconstrained (either dialect). Cross-model review surfaced the original
// over-acceptance.
func (p *Parser) checkObjectListGrantees(paths []ast.NamePath, grantees []*ast.Grantee) error {
	if len(paths) <= 1 {
		return nil
	}
	for _, g := range grantees {
		if g.Kind != ast.GranteeRole {
			return &ParseError{
				Loc: g.Loc,
				Msg: "syntax error: a GRANT/REVOKE over multiple objects requires ROLE grantees",
			}
		}
	}
	return nil
}

// parseRoleGranteeList parses the target of a role-to-role grant/revoke, which
// MUST be a `ROLE role_name [, ...]` list (the emulator REJECTS a string/parameter
// target there — `GRANT ROLE a TO 'x'` fails with "Expecting 'ROLE'"). cur is at
// the target `ROLE` keyword.
func (p *Parser) parseRoleGranteeList() ([]*ast.Grantee, error) {
	if !p.curIsRoleKeyword() {
		return nil, p.syntaxErrorAtCur()
	}
	roles, err := p.parseRoleNameList()
	if err != nil {
		return nil, err
	}
	return roleNamesToGrantees(roles), nil
}

// parseGrantRevokeHead parses the shared GRANT/REVOKE prefix up to (but not
// including) the TO / FROM keyword:
//
//	privileges ON (identifier identifier?)? path_expression (',' path_expression)*
//
// It returns (allPrivileges, privileges, objectType, paths, err) where paths is
// the (>= 1) ON object list. The legacy ZetaSQL grant names a single object; the
// Spanner role grant allows a comma-separated object list (`ON TABLE t1, t2 TO
// ROLE r`, emulator-verified). The verb keyword has already been consumed; on
// success cur is positioned at TO (GRANT) or FROM (REVOKE).
func (p *Parser) parseGrantRevokeHead() (bool, []*ast.Privilege, []string, []ast.NamePath, error) {
	allPrivs, privs, err := p.parsePrivileges()
	if err != nil {
		return false, nil, nil, nil, err
	}
	if _, err := p.expect(kwON); err != nil {
		return false, nil, nil, nil, err
	}
	objType, paths, err := p.parseObjectTypeAndPathList()
	if err != nil {
		return false, nil, nil, nil, err
	}
	return allPrivs, privs, objType, paths, nil
}

// parsePrivileges parses the grammar's `privileges`:
//
//	ALL PRIVILEGES? | privilege_list
//
// Returns (allPrivileges, privileges, err): exactly one path applies. ALL is the
// reserved keyword kwALL (so it can only mean "all privileges" here, never a
// privilege name), optionally followed by PRIVILEGES.
func (p *Parser) parsePrivileges() (bool, []*ast.Privilege, error) {
	if p.cur.Type == kwALL {
		p.advance() // consume ALL
		// PRIVILEGES is optional (`GRANT ALL ON ...` is valid).
		p.match(kwPRIVILEGES)
		return true, nil, nil
	}

	privs, err := p.parsePrivilegeList()
	if err != nil {
		return false, nil, err
	}
	return false, privs, nil
}

// parsePrivilegeList parses `privilege (',' privilege)*`. Requires at least one
// privilege; a trailing/leading comma is a syntax error (the loop expects a
// privilege after each comma).
func (p *Parser) parsePrivilegeList() ([]*ast.Privilege, error) {
	var privs []*ast.Privilege
	for {
		priv, err := p.parsePrivilege()
		if err != nil {
			return nil, err
		}
		privs = append(privs, priv)
		if _, ok := p.match(int(',')); ok {
			continue
		}
		break
	}
	return privs, nil
}

// parsePrivilege parses one `privilege`:
//
//	privilege_name path_expression_list_with_parens?
//	privilege_name: identifier | SELECT
//
// The optional parenthesized list is a column path list (e.g. the `(col1, col2)`
// in `insert(col1, col2)`).
func (p *Parser) parsePrivilege() (*ast.Privilege, error) {
	start := p.cur.Loc

	// privilege_name: identifier | SELECT. SELECT is reserved, so it would not
	// match the generic identifier predicate; the grammar lists it explicitly.
	var name string
	switch {
	case p.cur.Type == kwSELECT:
		name = p.advance().Str
	case isIdentifierStart(p.cur.Type):
		name = p.advance().Str
	default:
		return nil, p.syntaxErrorAtCur()
	}

	priv := &ast.Privilege{Name: name, Loc: ast.Loc{Start: start.Start, End: p.prev.Loc.End}}

	// Optional path_expression_list_with_parens: '(' path_expression_list ')'.
	if p.cur.Type == int('(') {
		cols, err := p.parsePathExpressionListWithParens()
		if err != nil {
			return nil, err
		}
		priv.Columns = cols
		priv.Loc.End = p.prev.Loc.End
	}

	return priv, nil
}

// parsePathExpressionListWithParens parses `'(' path_expression_list ')'`, the
// per-privilege column list. On entry cur is '('. Returns >= 1 path (the grammar
// requires at least one).
func (p *Parser) parsePathExpressionListWithParens() ([]ast.NamePath, error) {
	p.advance() // consume '('

	var paths []ast.NamePath
	for {
		path, err := p.parsePathExpression()
		if err != nil {
			return nil, err
		}
		paths = append(paths, path)
		if _, ok := p.match(int(',')); ok {
			continue
		}
		break
	}

	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}
	return paths, nil
}

// parseObjectTypeAndPathList parses the grammar's object-type + object-path list.
// The verb keyword and ON have already been consumed; on success cur is at
// TO/FROM. Returns (objType, paths, err) with len(paths) >= 1.
//
// TWO dialect shapes share this position:
//   - legacy ZetaSQL — `(identifier identifier?)? path_expression`: an OPTIONAL
//     0/1/2-word object type and EXACTLY ONE object path. `GRANT SELECT ON foo TO
//     'x'` (no type) is valid here.
//   - Spanner — `target_type path (',' path)*`: a comma-separated object list,
//     which REQUIRES a leading type keyword (TABLE/VIEW/CHANGE STREAM/TABLE
//     FUNCTION). `GRANT SELECT ON TABLE t1, t2 TO ROLE r` is valid; `ON foo, bar`
//     (no type) is NOT (the emulator: "Encountered 'foo' while parsing:
//     target_type").
//
// Therefore a comma-separated list is accepted ONLY when an object type was
// present. With no type keyword, exactly one path is allowed (the legacy form) and
// a trailing comma is a syntax error. (The complementary "multi-object ⇒ ROLE
// grantee" rule is enforced by checkObjectListGrantees after the grantee list.)
//
// Object-type disambiguation (type words and the first path are all bare
// identifiers): a leading word is a TYPE word only when followed by ANOTHER
// path-expression-start — NOT by `,` and NOT by TO/FROM:
//
//   - `ON foo TO`                 => no type,  paths [foo]
//   - `ON table foo TO`           => type [table], paths [foo]
//   - `ON table t1, t2 TO`        => type [table], paths [t1, t2]
//   - `ON materialized view v TO` => type [materialized, view], paths [v]
func (p *Parser) parseObjectTypeAndPathList() ([]string, []ast.NamePath, error) {
	objType, first, err := p.parseObjectTypeAndFirstPath()
	if err != nil {
		return nil, nil, err
	}
	paths := []ast.NamePath{first}

	// A comma-separated object list is the Spanner form, which requires a leading
	// object-type keyword. Without one, only the single legacy object is allowed,
	// so a comma here is a syntax error (left for the caller's expect(TO/FROM), or
	// reported directly so the diagnostic points at the comma).
	if p.cur.Type == int(',') && len(objType) == 0 {
		return nil, nil, p.syntaxErrorAtCur()
	}
	for p.cur.Type == int(',') {
		p.advance() // ','
		next, err := p.parsePathExpression()
		if err != nil {
			return nil, nil, err
		}
		paths = append(paths, next)
	}
	return objType, paths, nil
}

// parseObjectTypeAndFirstPath consumes the optional 0/1/2-word object type and the
// FIRST object path. A word is a type word only when followed by another bare
// identifier (not ',' and not TO/FROM); see parseObjectTypeAndPathList.
func (p *Parser) parseObjectTypeAndFirstPath() ([]string, ast.NamePath, error) {
	first, err := p.parsePathExpression()
	if err != nil {
		return nil, ast.NamePath{}, err
	}
	// `first` is the object path when the run ends here — at TO/FROM (no more
	// objects) or at ',' (more objects follow, so `first` is path #1, not a type).
	if p.atPathRunBoundary() {
		return nil, first, nil
	}

	// Another path_expression follows => `first` was the first object-type word
	// (which must be a single, non-dotted identifier).
	firstWord, err := singleIdentName(first)
	if err != nil {
		return nil, ast.NamePath{}, err
	}
	objType := []string{firstWord}

	path, err := p.parsePathExpression()
	if err != nil {
		return nil, ast.NamePath{}, err
	}
	if p.atPathRunBoundary() {
		return objType, path, nil
	}

	// A second object-type word: `path` was actually the second type word.
	secondWord, err := singleIdentName(path)
	if err != nil {
		return nil, ast.NamePath{}, err
	}
	objType = append(objType, secondWord)

	path, err = p.parsePathExpression()
	if err != nil {
		return nil, ast.NamePath{}, err
	}
	// After at most two type words the grammar requires the path-run boundary next;
	// anything else (e.g. a third bare word) is a syntax error, surfaced by the
	// caller's expect(TO/FROM) once the comma loop (if any) finishes.
	return objType, path, nil
}

// atPathRunBoundary reports whether cur ends the object-type/path run: a comma
// (the next comma-separated object path) or the TO/FROM grantee intro. A word
// before one of these is an object PATH, not an object-type word.
func (p *Parser) atPathRunBoundary() bool {
	return p.cur.Type == int(',') || p.atGranteeIntro()
}

// singleIdentName returns the single component of a path that must be a lone
// identifier (object-type word). It errors (syntax) if the path is dotted, which
// the grammar cannot produce for an object-type position.
func singleIdentName(path ast.NamePath) (string, error) {
	if len(path.Parts) != 1 {
		return "", &ParseError{
			Loc: path.Loc,
			Msg: "syntax error: object type in GRANT/REVOKE must be a single identifier, not a dotted path",
		}
	}
	return path.Parts[0], nil
}

// atGranteeIntro reports whether cur is the TO/FROM keyword that introduces the
// grantee list — i.e. the object-type/path run has ended.
func (p *Parser) atGranteeIntro() bool {
	return p.cur.Type == kwTO || p.cur.Type == kwFROM
}

// parsePathExpression parses the grammar's `path_expression: identifier ('.'
// identifier)*`. Each component is a GoogleSQL identifier (token_identifier or a
// non-reserved keyword-as-identifier; see isIdentifierStart). Requires at least
// one component.
func (p *Parser) parsePathExpression() (ast.NamePath, error) {
	if !isIdentifierStart(p.cur.Type) {
		return ast.NamePath{}, p.syntaxErrorAtCur()
	}
	start := p.cur.Loc
	first := p.advance()
	parts := []string{first.Str}
	end := first.Loc

	for p.cur.Type == int('.') {
		p.advance() // consume '.'
		if !isIdentifierStart(p.cur.Type) {
			return ast.NamePath{}, p.syntaxErrorAtCur()
		}
		next := p.advance()
		parts = append(parts, next.Str)
		end = next.Loc
	}

	return ast.NamePath{Parts: parts, Loc: ast.Loc{Start: start.Start, End: end.End}}, nil
}

// parseGranteeList parses the recipients after TO (GRANT) / FROM (REVOKE). It
// accepts the UNION of both dialects:
//
//   - legacy ZetaSQL `grantee_list: string_literal_or_parameter (',' …)*`
//     (string / @param / '?' / @@sysvar grantees), and
//   - the Spanner `ROLE role_name (',' role_name)*` form — a single `ROLE`
//     keyword introduces a comma-separated role-name list (parser-ddl-spanner;
//     emulator-authoritative). The two are NOT mixable: a list either starts with
//     `ROLE` (all role grantees) or with a string/parameter (all legacy grantees),
//     matching the emulator (`TO ROLE r1, ROLE r2` REJECTS — the `ROLE` keyword
//     appears once).
//
// Requires at least one grantee; a trailing/leading comma is a syntax error.
func (p *Parser) parseGranteeList() ([]*ast.Grantee, error) {
	// Spanner role grantee list: `ROLE role_name [, role_name ...]`.
	if p.curIsRoleKeyword() {
		roles, err := p.parseRoleNameList()
		if err != nil {
			return nil, err
		}
		return roleNamesToGrantees(roles), nil
	}

	var grantees []*ast.Grantee
	for {
		g, err := p.parseGrantee()
		if err != nil {
			return nil, err
		}
		grantees = append(grantees, g)
		if _, ok := p.match(int(',')); ok {
			continue
		}
		break
	}
	return grantees, nil
}

// parseGrantee parses one `string_literal_or_parameter`:
//
//	string_literal              -> GranteeString
//	'@' identifier              -> GranteeNamedParameter      (named_parameter_expression)
//	'?'                         -> GranteePositionalParameter (QUESTION)
//	'@@' path_expression        -> GranteeSystemVariable      (system_variable_expression)
func (p *Parser) parseGrantee() (*ast.Grantee, error) {
	start := p.cur.Loc

	switch p.cur.Type {
	case tokString:
		tok := p.advance()
		return &ast.Grantee{Kind: ast.GranteeString, Value: tok.Str, Loc: tok.Loc}, nil

	case int('?'):
		tok := p.advance()
		return &ast.Grantee{Kind: ast.GranteePositionalParameter, Loc: tok.Loc}, nil

	case int('@'):
		p.advance() // consume '@'
		// named_parameter_expression: '@' identifier.
		if !isIdentifierStart(p.cur.Type) {
			return nil, p.syntaxErrorAtCur()
		}
		name := p.advance()
		return &ast.Grantee{
			Kind:  ast.GranteeNamedParameter,
			Value: name.Str,
			Loc:   ast.Loc{Start: start.Start, End: name.Loc.End},
		}, nil

	case tokAtAt:
		p.advance() // consume '@@'
		// system_variable_expression: '@@' path_expression.
		path, err := p.parsePathExpression()
		if err != nil {
			return nil, err
		}
		return &ast.Grantee{
			Kind:  ast.GranteeSystemVariable,
			Value: path.String(),
			Loc:   ast.Loc{Start: start.Start, End: path.Loc.End},
		}, nil

	default:
		return nil, p.syntaxErrorAtCur()
	}
}
