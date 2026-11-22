package parser

func (p *Parser) parseExpression() error {
	return p.parseBooleanExpression()
}

func (p *Parser) parseExpressionSeq() error {
	if err := p.parseExpression(); err != nil {
		return err
	}
	for p.match(tokComma) {
		if err := p.parseExpression(); err != nil {
			return err
		}
	}
	return nil
}

// parseNamedExpression parses:
//
//	namedExpression : expression (AS? (errorCapturingIdentifier | identifierList))? ;
func (p *Parser) parseNamedExpression() error {
	if err := p.parseExpression(); err != nil {
		return err
	}
	return p.parseOptionalNamedExprAlias()
}

func (p *Parser) parseOptionalNamedExprAlias() error {
	if p.match(kwAS) {
		if p.cur().Kind == tokLParen {
			return p.parseIdentifierList()
		}
		_, err := p.parseErrorCapturingIdentifier()
		return err
	}
	if p.cur().Kind == tokLParen {
		cp := p.save()
		if err := p.parseIdentifierList(); err == nil {
			return nil
		}
		p.restore(cp)
		return nil
	}
	if p.canBeImplicitExprAlias() {
		_, err := p.parseErrorCapturingIdentifier()
		return err
	}
	return nil
}

// canBeImplicitExprAlias checks whether the current token can serve as an
// implicit (no AS) column alias after an expression, without stealing a
// keyword that introduces a following SQL clause.
func (p *Parser) canBeImplicitExprAlias() bool {
	k := p.cur().Kind
	if !k.isIdentifier() {
		return false
	}
	if k == tokIdent || k == tokBackquoted || k == tokDoubleQuoted {
		if k == tokIdent && isWord(p.cur(), "QUALIFY") && !p.isStatementEnd() && p.peek(1).Kind != tokRParen {
			return false
		}
		return true
	}
	// If the next token ends the current expression list / statement / paren,
	// a keyword can safely be an implicit alias (e.g. SELECT a FROM, SELECT a ALL,
	// SELECT x asc, SELECT x rows).
	next := p.peek(1).Kind
	if next == tokEOF || next == tokSemicolon || next == tokComma || next == tokRParen {
		switch k {
		case kwEND, kwELSE, kwELSEIF, kwTHEN, kwWHEN, kwUNTIL, kwDO:
			return false
		}
		return true
	}
	// Allow a keyword as implicit column alias right before FROM/WHERE/GROUP/HAVING/WINDOW/ORDER/LIMIT
	// when the keyword itself doesn't start a clause (e.g. `SELECT a IS FROM a`,
	// `SELECT x filter FROM t`, `SELECT x over FROM t`, `SELECT x rows FROM t`,
	// `SELECT x nulls FROM t`, `SELECT x asc FROM t`).
	if k == kwIS && next == kwFROM {
		return true
	}
	switch k {
	case kwOVER, kwFILTER, kwWITHIN, kwIGNORE, kwRESPECT,
		kwASC, kwDESC, kwNULLS, kwFIRST, kwLAST, kwROWS, kwRANGE, kwPRECEDING, kwFOLLOWING:
		switch next {
		case tokOperatorPipe, kwFROM, kwWHERE, kwGROUP, kwHAVING, kwWINDOW,
			kwORDER, kwCLUSTER, kwDISTRIBUTE, kwSORT, kwLIMIT, kwOFFSET,
			kwUNION, kwINTERSECT, kwEXCEPT, kwMINUS, kwLATERAL, kwFOR:
			return true
		case tokIdent:
			if isWord(p.peek(1), "QUALIFY") && p.peek(2).Kind != tokEOF && p.peek(2).Kind != tokSemicolon && p.peek(2).Kind != tokRParen {
				return true
			}
		}
		return false
	}
	// Otherwise disallow clause-introducing keywords as implicit aliases.
	switch k {
	case kwFROM, kwWHERE, kwGROUP, kwGROUPING, kwHAVING, kwWINDOW, kwORDER, kwCLUSTER, kwDISTRIBUTE,
		kwSORT, kwLIMIT, kwOFFSET, kwUNION, kwINTERSECT, kwEXCEPT, kwMINUS,
		kwLATERAL, kwPIVOT, kwUNPIVOT, kwJOIN, kwINNER, kwLEFT, kwRIGHT, kwFULL,
		kwCROSS, kwNATURAL, kwSEMI, kwANTI, kwON, kwUSING, kwINTO, kwSET, kwVALUES,
		kwSELECT, kwINSERT, kwWHEN, kwTHEN, kwELSE, kwELSEIF, kwEND, kwWITH, kwFOR, kwAS,
		kwDO, kwUNTIL, kwBY, kwIN, kwIS, kwAND, kwOR, kwNOT, kwBETWEEN, kwLIKE,
		kwILIKE, kwRLIKE, kwESCAPE, kwOVER, kwFILTER, kwWITHIN, kwIGNORE, kwRESPECT,
		kwASC, kwDESC, kwNULLS, kwROWS, kwRANGE, kwPRECEDING, kwFOLLOWING,
		kwRECORDWRITER, kwRECORDREADER, kwROW, kwPLACING, kwRETURNS, kwRETURN:
		return false
	}
	return true
}

