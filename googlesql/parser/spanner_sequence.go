package parser

import (
	"github.com/bytebase/omni/googlesql/ast"
)

// ---------------------------------------------------------------------------
// Spanner DDL — SEQUENCE + LOCALITY GROUP + PROTO BUNDLE (parser-ddl-spanner)
// ---------------------------------------------------------------------------
//
// Cloud Spanner sequences / locality groups / proto bundles (truth1
// DDL-027..029, 041..043, 046..047):
//
//	CREATE SEQUENCE [IF NOT EXISTS] name [OPTIONS ( … )]
//	ALTER  SEQUENCE [IF EXISTS] name SET OPTIONS ( … )
//	DROP   SEQUENCE [IF EXISTS] name
//
//	CREATE LOCALITY GROUP name [OPTIONS ( … )]
//	ALTER  LOCALITY GROUP name SET OPTIONS ( … )
//	DROP   LOCALITY GROUP name
//
//	CREATE PROTO BUNDLE ( proto_type_name [, ...] )
//	ALTER  PROTO BUNDLE { INSERT ( … ) | UPDATE ( … ) | DELETE ( … ) } [, ...]
//
// DIALECT NOTE — like CHANGE STREAM, none of these is a first-class rule in the
// legacy ANTLR grammar; the omni parser models them directly. The authoritative
// oracle is the live Cloud Spanner emulator, which ACCEPTS each form
// (spanner_ddl_oracle_test.go). `SEQUENCE` is a (non-reserved) keyword token;
// `LOCALITY` lexes as a bare identifier and `GROUP` as a reserved keyword, so the
// dispatcher recognizes LOCALITY by spelling; `PROTO` is a reserved keyword and
// `BUNDLE` a bare identifier.

// --- SEQUENCE ---

