package parser

import "fmt"

// parseAlterDDL handles AlterClauseParser.g + ALTER SCHEDULED QUERY +
// ResourcePlanParser.g ALTER statements.
func (p *parser) parseAlterDDL() error {
	if p.match(KW_OPTIMIZE) {
		if _, err := p.expect(KW_TABLE); err != nil {
			return err
		}
		if err := p.parseTableName(); err != nil {
			return err
		}
		if _, err := p.expect(KW_REWRITE); err != nil {
			return err
		}
		_, err := p.expect(KW_DATA)
		return err
	}
	if _, err := p.expect(KW_ALTER); err != nil {
		return err
	}
	switch p.peek().Type {
	case KW_TABLE:
		p.next()
		if err := p.parseTableName(); err != nil {
			return err
		}
		return p.parseAlterTableStatementSuffix()
	case KW_VIEW:
		p.next()
		if err := p.parseTableName(); err != nil {
			return err
		}
		p.match(KW_AS)
		return p.parseAlterViewStatementSuffix()
	case KW_MATERIALIZED:
		p.next()
		if _, err := p.expect(KW_VIEW); err != nil {
			return err
		}
		if err := p.parseTableName(); err != nil {
			return err
		}
		if p.match(KW_REBUILD) {
			return nil
		}
		if p.peek().Type == KW_ENABLE || p.peek().Type == KW_DISABLE {
			p.next()
			_, err := p.expect(KW_REWRITE)
			return err
		}
		return p.errAt(p.peek(), fmt.Sprintf("expected REBUILD or ENABLE/DISABLE REWRITE, got %q", p.peek().Str))
	case KW_DATABASE, KW_SCHEMA:
		p.next()
		return p.parseAlterDatabaseStatementSuffix()
	case KW_DATACONNECTOR:
		p.next()
		return p.parseAlterDataConnectorStatementSuffix()
	case KW_SCHEDULED:
		p.next()
		if _, err := p.expect(KW_QUERY); err != nil {
			return err
		}
		if _, err := p.parseIdentifier(); err != nil {
			return err
		}
		return p.parseAlterScheduledQueryChange()
	case KW_RESOURCE:
		return p.parseAlterResourcePlanStatement()
	case KW_TRIGGER:
		return p.parseAlterTriggerStatement()
	case KW_POOL:
		return p.parseAlterPoolStatement()
	case KW_USER, KW_GROUP, KW_APPLICATION:
		return p.parseAlterMappingStatement()
	default:
		return p.errAt(p.peek(), fmt.Sprintf("unexpected token %q after ALTER", p.peek().Str))
	}
}

func (p *parser) parseAlterScheduledQueryChange() error {
	switch p.peek().Type {
	case KW_CRON, KW_EVERY:
		return p.parseScheduleSpec()
	case KW_EXECUTED:
		return p.parseExecutedAsSpec()
	case KW_ENABLE, KW_DISABLE:
		p.next()
		return nil
	case KW_DEFINED, KW_AS:
		return p.parseDefinedAsSpec()
	case KW_EXECUTE:
		p.next()
		return nil
	default:
		return p.errAt(p.peek(), fmt.Sprintf("unexpected token %q in ALTER SCHEDULED QUERY", p.peek().Str))
	}
}

func (p *parser) parseAlterDatabaseStatementSuffix() error {
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	if _, err := p.expect(KW_SET); err != nil {
		return err
	}
	switch p.peek().Type {
	case KW_DBPROPERTIES:
		p.next()
		return p.parseDBProperties()
	case KW_OWNER:
		p.next()
		return p.parsePrincipalName()
	case KW_LOCATION, KW_MANAGEDLOCATION:
		p.next()
		_, err := p.expect(TokStringLiteral)
		return err
	default:
		return p.errAt(p.peek(), fmt.Sprintf("unexpected token %q after ALTER DATABASE SET", p.peek().Str))
	}
}

