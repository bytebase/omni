package parser

import (
	"fmt"
	"strings"
)

// checkpoint stores parser cursor state for LL(k) backtracking.
type checkpoint struct {
	pos        int
	gtPending  int
	gtStart    int
	gtEnd      int
	parenDepth int
}

// Parser is a hand-written recursive-descent parser for Databricks / Spark SQL.
type caseScanEntry struct {
	start       int
	end         int
	lastWhenPos int
	ends        int
}

type openCaseRun struct {
	parenDepth int
	count      int
}

type scriptScope struct {
	label      string              // normalized label name (empty for unlabeled compound block)
	isLoop     bool                // true for WHILE/REPEAT/LOOP/FOR; false for BEGIN..END
	forVar     string              // normalized FOR cursor variable name, if any
	conditions map[string]struct{} // declared condition names in this compound block
}

type Parser struct {
	src                    string
	toks                   []Token
	pos                    int
	gtPending              int
	gtStart                int
	gtEnd                  int
	depth                  int
	parenDepth             int
	openCaseRuns           []openCaseRun
	caseScanCache          []caseScanEntry
	caseScanConfirmedOpen  map[int]int
	matchingClose          []int
	lparenRunStart         int
	lparenRunEnd           int
	inCompound             bool
	inFromStatement        bool
	disallowParamMarkers   bool
	disallowColumnRefs     bool
	disallowGeneratedCols  bool
	relaxFunctionDataTypes bool
	scriptScopes           []scriptScope
}

