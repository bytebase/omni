// Package parser implements a recursive descent SQL parser for Oracle PL/SQL.
//
// This parser produces AST nodes from the oracle/ast package.
package parser

import (
	"errors"
	"fmt"

	nodes "github.com/bytebase/omni/oracle/ast"
)

// Parser is a recursive descent parser for Oracle SQL/PL/SQL.
//
// Oracle parser helpers follow the project-wide dual-return contract: required
// parse functions return (Node, error), optional probes use a tryParse or
// parseOptional name and may return (nil, nil), and nested parse errors must be
// returned to the caller instead of stored or discarded.
type Parser struct {
	lexer   *Lexer
	source  string // original SQL input
	cur     Token  // current token
	prev    Token  // previous token (for error reporting)
	nextBuf Token  // buffered next token for 2-token lookahead
	hasNext bool   // whether nextBuf is valid

	// caseKinds tells, for each CASE token from the first one parsePrimary
	// met to the end of the range, whether it opens a CASE expression (true)
	// or names a column (false). decideCaseTokens fills it in one pass.
	caseKinds map[int]bool

	// typeMods says what the parenthesized modifiers of the datatype being
	// parsed may hold; parsePLSQLDatatype sets it for PL/SQL datatypes.
	typeMods typeModMode

	// qualifiedCharsets holds the components of each item%CHARSET source
	// that is a qualified name (k.c), as opposed to one identifier that may
	// itself be a quoted "K.C".
	qualifiedCharsets map[*nodes.TypeName][]string

	// typeScopes holds, for each open PL/SQL declaration scope, the names
	// its TYPE and SUBTYPE declarations have declared so far. Oracle 23ai
	// lets one redeclare a predefined type name, NUMBER or SYS_REFCURSOR
	// alike, and the redeclaration then shadows STANDARD's type.
	typeScopes []map[string]bool

	// inPackageBody is set while a package body's items are parsed: its
	// specification may redeclare any predefined type name out of sight.
	inPackageBody bool

	// shadowedTypes records the datatypes spelled as a predefined type
	// where a redeclaration of that name is, or may be, in scope.
	shadowedTypes map[*nodes.TypeName]bool

	// inAnonymousBlock is set while a top-level anonymous block is parsed:
	// Oracle 23ai rejects RESULT_CACHE on any subprogram in one (PLS-00999).
	inAnonymousBlock bool

	// ccDirective, ccCondition, and ccBind index the first entries of the
	// lexer's conditional compilation record that no statement has checked
	// yet; ccAfterUnit says the last statement parsed is a PL/SQL unit, and
	// ccAfterBlock that it is an anonymous block.
	ccDirective  int
	ccCondition  int
	ccBind       int
	ccAfterUnit  bool
	ccAfterBlock bool
}

// checkConditionalCompilation checks the conditional compilation of stmt,
// whose text runs from start, after the previous statement, to end, its
// first token at first. Oracle preprocesses a PL/SQL unit alone, whose text
// runs on after its END: Oracle 23ai takes BEGIN NULL; END; $IF TRUE $THEN
// $END. It rejects a directive or an inquiry directive in any other
// statement, before one, or too early in one (firstDirectiveTooEarly) with
// ORA-00911, and one in an object type specification with PLS-00180. In a
// unit, the condition of each $IF and $ELSIF the lexer evaluated is a
// Boolean expression (checkDirectiveCondition).
func (p *Parser) checkConditionalCompilation(stmt nodes.StmtNode, start, first, end int) error {
	cc := p.lexer.cc
	if cc == nil {
		return nil
	}
	unit := isPreprocessedUnit(stmt)
	for ; p.ccDirective < len(cc.directives) && cc.directives[p.ccDirective] < end; p.ccDirective++ {
		loc := cc.directives[p.ccDirective]
		if loc < first && !p.ccAfterUnit || loc >= first && (!unit || p.firstDirectiveTooEarly(stmt, loc)) {
			return p.syntaxErrorAtDirective(loc)
		}
	}
	p.ccAfterUnit = unit
	// An anonymous block takes no bind variable inside a directive
	// (PLS-00110), in a selected branch or not; one before stmt's first
	// token is in the text of the statement before it.
	afterBlock := p.ccAfterBlock
	_, block := stmt.(*nodes.PLSQLBlock)
	p.ccAfterBlock = block
	for ; p.ccBind < len(cc.binds) && cc.binds[p.ccBind] < end; p.ccBind++ {
		if loc := cc.binds[p.ccBind]; loc < first && afterBlock || loc >= first && block {
			return p.syntaxErrorAtTok(NewLexerRange(p.source, loc, len(p.source)).nextRaw())
		}
	}
	return p.checkDirectiveConditions(end)
}

