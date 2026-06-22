package parser

import (
	"strings"

	"github.com/bytebase/omni/cassandra/ast"
)

// ---------------------------------------------------------------------------
// Constants
// ---------------------------------------------------------------------------

// parseConstant parses a CQL constant literal: string, integer, float, uuid,
// hex, boolean (true/false), null, or code block ($$...$$).
func (p *Parser) parseConstant() (ast.ExprNode, error) {
	tok := p.cur
	switch tok.Type {
	case tokSTRING:
		p.advance()
		return &ast.StringLit{Val: tok.Str, Loc: ast.Loc{Start: tok.Loc, End: tok.End}}, nil

	case tokINTEGER:
		p.advance()
		return &ast.IntegerLit{Val: tok.Str, Loc: ast.Loc{Start: tok.Loc, End: tok.End}}, nil

	case tokFLOAT:
		p.advance()
		return &ast.FloatLit{Val: tok.Str, Loc: ast.Loc{Start: tok.Loc, End: tok.End}}, nil

	case tokUUID:
		p.advance()
		return &ast.UUIDLit{Val: tok.Str, Loc: ast.Loc{Start: tok.Loc, End: tok.End}}, nil

	case tokHEX:
		p.advance()
		return &ast.HexLit{Val: tok.Str, Loc: ast.Loc{Start: tok.Loc, End: tok.End}}, nil

	case tokTRUE:
		p.advance()
		return &ast.BoolLit{Val: true, Loc: ast.Loc{Start: tok.Loc, End: tok.End}}, nil

	case tokFALSE:
		p.advance()
		return &ast.BoolLit{Val: false, Loc: ast.Loc{Start: tok.Loc, End: tok.End}}, nil

	case tokNULL:
		p.advance()
		return &ast.NullLit{Loc: ast.Loc{Start: tok.Loc, End: tok.End}}, nil

	case tokCODEBLOCK:
		p.advance()
		return &ast.CodeBlock{Val: tok.Str, Loc: ast.Loc{Start: tok.Loc, End: tok.End}}, nil

	case tokMINUS:
		// Negative numeric literal: -123 or -1.5
		start := tok.Loc
		p.advance()
		next := p.cur
		switch next.Type {
		case tokINTEGER:
			p.advance()
			return &ast.IntegerLit{Val: "-" + next.Str, Loc: ast.Loc{Start: start, End: next.End}}, nil
		case tokFLOAT:
			p.advance()
			return &ast.FloatLit{Val: "-" + next.Str, Loc: ast.Loc{Start: start, End: next.End}}, nil
		default:
			return nil, p.errorf("expected number after '-', got %s", p.tokenDesc())
		}

	default:
		return nil, p.errorf("expected constant, got %s", p.tokenDesc())
	}
}

// ---------------------------------------------------------------------------
// Expressions
// ---------------------------------------------------------------------------

// parseExpression parses a CQL expression: constant, function call,
// identifier, or collection literal (map, set, list, tuple).  After parsing
// the primary expression it checks for index access (e.g. col[0]).
func (p *Parser) parseExpression() (ast.ExprNode, error) {
	var expr ast.ExprNode
	var err error

	switch {
	// Collection literals.
	case p.cur.Type == tokLBRACE || p.cur.Type == tokLBRACK:
		expr, err = p.parseCollectionLiteral()

	// Tuple literal: (val, val, ...)
	case p.cur.Type == tokLPAREN:
		start := p.curLoc()
		p.advance() // (
		elems, err2 := p.parseExpressionList()
		if err2 != nil {
			return nil, err2
		}
		if _, err2 = p.expect(tokRPAREN); err2 != nil {
			return nil, err2
		}
		expr = &ast.TupleLit{Elements: elems, Loc: p.makeLoc(start)}

	// Identifier or function call.
	case isIdentLike(p.cur.Type):
		// Look ahead: if followed by '(' it is a function call.
		if p.peekNext().Type == tokLPAREN {
			expr, err = p.parseFunctionCall()
		} else {
			expr, err = p.parseIdentifier()
		}

	// Everything else: try as a constant.
	default:
		expr, err = p.parseConstant()
	}

	if err != nil {
		return nil, err
	}

	// Post-fix: index access  col[idx]
	for p.cur.Type == tokLBRACK {
		start := expr.GetLoc().Start
		p.advance() // [
		idx, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(tokRBRACK); err != nil {
			return nil, err
		}
		expr = &ast.IndexAccess{Collection: expr, Index: idx, Loc: p.makeLoc(start)}
	}

	return expr, nil
}

// parseExpressionList parses a comma-separated list of expressions.
func (p *Parser) parseExpressionList() ([]ast.ExprNode, error) {
	first, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	list := []ast.ExprNode{first}
	for p.match(tokCOMMA) {
		next, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		list = append(list, next)
	}
	return list, nil
}

// ---------------------------------------------------------------------------
// Function calls
// ---------------------------------------------------------------------------

