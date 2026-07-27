package parser

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/bytebase/omni/mongo/ast"
)

// parseExpression parses a value expression (document, array, literal,
// helper, identifier), including postfix .method(args) chains and constant
// arithmetic over literals. mongosh evaluates JavaScript, so expressions
// like 90 * 24 * 60 * 60 are legal wherever a value fits; this parser folds
// the pure-literal subset at parse time with JavaScript number semantics
// (all arithmetic in float64) and rejects everything else, keeping the
// accepted surface a closed whitelist:
//
//   - binary * / % + - over number literals, ( ) grouping, unary + -
//   - string + string concatenation
//   - operands that are not literals (identifiers, calls) do not fold and
//     fail exactly like before
//   - results that are Infinity or NaN are rejected: BSON can carry them
//     but a folded Infinity in DDL options is always a script bug, and the
//     literal text form cannot round-trip them
//
// The folded literal's text form encodes the BSON type downstream
// (mongosh-verified): integral results within int32 keep integer form
// (-> int32), everything else takes a float form with '.' or 'e'
// (-> double), matching what mongosh itself sends.
func (p *Parser) parseExpression() (ast.Node, error) {
	return p.parseAdditive()
}

// parseAdditive parses + and - over multiplicative expressions, folding
// literal operands. The lexer folds a sign into a following number, so
// "5 -2" arrives as two adjacent number tokens; an adjacent sign-prefixed
// number in operand position is JavaScript subtraction/addition and is
// treated as the corresponding binary operation.
func (p *Parser) parseAdditive() (ast.Node, error) {
	left, err := p.parseMultiplicative()
	if err != nil {
		return nil, err
	}
	for {
		var op byte
		switch {
		case p.cur.Type == '+':
			op = '+'
			p.advance()
		case p.cur.Type == '-':
			op = '-'
			p.advance()
		case p.cur.Type == tokNumber && len(p.cur.Str) > 0 && (p.cur.Str[0] == '-' || p.cur.Str[0] == '+'):
			// Adjacent signed number: implicit addition of a signed term.
			op = '+'
		default:
			return left, nil
		}
		right, err := p.parseMultiplicative()
		if err != nil {
			return nil, err
		}
		left, err = p.foldBinary(op, left, right)
		if err != nil {
			return nil, err
		}
	}
}

// parseMultiplicative parses *, / and % over unary expressions, folding
// literal operands.
func (p *Parser) parseMultiplicative() (ast.Node, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}
	for p.cur.Type == '*' || p.cur.Type == '/' || p.cur.Type == '%' {
		op := byte(p.cur.Type)
		p.advance()
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		left, err = p.foldBinary(op, left, right)
		if err != nil {
			return nil, err
		}
	}
	return left, nil
}

// parseUnary parses an optional +/- sign applied to a unary expression,
// then a postfix expression. The sign only folds over number literals.
func (p *Parser) parseUnary() (ast.Node, error) {
	if p.cur.Type == '-' || p.cur.Type == '+' {
		opTok := p.cur
		neg := p.cur.Type == '-'
		p.advance()
		operand, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		num, ok := operand.(*ast.NumberLiteral)
		if !ok {
			return nil, p.arithErrorAt(opTok.Loc, string(rune(opTok.Type)))
		}
		v, err := numberValue(num)
		if err != nil {
			return nil, p.arithErrorAt(opTok.Loc, string(rune(opTok.Type)))
		}
		if neg {
			v = -v
		}
		return p.foldedNumber(v, ast.Loc{Start: opTok.Loc, End: num.Loc.End})
	}
	return p.parsePostfix()
}

