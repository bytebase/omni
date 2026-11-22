package parser

import (
	"fmt"
	"strings"
)

const maxNestingDepth = 1000

// checkpoint stores parser cursor state for LL(k) backtracking.
type checkpoint struct {
	pos        int
	gtPending  int
	gtStart    int
	gtEnd      int
	parenDepth int
}

// Parser is a hand-written recursive-descent parser for Databricks / Spark SQL.
type Parser struct {
	src                 string
	toks                []Token
	pos                 int
	gtPending           int
	gtStart             int
	gtEnd               int
	depth               int
	parenDepth          int
	openCaseParenDepths []int
	caseScanStart       int
	caseScanEnd         int
	caseScanParenDepth  int
	caseScanHasWhen     bool
	caseScanEnds        int
}

// Parse parses a Databricks / Spark SQL input string (single statement or
// compound BEGIN ... END script block) and returns a ParseResult.
func Parse(sql string) (*ParseResult, error) {
	toks, err := lex(sql)
	if err != nil {
		return nil, err
	}
	p := &Parser{
		src:  sql,
		toks: toks,
	}
	stmt, err := p.parseCompoundOrSingleStatement()
	if err != nil {
		return nil, err
	}
	if p.cur().Kind != tokEOF {
		return nil, p.errorf("unexpected token %q", p.cur().Text)
	}
	return &ParseResult{Statements: []Statement{stmt}}, nil
}

func (p *Parser) save() checkpoint {
	return checkpoint{
		pos:        p.pos,
		gtPending:  p.gtPending,
		gtStart:    p.gtStart,
		gtEnd:      p.gtEnd,
		parenDepth: p.parenDepth,
	}
}

func (p *Parser) restore(cp checkpoint) {
	p.pos = cp.pos
	p.gtPending = cp.gtPending
	p.gtStart = cp.gtStart
	p.gtEnd = cp.gtEnd
	p.parenDepth = cp.parenDepth
}

func (p *Parser) enter() error {
	p.depth++
	if p.depth > maxNestingDepth {
		return p.errorf("maximum nesting depth exceeded")
	}
	return nil
}

func (p *Parser) leave() {
	p.depth--
}

func (p *Parser) cur() Token {
	if p.gtPending > 0 {
		return Token{Kind: tokGT, Text: ">", Start: p.gtStart, End: p.gtEnd}
	}
	if p.pos >= len(p.toks) {
		return p.toks[len(p.toks)-1]
	}
	return p.toks[p.pos]
}

func (p *Parser) peek(offset int) Token {
	if offset == 0 {
		return p.cur()
	}
	if p.gtPending > 0 {
		if offset < p.gtPending {
			return Token{Kind: tokGT, Text: ">", Start: p.gtStart, End: p.gtEnd}
		}
		offset -= p.gtPending
	}
	idx := p.pos + offset
	if idx < 0 {
		idx = 0
	}
	if idx >= len(p.toks) {
		return p.toks[len(p.toks)-1]
	}
	return p.toks[idx]
}

func (p *Parser) prevEnd() int {
	if p.gtPending > 0 {
		return p.gtEnd
	}
	if p.pos > 0 && p.pos-1 < len(p.toks) {
		return p.toks[p.pos-1].End
	}
	return 0
}

func (p *Parser) advance() Token {
	if p.gtPending > 0 {
		p.gtPending--
		return Token{Kind: tokGT, Text: ">", Start: p.gtStart, End: p.gtEnd}
	}
	tok := p.cur()
	if p.pos < len(p.toks)-1 {
		p.pos++
		if tok.Kind == tokLParen || tok.Kind == tokLBracket {
			p.parenDepth++
		} else if tok.Kind == tokRParen || tok.Kind == tokRBracket {
			p.parenDepth--
		}
	}
	return tok
}

