package parser

func (p *Parser) parseCreateStatement() (Statement, error) {
	start := p.cur().Start
	cp := p.save()

	// Dispatch by looking past CREATE [OR REPLACE] [GLOBAL] [TEMPORARY] [EXTERNAL].
	next := p.peek(1).Kind
	switch next {
	case kwNAMESPACE, kwDATABASE, kwSCHEMA:
		return p.parseCreateNamespace()
	case kwROLE:
		return p.parseUnsupportedHiveNativeCommand()
	case kwINDEX:
		// Bare `CREATE INDEX` is in Spark's unsupportedHiveNativeCommands; any arguments
		// must follow Spark's full CREATE INDEX grammar cleanly.
		if p.isStatementEndAt(2) {
			return p.parseUnsupportedHiveNativeCommand()
		}
		return p.parseCreateIndex()
	case kwTEMPORARY:
		if p.peek(2).Kind == kwMACRO {
			return p.parseUnsupportedHiveNativeCommand()
		}
	}

	// CREATE OR REPLACE TABLE ... -> replaceTable
	if next == kwOR && p.peek(2).Kind == kwREPLACE && p.peek(3).Kind == kwTABLE {
		return p.parseReplaceTable()
	}

	// Scan modifiers to decide TABLE vs VIEW vs FUNCTION.
	i := 1
	if p.peek(i).Kind == kwOR && p.peek(i+1).Kind == kwREPLACE {
		i += 2
	}
	if p.peek(i).Kind == kwGLOBAL {
		i++
	}
	if p.peek(i).Kind == kwTEMPORARY {
		i++
	}
	if p.peek(i).Kind == kwEXTERNAL {
		i++
	}

	switch p.peek(i).Kind {
	case kwTABLE:
		// Try CREATE TABLE ... LIKE first if LIKE appears after table name, else regular CREATE TABLE.
		stmt, err := p.parseCreateTableLike()
		if err == nil && p.isStatementEnd() {
			return stmt, nil
		}
		p.restore(cp)
		return p.parseCreateTable()
	case kwVIEW:
		// Try createView, fallback to createTempViewUsing.
		stmt, err := p.parseCreateView()
		if err == nil && p.isStatementEnd() {
			return stmt, nil
		}
		p.restore(cp)
		return p.parseCreateTempViewUsing()
	case kwFUNCTION:
		return p.parseCreateFunction()
	}

	return nil, &ParseError{Message: "unexpected CREATE statement", Position: start}
}

func (p *Parser) matchIfNotExists() bool {
	if p.cur().Kind == kwIF && p.isNotAt(1) && p.peek(2).Kind == kwEXISTS {
		p.advance()
		p.advance()
		p.advance()
		return true
	}
	return false
}

func (p *Parser) matchIfExists() bool {
	if p.cur().Kind == kwIF && p.peek(1).Kind == kwEXISTS {
		p.advance()
		p.advance()
		return true
	}
	return false
}