func (p *parser) parseAlterDataConnectorStatementSuffix() error {
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	if _, err := p.expect(KW_SET); err != nil {
		return err
	}
	switch p.peek().Type {
	case KW_DCPROPERTIES:
		p.next()
		return p.parseDBProperties()
	case KW_OWNER:
		p.next()
		return p.parsePrincipalName()
	case KW_URL:
		p.next()
		_, err := p.expect(TokStringLiteral)
		return err
	default:
		return p.errAt(p.peek(), fmt.Sprintf("unexpected token %q after ALTER CONNECTOR SET", p.peek().Str))
	}
}

func (p *parser) parseAlterViewStatementSuffix() error {
	switch p.peek().Type {
	case KW_SET:
		p.next()
		if _, err := p.expect(KW_TBLPROPERTIES); err != nil {
			return err
		}
		return p.parseTableProperties()
	case KW_UNSET:
		p.next()
		if _, err := p.expect(KW_TBLPROPERTIES); err != nil {
			return err
		}
		p.try(p.parseIfExists)
		return p.parseTableProperties()
	case KW_RENAME:
		p.next()
		if _, err := p.expect(KW_TO); err != nil {
			return err
		}
		return p.parseTableName()
	case KW_ADD:
		return p.parseAlterAddPartitions()
	case KW_DROP:
		return p.parseAlterDropPartitions()
	default:
		return p.parseSelectStatementWithCTE()
	}
}

func (p *parser) parseAlterTableStatementSuffix() error {
	switch p.peek().Type {
	case KW_RENAME:
		// RENAME TO tableName | RENAME TO partitionSpec
		if p.try(func() error {
			p.next()
			if _, err := p.expect(KW_TO); err != nil {
				return err
			}
			return p.parseTableName()
		}) {
			return nil
		}
		return p.parseAlterTblPartitionStatementSuffix()
	case KW_DROP:
		switch p.peekAt(1).Type {
		case KW_CONSTRAINT:
			p.pos += 2
			_, err := p.parseIdentifier()
			return err
		case KW_BRANCH, KW_TAG:
			p.pos += 2
			p.try(p.parseIfExists)
			_, err := p.parseIdentifier()
			return err
		default:
			return p.parseAlterDropPartitions()
		}
	case KW_ADD:
		if p.peekAt(1).Type == KW_COLUMNS {
			return p.parseAlterAddOrReplaceCols()
		}
		if p.peekAt(1).Type == KW_CONSTRAINT {
			p.next()
			if p.try(p.parseAlterForeignKeyWithName) {
				return nil
			}
			return p.parseAlterConstraintWithName()
		}
		return p.parseAlterAddPartitions()
	case KW_TOUCH, KW_ARCHIVE, KW_UNARCHIVE:
		p.next()
		for p.peek().Type == KW_PARTITION {
			if err := p.parsePartitionSpec(); err != nil {
				return err
			}
		}
		return nil
	case KW_SKEWED:
		return p.parseTableSkewed()
	case KW_NOT:
		if p.peekAt(1).Type == KW_SKEWED {
			p.pos += 2
			return nil
		}
		if p.peekAt(1).Type == KW_STORED && p.peekAt(2).Type == KW_AS && p.peekAt(3).Type == KW_DIRECTORIES {
			p.pos += 4
			return nil
		}
		return p.parseAlterTblPartitionStatementSuffix()
	case KW_EXCHANGE:
		p.next()
		if err := p.parsePartitionSpec(); err != nil {
			return err
		}
		if _, err := p.expect(KW_WITH); err != nil {
			return err
		}
		if _, err := p.expect(KW_TABLE); err != nil {
			return err
		}
		return p.parseTableName()
	case KW_PARTITION:
		if p.peekAt(1).Type == KW_COLUMN {
			p.pos += 2
			if _, err := p.expect(TokLParen); err != nil {
				return err
			}
			if err := p.parseColumnNameType(); err != nil {
				return err
			}
			_, err := p.expect(TokRParen)
			return err
		}
		if err := p.parsePartitionSpec(); err != nil {
			return err
		}
		return p.parseAlterTblPartitionStatementSuffix()
	case KW_SET:
		switch p.peekAt(1).Type {
		case KW_TBLPROPERTIES:
			p.pos += 2
			return p.parseTableProperties()
		case KW_OWNER:
			p.pos += 2
			return p.parsePrincipalName()
		case KW_PARTITION:
			p.pos += 2
			if _, err := p.expect(KW_SPEC); err != nil {
				return err
			}
			if _, err := p.expect(TokLParen); err != nil {
				return err
			}
			if err := p.parsePartitionTransformSpec(); err != nil {
				return err
			}
			_, err := p.expect(TokRParen)
			return err
		default:
			return p.parseAlterTblPartitionStatementSuffix()
		}
	case KW_UNSET:
		if p.peekAt(1).Type == KW_TBLPROPERTIES {
			p.pos += 2
			p.try(p.parseIfExists)
			return p.parseTableProperties()
		}
		return p.parseAlterTblPartitionStatementSuffix()
	case KW_EXECUTE:
		return p.parseAlterExecute()
	case KW_CREATE:
		if p.peekAt(1).Type == KW_BRANCH {
			return p.parseAlterCreateBranch()
		}
		if p.peekAt(1).Type == KW_TAG {
			return p.parseAlterCreateTag()
		}
		return p.errAt(p.peekAt(1), fmt.Sprintf("expected BRANCH or TAG after ALTER TABLE CREATE, got %q", p.peekAt(1).Str))
	case KW_CONVERT:
		p.next()
		if _, err := p.expect(KW_TO); err != nil {
			return err
		}
		if _, err := p.parseIdentifier(); err != nil {
			return err
		}
		if p.peek().Type == KW_TBLPROPERTIES {
			return p.parseTablePropertiesPrefixed()
		}
		return nil
	default:
		return p.parseAlterTblPartitionStatementSuffix()
	}
}

