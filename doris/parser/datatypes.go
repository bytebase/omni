package parser

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/doris/ast"
)

// isPrimitiveTypeKeyword reports whether kind is a keyword that names a
// primitive Doris column type (i.e., usable as a standalone type name without
// angle-bracket parameters).
func isPrimitiveTypeKeyword(kind TokenKind) bool {
	switch kind {
	// Integer types
	case kwTINYINT, kwSMALLINT, kwINT, kwINTEGER, kwBIGINT, kwLARGEINT:
		return true
	// Float/decimal types
	case kwFLOAT, kwDOUBLE, kwDECIMAL, kwDECIMALV2, kwDECIMALV3:
		return true
	// Boolean
	case kwBOOLEAN:
		return true
	// String types
	case kwSTRING, kwTEXT, kwVARCHAR, kwCHAR, kwVARBINARY:
		return true
	// Date/time types
	case kwDATE, kwDATETIME, kwTIME, kwDATEV1, kwDATEV2, kwDATETIMEV1, kwDATETIMEV2:
		return true
	// Special Doris types
	case kwBITMAP, kwHLL, kwQUANTILE_STATE, kwJSON, kwJSONB:
		return true
	case kwIPV4, kwIPV6, kwVARIANT, kwAGG_STATE:
		return true
	// ALL keyword can appear as a type name in some grammar positions
	case kwALL:
		return true
	default:
		return false
	}
}

// typeKeywordName returns the canonical uppercase name for a type keyword.
// Falls back to the token's Str field (which may be lowercase from the lexer)
// uppercased.
func typeKeywordName(tok Token) string {
	switch tok.Kind {
	case kwTINYINT:
		return "TINYINT"
	case kwSMALLINT:
		return "SMALLINT"
	case kwINT:
		return "INT"
	case kwINTEGER:
		return "INTEGER"
	case kwBIGINT:
		return "BIGINT"
	case kwLARGEINT:
		return "LARGEINT"
	case kwFLOAT:
		return "FLOAT"
	case kwDOUBLE:
		return "DOUBLE"
	case kwDECIMAL:
		return "DECIMAL"
	case kwDECIMALV2:
		return "DECIMALV2"
	case kwDECIMALV3:
		return "DECIMALV3"
	case kwBOOLEAN:
		return "BOOLEAN"
	case kwSTRING:
		return "STRING"
	case kwTEXT:
		return "TEXT"
	case kwVARCHAR:
		return "VARCHAR"
	case kwCHAR:
		return "CHAR"
	case kwVARBINARY:
		return "VARBINARY"
	case kwDATE:
		return "DATE"
	case kwDATETIME:
		return "DATETIME"
	case kwTIME:
		return "TIME"
	case kwDATEV1:
		return "DATEV1"
	case kwDATEV2:
		return "DATEV2"
	case kwDATETIMEV1:
		return "DATETIMEV1"
	case kwDATETIMEV2:
		return "DATETIMEV2"
	case kwBITMAP:
		return "BITMAP"
	case kwHLL:
		return "HLL"
	case kwQUANTILE_STATE:
		return "QUANTILE_STATE"
	case kwJSON:
		return "JSON"
	case kwJSONB:
		return "JSONB"
	case kwIPV4:
		return "IPV4"
	case kwIPV6:
		return "IPV6"
	case kwVARIANT:
		return "VARIANT"
	case kwAGG_STATE:
		return "AGG_STATE"
	case kwALL:
		return "ALL"
	case kwARRAY:
		return "ARRAY"
	case kwMAP:
		return "MAP"
	case kwSTRUCT:
		return "STRUCT"
	default:
		return strings.ToUpper(tok.Str)
	}
}

// parseDataType parses a Doris data type expression and returns an *ast.TypeName.
//
// Supported forms:
//
//	primitive_type                    — e.g. INT, VARCHAR, BOOLEAN
//	primitive_type(N)                 — e.g. VARCHAR(255), CHAR(1)
//	primitive_type(N, M)              — e.g. DECIMAL(10,2)
//	ARRAY<elementType>
//	MAP<keyType, valueType>
//	STRUCT<name1 type1, name2 type2, ...>
func (p *Parser) parseDataType() (*ast.TypeName, error) {
	tok := p.cur

	switch tok.Kind {
	case kwARRAY:
		return p.parseArrayType()
	case kwMAP:
		return p.parseMapType()
	case kwSTRUCT:
		return p.parseStructType()
	default:
		if isPrimitiveTypeKeyword(tok.Kind) {
			return p.parsePrimitiveType()
		}
		// Also accept unquoted identifiers as type names (for forward
		// compatibility with new type names not yet in the keyword list).
		if tok.Kind == tokIdent || tok.Kind == tokQuotedIdent {
			return p.parsePrimitiveType()
		}
		return nil, &ParseError{
			Loc: tok.Loc,
			Msg: fmt.Sprintf("expected a data type, got %q", tok.Str),
		}
	}
}

