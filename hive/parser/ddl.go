package parser

import "fmt"

// parseDDLStatement dispatches all ddlStatement alternatives in HiveParser.g.
func (p *parser) parseDDLStatement() error {
	switch p.peek().Type {
	case KW_CREATE:
		return p.parseCreateDDL()
	case KW_USE:
		p.next()
		_, err := p.parseIdentifier()
		return err
	case KW_DROP:
		return p.parseDropDDL()
	case KW_TRUNCATE:
		return p.parseTruncateTableStatement()
	case KW_ALTER, KW_OPTIMIZE:
		return p.parseAlterDDL()
	case KW_DESCRIBE, KW_DESC:
		return p.parseDescStatement()
	case KW_SHOW:
		return p.parseShowStatement()
	case KW_MSCK:
		return p.parseMetastoreCheck()
	case KW_RELOAD:
		p.next()
		if p.peek().Type == KW_FUNCTIONS || p.peek().Type == KW_FUNCTION {
			p.next()
			return nil
		}
		return p.errAt(p.peek(), fmt.Sprintf("expected FUNCTION or FUNCTIONS after RELOAD, got %q", p.peek().Str))
	case KW_ANALYZE:
		return p.parseAnalyzeStatement()
	case KW_LOCK:
		return p.parseLockStatement()
	case KW_UNLOCK:
		return p.parseUnlockStatement()
	case KW_GRANT:
		if p.try(p.parseGrantPrivileges) {
			return nil
		}
		return p.parseGrantRole()
	case KW_REVOKE:
		if p.try(p.parseRevokePrivileges) {
			return nil
		}
		return p.parseRevokeRole()
	case KW_SET:
		return p.parseSetRole()
	case KW_ABORT:
		return p.parseAbortStatement()
	case KW_KILL:
		return p.parseKillQueryStatement()
	case KW_ENABLE, KW_DISABLE:
		return p.parseGlobalWmStatement()
	case KW_REPLACE:
		return p.parseReplaceResourcePlanStatement()
	default:
		return p.errAt(p.peek(), fmt.Sprintf("unexpected DDL keyword %q", p.peek().Str))
	}
}

func (p *parser) parseIfExists() error {
	if _, err := p.expect(KW_IF); err != nil {
		return err
	}
	_, err := p.expect(KW_EXISTS)
	return err
}

func (p *parser) parseIfNotExists() error {
	if _, err := p.expect(KW_IF); err != nil {
		return err
	}
	if _, err := p.expect(KW_NOT); err != nil {
		return err
	}
	_, err := p.expect(KW_EXISTS)
	return err
}

func (p *parser) parseRestrictOrCascadeOpt() {
	if p.peek().Type == KW_RESTRICT || p.peek().Type == KW_CASCADE {
		p.next()
	}
}

// --- CREATE ---

func (p *parser) parseCreateDDL() error {
	switch p.peekAt(1).Type {
	case KW_DATABASE, KW_SCHEMA, KW_REMOTE:
		return p.parseCreateDatabaseStatement()
	case KW_DATACONNECTOR:
		return p.parseCreateDataConnectorStatement()
	case KW_ROLE:
		p.pos += 2
		_, err := p.parseIdentifier()
		return err
	case KW_VIEW, KW_OR:
		return p.parseCreateViewStatement()
	case KW_MATERIALIZED:
		return p.parseCreateMaterializedViewStatement()
	case KW_SCHEDULED:
		return p.parseCreateScheduledQueryStatement()
	case KW_FUNCTION:
		return p.parseCreateFunctionStatement()
	case KW_TEMPORARY:
		switch p.peekAt(2).Type {
		case KW_FUNCTION:
			return p.parseCreateFunctionStatement()
		case KW_MACRO:
			return p.parseCreateMacroStatement()
		default:
			return p.parseCreateTableStatement()
		}
	case KW_TABLE, KW_TRANSACTIONAL, KW_EXTERNAL, KW_MANAGED:
		return p.parseCreateTableStatement()
	case KW_RESOURCE:
		return p.parseCreateResourcePlanStatement()
	case KW_TRIGGER:
		return p.parseCreateTriggerStatement()
	case KW_POOL:
		return p.parseCreatePoolStatement()
	case KW_USER, KW_GROUP, KW_APPLICATION:
		return p.parseCreateMappingStatement()
	default:
		return p.errAt(p.peekAt(1), fmt.Sprintf("unexpected token %q after CREATE", p.peekAt(1).Str))
	}
}

func (p *parser) parseCreateDatabaseStatement() error {
	if _, err := p.expect(KW_CREATE); err != nil {
		return err
	}
	remote := p.match(KW_REMOTE)
	if p.peek().Type != KW_DATABASE && p.peek().Type != KW_SCHEMA {
		return p.errAt(p.peek(), fmt.Sprintf("expected DATABASE or SCHEMA, got %q", p.peek().Str))
	}
	p.next()
	p.try(p.parseIfNotExists)
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	if p.match(KW_COMMENT) {
		if _, err := p.expect(TokStringLiteral); err != nil {
			return err
		}
	}
	if remote {
		if _, err := p.expect(KW_USING); err != nil {
			return err
		}
		if _, err := p.parseIdentifier(); err != nil {
			return err
		}
	} else {
		if p.match(KW_LOCATION) {
			if _, err := p.expect(TokStringLiteral); err != nil {
				return err
			}
		}
		if p.match(KW_MANAGEDLOCATION) {
			if _, err := p.expect(TokStringLiteral); err != nil {
				return err
			}
		}
	}
	if p.match(KW_WITH) {
		if _, err := p.expect(KW_DBPROPERTIES); err != nil {
			return err
		}
		return p.parseDBProperties()
	}
	return nil
}

