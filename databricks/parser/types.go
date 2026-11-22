package parser

import "strings"

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
func (p *Parser) parseIdentifierReference() (string, error) {
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
	_, err := p.parseIdentifier()
	return err
}

func (p *Parser) parseCollateClause() error {
	if _, err := p.expect(kwCOLLATE); err != nil {
		return err
	}
	_, err := p.parseMultipartIdentifier()
	return err
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
	if err := p.enterNesting(); err != nil {
		return err
	}
	defer p.leaveNesting()

	switch p.cur().Kind {
	case kwARRAY:
		if p.peek(1).Kind == tokLT {
			p.advance()
			p.advance()
			if err := p.parseDataType(); err != nil {
				return err
			}
			return p.expectGT()
		}
	case kwMAP:
		if p.peek(1).Kind == tokLT {
			p.advance()
			p.advance()
			if err := p.parseDataType(); err != nil {
				return err
			}
			if _, err := p.expect(tokComma); err != nil {
				return err
			}
			if err := p.parseDataType(); err != nil {
				return err
			}
			return p.expectGT()
		}
	case kwSTRUCT:
		if p.peek(1).Kind == tokNeq {
			p.advance()
			p.advance()
			return nil
		}
		if p.peek(1).Kind == tokLT {
			p.advance()
			p.advance()
			if p.cur().Kind != tokGT && p.cur().Kind != tokShiftRight && p.cur().Kind != tokShiftRightU {
				if err := p.parseComplexColTypeList(); err != nil {
					return err
				}
			}
			return p.expectGT()
		}
	case kwINTERVAL:
		p.advance()
		if p.match(kwYEAR, kwMONTH) {
			if p.match(kwTO) {
				if _, err := p.expect(kwMONTH); err != nil {
					return err
				}
			}
			return nil
		}
		if p.match(kwDAY, kwHOUR, kwMINUTE, kwSECOND) {
			if p.match(kwTO) {
				if !p.match(kwHOUR, kwMINUTE, kwSECOND) {
					return p.errorf("expected HOUR, MINUTE, or SECOND after TO")
				}
			}
			return nil
		}
		// Plain INTERVAL primitive type (optionally with params).
		return p.parseOptionalTypeParams()
	}

	// Primitive / named type.
	if err := p.parsePrimitiveTypeName(); err != nil {
		return err
	}
	return p.parseOptionalTypeParams()
}

func (p *Parser) parsePrimitiveTypeName() error {
	switch p.cur().Kind {
	case kwSTRING:
		p.advance()
		if p.cur().Kind == kwCOLLATE {
			return p.parseCollateClause()
		}
		return nil
	case kwBOOLEAN, kwTINYINT, kwBYTE, kwSMALLINT, kwSHORT,
		kwINT, kwINTEGER, kwBIGINT, kwLONG, kwFLOAT, kwREAL, kwDOUBLE,
		kwDATE, kwTIMESTAMP, kwTIMESTAMP_NTZ, kwTIMESTAMP_LTZ,
		kwCHARACTER, kwCHAR, kwVARCHAR, kwBINARY,
		kwDECIMAL, kwDEC, kwNUMERIC, kwVOID, kwINTERVAL, kwVARIANT,
		kwARRAY, kwSTRUCT, kwMAP:
		p.advance()
		return nil
	default:
		_, err := p.parseIdentifier()
		return err
	}
}