// foldBinary folds one binary arithmetic operation over literal operands.
// Number op number uses JavaScript semantics (float64 throughout,
// math.Mod remainder). "+" additionally concatenates string literals.
// Any other operand combination is a syntax error at the current position:
// the folding surface stays a closed literal-only whitelist, and
// JavaScript's mixed-type coercion table is deliberately not modeled.
func (p *Parser) foldBinary(op byte, left, right ast.Node) (ast.Node, error) {
	loc := ast.Loc{Start: left.GetLoc().Start, End: right.GetLoc().End}

	if op == '+' {
		if ls, ok := left.(*ast.StringLiteral); ok {
			rs, ok := right.(*ast.StringLiteral)
			if !ok {
				return nil, p.syntaxErrorAtCur()
			}
			return &ast.StringLiteral{Value: ls.Value + rs.Value, Loc: loc}, nil
		}
	}

	ln, ok := left.(*ast.NumberLiteral)
	if !ok {
		return nil, p.syntaxErrorAtCur()
	}
	rn, ok := right.(*ast.NumberLiteral)
	if !ok {
		return nil, p.syntaxErrorAtCur()
	}
	lv, err := numberValue(ln)
	if err != nil {
		return nil, p.syntaxErrorAtCur()
	}
	rv, err := numberValue(rn)
	if err != nil {
		return nil, p.syntaxErrorAtCur()
	}

	var v float64
	switch op {
	case '+':
		v = lv + rv
	case '-':
		v = lv - rv
	case '*':
		v = lv * rv
	case '/':
		v = lv / rv
	case '%':
		v = math.Mod(lv, rv)
	default:
		return nil, p.syntaxErrorAtCur()
	}
	return p.foldedNumber(v, loc)
}

// arithErrorAt returns a ParseError for an arithmetic operator applied to
// a non-literal operand.
func (p *Parser) arithErrorAt(pos int, opText string) *ParseError {
	line, col := p.lineCol(pos)
	return &ParseError{
		Message:  fmt.Sprintf("syntax error at or near %q", opText),
		Position: pos,
		Line:     line,
		Column:   col,
	}
}

// numberValue evaluates a number literal as a JavaScript number (float64).
func numberValue(n *ast.NumberLiteral) (float64, error) {
	return strconv.ParseFloat(n.Value, 64)
}

// foldedNumber renders a folded value back to a literal whose text form
// selects the same BSON type mongosh would send: integral within int32
// keeps integer form (-> int32); everything else takes a float form with
// '.' or 'e' (-> double). Infinity and NaN are rejected.
func (p *Parser) foldedNumber(v float64, loc ast.Loc) (*ast.NumberLiteral, error) {
	if math.IsInf(v, 0) || math.IsNaN(v) {
		line, col := p.lineCol(loc.Start)
		return nil, &ParseError{
			Message:  "arithmetic result is not a finite number",
			Position: loc.Start,
			Line:     line,
			Column:   col,
		}
	}
	if v == math.Trunc(v) && v >= math.MinInt32 && v <= math.MaxInt32 {
		return &ast.NumberLiteral{
			Value:   strconv.FormatInt(int64(v), 10),
			IsFloat: false,
			Loc:     loc,
		}, nil
	}
	text := strconv.FormatFloat(v, 'g', -1, 64)
	if !strings.ContainsAny(text, ".eE") {
		text += ".0"
	}
	return &ast.NumberLiteral{
		Value:   text,
		IsFloat: true,
		Loc:     loc,
	}, nil
}

// parsePostfix parses a primary value with optional postfix .method(args)
// chains (e.g., Binary.createFromBase64("...")).
func (p *Parser) parsePostfix() (ast.Node, error) {
	node, err := p.parseValue()
	if err != nil {
		return nil, err
	}

	// Handle postfix .method(args) chains on identifiers/helpers
	for p.cur.Type == '.' {
		// Check that what follows is identifier + '(' (a method call)
		next := p.peekNext()
		if !p.isIdentLike(next.Type) {
			break
		}
		p.advance() // consume '.'
		mc, err := p.parseMethodCall()
		if err != nil {
			return nil, err
		}
		// Wrap as a HelperCall with qualified name
		start := node.GetLoc().Start
		name := nodeToName(node) + "." + mc.Name
		node = &ast.HelperCall{
			Name: name,
			Args: mc.Args,
			Loc:  ast.Loc{Start: start, End: mc.Loc.End},
		}
	}

	return node, nil
}

// nodeToName extracts a name string from an AST node for building qualified names.
func nodeToName(n ast.Node) string {
	switch v := n.(type) {
	case *ast.Identifier:
		return v.Name
	case *ast.HelperCall:
		return v.Name
	default:
		return "?"
	}
}