// checkDirectiveConditions checks the recorded conditions that start
// before end.
func (p *Parser) checkDirectiveConditions(end int) error {
	cc := p.lexer.cc
	for ; p.ccCondition < len(cc.conditions) && cc.conditions[p.ccCondition].Start < end; p.ccCondition++ {
		if err := p.checkDirectiveCondition(cc.conditions[p.ccCondition]); err != nil {
			return err
		}
	}
	return nil
}

// isPreprocessedUnit reports whether stmt is a PL/SQL unit, whose text
// Oracle preprocesses before compiling it.
func isPreprocessedUnit(stmt nodes.StmtNode) bool {
	switch s := stmt.(type) {
	case *nodes.PLSQLBlock, *nodes.CreateProcedureStmt, *nodes.CreateFunctionStmt,
		*nodes.CreatePackageStmt, *nodes.CreateTriggerStmt:
		return true
	case *nodes.CreateTypeStmt:
		return s.IsBody
	}
	return false
}

// syntaxErrorAtDirective returns a syntax error at the preprocessor token
// that starts at loc.
func (p *Parser) syntaxErrorAtDirective(loc int) *ParseError {
	end := loc + 1
	for end < len(p.source) && (p.source[end] == '$' || isIdentCont(p.source[end])) {
		end++
	}
	return p.syntaxErrorAtTok(Token{Type: tokIDENT, Loc: loc, End: end})
}

// parseAnonymousBlock parses a top-level anonymous block, noting that its
// subprograms are in one.
func (p *Parser) parseAnonymousBlock() (nodes.StmtNode, error) {
	saved := p.inAnonymousBlock
	p.inAnonymousBlock = true
	block, err := p.parsePLSQLBlock()
	p.inAnonymousBlock = saved
	if err != nil {
		return nil, err
	}
	if err := p.checkBlockCharsetScopes(block); err != nil {
		return nil, err
	}
	return block, nil
}

// Parse parses a SQL string into an AST list.
// Each statement is wrapped in a *RawStmt.
func Parse(sql string) (*nodes.List, error) {
	return ParseRange(sql, 0, len(sql))
}

// ParseRange parses source[start:end] in place. Every Loc and error Position
// is an offset into source, so a caller that split a script can parse each
// segment with absolute positions without copying or padding the text.
func ParseRange(source string, start, end int) (*nodes.List, error) {
	if err := validateBalancedDelimiters(source, start, end); err != nil {
		return nil, err
	}

	p := &Parser{
		lexer:  NewLexerRange(source, start, end),
		source: source,
	}
	p.lexer.cc = &ccRecord{}
	p.lexer.ccDecisions = &[]ccDecision{}
	p.lexer.literalLimit = 32767
	p.advance()

	var stmts []nodes.Node
	needSeparator := false
	// textStart is where the text of the next statement begins.
	textStart := start
	for p.cur.Type != tokEOF {
		// Skip semicolons
		if p.cur.Type == ';' {
			p.advance()
			needSeparator = false
			textStart = p.prev.End
			continue
		}
		if needSeparator {
			return nil, p.syntaxErrorAtCur()
		}
		if p.lexer.Err != nil {
			return nil, p.lexerError()
		}

		stmtLoc := p.pos()
		stmt, err := p.parseStmtSelections(textStart, stmtLoc)
		if err != nil {
			return nil, err
		}
		textStart = p.prev.End

		raw := &nodes.RawStmt{
			Stmt: stmt,
			Loc:  nodes.Loc{Start: stmtLoc, End: p.prev.End},
		}
		stmts = append(stmts, raw)
		needSeparator = true
	}
	if p.lexer.Err != nil {
		return nil, p.lexerError()
	}
	// A directive after the last statement belongs to it, if it is a unit.
	if cc := p.lexer.cc; p.ccDirective < len(cc.directives) && !p.ccAfterUnit {
		return nil, p.syntaxErrorAtDirective(cc.directives[p.ccDirective])
	}
	if err := p.checkDirectiveConditions(end); err != nil {
		return nil, err
	}
	// So are its bind variables, which an anonymous block takes none of
	// inside a directive (PLS-00110).
	if cc := p.lexer.cc; p.ccAfterBlock && p.ccBind < len(cc.binds) {
		loc := cc.binds[p.ccBind]
		return nil, p.syntaxErrorAtTok(NewLexerRange(p.source, loc, len(p.source)).nextRaw())
	}

	if len(stmts) == 0 {
		return &nodes.List{}, nil
	}
	return &nodes.List{Items: stmts}, nil
}

