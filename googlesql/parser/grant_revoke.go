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
func (p *Parser) parseGrantStmt() (ast.Node, error) {
	grant := p.advance() // consume GRANT

	allPrivs, privs, objType, path, err := p.parseGrantRevokeHead()
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

	return &ast.GrantStmt{
		AllPrivileges: allPrivs,
		Privileges:    privs,
		ObjectType:    objType,
		Path:          path,
		Grantees:      grantees,
		Loc:           ast.Loc{Start: grant.Loc.Start, End: p.prev.Loc.End},
	}, nil
}

// parseRevokeStmt parses a REVOKE statement. The REVOKE keyword has NOT yet been
// consumed when this is called.
func (p *Parser) parseRevokeStmt() (ast.Node, error) {
	revoke := p.advance() // consume REVOKE

	allPrivs, privs, objType, path, err := p.parseGrantRevokeHead()
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

	return &ast.RevokeStmt{
		AllPrivileges: allPrivs,
		Privileges:    privs,
		ObjectType:    objType,
		Path:          path,
		Grantees:      grantees,
		Loc:           ast.Loc{Start: revoke.Loc.Start, End: p.prev.Loc.End},
	}, nil
}

// parseGrantRevokeHead parses the shared GRANT/REVOKE prefix up to (but not
// including) the TO / FROM keyword:
//
//	privileges ON (identifier identifier?)? path_expression
//
// It returns (allPrivileges, privileges, objectType, path, err). The verb keyword
// has already been consumed; on success cur is positioned at TO (GRANT) or FROM
// (REVOKE).
func (p *Parser) parseGrantRevokeHead() (bool, []*ast.Privilege, []string, ast.NamePath, error) {
	allPrivs, privs, err := p.parsePrivileges()
	if err != nil {
		return false, nil, nil, ast.NamePath{}, err
	}
	if _, err := p.expect(kwON); err != nil {
		return false, nil, nil, ast.NamePath{}, err
	}
	objType, path, err := p.parseObjectTypeAndPath()
	if err != nil {
		return false, nil, nil, ast.NamePath{}, err
	}
	return allPrivs, privs, objType, path, nil
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

// parseObjectTypeAndPath parses the grammar's `(identifier identifier?)?
// path_expression` — the optional 0/1/2-word object type followed by the
// required object path. The verb keyword and ON have already been consumed; on
// success cur is at TO/FROM.
//
// Disambiguation (matching ANTLR's longest-type-then-backtrack over this rule):
// the object type words are SINGLE (non-dotted) identifiers and the path is one
// dotted path_expression, then TO/FROM. So:
//
//   - Parse a path_expression P1 (greedy on dots).
//   - If TO/FROM follows, there is NO object type and P1 is the path.
//     (e.g. `ON datascape.foo TO`, `ON foo TO`)
//   - Otherwise another identifier follows, so P1 was the FIRST type word and
//     must be a single (non-dotted) identifier. Parse the next path_expression
//     P2 and repeat once more for an optional SECOND type word.
//     (e.g. `ON table foo TO` => type [table], path foo;
//     `ON materialized view foo TO` => type [materialized, view], path foo)
//
// A would-be type word that is actually dotted (e.g. `ON a.b c TO`) has no valid
// parse in the grammar (type words are single identifiers) and is reported as a
// syntax error.
func (p *Parser) parseObjectTypeAndPath() ([]string, ast.NamePath, error) {
	first, err := p.parsePathExpression()
	if err != nil {
		return nil, ast.NamePath{}, err
	}
	// No object type: the first path_expression is the object name.
	if p.atGranteeIntro() {
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
	if p.atGranteeIntro() {
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
	// After at most two type words the grammar requires the grantee intro next;
	// anything else (e.g. a third bare word) is a syntax error, surfaced by the
	// caller's expect(TO/FROM).
	return objType, path, nil
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

// parseGranteeList parses the grammar's `grantee_list: string_literal_or_parameter
// (',' string_literal_or_parameter)*`. Requires at least one grantee; a
// trailing/leading comma is a syntax error.
func (p *Parser) parseGranteeList() ([]*ast.Grantee, error) {
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
