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

// parseObjectName parses a dotted object name with 1, 2, or 3 parts:
//
//	table
//	schema.table
//	database.schema.table
//
// Starts by parsing one identifier, then greedily consumes up to two more
// dot-separated identifiers if present. Returns *ast.ObjectName with the
// correct parts populated and Loc spanning all parts.
func (p *Parser) parseObjectName() (*ast.ObjectName, error) {
	first, err := p.parseIdent()
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

	second, err := p.parseIdent()
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

	third, err := p.parseIdent()
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
