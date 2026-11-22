package parser

import "fmt"

// queryStatementExpression : withClause? queryStatementExpressionBody ;
func (p *parser) parseQueryStatementExpression() error {
	if err := p.enter(); err != nil {
		return err
	}
	defer p.leave()

	if p.peek().Type == KW_WITH {
		if err := p.parseWithClause(); err != nil {
			return err
		}
	}
	return p.parseQueryStatementExpressionBody()
}

func (p *parser) parseWithClause() error {
	if _, err := p.expect(KW_WITH); err != nil {
		return err
	}
	if err := p.parseCTEStatement(); err != nil {
		return err
	}
	for p.match(TokComma) {
		if err := p.parseCTEStatement(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseCTEStatement() error {
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	if p.peek().Type == TokLParen {
		if err := p.parseColumnParenthesesList(); err != nil {
			return err
		}
	}
	if _, err := p.expect(KW_AS); err != nil {
		return err
	}
	if _, err := p.expect(TokLParen); err != nil {
		return err
	}
	if err := p.parseQueryStatementExpression(); err != nil {
		return err
	}
	_, err := p.expect(TokRParen)
	return err
}

func (p *parser) parseQueryStatementExpressionBody() error {
	if p.peek().Type == KW_FROM {
		return p.parseFromStatement()
	}
	return p.parseRegularBody()
}

// fromStatement : singleFromStatement (setOperator singleFromStatement)* ;
func (p *parser) parseFromStatement() error {
	if err := p.parseSingleFromStatement(); err != nil {
		return err
	}
	for p.isSetOperatorStart() {
		if err := p.parseSetOperator(); err != nil {
			return err
		}
		if err := p.parseSingleFromStatement(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseSingleFromStatement() error {
	if err := p.parseFromClause(); err != nil {
		return err
	}
	if err := p.parseBody(); err != nil {
		return err
	}
	for p.isBodyStart() {
		if err := p.parseBody(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) isBodyStart() bool {
	switch p.peek().Type {
	case KW_INSERT, KW_SELECT, KW_MAP, KW_REDUCE:
		return true
	}
	return false
}

func (p *parser) parseBody() error {
	if p.peek().Type == KW_INSERT {
		if err := p.parseInsertClause(); err != nil {
			return err
		}
	}
	if err := p.parseSelectClause(); err != nil {
		return err
	}
	p.try(p.parseLateralView)
	if p.peek().Type == KW_WHERE {
		if err := p.parseWhereClause(); err != nil {
			return err
		}
	}
	if p.peek().Type == KW_GROUP {
		if err := p.parseGroupByClause(); err != nil {
			return err
		}
	}
	if p.peek().Type == KW_HAVING {
		if err := p.parseHavingClause(); err != nil {
			return err
		}
	}
	if p.peek().Type == KW_WINDOW {
		if err := p.parseWindowClause(); err != nil {
			return err
		}
	}
	if p.peek().Type == KW_QUALIFY {
		if err := p.parseQualifyClause(); err != nil {
			return err
		}
	}
	if p.peek().Type == KW_ORDER {
		if err := p.parseOrderByClause(); err != nil {
			return err
		}
	}
	if p.peek().Type == KW_CLUSTER {
		if err := p.parseClusterByClause(); err != nil {
			return err
		}
	}
	if p.peek().Type == KW_DISTRIBUTE {
		if err := p.parseDistributeByClause(); err != nil {
			return err
		}
	}
	if p.peek().Type == KW_SORT {
		if err := p.parseSortByClause(); err != nil {
			return err
		}
	}
	if p.peek().Type == KW_LIMIT {
		if err := p.parseLimitClause(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseRegularBody() error {
	if p.peek().Type == KW_INSERT {
		if err := p.parseInsertClause(); err != nil {
			return err
		}
	}
	return p.parseSelectStatement()
}

func (p *parser) parseSelectStatementWithCTE() error {
	if p.peek().Type == KW_WITH {
		if err := p.parseWithClause(); err != nil {
			return err
		}
	}
	return p.parseSelectStatement()
}

func (p *parser) parseSelectStatement() error {
	if err := p.enter(); err != nil {
		return err
	}
	defer p.leave()

	if err := p.parseAtomSelectStatement(); err != nil {
		return err
	}
	for p.isSetOperatorStart() {
		if err := p.parseSetOperator(); err != nil {
			return err
		}
		if err := p.parseAtomSelectStatement(); err != nil {
			return err
		}
	}
	if p.peek().Type == KW_ORDER {
		if err := p.parseOrderByClause(); err != nil {
			return err
		}
	}
	if p.peek().Type == KW_CLUSTER {
		if err := p.parseClusterByClause(); err != nil {
			return err
		}
	}
	if p.peek().Type == KW_DISTRIBUTE {
		if err := p.parseDistributeByClause(); err != nil {
			return err
		}
	}
	if p.peek().Type == KW_SORT {
		if err := p.parseSortByClause(); err != nil {
			return err
		}
	}
	if p.peek().Type == KW_LIMIT {
		if err := p.parseLimitClause(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseAtomSelectStatement() error {
	if p.match(TokLParen) {
		if err := p.parseSelectStatement(); err != nil {
			return err
		}
		_, err := p.expect(TokRParen)
		return err
	}
	if p.peek().Type == KW_VALUES {
		return p.parseValuesClause()
	}
	if err := p.parseSelectClause(); err != nil {
		return err
	}
	if p.peek().Type == KW_FROM {
		if err := p.parseFromClause(); err != nil {
			return err
		}
	}
	if p.peek().Type == KW_WHERE {
		if err := p.parseWhereClause(); err != nil {
			return err
		}
	}
	if p.peek().Type == KW_GROUP {
		if err := p.parseGroupByClause(); err != nil {
			return err
		}
	}
	if p.peek().Type == KW_HAVING {
		if err := p.parseHavingClause(); err != nil {
			return err
		}
	}
	if p.peek().Type == KW_WINDOW {
		if err := p.parseWindowClause(); err != nil {
			return err
		}
	}
	if p.peek().Type == KW_QUALIFY {
		if err := p.parseQualifyClause(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) isSetOperatorStart() bool {
	switch p.peek().Type {
	case KW_UNION, KW_INTERSECT, KW_EXCEPT, KW_MINUS:
		return true
	}
	return false
}

func (p *parser) parseSetOperator() error {
	switch p.peek().Type {
	case KW_UNION, KW_INTERSECT, KW_EXCEPT, KW_MINUS:
		p.next()
		if p.peek().Type == KW_ALL || p.peek().Type == KW_DISTINCT {
			p.next()
		}
		return nil
	default:
		return p.errAt(p.peek(), fmt.Sprintf("expected set operator, got %q", p.peek().Str))
	}
}

func (p *parser) parseSelectClause() error {
	if p.peek().Type == KW_MAP || p.peek().Type == KW_REDUCE {
		return p.parseTrfmClause()
	}
	if _, err := p.expect(KW_SELECT); err != nil {
		return err
	}
	p.match(TokQueryHint)
	if p.match(KW_TRANSFORM) {
		return p.parseSelectTrfmClause()
	}
	if p.peek().Type == KW_ALL || p.peek().Type == KW_DISTINCT {
		p.next()
	}
	return p.parseSelectList()
}

func (p *parser) parseSelectList() error {
	if err := p.parseSelectItem(); err != nil {
		return err
	}
	for p.match(TokComma) {
		if err := p.parseSelectItem(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseSelectItem() error {
	if p.try(p.parseTableAllColumns) {
		return nil
	}
	if err := p.parseExpression(); err != nil {
		return err
	}
	if p.match(KW_AS) {
		if p.match(TokLParen) {
			if err := p.parseAliasList(); err != nil {
				return err
			}
			_, err := p.expect(TokRParen)
			return err
		}
		_, err := p.parseIdentifier()
		return err
	}
	if isIdentifier(p.peek().Type) && !p.isSelectItemAliasBoundary() {
		p.next()
	}
	return nil
}

func (p *parser) isSelectItemAliasBoundary() bool {
	switch p.peek().Type {
	case KW_LIMIT:
		return p.peekAt(1).Type == TokNumber
	case KW_CLUSTER, KW_DISTRIBUTE, KW_SORT:
		return p.peekAt(1).Type == KW_BY
	case KW_LATERAL:
		return p.peekAt(1).Type == KW_VIEW || p.peekAt(1).Type == KW_TABLE
	case KW_WINDOW:
		return isIdentifier(p.peekAt(1).Type) && p.peekAt(2).Type == KW_AS
	case KW_QUALIFY:
		return p.peekAt(1).Type != TokEOF && p.peekAt(1).Type != TokComma && p.peekAt(1).Type != TokRParen
	}
	return false
}

func (p *parser) parseSelectTrfmClause() error {
	if _, err := p.expect(TokLParen); err != nil {
		return err
	}
	if err := p.parseSelectExpressionList(); err != nil {
		return err
	}
	if _, err := p.expect(TokRParen); err != nil {
		return err
	}
	return p.parseTrfmTail()
}

func (p *parser) parseTrfmClause() error {
	if p.peek().Type != KW_MAP && p.peek().Type != KW_REDUCE {
		return p.errAt(p.peek(), fmt.Sprintf("expected MAP or REDUCE, got %q", p.peek().Str))
	}
	p.next()
	if err := p.parseSelectExpressionList(); err != nil {
		return err
	}
	return p.parseTrfmTail()
}

func (p *parser) parseTrfmTail() error {
	if err := p.parseOptionalRowFormat(); err != nil {
		return err
	}
	if p.match(KW_RECORDWRITER) {
		if _, err := p.expect(TokStringLiteral); err != nil {
			return err
		}
	}
	if _, err := p.expect(KW_USING); err != nil {
		return err
	}
	if _, err := p.expect(TokStringLiteral); err != nil {
		return err
	}
	if p.match(KW_AS) {
		if p.match(TokLParen) {
			if !p.try(p.parseColumnNameTypeList) {
				if err := p.parseAliasList(); err != nil {
					return err
				}
			}
			if _, err := p.expect(TokRParen); err != nil {
				return err
			}
		} else {
			if !p.try(p.parseColumnNameTypeList) {
				if err := p.parseAliasList(); err != nil {
					return err
				}
			}
		}
	}
	if err := p.parseOptionalRowFormat(); err != nil {
		return err
	}
	if p.match(KW_RECORDREADER) {
		if _, err := p.expect(TokStringLiteral); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseAliasList() error {
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

// --- FROM clause ---

func (p *parser) parseFromClause() error {
	if _, err := p.expect(KW_FROM); err != nil {
		return err
	}
	return p.parseFromSource()
}

func (p *parser) parseFromSource() error {
	if p.match(KW_UNIQUEJOIN) {
		if err := p.parseUniqueJoinSource(); err != nil {
			return err
		}
		if _, err := p.expect(TokComma); err != nil {
			return err
		}
		if err := p.parseUniqueJoinSource(); err != nil {
			return err
		}
		for p.match(TokComma) {
			if err := p.parseUniqueJoinSource(); err != nil {
				return err
			}
		}
		return nil
	}
	return p.parseJoinSource()
}

func (p *parser) parseUniqueJoinSource() error {
	p.match(KW_PRESERVE)
	if err := p.parseTableName(); err != nil {
		return err
	}
	if p.peek().Type == KW_TABLESAMPLE {
		if err := p.parseTableSample(); err != nil {
			return err
		}
	}
	if p.match(KW_AS) {
		if _, err := p.parseIdentifier(); err != nil {
			return err
		}
	} else if isIdentifier(p.peek().Type) {
		p.next()
	}
	if _, err := p.expect(TokLParen); err != nil {
		return err
	}
	if err := p.parseExpressionList(); err != nil {
		return err
	}
	_, err := p.expect(TokRParen)
	return err
}

func (p *parser) parseJoinSource() error {
	if err := p.parseAtomJoinSource(); err != nil {
		return err
	}
	for {
		isComma, ok := p.parseJoinToken()
		if !ok {
			break
		}
		if err := p.parseJoinSourcePart(); err != nil {
			return err
		}
		if !isComma {
			if p.match(KW_ON) {
				if err := p.parseExpression(); err != nil {
					return err
				}
			} else if p.match(KW_USING) {
				if err := p.parseColumnParenthesesList(); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (p *parser) parseJoinToken() (isComma bool, ok bool) {
	saved := p.save()
	switch p.peek().Type {
	case KW_JOIN:
		p.next()
		return false, true
	case TokComma:
		// Do not swallow comma if it belongs to COMMA LATERAL VIEW/TABLE.
		if p.peekAt(1).Type == KW_LATERAL {
			return false, false
		}
		p.next()
		return true, true
	case KW_INNER, KW_CROSS:
		p.next()
		if p.match(KW_JOIN) {
			return false, true
		}
	case KW_LEFT:
		p.next()
		if p.match(KW_SEMI) || p.match(KW_OUTER) {
		}
		if p.match(KW_JOIN) {
			return false, true
		}
	case KW_RIGHT, KW_FULL:
		p.next()
		p.match(KW_OUTER)
		if p.match(KW_JOIN) {
			return false, true
		}
	}
	p.restore(saved)
	return false, false
}

func (p *parser) parseAtomJoinSource() error {
	if p.peek().Type == KW_TABLE && p.peekAt(1).Type == TokLParen {
		if err := p.parseVirtualTableSource(); err != nil {
			return err
		}
		for p.try(p.parseLateralView) {
		}
		return nil
	}
	if p.peek().Type == TokLParen {
		if p.try(p.parseSubQuerySource) {
			for p.try(p.parseLateralView) {
			}
			return nil
		}
		p.next() // '('
		if err := p.parseJoinSource(); err != nil {
			return err
		}
		_, err := p.expect(TokRParen)
		return err
	}
	if isIdentifier(p.peek().Type) && p.peekAt(1).Type == TokLParen && p.peekAt(2).Type == KW_ON {
		if err := p.parsePartitionedTableFunction(); err != nil {
			return err
		}
		for p.try(p.parseLateralView) {
		}
		return nil
	}
	if err := p.parseTableSource(); err != nil {
		return err
	}
	for p.try(p.parseLateralView) {
	}
	return nil
}

func (p *parser) parseJoinSourcePart() error {
	if p.peek().Type == KW_TABLE && p.peekAt(1).Type == TokLParen {
		if err := p.parseVirtualTableSource(); err != nil {
			return err
		}
	} else if p.peek().Type == TokLParen {
		if err := p.parseSubQuerySource(); err != nil {
			return err
		}
	} else if isIdentifier(p.peek().Type) && p.peekAt(1).Type == TokLParen && p.peekAt(2).Type == KW_ON {
		if err := p.parsePartitionedTableFunction(); err != nil {
			return err
		}
	} else {
		if err := p.parseTableSource(); err != nil {
			return err
		}
	}
	for p.try(p.parseLateralView) {
	}
	return nil
}

func (p *parser) parseTableSource() error {
	if err := p.parseTableName(); err != nil {
		return err
	}
	if p.peek().Type == TokLParen {
		if err := p.parseTableProperties(); err != nil {
			return err
		}
	}
	if p.peek().Type == KW_TABLESAMPLE {
		if err := p.parseTableSample(); err != nil {
			return err
		}
	}
	if p.peek().Type == KW_FOR && (p.peekAt(1).Type == KW_SYSTEM_TIME || p.peekAt(1).Type == KW_SYSTEM_VERSION) {
		if err := p.parseAsOfClause(); err != nil {
			return err
		}
	}
	if p.match(KW_AS) {
		if _, err := p.parseIdentifier(); err != nil {
			return err
		}
	} else if isIdentifier(p.peek().Type) && !p.isClauseBoundaryAfterTableSource() && p.peekAt(1).Type != TokLParen {
		p.next()
	}
	return nil
}

func (p *parser) isClauseBoundaryAfterTableSource() bool {
	switch p.peek().Type {
	case KW_LATERAL, KW_WHERE, KW_GROUP, KW_HAVING, KW_WINDOW, KW_QUALIFY,
		KW_ORDER, KW_CLUSTER, KW_DISTRIBUTE, KW_SORT, KW_LIMIT,
		KW_JOIN, KW_INNER, KW_LEFT, KW_RIGHT, KW_FULL, KW_CROSS,
		KW_ON, KW_USING, KW_UNION, KW_INTERSECT, KW_EXCEPT, KW_MINUS,
		KW_INSERT, KW_SELECT, KW_MAP, KW_REDUCE, KW_SET, KW_WHEN:
		return true
	}
	return false
}

func (p *parser) parseAsOfClause() error {
	if _, err := p.expect(KW_FOR); err != nil {
		return err
	}
	if p.peek().Type != KW_SYSTEM_TIME && p.peek().Type != KW_SYSTEM_VERSION {
		return p.errAt(p.peek(), fmt.Sprintf("expected SYSTEM_TIME or SYSTEM_VERSION, got %q", p.peek().Str))
	}
	p.next()
	if _, err := p.expect(KW_AS); err != nil {
		return err
	}
	if _, err := p.expect(KW_OF); err != nil {
		return err
	}
	return p.parseExpression()
}

func (p *parser) parseTableSample() error {
	if _, err := p.expect(KW_TABLESAMPLE); err != nil {
		return err
	}
	if _, err := p.expect(TokLParen); err != nil {
		return err
	}
	if p.match(KW_BUCKET) {
		if _, err := p.expect(TokNumber); err != nil {
			return err
		}
		if _, err := p.expect(KW_OUT); err != nil {
			return err
		}
		if _, err := p.expect(KW_OF); err != nil {
			return err
		}
		if _, err := p.expect(TokNumber); err != nil {
			return err
		}
		if p.match(KW_ON) {
			if err := p.parseExpressionList(); err != nil {
				return err
			}
		}
		_, err := p.expect(TokRParen)
		return err
	}
	if p.match(TokByteLengthLiteral) {
		_, err := p.expect(TokRParen)
		return err
	}
	if _, err := p.expect(TokNumber); err != nil {
		return err
	}
	if p.peek().Type != KW_PERCENT && p.peek().Type != KW_ROWS {
		return p.errAt(p.peek(), fmt.Sprintf("expected PERCENT or ROWS in TABLESAMPLE, got %q", p.peek().Str))
	}
	p.next()
	_, err := p.expect(TokRParen)
	return err
}

func (p *parser) parseSubQuerySource() error {
	if _, err := p.expect(TokLParen); err != nil {
		return err
	}
	if err := p.parseQueryStatementExpression(); err != nil {
		return err
	}
	if _, err := p.expect(TokRParen); err != nil {
		return err
	}
	p.match(KW_AS)
	_, err := p.parseIdentifier()
	return err
}

func (p *parser) parseVirtualTableSource() error {
	if _, err := p.expect(KW_TABLE); err != nil {
		return err
	}
	if _, err := p.expect(TokLParen); err != nil {
		return err
	}
	if err := p.parseValuesClause(); err != nil {
		return err
	}
	if _, err := p.expect(TokRParen); err != nil {
		return err
	}
	p.match(KW_AS)
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	if p.match(TokLParen) {
		if err := p.parseAliasList(); err != nil {
			return err
		}
	}
	_, err := p.expect(TokRParen)
	return err
}

func (p *parser) parseLateralView() error {
	saved := p.save()
	p.match(TokComma)
	if !p.match(KW_LATERAL) {
		p.restore(saved)
		return p.errAt(p.peek(), "expected LATERAL")
	}
	if p.match(KW_TABLE) {
		if _, err := p.expect(TokLParen); err != nil {
			p.restore(saved)
			return err
		}
		if err := p.parseValuesClause(); err != nil {
			p.restore(saved)
			return err
		}
		if _, err := p.expect(TokRParen); err != nil {
			p.restore(saved)
			return err
		}
		p.match(KW_AS)
		if _, err := p.parseIdentifier(); err != nil {
			p.restore(saved)
			return err
		}
		if p.match(TokLParen) {
			if err := p.parseAliasList(); err != nil {
				p.restore(saved)
				return err
			}
			if _, err := p.expect(TokRParen); err != nil {
				p.restore(saved)
				return err
			}
		}
		return nil
	}
	if _, err := p.expect(KW_VIEW); err != nil {
		p.restore(saved)
		return err
	}
	p.match(KW_OUTER)
	if err := p.parseFunction(); err != nil {
		p.restore(saved)
		return err
	}
	if _, err := p.parseIdentifier(); err != nil {
		p.restore(saved)
		return err
	}
	if p.match(KW_AS) {
		if _, err := p.parseIdentifier(); err != nil {
			p.restore(saved)
			return err
		}
		for p.peek().Type == TokComma && isIdentifier(p.peekAt(1).Type) && p.peekAt(1).Type != KW_LATERAL {
			// Be careful: COMMA could also start the next table in FROM a, b!
			// In FromClauseParser.g:165, ((COMMA)=> COMMA identifier)* is greedy.
			p.next()
			if _, err := p.parseIdentifier(); err != nil {
				p.restore(saved)
				return err
			}
		}
	}
	return nil
}

func (p *parser) parsePartitionedTableFunction() error {
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	if _, err := p.expect(TokLParen); err != nil {
		return err
	}
	if _, err := p.expect(KW_ON); err != nil {
		return err
	}
	// partitionTableFunctionSource: subQuerySource | tableSource | partitionedTableFunction
	if p.peek().Type == TokLParen {
		if err := p.parseSubQuerySource(); err != nil {
			return err
		}
	} else if isIdentifier(p.peek().Type) && p.peekAt(1).Type == TokLParen && p.peekAt(2).Type == KW_ON {
		if err := p.parsePartitionedTableFunction(); err != nil {
			return err
		}
	} else {
		if err := p.parseTableSource(); err != nil {
			return err
		}
	}
	p.try(p.parsePartitioningSpec)
	// Optional PTF args: Identifier '(' expression ')' (',' Identifier '(' expression ')')*
	if p.peek().Type == TokIdentifier && p.peekAt(1).Type == TokLParen {
		for {
			if _, err := p.expect(TokIdentifier); err != nil {
				return err
			}
			if _, err := p.expect(TokLParen); err != nil {
				return err
			}
			if err := p.parseExpression(); err != nil {
				return err
			}
			if _, err := p.expect(TokRParen); err != nil {
				return err
			}
			if !p.match(TokComma) {
				break
			}
		}
	}
	if _, err := p.expect(TokRParen); err != nil {
		return err
	}
	if p.peek().Type == TokIdentifier {
		p.next()
	}
	return nil
}

func (p *parser) parsePartitioningSpec() error {
	switch p.peek().Type {
	case KW_PARTITION:
		if err := p.parsePartitionByClause(); err != nil {
			return err
		}
		if p.peek().Type == KW_ORDER {
			return p.parseOrderByClause()
		}
		return nil
	case KW_ORDER:
		return p.parseOrderByClause()
	case KW_DISTRIBUTE:
		if err := p.parseDistributeByClause(); err != nil {
			return err
		}
		if p.peek().Type == KW_SORT {
			return p.parseSortByClause()
		}
		return nil
	case KW_SORT:
		return p.parseSortByClause()
	case KW_CLUSTER:
		return p.parseClusterByClause()
	default:
		return p.errAt(p.peek(), "expected partitioning spec")
	}
}

// --- VALUES ---

func (p *parser) parseValuesClause() error {
	if _, err := p.expect(KW_VALUES); err != nil {
		return err
	}
	return p.parseValuesTableConstructor()
}

func (p *parser) parseValuesTableConstructor() error {
	if p.try(func() error {
		if err := p.parseValueRowConstructor(); err != nil {
			return err
		}
		for p.match(TokComma) {
			if err := p.parseValueRowConstructor(); err != nil {
				return err
			}
		}
		return nil
	}) {
		return nil
	}
	if err := p.parseFirstValueRowConstructor(); err != nil {
		return err
	}
	for p.match(TokComma) {
		if err := p.parseValueRowConstructor(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseValueRowConstructor() error {
	return p.parseExpressionsInParenthesis()
}

func (p *parser) parseFirstValueRowConstructor() error {
	if _, err := p.expect(TokLParen); err != nil {
		return err
	}
	if err := p.parseExpressionWithAlias(); err != nil {
		return err
	}
	for p.match(TokComma) {
		if err := p.parseExpressionWithAlias(); err != nil {
			return err
		}
	}
	_, err := p.expect(TokRParen)
	return err
}

func (p *parser) parseExpressionWithAlias() error {
	if err := p.parseExpression(); err != nil {
		return err
	}
	if p.match(KW_AS) {
		if isIdentifier(p.peek().Type) {
			p.next()
		}
	} else if isIdentifier(p.peek().Type) {
		p.next()
	}
	return nil
}

// --- WHERE / GROUP BY / HAVING / WINDOW / QUALIFY / ORDER BY / CLUSTER / DISTRIBUTE / SORT / LIMIT ---

func (p *parser) parseWhereClause() error {
	if _, err := p.expect(KW_WHERE); err != nil {
		return err
	}
	return p.parseExpression()
}

func (p *parser) parseGroupByClause() error {
	if _, err := p.expect(KW_GROUP); err != nil {
		return err
	}
	if _, err := p.expect(KW_BY); err != nil {
		return err
	}
	if p.peek().Type == TokLParen && p.peekAt(1).Type == TokRParen {
		p.pos += 2
		return nil
	}
	if p.peek().Type == KW_ROLLUP || p.peek().Type == KW_CUBE {
		p.next()
		if _, err := p.expect(TokLParen); err != nil {
			return err
		}
		if err := p.parseExpressionList(); err != nil {
			return err
		}
		_, err := p.expect(TokRParen)
		return err
	}
	if err := p.parseExpressionsNotInParenthesis(); err != nil {
		return err
	}
	if p.peek().Type == KW_WITH && (p.peekAt(1).Type == KW_ROLLUP || p.peekAt(1).Type == KW_CUBE) {
		p.pos += 2
	}
	if p.peek().Type == KW_GROUPING && p.peekAt(1).Type == KW_SETS {
		p.pos += 2
		if _, err := p.expect(TokLParen); err != nil {
			return err
		}
		if err := p.parseGroupingSetExpression(); err != nil {
			return err
		}
		for p.match(TokComma) {
			if err := p.parseGroupingSetExpression(); err != nil {
				return err
			}
		}
		_, err := p.expect(TokRParen)
		return err
	}
	return nil
}

func (p *parser) parseGroupingSetExpression() error {
	if p.try(func() error {
		if _, err := p.expect(TokLParen); err != nil {
			return err
		}
		if p.peek().Type != TokRParen && p.peek().Type != TokComma {
			if err := p.parseExpression(); err != nil {
				return err
			}
		}
		for p.match(TokComma) {
			if err := p.parseExpression(); err != nil {
				return err
			}
		}
		_, err := p.expect(TokRParen)
		return err
	}) {
		return nil
	}
	return p.parseExpression()
}

func (p *parser) parseHavingClause() error {
	if _, err := p.expect(KW_HAVING); err != nil {
		return err
	}
	return p.parseExpression()
}

func (p *parser) parseQualifyClause() error {
	if _, err := p.expect(KW_QUALIFY); err != nil {
		return err
	}
	return p.parseExpression()
}

func (p *parser) parseWindowClause() error {
	if _, err := p.expect(KW_WINDOW); err != nil {
		return err
	}
	if err := p.parseWindowDefn(); err != nil {
		return err
	}
	for p.match(TokComma) {
		if err := p.parseWindowDefn(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseWindowDefn() error {
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	if _, err := p.expect(KW_AS); err != nil {
		return err
	}
	return p.parseWindowSpecification()
}

func (p *parser) parseWindowSpecification() error {
	if p.match(TokLParen) {
		// Optional window name identifier (not starting partitioningSpec or window_frame).
		if isIdentifier(p.peek().Type) && !p.isWindowSpecClauseStart() {
			p.next()
		}
		p.try(p.parsePartitioningSpec)
		if p.peek().Type == KW_ROWS || p.peek().Type == KW_RANGE {
			if err := p.parseWindowFrame(); err != nil {
				return err
			}
		}
		_, err := p.expect(TokRParen)
		return err
	}
	_, err := p.parseIdentifier()
	return err
}

func (p *parser) isWindowSpecClauseStart() bool {
	switch p.peek().Type {
	case KW_PARTITION, KW_ORDER, KW_DISTRIBUTE, KW_SORT, KW_CLUSTER, KW_ROWS, KW_RANGE:
		return true
	}
	return false
}

func (p *parser) parseWindowFrame() error {
	if p.peek().Type != KW_ROWS && p.peek().Type != KW_RANGE {
		return p.errAt(p.peek(), "expected ROWS or RANGE in window frame")
	}
	p.next()
	if p.match(KW_BETWEEN) {
		if err := p.parseWindowFrameBoundary(); err != nil {
			return err
		}
		if _, err := p.expect(KW_AND); err != nil {
			return err
		}
		return p.parseWindowFrameBoundary()
	}
	return p.parseWindowFrameStartBoundary()
}

func (p *parser) parseWindowFrameStartBoundary() error {
	switch p.peek().Type {
	case KW_UNBOUNDED:
		p.next()
		_, err := p.expect(KW_PRECEDING)
		return err
	case KW_CURRENT:
		p.next()
		_, err := p.expect(KW_ROW)
		return err
	case TokNumber:
		p.next()
		_, err := p.expect(KW_PRECEDING)
		return err
	default:
		return p.errAt(p.peek(), fmt.Sprintf("expected window frame boundary, got %q", p.peek().Str))
	}
}

func (p *parser) parseWindowFrameBoundary() error {
	switch p.peek().Type {
	case KW_UNBOUNDED, TokNumber:
		p.next()
		if p.peek().Type == KW_PRECEDING || p.peek().Type == KW_FOLLOWING {
			p.next()
			return nil
		}
		return p.errAt(p.peek(), fmt.Sprintf("expected PRECEDING or FOLLOWING, got %q", p.peek().Str))
	case KW_CURRENT:
		p.next()
		_, err := p.expect(KW_ROW)
		return err
	default:
		return p.errAt(p.peek(), fmt.Sprintf("expected window frame boundary, got %q", p.peek().Str))
	}
}

func (p *parser) parseColumnRefOrder() error {
	if err := p.parseExpression(); err != nil {
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

func (p *parser) parseOrderByClause() error {
	if _, err := p.expect(KW_ORDER); err != nil {
		return err
	}
	if _, err := p.expect(KW_BY); err != nil {
		return err
	}
	if err := p.parseColumnRefOrder(); err != nil {
		return err
	}
	for p.match(TokComma) {
		if err := p.parseColumnRefOrder(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseClusterByClause() error {
	if _, err := p.expect(KW_CLUSTER); err != nil {
		return err
	}
	if _, err := p.expect(KW_BY); err != nil {
		return err
	}
	return p.parseExpressions()
}

func (p *parser) parsePartitionByClause() error {
	if _, err := p.expect(KW_PARTITION); err != nil {
		return err
	}
	if _, err := p.expect(KW_BY); err != nil {
		return err
	}
	return p.parseExpressions()
}

func (p *parser) parseDistributeByClause() error {
	if _, err := p.expect(KW_DISTRIBUTE); err != nil {
		return err
	}
	if _, err := p.expect(KW_BY); err != nil {
		return err
	}
	return p.parseExpressions()
}

func (p *parser) parseSortByClause() error {
	if _, err := p.expect(KW_SORT); err != nil {
		return err
	}
	if _, err := p.expect(KW_BY); err != nil {
		return err
	}
	if p.peek().Type == TokLParen {
		if p.try(func() error {
			if _, err := p.expect(TokLParen); err != nil {
				return err
			}
			if err := p.parseColumnRefOrder(); err != nil {
				return err
			}
			for p.match(TokComma) {
				if err := p.parseColumnRefOrder(); err != nil {
					return err
				}
			}
			_, err := p.expect(TokRParen)
			return err
		}) {
			return nil
		}
	}
	if err := p.parseColumnRefOrder(); err != nil {
		return err
	}
	for p.match(TokComma) {
		if err := p.parseColumnRefOrder(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseLimitClause() error {
	if _, err := p.expect(KW_LIMIT); err != nil {
		return err
	}
	if _, err := p.expect(TokNumber); err != nil {
		return err
	}
	if p.match(TokComma) {
		_, err := p.expect(TokNumber)
		return err
	}
	if p.match(KW_OFFSET) {
		_, err := p.expect(TokNumber)
		return err
	}
	return nil
}