// builtinFuncToken maps keyword token types that represent built-in CQL
// functions to the canonical (lowercase) function name used in the AST.
var builtinFuncToken = map[int]string{
	tokNOW:               "now",
	tokUUID_KW:           "uuid",
	tokFROMJSON:          "fromjson",
	tokTOJSON:            "tojson",
	tokMINTIMEUUID:       "mintimeuuid",
	tokMAXTIMEUUID:       "maxtimeuuid",
	tokDATETIMENOW:       "datetimenow",
	tokCURRENTTIMESTAMP:  "currenttimestamp",
	tokCURRENTDATE:       "currentdate",
	tokCURRENTTIME:       "currenttime",
	tokCURRENTTIMEUUID:   "currenttimeuuid",
}

// parseFunctionCall parses a function call.  It handles both built-in
// function keywords (now(), uuid(), fromJson(), ...) and generic
// identifier-based calls (token(col), writetime(col), ...).
func (p *Parser) parseFunctionCall() (ast.ExprNode, error) {
	tok := p.cur

	// Built-in function keywords.
	if name, ok := builtinFuncToken[tok.Type]; ok {
		start := tok.Loc
		ident := &ast.Identifier{Name: name, Loc: ast.Loc{Start: tok.Loc, End: tok.End}}
		p.advance() // consume keyword

		fc, err := p.parseFunctionCallWithName(ident)
		if err != nil {
			return nil, err
		}
		fc.Loc = p.makeLoc(start)
		return fc, nil
	}

	// Generic function call: name(args...) or ks.name(args...).
	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	// Handle dotted function name: keyspace.func(...)
	if p.cur.Type == tokDOT {
		p.advance()
		second, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		// Combine into "ks.func" identifier preserving full location.
		combined := &ast.Identifier{
			Name: name.Name + "." + second.Name,
			Loc:  ast.Loc{Start: name.Loc.Start, End: second.Loc.End},
		}
		return p.parseFunctionCallWithName(combined)
	}

	return p.parseFunctionCallWithName(name)
}

// parseFunctionCallWithName parses the (args...) or (*) part of a function
// call, given the already-consumed function name identifier.
func (p *Parser) parseFunctionCallWithName(name *ast.Identifier) (*ast.FunctionCall, error) {
	start := name.Loc.Start
	if _, err := p.expect(tokLPAREN); err != nil {
		return nil, err
	}

	// count(*) or similar star argument.
	if p.cur.Type == tokSTAR {
		p.advance()
		if _, err := p.expect(tokRPAREN); err != nil {
			return nil, err
		}
		return &ast.FunctionCall{
			Name: name,
			Star: true,
			Loc:  p.makeLoc(start),
		}, nil
	}

	// Empty argument list: func()
	if p.cur.Type == tokRPAREN {
		p.advance()
		return &ast.FunctionCall{
			Name: name,
			Loc:  p.makeLoc(start),
		}, nil
	}

	args, err := p.parseFunctionArgs()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tokRPAREN); err != nil {
		return nil, err
	}
	return &ast.FunctionCall{
		Name: name,
		Args: args,
		Loc:  p.makeLoc(start),
	}, nil
}

// parseFunctionArgs parses comma-separated function arguments.  Each argument
// can be a constant, an identifier, a collection literal, or a nested
// function call.
func (p *Parser) parseFunctionArgs() ([]ast.ExprNode, error) {
	var args []ast.ExprNode
	for {
		arg, err := p.parseFunctionArg()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
		if !p.match(tokCOMMA) {
			break
		}
	}
	return args, nil
}

// parseFunctionArg parses a single function argument.
func (p *Parser) parseFunctionArg() (ast.ExprNode, error) {
	// Collection literals.
	if p.cur.Type == tokLBRACE || p.cur.Type == tokLBRACK {
		return p.parseCollectionLiteral()
	}

	// Tuple literal inside function args.
	if p.cur.Type == tokLPAREN {
		start := p.curLoc()
		p.advance() // (
		elems, err := p.parseExpressionList()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(tokRPAREN); err != nil {
			return nil, err
		}
		return &ast.TupleLit{Elements: elems, Loc: p.makeLoc(start)}, nil
	}

	// Identifier or nested function call.
	if isIdentLike(p.cur.Type) {
		if p.peekNext().Type == tokLPAREN {
			return p.parseFunctionCall()
		}
		return p.parseIdentifier()
	}

	// Constant.
	return p.parseConstant()
}

// ---------------------------------------------------------------------------
// Collection literals
// ---------------------------------------------------------------------------

// parseCollectionLiteral dispatches to map/set literal for { ... } or list
// literal for [ ... ].
func (p *Parser) parseCollectionLiteral() (ast.ExprNode, error) {
	switch p.cur.Type {
	case tokLBRACE:
		return p.parseMapOrSetLiteral()
	case tokLBRACK:
		return p.parseListLiteral()
	default:
		return nil, p.errorf("expected '{' or '[' for collection literal, got %s", p.tokenDesc())
	}
}

