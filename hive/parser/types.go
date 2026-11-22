package parser

import "fmt"

// parseType parses HiveQL data types: primitive, ARRAY, STRUCT, MAP, UNIONTYPE.
func (p *parser) parseType() error {
	if err := p.enter(); err != nil {
		return err
	}
	defer p.leave()

	switch p.peek().Type {
	case KW_ARRAY:
		return p.parseListType()
	case KW_STRUCT:
		return p.parseStructType()
	case KW_MAP:
		return p.parseMapType()
	case KW_UNIONTYPE:
		return p.parseUnionType()
	default:
		return p.parsePrimitiveType()
	}
}

func (p *parser) parsePrimitiveType() error {
	t := p.peek()
	switch t.Type {
	case KW_TINYINT, KW_SMALLINT, KW_INT, KW_BIGINT, KW_BOOLEAN,
		KW_FLOAT, KW_REAL, KW_DATE, KW_DATETIME, KW_TIMESTAMPLOCALTZ,
		KW_STRING, KW_BINARY:
		p.next()
		return nil
	case KW_DOUBLE:
		p.next()
		p.match(KW_PRECISION)
		return nil
	case KW_TIMESTAMP:
		p.next()
		if p.peek().Type == KW_WITH && p.peekAt(1).Type == KW_LOCAL &&
			p.peekAt(2).Type == KW_TIME && p.peekAt(3).Type == KW_ZONE {
			p.pos += 4
		}
		return nil
	case KW_DECIMAL:
		p.next()
		if p.match(TokLParen) {
			if _, err := p.expect(TokNumber); err != nil {
				return err
			}
			if p.match(TokComma) {
				if _, err := p.expect(TokNumber); err != nil {
					return err
				}
			}
			if _, err := p.expect(TokRParen); err != nil {
				return err
			}
		}
		return nil
	case KW_VARCHAR, KW_CHAR:
		p.next()
		if _, err := p.expect(TokLParen); err != nil {
			return err
		}
		if _, err := p.expect(TokNumber); err != nil {
			return err
		}
		_, err := p.expect(TokRParen)
		return err
	default:
		return p.errAt(t, fmt.Sprintf("expected type name, got %q", t.Str))
	}
}

func (p *parser) parseListType() error {
	if _, err := p.expect(KW_ARRAY); err != nil {
		return err
	}
	if _, err := p.expect(TokLessThan); err != nil {
		return err
	}
	if err := p.parseType(); err != nil {
		return err
	}
	_, err := p.expect(TokGreaterThan)
	return err
}

func (p *parser) parseStructType() error {
	if _, err := p.expect(KW_STRUCT); err != nil {
		return err
	}
	if _, err := p.expect(TokLessThan); err != nil {
		return err
	}
	if err := p.parseColumnNameColonTypeList(); err != nil {
		return err
	}
	_, err := p.expect(TokGreaterThan)
	return err
}

func (p *parser) parseMapType() error {
	if _, err := p.expect(KW_MAP); err != nil {
		return err
	}
	if _, err := p.expect(TokLessThan); err != nil {
		return err
	}
	if err := p.parsePrimitiveType(); err != nil {
		return err
	}
	if _, err := p.expect(TokComma); err != nil {
		return err
	}
	if err := p.parseType(); err != nil {
		return err
	}
	_, err := p.expect(TokGreaterThan)
	return err
}

func (p *parser) parseUnionType() error {
	if _, err := p.expect(KW_UNIONTYPE); err != nil {
		return err
	}
	if _, err := p.expect(TokLessThan); err != nil {
		return err
	}
	if err := p.parseType(); err != nil {
		return err
	}
	for p.match(TokComma) {
		if err := p.parseType(); err != nil {
			return err
		}
	}
	_, err := p.expect(TokGreaterThan)
	return err
}