func (p *Parser) parseNamedExpressionSeq() error {
	if err := p.parseNamedExpression(); err != nil {
		return err
	}
	for p.match(tokComma) {
		if err := p.parseNamedExpression(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parseFunctionArgument() error {
	_, err := p.parseFunctionArgumentWithKind()
	return err
}

// parseFunctionArgumentWithKind parses a functionArgument and reports whether
// it is a plain valueExpression (not a named argument, logical NOT/AND/OR,
// EXISTS, or predicated expression). SQL-specific forms like POSITION(a IN b),
// SUBSTRING(s FROM pos [FOR len]), TRIM(c FROM s), and OVERLAY(s PLACING r ...)
// require valueExpression operands in SqlBaseParser.g4.
func (p *Parser) parseFunctionArgumentWithKind() (bool, error) {
	if p.cur().Kind.isIdentifier() && p.peek(1).Kind == tokFatArrow {
		p.advance()
		p.advance()
		return false, p.parseExpression()
	}
	return p.parseBooleanExpressionWithKind()
}

// booleanExpression precedence: OR < AND < NOT/! < EXISTS / predicated.
func (p *Parser) parseBooleanExpression() error {
	_, err := p.parseBooleanExpressionWithKind()
	return err
}

func (p *Parser) parseBooleanExpressionWithKind() (bool, error) {
	if err := p.enter(); err != nil {
		return false, err
	}
	defer p.leave()
	return p.parseOrExpressionWithKind()
}

func (p *Parser) parseOrExpressionWithKind() (bool, error) {
	isVal, err := p.parseAndExpressionWithKind()
	if err != nil {
		return false, err
	}
	return p.parseOrTailWithKind(isVal)
}

func (p *Parser) parseOrTailWithKind(isVal bool) (bool, error) {
	for p.cur().Kind == kwOR {
		cp := p.save()
		p.advance()
		if _, err := p.parseAndExpressionWithKind(); err != nil {
			p.restore(cp)
			break
		}
		isVal = false
	}
	return isVal, nil
}

func (p *Parser) parseAndExpressionWithKind() (bool, error) {
	isVal, err := p.parseNotExpressionWithKind()
	if err != nil {
		return false, err
	}
	return p.parseAndTailWithKind(isVal)
}

func (p *Parser) parseAndTailWithKind(isVal bool) (bool, error) {
	for p.cur().Kind == kwAND {
		cp := p.save()
		p.advance()
		if _, err := p.parseNotExpressionWithKind(); err != nil {
			p.restore(cp)
			break
		}
		isVal = false
	}
	return isVal, nil
}

func (p *Parser) parseNotExpressionWithKind() (bool, error) {
	ops := 0
	for p.cur().Kind == kwNOT || p.cur().Kind == tokBang {
		if p.cur().Kind == kwNOT && p.isNotIdentifierContinuation() {
			// NOT is in Spark's nonReserved list, so `not` followed by a clause
			// boundary, postfix operator (`.`, `[`, `::`, `COLLATE`), or
			// non-unary binary/predicate operator is an identifier expression
			// rather than unary NOT.
			break
		}
		ops++
		if p.depth+ops > maxNestingDepth {
			return false, p.errorf("maximum nesting depth exceeded")
		}
		p.advance()
	}
	if p.cur().Kind == kwEXISTS && p.peek(1).Kind == tokLParen {
		cp := p.save()
		p.advance()
		p.advance()
		if p.startsWithCompleteParenQuery() {
			if err := p.parseQuery(); err != nil {
				return false, err
			}
			_, err := p.expect(tokRParen)
			return false, err
		}
		p.restore(cp)
	}
	if err := p.parseValueExpression(); err != nil {
		return false, err
	}
	beforePred := p.pos
	if err := p.parseOptionalPredicate(); err != nil {
		return false, err
	}
	return ops == 0 && p.pos == beforePred, nil
}

func (p *Parser) isNotIdentifierContinuation() bool {
	next := p.peek(1).Kind
	switch next {
	case tokEOF, tokSemicolon, tokComma, tokRParen, tokRBracket, tokOperatorPipe,
		tokDot, tokLBracket, tokDoubleColon,
		tokAsterisk, tokSlash, tokPercent, kwDIV,
		tokConcatPipe, tokPipe, tokAmpersand, tokHat,
		tokShiftLeft, tokShiftRight, tokShiftRightU,
		tokEq, tokNeq, tokNeqBang, tokNullSafeEq, tokLT, tokLTE, tokGT, tokGTE,
		kwAND, kwOR, kwIS, kwIN, kwLIKE, kwILIKE, kwRLIKE, kwBETWEEN,
		tokArrow, tokFatArrow:
		return true
	case kwGROUP, kwORDER, kwCLUSTER, kwDISTRIBUTE, kwSORT:
		return p.peek(2).Kind == kwBY
	case kwWINDOW:
		return p.peek(2).Kind.isIdentifier() && p.peek(3).Kind == kwAS
	case kwUNION, kwINTERSECT, kwEXCEPT, kwMINUS:
		switch p.peek(2).Kind {
		case kwALL, kwDISTINCT, kwSELECT, kwTABLE, kwVALUES, kwMAP, kwREDUCE, tokLParen:
			return true
		}
		return false
	case kwUSING:
		after := p.peek(2).Kind
		return after == tokLParen || after == tokString || after == tokDoubleQuoted
	case kwEND:
		return p.openCasesAtCurrentParenDepth() > 0
	case kwAS:
		after := p.peek(2).Kind
		if after == tokLParen {
			return true
		}
		if !after.isIdentifier() {
			return false
		}
		return !isNotOperandClauseEnd(after)
	case kwFROM, kwWHERE, kwHAVING, kwLIMIT, kwOFFSET, kwTHEN, kwELSE, kwWHEN, kwON, kwINTO:
		return !isNotOperandFollow(p.peek(2).Kind)
	case kwCOLLATE:
		// `not COLLATE utf8_binary` is a postfix COLLATE on column `not`,
		// whereas `NOT collate FROM t` negates a column named `collate`.
		after := p.peek(2).Kind
		return after == tokIdent || after == tokBackquoted || after == tokDoubleQuoted
	}
	return false
}

func isNotOperandClauseEnd(k TokenKind) bool {
	switch k {
	case kwFROM, kwWHERE, kwGROUP, kwHAVING, kwWINDOW, kwORDER, kwCLUSTER,
		kwDISTRIBUTE, kwSORT, kwLIMIT, kwOFFSET, kwUNION, kwINTERSECT, kwEXCEPT, kwMINUS,
		kwTHEN, kwELSE, kwEND, kwWHEN, kwON, kwUSING, kwINTO:
		return true
	}
	return false
}

func isNotOperandFollow(k TokenKind) bool {
	if isNotOperandClauseEnd(k) {
		return true
	}
	switch k {
	case tokEOF, tokSemicolon, tokComma, tokRParen, tokRBracket, tokOperatorPipe,
		tokDot, tokLBracket, tokDoubleColon, kwCOLLATE,
		tokAsterisk, tokSlash, tokPercent, kwDIV,
		tokConcatPipe, tokPipe, tokAmpersand, tokHat,
		tokShiftLeft, tokShiftRight, tokShiftRightU,
		tokEq, tokNeq, tokNeqBang, tokNullSafeEq, tokLT, tokLTE, tokGT, tokGTE,
		kwAND, kwOR, kwIS, kwIN, kwLIKE, kwILIKE, kwRLIKE, kwBETWEEN,
		tokArrow, tokFatArrow:
		return true
	}
	return false
}

func (p *Parser) parseOptionalPredicate() error {
	// IS [NOT|!] (NULL | TRUE | FALSE | UNKNOWN | DISTINCT FROM valueExpression)
	if p.cur().Kind == kwIS {
		cp := p.save()
		p.advance()
		p.matchNot()
		if p.match(kwNULL, kwTRUE, kwFALSE, kwUNKNOWN) {
			return nil
		}
		if p.match(kwDISTINCT) {
			if _, err := p.expect(kwFROM); err != nil {
				return err
			}
			return p.parseValueExpression()
		}
		p.restore(cp)
		return nil
	}

	// [NOT|!] (BETWEEN | IN | RLIKE | LIKE | ILIKE)
	cp := p.save()
	hasNot := p.matchNot()
	switch p.cur().Kind {
	case kwBETWEEN:
		p.advance()
		if err := p.parseValueExpression(); err != nil {
			return err
		}
		if _, err := p.expect(kwAND); err != nil {
			return err
		}
		return p.parseValueExpression()

	case kwIN:
		if p.peek(1).Kind != tokLParen {
			if hasNot {
				p.restore(cp)
				return nil
			}
			return nil
		}
		p.advance()
		p.advance()
		// Could be IN (query) or IN (expression (, expression)*)
		if p.startsWithCompleteParenQuery() {
			if err := p.parseQuery(); err != nil {
				return err
			}
			_, err := p.expect(tokRParen)
			return err
		}
		if err := p.parseExpression(); err != nil {
			return err
		}
		for p.match(tokComma) {
			if err := p.parseExpression(); err != nil {
				return err
			}
		}
		_, err := p.expect(tokRParen)
		return err

	case kwRLIKE:
		p.advance()
		if err := p.parseValueExpression(); err != nil {
			if !hasNot {
				p.restore(cp)
				return nil
			}
			return err
		}
		return nil

	case kwLIKE, kwILIKE:
		p.advance()
		if p.cur().Kind == kwANY || p.cur().Kind == kwSOME || p.cur().Kind == kwALL {
			if p.peek(1).Kind == tokLParen {
				p.advance()
				p.advance()
				if !p.match(tokRParen) {
					if err := p.parseExpression(); err != nil {
						return err
					}
					for p.match(tokComma) {
						if err := p.parseExpression(); err != nil {
							return err
						}
					}
					if _, err := p.expect(tokRParen); err != nil {
						return err
					}
				}
				return nil
			}
		}
		if err := p.parseValueExpression(); err != nil {
			if !hasNot {
				p.restore(cp)
				return nil
			}
			return err
		}
		if p.cur().Kind == kwESCAPE && p.isStringLitAt(1) {
			p.advance()
			return p.parseStringLit()
		}
		return nil
	}

	if hasNot {
		p.restore(cp)
	}
	return nil
}

// valueExpression precedence (lowest to highest, matching SqlBaseParser.g4):
// comparison < | < ^ < & < shift (<< >> >>>) < + - || < * / % DIV < unary (+ - ~) < primary
func (p *Parser) parseValueExpression() error {
	if err := p.enter(); err != nil {
		return err
	}
	defer p.leave()
	return p.parseComparisonExpression()
}

func (p *Parser) parseComparisonExpression() error {
	if err := p.parseBitOrExpression(); err != nil {
		return err
	}
	for p.isComparisonOp() {
		p.advance()
		if err := p.parseBitOrExpression(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) isComparisonOp() bool {
	switch p.cur().Kind {
	case tokEq, tokNeq, tokNeqBang, tokLT, tokLTE, tokGT, tokGTE, tokNullSafeEq:
		return true
	}
	return false
}

func (p *Parser) parseBitOrExpression() error {
	if err := p.parseBitXorExpression(); err != nil {
		return err
	}
	for p.match(tokPipe) {
		if err := p.parseBitXorExpression(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parseBitXorExpression() error {
	if err := p.parseBitAndExpression(); err != nil {
		return err
	}
	for p.match(tokHat) {
		if err := p.parseBitAndExpression(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parseBitAndExpression() error {
	if err := p.parseShiftExpression(); err != nil {
		return err
	}
	for p.match(tokAmpersand) {
		if err := p.parseShiftExpression(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parseShiftExpression() error {
	if err := p.parseAddExpression(); err != nil {
		return err
	}
	for p.match(tokShiftLeft, tokShiftRight, tokShiftRightU) {
		if err := p.parseAddExpression(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parseAddExpression() error {
	if err := p.parseMulExpression(); err != nil {
		return err
	}
	for p.match(tokPlus, tokMinus, tokConcatPipe) {
		if err := p.parseMulExpression(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parseMulExpression() error {
	if err := p.parseUnaryExpression(); err != nil {
		return err
	}
	for p.match(tokAsterisk, tokSlash, tokPercent, kwDIV) {
		if err := p.parseUnaryExpression(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parseUnaryExpression() error {
	ops := 0
	for p.match(tokMinus, tokPlus, tokTilde) {
		ops++
		if p.depth+ops > maxNestingDepth {
			return p.errorf("maximum nesting depth exceeded")
		}
	}
	return p.parsePrimaryExpression()
}

// parsePrimaryExpression parses Spark's primaryExpression + postfix operators
// (COLLATE, :: dataType, [index], .field).
func (p *Parser) parsePrimaryExpression() error {
	if err := p.enter(); err != nil {
		return err
	}
	defer p.leave()

	isQualifiedIdent, err := p.parsePrimaryAtom()
	if err != nil {
		return err
	}

	// Postfix loop: COLLATE, :: dataType, [index], .field (or .* [EXCEPT (...)] on qualified identifiers).
	for {
		switch p.cur().Kind {
		case kwCOLLATE:
			cp := p.save()
			if err := p.parseCollateClause(); err != nil {
				p.restore(cp)
				return nil
			}
			isQualifiedIdent = false
		case tokDoubleColon:
			p.advance()
			if err := p.parseDataType(); err != nil {
				return err
			}
			isQualifiedIdent = false
		case tokLBracket:
			p.advance()
			if err := p.parseValueExpression(); err != nil {
				return err
			}
			if _, err := p.expect(tokRBracket); err != nil {
				return err
			}
			isQualifiedIdent = false
		case tokDot:
			p.advance()
			if p.cur().Kind == tokAsterisk {
				if !isQualifiedIdent {
					return p.errorf("qualified star '.*' is only allowed on a qualified identifier")
				}
				p.advance()
				if p.cur().Kind == kwEXCEPT && p.peek(1).Kind == tokLParen {
					if err := p.parseExceptClause(); err != nil {
						return err
					}
				}
				isQualifiedIdent = false
			} else {
				if _, err := p.parseIdentifier(); err != nil {
					return err
				}
			}
		default:
			return nil
		}
	}
}

func (p *Parser) parseExceptClause() error {
	if _, err := p.expect(kwEXCEPT); err != nil {
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

func (p *Parser) parsePrimaryAtom() (bool, error) {
	k := p.cur().Kind

	// Lambda: identifier -> expression
	if k.isIdentifier() && p.peek(1).Kind == tokArrow {
		p.advance()
		p.advance()
		return false, p.parseExpression()
	}

	// Parenthesized forms:
	// - (ident, ident+) -> expression   (multi-param lambda)
	// - (query)                         (subquery expression)
	// - (namedExpr, namedExpr+)         (row constructor)
	// - (expression)                    (parenthesized expression)
	if k == tokLParen {
		return false, p.parseParenthesizedPrimary()
	}

	// Star: * [EXCEPT (...)]
	if k == tokAsterisk {
		p.advance()
		if p.cur().Kind == kwEXCEPT && p.peek(1).Kind == tokLParen {
			return false, p.parseExceptClause()
		}
		return false, nil
	}

	// CASE ... END (searched or simple) vs ordinary function call CASE(...) or column `case`.
	if k == kwCASE {
		if p.peek(1).Kind == kwWHEN && (p.openCasesAtCurrentParenDepth() == 0 || p.hasCaseWhenAndEndAhead()) {
			return false, p.parseCaseExpression()
		}
		if p.peek(1).Kind == tokLParen {
			if p.hasCaseWhenAndEndAhead() {
				return false, p.parseCaseExpression()
			}
			return false, p.parseCaseOrCaseFuncCall()
		}
		if !p.isCaseColumnReference() && p.hasCaseWhenAndEndAhead() {
			return false, p.parseCaseExpression()
		}
	}

	// Special keyword expressions that take parentheses.
	if p.peek(1).Kind == tokLParen && isSpecialKeywordCall(k) {
		return false, p.parseSpecialKeywordCall(k)
	}

	// Function call: functionName (...) [WITHIN GROUP (...)] [FILTER (WHERE ...)] [(IGNORE|RESPECT) NULLS] [OVER windowSpec]
	if p.canStartFunctionCall() {
		cp := p.save()
		if err := p.parseFunctionCall(); err == nil {
			return false, nil
		}
		p.restore(cp)
	}

	// Constants (NULL, ?, :ident, INTERVAL ..., literalType '...', number, TRUE/FALSE, stringLit+).
	if p.canStartConstant() {
		cp := p.save()
		if err := p.parseConstant(); err == nil {
			return false, nil
		}
		p.restore(cp)
	}

	// currentLike keywords: CURRENT_DATE | CURRENT_TIMESTAMP | CURRENT_USER | USER | SESSION_USER
	if k == kwCURRENT_DATE || k == kwCURRENT_TIMESTAMP || k == kwCURRENT_USER || k == kwUSER || k == kwSESSION_USER {
		p.advance()
		return k.isIdentifier() && p.cur().Kind == tokDot, nil
	}

	// Column reference (identifier).
	if k.isIdentifier() {
		p.advance()
		return true, nil
	}

	return false, p.errorf("unexpected token %q in expression", p.cur().Text)
}

func isSpecialKeywordCall(k TokenKind) bool {
	switch k {
	case kwCAST, kwTRY_CAST, kwSTRUCT, kwFIRST, kwANY_VALUE, kwLAST,
		kwPOSITION, kwEXTRACT, kwSUBSTR, kwSUBSTRING, kwTRIM, kwOVERLAY,
		kwTIMESTAMPADD, kwDATEADD, kwDATE_ADD, kwTIMESTAMPDIFF, kwDATEDIFF, kwDATE_DIFF, kwTIMEDIFF:
		return true
	}
	return false
}

// isPositionInValueExpr reports whether `IN` at p.cur() inside POSITION(...)
// introduces SQL `POSITION(substr IN valueExpression)` rather than an
// `IN (list|query)` membership predicate of an ordinary function call argument.
// When `IN` is followed by `(...)` and then a valueExpression postfix or binary
// operator (such as `||`, `::`, `[`, `.`, `+`, `COLLATE`), the operand can only
// be a complete valueExpression in the SQL-specific POSITION form.
func (p *Parser) isPositionInValueExpr() bool {
	if p.cur().Kind != kwIN {
		return false
	}
	if p.peek(1).Kind != tokLParen {
		return true
	}
	depth := 0
	closeIdx := -1
	for i := 1; ; i++ {
		k := p.peek(i).Kind
		if k == tokEOF || k == tokSemicolon {
			return false
		}
		if k == tokLParen || k == tokLBracket {
			depth++
		} else if k == tokRParen || k == tokRBracket {
			depth--
			if depth == 0 {
				closeIdx = i
				break
			}
		}
	}
	switch p.peek(closeIdx + 1).Kind {
	case tokDot, tokLBracket, tokDoubleColon, kwCOLLATE,
		tokAsterisk, tokSlash, tokPercent, kwDIV,
		tokPlus, tokMinus, tokConcatPipe,
		tokShiftLeft, tokShiftRight, tokShiftRightU,
		tokAmpersand, tokHat, tokPipe,
		tokEq, tokNeq, tokNeqBang, tokNullSafeEq, tokLT, tokLTE, tokGT, tokGTE:
		return true
	}
	return false
}

func (p *Parser) parseSpecialKeywordCall(k TokenKind) error {
	switch k {
	case kwCAST, kwTRY_CAST:
		// Single-pass parsing for CAST/TRY_CAST(expr AS dataType) and ordinary function calls.
		p.advance()
		p.advance()
		if p.match(tokRParen) {
			return p.parseFunctionCallSuffixes()
		}
		hasQuantifier := p.isFunctionSetQuantifier()
		if hasQuantifier {
			p.advance()
		}
		isNamed := p.cur().Kind.isIdentifier() && p.peek(1).Kind == tokFatArrow
		if err := p.parseFunctionArgument(); err != nil {
			return err
		}
		if !hasQuantifier && !isNamed && p.match(kwAS) {
			if err := p.parseDataType(); err != nil {
				return err
			}
			_, err := p.expect(tokRParen)
			return err
		}
		for p.match(tokComma) {
			if err := p.parseFunctionArgument(); err != nil {
				return err
			}
		}
		if _, err := p.expect(tokRParen); err != nil {
			return err
		}
		return p.parseFunctionCallSuffixes()

	case kwSTRUCT:
		// Single-pass parsing for STRUCT(namedExpr | key => expr, ...) so nested
		// STRUCT(..., k => 1) calls never backtrack exponentially.
		p.advance()
		p.advance()
		if p.match(tokRParen) {
			return p.parseFunctionCallSuffixes()
		}
		hasQuantifier := p.isFunctionSetQuantifier()
		if hasQuantifier {
			p.advance()
		}
		hasAlias := false
		hasNamedArg := false
		for {
			if p.cur().Kind.isIdentifier() && p.peek(1).Kind == tokFatArrow {
				if hasAlias {
					return p.errorf("cannot mix field aliases and named arguments in STRUCT")
				}
				hasNamedArg = true
				p.advance()
				p.advance()
				if err := p.parseExpression(); err != nil {
					return err
				}
			} else {
				if err := p.parseExpression(); err != nil {
					return err
				}
				if !hasQuantifier && !hasNamedArg && p.cur().Kind != tokComma && p.cur().Kind != tokRParen {
					before := p.pos
					if err := p.parseOptionalNamedExprAlias(); err != nil {
						return err
					}
					if p.pos > before {
						hasAlias = true
					}
				}
			}
			if !p.match(tokComma) {
				break
			}
		}
		if _, err := p.expect(tokRParen); err != nil {
			return err
		}
		if hasAlias {
			return nil
		}
		return p.parseFunctionCallSuffixes()

	case kwFIRST, kwANY_VALUE, kwLAST:
		p.advance()
		p.advance()
		if p.match(tokRParen) {
			return p.parseFunctionCallSuffixes()
		}
		// Handle optional DISTINCT/ALL followed by arguments, or Spark's
		// `FIRST(expr IGNORE NULLS)` / `LAST(expr IGNORE NULLS)` / `ANY_VALUE(expr IGNORE NULLS)`,
		// without reparsing the argument expression.
		hasQuantifier := p.isFunctionSetQuantifier() &&
			!(p.peek(1).Kind == kwIGNORE && p.peek(2).Kind == kwNULLS && p.peek(3).Kind == tokRParen)
		if hasQuantifier {
			p.advance()
		}
		isNamed := p.cur().Kind.isIdentifier() && p.peek(1).Kind == tokFatArrow
		if err := p.parseFunctionArgument(); err != nil {
			return err
		}
		if !hasQuantifier && !isNamed && p.cur().Kind == kwIGNORE && p.peek(1).Kind == kwNULLS && p.peek(2).Kind == tokRParen {
			p.advance()
			p.advance()
			p.advance()
			return nil
		}
		for p.match(tokComma) {
			if err := p.parseFunctionArgument(); err != nil {
				return err
			}
		}
		if _, err := p.expect(tokRParen); err != nil {
			return err
		}
		return p.parseFunctionCallSuffixes()

	case kwPOSITION:
		// Single-pass parsing for SQL POSITION(substr IN str) and ordinary POSITION(a, b) calls.
		p.advance()
		p.advance()
		if p.match(tokRParen) {
			return p.parseFunctionCallSuffixes()
		}
		hasQuantifier := p.isFunctionSetQuantifier()
		if hasQuantifier {
			p.advance()
		}
		isNamed := p.cur().Kind.isIdentifier() && p.peek(1).Kind == tokFatArrow
		startsBoolOnly := p.cur().Kind == tokBang ||
			(p.cur().Kind == kwNOT && !p.isNotIdentifierContinuation())
		if !startsBoolOnly && p.cur().Kind == kwEXISTS && p.peek(1).Kind == tokLParen {
			cpExists := p.save()
			p.advance()
			p.advance()
			startsBoolOnly = p.startsWithCompleteParenQuery()
			p.restore(cpExists)
		}
		if hasQuantifier || isNamed || startsBoolOnly {
			if err := p.parseFunctionArgument(); err != nil {
				return err
			}
		} else {
			if err := p.parseValueExpression(); err != nil {
				return err
			}
			if p.isPositionInValueExpr() {
				p.advance() // IN
				if err := p.parseValueExpression(); err != nil {
					return err
				}
				_, err := p.expect(tokRParen)
				return err
			}
			if err := p.parseOptionalPredicate(); err != nil {
				return err
			}
			if _, err := p.parseAndTailWithKind(false); err != nil {
				return err
			}
			if _, err := p.parseOrTailWithKind(false); err != nil {
				return err
			}
		}
		for p.match(tokComma) {
			if err := p.parseFunctionArgument(); err != nil {
				return err
			}
		}
		if _, err := p.expect(tokRParen); err != nil {
			return err
		}
		return p.parseFunctionCallSuffixes()

	case kwEXTRACT:
		// Single-pass parsing for SQL EXTRACT(field FROM src) and ordinary EXTRACT(a, b) calls.
		p.advance()
		p.advance()
		if p.match(tokRParen) {
			return p.parseFunctionCallSuffixes()
		}
		if p.cur().Kind.isIdentifier() && p.peek(1).Kind == kwFROM {
			p.advance()
			p.advance()
			if err := p.parseValueExpression(); err != nil {
				return err
			}
			_, err := p.expect(tokRParen)
			return err
		}
		if p.isFunctionSetQuantifier() {
			p.advance()
		}
		if err := p.parseFunctionArgument(); err != nil {
			return err
		}
		for p.match(tokComma) {
			if err := p.parseFunctionArgument(); err != nil {
				return err
			}
		}
		if _, err := p.expect(tokRParen); err != nil {
			return err
		}
		return p.parseFunctionCallSuffixes()

	case kwSUBSTR, kwSUBSTRING:
		// Single-pass parsing for both SQL SUBSTRING(s FROM pos [FOR len])
		// and normal function calls SUBSTR(s, pos [, len]) without reparsing.
		p.advance()
		p.advance()
		if p.match(tokRParen) {
			return p.parseFunctionCallSuffixes()
		}
		hasQuantifier := p.isFunctionSetQuantifier()
		if hasQuantifier {
			p.advance()
		}
		isVal1, err := p.parseFunctionArgumentWithKind()
		if err != nil {
			return err
		}
		if !hasQuantifier && isVal1 && p.match(kwFROM) {
			if err := p.parseValueExpression(); err != nil {
				return err
			}
			if p.match(kwFOR, tokComma) {
				if err := p.parseValueExpression(); err != nil {
					return err
				}
			}
			_, err := p.expect(tokRParen)
			return err
		}
		if p.match(tokComma) {
			isVal2, err := p.parseFunctionArgumentWithKind()
			if err != nil {
				return err
			}
			if !hasQuantifier && isVal1 && isVal2 && p.match(kwFOR) {
				if err := p.parseValueExpression(); err != nil {
					return err
				}
				_, err := p.expect(tokRParen)
				return err
			}
			for p.match(tokComma) {
				if err := p.parseFunctionArgument(); err != nil {
					return err
				}
			}
		}
		if _, err := p.expect(tokRParen); err != nil {
			return err
		}
		return p.parseFunctionCallSuffixes()

	case kwTRIM:
		// Single-pass parsing for both SQL TRIM([BOTH|LEADING|TRAILING] [chars] FROM str)
		// and ordinary function calls TRIM(str [, chars]) to avoid exponential reparsing
		// on deeply nested TRIM(TRIM(...)) expressions and preserve columns named
		// `leading`, `trailing`, or `both`.
		p.advance()
		p.advance()
		if p.match(tokRParen) {
			return p.parseFunctionCallSuffixes()
		}
		if (p.cur().Kind == kwBOTH || p.cur().Kind == kwLEADING || p.cur().Kind == kwTRAILING) &&
			p.startsTrimOptionForm() {
			p.advance()
			if !p.startsTrimFromSeparator(true) {
				if err := p.parseValueExpression(); err != nil {
					return err
				}
			}
			if _, err := p.expect(kwFROM); err != nil {
				return err
			}
			if err := p.parseValueExpression(); err != nil {
				return err
			}
			_, err := p.expect(tokRParen)
			return err
		}
		if p.startsTrimFromSeparator(false) {
			p.advance()
			if err := p.parseValueExpression(); err != nil {
				return err
			}
			_, err := p.expect(tokRParen)
			return err
		}
		hasQuantifier := p.isFunctionSetQuantifier()
		if hasQuantifier {
			p.advance()
		}
		isVal, err := p.parseFunctionArgumentWithKind()
		if err != nil {
			return err
		}
		if !hasQuantifier && isVal && p.cur().Kind == kwFROM && p.peek(1).Kind != tokRParen {
			p.advance()
			if err := p.parseValueExpression(); err != nil {
				return err
			}
			_, err := p.expect(tokRParen)
			return err
		}
		for p.match(tokComma) {
			if err := p.parseFunctionArgument(); err != nil {
				return err
			}
		}
		if _, err := p.expect(tokRParen); err != nil {
			return err
		}
		return p.parseFunctionCallSuffixes()

	case kwOVERLAY:
		// Single-pass parsing for SQL OVERLAY(input PLACING replace FROM pos [FOR len])
		// and ordinary comma-separated OVERLAY(input, replace, pos [, len]) calls
		// without reparsing the first argument.
		p.advance()
		p.advance()
		if p.match(tokRParen) {
			return p.parseFunctionCallSuffixes()
		}
		hasQuantifier := p.isFunctionSetQuantifier() && p.peek(1).Kind != kwPLACING
		if hasQuantifier {
			p.advance()
		}
		isVal, err := p.parseFunctionArgumentWithKind()
		if err != nil {
			return err
		}
		if !hasQuantifier && isVal && p.match(kwPLACING) {
			if err := p.parseValueExpression(); err != nil {
				return err
			}
			if _, err := p.expect(kwFROM); err != nil {
				return err
			}
			if err := p.parseValueExpression(); err != nil {
				return err
			}
			if p.match(kwFOR) {
				if err := p.parseValueExpression(); err != nil {
					return err
				}
			}
			_, err := p.expect(tokRParen)
			return err
		}
		for p.match(tokComma) {
			if err := p.parseFunctionArgument(); err != nil {
				return err
			}
		}
		if _, err := p.expect(tokRParen); err != nil {
			return err
		}
		return p.parseFunctionCallSuffixes()

	case kwTIMESTAMPADD, kwDATEADD, kwDATE_ADD, kwTIMESTAMPDIFF, kwDATEDIFF, kwDATE_DIFF, kwTIMEDIFF:
		// Single-pass parsing for both 3-arg SQL TIMESTAMPADD/DATEDIFF(unit, a, b)
		// and 2-arg function calls DATE_ADD(d, n) / DATEDIFF(d1, d2) without reparsing.
		p.advance()
		p.advance()
		if p.match(tokRParen) {
			return p.parseFunctionCallSuffixes()
		}
		if p.isFunctionSetQuantifier() {
			p.advance()
		}
		if err := p.parseFunctionArgument(); err != nil {
			return err
		}
		for p.match(tokComma) {
			if err := p.parseFunctionArgument(); err != nil {
				return err
			}
		}
		if _, err := p.expect(tokRParen); err != nil {
			return err
		}
		return p.parseFunctionCallSuffixes()
	}
	return p.errorf("unexpected special keyword %q", p.cur().Text)
}

func (p *Parser) parseParenthesizedPrimary() error {
	cp := p.save()

	// 1. Multi-parameter lambda: (ident, ident+) -> expression
	p.advance() // '('
	if p.cur().Kind.isIdentifier() {
		p.advance()
		if p.match(tokComma) {
			ok := true
			for {
				if !p.cur().Kind.isIdentifier() {
					ok = false
					break
				}
				p.advance()
				if !p.match(tokComma) {
					break
				}
			}
			if ok && p.match(tokRParen) && p.match(tokArrow) {
				return p.parseExpression()
			}
		}
	}
	p.restore(cp)

	// 2. Subquery: (query) — once identified as a complete parenthesized query,
	// commit directly without restoring/retrying as an expression so nested
	// malformed queries like `SELECT ((SELECT 1 +))` fail in linear time.
	p.advance() // '('
	if p.startsWithCompleteParenQuery() {
		if err := p.parseQuery(); err != nil {
			return err
		}
		_, err := p.expect(tokRParen)
		return err
	}

	// 3. Single-pass parenthesized expression `(expr)` OR row constructor
	// `(namedExpr, namedExpr+)` without reparsing the first element.
	if err := p.parseExpression(); err != nil {
		return err
	}
	if p.match(tokRParen) {
		return nil
	}
	if err := p.parseOptionalNamedExprAlias(); err != nil {
		return err
	}
	if !p.match(tokComma) {
		return p.errorf("expected ')' or ',' in parenthesized expression, got %q", p.cur().Text)
	}
	for {
		if err := p.parseNamedExpression(); err != nil {
			return err
		}
		if !p.match(tokComma) {
			break
		}
	}
	_, err := p.expect(tokRParen)
	return err
}

func (p *Parser) isCaseColumnReference() bool {
	next := p.peek(1).Kind
	if next == kwWHEN {
		return false
	}
	if p.peek(2).Kind == kwWHEN {
		// `CASE <keyword> WHEN ...` uses the non-reserved keyword as the operand.
		return !next.isIdentifier()
	}
	return isCaseColumnNext(next)
}

func isCaseColumnNext(next TokenKind) bool {
	switch next {
	case tokEOF, tokSemicolon, tokComma, tokRParen, tokRBracket,
		tokDot, tokLBracket, tokDoubleColon, kwCOLLATE,
		tokAsterisk, tokSlash, tokPercent, kwDIV,
		tokConcatPipe, tokPipe, tokAmpersand, tokHat,
		tokShiftLeft, tokShiftRight, tokShiftRightU,
		tokEq, tokNeq, tokNeqBang, tokNullSafeEq, tokLT, tokLTE, tokGT, tokGTE,
		kwAND, kwOR, kwIS, kwIN, kwLIKE, kwILIKE, kwRLIKE, kwBETWEEN,
		tokArrow, tokFatArrow, kwAS,
		kwTHEN, kwELSE, kwEND:
		return true
	}
	return false
}

func (p *Parser) isCaseClauseBoundaryAt(i int) bool {
	k := p.peek(i).Kind
	switch k {
	case tokComma, tokOperatorPipe:
		return true
	case kwGROUP, kwORDER, kwCLUSTER, kwDISTRIBUTE, kwSORT:
		return p.peek(i+1).Kind == kwBY
	case kwWINDOW:
		return p.peek(i+1).Kind.isIdentifier() && p.peek(i+2).Kind == kwAS
	case kwUNION, kwINTERSECT, kwEXCEPT, kwMINUS:
		switch p.peek(i + 1).Kind {
		case kwALL, kwDISTINCT, kwSELECT, kwFROM, kwTABLE, kwVALUES, kwMAP, kwREDUCE, tokLParen:
			return true
		}
		return false
	case kwFROM, kwWHERE, kwHAVING, kwLIMIT, kwOFFSET:
		if i > 0 && p.peek(i-1).Kind == tokDot {
			return false
		}
		switch p.peek(i + 1).Kind {
		case kwWHEN, kwTHEN, kwELSE, kwEND,
			tokDot, tokLBracket, tokDoubleColon, kwCOLLATE,
			tokPlus, tokMinus, tokAsterisk, tokSlash, tokPercent, kwDIV,
			tokConcatPipe, tokPipe, tokAmpersand, tokHat,
			tokShiftLeft, tokShiftRight, tokShiftRightU,
			tokEq, tokNeq, tokNeqBang, tokNullSafeEq, tokLT, tokLTE, tokGT, tokGTE,
			kwAND, kwOR, kwIS, kwIN, kwLIKE, kwILIKE, kwRLIKE, kwBETWEEN,
			tokArrow, tokFatArrow:
			return false
		}
		return true
	}
	return false
}

func (p *Parser) openCasesAtCurrentParenDepth() int {
	count := 0
	for i := len(p.openCaseParenDepths) - 1; i >= 0; i-- {
		if p.openCaseParenDepths[i] == p.parenDepth {
			count++
		} else {
			break
		}
	}
	return count
}

// hasCaseWhenAndEndAhead scans forward from p.pos+1 at paren/bracket depth 0
// to decide whether a `CASE` token starts a simple `CASE <expr> WHEN ... END`
// expression rather than a column reference `case` or a `CASE(...)` function
// call. When inside an outer open CASE at the same paren depth, it requires
// more unmatched `END` tokens than currently open CASEs at this depth so
// `SELECT CASE case + case + 1 WHEN 1 THEN 'ok' END` keeps the inner `case`s as
// columns while `SELECT CASE CASE -1 WHEN -1 THEN 2 END WHEN 2 THEN 3 END` and
// `SELECT CASE abs(CASE -1 WHEN -1 THEN 2 END) WHEN 2 THEN 3 END` parse nested
// simple CASE expressions cleanly.
func (p *Parser) hasCaseWhenAndEndAhead() bool {
	open := p.openCasesAtCurrentParenDepth()
	if maxOpen, ok := p.caseScanConfirmedOpen[p.pos]; ok && open <= maxOpen {
		return true
	}
	if p.parenDepth >= 0 && p.parenDepth < len(p.caseScanCache) {
		entry := p.caseScanCache[p.parenDepth]
		if p.pos >= entry.start && p.pos < entry.end {
			if p.pos >= entry.lastWhenPos {
				return false
			}
			return open == 0 || entry.ends > open
		}
	}
	depth := 0
	sawWhen := false
	lastWhenPos := -1
	ends := 0
	firstEndPos := -1
	type siblingCase struct {
		pos     int
		hasWhen bool
	}
	var siblingStack []siblingCase
	var closedSiblings []int
	record := func(stopPos int) {
		if p.parenDepth < 0 {
			return
		}
		endPos := stopPos
		if firstEndPos >= 0 && firstEndPos < endPos {
			endPos = firstEndPos
		}
		if p.parenDepth >= len(p.caseScanCache) {
			grown := make([]caseScanEntry, p.parenDepth+1)
			copy(grown, p.caseScanCache)
			p.caseScanCache = grown
		}
		p.caseScanCache[p.parenDepth] = caseScanEntry{
			start:       p.pos,
			end:         endPos,
			lastWhenPos: lastWhenPos,
			ends:        ends,
		}
	}
	confirmSiblings := func() {
		if len(closedSiblings) == 0 {
			return
		}
		if p.caseScanConfirmedOpen == nil {
			p.caseScanConfirmedOpen = make(map[int]int, len(closedSiblings))
		}
		for _, spos := range closedSiblings {
			if prev, ok := p.caseScanConfirmedOpen[spos]; !ok || open > prev {
				p.caseScanConfirmedOpen[spos] = open
			}
		}
	}
	for i := 1; ; i++ {
		k := p.peek(i).Kind
		switch k {
		case tokEOF, tokSemicolon:
			record(p.pos + i)
			return sawWhen && (open == 0 || ends > open)
		case tokLParen, tokLBracket:
			depth++
		case tokRParen, tokRBracket:
			if depth == 0 {
				record(p.pos + i)
				return sawWhen && (open == 0 || ends > open)
			}
			depth--
		case kwWHEN:
			if depth == 0 {
				if open == 0 {
					return true
				}
				sawWhen = true
				lastWhenPos = p.pos + i
				if len(siblingStack) > 0 {
					siblingStack[len(siblingStack)-1].hasWhen = true
				}
			}
		case kwCASE:
			// After an END has closed the first CASE, a subsequent CASE at depth 0
			// may start a sibling CASE whose own END should not count toward closing
			// an enclosing open CASE.
			if depth == 0 && ends > 0 {
				next := p.peek(i + 1).Kind
				if next == kwWHEN || !isCaseColumnNext(next) {
					siblingStack = append(siblingStack, siblingCase{pos: p.pos + i})
				}
			}
		case kwEND:
			if depth == 0 {
				next := p.peek(i + 1).Kind
				if !(i == 1 && next == kwWHEN) &&
					next != kwIF && next != kwCASE && next != kwLOOP && next != kwWHILE && next != kwREPEAT && next != kwFOR {
					// Discard any speculatively pushed sibling CASEs that never saw a WHEN
					// (they were column references like `case + 1`, not CASE expressions).
					for len(siblingStack) > 0 && !siblingStack[len(siblingStack)-1].hasWhen {
						siblingStack = siblingStack[:len(siblingStack)-1]
					}
					if len(siblingStack) > 0 {
						top := siblingStack[len(siblingStack)-1]
						siblingStack = siblingStack[:len(siblingStack)-1]
						if len(siblingStack) == 0 {
							closedSiblings = append(closedSiblings, top.pos)
						}
					} else {
						if firstEndPos < 0 {
							firstEndPos = p.pos + i
						}
						ends++
						if sawWhen && ends > open {
							confirmSiblings()
							return true
						}
					}
				}
			}
		default:
			if depth == 0 && p.isCaseClauseBoundaryAt(i) {
				record(p.pos + i)
				return sawWhen && (open == 0 || ends > open)
			}
		}
	}
}

// parseCaseOrCaseFuncCall parses `CASE(...)` (either a function call or a simple
// `CASE (expr) WHEN ... END` expression) in a single pass without reparsing the
// parenthesized argument.
func (p *Parser) parseCaseOrCaseFuncCall() error {
	if _, err := p.expect(kwCASE); err != nil {
		return err
	}
	if _, err := p.expect(tokLParen); err != nil {
		return err
	}
	if p.match(tokRParen) {
		return p.parseFunctionCallSuffixes()
	}
	if p.isFunctionSetQuantifier() {
		p.advance()
		if err := p.parseFunctionArgument(); err != nil {
			return err
		}
		for p.match(tokComma) {
			if err := p.parseFunctionArgument(); err != nil {
				return err
			}
		}
		if _, err := p.expect(tokRParen); err != nil {
			return err
		}
		return p.parseFunctionCallSuffixes()
	}
	isNamed := p.cur().Kind.isIdentifier() && p.peek(1).Kind == tokFatArrow
	if isNamed {
		if err := p.parseFunctionArgument(); err != nil {
			return err
		}
		for p.match(tokComma) {
			if err := p.parseFunctionArgument(); err != nil {
				return err
			}
		}
		if _, err := p.expect(tokRParen); err != nil {
			return err
		}
		return p.parseFunctionCallSuffixes()
	}
	if err := p.parseExpression(); err != nil {
		return err
	}
	if p.match(tokRParen) {
		if p.cur().Kind == kwWHEN {
			return p.parseCaseWhenTail()
		}
		return p.parseFunctionCallSuffixes()
	}
	// Multi-argument CASE(a, b, ...) function call or row constructor CASE (a, b) WHEN ... END.
	for p.match(tokComma) {
		if err := p.parseFunctionArgument(); err != nil {
			return err
		}
	}
	if _, err := p.expect(tokRParen); err != nil {
		return err
	}
	if p.cur().Kind == kwWHEN {
		return p.parseCaseWhenTail()
	}
	return p.parseFunctionCallSuffixes()
}

func (p *Parser) parseCaseExpression() error {
	if _, err := p.expect(kwCASE); err != nil {
		return err
	}
	p.openCaseParenDepths = append(p.openCaseParenDepths, p.parenDepth)
	defer func() {
		p.openCaseParenDepths = p.openCaseParenDepths[:len(p.openCaseParenDepths)-1]
	}()
	if p.cur().Kind != kwWHEN {
		if err := p.parseExpression(); err != nil {
			return err
		}
	}
	return p.parseCaseWhenTail()
}

func (p *Parser) parseCaseWhenTail() error {
	whenCount := 0
	for p.match(kwWHEN) {
		if err := p.parseExpression(); err != nil {
			return err
		}
		if _, err := p.expect(kwTHEN); err != nil {
			return err
		}
		if err := p.parseExpression(); err != nil {
			return err
		}
		whenCount++
	}
	if whenCount == 0 {
		return p.errorf("expected WHEN in CASE expression")
	}
	if p.match(kwELSE) {
		if err := p.parseExpression(); err != nil {
			return err
		}
	}
	_, err := p.expect(kwEND)
	return err
}

func (p *Parser) canStartFunctionCall() bool {
	if p.cur().Kind == kwIDENTIFIER && p.peek(1).Kind == tokLParen {
		return true
	}
	if !p.cur().Kind.isIdentifier() {
		return false
	}
	// Look ahead past optional .ident segments to see if '(' follows.
	i := 1
	for p.peek(i).Kind == tokDot && p.peek(i+1).Kind.isIdentifier() {
		i += 2
	}
	return p.peek(i).Kind == tokLParen
}

// startsTrimOptionForm reports whether a leading BOTH/LEADING/TRAILING token
// inside TRIM(...) introduces the SQL trim option syntax (e.g.
// `TRIM(BOTH FROM x)` or `TRIM(LEADING '0' FROM x)`) rather than an ordinary
// column/expression argument (e.g. `TRIM(leading)` or `TRIM(leading, 'x')` or
// `TRIM(leading || 'a')`).
func (p *Parser) startsTrimOptionForm() bool {
	next := p.peek(1).Kind
	if next == kwFROM {
		return true
	}
	switch next {
	case tokRParen, tokComma, tokEOF, tokSemicolon,
		tokDot, tokLParen, tokLBracket, tokDoubleColon, kwCOLLATE,
		tokPlus, tokMinus, tokAsterisk, tokSlash, tokPercent, kwDIV,
		tokConcatPipe, tokPipe, tokAmpersand, tokHat,
		tokShiftLeft, tokShiftRight, tokShiftRightU,
		tokEq, tokNeq, tokNeqBang, tokNullSafeEq, tokLT, tokLTE, tokGT, tokGTE,
		kwAND, kwOR, kwIS, kwIN, kwLIKE, kwILIKE, kwRLIKE, kwBETWEEN, kwNOT,
		tokArrow, tokFatArrow:
		return false
	}
	return true
}

// startsTrimFromSeparator reports whether `FROM` at p.cur() inside `TRIM(...)`
// acts as the SQL-specific `FROM srcStr` separator (with omitted trimStr)
// rather than the start of an argument expression beginning with a column or
// function named `from` (e.g. `TRIM(from, s)`, `TRIM(from || x)`,
// `TRIM(from.x)`, `TRIM(from FROM s)`). When hasOption is true, a trim option
// (BOTH/LEADING/TRAILING) was already consumed so parenthesized or unary
// source expressions like `TRIM(BOTH FROM (s))` and `TRIM(BOTH FROM -s)` treat
// `FROM` as the separator unless another depth-0 `FROM` appears before ')'.
func (p *Parser) startsTrimFromSeparator(hasOption bool) bool {
	if p.cur().Kind != kwFROM {
		return false
	}
	next := p.peek(1).Kind
	if next == kwFROM {
		after := p.peek(2).Kind
		switch after {
		case tokRParen, tokDot, tokLParen, tokLBracket, tokDoubleColon, kwCOLLATE,
			tokPlus, tokMinus, tokAsterisk, tokSlash, tokPercent, kwDIV,
			tokConcatPipe, tokPipe, tokAmpersand, tokHat,
			tokShiftLeft, tokShiftRight, tokShiftRightU:
			return true
		}
		return false
	}
	if hasOption {
		switch next {
		case tokRParen, tokEOF, tokSemicolon:
			return false
		case tokLParen, tokPlus, tokMinus, tokDot, tokLBracket, tokDoubleColon, kwCOLLATE,
			tokAsterisk, tokSlash, tokPercent, kwDIV,
			tokConcatPipe, tokPipe, tokAmpersand, tokHat,
			tokShiftLeft, tokShiftRight, tokShiftRightU:
			return !p.hasKeywordAtDepthZero(1, kwFROM)
		}
		return true
	}
	if next == tokLParen {
		return p.startsWithQueryKeywordAt(2)
	}
	switch next {
	case tokRParen, tokComma, tokEOF, tokSemicolon,
		tokDot, tokLBracket, tokDoubleColon, kwCOLLATE,
		tokPlus, tokMinus, tokAsterisk, tokSlash, tokPercent, kwDIV,
		tokConcatPipe, tokPipe, tokAmpersand, tokHat,
		tokShiftLeft, tokShiftRight, tokShiftRightU,
		tokEq, tokNeq, tokNeqBang, tokNullSafeEq, tokLT, tokLTE, tokGT, tokGTE,
		kwAND, kwOR, kwIS, kwIN, kwLIKE, kwILIKE, kwRLIKE, kwBETWEEN, kwNOT,
		tokArrow, tokFatArrow:
		return false
	}
	return true
}

func (p *Parser) parseFunctionCall() error {
	// Single-pass handling for IDENTIFIER(expr) vs IDENTIFIER(expr)(args) so
	// deeply nested IDENTIFIER(IDENTIFIER(...)) never reparses its argument.
	if p.cur().Kind == kwIDENTIFIER && p.peek(1).Kind == tokLParen {
		p.advance()
		p.advance()
		singlePlainArg := false
		if !p.match(tokRParen) {
			if p.isFunctionSetQuantifier() {
				p.advance()
				if err := p.parseFunctionArgument(); err != nil {
					return err
				}
				for p.match(tokComma) {
					if err := p.parseFunctionArgument(); err != nil {
						return err
					}
				}
			} else {
				isNamed := p.cur().Kind.isIdentifier() && p.peek(1).Kind == tokFatArrow
				if err := p.parseFunctionArgument(); err != nil {
					return err
				}
				argCount := 1
				for p.match(tokComma) {
					if err := p.parseFunctionArgument(); err != nil {
						return err
					}
					argCount++
				}
				singlePlainArg = !isNamed && argCount == 1
			}
			if _, err := p.expect(tokRParen); err != nil {
				return err
			}
		}
		if singlePlainArg && p.cur().Kind == tokLParen {
			// It was `IDENTIFIER(expr)` as functionName; now parse the actual call argument list.
			if err := p.parseFunctionCallArgList(); err != nil {
				return err
			}
		}
		return p.parseFunctionCallSuffixes()
	}

	if _, err := p.parseQualifiedName(); err != nil {
		return err
	}
	if err := p.parseFunctionCallArgList(); err != nil {
		return err
	}
	return p.parseFunctionCallSuffixes()
}

func (p *Parser) parseFunctionCallArgList() error {
	if _, err := p.expect(tokLParen); err != nil {
		return err
	}
	if p.match(tokRParen) {
		return nil
	}
	if p.isFunctionSetQuantifier() {
		p.advance()
	}
	if err := p.parseFunctionArgument(); err != nil {
		return err
	}
	for p.match(tokComma) {
		if err := p.parseFunctionArgument(); err != nil {
			return err
		}
	}
	_, err := p.expect(tokRParen)
	return err
}

// isFunctionSetQuantifier reports whether a leading DISTINCT/ALL inside a
// function call argument list acts as a set quantifier (e.g. COUNT(DISTINCT x))
// rather than an ordinary column/expression named `distinct` or `all`.
func (p *Parser) isFunctionSetQuantifier() bool {
	return p.isFunctionSetQuantifierAt(0)
}

func (p *Parser) isFunctionSetQuantifierAt(offset int) bool {
	k := p.peek(offset).Kind
	if k != kwDISTINCT && k != kwALL {
		return false
	}
	next := p.peek(offset + 1).Kind
	switch next {
	case tokAsterisk:
		// Distinguish a standalone star argument (`count(ALL *)`,
		// `f(DISTINCT *, x)`) from multiplication on a column named `all` or
		// `distinct` (`f(all * 2)`).
		afterStar := p.peek(offset + 2).Kind
		return afterStar == tokRParen || afterStar == tokComma
	case tokRParen, tokComma, tokEOF, tokSemicolon,
		tokDot, tokLBracket, tokDoubleColon, kwCOLLATE,
		tokSlash, tokPercent, kwDIV,
		tokConcatPipe, tokPipe, tokAmpersand, tokHat,
		tokShiftLeft, tokShiftRight, tokShiftRightU,
		tokEq, tokNeq, tokNeqBang, tokNullSafeEq, tokLT, tokLTE, tokGT, tokGTE,
		kwAND, kwOR, kwIS, kwIN, kwLIKE, kwILIKE, kwRLIKE, kwBETWEEN,
		kwFROM, kwAS,
		tokArrow, tokFatArrow:
		return false
	}
	return true
}

func (p *Parser) parseFunctionCallSuffixes() error {
	// Optional WITHIN GROUP (ORDER BY sortItem (COMMA sortItem)*)
	if p.cur().Kind == kwWITHIN && p.peek(1).Kind == kwGROUP {
		p.advance()
		p.advance()
		if _, err := p.expect(tokLParen); err != nil {
			return err
		}
		if _, err := p.expect(kwORDER); err != nil {
			return err
		}
		if _, err := p.expect(kwBY); err != nil {
			return err
		}
		if err := p.parseSortItemList(); err != nil {
			return err
		}
		if _, err := p.expect(tokRParen); err != nil {
			return err
		}
	}
	// Optional FILTER (WHERE booleanExpression)
	if p.cur().Kind == kwFILTER && p.peek(1).Kind == tokLParen {
		p.advance()
		p.advance()
		if _, err := p.expect(kwWHERE); err != nil {
			return err
		}
		if err := p.parseBooleanExpression(); err != nil {
			return err
		}
		if _, err := p.expect(tokRParen); err != nil {
			return err
		}
	}
	// Optional (IGNORE | RESPECT) NULLS
	if (p.cur().Kind == kwIGNORE || p.cur().Kind == kwRESPECT) && p.peek(1).Kind == kwNULLS {
		p.advance()
		p.advance()
	}
	// Optional OVER windowSpec (skip when OVER is immediately followed by a
	// clause boundary or delimiter, e.g. `SELECT count(*) over FROM t`, while
	// still allowing non-reserved keywords like `window` or `order` as named
	// window references such as `SELECT count(*) OVER window FROM t WINDOW window AS ()`).
	if p.cur().Kind == kwOVER {
		next := p.peek(1).Kind
		canStartWindow := next == tokLParen || (next.isIdentifier() && !p.isClauseAfterOverAlias())
		if canStartWindow {
			cpOver := p.save()
			p.advance()
			if err := p.parseWindowSpec(); err != nil {
				if p.isStatementEnd() || p.cur().Kind == tokComma || p.cur().Kind == tokRParen {
					p.restore(cpOver)
					return nil
				}
				return err
			}
		}
	}
	return nil
}

// isClauseAfterOverAlias reports whether the tokens starting at p.peek(1)
// (immediately after `OVER`) introduce a following SQL clause — meaning `OVER`
// serves as an implicit column alias (e.g. `SELECT count(*) over FROM t` or
// `SELECT count(*) over ORDER BY 1`) rather than `OVER windowRef` (e.g.
// `SELECT count(*) OVER window FROM t WINDOW window AS ()`).
func (p *Parser) isClauseAfterOverAlias() bool {
	next := p.peek(1).Kind
	after := p.peek(2).Kind
	isBoundaryOrClause := func(k TokenKind, tok Token) bool {
		switch k {
		case tokEOF, tokSemicolon, tokComma, tokRParen, tokRBracket, tokOperatorPipe,
			kwAS, kwFROM, kwWHERE, kwGROUP, kwHAVING, kwWINDOW, kwORDER,
			kwCLUSTER, kwDISTRIBUTE, kwSORT, kwLIMIT, kwOFFSET,
			kwUNION, kwINTERSECT, kwEXCEPT, kwMINUS, kwLATERAL,
			kwTHEN, kwELSE, kwELSEIF, kwEND, kwWHEN:
			return true
		case tokIdent:
			return isWord(tok, "QUALIFY")
		}
		return false
	}
	isBinaryOrPostfixOp := func(k TokenKind) bool {
		switch k {
		case tokDot, tokLBracket, tokDoubleColon, kwCOLLATE,
			tokAsterisk, tokSlash, tokPercent, kwDIV,
			tokConcatPipe, tokPipe, tokAmpersand, tokHat,
			tokShiftLeft, tokShiftRight, tokShiftRightU,
			tokEq, tokNeq, tokNeqBang, tokNullSafeEq, tokLT, tokLTE, tokGT, tokGTE,
			kwAND, kwOR, kwIS, kwIN, kwLIKE, kwILIKE, kwRLIKE, kwBETWEEN:
			return true
		}
		return false
	}
	switch next {
	case kwORDER, kwGROUP, kwCLUSTER, kwDISTRIBUTE, kwSORT:
		return after == kwBY
	case kwLATERAL:
		return after == kwVIEW
	case kwWINDOW:
		return after.isIdentifier() && p.peek(3).Kind == kwAS
	case kwUNION, kwINTERSECT, kwEXCEPT, kwMINUS:
		switch after {
		case kwSELECT, kwTABLE, kwVALUES, kwWITH, tokLParen, kwALL, kwDISTINCT:
			return true
		case kwBY:
			return p.peek(3).Kind == kwNAME
		}
		return false
	case kwFROM:
		if isBoundaryOrClause(after, p.peek(2)) || isBinaryOrPostfixOp(after) ||
			after == tokPlus || after == tokMinus || after == tokTilde || after == tokBang {
			return false
		}
		if after.isIdentifier() && p.peek(3).Kind == kwFROM {
			return false
		}
		return true
	case kwWHERE, kwHAVING, kwLIMIT, kwOFFSET:
		if isBoundaryOrClause(after, p.peek(2)) || isBinaryOrPostfixOp(after) {
			return false
		}
		if after.isIdentifier() && p.peek(3).Kind == kwFROM {
			return false
		}
		return true
	case tokIdent:
		if !isWord(p.peek(1), "QUALIFY") {
			return false
		}
		if isBoundaryOrClause(after, p.peek(2)) || isBinaryOrPostfixOp(after) {
			return false
		}
		if after.isIdentifier() && p.peek(3).Kind == kwFROM {
			return false
		}
		return true
	}
	return false
}

func (p *Parser) isStringLitAt(offset int) bool {
	k := p.peek(offset).Kind
	return k == tokString || k == tokDoubleQuoted
}

func (p *Parser) canStartConstant() bool {
	switch p.cur().Kind {
	case kwNULL, tokQuestion, tokColon, kwINTERVAL,
		tokInt, tokDecimal, tokExponent, tokBigInt, tokSmallInt, tokTinyInt,
		tokFloat, tokDouble, tokBigDecimal, tokMinus,
		kwTRUE, kwFALSE, tokString, tokDoubleQuoted:
		return true
	}
	// literalType stringLit: DATE '...', TIMESTAMP '...', X '...', ident '...'
	if p.cur().Kind.isIdentifier() && (p.peek(1).Kind == tokString || p.peek(1).Kind == tokDoubleQuoted) {
		return true
	}
	return false
}

func (p *Parser) parseConstant() error {
	switch p.cur().Kind {
	case kwNULL, tokQuestion, kwTRUE, kwFALSE:
		p.advance()
		return nil
	case tokColon:
		p.advance()
		_, err := p.parseIdentifier()
		return err
	case kwINTERVAL:
		// Could be `INTERVAL 's'` (typeConstructor) or `INTERVAL 1 DAY` (interval literal).
		cp := p.save()
		if err := p.parseInterval(); err == nil {
			return nil
		}
		p.restore(cp)
		if p.peek(1).Kind == tokString || p.peek(1).Kind == tokDoubleQuoted {
			p.advance()
			p.advance()
			return nil
		}
		return p.errorf("invalid INTERVAL literal")
	case tokString, tokDoubleQuoted:
		p.advance()
		for p.cur().Kind == tokString || p.cur().Kind == tokDoubleQuoted {
			p.advance()
		}
		return nil
	case tokMinus, tokInt, tokDecimal, tokExponent, tokBigInt, tokSmallInt, tokTinyInt,
		tokFloat, tokDouble, tokBigDecimal:
		return p.parseNumber()
	}

	// literalType stringLit (DATE '...', TIMESTAMP '...', TIMESTAMP_LTZ '...', TIMESTAMP_NTZ '...', X '...', ident '...')
	if p.cur().Kind.isIdentifier() && (p.peek(1).Kind == tokString || p.peek(1).Kind == tokDoubleQuoted) {
		p.advance()
		p.advance()
		return nil
	}
	return p.errorf("expected constant, got %q", p.cur().Text)
}

func (p *Parser) parseNumber() error {
	cp := p.save()
	p.match(tokMinus)
	switch p.cur().Kind {
	case tokInt, tokDecimal, tokExponent, tokBigInt, tokSmallInt, tokTinyInt,
		tokFloat, tokDouble, tokBigDecimal:
		p.advance()
		return nil
	}
	p.restore(cp)
	return p.errorf("expected numeric literal, got %q", p.cur().Text)
}

// parseInterval parses:
//
//	interval : INTERVAL (errorCapturingMultiUnitsInterval | errorCapturingUnitToUnitInterval) ;
func (p *Parser) parseInterval() error {
	if _, err := p.expect(kwINTERVAL); err != nil {
		return err
	}
	// Try unitToUnitInterval first if we see: intervalValue unitInUnitToUnit TO unitInUnitToUnit
	cp := p.save()
	if err := p.parseUnitToUnitInterval(); err == nil {
		// errorCapturingUnitToUnitInterval allows trailing error1=multiUnitsInterval | error2=unitToUnitInterval
		cpExtra := p.save()
		if err := p.parseUnitToUnitInterval(); err != nil {
			p.restore(cpExtra)
			if err := p.parseMultiUnitsInterval(); err != nil {
				p.restore(cpExtra)
			}
		}
		return nil
	}
	p.restore(cp)

	if err := p.parseMultiUnitsInterval(); err != nil {
		return err
	}
	// errorCapturingMultiUnitsInterval allows optional trailing unitToUnitInterval
	cpTrail := p.save()
	if err := p.parseUnitToUnitInterval(); err != nil {
		p.restore(cpTrail)
	}
	return nil
}

func (p *Parser) parseMultiUnitsInterval() error {
	count := 0
	for {
		cp := p.save()
		if err := p.parseIntervalValue(); err != nil {
			p.restore(cp)
			break
		}
		if !p.isUnitInMultiUnits() {
			p.restore(cp)
			break
		}
		p.advance()
		// Don't steal the start of a unitToUnitInterval (value unit TO unit) if we already have >=1 unit.
		if count > 0 && p.cur().Kind == kwTO && p.isUnitInUnitToUnitAt(1) {
			p.restore(cp)
			break
		}
		count++
	}
	if count == 0 {
		return p.errorf("expected interval value and unit")
	}
	return nil
}

func (p *Parser) isUnitInUnitToUnitAt(offset int) bool {
	switch p.peek(offset).Kind {
	case kwSECOND, kwMINUTE, kwHOUR, kwDAY, kwMONTH, kwYEAR:
		return true
	}
	return false
}

func (p *Parser) parseUnitToUnitInterval() error {
	if err := p.parseIntervalValue(); err != nil {
		return err
	}
	if !p.isUnitInUnitToUnit() {
		return p.errorf("expected interval unit")
	}
	p.advance()
	if _, err := p.expect(kwTO); err != nil {
		return err
	}
	if !p.isUnitInUnitToUnit() {
		return p.errorf("expected interval unit after TO")
	}
	p.advance()
	return nil
}

func (p *Parser) parseIntervalValue() error {
	p.match(tokPlus, tokMinus)
	switch p.cur().Kind {
	case tokInt, tokDecimal, tokString, tokDoubleQuoted:
		p.advance()
		return nil
	default:
		return p.errorf("expected interval value")
	}
}

func (p *Parser) isUnitInMultiUnits() bool {
	switch p.cur().Kind {
	case kwNANOSECOND, kwNANOSECONDS, kwMICROSECOND, kwMICROSECONDS,
		kwMILLISECOND, kwMILLISECONDS, kwSECOND, kwSECONDS,
		kwMINUTE, kwMINUTES, kwHOUR, kwHOURS,
		kwDAY, kwDAYS, kwWEEK, kwWEEKS,
		kwMONTH, kwMONTHS, kwYEAR, kwYEARS:
		return true
	}
	return false
}

func (p *Parser) isUnitInUnitToUnit() bool {
	switch p.cur().Kind {
	case kwSECOND, kwMINUTE, kwHOUR, kwDAY, kwMONTH, kwYEAR:
		return true
	}
	return false
}

func (p *Parser) isDatetimeUnit() bool {
	switch p.cur().Kind {
	case kwYEAR, kwQUARTER, kwMONTH, kwWEEK, kwDAY, kwDAYOFYEAR,
		kwHOUR, kwMINUTE, kwSECOND, kwMILLISECOND, kwMICROSECOND:
		return true
	}
	return false
}
