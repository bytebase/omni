package parser

import "strings"

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
	afterSetPos := p.prevEnd()
	if p.isStatementEnd() {
		return &ParsedStmt{Kind: StmtSet, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
	}

	// SET ROLE ... is rejected by Spark v4.0.5 (INVALID_STATEMENT_OR_CLAUSE).
	if p.cur().Kind == kwROLE {
		return nil, p.errorf("SET ROLE is not valid")
	}

	// SET TIME ZONE (interval | stringLit | LOCAL)
	if p.cur().Kind == kwTIME && p.peek(1).Kind == kwZONE {
		p.advance()
		p.advance()
		if p.isStatementEnd() {
			return nil, p.errorf("expected timezone or interval after SET TIME ZONE")
		}
		if p.cur().Kind == kwINTERVAL {
			secs, err := p.parseIntervalWithSeconds()
			if err != nil {
				return nil, err
			}
			if secs < -18*3600 || secs > 18*3600 {
				return nil, p.errorf("timezone interval offset out of range [-18:00, +18:00]")
			}
			return &ParsedStmt{Kind: StmtSet, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
		}
		if p.isStringLit() || p.cur().Kind == kwLOCAL {
			p.advance()
			return &ParsedStmt{Kind: StmtSet, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
		}
		return nil, p.errorf("invalid timezone in SET TIME ZONE")
	}

	// SET (VARIABLE | VAR) (assignmentList | (multipartIdentifierList) = (query))
	// Only `SET VARIABLE` / `SET VAR` (bare) or `SET VARIABLE = ...` / `SET VAR = ...`
	// fall through to general SET configuration; any other `SET VARIABLE ...` or
	// `SET VAR ...` commits to the SQL variable assignment form.
	if (p.cur().Kind == kwVARIABLE || p.cur().Kind == kwVAR) &&
		p.peek(1).Kind != tokEq && !p.isStatementEndAt(1) {
		p.advance()
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

	// General Spark SET configuration:
	//   SET                 (already handled)
	//   SET -v
	//   SET key             (single quoted key or unquoted text without whitespace)
	//   SET key = value     (non-empty key without unquoted whitespace; value may be arbitrary)
	for !p.isStatementEnd() {
		p.advance()
	}
	endPos := p.prevEnd()
	raw := strings.TrimSpace(p.src[afterSetPos:endPos])
	if raw == "" || raw == "-v" {
		return &ParsedStmt{Kind: StmtSet, Loc: Loc{Start: start, End: endPos}}, nil
	}
	if eqIdx := strings.IndexByte(raw, '='); eqIdx >= 0 {
		key := strings.TrimSpace(raw[:eqIdx])
		if !isValidSparkSetKey(key) {
			return nil, p.errorf("invalid configuration key %q in SET statement", key)
		}
		return &ParsedStmt{Kind: StmtSet, Loc: Loc{Start: start, End: endPos}}, nil
	}
	if !isValidSparkSetKey(raw) {
		return nil, p.errorf("invalid SET syntax %q; expected 'SET', 'SET key', or 'SET key=value'", raw)
	}
	return &ParsedStmt{Kind: StmtSet, Loc: Loc{Start: start, End: endPos}}, nil
}

func isValidSparkSetKey(key string) bool {
	if key == "" {
		return false
	}
	if key[0] == '`' {
		if len(key) < 2 || key[len(key)-1] != '`' {
			return false
		}
		for i := 1; i < len(key)-1; i++ {
			if key[i] == '`' {
				if i+1 < len(key)-1 && key[i+1] == '`' {
					i++
					continue
				}
				return false
			}
		}
		return true
	}
	for i := 0; i < len(key); i++ {
		c := key[i]
		ok := (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') ||
			c == '_' || c == '.' || c == ':'
		if !ok {
			return false
		}
	}
	return true
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
		if err := p.parsePartitionSpec(false); err != nil {
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
	if p.cur().Kind == kwLOGICAL {
		return nil, p.errorf("EXPLAIN LOGICAL is not valid")
	}
	p.match(kwFORMATTED, kwEXTENDED, kwCODEGEN, kwCOST)
	if _, err := p.parseStatementOrSetReset(); err != nil {
		return nil, err
	}
	return &ParsedStmt{Kind: StmtExplain, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
}

func (p *Parser) parseShow() (Statement, error) {
	start := p.cur().Start
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
			if err := p.parsePartitionSpec(true); err != nil {
				return nil, err
			}
		}

	case kwTBLPROPERTIES:
		p.advance()
		if _, err := p.parseIdentifierReference(); err != nil {
			return nil, err
		}
		if p.match(tokLParen) {
			if _, err := p.parsePropertyKey(); err != nil {
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
			if err := p.parsePartitionSpec(true); err != nil {
				return nil, err
			}
		}

	case kwCREATE:
		p.advance()
		if _, err := p.expect(kwTABLE); err != nil {
			return nil, err
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

	case kwFUNCTIONS:
		return p.parseShowFunctionsTail(start)

	default:
		// SHOW (USER | SYSTEM | ALL) FUNCTIONS ...
		if p.cur().Kind.isIdentifier() && p.peek(1).Kind == kwFUNCTIONS {
			scope := p.advance().Text
			if !isValidShowFunctionsScope(scope) {
				return nil, p.errorf("SHOW %s FUNCTIONS is not supported; expected USER, SYSTEM, or ALL", scope)
			}
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
	if p.cur().Kind == kwLIKE {
		if p.isStatementEndAt(1) {
			p.advance()
		} else {
			p.advance()
			if p.isStringLit() {
				p.advance()
			} else if _, err := p.parseMultipartIdentifier(); err != nil {
				return nil, err
			}
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
	if p.looksLikeQueryStart() && p.cur().Kind != kwTABLE &&
		!(p.peek(1).Kind == kwPARTITION && p.peek(2).Kind == tokLParen) {
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
		if target, err := p.parseDescribeRelationTail(true); err == nil && p.isStatementEnd() {
			return target, true
		}
		p.restore(cpOpt)
	}
	if target, err := p.parseDescribeRelationTail(false); err == nil && p.isStatementEnd() {
		return target, true
	}
	return "", false
}

func (p *Parser) parseDescribeRelationTail(isExtended bool) (string, error) {
	target, err := p.parseIdentifierReference()
	if err != nil {
		return "", err
	}
	hasPartition := false
	if p.cur().Kind == kwPARTITION {
		if p.peek(1).Kind == tokLParen {
			if err := p.parsePartitionSpec(true); err != nil {
				return "", err
			}
			hasPartition = true
		}
	}
	if p.cur().Kind == kwAS && p.peek(1).Kind == kwJSON && p.isStatementEndAt(2) {
		if !isExtended {
			return "", p.errorf("DESCRIBE TABLE ... AS JSON requires EXTENDED or FORMATTED")
		}
		p.advance()
		p.advance()
		return target, nil
	}
	hasCol := false
	if !p.isStatementEnd() {
		if hasPartition {
			return "", p.errorf("DESCRIBE TABLE column for a specific partition is not supported")
		}
		if _, err := p.parseMultipartIdentifier(); err != nil {
			return "", err
		}
		hasCol = true
	}
	if p.cur().Kind == kwAS && p.peek(1).Kind == kwJSON {
		if hasCol {
			return "", p.errorf("DESCRIBE TABLE column AS JSON is not supported")
		}
		if !isExtended {
			return "", p.errorf("DESCRIBE TABLE ... AS JSON requires EXTENDED or FORMATTED")
		}
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
	if (p.cur().Kind == kwTABLE || p.cur().Kind == kwFUNCTION) && !p.isStatementEndAt(1) {
		p.advance()
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
	// Spark v4.0.5 `REFRESH .*?` accepts trailing top-level `;` (e.g. `REFRESH;`,
	// `REFRESH;;`) by capturing `;` as the resource path, while rejecting bare
	// `REFRESH` at EOF and unquoted paths containing whitespace (including inside
	// BEGIN ... END compound bodies).
	if !p.inCompound && p.cur().Kind == tokSemicolon {
		p.advance()
		return &ParsedStmt{Kind: StmtRefresh, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
	}
	if p.isStatementEnd() {
		return nil, p.errorf("resource path cannot be empty in REFRESH statement")
	}
	if p.inCompound {
		return nil, p.errorf("unquoted REFRESH resource path is not allowed inside compound block")
	}
	pathStart := p.cur().Start
	for !p.isStatementEnd() {
		p.advance()
	}
	endPos := p.prevEnd()
	rawPath := strings.TrimSpace(p.src[pathStart:endPos])
	if rawPath == "" || strings.ContainsAny(rawPath, " \t\n\r") {
		return nil, p.errorf("unquoted REFRESH resource path %q cannot contain whitespace", rawPath)
	}
	return &ParsedStmt{Kind: StmtRefresh, Loc: Loc{Start: start, End: endPos}}, nil
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
	resType, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	if !isValidManageResourceType(resType) {
		return nil, p.errorf("unsupported resource type %q; expected JAR, FILE, or ARCHIVE", resType)
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

// parseGrantRevoke parses Databricks SQL Unity Catalog / Hive Metastore
// GRANT <privileges> ON <securable> TO <principal> and
// REVOKE <privileges> ON <securable> FROM <principal> statements.
// This is an intentional Databricks SQL dialect delta beyond Spark 4.0.5
// (which rejects GRANT/REVOKE with INVALID_STATEMENT_OR_CLAUSE). Bare
// `GRANT`, `GRANT zzz`, `GRANT ROLE`, `REVOKE`, `REVOKE zzz`, `REVOKE ROLE`
// without `ON ... TO/FROM ...` are still rejected.
func (p *Parser) parseGrantRevoke() (Statement, error) {
	start := p.cur().Start
	isGrant := p.match(kwGRANT)
	if !isGrant {
		if _, err := p.expect(kwREVOKE); err != nil {
			return nil, err
		}
	}
	if p.isStatementEnd() || p.cur().Kind == kwROLE {
		return nil, p.errorf("expected privilege list after GRANT/REVOKE")
	}
	// Require privilege tokens before ON, securable tokens after ON, and TO (GRANT) or FROM (REVOKE) + principal.
	sawPriv := false
	for !p.isStatementEnd() && p.cur().Kind != kwON {
		p.advance()
		sawPriv = true
	}
	if !sawPriv || !p.match(kwON) {
		return nil, p.errorf("expected ON <securable> in GRANT/REVOKE")
	}
	dirKw, dirName := kwTO, "TO"
	if !isGrant {
		dirKw, dirName = kwFROM, "FROM"
	}
	sawObj := false
	for !p.isStatementEnd() && p.cur().Kind != dirKw {
		p.advance()
		sawObj = true
	}
	if !sawObj || !p.match(dirKw) {
		return nil, p.errorf("expected %s <principal> in GRANT/REVOKE", dirName)
	}
	if p.isStatementEnd() {
		return nil, p.errorf("expected principal in GRANT/REVOKE")
	}
	// Parse principal: optional USER/GROUP/ROLE prefix followed by a single
	// identifier or string literal, leaving any trailing tokens for strict
	// top-level validation.
	if (p.cur().Kind == kwUSER || p.cur().Kind == kwGROUP || p.cur().Kind == kwROLE) && !p.isStatementEndAt(1) {
		p.advance()
	}
	if p.isStringLit() {
		p.advance()
	} else if _, err := p.parseIdentifier(); err != nil {
		return nil, err
	}
	return &ParsedStmt{Kind: StmtHiveNative, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
}
