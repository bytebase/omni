package parser

import (
	"strconv"
	"strings"
)

// normalizeIdentName unquotes a backtick-quoted identifier (unescaping “) and
// lowercases it for case-insensitive duplicate checks (e.g. CTE names).
func normalizeIdentName(s string) string {
	if len(s) >= 2 && s[0] == '`' && s[len(s)-1] == '`' {
		s = strings.ReplaceAll(s[1:len(s)-1], "``", "`")
	}
	return strings.ToLower(s)
}

// parseIdentifier parses Spark's identifier rule:
//
//	identifier : strictIdentifier | strictNonReserved ;
func (p *Parser) parseIdentifier() (string, error) {
	if !p.cur().Kind.isIdentifier() {
		return "", p.errorf("expected identifier, got %q", p.cur().Text)
	}
	tok := p.advance()
	return tok.Text, nil
}

// parseStrictIdentifier parses Spark's strictIdentifier rule.
func (p *Parser) parseStrictIdentifier() (string, error) {
	if !p.cur().Kind.isStrictIdentifier() {
		return "", p.errorf("expected identifier, got %q", p.cur().Text)
	}
	tok := p.advance()
	return tok.Text, nil
}

// parseErrorCapturingIdentifier parses Spark's errorCapturingIdentifier rule
// (SqlBaseParser.g4:1479): an identifier, rejecting the error-capturing
// unquoted hyphenated alternative (`a-b`) with a positioned ParseError just as
// Spark's AstBuilder does.
func (p *Parser) parseErrorCapturingIdentifier() (string, error) {
	first, err := p.parseIdentifier()
	if err != nil {
		return "", err
	}
	if p.cur().Kind == tokMinus && p.peek(1).Kind.isIdentifier() {
		return "", p.errorf("unquoted identifier cannot contain '-'; quote it with backticks")
	}
	return first, nil
}

// parseQualifiedName parses: identifier (DOT identifier)*
func (p *Parser) parseQualifiedName() (string, error) {
	first, err := p.parseIdentifier()
	if err != nil {
		return "", err
	}
	if p.cur().Kind != tokDot {
		return first, nil
	}
	var b strings.Builder
	b.WriteString(first)
	for p.match(tokDot) {
		part, err := p.parseIdentifier()
		if err != nil {
			return "", err
		}
		b.WriteByte('.')
		b.WriteString(part)
	}
	return b.String(), nil
}

// parseMultipartIdentifier parses: errorCapturingIdentifier (DOT errorCapturingIdentifier)*
func (p *Parser) parseMultipartIdentifier() (string, error) {
	first, err := p.parseErrorCapturingIdentifier()
	if err != nil {
		return "", err
	}
	if p.cur().Kind != tokDot {
		return first, nil
	}
	var b strings.Builder
	b.WriteString(first)
	for p.match(tokDot) {
		part, err := p.parseErrorCapturingIdentifier()
		if err != nil {
			return "", err
		}
		b.WriteByte('.')
		b.WriteString(part)
	}
	return b.String(), nil
}

// parseMultipartIdentifierList parses: multipartIdentifier (COMMA multipartIdentifier)*
func (p *Parser) parseMultipartIdentifierList() error {
	if _, err := p.parseMultipartIdentifier(); err != nil {
		return err
	}
	for p.match(tokComma) {
		if _, err := p.parseMultipartIdentifier(); err != nil {
			return err
		}
	}
	return nil
}

// parseTableIdentifier parses: (errorCapturingIdentifier DOT)? errorCapturingIdentifier
func (p *Parser) parseTableIdentifier() (string, error) {
	first, err := p.parseErrorCapturingIdentifier()
	if err != nil {
		return "", err
	}
	if p.match(tokDot) {
		second, err := p.parseErrorCapturingIdentifier()
		if err != nil {
			return "", err
		}
		return first + "." + second, nil
	}
	return first, nil
}

// parseIdentifierReference parses:
//
//	identifierReference
//	    : IDENTIFIER_KW LEFT_PAREN expression RIGHT_PAREN
//	    | multipartIdentifier
//	    ;
//
// When `IDENTIFIER` is followed by `(` that does not form a single parenthesized
// expression (e.g. `CREATE TABLE identifier (a INT)`, `CREATE FUNCTION identifier(x INT) ...`,
// `INSERT INTO identifier (a, b) ...`), it falls back to multipartIdentifier so
// objects named `identifier` do not require quoting.
func (p *Parser) parseIdentifierReference() (string, error) {
	if p.cur().Kind == kwIDENTIFIER && p.peek(1).Kind == tokLParen {
		cp := p.save()
		p.advance()
		p.advance()
		if err := p.parseExpression(); err == nil && p.match(tokRParen) {
			return "IDENTIFIER(...)", nil
		}
		p.restore(cp)
	}
	return p.parseMultipartIdentifier()
}

