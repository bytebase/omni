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
		return nil, p.errorf("CREATE ROLE is not valid")
	case kwINDEX:
		return p.parseCreateIndex()
	case kwTEMPORARY:
		if p.peek(2).Kind == kwMACRO {
			return nil, p.errorf("CREATE TEMPORARY MACRO is not valid")
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
	if p.cur().Kind == kwIF && p.peek(1).Kind == kwNOT && p.peek(2).Kind == kwEXISTS {
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
	var seenComment, seenLocation, seenProps bool
	for {
		switch p.cur().Kind {
		case kwCOMMENT:
			if seenComment {
				return nil, p.errorf("duplicate COMMENT clause in CREATE NAMESPACE")
			}
			seenComment = true
			if err := p.parseCommentSpec(); err != nil {
				return nil, err
			}
		case kwLOCATION:
			if seenLocation {
				return nil, p.errorf("duplicate LOCATION clause in CREATE NAMESPACE")
			}
			seenLocation = true
			if err := p.parseLocationSpec(); err != nil {
				return nil, err
			}
		case kwWITH:
			if seenProps {
				return nil, p.errorf("duplicate PROPERTIES clause in CREATE NAMESPACE")
			}
			seenProps = true
			p.advance()
			if !p.match(kwDBPROPERTIES, kwPROPERTIES) {
				return nil, p.errorf("expected DBPROPERTIES or PROPERTIES after WITH")
			}
			if err := p.parseNamespacePropertyList(); err != nil {
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
	isTemporary := p.match(kwTEMPORARY)
	isExternal := p.match(kwEXTERNAL)
	if isTemporary && isExternal {
		return nil, p.errorf("CREATE TEMPORARY EXTERNAL TABLE is not supported")
	}
	if _, err := p.expect(kwTABLE); err != nil {
		return nil, err
	}
	ifNotExists := p.matchIfNotExists()
	if isTemporary && ifNotExists {
		return nil, p.errorf("CREATE TEMPORARY TABLE cannot specify IF NOT EXISTS")
	}
	target, err := p.parseIdentifierReference()
	if err != nil {
		return nil, err
	}
	if err := p.parseCreateOrReplaceTableBody(isTemporary); err != nil {
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
	if err := p.parseCreateOrReplaceTableBody(false); err != nil {
		return nil, err
	}
	return &ParsedStmt{Kind: StmtReplaceTable, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
}

func (p *Parser) parseCreateOrReplaceTableBody(isTemporary bool) error {
	// Optional (colDefinitionList), careful not to mistake (query) for colDefinitionList!
	hasColDefs := false
	if p.cur().Kind == tokLParen {
		cp := p.save()
		p.advance()
		if err := p.parseColDefinitionList(); err == nil && p.match(tokRParen) {
			hasColDefs = true
		} else {
			p.restore(cp)
		}
	}
	hasUsing := false
	if p.cur().Kind == kwUSING {
		hasUsing = true
		if err := p.parseTableProvider(); err != nil {
			return err
		}
	}
	hasTypedPartitionCols, err := p.parseCreateTableClauses(hasUsing)
	if err != nil {
		return err
	}
	hasQuery := false
	if p.match(kwAS) {
		hasQuery = true
		if hasColDefs {
			return p.errorf("CREATE/REPLACE TABLE AS SELECT cannot specify a column list")
		}
		if hasTypedPartitionCols {
			return p.errorf("CREATE/REPLACE TABLE AS SELECT cannot specify partition column types")
		}
		if isTemporary {
			return p.errorf("CREATE TEMPORARY TABLE AS SELECT is not supported")
		}
		if err := p.parseQuery(); err != nil {
			return err
		}
	} else if p.looksLikeQueryStart() {
		hasQuery = true
		if hasColDefs {
			return p.errorf("CREATE/REPLACE TABLE AS SELECT cannot specify a column list")
		}
		if hasTypedPartitionCols {
			return p.errorf("CREATE/REPLACE TABLE AS SELECT cannot specify partition column types")
		}
		if isTemporary {
			return p.errorf("CREATE TEMPORARY TABLE AS SELECT is not supported")
		}
		if err := p.parseQuery(); err != nil {
			return err
		}
	}
	if isTemporary && !hasUsing && !hasQuery {
		return p.errorf("CREATE TEMPORARY TABLE requires a USING provider")
	}
	return nil
}

func (p *Parser) parseCreateTableClauses(hasUsing bool) (bool, error) {
	var (
		hasOptions            bool
		hasOptionsPath        bool
		hasPartitionedBy      bool
		hasTypedPartitionCols bool
		hasClusterBy          bool
		hasClusteredBy        bool
		hasRowFormat          bool
		rowFmt                rowFormatKind
		hasStoredAs           bool
		storedAsFormat        string
		hasLocation           bool
		hasComment            bool
		hasCollation          bool
		hasTblProps           bool
	)
	for {
		switch p.cur().Kind {
		case kwOPTIONS:
			if hasOptions {
				return false, p.errorf("duplicate OPTIONS clause in CREATE TABLE")
			}
			hasOptions = true
			p.advance()
			optPath, err := p.parseExpressionPropertyList()
			if err != nil {
				return false, err
			}
			if optPath {
				hasOptionsPath = true
			}
		case kwPARTITIONED:
			if hasPartitionedBy {
				return false, p.errorf("duplicate PARTITIONED BY clause in CREATE TABLE")
			}
			if hasClusterBy {
				return false, p.errorf("PARTITIONED BY cannot be combined with CLUSTER BY")
			}
			hasPartitionedBy = true
			p.advance()
			if _, err := p.expect(kwBY); err != nil {
				return false, err
			}
			typed, err := p.parsePartitionFieldList()
			if err != nil {
				return false, err
			}
			if typed {
				hasTypedPartitionCols = true
			}
		case kwSKEWED:
			return false, p.errorf("CREATE TABLE ... SKEWED BY is not valid")
		case kwCLUSTER:
			if hasClusterBy {
				return false, p.errorf("duplicate CLUSTER BY clause in CREATE TABLE")
			}
			if hasPartitionedBy {
				return false, p.errorf("CLUSTER BY cannot be combined with PARTITIONED BY")
			}
			if hasClusteredBy {
				return false, p.errorf("CLUSTER BY cannot be combined with CLUSTERED BY")
			}
			hasClusterBy = true
			if err := p.parseClusterBySpec(); err != nil {
				return false, err
			}
		case kwCLUSTERED:
			if hasClusteredBy {
				return false, p.errorf("duplicate CLUSTERED BY clause in CREATE TABLE")
			}
			if hasClusterBy {
				return false, p.errorf("CLUSTERED BY cannot be combined with CLUSTER BY")
			}
			hasClusteredBy = true
			if err := p.parseBucketSpec(); err != nil {
				return false, err
			}
		case kwROW:
			if hasUsing {
				return false, p.errorf("ROW FORMAT cannot be combined with USING")
			}
			if hasRowFormat {
				return false, p.errorf("duplicate ROW FORMAT clause in CREATE TABLE")
			}
			hasRowFormat = true
			rf, err := p.parseRowFormatKind()
			if err != nil {
				return false, err
			}
			rowFmt = rf
		case kwSTORED:
			if hasUsing {
				return false, p.errorf("STORED AS cannot be combined with USING")
			}
			if hasStoredAs {
				return false, p.errorf("duplicate STORED AS clause in CREATE TABLE")
			}
			hasStoredAs = true
			ff, err := p.parseCreateFileFormat(false)
			if err != nil {
				return false, err
			}
			storedAsFormat = ff
		case kwLOCATION:
			if hasLocation {
				return false, p.errorf("duplicate LOCATION clause in CREATE TABLE")
			}
			hasLocation = true
			if err := p.parseLocationSpec(); err != nil {
				return false, err
			}
		case kwCOMMENT:
			if hasComment {
				return false, p.errorf("duplicate COMMENT clause in CREATE TABLE")
			}
			hasComment = true
			if err := p.parseCommentSpec(); err != nil {
				return false, err
			}
		case kwDEFAULT:
			if hasCollation {
				return false, p.errorf("duplicate DEFAULT COLLATION clause in CREATE TABLE")
			}
			hasCollation = true
			if err := p.parseCollationSpec(); err != nil {
				return false, err
			}
		case kwTBLPROPERTIES:
			if hasTblProps {
				return false, p.errorf("duplicate TBLPROPERTIES clause in CREATE TABLE")
			}
			hasTblProps = true
			p.advance()
			if err := p.parseTablePropertyList(); err != nil {
				return false, err
			}
		default:
			if hasOptionsPath && hasLocation {
				return false, p.errorf("duplicated table paths: LOCATION and OPTIONS path cannot both be specified")
			}
			if msg := validateRowAndFileFormat(rowFmt, storedAsFormat); msg != "" {
				return false, p.errorf("%s", msg)
			}
			return hasTypedPartitionCols, nil
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
	var hasUsing, hasRowFormat, hasStoredAs, hasLocation, hasTblProps bool
	var rowFmt rowFormatKind
	var storedAsFormat string
	for {
		switch p.cur().Kind {
		case kwUSING:
			if hasUsing {
				return nil, p.errorf("duplicate USING clause in CREATE TABLE LIKE")
			}
			hasUsing = true
			if err := p.parseTableProvider(); err != nil {
				return nil, err
			}
		case kwROW:
			if hasRowFormat {
				return nil, p.errorf("duplicate ROW FORMAT clause in CREATE TABLE LIKE")
			}
			hasRowFormat = true
			rf, err := p.parseRowFormatKind()
			if err != nil {
				return nil, err
			}
			rowFmt = rf
		case kwSTORED:
			if hasStoredAs {
				return nil, p.errorf("duplicate STORED AS clause in CREATE TABLE LIKE")
			}
			hasStoredAs = true
			ff, err := p.parseCreateFileFormat(true)
			if err != nil {
				return nil, err
			}
			storedAsFormat = ff
		case kwLOCATION:
			if hasLocation {
				return nil, p.errorf("duplicate LOCATION clause in CREATE TABLE LIKE")
			}
			hasLocation = true
			if err := p.parseLocationSpec(); err != nil {
				return nil, err
			}
		case kwTBLPROPERTIES:
			if hasTblProps {
				return nil, p.errorf("duplicate TBLPROPERTIES clause in CREATE TABLE LIKE")
			}
			hasTblProps = true
			p.advance()
			if err := p.parseTablePropertyList(); err != nil {
				return nil, err
			}
		default:
			if hasUsing && (hasRowFormat || hasStoredAs) {
				return nil, p.errorf("CREATE TABLE LIKE cannot combine USING with ROW FORMAT or STORED AS")
			}
			if msg := validateRowAndFileFormat(rowFmt, storedAsFormat); msg != "" {
				return nil, p.errorf("%s", msg)
			}
			return &ParsedStmt{Kind: StmtCreateTableLike, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
		}
	}
}

func (p *Parser) parseCreateView() (Statement, error) {
	start := p.cur().Start
	if _, err := p.expect(kwCREATE); err != nil {
		return nil, err
	}
	isReplace := false
	if p.match(kwOR) {
		if _, err := p.expect(kwREPLACE); err != nil {
			return nil, err
		}
		isReplace = true
	}
	isTemp := false
	if p.match(kwGLOBAL) {
		if _, err := p.expect(kwTEMPORARY); err != nil {
			return nil, err
		}
		isTemp = true
	} else if p.match(kwTEMPORARY) {
		isTemp = true
	}
	if _, err := p.expect(kwVIEW); err != nil {
		return nil, err
	}
	ifNotExists := p.matchIfNotExists()
	if isReplace && ifNotExists {
		return nil, p.errorf("CREATE OR REPLACE VIEW cannot specify IF NOT EXISTS")
	}
	if isTemp && ifNotExists {
		return nil, p.errorf("CREATE TEMPORARY VIEW cannot specify IF NOT EXISTS")
	}
	target, err := p.parseIdentifierReference()
	if err != nil {
		return nil, err
	}
	if isTemp && isQualifiedIdent(target) {
		return nil, p.errorf("temporary view name %q cannot be qualified", target)
	}
	if p.cur().Kind == tokLParen {
		if err := p.parseIdentifierCommentList(); err != nil {
			return nil, err
		}
	}
	var seenComment, seenBinding, seenCollation, seenTblProps bool
	for {
		switch p.cur().Kind {
		case kwCOMMENT:
			if seenComment {
				return nil, p.errorf("duplicate COMMENT clause in CREATE VIEW")
			}
			seenComment = true
			if err := p.parseCommentSpec(); err != nil {
				return nil, err
			}
		case kwWITH:
			if isTemp {
				return nil, p.errorf("WITH SCHEMA clause is not supported on TEMPORARY VIEW")
			}
			if seenBinding {
				return nil, p.errorf("duplicate WITH SCHEMA clause in CREATE VIEW")
			}
			seenBinding = true
			if err := p.parseSchemaBinding(); err != nil {
				return nil, err
			}
		case kwDEFAULT:
			if seenCollation {
				return nil, p.errorf("duplicate DEFAULT COLLATION clause in CREATE VIEW")
			}
			seenCollation = true
			if err := p.parseCollationSpec(); err != nil {
				return nil, err
			}
		case kwPARTITIONED:
			return nil, p.errorf("CREATE VIEW ... PARTITIONED ON is not valid")
		case kwTBLPROPERTIES:
			if isTemp {
				return nil, p.errorf("TBLPROPERTIES is not allowed on CREATE TEMPORARY VIEW")
			}
			if seenTblProps {
				return nil, p.errorf("duplicate TBLPROPERTIES clause in CREATE VIEW")
			}
			seenTblProps = true
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
	prevDisallow := p.disallowParamMarkers
	p.disallowParamMarkers = true
	err = p.parseQuery()
	p.disallowParamMarkers = prevDisallow
	if err != nil {
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
	if isQualifiedIdent(target) {
		return nil, p.errorf("temporary view name %q cannot be qualified", target)
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
	isReplace := false
	if p.match(kwOR) {
		if _, err := p.expect(kwREPLACE); err != nil {
			return nil, err
		}
		isReplace = true
	}
	isTemp := p.match(kwTEMPORARY)
	if _, err := p.expect(kwFUNCTION); err != nil {
		return nil, err
	}
	ifNotExists := p.matchIfNotExists()
	if isReplace && ifNotExists {
		return nil, p.errorf("CREATE OR REPLACE FUNCTION cannot specify IF NOT EXISTS")
	}
	if isTemp && ifNotExists {
		return nil, p.errorf("CREATE TEMPORARY FUNCTION cannot specify IF NOT EXISTS")
	}
	target, err := p.parseIdentifierReference()
	if err != nil {
		return nil, err
	}
	if isTemp && isQualifiedIdent(target) {
		return nil, p.errorf("temporary function name %q cannot be qualified", target)
	}

	// SQL UDF: (colDefinitionList?) (RETURNS (dataType | TABLE (colTypeList)))? routineCharacteristics RETURN (query | expression)
	if p.match(tokLParen) {
		// Spark v4.0.5 AstBuilder defers SQL UDF parameter/return dataType
		// validation to catalog analysis rather than parsePlan, while rejecting
		// GENERATED column expressions in function parameter definitions.
		prevRelax := p.relaxFunctionDataTypes
		prevNoGen := p.disallowGeneratedCols
		p.relaxFunctionDataTypes = true
		p.disallowGeneratedCols = true
		if !p.match(tokRParen) {
			if err := p.parseColDefinitionList(); err != nil {
				p.relaxFunctionDataTypes = prevRelax
				p.disallowGeneratedCols = prevNoGen
				return nil, err
			}
			if _, err := p.expect(tokRParen); err != nil {
				p.relaxFunctionDataTypes = prevRelax
				p.disallowGeneratedCols = prevNoGen
				return nil, err
			}
		}
		p.disallowGeneratedCols = prevNoGen
		if p.cur().Kind == kwRETURNS && p.peek(1).Kind != kwNULL {
			p.advance()
			if p.cur().Kind == kwTABLE && p.peek(1).Kind == tokLParen {
				p.advance()
				p.advance()
				if err := p.parseColTypeList(); err != nil {
					p.relaxFunctionDataTypes = prevRelax
					return nil, err
				}
				if _, err := p.expect(tokRParen); err != nil {
					p.relaxFunctionDataTypes = prevRelax
					return nil, err
				}
			} else {
				if err := p.parseDataType(); err != nil {
					p.relaxFunctionDataTypes = prevRelax
					return nil, err
				}
			}
		}
		p.relaxFunctionDataTypes = prevRelax
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
			resType, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			if !isValidFunctionResourceType(resType) {
				return nil, p.errorf("unsupported function resource type %q; expected JAR, FILE, or ARCHIVE", resType)
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
	if p.cur().Kind == kwIDENTIFIER && p.peek(1).Kind == tokLParen {
		closeOff := p.matchingCloseOffset(1)
		if closeOff > 0 && p.peek(closeOff+1).Kind != tokLParen && p.peek(closeOff+1).Kind != kwUSING {
			if _, err := p.parseMultipartIdentifier(); err != nil {
				return nil, err
			}
		} else if _, err := p.parseIdentifierReference(); err != nil {
			return nil, err
		}
	} else if _, err := p.parseIdentifierReference(); err != nil {
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
		if err := p.parseNamespacePropertyList(); err != nil {
			return nil, err
		}
		return &ParsedStmt{Kind: StmtAlterNamespace, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
	}
	if p.match(kwUNSET) {
		if !p.match(kwDBPROPERTIES, kwPROPERTIES) {
			return nil, p.errorf("expected DBPROPERTIES or PROPERTIES")
		}
		if err := p.parseNamespacePropertyKeyList(); err != nil {
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
		if err := p.parsePartitionSpec(true); err != nil {
			return nil, err
		}
		for p.match(tokComma) {
			if err := p.parsePartitionSpec(true); err != nil {
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
	// setTableSerDe, renameTablePartition, setTableLocation).
	if p.cur().Kind == kwPARTITION {
		if err := p.parsePartitionSpec(true); err != nil {
			return nil, err
		}
		switch p.cur().Kind {
		case kwRENAME:
			p.advance()
			if _, err := p.expect(kwTO); err != nil {
				return nil, err
			}
			if err := p.parsePartitionSpec(true); err != nil {
				return nil, err
			}
			return &ParsedStmt{Kind: StmtAlterTable, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
		case kwCHANGE:
			if err := p.parseHiveChangeColumnTail(); err != nil {
				return nil, err
			}
			return &ParsedStmt{Kind: StmtAlterTable, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
		case kwREPLACE:
			if err := p.parseHiveReplaceColumnsTail(); err != nil {
				return nil, err
			}
			return &ParsedStmt{Kind: StmtAlterTable, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
		case kwSET:
			if err := p.parseAlterTableSetWithOptionalPart(true); err != nil {
				return nil, err
			}
			return &ParsedStmt{Kind: StmtAlterTable, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
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
			if _, err := p.expect(kwNOT); err != nil {
				return nil, err
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
		if err := p.parsePartitionSpec(true); err != nil {
			return nil, err
		}
		for p.match(tokComma) {
			if err := p.parsePartitionSpec(true); err != nil {
				return nil, err
			}
		}
		p.match(kwPURGE)

	case kwSET:
		if err := p.parseAlterTableSetWithOptionalPart(false); err != nil {
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
		if err := p.parseTablePropertyKeyList(); err != nil {
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

	default:
		return nil, p.errorf("unexpected token %q in ALTER TABLE", p.cur().Text)
	}

	return &ParsedStmt{Kind: StmtAlterTable, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
}

func (p *Parser) parseAlterTableSetWithOptionalPart(hasPartition bool) error {
	if _, err := p.expect(kwSET); err != nil {
		return err
	}
	switch p.cur().Kind {
	case kwTBLPROPERTIES:
		if hasPartition {
			return p.errorf("SET TBLPROPERTIES is not allowed after PARTITION in ALTER TABLE")
		}
		p.advance()
		return p.parseTablePropertyList()
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
	return p.errorf("expected TYPE, COMMENT, FIRST, AFTER, SET, or DROP in ALTER COLUMN specification, got %q", p.cur().Text)
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
	oldCol, err := p.parseMultipartIdentifier()
	if err != nil {
		return err
	}
	newCol, err := p.parseErrorCapturingIdentifier()
	if err != nil {
		return err
	}
	// Spark v4.0.5 AstBuilder.visitHiveChangeColumn rejects column renaming
	// (newCol must match the leaf name of oldCol, ignoring dots inside backticks).
	parts := splitMultipartIdent(oldCol)
	leaf := parts[len(parts)-1]
	if exactIdentName(leaf) != exactIdentName(newCol) {
		return p.errorf("ALTER TABLE CHANGE COLUMN cannot rename column %q to %q", oldCol, newCol)
	}
	if err := p.parseDataType(); err != nil {
		return err
	}
	if p.cur().Kind == kwNOT && p.peek(1).Kind == kwNULL {
		p.advance()
		p.advance()
	}
	if p.cur().Kind == kwCOMMENT {
		if err := p.parseCommentSpec(); err != nil {
			return err
		}
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
	// Spark v4.0.5 AstBuilder.visitHiveReplaceColumns rejects nested column
	// names, NOT NULL, DEFAULT, and column positions (FIRST/AFTER), accepting
	// only simple column name + dataType + optional COMMENT.
	for {
		if _, err := p.parseErrorCapturingIdentifier(); err != nil {
			return err
		}
		if p.cur().Kind == tokDot {
			return p.errorf("nested column names are not supported in REPLACE COLUMNS")
		}
		if err := p.parseDataType(); err != nil {
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

func (p *Parser) parseDropStatement() (Statement, error) {
	start := p.cur().Start
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
			return nil, p.errorf("DROP TEMPORARY MACRO is not valid")
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
	if p.cur().Kind == kwVARIABLE || p.cur().Kind == kwVAR {
		// Match Spark ANTLR's greedy `variable?` consumption: VAR/VARIABLE acts
		// as the variable name only when the next token forces it (e.g. `=`,
		// `.`, `DEFAULT <literal/expr not starting with ident or '('>`, or a
		// parameterized/complex dataType like `VARCHAR(10)` or `ARRAY<INT>`).
		next := p.peek(1).Kind
		after := p.peek(2).Kind
		isVarAsName := next == tokEq || next == tokDot ||
			(next == kwDEFAULT && !p.isStatementEndAt(2) && !after.isIdentifier() && after != tokLParen && after != tokEq && after != tokDot) ||
			(next.isIdentifier() && (after == tokLParen || after == tokLT))
		if !isVarAsName {
			p.advance()
		}
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
	// Spark requires either a dataType or a DEFAULT/= initializer (or both).
	if p.isStatementEnd() {
		return "", p.errorf("variable %q requires a data type or DEFAULT expression", target)
	}
	parseVarInit := func() error {
		prevDisallow := p.disallowParamMarkers
		p.disallowParamMarkers = true
		err := p.parseExpression()
		p.disallowParamMarkers = prevDisallow
		return err
	}
	if p.cur().Kind == tokEq {
		p.advance()
		if err := parseVarInit(); err != nil {
			return "", err
		}
		return target, nil
	}
	if p.cur().Kind == kwDEFAULT && !p.isStatementEndAt(1) {
		p.advance()
		if err := parseVarInit(); err != nil {
			return "", err
		}
		return target, nil
	}
	if err := p.parseDataType(); err != nil {
		return "", err
	}
	if p.match(kwDEFAULT, tokEq) {
		if err := parseVarInit(); err != nil {
			return "", err
		}
	}
	return target, nil
}