// ccMaxAttempts bounds the branch selections parseStmtSelections tries
// for one statement.
const ccMaxAttempts = 64

// parseStmtSelections parses the statement at the current token, with the
// checks that follow it, under each selection of the conditional
// compilation branches whose conditions the text does not settle, until
// one parses. Oracle 23ai compiles a unit under the selection its
// environment settles (the database version, PLSQL_CCFLAGS, package
// constants), which may be any of them, so a unit fails only when every
// selection does: $IF $$x $THEN garbage $ELSE NULL; $END parses, and
// $IF $$x $THEN $ERROR 'a' $END $ELSE $ERROR 'b' $END $END fails. The
// first selection takes each undecided $IF's first undecided branch; the
// next ones pass over more of them, the last $IF met first. Of the
// failures, the one that got furthest is reported.
func (p *Parser) parseStmtSelections(textStart, stmtLoc int) (nodes.StmtNode, error) {
	snap := p.snapshot()
	defer func() { p.lexer.ccChoices = nil }()
	var failure error
	furthest := -1
	for attempt := 1; ; attempt++ {
		stmt, err := p.parseCheckedStmt(textStart, stmtLoc)
		if err == nil {
			return stmt, nil
		}
		var pe *ParseError
		switch {
		case errors.As(err, &pe) && pe.Position > furthest:
			failure, furthest = err, pe.Position
		case failure == nil:
			failure = err
		}
		choices := nextSelection((*p.lexer.ccDecisions)[snap.decisions:])
		if choices == nil || attempt == ccMaxAttempts {
			return nil, failure
		}
		p.restore(snap)
		p.lexer.ccChoices = choices
	}
}

// parseCheckedStmt parses one statement and checks it: its lexing, its
// end, its conditional compilation, and that a ';' or the end follows.
func (p *Parser) parseCheckedStmt(textStart, stmtLoc int) (nodes.StmtNode, error) {
	stmt, err := p.parseStmt()
	if err != nil {
		return nil, err
	}
	if p.lexer.Err != nil {
		return nil, p.lexerError()
	}
	if p.cur.Type == tokEOF && (p.isIncompleteStatementEnd(stmt) || isIncompleteParsedStmt(stmt)) {
		return nil, p.syntaxErrorAtCur()
	}
	if stmt == nil {
		return nil, p.syntaxErrorAtCur()
	}
	if err := p.checkConditionalCompilation(stmt, textStart, stmtLoc, p.prev.End); err != nil {
		return nil, err
	}
	if p.cur.Type != ';' && p.cur.Type != tokEOF {
		return nil, p.syntaxErrorAtCur()
	}
	return stmt, nil
}

// nextSelection returns the branch choices of the selection that follows
// the one whose undecided $IF directives decided records in order, or nil
// when none is left: the last $IF that can pass over one more undecided
// branch does, those before it keep their choices, and those after it,
// which the new choice may move or remove, start over.
func nextSelection(decided []ccDecision) map[int]int {
	for i := len(decided) - 1; i >= 0; i-- {
		if decided[i].more {
			choices := map[int]int{decided[i].at: decided[i].passed + 1}
			for _, d := range decided[:i] {
				choices[d.at] = d.passed
			}
			return choices
		}
	}
	return nil
}

// parseSnapshot is the parser's state at the start of a statement, for
// parseStmtSelections to parse it again.
type parseSnapshot struct {
	parser                        Parser
	lexer                         Lexer
	directives, conditions, binds int
	decisions                     int
}

// snapshot returns the parser's state.
func (p *Parser) snapshot() parseSnapshot {
	s := parseSnapshot{parser: *p, lexer: *p.lexer}
	if cc := p.lexer.cc; cc != nil {
		s.directives, s.conditions, s.binds = len(cc.directives), len(cc.conditions), len(cc.binds)
	}
	if d := p.lexer.ccDecisions; d != nil {
		s.decisions = len(*d)
	}
	return s
}