// parseCatalogIdentifierReference parses:
//
//	catalogIdentifierReference
//	    : IDENTIFIER_KW LEFT_PAREN expression RIGHT_PAREN
//	    | errorCapturingIdentifier
//	    | stringLit
//	    ;
func (p *Parser) parseCatalogIdentifierReference() (string, error) {
	if p.cur().Kind == kwIDENTIFIER && p.peek(1).Kind == tokLParen {
		p.advance()
		p.advance()
		if err := p.parseExpression(); err != nil {
			return "", err
		}
		if _, err := p.expect(tokRParen); err != nil {
			return "", err
		}
		return "IDENTIFIER(...)", nil
	}
	if p.cur().Kind == tokString {
		tok := p.advance()
		return tok.Text, nil
	}
	return p.parseErrorCapturingIdentifier()
}

// parseIdentifierSeq parses: errorCapturingIdentifier (COMMA errorCapturingIdentifier)*
func (p *Parser) parseIdentifierSeq() error {
	if _, err := p.parseErrorCapturingIdentifier(); err != nil {
		return err
	}
	for p.match(tokComma) {
		if _, err := p.parseErrorCapturingIdentifier(); err != nil {
			return err
		}
	}
	return nil
}

// parseIdentifierList parses: LEFT_PAREN identifierSeq RIGHT_PAREN
func (p *Parser) parseIdentifierList() error {
	if _, err := p.expect(tokLParen); err != nil {
		return err
	}
	if err := p.parseIdentifierSeq(); err != nil {
		return err
	}
	_, err := p.expect(tokRParen)
	return err
}

func (p *Parser) parseOrderedIdentifierList() error {
	if _, err := p.expect(tokLParen); err != nil {
		return err
	}
	for {
		if _, err := p.parseErrorCapturingIdentifier(); err != nil {
			return err
		}
		p.match(kwASC, kwDESC)
		if !p.match(tokComma) {
			break
		}
	}
	_, err := p.expect(tokRParen)
	return err
}

func (p *Parser) parseIdentifierCommentList() error {
	if _, err := p.expect(tokLParen); err != nil {
		return err
	}
	for {
		if _, err := p.parseIdentifier(); err != nil {
			return err
		}
		if p.cur().Kind == kwCOMMENT {
			if err := p.parseCommentSpec(); err != nil {
				return err
			}
		}
		if !p.match(tokComma) {
			break
		}
	}
	_, err := p.expect(tokRParen)
	return err
}

func (p *Parser) parseStringLit() error {
	if p.cur().Kind == tokString || p.cur().Kind == tokDoubleQuoted {
		p.advance()
		return nil
	}
	return p.errorf("expected string literal, got %q", p.cur().Text)
}

func (p *Parser) isStringLit() bool {
	return p.cur().Kind == tokString || p.cur().Kind == tokDoubleQuoted
}

func (p *Parser) matchNot() bool {
	return p.match(kwNOT, tokBang)
}

func (p *Parser) isNot() bool {
	return p.cur().Kind == kwNOT || p.cur().Kind == tokBang
}

func (p *Parser) matchNamespace() bool {
	return p.match(kwNAMESPACE, kwDATABASE, kwSCHEMA)
}

func (p *Parser) isNamespace() bool {
	k := p.cur().Kind
	return k == kwNAMESPACE || k == kwDATABASE || k == kwSCHEMA
}

func (p *Parser) matchNamespaces() bool {
	return p.match(kwNAMESPACES, kwDATABASES, kwSCHEMAS)
}

func (p *Parser) matchVariable() bool {
	return p.match(kwVARIABLE, kwVAR)
}

func (p *Parser) parseCommentSpec() error {
	if _, err := p.expect(kwCOMMENT); err != nil {
		return err
	}
	return p.parseStringLit()
}

func (p *Parser) parseLocationSpec() error {
	if _, err := p.expect(kwLOCATION); err != nil {
		return err
	}
	return p.parseStringLit()
}

func (p *Parser) parseCollationSpec() error {
	if _, err := p.expect(kwDEFAULT); err != nil {
		return err
	}
	if _, err := p.expect(kwCOLLATION); err != nil {
		return err
	}
	name, err := p.parseIdentifier()
	if err != nil {
		return err
	}
	if !isValidSparkCollation(name) {
		return p.errorf("invalid collation name %q", name)
	}
	return nil
}

func (p *Parser) parseCollateClause(strictName bool) error {
	if _, err := p.expect(kwCOLLATE); err != nil {
		return err
	}
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return err
	}
	if isTrimSparkCollation(name) {
		return p.errorf("TRIM specifier in collation %q is not supported", name)
	}
	if strictName && !isValidSparkCollation(name) {
		return p.errorf("invalid collation name %q", name)
	}
	return nil
}

func (p *Parser) parseSchemaBinding() error {
	if _, err := p.expect(kwWITH); err != nil {
		return err
	}
	if _, err := p.expect(kwSCHEMA); err != nil {
		return err
	}
	if p.match(kwBINDING, kwCOMPENSATION, kwEVOLUTION) {
		return nil
	}
	if p.match(kwTYPE) {
		_, err := p.expect(kwEVOLUTION)
		return err
	}
	return p.errorf("expected BINDING, COMPENSATION, EVOLUTION, or TYPE EVOLUTION")
}

