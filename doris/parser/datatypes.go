package parser

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/doris/ast"
)

// isPrimitiveTypeKeyword reports whether kind is a keyword of primitiveColType
// (DorisParser.g4:1717), the types a column may name without angle brackets.
// VARIANT is variantTypeDefinitions (:1759), read on its own.
func isPrimitiveTypeKeyword(kind TokenKind) bool {
	switch kind {
	case kwTINYINT, kwSMALLINT, kwINT, kwINTEGER, kwBIGINT, kwLARGEINT,
		kwBOOLEAN, kwFLOAT, kwDOUBLE,
		kwDATE, kwDATETIME, kwTIME, kwDATEV2, kwDATETIMEV2, kwDATEV1, kwDATETIMEV1,
		kwBITMAP, kwQUANTILE_STATE, kwHLL, kwAGG_STATE,
		kwSTRING, kwJSON, kwJSONB, kwTEXT, kwVARCHAR, kwCHAR,
		kwDECIMAL, kwDECIMALV2, kwDECIMALV3, kwIPV4, kwIPV6, kwALL:
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

// parseDataType reads a data type as DorisParser.g4's dataType (:1705) does:
//
//	dataType
//	    : ARRAY LT dataType GT
//	    | MAP LT dataType COMMA dataType GT
//	    | STRUCT LT complexColTypeList GT
//	    | variantTypeDefinitions
//	    | AGG_STATE LT functionNameIdentifier
//	        LEFT_PAREN dataTypeWithNullable (COMMA dataTypeWithNullable)* RIGHT_PAREN GT
//	    | primitiveColType (LEFT_PAREN (INTEGER_VALUE | ASTERISK)
//	        (COMMA INTEGER_VALUE)* RIGHT_PAREN)?
//	    ;
//
// It is the type of the statements only the ANTLR parser reads (CREATE TABLE,
// CAST). A statement Doris also hands to its CUP parser takes
// parseLegacyDataType instead.
func (p *Parser) parseDataType() (*ast.TypeName, error) {
	return p.readDataType(false)
}

// parseLegacyDataType reads a data type in a statement that DorisParser.g4
// derives from an unsupportedStatement, which Doris parses with both of its
// parsers: the part of dataType (:1705) that the CUP parser's type
// (DorisLegacyParser.g4:2528) also takes. That narrows each primitive to
// the parameters its CUP alternative lists, drops a bare AGG_STATE and the
// MATCH_NAME of a VARIANT sub-column, and takes VARIANT properties as
// string pairs only.
func (p *Parser) parseLegacyDataType() (*ast.TypeName, error) {
	return p.readDataType(true)
}

// readDataType reads dataType; legacy narrows it to what the CUP parser's
// type also takes.
func (p *Parser) readDataType(legacy bool) (*ast.TypeName, error) {
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()

	switch p.cur.Kind {
	case kwARRAY:
		return p.parseArrayType(legacy)
	case kwMAP:
		return p.parseMapType(legacy)
	case kwSTRUCT:
		return p.parseStructType(legacy)
	case kwVARIANT:
		return p.parseVariantType(legacy)
	case kwAGG_STATE:
		if p.peekNext().Kind == int('<') {
			return p.parseAggStateType(legacy)
		}
		if legacy {
			// The CUP parser has no bare AGG_STATE.
			p.advance()
			return nil, p.syntaxErrorAtCur()
		}
	}
	if !isPrimitiveTypeKeyword(p.cur.Kind) {
		return nil, &ParseError{
			Position: p.cur.Loc.Start, End: p.cur.Loc.End,
			Message: fmt.Sprintf("expected a data type, got %q", p.cur.Str),
		}
	}
	return p.parsePrimitiveType(legacy)
}

// legacyTypeParams says which parameter lists the CUP parser's type takes
// after each primitive: none (absent from the map), "(" N ")" (one), or
// also "(" N "," M ")" (two). VARCHAR also takes "(*)".
var legacyTypeParams = map[TokenKind]int{
	kwTINYINT: 1, kwSMALLINT: 1, kwINT: 1, kwINTEGER: 1, kwBIGINT: 1, kwLARGEINT: 1,
	kwDATETIME: 1, kwTIME: 1, kwDATETIMEV2: 1, kwCHAR: 1, kwVARCHAR: 1,
	kwDECIMAL: 2, kwDECIMALV2: 2, kwDECIMALV3: 2,
}

// parsePrimitiveType reads primitiveColType and its optional parameter list:
//
//	'(' (INTEGER_VALUE | '*') (',' INTEGER_VALUE)* ')'
//
// legacy holds the list to what the CUP parser's type takes for that
// primitive (legacyTypeParams).
func (p *Parser) parsePrimitiveType(legacy bool) (*ast.TypeName, error) {
	tok := p.advance()
	typName := &ast.TypeName{
		Name: typeKeywordName(tok),
		Loc:  tok.Loc,
	}
	if p.cur.Kind != int('(') {
		return typName, nil
	}
	maxParams := -1 // no limit
	if legacy {
		maxParams = legacyTypeParams[tok.Kind]
		if maxParams == 0 {
			return nil, p.syntaxErrorAtCur()
		}
	}
	p.advance() // consume '('
	switch {
	case p.cur.Kind == int('*') && (!legacy || tok.Kind == kwVARCHAR):
		typName.Star = true
		p.advance()
		if legacy {
			maxParams = 0 // the CUP parser takes VARCHAR(*) alone
		}
	case p.cur.Kind == tokInt:
		typName.Params = append(typName.Params, int(p.advance().Ival))
	default:
		return nil, &ParseError{
			Position: p.cur.Loc.Start, End: p.cur.Loc.End,
			Message: fmt.Sprintf("expected integer parameter, got %q", p.cur.Str),
		}
	}
	for p.cur.Kind == int(',') && (maxParams < 0 || len(typName.Params) < maxParams) {
		p.advance() // consume ','
		if p.cur.Kind != tokInt {
			return nil, &ParseError{
				Position: p.cur.Loc.Start, End: p.cur.Loc.End,
				Message: fmt.Sprintf("expected integer parameter, got %q", p.cur.Str),
			}
		}
		typName.Params = append(typName.Params, int(p.advance().Ival))
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	typName.Loc.End = closeTok.Loc.End
	return typName, nil
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
		Position: p.cur.Loc.Start, End: p.cur.Loc.End,
		Message: fmt.Sprintf("expected '>', got %q", p.cur.Str),
	}
}

// expectLT consumes the '<' that opens a complex type's arguments.
func (p *Parser) expectLT(typeName string) error {
	if p.cur.Kind != int('<') {
		return &ParseError{
			Position: p.cur.Loc.Start, End: p.cur.Loc.End,
			Message: "expected '<' after " + typeName,
		}
	}
	p.advance()
	return nil
}

// parseArrayType parses: ARRAY '<' dataType '>'
func (p *Parser) parseArrayType(legacy bool) (*ast.TypeName, error) {
	startTok := p.advance() // consume ARRAY
	typName := &ast.TypeName{Name: "ARRAY", Loc: startTok.Loc}
	if err := p.expectLT("ARRAY"); err != nil {
		return nil, err
	}
	elemType, err := p.readDataType(legacy)
	if err != nil {
		return nil, err
	}
	typName.ElementType = elemType
	// expectGT splits the ">>" that closes ARRAY<ARRAY<INT>>.
	closeTok, err := p.expectGT()
	if err != nil {
		return nil, err
	}
	typName.Loc.End = closeTok.Loc.End
	return typName, nil
}

// parseMapType parses: MAP '<' dataType ',' dataType '>'
func (p *Parser) parseMapType(legacy bool) (*ast.TypeName, error) {
	startTok := p.advance() // consume MAP
	typName := &ast.TypeName{Name: "MAP", Loc: startTok.Loc}
	if err := p.expectLT("MAP"); err != nil {
		return nil, err
	}
	keyType, err := p.readDataType(legacy)
	if err != nil {
		return nil, err
	}
	typName.ElementType = keyType
	if p.cur.Kind != int(',') {
		return nil, &ParseError{
			Position: p.cur.Loc.Start, End: p.cur.Loc.End,
			Message: "expected ',' between MAP key and value types",
		}
	}
	p.advance() // consume ','
	valType, err := p.readDataType(legacy)
	if err != nil {
		return nil, err
	}
	typName.ValueType = valType
	closeTok, err := p.expectGT()
	if err != nil {
		return nil, err
	}
	typName.Loc.End = closeTok.Loc.End
	return typName, nil
}

// parseStructType parses STRUCT '<' complexColTypeList '>' (DorisParser.g4:1751):
//
//	complexColType : identifier COLON dataType commentSpec?
//
// The CUP parser's struct_field (DorisLegacyParser.g4:2708) is the same.
func (p *Parser) parseStructType(legacy bool) (*ast.TypeName, error) {
	startTok := p.advance() // consume STRUCT
	typName := &ast.TypeName{Name: "STRUCT", Loc: startTok.Loc}
	if err := p.expectLT("STRUCT"); err != nil {
		return nil, err
	}
	for {
		nameTok := p.cur
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		field, err := p.parseTypeField(name, nameTok.Loc, legacy)
		if err != nil {
			return nil, err
		}
		typName.Fields = append(typName.Fields, field)
		if p.cur.Kind != int(',') {
			break
		}
		p.advance() // consume ','
	}
	closeTok, err := p.expectGT()
	if err != nil {
		return nil, err
	}
	typName.Loc.End = closeTok.Loc.End
	return typName, nil
}

// parseTypeField reads the rest of a STRUCT field or a VARIANT sub-column
// after its name: COLON dataType (COMMENT STRING_LITERAL)?.
func (p *Parser) parseTypeField(name string, nameLoc ast.Loc, legacy bool) (*ast.StructField, error) {
	if _, err := p.expect(int(':')); err != nil {
		return nil, err
	}
	fieldType, err := p.readDataType(legacy)
	if err != nil {
		return nil, err
	}
	field := &ast.StructField{
		Name: name,
		Type: fieldType,
		Loc:  ast.Loc{Start: nameLoc.Start, End: fieldType.Loc.End},
	}
	if p.cur.Kind == kwCOMMENT {
		p.advance() // consume COMMENT
		if p.cur.Kind != tokString {
			return nil, p.syntaxErrorAtCur()
		}
		tok := p.advance()
		field.Comment = tok.Str
		field.Loc.End = tok.Loc.End
	}
	return field, nil
}

// parseVariantType parses variantTypeDefinitions (DorisParser.g4:1759):
//
//	VARIANT LT variantSubColTypeList COMMA propertyClause GT
//	| VARIANT LT variantSubColTypeList GT
//	| VARIANT LT propertyClause GT
//	| VARIANT
//	variantSubColType : variantSubColMatchType? STRING_LITERAL COLON dataType commentSpec?
//
// legacy reads the CUP parser's forms (DorisLegacyParser.g4:2564): no
// MATCH_NAME, and properties as 'k' = 'v' pairs.
func (p *Parser) parseVariantType(legacy bool) (*ast.TypeName, error) {
	startTok := p.advance() // consume VARIANT
	typName := &ast.TypeName{Name: "VARIANT", Loc: startTok.Loc}
	if p.cur.Kind != int('<') {
		return typName, nil
	}
	p.advance() // consume '<'
	// variantSubColTypeList COMMA propertyClause: properties after the
	// sub-columns need the comma.
	propsAllowed := p.cur.Kind == kwPROPERTIES
	if !propsAllowed {
		for {
			start := p.cur.Loc
			var match string
			if !legacy && (p.cur.Kind == kwMATCH_NAME || p.cur.Kind == kwMATCH_NAME_GLOB) {
				match = strings.ToUpper(p.advance().Str)
			}
			if p.cur.Kind != tokString {
				return nil, p.syntaxErrorAtCur()
			}
			path := p.advance()
			field, err := p.parseTypeField(path.Str, start, legacy)
			if err != nil {
				return nil, err
			}
			field.MatchType = match
			typName.Fields = append(typName.Fields, field)
			if p.cur.Kind != int(',') {
				break
			}
			p.advance() // consume ','
			if p.cur.Kind == kwPROPERTIES {
				propsAllowed = true
				break
			}
		}
	}
	if propsAllowed {
		parseItem := p.parsePropertyItem
		if legacy {
			parseItem = p.parseStringProperty
		}
		p.advance() // consume PROPERTIES
		if _, err := p.expect(int('(')); err != nil {
			return nil, err
		}
		props, err := p.parsePropertyItems(parseItem)
		if err != nil {
			return nil, err
		}
		typName.Properties = props
	}
	closeTok, err := p.expectGT()
	if err != nil {
		return nil, err
	}
	typName.Loc.End = closeTok.Loc.End
	return typName, nil
}

// parseAggStateType parses
//
//	AGG_STATE LT functionNameIdentifier
//	    LEFT_PAREN dataTypeWithNullable (COMMA dataTypeWithNullable)* RIGHT_PAREN GT
//	dataTypeWithNullable : dataType ((NOT)? NULL)?
//
// (DorisParser.g4:1710, :1701). The CUP parser's form
// (DorisLegacyParser.g4:2585) names the function by an IDENTIFIER, so legacy
// takes no reserved word there.
func (p *Parser) parseAggStateType(legacy bool) (*ast.TypeName, error) {
	startTok := p.advance() // consume AGG_STATE
	typName := &ast.TypeName{Name: "AGG_STATE", Loc: startTok.Loc}
	p.advance() // consume '<'
	switch {
	case isIdentifierToken(p.cur.Kind):
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		typName.Func = name
	case !legacy && isFunctionNameKeyword(p.cur.Kind):
		typName.Func = p.advance().Str
	default:
		return nil, p.syntaxErrorAtCur()
	}
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	for {
		arg, err := p.readDataType(legacy)
		if err != nil {
			return nil, err
		}
		switch {
		case p.cur.Kind == kwNULL:
			arg.Nullable = "NULL"
			arg.Loc.End = p.advance().Loc.End
		case p.cur.Kind == kwNOT && p.peekNext().Kind == kwNULL:
			p.advance() // consume NOT
			arg.Nullable = "NOT NULL"
			arg.Loc.End = p.advance().Loc.End
		}
		typName.Args = append(typName.Args, arg)
		if p.cur.Kind != int(',') {
			break
		}
		p.advance() // consume ','
	}
	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}
	closeTok, err := p.expectGT()
	if err != nil {
		return nil, err
	}
	typName.Loc.End = closeTok.Loc.End
	return typName, nil
}
