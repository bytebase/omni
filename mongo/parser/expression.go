package parser

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/mongo/ast"
)

// parseExpression parses a value expression (document, array, literal, helper, identifier).
func (p *Parser) parseExpression() (ast.Node, error) {
	return p.parseValue()
}

// parseValue parses a value expression (document, array, literal, helper, identifier).
func (p *Parser) parseValue() (ast.Node, error) {
	switch p.cur.Type {
	case '{':
		return p.parseDocument()
	case '[':
		return p.parseArray()

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
