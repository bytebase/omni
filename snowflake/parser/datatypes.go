package parser

import (
	"github.com/bytebase/omni/snowflake/ast"
)

// parseDataType parses a Snowflake data type. Called from type positions
// like CAST(x AS type), CREATE TABLE column definitions, function
// signatures, etc.
//
// Handles all forms from the legacy SnowflakeParser.g4 data_type rule.
func (p *Parser) parseDataType() (*ast.TypeName, error) {
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()
	tok := p.cur

	switch {
	// Integer types — no parameters.
	case p.curIs(kwINT), p.curIs(kwINTEGER), p.curIs(kwSMALLINT), p.curIs(kwTINYINT), p.curIs(kwBYTEINT), p.curIs(kwBIGINT):
		p.advance()
		return &ast.TypeName{Kind: ast.TypeInt, Name: tok.Str, VectorDim: -1, Loc: tok.Loc}, nil

	// Numeric types — optional (precision [, scale]).
	case p.curIs(kwNUMBER), p.curIs(kwNUMERIC), p.curIs(kwDECIMAL):
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
	case p.curIs(kwFLOAT), p.curIs(kwFLOAT4), p.curIs(kwFLOAT8), p.curIs(kwREAL):
		p.advance()
		return &ast.TypeName{Kind: ast.TypeFloat, Name: tok.Str, VectorDim: -1, Loc: tok.Loc}, nil

	case p.curIs(kwDOUBLE):
		p.advance()
		// Lookahead for PRECISION (comes as tokIdent since it's not a keyword).
		if p.curIsUnquotedWord("PRECISION") {
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
	case p.curIs(kwBOOLEAN):
		p.advance()
		return &ast.TypeName{Kind: ast.TypeBoolean, Name: tok.Str, VectorDim: -1, Loc: tok.Loc}, nil

	// Date.
	case p.curIs(kwDATE):
		p.advance()
		return &ast.TypeName{Kind: ast.TypeDate, Name: tok.Str, VectorDim: -1, Loc: tok.Loc}, nil

	// DateTime — optional (precision).
	case p.curIs(kwDATETIME):
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
	case p.curIs(kwTIME):
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
	case p.curIs(kwTIMESTAMP):
		return p.parseTimestampType(ast.TypeTimestamp, tok)
	case p.curIs(kwTIMESTAMP_LTZ):
		return p.parseTimestampType(ast.TypeTimestampLTZ, tok)
	case p.curIs(kwTIMESTAMP_NTZ):
		return p.parseTimestampType(ast.TypeTimestampNTZ, tok)
	case p.curIs(kwTIMESTAMP_TZ):
		return p.parseTimestampType(ast.TypeTimestampTZ, tok)

	// Char types — optional (length). Lookahead for VARYING → TypeVarchar.
	case p.curIs(kwCHAR), p.curIs(kwNCHAR):
		p.advance()
		// Lookahead for VARYING (comes as tokIdent).
		if p.curIsUnquotedWord("VARYING") {
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

	case p.curIs(kwCHARACTER):
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
	case p.curIs(kwVARCHAR), p.curIs(kwNVARCHAR), p.curIs(kwNVARCHAR2), p.curIs(kwSTRING), p.curIs(kwTEXT):
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
	case p.curIs(kwBINARY):
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

	case p.curIs(kwVARBINARY):
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
	case p.curIs(kwVARIANT):
		p.advance()
		return &ast.TypeName{Kind: ast.TypeVariant, Name: tok.Str, VectorDim: -1, Loc: tok.Loc}, nil

	case p.curIs(kwOBJECT):
		p.advance()
		return &ast.TypeName{Kind: ast.TypeObject, Name: tok.Str, VectorDim: -1, Loc: tok.Loc}, nil

	case p.curIs(kwARRAY):
		p.advance()
		if p.curIs('(') {
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
	case p.curIs(kwGEOGRAPHY):
		p.advance()
		return &ast.TypeName{Kind: ast.TypeGeography, Name: tok.Str, VectorDim: -1, Loc: tok.Loc}, nil

	case p.curIs(kwGEOMETRY):
		p.advance()
		return &ast.TypeName{Kind: ast.TypeGeometry, Name: tok.Str, VectorDim: -1, Loc: tok.Loc}, nil

	// VECTOR(element_type, dimensions).
	case p.curIs(kwVECTOR):
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
		if !p.curIs(tokInt) {
			return nil, &ParseError{Position: dimTok.Loc.Start, End: dimTok.Loc.End, Message: "expected integer dimension for VECTOR"}
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

	return nil, &ParseError{Position: tok.Loc.Start, End: tok.Loc.End, Message: "expected data type"}
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
	if !p.curIs('(') {
		return nil, ast.NoLoc(), nil
	}
	openTok := p.advance() // consume (

	if !p.curIs(tokInt) {
		return nil, ast.NoLoc(), &ParseError{Position: p.cur.Loc.Start, End: p.cur.Loc.End, Message: "expected integer type parameter"}
	}
	first := int(p.cur.Ival)
	p.advance()

	params := []int{first}

	if p.curIs(',') {
		p.advance() // consume ,
		if !p.curIs(tokInt) {
			return nil, ast.NoLoc(), &ParseError{Position: p.cur.Loc.Start, End: p.cur.Loc.End, Message: "expected integer type parameter"}
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
	switch {
	case p.curIs(kwINT), p.curIs(kwINTEGER):
		p.advance()
		return &ast.TypeName{Kind: ast.TypeInt, Name: tok.Str, VectorDim: -1, Loc: tok.Loc}, nil
	case p.curIs(kwFLOAT), p.curIs(kwFLOAT4), p.curIs(kwFLOAT8):
		p.advance()
		return &ast.TypeName{Kind: ast.TypeFloat, Name: tok.Str, VectorDim: -1, Loc: tok.Loc}, nil
	}
	return nil, &ParseError{Position: tok.Loc.Start, End: tok.Loc.End, Message: "expected VECTOR element type (INT, INTEGER, FLOAT, FLOAT4, or FLOAT8)"}
}

// ParseDataType parses a data type from a standalone string. Useful for
// tests and catalog integration.
func ParseDataType(input string) (*ast.TypeName, []ParseError) {
	p := &Parser{
		lexer: NewLexer(input),
		input: input,
	}
	p.advance()

	first := p.cur.Loc
	dt, err := p.parseDataType()
	if err == nil && treeTooHigh(dt, p.prev.Loc.End-first.Start) {
		dt, err = nil, nestedTooDeeply(first)
	}
	if err != nil {
		if pe, ok := err.(*ParseError); ok {
			return nil, []ParseError{*pe}
		}
		return nil, []ParseError{{Message: err.Error()}}
	}

	if !p.curIs(tokEOF) {
		return dt, []ParseError{{
			Position: p.cur.Loc.Start, End: p.cur.Loc.End,
			Message: "unexpected token after data type",
		}}
	}

	return dt, nil
}