func (p *parser) parseAlterTblPartitionStatementSuffix() error {
	switch p.peek().Type {
	case KW_SET:
		p.next()
		switch p.peek().Type {
		case KW_FILEFORMAT:
			p.next()
			return p.parseFileFormat()
		case KW_LOCATION:
			p.next()
			_, err := p.expect(TokStringLiteral)
			return err
		case KW_SERDE:
			p.next()
			if _, err := p.expect(TokStringLiteral); err != nil {
				return err
			}
			if p.match(KW_WITH) {
				if _, err := p.expect(KW_SERDEPROPERTIES); err != nil {
					return err
				}
				return p.parseTableProperties()
			}
			return nil
		case KW_SERDEPROPERTIES:
			p.next()
			return p.parseTableProperties()
		case KW_SKEWED:
			p.next()
			if _, err := p.expect(KW_LOCATION); err != nil {
				return err
			}
			return p.parseSkewedLocations()
		default:
			return p.errAt(p.peek(), fmt.Sprintf("unexpected token %q after ALTER TABLE SET", p.peek().Str))
		}
	case KW_UNSET:
		p.next()
		if _, err := p.expect(KW_SERDEPROPERTIES); err != nil {
			return err
		}
		return p.parseTableProperties()
	case KW_CONCATENATE:
		p.next()
		return nil
	case KW_RENAME:
		p.next()
		if _, err := p.expect(KW_TO); err != nil {
			return err
		}
		return p.parsePartitionSpec()
	case KW_INTO:
		p.next()
		if _, err := p.expect(TokNumber); err != nil {
			return err
		}
		_, err := p.expect(KW_BUCKETS)
		return err
	case KW_NOT:
		p.next()
		if p.peek().Type == KW_CLUSTERED || p.peek().Type == KW_SORTED {
			p.next()
			return nil
		}
		return p.errAt(p.peek(), fmt.Sprintf("expected CLUSTERED or SORTED after NOT, got %q", p.peek().Str))
	case KW_CLUSTERED:
		return p.parseTableBuckets()
	case KW_COMPACT:
		p.next()
		if _, err := p.expect(TokStringLiteral); err != nil {
			return err
		}
		if p.peek().Type == KW_CLUSTERED {
			if err := p.parseTableImplBuckets(); err != nil {
				return err
			}
		}
		if p.peek().Type == KW_AND && p.peekAt(1).Type == KW_WAIT {
			p.pos += 2
		}
		if p.match(KW_POOL) {
			if _, err := p.expect(TokStringLiteral); err != nil {
				return err
			}
		}
		if p.match(KW_WITH) {
			if _, err := p.expect(KW_OVERWRITE); err != nil {
				return err
			}
			if _, err := p.expect(KW_TBLPROPERTIES); err != nil {
				return err
			}
			if err := p.parseTableProperties(); err != nil {
				return err
			}
		}
		if p.peek().Type == KW_ORDER {
			return p.parseOrderByClause()
		}
		return nil
	case KW_UPDATE:
		p.next()
		if p.match(KW_COLUMNS) {
			p.parseRestrictOrCascadeOpt()
			return nil
		}
		if _, err := p.expect(KW_STATISTICS); err != nil {
			return err
		}
		if p.match(KW_FOR) {
			p.match(KW_COLUMN)
			if _, err := p.parseIdentifier(); err != nil {
				return err
			}
			if _, err := p.expect(KW_SET); err != nil {
				return err
			}
			if err := p.parseTableProperties(); err != nil {
				return err
			}
			if p.match(KW_COMMENT) {
				_, err := p.expect(TokStringLiteral)
				return err
			}
			return nil
		}
		if _, err := p.expect(KW_SET); err != nil {
			return err
		}
		return p.parseTableProperties()
	case KW_CHANGE:
		p.next()
		p.match(KW_COLUMN)
		if _, err := p.parseIdentifier(); err != nil {
			return err
		}
		if _, err := p.parseIdentifier(); err != nil {
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
		if p.match(KW_FIRST) {
		} else if p.match(KW_AFTER) {
			if _, err := p.parseIdentifier(); err != nil {
				return err
			}
		}
		p.parseRestrictOrCascadeOpt()
		return nil
	case KW_ADD, KW_REPLACE:
		return p.parseAlterAddOrReplaceCols()
	default:
		return p.errAt(p.peek(), fmt.Sprintf("unexpected token %q in ALTER TABLE", p.peek().Str))
	}
}

func (p *parser) parseAlterAddOrReplaceCols() error {
	if p.peek().Type != KW_ADD && p.peek().Type != KW_REPLACE {
		return p.errAt(p.peek(), "expected ADD or REPLACE")
	}
	p.next()
	if _, err := p.expect(KW_COLUMNS); err != nil {
		return err
	}
	if _, err := p.expect(TokLParen); err != nil {
		return err
	}
	if err := p.parseColumnNameTypeList(); err != nil {
		return err
	}
	if _, err := p.expect(TokRParen); err != nil {
		return err
	}
	p.parseRestrictOrCascadeOpt()
	return nil
}

func (p *parser) parseAlterAddPartitions() error {
	if _, err := p.expect(KW_ADD); err != nil {
		return err
	}
	p.try(p.parseIfNotExists)
	if err := p.parsePartitionSpec(); err != nil {
		return err
	}
	if p.match(KW_LOCATION) {
		if _, err := p.expect(TokStringLiteral); err != nil {
			return err
		}
	}
	for p.peek().Type == KW_PARTITION {
		if err := p.parsePartitionSpec(); err != nil {
			return err
		}
		if p.match(KW_LOCATION) {
			if _, err := p.expect(TokStringLiteral); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *parser) parseAlterDropPartitions() error {
	if _, err := p.expect(KW_DROP); err != nil {
		return err
	}
	p.try(p.parseIfExists)
	if _, err := p.expect(KW_PARTITION); err != nil {
		return err
	}
	if err := p.parsePartitionSelectorSpec(); err != nil {
		return err
	}
	for p.match(TokComma) {
		if _, err := p.expect(KW_PARTITION); err != nil {
			return err
		}
		if err := p.parsePartitionSelectorSpec(); err != nil {
			return err
		}
	}
	p.match(KW_PURGE)
	if p.peek().Type == KW_FOR {
		return p.parseReplicationClause()
	}
	return nil
}

func (p *parser) parseFileFormat() error {
	if p.match(KW_INPUTFORMAT) {
		if _, err := p.expect(TokStringLiteral); err != nil {
			return err
		}
		if _, err := p.expect(KW_OUTPUTFORMAT); err != nil {
			return err
		}
		if _, err := p.expect(TokStringLiteral); err != nil {
			return err
		}
		if _, err := p.expect(KW_SERDE); err != nil {
			return err
		}
		if _, err := p.expect(TokStringLiteral); err != nil {
			return err
		}
		if p.match(KW_INPUTDRIVER) {
			if _, err := p.expect(TokStringLiteral); err != nil {
				return err
			}
			if _, err := p.expect(KW_OUTPUTDRIVER); err != nil {
				return err
			}
			if _, err := p.expect(TokStringLiteral); err != nil {
				return err
			}
		}
		return nil
	}
	_, err := p.parseIdentifier()
	return err
}

func (p *parser) parseSkewedLocations() error {
	if _, err := p.expect(TokLParen); err != nil {
		return err
	}
	if err := p.parseSkewedLocationMap(); err != nil {
		return err
	}
	for p.match(TokComma) {
		if err := p.parseSkewedLocationMap(); err != nil {
			return err
		}
	}
	_, err := p.expect(TokRParen)
	return err
}

func (p *parser) parseSkewedLocationMap() error {
	if p.peek().Type == TokLParen {
		if err := p.parseSkewedColumnValuePair(); err != nil {
			return err
		}
	} else {
		if err := p.parseConstant(); err != nil {
			return err
		}
	}
	if _, err := p.expect(TokEqual); err != nil {
		return err
	}
	_, err := p.expect(TokStringLiteral)
	return err
}

func (p *parser) parseAlterExecute() error {
	if _, err := p.expect(KW_EXECUTE); err != nil {
		return err
	}
	switch p.peek().Type {
	case KW_ROLLBACK:
		p.next()
		if _, err := p.expect(TokLParen); err != nil {
			return err
		}
		if p.peek().Type != TokStringLiteral && p.peek().Type != TokNumber {
			return p.errAt(p.peek(), fmt.Sprintf("expected string or number in ROLLBACK, got %q", p.peek().Str))
		}
		p.next()
		_, err := p.expect(TokRParen)
		return err
	case KW_EXPIRE_SNAPSHOTS:
		p.next()
		if p.match(TokLParen) {
			if _, err := p.expect(TokStringLiteral); err != nil {
				return err
			}
			_, err := p.expect(TokRParen)
			return err
		}
		if p.match(KW_BETWEEN) {
			if _, err := p.expect(TokStringLiteral); err != nil {
				return err
			}
			if _, err := p.expect(KW_AND); err != nil {
				return err
			}
			_, err := p.expect(TokStringLiteral)
			return err
		}
		if p.match(KW_RETAIN) {
			if _, err := p.expect(KW_LAST); err != nil {
				return err
			}
			_, err := p.expect(TokNumber)
			return err
		}
		return nil
	case KW_SET_CURRENT_SNAPSHOT:
		p.next()
		if _, err := p.expect(TokLParen); err != nil {
			return err
		}
		if err := p.parseExpression(); err != nil {
			return err
		}
		_, err := p.expect(TokRParen)
		return err
	case KW_FAST_FORWARD:
		p.next()
		if _, err := p.expect(TokStringLiteral); err != nil {
			return err
		}
		p.match(TokStringLiteral)
		return nil
	case KW_CHERRY_PICK:
		p.next()
		_, err := p.expect(TokNumber)
		return err
	case KW_DELETE:
		p.next()
		if _, err := p.expect(KW_ORPHAN_FILES); err != nil {
			return err
		}
		if p.match(KW_OLDER) {
			if _, err := p.expect(KW_THAN); err != nil {
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
		return nil
	default:
		return p.errAt(p.peek(), fmt.Sprintf("unexpected token %q after ALTER TABLE EXECUTE", p.peek().Str))
	}
}

func (p *parser) parseAlterCreateBranch() error {
	if _, err := p.expect(KW_CREATE); err != nil {
		return err
	}
	if _, err := p.expect(KW_BRANCH); err != nil {
		return err
	}
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	if p.peek().Type == KW_FOR {
		if err := p.parseSnapshotIDOfRef(); err != nil {
			return err
		}
	}
	if p.peek().Type == KW_RETAIN {
		if err := p.parseRefRetain(); err != nil {
			return err
		}
	}
	if p.peek().Type == KW_WITH {
		return p.parseRetentionOfSnapshots()
	}
	return nil
}

func (p *parser) parseAlterCreateTag() error {
	if _, err := p.expect(KW_CREATE); err != nil {
		return err
	}
	if _, err := p.expect(KW_TAG); err != nil {
		return err
	}
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	if p.peek().Type == KW_FOR {
		if err := p.parseSnapshotIDOfRef(); err != nil {
			return err
		}
	}
	if p.peek().Type == KW_RETAIN {
		return p.parseRefRetain()
	}
	return nil
}

func (p *parser) parseSnapshotIDOfRef() error {
	if _, err := p.expect(KW_FOR); err != nil {
		return err
	}
	switch p.peek().Type {
	case KW_SYSTEM_VERSION:
		p.next()
		if _, err := p.expect(KW_AS); err != nil {
			return err
		}
		if _, err := p.expect(KW_OF); err != nil {
			return err
		}
		_, err := p.expect(TokNumber)
		return err
	case KW_SYSTEM_TIME:
		p.next()
		if _, err := p.expect(KW_AS); err != nil {
			return err
		}
		if _, err := p.expect(KW_OF); err != nil {
			return err
		}
		_, err := p.expect(TokStringLiteral)
		return err
	case KW_TAG:
		p.next()
		if _, err := p.expect(KW_AS); err != nil {
			return err
		}
		if _, err := p.expect(KW_OF); err != nil {
			return err
		}
		_, err := p.parseIdentifier()
		return err
	default:
		return p.errAt(p.peek(), fmt.Sprintf("expected SYSTEM_VERSION, SYSTEM_TIME, or TAG, got %q", p.peek().Str))
	}
}

func (p *parser) parseRefRetain() error {
	if _, err := p.expect(KW_RETAIN); err != nil {
		return err
	}
	if _, err := p.expect(TokNumber); err != nil {
		return err
	}
	return p.parseTimeUnitQualifiers()
}

func (p *parser) parseRetentionOfSnapshots() error {
	if _, err := p.expect(KW_WITH); err != nil {
		return err
	}
	if _, err := p.expect(KW_SNAPSHOT); err != nil {
		return err
	}
	if _, err := p.expect(KW_RETENTION); err != nil {
		return err
	}
	if _, err := p.expect(TokNumber); err != nil {
		return err
	}
	if _, err := p.expect(KW_SNAPSHOTS); err != nil {
		return err
	}
	if p.match(TokNumber) {
		return p.parseTimeUnitQualifiers()
	}
	return nil
}
