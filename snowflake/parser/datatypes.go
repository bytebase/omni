package parser

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// parseDataType parses a Snowflake data type. Called from type positions
// like CAST(x AS type), CREATE TABLE column definitions, function
// signatures, etc.
//
// Handles all forms from the legacy SnowflakeParser.g4 data_type rule.
func (p *Parser) parseDataType() (*ast.TypeName, error) {
	tok := p.cur

	switch tok.Type {
	// Integer types — no parameters.
	case kwINT, kwINTEGER, kwSMALLINT, kwTINYINT, kwBYTEINT, kwBIGINT:
		p.advance()
		return &ast.TypeName{Kind: ast.TypeInt, Name: tok.Str, VectorDim: -1, Loc: tok.Loc}, nil

	// Numeric types — optional (precision [, scale]).
	case kwNUMBER, kwNUMERIC, kwDECIMAL:
		p.advance()
		params, endLoc, err := p.parseOptionalTypeParams()
		if err != nil {
			return nil, err
		}
		loc := tok.Loc
		if endLoc.End > loc.End {
			loc.End = endLoc.End
		}
		return &ast.TypeName{Kind: ast.TypeNumber, Name: tok.Str, Params: params, VectorDim: -1, Loc: loc}, nil

	// Float types — no parameters. DOUBLE has PRECISION lookahead.
	case kwFLOAT, kwFLOAT4, kwFLOAT8, kwREAL:
		p.advance()
		return &ast.TypeName{Kind: ast.TypeFloat, Name: tok.Str, VectorDim: -1, Loc: tok.Loc}, nil

	case kwDOUBLE:
		p.advance()
		// Lookahead for PRECISION (comes as tokIdent since it's not a keyword).
		if p.cur.Type == tokIdent && strings.ToUpper(p.cur.Str) == "PRECISION" {
			precTok := p.advance()
			return &ast.TypeName{
				Kind:      ast.TypeFloat,
				Name:      "DOUBLE PRECISION",
				VectorDim: -1,
				Loc:       ast.Loc{Start: tok.Loc.Start, End: precTok.Loc.End},
			}, nil
		}
		return &ast.TypeName{Kind: ast.TypeFloat, Name: "DOUBLE", VectorDim: -1, Loc: tok.Loc}, nil

	// Boolean.
	case kwBOOLEAN:
		p.advance()
		return &ast.TypeName{Kind: ast.TypeBoolean, Name: tok.Str, VectorDim: -1, Loc: tok.Loc}, nil

	// Date.
	case kwDATE:
		p.advance()
		return &ast.TypeName{Kind: ast.TypeDate, Name: tok.Str, VectorDim: -1, Loc: tok.Loc}, nil

	// DateTime — optional (precision).
	case kwDATETIME:
		p.advance()
		params, endLoc, err := p.parseOptionalTypeParams()
		if err != nil {
			return nil, err
		}
		loc := tok.Loc
		if endLoc.End > loc.End {
			loc.End = endLoc.End
		}
		return &ast.TypeName{Kind: ast.TypeDateTime, Name: tok.Str, Params: params, VectorDim: -1, Loc: loc}, nil

	// Time — optional (precision).
	case kwTIME:
		p.advance()
		params, endLoc, err := p.parseOptionalTypeParams()
		if err != nil {
			return nil, err
		}
		loc := tok.Loc
		if endLoc.End > loc.End {
			loc.End = endLoc.End
		}
		return &ast.TypeName{Kind: ast.TypeTime, Name: tok.Str, Params: params, VectorDim: -1, Loc: loc}, nil

	// Timestamp variants — optional (precision).
	case kwTIMESTAMP:
		return p.parseTimestampType(ast.TypeTimestamp, tok)
	case kwTIMESTAMP_LTZ:
		return p.parseTimestampType(ast.TypeTimestampLTZ, tok)
	case kwTIMESTAMP_NTZ:
		return p.parseTimestampType(ast.TypeTimestampNTZ, tok)
	case kwTIMESTAMP_TZ:
		return p.parseTimestampType(ast.TypeTimestampTZ, tok)

	// Char types — optional (length). Lookahead for VARYING → TypeVarchar.
	case kwCHAR, kwNCHAR:
		p.advance()
		// Lookahead for VARYING (comes as tokIdent).
		if p.cur.Type == tokIdent && strings.ToUpper(p.cur.Str) == "VARYING" {
			varyTok := p.advance()
			name := tok.Str + " VARYING"
			params, endLoc, err := p.parseOptionalTypeParams()
			if err != nil {
				return nil, err
			}
			loc := ast.Loc{Start: tok.Loc.Start, End: varyTok.Loc.End}
			if endLoc.End > loc.End {
				loc.End = endLoc.End
			}
			return &ast.TypeName{Kind: ast.TypeVarchar, Name: name, Params: params, VectorDim: -1, Loc: loc}, nil
		}
		params, endLoc, err := p.parseOptionalTypeParams()
		if err != nil {
			return nil, err
		}
		loc := tok.Loc
		if endLoc.End > loc.End {
			loc.End = endLoc.End
		}
		return &ast.TypeName{Kind: ast.TypeChar, Name: tok.Str, Params: params, VectorDim: -1, Loc: loc}, nil

	case kwCHARACTER:
		p.advance()
		params, endLoc, err := p.parseOptionalTypeParams()
		if err != nil {
			return nil, err
		}
		loc := tok.Loc
		if endLoc.End > loc.End {
			loc.End = endLoc.End
		}
		return &ast.TypeName{Kind: ast.TypeChar, Name: tok.Str, Params: params, VectorDim: -1, Loc: loc}, nil

	// Varchar-family types — optional (length).
	case kwVARCHAR, kwNVARCHAR, kwNVARCHAR2, kwSTRING, kwTEXT:
		p.advance()
		params, endLoc, err := p.parseOptionalTypeParams()
		if err != nil {
			return nil, err
		}
		loc := tok.Loc
		if endLoc.End > loc.End {
			loc.End = endLoc.End
		}
		return &ast.TypeName{Kind: ast.TypeVarchar, Name: tok.Str, Params: params, VectorDim: -1, Loc: loc}, nil

	// Binary types — optional (length).
	case kwBINARY:
		p.advance()
		params, endLoc, err := p.parseOptionalTypeParams()
		if err != nil {
			return nil, err
		}
		loc := tok.Loc
		if endLoc.End > loc.End {
			loc.End = endLoc.End
		}
		return &ast.TypeName{Kind: ast.TypeBinary, Name: tok.Str, Params: params, VectorDim: -1, Loc: loc}, nil

	case kwVARBINARY:
		p.advance()
		params, endLoc, err := p.parseOptionalTypeParams()
		if err != nil {
			return nil, err
		}
		loc := tok.Loc
		if endLoc.End > loc.End {
			loc.End = endLoc.End
		}
		return &ast.TypeName{Kind: ast.TypeVarbinary, Name: tok.Str, Params: params, VectorDim: -1, Loc: loc}, nil

	// Semi-structured types.
	case kwVARIANT:
		p.advance()
		return &ast.TypeName{Kind: ast.TypeVariant, Name: tok.Str, VectorDim: -1, Loc: tok.Loc}, nil

	case kwOBJECT:
		p.advance()
		return &ast.TypeName{Kind: ast.TypeObject, Name: tok.Str, VectorDim: -1, Loc: tok.Loc}, nil

	case kwARRAY:
		p.advance()
		if p.cur.Type == '(' {
			openTok := p.advance() // consume (
			_ = openTok
			elem, err := p.parseDataType()
			if err != nil {
				return nil, err
			}
			closeTok, err := p.expect(')')
			if err != nil {
				return nil, err
			}
			return &ast.TypeName{
				Kind:        ast.TypeArray,
				Name:        "ARRAY",
				ElementType: elem,
				VectorDim:   -1,
				Loc:         ast.Loc{Start: tok.Loc.Start, End: closeTok.Loc.End},
			}, nil
		}
		return &ast.TypeName{Kind: ast.TypeArray, Name: "ARRAY", VectorDim: -1, Loc: tok.Loc}, nil

	// Geospatial types.
	case kwGEOGRAPHY:
		p.advance()
		return &ast.TypeName{Kind: ast.TypeGeography, Name: tok.Str, VectorDim: -1, Loc: tok.Loc}, nil

	case kwGEOMETRY:
		p.advance()
		return &ast.TypeName{Kind: ast.TypeGeometry, Name: tok.Str, VectorDim: -1, Loc: tok.Loc}, nil

	// VECTOR(element_type, dimensions).
	case kwVECTOR:
		p.advance()
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		elemType, err := p.parseVectorElementType()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(','); err != nil {
			return nil, err
		}
		dimTok := p.cur
		if dimTok.Type != tokInt {
			return nil, &ParseError{Loc: dimTok.Loc, Msg: "expected integer dimension for VECTOR"}
		}
		p.advance()
		closeTok, err := p.expect(')')
		if err != nil {
			return nil, err
		}
		return &ast.TypeName{
			Kind:        ast.TypeVector,
			Name:        "VECTOR",
			ElementType: elemType,
			VectorDim:   int(dimTok.Ival),
			Loc:         ast.Loc{Start: tok.Loc.Start, End: closeTok.Loc.End},
		}, nil
	}

	return nil, &ParseError{Loc: tok.Loc, Msg: "expected data type"}
}

