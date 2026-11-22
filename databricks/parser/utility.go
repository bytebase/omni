package parser

func (p *Parser) parseUse() (Statement, error) {
	start := p.cur().Start
	if _, err := p.expect(kwUSE); err != nil {
		return nil, err
	}
	if p.isNamespace() && !p.isStatementEndAt(1) {
		cp := p.save()
		p.advance()
		if target, err := p.parseIdentifierReference(); err == nil && p.isStatementEnd() {
			return &ParsedStmt{Kind: StmtUse, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
		}
		p.restore(cp)
	}
	target, err := p.parseIdentifierReference()
	if err != nil {
		return nil, err
	}
	return &ParsedStmt{Kind: StmtUse, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
}

func (p *Parser) isStatementEndAt(offset int) bool {
	k := p.peek(offset).Kind
	return k == tokEOF || k == tokSemicolon
}

func (p *Parser) parseSetCatalog() (Statement, error) {
	start := p.cur().Start
	if _, err := p.expect(kwSET); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwCATALOG); err != nil {
		return nil, err
	}
	target, err := p.parseCatalogIdentifierReference()
	if err != nil {
		return nil, err
	}
	return &ParsedStmt{Kind: StmtSetCatalog, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
}

// parseSetResetStatement implements setResetStatement from SqlBaseParser.g4:
//
//	setResetStatement
//	    : SET ROLE .*?
//	    | SET TIME ZONE interval
//	    | SET TIME ZONE timezone
//	    | SET TIME ZONE .*?
//	    | SET variable assignmentList
//	    | SET variable ( multipartIdentifierList ) = ( query )
//	    | SET configKey = configValue
//	    | SET configKey (= .*?)?
//	    | SET .*? = configValue
//	    | SET .*?
//	    | RESET configKey
//	    | RESET .*?
//	    ;
func (p *Parser) parseSetResetStatement() (Statement, error) {
	start := p.cur().Start
	if p.match(kwRESET) {
		if p.isStatementEnd() {
			return &ParsedStmt{Kind: StmtReset, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
		}
		if p.cur().Kind == tokBackquoted || p.cur().Kind == tokDoubleQuoted {
			p.advance()
			return &ParsedStmt{Kind: StmtReset, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
		}
		// Unquoted dotted configuration key (e.g. RESET spark.sql.shuffle.partitions)
		if _, err := p.parseConfigKeyUnquoted(); err != nil {
			return nil, err
		}
		return &ParsedStmt{Kind: StmtReset, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
	}

	if _, err := p.expect(kwSET); err != nil {
		return nil, err
	}
	if p.isStatementEnd() {
		return &ParsedStmt{Kind: StmtSet, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
	}

	// SET ROLE ...
	if p.match(kwROLE) {
		for !p.isStatementEnd() {
			p.advance()
		}
		return &ParsedStmt{Kind: StmtSet, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
	}

	// SET TIME ZONE (interval | timezone | unquotedIdent*)
	if p.cur().Kind == kwTIME && p.peek(1).Kind == kwZONE {
		p.advance()
		p.advance()
		if p.isStatementEnd() {
			return &ParsedStmt{Kind: StmtSet, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
		}
		if p.cur().Kind == kwINTERVAL {
			if err := p.parseInterval(); err != nil {
				return nil, err
			}
			return &ParsedStmt{Kind: StmtSet, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
		}
		if p.isStringLit() || p.cur().Kind == kwLOCAL {
			p.advance()
			return &ParsedStmt{Kind: StmtSet, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
		}
		if p.cur().Kind == tokIdent {
			for p.cur().Kind == tokIdent {
				p.advance()
			}
			return &ParsedStmt{Kind: StmtSet, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
		}
		return nil, p.errorf("invalid timezone in SET TIME ZONE")
	}

	// SET (VARIABLE | VAR) (assignmentList | (multipartIdentifierList) = (query))
	if p.match(kwVARIABLE, kwVAR) {
		if p.match(tokLParen) {
			if err := p.parseMultipartIdentifierList(); err != nil {
				return nil, err
			}
			if _, err := p.expect(tokRParen); err != nil {
				return nil, err
			}
			if _, err := p.expect(tokEq); err != nil {
				return nil, err
			}
			if _, err := p.expect(tokLParen); err != nil {
				return nil, err
			}
			if err := p.parseQuery(); err != nil {
				return nil, err
			}
			if _, err := p.expect(tokRParen); err != nil {
				return nil, err
			}
			return &ParsedStmt{Kind: StmtSet, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
		}
		if err := p.parseAssignmentList(); err != nil {
			return nil, err
		}
		return &ParsedStmt{Kind: StmtSet, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
	}

	// SET configKey (= configValue | = .*?)?
	if p.cur().Kind == tokBackquoted || p.cur().Kind == tokDoubleQuoted {
		p.advance()
		if p.match(tokEq) {
			if p.cur().Kind == tokBackquoted {
				p.advance()
			} else {
				for !p.isStatementEnd() {
					p.advance()
				}
			}
		}
		return &ParsedStmt{Kind: StmtSet, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
	}

	// SET = `configValue`
	if p.match(tokEq) {
		if _, err := p.expect(tokBackquoted); err != nil {
			return nil, err
		}
		return &ParsedStmt{Kind: StmtSet, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
	}

	// SET -v (Spark verbose listing)
	if p.cur().Kind == tokMinus && p.peek(1).Kind == tokIdent && isWord(p.peek(1), "v") {
		p.advance()
		p.advance()
		return &ParsedStmt{Kind: StmtSet, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
	}

	// SET <unquoted config key> [= value]
	if p.cur().Kind == tokIdent {
		if _, err := p.parseConfigKeyUnquoted(); err != nil {
			return nil, err
		}
		if p.match(tokEq) {
			if p.cur().Kind == tokBackquoted {
				p.advance()
			} else {
				for !p.isStatementEnd() {
					p.advance()
				}
			}
		} else {
			for p.cur().Kind == tokIdent {
				p.advance()
			}
		}
		return &ParsedStmt{Kind: StmtSet, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
	}

	return nil, p.errorf("unexpected token %q in SET statement", p.cur().Text)
}

func (p *Parser) parseConfigKeyUnquoted() (bool, error) {
	if p.cur().Kind != tokIdent {
		return false, p.errorf("expected configuration key, got %q", p.cur().Text)
	}
	p.advance()
	hasDot := false
	for p.match(tokDot) {
		hasDot = true
		if !p.cur().Kind.isIdentifier() {
			return hasDot, p.errorf("expected identifier after '.' in configuration key")
		}
		p.advance()
	}
	return hasDot, nil
}

func (p *Parser) parseAnalyze() (Statement, error) {
	start := p.cur().Start
	if _, err := p.expect(kwANALYZE); err != nil {
		return nil, err
	}
	if p.match(kwTABLES) {
		if p.match(kwFROM, kwIN) {
			if _, err := p.parseIdentifierReference(); err != nil {
				return nil, err
			}
		}
		if _, err := p.expect(kwCOMPUTE); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwSTATISTICS); err != nil {
			return nil, err
		}
		if !p.isStatementEnd() {
			if p.cur().Kind == tokIdent && isWord(p.cur(), "NOSCAN") {
				p.advance()
			} else {
				return nil, p.errorf("expected NOSCAN after COMPUTE STATISTICS, got %q", p.cur().Text)
			}
		}
		return &ParsedStmt{Kind: StmtAnalyze, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
	}
	if _, err := p.expect(kwTABLE); err != nil {
		return nil, err
	}
	target, err := p.parseIdentifierReference()
	if err != nil {
		return nil, err
	}
	if p.cur().Kind == kwPARTITION {
		if err := p.parsePartitionSpec(); err != nil {
			return nil, err
		}
	}
	if _, err := p.expect(kwCOMPUTE); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwSTATISTICS); err != nil {
		return nil, err
	}
	if p.cur().Kind == kwFOR {
		cp := p.save()
		p.advance()
		if p.match(kwALL) {
			if _, err := p.expect(kwCOLUMNS); err == nil {
				return &ParsedStmt{Kind: StmtAnalyze, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
			}
		} else if p.match(kwCOLUMNS) {
			if err := p.parseIdentifierSeq(); err == nil {
				return &ParsedStmt{Kind: StmtAnalyze, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
			}
		}
		p.restore(cp)
	}
	if !p.isStatementEnd() {
		if p.cur().Kind == tokIdent && isWord(p.cur(), "NOSCAN") {
			p.advance()
		} else {
			return nil, p.errorf("expected NOSCAN after COMPUTE STATISTICS, got %q", p.cur().Text)
		}
	}
	return &ParsedStmt{Kind: StmtAnalyze, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
}

func (p *Parser) parseExplain() (Statement, error) {
	start := p.cur().Start
	if _, err := p.expect(kwEXPLAIN); err != nil {
		return nil, err
	}
	p.match(kwLOGICAL, kwFORMATTED, kwEXTENDED, kwCODEGEN, kwCOST)
	if _, err := p.parseStatementOrSetReset(); err != nil {
		return nil, err
	}
	return &ParsedStmt{Kind: StmtExplain, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
}

func (p *Parser) parseShow() (Statement, error) {
	start := p.cur().Start
	cp := p.save()
	if _, err := p.expect(kwSHOW); err != nil {
		return nil, err
	}
	switch p.cur().Kind {
	case kwNAMESPACES, kwDATABASES, kwSCHEMAS:
		p.advance()
		if p.match(kwFROM, kwIN) {
			if _, err := p.parseMultipartIdentifier(); err != nil {
				return nil, err
			}
		}
		if p.match(kwLIKE) {
			if err := p.parseStringLit(); err != nil {
				return nil, err
			}
		} else if p.isStringLit() {
			p.advance()
		}

	case kwTABLES:
		p.advance()
		if p.match(kwFROM, kwIN) {
			if _, err := p.parseIdentifierReference(); err != nil {
				return nil, err
			}
		}
		if p.match(kwLIKE) {
			if err := p.parseStringLit(); err != nil {
				return nil, err
			}
		} else if p.isStringLit() {
			p.advance()
		}

	case kwTABLE:
		p.advance()
		if _, err := p.expect(kwEXTENDED); err != nil {
			return nil, err
		}
		if p.match(kwFROM, kwIN) {
			if _, err := p.parseIdentifierReference(); err != nil {
				return nil, err
			}
		}
		if _, err := p.expect(kwLIKE); err != nil {
			return nil, err
		}
		if err := p.parseStringLit(); err != nil {
			return nil, err
		}
		if p.cur().Kind == kwPARTITION {
			if err := p.parsePartitionSpec(); err != nil {
				return nil, err
			}
		}

	case kwTBLPROPERTIES:
		p.advance()
		if _, err := p.parseIdentifierReference(); err != nil {
			return nil, err
		}
		if p.match(tokLParen) {
			if err := p.parsePropertyKey(); err != nil {
				return nil, err
			}
			if _, err := p.expect(tokRParen); err != nil {
				return nil, err
			}
		}

	case kwCOLUMNS:
		p.advance()
		if !p.match(kwFROM, kwIN) {
			return nil, p.errorf("expected FROM or IN after SHOW COLUMNS")
		}
		if _, err := p.parseIdentifierReference(); err != nil {
			return nil, err
		}
		if p.match(kwFROM, kwIN) {
			if _, err := p.parseMultipartIdentifier(); err != nil {
				return nil, err
			}
		}

	case kwVIEWS:
		p.advance()
		if p.match(kwFROM, kwIN) {
			if _, err := p.parseIdentifierReference(); err != nil {
				return nil, err
			}
		}
		if p.match(kwLIKE) {
			if err := p.parseStringLit(); err != nil {
				return nil, err
			}
		} else if p.isStringLit() {
			p.advance()
		}

	case kwPARTITIONS:
		p.advance()
		if _, err := p.parseIdentifierReference(); err != nil {
			return nil, err
		}
		if p.cur().Kind == kwPARTITION {
			if err := p.parsePartitionSpec(); err != nil {
				return nil, err
			}
		}

	case kwCREATE:
		p.advance()
		if _, err := p.expect(kwTABLE); err != nil {
			return nil, err
		}
		if p.isStatementEnd() {
			return &ParsedStmt{Kind: StmtHiveNative, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
		}
		if _, err := p.parseIdentifierReference(); err != nil {
			return nil, err
		}
		if p.match(kwAS) {
			if _, err := p.expect(kwSERDE); err != nil {
				return nil, err
			}
		}

	case kwCURRENT:
		p.advance()
		if p.cur().Kind == kwROLES {
			p.restore(cp)
			return p.parseUnsupportedHiveNativeCommand()
		}
		if !p.matchNamespace() {
			return nil, p.errorf("expected NAMESPACE, DATABASE, or SCHEMA after SHOW CURRENT")
		}

	case kwCATALOGS:
		p.advance()
		if p.match(kwLIKE) {
			if err := p.parseStringLit(); err != nil {
				return nil, err
			}
		} else if p.isStringLit() {
			p.advance()
		}

	case kwGRANT, kwROLE, kwPRINCIPALS, kwROLES, kwCOMPACTIONS, kwTRANSACTIONS, kwINDEXES, kwLOCKS:
		// Wait: SHOW <ident> FUNCTIONS could start with one of these keywords if followed by FUNCTIONS!
		if p.peek(1).Kind == kwFUNCTIONS {
			p.advance()
			return p.parseShowFunctionsTail(start)
		}
		p.restore(cp)
		return p.parseUnsupportedHiveNativeCommand()

	case kwFUNCTIONS:
		return p.parseShowFunctionsTail(start)

	default:
		// SHOW identifier FUNCTIONS ...
		if p.cur().Kind.isIdentifier() && p.peek(1).Kind == kwFUNCTIONS {
			p.advance()
			return p.parseShowFunctionsTail(start)
		}
		return nil, p.errorf("unexpected token %q after SHOW", p.cur().Text)
	}

	return &ParsedStmt{Kind: StmtShow, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
}

func (p *Parser) parseShowFunctionsTail(start int) (Statement, error) {
	if _, err := p.expect(kwFUNCTIONS); err != nil {
		return nil, err
	}
	if p.match(kwFROM, kwIN) {
		if _, err := p.parseIdentifierReference(); err != nil {
			return nil, err
		}
	}
	if p.match(kwLIKE) {
		if p.isStringLit() {
			p.advance()
		} else if _, err := p.parseMultipartIdentifier(); err != nil {
			return nil, err
		}
	} else if !p.isStatementEnd() {
		if p.isStringLit() {
			p.advance()
		} else if _, err := p.parseMultipartIdentifier(); err != nil {
			return nil, err
		}
	}
	return &ParsedStmt{Kind: StmtShow, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
}

func (p *Parser) parseDescribe() (Statement, error) {
	start := p.cur().Start
	if !p.match(kwDESC, kwDESCRIBE) {
		return nil, p.errorf("expected DESC or DESCRIBE")
	}
	cp := p.save()

	// 1. DESCRIBE FUNCTION EXTENDED? describeFuncName
	if p.cur().Kind == kwFUNCTION {
		p.advance()
		if p.cur().Kind == kwEXTENDED && !p.isStatementEndAt(1) {
			p.advance()
		}
		if err := p.parseDescribeFuncName(); err == nil && p.isStatementEnd() {
			return &ParsedStmt{Kind: StmtDescribe, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
		}
		p.restore(cp)
	}

	// 2. DESCRIBE namespace EXTENDED? identifierReference
	if p.isNamespace() && !p.isStatementEndAt(1) {
		p.advance()
		if p.cur().Kind == kwEXTENDED && !p.isStatementEndAt(1) {
			p.advance()
		}
		if target, err := p.parseIdentifierReference(); err == nil && p.isStatementEnd() {
			return &ParsedStmt{Kind: StmtDescribe, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
		}
		p.restore(cp)
	}

	// 3. DESCRIBE QUERY? query
	if p.cur().Kind == kwQUERY {
		p.advance()
		if err := p.parseQuery(); err == nil && p.isStatementEnd() {
			return &ParsedStmt{Kind: StmtDescribe, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
		}
		p.restore(cp)
	}
	if p.looksLikeQueryStart() && p.cur().Kind != kwTABLE {
		if err := p.parseQuery(); err == nil && p.isStatementEnd() {
			return &ParsedStmt{Kind: StmtDescribe, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
		}
		p.restore(cp)
	}

	// 4. DESCRIBE TABLE? (EXTENDED | FORMATTED)? identifierReference partitionSpec? describeColName? (AS JSON)?
	hasTable := p.cur().Kind == kwTABLE
	if hasTable && !p.isStatementEndAt(1) && p.peek(1).Kind != tokDot &&
		!(p.peek(1).Kind == kwAS && p.peek(2).Kind == kwJSON && p.isStatementEndAt(3)) &&
		!(p.peek(1).Kind == kwPARTITION && p.peek(2).Kind == tokLParen) {
		p.advance()
		if target, ok := p.tryParseDescribeRelationBody(); ok {
			return &ParsedStmt{Kind: StmtDescribe, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
		}
		p.restore(cp)
	}
	if target, ok := p.tryParseDescribeRelationBody(); ok {
		return &ParsedStmt{Kind: StmtDescribe, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
	}
	p.restore(cp)

	// Fallback: DESCRIBE TABLE identRef ... as query (`TABLE identRef` with set/pipe/organization clauses).
	if hasTable {
		if err := p.parseQuery(); err == nil && p.isStatementEnd() {
			return &ParsedStmt{Kind: StmtDescribe, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
		}
		p.restore(cp)
	}

	return nil, p.errorf("invalid DESCRIBE statement")
}

func (p *Parser) tryParseDescribeRelationBody() (string, bool) {
	if (p.cur().Kind == kwEXTENDED || p.cur().Kind == kwFORMATTED) && !p.isStatementEndAt(1) {
		cpOpt := p.save()
		p.advance()
		if target, err := p.parseDescribeRelationTail(); err == nil && p.isStatementEnd() {
			return target, true
		}
		p.restore(cpOpt)
	}
	if target, err := p.parseDescribeRelationTail(); err == nil && p.isStatementEnd() {
		return target, true
	}
	return "", false
}

func (p *Parser) parseDescribeRelationTail() (string, error) {
	target, err := p.parseIdentifierReference()
	if err != nil {
		return "", err
	}
	if p.cur().Kind == kwPARTITION {
		cpPart := p.save()
		if err := p.parsePartitionSpec(); err != nil {
			p.restore(cpPart)
		}
	}
	if p.cur().Kind == kwAS && p.peek(1).Kind == kwJSON && p.isStatementEndAt(2) {
		p.advance()
		p.advance()
		return target, nil
	}
	if !p.isStatementEnd() {
		if _, err := p.parseMultipartIdentifier(); err != nil {
			return "", err
		}
	}
	if p.cur().Kind == kwAS && p.peek(1).Kind == kwJSON {
		p.advance()
		p.advance()
	}
	return target, nil
}

func (p *Parser) parseDescribeFuncName() error {
	if p.isStringLit() {
		p.advance()
		return nil
	}
	if p.isComparisonOp() || p.isDescribeOp() {
		p.advance()
		return nil
	}
	_, err := p.parseIdentifierReference()
	return err
}

func (p *Parser) isDescribeOp() bool {
	switch p.cur().Kind {
	case tokPlus, tokMinus, tokAsterisk, tokSlash, tokPercent, kwDIV,
		tokTilde, tokAmpersand, tokPipe, tokConcatPipe, tokHat,
		kwOR, kwAND, kwIN, kwNOT,
		tokShiftLeft, tokShiftRight, tokShiftRightU, tokBang:
		return true
	}
	return false
}

func (p *Parser) parseCommentOn() (Statement, error) {
	start := p.cur().Start
	if _, err := p.expect(kwCOMMENT); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	if !p.matchNamespace() && !p.match(kwTABLE) {
		return nil, p.errorf("expected TABLE, DATABASE, SCHEMA, or NAMESPACE after COMMENT ON")
	}
	target, err := p.parseIdentifierReference()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(kwIS); err != nil {
		return nil, err
	}
	if !p.match(kwNULL) {
		if err := p.parseStringLit(); err != nil {
			return nil, err
		}
	}
	return &ParsedStmt{Kind: StmtComment, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
}

func (p *Parser) parseRefresh() (Statement, error) {
	start := p.cur().Start
	if _, err := p.expect(kwREFRESH); err != nil {
		return nil, err
	}
	if p.match(kwTABLE, kwFUNCTION) {
		target, err := p.parseIdentifierReference()
		if err != nil {
			return nil, err
		}
		return &ParsedStmt{Kind: StmtRefresh, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
	}
	if p.isStringLit() {
		p.advance()
		return &ParsedStmt{Kind: StmtRefresh, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
	}
	// REFRESH <unquoted resource path> (or bare REFRESH)
	for !p.isStatementEnd() {
		p.advance()
	}
	return &ParsedStmt{Kind: StmtRefresh, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
}

func (p *Parser) parseCacheTable() (Statement, error) {
	start := p.cur().Start
	if _, err := p.expect(kwCACHE); err != nil {
		return nil, err
	}
	p.match(kwLAZY)
	if _, err := p.expect(kwTABLE); err != nil {
		return nil, err
	}
	target, err := p.parseIdentifierReference()
	if err != nil {
		return nil, err
	}
	if p.match(kwOPTIONS) {
		if err := p.parsePropertyList(); err != nil {
			return nil, err
		}
	}
	if p.match(kwAS) {
		if err := p.parseQuery(); err != nil {
			return nil, err
		}
	} else if !p.isStatementEnd() {
		if err := p.parseQuery(); err != nil {
			return nil, err
		}
	}
	return &ParsedStmt{Kind: StmtCache, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
}

func (p *Parser) parseUncacheTable() (Statement, error) {
	start := p.cur().Start
	if _, err := p.expect(kwUNCACHE); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwTABLE); err != nil {
		return nil, err
	}
	p.matchIfExists()
	target, err := p.parseIdentifierReference()
	if err != nil {
		return nil, err
	}
	return &ParsedStmt{Kind: StmtUncache, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
}

func (p *Parser) parseClearCache() (Statement, error) {
	start := p.cur().Start
	if _, err := p.expect(kwCLEAR); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwCACHE); err != nil {
		return nil, err
	}
	return &ParsedStmt{Kind: StmtClearCache, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
}

func (p *Parser) parseRepairTable() (Statement, error) {
	start := p.cur().Start
	p.match(kwMSCK)
	if _, err := p.expect(kwREPAIR); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwTABLE); err != nil {
		return nil, err
	}
	target, err := p.parseIdentifierReference()
	if err != nil {
		return nil, err
	}
	if p.match(kwADD, kwDROP, kwSYNC) {
		if _, err := p.expect(kwPARTITIONS); err != nil {
			return nil, err
		}
	}
	return &ParsedStmt{Kind: StmtRepair, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
}

func (p *Parser) parseManageResource() (Statement, error) {
	start := p.cur().Start
	if !p.match(kwADD, kwLIST) {
		return nil, p.errorf("expected ADD or LIST")
	}
	if _, err := p.parseIdentifier(); err != nil {
		return nil, err
	}
	for !p.isStatementEnd() {
		p.advance()
	}
	return &ParsedStmt{Kind: StmtManageResource, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
}

func (p *Parser) parseCall() (Statement, error) {
	start := p.cur().Start
	if _, err := p.expect(kwCALL); err != nil {
		return nil, err
	}
	var target string
	var err error
	if p.cur().Kind == kwIDENTIFIER && p.peek(1).Kind == tokLParen {
		closeOff := p.matchingCloseOffset(1)
		if closeOff > 0 && p.peek(closeOff+1).Kind != tokLParen {
			target, err = p.parseMultipartIdentifier()
		} else {
			target, err = p.parseIdentifierReference()
		}
	} else {
		target, err = p.parseIdentifierReference()
	}
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tokLParen); err != nil {
		return nil, err
	}
	if !p.match(tokRParen) {
		if err := p.parseFunctionArgument(); err != nil {
			return nil, err
		}
		for p.match(tokComma) {
			if err := p.parseFunctionArgument(); err != nil {
				return nil, err
			}
		}
		if _, err := p.expect(tokRParen); err != nil {
			return nil, err
		}
	}
	return &ParsedStmt{Kind: StmtCall, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
}

func (p *Parser) parseExecuteImmediate() (Statement, error) {
	start := p.cur().Start
	if _, err := p.expect(kwEXECUTE); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwIMMEDIATE); err != nil {
		return nil, err
	}
	if p.isStringLit() {
		p.advance()
	} else if _, err := p.parseMultipartIdentifier(); err != nil {
		return nil, err
	}
	if p.match(kwINTO) {
		if err := p.parseMultipartIdentifierList(); err != nil {
			return nil, err
		}
	}
	if p.match(kwUSING) {
		cp := p.save()
		if p.match(tokLParen) {
			if err := p.parseNamedExpressionSeq(); err == nil && p.match(tokRParen) && p.isStatementEnd() {
				return &ParsedStmt{Kind: StmtExecuteImmediate, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
			}
			p.restore(cp)
		}
		if err := p.parseNamedExpressionSeq(); err != nil {
			return nil, err
		}
	}
	return &ParsedStmt{Kind: StmtExecuteImmediate, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
}

// parseUnsupportedHiveNativeCommand matches Spark's unsupportedHiveNativeCommands .*?
// prefix productions exactly so valid grammar rows agree and malformed prefixes reject.
func (p *Parser) parseUnsupportedHiveNativeCommand() (Statement, error) {
	start := p.cur().Start
	ok := false
	consumeAnyTail := false
	switch p.cur().Kind {
	case kwCREATE:
		p.advance()
		if p.match(kwROLE, kwINDEX) {
			ok = true
		} else if p.match(kwTEMPORARY) && p.match(kwMACRO) {
			ok = true
		}
	case kwDROP:
		p.advance()
		if p.match(kwROLE, kwINDEX) {
			ok = true
		} else if p.match(kwTEMPORARY) && p.match(kwMACRO) {
			ok = true
		}
	case kwGRANT:
		p.advance()
		hasRole := p.match(kwROLE)
		if !hasRole && p.cur().Kind == kwROLES && p.isStatementEndAt(1) {
			return nil, p.errorf("unexpected token %q after GRANT", p.cur().Text)
		}
		ok = true
		consumeAnyTail = true
	case kwREVOKE:
		p.advance()
		p.match(kwROLE)
		ok = true
		consumeAnyTail = true
	case kwSHOW:
		p.advance()
		switch p.cur().Kind {
		case kwGRANT, kwPRINCIPALS, kwROLES, kwCOMPACTIONS, kwTRANSACTIONS, kwINDEXES, kwLOCKS:
			p.advance()
			ok = true
		case kwROLE:
			p.advance()
			p.match(kwGRANT)
			ok = true
		case kwCURRENT:
			p.advance()
			ok = p.match(kwROLES)
		case kwCREATE:
			p.advance()
			ok = p.match(kwTABLE)
		}
	case kwEXPORT, kwIMPORT:
		p.advance()
		ok = p.match(kwTABLE)
	case kwLOCK, kwUNLOCK:
		p.advance()
		ok = p.match(kwTABLE, kwDATABASE)
	case kwALTER:
		p.advance()
		if p.match(kwINDEX) {
			ok = true
		} else if p.match(kwTABLE) {
			if _, err := p.parseTableIdentifier(); err == nil {
				if p.cur().Kind == kwPARTITION {
					if err := p.parsePartitionSpec(); err != nil {
						return nil, err
					}
					switch {
					case p.match(kwCOMPACT, kwCONCATENATE):
						ok = true
					case p.match(kwSET) && p.match(kwFILEFORMAT):
						ok = true
					case p.match(kwREPLACE) && p.match(kwCOLUMNS):
						ok = true
					}
				} else {
					switch {
					case p.match(kwNOT):
						if p.match(kwCLUSTERED, kwSORTED, kwSKEWED) {
							ok = true
						} else if p.match(kwSTORED) && p.match(kwAS) && p.match(kwDIRECTORIES) {
							ok = true
						}
					case p.match(kwCLUSTERED) && p.match(kwBY):
						ok = true
					case p.match(kwSKEWED) && p.match(kwBY):
						ok = true
					case p.match(kwSET):
						if p.match(kwSKEWED) && p.match(kwLOCATION) {
							ok = true
						} else if p.match(kwFILEFORMAT) {
							ok = true
						}
					case p.match(kwEXCHANGE, kwARCHIVE, kwUNARCHIVE) && p.match(kwPARTITION):
						ok = true
					case p.match(kwTOUCH, kwCOMPACT, kwCONCATENATE):
						ok = true
					case p.match(kwREPLACE) && p.match(kwCOLUMNS):
						ok = true
					}
				}
			}
		}
	case kwSTART:
		p.advance()
		ok = p.match(kwTRANSACTION)
	case kwCOMMIT, kwROLLBACK, kwDFS:
		p.advance()
		ok = true
	}
	if !ok {
		return nil, p.errorf("unexpected Hive command")
	}
	if consumeAnyTail {
		for !p.isStatementEnd() {
			p.advance()
		}
	} else {
		for p.cur().Kind == tokIdent {
			p.advance()
		}
	}
	return &ParsedStmt{Kind: StmtHiveNative, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
}
