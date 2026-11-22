package parser

import "fmt"

// --- INSERT / DESTINATION ---

func (p *parser) parseInsertClause() error {
	if _, err := p.expect(KW_INSERT); err != nil {
		return err
	}
	if p.match(KW_OVERWRITE) {
		if err := p.parseDestination(); err != nil {
			return err
		}
		p.try(p.parseIfNotExists)
		return nil
	}
	if _, err := p.expect(KW_INTO); err != nil {
		return err
	}
	p.match(KW_TABLE)
	if err := p.parseTableOrPartition(); err != nil {
		return err
	}
	if p.peek().Type == TokLParen {
		if err := p.parseColumnParenthesesList(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseDestination() error {
	if p.match(KW_TABLE) {
		return p.parseTableOrPartition()
	}
	p.match(KW_LOCAL)
	if _, err := p.expect(KW_DIRECTORY); err != nil {
		return err
	}
	if _, err := p.expect(TokStringLiteral); err != nil {
		return err
	}
	if p.peek().Type == KW_ROW {
		if err := p.parseTableRowFormat(); err != nil {
			return err
		}
	}
	if p.peek().Type == KW_STORED {
		if err := p.parseTableFileFormat(); err != nil {
			return err
		}
	}
	return nil
}

// --- DELETE / UPDATE / MERGE ---

func (p *parser) parseDeleteStatement() error {
	if _, err := p.expect(KW_DELETE); err != nil {
		return err
	}
	if _, err := p.expect(KW_FROM); err != nil {
		return err
	}
	if err := p.parseTableName(); err != nil {
		return err
	}
	if p.peek().Type == KW_WHERE {
		return p.parseWhereClause()
	}
	return nil
}

func (p *parser) parseUpdateStatement() error {
	if _, err := p.expect(KW_UPDATE); err != nil {
		return err
	}
	if err := p.parseTableName(); err != nil {
		return err
	}
	if err := p.parseSetColumnsClause(); err != nil {
		return err
	}
	if p.peek().Type == KW_WHERE {
		return p.parseWhereClause()
	}
	return nil
}

func (p *parser) parseSetColumnsClause() error {
	if _, err := p.expect(KW_SET); err != nil {
		return err
	}
	if err := p.parseColumnAssignmentClause(); err != nil {
		return err
	}
	for p.match(TokComma) {
		if err := p.parseColumnAssignmentClause(); err != nil {
			return err
		}
	}
	return nil
}

// columnAssignmentClause : /* empty */ | tableOrColumn EQUAL precedencePlusExpressionOrDefault ;
// Note HiveParser.g:2818 has a leading `|` making columnAssignmentClause nullable!
func (p *parser) parseColumnAssignmentClause() error {
	if isIdentifier(p.peek().Type) && p.peekAt(1).Type == TokEqual {
		p.pos += 2
		return p.parsePrecedencePlusExpressionOrDefault()
	}
	return nil
}

func (p *parser) parsePrecedencePlusExpressionOrDefault() error {
	if p.peek().Type == KW_DEFAULT && p.peekAt(1).Type != TokDot {
		p.next()
		return nil
	}
	return p.parsePrecedencePlusExpression()
}

func (p *parser) parseMergeStatement() error {
	if _, err := p.expect(KW_MERGE); err != nil {
		return err
	}
	p.match(TokQueryHint)
	if _, err := p.expect(KW_INTO); err != nil {
		return err
	}
	if err := p.parseTableName(); err != nil {
		return err
	}
	if p.match(KW_AS) {
		if _, err := p.parseIdentifier(); err != nil {
			return err
		}
	} else if isIdentifier(p.peek().Type) && p.peek().Type != KW_USING {
		p.next()
	}
	if _, err := p.expect(KW_USING); err != nil {
		return err
	}
	if err := p.parseJoinSourcePart(); err != nil {
		return err
	}
	if _, err := p.expect(KW_ON); err != nil {
		return err
	}
	if err := p.parseExpression(); err != nil {
		return err
	}
	return p.parseWhenClauses()
}

func (p *parser) parseWhenClauses() error {
	count := 0
	for p.peek().Type == KW_WHEN && p.peekAt(1).Type == KW_MATCHED {
		count++
		p.pos += 2
		if p.match(KW_AND) {
			if err := p.parseExpression(); err != nil {
				return err
			}
		}
		if _, err := p.expect(KW_THEN); err != nil {
			return err
		}
		if p.match(KW_UPDATE) {
			if err := p.parseSetColumnsClause(); err != nil {
				return err
			}
		} else if !p.match(KW_DELETE) {
			return p.errAt(p.peek(), fmt.Sprintf("expected UPDATE or DELETE in MERGE WHEN MATCHED, got %q", p.peek().Str))
		}
	}
	if p.peek().Type == KW_WHEN && p.peekAt(1).Type == KW_NOT && p.peekAt(2).Type == KW_MATCHED {
		count++
		p.pos += 3
		if p.match(KW_AND) {
			if err := p.parseExpression(); err != nil {
				return err
			}
		}
		if _, err := p.expect(KW_THEN); err != nil {
			return err
		}
		if _, err := p.expect(KW_INSERT); err != nil {
			return err
		}
		if p.peek().Type == TokLParen {
			if err := p.parseColumnParenthesesList(); err != nil {
				return err
			}
		}
		if _, err := p.expect(KW_VALUES); err != nil {
			return err
		}
		return p.parseValueRowConstructor()
	}
	if count == 0 {
		return p.errAt(p.peek(), "MERGE requires at least one WHEN clause")
	}
	return nil
}

// --- LOAD / EXPORT / IMPORT / REPL ---

func (p *parser) parseLoadStatement() error {
	if _, err := p.expect(KW_LOAD); err != nil {
		return err
	}
	if _, err := p.expect(KW_DATA); err != nil {
		return err
	}
	p.match(KW_LOCAL)
	if _, err := p.expect(KW_INPATH); err != nil {
		return err
	}
	if _, err := p.expect(TokStringLiteral); err != nil {
		return err
	}
	p.match(KW_OVERWRITE)
	if _, err := p.expect(KW_INTO); err != nil {
		return err
	}
	if _, err := p.expect(KW_TABLE); err != nil {
		return err
	}
	if err := p.parseTableOrPartition(); err != nil {
		return err
	}
	if p.match(KW_INPUTFORMAT) {
		if _, err := p.expect(TokStringLiteral); err != nil {
			return err
		}
		if _, err := p.expect(KW_SERDE); err != nil {
			return err
		}
		if _, err := p.expect(TokStringLiteral); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseReplicationClause() error {
	if _, err := p.expect(KW_FOR); err != nil {
		return err
	}
	p.match(KW_METADATA)
	if _, err := p.expect(KW_REPLICATION); err != nil {
		return err
	}
	if _, err := p.expect(TokLParen); err != nil {
		return err
	}
	if _, err := p.expect(TokStringLiteral); err != nil {
		return err
	}
	_, err := p.expect(TokRParen)
	return err
}

func (p *parser) parseExportStatement() error {
	if _, err := p.expect(KW_EXPORT); err != nil {
		return err
	}
	if _, err := p.expect(KW_TABLE); err != nil {
		return err
	}
	if err := p.parseTableOrPartition(); err != nil {
		return err
	}
	if _, err := p.expect(KW_TO); err != nil {
		return err
	}
	if _, err := p.expect(TokStringLiteral); err != nil {
		return err
	}
	if p.peek().Type == KW_FOR {
		return p.parseReplicationClause()
	}
	return nil
}

func (p *parser) parseImportStatement() error {
	if _, err := p.expect(KW_IMPORT); err != nil {
		return err
	}
	if p.peek().Type == KW_EXTERNAL || p.peek().Type == KW_TABLE {
		p.match(KW_EXTERNAL)
		if _, err := p.expect(KW_TABLE); err != nil {
			return err
		}
		if err := p.parseTableOrPartition(); err != nil {
			return err
		}
	}
	if _, err := p.expect(KW_FROM); err != nil {
		return err
	}
	if _, err := p.expect(TokStringLiteral); err != nil {
		return err
	}
	if p.peek().Type == KW_LOCATION {
		return p.parseTableLocation()
	}
	return nil
}

func (p *parser) parseReplStatement() error {
	if _, err := p.expect(KW_REPL); err != nil {
		return err
	}
	switch p.peek().Type {
	case KW_DUMP:
		p.next()
		if err := p.parseReplDbPolicy(); err != nil {
			return err
		}
		if p.match(KW_REPLACE) {
			if err := p.parseReplDbPolicy(); err != nil {
				return err
			}
		}
		if p.match(KW_WITH) {
			return p.parseReplConfigs()
		}
		return nil
	case KW_LOAD:
		p.next()
		if err := p.parseReplDbPolicy(); err != nil {
			return err
		}
		if p.match(KW_INTO) {
			if _, err := p.parseIdentifier(); err != nil {
				return err
			}
		}
		if p.match(KW_WITH) {
			return p.parseReplConfigs()
		}
		return nil
	case KW_STATUS:
		p.next()
		if _, err := p.parseIdentifier(); err != nil {
			return err
		}
		if p.match(KW_WITH) {
			return p.parseReplConfigs()
		}
		return nil
	default:
		return p.errAt(p.peek(), fmt.Sprintf("expected DUMP, LOAD, or STATUS after REPL, got %q", p.peek().Str))
	}
}

func (p *parser) parseReplDbPolicy() error {
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	if p.match(TokDot) {
		if _, err := p.expect(TokStringLiteral); err != nil {
			return err
		}
		if p.match(TokDot) {
			if _, err := p.expect(TokStringLiteral); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *parser) parseReplConfigs() error {
	if _, err := p.expect(TokLParen); err != nil {
		return err
	}
	if err := p.parseKeyValueProperty(); err != nil {
		return err
	}
	for p.match(TokComma) {
		if err := p.parseKeyValueProperty(); err != nil {
			return err
		}
	}
	_, err := p.expect(TokRParen)
	return err
}

// --- TRANSACTIONS ---

func (p *parser) parseSQLTransactionStatement() error {
	switch p.peek().Type {
	case KW_START:
		p.next()
		if _, err := p.expect(KW_TRANSACTION); err != nil {
			return err
		}
		if p.peek().Type == KW_ISOLATION || p.peek().Type == KW_READ {
			if err := p.parseTransactionMode(); err != nil {
				return err
			}
			for p.match(TokComma) {
				if err := p.parseTransactionMode(); err != nil {
					return err
				}
			}
		}
		return nil
	case KW_COMMIT, KW_ROLLBACK:
		p.next()
		p.match(KW_WORK)
		return nil
	default:
		return p.errAt(p.peek(), fmt.Sprintf("unexpected transaction token %q", p.peek().Str))
	}
}

func (p *parser) parseTransactionMode() error {
	if p.match(KW_ISOLATION) {
		if _, err := p.expect(KW_LEVEL); err != nil {
			return err
		}
		_, err := p.expect(KW_SNAPSHOT)
		return err
	}
	if _, err := p.expect(KW_READ); err != nil {
		return err
	}
	if p.peek().Type == KW_ONLY || p.peek().Type == KW_WRITE {
		p.next()
		return nil
	}
	return p.errAt(p.peek(), fmt.Sprintf("expected ONLY or WRITE after READ, got %q", p.peek().Str))
}

func (p *parser) parseSetAutoCommitStatement() error {
	if _, err := p.expect(KW_SET); err != nil {
		return err
	}
	if _, err := p.expect(KW_AUTOCOMMIT); err != nil {
		return err
	}
	if p.peek().Type == KW_TRUE || p.peek().Type == KW_FALSE {
		p.next()
		return nil
	}
	return p.errAt(p.peek(), fmt.Sprintf("expected TRUE or FALSE after SET AUTOCOMMIT, got %q", p.peek().Str))
}

// --- PREPARE / EXECUTE ---

func (p *parser) parsePrepareStatement() error {
	if _, err := p.expect(KW_PREPARE); err != nil {
		return err
	}
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	if _, err := p.expect(KW_FROM); err != nil {
		return err
	}
	return p.parseQueryStatementExpression()
}

func (p *parser) parseExecuteStatement() error {
	if _, err := p.expect(KW_EXECUTE); err != nil {
		return err
	}
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	if _, err := p.expect(KW_USING); err != nil {
		return err
	}
	if err := p.parseConstant(); err != nil {
		return err
	}
	for p.match(TokComma) {
		if err := p.parseConstant(); err != nil {
			return err
		}
	}
	return nil
}
