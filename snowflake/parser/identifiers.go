package parser

import "github.com/bytebase/omni/snowflake/ast"

// parseIdent parses one identifier from the token stream. Accepts:
//   - tokIdent (bare identifier)
//   - tokQuotedIdent (double-quoted identifier)
//   - Any non-reserved keyword token (type >= 700 and not in keywordReserved)
//
// Returns a ParseError if the current token is none of the above.
func (p *Parser) parseIdent() (ast.Ident, error) {
	switch {
	case p.cur.Type == tokIdent:
		tok := p.advance()
		return ast.Ident{Name: tok.Str, Quoted: false, Loc: tok.Loc}, nil
	case p.cur.Type == tokQuotedIdent:
		tok := p.advance()
		return ast.Ident{Name: tok.Str, Quoted: true, Loc: tok.Loc}, nil
	case p.cur.Type >= 700 && !keywordReserved[p.cur.Type]:
		// Non-reserved keyword used as an identifier.
		tok := p.advance()
		return ast.Ident{Name: tok.Str, Quoted: false, Loc: tok.Loc}, nil
	}
	return ast.Ident{}, &ParseError{
		Loc: p.cur.Loc,
		Msg: "expected identifier",
	}
}

// parseIdentStrict parses one identifier but ONLY accepts tokIdent or
// tokQuotedIdent — NOT non-reserved keywords. Used in contexts where
// keyword-as-identifier is not allowed (rare in Snowflake).
func (p *Parser) parseIdentStrict() (ast.Ident, error) {
	switch p.cur.Type {
	case tokIdent:
		tok := p.advance()
		return ast.Ident{Name: tok.Str, Quoted: false, Loc: tok.Loc}, nil
	case tokQuotedIdent:
		tok := p.advance()
		return ast.Ident{Name: tok.Str, Quoted: true, Loc: tok.Loc}, nil
	}
	return ast.Ident{}, &ParseError{
		Loc: p.cur.Loc,
		Msg: "expected identifier",
	}
}

// parseNamePart parses one component of a dotted object name. It is more
// permissive than parseIdent: it accepts any keyword token (including reserved
// keywords) in addition to bare and quoted identifiers. This matches Snowflake
// behaviour where reserved words such as SCHEMA or DATABASE are valid when
// used as object-name components in a dotted path (e.g. mydb.schema.t1).
func (p *Parser) parseNamePart() (ast.Ident, error) {
	switch {
	case p.cur.Type == tokIdent:
		tok := p.advance()
		return ast.Ident{Name: tok.Str, Quoted: false, Loc: tok.Loc}, nil
	case p.cur.Type == tokQuotedIdent:
		tok := p.advance()
		return ast.Ident{Name: tok.Str, Quoted: true, Loc: tok.Loc}, nil
	case p.cur.Type >= 700:
		// Any keyword — reserved or not — is accepted as a name part.
		tok := p.advance()
		return ast.Ident{Name: tok.Str, Quoted: false, Loc: tok.Loc}, nil
	}
	return ast.Ident{}, &ParseError{
		Loc: p.cur.Loc,
		Msg: "expected identifier",
	}
}

// parseObjectName parses a dotted object name with 1, 2, or 3 parts:
//
//	table
//	schema.table
//	database.schema.table
//
// Starts by parsing one identifier, then greedily consumes up to two more
// dot-separated identifiers if present. Returns *ast.ObjectName with the
// correct parts populated and Loc spanning all parts.
//
// Each name part is parsed with parseNamePart, which accepts any keyword
// (including reserved keywords such as SCHEMA or DATABASE) as an identifier.
//
// The first part may also be Snowflake's IDENTIFIER(<expr>) literal form —
// USE WAREHOUSE IDENTIFIER($wh), UNDROP TABLE IDENTIFIER(408578) — which is
// captured verbatim as a single (unquoted) Ident part; see
// parseIdentifierFuncPart.
func (p *Parser) parseObjectName() (*ast.ObjectName, error) {
	var first ast.Ident
	var err error
	if p.cur.Type == kwIDENTIFIER && p.peekNext().Type == '(' {
		first, err = p.parseIdentifierFuncPart()
	} else {
		first, err = p.parseNamePart()
	}
	if err != nil {
		return nil, err
	}

	if p.cur.Type != '.' {
		// 1-part name.
		return &ast.ObjectName{
			Name: first,
			Loc:  first.Loc,
		}, nil
	}
	p.advance() // consume first dot

	second, err := p.parseNamePart()
	if err != nil {
		return nil, err
	}

	if p.cur.Type != '.' {
		// 2-part name: schema.table.
		return &ast.ObjectName{
			Schema: first,
			Name:   second,
			Loc:    ast.Loc{Start: first.Loc.Start, End: second.Loc.End},
		}, nil
	}
	p.advance() // consume second dot

	third, err := p.parseNamePart()
	if err != nil {
		return nil, err
	}

	// 3-part name: database.schema.table.
	return &ast.ObjectName{
		Database: first,
		Schema:   second,
		Name:     third,
		Loc:      ast.Loc{Start: first.Loc.Start, End: third.Loc.End},
	}, nil
}

// parseIdentifierFuncPart parses Snowflake's IDENTIFIER(<expr>) name form,
// used wherever an object name is expected with the name supplied as a
// string / session variable / bind / object id:
//
//	USE WAREHOUSE IDENTIFIER($current_wh_name)
//	UNDROP TABLE IDENTIFIER(408578)
//	SELECT * FROM IDENTIFIER('db.sch.t')
//
// The whole construct is captured VERBATIM (source text, including the
// IDENTIFIER keyword and parentheses) as one unquoted Ident part, so
// consumers and deparse round-trip it unchanged; the argument's structure is
// validated by parsing it as an expression but not retained.
//
// The caller must have checked cur == kwIDENTIFIER and peekNext == '('.
func (p *Parser) parseIdentifierFuncPart() (ast.Ident, error) {
	startTok := p.advance() // consume IDENTIFIER
	p.advance()             // consume '('
	if _, err := p.parseExpr(); err != nil {
		return ast.Ident{}, err
	}
	closeTok, err := p.expect(')')
	if err != nil {
		return ast.Ident{}, err
	}
	loc := ast.Loc{Start: startTok.Loc.Start, End: closeTok.Loc.End}
	return ast.Ident{
		Name: p.srcSlice(loc.Start, loc.End),
		Loc:  loc,
	}, nil
}

// ParseObjectName parses an object name from a standalone string. Returns
// the ObjectName and any ParseErrors encountered. Useful for catalog
// lookups, tests, and callers that have a name string but not a token
// stream.
//
// Examples:
//
//	ParseObjectName("my_table")
//	ParseObjectName("schema.table")
//	ParseObjectName(`"My DB".schema."Quoted Table"`)
func ParseObjectName(input string) (*ast.ObjectName, []ParseError) {
	p := &Parser{
		lexer: NewLexer(input),
		input: input,
	}
	p.advance()

	obj, err := p.parseObjectName()
	if err != nil {
		if pe, ok := err.(*ParseError); ok {
			return nil, []ParseError{*pe}
		}
		return nil, []ParseError{{Msg: err.Error()}}
	}

	// Check for trailing tokens (should be EOF).
	if p.cur.Type != tokEOF {
		return obj, []ParseError{{
			Loc: p.cur.Loc,
			Msg: "unexpected token after object name",
		}}
	}

	return obj, nil
}