// parsePrimitiveType parses a primitive type name and its optional numeric
// parameter list: type_name [ '(' N [ ',' M ] ')' ]
func (p *Parser) parsePrimitiveType() (*ast.TypeName, error) {
	tok := p.advance()
	typName := &ast.TypeName{
		Name: typeKeywordName(tok),
		Loc:  tok.Loc,
	}

	// Optional parameter list: (N) or (N, M)
	if p.cur.Kind == int('(') {
		params, endLoc, err := p.parseTypeParams()
		if err != nil {
			return nil, err
		}
		typName.Params = params
		typName.Loc.End = endLoc
	}

	return typName, nil
}

// parseTypeParams parses the optional numeric parameter list for a type:
//
//	'(' integer [',' integer] ')'
//
// Returns the list of param values and the End byte offset of the closing ')'.
func (p *Parser) parseTypeParams() ([]int, int, error) {
	// consume '('
	if _, err := p.expect(int('(')); err != nil {
		return nil, 0, err
	}

	firstTok := p.cur
	if firstTok.Kind != tokInt {
		return nil, 0, &ParseError{
			Loc: firstTok.Loc,
			Msg: fmt.Sprintf("expected integer parameter, got %q", firstTok.Str),
		}
	}
	p.advance()
	params := []int{int(firstTok.Ival)}

	// Optional second param
	if p.cur.Kind == int(',') {
		p.advance() // consume ','
		secondTok := p.cur
		if secondTok.Kind != tokInt {
			return nil, 0, &ParseError{
				Loc: secondTok.Loc,
				Msg: fmt.Sprintf("expected integer parameter, got %q", secondTok.Str),
			}
		}
		p.advance()
		params = append(params, int(secondTok.Ival))
	}

	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, 0, err
	}
	return params, closeTok.Loc.End, nil
}

// expectGT consumes a '>' token in a type-argument context, handling the case
// where the lexer emits tokShiftRight (>>) for the terminal ">>" in nested
// types like ARRAY<ARRAY<INT>>. When tokShiftRight is seen, we "split" it:
// return a synthetic first '>' and leave the parser positioned at a synthetic
// second '>' so the outer type parser can also consume its closing '>'.
func (p *Parser) expectGT() (Token, error) {
	if p.cur.Kind == int('>') {
		return p.advance(), nil
	}
	if p.cur.Kind == tokShiftRight {
		tok := p.cur // the ">>" token
		firstGT := Token{
			Kind: int('>'),
			Str:  ">",
			Loc:  ast.Loc{Start: tok.Loc.Start, End: tok.Loc.Start + 1},
		}
		secondGT := Token{
			Kind: int('>'),
			Str:  ">",
			Loc:  ast.Loc{Start: tok.Loc.Start + 1, End: tok.Loc.End},
		}
		// advance() puts the token after ">>" into p.cur and clears hasNext.
		p.advance()
		// Now p.cur = token-after->>, p.hasNext = false.
		// We want cur = secondGT, and next = token-after->>.
		realNext := p.cur
		p.cur = secondGT
		p.nextBuf = realNext
		p.hasNext = true
		return firstGT, nil
	}
	return Token{}, &ParseError{
		Loc: p.cur.Loc,
		Msg: fmt.Sprintf("expected '>', got %q", p.cur.Str),
	}
}

