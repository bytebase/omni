package parser

import (
	"strconv"
	"strings"
)

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
	_, err := p.parseNamedExpressionWithAliasFlag()
	return err
}

func (p *Parser) parseNamedExpressionWithAliasFlag() (bool, error) {
	if err := p.parseExpression(); err != nil {
		return false, err
	}
	return p.parseOptionalNamedExprAliasFlag()
}

func (p *Parser) parseOptionalNamedExprAlias() error {
	_, err := p.parseOptionalNamedExprAliasFlag()
	return err
}

func (p *Parser) parseOptionalNamedExprAliasFlag() (bool, error) {
	if p.cur().Kind == kwAS {
		next := p.peek(1).Kind
		// In Spark default mode, `AS` is non-reserved, so `SELECT 1 AS`,
		// `SELECT a AS FROM t`, and `SELECT a AS, b FROM t` use `AS` itself as
		// an implicit alias when a clause or delimiter immediately follows.
		switch next {
		case tokEOF, tokSemicolon, tokComma, tokRParen, tokOperatorPipe,
			kwWHERE, kwGROUP, kwHAVING, kwWINDOW, kwORDER,
			kwCLUSTER, kwDISTRIBUTE, kwSORT, kwLIMIT, kwOFFSET,
			kwUNION, kwINTERSECT, kwEXCEPT, kwMINUS:
			p.advance()
			return true, nil
		case kwFROM:
			if p.isFromQueryAt(1) {
				p.advance()
				return true, nil
			}
		}
		p.advance()
		if p.cur().Kind == tokLParen {
			return true, p.parseIdentifierList()
		}
		_, err := p.parseErrorCapturingIdentifier()
		return true, err
	}
	if p.cur().Kind == tokLParen {
		cp := p.save()
		if err := p.parseIdentifierList(); err == nil {
			return true, nil
		}
		p.restore(cp)
		return false, nil
	}
	if p.canBeImplicitExprAlias() {
		_, err := p.parseErrorCapturingIdentifier()
		return true, err
	}
	return false, nil
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
	// If the next token ends the current expression list / statement / paren / pipe,
	// a non-reserved keyword can safely be an implicit alias (e.g. SELECT a FROM,
	// SELECT a ALL, SELECT x asc, SELECT x rows, SELECT x limit |> SELECT *).
	next := p.peek(1).Kind
	if next == tokEOF || next == tokSemicolon || next == tokComma || next == tokRParen || next == tokOperatorPipe {
		switch k {
		case kwEND:
			return p.openCasesAtCurrentParenDepth() == 0
		case kwELSEIF, kwTHEN, kwDO:
			return false
		}
		return true
	}
	// Never treat CASE/control-flow terminators as implicit aliases.
	switch k {
	case kwEND, kwELSE, kwELSEIF, kwTHEN, kwWHEN, kwUNTIL, kwDO, kwAS:
		return false
	}
	// Clauses that require a specific second keyword can serve as implicit aliases
	// whenever that second keyword does not follow.
	switch k {
	case kwGROUP, kwORDER, kwCLUSTER, kwDISTRIBUTE, kwSORT:
		return next != kwBY
	case kwGROUPING:
		return !(next == kwSETS && p.peek(2).Kind == tokLParen)
	case kwWINDOW:
		return !(next.isIdentifier() && p.peek(2).Kind == kwAS)
	case kwWITH:
		return next != kwROLLUP && next != kwCUBE
	case kwLATERAL:
		return next != kwVIEW
	case kwLIMIT, kwOFFSET, kwWHERE, kwHAVING, kwFROM:
		if next == kwFROM {
			return p.isFromQueryAt(1)
		}
		if (next == kwORDER || next == kwGROUP || next == kwCLUSTER || next == kwDISTRIBUTE || next == kwSORT) && p.peek(2).Kind == kwBY {
			return true
		}
		if next == kwWINDOW && p.peek(2).Kind.isIdentifier() && p.peek(3).Kind == kwAS {
			return true
		}
		if next == kwWITH && (p.peek(2).Kind == kwROLLUP || p.peek(2).Kind == kwCUBE) {
			return true
		}
		if next == kwGROUPING && p.peek(2).Kind == kwSETS && p.peek(3).Kind == tokLParen {
			return true
		}
		if next == kwLATERAL && p.peek(2).Kind == kwVIEW {
			return true
		}
		if next == kwLIMIT || next == kwOFFSET || next == kwWHERE || next == kwHAVING {
			after := p.peek(2).Kind
			if !p.isClauseStartAt(2) && !isExprContinuationOp(after) {
				return true
			}
		}
		if next == tokIdent && isWord(p.peek(1), "QUALIFY") && p.peek(2).Kind != tokEOF && p.peek(2).Kind != tokSemicolon && p.peek(2).Kind != tokRParen {
			return true
		}
		return false
	}
	// Allow a non-reserved keyword as implicit column alias right before a FROM
	// clause (e.g. `SELECT a IS FROM a`, `SELECT x pivot FROM t`, `SELECT x join FROM t`),
	// except set-operator keywords which can precede a `FROM ...` query term.
	if next == kwFROM && p.isFromQueryAt(1) {
		switch k {
		case kwUNION, kwINTERSECT, kwEXCEPT, kwMINUS:
			return false
		}
		return true
	}
	switch k {
	case kwOVER, kwFILTER, kwWITHIN, kwIGNORE, kwRESPECT,
		kwASC, kwDESC, kwNULLS, kwFIRST, kwLAST, kwROWS, kwRANGE, kwPRECEDING, kwFOLLOWING:
		switch next {
		case kwWHERE, kwGROUP, kwHAVING, kwWINDOW,
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
	case kwUNION, kwINTERSECT, kwEXCEPT, kwMINUS,
		kwPIVOT, kwUNPIVOT, kwJOIN, kwINNER, kwLEFT, kwRIGHT, kwFULL,
		kwCROSS, kwNATURAL, kwSEMI, kwANTI, kwON, kwUSING, kwINTO, kwSET, kwVALUES,
		kwSELECT, kwINSERT, kwFOR,
		kwBY, kwIN, kwIS, kwAND, kwOR, kwNOT, kwBETWEEN, kwLIKE,
		kwILIKE, kwRLIKE, kwESCAPE,
		kwRECORDWRITER, kwRECORDREADER, kwROW, kwPLACING, kwRETURNS, kwRETURN:
		return false
	}
	return true
}

func isExprContinuationOp(k TokenKind) bool {
	switch k {
	case tokDot, tokLBracket, tokDoubleColon, kwCOLLATE,
		tokPlus, tokMinus, tokAsterisk, tokSlash, tokPercent, kwDIV,
		tokConcatPipe, tokPipe, tokAmpersand, tokHat,
		tokShiftLeft, tokShiftRight, tokShiftRightU,
		tokEq, tokNeq, tokNeqBang, tokNullSafeEq, tokLT, tokLTE, tokGT, tokGTE,
		kwAND, kwOR, kwIS, kwIN, kwLIKE, kwILIKE, kwRLIKE, kwBETWEEN,
		tokArrow, tokFatArrow:
		return true
	}
	return false
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
	case tokLParen:
		return p.peek(2).Kind == tokRParen
	case kwLATERAL:
		return p.peek(2).Kind == kwVIEW
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
	case kwNULLS:
		after := p.peek(2).Kind
		return after == kwFIRST || after == kwLAST
	case kwASC, kwDESC:
		after := p.peek(2).Kind
		return after == tokEOF || after == tokSemicolon || after == tokComma || after == tokRParen ||
			after == tokOperatorPipe || after == kwNULLS || after == kwROWS || after == kwRANGE ||
			after == kwLIMIT || after == kwOFFSET
	case kwROWS, kwRANGE:
		after := p.peek(2).Kind
		if after == kwBETWEEN || after == kwUNBOUNDED || (after == kwCURRENT && p.peek(3).Kind == kwROW) {
			return true
		}
		return p.hasKeywordAtDepthZero(2, kwPRECEDING) || p.hasKeywordAtDepthZero(2, kwFOLLOWING)
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
	case kwFROM, kwWHERE, kwHAVING, kwLIMIT, kwOFFSET, kwTHEN, kwELSE, kwWHEN, kwON, kwINTO, kwFOR:
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
	// IS [NOT] (NULL | TRUE | FALSE | UNKNOWN | DISTINCT FROM valueExpression)
	// Spark v4.0.5 discontinues `!` in place of NOT here (SYNTAX_DISCONTINUED.BANG_EQUALS_NOT).
	if p.cur().Kind == kwIS {
		cp := p.save()
		p.advance()
		p.match(kwNOT)
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

	// [NOT] (BETWEEN | IN | RLIKE | LIKE | ILIKE)
	// Spark v4.0.5 discontinues `!` in place of NOT here (SYNTAX_DISCONTINUED.BANG_EQUALS_NOT).
	cp := p.save()
	hasNot := p.match(kwNOT)
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
				if p.cur().Kind == tokRParen {
					return p.errorf("expected at least one pattern inside LIKE ANY/SOME/ALL (...)")
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
	negated := false
	for p.cur().Kind == tokMinus || p.cur().Kind == tokPlus || p.cur().Kind == tokTilde {
		k := p.advance().Kind
		ops++
		if p.depth+ops > maxNestingDepth {
			return p.errorf("maximum nesting depth exceeded")
		}
		negated = k == tokMinus && ops == 1
	}
	return p.parsePrimaryExpressionNeg(negated)
}

// parsePrimaryExpression parses Spark's primaryExpression + postfix operators
// (COLLATE, :: dataType, [index], .field).
func (p *Parser) parsePrimaryExpression() error {
	return p.parsePrimaryExpressionNeg(false)
}

func (p *Parser) parsePrimaryExpressionNeg(negated bool) error {
	if err := p.enterNesting(); err != nil {
		return err
	}
	defer p.leaveNesting()

	isQualifiedIdent, err := p.parsePrimaryAtom(negated)
	if err != nil {
		return err
	}

	// Postfix loop: COLLATE, :: dataType, [index], .field (or .* [EXCEPT (...)] on qualified identifiers).
	for {
		switch p.cur().Kind {
		case kwCOLLATE:
			if !p.peek(1).Kind.isIdentifier() {
				return nil
			}
			if err := p.parseCollateClause(false); err != nil {
				return err
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

func (p *Parser) parsePrimaryAtom(negated bool) (bool, error) {
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
			// O(1) fast check: if the matching ')' is absent or immediately
			// followed by a delimiter/EOF, `CASE(...)` cannot be `CASE (x) WHEN ... END`.
			closeOff := p.matchingCloseOffset(1)
			if closeOff > 0 {
				afterClose := p.peek(closeOff + 1).Kind
				if afterClose != tokEOF && afterClose != tokSemicolon &&
					afterClose != tokRParen && afterClose != tokRBracket &&
					afterClose != tokComma && afterClose != tokOperatorPipe {
					if p.hasCaseWhenAndEndAhead() {
						return false, p.parseCaseExpression()
					}
				}
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
	// In Spark default mode, TRUE/FALSE/NULL followed by '.' are qualified
	// identifiers (e.g. `SELECT FALSE.*`, `SELECT TRUE.x`), not constants.
	// Bare `INTERVAL` without a following numeric/string/signed literal (e.g.
	// `SELECT INTERVAL`, `SELECT INTERVAL zzz`, `SELECT INTERVAL NANOSECOND`)
	// falls back to a column reference, while malformed interval literals
	// (e.g. `INTERVAL 'abc' DAY`) fail immediately.
	if p.canStartConstant() && !((k == kwTRUE || k == kwFALSE || k == kwNULL) && p.peek(1).Kind == tokDot) {
		if k == kwINTERVAL {
			cp := p.save()
			next := p.peek(1).Kind
			if err := p.parseConstantWithNeg(negated); err != nil {
				isLitTok := func(tk TokenKind) bool {
					return tk == tokString || tk == tokInt || tk == tokDecimal ||
						tk == tokFloat || tk == tokDouble || tk == tokBigInt || tk == tokSmallInt ||
						tk == tokTinyInt || tk == tokBigDecimal || tk == tokExponent
				}
				isLiteralStart := isLitTok(next) || ((next == tokPlus || next == tokMinus) && isLitTok(p.peek(2).Kind))
				if isLiteralStart {
					return false, err
				}
				p.restore(cp)
			} else {
				return false, nil
			}
		} else {
			if err := p.parseConstantWithNeg(negated); err != nil {
				return false, err
			}
			return false, nil
		}
	}

	// currentLike keywords: CURRENT_DATE | CURRENT_TIMESTAMP | CURRENT_USER | USER | SESSION_USER
	if k == kwCURRENT_DATE || k == kwCURRENT_TIMESTAMP || k == kwCURRENT_USER || k == kwUSER || k == kwSESSION_USER {
		if p.disallowColumnRefs {
			return false, p.errorf("column reference %q is not allowed in time travel expression", p.cur().Text)
		}
		p.advance()
		return k.isIdentifier() && p.cur().Kind == tokDot, nil
	}

	// Column reference (identifier).
	if k.isIdentifier() {
		if p.disallowColumnRefs {
			return false, p.errorf("column reference %q is not allowed in time travel expression", p.cur().Text)
		}
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
	closeIdx := p.matchingCloseOffset(1)
	if closeIdx < 0 {
		return false
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
		hasQuantifier := p.isFunctionSetQuantifier() &&
			!(p.isQuantifierAmbiguousWithOperand() && p.hasKeywordAtDepthZero(1, kwAS))
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
		hasQuantifier := p.isFunctionSetQuantifier() &&
			!(p.isQuantifierAmbiguousWithOperand() && p.hasKeywordAtDepthZero(1, kwAS))
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
			!(p.peek(1).Kind == kwIGNORE && p.peek(2).Kind == kwNULLS && p.peek(3).Kind == tokRParen) &&
			!(p.isQuantifierAmbiguousWithOperand() && p.endsParenWithIgnoreNulls())
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
		hasQuantifier := p.isFunctionSetQuantifier() &&
			!(p.isQuantifierAmbiguousWithOperand() && p.hasKeywordAtDepthZero(1, kwIN))
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
		hasQuantifier := p.isFunctionSetQuantifier() &&
			!(p.isQuantifierAmbiguousWithOperand() && (p.hasKeywordAtDepthZero(1, kwFROM) || p.hasKeywordAtDepthZero(1, kwFOR)))
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
		hasQuantifier := p.isFunctionSetQuantifier() &&
			!(p.isQuantifierAmbiguousWithOperand() && p.hasKeywordAtDepthZero(1, kwFROM))
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
		hasQuantifier := p.isFunctionSetQuantifier() &&
			p.peek(1).Kind != kwPLACING &&
			!(p.isQuantifierAmbiguousWithOperand() && p.hasKeywordAtDepthZero(1, kwPLACING))
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
		// Spark v4.0.5 rejects quoted string literals as the datetime unit in 3-arg forms
		// (INVALID_PARAMETER_VALUE.DATETIME_UNIT).
		p.advance()
		p.advance()
		if p.match(tokRParen) {
			return p.parseFunctionCallSuffixes()
		}
		hasQuant := p.isFunctionSetQuantifier()
		if hasQuant {
			p.advance()
		}
		firstIsStringLit := !hasQuant && p.isStringLit() && p.peek(1).Kind == tokComma
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
		if _, err := p.expect(tokRParen); err != nil {
			return err
		}
		if argCount == 3 && firstIsStringLit {
			return p.errorf("datetime unit parameter must be an unquoted unit keyword, not a string literal")
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
	after := p.peek(2).Kind
	if after == kwWHEN {
		// `CASE <keyword> WHEN ...` uses the non-reserved keyword as the operand.
		return !next.isIdentifier()
	}
	if next.isIdentifier() {
		// Compound expressions beginning with an operator-like non-reserved
		// identifier (e.g. `CASE and + 1 WHEN ...`, `CASE collate + 1 WHEN ...`,
		// `CASE and(1) WHEN ...`, `CASE in(1) WHEN ...`) are CASE operands.
		if isExprContinuationOp(after) {
			return false
		}
		if after == tokLParen {
			if next != kwIN {
				return false
			}
			closeOff := p.matchingCloseOffset(2)
			if closeOff > 0 {
				afterClose := p.peek(closeOff + 1).Kind
				if afterClose == kwWHEN || isExprContinuationOp(afterClose) {
					return false
				}
			}
		}
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
	if n := len(p.openCaseRuns); n > 0 && p.openCaseRuns[n-1].parenDepth == p.parenDepth {
		return p.openCaseRuns[n-1].count
	}
	return 0
}

// hasCaseWhenAndEndAhead scans forward from p.pos+1 at paren/bracket depth 0
// to decide whether a `CASE` token starts a `CASE ... WHEN ... END` expression
// rather than a column reference `case` or a `CASE(...)` function call. When
// inside an outer open CASE at the same paren depth, it requires more unmatched
// `END` tokens than currently open CASEs at this depth so
// `SELECT CASE case + case + 1 WHEN 1 THEN 'ok' END` keeps the inner `case`s as
// columns while `SELECT CASE CASE -1 WHEN -1 THEN 2 END WHEN 2 THEN 3 END` and
// `SELECT CASE abs(CASE -1 WHEN -1 THEN 2 END) WHEN 2 THEN 3 END` parse nested
// CASE expressions cleanly. Results for the segment before the first unmatched
// END (including nested CASEs) and for all closed sibling/nested CASEs after it
// are memoized so arbitrarily deep nested or sibling CASEs parse in linear time.
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
	sawWhen := false
	lastWhenPos := -1
	ends := 0
	firstEndPos := -1
	type siblingNode struct {
		pos       int
		parent    int
		closedAnc int
		hasWhen   bool
		closed    bool
		endsAtPop int
		maxOpen   int
	}
	var nodes []siblingNode
	var siblingStack []int
	finish := func(stopPos int) bool {
		if p.parenDepth >= 0 {
			endPos := stopPos
			if firstEndPos >= 0 && firstEndPos < endPos {
				endPos = firstEndPos
			}
			if p.parenDepth >= len(p.caseScanCache) {
				newCap := cap(p.caseScanCache) * 2
				if newCap <= p.parenDepth {
					newCap = p.parenDepth + 1
				}
				if newCap < 16 {
					newCap = 16
				}
				grown := make([]caseScanEntry, p.parenDepth+1, newCap)
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
		if len(nodes) > 0 {
			for i := range nodes {
				par := nodes[i].parent
				if par < 0 {
					nodes[i].closedAnc = -1
				} else if nodes[par].closed {
					nodes[i].closedAnc = par
				} else {
					nodes[i].closedAnc = nodes[par].closedAnc
				}
				if !nodes[i].closed {
					continue
				}
				if anc := nodes[i].closedAnc; anc >= 0 {
					nodes[i].maxOpen = nodes[anc].maxOpen + 1
				} else {
					nodes[i].maxOpen = ends - nodes[i].endsAtPop
				}
				if p.caseScanConfirmedOpen == nil {
					p.caseScanConfirmedOpen = make(map[int]int, len(nodes))
				}
				spos := nodes[i].pos
				if prev, ok := p.caseScanConfirmedOpen[spos]; !ok || nodes[i].maxOpen > prev {
					p.caseScanConfirmedOpen[spos] = nodes[i].maxOpen
				}
			}
		}
		return sawWhen && (open == 0 || ends > open)
	}
	for i := 1; ; i++ {
		k := p.peek(i).Kind
		switch k {
		case tokEOF, tokSemicolon:
			return finish(p.pos + i)
		case tokLParen, tokLBracket:
			if closeOff := p.matchingCloseOffset(i); closeOff > i {
				i = closeOff
				continue
			}
			return finish(p.pos + i)
		case tokRParen, tokRBracket:
			return finish(p.pos + i)
		case kwWHEN:
			if open == 0 {
				return true
			}
			sawWhen = true
			lastWhenPos = p.pos + i
			if len(siblingStack) > 0 {
				nodes[siblingStack[len(siblingStack)-1]].hasWhen = true
			}
		case kwCASE:
			// After an END has closed the first CASE, a subsequent CASE at depth 0
			// may start a sibling CASE whose own END should not count toward closing
			// an enclosing open CASE.
			if ends > 0 {
				next := p.peek(i + 1).Kind
				if next == kwWHEN || (p.peek(i+2).Kind == kwWHEN && next.isIdentifier()) || !isCaseColumnNext(next) {
					parent := -1
					if len(siblingStack) > 0 {
						parent = siblingStack[len(siblingStack)-1]
					}
					siblingStack = append(siblingStack, len(nodes))
					nodes = append(nodes, siblingNode{pos: p.pos + i, parent: parent})
				}
			}
		case kwEND:
			next := p.peek(i + 1).Kind
			if !(p.peek(i-1).Kind == kwCASE && (next == kwWHEN || isExprContinuationOp(next) || next == tokLParen)) &&
				next != kwIF && next != kwCASE && next != kwLOOP && next != kwWHILE && next != kwREPEAT && next != kwFOR {
				// Discard any speculatively pushed sibling CASEs that never saw a WHEN
				// (they were column references like `case + 1`, not CASE expressions).
				for len(siblingStack) > 0 && !nodes[siblingStack[len(siblingStack)-1]].hasWhen {
					siblingStack = siblingStack[:len(siblingStack)-1]
				}
				if len(siblingStack) > 0 {
					topIdx := siblingStack[len(siblingStack)-1]
					siblingStack = siblingStack[:len(siblingStack)-1]
					nodes[topIdx].closed = true
					nodes[topIdx].endsAtPop = ends
				} else {
					if firstEndPos < 0 {
						firstEndPos = p.pos + i
					}
					ends++
				}
			}
		default:
			if p.isCaseClauseBoundaryAt(i) {
				return finish(p.pos + i)
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
	if n := len(p.openCaseRuns); n > 0 && p.openCaseRuns[n-1].parenDepth == p.parenDepth {
		p.openCaseRuns[n-1].count++
	} else {
		p.openCaseRuns = append(p.openCaseRuns, openCaseRun{parenDepth: p.parenDepth, count: 1})
	}
	defer func() {
		n := len(p.openCaseRuns)
		if p.openCaseRuns[n-1].count > 1 {
			p.openCaseRuns[n-1].count--
		} else {
			p.openCaseRuns = p.openCaseRuns[:n-1]
		}
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
		return p.startsWithQueryKeywordAt(2, true)
	}
	switch next {
	case tokRParen, tokComma, tokEOF, tokSemicolon,
		tokDot, tokLBracket, tokDoubleColon, kwCOLLATE,
		tokPlus, tokMinus, tokAsterisk, tokSlash, tokPercent, kwDIV,
		tokConcatPipe, tokPipe, tokAmpersand, tokHat,
		tokShiftLeft, tokShiftRight, tokShiftRightU,
		tokEq, tokNeq, tokNeqBang, tokNullSafeEq, tokLT, tokLTE, tokGT, tokGTE,
		tokArrow, tokFatArrow:
		return false
	case kwAND, kwOR, kwIS, kwIN, kwLIKE, kwILIKE, kwRLIKE, kwBETWEEN, kwNOT:
		return p.peek(2).Kind == tokRParen
	}
	return true
}

func (p *Parser) parseFunctionCall() error {
	// Single-pass handling for IDENTIFIER(expr) vs IDENTIFIER(expr)(args) so
	// deeply nested IDENTIFIER(IDENTIFIER(...)) never reparses its argument.
	// Spark v4.0.5 rejects empty `IDENTIFIER()` with PARSE_SYNTAX_ERROR.
	if p.cur().Kind == kwIDENTIFIER && p.peek(1).Kind == tokLParen {
		p.advance()
		p.advance()
		if p.cur().Kind == tokRParen {
			return p.errorf("IDENTIFIER() requires an expression argument")
		}
		singlePlainArg := false
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
		tokArrow, tokFatArrow:
		return false
	case kwAND, kwOR, kwIS, kwIN, kwLIKE, kwILIKE, kwRLIKE, kwBETWEEN, kwFROM, kwAS:
		after := p.peek(offset + 2).Kind
		return after == tokRParen || after == tokComma
	}
	return true
}

// isQuantifierAmbiguousWithOperand reports whether a leading ALL/DISTINCT is
// followed by a token (`+`, `-`, `(`) that can either start a quantified unary/
// parenthesized expression (`ALL -1`, `DISTINCT (x)`) or continue an identifier
// operand (`all + 1`, `distinct - 1`, `all(1)`).
func (p *Parser) isQuantifierAmbiguousWithOperand() bool {
	next := p.peek(1).Kind
	return next == tokPlus || next == tokMinus || next == tokLParen
}

// endsParenWithIgnoreNulls reports whether the enclosing `(...)` whose '(' is
// at offset -1 ends with `IGNORE NULLS)`.
func (p *Parser) endsParenWithIgnoreNulls() bool {
	closeOff := p.matchingCloseOffset(-1)
	return closeOff >= 3 &&
		p.peek(closeOff-2).Kind == kwIGNORE &&
		p.peek(closeOff-1).Kind == kwNULLS
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

func isSupportedLiteralType(tok Token) bool {
	switch tok.Kind {
	case kwDATE, kwTIMESTAMP, kwTIMESTAMP_LTZ, kwTIMESTAMP_NTZ, kwINTERVAL, kwBINARY_HEX:
		return true
	}
	return false
}

func (p *Parser) canStartConstant() bool {
	switch p.cur().Kind {
	case kwNULL, tokQuestion, tokColon, kwINTERVAL,
		tokInt, tokDecimal, tokExponent, tokBigInt, tokSmallInt, tokTinyInt,
		tokFloat, tokDouble, tokBigDecimal, tokMinus,
		kwTRUE, kwFALSE, tokString, tokDoubleQuoted:
		return true
	}
	// literalType stringLit: DATE '...', TIMESTAMP '...', TIMESTAMP_LTZ '...', TIMESTAMP_NTZ '...', X '...'
	if isSupportedLiteralType(p.cur()) && (p.peek(1).Kind == tokString || p.peek(1).Kind == tokDoubleQuoted) {
		return true
	}
	return false
}

func (p *Parser) parseConstant() error {
	return p.parseConstantWithNeg(false)
}

func (p *Parser) parseConstantWithNeg(negated bool) error {
	switch p.cur().Kind {
	case kwNULL, kwTRUE, kwFALSE:
		p.advance()
		return nil
	case tokQuestion:
		if p.disallowParamMarkers {
			return p.errorf("parameter markers are not allowed in this statement")
		}
		p.advance()
		return nil
	case tokColon:
		if p.disallowParamMarkers {
			return p.errorf("parameter markers are not allowed in this statement")
		}
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
			litTok := p.advance()
			if !validateTypedLiteral(kwINTERVAL, litTok.Text) {
				return p.errorf("invalid INTERVAL typed literal %s", litTok.Text)
			}
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
		return p.parseNumberWithNeg(negated)
	}

	// literalType stringLit (DATE '...', TIMESTAMP '...', TIMESTAMP_LTZ '...', TIMESTAMP_NTZ '...', X '...')
	if isSupportedLiteralType(p.cur()) && (p.peek(1).Kind == tokString || p.peek(1).Kind == tokDoubleQuoted) {
		typeTok := p.advance()
		litTok := p.advance()
		if !validateTypedLiteral(typeTok.Kind, litTok.Text) {
			return p.errorf("invalid %s typed literal %s", typeTok.Text, litTok.Text)
		}
		return nil
	}
	return p.errorf("expected constant, got %q", p.cur().Text)
}

func (p *Parser) parseNumber() error {
	return p.parseNumberWithNeg(false)
}

func (p *Parser) parseNumberWithNeg(negated bool) error {
	cp := p.save()
	if p.match(tokMinus) {
		negated = true
	}
	tok := p.cur()
	switch tok.Kind {
	case tokInt, tokDecimal, tokExponent, tokBigInt, tokSmallInt, tokTinyInt,
		tokFloat, tokDouble, tokBigDecimal:
		if !validateNumericLiteral(tok.Kind, tok.Text, negated) {
			return p.errorf("numeric literal %q out of range", tok.Text)
		}
		p.advance()
		return nil
	}
	p.restore(cp)
	return p.errorf("expected numeric literal, got %q", p.cur().Text)
}

// parseInterval parses:
//
//	interval : INTERVAL (multiUnitsInterval | unitToUnitInterval) ;
//
// Spark's errorCapturingMultiUnitsInterval and errorCapturingUnitToUnitInterval
// grammar rules exist only so AstBuilder can throw tailored diagnostics for
// mixed/repeated interval continuations; strict Parse rejects them directly.
func (p *Parser) parseInterval() error {
	_, err := p.parseIntervalWithSeconds()
	return err
}

func (p *Parser) parseIntervalWithSeconds() (int64, error) {
	if _, err := p.expect(kwINTERVAL); err != nil {
		return 0, err
	}
	// Try unitToUnitInterval first if we see: intervalValue unitInUnitToUnit TO unitInUnitToUnit
	cp := p.save()
	if secs, err := p.parseUnitToUnitInterval(); err == nil {
		return secs, nil
	}
	p.restore(cp)
	return p.parseMultiUnitsInterval()
}

func (p *Parser) parseMultiUnitsInterval() (int64, error) {
	count := 0
	var totalSecs int64
	for {
		cp := p.save()
		sign, rawVal, isStr, err := p.parseIntervalValueToken()
		if err != nil {
			p.restore(cp)
			break
		}
		if !p.isUnitInMultiUnits() {
			p.restore(cp)
			break
		}
		unitTok := p.advance()
		// Don't steal the start of a unitToUnitInterval (value unit TO unit) if we already have >=1 unit.
		if count > 0 && p.cur().Kind == kwTO && p.isUnitInUnitToUnitAt(1) {
			p.restore(cp)
			break
		}
		// Spark v4.0.5 rejects NANOSECOND/NANOSECONDS in multi-unit intervals
		// (INVALID_INTERVAL_FORMAT.INVALID_UNIT).
		if !isSupportedMultiUnit(unitTok.Kind) {
			return 0, p.errorf("unsupported interval unit %q", unitTok.Text)
		}
		numStr := rawVal
		if isStr {
			numStr = strings.TrimSpace(unquoteStringLit(rawVal))
		}
		if !isValidMultiUnitNumber(numStr, unitTok.Kind) {
			return 0, p.errorf("invalid interval value %q for unit %q", rawVal, unitTok.Text)
		}
		totalSecs += intervalUnitSeconds(sign, numStr, unitTok.Kind)
		count++
	}
	if count == 0 {
		return 0, p.errorf("expected interval value and unit")
	}
	return totalSecs, nil
}

func intervalUnitSeconds(sign int64, numStr string, unit TokenKind) int64 {
	s := numStr
	if strings.HasPrefix(s, "+") {
		s = s[1:]
	} else if strings.HasPrefix(s, "-") {
		sign = -sign
		s = s[1:]
	}
	intPart, _, _ := strings.Cut(s, ".")
	v, err := strconv.ParseInt(intPart, 10, 64)
	if err != nil {
		return sign * 1_000_000_000 // overflow sentinel (outside +/-18h)
	}
	var mult int64
	switch unit {
	case kwMICROSECOND, kwMICROSECONDS, kwMILLISECOND, kwMILLISECONDS:
		return 0
	case kwSECOND, kwSECONDS:
		mult = 1
	case kwMINUTE, kwMINUTES:
		mult = 60
	case kwHOUR, kwHOURS:
		mult = 3600
	case kwDAY, kwDAYS:
		mult = 86400
	case kwWEEK, kwWEEKS:
		mult = 7 * 86400
	case kwMONTH, kwMONTHS:
		mult = 30 * 86400
	case kwYEAR, kwYEARS:
		mult = 365 * 86400
	}
	return sign * v * mult
}

func (p *Parser) isUnitInUnitToUnitAt(offset int) bool {
	switch p.peek(offset).Kind {
	case kwSECOND, kwMINUTE, kwHOUR, kwDAY, kwMONTH, kwYEAR:
		return true
	}
	return false
}

func (p *Parser) parseUnitToUnitInterval() (int64, error) {
	sign, rawVal, isStr, err := p.parseIntervalValueToken()
	if err != nil {
		return 0, err
	}
	if !p.isUnitInUnitToUnit() {
		return 0, p.errorf("expected interval unit")
	}
	fromUnit := p.advance().Kind
	if _, err := p.expect(kwTO); err != nil {
		return 0, err
	}
	if !p.isUnitInUnitToUnit() {
		return 0, p.errorf("expected interval unit after TO")
	}
	toUnit := p.advance().Kind
	// Spark v4.0.5 requires a string literal for from-to intervals and validates
	// both the unit pair and the literal format.
	if !isStr {
		return 0, p.errorf("from-to interval value must be a string literal")
	}
	secs, ok := parseUnitToUnitIntervalValue(unquoteStringLit(rawVal), fromUnit, toUnit, sign)
	if !ok {
		return 0, p.errorf("invalid from-to interval %s", rawVal)
	}
	return secs, nil
}

func (p *Parser) parseIntervalValueToken() (sign int64, raw string, isStr bool, err error) {
	sign = 1
	if p.match(tokPlus) {
		sign = 1
	} else if p.match(tokMinus) {
		sign = -1
	}
	tok := p.cur()
	switch tok.Kind {
	case tokInt, tokDecimal:
		p.advance()
		return sign, tok.Text, false, nil
	case tokString, tokDoubleQuoted:
		p.advance()
		return sign, tok.Text, true, nil
	default:
		return 0, "", false, p.errorf("expected interval value")
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