// parseDataType parses Spark/Databricks data types:
//
//	dataType
//	    : ARRAY < dataType >
//	    | MAP < dataType , dataType >
//	    | STRUCT ( < complexColTypeList? > | <> )
//	    | INTERVAL (YEAR | MONTH) (TO MONTH)?
//	    | INTERVAL (DAY | HOUR | MINUTE | SECOND) (TO (HOUR | MINUTE | SECOND))?
//	    | type ( ( INTEGER_VALUE (, INTEGER_VALUE)* ) )?
//	    ;
func (p *Parser) parseDataType() error {
	_, err := p.parseDataTypeWithKind()
	return err
}

func (p *Parser) parseDataTypeWithKind() (TokenKind, error) {
	if err := p.enterNesting(); err != nil {
		return tokEOF, err
	}
	defer p.leaveNesting()

	switch p.cur().Kind {
	case kwARRAY:
		if p.peek(1).Kind == tokLT {
			p.advance()
			p.advance()
			if err := p.parseDataType(); err != nil {
				return tokEOF, err
			}
			return kwARRAY, p.expectGT()
		}
		if !p.relaxFunctionDataTypes {
			return tokEOF, p.errorf("incomplete ARRAY type definition; expected ARRAY<elementType>")
		}
	case kwMAP:
		if p.peek(1).Kind == tokLT {
			p.advance()
			p.advance()
			if err := p.parseDataType(); err != nil {
				return tokEOF, err
			}
			if _, err := p.expect(tokComma); err != nil {
				return tokEOF, err
			}
			if err := p.parseDataType(); err != nil {
				return tokEOF, err
			}
			return kwMAP, p.expectGT()
		}
		if !p.relaxFunctionDataTypes {
			return tokEOF, p.errorf("incomplete MAP type definition; expected MAP<keyType, valueType>")
		}
	case kwSTRUCT:
		if p.peek(1).Kind == tokNeq {
			p.advance()
			p.advance()
			return kwSTRUCT, nil
		}
		if p.peek(1).Kind == tokLT {
			p.advance()
			p.advance()
			if p.cur().Kind != tokGT && p.cur().Kind != tokShiftRight && p.cur().Kind != tokShiftRightU {
				if err := p.parseComplexColTypeList(); err != nil {
					return tokEOF, err
				}
			}
			return kwSTRUCT, p.expectGT()
		}
		if !p.relaxFunctionDataTypes {
			return tokEOF, p.errorf("incomplete STRUCT type definition; expected STRUCT<...>")
		}
	case kwINTERVAL:
		p.advance()
		if p.match(kwYEAR, kwMONTH) {
			if p.match(kwTO) {
				if _, err := p.expect(kwMONTH); err != nil {
					return tokEOF, err
				}
			}
			return kwINTERVAL, nil
		}
		if p.match(kwDAY, kwHOUR, kwMINUTE, kwSECOND) {
			if p.match(kwTO) {
				if !p.match(kwHOUR, kwMINUTE, kwSECOND) {
					return tokEOF, p.errorf("expected HOUR, MINUTE, or SECOND after TO")
				}
			}
			return kwINTERVAL, nil
		}
		// Plain INTERVAL primitive type (no parameters).
		return kwINTERVAL, p.parseOptionalTypeParams(kwINTERVAL, 0, false, "INTERVAL")
	}

	// Primitive / named type.
	kind, maxParams, requireParams, typeName, err := p.parsePrimitiveTypeName()
	if err != nil {
		return tokEOF, err
	}
	return kind, p.parseOptionalTypeParams(kind, maxParams, requireParams, typeName)
}

func (p *Parser) parsePrimitiveTypeName() (TokenKind, int, bool, string, error) {
	tok := p.cur()
	switch tok.Kind {
	case kwSTRING:
		p.advance()
		if p.cur().Kind == kwCOLLATE {
			if err := p.parseCollateClause(true); err != nil {
				return tokEOF, 0, false, "", err
			}
		}
		return kwSTRING, 0, false, tok.Text, nil
	case kwCHARACTER, kwCHAR, kwVARCHAR:
		p.advance()
		return tok.Kind, 1, !p.relaxFunctionDataTypes, tok.Text, nil
	case kwDECIMAL, kwDEC, kwNUMERIC:
		p.advance()
		return tok.Kind, 2, false, tok.Text, nil
	case kwBOOLEAN, kwTINYINT, kwBYTE, kwSMALLINT, kwSHORT,
		kwINT, kwINTEGER, kwBIGINT, kwLONG, kwFLOAT, kwREAL, kwDOUBLE,
		kwDATE, kwTIMESTAMP, kwTIMESTAMP_NTZ, kwTIMESTAMP_LTZ,
		kwBINARY, kwVOID, kwINTERVAL, kwVARIANT:
		p.advance()
		return tok.Kind, 0, false, tok.Text, nil
	default:
		if p.relaxFunctionDataTypes && tok.Kind.isIdentifier() {
			p.advance()
			return tok.Kind, 2, false, tok.Text, nil
		}
		return tokEOF, 0, false, "", p.errorf("unsupported data type %q", tok.Text)
	}
}

