package parser

func (p *Parser) parseDMLWithCTEs() (Statement, error) {
	start := p.cur().Start
	if err := p.parseCTEs(); err != nil {
		return nil, err
	}
	stmt, err := p.parseDMLStatementNoWith()
	if err != nil {
		return nil, err
	}
	if ps, ok := stmt.(*ParsedStmt); ok {
		ps.Loc.Start = start
	}
	return stmt, nil
}

func (p *Parser) parseDMLStatementNoWith() (Statement, error) {
	switch p.cur().Kind {
	case kwINSERT:
		return p.parseSingleInsertQuery()
	case kwFROM:
		return p.parseMultiInsertQuery()
	case kwDELETE:
		return p.parseDeleteFromTable()
	case kwUPDATE:
		return p.parseUpdateTable()
	case kwMERGE:
		return p.parseMergeIntoTable()
	default:
		return nil, p.errorf("expected DML statement, got %q", p.cur().Text)
	}
}

func (p *Parser) parseSingleInsertQuery() (Statement, error) {
	start := p.cur().Start
	cp := p.save()
	ambiguous := p.hasAmbiguousInsertTableTarget()
	target, err := p.parseInsertIntoWithTable(true)
	if err == nil {
		if err = p.parseQuery(); err == nil && p.isStatementEnd() {
			return &ParsedStmt{Kind: StmtInsert, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
		}
	}
	if !ambiguous {
		if err != nil {
			return nil, err
		}
		return &ParsedStmt{Kind: StmtInsert, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
	}
	p.restore(cp)
	target, err = p.parseInsertIntoWithTable(false)
	if err != nil {
		return nil, err
	}
	if err := p.parseQuery(); err != nil {
		return nil, err
	}
	return &ParsedStmt{Kind: StmtInsert, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
}

func (p *Parser) parseMultiInsertQuery() (Statement, error) {
	start := p.cur().Start
	if err := p.parseFromClause(); err != nil {
		return nil, err
	}
	count := 0
	for p.cur().Kind == kwINSERT {
		cp := p.save()
		ambiguous := p.hasAmbiguousInsertTableTarget()
		_, err := p.parseInsertIntoWithTable(true)
		if err == nil {
			err = p.parseFromStatementBody()
		}
		if (err != nil || (!p.isStatementEnd() && p.cur().Kind != kwINSERT)) && ambiguous {
			p.restore(cp)
			if _, err = p.parseInsertIntoWithTable(false); err != nil {
				return nil, err
			}
			if err = p.parseFromStatementBody(); err != nil {
				return nil, err
			}
		} else if err != nil {
			return nil, err
		}
		count++
	}
	if count == 0 {
		return nil, p.errorf("expected INSERT in multi-insert query")
	}
	return &ParsedStmt{Kind: StmtMultiInsert, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
}

// hasAmbiguousInsertTableTarget reports whether `INSERT (INTO|OVERWRITE) TABLE <kw>`
// is followed by a keyword (`SELECT`, `FROM`, `WITH`, `VALUES`, `TABLE`, `MAP`, `REDUCE`)
// where `TABLE` could be either the optional keyword before a keyword-named target
// (`INSERT INTO TABLE select SELECT 1`) or the target table itself
// (`INSERT INTO table SELECT select FROM t`, `INSERT INTO table TABLE t`).
func (p *Parser) hasAmbiguousInsertTableTarget() bool {
	if p.cur().Kind != kwINSERT {
		return false
	}
	if p.peek(1).Kind != kwINTO && p.peek(1).Kind != kwOVERWRITE {
		return false
	}
	if p.peek(2).Kind != kwTABLE {
		return false
	}
	switch p.peek(3).Kind {
	case kwSELECT, kwFROM, kwWITH, kwVALUES, kwTABLE, kwMAP, kwREDUCE:
		return true
	}
	return false
}

// parseInsertInto implements all 5 alternatives of insertInto in SqlBaseParser.g4:
//  1. INSERT OVERWRITE TABLE? identifierReference optionsClause? (partitionSpec (IF NOT EXISTS)?)? ((BY NAME) | identifierList)?
//  2. INSERT INTO TABLE? identifierReference optionsClause? partitionSpec? (IF NOT EXISTS)? ((BY NAME) | identifierList)?
//  3. INSERT INTO TABLE? identifierReference optionsClause? REPLACE whereClause
//  4. INSERT OVERWRITE LOCAL? DIRECTORY stringLit rowFormat? createFileFormat?
//  5. INSERT OVERWRITE LOCAL? DIRECTORY stringLit? tableProvider (OPTIONS propertyList)?
func (p *Parser) parseInsertIntoWithTable(allowOptionalTable bool) (string, error) {
	if _, err := p.expect(kwINSERT); err != nil {
		return "", err
	}
	if p.match(kwOVERWRITE) {
		if p.cur().Kind == kwLOCAL || p.cur().Kind == kwDIRECTORY {
			cpDir := p.save()
			p.match(kwLOCAL)
			if p.match(kwDIRECTORY) {
				if p.cur().Kind == kwUSING {
					if err := p.parseTableProvider(); err == nil {
						if p.match(kwOPTIONS) {
							if err := p.parsePropertyList(); err != nil {
								return "", err
							}
						}
						return "", nil
					}
				} else if p.isStringLit() {
					p.advance()
					if p.cur().Kind == kwUSING {
						if err := p.parseTableProvider(); err != nil {
							return "", err
						}
						if p.match(kwOPTIONS) {
							if err := p.parsePropertyList(); err != nil {
								return "", err
							}
						}
						return "", nil
					}
					if p.cur().Kind == kwROW && p.peek(1).Kind == kwFORMAT {
						if err := p.parseRowFormat(); err != nil {
							return "", err
						}
					}
					if p.cur().Kind == kwSTORED {
						if err := p.parseCreateFileFormat(); err != nil {
							return "", err
						}
					}
					return "", nil
				}
			}
			// Not a directory insert (e.g. `INSERT OVERWRITE DIRECTORY SELECT a` where DIRECTORY is a table name).
			p.restore(cpDir)
		}
		if allowOptionalTable && p.isOptionalInsertTableKeyword() {
			p.advance()
		}
		target, err := p.parseIdentifierReference()
		if err != nil {
			return "", err
		}
		if p.cur().Kind == kwWITH && p.peek(1).Kind == tokLParen {
			if err := p.parseOptionsClause(); err != nil {
				return "", err
			}
		}
		if p.cur().Kind == kwPARTITION {
			if err := p.parsePartitionSpec(); err != nil {
				return "", err
			}
			if p.cur().Kind == kwIF {
				p.advance()
				if !p.matchNot() {
					return "", p.errorf("expected NOT after IF")
				}
				if _, err := p.expect(kwEXISTS); err != nil {
					return "", err
				}
			}
		}
		if p.cur().Kind == kwBY && p.peek(1).Kind == kwNAME {
			p.advance()
			p.advance()
		} else if p.cur().Kind == tokLParen {
			cp := p.save()
			if err := p.parseIdentifierList(); err == nil && p.looksLikeQueryStart() {
				// column list before query
			} else {
				p.restore(cp)
			}
		}
		return target, nil
	}

	if _, err := p.expect(kwINTO); err != nil {
		return "", err
	}
	if allowOptionalTable && p.isOptionalInsertTableKeyword() {
		p.advance()
	}
	target, err := p.parseIdentifierReference()
	if err != nil {
		return "", err
	}
	if p.cur().Kind == kwWITH && p.peek(1).Kind == tokLParen {
		if err := p.parseOptionsClause(); err != nil {
			return "", err
		}
	}
	if p.match(kwREPLACE) {
		if err := p.parseWhereClause(); err != nil {
			return "", err
		}
		return target, nil
	}
	if p.cur().Kind == kwPARTITION {
		if err := p.parsePartitionSpec(); err != nil {
			return "", err
		}
	}
	if p.cur().Kind == kwIF {
		p.advance()
		if !p.matchNot() {
			return "", p.errorf("expected NOT after IF")
		}
		if _, err := p.expect(kwEXISTS); err != nil {
			return "", err
		}
	}
	if p.cur().Kind == kwBY && p.peek(1).Kind == kwNAME {
		p.advance()
		p.advance()
	} else if p.cur().Kind == tokLParen {
		cp := p.save()
		if err := p.parseIdentifierList(); err == nil && p.looksLikeQueryStart() {
			// column list before query
		} else {
			p.restore(cp)
		}
	}
	return target, nil
}

// isOptionalInsertTableKeyword reports whether `TABLE` after `INSERT INTO` or
// `INSERT OVERWRITE` acts as the optional `TABLE` keyword rather than a target
// table or namespace named `table` (e.g. `INSERT INTO table VALUES (1)`,
// `INSERT INTO table.t VALUES (1)`, `INSERT OVERWRITE table SELECT 1`).
func (p *Parser) isOptionalInsertTableKeyword() bool {
	if p.cur().Kind != kwTABLE {
		return false
	}
	next := p.peek(1).Kind
	switch next {
	case tokEOF, tokSemicolon, tokDot, tokLParen:
		return false
	case kwSELECT, kwFROM, kwWITH, kwVALUES:
		after := p.peek(2).Kind
		return after == tokDot || after == tokLParen ||
			after == kwSELECT || after == kwFROM || after == kwWITH ||
			after == kwVALUES || after == kwTABLE || after == kwMAP || after == kwREDUCE ||
			after == kwPARTITION || after == kwBY || after == kwREPLACE || after == kwIF
	case kwPARTITION:
		return p.peek(2).Kind != tokLParen
	case kwBY:
		return p.peek(2).Kind != kwNAME
	case kwREPLACE:
		return p.peek(2).Kind != kwWHERE
	case kwIF:
		return p.peek(2).Kind != kwNOT && p.peek(2).Kind != tokBang
	}
	return next.isIdentifier() || next == kwIDENTIFIER
}

func (p *Parser) parseDeleteFromTable() (Statement, error) {
	start := p.cur().Start
	if _, err := p.expect(kwDELETE); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}
	target, err := p.parseIdentifierReference()
	if err != nil {
		return nil, err
	}
	if err := p.parseTableAlias(); err != nil {
		return nil, err
	}
	if p.cur().Kind == kwWHERE {
		if err := p.parseWhereClause(); err != nil {
			return nil, err
		}
	}
	return &ParsedStmt{Kind: StmtDelete, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
}

func (p *Parser) parseUpdateTable() (Statement, error) {
	start := p.cur().Start
	if _, err := p.expect(kwUPDATE); err != nil {
		return nil, err
	}
	target, err := p.parseIdentifierReference()
	if err != nil {
		return nil, err
	}
	if err := p.parseTableAlias(); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwSET); err != nil {
		return nil, err
	}
	if err := p.parseAssignmentList(); err != nil {
		return nil, err
	}
	if p.cur().Kind == kwWHERE {
		if err := p.parseWhereClause(); err != nil {
			return nil, err
		}
	}
	return &ParsedStmt{Kind: StmtUpdate, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
}

func (p *Parser) parseAssignmentList() error {
	if err := p.parseAssignment(); err != nil {
		return err
	}
	for p.match(tokComma) {
		if err := p.parseAssignment(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parseAssignment() error {
	if _, err := p.parseMultipartIdentifier(); err != nil {
		return err
	}
	if _, err := p.expect(tokEq); err != nil {
		return err
	}
	return p.parseExpression()
}

func (p *Parser) parseMergeIntoTable() (Statement, error) {
	start := p.cur().Start
	if _, err := p.expect(kwMERGE); err != nil {
		return nil, err
	}
	if p.match(kwWITH) {
		if _, err := p.expect(kwSCHEMA); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwEVOLUTION); err != nil {
			return nil, err
		}
	}
	if _, err := p.expect(kwINTO); err != nil {
		return nil, err
	}
	target, err := p.parseIdentifierReference()
	if err != nil {
		return nil, err
	}
	if err := p.parseTableAlias(); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwUSING); err != nil {
		return nil, err
	}
	if p.match(tokLParen) {
		if err := p.parseQuery(); err != nil {
			return nil, err
		}
		if _, err := p.expect(tokRParen); err != nil {
			return nil, err
		}
	} else {
		if _, err := p.parseIdentifierReference(); err != nil {
			return nil, err
		}
	}
	if err := p.parseTableAlias(); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	if err := p.parseBooleanExpression(); err != nil {
		return nil, err
	}
	// matchedClause* notMatchedClause* notMatchedBySourceClause*
	for p.cur().Kind == kwWHEN && p.peek(1).Kind == kwMATCHED {
		p.advance()
		p.advance()
		if p.match(kwAND) {
			if err := p.parseBooleanExpression(); err != nil {
				return nil, err
			}
		}
		if _, err := p.expect(kwTHEN); err != nil {
			return nil, err
		}
		if err := p.parseMatchedAction(); err != nil {
			return nil, err
		}
	}
	for p.cur().Kind == kwWHEN && p.isNotAt(1) && p.peek(2).Kind == kwMATCHED &&
		!(p.peek(3).Kind == kwBY && p.peek(4).Kind == kwSOURCE) {
		p.advance()
		p.advance()
		p.advance()
		if p.cur().Kind == kwBY && p.peek(1).Kind == kwTARGET {
			p.advance()
			p.advance()
		}
		if p.match(kwAND) {
			if err := p.parseBooleanExpression(); err != nil {
				return nil, err
			}
		}
		if _, err := p.expect(kwTHEN); err != nil {
			return nil, err
		}
		if err := p.parseNotMatchedAction(); err != nil {
			return nil, err
		}
	}
	for p.cur().Kind == kwWHEN && p.isNotAt(1) && p.peek(2).Kind == kwMATCHED &&
		p.peek(3).Kind == kwBY && p.peek(4).Kind == kwSOURCE {
		p.advance()
		p.advance()
		p.advance()
		p.advance()
		p.advance()
		if p.match(kwAND) {
			if err := p.parseBooleanExpression(); err != nil {
				return nil, err
			}
		}
		if _, err := p.expect(kwTHEN); err != nil {
			return nil, err
		}
		if err := p.parseNotMatchedBySourceAction(); err != nil {
			return nil, err
		}
	}
	return &ParsedStmt{Kind: StmtMerge, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
}

func (p *Parser) isNotAt(offset int) bool {
	k := p.peek(offset).Kind
	return k == kwNOT || k == tokBang
}

func (p *Parser) parseMatchedAction() error {
	if p.match(kwDELETE) {
		return nil
	}
	if _, err := p.expect(kwUPDATE); err != nil {
		return err
	}
	if _, err := p.expect(kwSET); err != nil {
		return err
	}
	if p.match(tokAsterisk) {
		return nil
	}
	return p.parseAssignmentList()
}

func (p *Parser) parseNotMatchedAction() error {
	if _, err := p.expect(kwINSERT); err != nil {
		return err
	}
	if p.match(tokAsterisk) {
		return nil
	}
	if _, err := p.expect(tokLParen); err != nil {
		return err
	}
	if err := p.parseMultipartIdentifierList(); err != nil {
		return err
	}
	if _, err := p.expect(tokRParen); err != nil {
		return err
	}
	if _, err := p.expect(kwVALUES); err != nil {
		return err
	}
	if _, err := p.expect(tokLParen); err != nil {
		return err
	}
	if err := p.parseExpressionSeq(); err != nil {
		return err
	}
	_, err := p.expect(tokRParen)
	return err
}

func (p *Parser) parseNotMatchedBySourceAction() error {
	if p.match(kwDELETE) {
		return nil
	}
	if _, err := p.expect(kwUPDATE); err != nil {
		return err
	}
	if _, err := p.expect(kwSET); err != nil {
		return err
	}
	return p.parseAssignmentList()
}

func (p *Parser) parseLoadData() (Statement, error) {
	start := p.cur().Start
	if _, err := p.expect(kwLOAD); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwDATA); err != nil {
		return nil, err
	}
	p.match(kwLOCAL)
	if _, err := p.expect(kwINPATH); err != nil {
		return nil, err
	}
	if err := p.parseStringLit(); err != nil {
		return nil, err
	}
	p.match(kwOVERWRITE)
	if _, err := p.expect(kwINTO); err != nil {
		return nil, err
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
	return &ParsedStmt{Kind: StmtLoadData, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
}

func (p *Parser) parseTruncateTable() (Statement, error) {
	start := p.cur().Start
	if _, err := p.expect(kwTRUNCATE); err != nil {
		return nil, err
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
	return &ParsedStmt{Kind: StmtTruncate, Target: target, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
}