// parseMapOrSetLiteral parses a { ... } literal.  After consuming the
// opening brace it looks ahead: if the first element is followed by a colon,
// it is a map literal {k:v, ...}; otherwise it is a set literal {v, ...}.
// An empty {} is treated as an empty map.
func (p *Parser) parseMapOrSetLiteral() (ast.ExprNode, error) {
	start := p.curLoc()
	p.advance() // {

	// Empty braces: empty map.
	if p.cur.Type == tokRBRACE {
		p.advance()
		return &ast.MapLit{Loc: p.makeLoc(start)}, nil
	}

	// Parse the first element to decide map vs set.
	first, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if p.cur.Type == tokCOLON {
		// Map literal.
		return p.finishMapLiteral(start, first)
	}

	// Set literal.
	return p.finishSetLiteral(start, first)
}

// finishMapLiteral continues parsing a map literal after the first key has
// been consumed and the colon has been seen (but not consumed).
func (p *Parser) finishMapLiteral(start int, firstKey ast.ExprNode) (ast.ExprNode, error) {
	p.advance() // consume :
	firstVal, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	keys := []ast.ExprNode{firstKey}
	values := []ast.ExprNode{firstVal}

	for p.match(tokCOMMA) {
		k, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(tokCOLON); err != nil {
			return nil, err
		}
		v, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		keys = append(keys, k)
		values = append(values, v)
	}

	if _, err := p.expect(tokRBRACE); err != nil {
		return nil, err
	}
	return &ast.MapLit{Keys: keys, Values: values, Loc: p.makeLoc(start)}, nil
}

// finishSetLiteral continues parsing a set literal after the first element
// has been consumed.
func (p *Parser) finishSetLiteral(start int, firstElem ast.ExprNode) (ast.ExprNode, error) {
	elems := []ast.ExprNode{firstElem}

	for p.match(tokCOMMA) {
		e, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		elems = append(elems, e)
	}

	if _, err := p.expect(tokRBRACE); err != nil {
		return nil, err
	}
	return &ast.SetLit{Elements: elems, Loc: p.makeLoc(start)}, nil
}

// parseListLiteral parses a [ elem, elem, ... ] list literal.
func (p *Parser) parseListLiteral() (ast.ExprNode, error) {
	start := p.curLoc()
	p.advance() // [

	// Empty list.
	if p.cur.Type == tokRBRACK {
		p.advance()
		return &ast.ListLit{Loc: p.makeLoc(start)}, nil
	}

	elems, err := p.parseExpressionList()
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(tokRBRACK); err != nil {
		return nil, err
	}
	return &ast.ListLit{Elements: elems, Loc: p.makeLoc(start)}, nil
}

// ---------------------------------------------------------------------------
// Option hash (used in table / index options)
// ---------------------------------------------------------------------------

// parseOptionHash parses a { 'key' : 'value', ... } option hash used in
// WITH clauses for table and index options.
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

	for {
		item, err := p.parseOptionHashItem()
		if err != nil {
			return nil, err
		}
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
// Keys and values are typically string literals but may also be identifiers,
// integers, floats, or booleans.
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

// parseOptionHashValue parses a value that can appear as a key or value in an
// option hash.  This is more permissive than parseConstant because option
// hashes accept unquoted identifiers as values (e.g. class : 'SimpleStrategy').
func (p *Parser) parseOptionHashValue() (ast.ExprNode, error) {
	tok := p.cur
	switch {
	case tok.Type == tokSTRING:
		p.advance()
		return &ast.StringLit{Val: tok.Str, Loc: ast.Loc{Start: tok.Loc, End: tok.End}}, nil
	case tok.Type == tokINTEGER:
		p.advance()
		return &ast.IntegerLit{Val: tok.Str, Loc: ast.Loc{Start: tok.Loc, End: tok.End}}, nil
	case tok.Type == tokFLOAT:
		p.advance()
		return &ast.FloatLit{Val: tok.Str, Loc: ast.Loc{Start: tok.Loc, End: tok.End}}, nil
	case tok.Type == tokTRUE:
		p.advance()
		return &ast.BoolLit{Val: true, Loc: ast.Loc{Start: tok.Loc, End: tok.End}}, nil
	case tok.Type == tokFALSE:
		p.advance()
		return &ast.BoolLit{Val: false, Loc: ast.Loc{Start: tok.Loc, End: tok.End}}, nil
	case tok.Type == tokNULL:
		p.advance()
		return &ast.NullLit{Loc: ast.Loc{Start: tok.Loc, End: tok.End}}, nil
	case isIdentLike(tok.Type):
		p.advance()
		return &ast.Identifier{
			Name:   strings.ToLower(tok.Str),
			Quoted: tok.Type == tokQUOTED,
			Loc:    ast.Loc{Start: tok.Loc, End: tok.End},
		}, nil
	default:
		return nil, p.errorf("expected value in option hash, got %s", p.tokenDesc())
	}
}