// parseTimestampType handles TIMESTAMP, TIMESTAMP_LTZ, TIMESTAMP_NTZ,
// TIMESTAMP_TZ — all with optional (precision).
func (p *Parser) parseTimestampType(kind ast.TypeKind, tok Token) (*ast.TypeName, error) {
	p.advance()
	params, endLoc, err := p.parseOptionalTypeParams()
	if err != nil {
		return nil, err
	}
	loc := tok.Loc
	if endLoc.End > loc.End {
		loc.End = endLoc.End
	}
	return &ast.TypeName{Kind: kind, Name: tok.Str, Params: params, VectorDim: -1, Loc: loc}, nil
}

// parseOptionalTypeParams parses an optional parenthesized parameter list
// after a type keyword: (n) or (n, m). Returns nil if no opening paren.
// The returned ast.Loc is the span of the parenthesized expression (or
// NoLoc if no params).
func (p *Parser) parseOptionalTypeParams() ([]int, ast.Loc, error) {
	if p.cur.Type != '(' {
		return nil, ast.NoLoc(), nil
	}
	openTok := p.advance() // consume (

	if p.cur.Type != tokInt {
		return nil, ast.NoLoc(), &ParseError{Loc: p.cur.Loc, Msg: "expected integer type parameter"}
	}
	first := int(p.cur.Ival)
	p.advance()

	params := []int{first}

	if p.cur.Type == ',' {
		p.advance() // consume ,
		if p.cur.Type != tokInt {
			return nil, ast.NoLoc(), &ParseError{Loc: p.cur.Loc, Msg: "expected integer type parameter"}
		}
		second := int(p.cur.Ival)
		p.advance()
		params = append(params, second)
	}

	closeTok, err := p.expect(')')
	if err != nil {
		return nil, ast.NoLoc(), err
	}
	_ = openTok
	return params, ast.Loc{Start: openTok.Loc.Start, End: closeTok.Loc.End}, nil
}