func (p *parser) parseColumnNameColonTypeList() error {
	if err := p.parseColumnNameColonType(); err != nil {
		return err
	}
	for p.match(TokComma) {
		if err := p.parseColumnNameColonType(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseColumnNameColonType() error {
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	if _, err := p.expect(TokColon); err != nil {
		return err
	}
	if err := p.parseType(); err != nil {
		return err
	}
	if p.match(KW_COMMENT) {
		if _, err := p.expect(TokStringLiteral); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseColumnNameTypeList() error {
	if err := p.parseColumnNameType(); err != nil {
		return err
	}
	for p.match(TokComma) {
		if err := p.parseColumnNameType(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseColumnNameType() error {
	col, err := p.parseIdentifier()
	if err != nil {
		return err
	}
	if err := p.parseType(); err != nil {
		return err
	}
	if p.match(KW_COMMENT) {
		if _, err := p.expect(TokStringLiteral); err != nil {
			return err
		}
	}
	return p.checkCreateTableColName(col)
}

func (p *parser) parseColumnNameTypeOrConstraintList() error {
	if err := p.parseColumnNameTypeOrConstraint(); err != nil {
		return err
	}
	for p.match(TokComma) {
		if err := p.parseColumnNameTypeOrConstraint(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseColumnNameTypeOrConstraint() error {
	if p.try(p.parseTableConstraint) {
		return nil
	}
	return p.parseColumnNameTypeConstraint()
}

func (p *parser) parseTableConstraint() error {
	if p.try(p.parseCreateForeignKey) {
		return nil
	}
	return p.parseCreateConstraint()
}

func (p *parser) parseColumnNameTypeConstraint() error {
	col, err := p.parseIdentifier()
	if err != nil {
		return err
	}
	if err := p.parseType(); err != nil {
		return err
	}
	p.try(p.parseColumnConstraint)
	if p.match(KW_COMMENT) {
		if _, err := p.expect(TokStringLiteral); err != nil {
			return err
		}
	}
	return p.checkCreateTableColName(col)
}

func (p *parser) parseColumnConstraint() error {
	if p.try(p.parseForeignKeyConstraint) {
		return nil
	}
	return p.parseColConstraint()
}

func (p *parser) parseForeignKeyConstraint() error {
	if p.match(KW_CONSTRAINT) {
		if _, err := p.parseIdentifier(); err != nil {
			return err
		}
	}
	if _, err := p.expect(KW_REFERENCES); err != nil {
		return err
	}
	if err := p.parseTableName(); err != nil {
		return err
	}
	if _, err := p.expect(TokLParen); err != nil {
		return err
	}
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	if _, err := p.expect(TokRParen); err != nil {
		return err
	}
	p.try(p.parseConstraintOpts)
	return nil
}

func (p *parser) parseColConstraint() error {
	if p.match(KW_CONSTRAINT) {
		if _, err := p.parseIdentifier(); err != nil {
			return err
		}
	}
	if err := p.parseColumnConstraintType(); err != nil {
		return err
	}
	p.try(p.parseConstraintOpts)
	return nil
}

func (p *parser) parseColumnConstraintType() error {
	switch p.peek().Type {
	case KW_NOT:
		p.next()
		_, err := p.expect(KW_NULL)
		return err
	case KW_DEFAULT:
		p.next()
		return p.parseDefaultVal()
	case KW_CHECK:
		return p.parseCheckConstraint()
	case KW_PRIMARY, KW_UNIQUE:
		return p.parseTableConstraintType()
	default:
		return p.errAt(p.peek(), fmt.Sprintf("expected column constraint, got %q", p.peek().Str))
	}
}

func (p *parser) parseDefaultVal() error {
	if p.peek().Type == KW_CAST {
		return p.parseCastExpression()
	}
	if p.try(p.parseFunction) {
		return nil
	}
	return p.parseConstant()
}

func (p *parser) parseTableConstraintType() error {
	if p.match(KW_PRIMARY) {
		_, err := p.expect(KW_KEY)
		return err
	}
	_, err := p.expect(KW_UNIQUE)
	return err
}

func (p *parser) parseCheckConstraint() error {
	if _, err := p.expect(KW_CHECK); err != nil {
		return err
	}
	if _, err := p.expect(TokLParen); err != nil {
		return err
	}
	if err := p.parseExpression(); err != nil {
		return err
	}
	_, err := p.expect(TokRParen)
	return err
}

func (p *parser) parseCreateConstraint() error {
	if p.match(KW_CONSTRAINT) {
		if _, err := p.parseIdentifier(); err != nil {
			return err
		}
	}
	if err := p.parseTableLevelConstraint(); err != nil {
		return err
	}
	p.try(p.parseConstraintOpts)
	return nil
}

func (p *parser) parseAlterConstraintWithName() error {
	if _, err := p.expect(KW_CONSTRAINT); err != nil {
		return err
	}
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	if err := p.parseTableLevelConstraint(); err != nil {
		return err
	}
	p.try(p.parseConstraintOpts)
	return nil
}

func (p *parser) parseTableLevelConstraint() error {
	if p.peek().Type == KW_CHECK {
		return p.parseCheckConstraint()
	}
	if err := p.parseTableConstraintType(); err != nil {
		return err
	}
	return p.parseColumnParenthesesList()
}

func (p *parser) parseCreateForeignKey() error {
	if p.match(KW_CONSTRAINT) {
		if _, err := p.parseIdentifier(); err != nil {
			return err
		}
	}
	return p.parseForeignKeyBody()
}

func (p *parser) parseAlterForeignKeyWithName() error {
	if _, err := p.expect(KW_CONSTRAINT); err != nil {
		return err
	}
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	return p.parseForeignKeyBody()
}

func (p *parser) parseForeignKeyBody() error {
	if _, err := p.expect(KW_FOREIGN); err != nil {
		return err
	}
	if _, err := p.expect(KW_KEY); err != nil {
		return err
	}
	if err := p.parseColumnParenthesesList(); err != nil {
		return err
	}
	if _, err := p.expect(KW_REFERENCES); err != nil {
		return err
	}
	if err := p.parseTableName(); err != nil {
		return err
	}
	if err := p.parseColumnParenthesesList(); err != nil {
		return err
	}
	p.try(p.parseConstraintOpts)
	return nil
}

func (p *parser) parseConstraintOpts() error {
	if err := p.parseEnableValidateSpecification(); err != nil {
		return err
	}
	if p.peek().Type == KW_RELY || p.peek().Type == KW_NORELY {
		p.next()
	}
	return nil
}

func (p *parser) parseEnableValidateSpecification() error {
	switch p.peek().Type {
	case KW_ENABLE, KW_DISABLE:
		p.next()
		if p.peek().Type == KW_VALIDATE || p.peek().Type == KW_NOVALIDATE {
			p.next()
		}
		return nil
	case KW_ENFORCED:
		p.next()
		return nil
	case KW_NOT:
		if p.peekAt(1).Type == KW_ENFORCED {
			p.pos += 2
			return nil
		}
	}
	return p.errAt(p.peek(), fmt.Sprintf("expected ENABLE/DISABLE/ENFORCED/NOT ENFORCED, got %q", p.peek().Str))
}

func (p *parser) parseColumnNameList() error {
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	for p.match(TokComma) {
		if _, err := p.parseIdentifier(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseColumnParenthesesList() error {
	if _, err := p.expect(TokLParen); err != nil {
		return err
	}
	if err := p.parseColumnNameList(); err != nil {
		return err
	}
	_, err := p.expect(TokRParen)
	return err
}

func (p *parser) parseColumnNameCommentList() error {
	if err := p.parseColumnNameComment(); err != nil {
		return err
	}
	for p.match(TokComma) {
		if err := p.parseColumnNameComment(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseColumnNameComment() error {
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	if p.match(KW_COMMENT) {
		if _, err := p.expect(TokStringLiteral); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseColumnNameOrderList() error {
	if err := p.parseColumnNameOrder(); err != nil {
		return err
	}
	for p.match(TokComma) {
		if err := p.parseColumnNameOrder(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseColumnNameOrder() error {
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	if p.peek().Type == KW_ASC || p.peek().Type == KW_DESC {
		p.next()
	}
	if p.peek().Type == KW_NULLS {
		if err := p.parseNullOrdering(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseNullOrdering() error {
	if _, err := p.expect(KW_NULLS); err != nil {
		return err
	}
	if p.peek().Type == KW_FIRST || p.peek().Type == KW_LAST {
		p.next()
		return nil
	}
	return p.errAt(p.peek(), fmt.Sprintf("expected FIRST or LAST after NULLS, got %q", p.peek().Str))
}