func (p *Parser) match(kinds ...TokenKind) bool {
	k := p.cur().Kind
	for _, want := range kinds {
		if k == want {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) expect(kind TokenKind) (Token, error) {
	if p.cur().Kind == kind {
		return p.advance(), nil
	}
	return Token{}, p.errorf("unexpected token %q", p.cur().Text)
}

// expectGT consumes a '>' closing angle bracket in complex data types, splitting
// '>>' or '>>>' tokens cleanly when nested types like ARRAY<ARRAY<INT>> occur.
func (p *Parser) expectGT() error {
	if p.gtPending > 0 {
		p.gtPending--
		return nil
	}
	tok := p.cur()
	switch tok.Kind {
	case tokGT:
		p.advance()
		return nil
	case tokShiftRight:
		p.advance()
		p.gtPending = 1
		p.gtStart = tok.Start + 1
		p.gtEnd = tok.End
		return nil
	case tokShiftRightU:
		p.advance()
		p.gtPending = 2
		p.gtStart = tok.Start + 1
		p.gtEnd = tok.End
		return nil
	default:
		return p.errorf("expected '>', got %q", tok.Text)
	}
}

func (p *Parser) errorf(format string, args ...any) error {
	return &ParseError{
		Message:  fmt.Sprintf(format, args...),
		Position: p.cur().Start,
	}
}

// parseCompoundOrSingleStatement implements:
//
//	compoundOrSingleStatement
//	    : singleStatement
//	    | singleCompoundStatement
//	    ;
func (p *Parser) parseCompoundOrSingleStatement() (Statement, error) {
	if p.cur().Kind == kwBEGIN {
		cp := p.save()
		stmt, err := p.parseSingleCompoundStatement()
		if err == nil && p.cur().Kind == tokEOF {
			return stmt, nil
		}
		p.restore(cp)
	}
	return p.parseSingleStatement()
}

// parseSingleCompoundStatement implements:
//
//	singleCompoundStatement
//	    : BEGIN (NOT ATOMIC)? compoundBody? END SEMICOLON? EOF
//	    ;
func (p *Parser) parseSingleCompoundStatement() (Statement, error) {
	if err := p.enter(); err != nil {
		return nil, err
	}
	defer p.leave()
	start := p.cur().Start
	if _, err := p.expect(kwBEGIN); err != nil {
		return nil, err
	}
	if p.cur().Kind == kwNOT && p.peek(1).Kind == kwATOMIC {
		p.advance()
		p.advance()
	}
	if p.cur().Kind != kwEND {
		if err := p.parseCompoundBody(); err != nil {
			return nil, err
		}
	}
	if _, err := p.expect(kwEND); err != nil {
		return nil, err
	}
	p.match(tokSemicolon)
	if p.cur().Kind != tokEOF {
		return nil, p.errorf("unexpected token %q after END", p.cur().Text)
	}
	return &ParsedStmt{Kind: StmtCompound, Loc: Loc{Start: start, End: p.prevEnd()}}, nil
}

// parseBeginEndCompoundBlock implements:
//
//	beginEndCompoundBlock
//	    : beginLabel? BEGIN (NOT ATOMIC)? compoundBody? END endLabel?
//	    ;
func (p *Parser) parseBeginEndCompoundBlock() error {
	if err := p.enter(); err != nil {
		return err
	}
	defer p.leave()
	p.tryParseBeginLabel()
	if _, err := p.expect(kwBEGIN); err != nil {
		return err
	}
	if p.cur().Kind == kwNOT && p.peek(1).Kind == kwATOMIC {
		p.advance()
		p.advance()
	}
	if p.cur().Kind != kwEND {
		if err := p.parseCompoundBody(); err != nil {
			return err
		}
	}
	if _, err := p.expect(kwEND); err != nil {
		return err
	}
	if p.cur().Kind != tokSemicolon && p.cur().Kind != tokEOF && p.cur().Kind.isIdentifier() {
		if _, err := p.parseMultipartIdentifier(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) tryParseBeginLabel() bool {
	cp := p.save()
	if _, err := p.parseMultipartIdentifier(); err != nil {
		p.restore(cp)
		return false
	}
	if p.cur().Kind == tokColon {
		p.advance()
		return true
	}
	p.restore(cp)
	return false
}

// parseCompoundBody implements:
//
//	compoundBody : (compoundStatements+=compoundStatement SEMICOLON)+ ;
func (p *Parser) parseCompoundBody() error {
	count := 0
	for {
		if p.isCompoundBodyEnd() {
			break
		}
		if err := p.parseCompoundStatement(); err != nil {
			return err
		}
		if _, err := p.expect(tokSemicolon); err != nil {
			return err
		}
		count++
	}
	if count == 0 {
		return p.errorf("expected statement in compound body")
	}
	return nil
}

func (p *Parser) isCompoundBodyEnd() bool {
	switch p.cur().Kind {
	case kwEND, kwELSE, kwELSEIF, kwUNTIL, kwWHEN, tokEOF:
		return true
	}
	return false
}

func (p *Parser) parseCompoundStatement() error {
	if err := p.enter(); err != nil {
		return err
	}
	defer p.leave()

	// Try labeled compound statements first: beginLabel followed by BEGIN/WHILE/REPEAT/LOOP/FOR.
	cpLabel := p.save()
	if p.tryParseBeginLabel() {
		k := p.cur().Kind
		p.restore(cpLabel)
		switch k {
		case kwBEGIN:
			return p.parseBeginEndCompoundBlock()
		case kwWHILE:
			return p.parseWhileStatement()
		case kwREPEAT:
			return p.parseRepeatStatement()
		case kwLOOP:
			return p.parseLoopStatement()
		case kwFOR:
			return p.parseForStatement()
		}
	}

	switch p.cur().Kind {
	case kwBEGIN:
		return p.parseBeginEndCompoundBlock()
	case kwIF:
		return p.parseIfElseStatement()
	case kwCASE:
		return p.parseCaseStatement()
	case kwWHILE:
		return p.parseWhileStatement()
	case kwREPEAT:
		return p.parseRepeatStatement()
	case kwLEAVE:
		p.advance()
		_, err := p.parseMultipartIdentifier()
		return err
	case kwITERATE:
		p.advance()
		_, err := p.parseMultipartIdentifier()
		return err
	case kwLOOP:
		return p.parseLoopStatement()
	case kwFOR:
		return p.parseForStatement()
	case kwDECLARE:
		cp := p.save()
		if err := p.parseDeclareConditionStatement(); err == nil && p.cur().Kind == tokSemicolon {
			return nil
		}
		p.restore(cp)
		if err := p.parseDeclareHandlerStatement(); err == nil && p.cur().Kind == tokSemicolon {
			return nil
		}
		p.restore(cp)
	case kwSET:
		cp := p.save()
		if err := p.parseSetStatementInsideSqlScript(); err == nil && p.cur().Kind == tokSemicolon {
			return nil
		}
		p.restore(cp)
	}

	_, err := p.parseStatement()
	return err
}

func (p *Parser) parseSetStatementInsideSqlScript() error {
	if _, err := p.expect(kwSET); err != nil {
		return err
	}
	if p.cur().Kind == tokLParen {
		p.advance()
		if err := p.parseMultipartIdentifierList(); err != nil {
			return err
		}
		if _, err := p.expect(tokRParen); err != nil {
			return err
		}
		if _, err := p.expect(tokEq); err != nil {
			return err
		}
		if _, err := p.expect(tokLParen); err != nil {
			return err
		}
		if err := p.parseQuery(); err != nil {
			return err
		}
		_, err := p.expect(tokRParen)
		return err
	}
	return p.parseAssignmentList()
}

func (p *Parser) parseDeclareConditionStatement() error {
	if _, err := p.expect(kwDECLARE); err != nil {
		return err
	}
	if _, err := p.parseMultipartIdentifier(); err != nil {
		return err
	}
	if _, err := p.expect(kwCONDITION); err != nil {
		return err
	}
	if p.match(kwFOR) {
		if _, err := p.expect(kwSQLSTATE); err != nil {
			return err
		}
		p.match(kwVALUE)
		if err := p.parseStringLit(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parseDeclareHandlerStatement() error {
	if _, err := p.expect(kwDECLARE); err != nil {
		return err
	}
	if !p.match(kwCONTINUE, kwEXIT) {
		return p.errorf("expected CONTINUE or EXIT")
	}
	if _, err := p.expect(kwHANDLER); err != nil {
		return err
	}
	if _, err := p.expect(kwFOR); err != nil {
		return err
	}
	// conditionValues
	for {
		switch {
		case p.cur().Kind == kwSQLSTATE:
			cp := p.save()
			p.advance()
			p.match(kwVALUE)
			if err := p.parseStringLit(); err != nil {
				p.restore(cp)
				if _, err := p.parseMultipartIdentifier(); err != nil {
					return err
				}
			}
		case p.match(kwSQLEXCEPTION):
		case p.cur().Kind == kwNOT && p.peek(1).Kind == kwFOUND:
			p.advance()
			p.advance()
		default:
			if _, err := p.parseMultipartIdentifier(); err != nil {
				return err
			}
		}
		if !p.match(tokComma) {
			break
		}
	}
	// Handler action: beginEndCompoundBlock | statement | setStatementInsideSqlScript
	cp := p.save()
	if err := p.parseBeginEndCompoundBlock(); err == nil && p.cur().Kind == tokSemicolon {
		return nil
	}
	p.restore(cp)
	if p.cur().Kind == kwSET {
		if err := p.parseSetStatementInsideSqlScript(); err == nil && p.cur().Kind == tokSemicolon {
			return nil
		}
		p.restore(cp)
	}
	_, err := p.parseStatement()
	return err
}

func (p *Parser) parseWhileStatement() error {
	p.tryParseBeginLabel()
	if _, err := p.expect(kwWHILE); err != nil {
		return err
	}
	if err := p.parseBooleanExpression(); err != nil {
		return err
	}
	if _, err := p.expect(kwDO); err != nil {
		return err
	}
	if err := p.parseCompoundBody(); err != nil {
		return err
	}
	if _, err := p.expect(kwEND); err != nil {
		return err
	}
	if _, err := p.expect(kwWHILE); err != nil {
		return err
	}
	if p.cur().Kind != tokSemicolon && p.cur().Kind != tokEOF && p.cur().Kind.isIdentifier() {
		if _, err := p.parseMultipartIdentifier(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parseIfElseStatement() error {
	if _, err := p.expect(kwIF); err != nil {
		return err
	}
	if err := p.parseBooleanExpression(); err != nil {
		return err
	}
	if _, err := p.expect(kwTHEN); err != nil {
		return err
	}
	if err := p.parseCompoundBody(); err != nil {
		return err
	}
	for p.match(kwELSEIF) {
		if err := p.parseBooleanExpression(); err != nil {
			return err
		}
		if _, err := p.expect(kwTHEN); err != nil {
			return err
		}
		if err := p.parseCompoundBody(); err != nil {
			return err
		}
	}
	if p.match(kwELSE) {
		if err := p.parseCompoundBody(); err != nil {
			return err
		}
	}
	if _, err := p.expect(kwEND); err != nil {
		return err
	}
	_, err := p.expect(kwIF)
	return err
}

func (p *Parser) parseRepeatStatement() error {
	p.tryParseBeginLabel()
	if _, err := p.expect(kwREPEAT); err != nil {
		return err
	}
	if err := p.parseCompoundBody(); err != nil {
		return err
	}
	if _, err := p.expect(kwUNTIL); err != nil {
		return err
	}
	if err := p.parseBooleanExpression(); err != nil {
		return err
	}
	if _, err := p.expect(kwEND); err != nil {
		return err
	}
	if _, err := p.expect(kwREPEAT); err != nil {
		return err
	}
	if p.cur().Kind != tokSemicolon && p.cur().Kind != tokEOF && p.cur().Kind.isIdentifier() {
		if _, err := p.parseMultipartIdentifier(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parseCaseStatement() error {
	if _, err := p.expect(kwCASE); err != nil {
		return err
	}
	if p.cur().Kind != kwWHEN {
		if err := p.parseExpression(); err != nil {
			return err
		}
	}
	whenCount := 0
	for p.match(kwWHEN) {
		if err := p.parseExpression(); err != nil {
			return err
		}
		if _, err := p.expect(kwTHEN); err != nil {
			return err
		}
		if err := p.parseCompoundBody(); err != nil {
			return err
		}
		whenCount++
	}
	if whenCount == 0 {
		return p.errorf("expected WHEN in CASE statement")
	}
	if p.match(kwELSE) {
		if err := p.parseCompoundBody(); err != nil {
			return err
		}
	}
	if _, err := p.expect(kwEND); err != nil {
		return err
	}
	_, err := p.expect(kwCASE)
	return err
}

func (p *Parser) parseLoopStatement() error {
	p.tryParseBeginLabel()
	if _, err := p.expect(kwLOOP); err != nil {
		return err
	}
	if err := p.parseCompoundBody(); err != nil {
		return err
	}
	if _, err := p.expect(kwEND); err != nil {
		return err
	}
	if _, err := p.expect(kwLOOP); err != nil {
		return err
	}
	if p.cur().Kind != tokSemicolon && p.cur().Kind != tokEOF && p.cur().Kind.isIdentifier() {
		if _, err := p.parseMultipartIdentifier(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parseForStatement() error {
	p.tryParseBeginLabel()
	if _, err := p.expect(kwFOR); err != nil {
		return err
	}
	// Optional (multipartIdentifier AS)? before query.
	cp := p.save()
	if _, err := p.parseMultipartIdentifier(); err == nil && p.match(kwAS) {
		// consumed label AS
	} else {
		p.restore(cp)
	}
	if err := p.parseQuery(); err != nil {
		return err
	}
	if _, err := p.expect(kwDO); err != nil {
		return err
	}
	if err := p.parseCompoundBody(); err != nil {
		return err
	}
	if _, err := p.expect(kwEND); err != nil {
		return err
	}
	if _, err := p.expect(kwFOR); err != nil {
		return err
	}
	if p.cur().Kind != tokSemicolon && p.cur().Kind != tokEOF && p.cur().Kind.isIdentifier() {
		if _, err := p.parseMultipartIdentifier(); err != nil {
			return err
		}
	}
	return nil
}

// parseSingleStatement implements:
//
//	singleStatement : (statement | setResetStatement) SEMICOLON* EOF ;
func (p *Parser) parseSingleStatement() (Statement, error) {
	if p.cur().Kind == tokEOF {
		return nil, p.errorf("empty SQL input")
	}
	stmt, err := p.parseStatementOrSetReset()
	if err != nil {
		return nil, err
	}
	for p.match(tokSemicolon) {
	}
	if p.cur().Kind != tokEOF {
		return nil, p.errorf("unexpected trailing token %q", p.cur().Text)
	}
	return stmt, nil
}

func (p *Parser) parseStatementOrSetReset() (Statement, error) {
	if p.cur().Kind == kwSET {
		if p.peek(1).Kind == kwCATALOG {
			return p.parseSetCatalog()
		}
		return p.parseSetResetStatement()
	}
	if p.cur().Kind == kwRESET {
		return p.parseSetResetStatement()
	}
	return p.parseStatement()
}

func (p *Parser) isStatementEnd() bool {
	return p.cur().Kind == tokEOF || p.cur().Kind == tokSemicolon
}

func (p *Parser) parseStatement() (Statement, error) {
	if err := p.enter(); err != nil {
		return nil, err
	}
	defer p.leave()

	start := p.cur().Start
	switch p.cur().Kind {
	case kwSELECT, kwVALUES, kwTABLE, tokLParen, kwMAP, kwREDUCE:
		if err := p.parseQuery(); err != nil {
			return nil, err
		}
		return &ParsedStmt{Kind: StmtQuery, Loc: Loc{Start: start, End: p.prevEnd()}}, nil

	case kwFROM:
		// Could be multiInsertQuery (fromClause multiInsertQueryBody+) or query (fromStatement).
		cp := p.save()
		if stmt, err := p.parseMultiInsertQuery(); err == nil && p.isStatementEnd() {
			return stmt, nil
		}
		p.restore(cp)
		if err := p.parseQuery(); err != nil {
			return nil, err
		}
		return &ParsedStmt{Kind: StmtQuery, Loc: Loc{Start: start, End: p.prevEnd()}}, nil

	case kwWITH:
		// Could be ctes + dmlStatementNoWith OR query (ctes + queryTerm + queryOrganization).
		cp := p.save()
		if stmt, err := p.parseDMLWithCTEs(); err == nil && p.isStatementEnd() {
			return stmt, nil
		}
		p.restore(cp)
		if err := p.parseQuery(); err != nil {
			return nil, err
		}
		return &ParsedStmt{Kind: StmtQuery, Loc: Loc{Start: start, End: p.prevEnd()}}, nil

	case kwINSERT, kwUPDATE, kwDELETE, kwMERGE:
		return p.parseDMLStatementNoWith()

	case kwEXECUTE:
		return p.parseExecuteImmediate()

	case kwUSE:
		return p.parseUse()

	case kwSET:
		if p.peek(1).Kind == kwCATALOG {
			return p.parseSetCatalog()
		}
		return nil, p.errorf("unexpected SET")

	case kwCREATE:
		return p.parseCreateStatement()

	case kwREPLACE:
		return p.parseReplaceTable()

	case kwALTER:
		return p.parseAlterStatement()

	case kwDROP:
		return p.parseDropStatement()

	case kwANALYZE:
		return p.parseAnalyze()

	case kwDECLARE:
		return p.parseCreateVariable()

	case kwEXPLAIN:
		return p.parseExplain()

	case kwSHOW:
		return p.parseShow()

	case kwDESC, kwDESCRIBE:
		return p.parseDescribe()

	case kwCOMMENT:
		return p.parseCommentOn()

	case kwREFRESH:
		return p.parseRefresh()

	case kwCACHE:
		return p.parseCacheTable()

	case kwUNCACHE:
		return p.parseUncacheTable()

	case kwCLEAR:
		return p.parseClearCache()

	case kwLOAD:
		return p.parseLoadData()

	case kwTRUNCATE:
		return p.parseTruncateTable()

	case kwMSCK, kwREPAIR:
		return p.parseRepairTable()

	case kwADD, kwLIST:
		return p.parseManageResource()

	case kwCALL:
		return p.parseCall()

	case kwGRANT, kwREVOKE, kwEXPORT, kwIMPORT, kwLOCK, kwUNLOCK, kwSTART, kwCOMMIT, kwROLLBACK, kwDFS:
		return p.parseUnsupportedHiveNativeCommand()
	}

	return nil, p.errorf("unexpected token %q at start of statement", p.cur().Text)
}

func isWord(tok Token, upper string) bool {
	return strings.EqualFold(tok.Text, upper)
}
