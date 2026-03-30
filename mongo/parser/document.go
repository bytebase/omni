package parser

import "github.com/bytebase/omni/mongo/ast"

// parseDocument parses a { key: value, ... } document literal.
func (p *Parser) parseDocument() (*ast.Document, error) {
	openTok, err := p.expect('{')
	if err != nil {
		return nil, err
	}
	start := openTok.Loc

	var pairs []ast.KeyValue

	if p.cur.Type == '}' {
		closeTok := p.advance()
		return &ast.Document{
			Pairs: pairs,
			Loc:   ast.Loc{Start: start, End: closeTok.End},
		}, nil
	}

	for {
		kv, err := p.parsePair()
		if err != nil {
			return nil, err
		}
		pairs = append(pairs, kv)

		// Trailing comma support
		if p.cur.Type == ',' {
			p.advance()
		}

		if p.cur.Type == '}' {
			closeTok := p.advance()
			return &ast.Document{
				Pairs: pairs,
				Loc:   ast.Loc{Start: start, End: closeTok.End},
			}, nil
		}
	}
}

// parsePair parses a single key: value pair in a document.
// Keys can be: quoted strings, identifiers, keywords-as-identifiers, or $-prefixed identifiers.
func (p *Parser) parsePair() (ast.KeyValue, error) {
	var key string
	var keyLoc ast.Loc

	switch {
	case p.cur.Type == tokString:
		tok := p.advance()
		key = tok.Str
		keyLoc = ast.Loc{Start: tok.Loc, End: tok.End}
	case p.cur.Type == tokIdent:
		tok := p.advance()
		key = tok.Str
		keyLoc = ast.Loc{Start: tok.Loc, End: tok.End}
	case p.isIdentLike(p.cur.Type):
		tok := p.advance()
		key = tok.Str
		keyLoc = ast.Loc{Start: tok.Loc, End: tok.End}
	default:
		return ast.KeyValue{}, p.syntaxErrorAtCur()
	}

	if _, err := p.expect(':'); err != nil {
		return ast.KeyValue{}, err
	}

	value, err := p.parseExpression()
	if err != nil {
		return ast.KeyValue{}, err
	}

	return ast.KeyValue{
		Key:    key,
		KeyLoc: keyLoc,
		Value:  value,
	}, nil
}

// parseArray parses a [ elem, ... ] array literal.
func (p *Parser) parseArray() (ast.Node, error) {
	openTok, err := p.expect('[')
	if err != nil {
		return nil, err
	}
	start := openTok.Loc

	var elements []ast.Node

	if p.cur.Type == ']' {
		closeTok := p.advance()
		return &ast.Array{
			Elements: elements,
			Loc:      ast.Loc{Start: start, End: closeTok.End},
		}, nil
	}

	for {
		elem, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		elements = append(elements, elem)

		// Trailing comma support
		if p.cur.Type == ',' {
			p.advance()
		}

		if p.cur.Type == ']' {
			closeTok := p.advance()
			return &ast.Array{
				Elements: elements,
				Loc:      ast.Loc{Start: start, End: closeTok.End},
			}, nil
		}
	}
}