// restore returns the parser to the state s.
func (p *Parser) restore(s parseSnapshot) {
	lexer := p.lexer
	*p = s.parser
	p.lexer = lexer
	*lexer = s.lexer
	if cc := lexer.cc; cc != nil {
		cc.directives, cc.conditions, cc.binds = cc.directives[:s.directives], cc.conditions[:s.conditions], cc.binds[:s.binds]
	}
	if d := lexer.ccDecisions; d != nil {
		*d = (*d)[:s.decisions]
	}
}

func validateBalancedDelimiters(sql string, start, end int) error {
	// Under a selection of undecided branches the delimiters need not
	// balance, as under another they may: the statement's own parse,
	// which tries them all, settles it.
	var decided []ccDecision
	lexer := NewLexerRange(sql, start, end)
	lexer.ccDecisions = &decided
	depth := 0
	var javaHead embeddedSourceHead
	var prev1, prev2, prev3 Token
	for {
		tok := lexer.NextToken()
		if lexer.Err != nil {
			return nil
		}
		if javaHead.observe(tok) {
			// The rest of the range is Java or JavaScript source, not SQL.
			return nil
		}
		if (prev3.Type == kwIS || prev3.Type == kwAS) &&
			prev2.Type == tokIDENT && prev2.Str == "MLE" && prev1.Type == tokIDENT && prev1.Str == "LANGUAGE" {
			// tok names the language of an inline MLE call spec; its
			// delimited code is not SQL.
			if _, codeEnd, _, ok := mleInlineCodeEnd(sql, tok.End, end); ok {
				lexer.pos = codeEnd
			}
		}
		prev3, prev2, prev1 = prev2, prev1, tok
		switch tok.Type {
		case ';':
			javaHead = embeddedSourceHead{}
		case tokEOF:
			if depth > 0 && len(decided) == 0 {
				return &ParseError{Message: "syntax error at end of input", Position: tok.Loc}
			}
			return nil
		case '(':
			depth++
		case ')':
			depth--
			if depth < 0 && len(decided) == 0 {
				text := ")"
				return &ParseError{Message: fmt.Sprintf("syntax error at or near %q", text), Position: tok.Loc}
			}
		}
	}
}

func (p *Parser) isIncompleteStatementEnd(stmt nodes.StmtNode) bool {
	switch p.prev.Type {
	case '(', ',', '.', '=', '+', '-', '*', '/', tokASSIGN, tokASSOC, tokCONCAT:
		return true
	}

	// Keywords Oracle does not reserve (CONTENT, JOIN, USING) are left out: a
	// statement can end with a column of that name, and a clause they open
	// without its operand ("FROM t JOIN") fails in the clause's own parser.
	switch stmt.(type) {
	case *nodes.SelectStmt:
		switch p.prev.Type {
		case kwAND, kwAS, kwBY, kwFROM, kwGROUP, kwHAVING,
			kwIS, kwNOT, kwON, kwOR, kwORDER, kwTHEN, kwUNION,
			kwWHERE:
			return true
		}
	case *nodes.InsertStmt:
		return p.prev.Type == kwINTO || p.prev.Type == kwVALUES
	case *nodes.UpdateStmt:
		return p.prev.Type == kwUPDATE || p.prev.Type == kwSET || p.prev.Type == kwWHERE
	case *nodes.DeleteStmt:
		return p.prev.Type == kwFROM || p.prev.Type == kwWHERE
	case *nodes.MergeStmt:
		switch p.prev.Type {
		case kwINTO, kwON, kwSET, kwTHEN, kwWHERE:
			return true
		}
	case *nodes.CreateTableStmt:
		return p.prev.Type == kwTABLE || p.prev.Type == kwDEFAULT
	case *nodes.CreateIndexStmt:
		return p.prev.Type == kwINDEX
	case *nodes.CreateViewStmt:
		return p.prev.Type == kwAS
	case *nodes.AlterTableStmt:
		return p.prev.Type == kwTABLE || p.prev.Type == kwADD || p.prev.Type == kwMODIFY
	case *nodes.CreateProcedureStmt, *nodes.CreateFunctionStmt, *nodes.CreatePackageStmt,
		*nodes.CreateTriggerStmt, *nodes.PLSQLBlock:
		return p.prev.Type == kwBEGIN || p.prev.Type == kwTHEN
	case *nodes.CreateUserStmt:
		return p.prev.Type == kwUSER
	case *nodes.GrantStmt:
		return p.prev.Type == kwON || p.prev.Type == kwTO
	case *nodes.DropStmt:
		return p.prev.Type == kwDROP
	}

	return false
}

