package parser

import "fmt"

// parseExpression parses IdentifiersParser.g's expression rule.
func (p *parser) parseExpression() error {
	if err := p.enter(); err != nil {
		return err
	}
	defer p.leave()
	return p.parsePrecedenceOrExpression()
}

func (p *parser) parsePrecedenceOrExpression() error {
	if err := p.parsePrecedenceAndExpression(); err != nil {
		return err
	}
	for p.match(KW_OR) {
		if err := p.parsePrecedenceAndExpression(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parsePrecedenceAndExpression() error {
	if err := p.parsePrecedenceNotExpression(); err != nil {
		return err
	}
	for p.match(KW_AND) {
		if err := p.parsePrecedenceNotExpression(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parsePrecedenceNotExpression() error {
	for p.match(KW_NOT) {
	}
	return p.parsePrecedenceUnarySuffixExpression()
}

func (p *parser) parsePrecedenceUnarySuffixExpression() error {
	if err := p.parsePrecedenceEqualExpression(); err != nil {
		return err
	}
	if p.peek().Type == KW_IS {
		// IS [NOT] (NULL | TRUE | FALSE | UNKNOWN)
		saved := p.save()
		p.next()
		if err := p.parseIsCondition(); err == nil {
			return nil
		}
		p.restore(saved)
	}
	return nil
}

func (p *parser) parseIsCondition() error {
	p.match(KW_NOT)
	switch p.peek().Type {
	case KW_NULL, KW_TRUE, KW_FALSE, KW_UNKNOWN:
		p.next()
		return nil
	default:
		return p.errAt(p.peek(), fmt.Sprintf("expected NULL/TRUE/FALSE/UNKNOWN after IS, got %q", p.peek().Str))
	}
}

func (p *parser) parsePrecedenceEqualExpression() error {
	if err := p.parsePrecedenceSimilarExpression(); err != nil {
		return err
	}
	for {
		switch p.peek().Type {
		case TokEqual, TokEqualNS, TokNotEqual:
			p.next()
			if err := p.parsePrecedenceSimilarExpression(); err != nil {
				return err
			}
		case KW_IS:
			// KW_IS KW_NOT KW_DISTINCT KW_FROM | KW_IS KW_DISTINCT KW_FROM
			if p.peekAt(1).Type == KW_NOT && p.peekAt(2).Type == KW_DISTINCT && p.peekAt(3).Type == KW_FROM {
				p.pos += 4
				if err := p.parsePrecedenceSimilarExpression(); err != nil {
					return err
				}
			} else if p.peekAt(1).Type == KW_DISTINCT && p.peekAt(2).Type == KW_FROM {
				p.pos += 3
				if err := p.parsePrecedenceSimilarExpression(); err != nil {
					return err
				}
			} else {
				return nil
			}
		default:
			return nil
		}
	}
}

func (p *parser) parsePrecedenceSimilarExpression() error {
	if p.match(KW_EXISTS) {
		return p.parseSubQueryExpression()
	}
	if err := p.parsePrecedenceBitwiseOrExpression(); err != nil {
		return err
	}
	return p.parsePrecedenceSimilarExpressionPartOpt()
}

func isSubQuerySelectorOperator(tt TokenType) bool {
	switch tt {
	case TokEqual, TokNotEqual, TokLessThanOrEqualTo, TokLessThan, TokGreaterThanOrEqualTo, TokGreaterThan:
		return true
	}
	return false
}

func isQuantifierType(tt TokenType) bool {
	return tt == KW_ANY || tt == KW_SOME || tt == KW_ALL
}

func isRegexpOperator(tt TokenType) bool {
	return tt == KW_LIKE || tt == KW_RLIKE || tt == KW_REGEXP
}

func (p *parser) parsePrecedenceSimilarExpressionPartOpt() error {
	t := p.peek().Type
	// Check quantifier predicate first: subQuerySelectorOperator quantifierType subQueryExpression
	if isSubQuerySelectorOperator(t) && isQuantifierType(p.peekAt(1).Type) && p.peekAt(2).Type == TokLParen {
		p.pos += 2
		return p.parseSubQueryExpression()
	}
	// LIKE ANY (...) / LIKE ALL (...)
	if t == KW_LIKE && (p.peekAt(1).Type == KW_ANY || p.peekAt(1).Type == KW_ALL) && p.peekAt(2).Type == TokLParen {
		p.pos += 2
		return p.parseExpressionsInParenthesis()
	}
	// precedenceSimilarOperator: LIKE | RLIKE | REGEXP | <= | < | >= | >
	switch t {
	case KW_LIKE, KW_RLIKE, KW_REGEXP, TokLessThanOrEqualTo, TokLessThan, TokGreaterThanOrEqualTo, TokGreaterThan:
		p.next()
		return p.parsePrecedenceBitwiseOrExpression()
	case KW_IN, KW_BETWEEN:
		return p.parsePrecedenceSimilarExpressionAtom()
	case KW_NOT:
		// Only consume NOT if it introduces a valid similar-part-not.
		next := p.peekAt(1).Type
		if isRegexpOperator(next) || next == KW_IN || next == KW_BETWEEN ||
			(isSubQuerySelectorOperator(next) && isQuantifierType(p.peekAt(2).Type)) {
			p.next()
			return p.parsePrecedenceSimilarExpressionPartNot()
		}
	}
	return nil
}

func (p *parser) parsePrecedenceSimilarExpressionAtom() error {
	switch p.peek().Type {
	case KW_IN:
		p.next()
		return p.parsePrecedenceSimilarExpressionIn()
	case KW_BETWEEN:
		p.next()
		if err := p.parsePrecedenceBitwiseOrExpression(); err != nil {
			return err
		}
		if _, err := p.expect(KW_AND); err != nil {
			return err
		}
		return p.parsePrecedenceBitwiseOrExpression()
	case KW_LIKE:
		p.next()
		if p.peek().Type == KW_ANY || p.peek().Type == KW_ALL {
			p.next()
			return p.parseExpressionsInParenthesis()
		}
		return p.errAt(p.peek(), fmt.Sprintf("expected ANY or ALL after LIKE, got %q", p.peek().Str))
	default:
		if isSubQuerySelectorOperator(p.peek().Type) {
			p.next()
			if !isQuantifierType(p.peek().Type) {
				return p.errAt(p.peek(), fmt.Sprintf("expected ANY, SOME, or ALL, got %q", p.peek().Str))
			}
			p.next()
			return p.parseSubQueryExpression()
		}
		return p.errAt(p.peek(), fmt.Sprintf("unexpected token %q in predicate", p.peek().Str))
	}
}

func (p *parser) parsePrecedenceSimilarExpressionIn() error {
	if p.try(p.parseSubQueryExpression) {
		return nil
	}
	return p.parseExpressionsInParenthesis()
}

func (p *parser) parsePrecedenceSimilarExpressionPartNot() error {
	t := p.peek().Type
	if t == KW_LIKE && (p.peekAt(1).Type == KW_ANY || p.peekAt(1).Type == KW_ALL) && p.peekAt(2).Type == TokLParen {
		return p.parsePrecedenceSimilarExpressionAtom()
	}
	if isRegexpOperator(t) {
		p.next()
		return p.parsePrecedenceBitwiseOrExpression()
	}
	return p.parsePrecedenceSimilarExpressionAtom()
}

func (p *parser) parseSubQueryExpression() error {
	if _, err := p.expect(TokLParen); err != nil {
		return err
	}
	if err := p.parseSelectStatement(); err != nil {
		return err
	}
	_, err := p.expect(TokRParen)
	return err
}

func (p *parser) parsePrecedenceBitwiseOrExpression() error {
	if err := p.parsePrecedenceAmpersandExpression(); err != nil {
		return err
	}
	for p.match(TokBitwiseOr) {
		if err := p.parsePrecedenceAmpersandExpression(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parsePrecedenceAmpersandExpression() error {
	if err := p.parsePrecedenceConcatenateExpression(); err != nil {
		return err
	}
	for p.match(TokAmpersand) {
		if err := p.parsePrecedenceConcatenateExpression(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parsePrecedenceConcatenateExpression() error {
	if err := p.parsePrecedencePlusExpression(); err != nil {
		return err
	}
	for p.match(TokConcatenate) {
		if err := p.parsePrecedencePlusExpression(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parsePrecedencePlusExpression() error {
	if err := p.parsePrecedenceStarExpression(); err != nil {
		return err
	}
	for p.peek().Type == TokPlus || p.peek().Type == TokMinus {
		p.next()
		if err := p.parsePrecedenceStarExpression(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parsePrecedenceStarExpression() error {
	if err := p.parsePrecedenceBitwiseXorExpression(); err != nil {
		return err
	}
	for {
		switch p.peek().Type {
		case TokStar, TokDivide, TokMod, TokDiv:
			p.next()
			if err := p.parsePrecedenceBitwiseXorExpression(); err != nil {
				return err
			}
		default:
			return nil
		}
	}
}

func (p *parser) parsePrecedenceBitwiseXorExpression() error {
	if err := p.parsePrecedenceUnaryPrefixExpression(); err != nil {
		return err
	}
	for p.match(TokBitwiseXor) {
		if err := p.parsePrecedenceUnaryPrefixExpression(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parsePrecedenceUnaryPrefixExpression() error {
	for p.peek().Type == TokPlus || p.peek().Type == TokMinus || p.peek().Type == TokTilde {
		p.next()
	}
	return p.parsePrecedenceFieldExpression()
}

func (p *parser) parsePrecedenceFieldExpression() error {
	if err := p.parseAtomExpression(); err != nil {
		return err
	}
	for {
		if p.match(TokLSquare) {
			if err := p.parseExpression(); err != nil {
				return err
			}
			if _, err := p.expect(TokRSquare); err != nil {
				return err
			}
		} else if p.match(TokDot) {
			if _, err := p.parseIdentifier(); err != nil {
				return err
			}
		} else {
			return nil
		}
	}
}

func (p *parser) parseAtomExpression() error {
	t := p.peek()
	switch t.Type {
	case KW_CAST:
		return p.parseCastExpression()
	case KW_EXTRACT:
		return p.parseExtractExpression()
	case KW_FLOOR:
		return p.parseFloorExpression()
	case KW_CASE:
		return p.parseCaseOrWhenExpression()
	case KW_INTERVAL:
		return p.parseIntervalExpression()
	case KW_DATE:
		if p.peekAt(1).Type == TokStringLiteral {
			p.pos += 2
			return nil
		}
		if p.peekAt(1).Type == TokLParen {
			return p.parseFunction()
		}
		return p.errAt(t, "unexpected DATE in expression")
	case KW_TIMESTAMP:
		if p.peekAt(1).Type == TokStringLiteral {
			p.pos += 2
			return nil
		}
		if p.peekAt(1).Type == TokLParen {
			return p.parseFunction()
		}
		return p.errAt(t, "unexpected TIMESTAMP in expression")
	case KW_TIMESTAMPLOCALTZ:
		if p.peekAt(1).Type == TokStringLiteral {
			p.pos += 2
			return nil
		}
		return p.errAt(t, "expected string literal after TIMESTAMPLOCALTZ")
	case KW_CURRENT_DATE, KW_CURRENT_TIMESTAMP:
		if p.peekAt(1).Type == TokLParen {
			return p.parseFunction()
		}
		p.next()
		return nil
	case TokLParen:
		// Could be: intervalExpression '(val) qual', subQueryExpression '(selectStatement)', or expressionsInParenthesis.
		if p.try(p.parseParenIntervalExpression) {
			return nil
		}
		if p.try(p.parseSubQueryExpression) {
			return nil
		}
		return p.parseExpressionsInParenthesis()
	case TokNumber, TokStringLiteral, TokIntegralLiteral, TokNumberLiteral, TokCharSetName, KW_TRUE, KW_FALSE, KW_NULL, TokQuestion:
		return p.parseConstant()
	}

	// TRIM(...) special form or function call.
	if t.Type == KW_TRIM && p.peekAt(1).Type == TokLParen {
		if p.try(p.parseTrimFunction) {
			return nil
		}
	}

	// (functionName LPAREN) => function
	if isFunctionNameStart(t.Type) {
		if p.peekAt(1).Type == TokLParen || (isIdentifier(t.Type) && p.peekAt(1).Type == TokDot && isIdentifier(p.peekAt(2).Type) && p.peekAt(3).Type == TokLParen) {
			return p.parseFunction()
		}
	}

	// tableOrColumn
	if isIdentifier(t.Type) {
		p.next()
		return nil
	}
	return p.errAt(t, fmt.Sprintf("unexpected token %q in expression", t.Str))
}

func (p *parser) parseConstant() error {
	t := p.peek()
	switch t.Type {
	case TokNumber:
		// Check intervalLiteral: Number intervalQualifiers
		if p.try(p.parseIntervalLiteral) {
			return nil
		}
		p.next()
		return nil
	case TokStringLiteral:
		// Check intervalLiteral: StringLiteral intervalQualifiers
		if p.try(p.parseIntervalLiteral) {
			return nil
		}
		p.next()
		// stringLiteralSequence: StringLiteral StringLiteral+
		for p.peek().Type == TokStringLiteral {
			p.next()
		}
		return nil
	case TokIntegralLiteral, TokNumberLiteral, KW_TRUE, KW_FALSE, KW_NULL, TokQuestion, KW_CURRENT_DATE, KW_CURRENT_TIMESTAMP:
		p.next()
		return nil
	case KW_DATE, KW_TIMESTAMP, KW_TIMESTAMPLOCALTZ:
		p.next()
		_, err := p.expect(TokStringLiteral)
		return err
	case TokCharSetName:
		p.next()
		if p.peek().Type == TokStringLiteral || p.peek().Type == TokCharSetLiteral {
			p.next()
			return nil
		}
		return p.errAt(p.peek(), fmt.Sprintf("expected string or hex literal after charset %q", t.Str))
	default:
		return p.errAt(t, fmt.Sprintf("expected constant, got %q", t.Str))
	}
}

func (p *parser) parseIntervalLiteral() error {
	if p.peek().Type != TokStringLiteral && p.peek().Type != TokNumber {
		return p.errAt(p.peek(), "expected interval value")
	}
	p.next()
	return p.parseIntervalQualifiers()
}

func (p *parser) parseParenIntervalExpression() error {
	if _, err := p.expect(TokLParen); err != nil {
		return err
	}
	if p.peek().Type != TokStringLiteral && p.peek().Type != TokNumber {
		return p.errAt(p.peek(), "expected interval value")
	}
	p.next()
	if _, err := p.expect(TokRParen); err != nil {
		return err
	}
	return p.parseIntervalQualifiers()
}

func (p *parser) parseIntervalExpression() error {
	if _, err := p.expect(KW_INTERVAL); err != nil {
		return err
	}
	if p.match(TokLParen) {
		if err := p.parseExpression(); err != nil {
			return err
		}
		if _, err := p.expect(TokRParen); err != nil {
			return err
		}
		return p.parseIntervalQualifiers()
	}
	if p.peek().Type != TokStringLiteral && p.peek().Type != TokNumber {
		return p.errAt(p.peek(), fmt.Sprintf("expected interval value, got %q", p.peek().Str))
	}
	p.next()
	return p.parseIntervalQualifiers()
}

func (p *parser) parseIntervalQualifiers() error {
	switch p.peek().Type {
	case KW_YEAR:
		p.next()
		if p.peek().Type == KW_TO && p.peekAt(1).Type == KW_MONTH {
			p.pos += 2
		}
		return nil
	case KW_DAY:
		p.next()
		if p.peek().Type == KW_TO && p.peekAt(1).Type == KW_SECOND {
			p.pos += 2
		}
		return nil
	case KW_MONTH, KW_HOUR, KW_MINUTE, KW_SECOND:
		p.next()
		return nil
	default:
		return p.errAt(p.peek(), fmt.Sprintf("expected interval qualifier, got %q", p.peek().Str))
	}
}

func (p *parser) parseCastExpression() error {
	if _, err := p.expect(KW_CAST); err != nil {
		return err
	}
	if _, err := p.expect(TokLParen); err != nil {
		return err
	}
	if err := p.parseExpression(); err != nil {
		return err
	}
	if _, err := p.expect(KW_AS); err != nil {
		return err
	}
	if err := p.parseType(); err != nil {
		return err
	}
	if p.match(KW_FORMAT) {
		if _, err := p.expect(TokStringLiteral); err != nil {
			return err
		}
	}
	_, err := p.expect(TokRParen)
	return err
}

func (p *parser) parseExtractExpression() error {
	if _, err := p.expect(KW_EXTRACT); err != nil {
		return err
	}
	if _, err := p.expect(TokLParen); err != nil {
		return err
	}
	if err := p.parseTimeQualifiers(); err != nil {
		return err
	}
	if _, err := p.expect(KW_FROM); err != nil {
		return err
	}
	if err := p.parseExpression(); err != nil {
		return err
	}
	_, err := p.expect(TokRParen)
	return err
}

func (p *parser) parseTimeQualifiers() error {
	switch p.peek().Type {
	case KW_YEAR, KW_QUARTER, KW_MONTH, KW_WEEK, KW_DAY, KW_DOW, KW_HOUR, KW_MINUTE, KW_SECOND:
		p.next()
		return nil
	default:
		return p.errAt(p.peek(), fmt.Sprintf("expected time qualifier, got %q", p.peek().Str))
	}
}

func (p *parser) parseTimeUnitQualifiers() error {
	switch p.peek().Type {
	case KW_DAY, KW_HOUR, KW_MINUTE:
		p.next()
		return nil
	default:
		return p.errAt(p.peek(), fmt.Sprintf("expected time unit qualifier (DAY/HOUR/MINUTE), got %q", p.peek().Str))
	}
}

func (p *parser) parseFloorExpression() error {
	if _, err := p.expect(KW_FLOOR); err != nil {
		return err
	}
	if _, err := p.expect(TokLParen); err != nil {
		return err
	}
	if err := p.parseExpression(); err != nil {
		return err
	}
	if p.match(KW_TO) {
		switch p.peek().Type {
		case KW_YEAR, KW_QUARTER, KW_MONTH, KW_WEEK, KW_DAY, KW_HOUR, KW_MINUTE, KW_SECOND:
			p.next()
		default:
			return p.errAt(p.peek(), fmt.Sprintf("expected floor date qualifier, got %q", p.peek().Str))
		}
	}
	_, err := p.expect(TokRParen)
	return err
}

func (p *parser) parseCaseOrWhenExpression() error {
	if _, err := p.expect(KW_CASE); err != nil {
		return err
	}
	if p.peek().Type != KW_WHEN {
		if err := p.parseExpression(); err != nil {
			return err
		}
	}
	if p.peek().Type != KW_WHEN {
		return p.errAt(p.peek(), fmt.Sprintf("expected WHEN in CASE expression, got %q", p.peek().Str))
	}
	for p.match(KW_WHEN) {
		if err := p.parseExpression(); err != nil {
			return err
		}
		if _, err := p.expect(KW_THEN); err != nil {
			return err
		}
		if err := p.parseExpression(); err != nil {
			return err
		}
	}
	if p.match(KW_ELSE) {
		if err := p.parseExpression(); err != nil {
			return err
		}
	}
	_, err := p.expect(KW_END)
	return err
}

func (p *parser) parseTrimFunction() error {
	if _, err := p.expect(KW_TRIM); err != nil {
		return err
	}
	if _, err := p.expect(TokLParen); err != nil {
		return err
	}
	switch p.peek().Type {
	case KW_LEADING, KW_TRAILING, KW_BOTH:
		p.next()
	}
	if p.peek().Type != KW_FROM {
		if err := p.parseSelectExpression(); err != nil {
			return err
		}
	}
	if _, err := p.expect(KW_FROM); err != nil {
		return err
	}
	if err := p.parseSelectExpression(); err != nil {
		return err
	}
	_, err := p.expect(TokRParen)
	return err
}

func (p *parser) parseFunctionName() error {
	if sql11ReservedFunctionNames[p.peek().Type] {
		p.next()
		return nil
	}
	return p.parseFunctionIdentifier()
}

func (p *parser) parseFunction() error {
	if p.peek().Type == KW_TRIM && p.peekAt(1).Type == TokLParen {
		if p.try(p.parseTrimFunction) {
			return nil
		}
	}
	if err := p.parseFunctionName(); err != nil {
		return err
	}
	if _, err := p.expect(TokLParen); err != nil {
		return err
	}
	if p.peek().Type == TokStar {
		p.next()
	} else {
		if p.peek().Type == KW_DISTINCT || p.peek().Type == KW_ALL {
			p.next()
		}
		if p.peek().Type != TokRParen && !p.isNullTreatment() {
			if err := p.parseSelectExpression(); err != nil {
				return err
			}
			for p.match(TokComma) {
				if err := p.parseSelectExpression(); err != nil {
					return err
				}
			}
		}
	}
	// Closing paren and optional null treatment / WITHIN GROUP / OVER window_specification.
	if p.isNullTreatment() {
		p.pos += 2
		if _, err := p.expect(TokRParen); err != nil {
			return err
		}
		if _, err := p.expect(KW_OVER); err != nil {
			return err
		}
		return p.parseWindowSpecification()
	}
	if _, err := p.expect(TokRParen); err != nil {
		return err
	}
	if p.match(KW_WITHIN) {
		if _, err := p.expect(KW_GROUP); err != nil {
			return err
		}
		if _, err := p.expect(TokLParen); err != nil {
			return err
		}
		if err := p.parseOrderByClause(); err != nil {
			return err
		}
		_, err := p.expect(TokRParen)
		return err
	}
	if p.isNullTreatment() && p.peekAt(2).Type == KW_OVER {
		p.pos += 3
		return p.parseWindowSpecification()
	}
	if p.match(KW_OVER) {
		return p.parseWindowSpecification()
	}
	return nil
}

func (p *parser) isNullTreatment() bool {
	return (p.peek().Type == KW_RESPECT || p.peek().Type == KW_IGNORE) && p.peekAt(1).Type == KW_NULLS
}

func (p *parser) parseExpressionOrDefault() error {
	if p.peek().Type == KW_DEFAULT && p.peekAt(1).Type != TokDot {
		p.next()
		return nil
	}
	return p.parseExpression()
}

func (p *parser) parseExpressionsInParenthesis() error {
	if _, err := p.expect(TokLParen); err != nil {
		return err
	}
	if err := p.parseExpressionsNotInParenthesis(); err != nil {
		return err
	}
	_, err := p.expect(TokRParen)
	return err
}

func (p *parser) parseExpressionsNotInParenthesis() error {
	if err := p.parseExpressionOrDefault(); err != nil {
		return err
	}
	for p.match(TokComma) {
		if err := p.parseExpressionOrDefault(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseExpressions() error {
	if p.try(p.parseExpressionsInParenthesis) {
		return nil
	}
	return p.parseExpressionsNotInParenthesis()
}

func (p *parser) parseExpressionList() error {
	if err := p.parseExpression(); err != nil {
		return err
	}
	for p.match(TokComma) {
		if err := p.parseExpression(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseSelectExpression() error {
	if p.try(p.parseTableAllColumns) {
		return nil
	}
	return p.parseExpression()
}

func (p *parser) parseSelectExpressionList() error {
	if err := p.parseSelectExpression(); err != nil {
		return err
	}
	for p.match(TokComma) {
		if err := p.parseSelectExpression(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseTableAllColumns() error {
	if p.match(TokStar) {
		return nil
	}
	if err := p.parseTableName(); err != nil {
		return err
	}
	if _, err := p.expect(TokDot); err != nil {
		return err
	}
	_, err := p.expect(TokStar)
	return err
}