// Parse parses a Databricks / Spark SQL input string (one or more
// semicolon-separated statements and/or compound BEGIN ... END script blocks)
// and returns a ParseResult with byte offsets relative to sql.
func Parse(sql string) (*ParseResult, error) {
	toks, err := lex(sql)
	if err != nil {
		return nil, err
	}
	p := &Parser{
		src:  sql,
		toks: toks,
	}
	if p.cur().Kind == tokEOF {
		return nil, p.errorf("empty SQL input")
	}
	var stmts []Statement
	for p.cur().Kind != tokEOF {
		stmt, err := p.parseTopLevelStatement()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
		if p.cur().Kind == tokEOF {
			break
		}
		if !p.match(tokSemicolon) {
			return nil, p.errorf("unexpected trailing token %q", p.cur().Text)
		}
		// Spark's singleStatement permits SEMICOLON* whereas
		// singleCompoundStatement permits at most one trailing ';'.
		if ps, ok := stmt.(*ParsedStmt); !ok || ps.Kind != StmtCompound {
			for p.match(tokSemicolon) {
			}
		}
	}
	return &ParseResult{Statements: stmts}, nil
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

func (p *Parser) ensureMatchingClose() {
	if p.matchingClose != nil {
		return
	}
	mc := make([]int, len(p.toks))
	var stack []int
	for i, tok := range p.toks {
		switch tok.Kind {
		case tokLParen, tokLBracket:
			stack = append(stack, i)
		case tokRParen, tokRBracket:
			if len(stack) > 0 {
				openIdx := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				mc[openIdx] = i
			}
		case tokSemicolon:
			stack = stack[:0]
		}
	}
	p.matchingClose = mc
}

// matchingCloseOffset returns the peek offset of the ')' or ']' matching the
// '(' or '[' at openOffset (relative to p.cur()), or -1 if unmatched before
// ';' or EOF.
func (p *Parser) matchingCloseOffset(openOffset int) int {
	if p.gtPending > 0 {
		depth := 0
		for i := openOffset; ; i++ {
			k := p.peek(i).Kind
			if k == tokEOF || k == tokSemicolon {
				return -1
			}
			if k == tokLParen || k == tokLBracket {
				depth++
			} else if k == tokRParen || k == tokRBracket {
				depth--
				if depth == 0 {
					return i
				}
			}
		}
	}
	p.ensureMatchingClose()
	idx := p.pos + openOffset
	if idx < 0 || idx >= len(p.matchingClose) {
		return -1
	}
	closeIdx := p.matchingClose[idx]
	if closeIdx <= idx {
		return -1
	}
	return closeIdx - p.pos
}

// skipLParensOffset returns the first peek offset >= startOffset whose token is
// not '(', caching consecutive '(' runs so deeply nested parenthesized
// expressions do not rescan leading parens quadratically.
func (p *Parser) skipLParensOffset(startOffset int) int {
	if p.gtPending == 0 {
		abs := p.pos + startOffset
		if abs >= p.lparenRunStart && abs < p.lparenRunEnd {
			return p.lparenRunEnd - p.pos
		}
	}
	i := startOffset
	for p.peek(i).Kind == tokLParen {
		i++
	}
	if p.gtPending == 0 && i > startOffset {
		p.lparenRunStart = p.pos + startOffset
		p.lparenRunEnd = p.pos + i
	}
	return i
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

// parseTopLevelStatement parses one top-level singleStatement or
// singleCompoundStatement without consuming its trailing semicolon(s).
func (p *Parser) parseTopLevelStatement() (Statement, error) {
	if p.cur().Kind == kwBEGIN {
		return p.parseSingleCompoundStatement()
	}
	return p.parseStatementOrSetReset()
}

// parseSingleCompoundStatement implements:
//
//	singleCompoundStatement
//	    : BEGIN (NOT ATOMIC)? compoundBody? END SEMICOLON? EOF
//	    ;
func (p *Parser) parseSingleCompoundStatement() (Statement, error) {
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()
	start := p.cur().Start
	if _, err := p.expect(kwBEGIN); err != nil {
		return nil, err
	}
	if p.cur().Kind == kwNOT && p.peek(1).Kind == kwATOMIC {
		p.advance()
		p.advance()
	}
	prevCompound := p.inCompound
	p.inCompound = true
	p.scriptScopes = append(p.scriptScopes, scriptScope{conditions: make(map[string]struct{})})
	defer func() {
		p.inCompound = prevCompound
		p.scriptScopes = p.scriptScopes[:len(p.scriptScopes)-1]
	}()
	if p.cur().Kind != kwEND || p.peek(1).Kind == tokColon {
		if err := p.parseCompoundBody(true); err != nil {
			return nil, err
		}
	}
	if _, err := p.expect(kwEND); err != nil {
		return nil, err
	}
	end := p.prevEnd()
	return &ParsedStmt{Kind: StmtCompound, Loc: Loc{Start: start, End: end}}, nil
}

// parseBeginEndCompoundBlock implements:
//
//	beginEndCompoundBlock
//	    : beginLabel? BEGIN (NOT ATOMIC)? compoundBody? END endLabel?
//	    ;
func (p *Parser) parseBeginEndCompoundBlock() error {
	if err := p.enterNesting(); err != nil {
		return err
	}
	defer p.leaveNesting()
	beginLabel, err := p.parseOptionalBeginLabel()
	if err != nil {
		return err
	}
	if _, err := p.expect(kwBEGIN); err != nil {
		return err
	}
	if p.cur().Kind == kwNOT && p.peek(1).Kind == kwATOMIC {
		p.advance()
		p.advance()
	}
	prevCompound := p.inCompound
	p.inCompound = true
	p.scriptScopes = append(p.scriptScopes, scriptScope{
		label:      beginLabel,
		isLoop:     false,
		conditions: make(map[string]struct{}),
	})
	defer func() {
		p.inCompound = prevCompound
		p.scriptScopes = p.scriptScopes[:len(p.scriptScopes)-1]
	}()
	if p.cur().Kind != kwEND || p.peek(1).Kind == tokColon {
		if err := p.parseCompoundBody(true); err != nil {
			return err
		}
	}
	if _, err := p.expect(kwEND); err != nil {
		return err
	}
	return p.parseOptionalEndLabel(beginLabel)
}

func (p *Parser) hasActiveLabelOrForVar(norm string) bool {
	for _, sc := range p.scriptScopes {
		if sc.label == norm || sc.forVar == norm {
			return true
		}
	}
	return false
}

func (p *Parser) findActiveLabel(norm string) (found bool, isLoop bool) {
	for i := len(p.scriptScopes) - 1; i >= 0; i-- {
		if p.scriptScopes[i].label == norm {
			return true, p.scriptScopes[i].isLoop
		}
	}
	return false, false
}

func (p *Parser) hasConditionInScope(norm string) bool {
	for i := len(p.scriptScopes) - 1; i >= 0; i-- {
		if _, ok := p.scriptScopes[i].conditions[norm]; ok {
			return true
		}
	}
	return false
}

func (p *Parser) hasBeginLabelAhead() bool {
	cp := p.save()
	defer p.restore(cp)
	if _, err := p.parseMultipartIdentifier(); err != nil {
		return false
	}
	return p.cur().Kind == tokColon
}

func (p *Parser) parseOptionalBeginLabel() (string, error) {
	cp := p.save()
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		p.restore(cp)
		return "", nil
	}
	if p.cur().Kind != tokColon {
		p.restore(cp)
		return "", nil
	}
	p.advance()
	if isQualifiedIdent(name) {
		return "", p.errorf("qualified label name %q is not allowed", name)
	}
	norm := normalizeIdentName(name)
	if p.hasActiveLabelOrForVar(norm) {
		return "", p.errorf("label %q already exists in enclosing scope", name)
	}
	return norm, nil
}

func (p *Parser) parseOptionalEndLabel(beginLabel string) error {
	if p.cur().Kind == tokSemicolon || p.cur().Kind == tokEOF || !p.cur().Kind.isIdentifier() {
		return nil
	}
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return err
	}
	if isQualifiedIdent(name) {
		return p.errorf("qualified end label %q is not allowed", name)
	}
	if beginLabel == "" {
		return p.errorf("end label %q cannot exist without begin label", name)
	}
	if normalizeIdentName(name) != beginLabel {
		return p.errorf("end label %q does not match begin label %q", name, beginLabel)
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
//
// When allowDecls is true (direct body of BEGIN ... END), DECLARE variable/condition
// statements must precede DECLARE HANDLER statements, which must precede executable
// statements, matching Spark v4.0.5 AstBuilder.visitCompoundBody. When allowDecls is
// false (body of IF/CASE/WHILE/REPEAT/LOOP/FOR), DECLARE statements are rejected.
func (p *Parser) parseCompoundBody(allowDecls bool) error {
	count := 0
	// phase: 0 = variable/condition declarations, 1 = handler declarations, 2 = statements
	phase := 0
	for {
		if p.isCompoundBodyEnd() {
			break
		}
		stmtPhase, err := p.parseCompoundStatement(allowDecls, phase)
		if err != nil {
			return err
		}
		if stmtPhase > phase {
			phase = stmtPhase
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
	// A compound-body terminator keyword followed by ':' is a beginLabel
	// (e.g. `BEGIN ELSE: BEGIN END ELSE; END`), not the end of the body.
	switch p.cur().Kind {
	case kwEND, kwELSE, kwELSEIF, kwUNTIL, kwWHEN:
		return p.peek(1).Kind != tokColon
	case tokEOF:
		return true
	}
	return false
}

func (p *Parser) parseCompoundStatement(allowDecls bool, phase int) (int, error) {
	if err := p.enterNesting(); err != nil {
		return 2, err
	}
	defer p.leaveNesting()

	// Try labeled compound statements first: beginLabel followed by BEGIN/WHILE/REPEAT/LOOP/FOR.
	cpLabel := p.save()
	if p.tryParseBeginLabel() {
		k := p.cur().Kind
		p.restore(cpLabel)
		switch k {
		case kwBEGIN:
			return 2, p.parseBeginEndCompoundBlock()
		case kwWHILE:
			return 2, p.parseWhileStatement()
		case kwREPEAT:
			return 2, p.parseRepeatStatement()
		case kwLOOP:
			return 2, p.parseLoopStatement()
		case kwFOR:
			return 2, p.parseForStatement()
		}
	}

	switch p.cur().Kind {
	case kwBEGIN:
		return 2, p.parseBeginEndCompoundBlock()
	case kwIF:
		return 2, p.parseIfElseStatement()
	case kwCASE:
		return 2, p.parseCaseStatement()
	case kwWHILE:
		return 2, p.parseWhileStatement()
	case kwREPEAT:
		return 2, p.parseRepeatStatement()
	case kwLEAVE:
		p.advance()
		name, err := p.parseMultipartIdentifier()
		if err != nil {
			return 2, err
		}
		if isQualifiedIdent(name) {
			return 2, p.errorf("qualified label %q is not allowed in LEAVE", name)
		}
		if found, _ := p.findActiveLabel(normalizeIdentName(name)); !found {
			return 2, p.errorf("label %q does not belong to any surrounding block", name)
		}
		return 2, nil
	case kwITERATE:
		p.advance()
		name, err := p.parseMultipartIdentifier()
		if err != nil {
			return 2, err
		}
		if isQualifiedIdent(name) {
			return 2, p.errorf("qualified label %q is not allowed in ITERATE", name)
		}
		found, isLoop := p.findActiveLabel(normalizeIdentName(name))
		if !found {
			return 2, p.errorf("label %q does not belong to any surrounding block", name)
		}
		if !isLoop {
			return 2, p.errorf("ITERATE cannot reference compound block label %q", name)
		}
		return 2, nil
	case kwLOOP:
		return 2, p.parseLoopStatement()
	case kwFOR:
		return 2, p.parseForStatement()
	case kwDECLARE:
		if !allowDecls {
			return 0, p.errorf("DECLARE is only allowed at the beginning of a BEGIN ... END compound block")
		}
		// Distinguish DECLARE CONDITION, DECLARE HANDLER, and DECLARE variable.
		if p.isDeclareHandlerAhead() {
			if phase > 1 {
				return 1, p.errorf("DECLARE HANDLER must appear before executable statements")
			}
			return 1, p.parseDeclareHandlerStatement()
		}
		if p.isDeclareConditionAhead() {
			if phase > 0 {
				return 0, p.errorf("DECLARE CONDITION must appear at the beginning of a compound block")
			}
			return 0, p.parseDeclareConditionStatement()
		}
		if phase > 0 {
			return 0, p.errorf("DECLARE variable must appear at the beginning of a compound block")
		}
		_, err := p.parseCreateVariable()
		return 0, err
	case kwSET:
		cp := p.save()
		if err := p.parseSetStatementInsideSqlScript(); err == nil && p.cur().Kind == tokSemicolon {
			return 2, nil
		}
		p.restore(cp)
	}

	_, err := p.parseStatement()
	return 2, err
}

func (p *Parser) isDeclareHandlerAhead() bool {
	k := p.peek(1).Kind
	return (k == kwCONTINUE || k == kwEXIT) && p.peek(2).Kind == kwHANDLER
}

func (p *Parser) isDeclareConditionAhead() bool {
	cp := p.save()
	defer p.restore(cp)
	if !p.match(kwDECLARE) {
		return false
	}
	if _, err := p.parseMultipartIdentifier(); err != nil {
		return false
	}
	return p.cur().Kind == kwCONDITION
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
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return err
	}
	if !isValidSparkConditionName(name) {
		return p.errorf("invalid condition name %q; must match [A-Za-z0-9_]+", name)
	}
	if _, err := p.expect(kwCONDITION); err != nil {
		return err
	}
	if p.match(kwFOR) {
		if _, err := p.expect(kwSQLSTATE); err != nil {
			return err
		}
		p.match(kwVALUE)
		if !p.isStringLit() {
			return p.errorf("expected SQLSTATE string literal, got %q", p.cur().Text)
		}
		tok := p.advance()
		if !validateSQLState(unquoteStringLit(tok.Text)) {
			return p.errorf("invalid SQLSTATE value %s", tok.Text)
		}
	}
	norm := normalizeIdentName(name)
	if len(p.scriptScopes) > 0 {
		if p.hasConditionInScope(norm) {
			return p.errorf("duplicate condition %q in scope", name)
		}
		p.scriptScopes[len(p.scriptScopes)-1].conditions[norm] = struct{}{}
	}
	return nil
}

func (p *Parser) parseDeclareHandlerStatement() error {
	if _, err := p.expect(kwDECLARE); err != nil {
		return err
	}
	if p.match(kwCONTINUE) {
		return p.errorf("CONTINUE exception handler is not supported; use EXIT handler")
	}
	if _, err := p.expect(kwEXIT); err != nil {
		return err
	}
	if _, err := p.expect(kwHANDLER); err != nil {
		return err
	}
	if _, err := p.expect(kwFOR); err != nil {
		return err
	}
	// conditionValues: Spark forbids mixing general conditions (SQLEXCEPTION / NOT FOUND)
	// with specific conditions (named condition / SQLSTATE) and forbids duplicate entries.
	var seenSQLException, seenNotFound, hasGeneral, hasSpecific bool
	seenStates := make(map[string]struct{})
	seenConds := make(map[string]struct{})
	for {
		switch {
		case p.cur().Kind == kwSQLSTATE:
			cp := p.save()
			p.advance()
			p.match(kwVALUE)
			if p.isStringLit() {
				tok := p.advance()
				val := unquoteStringLit(tok.Text)
				if !validateSQLState(val) {
					return p.errorf("invalid SQLSTATE value %s", tok.Text)
				}
				if _, dup := seenStates[val]; dup {
					return p.errorf("duplicate SQLSTATE %s in handler declaration", tok.Text)
				}
				seenStates[val] = struct{}{}
				hasSpecific = true
			} else {
				p.restore(cp)
				name, err := p.parseMultipartIdentifier()
				if err != nil {
					return err
				}
				norm := normalizeIdentName(name)
				if !p.hasConditionInScope(norm) && !isBuiltinSparkErrorCondition(name) {
					return p.errorf("condition %q is not declared", name)
				}
				if _, dup := seenConds[norm]; dup {
					return p.errorf("duplicate condition %q in handler declaration", name)
				}
				seenConds[norm] = struct{}{}
				hasSpecific = true
			}
		case p.match(kwSQLEXCEPTION):
			if seenSQLException {
				return p.errorf("duplicate SQLEXCEPTION in handler declaration")
			}
			seenSQLException = true
			hasGeneral = true
		case p.cur().Kind == kwNOT && p.peek(1).Kind == kwFOUND:
			p.advance()
			p.advance()
			if seenNotFound {
				return p.errorf("duplicate NOT FOUND in handler declaration")
			}
			seenNotFound = true
			hasGeneral = true
		default:
			name, err := p.parseMultipartIdentifier()
			if err != nil {
				return err
			}
			norm := normalizeIdentName(name)
			if !p.hasConditionInScope(norm) && !isBuiltinSparkErrorCondition(name) {
				return p.errorf("condition %q is not declared", name)
			}
			if _, dup := seenConds[norm]; dup {
				return p.errorf("duplicate condition %q in handler declaration", name)
			}
			seenConds[norm] = struct{}{}
			hasSpecific = true
		}
		if hasGeneral && hasSpecific {
			return p.errorf("cannot mix SQLEXCEPTION/NOT FOUND with specific conditions or SQLSTATE in handler declaration")
		}
		if !p.match(tokComma) {
			break
		}
	}
	// Handler action: beginEndCompoundBlock | statement | setStatementInsideSqlScript
	cp := p.save()
	if p.cur().Kind == kwBEGIN || p.hasBeginLabelAhead() {
		if err := p.parseBeginEndCompoundBlock(); err == nil && p.cur().Kind == tokSemicolon {
			return nil
		}
		p.restore(cp)
	}
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
	beginLabel, err := p.parseOptionalBeginLabel()
	if err != nil {
		return err
	}
	if _, err := p.expect(kwWHILE); err != nil {
		return err
	}
	if err := p.parseBooleanExpression(); err != nil {
		return err
	}
	if _, err := p.expect(kwDO); err != nil {
		return err
	}
	p.scriptScopes = append(p.scriptScopes, scriptScope{label: beginLabel, isLoop: true})
	defer func() { p.scriptScopes = p.scriptScopes[:len(p.scriptScopes)-1] }()
	if err := p.parseCompoundBody(false); err != nil {
		return err
	}
	if _, err := p.expect(kwEND); err != nil {
		return err
	}
	if _, err := p.expect(kwWHILE); err != nil {
		return err
	}
	return p.parseOptionalEndLabel(beginLabel)
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
	if err := p.parseCompoundBody(false); err != nil {
		return err
	}
	for p.match(kwELSEIF) {
		if err := p.parseBooleanExpression(); err != nil {
			return err
		}
		if _, err := p.expect(kwTHEN); err != nil {
			return err
		}
		if err := p.parseCompoundBody(false); err != nil {
			return err
		}
	}
	if p.match(kwELSE) {
		if err := p.parseCompoundBody(false); err != nil {
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
	beginLabel, err := p.parseOptionalBeginLabel()
	if err != nil {
		return err
	}
	if _, err := p.expect(kwREPEAT); err != nil {
		return err
	}
	p.scriptScopes = append(p.scriptScopes, scriptScope{label: beginLabel, isLoop: true})
	defer func() { p.scriptScopes = p.scriptScopes[:len(p.scriptScopes)-1] }()
	if err := p.parseCompoundBody(false); err != nil {
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
	return p.parseOptionalEndLabel(beginLabel)
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
		if err := p.parseCompoundBody(false); err != nil {
			return err
		}
		whenCount++
	}
	if whenCount == 0 {
		return p.errorf("expected WHEN in CASE statement")
	}
	if p.match(kwELSE) {
		if err := p.parseCompoundBody(false); err != nil {
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
	beginLabel, err := p.parseOptionalBeginLabel()
	if err != nil {
		return err
	}
	if _, err := p.expect(kwLOOP); err != nil {
		return err
	}
	p.scriptScopes = append(p.scriptScopes, scriptScope{label: beginLabel, isLoop: true})
	defer func() { p.scriptScopes = p.scriptScopes[:len(p.scriptScopes)-1] }()
	if err := p.parseCompoundBody(false); err != nil {
		return err
	}
	if _, err := p.expect(kwEND); err != nil {
		return err
	}
	if _, err := p.expect(kwLOOP); err != nil {
		return err
	}
	return p.parseOptionalEndLabel(beginLabel)
}

func (p *Parser) parseForStatement() error {
	beginLabel, err := p.parseOptionalBeginLabel()
	if err != nil {
		return err
	}
	if _, err := p.expect(kwFOR); err != nil {
		return err
	}
	// Optional (multipartIdentifier AS)? before query.
	forVar := ""
	cp := p.save()
	if name, err := p.parseMultipartIdentifier(); err == nil && p.match(kwAS) {
		if !isQualifiedIdent(name) {
			norm := normalizeIdentName(name)
			if norm == beginLabel || p.hasActiveLabelOrForVar(norm) {
				return p.errorf("FOR variable %q already exists in scope", name)
			}
			forVar = norm
		}
	} else {
		p.restore(cp)
	}
	if err := p.parseQuery(); err != nil {
		return err
	}
	if _, err := p.expect(kwDO); err != nil {
		return err
	}
	p.scriptScopes = append(p.scriptScopes, scriptScope{label: beginLabel, isLoop: true, forVar: forVar})
	defer func() { p.scriptScopes = p.scriptScopes[:len(p.scriptScopes)-1] }()
	if err := p.parseCompoundBody(false); err != nil {
		return err
	}
	if _, err := p.expect(kwEND); err != nil {
		return err
	}
	if _, err := p.expect(kwFOR); err != nil {
		return err
	}
	return p.parseOptionalEndLabel(beginLabel)
}

func (p *Parser) parseStatementOrSetReset() (Statement, error) {
	if p.cur().Kind == kwSET {
		if p.peek(1).Kind == kwCATALOG && !p.isStatementEndAt(2) {
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
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()

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
		if p.peek(1).Kind == kwCATALOG && !p.isStatementEndAt(2) {
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

	case kwGRANT, kwREVOKE:
		return p.parseGrantRevoke()
	}

	return nil, p.errorf("unexpected token %q at start of statement", p.cur().Text)
}

func isWord(tok Token, upper string) bool {
	return strings.EqualFold(tok.Text, upper)
}
