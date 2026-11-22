package parser

import "strconv"

// parseQuery parses:
//
//	query : ctes? queryTerm queryOrganization ;
func (p *Parser) parseQuery() error {
	if err := p.enterNesting(); err != nil {
		return err
	}
	defer p.leaveNesting()

	if p.cur().Kind == kwWITH {
		if err := p.parseCTEs(); err != nil {
			return err
		}
	}
	if err := p.parseQueryTerm(); err != nil {
		return err
	}
	return p.parseQueryOrganization(false)
}

func (p *Parser) looksLikeQueryStart() bool {
	return p.startsWithQueryKeywordAt(0, false)
}

// startsWithQueryKeywordAfterParens scans iteratively past any leading '('
// tokens and reports whether a query-introducing keyword follows. Scanning
// iteratively avoids exponential backtracking on deeply nested parenthesized
// expressions or parenthesized table relations.
func (p *Parser) startsWithQueryKeywordAfterParens() bool {
	return p.startsWithQueryKeywordAt(0, true)
}

// startsWithCompleteParenQuery reports whether the tokens starting at p.cur()
// (immediately after an opening '(') can form a complete parenthesized query
// up to its matching ')'. When p.cur() is itself '(', it checks that the
// token after the inner matching ')' is either ')' or a query continuation
// (set operator, pipe operator, or query organization clause), so expressions
// and row constructors like `((SELECT 1), 1)` or `((SELECT 1) + 1)` skip
// parseQuery() without reparsing the inner subquery.
func (p *Parser) startsWithCompleteParenQuery() bool {
	if !p.startsWithQueryKeywordAfterParens() {
		return false
	}
	if p.cur().Kind != tokLParen {
		return true
	}
	leadingParens := p.skipLParensOffset(0)
	checkClose := func(d int) bool {
		closeOff := p.matchingCloseOffset(d)
		if closeOff < 0 || p.peek(closeOff).Kind != tokRParen {
			return false
		}
		switch p.peek(closeOff + 1).Kind {
		case tokRParen, tokOperatorPipe,
			kwUNION, kwINTERSECT, kwEXCEPT, kwMINUS,
			kwORDER, kwCLUSTER, kwDISTRIBUTE, kwSORT, kwWINDOW, kwLIMIT, kwOFFSET:
			return true
		default:
			return false
		}
	}
	if !checkClose(0) {
		return false
	}
	for d := leadingParens - 1; d >= 1; d-- {
		if !checkClose(d) {
			return false
		}
	}
	return true
}

// isValuesQueryAt reports whether `VALUES` at offset i can introduce an
// inlineTable query (`VALUES expr, ... [tableAlias]`) rather than a column
// reference `values` or a `VALUES(...)` function call (with empty/star/quantified/
// named args, window/filter/within-group/nulls suffixes, or row-constructor comma
// after an alias). This avoids exponential reparsing on nested parenthesized
// `(VALUES(<expr>) OVER ())` calls.
func (p *Parser) isValuesQueryAt(i int, inParens bool) bool {
	next := p.peek(i + 1).Kind
	switch next {
	case tokEOF, tokSemicolon, tokRParen, tokRBracket, tokComma, tokOperatorPipe,
		tokDot, tokLBracket, tokDoubleColon, kwCOLLATE,
		tokAsterisk, tokSlash, tokPercent, kwDIV,
		tokConcatPipe, tokPipe, tokAmpersand, tokHat,
		tokShiftLeft, tokShiftRight, tokShiftRightU,
		tokEq, tokNeq, tokNeqBang, tokNullSafeEq, tokLT, tokLTE, tokGT, tokGTE,
		tokArrow, tokFatArrow, kwAS:
		return false
	case kwAND, kwOR, kwIS, kwIN, kwLIKE, kwILIKE, kwRLIKE, kwBETWEEN:
		if !inParens {
			return true
		}
		return false
	case tokLParen:
		inner := p.peek(i + 2).Kind
		if inner == tokRParen || inner == tokAsterisk {
			return false
		}
		if (inner == kwDISTINCT || inner == kwALL) && p.peek(i+3).Kind != tokRParen && p.peek(i+3).Kind != tokComma {
			return false
		}
		// Scan depth-1 tokens inside VALUES(...), checking for named arguments
		// (`k => v`), TABLE arguments (`TABLE ...`), and function-call suffixes
		// immediately after ')'.
		closeIdx := p.matchingCloseOffset(i + 1)
		if closeIdx < 0 {
			return false
		}
		for j := i + 2; j < closeIdx; j++ {
			k := p.peek(j).Kind
			if k == tokLParen || k == tokLBracket {
				c := p.matchingCloseOffset(j)
				if c < 0 {
					return false
				}
				j = c
				continue
			}
			if k == tokFatArrow || (k == kwTABLE && (j == i+2 || p.peek(j-1).Kind == tokComma)) {
				return false
			}
		}
		after := p.peek(closeIdx + 1).Kind
		switch after {
		case kwOVER:
			return false
		case kwFILTER:
			if p.peek(closeIdx+2).Kind == tokLParen {
				return false
			}
		case kwWITHIN:
			if p.peek(closeIdx+2).Kind == kwGROUP {
				return false
			}
		case kwIGNORE, kwRESPECT:
			if p.peek(closeIdx+2).Kind == kwNULLS {
				return false
			}
		}
	}
	if !inParens {
		return true
	}
	// In `(VALUES(...) AS a, 2)` or `(VALUES(...) a, 2)`, a depth-0 comma after
	// an alias indicates a row constructor rather than an inlineTable (whose
	// optional tableAlias only appears after all comma-separated expressions).
	for j := i + 1; ; j++ {
		k := p.peek(j).Kind
		switch k {
		case tokEOF, tokSemicolon, tokRParen, tokRBracket:
			return true
		case tokLParen, tokLBracket:
			c := p.matchingCloseOffset(j)
			if c < 0 {
				return true
			}
			j = c
			// `(...) alias ,` where alias is an identifier (not a query clause keyword)
			a := p.peek(j + 1).Kind
			if a == kwAS && p.peek(j+2).Kind.isIdentifier() && p.peek(j+3).Kind == tokComma {
				return false
			}
			if a.isIdentifier() && a != kwAS && a != kwUNION && a != kwINTERSECT && a != kwEXCEPT && a != kwMINUS &&
				a != kwORDER && a != kwCLUSTER && a != kwDISTRIBUTE && a != kwSORT && a != kwWINDOW && a != kwLIMIT && a != kwOFFSET &&
				p.peek(j+2).Kind == tokComma {
				return false
			}
		case kwAS:
			if p.peek(j+1).Kind.isIdentifier() && p.peek(j+2).Kind == tokComma {
				return false
			}
		}
	}
}

func (p *Parser) hasKeywordAtDepthZero(startOffset int, target TokenKind) bool {
	for i := startOffset; ; i++ {
		k := p.peek(i).Kind
		switch k {
		case tokEOF, tokSemicolon, tokRParen, tokRBracket:
			return false
		case tokLParen, tokLBracket:
			c := p.matchingCloseOffset(i)
			if c < 0 {
				return false
			}
			i = c
		default:
			if k == target {
				return true
			}
		}
	}
}

func (p *Parser) parseCTEs() error {
	if _, err := p.expect(kwWITH); err != nil {
		return err
	}
	// RECURSIVE is non-reserved, so `WITH recursive AS (...)` and
	// `WITH recursive(x) AS (...)` name the CTE `recursive` rather than
	// using the RECURSIVE modifier.
	if p.cur().Kind == kwRECURSIVE && p.peek(1).Kind != kwAS && p.peek(1).Kind != tokLParen {
		p.advance()
	}
	seen := make(map[string]struct{})
	name, err := p.parseNamedQuery()
	if err != nil {
		return err
	}
	seen[exactIdentName(name)] = struct{}{}
	for p.match(tokComma) {
		name, err := p.parseNamedQuery()
		if err != nil {
			return err
		}
		key := exactIdentName(name)
		if _, dup := seen[key]; dup {
			return p.errorf("duplicate CTE name %q in WITH clause", name)
		}
		seen[key] = struct{}{}
	}
	return nil
}

func (p *Parser) parseNamedQuery() (string, error) {
	name, err := p.parseErrorCapturingIdentifier()
	if err != nil {
		return "", err
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
		return "", err
	}
	if err := p.parseQuery(); err != nil {
		return "", err
	}
	_, err = p.expect(tokRParen)
	return name, err
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
		// Spark v4.0.5 requires at least one aggregate expression or a GROUP BY
		// clause after |> AGGREGATE.
		if p.cur().Kind == kwGROUP && p.peek(1).Kind == kwBY {
			cpAgg := p.save()
			if err := p.parseAggregationClause(true); err == nil && (p.isStatementEnd() || p.cur().Kind == tokRParen || p.cur().Kind == tokOperatorPipe || p.isQueryOrgStart()) {
				return nil
			}
			p.restore(cpAgg)
		}
		hasAggExprs := false
		if p.cur().Kind != tokEOF && p.cur().Kind != tokSemicolon &&
			p.cur().Kind != tokRParen && p.cur().Kind != tokOperatorPipe && !p.isQueryOrgStart() {
			if err := p.parseNamedExpressionSeq(); err != nil {
				return err
			}
			hasAggExprs = true
		}
		if p.cur().Kind == kwGROUP {
			return p.parseAggregationClause(true)
		}
		if !hasAggExprs {
			return p.errorf("AGGREGATE pipe operator requires aggregate expressions or a GROUP BY clause")
		}
		return nil
	case kwJOIN, kwINNER, kwCROSS, kwLEFT, kwRIGHT, kwFULL, kwSEMI, kwANTI, kwNATURAL:
		return p.parseJoinRelation()
	default:
		// queryOrganization in a pipe operator forbids WINDOW and combining
		// ORDER/CLUSTER/DISTRIBUTE/SORT BY with LIMIT/OFFSET in one pipe stage.
		return p.parseQueryOrganization(true)
	}
}

