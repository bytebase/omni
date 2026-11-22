package parser

// parseQuery parses:
//
//	query : ctes? queryTerm queryOrganization ;
func (p *Parser) parseQuery() error {
	if err := p.enter(); err != nil {
		return err
	}
	defer p.leave()

	if p.cur().Kind == kwWITH {
		if err := p.parseCTEs(); err != nil {
			return err
		}
	}
	if err := p.parseQueryTerm(); err != nil {
		return err
	}
	return p.parseQueryOrganization()
}

func (p *Parser) looksLikeQueryStart() bool {
	return p.startsWithQueryKeywordAfterParens()
}

// startsWithQueryKeywordAfterParens scans iteratively past any leading '('
// tokens and reports whether a query-introducing keyword follows. Scanning
// iteratively avoids exponential backtracking on deeply nested parenthesized
// expressions or parenthesized table relations.
func (p *Parser) startsWithQueryKeywordAfterParens() bool {
	i := 0
	for p.peek(i).Kind == tokLParen {
		i++
	}
	switch p.peek(i).Kind {
	case kwWITH, kwSELECT, kwFROM, kwTABLE, kwVALUES, kwMAP, kwREDUCE:
		return true
	}
	return false
}

func (p *Parser) parseCTEs() error {
	if _, err := p.expect(kwWITH); err != nil {
		return err
	}
	p.match(kwRECURSIVE)
	if err := p.parseNamedQuery(); err != nil {
		return err
	}
	for p.match(tokComma) {
		if err := p.parseNamedQuery(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parseNamedQuery() error {
	if _, err := p.parseErrorCapturingIdentifier(); err != nil {
		return err
	}
	if p.cur().Kind == tokLParen {
		// Could be optional columnAliases=identifierList followed by AS? (query), OR directly (query)!
		cp := p.save()
		if err := p.parseIdentifierList(); err == nil {
			if p.cur().Kind == kwAS || p.cur().Kind == tokLParen {
				// Yes, it was columnAliases!
				goto afterAliases
			}
		}
		p.restore(cp)
	}
afterAliases:
	p.match(kwAS)
	if _, err := p.expect(tokLParen); err != nil {
		return err
	}
	if err := p.parseQuery(); err != nil {
		return err
	}
	_, err := p.expect(tokRParen)
	return err
}

// parseQueryTerm parses queryPrimary, set operations (UNION/INTERSECT/EXCEPT/MINUS),
// and the Spark 4.0 SQL pipe operator (|> operatorPipeRightSide).
func (p *Parser) parseQueryTerm() error {
	if err := p.parseQueryPrimary(); err != nil {
		return err
	}
	for {
		switch p.cur().Kind {
		case kwUNION, kwINTERSECT, kwEXCEPT, kwMINUS:
			p.advance()
			p.match(kwDISTINCT, kwALL)
			if err := p.parseQueryPrimary(); err != nil {
				return err
			}
		case tokOperatorPipe:
			p.advance()
			if err := p.parseOperatorPipeRightSide(); err != nil {
				return err
			}
		default:
			return nil
		}
	}
}

func (p *Parser) parseOperatorPipeRightSide() error {
	switch p.cur().Kind {
	case kwSELECT:
		if err := p.parseSelectClause(); err != nil {
			return err
		}
		if p.cur().Kind == kwWINDOW {
			return p.parseWindowClause()
		}
		return nil
	case kwEXTEND:
		p.advance()
		return p.parseNamedExpressionSeq()
	case kwSET:
		p.advance()
		for {
			// SqlBaseParser.g4:1579 notes dotted names in `|> SET` are invalid syntax
			// captured only for diagnostics; strict Parse accepts simple identifiers only.
			if _, err := p.parseErrorCapturingIdentifier(); err != nil {
				return err
			}
			if _, err := p.expect(tokEq); err != nil {
				return err
			}
			if err := p.parseExpression(); err != nil {
				return err
			}
			if !p.match(tokComma) {
				break
			}
		}
		return nil
	case kwDROP:
		p.advance()
		return p.parseIdentifierSeq()
	case kwAS:
		p.advance()
		_, err := p.parseErrorCapturingIdentifier()
		return err
	case kwWHERE:
		// SqlBaseParser.g4:1561 notes WINDOW is not allowed in `|> WHERE`
		// (only captured in the grammar for diagnostics); strict Parse rejects it.
		if err := p.parseWhereClause(); err != nil {
			return err
		}
		if p.cur().Kind == kwWINDOW {
			return p.errorf("WINDOW clause is not allowed in WHERE pipe operator")
		}
		return nil
	case kwPIVOT:
		// SqlBaseParser.g4:1565 notes combining PIVOT and UNPIVOT in one pipe operator
		// is invalid syntax captured only for diagnostics; strict Parse rejects it.
		if err := p.parsePivotClause(); err != nil {
			return err
		}
		if p.cur().Kind == kwUNPIVOT {
			return p.errorf("UNPIVOT cannot follow PIVOT in the same pipe operator")
		}
		return nil
	case kwUNPIVOT:
		if err := p.parseUnpivotClause(); err != nil {
			return err
		}
		if p.cur().Kind == kwPIVOT {
			return p.errorf("PIVOT cannot follow UNPIVOT in the same pipe operator")
		}
		return nil
	case kwTABLESAMPLE:
		return p.parseSample()
	case kwUNION, kwEXCEPT, kwMINUS, kwINTERSECT:
		p.advance()
		p.match(kwDISTINCT, kwALL)
		return p.parseQueryPrimary()
	case kwAGGREGATE:
		p.advance()
		// Try optional namedExpressionSeq followed by optional aggregationClause.
		// If `GROUP BY ...` is immediately present, parse it directly as aggregationClause;
		// otherwise try namedExpressionSeq first (so `AGGREGATE GROUP a` / `AGGREGATE GROUP BY` work).
		if p.cur().Kind == kwGROUP && p.peek(1).Kind == kwBY {
			cpAgg := p.save()
			if err := p.parseAggregationClause(); err == nil && (p.isStatementEnd() || p.cur().Kind == tokRParen || p.cur().Kind == tokOperatorPipe || p.isQueryOrgStart()) {
				return nil
			}
			p.restore(cpAgg)
		}
		if p.cur().Kind != tokEOF && p.cur().Kind != tokSemicolon &&
			p.cur().Kind != tokRParen && p.cur().Kind != tokOperatorPipe && !p.isQueryOrgStart() {
			if err := p.parseNamedExpressionSeq(); err != nil {
				return err
			}
		}
		if p.cur().Kind == kwGROUP {
			return p.parseAggregationClause()
		}
		return nil
	case kwJOIN, kwINNER, kwCROSS, kwLEFT, kwRIGHT, kwFULL, kwSEMI, kwANTI, kwNATURAL:
		return p.parseJoinRelation()
	default:
		// queryOrganization (may be empty if followed by EOF/;/etc., or non-empty)
		return p.parseQueryOrganization()
	}
}

func (p *Parser) isQueryOrgStart() bool {
	switch p.cur().Kind {
	case kwORDER, kwCLUSTER, kwDISTRIBUTE, kwSORT, kwWINDOW, kwLIMIT, kwOFFSET:
		return true
	}
	return false
}

func (p *Parser) parseQueryPrimary() error {
	switch p.cur().Kind {
	case kwSELECT, kwMAP, kwREDUCE:
		return p.parseQuerySpecification()
	case kwFROM:
		return p.parseFromStatement()
	case kwTABLE:
		p.advance()
		_, err := p.parseIdentifierReference()
		return err
	case kwVALUES:
		return p.parseInlineTable()
	case tokLParen:
		p.advance()
		if err := p.parseQuery(); err != nil {
			return err
		}
		_, err := p.expect(tokRParen)
		return err
	default:
		return p.errorf("expected query, got %q", p.cur().Text)
	}
}

func (p *Parser) parseFromStatement() error {
	if err := p.parseFromClause(); err != nil {
		return err
	}
	for p.cur().Kind == kwSELECT || p.cur().Kind == kwMAP || p.cur().Kind == kwREDUCE {
		if err := p.parseFromStatementBody(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parseFromStatementBody() error {
	if p.isTransformStart() {
		if err := p.parseTransformClause(); err != nil {
			return err
		}
		if p.cur().Kind == kwWHERE {
			if err := p.parseWhereClause(); err != nil {
				return err
			}
		}
		return p.parseQueryOrganization()
	}
	if err := p.parseSelectClause(); err != nil {
		return err
	}
	for p.cur().Kind == kwLATERAL && p.peek(1).Kind == kwVIEW {
		if err := p.parseLateralView(); err != nil {
			return err
		}
	}
	if p.cur().Kind == kwWHERE {
		if err := p.parseWhereClause(); err != nil {
			return err
		}
	}
	if p.cur().Kind == kwGROUP {
		if err := p.parseAggregationClause(); err != nil {
			return err
		}
	}
	if p.cur().Kind == kwHAVING {
		if err := p.parseHavingClause(); err != nil {
			return err
		}
	}
	if p.cur().Kind == kwWINDOW {
		if err := p.parseWindowClause(); err != nil {
			return err
		}
	}
	return p.parseQueryOrganization()
}

func (p *Parser) parseQuerySpecification() error {
	if p.isTransformStart() {
		if err := p.parseTransformClause(); err != nil {
			return err
		}
	} else {
		if err := p.parseSelectClause(); err != nil {
			return err
		}
	}
	if p.cur().Kind == kwFROM {
		if err := p.parseFromClause(); err != nil {
			return err
		}
	}
	for p.cur().Kind == kwLATERAL && p.peek(1).Kind == kwVIEW {
		if err := p.parseLateralView(); err != nil {
			return err
		}
	}
	if p.cur().Kind == kwWHERE {
		if err := p.parseWhereClause(); err != nil {
			return err
		}
	}
	if p.cur().Kind == kwGROUP {
		if err := p.parseAggregationClause(); err != nil {
			return err
		}
	}
	if p.cur().Kind == kwHAVING {
		if err := p.parseHavingClause(); err != nil {
			return err
		}
	}
	if p.cur().Kind == kwWINDOW {
		if err := p.parseWindowClause(); err != nil {
			return err
		}
	}
	// Databricks SQL QUALIFY extension: allow `QUALIFY <booleanExpression>` when
	// QUALIFY is followed by an expression (outside the strict Spark grammar, only
	// accepted when an expression follows, never as a stray trailing word).
	if p.cur().Kind == tokIdent && isWord(p.cur(), "QUALIFY") && !p.isStatementEnd() && p.peek(1).Kind != tokRParen {
		cp := p.save()
		p.advance()
		if err := p.parseBooleanExpression(); err != nil {
			p.restore(cp)
		}
	}
	return nil
}

func (p *Parser) isTransformStart() bool {
	if p.cur().Kind == kwMAP || p.cur().Kind == kwREDUCE {
		return true
	}
	if p.cur().Kind != kwSELECT || p.peek(1).Kind != kwTRANSFORM || p.peek(2).Kind != tokLParen {
		return false
	}
	return p.hasTransformClauseAfterParens(2)
}

// hasTransformClauseAfterParens scans from the '(' at lparenOffset to its
// matching ')' and reports whether the following tokens start a Hive
// transformClause (`ROW FORMAT`, `RECORDWRITER`, or `USING`), distinguishing
// script transformations from higher-order `SELECT TRANSFORM(arr, x -> x)`
// function calls without reparsing the arguments.
func (p *Parser) hasTransformClauseAfterParens(lparenOffset int) bool {
	depth := 0
	for i := lparenOffset; ; i++ {
		k := p.peek(i).Kind
		if k == tokEOF {
			return false
		}
		if k == tokLParen {
			depth++
		} else if k == tokRParen {
			depth--
			if depth == 0 {
				after := p.peek(i + 1).Kind
				return after == kwUSING ||
					(after == kwROW && p.peek(i+2).Kind == kwFORMAT) ||
					after == kwRECORDWRITER
			}
		}
	}
}

func (p *Parser) parseTransformClause() error {
	if p.match(kwSELECT) {
		if _, err := p.expect(kwTRANSFORM); err != nil {
			return err
		}
		if _, err := p.expect(tokLParen); err != nil {
			return err
		}
		p.matchSetQuantifierIfFollowedByExpr()
		if err := p.parseExpressionSeq(); err != nil {
			return err
		}
		if _, err := p.expect(tokRParen); err != nil {
			return err
		}
	} else if p.match(kwMAP, kwREDUCE) {
		p.matchSetQuantifierIfFollowedByExpr()
		if err := p.parseExpressionSeq(); err != nil {
			return err
		}
	} else {
		return p.errorf("expected TRANSFORM, MAP, or REDUCE")
	}

	if p.cur().Kind == kwROW && p.peek(1).Kind == kwFORMAT {
		if err := p.parseRowFormat(); err != nil {
			return err
		}
	}
	if p.match(kwRECORDWRITER) {
		if err := p.parseStringLit(); err != nil {
			return err
		}
	}
	if _, err := p.expect(kwUSING); err != nil {
		return err
	}
	if err := p.parseStringLit(); err != nil {
		return err
	}
	if p.match(kwAS) {
		if p.match(tokLParen) {
			cp := p.save()
			if err := p.parseColTypeList(); err == nil && p.cur().Kind == tokRParen {
				p.advance()
			} else {
				p.restore(cp)
				if err := p.parseIdentifierSeq(); err != nil {
					return err
				}
				if _, err := p.expect(tokRParen); err != nil {
					return err
				}
			}
		} else {
			cp := p.save()
			if err := p.parseColTypeList(); err != nil {
				p.restore(cp)
				if err := p.parseIdentifierSeq(); err != nil {
					return err
				}
			}
		}
	}
	if p.cur().Kind == kwROW && p.peek(1).Kind == kwFORMAT {
		if err := p.parseRowFormat(); err != nil {
			return err
		}
	}
	if p.match(kwRECORDREADER) {
		if err := p.parseStringLit(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) matchSetQuantifierIfFollowedByExpr() {
	if p.cur().Kind != kwDISTINCT && p.cur().Kind != kwALL {
		return
	}
	next := p.peek(1).Kind
	switch next {
	case tokAsterisk:
		// Distinguish a standalone star (`SELECT ALL * FROM t`,
		// `SELECT DISTINCT *, a FROM t`, `SELECT ALL * EXCEPT (a) FROM t`)
		// from multiplication on a column named `all` or `distinct`
		// (`SELECT all * 2 FROM t`, `SELECT distinct * 2 FROM t`).
		afterStar := p.peek(2).Kind
		switch afterStar {
		case tokEOF, tokSemicolon, tokComma, tokRParen, tokOperatorPipe,
			kwFROM, kwWHERE, kwUSING, kwROW, kwRECORDWRITER,
			kwGROUP, kwHAVING, kwWINDOW, kwORDER, kwCLUSTER, kwDISTRIBUTE,
			kwSORT, kwLIMIT, kwOFFSET, kwUNION, kwINTERSECT, kwEXCEPT, kwMINUS:
			p.advance()
		case tokIdent:
			if isWord(p.peek(2), "QUALIFY") && p.peek(3).Kind != tokEOF && p.peek(3).Kind != tokSemicolon && p.peek(3).Kind != tokRParen {
				p.advance()
			}
		}
		return
	case tokEOF, tokSemicolon, tokComma, tokRParen,
		tokDot, tokLBracket, tokDoubleColon, kwCOLLATE,
		tokSlash, tokPercent, kwDIV,
		tokConcatPipe, tokPipe, tokAmpersand, tokHat,
		tokShiftLeft, tokShiftRight, tokShiftRightU,
		tokEq, tokNeq, tokNeqBang, tokNullSafeEq, tokLT, tokLTE, tokGT, tokGTE,
		kwAND, kwOR, kwIS, kwIN, kwLIKE, kwILIKE, kwRLIKE, kwBETWEEN,
		tokArrow, tokFatArrow, kwAS,
		kwFROM, kwWHERE, kwUSING, kwROW, kwRECORDWRITER,
		kwGROUP, kwHAVING, kwWINDOW, kwORDER, kwCLUSTER, kwDISTRIBUTE,
		kwSORT, kwLIMIT, kwOFFSET, kwUNION, kwINTERSECT, kwEXCEPT, kwMINUS:
		return
	}
	p.advance()
}

func (p *Parser) parseSelectClause() error {
	if _, err := p.expect(kwSELECT); err != nil {
		return err
	}
	for p.cur().Kind == tokHintStart {
		if err := p.parseHint(); err != nil {
			return err
		}
	}
	p.matchSetQuantifierIfFollowedByExpr()
	return p.parseNamedExpressionSeq()
}

func (p *Parser) parseHint() error {
	if _, err := p.expect(tokHintStart); err != nil {
		return err
	}
	if err := p.parseHintStatement(); err != nil {
		return err
	}
	for p.cur().Kind != tokHintEnd && p.cur().Kind != tokEOF {
		p.match(tokComma)
		if p.cur().Kind == tokHintEnd {
			return p.errorf("unexpected trailing comma in hint")
		}
		if err := p.parseHintStatement(); err != nil {
			return err
		}
	}
	_, err := p.expect(tokHintEnd)
	return err
}

func (p *Parser) parseHintStatement() error {
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	if p.match(tokLParen) {
		if err := p.parsePrimaryExpression(); err != nil {
			return err
		}
		for p.match(tokComma) {
			if err := p.parsePrimaryExpression(); err != nil {
				return err
			}
		}
		if _, err := p.expect(tokRParen); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parseFromClause() error {
	if _, err := p.expect(kwFROM); err != nil {
		return err
	}
	if err := p.parseRelation(); err != nil {
		return err
	}
	for p.match(tokComma) {
		if err := p.parseRelation(); err != nil {
			return err
		}
	}
	for p.cur().Kind == kwLATERAL && p.peek(1).Kind == kwVIEW {
		if err := p.parseLateralView(); err != nil {
			return err
		}
	}
	if p.cur().Kind == kwPIVOT {
		if err := p.parsePivotClause(); err != nil {
			return err
		}
	}
	if p.cur().Kind == kwUNPIVOT {
		if err := p.parseUnpivotClause(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parseRelation() error {
	if err := p.enter(); err != nil {
		return err
	}
	defer p.leave()

	if p.isLateralRelationModifier() {
		p.advance()
	}
	if err := p.parseRelationPrimary(); err != nil {
		return err
	}
	for {
		switch {
		case p.isJoinStart():
			if err := p.parseJoinRelation(); err != nil {
				return err
			}
		case p.cur().Kind == kwPIVOT:
			if err := p.parsePivotClause(); err != nil {
				return err
			}
		case p.cur().Kind == kwUNPIVOT:
			if err := p.parseUnpivotClause(); err != nil {
				return err
			}
		default:
			return nil
		}
	}
}

// isLateralRelationModifier reports whether a leading LATERAL keyword before a
// relationPrimary acts as the LATERAL modifier (e.g. `FROM LATERAL (SELECT ...)`,
// `FROM LATERAL explode(arr)`) rather than a table or schema named `lateral`
// (e.g. `FROM lateral`, `FROM lateral.t`, `JOIN lateral ON ...`).
func (p *Parser) isLateralRelationModifier() bool {
	if p.cur().Kind != kwLATERAL {
		return false
	}
	next := p.peek(1).Kind
	switch next {
	case tokEOF, tokSemicolon, tokComma, tokRParen, tokOperatorPipe,
		tokDot, kwAS, kwVIEW,
		kwON, kwUSING,
		kwTABLESAMPLE, kwPIVOT, kwUNPIVOT,
		kwJOIN, kwINNER, kwCROSS, kwLEFT, kwRIGHT, kwFULL, kwSEMI, kwANTI, kwNATURAL,
		kwWHERE, kwGROUP, kwHAVING, kwWINDOW, kwORDER, kwCLUSTER, kwDISTRIBUTE,
		kwSORT, kwLIMIT, kwOFFSET, kwUNION, kwINTERSECT, kwEXCEPT, kwMINUS,
		kwSELECT, kwMAP, kwREDUCE, kwSET, kwINTO:
		return false
	case kwFOR:
		return p.peek(2).Kind != kwSYSTEM_TIME && p.peek(2).Kind != kwSYSTEM_VERSION && p.peek(2).Kind != kwVERSION
	case kwTIMESTAMP, kwVERSION:
		return !(p.peek(2).Kind == kwAS && p.peek(3).Kind == kwOF)
	case kwWITH:
		return p.peek(2).Kind != tokLParen
	case kwLATERAL:
		return p.peek(2).Kind != kwVIEW
	case tokIdent:
		if isWord(p.peek(1), "QUALIFY") && p.peek(2).Kind != tokEOF && p.peek(2).Kind != tokSemicolon && p.peek(2).Kind != tokRParen {
			return false
		}
	}
	return true
}

func (p *Parser) isJoinStart() bool {
	k := p.cur().Kind
	switch k {
	case kwJOIN:
		return true
	case kwINNER, kwCROSS, kwSEMI, kwANTI:
		return p.peek(1).Kind == kwJOIN
	case kwLEFT, kwRIGHT, kwFULL:
		n := p.peek(1).Kind
		return n == kwJOIN || ((n == kwOUTER || n == kwSEMI || n == kwANTI) && p.peek(2).Kind == kwJOIN)
	case kwNATURAL:
		return true
	}
	return false
}

func (p *Parser) parseJoinRelation() error {
	natural := p.match(kwNATURAL)
	if err := p.parseJoinType(); err != nil {
		return err
	}
	if _, err := p.expect(kwJOIN); err != nil {
		return err
	}
	if p.isLateralRelationModifier() {
		p.advance()
	}
	if err := p.parseRelationPrimary(); err != nil {
		return err
	}
	if !natural {
		if p.match(kwON) {
			return p.parseBooleanExpression()
		}
		if p.match(kwUSING) {
			return p.parseIdentifierList()
		}
	}
	return nil
}

func (p *Parser) parseJoinType() error {
	switch p.cur().Kind {
	case kwINNER:
		p.advance()
	case kwCROSS:
		p.advance()
	case kwLEFT:
		p.advance()
		p.match(kwOUTER, kwSEMI, kwANTI)
	case kwRIGHT:
		p.advance()
		p.match(kwOUTER)
	case kwFULL:
		p.advance()
		p.match(kwOUTER)
	case kwSEMI, kwANTI:
		p.advance()
	}
	return nil
}

func (p *Parser) parseRelationPrimary() error {
	switch p.cur().Kind {
	case kwVALUES:
		cp := p.save()
		if err := p.parseInlineTable(); err == nil {
			return nil
		}
		p.restore(cp)

	case tokLParen:
		cp := p.save()
		p.advance()
		// Try (query) only when a query-introducing keyword follows any leading '('.
		if p.startsWithQueryKeywordAfterParens() {
			if err := p.parseQuery(); err == nil && p.cur().Kind == tokRParen {
				p.advance()
				if p.cur().Kind == kwTABLESAMPLE {
					if err := p.parseSample(); err != nil {
						return err
					}
				}
				return p.parseTableAlias()
			}
			p.restore(cp)
			p.advance()
		}
		if err := p.parseRelation(); err != nil {
			return err
		}
		if _, err := p.expect(tokRParen); err != nil {
			return err
		}
		if p.cur().Kind == kwTABLESAMPLE {
			if err := p.parseSample(); err != nil {
				return err
			}
		}
		return p.parseTableAlias()
	}

	// Could be functionTable (funcName(...) tableAlias) or tableName (identifierReference ...).
	if p.canStartFunctionCall() {
		cp := p.save()
		if err := p.parseFunctionTable(); err == nil {
			return nil
		}
		p.restore(cp)
	}

	// tableName: identifierReference temporalClause? optionsClause? sample? tableAlias
	if _, err := p.parseIdentifierReference(); err != nil {
		return err
	}
	if p.isTemporalClauseStart() {
		if err := p.parseTemporalClause(); err != nil {
			return err
		}
	}
	if p.cur().Kind == kwWITH && p.peek(1).Kind == tokLParen {
		if err := p.parseOptionsClause(); err != nil {
			return err
		}
	}
	if p.cur().Kind == kwTABLESAMPLE {
		if err := p.parseSample(); err != nil {
			return err
		}
	}
	return p.parseTableAlias()
}

func (p *Parser) isTemporalClauseStart() bool {
	switch p.cur().Kind {
	case kwFOR:
		n := p.peek(1).Kind
		return n == kwSYSTEM_VERSION || n == kwVERSION || n == kwSYSTEM_TIME || n == kwTIMESTAMP
	case kwSYSTEM_VERSION, kwVERSION, kwSYSTEM_TIME, kwTIMESTAMP:
		return p.peek(1).Kind == kwAS && p.peek(2).Kind == kwOF
	}
	return false
}

func (p *Parser) parseTemporalClause() error {
	p.match(kwFOR)
	if p.match(kwSYSTEM_VERSION, kwVERSION) {
		if _, err := p.expect(kwAS); err != nil {
			return err
		}
		if _, err := p.expect(kwOF); err != nil {
			return err
		}
		if p.match(tokInt) {
			return nil
		}
		return p.parseStringLit()
	}
	if p.match(kwSYSTEM_TIME, kwTIMESTAMP) {
		if _, err := p.expect(kwAS); err != nil {
			return err
		}
		if _, err := p.expect(kwOF); err != nil {
			return err
		}
		return p.parseValueExpression()
	}
	return p.errorf("expected temporal clause")
}

func (p *Parser) parseOptionsClause() error {
	if _, err := p.expect(kwWITH); err != nil {
		return err
	}
	return p.parsePropertyList()
}

func (p *Parser) parseInlineTable() error {
	if _, err := p.expect(kwVALUES); err != nil {
		return err
	}
	if err := p.parseExpression(); err != nil {
		return err
	}
	for p.match(tokComma) {
		if err := p.parseExpression(); err != nil {
			return err
		}
	}
	return p.parseTableAlias()
}

func (p *Parser) parseFunctionTable() error {
	// Single-pass handling for IDENTIFIER(expr)(tableArgs) vs IDENTIFIER(tableArgs)
	// without reparsing nested IDENTIFIER(...) arguments.
	if p.cur().Kind == kwIDENTIFIER && p.peek(1).Kind == tokLParen {
		p.advance()
		p.advance()
		singlePlainArg := false
		if !p.match(tokRParen) {
			isTableOrNamed := p.cur().Kind == kwTABLE ||
				(p.cur().Kind.isIdentifier() && p.peek(1).Kind == tokFatArrow)
			if err := p.parseFunctionTableArgument(); err != nil {
				return err
			}
			argCount := 1
			for p.match(tokComma) {
				if err := p.parseFunctionTableArgument(); err != nil {
					return err
				}
				argCount++
			}
			if _, err := p.expect(tokRParen); err != nil {
				return err
			}
			singlePlainArg = !isTableOrNamed && argCount == 1
		}
		if singlePlainArg && p.cur().Kind == tokLParen {
			p.advance()
			if !p.match(tokRParen) {
				for {
					if err := p.parseFunctionTableArgument(); err != nil {
						return err
					}
					if !p.match(tokComma) {
						break
					}
				}
				if _, err := p.expect(tokRParen); err != nil {
					return err
				}
			}
		}
		return p.parseTableAlias()
	}

	if _, err := p.parseQualifiedName(); err != nil {
		return err
	}
	if _, err := p.expect(tokLParen); err != nil {
		return err
	}
	if !p.match(tokRParen) {
		for {
			if err := p.parseFunctionTableArgument(); err != nil {
				return err
			}
			if !p.match(tokComma) {
				break
			}
		}
		if _, err := p.expect(tokRParen); err != nil {
			return err
		}
	}
	return p.parseTableAlias()
}

func (p *Parser) parseFunctionTableArgument() error {
	if p.cur().Kind.isIdentifier() && p.peek(1).Kind == tokFatArrow && p.peek(2).Kind == kwTABLE {
		cp := p.save()
		p.advance()
		p.advance()
		if err := p.parseFunctionTableSubqueryArgument(); err == nil {
			return nil
		}
		p.restore(cp)
	}
	if p.cur().Kind == kwTABLE {
		cp := p.save()
		if err := p.parseFunctionTableSubqueryArgument(); err == nil {
			return nil
		}
		p.restore(cp)
	}
	return p.parseFunctionArgument()
}

func (p *Parser) parseFunctionTableSubqueryArgument() error {
	if _, err := p.expect(kwTABLE); err != nil {
		return err
	}
	if p.match(tokLParen) {
		cp := p.save()
		if err := p.parseQuery(); err == nil && p.cur().Kind == tokRParen {
			p.advance()
		} else {
			p.restore(cp)
			if _, err := p.parseIdentifierReference(); err != nil {
				return err
			}
			if _, err := p.expect(tokRParen); err != nil {
				return err
			}
		}
	} else {
		if _, err := p.parseIdentifierReference(); err != nil {
			return err
		}
	}
	if p.isTableArgPartitioningStart() {
		return p.parseTableArgumentPartitioning()
	}
	return nil
}

func (p *Parser) isTableArgPartitioningStart() bool {
	if p.cur().Kind == kwWITH && p.peek(1).Kind == kwSINGLE && p.peek(2).Kind == kwPARTITION {
		return true
	}
	if (p.cur().Kind == kwPARTITION || p.cur().Kind == kwDISTRIBUTE) && p.peek(1).Kind == kwBY {
		return true
	}
	return false
}

func (p *Parser) parseTableArgumentPartitioning() error {
	if p.match(kwWITH) {
		if _, err := p.expect(kwSINGLE); err != nil {
			return err
		}
		if _, err := p.expect(kwPARTITION); err != nil {
			return err
		}
	} else {
		if !p.match(kwPARTITION, kwDISTRIBUTE) {
			return p.errorf("expected PARTITION or DISTRIBUTE")
		}
		if _, err := p.expect(kwBY); err != nil {
			return err
		}
		if p.cur().Kind == tokLParen {
			cpParen := p.save()
			p.advance()
			if err := p.parseExpressionSeq(); err == nil && p.match(tokRParen) {
				// parenthesized partition expression list
			} else {
				p.restore(cpParen)
				if err := p.parseExpression(); err != nil {
					return err
				}
			}
		} else {
			if err := p.parseExpression(); err != nil {
				return err
			}
			for {
				cp := p.save()
				if !p.match(tokComma) {
					break
				}
				if err := p.parseExpression(); err != nil {
					p.restore(cp)
					break
				}
			}
		}
	}
	if (p.cur().Kind == kwORDER || p.cur().Kind == kwSORT) && p.peek(1).Kind == kwBY {
		p.advance()
		p.advance()
		if p.cur().Kind == tokLParen {
			cpParen := p.save()
			p.advance()
			if err := p.parseSortItemList(); err == nil && p.match(tokRParen) {
				// parenthesized sort item list
			} else {
				p.restore(cpParen)
				if err := p.parseSortItem(); err != nil {
					return err
				}
			}
		} else {
			if err := p.parseSortItem(); err != nil {
				return err
			}
			for {
				cp := p.save()
				if !p.match(tokComma) {
					break
				}
				if err := p.parseSortItem(); err != nil {
					p.restore(cp)
					break
				}
			}
		}
	}
	return nil
}

// parseTableAlias parses: (AS? strictIdentifier identifierList?)?
func (p *Parser) parseTableAlias() error {
	if p.match(kwAS) {
		if _, err := p.parseStrictIdentifier(); err != nil {
			return err
		}
		if p.cur().Kind == tokLParen {
			return p.parseIdentifierList()
		}
		return nil
	}
	if !p.canBeImplicitTableAlias() {
		return nil
	}
	if _, err := p.parseStrictIdentifier(); err != nil {
		return err
	}
	if p.cur().Kind == tokLParen {
		cp := p.save()
		if err := p.parseIdentifierList(); err != nil {
			p.restore(cp)
		}
	}
	return nil
}

func (p *Parser) canBeImplicitTableAlias() bool {
	k := p.cur().Kind
	if !k.isStrictIdentifier() {
		return false
	}
	if k == tokIdent || k == tokBackquoted || k == tokDoubleQuoted {
		// Don't swallow QUALIFY when followed by an expression.
		if k == tokIdent && isWord(p.cur(), "QUALIFY") && !p.isStatementEnd() && p.peek(1).Kind != tokRParen {
			return false
		}
		return true
	}
	next := p.peek(1).Kind
	// If keyword is at end of relation (EOF, ';', ',', ')'), check if it's a valid strict alias.
	if next == tokEOF || next == tokSemicolon || next == tokComma || next == tokRParen {
		return true
	}
	switch k {
	case kwWHERE, kwGROUP, kwHAVING, kwWINDOW, kwORDER, kwCLUSTER, kwDISTRIBUTE,
		kwSORT, kwLIMIT, kwOFFSET, kwPIVOT, kwUNPIVOT, kwTABLESAMPLE, kwSET,
		kwSELECT, kwINSERT, kwMAP, kwREDUCE, kwWHEN, kwTHEN, kwELSE, kwEND,
		kwVALUES, kwFROM, kwWITH, kwFOR, kwDO, kwUNTIL, kwINTO, kwUSING, kwON,
		kwANTI, kwSEMI, kwAS, kwBY:
		return false
	}
	return true
}

func (p *Parser) parseSample() error {
	if _, err := p.expect(kwTABLESAMPLE); err != nil {
		return err
	}
	if _, err := p.expect(tokLParen); err != nil {
		return err
	}
	if !p.match(tokRParen) {
		if err := p.parseSampleMethod(); err != nil {
			return err
		}
		if _, err := p.expect(tokRParen); err != nil {
			return err
		}
	}
	if p.match(kwREPEATABLE) {
		if _, err := p.expect(tokLParen); err != nil {
			return err
		}
		if _, err := p.expect(tokInt); err != nil {
			return err
		}
		if _, err := p.expect(tokRParen); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parseSampleMethod() error {
	// BUCKET int OUT OF int (ON (ident | qualName()))?
	if p.cur().Kind == kwBUCKET && p.peek(1).Kind == tokInt {
		p.advance()
		p.advance()
		if _, err := p.expect(kwOUT); err != nil {
			return err
		}
		if _, err := p.expect(kwOF); err != nil {
			return err
		}
		if _, err := p.expect(tokInt); err != nil {
			return err
		}
		if p.match(kwON) {
			if _, err := p.parseQualifiedName(); err != nil {
				return err
			}
			if p.match(tokLParen) {
				if _, err := p.expect(tokRParen); err != nil {
					return err
				}
			}
		}
		return nil
	}
	// MINUS? (INT | DECIMAL) PERCENT
	cp := p.save()
	p.match(tokMinus)
	if p.match(tokInt, tokDecimal) && p.match(kwPERCENT) {
		return nil
	}
	p.restore(cp)

	// expression ROWS | expression
	if err := p.parseExpression(); err != nil {
		return err
	}
	p.match(kwROWS)
	return nil
}

func (p *Parser) parseLateralView() error {
	if _, err := p.expect(kwLATERAL); err != nil {
		return err
	}
	if _, err := p.expect(kwVIEW); err != nil {
		return err
	}
	p.match(kwOUTER)
	if _, err := p.parseQualifiedName(); err != nil {
		return err
	}
	if _, err := p.expect(tokLParen); err != nil {
		return err
	}
	if !p.match(tokRParen) {
		if err := p.parseExpressionSeq(); err != nil {
			return err
		}
		if _, err := p.expect(tokRParen); err != nil {
			return err
		}
	}
	if _, err := p.parseIdentifier(); err != nil {
		return err
	}
	if p.match(kwAS) {
		if _, err := p.parseIdentifier(); err != nil {
			return err
		}
		for p.match(tokComma) {
			if _, err := p.parseIdentifier(); err != nil {
				return err
			}
		}
		return nil
	}
	// Optional colName without AS (careful not to steal clause keywords).
	if p.canBeImplicitTableAlias() {
		if _, err := p.parseIdentifier(); err != nil {
			return err
		}
		for p.match(tokComma) {
			if _, err := p.parseIdentifier(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *Parser) parsePivotClause() error {
	if _, err := p.expect(kwPIVOT); err != nil {
		return err
	}
	if _, err := p.expect(tokLParen); err != nil {
		return err
	}
	if err := p.parseNamedExpressionSeq(); err != nil {
		return err
	}
	if _, err := p.expect(kwFOR); err != nil {
		return err
	}
	if p.match(tokLParen) {
		if err := p.parseIdentifierSeq(); err != nil {
			return err
		}
		if _, err := p.expect(tokRParen); err != nil {
			return err
		}
	} else {
		if _, err := p.parseErrorCapturingIdentifier(); err != nil {
			return err
		}
	}
	if _, err := p.expect(kwIN); err != nil {
		return err
	}
	if _, err := p.expect(tokLParen); err != nil {
		return err
	}
	for {
		if err := p.parseExpression(); err != nil {
			return err
		}
		if p.match(kwAS) {
			if _, err := p.parseErrorCapturingIdentifier(); err != nil {
				return err
			}
		} else if p.cur().Kind != tokComma && p.cur().Kind != tokRParen && p.cur().Kind.isIdentifier() {
			if _, err := p.parseErrorCapturingIdentifier(); err != nil {
				return err
			}
		}
		if !p.match(tokComma) {
			break
		}
	}
	if _, err := p.expect(tokRParen); err != nil {
		return err
	}
	_, err := p.expect(tokRParen)
	return err
}

func (p *Parser) parseUnpivotClause() error {
	if _, err := p.expect(kwUNPIVOT); err != nil {
		return err
	}
	if p.cur().Kind == kwINCLUDE || p.cur().Kind == kwEXCLUDE {
		p.advance()
		if _, err := p.expect(kwNULLS); err != nil {
			return err
		}
	}
	if _, err := p.expect(tokLParen); err != nil {
		return err
	}
	if p.match(tokLParen) {
		// Multi-value column clause.
		if _, err := p.parseIdentifier(); err != nil {
			return err
		}
		for p.match(tokComma) {
			if _, err := p.parseIdentifier(); err != nil {
				return err
			}
		}
		if _, err := p.expect(tokRParen); err != nil {
			return err
		}
		if _, err := p.expect(kwFOR); err != nil {
			return err
		}
		if _, err := p.parseIdentifier(); err != nil {
			return err
		}
		if _, err := p.expect(kwIN); err != nil {
			return err
		}
		if _, err := p.expect(tokLParen); err != nil {
			return err
		}
		for {
			if _, err := p.expect(tokLParen); err != nil {
				return err
			}
			if err := p.parseMultipartIdentifierList(); err != nil {
				return err
			}
			if _, err := p.expect(tokRParen); err != nil {
				return err
			}
			if err := p.parseOptionalUnpivotAlias(); err != nil {
				return err
			}
			if !p.match(tokComma) {
				break
			}
		}
		if _, err := p.expect(tokRParen); err != nil {
			return err
		}
	} else {
		// Single-value column clause.
		if _, err := p.parseIdentifier(); err != nil {
			return err
		}
		if _, err := p.expect(kwFOR); err != nil {
			return err
		}
		if _, err := p.parseIdentifier(); err != nil {
			return err
		}
		if _, err := p.expect(kwIN); err != nil {
			return err
		}
		if _, err := p.expect(tokLParen); err != nil {
			return err
		}
		for {
			if _, err := p.parseMultipartIdentifier(); err != nil {
				return err
			}
			if err := p.parseOptionalUnpivotAlias(); err != nil {
				return err
			}
			if !p.match(tokComma) {
				break
			}
		}
		if _, err := p.expect(tokRParen); err != nil {
			return err
		}
	}
	if _, err := p.expect(tokRParen); err != nil {
		return err
	}
	if p.match(kwAS) {
		_, err := p.parseErrorCapturingIdentifier()
		return err
	}
	if p.canBeImplicitTableAlias() {
		_, err := p.parseErrorCapturingIdentifier()
		return err
	}
	return nil
}

func (p *Parser) parseOptionalUnpivotAlias() error {
	if p.match(kwAS) {
		_, err := p.parseErrorCapturingIdentifier()
		return err
	}
	if p.cur().Kind != tokComma && p.cur().Kind != tokRParen && p.cur().Kind.isIdentifier() {
		_, err := p.parseErrorCapturingIdentifier()
		return err
	}
	return nil
}

func (p *Parser) parseWhereClause() error {
	if _, err := p.expect(kwWHERE); err != nil {
		return err
	}
	return p.parseBooleanExpression()
}

func (p *Parser) parseHavingClause() error {
	if _, err := p.expect(kwHAVING); err != nil {
		return err
	}
	return p.parseBooleanExpression()
}

func (p *Parser) parseAggregationClause() error {
	if _, err := p.expect(kwGROUP); err != nil {
		return err
	}
	if _, err := p.expect(kwBY); err != nil {
		return err
	}
	// First item could be groupingAnalytics or namedExpression / expression.
	for {
		if p.isGroupingAnalyticsStart() {
			cpGA := p.save()
			if err := p.parseGroupingAnalytics(); err == nil && !p.canBeImplicitExprAlias() {
				// parsed groupingAnalytics cleanly
			} else {
				p.restore(cpGA)
				if err := p.parseNamedExpression(); err != nil {
					return err
				}
			}
		} else {
			if err := p.parseNamedExpression(); err != nil {
				return err
			}
		}
		if !p.match(tokComma) {
			break
		}
	}
	// Optional WITH ROLLUP | WITH CUBE | GROUPING SETS (...)
	if p.cur().Kind == kwWITH && (p.peek(1).Kind == kwROLLUP || p.peek(1).Kind == kwCUBE) {
		p.advance()
		p.advance()
	} else if p.cur().Kind == kwGROUPING && p.peek(1).Kind == kwSETS {
		p.advance()
		p.advance()
		if _, err := p.expect(tokLParen); err != nil {
			return err
		}
		if err := p.parseGroupingSet(); err != nil {
			return err
		}
		for p.match(tokComma) {
			if err := p.parseGroupingSet(); err != nil {
				return err
			}
		}
		if _, err := p.expect(tokRParen); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) isGroupingAnalyticsStart() bool {
	if (p.cur().Kind == kwROLLUP || p.cur().Kind == kwCUBE) && p.peek(1).Kind == tokLParen {
		return true
	}
	return p.cur().Kind == kwGROUPING && p.peek(1).Kind == kwSETS && p.peek(2).Kind == tokLParen
}

func (p *Parser) parseGroupingAnalytics() error {
	if err := p.enter(); err != nil {
		return err
	}
	defer p.leave()

	if p.match(kwROLLUP, kwCUBE) {
		if _, err := p.expect(tokLParen); err != nil {
			return err
		}
		if err := p.parseGroupingSet(); err != nil {
			return err
		}
		for p.match(tokComma) {
			if err := p.parseGroupingSet(); err != nil {
				return err
			}
		}
		_, err := p.expect(tokRParen)
		return err
	}
	if _, err := p.expect(kwGROUPING); err != nil {
		return err
	}
	if _, err := p.expect(kwSETS); err != nil {
		return err
	}
	if _, err := p.expect(tokLParen); err != nil {
		return err
	}
	if err := p.parseGroupingElement(); err != nil {
		return err
	}
	for p.match(tokComma) {
		if err := p.parseGroupingElement(); err != nil {
			return err
		}
	}
	_, err := p.expect(tokRParen)
	return err
}

func (p *Parser) parseGroupingElement() error {
	if p.isGroupingAnalyticsStart() {
		cp := p.save()
		if err := p.parseGroupingAnalytics(); err == nil {
			return nil
		}
		p.restore(cp)
	}
	return p.parseGroupingSet()
}

func (p *Parser) parseGroupingSet() error {
	if p.cur().Kind == tokLParen {
		cp := p.save()
		p.advance()
		if p.match(tokRParen) {
			return nil
		}
		if err := p.parseExpressionSeq(); err == nil && p.cur().Kind == tokRParen {
			p.advance()
			return nil
		}
		p.restore(cp)
	}
	return p.parseExpression()
}

func (p *Parser) parseWindowClause() error {
	if _, err := p.expect(kwWINDOW); err != nil {
		return err
	}
	for {
		if _, err := p.parseErrorCapturingIdentifier(); err != nil {
			return err
		}
		if _, err := p.expect(kwAS); err != nil {
			return err
		}
		if err := p.parseWindowSpec(); err != nil {
			return err
		}
		if !p.match(tokComma) {
			break
		}
	}
	return nil
}

func (p *Parser) parseWindowSpec() error {
	if !p.match(tokLParen) {
		_, err := p.parseErrorCapturingIdentifier()
		return err
	}
	if p.match(tokRParen) {
		return nil
	}
	// Could be windowRef: (name) or windowDef: (CLUSTER BY ... | (PARTITION|DISTRIBUTE BY ...)? (ORDER|SORT BY ...)? windowFrame?)
	cp := p.save()
	if _, err := p.parseErrorCapturingIdentifier(); err == nil && p.cur().Kind == tokRParen {
		p.advance()
		return nil
	}
	p.restore(cp)

	if p.cur().Kind == kwCLUSTER && p.peek(1).Kind == kwBY {
		p.advance()
		p.advance()
		if err := p.parseExpressionSeq(); err != nil {
			return err
		}
	} else {
		if (p.cur().Kind == kwPARTITION || p.cur().Kind == kwDISTRIBUTE) && p.peek(1).Kind == kwBY {
			p.advance()
			p.advance()
			if err := p.parseExpressionSeq(); err != nil {
				return err
			}
		}
		if (p.cur().Kind == kwORDER || p.cur().Kind == kwSORT) && p.peek(1).Kind == kwBY {
			p.advance()
			p.advance()
			if err := p.parseSortItemList(); err != nil {
				return err
			}
		}
	}
	if p.cur().Kind == kwRANGE || p.cur().Kind == kwROWS {
		if err := p.parseWindowFrame(); err != nil {
			return err
		}
	}
	_, err := p.expect(tokRParen)
	return err
}

func (p *Parser) parseWindowFrame() error {
	if !p.match(kwRANGE, kwROWS) {
		return p.errorf("expected RANGE or ROWS")
	}
	if p.match(kwBETWEEN) {
		if err := p.parseFrameBound(); err != nil {
			return err
		}
		if _, err := p.expect(kwAND); err != nil {
			return err
		}
		return p.parseFrameBound()
	}
	return p.parseFrameBound()
}

func (p *Parser) parseFrameBound() error {
	if p.match(kwUNBOUNDED) {
		if !p.match(kwPRECEDING, kwFOLLOWING) {
			return p.errorf("expected PRECEDING or FOLLOWING after UNBOUNDED")
		}
		return nil
	}
	if p.cur().Kind == kwCURRENT && p.peek(1).Kind == kwROW {
		p.advance()
		p.advance()
		return nil
	}
	if err := p.parseExpression(); err != nil {
		return err
	}
	if !p.match(kwPRECEDING, kwFOLLOWING) {
		return p.errorf("expected PRECEDING or FOLLOWING")
	}
	return nil
}

func (p *Parser) parseSortItemList() error {
	if err := p.parseSortItem(); err != nil {
		return err
	}
	for p.match(tokComma) {
		if err := p.parseSortItem(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parseSortItem() error {
	if err := p.parseExpression(); err != nil {
		return err
	}
	p.match(kwASC, kwDESC)
	if p.match(kwNULLS) {
		if !p.match(kwLAST, kwFIRST) {
			return p.errorf("expected FIRST or LAST after NULLS")
		}
	}
	return nil
}

func (p *Parser) parseQueryOrganization() error {
	if p.cur().Kind == kwORDER && p.peek(1).Kind == kwBY {
		p.advance()
		p.advance()
		if err := p.parseSortItemList(); err != nil {
			return err
		}
	}
	if p.cur().Kind == kwCLUSTER && p.peek(1).Kind == kwBY {
		p.advance()
		p.advance()
		if err := p.parseExpressionSeq(); err != nil {
			return err
		}
	}
	if p.cur().Kind == kwDISTRIBUTE && p.peek(1).Kind == kwBY {
		p.advance()
		p.advance()
		if err := p.parseExpressionSeq(); err != nil {
			return err
		}
	}
	if p.cur().Kind == kwSORT && p.peek(1).Kind == kwBY {
		p.advance()
		p.advance()
		if err := p.parseSortItemList(); err != nil {
			return err
		}
	}
	if p.cur().Kind == kwWINDOW {
		if err := p.parseWindowClause(); err != nil {
			return err
		}
	}
	if p.match(kwLIMIT) {
		if !p.match(kwALL) {
			if err := p.parseExpression(); err != nil {
				return err
			}
		}
	}
	if p.match(kwOFFSET) {
		if err := p.parseExpression(); err != nil {
			return err
		}
	}
	return nil
}