// parseArrayType parses: ARRAY '<' elementType '>'
func (p *Parser) parseArrayType() (*ast.TypeName, error) {
	startTok := p.advance() // consume ARRAY
	typName := &ast.TypeName{
		Name: "ARRAY",
		Loc:  startTok.Loc,
	}

	// Expect '<'
	if _, err := p.expect(int('<')); err != nil {
		return nil, &ParseError{
			Loc: p.cur.Loc,
			Msg: "expected '<' after ARRAY",
		}
	}

	// Parse element type recursively
	elemType, err := p.parseDataType()
	if err != nil {
		return nil, err
	}
	typName.ElementType = elemType

	// Expect '>' — uses expectGT to handle ">>" (tokShiftRight) in nested types
	// like ARRAY<ARRAY<INT>> where the lexer produces a single tokShiftRight.
	closeTok, err := p.expectGT()
	if err != nil {
		return nil, err
	}
	typName.Loc.End = closeTok.Loc.End

	return typName, nil
}

// parseMapType parses: MAP '<' keyType ',' valueType '>'
func (p *Parser) parseMapType() (*ast.TypeName, error) {
	startTok := p.advance() // consume MAP
	typName := &ast.TypeName{
		Name: "MAP",
		Loc:  startTok.Loc,
	}

	// Expect '<'
	if _, err := p.expect(int('<')); err != nil {
		return nil, &ParseError{
			Loc: p.cur.Loc,
			Msg: "expected '<' after MAP",
		}
	}

	// Parse key type
	keyType, err := p.parseDataType()
	if err != nil {
		return nil, err
	}
	typName.ElementType = keyType

	// Expect ','
	if _, err := p.expect(int(',')); err != nil {
		return nil, &ParseError{
			Loc: p.cur.Loc,
			Msg: "expected ',' between MAP key and value types",
		}
	}

	// Parse value type
	valType, err := p.parseDataType()
	if err != nil {
		return nil, err
	}
	typName.ValueType = valType

	// Expect '>' — handle ">>" from nested types.
	closeTok, err := p.expectGT()
	if err != nil {
		return nil, err
	}
	typName.Loc.End = closeTok.Loc.End

	return typName, nil
}

// parseStructType parses: STRUCT '<' name1 type1 [',' name2 type2 ...] '>'
func (p *Parser) parseStructType() (*ast.TypeName, error) {
	startTok := p.advance() // consume STRUCT
	typName := &ast.TypeName{
		Name: "STRUCT",
		Loc:  startTok.Loc,
	}

	// Expect '<'
	if _, err := p.expect(int('<')); err != nil {
		return nil, &ParseError{
			Loc: p.cur.Loc,
			Msg: "expected '<' after STRUCT",
		}
	}

	// Parse one or more field definitions: name type [, name type ...]
	for {
		field, err := p.parseStructField()
		if err != nil {
			return nil, err
		}
		typName.Fields = append(typName.Fields, field)

		// Stop at '>' (or '>>' which will be split by expectGT)
		if p.cur.Kind == int('>') || p.cur.Kind == tokShiftRight {
			break
		}
		if p.cur.Kind != int(',') {
			return nil, &ParseError{
				Loc: p.cur.Loc,
				Msg: "expected ',' or '>' in STRUCT type field list",
			}
		}
		p.advance() // consume ','

		// Allow trailing comma before '>'
		if p.cur.Kind == int('>') || p.cur.Kind == tokShiftRight {
			break
		}
	}

	// Expect '>' — handle ">>" from nested types.
	closeTok, err := p.expectGT()
	if err != nil {
		return nil, err
	}
	typName.Loc.End = closeTok.Loc.End

	return typName, nil
}

// parseStructField parses a single STRUCT field definition: name dataType
//
// The field name may be an identifier or a non-reserved keyword used as a
// field name.
func (p *Parser) parseStructField() (*ast.StructField, error) {
	// Field name: accept identifiers AND keywords (field names can shadow keywords)
	nameTok := p.cur
	var fieldName string
	switch {
	case nameTok.Kind == tokIdent || nameTok.Kind == tokQuotedIdent:
		fieldName = nameTok.Str
		p.advance()
	case nameTok.Kind >= 700:
		// Any keyword token is valid as a struct field name
		fieldName = nameTok.Str
		p.advance()
	default:
		return nil, &ParseError{
			Loc: nameTok.Loc,
			Msg: fmt.Sprintf("expected field name in STRUCT, got %q", nameTok.Str),
		}
	}

	// Field type
	fieldType, err := p.parseDataType()
	if err != nil {
		return nil, err
	}

	field := &ast.StructField{
		Name: fieldName,
		Type: fieldType,
		Loc:  ast.Loc{Start: nameTok.Loc.Start, End: fieldType.Loc.End},
	}
	return field, nil
}