func (p *Parser) parseOptionalTypeParams(kind TokenKind, maxParams int, requireParams bool, typeName string) error {
	if p.cur().Kind != tokLParen {
		if requireParams {
			return p.errorf("data type %q requires a length parameter", typeName)
		}
		return nil
	}
	cp := p.save()
	p.advance()
	if p.cur().Kind != tokInt {
		p.restore(cp)
		if requireParams {
			return p.errorf("data type %q requires a length parameter", typeName)
		}
		return nil
	}
	if !p.relaxFunctionDataTypes && maxParams == 0 {
		return p.errorf("data type %q does not accept parameters", typeName)
	}
	firstTok := p.advance()
	params := []string{firstTok.Text}
	for p.match(tokComma) {
		tok, err := p.expect(tokInt)
		if err != nil {
			p.restore(cp)
			return err
		}
		params = append(params, tok.Text)
		if !p.relaxFunctionDataTypes && len(params) > maxParams {
			return p.errorf("data type %q accepts at most %d parameter(s)", typeName, maxParams)
		}
	}
	if _, err := p.expect(tokRParen); err != nil {
		return err
	}
	if !p.relaxFunctionDataTypes && (kind == kwDECIMAL || kind == kwDEC || kind == kwNUMERIC) {
		prec, err := strconv.Atoi(params[0])
		if err != nil || prec > 38 {
			return p.errorf("DECIMAL precision %s exceeds maximum precision 38", params[0])
		}
		if len(params) == 2 {
			scale, err := strconv.Atoi(params[1])
			if err != nil || scale > prec {
				return p.errorf("DECIMAL scale %s cannot exceed precision %d", params[1], prec)
			}
		}
	}
	return nil
}