func (p *Parser) parseCreateNamespace() (Statement, error) {
	start := p.cur().Start
	if _, err := p.expect(kwCREATE); err != nil {
		return nil, err
	}
	if !p.matchNamespace() {
		return nil, p.errorf("expected NAMESPACE, DATABASE, or SCHEMA")
	}
	p.matchIfNotExists()
	target, err := p.parseIdentifierReference()
	if err != nil {
		return nil, err
	}
	for {
		switch p.cur().Kind {
		case kwCOMMENT:
			if err := p.parseCommentSpec(); err != nil {
				return nil, err
			}
		case kwLOCATION:
			if err := p.parseLocationSpec(); err != nil {
				return nil, err
			}
		case kwWITH:
			p.advance()
			if !p.match(kwDBPROPERTIES, kwPROPERTIES) {
				return nil, p.errorf("expected DBPROPERTIES or PROPERTIES after WITH")
			}
			if err := p.parsePropertyList(); err != nil {
				return nil, err
			}
		default:
			return &ParsedStmt{Kind: StmtCreateNamespace, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
		}
	}
}

func (p *Parser) parseCreateTable() (Statement, error) {
	start := p.cur().Start
	if _, err := p.expect(kwCREATE); err != nil {
		return nil, err
	}
	p.match(kwTEMPORARY)
	p.match(kwEXTERNAL)
	if _, err := p.expect(kwTABLE); err != nil {
		return nil, err
	}
	p.matchIfNotExists()
	target, err := p.parseIdentifierReference()
	if err != nil {
		return nil, err
	}
	if err := p.parseCreateOrReplaceTableBody(); err != nil {
		return nil, err
	}
	return &ParsedStmt{Kind: StmtCreateTable, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
}

func (p *Parser) parseReplaceTable() (Statement, error) {
	start := p.cur().Start
	if p.match(kwCREATE) {
		if _, err := p.expect(kwOR); err != nil {
			return nil, err
		}
	}
	if _, err := p.expect(kwREPLACE); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwTABLE); err != nil {
		return nil, err
	}
	target, err := p.parseIdentifierReference()
	if err != nil {
		return nil, err
	}
	if err := p.parseCreateOrReplaceTableBody(); err != nil {
		return nil, err
	}
	return &ParsedStmt{Kind: StmtReplaceTable, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
}

func (p *Parser) parseCreateOrReplaceTableBody() error {
	// Optional (colDefinitionList), careful not to mistake (query) for colDefinitionList!
	if p.cur().Kind == tokLParen {
		cp := p.save()
		p.advance()
		if err := p.parseColDefinitionList(); err == nil && p.match(tokRParen) {
			// parsed colDefinitionList
		} else {
			p.restore(cp)
		}
	}
	if p.cur().Kind == kwUSING {
		if err := p.parseTableProvider(); err != nil {
			return err
		}
	}
	if err := p.parseCreateTableClauses(); err != nil {
		return err
	}
	if p.match(kwAS) {
		return p.parseQuery()
	}
	if p.looksLikeQueryStart() {
		return p.parseQuery()
	}
	return nil
}

func (p *Parser) parseCreateTableClauses() error {
	for {
		switch p.cur().Kind {
		case kwOPTIONS:
			p.advance()
			if err := p.parseExpressionPropertyList(); err != nil {
				return err
			}
		case kwPARTITIONED:
			p.advance()
			if _, err := p.expect(kwBY); err != nil {
				return err
			}
			if err := p.parsePartitionFieldList(); err != nil {
				return err
			}
		case kwSKEWED:
			if err := p.parseSkewSpec(); err != nil {
				return err
			}
		case kwCLUSTER:
			if err := p.parseClusterBySpec(); err != nil {
				return err
			}
		case kwCLUSTERED:
			if err := p.parseBucketSpec(); err != nil {
				return err
			}
		case kwROW:
			if err := p.parseRowFormat(); err != nil {
				return err
			}
		case kwSTORED:
			if err := p.parseCreateFileFormat(); err != nil {
				return err
			}
		case kwLOCATION:
			if err := p.parseLocationSpec(); err != nil {
				return err
			}
		case kwCOMMENT:
			if err := p.parseCommentSpec(); err != nil {
				return err
			}
		case kwDEFAULT:
			if err := p.parseCollationSpec(); err != nil {
				return err
			}
		case kwTBLPROPERTIES:
			p.advance()
			if err := p.parsePropertyList(); err != nil {
				return err
			}
		default:
			return nil
		}
	}
}

func (p *Parser) parseCreateTableLike() (Statement, error) {
	start := p.cur().Start
	if _, err := p.expect(kwCREATE); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwTABLE); err != nil {
		return nil, err
	}
	p.matchIfNotExists()
	target, err := p.parseTableIdentifier()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(kwLIKE); err != nil {
		return nil, err
	}
	if _, err := p.parseTableIdentifier(); err != nil {
		return nil, err
	}
	for {
		switch p.cur().Kind {
		case kwUSING:
			if err := p.parseTableProvider(); err != nil {
				return nil, err
			}
		case kwROW:
			if err := p.parseRowFormat(); err != nil {
				return nil, err
			}
		case kwSTORED:
			if err := p.parseCreateFileFormat(); err != nil {
				return nil, err
			}
		case kwLOCATION:
			if err := p.parseLocationSpec(); err != nil {
				return nil, err
			}
		case kwTBLPROPERTIES:
			p.advance()
			if err := p.parsePropertyList(); err != nil {
				return nil, err
			}
		default:
			return &ParsedStmt{Kind: StmtCreateTableLike, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
		}
	}
}

func (p *Parser) parseCreateView() (Statement, error) {
	start := p.cur().Start
	if _, err := p.expect(kwCREATE); err != nil {
		return nil, err
	}
	if p.match(kwOR) {
		if _, err := p.expect(kwREPLACE); err != nil {
			return nil, err
		}
	}
	if p.match(kwGLOBAL) {
		if _, err := p.expect(kwTEMPORARY); err != nil {
			return nil, err
		}
	} else {
		p.match(kwTEMPORARY)
	}
	if _, err := p.expect(kwVIEW); err != nil {
		return nil, err
	}
	p.matchIfNotExists()
	target, err := p.parseIdentifierReference()
	if err != nil {
		return nil, err
	}
	if p.cur().Kind == tokLParen {
		if err := p.parseIdentifierCommentList(); err != nil {
			return nil, err
		}
	}
	for {
		switch p.cur().Kind {
		case kwCOMMENT:
			if err := p.parseCommentSpec(); err != nil {
				return nil, err
			}
		case kwWITH:
			if err := p.parseSchemaBinding(); err != nil {
				return nil, err
			}
		case kwDEFAULT:
			if err := p.parseCollationSpec(); err != nil {
				return nil, err
			}
		case kwPARTITIONED:
			p.advance()
			if _, err := p.expect(kwON); err != nil {
				return nil, err
			}
			if err := p.parseIdentifierList(); err != nil {
				return nil, err
			}
		case kwTBLPROPERTIES:
			p.advance()
			if err := p.parsePropertyList(); err != nil {
				return nil, err
			}
		default:
			goto expectAs
		}
	}
expectAs:
	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}
	if err := p.parseQuery(); err != nil {
		return nil, err
	}
	return &ParsedStmt{Kind: StmtCreateView, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
}

func (p *Parser) parseCreateTempViewUsing() (Statement, error) {
	start := p.cur().Start
	if _, err := p.expect(kwCREATE); err != nil {
		return nil, err
	}
	if p.match(kwOR) {
		if _, err := p.expect(kwREPLACE); err != nil {
			return nil, err
		}
	}
	p.match(kwGLOBAL)
	if _, err := p.expect(kwTEMPORARY); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwVIEW); err != nil {
		return nil, err
	}
	target, err := p.parseTableIdentifier()
	if err != nil {
		return nil, err
	}
	if p.match(tokLParen) {
		if err := p.parseColTypeList(); err != nil {
			return nil, err
		}
		if _, err := p.expect(tokRParen); err != nil {
			return nil, err
		}
	}
	if err := p.parseTableProvider(); err != nil {
		return nil, err
	}
	if p.match(kwOPTIONS) {
		if err := p.parsePropertyList(); err != nil {
			return nil, err
		}
	}
	return &ParsedStmt{Kind: StmtCreateView, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
}

func (p *Parser) parseCreateFunction() (Statement, error) {
	start := p.cur().Start
	if _, err := p.expect(kwCREATE); err != nil {
		return nil, err
	}
	if p.match(kwOR) {
		if _, err := p.expect(kwREPLACE); err != nil {
			return nil, err
		}
	}
	p.match(kwTEMPORARY)
	if _, err := p.expect(kwFUNCTION); err != nil {
		return nil, err
	}
	p.matchIfNotExists()
	target, err := p.parseIdentifierReference()
	if err != nil {
		return nil, err
	}

	// SQL UDF: (colDefinitionList?) (RETURNS (dataType | TABLE (colTypeList)))? routineCharacteristics RETURN (query | expression)
	if p.match(tokLParen) {
		if !p.match(tokRParen) {
			if err := p.parseColDefinitionList(); err != nil {
				return nil, err
			}
			if _, err := p.expect(tokRParen); err != nil {
				return nil, err
			}
		}
		if p.cur().Kind == kwRETURNS && p.peek(1).Kind != kwNULL {
			p.advance()
			if p.cur().Kind == kwTABLE && p.peek(1).Kind == tokLParen {
				p.advance()
				p.advance()
				if err := p.parseColTypeList(); err != nil {
					return nil, err
				}
				if _, err := p.expect(tokRParen); err != nil {
					return nil, err
				}
			} else {
				if err := p.parseDataType(); err != nil {
					return nil, err
				}
			}
		}
		if err := p.parseRoutineCharacteristics(); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwRETURN); err != nil {
			return nil, err
		}
		cp := p.save()
		if err := p.parseQuery(); err == nil && p.isStatementEnd() {
			return &ParsedStmt{Kind: StmtCreateFunction, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
		}
		p.restore(cp)
		if err := p.parseExpression(); err != nil {
			return nil, err
		}
		return &ParsedStmt{Kind: StmtCreateFunction, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
	}

	// Class-based function: AS stringLit (USING resource (, resource)*)?
	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}
	if err := p.parseStringLit(); err != nil {
		return nil, err
	}
	if p.match(kwUSING) {
		for {
			if _, err := p.parseIdentifier(); err != nil {
				return nil, err
			}
			if err := p.parseStringLit(); err != nil {
				return nil, err
			}
			if !p.match(tokComma) {
				break
			}
		}
	}
	return &ParsedStmt{Kind: StmtCreateFunction, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
}

func (p *Parser) parseCreateIndex() (Statement, error) {
	start := p.cur().Start
	if _, err := p.expect(kwCREATE); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwINDEX); err != nil {
		return nil, err
	}
	p.matchIfNotExists()
	target, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	if p.cur().Kind == kwTABLE && p.peek(1).Kind != tokDot && p.peek(1).Kind != tokLParen && p.peek(1).Kind != kwUSING && p.peek(1).Kind != tokEOF && p.peek(1).Kind != tokSemicolon {
		p.advance()
	}
	if _, err := p.parseIdentifierReference(); err != nil {
		return nil, err
	}
	if p.match(kwUSING) {
		if _, err := p.parseIdentifier(); err != nil {
			return nil, err
		}
	}
	if _, err := p.expect(tokLParen); err != nil {
		return nil, err
	}
	for {
		if _, err := p.parseMultipartIdentifier(); err != nil {
			return nil, err
		}
		if p.match(kwOPTIONS) {
			if err := p.parsePropertyList(); err != nil {
				return nil, err
			}
		}
		if !p.match(tokComma) {
			break
		}
	}
	if _, err := p.expect(tokRParen); err != nil {
		return nil, err
	}
	if p.match(kwOPTIONS) {
		if err := p.parsePropertyList(); err != nil {
			return nil, err
		}
	}
	return &ParsedStmt{Kind: StmtCreateIndex, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
}

func (p *Parser) parseAlterStatement() (Statement, error) {
	switch p.peek(1).Kind {
	case kwNAMESPACE, kwDATABASE, kwSCHEMA:
		return p.parseAlterNamespace()
	case kwTABLE:
		return p.parseAlterTable()
	case kwVIEW:
		return p.parseAlterView()
	case kwINDEX:
		return p.parseUnsupportedHiveNativeCommand()
	default:
		return nil, p.errorf("unexpected token %q after ALTER", p.peek(1).Text)
	}
}

func (p *Parser) parseAlterNamespace() (Statement, error) {
	start := p.cur().Start
	if _, err := p.expect(kwALTER); err != nil {
		return nil, err
	}
	if !p.matchNamespace() {
		return nil, p.errorf("expected NAMESPACE, DATABASE, or SCHEMA")
	}
	target, err := p.parseIdentifierReference()
	if err != nil {
		return nil, err
	}
	if p.match(kwSET) {
		if p.cur().Kind == kwLOCATION {
			if err := p.parseLocationSpec(); err != nil {
				return nil, err
			}
			return &ParsedStmt{Kind: StmtAlterNamespace, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
		}
		if !p.match(kwDBPROPERTIES, kwPROPERTIES) {
			return nil, p.errorf("expected DBPROPERTIES, PROPERTIES, or LOCATION")
		}
		if err := p.parsePropertyList(); err != nil {
			return nil, err
		}
		return &ParsedStmt{Kind: StmtAlterNamespace, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
	}
	if p.match(kwUNSET) {
		if !p.match(kwDBPROPERTIES, kwPROPERTIES) {
			return nil, p.errorf("expected DBPROPERTIES or PROPERTIES")
		}
		if err := p.parsePropertyKeyList(); err != nil {
			return nil, err
		}
		return &ParsedStmt{Kind: StmtAlterNamespace, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
	}
	return nil, p.errorf("expected SET or UNSET in ALTER NAMESPACE")
}

func (p *Parser) parseAlterView() (Statement, error) {
	start := p.cur().Start
	if _, err := p.expect(kwALTER); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwVIEW); err != nil {
		return nil, err
	}
	target, err := p.parseIdentifierReference()
	if err != nil {
		return nil, err
	}
	switch p.cur().Kind {
	case kwRENAME:
		p.advance()
		if _, err := p.expect(kwTO); err != nil {
			return nil, err
		}
		if _, err := p.parseMultipartIdentifier(); err != nil {
			return nil, err
		}
	case kwSET:
		p.advance()
		if _, err := p.expect(kwTBLPROPERTIES); err != nil {
			return nil, err
		}
		if err := p.parsePropertyList(); err != nil {
			return nil, err
		}
	case kwUNSET:
		p.advance()
		if _, err := p.expect(kwTBLPROPERTIES); err != nil {
			return nil, err
		}
		if p.match(kwIF) {
			if _, err := p.expect(kwEXISTS); err != nil {
				return nil, err
			}
		}
		if err := p.parsePropertyKeyList(); err != nil {
			return nil, err
		}
	case kwADD:
		p.advance()
		if p.match(kwIF) {
			if !p.matchNot() {
				return nil, p.errorf("expected NOT after IF")
			}
			if _, err := p.expect(kwEXISTS); err != nil {
				return nil, err
			}
		}
		if err := p.parsePartitionSpecLocation(); err != nil {
			return nil, err
		}
		for p.cur().Kind == kwPARTITION {
			if err := p.parsePartitionSpecLocation(); err != nil {
				return nil, err
			}
		}
	case kwDROP:
		p.advance()
		if p.match(kwIF) {
			if _, err := p.expect(kwEXISTS); err != nil {
				return nil, err
			}
		}
		if err := p.parsePartitionSpec(); err != nil {
			return nil, err
		}
		for p.match(tokComma) {
			if err := p.parsePartitionSpec(); err != nil {
				return nil, err
			}
		}
		p.match(kwPURGE)
	case kwWITH:
		// Could be WITH SCHEMA ... (schemaBinding) or WITH cte ... (alterViewQuery)!
		cp := p.save()
		if err := p.parseSchemaBinding(); err == nil && p.isStatementEnd() {
			return &ParsedStmt{Kind: StmtAlterView, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
		}
		p.restore(cp)
		if err := p.parseQuery(); err != nil {
			return nil, err
		}
	case kwAS:
		p.advance()
		if err := p.parseQuery(); err != nil {
			return nil, err
		}
	default:
		if err := p.parseQuery(); err != nil {
			return nil, err
		}
	}
	return &ParsedStmt{Kind: StmtAlterView, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
}

func (p *Parser) parseAlterTable() (Statement, error) {
	start := p.cur().Start
	cpAll := p.save()

	if _, err := p.expect(kwALTER); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwTABLE); err != nil {
		return nil, err
	}
	target, err := p.parseIdentifierReference()
	if err != nil {
		return nil, err
	}

	// Optional leading partitionSpec (used by hiveChangeColumn, hiveReplaceColumns,
	// setTableSerDe, renameTablePartition, setTableLocation, and unsupportedHiveNativeCommands).
	if p.cur().Kind == kwPARTITION {
		if err := p.parsePartitionSpec(); err != nil {
			return nil, err
		}
		switch p.cur().Kind {
		case kwRENAME:
			p.advance()
			if _, err := p.expect(kwTO); err != nil {
				return nil, err
			}
			if err := p.parsePartitionSpec(); err != nil {
				return nil, err
			}
			return &ParsedStmt{Kind: StmtAlterTable, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
		case kwCHANGE:
			if err := p.parseHiveChangeColumnTail(); err != nil {
				return nil, err
			}
			return &ParsedStmt{Kind: StmtAlterTable, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
		case kwREPLACE:
			if p.peek(1).Kind == kwCOLUMNS && p.isStatementEndAt(2) {
				p.restore(cpAll)
				return p.parseUnsupportedHiveNativeCommand()
			}
			if err := p.parseHiveReplaceColumnsTail(); err != nil {
				return nil, err
			}
			return &ParsedStmt{Kind: StmtAlterTable, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
		case kwSET:
			if err := p.parseAlterTableSetWithOptionalPart(cpAll, true); err != nil {
				return nil, err
			}
			return &ParsedStmt{Kind: StmtAlterTable, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
		case kwCOMPACT, kwCONCATENATE:
			p.restore(cpAll)
			return p.parseUnsupportedHiveNativeCommand()
		default:
			return nil, p.errorf("unexpected token %q after PARTITION in ALTER TABLE", p.cur().Text)
		}
	}

	switch p.cur().Kind {
	case kwADD:
		p.advance()
		if p.match(kwCOLUMN, kwCOLUMNS) {
			if p.match(tokLParen) {
				if err := p.parseQualifiedColTypeWithPositionList(); err != nil {
					return nil, err
				}
				if _, err := p.expect(tokRParen); err != nil {
					return nil, err
				}
			} else {
				if err := p.parseQualifiedColTypeWithPositionList(); err != nil {
					return nil, err
				}
			}
			return &ParsedStmt{Kind: StmtAlterTable, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
		}
		if p.match(kwIF) {
			if !p.matchNot() {
				return nil, p.errorf("expected NOT after IF")
			}
			if _, err := p.expect(kwEXISTS); err != nil {
				return nil, err
			}
		}
		if err := p.parsePartitionSpecLocation(); err != nil {
			return nil, err
		}
		for p.cur().Kind == kwPARTITION {
			if err := p.parsePartitionSpecLocation(); err != nil {
				return nil, err
			}
		}

	case kwRENAME:
		p.advance()
		if p.match(kwCOLUMN) {
			if _, err := p.parseMultipartIdentifier(); err != nil {
				return nil, err
			}
			if _, err := p.expect(kwTO); err != nil {
				return nil, err
			}
			if _, err := p.parseErrorCapturingIdentifier(); err != nil {
				return nil, err
			}
		} else {
			if _, err := p.expect(kwTO); err != nil {
				return nil, err
			}
			if _, err := p.parseMultipartIdentifier(); err != nil {
				return nil, err
			}
		}

	case kwDROP:
		p.advance()
		if p.match(kwCOLUMN, kwCOLUMNS) {
			p.matchIfExists()
			if p.match(tokLParen) {
				if err := p.parseMultipartIdentifierList(); err != nil {
					return nil, err
				}
				if _, err := p.expect(tokRParen); err != nil {
					return nil, err
				}
			} else {
				if err := p.parseMultipartIdentifierList(); err != nil {
					return nil, err
				}
			}
			return &ParsedStmt{Kind: StmtAlterTable, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
		}
		if p.match(kwIF) {
			if _, err := p.expect(kwEXISTS); err != nil {
				return nil, err
			}
		}
		if err := p.parsePartitionSpec(); err != nil {
			return nil, err
		}
		for p.match(tokComma) {
			if err := p.parsePartitionSpec(); err != nil {
				return nil, err
			}
		}
		p.match(kwPURGE)

	case kwSET:
		if err := p.parseAlterTableSetWithOptionalPart(cpAll, false); err != nil {
			return nil, err
		}

	case kwUNSET:
		p.advance()
		if _, err := p.expect(kwTBLPROPERTIES); err != nil {
			return nil, err
		}
		if p.match(kwIF) {
			if _, err := p.expect(kwEXISTS); err != nil {
				return nil, err
			}
		}
		if err := p.parsePropertyKeyList(); err != nil {
			return nil, err
		}

	case kwALTER, kwCHANGE:
		cpCol := p.save()
		// Try alterTableAlterColumn first, then hiveChangeColumn (if CHANGE).
		if err := p.parseAlterColumnSpecListTail(); err == nil && p.isStatementEnd() {
			return &ParsedStmt{Kind: StmtAlterTable, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
		}
		p.restore(cpCol)
		if p.cur().Kind == kwCHANGE {
			if err := p.parseHiveChangeColumnTail(); err != nil {
				return nil, err
			}
			return &ParsedStmt{Kind: StmtAlterTable, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
		}
		return nil, p.errorf("invalid ALTER TABLE ALTER COLUMN clause")

	case kwREPLACE:
		if p.peek(1).Kind == kwCOLUMNS && p.isStatementEndAt(2) {
			p.restore(cpAll)
			return p.parseUnsupportedHiveNativeCommand()
		}
		if err := p.parseHiveReplaceColumnsTail(); err != nil {
			return nil, err
		}
		return &ParsedStmt{Kind: StmtAlterTable, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil

	case kwRECOVER:
		p.advance()
		if _, err := p.expect(kwPARTITIONS); err != nil {
			return nil, err
		}

	case kwCLUSTER:
		if p.peek(1).Kind == kwBY && p.peek(2).Kind == kwNONE {
			p.advance()
			p.advance()
			p.advance()
		} else {
			if err := p.parseClusterBySpec(); err != nil {
				return nil, err
			}
		}

	case kwDEFAULT:
		if err := p.parseCollationSpec(); err != nil {
			return nil, err
		}

	case kwNOT, kwCLUSTERED, kwSKEWED, kwEXCHANGE, kwARCHIVE, kwUNARCHIVE, kwTOUCH, kwCOMPACT, kwCONCATENATE:
		p.restore(cpAll)
		return p.parseUnsupportedHiveNativeCommand()

	default:
		return nil, p.errorf("unexpected token %q in ALTER TABLE", p.cur().Text)
	}

	return &ParsedStmt{Kind: StmtAlterTable, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
}

func (p *Parser) parseAlterTableSetWithOptionalPart(cpAll checkpoint, hasPartition bool) error {
	if _, err := p.expect(kwSET); err != nil {
		return err
	}
	switch p.cur().Kind {
	case kwTBLPROPERTIES:
		if hasPartition {
			return p.errorf("SET TBLPROPERTIES is not allowed after PARTITION in ALTER TABLE")
		}
		p.advance()
		return p.parsePropertyList()
	case kwSERDE:
		p.advance()
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
	case kwSERDEPROPERTIES:
		p.advance()
		return p.parsePropertyList()
	case kwLOCATION:
		return p.parseLocationSpec()
	case kwSKEWED, kwFILEFORMAT:
		p.restore(cpAll)
		_, err := p.parseUnsupportedHiveNativeCommand()
		return err
	default:
		return p.errorf("unexpected token %q after ALTER TABLE ... SET", p.cur().Text)
	}
}

func (p *Parser) parseAlterColumnSpecListTail() error {
	if !p.match(kwALTER, kwCHANGE) {
		return p.errorf("expected ALTER or CHANGE")
	}
	if p.cur().Kind == kwCOLUMN {
		cp := p.save()
		p.advance()
		if err := p.parseAlterColumnSpecs(); err == nil && p.isStatementEnd() {
			return nil
		}
		p.restore(cp)
	}
	return p.parseAlterColumnSpecs()
}

func (p *Parser) parseAlterColumnSpecs() error {
	if err := p.parseAlterColumnSpec(); err != nil {
		return err
	}
	for p.match(tokComma) {
		if err := p.parseAlterColumnSpec(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parseAlterColumnSpec() error {
	if _, err := p.parseMultipartIdentifier(); err != nil {
		return err
	}
	switch p.cur().Kind {
	case kwTYPE:
		p.advance()
		return p.parseDataType()
	case kwCOMMENT:
		return p.parseCommentSpec()
	case kwFIRST, kwAFTER:
		return p.parseColPosition()
	case kwSET:
		p.advance()
		if p.cur().Kind == kwDEFAULT {
			return p.parseDefaultExpression()
		}
		if !p.matchNot() {
			return p.errorf("expected NOT NULL or DEFAULT after SET")
		}
		_, err := p.expect(kwNULL)
		return err
	case kwDROP:
		p.advance()
		if p.match(kwDEFAULT) {
			return nil
		}
		if !p.matchNot() {
			return p.errorf("expected NOT NULL or DEFAULT after DROP")
		}
		_, err := p.expect(kwNULL)
		return err
	}
	return nil
}

func (p *Parser) parseHiveChangeColumnTail() error {
	if _, err := p.expect(kwCHANGE); err != nil {
		return err
	}
	if p.cur().Kind == kwCOLUMN {
		cp := p.save()
		p.advance()
		if err := p.parseHiveChangeColumnBody(); err == nil && p.isStatementEnd() {
			return nil
		}
		p.restore(cp)
	}
	return p.parseHiveChangeColumnBody()
}

func (p *Parser) parseHiveChangeColumnBody() error {
	if _, err := p.parseMultipartIdentifier(); err != nil {
		return err
	}
	if err := p.parseColType(); err != nil {
		return err
	}
	if p.cur().Kind == kwFIRST || p.cur().Kind == kwAFTER {
		return p.parseColPosition()
	}
	return nil
}

func (p *Parser) parseHiveReplaceColumnsTail() error {
	if _, err := p.expect(kwREPLACE); err != nil {
		return err
	}
	if _, err := p.expect(kwCOLUMNS); err != nil {
		return err
	}
	if _, err := p.expect(tokLParen); err != nil {
		return err
	}
	if err := p.parseQualifiedColTypeWithPositionList(); err != nil {
		return err
	}
	_, err := p.expect(tokRParen)
	return err
}

func (p *Parser) parseDropStatement() (Statement, error) {
	start := p.cur().Start
	cp := p.save()
	if _, err := p.expect(kwDROP); err != nil {
		return nil, err
	}
	matchIfExists := func() {
		if p.cur().Kind == kwIF && p.peek(1).Kind == kwEXISTS {
			p.advance()
			p.advance()
		}
	}
	switch p.cur().Kind {
	case kwNAMESPACE, kwDATABASE, kwSCHEMA:
		p.advance()
		matchIfExists()
		target, err := p.parseIdentifierReference()
		if err != nil {
			return nil, err
		}
		p.match(kwRESTRICT, kwCASCADE)
		return &ParsedStmt{Kind: StmtDropNamespace, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil

	case kwTABLE:
		p.advance()
		matchIfExists()
		target, err := p.parseIdentifierReference()
		if err != nil {
			return nil, err
		}
		p.match(kwPURGE)
		return &ParsedStmt{Kind: StmtDropTable, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil

	case kwVIEW:
		p.advance()
		matchIfExists()
		target, err := p.parseIdentifierReference()
		if err != nil {
			return nil, err
		}
		return &ParsedStmt{Kind: StmtDropView, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil

	case kwTEMPORARY:
		p.advance()
		if p.cur().Kind == kwMACRO {
			p.restore(cp)
			return p.parseUnsupportedHiveNativeCommand()
		}
		if p.matchVariable() {
			matchIfExists()
			target, err := p.parseIdentifierReference()
			if err != nil {
				return nil, err
			}
			return &ParsedStmt{Kind: StmtDropVariable, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
		}
		if _, err := p.expect(kwFUNCTION); err != nil {
			return nil, err
		}
		matchIfExists()
		target, err := p.parseIdentifierReference()
		if err != nil {
			return nil, err
		}
		return &ParsedStmt{Kind: StmtDropFunction, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil

	case kwFUNCTION:
		p.advance()
		matchIfExists()
		target, err := p.parseIdentifierReference()
		if err != nil {
			return nil, err
		}
		return &ParsedStmt{Kind: StmtDropFunction, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil

	case kwINDEX:
		if p.isStatementEndAt(1) {
			p.restore(cp)
			return p.parseUnsupportedHiveNativeCommand()
		}
		p.advance()
		matchIfExists()
		target, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(kwON); err != nil {
			return nil, err
		}
		if p.cur().Kind == kwTABLE && p.peek(1).Kind != tokDot && p.peek(1).Kind != tokEOF && p.peek(1).Kind != tokSemicolon {
			p.advance()
		}
		if _, err := p.parseIdentifierReference(); err != nil {
			return nil, err
		}
		return &ParsedStmt{Kind: StmtDropIndex, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil

	case kwROLE:
		p.restore(cp)
		return p.parseUnsupportedHiveNativeCommand()

	default:
		return nil, p.errorf("unexpected token %q after DROP", p.cur().Text)
	}
}

func (p *Parser) parseCreateVariable() (Statement, error) {
	start := p.cur().Start
	if _, err := p.expect(kwDECLARE); err != nil {
		return nil, err
	}
	if p.cur().Kind == kwOR && p.peek(1).Kind == kwREPLACE {
		p.advance()
		p.advance()
	}
	if (p.cur().Kind == kwVARIABLE || p.cur().Kind == kwVAR) &&
		!p.isStatementEndAt(1) && p.peek(1).Kind != tokDot && p.peek(1).Kind != tokEq &&
		!(p.peek(1).Kind == kwDEFAULT && p.peek(2).Kind != tokEq && p.peek(2).Kind != kwDEFAULT && p.peek(2).Kind != tokDot && !p.isStatementEndAt(2)) {
		cp := p.save()
		p.advance()
		if target, err := p.parseCreateVariableBody(); err == nil && p.isStatementEnd() {
			return &ParsedStmt{Kind: StmtDeclare, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
		}
		p.restore(cp)
	}
	target, err := p.parseCreateVariableBody()
	if err != nil {
		return nil, err
	}
	return &ParsedStmt{Kind: StmtDeclare, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
}

func (p *Parser) parseCreateVariableBody() (string, error) {
	target, err := p.parseIdentifierReference()
	if err != nil {
		return "", err
	}
	// Optional dataType followed by optional (DEFAULT | =) expression.
	if !p.isStatementEnd() {
		if p.cur().Kind == tokEq {
			p.advance()
			if err := p.parseExpression(); err != nil {
				return "", err
			}
		} else if p.cur().Kind == kwDEFAULT && !p.isStatementEndAt(1) {
			p.advance()
			if err := p.parseExpression(); err != nil {
				return "", err
			}
		} else {
			if err := p.parseDataType(); err != nil {
				return "", err
			}
			if p.match(kwDEFAULT, tokEq) {
				if err := p.parseExpression(); err != nil {
					return "", err
				}
			}
		}
	}
	return target, nil
}