func (p *parser) parseDBProperties() error {
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

func (p *parser) parseCreateDataConnectorStatement() error {
	if _, err := p.expect(KW_CREATE); err != nil {
		return err
	}
	if _, err := p.expect(KW_DATACONNECTOR); err != nil {
		return err
	}
	p.try(p.parseIfNotExists)
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	if _, err := p.expect(KW_TYPE); err != nil {
		return err
	}
	if _, err := p.expect(TokStringLiteral); err != nil {
		return err
	}
	if _, err := p.expect(KW_URL); err != nil {
		return err
	}
	if _, err := p.expect(TokStringLiteral); err != nil {
		return err
	}
	if p.match(KW_COMMENT) {
		if _, err := p.expect(TokStringLiteral); err != nil {
			return err
		}
	}
	if p.match(KW_WITH) {
		if _, err := p.expect(KW_DCPROPERTIES); err != nil {
			return err
		}
		return p.parseDBProperties()
	}
	return nil
}

func (p *parser) parseCreateTableStatement() error {
	if _, err := p.expect(KW_CREATE); err != nil {
		return err
	}
	managed := false
	if p.match(KW_MANAGED) {
		managed = true
	} else {
		p.match(KW_TEMPORARY)
		p.match(KW_TRANSACTIONAL)
		p.match(KW_EXTERNAL)
	}
	if _, err := p.expect(KW_TABLE); err != nil {
		return err
	}
	p.try(p.parseIfNotExists)
	if err := p.parseTableName(); err != nil {
		return err
	}
	if p.peek().Type == KW_LIKE {
		if err := p.parseLikeTableOrFile(); err != nil {
			return err
		}
		if !managed && p.peek().Type == KW_PARTITIONED {
			if err := p.parseCreateTablePartitionSpec(); err != nil {
				return err
			}
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
		if p.peek().Type == KW_LOCATION {
			if err := p.parseTableLocation(); err != nil {
				return err
			}
		}
		if p.peek().Type == KW_TBLPROPERTIES {
			if err := p.parseTablePropertiesPrefixed(); err != nil {
				return err
			}
		}
		return nil
	}
	if p.match(TokLParen) {
		if err := p.parseColumnNameTypeOrConstraintList(); err != nil {
			return err
		}
		if _, err := p.expect(TokRParen); err != nil {
			return err
		}
	}
	if p.match(KW_COMMENT) {
		if _, err := p.expect(TokStringLiteral); err != nil {
			return err
		}
	}
	if p.peek().Type == KW_PARTITIONED {
		if err := p.parseCreateTablePartitionSpec(); err != nil {
			return err
		}
	}
	if p.peek().Type == KW_CLUSTERED {
		if err := p.parseTableBuckets(); err != nil {
			return err
		}
	}
	if p.peek().Type == KW_SKEWED {
		if err := p.parseTableSkewed(); err != nil {
			return err
		}
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
	if p.peek().Type == KW_LOCATION {
		if err := p.parseTableLocation(); err != nil {
			return err
		}
	}
	if p.peek().Type == KW_TBLPROPERTIES {
		if err := p.parseTablePropertiesPrefixed(); err != nil {
			return err
		}
	}
	if p.match(KW_AS) {
		return p.parseSelectStatementWithCTE()
	}
	return nil
}

// likeTableOrFile in CreateDDLParser.g:46:
//
//	(KW_LIKE KW_FILE) => KW_LIKE KW_FILE
//	| KW_LIKE KW_FILE format=identifier uri=StringLiteral
//	| KW_LIKE likeName=tableName
//
// HiveParser.g decides with k=3 and no backtracking (HiveParser.g:24-25), so
// the token after LIKE FILE picks the form: one that can be an identifier
// commits to `format uri`, even where it would begin the next clause
// (Hive 4.0.1 refuses LIKE FILE STORED BY 'x' and LIKE FILE TBLPROPERTIES
// (...)); anything else, or the end, is bare LIKE FILE.
func (p *parser) parseLikeTableOrFile() error {
	if _, err := p.expect(KW_LIKE); err != nil {
		return err
	}
	if p.peek().Type == KW_FILE && p.peekAt(1).Type != TokDot {
		p.next()
		if isIdentifier(p.peek().Type) {
			p.next()
			_, err := p.expect(TokStringLiteral)
			return err
		}
		return nil
	}
	return p.parseTableName()
}

func (p *parser) parseCreateTablePartitionSpec() error {
	if _, err := p.expect(KW_PARTITIONED); err != nil {
		return err
	}
	if _, err := p.expect(KW_BY); err != nil {
		return err
	}
	if p.match(KW_SPEC) {
		if _, err := p.expect(TokLParen); err != nil {
			return err
		}
		if err := p.parsePartitionTransformSpec(); err != nil {
			return err
		}
		_, err := p.expect(TokRParen)
		return err
	}
	if _, err := p.expect(TokLParen); err != nil {
		return err
	}
	if !p.try(func() error {
		if err := p.parseColumnNameTypeConstraint(); err != nil {
			return err
		}
		for p.match(TokComma) {
			if err := p.parseColumnNameTypeConstraint(); err != nil {
				return err
			}
		}
		return nil
	}) {
		if err := p.parseColumnNameList(); err != nil {
			return err
		}
	}
	_, err := p.expect(TokRParen)
	return err
}

func (p *parser) parsePartitionTransformSpec() error {
	if err := p.parsePartitionTransformType(); err != nil {
		return err
	}
	for p.match(TokComma) {
		if err := p.parsePartitionTransformType(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parsePartitionTransformType() error {
	switch p.peek().Type {
	case KW_YEAR, KW_MONTH, KW_DAY, KW_HOUR:
		if p.peekAt(1).Type == TokLParen {
			p.pos += 2
			col, err := p.parseIdentifier()
			if err != nil {
				return err
			}
			if err := p.checkCreateTableColName(col); err != nil {
				return err
			}
			_, err = p.expect(TokRParen)
			return err
		}
	case KW_TRUNCATE, KW_BUCKET:
		if p.peekAt(1).Type == TokLParen {
			p.pos += 2
			if _, err := p.expect(TokNumber); err != nil {
				return err
			}
			if _, err := p.expect(TokComma); err != nil {
				return err
			}
			col, err := p.parseIdentifier()
			if err != nil {
				return err
			}
			if err := p.checkCreateTableColName(col); err != nil {
				return err
			}
			_, err = p.expect(TokRParen)
			return err
		}
	}
	col, err := p.parseIdentifier()
	if err != nil {
		return err
	}
	return p.checkCreateTableColName(col)
}

func (p *parser) parseTableBuckets() error {
	if _, err := p.expect(KW_CLUSTERED); err != nil {
		return err
	}
	if _, err := p.expect(KW_BY); err != nil {
		return err
	}
	if err := p.parseColumnParenthesesList(); err != nil {
		return err
	}
	if p.match(KW_SORTED) {
		if _, err := p.expect(KW_BY); err != nil {
			return err
		}
		if _, err := p.expect(TokLParen); err != nil {
			return err
		}
		if err := p.parseColumnNameOrderList(); err != nil {
			return err
		}
		if _, err := p.expect(TokRParen); err != nil {
			return err
		}
	}
	if _, err := p.expect(KW_INTO); err != nil {
		return err
	}
	if _, err := p.expect(TokNumber); err != nil {
		return err
	}
	_, err := p.expect(KW_BUCKETS)
	return err
}

func (p *parser) parseTableImplBuckets() error {
	if _, err := p.expect(KW_CLUSTERED); err != nil {
		return err
	}
	if _, err := p.expect(KW_INTO); err != nil {
		return err
	}
	if _, err := p.expect(TokNumber); err != nil {
		return err
	}
	_, err := p.expect(KW_BUCKETS)
	return err
}

func (p *parser) parseTableSkewed() error {
	if _, err := p.expect(KW_SKEWED); err != nil {
		return err
	}
	if _, err := p.expect(KW_BY); err != nil {
		return err
	}
	if err := p.parseColumnParenthesesList(); err != nil {
		return err
	}
	if _, err := p.expect(KW_ON); err != nil {
		return err
	}
	if _, err := p.expect(TokLParen); err != nil {
		return err
	}
	if p.peek().Type == TokLParen {
		if err := p.parseSkewedColumnValuePair(); err != nil {
			return err
		}
		for p.match(TokComma) {
			if err := p.parseSkewedColumnValuePair(); err != nil {
				return err
			}
		}
	} else {
		if err := p.parseSkewedColumnValues(); err != nil {
			return err
		}
	}
	if _, err := p.expect(TokRParen); err != nil {
		return err
	}
	if p.peek().Type == KW_STORED && p.peekAt(1).Type == KW_AS && p.peekAt(2).Type == KW_DIRECTORIES {
		p.pos += 3
	}
	return nil
}

func (p *parser) parseSkewedColumnValuePair() error {
	if _, err := p.expect(TokLParen); err != nil {
		return err
	}
	if err := p.parseSkewedColumnValues(); err != nil {
		return err
	}
	_, err := p.expect(TokRParen)
	return err
}

func (p *parser) parseSkewedColumnValues() error {
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

func (p *parser) parseOptionalRowFormat() error {
	if p.peek().Type == KW_ROW {
		return p.parseTableRowFormat()
	}
	return nil
}

func (p *parser) parseTableRowFormat() error {
	if _, err := p.expect(KW_ROW); err != nil {
		return err
	}
	if _, err := p.expect(KW_FORMAT); err != nil {
		return err
	}
	if p.match(KW_SERDE) {
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
	}
	if _, err := p.expect(KW_DELIMITED); err != nil {
		return err
	}
	if p.match(KW_FIELDS) {
		if _, err := p.expect(KW_TERMINATED); err != nil {
			return err
		}
		if _, err := p.expect(KW_BY); err != nil {
			return err
		}
		if _, err := p.expect(TokStringLiteral); err != nil {
			return err
		}
		if p.match(KW_ESCAPED) {
			if _, err := p.expect(KW_BY); err != nil {
				return err
			}
			if _, err := p.expect(TokStringLiteral); err != nil {
				return err
			}
		}
	}
	if p.match(KW_COLLECTION) {
		if _, err := p.expect(KW_ITEMS); err != nil {
			return err
		}
		if _, err := p.expect(KW_TERMINATED); err != nil {
			return err
		}
		if _, err := p.expect(KW_BY); err != nil {
			return err
		}
		if _, err := p.expect(TokStringLiteral); err != nil {
			return err
		}
	}
	if p.match(KW_MAP) {
		if _, err := p.expect(KW_KEYS); err != nil {
			return err
		}
		if _, err := p.expect(KW_TERMINATED); err != nil {
			return err
		}
		if _, err := p.expect(KW_BY); err != nil {
			return err
		}
		if _, err := p.expect(TokStringLiteral); err != nil {
			return err
		}
	}
	if p.match(KW_LINES) {
		if _, err := p.expect(KW_TERMINATED); err != nil {
			return err
		}
		if _, err := p.expect(KW_BY); err != nil {
			return err
		}
		if _, err := p.expect(TokStringLiteral); err != nil {
			return err
		}
	}
	if p.match(KW_NULL) {
		if _, err := p.expect(KW_DEFINED); err != nil {
			return err
		}
		if _, err := p.expect(KW_AS); err != nil {
			return err
		}
		if _, err := p.expect(TokStringLiteral); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseTableFileFormat() error {
	if _, err := p.expect(KW_STORED); err != nil {
		return err
	}
	if p.match(KW_BY) {
		if p.peek().Type == TokStringLiteral {
			p.next()
		} else if _, err := p.parseIdentifier(); err != nil {
			return err
		}
		if p.match(KW_WITH) {
			if _, err := p.expect(KW_SERDEPROPERTIES); err != nil {
				return err
			}
			if err := p.parseTableProperties(); err != nil {
				return err
			}
		}
		if p.match(KW_STORED) {
			if _, err := p.expect(KW_AS); err != nil {
				return err
			}
			_, err := p.parseIdentifier()
			return err
		}
		return nil
	}
	if _, err := p.expect(KW_AS); err != nil {
		return err
	}
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

func (p *parser) parseTableLocation() error {
	if _, err := p.expect(KW_LOCATION); err != nil {
		return err
	}
	_, err := p.expect(TokStringLiteral)
	return err
}

func (p *parser) parseTablePropertiesPrefixed() error {
	if _, err := p.expect(KW_TBLPROPERTIES); err != nil {
		return err
	}
	return p.parseTableProperties()
}

func (p *parser) parseTableProperties() error {
	if _, err := p.expect(TokLParen); err != nil {
		return err
	}
	if p.peek().Type == TokStringLiteral && p.peekAt(1).Type == TokEqual {
		if err := p.parseKeyValueProperty(); err != nil {
			return err
		}
		for p.match(TokComma) {
			if err := p.parseKeyValueProperty(); err != nil {
				return err
			}
		}
	} else {
		if _, err := p.expect(TokStringLiteral); err != nil {
			return err
		}
		for p.match(TokComma) {
			if _, err := p.expect(TokStringLiteral); err != nil {
				return err
			}
		}
	}
	_, err := p.expect(TokRParen)
	return err
}

func (p *parser) parseKeyValueProperty() error {
	if _, err := p.expect(TokStringLiteral); err != nil {
		return err
	}
	if _, err := p.expect(TokEqual); err != nil {
		return err
	}
	_, err := p.expect(TokStringLiteral)
	return err
}

// --- CREATE VIEW / MV / FUNCTION / MACRO / SCHEDULED QUERY ---

func (p *parser) parseCreateViewStatement() error {
	if _, err := p.expect(KW_CREATE); err != nil {
		return err
	}
	if p.match(KW_OR) {
		if _, err := p.expect(KW_REPLACE); err != nil {
			return err
		}
	}
	if _, err := p.expect(KW_VIEW); err != nil {
		return err
	}
	p.try(p.parseIfNotExists)
	if err := p.parseTableName(); err != nil {
		return err
	}
	if p.match(TokLParen) {
		if err := p.parseColumnNameCommentList(); err != nil {
			return err
		}
		if _, err := p.expect(TokRParen); err != nil {
			return err
		}
	}
	if p.match(KW_COMMENT) {
		if _, err := p.expect(TokStringLiteral); err != nil {
			return err
		}
	}
	if p.peek().Type == KW_PARTITIONED {
		if err := p.parseViewPartition(); err != nil {
			return err
		}
	}
	if p.peek().Type == KW_TBLPROPERTIES {
		if err := p.parseTablePropertiesPrefixed(); err != nil {
			return err
		}
	}
	if _, err := p.expect(KW_AS); err != nil {
		return err
	}
	return p.parseSelectStatementWithCTE()
}

func (p *parser) parseViewPartition() error {
	if _, err := p.expect(KW_PARTITIONED); err != nil {
		return err
	}
	if _, err := p.expect(KW_ON); err != nil {
		return err
	}
	if p.match(KW_SPEC) {
		if _, err := p.expect(TokLParen); err != nil {
			return err
		}
		if err := p.parsePartitionTransformSpec(); err != nil {
			return err
		}
		_, err := p.expect(TokRParen)
		return err
	}
	return p.parseColumnParenthesesList()
}

func (p *parser) parseCreateMaterializedViewStatement() error {
	if _, err := p.expect(KW_CREATE); err != nil {
		return err
	}
	if _, err := p.expect(KW_MATERIALIZED); err != nil {
		return err
	}
	if _, err := p.expect(KW_VIEW); err != nil {
		return err
	}
	p.try(p.parseIfNotExists)
	if err := p.parseTableName(); err != nil {
		return err
	}
	if p.peek().Type == KW_DISABLE && p.peekAt(1).Type == KW_REWRITE {
		p.pos += 2
	}
	if p.match(KW_COMMENT) {
		if _, err := p.expect(TokStringLiteral); err != nil {
			return err
		}
	}
	if p.peek().Type == KW_PARTITIONED {
		if err := p.parseViewPartition(); err != nil {
			return err
		}
	}
	if p.peek().Type == KW_CLUSTERED {
		p.next()
		if _, err := p.expect(KW_ON); err != nil {
			return err
		}
		if err := p.parseColumnParenthesesList(); err != nil {
			return err
		}
	} else if p.peek().Type == KW_DISTRIBUTED {
		p.next()
		if _, err := p.expect(KW_ON); err != nil {
			return err
		}
		if err := p.parseColumnParenthesesList(); err != nil {
			return err
		}
		if _, err := p.expect(KW_SORTED); err != nil {
			return err
		}
		if _, err := p.expect(KW_ON); err != nil {
			return err
		}
		if err := p.parseColumnParenthesesList(); err != nil {
			return err
		}
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
	if p.peek().Type == KW_LOCATION {
		if err := p.parseTableLocation(); err != nil {
			return err
		}
	}
	if p.peek().Type == KW_TBLPROPERTIES {
		if err := p.parseTablePropertiesPrefixed(); err != nil {
			return err
		}
	}
	if _, err := p.expect(KW_AS); err != nil {
		return err
	}
	return p.parseSelectStatementWithCTE()
}

func (p *parser) parseCreateFunctionStatement() error {
	if _, err := p.expect(KW_CREATE); err != nil {
		return err
	}
	p.match(KW_TEMPORARY)
	if _, err := p.expect(KW_FUNCTION); err != nil {
		return err
	}
	if err := p.parseFunctionIdentifier(); err != nil {
		return err
	}
	if _, err := p.expect(KW_AS); err != nil {
		return err
	}
	if _, err := p.expect(TokStringLiteral); err != nil {
		return err
	}
	if p.match(KW_USING) {
		if err := p.parseResource(); err != nil {
			return err
		}
		for p.match(TokComma) {
			if err := p.parseResource(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *parser) parseResource() error {
	switch p.peek().Type {
	case KW_JAR, KW_FILE, KW_ARCHIVE:
		p.next()
	default:
		return p.errAt(p.peek(), fmt.Sprintf("expected JAR, FILE, or ARCHIVE, got %q", p.peek().Str))
	}
	_, err := p.expect(TokStringLiteral)
	return err
}

func (p *parser) parseCreateMacroStatement() error {
	if _, err := p.expect(KW_CREATE); err != nil {
		return err
	}
	if _, err := p.expect(KW_TEMPORARY); err != nil {
		return err
	}
	if _, err := p.expect(KW_MACRO); err != nil {
		return err
	}
	if _, err := p.expect(TokIdentifier); err != nil {
		return err
	}
	if _, err := p.expect(TokLParen); err != nil {
		return err
	}
	if p.peek().Type != TokRParen {
		if err := p.parseColumnNameTypeList(); err != nil {
			return err
		}
	}
	if _, err := p.expect(TokRParen); err != nil {
		return err
	}
	return p.parseExpression()
}

func (p *parser) parseCreateScheduledQueryStatement() error {
	if _, err := p.expect(KW_CREATE); err != nil {
		return err
	}
	if _, err := p.expect(KW_SCHEDULED); err != nil {
		return err
	}
	if _, err := p.expect(KW_QUERY); err != nil {
		return err
	}
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	if err := p.parseScheduleSpec(); err != nil {
		return err
	}
	if p.peek().Type == KW_EXECUTED {
		if err := p.parseExecutedAsSpec(); err != nil {
			return err
		}
	}
	if p.peek().Type == KW_ENABLE || p.peek().Type == KW_DISABLE {
		p.next()
	}
	return p.parseDefinedAsSpec()
}

func (p *parser) parseScheduleSpec() error {
	if p.match(KW_CRON) {
		_, err := p.expect(TokStringLiteral)
		return err
	}
	if _, err := p.expect(KW_EVERY); err != nil {
		return err
	}
	p.match(TokNumber)
	if err := p.parseIntervalQualifiers(); err != nil {
		return err
	}
	if p.match(KW_AT) {
		_, err := p.expect(TokStringLiteral)
		return err
	}
	if p.match(KW_OFFSET) {
		if _, err := p.expect(KW_BY); err != nil {
			return err
		}
		_, err := p.expect(TokStringLiteral)
		return err
	}
	return nil
}

func (p *parser) parseExecutedAsSpec() error {
	if _, err := p.expect(KW_EXECUTED); err != nil {
		return err
	}
	if _, err := p.expect(KW_AS); err != nil {
		return err
	}
	_, err := p.expect(TokStringLiteral)
	return err
}

func (p *parser) parseDefinedAsSpec() error {
	p.match(KW_DEFINED)
	if _, err := p.expect(KW_AS); err != nil {
		return err
	}
	_, _, err := p.parseStatement()
	return err
}

// --- DROP / TRUNCATE ---

func (p *parser) parseDropDDL() error {
	switch p.peekAt(1).Type {
	case KW_DATABASE, KW_SCHEMA:
		p.pos += 2
		p.try(p.parseIfExists)
		if _, err := p.parseIdentifier(); err != nil {
			return err
		}
		p.parseRestrictOrCascadeOpt()
		return nil
	case KW_TABLE:
		p.pos += 2
		p.try(p.parseIfExists)
		if err := p.parseTableName(); err != nil {
			return err
		}
		p.match(KW_PURGE)
		if p.peek().Type == KW_FOR {
			return p.parseReplicationClause()
		}
		return nil
	case KW_VIEW:
		p.pos += 2
		p.try(p.parseIfExists)
		return p.parseViewName()
	case KW_MATERIALIZED:
		p.pos += 2
		if _, err := p.expect(KW_VIEW); err != nil {
			return err
		}
		p.try(p.parseIfExists)
		return p.parseViewName()
	case KW_SCHEDULED:
		p.pos += 2
		if _, err := p.expect(KW_QUERY); err != nil {
			return err
		}
		_, err := p.parseIdentifier()
		return err
	case KW_FUNCTION:
		p.pos += 2
		p.try(p.parseIfExists)
		return p.parseFunctionIdentifier()
	case KW_TEMPORARY:
		p.pos += 2
		if p.match(KW_FUNCTION) {
			p.try(p.parseIfExists)
			return p.parseFunctionIdentifier()
		}
		if _, err := p.expect(KW_MACRO); err != nil {
			return err
		}
		p.try(p.parseIfExists)
		_, err := p.expect(TokIdentifier)
		return err
	case KW_ROLE:
		p.pos += 2
		_, err := p.parseIdentifier()
		return err
	case KW_DATACONNECTOR:
		p.pos += 2
		p.try(p.parseIfExists)
		_, err := p.parseIdentifier()
		return err
	case KW_RESOURCE:
		return p.parseDropResourcePlanStatement()
	case KW_TRIGGER:
		return p.parseDropTriggerStatement()
	case KW_POOL:
		return p.parseDropPoolStatement()
	case KW_USER, KW_GROUP, KW_APPLICATION:
		return p.parseDropMappingStatement()
	default:
		return p.errAt(p.peekAt(1), fmt.Sprintf("unexpected token %q after DROP", p.peekAt(1).Str))
	}
}

func (p *parser) parseTruncateTableStatement() error {
	if _, err := p.expect(KW_TRUNCATE); err != nil {
		return err
	}
	p.match(KW_TABLE)
	if err := p.parseTableOrPartition(); err != nil {
		return err
	}
	if p.match(KW_COLUMNS) {
		if err := p.parseColumnParenthesesList(); err != nil {
			return err
		}
	}
	p.match(KW_FORCE)
	return nil
}

// --- PARTITIONS ---

func (p *parser) parseTableOrPartition() error {
	if err := p.parseTableName(); err != nil {
		return err
	}
	if p.peek().Type == KW_PARTITION {
		return p.parsePartitionSpec()
	}
	return nil
}

func (p *parser) parsePartitionSpec() error {
	if _, err := p.expect(KW_PARTITION); err != nil {
		return err
	}
	if _, err := p.expect(TokLParen); err != nil {
		return err
	}
	if err := p.parsePartitionVal(); err != nil {
		return err
	}
	for p.match(TokComma) {
		if err := p.parsePartitionVal(); err != nil {
			return err
		}
	}
	_, err := p.expect(TokRParen)
	return err
}

func (p *parser) parsePartitionVal() error {
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	if p.match(TokEqual) {
		return p.parseConstant()
	}
	return nil
}

func (p *parser) parsePartitionSelectorSpec() error {
	if _, err := p.expect(TokLParen); err != nil {
		return err
	}
	if err := p.parsePartitionSelectorVal(); err != nil {
		return err
	}
	for p.match(TokComma) {
		if err := p.parsePartitionSelectorVal(); err != nil {
			return err
		}
	}
	_, err := p.expect(TokRParen)
	return err
}

func (p *parser) parsePartitionSelectorVal() error {
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	if p.peek().Type == KW_LIKE || isSubQuerySelectorOperator(p.peek().Type) {
		p.next()
	} else {
		return p.errAt(p.peek(), fmt.Sprintf("expected partition selector operator, got %q", p.peek().Str))
	}
	return p.parseConstant()
}

// --- DESCRIBE / SHOW / MSCK / ANALYZE ---

func (p *parser) parseDescStatement() error {
	if p.peek().Type != KW_DESCRIBE && p.peek().Type != KW_DESC {
		return p.errAt(p.peek(), "expected DESCRIBE or DESC")
	}
	p.next()
	switch p.peek().Type {
	case KW_DATABASE, KW_SCHEMA:
		p.next()
		p.match(KW_EXTENDED)
		_, err := p.parseIdentifier()
		return err
	case KW_DATACONNECTOR:
		p.next()
		p.match(KW_EXTENDED)
		_, err := p.parseIdentifier()
		return err
	case KW_FUNCTION:
		p.next()
		p.match(KW_EXTENDED)
		return p.parseDescFuncNames()
	case KW_FORMATTED, KW_EXTENDED:
		if p.try(func() error {
			p.next()
			return p.parseTabPartColTypeExpr()
		}) {
			return nil
		}
		return p.parseTabPartColTypeExpr()
	default:
		return p.parseTabPartColTypeExpr()
	}
}

func (p *parser) parseDescFuncNames() error {
	if sysFuncNames[p.peek().Type] {
		p.next()
		return nil
	}
	if p.match(TokStringLiteral) {
		return nil
	}
	return p.parseFunctionIdentifier()
}

func (p *parser) parseTabPartColTypeExpr() error {
	if err := p.parseTableName(); err != nil {
		return err
	}
	if p.peek().Type == KW_PARTITION {
		if err := p.parsePartitionSpec(); err != nil {
			return err
		}
	}
	if isIdentifier(p.peek().Type) {
		return p.parseExtColumnName()
	}
	return nil
}

func (p *parser) parseExtColumnName() error {
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	for p.match(TokDot) {
		switch p.peek().Type {
		case KW_ELEM_TYPE, KW_KEY_TYPE, KW_VALUE_TYPE:
			p.next()
		default:
			if _, err := p.parseIdentifier(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *parser) parseTabTypeExpr() error {
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	if p.match(TokDot) {
		if _, err := p.parseIdentifier(); err != nil {
			return err
		}
	}
	if isIdentifier(p.peek().Type) && !p.isAfterTabTypeBoundary() {
		return p.parseExtColumnName()
	}
	return nil
}

func (p *parser) isAfterTabTypeBoundary() bool {
	switch p.peek().Type {
	case KW_PARTITION, KW_EXTENDED, KW_ORDER:
		return true
	case KW_LIMIT:
		// HiveParser.g:2805 limitClause: KW_LIMIT Number. Hive decides
		// tabTypeExpr's (identifier ...)? (HiveParser.g:1272) on more than one
		// token, so LIMIT not followed by a number is a column name, as in
		// SHOW LOCKS a LIMIT.
		return p.peekAt(1).Type == TokNumber
	case KW_POOL, KW_TYPE, KW_STATUS:
		// HiveParser.g:2984-2991: POOL, TYPE and STATUS take a string; without
		// one they are column names.
		return p.peekAt(1).Type == TokStringLiteral
	}
	return false
}

func (p *parser) parsePartTypeExpr() error {
	if err := p.parseTabTypeExpr(); err != nil {
		return err
	}
	if p.peek().Type == KW_PARTITION {
		return p.parsePartitionSpec()
	}
	return nil
}

func (p *parser) parseAnalyzeStatement() error {
	if _, err := p.expect(KW_ANALYZE); err != nil {
		return err
	}
	if _, err := p.expect(KW_TABLE); err != nil {
		return err
	}
	if err := p.parseTableOrPartition(); err != nil {
		return err
	}
	if p.match(KW_CACHE) {
		_, err := p.expect(KW_METADATA)
		return err
	}
	if _, err := p.expect(KW_COMPUTE); err != nil {
		return err
	}
	if _, err := p.expect(KW_STATISTICS); err != nil {
		return err
	}
	if p.match(KW_NOSCAN) {
		return nil
	}
	if p.match(KW_FOR) {
		if _, err := p.expect(KW_COLUMNS); err != nil {
			return err
		}
		if isIdentifier(p.peek().Type) {
			return p.parseColumnNameList()
		}
	}
	return nil
}

func (p *parser) parseShowStatement() error {
	if _, err := p.expect(KW_SHOW); err != nil {
		return err
	}
	switch p.peek().Type {
	case KW_DATABASES, KW_SCHEMAS:
		p.next()
		if p.match(KW_LIKE) {
			return p.parseShowStmtIdentifier()
		}
		return nil
	case KW_EXTENDED:
		p.next()
		return p.parseShowTablesRest()
	case KW_TABLES:
		return p.parseShowTablesRest()
	case KW_VIEWS:
		p.next()
		return p.parseShowFromLikeRest()
	case KW_MATERIALIZED:
		p.next()
		if _, err := p.expect(KW_VIEWS); err != nil {
			return err
		}
		return p.parseShowFromLikeRest()
	case KW_SORTED, KW_COLUMNS:
		p.match(KW_SORTED)
		if _, err := p.expect(KW_COLUMNS); err != nil {
			return err
		}
		if p.peek().Type != KW_FROM && p.peek().Type != KW_IN {
			return p.errAt(p.peek(), fmt.Sprintf("expected FROM or IN, got %q", p.peek().Str))
		}
		p.next()
		if err := p.parseTableName(); err != nil {
			return err
		}
		return p.parseShowFromLikeRest()
	case KW_FUNCTIONS:
		p.next()
		if p.match(KW_LIKE) {
			if p.match(TokStringLiteral) {
				return nil
			}
			return p.parseFunctionIdentifier()
		}
		return nil
	case KW_PARTITIONS:
		p.next()
		if err := p.parseTableName(); err != nil {
			return err
		}
		if p.peek().Type == KW_PARTITION {
			if err := p.parsePartitionSpec(); err != nil {
				return err
			}
		}
		if p.peek().Type == KW_WHERE {
			if err := p.parseWhereClause(); err != nil {
				return err
			}
		}
		if p.peek().Type == KW_ORDER {
			if err := p.parseOrderByClause(); err != nil {
				return err
			}
		}
		if p.peek().Type == KW_LIMIT {
			return p.parseLimitClause()
		}
		return nil
	case KW_CREATE:
		p.next()
		if p.peek().Type == KW_DATABASE || p.peek().Type == KW_SCHEMA {
			p.next()
			_, err := p.parseIdentifier()
			return err
		}
		if _, err := p.expect(KW_TABLE); err != nil {
			return err
		}
		return p.parseTableName()
	case KW_TABLE:
		p.next()
		if _, err := p.expect(KW_EXTENDED); err != nil {
			return err
		}
		if p.peek().Type == KW_FROM || p.peek().Type == KW_IN {
			p.next()
			if _, err := p.parseIdentifier(); err != nil {
				return err
			}
		}
		if _, err := p.expect(KW_LIKE); err != nil {
			return err
		}
		if err := p.parseShowStmtIdentifier(); err != nil {
			return err
		}
		if p.peek().Type == KW_PARTITION {
			return p.parsePartitionSpec()
		}
		return nil
	case KW_TBLPROPERTIES:
		p.next()
		if err := p.parseTableName(); err != nil {
			return err
		}
		if p.match(TokLParen) {
			if _, err := p.expect(TokStringLiteral); err != nil {
				return err
			}
			_, err := p.expect(TokRParen)
			return err
		}
		return nil
	case KW_LOCKS:
		p.next()
		if p.peek().Type == KW_DATABASE || p.peek().Type == KW_SCHEMA {
			p.next()
			if _, err := p.parseIdentifier(); err != nil {
				return err
			}
			p.match(KW_EXTENDED)
			return nil
		}
		if isIdentifier(p.peek().Type) && p.peek().Type != KW_EXTENDED {
			if err := p.parsePartTypeExpr(); err != nil {
				return err
			}
		}
		p.match(KW_EXTENDED)
		return nil
	case KW_COMPACTIONS:
		p.next()
		if p.match(KW_COMPACT_ID) {
			if _, err := p.expect(TokEqual); err != nil {
				return err
			}
			_, err := p.expect(TokNumber)
			return err
		}
		if p.peek().Type == KW_DATABASE || p.peek().Type == KW_SCHEMA {
			p.next()
			if _, err := p.parseIdentifier(); err != nil {
				return err
			}
		} else if isIdentifier(p.peek().Type) && !p.isCompactionOptionStart() {
			if err := p.parsePartTypeExpr(); err != nil {
				return err
			}
		}
		if p.match(KW_POOL) {
			if _, err := p.expect(TokStringLiteral); err != nil {
				return err
			}
		}
		if p.match(KW_TYPE) {
			if _, err := p.expect(TokStringLiteral); err != nil {
				return err
			}
		}
		if p.match(KW_STATUS) {
			if _, err := p.expect(TokStringLiteral); err != nil {
				return err
			}
		}
		if p.peek().Type == KW_ORDER {
			if err := p.parseOrderByClause(); err != nil {
				return err
			}
		}
		if p.peek().Type == KW_LIMIT {
			return p.parseLimitClause()
		}
		return nil
	case KW_TRANSACTIONS, KW_ROLES, KW_DATACONNECTORS:
		p.next()
		return nil
	case KW_CONF:
		p.next()
		_, err := p.expect(TokStringLiteral)
		return err
	case KW_RESOURCE:
		p.next()
		if p.match(KW_PLANS) {
			return nil
		}
		if _, err := p.expect(KW_PLAN); err != nil {
			return err
		}
		_, err := p.parseIdentifier()
		return err
	case KW_ROLE:
		p.next()
		if _, err := p.expect(KW_GRANT); err != nil {
			return err
		}
		return p.parsePrincipalName()
	case KW_CURRENT:
		p.next()
		_, err := p.expect(KW_ROLES)
		return err
	case KW_GRANT:
		p.next()
		if p.peek().Type == KW_USER || p.peek().Type == KW_GROUP || p.peek().Type == KW_ROLE {
			if err := p.parsePrincipalName(); err != nil {
				return err
			}
		}
		if p.match(KW_ON) {
			if p.match(KW_ALL) {
				return nil
			}
			return p.parsePrivObject(true)
		}
		return nil
	case KW_PRINCIPALS:
		p.next()
		_, err := p.parseIdentifier()
		return err
	default:
		return p.errAt(p.peek(), fmt.Sprintf("unexpected token %q after SHOW", p.peek().Str))
	}
}

func (p *parser) isCompactionOptionStart() bool {
	switch p.peek().Type {
	case KW_POOL, KW_TYPE, KW_STATUS:
		return p.peekAt(1).Type == TokStringLiteral
	case KW_ORDER:
		return p.peekAt(1).Type == KW_BY
	case KW_LIMIT:
		return p.peekAt(1).Type == TokNumber
	}
	return false
}

func (p *parser) parseShowTablesRest() error {
	if _, err := p.expect(KW_TABLES); err != nil {
		return err
	}
	if p.peek().Type == KW_FROM || p.peek().Type == KW_IN {
		p.next()
		if _, err := p.parseIdentifier(); err != nil {
			return err
		}
	}
	if p.match(KW_WHERE) {
		if _, err := p.parseIdentifier(); err != nil {
			return err
		}
		if _, err := p.expect(TokEqual); err != nil {
			return err
		}
		_, err := p.expect(TokStringLiteral)
		return err
	}
	if p.match(KW_LIKE) {
		return p.parseShowStmtIdentifier()
	}
	if p.peek().Type == TokStringLiteral || isIdentifier(p.peek().Type) {
		return p.parseShowStmtIdentifier()
	}
	return nil
}

func (p *parser) parseShowFromLikeRest() error {
	if p.peek().Type == KW_FROM || p.peek().Type == KW_IN {
		p.next()
		if _, err := p.parseIdentifier(); err != nil {
			return err
		}
	}
	if p.match(KW_LIKE) {
		return p.parseShowStmtIdentifier()
	}
	if p.peek().Type == TokStringLiteral || isIdentifier(p.peek().Type) {
		return p.parseShowStmtIdentifier()
	}
	return nil
}

func (p *parser) parseShowStmtIdentifier() error {
	if p.match(TokStringLiteral) {
		return nil
	}
	_, err := p.parseIdentifier()
	return err
}

func (p *parser) parseMetastoreCheck() error {
	if _, err := p.expect(KW_MSCK); err != nil {
		return err
	}
	p.match(KW_REPAIR)
	if _, err := p.expect(KW_TABLE); err != nil {
		return err
	}
	if err := p.parseTableName(); err != nil {
		return err
	}
	if p.peek().Type == KW_ADD || p.peek().Type == KW_DROP || p.peek().Type == KW_SYNC {
		p.next()
		if _, err := p.expect(KW_PARTITIONS); err != nil {
			return err
		}
		if p.peek().Type == TokLParen {
			return p.parsePartitionSelectorSpec()
		}
	}
	return nil
}

// --- LOCK / UNLOCK / GRANT / REVOKE / SET ROLE / ABORT / KILL ---

func (p *parser) parseLockStatement() error {
	if _, err := p.expect(KW_LOCK); err != nil {
		return err
	}
	if p.peek().Type == KW_DATABASE || p.peek().Type == KW_SCHEMA {
		p.next()
		if _, err := p.parseIdentifier(); err != nil {
			return err
		}
		return p.parseLockMode()
	}
	if _, err := p.expect(KW_TABLE); err != nil {
		return err
	}
	if err := p.parseTableName(); err != nil {
		return err
	}
	if p.peek().Type == KW_PARTITION {
		if err := p.parsePartitionSpec(); err != nil {
			return err
		}
	}
	return p.parseLockMode()
}

func (p *parser) parseLockMode() error {
	if p.peek().Type == KW_SHARED || p.peek().Type == KW_EXCLUSIVE {
		p.next()
		return nil
	}
	return p.errAt(p.peek(), fmt.Sprintf("expected SHARED or EXCLUSIVE, got %q", p.peek().Str))
}

func (p *parser) parseUnlockStatement() error {
	if _, err := p.expect(KW_UNLOCK); err != nil {
		return err
	}
	if p.peek().Type == KW_DATABASE || p.peek().Type == KW_SCHEMA {
		p.next()
		_, err := p.parseIdentifier()
		return err
	}
	if _, err := p.expect(KW_TABLE); err != nil {
		return err
	}
	if err := p.parseTableName(); err != nil {
		return err
	}
	if p.peek().Type == KW_PARTITION {
		return p.parsePartitionSpec()
	}
	return nil
}

func (p *parser) parseGrantPrivileges() error {
	if _, err := p.expect(KW_GRANT); err != nil {
		return err
	}
	if err := p.parsePrivilegeList(); err != nil {
		return err
	}
	if p.match(KW_ON) {
		if err := p.parsePrivObject(false); err != nil {
			return err
		}
	}
	if _, err := p.expect(KW_TO); err != nil {
		return err
	}
	if err := p.parsePrincipalSpecification(); err != nil {
		return err
	}
	if p.match(KW_WITH) {
		if _, err := p.expect(KW_GRANT); err != nil {
			return err
		}
		_, err := p.expect(KW_OPTION)
		return err
	}
	return nil
}

func (p *parser) parseRevokePrivileges() error {
	if _, err := p.expect(KW_REVOKE); err != nil {
		return err
	}
	if p.peek().Type == KW_GRANT && p.peekAt(1).Type == KW_OPTION && p.peekAt(2).Type == KW_FOR {
		p.pos += 3
	}
	if err := p.parsePrivilegeList(); err != nil {
		return err
	}
	if p.match(KW_ON) {
		if err := p.parsePrivObject(false); err != nil {
			return err
		}
	}
	if _, err := p.expect(KW_FROM); err != nil {
		return err
	}
	return p.parsePrincipalSpecification()
}

func (p *parser) parseGrantRole() error {
	if _, err := p.expect(KW_GRANT); err != nil {
		return err
	}
	p.match(KW_ROLE)
	if err := p.parseAliasList(); err != nil {
		return err
	}
	if _, err := p.expect(KW_TO); err != nil {
		return err
	}
	if err := p.parsePrincipalSpecification(); err != nil {
		return err
	}
	if p.match(KW_WITH) {
		if _, err := p.expect(KW_ADMIN); err != nil {
			return err
		}
		_, err := p.expect(KW_OPTION)
		return err
	}
	return nil
}

func (p *parser) parseRevokeRole() error {
	if _, err := p.expect(KW_REVOKE); err != nil {
		return err
	}
	if p.peek().Type == KW_ADMIN && p.peekAt(1).Type == KW_OPTION && p.peekAt(2).Type == KW_FOR {
		p.pos += 3
	}
	p.match(KW_ROLE)
	if err := p.parseAliasList(); err != nil {
		return err
	}
	if _, err := p.expect(KW_FROM); err != nil {
		return err
	}
	return p.parsePrincipalSpecification()
}

func (p *parser) parsePrivilegeList() error {
	if err := p.parsePrivilegeDef(); err != nil {
		return err
	}
	for p.match(TokComma) {
		if err := p.parsePrivilegeDef(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parsePrivilegeDef() error {
	switch p.peek().Type {
	case KW_ALL, KW_ALTER, KW_UPDATE, KW_CREATE, KW_DROP, KW_LOCK, KW_SELECT, KW_SHOW_DATABASE, KW_INSERT, KW_DELETE:
		p.next()
	default:
		return p.errAt(p.peek(), fmt.Sprintf("expected privilege type, got %q", p.peek().Str))
	}
	if p.peek().Type == TokLParen {
		return p.parseColumnParenthesesList()
	}
	return nil
}

func (p *parser) parsePrivObject(allowCols bool) error {
	switch p.peek().Type {
	case KW_DATABASE, KW_SCHEMA:
		p.next()
		_, err := p.parseIdentifier()
		return err
	case KW_URI:
		if p.peekAt(1).Type == TokStringLiteral {
			p.pos += 2
			return nil
		}
	case KW_SERVER:
		if isIdentifier(p.peekAt(1).Type) {
			p.pos += 2
			return nil
		}
	}
	p.match(KW_TABLE)
	if err := p.parseTableName(); err != nil {
		return err
	}
	if allowCols && p.peek().Type == TokLParen {
		if err := p.parseColumnParenthesesList(); err != nil {
			return err
		}
	}
	if p.peek().Type == KW_PARTITION {
		return p.parsePartitionSpec()
	}
	return nil
}

func (p *parser) parsePrincipalSpecification() error {
	if err := p.parsePrincipalName(); err != nil {
		return err
	}
	for p.match(TokComma) {
		if err := p.parsePrincipalName(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parsePrincipalName() error {
	switch p.peek().Type {
	case KW_USER, KW_GROUP, KW_ROLE:
		p.next()
		_, err := p.parseIdentifier()
		return err
	default:
		return p.errAt(p.peek(), fmt.Sprintf("expected USER, GROUP, or ROLE, got %q", p.peek().Str))
	}
}

func (p *parser) parseSetRole() error {
	if _, err := p.expect(KW_SET); err != nil {
		return err
	}
	if _, err := p.expect(KW_ROLE); err != nil {
		return err
	}
	if p.peek().Type == KW_ALL || p.peek().Type == KW_NONE {
		p.next()
		return nil
	}
	_, err := p.parseIdentifier()
	return err
}

func (p *parser) parseAbortStatement() error {
	if _, err := p.expect(KW_ABORT); err != nil {
		return err
	}
	if p.peek().Type != KW_TRANSACTIONS && p.peek().Type != KW_COMPACTIONS {
		return p.errAt(p.peek(), fmt.Sprintf("expected TRANSACTIONS or COMPACTIONS after ABORT, got %q", p.peek().Str))
	}
	p.next()
	if _, err := p.expect(TokNumber); err != nil {
		return err
	}
	for p.match(TokNumber) {
	}
	return nil
}

func (p *parser) parseKillQueryStatement() error {
	if _, err := p.expect(KW_KILL); err != nil {
		return err
	}
	if _, err := p.expect(KW_QUERY); err != nil {
		return err
	}
	if _, err := p.expect(TokStringLiteral); err != nil {
		return err
	}
	for p.match(TokStringLiteral) {
	}
	return nil
}