func isIncompleteParsedStmt(stmt nodes.StmtNode) bool {
	switch s := stmt.(type) {
	case *nodes.SelectStmt:
		return (s.FromClause == nil || s.FromClause.Len() == 0) && selectTargetListContainsStar(s.TargetList)
	case *nodes.DropStmt:
		if s.ObjectType == nodes.OBJECT_DATABASE && (s.Names == nil || s.Names.Len() == 0) {
			return false
		}
		if s.Names == nil || s.Names.Len() == 0 {
			return true
		}
		for _, item := range s.Names.Items {
			name, ok := item.(*nodes.ObjectName)
			if ok && name.Name == "" {
				return true
			}
		}
		return false
	default:
		return false
	}
}

func selectTargetListContainsStar(list *nodes.List) bool {
	if list == nil {
		return false
	}
	for _, item := range list.Items {
		target, ok := item.(*nodes.ResTarget)
		if !ok {
			continue
		}
		if _, ok := target.Expr.(*nodes.Star); ok {
			return true
		}
	}
	return false
}

// parseStmt dispatches to statement-specific parsers.
// Statement parsers must return syntax and lexer errors directly so Parse can
// fail at the earliest invalid token while preserving partial-node Loc data.
func (p *Parser) parseStmt() (nodes.StmtNode, error) {
	switch p.cur.Type {
	case kwSELECT, kwWITH:
		return p.parseSelectStmt()
	case kwINSERT:
		return p.parseInsertStmt()
	case kwUPDATE:
		return p.parseUpdateStmt()
	case kwDELETE:
		return p.parseDeleteStmt()
	case kwMERGE:
		return p.parseMergeStmt()
	case kwDROP:
		return p.parseDropStmt()
	case kwCOMMIT:
		return p.parseCommitStmt()
	case kwROLLBACK:
		return p.parseRollbackStmt()
	case kwSAVEPOINT:
		return p.parseSavepointStmt()
	case kwSET:
		next := p.peekNext()
		if next.Type == kwTRANSACTION {
			return p.parseSetTransactionAtSet()
		}
		if next.Type == kwROLE {
			return p.parseSetRoleStmt()
		}
		if next.Type == kwCONSTRAINT || next.Type == kwCONSTRAINTS {
			return p.parseSetConstraintsStmt()
		}
		return nil, p.syntaxErrorAtCur()
	case kwAUDIT:
		return p.parseAuditStmt()
	case kwNOAUDIT:
		return p.parseNoauditStmt()
	case kwASSOCIATE:
		return p.parseAssociateStatisticsStmt()
	case kwDISASSOCIATE:
		return p.parseDisassociateStatisticsStmt()
	case kwCOMMENT:
		return p.parseCommentStmt()
	case kwTRUNCATE:
		return p.parseTruncateStmt()
	case kwANALYZE:
		return p.parseAnalyzeStmt()
	case kwEXPLAIN:
		return p.parseExplainPlanStmt()
	case kwFLASHBACK:
		next := p.peekNext()
		if next.Type == kwDATABASE ||
			(next.Type == tokIDENT && (next.Str == "STANDBY" || next.Str == "PLUGGABLE")) {
			return p.parseFlashbackDatabaseStmt()
		}
		return p.parseFlashbackTableStmt()
	case kwPURGE:
		return p.parsePurgeStmt()
	case kwLOCK:
		return p.parseLockTableStmt()
	case kwCALL:
		return p.parseCallStmt()
	case kwRENAME:
		return p.parseRenameStmt()
	case kwALTER:
		return p.parseAlterStmt()
	case kwGRANT:
		return p.parseGrantStmt()
	case kwREVOKE:
		return p.parseRevokeStmt()
	case kwCREATE:
		return p.parseCreateStmt()
	case kwDECLARE, kwBEGIN, tokLABELOPEN:
		return p.parseAnonymousBlock()
	default:
		// Handle identifier-based statements
		if p.isIdentLike() {
			switch p.cur.Str {
			case "ADMINISTER":
				return p.parseAdministerKeyManagementStmt()
			}
		}
		return nil, p.syntaxErrorAtCur()
	}
}