// parseCreateSequence parses a CREATE SEQUENCE body after the shared CREATE prefix.
// cur is at the SEQUENCE keyword. create_sequence has no opt_or_replace /
// opt_create_scope, so a leading OR REPLACE / scope rejects.
func (p *Parser) parseCreateSequence(create Token, orReplace bool, scope string) (ast.Node, error) {
	if orReplace || scope != "" {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // SEQUENCE

	stmt := &ast.CreateSequenceStmt{}
	stmt.Loc.Start = create.Loc.Start

	ifNotExists, err := p.parseIfNotExists()
	if err != nil {
		return nil, err
	}
	stmt.IfNotExists = ifNotExists

	name, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	if p.cur.Type == kwOPTIONS {
		opts, err := p.parseOptionsList()
		if err != nil {
			return nil, err
		}
		stmt.Options = opts
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseAlterSequence parses an ALTER SEQUENCE body. The ALTER keyword has been
// consumed by parseAlterStmt; cur is at the SEQUENCE keyword.
func (p *Parser) parseAlterSequence(alter Token) (ast.Node, error) {
	p.advance() // SEQUENCE

	stmt := &ast.AlterSequenceStmt{}
	stmt.Loc.Start = alter.Loc.Start

	ifExists, err := p.parseIfExists()
	if err != nil {
		return nil, err
	}
	stmt.IfExists = ifExists

	name, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// SET OPTIONS ( … ) is the only ALTER SEQUENCE action.
	if _, err := p.expect(kwSET); err != nil {
		return nil, err
	}
	if p.cur.Type != kwOPTIONS {
		return nil, p.syntaxErrorAtCur()
	}
	opts, err := p.parseOptionsList()
	if err != nil {
		return nil, err
	}
	stmt.SetOptions = opts

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseDropSequence parses a DROP SEQUENCE body. The DROP keyword has been
// consumed by parseDropStmt; cur is at the SEQUENCE keyword.
func (p *Parser) parseDropSequence(drop Token) (ast.Node, error) {
	p.advance() // SEQUENCE

	stmt := &ast.DropSequenceStmt{}
	stmt.Loc.Start = drop.Loc.Start

	ifExists, err := p.parseIfExists()
	if err != nil {
		return nil, err
	}
	stmt.IfExists = ifExists

	name, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// --- LOCALITY GROUP ---

// parseCreateLocalityGroup parses a CREATE LOCALITY GROUP body after the shared
// CREATE prefix. cur is at `LOCALITY` (dispatcher confirmed `LOCALITY GROUP`).
func (p *Parser) parseCreateLocalityGroup(create Token, orReplace bool, scope string) (ast.Node, error) {
	if orReplace || scope != "" {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // LOCALITY
	p.advance() // GROUP

	stmt := &ast.CreateLocalityGroupStmt{}
	stmt.Loc.Start = create.Loc.Start

	name, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	if p.cur.Type == kwOPTIONS {
		opts, err := p.parseOptionsList()
		if err != nil {
			return nil, err
		}
		stmt.Options = opts
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseAlterLocalityGroup parses an ALTER LOCALITY GROUP body. The ALTER keyword
// has been consumed; cur is at `LOCALITY` (dispatcher confirmed `LOCALITY GROUP`).
//
// The `SET OPTIONS ( … )` clause is OPTIONAL: the live emulator ACCEPTS a bare
// `ALTER LOCALITY GROUP name` (verdict accept, semantic "Locality group not
// found"), unlike ALTER SEQUENCE / ALTER CHANGE STREAM whose action is required
// (those REJECT a bare form). When present, the clause must be `SET OPTIONS` (a
// `SET` with no `OPTIONS`, or any other trailing token, REJECTS — emulator-
// verified). See spanner_ddl_oracle_test.go.
func (p *Parser) parseAlterLocalityGroup(alter Token) (ast.Node, error) {
	p.advance() // LOCALITY
	p.advance() // GROUP

	stmt := &ast.AlterLocalityGroupStmt{}
	stmt.Loc.Start = alter.Loc.Start

	name, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Optional SET OPTIONS ( … ).
	if p.cur.Type == kwSET {
		p.advance() // SET
		if p.cur.Type != kwOPTIONS {
			return nil, p.syntaxErrorAtCur()
		}
		opts, err := p.parseOptionsList()
		if err != nil {
			return nil, err
		}
		stmt.SetOptions = opts
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseDropLocalityGroup parses a DROP LOCALITY GROUP body. The DROP keyword has
// been consumed; cur is at `LOCALITY` (dispatcher confirmed `LOCALITY GROUP`).
func (p *Parser) parseDropLocalityGroup(drop Token) (ast.Node, error) {
	p.advance() // LOCALITY
	p.advance() // GROUP

	stmt := &ast.DropLocalityGroupStmt{}
	stmt.Loc.Start = drop.Loc.Start

	name, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// --- PROTO BUNDLE ---

// parseCreateProtoBundle parses a CREATE PROTO BUNDLE body after the shared CREATE
// prefix. cur is at the PROTO keyword (dispatcher confirmed `PROTO BUNDLE`).
// create_proto_bundle has no opt_or_replace / opt_create_scope.
func (p *Parser) parseCreateProtoBundle(create Token, orReplace bool, scope string) (ast.Node, error) {
	if orReplace || scope != "" {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // PROTO
	p.advance() // BUNDLE

	stmt := &ast.CreateProtoBundleStmt{}
	stmt.Loc.Start = create.Loc.Start

	// `( proto_type_name [, ...] )` — at least one type (an empty bundle rejects).
	types, err := p.parseProtoTypeNameListWithParens()
	if err != nil {
		return nil, err
	}
	stmt.Types = types

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseAlterProtoBundle parses an ALTER PROTO BUNDLE body. The ALTER keyword has
// been consumed; cur is at the PROTO keyword (dispatcher confirmed `PROTO
// BUNDLE`).
//
//	ALTER PROTO BUNDLE [INSERT ( … )] [UPDATE ( … )] [DELETE ( … )]
//
// The three groups are OPTIONAL, appear AT MOST ONCE EACH, in the FIXED order
// INSERT → UPDATE → DELETE, SPACE-separated (NO commas between groups), and at
// least one must be present. This is exactly the emulator's grammar (verified in
// spanner_ddl_oracle_test.go): a repeated group (`INSERT (…) INSERT (…)`), an
// out-of-order group (`UPDATE (…) INSERT (…)`), and a comma between groups
// (`INSERT (…), UPDATE (…)`) all REJECT — so a naive "any order, any repeat,
// optional comma" loop would be over-permissive.
func (p *Parser) parseAlterProtoBundle(alter Token) (ast.Node, error) {
	p.advance() // PROTO
	p.advance() // BUNDLE

	stmt := &ast.AlterProtoBundleStmt{}
	stmt.Loc.Start = alter.Loc.Start

	seen := false
	// INSERT ( … )?
	if p.cur.Type == kwINSERT {
		p.advance()
		types, err := p.parseProtoTypeNameListWithParens()
		if err != nil {
			return nil, err
		}
		stmt.Insert = types
		seen = true
	}
	// UPDATE ( … )?
	if p.cur.Type == kwUPDATE {
		p.advance()
		types, err := p.parseProtoTypeNameListWithParens()
		if err != nil {
			return nil, err
		}
		stmt.Update = types
		seen = true
	}
	// DELETE ( … )?
	if p.cur.Type == kwDELETE {
		p.advance()
		types, err := p.parseProtoTypeNameListWithParens()
		if err != nil {
			return nil, err
		}
		stmt.Delete = types
		seen = true
	}
	if !seen {
		// No INSERT/UPDATE/DELETE group at all — the grammar requires at least one.
		// (Also catches an out-of-order or comma-separated group: a leading UPDATE/
		// DELETE is consumed above, but a later out-of-order INSERT / a stray comma
		// is left for parseSingle's EOF check to reject.)
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseProtoTypeNameListWithParens parses `'(' proto_type_name (',' proto_type_name)* ','? ')'`
// where each proto_type_name is a (typically backtick-quoted) dotted path. cur is
// at '('. Requires at least one type name (an empty `()` rejects), but a TRAILING
// comma is allowed: the emulator ACCEPTS `CREATE PROTO BUNDLE (\`a.B\`,)` and
// `ALTER PROTO BUNDLE INSERT (\`a.B\`,)` (unlike a change-stream column list, which
// rejects a trailing comma). Verified in spanner_ddl_oracle_test.go.
func (p *Parser) parseProtoTypeNameListWithParens() ([]ast.NamePath, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	var types []ast.NamePath
	for {
		path, err := p.parseProtoTypeName()
		if err != nil {
			return nil, err
		}
		types = append(types, path)
		if _, ok := p.match(int(',')); !ok {
			break
		}
		// A comma was consumed; a closing ')' here is a permitted trailing comma.
		if p.cur.Type == int(')') {
			break
		}
	}
	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}
	return types, nil
}

// parseProtoTypeName parses one proto_type_name: a dotted name `identifier ('.'
// identifier)*`. A backtick-quoted name (the common `\`my.package.Msg\`` form) is
// a single identifier whose Str already holds the dotted body, so it lands in one
// part; an unquoted dotted name spreads across parts.
func (p *Parser) parseProtoTypeName() (ast.NamePath, error) {
	if !isIdentifierStart(p.cur.Type) {
		return ast.NamePath{}, p.syntaxErrorAtCur()
	}
	start := p.cur.Loc
	first := p.advance()
	firstText, err := p.identifierText(first)
	if err != nil {
		return ast.NamePath{}, err
	}
	parts := []string{firstText}
	end := first.Loc

	for p.cur.Type == int('.') {
		if !isAnyKeywordIdentifier(p.peekNext().Type) {
			break
		}
		p.advance() // '.'
		next := p.advance()
		text, err := p.identifierText(next)
		if err != nil {
			return ast.NamePath{}, err
		}
		parts = append(parts, text)
		end = next.Loc
	}

	return ast.NamePath{Parts: parts, Loc: ast.Loc{Start: start.Start, End: end.End}}, nil
}