func (p *Parser) isQueryOrgStart() bool {
	switch p.cur().Kind {
	case kwORDER, kwCLUSTER, kwDISTRIBUTE, kwSORT:
		return p.peek(1).Kind == kwBY
	case kwWINDOW:
		return p.peek(1).Kind.isIdentifier() && p.peek(2).Kind == kwAS
	case kwLIMIT, kwOFFSET:
		switch p.peek(1).Kind {
		case tokEOF, tokSemicolon, tokRParen, tokRBracket, tokComma, tokOperatorPipe,
			tokDot, tokLBracket, tokDoubleColon, kwCOLLATE,
			tokAsterisk, tokSlash, tokPercent, kwDIV,
			tokConcatPipe, tokPipe, tokAmpersand, tokHat,
			tokShiftLeft, tokShiftRight, tokShiftRightU,
			tokEq, tokNeq, tokNeqBang, tokNullSafeEq, tokLT, tokLTE, tokGT, tokGTE,
			kwAND, kwOR, kwIS, kwIN, kwLIKE, kwILIKE, kwRLIKE, kwBETWEEN,
			kwAS, kwGROUP, tokArrow, tokFatArrow:
			return false
		case tokPlus, tokMinus:
			return !p.hasKeywordAtDepthZero(2, kwGROUP)
		}
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
	prev := p.inFromStatement
	p.inFromStatement = true
	err := p.parseFromClause()
	p.inFromStatement = prev
	if err != nil {
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
		return p.parseQueryOrganization(false)
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
		if err := p.parseAggregationClause(false); err != nil {
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
	return p.parseQueryOrganization(false)
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
		if err := p.parseAggregationClause(false); err != nil {
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
	closeOff := p.matchingCloseOffset(lparenOffset)
	if closeOff < 0 || p.peek(closeOff).Kind != tokRParen {
		return false
	}
	after := p.peek(closeOff + 1).Kind
	return after == kwUSING ||
		(after == kwROW && p.peek(closeOff+2).Kind == kwFORMAT) ||
		after == kwRECORDWRITER
}

func (p *Parser) parseTransformClause() error {
	if p.match(kwSELECT) {
		if _, err := p.expect(kwTRANSFORM); err != nil {
			return err
		}
		if _, err := p.expect(tokLParen); err != nil {
			return err
		}
		if p.matchSetQuantifierIfFollowedByExpr() {
			return p.errorf("TRANSFORM with DISTINCT/ALL is not supported")
		}
		if err := p.parseExpressionSeq(); err != nil {
			return err
		}
		if _, err := p.expect(tokRParen); err != nil {
			return err
		}
	} else if p.match(kwMAP, kwREDUCE) {
		if p.matchSetQuantifierIfFollowedByExpr() {
			return p.errorf("TRANSFORM with DISTINCT/ALL is not supported")
		}
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
	if p.cur().Kind == kwRECORDWRITER {
		return p.errorf("custom RECORDWRITER is not supported in TRANSFORM")
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
	if p.cur().Kind == kwRECORDREADER {
		return p.errorf("custom RECORDREADER is not supported in TRANSFORM")
	}
	return nil
}

func (p *Parser) matchSetQuantifierIfFollowedByExpr() bool {
	if p.cur().Kind != kwDISTINCT && p.cur().Kind != kwALL {
		return false
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
			return true
		case tokIdent:
			if isWord(p.peek(2), "QUALIFY") && p.peek(3).Kind != tokEOF && p.peek(3).Kind != tokSemicolon && p.peek(3).Kind != tokRParen {
				p.advance()
				return true
			}
		}
		return false
	case tokLParen:
		if p.peek(2).Kind == tokRParen {
			return false
		}
	case kwLATERAL:
		if p.peek(2).Kind == kwVIEW {
			return false
		}
	case tokEOF, tokSemicolon, tokComma, tokRParen, tokOperatorPipe,
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
		return false
	}
	p.advance()
	return true
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
	if err := p.parseRelation(false); err != nil {
		return err
	}
	for p.match(tokComma) {
		if err := p.parseRelation(true); err != nil {
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

func (p *Parser) parseRelation(isJoinedOrComma bool) error {
	if err := p.enterNesting(); err != nil {
		return err
	}
	defer p.leaveNesting()

	isLateral := false
	if p.isLateralRelationModifier(isJoinedOrComma) {
		isLateral = true
		p.advance()
	}
	isSubqueryOrTVF, err := p.parseRelationPrimary()
	if err != nil {
		return err
	}
	// Spark v4.0.5 AstBuilder only enforces LATERAL_WITHOUT_SUBQUERY_OR_TABLE_VALUED_FUNC
	// when LATERAL joins onto a preceding relation (via `, LATERAL ...` or `JOIN LATERAL ...`).
	if isJoinedOrComma && isLateral && !isSubqueryOrTVF {
		return p.errorf("LATERAL can only be used with subqueries and table-valued functions")
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
// relationPrimary acts as the LATERAL modifier (e.g. `FROM t, LATERAL (SELECT ...)`,
// `JOIN LATERAL explode(arr)`, `JOIN LATERAL x`) rather than a table or schema
// named `lateral` (e.g. `JOIN lateral ON ...`, `JOIN lateral.t ON ...`).
func (p *Parser) isLateralRelationModifier(inJoinOrComma bool) bool {
	if p.cur().Kind != kwLATERAL {
		return false
	}
	next := p.peek(1).Kind
	switch next {
	case tokEOF, tokSemicolon, tokComma, tokRParen, tokOperatorPipe,
		tokDot, kwVIEW,
		kwON, kwUSING,
		kwPIVOT, kwUNPIVOT,
		kwJOIN, kwINNER, kwCROSS, kwLEFT, kwRIGHT, kwFULL, kwSEMI, kwANTI, kwNATURAL,
		kwWHERE, kwGROUP, kwHAVING, kwWINDOW, kwORDER, kwCLUSTER, kwDISTRIBUTE,
		kwSORT, kwLIMIT, kwOFFSET, kwUNION, kwINTERSECT, kwEXCEPT, kwMINUS,
		kwSELECT, kwINSERT, kwMAP, kwREDUCE, kwSET, kwINTO:
		return false
	case kwTABLESAMPLE:
		// `FROM LATERAL TABLESAMPLE(10 PERCENT)` parses LATERAL as table + TABLESAMPLE clause,
		// whereas `FROM LATERAL TABLESAMPLE(expr)` (non-sample args) parses LATERAL + TVF `tablesample(expr)`.
		if p.peek(2).Kind != tokLParen {
			return false
		}
		cpSample := p.save()
		p.advance() // LATERAL
		err := p.parseSample()
		p.restore(cpSample)
		return err != nil
	case kwAS:
		// In Spark's joinRelation/comma relation, ANTLR parses `JOIN lateral AS l`
		// as LATERAL + table `AS` + alias `l` (which AstBuilder then rejects as
		// LATERAL_WITHOUT_SUBQUERY_OR_TABLE_VALUED_FUNC).
		return inJoinOrComma && p.peek(2).Kind.isStrictIdentifier()
	case kwFOR:
		return p.peek(2).Kind != kwSYSTEM_TIME && p.peek(2).Kind != kwSYSTEM_VERSION && p.peek(2).Kind != kwVERSION
	case kwTIMESTAMP, kwVERSION, kwSYSTEM_TIME, kwSYSTEM_VERSION:
		return !(p.peek(2).Kind == kwAS && p.peek(3).Kind == kwOF)
	case kwWITH:
		return p.peek(2).Kind != tokLParen
	case kwLATERAL:
		return p.peek(2).Kind != kwVIEW
	case tokLParen:
		return p.isLateralParenRelation()
	case tokIdent:
		if isWord(p.peek(1), "QUALIFY") && p.peek(2).Kind != tokEOF && p.peek(2).Kind != tokSemicolon && p.peek(2).Kind != tokRParen {
			return false
		}
	}
	return true
}

// isLateralParenRelation reports whether `LATERAL (...)` (with p.cur() at
// kwLATERAL and p.peek(1) at tokLParen) is the LATERAL modifier before a
// parenthesized query or relation rather than a table-valued function call
// named `lateral(...)`.
func (p *Parser) isLateralParenRelation() bool {
	if p.peek(2).Kind == tokRParen {
		return false
	}
	if p.peek(2).Kind == kwLATERAL {
		return true
	}
	cp := p.save()
	p.advance() // LATERAL
	p.advance() // '('
	isQuery := p.startsWithCompleteParenQuery()
	p.restore(cp)
	if isQuery {
		return true
	}

	// Find the matching ')' for the opening '(' at offset 1 and determine how
	// many outermost paren pairs wrap the entire content.
	closeIdx := p.matchingCloseOffset(1)
	if closeIdx < 0 || p.peek(closeIdx).Kind != tokRParen {
		return false
	}
	if p.peek(closeIdx+1).Kind == kwTABLESAMPLE {
		return true
	}
	leadingParens := p.skipLParensOffset(1) - 1
	minDepth := 1
	for minDepth < leadingParens &&
		p.matchingCloseOffset(1+minDepth) == closeIdx-minDepth &&
		p.peek(closeIdx-minDepth).Kind == tokRParen {
		minDepth++
	}

	lo := 1 + minDepth
	hi := closeIdx - (minDepth - 1)
	if lo >= hi {
		return false
	}
	if p.peek(lo).Kind == kwLATERAL {
		return true
	}

	// Scan top-level tokens inside [lo, hi) for relation-only or function-arg-only markers.
	hasRelationMarker := false
	for j := lo; j < hi; j++ {
		k := p.peek(j).Kind
		if k == tokLParen || k == tokLBracket {
			c := p.matchingCloseOffset(j)
			if c < 0 || c >= hi {
				break
			}
			j = c
			continue
		}
		switch k {
		case tokComma, tokFatArrow:
			return false
		case kwTABLE:
			if j == lo {
				return false
			}
		case kwJOIN, kwINNER, kwCROSS, kwLEFT, kwRIGHT, kwFULL, kwSEMI, kwANTI, kwNATURAL,
			kwPIVOT, kwUNPIVOT, kwTABLESAMPLE, kwAS:
			hasRelationMarker = true
		case kwFOR:
			n := p.peek(j + 1).Kind
			if n == kwSYSTEM_VERSION || n == kwVERSION || n == kwSYSTEM_TIME || n == kwTIMESTAMP {
				hasRelationMarker = true
			}
		case kwWITH:
			if p.peek(j+1).Kind == tokLParen {
				hasRelationMarker = true
			}
		}
	}
	if hasRelationMarker {
		return true
	}
	// Bare `LATERAL (t alias)` with two plain identifiers and no operator between them.
	isPlainIdent := func(k TokenKind) bool {
		return k == tokIdent || k == tokBackquoted || k == tokDoubleQuoted
	}
	pos := lo
	if isPlainIdent(p.peek(pos).Kind) {
		pos++
		for p.peek(pos).Kind == tokDot && p.peek(pos+1).Kind.isIdentifier() {
			pos += 2
		}
		if isPlainIdent(p.peek(pos).Kind) && pos+1 == hi {
			return true
		}
	}
	return false
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

type joinKind int

const (
	joinInner joinKind = iota
	joinCross
	joinLeftOuter
	joinRightOuter
	joinFullOuter
	joinSemi
	joinAnti
)

func (p *Parser) parseJoinRelation() error {
	natural := p.match(kwNATURAL)
	jt, err := p.parseJoinType()
	if err != nil {
		return err
	}
	if natural && jt == joinCross {
		return p.errorf("NATURAL cannot be combined with CROSS JOIN")
	}
	if _, err := p.expect(kwJOIN); err != nil {
		return err
	}
	isLateral := false
	if p.isLateralRelationModifier(true) {
		isLateral = true
		p.advance()
	}
	if natural && isLateral {
		return p.errorf("NATURAL JOIN cannot be combined with LATERAL")
	}
	if isLateral && jt != joinInner && jt != joinCross && jt != joinLeftOuter {
		return p.errorf("LATERAL join is only supported for INNER, CROSS, and LEFT OUTER joins")
	}
	isSubqueryOrTVF, err := p.parseRelationPrimary()
	if err != nil {
		return err
	}
	if isLateral && !isSubqueryOrTVF {
		return p.errorf("LATERAL can only be used with subqueries and table-valued functions")
	}
	if !natural {
		if p.match(kwON) {
			return p.parseBooleanExpression()
		}
		if p.match(kwUSING) {
			if isLateral {
				return p.errorf("LATERAL join with USING criteria is not supported")
			}
			return p.parseIdentifierList()
		}
	}
	return nil
}

func (p *Parser) parseJoinType() (joinKind, error) {
	switch p.cur().Kind {
	case kwINNER:
		p.advance()
		return joinInner, nil
	case kwCROSS:
		p.advance()
		return joinCross, nil
	case kwLEFT:
		p.advance()
		if p.match(kwSEMI) {
			return joinSemi, nil
		}
		if p.match(kwANTI) {
			return joinAnti, nil
		}
		p.match(kwOUTER)
		return joinLeftOuter, nil
	case kwRIGHT:
		p.advance()
		p.match(kwOUTER)
		return joinRightOuter, nil
	case kwFULL:
		p.advance()
		p.match(kwOUTER)
		return joinFullOuter, nil
	case kwSEMI:
		p.advance()
		return joinSemi, nil
	case kwANTI:
		p.advance()
		return joinAnti, nil
	}
	return joinInner, nil
}

// parseRelationPrimary parses a relationPrimary and reports whether it is a
// subquery or table-valued function (the only forms Spark v4.0.5 allows after LATERAL).
func (p *Parser) parseRelationPrimary() (bool, error) {
	switch p.cur().Kind {
	case kwVALUES:
		if !p.isValuesFollowedByRelationClause() && p.isValuesQueryAt(0, false) {
			if p.peek(1).Kind == tokLParen {
				return false, p.parseInlineTable()
			}
			cp := p.save()
			if err := p.parseInlineTable(); err == nil && p.isValidRelationPrimaryEnd() {
				return false, nil
			}
			p.restore(cp)
		}

	case tokLParen:
		p.advance()
		// Commit directly to (query) once identified as a complete parenthesized query
		// so nested malformed subqueries in FROM never reparse exponentially.
		if p.startsWithCompleteParenQuery() {
			if err := p.parseQuery(); err != nil {
				return false, err
			}
			if _, err := p.expect(tokRParen); err != nil {
				return false, err
			}
			if p.cur().Kind == kwTABLESAMPLE && p.peek(1).Kind == tokLParen {
				if err := p.parseSample(); err != nil {
					return false, err
				}
			}
			return true, p.parseTableAlias()
		}
		if err := p.parseRelation(true); err != nil {
			return false, err
		}
		if _, err := p.expect(tokRParen); err != nil {
			return false, err
		}
		if p.cur().Kind == kwTABLESAMPLE && p.peek(1).Kind == tokLParen {
			if err := p.parseSample(); err != nil {
				return false, err
			}
		}
		return false, p.parseTableAlias()
	}

	// functionTable (`funcName(...)` or `IDENTIFIER(expr)(...)`) or dynamic tableName (`IDENTIFIER(expr) ...`).
	// parseFunctionTable handles all of these in a single pass, so commit directly without retrying.
	if p.canStartFunctionCall() {
		return p.parseFunctionTable()
	}

	// tableName: identifierReference temporalClause? optionsClause? sample? tableAlias
	if _, err := p.parseIdentifierReference(); err != nil {
		return false, err
	}
	if p.isTemporalClauseStart() {
		if err := p.parseTemporalClause(); err != nil {
			return false, err
		}
	}
	if p.cur().Kind == kwWITH && p.peek(1).Kind == tokLParen {
		if err := p.parseOptionsClause(); err != nil {
			return false, err
		}
	}
	if p.cur().Kind == kwTABLESAMPLE && p.peek(1).Kind == tokLParen {
		if err := p.parseSample(); err != nil {
			return false, err
		}
	}
	return false, p.parseTableAlias()
}

func (p *Parser) isValuesFollowedByRelationClause() bool {
	next := p.peek(1).Kind
	switch next {
	case kwFOR:
		n := p.peek(2).Kind
		return n == kwSYSTEM_TIME || n == kwSYSTEM_VERSION || n == kwVERSION || n == kwTIMESTAMP
	case kwSYSTEM_TIME, kwSYSTEM_VERSION, kwVERSION, kwTIMESTAMP:
		return p.peek(2).Kind == kwAS && p.peek(3).Kind == kwOF
	case kwWITH:
		return p.peek(2).Kind == tokLParen
	case kwLATERAL:
		return p.peek(2).Kind == kwVIEW
	case kwGROUP, kwORDER, kwCLUSTER, kwDISTRIBUTE, kwSORT:
		return p.peek(2).Kind == kwBY
	case kwWINDOW:
		return p.peek(2).Kind.isIdentifier() && p.peek(3).Kind == kwAS
	case kwPIVOT, kwUNPIVOT, kwTABLESAMPLE:
		return p.peek(2).Kind == tokLParen || (next == kwUNPIVOT && (p.peek(2).Kind == kwINCLUDE || p.peek(2).Kind == kwEXCLUDE))
	case kwJOIN, kwINNER, kwCROSS, kwLEFT, kwRIGHT, kwFULL, kwSEMI, kwANTI, kwNATURAL:
		return !p.isStatementEndAt(2) && p.peek(2).Kind != tokRParen && p.peek(2).Kind != tokComma
	case kwON, kwUSING, kwWHERE, kwHAVING, kwLIMIT, kwOFFSET,
		kwUNION, kwINTERSECT, kwEXCEPT, kwMINUS,
		kwSELECT, kwFROM, kwMAP, kwREDUCE, kwSET, kwINTO:
		return !p.isStatementEndAt(2) && p.peek(2).Kind != tokRParen && p.peek(2).Kind != tokComma
	case tokIdent:
		if isWord(p.peek(1), "QUALIFY") && !p.isStatementEndAt(2) && p.peek(2).Kind != tokRParen && p.peek(2).Kind != tokComma {
			return true
		}
	}
	return false
}

func (p *Parser) isValidRelationPrimaryEnd() bool {
	k := p.cur().Kind
	switch k {
	case tokEOF, tokSemicolon, tokComma, tokRParen, tokOperatorPipe,
		kwJOIN, kwINNER, kwCROSS, kwLEFT, kwRIGHT, kwFULL, kwSEMI, kwANTI, kwNATURAL,
		kwPIVOT, kwUNPIVOT, kwLATERAL, kwON, kwUSING,
		kwWHERE, kwGROUP, kwHAVING, kwWINDOW, kwORDER, kwCLUSTER, kwDISTRIBUTE,
		kwSORT, kwLIMIT, kwOFFSET, kwUNION, kwINTERSECT, kwEXCEPT, kwMINUS,
		kwSELECT, kwFROM, kwMAP, kwREDUCE, kwSET, kwINTO:
		return true
	case tokIdent:
		return isWord(p.cur(), "QUALIFY")
	}
	return false
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
		prevDisallow := p.disallowColumnRefs
		p.disallowColumnRefs = true
		err := p.parseValueExpression()
		p.disallowColumnRefs = prevDisallow
		return err
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

func (p *Parser) parseFunctionTable() (bool, error) {
	// Single-pass handling for IDENTIFIER(expr)(tableArgs) vs IDENTIFIER(tableArgs)
	// without reparsing nested IDENTIFIER(...) arguments.
	// Spark v4.0.5 rejects empty `IDENTIFIER()` with PARSE_SYNTAX_ERROR.
	if p.cur().Kind == kwIDENTIFIER && p.peek(1).Kind == tokLParen {
		p.advance()
		p.advance()
		if p.cur().Kind == tokRParen {
			return false, p.errorf("IDENTIFIER() requires an expression argument")
		}
		isTableOrNamed := p.isTableSubqueryArgAt(0) ||
			(p.cur().Kind.isIdentifier() && p.peek(1).Kind == tokFatArrow)
		if err := p.parseFunctionTableArgument(); err != nil {
			return false, err
		}
		argCount := 1
		for p.match(tokComma) {
			if err := p.parseFunctionTableArgument(); err != nil {
				return false, err
			}
			argCount++
		}
		if _, err := p.expect(tokRParen); err != nil {
			return false, err
		}
		singlePlainArg := !isTableOrNamed && argCount == 1
		if singlePlainArg && p.cur().Kind == tokLParen {
			p.advance()
			if !p.match(tokRParen) {
				for {
					if err := p.parseFunctionTableArgument(); err != nil {
						return false, err
					}
					if !p.match(tokComma) {
						break
					}
				}
				if _, err := p.expect(tokRParen); err != nil {
					return false, err
				}
			}
			return true, p.parseTableAlias()
		}
		if singlePlainArg {
			// IDENTIFIER(expr) without a second `(...)` is an identifierReference
			// tableName, which accepts temporalClause?, optionsClause?, sample?, tableAlias.
			if p.isTemporalClauseStart() {
				if err := p.parseTemporalClause(); err != nil {
					return false, err
				}
			}
			if p.cur().Kind == kwWITH && p.peek(1).Kind == tokLParen {
				if err := p.parseOptionsClause(); err != nil {
					return false, err
				}
			}
			if p.cur().Kind == kwTABLESAMPLE && p.peek(1).Kind == tokLParen {
				if err := p.parseSample(); err != nil {
					return false, err
				}
			}
			return false, p.parseTableAlias()
		}
		return true, p.parseTableAlias()
	}

	// Spark v4.0.5 AstBuilder.visitTableValuedFunction rejects qualified TVF names
	// (e.g. `FROM lateral.f(1)`, `FROM db.f(1)`): only an unqualified function name
	// is valid here.
	if _, err := p.parseIdentifier(); err != nil {
		return false, err
	}
	if p.cur().Kind == tokDot {
		return false, p.errorf("table-valued function name cannot be qualified")
	}
	if _, err := p.expect(tokLParen); err != nil {
		return false, err
	}
	if !p.match(tokRParen) {
		for {
			if err := p.parseFunctionTableArgument(); err != nil {
				return false, err
			}
			if !p.match(tokComma) {
				break
			}
		}
		if _, err := p.expect(tokRParen); err != nil {
			return false, err
		}
	}
	return true, p.parseTableAlias()
}

func (p *Parser) parseFunctionTableArgument() error {
	if p.cur().Kind.isIdentifier() && p.peek(1).Kind == tokFatArrow && p.isTableSubqueryArgAt(2) {
		p.advance()
		p.advance()
		return p.parseFunctionTableSubqueryArgument()
	}
	if p.isTableSubqueryArgAt(0) {
		return p.parseFunctionTableSubqueryArgument()
	}
	return p.parseFunctionArgument()
}

// isTableSubqueryArgAt reports whether `TABLE` at offset tOff starts a
// functionTableSubqueryArgument (`TABLE(identRef)` or `TABLE(query)`, with
// optional partitioning) rather than a bare column `table` or an ordinary
// `TABLE(...)` function call (e.g. `TABLE((query), k => 1)` or `TABLE(x) + 1`).
// Spark v4.0.5 rejects unparenthesized `TABLE identRef` in AstBuilder.
func (p *Parser) isTableSubqueryArgAt(tOff int) bool {
	if p.peek(tOff).Kind != kwTABLE {
		return false
	}
	if p.peek(tOff+1).Kind != tokLParen {
		return false
	}

	// `TABLE ( ... )`: find matching ')' of TABLE(...).
	outerClose := p.matchingCloseOffset(tOff + 1)
	if outerClose < 0 {
		return false
	}
	if outerClose == tOff+2 {
		return false // empty `TABLE()`
	}

	// Check that what follows `TABLE(...)` can end a functionTableSubqueryArgument.
	afterOuter := p.peek(outerClose + 1).Kind
	switch afterOuter {
	case tokRParen, tokComma:
	case kwWITH:
		if p.peek(outerClose+2).Kind != kwSINGLE || p.peek(outerClose+3).Kind != kwPARTITION {
			return false
		}
	case kwPARTITION, kwDISTRIBUTE:
		if p.peek(outerClose+2).Kind != kwBY {
			return false
		}
	default:
		return false
	}

	// Validate inside `TABLE ( <inner> )`: must be either `query` or `identifierReference`.
	innerOff := tOff + 2
	innerKind := p.peek(innerOff).Kind
	switch innerKind {
	case tokLParen:
		// Parenthesized query `TABLE((SELECT ...) [UNION ...])`: verify matching ')'
		// of the inner '(' is followed by a query continuation (or outerClose),
		// rejecting multi-arg calls like `TABLE((SELECT 1), k => 1)`.
		innerClose := p.matchingCloseOffset(innerOff)
		if innerClose < 0 || innerClose >= outerClose {
			return false
		}
		if innerClose+1 == outerClose {
			// `TABLE((...))`: confirm inner '(' starts a query.
			return p.startsWithQueryKeywordAt(innerOff, true)
		}
		switch p.peek(innerClose + 1).Kind {
		case kwUNION, kwINTERSECT, kwEXCEPT, kwMINUS, tokOperatorPipe,
			kwORDER, kwCLUSTER, kwDISTRIBUTE, kwSORT, kwWINDOW, kwLIMIT, kwOFFSET:
			return p.startsWithQueryKeywordAt(innerOff, true)
		default:
			return false
		}
	case kwSELECT, kwFROM, kwWITH, kwTABLE, kwVALUES, kwMAP, kwREDUCE:
		if p.startsWithQueryKeywordAt(innerOff, true) {
			return true
		}
		// Could still be a table named `values`/`table`/`map`/`reduce` as `TABLE(values)` or `TABLE(values.t)`.
		return p.isIdentRefSpan(innerOff, outerClose)
	case kwIDENTIFIER:
		if p.peek(innerOff+1).Kind == tokLParen {
			// `TABLE(IDENTIFIER(expr))`: ensure the matching ')' is outerClose-1.
			return p.matchingCloseOffset(innerOff+1)+1 == outerClose
		}
		return p.isIdentRefSpan(innerOff, outerClose)
	default:
		return p.isIdentRefSpan(innerOff, outerClose)
	}
}

func (p *Parser) startsWithQueryKeywordAt(startOffset int, inParens bool) bool {
	i := p.skipLParensOffset(startOffset)
	if i > startOffset {
		inParens = true
	}
	switch p.peek(i).Kind {
	case kwSELECT:
		return p.isSelectQueryAt(i, inParens)
	case kwFROM:
		return p.isFromQueryAt(i)
	case kwVALUES:
		return p.isValuesQueryAt(i, inParens)
	case kwWITH:
		pos := i + 1
		if p.peek(pos).Kind == kwRECURSIVE && p.peek(pos+1).Kind != kwAS && p.peek(pos+1).Kind != tokLParen {
			pos++
		}
		if !p.peek(pos).Kind.isIdentifier() {
			return false
		}
		after := p.peek(pos + 1).Kind
		if after == kwAS {
			return p.peek(pos+2).Kind == tokLParen
		}
		return after == tokLParen && p.peek(pos).Kind != kwIN
	case kwTABLE:
		return p.isTableQueryAt(i)
	case kwMAP, kwREDUCE:
		return p.hasKeywordAtDepthZero(i+1, kwUSING)
	}
	return false
}

func (p *Parser) isSelectQueryAt(selectOffset int, inParens bool) bool {
	next := p.peek(selectOffset + 1).Kind
	switch next {
	case tokEOF, tokSemicolon, tokRParen, tokRBracket, tokComma, tokOperatorPipe,
		tokDot, tokLBracket, tokDoubleColon,
		tokSlash, tokPercent,
		tokConcatPipe, tokPipe, tokAmpersand, tokHat,
		tokShiftLeft, tokShiftRight, tokShiftRightU,
		tokEq, tokNeq, tokNeqBang, tokNullSafeEq, tokLT, tokLTE, tokGT, tokGTE,
		tokArrow, tokFatArrow:
		return false
	case tokAsterisk:
		// `SELECT * ...` vs multiplication `(select * 2)` / `(select * x)`.
		switch p.peek(selectOffset + 2).Kind {
		case tokEOF, tokSemicolon, tokRParen, tokComma, tokOperatorPipe,
			kwEXCEPT, kwFROM, kwWHERE, kwGROUP, kwHAVING, kwWINDOW,
			kwORDER, kwCLUSTER, kwDISTRIBUTE, kwSORT, kwLIMIT, kwOFFSET,
			kwUNION, kwINTERSECT, kwMINUS, kwLATERAL:
			return true
		case tokIdent:
			return isWord(p.peek(selectOffset+2), "QUALIFY")
		}
		return false
	case kwAND, kwOR, kwIS, kwIN, kwLIKE, kwILIKE, kwRLIKE, kwBETWEEN:
		if !inParens {
			return true
		}
		switch p.peek(selectOffset + 2).Kind {
		case tokEOF, tokSemicolon, tokRParen, tokComma, tokOperatorPipe,
			kwAS, kwFROM, kwWHERE, kwGROUP, kwHAVING, kwWINDOW,
			kwORDER, kwCLUSTER, kwDISTRIBUTE, kwSORT, kwLIMIT, kwOFFSET,
			kwUNION, kwINTERSECT, kwEXCEPT, kwMINUS, kwLATERAL:
			return true
		}
		return false
	}
	return true
}

func (p *Parser) isFromQueryAt(fromOffset int) bool {
	next := p.peek(fromOffset + 1).Kind
	switch next {
	case tokEOF, tokSemicolon, tokRParen, tokRBracket, tokComma, tokOperatorPipe,
		tokDot, tokLBracket, tokDoubleColon, kwCOLLATE,
		tokPlus, tokMinus, tokAsterisk, tokSlash, tokPercent, kwDIV,
		tokConcatPipe, tokPipe, tokAmpersand, tokHat,
		tokShiftLeft, tokShiftRight, tokShiftRightU,
		tokEq, tokNeq, tokNeqBang, tokNullSafeEq, tokLT, tokLTE, tokGT, tokGTE,
		kwAND, kwOR, kwIS, kwIN, kwLIKE, kwILIKE, kwRLIKE, kwBETWEEN,
		tokArrow, tokFatArrow, kwAS:
		return false
	case tokLParen:
		inner := p.peek(fromOffset + 2).Kind
		if inner == tokRParen || (!inner.isIdentifier() && inner != tokLParen && inner != kwIDENTIFIER && inner != kwLATERAL) {
			return false
		}
		// Reject function calls like `(from(x) + 1)` or `(from(x), 1)` where an
		// expression operator immediately follows `from(...)`.
		closeOff := p.matchingCloseOffset(fromOffset + 1)
		if closeOff < 0 {
			return false
		}
		switch p.peek(closeOff + 1).Kind {
		case tokDot, tokLBracket, tokDoubleColon, kwCOLLATE,
			tokPlus, tokMinus, tokAsterisk, tokSlash, tokPercent, kwDIV,
			tokConcatPipe, tokPipe, tokAmpersand, tokHat,
			tokShiftLeft, tokShiftRight, tokShiftRightU,
			tokEq, tokNeq, tokNeqBang, tokNullSafeEq, tokLT, tokLTE, tokGT, tokGTE,
			kwAND, kwOR, kwIS, kwIN, kwLIKE, kwILIKE, kwRLIKE, kwBETWEEN,
			tokArrow, tokFatArrow:
			return false
		}
		return true
	}
	return true
}

// isTableQueryAt reports whether `TABLE` at tableOffset introduces a
// `TABLE identifierReference` queryPrimary (followed by end-of-query or a
// queryTerm/queryOrganization continuation) rather than a column/table
// identifier named `table` in a row constructor `(table a, b)` or joined
// relation `(table a JOIN t2)`.
func (p *Parser) isTableQueryAt(tableOffset int) bool {
	pos := tableOffset + 1
	if p.peek(pos).Kind == kwIDENTIFIER && p.peek(pos+1).Kind == tokLParen {
		closeOff := p.matchingCloseOffset(pos + 1)
		if closeOff < 0 || p.peek(closeOff).Kind != tokRParen || closeOff == pos+2 {
			return false
		}
		pos = closeOff + 1
	} else {
		if !p.peek(pos).Kind.isIdentifier() {
			return false
		}
		pos++
		for p.peek(pos).Kind == tokDot && p.peek(pos+1).Kind.isIdentifier() {
			pos += 2
		}
	}

	switch p.peek(pos).Kind {
	case tokEOF, tokSemicolon, tokRParen, tokOperatorPipe:
		return true
	case kwUNION, kwINTERSECT, kwEXCEPT, kwMINUS:
		switch p.peek(pos + 1).Kind {
		case kwALL, kwDISTINCT, kwSELECT, kwFROM, kwTABLE, kwVALUES, kwMAP, kwREDUCE, tokLParen:
			return true
		}
		return false
	case kwLIMIT, kwOFFSET:
		switch p.peek(pos + 1).Kind {
		case tokEOF, tokSemicolon, tokRParen, tokRBracket, tokComma, tokOperatorPipe,
			tokDot, tokLBracket, tokDoubleColon,
			tokAsterisk, tokSlash, tokPercent, kwDIV,
			tokConcatPipe, tokPipe, tokAmpersand, tokHat,
			tokShiftLeft, tokShiftRight, tokShiftRightU,
			tokEq, tokNeq, tokNeqBang, tokNullSafeEq, tokLT, tokLTE, tokGT, tokGTE,
			kwAND, kwOR, kwIS, kwIN, kwLIKE, kwILIKE, kwRLIKE, kwBETWEEN,
			tokArrow, tokFatArrow:
			return false
		}
		return true
	case kwORDER, kwCLUSTER, kwDISTRIBUTE, kwSORT:
		return p.peek(pos+1).Kind == kwBY
	case kwWINDOW:
		return p.peek(pos+1).Kind.isIdentifier() && p.peek(pos+2).Kind == kwAS
	default:
		return false
	}
}

func (p *Parser) isIdentRefSpan(startOffset, endOffset int) bool {
	if startOffset >= endOffset {
		return false
	}
	if !p.peek(startOffset).Kind.isIdentifier() {
		return false
	}
	for j := startOffset + 1; j < endOffset; j += 2 {
		if p.peek(j).Kind != tokDot || j+1 >= endOffset || !p.peek(j+1).Kind.isIdentifier() {
			return false
		}
	}
	return true
}

func (p *Parser) parseFunctionTableSubqueryArgument() error {
	if _, err := p.expect(kwTABLE); err != nil {
		return err
	}
	if _, err := p.expect(tokLParen); err != nil {
		return err
	}
	if p.startsWithQueryKeywordAfterParens() {
		cp := p.save()
		if err := p.parseQuery(); err == nil && p.cur().Kind == tokRParen {
			p.advance()
			if p.isTableArgPartitioningStart() {
				return p.parseTableArgumentPartitioning()
			}
			return nil
		}
		p.restore(cp)
	}
	if _, err := p.parseIdentifierReference(); err != nil {
		return err
	}
	if _, err := p.expect(tokRParen); err != nil {
		return err
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
		if p.cur().Kind == tokLParen && p.isParenthesizedPartitionList(true) {
			if err := p.parseParenthesizedExprOrRow(false); err != nil {
				return err
			}
		} else {
			if err := p.parseExpression(); err != nil {
				return err
			}
			// Spark v4.0.5 rejects unparenthesized multi-expression PARTITION BY lists
			// (invalidMultiPartitionExpression); only a subsequent named or TABLE(...)
			// argument may follow a comma after an unparenthesized partition expression.
			if p.cur().Kind == tokComma && !p.startsNextFunctionTableArgAfterComma() {
				return p.errorf("multiple PARTITION BY expressions in TABLE argument must be parenthesized")
			}
		}
	}
	if (p.cur().Kind == kwORDER || p.cur().Kind == kwSORT) && p.peek(1).Kind == kwBY {
		p.advance()
		p.advance()
		if p.cur().Kind == tokLParen && p.isParenthesizedPartitionList(false) {
			if err := p.parseParenthesizedSortOrRow(); err != nil {
				return err
			}
		} else {
			if err := p.parseSortItem(); err != nil {
				return err
			}
			// Spark v4.0.5 rejects unparenthesized multi-item ORDER/SORT BY lists
			// (invalidMultiSortItem); only a subsequent named or TABLE(...) argument
			// may follow a comma after an unparenthesized sort item.
			if p.cur().Kind == tokComma && !p.startsNextFunctionTableArgAfterComma() {
				return p.errorf("multiple ORDER BY items in TABLE argument must be parenthesized")
			}
		}
	}
	return nil
}

// startsNextFunctionTableArgAfterComma reports whether the token(s) after a
// comma in an unparenthesized tableArgumentPartitioning list unambiguously
// begin a subsequent functionTableArgument (a named argument `k => ...` or a
// `TABLE ...` argument) rather than another partition/sort expression.
func (p *Parser) startsNextFunctionTableArgAfterComma() bool {
	if p.peek(1).Kind.isIdentifier() && p.peek(2).Kind == tokFatArrow {
		return true
	}
	return p.isTableSubqueryArgAt(1)
}

// isParenthesizedPartitionList reports whether a '(' after PARTITION/DISTRIBUTE BY
// or ORDER/SORT BY is a parenthesized list (terminated at its matching ')' by
// ')', ',', or — for PARTITION BY — ORDER/SORT BY) rather than the start of a
// parenthesized expression with a continuation like `(x) + 1` or `(x) DESC`.
func (p *Parser) isParenthesizedPartitionList(allowOrderBy bool) bool {
	if p.cur().Kind != tokLParen {
		return false
	}
	closeOff := p.matchingCloseOffset(0)
	if closeOff < 0 {
		return false
	}
	after := p.peek(closeOff + 1).Kind
	if after == tokRParen || after == tokComma {
		return true
	}
	if allowOrderBy && (after == kwORDER || after == kwSORT) && p.peek(closeOff+2).Kind == kwBY {
		return true
	}
	return false
}

// parseTableAlias parses: (AS? strictIdentifier identifierList?)?
func (p *Parser) parseTableAlias() error {
	if p.cur().Kind == kwAS {
		next := p.peek(1).Kind
		switch next {
		case tokEOF, tokSemicolon, tokComma, tokRParen, tokOperatorPipe,
			kwWHERE, kwGROUP, kwHAVING, kwWINDOW, kwORDER, kwCLUSTER, kwDISTRIBUTE,
			kwSORT, kwLIMIT, kwOFFSET, kwUNION, kwINTERSECT, kwEXCEPT, kwMINUS,
			kwJOIN, kwINNER, kwCROSS, kwLEFT, kwRIGHT, kwFULL, kwSEMI, kwANTI, kwNATURAL,
			kwON, kwUSING, kwLATERAL, kwPIVOT, kwUNPIVOT:
			// `AS` itself as an implicit table alias (e.g. `FROM t AS`, `FROM t AS WHERE a = 1`).
			p.advance()
			return nil
		}
		p.advance()
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

// parseDMLTableAlias parses: (AS? strictIdentifier)?
// Spark v4.0.5 getTableAliasWithoutColumnAlias (used by DELETE, UPDATE, and
// MERGE) forbids column alias lists on DML table aliases.
func (p *Parser) parseDMLTableAlias() error {
	if p.match(kwAS) {
		if _, err := p.parseStrictIdentifier(); err != nil {
			return err
		}
		if p.cur().Kind == tokLParen {
			return p.errorf("column aliases are not allowed in DML table alias")
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
		return p.errorf("column aliases are not allowed in DML table alias")
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
	// If keyword is at end of relation (EOF, ';', ',', ')', '|>'), it is a valid strict alias.
	if next == tokEOF || next == tokSemicolon || next == tokComma || next == tokRParen || next == tokOperatorPipe {
		return true
	}
	switch k {
	case kwGROUP, kwORDER, kwCLUSTER, kwDISTRIBUTE, kwSORT:
		return next != kwBY
	case kwWINDOW:
		return !(next.isIdentifier() && p.peek(2).Kind == kwAS)
	case kwPIVOT, kwTABLESAMPLE:
		return next != tokLParen
	case kwUNPIVOT:
		return next != tokLParen && next != kwINCLUDE && next != kwEXCLUDE
	case kwANTI, kwSEMI:
		return next != kwJOIN
	case kwWITH:
		return next != tokLParen && next != kwROLLUP && next != kwCUBE
	case kwFOR:
		return next != kwSYSTEM_TIME && next != kwSYSTEM_VERSION && next != kwVERSION && next != kwTIMESTAMP
	case kwWHERE, kwHAVING, kwLIMIT, kwOFFSET, kwON, kwSET,
		kwSELECT, kwINSERT, kwMAP, kwREDUCE, kwWHEN, kwTHEN, kwELSE, kwEND,
		kwVALUES, kwFROM, kwDO, kwUNTIL, kwINTO, kwAS, kwBY:
		if k == kwAS {
			return false
		}
		if p.inFromStatement && (k == kwSELECT || k == kwMAP || k == kwREDUCE) {
			return false
		}
		// Allow as implicit alias when followed by an unmistakable clause boundary.
		after := p.peek(2).Kind
		isClauseExprArg := func() bool {
			return !p.isClauseStartAt(2) && !isExprContinuationOp(after)
		}
		switch next {
		case kwGROUP, kwORDER, kwCLUSTER, kwDISTRIBUTE, kwSORT:
			return after == kwBY
		case kwWINDOW:
			return after.isIdentifier() && p.peek(3).Kind == kwAS
		case kwLATERAL:
			return after == kwVIEW
		case kwWHERE, kwHAVING, kwLIMIT, kwOFFSET, kwON:
			return isClauseExprArg()
		case kwSET:
			return after.isIdentifier() && (p.peek(3).Kind == tokEq || p.peek(3).Kind == tokDot)
		case kwUSING:
			return after == tokLParen || isClauseExprArg()
		case kwJOIN:
			return isClauseExprArg()
		case kwINNER, kwCROSS, kwNATURAL, kwSEMI, kwANTI:
			return after == kwJOIN
		case kwLEFT, kwRIGHT, kwFULL:
			return after == kwJOIN || after == kwOUTER || after == kwSEMI || after == kwANTI
		case kwPIVOT, kwTABLESAMPLE:
			return after == tokLParen
		case kwUNPIVOT:
			return after == tokLParen || after == kwINCLUDE || after == kwEXCLUDE
		case kwUNION, kwINTERSECT, kwEXCEPT, kwMINUS:
			switch after {
			case kwALL, kwDISTINCT, kwSELECT, kwFROM, kwTABLE, kwVALUES, kwMAP, kwREDUCE, tokLParen:
				return true
			}
		case tokIdent:
			return isWord(p.peek(1), "QUALIFY") && isClauseExprArg()
		}
		return false
	}
	return true
}

func (p *Parser) isClauseStartAt(offset int) bool {
	k := p.peek(offset).Kind
	n := p.peek(offset + 1).Kind
	switch k {
	case tokEOF, tokSemicolon, tokRParen, tokComma, tokOperatorPipe:
		return true
	case kwGROUP, kwORDER, kwCLUSTER, kwDISTRIBUTE, kwSORT:
		return n == kwBY
	case kwWINDOW:
		return n.isIdentifier() && p.peek(offset+2).Kind == kwAS
	case kwLATERAL:
		return n == kwVIEW
	case kwWITH:
		return n == kwROLLUP || n == kwCUBE
	case kwGROUPING:
		return n == kwSETS && p.peek(offset+2).Kind == tokLParen
	case kwUNION, kwINTERSECT, kwEXCEPT, kwMINUS:
		switch n {
		case kwALL, kwDISTINCT, kwSELECT, kwFROM, kwTABLE, kwVALUES, kwMAP, kwREDUCE, tokLParen:
			return true
		}
		return false
	case kwWHERE, kwHAVING, kwLIMIT, kwOFFSET, kwFROM:
		return n != tokEOF && n != tokSemicolon && n != tokRParen && n != tokComma && !isExprContinuationOp(n)
	case tokIdent:
		return isWord(p.peek(offset), "QUALIFY") && n != tokEOF && n != tokSemicolon && n != tokRParen && n != tokComma && !isExprContinuationOp(n)
	}
	return false
}

func (p *Parser) parseSample() error {
	if _, err := p.expect(kwTABLESAMPLE); err != nil {
		return err
	}
	if _, err := p.expect(tokLParen); err != nil {
		return err
	}
	if err := p.parseSampleMethod(); err != nil {
		return err
	}
	if _, err := p.expect(tokRParen); err != nil {
		return err
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
	// BUCKET int OUT OF int (Spark v4.0.5 rejects ON col/func, y <= 0, and x > y)
	if p.cur().Kind == kwBUCKET && p.peek(1).Kind == tokInt {
		p.advance()
		numTok := p.advance()
		if _, err := p.expect(kwOUT); err != nil {
			return err
		}
		if _, err := p.expect(kwOF); err != nil {
			return err
		}
		denTok, err := p.expect(tokInt)
		if err != nil {
			return err
		}
		num, err1 := strconv.ParseInt(numTok.Text, 10, 64)
		den, err2 := strconv.ParseInt(denTok.Text, 10, 64)
		if err1 != nil || err2 != nil || den <= 0 || num < 0 || num > den {
			return p.errorf("invalid TABLESAMPLE BUCKET %s OUT OF %s", numTok.Text, denTok.Text)
		}
		if p.cur().Kind == kwON {
			return p.errorf("TABLESAMPLE BUCKET ... ON is not supported")
		}
		return nil
	}
	// MINUS? (INT | DECIMAL) PERCENT (Spark requires percentage in [0, 100])
	cp := p.save()
	neg := p.match(tokMinus)
	numTok := p.cur()
	if p.match(tokInt, tokDecimal) && p.match(kwPERCENT) {
		val, err := strconv.ParseFloat(numTok.Text, 64)
		if neg {
			val = -val
		}
		if err != nil || val < 0 || val > 100 {
			return p.errorf("TABLESAMPLE percentage %v must be in [0, 100]", val)
		}
		return nil
	}
	p.restore(cp)

	// expression ROWS (Spark v4.0.5 rejects bare byte-length / unannotated expressions)
	if err := p.parseExpression(); err != nil {
		return err
	}
	if _, err := p.expect(kwROWS); err != nil {
		return err
	}
	return nil
}

func (p *Parser) parseLateralView() error {
	if _, err := p.expect(kwLATERAL); err != nil {
		return err
	}
	if _, err := p.expect(kwVIEW); err != nil {
		return err
	}
	if p.cur().Kind == kwOUTER && p.peek(1).Kind != tokLParen && p.peek(1).Kind != tokDot {
		p.advance()
	}
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

func (p *Parser) parseAggregationClause(inPipeAggregate bool) error {
	if _, err := p.expect(kwGROUP); err != nil {
		return err
	}
	if _, err := p.expect(kwBY); err != nil {
		return err
	}
	if inPipeAggregate && p.cur().Kind == kwALL {
		next := p.peek(1).Kind
		if next == tokEOF || next == tokSemicolon || next == tokRParen || next == tokOperatorPipe || p.isQueryOrgStart() {
			return p.errorf("GROUP BY ALL is not supported in AGGREGATE pipe operator")
		}
	}
	// Spark v4.0.5 parses regular GROUP BY via groupByClause (no aliases allowed),
	// whereas pipe `|> AGGREGATE ... GROUP BY` allows named expressions (aliases)
	// while forbidding GROUP BY ALL, ROLLUP, CUBE, and GROUPING SETS.
	hasGroupingAnalytics := false
	for {
		if p.isGroupingAnalyticsStart() {
			if inPipeAggregate {
				return p.errorf("ROLLUP, CUBE, and GROUPING SETS are not supported in AGGREGATE pipe operator")
			}
			if _, err := p.parseGroupingAnalytics(); err != nil {
				return err
			}
			hasGroupingAnalytics = true
		} else if inPipeAggregate {
			if err := p.parseNamedExpression(); err != nil {
				return err
			}
		} else {
			if err := p.parseExpression(); err != nil {
				return err
			}
		}
		if !p.match(tokComma) {
			break
		}
	}
	if hasGroupingAnalytics {
		return nil
	}
	// Optional WITH ROLLUP | WITH CUBE | GROUPING SETS (...).
	if p.cur().Kind == kwWITH && (p.peek(1).Kind == kwROLLUP || p.peek(1).Kind == kwCUBE) {
		if inPipeAggregate {
			return p.errorf("WITH ROLLUP/CUBE is not supported in AGGREGATE pipe operator")
		}
		p.advance()
		p.advance()
	} else if p.cur().Kind == kwGROUPING && p.peek(1).Kind == kwSETS {
		if inPipeAggregate {
			return p.errorf("GROUPING SETS is not supported in AGGREGATE pipe operator")
		}
		p.advance()
		p.advance()
		if _, err := p.expect(tokLParen); err != nil {
			return err
		}
		if _, err := p.parseGroupingSet(); err != nil {
			return err
		}
		for p.match(tokComma) {
			if _, err := p.parseGroupingSet(); err != nil {
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
	if p.cur().Kind == kwGROUPING && p.peek(1).Kind == kwSETS && p.peek(2).Kind == tokLParen {
		return true
	}
	if (p.cur().Kind != kwROLLUP && p.cur().Kind != kwCUBE) || p.peek(1).Kind != tokLParen {
		return false
	}
	// Distinguish ROLLUP(...)/CUBE(...) grouping analytics from ordinary
	// function calls (e.g. `ROLLUP()`, `ROLLUP(k => x)`, `ROLLUP(x) + 1`,
	// `ROLLUP(x) AS r`, `ROLLUP(x) OVER ()`) before recursively parsing arguments.
	inner := p.peek(2).Kind
	if inner == tokRParen || inner == tokAsterisk {
		return false
	}
	if p.isFunctionSetQuantifierAt(2) {
		return false
	}
	closeIdx := p.matchingCloseOffset(1)
	if closeIdx < 0 {
		return false
	}
	for i := 2; i < closeIdx; i++ {
		k := p.peek(i).Kind
		if k == tokLParen || k == tokLBracket {
			c := p.matchingCloseOffset(i)
			if c < 0 {
				return false
			}
			i = c
			continue
		}
		if k == tokFatArrow {
			return false
		}
	}
	after := p.peek(closeIdx + 1).Kind
	after2 := p.peek(closeIdx + 2).Kind
	switch after {
	case tokEOF, tokSemicolon, tokComma, tokRParen, tokOperatorPipe:
		return true
	case kwWITH:
		return after2 == kwROLLUP || after2 == kwCUBE
	case kwGROUPING:
		return after2 == kwSETS && p.peek(closeIdx+3).Kind == tokLParen
	case kwORDER, kwCLUSTER, kwDISTRIBUTE, kwSORT:
		return after2 == kwBY
	case kwHAVING, kwWINDOW, kwLIMIT, kwOFFSET, kwUNION, kwINTERSECT, kwEXCEPT, kwMINUS:
		return after2 != tokEOF && after2 != tokSemicolon && after2 != tokComma && after2 != tokRParen
	case tokIdent:
		return isWord(p.peek(closeIdx+1), "QUALIFY") && after2 != tokEOF && after2 != tokSemicolon && after2 != tokComma && after2 != tokRParen
	}
	return false
}

func (p *Parser) parseGroupingAnalytics() (bool, error) {
	if err := p.enterNesting(); err != nil {
		return false, err
	}
	defer p.leaveNesting()

	if p.match(kwROLLUP, kwCUBE) {
		if _, err := p.expect(tokLParen); err != nil {
			return false, err
		}
		hasEmpty, err := p.parseGroupingSet()
		if err != nil {
			return false, err
		}
		if hasEmpty {
			return true, p.errorf("empty set () in ROLLUP/CUBE is not supported")
		}
		for p.match(tokComma) {
			empty, err := p.parseGroupingSet()
			if err != nil {
				return false, err
			}
			if empty {
				return true, p.errorf("empty set () in ROLLUP/CUBE is not supported")
			}
		}
		_, err = p.expect(tokRParen)
		return false, err
	}
	if _, err := p.expect(kwGROUPING); err != nil {
		return true, err
	}
	if _, err := p.expect(kwSETS); err != nil {
		return true, err
	}
	if _, err := p.expect(tokLParen); err != nil {
		return true, err
	}
	if err := p.parseGroupingElement(); err != nil {
		return true, err
	}
	for p.match(tokComma) {
		if err := p.parseGroupingElement(); err != nil {
			return true, err
		}
	}
	_, err := p.expect(tokRParen)
	return true, err
}

func (p *Parser) parseGroupingElement() error {
	if p.isGroupingAnalyticsStart() {
		_, err := p.parseGroupingAnalytics()
		return err
	}
	_, err := p.parseGroupingSet()
	return err
}

func (p *Parser) parseGroupingSet() (bool, error) {
	if p.cur().Kind == tokLParen && p.isParenthesizedGroupingSet() {
		if p.peek(1).Kind == tokRParen {
			p.advance()
			p.advance()
			return true, nil
		}
		return false, p.parseParenthesizedExprOrRow(true)
	}
	return false, p.parseExpression()
}

// parseParenthesizedExprOrRow parses `()` (when allowEmpty), `(query)`,
// `(expr)`, `(expr, expr+)`, or a multi-element row constructor
// `(namedExpr, namedExpr+)` in a single pass without ever reparsing the first
// element.
func (p *Parser) parseParenthesizedExprOrRow(allowEmpty bool) error {
	if _, err := p.expect(tokLParen); err != nil {
		return err
	}
	if allowEmpty && p.match(tokRParen) {
		return nil
	}
	if p.startsWithCompleteParenQuery() {
		if err := p.parseQuery(); err != nil {
			return err
		}
		_, err := p.expect(tokRParen)
		return err
	}
	if err := p.parseExpression(); err != nil {
		return err
	}
	if p.match(tokRParen) {
		return nil
	}
	if p.cur().Kind != tokComma {
		before := p.pos
		if err := p.parseOptionalNamedExprAlias(); err != nil {
			return err
		}
		if p.pos == before {
			return p.errorf("expected ')' or ',', got %q", p.cur().Text)
		}
		if !p.match(tokComma) {
			return p.errorf("expected ',' after aliased row constructor element")
		}
	} else {
		p.advance() // ','
	}
	for {
		if err := p.parseNamedExpression(); err != nil {
			return err
		}
		if !p.match(tokComma) {
			break
		}
	}
	_, err := p.expect(tokRParen)
	return err
}

// parseParenthesizedSortOrRow parses `(sortItem, sortItem*)` or a row
// constructor `(namedExpr, namedExpr+)` inside `ORDER BY (...)` in a single
// pass without reparsing the first element.
func (p *Parser) parseParenthesizedSortOrRow() error {
	if _, err := p.expect(tokLParen); err != nil {
		return err
	}
	if p.startsWithCompleteParenQuery() {
		if err := p.parseQuery(); err != nil {
			return err
		}
		_, err := p.expect(tokRParen)
		return err
	}
	if err := p.parseExpression(); err != nil {
		return err
	}
	hasSortDir := false
	hasAlias := false
	parseSortDir := func() (bool, error) {
		matched := false
		if p.match(kwASC, kwDESC) {
			matched = true
		}
		if p.match(kwNULLS) {
			if !p.match(kwFIRST, kwLAST) {
				return false, p.errorf("expected FIRST or LAST after NULLS")
			}
			matched = true
		}
		return matched, nil
	}
	var err error
	if hasSortDir, err = parseSortDir(); err != nil {
		return err
	}
	if p.match(tokRParen) {
		return nil
	}
	if !hasSortDir && p.cur().Kind != tokComma {
		before := p.pos
		if err := p.parseOptionalNamedExprAlias(); err != nil {
			return err
		}
		if p.pos == before {
			return p.errorf("expected ')' or ',', got %q", p.cur().Text)
		}
		hasAlias = true
		if !p.match(tokComma) {
			return p.errorf("expected ',' after aliased row constructor element")
		}
	} else if !p.match(tokComma) {
		return p.errorf("expected ')' or ',', got %q", p.cur().Text)
	}
	for {
		if err := p.parseExpression(); err != nil {
			return err
		}
		if !hasAlias {
			sd, err := parseSortDir()
			if err != nil {
				return err
			}
			if sd {
				hasSortDir = true
			} else if !hasSortDir && p.cur().Kind != tokComma && p.cur().Kind != tokRParen {
				before := p.pos
				if err := p.parseOptionalNamedExprAlias(); err != nil {
					return err
				}
				if p.pos > before {
					hasAlias = true
				}
			}
		} else if p.cur().Kind != tokComma && p.cur().Kind != tokRParen {
			if err := p.parseOptionalNamedExprAlias(); err != nil {
				return err
			}
		}
		if !p.match(tokComma) {
			break
		}
	}
	_, err = p.expect(tokRParen)
	return err
}

// isParenthesizedGroupingSet reports whether '(' at p.cur() closes directly
// before ',' or ')' (as a parenthesized groupingSet `(e1, e2)` or `()`) rather
// than beginning a parenthesized expression with an operator continuation such
// as `ROLLUP(((<query>)) + 1)`.
func (p *Parser) isParenthesizedGroupingSet() bool {
	if p.cur().Kind != tokLParen {
		return false
	}
	closeOff := p.matchingCloseOffset(0)
	if closeOff < 0 {
		return false
	}
	after := p.peek(closeOff + 1).Kind
	return after == tokRParen || after == tokComma
}

func (p *Parser) parseWindowClause() error {
	if _, err := p.expect(kwWINDOW); err != nil {
		return err
	}
	defined := make(map[string]struct{})
	var refs []string
	for {
		name, err := p.parseErrorCapturingIdentifier()
		if err != nil {
			return err
		}
		norm := normalizeIdentName(name)
		if _, dup := defined[norm]; dup {
			return p.errorf("repetitive window definition %q", name)
		}
		defined[norm] = struct{}{}
		if _, err := p.expect(kwAS); err != nil {
			return err
		}
		// Spark v4.0.5 AstBuilder.withWindows requires parenthesized windowSpec
		// (`WINDOW w AS (...)`) and rejects bare window references (`WINDOW w AS w2`).
		if p.cur().Kind != tokLParen {
			return p.errorf("window specification in WINDOW clause must be parenthesized")
		}
		ref, err := p.parseWindowSpecRef()
		if err != nil {
			return err
		}
		if ref != "" {
			refs = append(refs, ref)
		}
		if !p.match(tokComma) {
			break
		}
	}
	for _, ref := range refs {
		if _, ok := defined[normalizeIdentName(ref)]; !ok {
			return p.errorf("unresolved window reference %q", ref)
		}
	}
	return nil
}

func (p *Parser) parseWindowSpec() error {
	_, err := p.parseWindowSpecRef()
	return err
}

func (p *Parser) parseWindowSpecRef() (string, error) {
	if !p.match(tokLParen) {
		_, err := p.parseErrorCapturingIdentifier()
		return "", err
	}
	if p.match(tokRParen) {
		return "", nil
	}
	// Could be windowRef: (name) or windowDef: (CLUSTER BY ... | (PARTITION|DISTRIBUTE BY ...)? (ORDER|SORT BY ...)? windowFrame?)
	cp := p.save()
	if ref, err := p.parseErrorCapturingIdentifier(); err == nil && p.cur().Kind == tokRParen {
		p.advance()
		return ref, nil
	}
	p.restore(cp)

	if p.cur().Kind == kwCLUSTER && p.peek(1).Kind == kwBY {
		p.advance()
		p.advance()
		if err := p.parseExpressionSeq(); err != nil {
			return "", err
		}
	} else {
		if (p.cur().Kind == kwPARTITION || p.cur().Kind == kwDISTRIBUTE) && p.peek(1).Kind == kwBY {
			p.advance()
			p.advance()
			if err := p.parseExpressionSeq(); err != nil {
				return "", err
			}
		}
		if (p.cur().Kind == kwORDER || p.cur().Kind == kwSORT) && p.peek(1).Kind == kwBY {
			p.advance()
			p.advance()
			if err := p.parseSortItemList(); err != nil {
				return "", err
			}
		}
	}
	if p.cur().Kind == kwRANGE || p.cur().Kind == kwROWS {
		if err := p.parseWindowFrame(); err != nil {
			return "", err
		}
	}
	_, err := p.expect(tokRParen)
	return "", err
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

func (p *Parser) parseQueryOrganization(inPipe bool) error {
	hasOrder := false
	hasCluster := false
	hasDistributeOrSort := false
	if p.cur().Kind == kwORDER && p.peek(1).Kind == kwBY {
		hasOrder = true
		p.advance()
		p.advance()
		if err := p.parseSortItemList(); err != nil {
			return err
		}
	}
	if p.cur().Kind == kwCLUSTER && p.peek(1).Kind == kwBY {
		if hasOrder {
			return p.errorf("CLUSTER BY cannot be combined with ORDER BY")
		}
		hasCluster = true
		p.advance()
		p.advance()
		if err := p.parseExpressionSeq(); err != nil {
			return err
		}
	}
	if p.cur().Kind == kwDISTRIBUTE && p.peek(1).Kind == kwBY {
		if hasOrder || hasCluster {
			return p.errorf("DISTRIBUTE BY cannot be combined with ORDER BY or CLUSTER BY")
		}
		hasDistributeOrSort = true
		p.advance()
		p.advance()
		if err := p.parseExpressionSeq(); err != nil {
			return err
		}
	}
	if p.cur().Kind == kwSORT && p.peek(1).Kind == kwBY {
		if hasOrder || hasCluster {
			return p.errorf("SORT BY cannot be combined with ORDER BY or CLUSTER BY")
		}
		hasDistributeOrSort = true
		p.advance()
		p.advance()
		if err := p.parseSortItemList(); err != nil {
			return err
		}
	}
	if p.cur().Kind == kwWINDOW {
		if inPipe {
			return p.errorf("WINDOW clause is not supported directly after pipe operator")
		}
		if err := p.parseWindowClause(); err != nil {
			return err
		}
	}
	hasOrdering := hasOrder || hasCluster || hasDistributeOrSort
	if p.match(kwLIMIT) {
		if inPipe && hasOrdering {
			return p.errorf("cannot combine ORDER/CLUSTER/DISTRIBUTE/SORT BY with LIMIT in the same pipe operator")
		}
		if p.isStandaloneLimitAll() {
			p.advance()
		} else {
			if err := p.parseExpression(); err != nil {
				return err
			}
		}
	}
	if p.match(kwOFFSET) {
		if inPipe && hasOrdering {
			return p.errorf("cannot combine ORDER/CLUSTER/DISTRIBUTE/SORT BY with OFFSET in the same pipe operator")
		}
		if err := p.parseExpression(); err != nil {
			return err
		}
	}
	return nil
}

// isStandaloneLimitAll reports whether `ALL` after `LIMIT` is the standalone
// `LIMIT ALL` form rather than the start of an expression (such as `all + 1`,
// `all.x`, or `all(1)`).
func (p *Parser) isStandaloneLimitAll() bool {
	if p.cur().Kind != kwALL {
		return false
	}
	switch p.peek(1).Kind {
	case tokEOF, tokSemicolon, tokRParen, tokComma, tokOperatorPipe,
		kwOFFSET, kwUNION, kwINTERSECT, kwEXCEPT, kwMINUS,
		kwINSERT, kwEND, kwELSE, kwELSEIF, kwWHEN, kwTHEN:
		return true
	}
	return false
}