func (p *Parser) parseComplexColTypeList() error {
	if err := p.parseComplexColType(); err != nil {
		return err
	}
	for p.match(tokComma) {
		if err := p.parseComplexColType(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parseComplexColType() error {
	if _, err := p.parseErrorCapturingIdentifier(); err != nil {
		return err
	}
	p.match(tokColon)
	if err := p.parseDataType(); err != nil {
		return err
	}
	if p.isNot() && p.peek(1).Kind == kwNULL {
		p.advance()
		p.advance()
	}
	if p.cur().Kind == kwCOMMENT {
		if err := p.parseCommentSpec(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parseColTypeList() error {
	if err := p.parseColType(); err != nil {
		return err
	}
	for p.match(tokComma) {
		if err := p.parseColType(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parseColType() error {
	if _, err := p.parseErrorCapturingIdentifier(); err != nil {
		return err
	}
	if err := p.parseDataType(); err != nil {
		return err
	}
	if p.isNot() && p.peek(1).Kind == kwNULL {
		p.advance()
		p.advance()
	}
	if p.cur().Kind == kwCOMMENT {
		if err := p.parseCommentSpec(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parseColDefinitionList() error {
	if err := p.parseColDefinition(); err != nil {
		return err
	}
	for p.match(tokComma) {
		if err := p.parseColDefinition(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parseColDefinition() error {
	if _, err := p.parseErrorCapturingIdentifier(); err != nil {
		return err
	}
	typeKind, err := p.parseDataTypeWithKind()
	if err != nil {
		return err
	}
	var seenNotNull, seenDefault, seenGenerated, seenComment bool
	for {
		switch {
		case p.cur().Kind == kwNOT && p.peek(1).Kind == kwNULL:
			if seenNotNull {
				return p.errorf("duplicate NOT NULL descriptor in column definition")
			}
			seenNotNull = true
			p.advance()
			p.advance()
		case p.cur().Kind == kwDEFAULT:
			if seenDefault {
				return p.errorf("duplicate DEFAULT descriptor in column definition")
			}
			seenDefault = true
			if err := p.parseDefaultExpression(); err != nil {
				return err
			}
		case p.cur().Kind == kwGENERATED:
			if seenGenerated {
				return p.errorf("duplicate GENERATED descriptor in column definition")
			}
			seenGenerated = true
			if err := p.parseGenerationExpression(typeKind); err != nil {
				return err
			}
		case p.cur().Kind == kwCOMMENT:
			if seenComment {
				return p.errorf("duplicate COMMENT descriptor in column definition")
			}
			seenComment = true
			if err := p.parseCommentSpec(); err != nil {
				return err
			}
		default:
			return nil
		}
	}
}

func (p *Parser) parseDefaultExpression() error {
	if _, err := p.expect(kwDEFAULT); err != nil {
		return err
	}
	return p.parseExpression()
}

func isIdentityColumnType(k TokenKind) bool {
	switch k {
	case kwINT, kwINTEGER, kwBIGINT, kwLONG:
		return true
	}
	return false
}

func (p *Parser) parseGenerationExpression(colType TokenKind) error {
	if _, err := p.expect(kwGENERATED); err != nil {
		return err
	}
	if p.match(kwALWAYS) {
		if _, err := p.expect(kwAS); err != nil {
			return err
		}
		if p.match(kwIDENTITY) {
			if !isIdentityColumnType(colType) {
				return p.errorf("IDENTITY column requires INT or BIGINT data type")
			}
			return p.parseOptionalIdentityColSpec()
		}
		if _, err := p.expect(tokLParen); err != nil {
			return err
		}
		if err := p.parseExpression(); err != nil {
			return err
		}
		_, err := p.expect(tokRParen)
		return err
	}
	if _, err := p.expect(kwBY); err != nil {
		return err
	}
	if _, err := p.expect(kwDEFAULT); err != nil {
		return err
	}
	if _, err := p.expect(kwAS); err != nil {
		return err
	}
	if _, err := p.expect(kwIDENTITY); err != nil {
		return err
	}
	if !isIdentityColumnType(colType) {
		return p.errorf("IDENTITY column requires INT or BIGINT data type")
	}
	return p.parseOptionalIdentityColSpec()
}

func (p *Parser) parseOptionalIdentityColSpec() error {
	if !p.match(tokLParen) {
		return nil
	}
	var seenStart, seenIncrement bool
	for p.cur().Kind != tokRParen && p.cur().Kind != tokEOF {
		if p.match(kwSTART) {
			if seenStart {
				return p.errorf("duplicate START WITH option in IDENTITY column")
			}
			seenStart = true
			if _, err := p.expect(kwWITH); err != nil {
				return err
			}
			neg := p.match(tokMinus)
			tok := p.cur()
			if !p.match(tokInt, tokBigInt) {
				return p.errorf("expected integer or bigint for START WITH")
			}
			if !validateNumericLiteral(tok.Kind, tok.Text, neg) {
				return p.errorf("START WITH value %q out of range", tok.Text)
			}
		} else if p.match(kwINCREMENT) {
			if seenIncrement {
				return p.errorf("duplicate INCREMENT BY option in IDENTITY column")
			}
			seenIncrement = true
			if _, err := p.expect(kwBY); err != nil {
				return err
			}
			neg := p.match(tokMinus)
			tok := p.cur()
			if !p.match(tokInt, tokBigInt) {
				return p.errorf("expected integer or bigint for INCREMENT BY")
			}
			if !validateNumericLiteral(tok.Kind, tok.Text, neg) {
				return p.errorf("INCREMENT BY value %q out of range", tok.Text)
			}
			numStr := tok.Text
			if tok.Kind == tokBigInt {
				numStr = numStr[:len(numStr)-1]
			}
			if strings.TrimLeft(numStr, "0") == "" {
				return p.errorf("IDENTITY column INCREMENT BY step cannot be 0")
			}
		} else {
			return p.errorf("expected START WITH or INCREMENT BY, got %q", p.cur().Text)
		}
	}
	_, err := p.expect(tokRParen)
	return err
}

func (p *Parser) parseQualifiedColTypeWithPositionList() error {
	if err := p.parseQualifiedColTypeWithPosition(); err != nil {
		return err
	}
	for p.match(tokComma) {
		if err := p.parseQualifiedColTypeWithPosition(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parseQualifiedColTypeWithPosition() error {
	if _, err := p.parseMultipartIdentifier(); err != nil {
		return err
	}
	if err := p.parseDataType(); err != nil {
		return err
	}
	var seenNotNull, seenDefault, seenComment, seenPosition bool
	for {
		switch {
		case p.isNot() && p.peek(1).Kind == kwNULL:
			if seenNotNull {
				return p.errorf("duplicate NOT NULL descriptor in column definition")
			}
			seenNotNull = true
			p.advance()
			p.advance()
		case p.cur().Kind == kwDEFAULT:
			if seenDefault {
				return p.errorf("duplicate DEFAULT descriptor in column definition")
			}
			seenDefault = true
			if err := p.parseDefaultExpression(); err != nil {
				return err
			}
		case p.cur().Kind == kwCOMMENT:
			if seenComment {
				return p.errorf("duplicate COMMENT descriptor in column definition")
			}
			seenComment = true
			if err := p.parseCommentSpec(); err != nil {
				return err
			}
		case p.cur().Kind == kwFIRST || p.cur().Kind == kwAFTER:
			if seenPosition {
				return p.errorf("duplicate column position descriptor in column definition")
			}
			seenPosition = true
			if err := p.parseColPosition(); err != nil {
				return err
			}
		default:
			return nil
		}
	}
}

func (p *Parser) parseColPosition() error {
	if p.match(kwFIRST) {
		return nil
	}
	if p.match(kwAFTER) {
		_, err := p.parseErrorCapturingIdentifier()
		return err
	}
	return p.errorf("expected FIRST or AFTER")
}

// parsePropertyList parses setter key-value property lists
// (`TBLPROPERTIES`, `DBPROPERTIES`, `SERDEPROPERTIES`, `OPTIONS`, `WITH (...)`),
// requiring a value for every property and rejecting duplicate keys just as
// Spark v4.0.5 visitPropertyList / visitPropertyKeyValues does.
func (p *Parser) parsePropertyList() error {
	_, err := p.parsePropertyListWithKeys()
	return err
}

func (p *Parser) parsePropertyListWithKeys() (map[string]struct{}, error) {
	if _, err := p.expect(tokLParen); err != nil {
		return nil, err
	}
	seen := make(map[string]struct{})
	for {
		key, err := p.parseProperty()
		if err != nil {
			return nil, err
		}
		if _, dup := seen[key]; dup {
			return nil, p.errorf("duplicate property key %q", key)
		}
		seen[key] = struct{}{}
		if !p.match(tokComma) {
			break
		}
	}
	if _, err := p.expect(tokRParen); err != nil {
		return nil, err
	}
	return seen, nil
}

// parsePropertyKeyList parses unsetter key-only property lists
// (`UNSET TBLPROPERTIES`, `UNSET DBPROPERTIES/PROPERTIES`), forbidding values
// and duplicate keys just as Spark v4.0.5 visitPropertyKeys does.
func (p *Parser) parsePropertyKeyList() error {
	if _, err := p.expect(tokLParen); err != nil {
		return err
	}
	seen := make(map[string]struct{})
	for {
		key, err := p.parsePropertyKey()
		if err != nil {
			return err
		}
		if _, dup := seen[key]; dup {
			return p.errorf("duplicate property key %q", key)
		}
		seen[key] = struct{}{}
		if !p.match(tokComma) {
			break
		}
	}
	_, err := p.expect(tokRParen)
	return err
}

func (p *Parser) parseProperty() (string, error) {
	key, err := p.parsePropertyKey()
	if err != nil {
		return "", err
	}
	p.match(tokEq)
	if err := p.parsePropertyValue(); err != nil {
		return "", err
	}
	return key, nil
}

func (p *Parser) parsePropertyKey() (string, error) {
	if p.cur().Kind == tokString {
		tok := p.advance()
		return decodePropertyKey(tok.Text, true), nil
	}
	if p.cur().Kind == tokDoubleQuoted && p.peek(1).Kind != tokDot {
		tok := p.advance()
		return decodePropertyKey(tok.Text, true), nil
	}
	raw, err := p.parseMultipartIdentifier()
	if err != nil {
		return "", err
	}
	parts := strings.Split(raw, ".")
	for i := range parts {
		parts[i] = exactIdentName(parts[i])
	}
	return strings.Join(parts, "."), nil
}

func (p *Parser) parsePropertyValue() error {
	switch p.cur().Kind {
	case tokInt, tokDecimal, kwTRUE, kwFALSE, tokString, tokDoubleQuoted:
		p.advance()
		return nil
	default:
		return p.errorf("expected property value, got %q", p.cur().Text)
	}
}

func (p *Parser) parseExpressionPropertyList() error {
	if _, err := p.expect(tokLParen); err != nil {
		return err
	}
	if err := p.parseExpressionProperty(); err != nil {
		return err
	}
	for p.match(tokComma) {
		if err := p.parseExpressionProperty(); err != nil {
			return err
		}
	}
	_, err := p.expect(tokRParen)
	return err
}

func (p *Parser) parseExpressionProperty() error {
	if _, err := p.parsePropertyKey(); err != nil {
		return err
	}
	if p.match(tokEq) {
		return p.parseExpression()
	}
	if p.cur().Kind == tokComma || p.cur().Kind == tokRParen {
		return p.errorf("a value must be specified for OPTIONS property key")
	}
	return p.parseExpression()
}

func (p *Parser) parsePartitionSpec(requireValues bool) error {
	if _, err := p.expect(kwPARTITION); err != nil {
		return err
	}
	if _, err := p.expect(tokLParen); err != nil {
		return err
	}
	seen := make(map[string]struct{})
	for {
		key, err := p.parseIdentifier()
		if err != nil {
			return err
		}
		norm := normalizeIdentName(key)
		if _, dup := seen[norm]; dup {
			return p.errorf("duplicate partition key %q in PARTITION specification", key)
		}
		seen[norm] = struct{}{}
		if p.match(tokEq) {
			if p.cur().Kind == kwDEFAULT {
				return p.errorf("DEFAULT is not allowed in PARTITION specification")
			}
			if err := p.parseConstant(); err != nil {
				return err
			}
		} else if requireValues {
			return p.errorf("partition key %q must specify a value", key)
		}
		if !p.match(tokComma) {
			break
		}
	}
	_, err := p.expect(tokRParen)
	return err
}

func (p *Parser) parsePartitionSpecLocation() error {
	if err := p.parsePartitionSpec(true); err != nil {
		return err
	}
	if p.cur().Kind == kwLOCATION {
		return p.parseLocationSpec()
	}
	return nil
}

func (p *Parser) parseTableProvider() error {
	if _, err := p.expect(kwUSING); err != nil {
		return err
	}
	_, err := p.parseMultipartIdentifier()
	return err
}

func (p *Parser) parseClusterBySpec() error {
	if _, err := p.expect(kwCLUSTER); err != nil {
		return err
	}
	if _, err := p.expect(kwBY); err != nil {
		return err
	}
	if _, err := p.expect(tokLParen); err != nil {
		return err
	}
	if err := p.parseMultipartIdentifierList(); err != nil {
		return err
	}
	_, err := p.expect(tokRParen)
	return err
}

func (p *Parser) parseBucketSpec() error {
	if _, err := p.expect(kwCLUSTERED); err != nil {
		return err
	}
	if _, err := p.expect(kwBY); err != nil {
		return err
	}
	if err := p.parseIdentifierList(); err != nil {
		return err
	}
	if p.match(kwSORTED) {
		if _, err := p.expect(kwBY); err != nil {
			return err
		}
		if err := p.parseBucketOrderedIdentifierList(); err != nil {
			return err
		}
	}
	if _, err := p.expect(kwINTO); err != nil {
		return err
	}
	numTok, err := p.expect(tokInt)
	if err != nil {
		return err
	}
	n, convErr := strconv.Atoi(numTok.Text)
	if convErr != nil || n <= 0 {
		return p.errorf("number of buckets must be positive, got %q", numTok.Text)
	}
	_, err = p.expect(kwBUCKETS)
	return err
}

func (p *Parser) parseBucketOrderedIdentifierList() error {
	if _, err := p.expect(tokLParen); err != nil {
		return err
	}
	for {
		if _, err := p.parseErrorCapturingIdentifier(); err != nil {
			return err
		}
		if p.match(kwDESC) {
			return p.errorf("bucket sort column ordering must be ASC")
		}
		p.match(kwASC)
		if !p.match(tokComma) {
			break
		}
	}
	_, err := p.expect(tokRParen)
	return err
}

func (p *Parser) parseSkewSpec() error {
	if _, err := p.expect(kwSKEWED); err != nil {
		return err
	}
	if _, err := p.expect(kwBY); err != nil {
		return err
	}
	if err := p.parseIdentifierList(); err != nil {
		return err
	}
	if _, err := p.expect(kwON); err != nil {
		return err
	}
	if p.peek(1).Kind == tokLParen {
		if err := p.parseNestedConstantList(); err != nil {
			return err
		}
	} else {
		if err := p.parseConstantList(); err != nil {
			return err
		}
	}
	if p.match(kwSTORED) {
		if _, err := p.expect(kwAS); err != nil {
			return err
		}
		if _, err := p.expect(kwDIRECTORIES); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parseConstantList() error {
	if _, err := p.expect(tokLParen); err != nil {
		return err
	}
	if err := p.parseConstant(); err != nil {
		return err
	}
	for p.match(tokComma) {
		if err := p.parseConstant(); err != nil {
			return err
		}
	}
	_, err := p.expect(tokRParen)
	return err
}

func (p *Parser) parseNestedConstantList() error {
	if _, err := p.expect(tokLParen); err != nil {
		return err
	}
	if err := p.parseConstantList(); err != nil {
		return err
	}
	for p.match(tokComma) {
		if err := p.parseConstantList(); err != nil {
			return err
		}
	}
	_, err := p.expect(tokRParen)
	return err
}

func (p *Parser) parseRowFormat() error {
	if _, err := p.expect(kwROW); err != nil {
		return err
	}
	if _, err := p.expect(kwFORMAT); err != nil {
		return err
	}
	if p.match(kwSERDE) {
		if err := p.parseStringLit(); err != nil {
			return err
		}
		if p.match(kwWITH) {
			if _, err := p.expect(kwSERDEPROPERTIES); err != nil {
				return err
			}
			return p.parsePropertyList()
		}
		return nil
	}
	if _, err := p.expect(kwDELIMITED); err != nil {
		return err
	}
	if p.match(kwFIELDS) {
		if _, err := p.expect(kwTERMINATED); err != nil {
			return err
		}
		if _, err := p.expect(kwBY); err != nil {
			return err
		}
		if err := p.parseStringLit(); err != nil {
			return err
		}
		if p.match(kwESCAPED) {
			if _, err := p.expect(kwBY); err != nil {
				return err
			}
			if err := p.parseStringLit(); err != nil {
				return err
			}
		}
	}
	if p.match(kwCOLLECTION) {
		if _, err := p.expect(kwITEMS); err != nil {
			return err
		}
		if _, err := p.expect(kwTERMINATED); err != nil {
			return err
		}
		if _, err := p.expect(kwBY); err != nil {
			return err
		}
		if err := p.parseStringLit(); err != nil {
			return err
		}
	}
	if p.match(kwMAP) {
		if _, err := p.expect(kwKEYS); err != nil {
			return err
		}
		if _, err := p.expect(kwTERMINATED); err != nil {
			return err
		}
		if _, err := p.expect(kwBY); err != nil {
			return err
		}
		if err := p.parseStringLit(); err != nil {
			return err
		}
	}
	if p.match(kwLINES) {
		if _, err := p.expect(kwTERMINATED); err != nil {
			return err
		}
		if _, err := p.expect(kwBY); err != nil {
			return err
		}
		if !p.isStringLit() {
			return p.errorf("expected string literal after LINES TERMINATED BY")
		}
		tok := p.advance()
		if unquoteStringLit(tok.Text) != "\n" {
			return p.errorf("LINES TERMINATED BY only supports newline '\\n'")
		}
	}
	if p.match(kwNULL) {
		if _, err := p.expect(kwDEFINED); err != nil {
			return err
		}
		if _, err := p.expect(kwAS); err != nil {
			return err
		}
		if err := p.parseStringLit(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parseCreateFileFormat(validateHiveSerDe bool) error {
	if _, err := p.expect(kwSTORED); err != nil {
		return err
	}
	if p.match(kwAS) {
		if p.match(kwINPUTFORMAT) {
			if err := p.parseStringLit(); err != nil {
				return err
			}
			if _, err := p.expect(kwOUTPUTFORMAT); err != nil {
				return err
			}
			return p.parseStringLit()
		}
		name, err := p.parseIdentifier()
		if err != nil {
			return err
		}
		if validateHiveSerDe && !isValidHiveSerDeFormat(name) {
			return p.errorf("unsupported STORED AS file format %q", name)
		}
		return nil
	}
	return p.errorf("STORED BY is not valid")
}

func (p *Parser) parsePartitionFieldList() error {
	if _, err := p.expect(tokLParen); err != nil {
		return err
	}
	for {
		if err := p.parsePartitionField(); err != nil {
			return err
		}
		if !p.match(tokComma) {
			break
		}
	}
	_, err := p.expect(tokRParen)
	return err
}

func (p *Parser) parsePartitionField() error {
	// partitionField : transform | colType ;
	// Try transform first if it's ident(args) or ident(.ident)* followed by ',' or ')'.
	cp := p.save()
	if err := p.parseTransform(); err == nil && (p.cur().Kind == tokComma || p.cur().Kind == tokRParen) {
		return nil
	}
	p.restore(cp)
	return p.parseColType()
}

func (p *Parser) parseTransform() error {
	if p.cur().Kind.isIdentifier() && p.peek(1).Kind == tokLParen {
		p.advance()
		p.advance()
		for {
			cp := p.save()
			if err := p.parseConstant(); err == nil && (p.cur().Kind == tokComma || p.cur().Kind == tokRParen) {
				// constant transform argument
			} else {
				p.restore(cp)
				if _, err := p.parseQualifiedName(); err != nil {
					return err
				}
			}
			if !p.match(tokComma) {
				break
			}
		}
		_, err := p.expect(tokRParen)
		return err
	}
	_, err := p.parseQualifiedName()
	return err
}

func (p *Parser) parseRoutineCharacteristics() error {
	var seenLanguage, seenDet, seenDataAccess, seenNullCall, seenComment, seenSecurity bool
	for {
		switch {
		case p.match(kwLANGUAGE):
			if seenLanguage {
				return p.errorf("duplicate LANGUAGE clause in CREATE FUNCTION")
			}
			seenLanguage = true
			if !p.match(kwSQL) {
				return p.errorf("unsupported language for SQL user-defined function: %q", p.cur().Text)
			}
		case p.cur().Kind == kwSPECIFIC:
			return p.errorf("SQL function with SPECIFIC name is not supported")
		case p.match(kwDETERMINISTIC):
			if seenDet {
				return p.errorf("duplicate DETERMINISTIC clause in CREATE FUNCTION")
			}
			seenDet = true
		case p.cur().Kind == kwNOT && p.peek(1).Kind == kwDETERMINISTIC:
			if seenDet {
				return p.errorf("duplicate DETERMINISTIC clause in CREATE FUNCTION")
			}
			seenDet = true
			p.advance()
			p.advance()
		case p.cur().Kind == kwNO:
			return p.errorf("SQL function with NO SQL is not supported")
		case p.match(kwCONTAINS):
			if seenDataAccess {
				return p.errorf("duplicate SQL data access clause in CREATE FUNCTION")
			}
			seenDataAccess = true
			if _, err := p.expect(kwSQL); err != nil {
				return err
			}
		case p.match(kwREADS):
			if seenDataAccess {
				return p.errorf("duplicate SQL data access clause in CREATE FUNCTION")
			}
			seenDataAccess = true
			if _, err := p.expect(kwSQL); err != nil {
				return err
			}
			if _, err := p.expect(kwDATA); err != nil {
				return err
			}
		case p.cur().Kind == kwMODIFIES:
			return p.errorf("SQL function with MODIFIES SQL DATA is not supported")
		case p.cur().Kind == kwRETURNS && p.peek(1).Kind == kwNULL:
			return p.errorf("SQL function with RETURNS NULL ON NULL INPUT is not supported")
		case p.match(kwCALLED):
			if seenNullCall {
				return p.errorf("duplicate CALLED ON NULL INPUT clause in CREATE FUNCTION")
			}
			seenNullCall = true
			if _, err := p.expect(kwON); err != nil {
				return err
			}
			if _, err := p.expect(kwNULL); err != nil {
				return err
			}
			if _, err := p.expect(kwINPUT); err != nil {
				return err
			}
		case p.cur().Kind == kwCOMMENT:
			if seenComment {
				return p.errorf("duplicate COMMENT clause in CREATE FUNCTION")
			}
			seenComment = true
			if err := p.parseCommentSpec(); err != nil {
				return err
			}
		case p.cur().Kind == kwSQL && p.peek(1).Kind == kwSECURITY:
			if seenSecurity {
				return p.errorf("duplicate SQL SECURITY clause in CREATE FUNCTION")
			}
			seenSecurity = true
			p.advance()
			p.advance()
			if !p.match(kwDEFINER) {
				return p.errorf("SQL function with SQL SECURITY %q is not supported", p.cur().Text)
			}
		default:
			return nil
		}
	}
}