// parseVectorElementType parses the element type for VECTOR, which is
// restricted to INT | INTEGER | FLOAT | FLOAT4 | FLOAT8 per the grammar.
func (p *Parser) parseVectorElementType() (*ast.TypeName, error) {
	tok := p.cur
	switch tok.Type {
	case kwINT, kwINTEGER:
		p.advance()
		return &ast.TypeName{Kind: ast.TypeInt, Name: tok.Str, VectorDim: -1, Loc: tok.Loc}, nil
	case kwFLOAT, kwFLOAT4, kwFLOAT8:
		p.advance()
		return &ast.TypeName{Kind: ast.TypeFloat, Name: tok.Str, VectorDim: -1, Loc: tok.Loc}, nil
	}
	return nil, &ParseError{Loc: tok.Loc, Msg: "expected VECTOR element type (INT, INTEGER, FLOAT, FLOAT4, or FLOAT8)"}
}

// ParseDataType parses a data type from a standalone string. Useful for
// tests and catalog integration.
func ParseDataType(input string) (*ast.TypeName, []ParseError) {
	p := &Parser{
		lexer: NewLexer(input),
		input: input,
	}
	p.advance()

	dt, err := p.parseDataType()
	if err != nil {
		if pe, ok := err.(*ParseError); ok {
			return nil, []ParseError{*pe}
		}
		return nil, []ParseError{{Msg: err.Error()}}
	}

	if p.cur.Type != tokEOF {
		return dt, []ParseError{{
			Loc: p.cur.Loc,
			Msg: "unexpected token after data type",
		}}
	}

	return dt, nil
}
