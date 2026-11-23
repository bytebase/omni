package parser

import (
	"github.com/bytebase/omni/cassandra/ast"
)

// parseProperties parses properties (Parser.g:1750): property (AND
// property)*.
func (p *Parser) parseProperties() ([]*ast.TableOption, error) {
	seen := map[string]bool{}
	var opts []*ast.TableOption
	for {
		opt, err := p.parseProperty(seen)
		if err != nil {
			return nil, err
		}
		opts = append(opts, opt)
		if !p.match(tokAND) {
			return opts, nil
		}
	}
}

// parseProperty parses property (Parser.g:1754): name '=' (propertyValue |
// fullMapLiteral). seen holds the names given so far in the same WITH; a
// property given twice is refused while Cassandra parses
// (PropertyDefinitions.addProperty: "Multiple definitions for property").
func (p *Parser) parseProperty(seen map[string]bool) (*ast.TableOption, error) {
	start := p.curLoc()
	nameTok := p.cur
	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	if seen[identKey(name)] {
		return nil, p.errorAt(nameTok, "Multiple definitions for property '%s'", identKey(name))
	}
	seen[identKey(name)] = true
	if _, err := p.expect(tokEQ); err != nil {
		return nil, err
	}
	var value ast.ExprNode
	if p.cur.Type == tokLBRACE {
		value, err = p.parseOptionHash()
	} else {
		value, err = p.parsePropertyValue()
	}
	if err != nil {
		return nil, err
	}
	return &ast.TableOption{Name: name, Value: value, Loc: p.makeLoc(start)}, nil
}

// parsePropertyValue parses propertyValue (Parser.g:1759): a constant or an
// unreserved keyword.
func (p *Parser) parsePropertyValue() (ast.ExprNode, error) {
	if isLiteralStart(p.cur.Type) {
		return p.parseConstant()
	}
	if isUnreservedKeyword(p.cur) {
		tok := p.cur
		p.advance()
		return &ast.Identifier{Name: tok.Str, Loc: ast.Loc{Start: tok.Loc, End: tok.End}}, nil
	}
	return nil, p.errorf("expected property value, got %s", p.tokenDesc())
}

// isLiteralStart reports whether a token starts a constant (Parser.g:1527)
// that is a Constants.Literal: a string, a number, a boolean, a UUID, a blob,
// a duration, NaN or Infinity, with or without '-'. NULL is no such literal.
func isLiteralStart(typ int) bool {
	switch typ {
	case tokSTRING, tokCODEBLOCK, tokINTEGER, tokFLOAT, tokUUID, tokHEX,
		tokTRUE, tokFALSE, tokNAN, tokINFINITY, tokDURATIONLIT, tokMINUS:
		return true
	}
	return false
}

// parseOptionHash parses fullMapLiteral (Parser.g:1540) where a property
// takes it: { key : value, ... }. convertPropertyMap (Parser.g:98) refuses
// a key or value that is not a literal, and a key given twice.
func (p *Parser) parseOptionHash() (*ast.OptionHash, error) {
	start := p.curLoc()
	if _, err := p.expect(tokLBRACE); err != nil {
		return nil, err
	}

	var items []*ast.OptionHashItem

	// Empty hash.
	if p.cur.Type == tokRBRACE {
		p.advance()
		return &ast.OptionHash{Items: items, Loc: p.makeLoc(start)}, nil
	}

	seen := map[string]bool{}
	for {
		keyTok := p.cur
		item, err := p.parseOptionHashItem()
		if err != nil {
			return nil, err
		}
		key := p.literalText(item.Key)
		if seen[key] {
			return nil, p.errorAt(keyTok, "Multiple definition for property %s", key)
		}
		seen[key] = true
		items = append(items, item)
		if !p.match(tokCOMMA) {
			break
		}
	}

	if _, err := p.expect(tokRBRACE); err != nil {
		return nil, err
	}
	return &ast.OptionHash{Items: items, Loc: p.makeLoc(start)}, nil
}

// parseOptionHashItem parses a single key : value pair inside an option hash.
func (p *Parser) parseOptionHashItem() (*ast.OptionHashItem, error) {
	start := p.curLoc()
	key, err := p.parseOptionHashValue()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tokCOLON); err != nil {
		return nil, err
	}
	val, err := p.parseOptionHashValue()
	if err != nil {
		return nil, err
	}
	return &ast.OptionHashItem{Key: key, Value: val, Loc: p.makeLoc(start)}, nil
}

// parseOptionHashValue parses a key or value of an option hash, which
// convertPropertyMap takes only as a literal.
func (p *Parser) parseOptionHashValue() (ast.ExprNode, error) {
	if !isLiteralStart(p.cur.Type) {
		return nil, p.errorf("expected literal in property map, got %s", p.tokenDesc())
	}
	return p.parseConstant()
}

// literalText is a literal's raw text, which Cassandra compares map keys
// by: a string's content, any other literal as written.
func (p *Parser) literalText(e ast.ExprNode) string {
	switch v := e.(type) {
	case *ast.StringLit:
		return v.Val
	case *ast.CodeBlock:
		return v.Val
	}
	loc := e.GetLoc()
	if loc.Start >= 0 && loc.End <= len(p.source) && loc.Start <= loc.End {
		return p.source[loc.Start:loc.End]
	}
	return ""
}