// advance consumes the current token and moves to the next one.
func (p *Parser) advance() Token {
	p.prev = p.cur
	if p.hasNext {
		p.cur = p.nextBuf
		p.hasNext = false
	} else {
		p.cur = p.lexer.NextToken()
	}
	return p.prev
}

// peekNext returns the next token after cur without consuming it.
func (p *Parser) peekNext() Token {
	if !p.hasNext {
		p.nextBuf = p.lexer.NextToken()
		p.hasNext = true
	}
	return p.nextBuf
}

// peekAhead returns the n tokens after cur without consuming them; tokens
// past the end of input are tokEOF. It lexes a copy of the lexer, so the
// parser's position is untouched.
func (p *Parser) peekAhead(n int) []Token {
	toks := make([]Token, 0, n)
	if n == 0 {
		return toks
	}
	toks = append(toks, p.peekNext())
	lexer := *p.lexer
	for len(toks) < n {
		toks = append(toks, lexer.NextToken())
	}
	return toks
}

// peek returns the current token without consuming it.
func (p *Parser) peek() Token {
	return p.cur
}

// match checks if the current token type matches any of the given types.
// If it matches, the token is consumed and returned with ok=true.
func (p *Parser) match(types ...int) (Token, bool) {
	for _, t := range types {
		if p.cur.Type == t {
			return p.advance(), true
		}
	}
	return Token{}, false
}

// expect consumes the current token if it matches the expected type.
// Returns an error if the token does not match.
func (p *Parser) expect(tokenType int) (Token, error) {
	if p.cur.Type == tokenType {
		return p.advance(), nil
	}
	return Token{}, p.syntaxErrorAtCur()
}

// pos returns the byte position of the current token.
func (p *Parser) pos() int {
	return p.cur.Loc
}

// isKeyword checks whether the current token is a specific keyword.
func (p *Parser) isKeyword(kw int) bool {
	return p.cur.Type == kw
}

// matchKeyword consumes the token if it is the given keyword.
func (p *Parser) matchKeyword(kw int) (Token, bool) {
	if p.cur.Type == kw {
		return p.advance(), true
	}
	return Token{}, false
}

// ParseError represents a parse error with position information.
type ParseError struct {
	Severity string // e.g., "ERROR", "WARNING"; defaults to "ERROR"
	Code     string // SQLSTATE code; defaults to "42601"
	Message  string
	Position int
}

func (e *ParseError) Error() string {
	sev := e.Severity
	if sev == "" {
		sev = "ERROR"
	}
	code := e.Code
	if code == "" {
		code = "42601"
	}
	return fmt.Sprintf("%s: %s (SQLSTATE %s)", sev, e.Message, code)
}

// tokenText returns the original source text covered by tok.
func (p *Parser) tokenText(tok Token) string {
	if tok.Type == tokEOF {
		return ""
	}
	if tok.Loc >= 0 && tok.End >= tok.Loc && tok.End <= len(p.source) {
		return p.source[tok.Loc:tok.End]
	}
	if tok.Str != "" {
		return tok.Str
	}
	if tok.Type > 0 && tok.Type < 256 {
		return string(rune(tok.Type))
	}
	return ""
}

// syntaxErrorAtCur returns a syntax error for the current token.
func (p *Parser) syntaxErrorAtCur() *ParseError {
	if p.lexer.Err != nil {
		return p.lexerError()
	}
	return p.syntaxErrorAtTok(p.cur)
}

// syntaxErrorAtTok returns a PG-style syntax error for tok.
func (p *Parser) syntaxErrorAtTok(tok Token) *ParseError {
	text := p.tokenText(tok)
	msg := "syntax error at end of input"
	if text != "" {
		msg = fmt.Sprintf("syntax error at or near %q", text)
	}
	return &ParseError{
		Message:  msg,
		Position: tok.Loc,
	}
}

// lexerError returns the lexer's error with token context.
func (p *Parser) lexerError() *ParseError {
	tok := p.cur
	if tok.Type == tokEOF && p.prev.Type != 0 {
		tok = p.prev
	}
	var directive *directiveError
	if errors.As(p.lexer.Err, &directive) {
		tok = directive.tok
	}
	text := p.tokenText(tok)
	msg := p.lexer.Err.Error()
	if text != "" {
		msg = fmt.Sprintf("%s at or near %q", msg, text)
	}
	return &ParseError{
		Message:  msg,
		Position: tok.Loc,
	}
}