// parseValue parses a value expression (document, array, literal, helper, identifier).
func (p *Parser) parseValue() (ast.Node, error) {
	switch p.cur.Type {
	case '{':
		return p.parseDocument()
	case '[':
		return p.parseArray()

	case '(':
		// Parenthesized constant expression: (1 + 2) * 3.
		p.advance()
		inner, err := p.parseAdditive()
		if err != nil {
			return nil, err
		}
		if p.cur.Type != ')' {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		return inner, nil

	case tokRegex:
		tok := p.advance()
		pattern, flags := splitRegex(tok.Str)
		return &ast.RegexLiteral{
			Pattern: pattern,
			Flags:   flags,
			Loc:     ast.Loc{Start: tok.Loc, End: tok.End},
		}, nil

	case tokString:
		tok := p.advance()
		return &ast.StringLiteral{
			Value: tok.Str,
			Loc:   ast.Loc{Start: tok.Loc, End: tok.End},
		}, nil

	case tokNumber:
		tok := p.advance()
		isFloat := strings.ContainsAny(tok.Str, ".eE")
		return &ast.NumberLiteral{
			Value:   tok.Str,
			IsFloat: isFloat,
			Loc:     ast.Loc{Start: tok.Loc, End: tok.End},
		}, nil

	case kwTrue:
		tok := p.advance()
		return &ast.BoolLiteral{Value: true, Loc: ast.Loc{Start: tok.Loc, End: tok.End}}, nil

	case kwFalse:
		tok := p.advance()
		return &ast.BoolLiteral{Value: false, Loc: ast.Loc{Start: tok.Loc, End: tok.End}}, nil

	case kwNull:
		tok := p.advance()
		return &ast.NullLiteral{Loc: ast.Loc{Start: tok.Loc, End: tok.End}}, nil

	case kwNew:
		line, col := p.lineCol(p.cur.Loc)
		return nil, &ParseError{
			Message:  `"new" keyword is not supported`,
			Position: p.cur.Loc,
			Line:     line,
			Column:   col,
		}

	// BSON helper functions
	case kwObjectId, kwISODate, kwDate, kwUUID,
		kwNumberLong, kwNumberInt, kwNumberDecimal,
		kwTimestamp, kwRegExp, kwBinData, kwHexData,
		kwMinKey, kwMaxKey, kwCode, kwDBRef, kwSymbol, kwMD5:
		return p.parseHelperCall()

	default:
		// Handle identifier-like tokens that could be helper names
		// (e.g. Long, Int32, Double, Decimal128, Binary, BSONRegExp that aren't keywords)
		if p.isIdentLike(p.cur.Type) && p.peekNext().Type == '(' {
			return p.parseHelperCall()
		}

		// Plain identifier reference
		if p.isIdentLike(p.cur.Type) {
			tok := p.advance()
			return &ast.Identifier{
				Name: tok.Str,
				Loc:  ast.Loc{Start: tok.Loc, End: tok.End},
			}, nil
		}

		return nil, p.syntaxErrorAtCur()
	}
}

// parseHelperCall parses a BSON helper call like ObjectId("..."), NumberLong(1), etc.
func (p *Parser) parseHelperCall() (*ast.HelperCall, error) {
	nameTok := p.advance()
	start := nameTok.Loc

	args, err := p.parseArguments()
	if err != nil {
		return nil, err
	}

	return &ast.HelperCall{
		Name: nameTok.Str,
		Args: args,
		Loc:  ast.Loc{Start: start, End: p.prev.End},
	}, nil
}

// splitRegex splits a regex token string "pattern/flags" into pattern and flags.
// The lexer stores regex tokens as "pattern/flags" (without the leading /).
func splitRegex(s string) (string, string) {
	idx := strings.LastIndex(s, "/")
	if idx < 0 {
		return s, ""
	}
	return s[:idx], s[idx+1:]
}

// errorAtCur creates a ParseError with a custom message at the current token position.
func (p *Parser) errorAtCur(format string, args ...any) *ParseError {
	line, col := p.lineCol(p.cur.Loc)
	return &ParseError{
		Message:  fmt.Sprintf(format, args...),
		Position: p.cur.Loc,
		Line:     line,
		Column:   col,
	}
}