func (p *Parser) parseOptionalTypeParams() error {
	if p.cur().Kind != tokLParen {
		return nil
	}
	cp := p.save()
	p.advance()
	if p.cur().Kind != tokInt {
		p.restore(cp)
		return nil
	}
	p.advance()
	for p.match(tokComma) {
		if _, err := p.expect(tokInt); err != nil {
			p.restore(cp)
			return err
		}
	}
	if _, err := p.expect(tokRParen); err != nil {
		return err
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
	if err := p.parseDataType(); err != nil {
		return err
	}
	var seenNotNull, seenDefault, seenGenerated, seenComment bool
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
		case p.cur().Kind == kwGENERATED:
			if seenGenerated {
				return p.errorf("duplicate GENERATED descriptor in column definition")
			}
			seenGenerated = true
			if err := p.parseGenerationExpression(); err != nil {
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

func (p *Parser) parseGenerationExpression() error {
	if _, err := p.expect(kwGENERATED); err != nil {
		return err
	}
	if p.match(kwALWAYS) {
		if _, err := p.expect(kwAS); err != nil {
			return err
		}
		if p.match(kwIDENTITY) {
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
	return p.parseOptionalIdentityColSpec()
}

func (p *Parser) parseOptionalIdentityColSpec() error {
	if !p.match(tokLParen) {
		return nil
	}
	for p.cur().Kind != tokRParen && p.cur().Kind != tokEOF {
		if p.match(kwSTART) {
			if _, err := p.expect(kwWITH); err != nil {
				return err
			}
			p.match(tokMinus)
			if !p.match(tokInt, tokBigInt) {
				return p.errorf("expected integer or bigint for START WITH")
			}
		} else if p.match(kwINCREMENT) {
			if _, err := p.expect(kwBY); err != nil {
				return err
			}
			p.match(tokMinus)
			if !p.match(tokInt, tokBigInt) {
				return p.errorf("expected integer or bigint for INCREMENT BY")
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

// parsePropertyList parses: LEFT_PAREN property (COMMA property)* RIGHT_PAREN
func (p *Parser) parsePropertyList() error {
	if _, err := p.expect(tokLParen); err != nil {
		return err
	}
	if err := p.parseProperty(); err != nil {
		return err
	}
	for p.match(tokComma) {
		if err := p.parseProperty(); err != nil {
			return err
		}
	}
	_, err := p.expect(tokRParen)
	return err
}

func (p *Parser) parseProperty() error {
	if err := p.parsePropertyKey(); err != nil {
		return err
	}
	if p.match(tokEq) {
		return p.parsePropertyValue()
	}
	if p.isPropertyValueStart() {
		return p.parsePropertyValue()
	}
	return nil
}

func (p *Parser) parsePropertyKey() error {
	if p.cur().Kind == tokString {
		p.advance()
		return nil
	}
	if p.cur().Kind == tokDoubleQuoted && p.peek(1).Kind != tokDot {
		p.advance()
		return nil
	}
	_, err := p.parseMultipartIdentifier()
	return err
}

func (p *Parser) isPropertyValueStart() bool {
	switch p.cur().Kind {
	case tokInt, tokDecimal, kwTRUE, kwFALSE, tokString, tokDoubleQuoted:
		return true
	}
	return false
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
	if err := p.parsePropertyKey(); err != nil {
		return err
	}
	if p.match(tokEq) {
		return p.parseExpression()
	}
	if p.cur().Kind != tokComma && p.cur().Kind != tokRParen {
		return p.parseExpression()
	}
	return nil
}

func (p *Parser) parsePartitionSpec() error {
	if _, err := p.expect(kwPARTITION); err != nil {
		return err
	}
	if _, err := p.expect(tokLParen); err != nil {
		return err
	}
	for {
		if _, err := p.parseIdentifier(); err != nil {
			return err
		}
		if p.match(tokEq) {
			if !p.match(kwDEFAULT) {
				if err := p.parseConstant(); err != nil {
					return err
				}
			}
		}
		if !p.match(tokComma) {
			break
		}
	}
	_, err := p.expect(tokRParen)
	return err
}

func (p *Parser) parsePartitionSpecLocation() error {
	if err := p.parsePartitionSpec(); err != nil {
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
		if err := p.parseOrderedIdentifierList(); err != nil {
			return err
		}
	}
	if _, err := p.expect(kwINTO); err != nil {
		return err
	}
	if _, err := p.expect(tokInt); err != nil {
		return err
	}
	_, err := p.expect(kwBUCKETS)
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
		if err := p.parseStringLit(); err != nil {
			return err
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

func (p *Parser) parseCreateFileFormat() error {
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
		_, err := p.parseIdentifier()
		return err
	}
	if _, err := p.expect(kwBY); err != nil {
		return err
	}
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
	for {
		switch {
		case p.match(kwLANGUAGE):
			if !p.match(kwSQL, tokIdent) {
				return p.errorf("expected SQL or language identifier")
			}
		case p.match(kwSPECIFIC):
			if _, err := p.parseErrorCapturingIdentifier(); err != nil {
				return err
			}
		case p.match(kwDETERMINISTIC):
		case p.isNot() && p.peek(1).Kind == kwDETERMINISTIC:
			p.advance()
			p.advance()
		case p.match(kwNO):
			if _, err := p.expect(kwSQL); err != nil {
				return err
			}
		case p.match(kwCONTAINS):
			if _, err := p.expect(kwSQL); err != nil {
				return err
			}
		case p.match(kwREADS, kwMODIFIES):
			if _, err := p.expect(kwSQL); err != nil {
				return err
			}
			if _, err := p.expect(kwDATA); err != nil {
				return err
			}
		case p.cur().Kind == kwRETURNS && p.peek(1).Kind == kwNULL:
			p.advance()
			p.advance()
			if _, err := p.expect(kwON); err != nil {
				return err
			}
			if _, err := p.expect(kwNULL); err != nil {
				return err
			}
			if _, err := p.expect(kwINPUT); err != nil {
				return err
			}
		case p.match(kwCALLED):
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
			if err := p.parseCommentSpec(); err != nil {
				return err
			}
		case p.match(kwSQL):
			if _, err := p.expect(kwSECURITY); err != nil {
				return err
			}
			if !p.match(kwINVOKER, kwDEFINER) {
				return p.errorf("expected INVOKER or DEFINER")
			}
		default:
			return nil
		}
	}
}
