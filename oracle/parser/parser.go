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

	// parenQueries tells, for each '(' from the first one a query lookup
	// met to the end of the range, how its group reads as a query.
	// decideParenQueries fills it in one pass.
	parenQueries map[int]parenQuery

	// groupFroms tells, for each parenthesized group from the first TRIM
	// call atTrimSyntax met to the end of the range, by the offset of the
	// token after its '(', whether a FROM stands at the group's own depth.
	// decideGroupFroms fills it in one pass.
	groupFroms map[int]bool

	// queryIterand is set while the iteration control of a FOR loop that
	// opens with '(' is parsed: a query there is a cursor's, which may end
	// in a plain ORDER BY (Oracle 23ai runs FOR r IN (SELECT a FROM t ORDER
	// BY a) LOOP). parseParenExpr takes it for its own parentheses alone,
	// and records the query in iterandQuery and the token after it in
	// iterandAfter, for parseIterandExpr to check it stands alone.
	queryIterand bool
	iterandQuery *nodes.SubqueryExpr
	iterandAfter Token

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

	// addedConstraints holds the constraints of the ALTER TABLE ADD ( ... )
	// list just parsed, actions of their own after its columns'.
	addedConstraints []*nodes.AlterTableCmd

	// attributePaths records the names of each constraint column that
	// names an attribute of an object column (o.r), which a quoted column
	// name ("O.R") spells alike.
	attributePaths map[*nodes.String][]string

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

	// ccDelimiters is set when the parentheses of the input balance under
	// no single selection of conditional compilation branches the text
	// settles, so each statement checks them under the selection it parses
	// with.
	ccDelimiters bool

	// plsqlCode is set while a PL/SQL declaration or statement is parsed,
	// and cleared in the SQL statements it embeds: the expressions of
	// PL/SQL alone take the choices of a qualified expression.
	plsqlCode bool
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

// setPLSQLCode sets plsqlCode and returns its previous value, for
// restorePLSQLCode: defer p.restorePLSQLCode(p.setPLSQLCode(true)).
func (p *Parser) setPLSQLCode(plsql bool) bool {
	saved := p.plsqlCode
	p.plsqlCode = plsql
	return saved
}

// restorePLSQLCode restores plsqlCode to what setPLSQLCode returned.
func (p *Parser) restorePLSQLCode(saved bool) {
	p.plsqlCode = saved
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
	undecided, err := validateBalancedDelimiters(source, start, end)
	if err != nil {
		return nil, err
	}

	p := &Parser{
		lexer:        NewLexerRange(source, start, end),
		source:       source,
		ccDelimiters: undecided,
	}
	p.lexer.cc = &ccRecord{}
	p.lexer.ccDecisions = &[]ccDecision{}
	p.lexer.newlines = &newlineIndex{base: p.lexer.begin, limit: p.lexer.end}
	p.lexer.literalLimit = 32767
	p.lexer.numberLimit = true
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

// The branch selections parseStmtSelections tries for one statement are
// bounded: ccMinAttempts, then up to ccExhaustiveAttempts or
// ccAttemptsPerDecision for each undecided $IF an attempt met, whichever
// is more, while those attempts have read no more than ccMaxSearchBytes of
// text in all. The selections grow exponentially with the undecided $IF
// directives, so some bound is needed. Within it, a statement with up to
// ten undecided two-way $IF directives is searched exhaustively, a failure
// the search can place takes one attempt per $IF to vary, and the byte
// budget keeps the search's time linear in the text: about a second and a
// half for a statement of a thousand undecided $IF directives that no
// selection parses.
const (
	ccMinAttempts         = 64
	ccExhaustiveAttempts  = 1024
	ccAttemptsPerDecision = 2
	ccMaxSearchBytes      = 16 << 20
)

// ccSelection is a branch selection for parseStmtSelections to try: the
// choices to make, and the first decision it may vary for the selections
// after it. The ones before are its ancestors' to vary.
type ccSelection struct {
	choices *ccChoices
	from    int
}

// parseStmtSelections parses the statement at the current token, with the
// checks that follow it, under each selection of the conditional
// compilation branches whose conditions the text does not settle, until
// one parses. Oracle 23ai compiles a unit under the selection its
// environment settles (the database version, PLSQL_CCFLAGS, package
// constants), which may be any of them, so a unit fails only when every
// selection does: $IF $$x $THEN garbage $ELSE NULL; $END parses, and
// $IF $$x $THEN $ERROR 'a' $END $ELSE $ERROR 'b' $END $END fails. One
// environment gives a condition one value, so a selection decides a
// repeated condition once (ccFact).
//
// The selections form a tree, searched depth first: the first takes each
// undecided $IF's first undecided branch, and a failed one leads to those
// that pass over one more branch of one of its decisions, keeping the
// choices before it. Those that the selections of the levels holding the
// failure depend on come first, the last first, then the others, the last
// first. Every selection that keeps those decisions meets an $ERROR the
// levels hold, so a failure at one leads to none that vary only later
// decisions. Of the failures, the one that got furthest is reported.
func (p *Parser) parseStmtSelections(textStart, stmtLoc int) (nodes.StmtNode, error) {
	p.lexer.ccFacts = nil
	snap := p.snapshot()
	defer func() { p.lexer.ccChoices = nil }()
	// Each attempt records its own facts, which the lexer's copies share.
	p.lexer.ccFacts = &ccFacts{byCond: map[string]ccFact{}}
	var failure error
	furthest := -1
	pending := []ccSelection{{}}
	limit, read := ccExhaustiveAttempts, 0
	for attempt := 0; len(pending) > 0 && attempt < limit; attempt++ {
		if attempt >= ccMinAttempts && read > ccMaxSearchBytes {
			break
		}
		sel := pending[len(pending)-1]
		pending = pending[:len(pending)-1]
		if attempt > 0 {
			p.restore(snap)
			p.lexer.ccChoices = sel.choices
			p.lexer.ccFacts = &ccFacts{byCond: map[string]ccFact{}}
			// The CASE, TRIM, and parenthesized-query pre-scans read
			// another selection.
			p.caseKinds, p.groupFroms, p.parenQueries = nil, nil, nil
		}
		stmt, err := p.parseCheckedStmt(textStart, stmtLoc)
		if err == nil && p.ccDelimiters {
			err = checkSelectionDelimiters(p.source, stmtLoc, p.prev.End, sel.choices)
		}
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
		// The attempt read the text up to where it failed, which may lie
		// past the last token it parsed: comments and the code a directive
		// skips before an $ERROR.
		reached := p.prev.End
		var failed *ParseError
		if errors.As(err, &failed) {
			reached = max(reached, failed.Position)
		}
		read += max(reached-stmtLoc, 0)
		decided := (*p.lexer.ccDecisions)[snap.decisions:]
		limit = max(limit, ccMinAttempts+ccAttemptsPerDecision*len(decided))
		open := p.lexer.openDecisions()
		if de, ok := p.lexer.Err.(*directiveError); ok && de.raised {
			last := -1
			for _, at := range open {
				last = max(last, at)
			}
			n := 0
			for n < len(decided) && decided[n].at <= last {
				n++
			}
			decided = decided[:n]
		}
		pending = append(pending, nextSelections(decided, sel.from, open)...)
	}
	// Outside a PL/SQL unit, a directive is the first error, before its
	// condition and the code after it (ORA-00911 on Oracle 23ai for SELECT
	// $IF 1 $THEN 1 $END FROM dual).
	if cc := p.lexer.cc; cc != nil && len(cc.directives) > snap.directives {
		if loc := cc.directives[snap.directives]; loc >= stmtLoc && loc < furthest && !startsPreprocessedUnit(p.source, stmtLoc, p.lexer.end) {
			return nil, p.syntaxErrorAtDirective(loc)
		}
	}
	return nil, failure
}

// startsPreprocessedUnit reports whether the statement at source[start:end]
// is a PL/SQL unit Oracle preprocesses (isPreprocessedUnit), by its first
// words: an anonymous block, or CREATE of a procedure, function, package,
// trigger, or type body.
func startsPreprocessedUnit(source string, start, end int) bool {
	lex := NewLexerRange(source, start, end)
	word := func(tok Token, w string) bool { return tok.Type != tokQIDENT && tok.Str == w }
	tok := lex.nextRaw()
	switch tok.Type {
	case kwDECLARE, kwBEGIN, tokLABELOPEN:
		return true
	case kwCREATE:
	default:
		return false
	}
	tok = lex.nextRaw()
	if tok.Type == kwOR {
		lex.nextRaw() // REPLACE
		tok = lex.nextRaw()
	}
	if word(tok, "EDITIONABLE") || word(tok, "NONEDITIONABLE") {
		tok = lex.nextRaw()
	}
	switch tok.Type {
	case kwPROCEDURE, kwFUNCTION, kwPACKAGE, kwTRIGGER:
		return true
	case kwTYPE:
		return word(lex.nextRaw(), "BODY")
	}
	return false
}

// nextSelections returns the selections that follow a failed one whose
// undecided $IF directives decided records in order, from the decision
// from on: for each decision that can pass over one more undecided branch,
// the selection that does, keeping the choices before it; those after it,
// which the new choice may move or remove, start over. They come in the
// reverse of the order to try them: the decisions at offsets open, which
// the selections of the levels holding the failure depend on, are tried
// first, the last of them first, then the others, the last first.
func nextSelections(decided []ccDecision, from int, open []int) []ccSelection {
	isOpen := make(map[int]bool, len(open))
	for _, at := range open {
		isOpen[at] = true
	}
	// The decisions are recorded once, for every selection to share; the
	// lexer reuses the slice it records them in.
	decided = append([]ccDecision(nil), decided...)
	var others, held []ccSelection
	for i := from; i < len(decided); i++ {
		if !decided[i].more {
			continue
		}
		choices := &ccChoices{decided: decided[:i], at: decided[i].at, pass: decided[i].passed + 1}
		if isOpen[decided[i].at] {
			held = append(held, ccSelection{choices: choices, from: i})
		} else {
			others = append(others, ccSelection{choices: choices, from: i})
		}
	}
	return append(others, held...)
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
	if p.cur.Type == tokEOF && (p.isIncompleteStatementEnd(stmt) || p.isIncompleteParsedStmt(stmt)) && !p.failsBeforeEnd(stmt) {
		return nil, p.syntaxErrorAtCur()
	}
	if stmt == nil {
		return nil, p.syntaxErrorAtCur()
	}
	if lock := lockedSubquery(stmt); lock != nil {
		return nil, p.syntaxErrorAtNode(lock)
	}
	if sel := sourcelessStar(stmt); sel != nil {
		// At the block's start, its SELECT or its WITH, past the
		// parentheses its span takes.
		lex := p.lexer.lookahead(sel.Loc.Start, sel.Loc.End)
		tok := lex.NextToken()
		for tok.Type == '(' {
			tok = lex.NextToken()
		}
		return nil, p.syntaxErrorAtTok(tok)
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

// validateBalancedDelimiters checks that the parentheses of sql[start:end]
// balance, under the first selection of conditional compilation branches.
// Past a condition the text does not settle, they need balance under the
// selection each statement parses with, which may be another: it reports
// undecided, and parseStmtSelections checks each statement
// (checkSelectionDelimiters).
func validateBalancedDelimiters(sql string, start, end int) (undecided bool, err error) {
	var decided []ccDecision
	lexer := NewLexerRange(sql, start, end)
	lexer.ccDecisions = &decided
	err = scanDelimiters(lexer, sql, end, func() bool { return len(decided) == 0 })
	return len(decided) > 0 || lexer.Err != nil, err
}

// checkSelectionDelimiters checks that the parentheses of the statement at
// sql[start:end] balance under the selection choices names.
func checkSelectionDelimiters(sql string, start, end int, choices *ccChoices) error {
	lexer := NewLexerRange(sql, start, end)
	lexer.ccChoices = choices
	return scanDelimiters(lexer, sql, end, func() bool { return true })
}

// scanDelimiters checks the parentheses of the tokens lexer returns up to
// end, reporting an imbalance only while settled says the selection the
// lexer reads is the one to check.
func scanDelimiters(lexer *Lexer, sql string, end int, settled func() bool) error {
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
			if depth > 0 && settled() {
				return &ParseError{Message: "syntax error at end of input", Position: tok.Loc}
			}
			return nil
		case '(':
			depth++
		case ')':
			depth--
			if depth < 0 && settled() {
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

func (p *Parser) isIncompleteParsedStmt(stmt nodes.StmtNode) bool {
	switch s := stmt.(type) {
	case *nodes.SelectStmt:
		return p.endsAwaitingFrom(s)
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

// sourcelessStar returns the first query block in stmt that selects *
// without a FROM clause, which Oracle 23ai rejects with ORA-54058 at the
// block's start, its SELECT or its WITH, wherever it stands: CREATE TABLE ... AS (SELECT *), a view,
// an INSERT, a subquery, an inline view, a WITH element, MERGE USING, or a
// cursor FOR loop. The subquery of EXISTS is the exception when it is one
// query block with no row-limiting clause, as Oracle does not expand its
// select list: WHERE EXISTS (SELECT *) runs, and EXISTS (SELECT * UNION
// SELECT 1 FROM dual) and EXISTS (SELECT * FETCH FIRST 1 ROWS ONLY) fail.
func sourcelessStar(stmt nodes.Node) *nodes.SelectStmt {
	exempt := map[*nodes.SelectStmt]bool{}
	var bad *nodes.SelectStmt
	nodes.Inspect(stmt, func(n nodes.Node) bool {
		if bad != nil {
			return false
		}
		switch n := n.(type) {
		case *nodes.ExistsExpr:
			sel, ok := n.Subquery.(*nodes.SelectStmt)
			// A WITH clause before a parenthesized block wraps it, one
			// block still: Oracle 23ai runs EXISTS (WITH w AS (...)
			// (SELECT *)).
			for ok && sel.Larg != nil && sel.Op == nodes.SETOP_NONE && sel.FetchFirst == nil {
				sel = sel.Larg
			}
			if ok && sel.Larg == nil && sel.Op == nodes.SETOP_NONE && sel.FetchFirst == nil {
				exempt[sel] = true
			}
		case *nodes.SelectStmt:
			// A set operation holds its first block itself, Larg nil, or
			// wraps a parenthesized one in Larg.
			if n.Larg == nil && !exempt[n] &&
				(n.FromClause == nil || n.FromClause.Len() == 0) && selectTargetListContainsStar(n.TargetList) {
				bad = n
			}
		}
		return true
	})
	return bad
}

// lockedSubquery returns the FOR UPDATE clause that ends a subquery in
// stmt. Oracle 23ai locks rows for a top-level query and a cursor's alone:
// it rejects FOR UPDATE in an expression subquery, an inline view (LATERAL
// and APPLY ones too) or a DML target, a WITH element, a MODEL REFERENCE,
// MERGE USING, CURSOR(...), INSERT ... SELECT, and CREATE TABLE, VIEW, JSON
// DUALITY VIEW, or MATERIALIZED ZONEMAP ... AS (ORA-00907, ORA-03048), and
// runs it in a cursor declaration, OPEN FOR, and a FOR loop's query,
// VALUES OF (...) too. A subquery it does not list keeps its FOR UPDATE.
func lockedSubquery(stmt nodes.Node) nodes.Node {
	var lock nodes.Node
	check := func(query nodes.StmtNode) {
		if sel, ok := query.(*nodes.SelectStmt); ok && lock == nil {
			lock = trailingForUpdate(sel)
		}
	}
	iterands := map[*nodes.SubqueryExpr]bool{}
	nodes.Inspect(stmt, func(n nodes.Node) bool {
		if lock != nil {
			return false
		}
		switch n := n.(type) {
		case *nodes.PLSQLIterationControl:
			if sub, ok := n.Expr.(*nodes.SubqueryExpr); ok {
				iterands[sub] = true
			}
		case *nodes.SubqueryExpr:
			if !iterands[n] {
				check(n.Subquery)
			}
		case *nodes.ExistsExpr:
			check(n.Subquery)
		case *nodes.CursorExpr:
			check(n.Subquery)
		case *nodes.SubqueryRef:
			check(n.Subquery)
		case *nodes.LateralRef:
			check(n.Subquery)
		case *nodes.ModelRefModel:
			check(n.Subquery)
		case *nodes.CTE:
			check(n.Query)
		case *nodes.InsertStmt:
			if n.Select != nil {
				check(n.Select)
			}
			check(n.Subquery)
		case *nodes.CreateTableStmt:
			check(n.AsQuery)
		case *nodes.CreateViewStmt:
			check(n.Query)
		case *nodes.CreateMaterializedZonemapStmt:
			check(n.AsQuery)
		case *nodes.CreateJsonDualityViewStmt:
			check(n.Query)
		}
		return true
	})
	return lock
}

// endsAwaitingFrom reports whether the last query block of query, where
// the text ends, selects * with no FROM after it yet: SELECT * and SELECT
// 1 FROM dual UNION SELECT * end unfinished. A parenthesized block is
// whole, and sourcelessStar reports it: SELECT 1 FROM dual UNION (SELECT
// *) fails at its SELECT.
func (p *Parser) endsAwaitingFrom(query *nodes.SelectStmt) bool {
	for query.Rarg != nil {
		query = query.Rarg
	}
	if query.Larg != nil || p.parenthesizedBlock(query) {
		return false
	}
	return (query.FromClause == nil || query.FromClause.Len() == 0) && selectTargetListContainsStar(query.TargetList)
}

// parenthesizedBlock reports whether the span of the query block sel opens
// with its parentheses, as a parenthesized operand's does.
func (p *Parser) parenthesizedBlock(sel *nodes.SelectStmt) bool {
	return sel.Loc.Start >= 0 && sel.Loc.Start < len(p.source) && p.source[sel.Loc.Start] == '('
}

// failsBeforeEnd reports whether stmt holds a query block that selects *
// without FROM (sourcelessStar) other than one still open where the text
// ends, which may yet get its FROM: (SELECT *) UNION SELECT * fails at its
// first block, unfinished or not.
func (p *Parser) failsBeforeEnd(stmt nodes.StmtNode) bool {
	if stmt == nil {
		return false
	}
	bad := sourcelessStar(stmt)
	if bad == nil {
		return false
	}
	open := bad.Op == nodes.SETOP_NONE && bad.Larg == nil && !p.parenthesizedBlock(bad) && bad.Loc.End == p.prev.End
	return !open
}

// selectTargetListContainsStar reports whether list selects * or table.*:
// Oracle 23ai rejects (SELECT x.*) as it does (SELECT *) (ORA-00904), and
// runs EXISTS (SELECT x.*).
func selectTargetListContainsStar(list *nodes.List) bool {
	return list != nil && selectListHasStar(list)
}

// parseStmt dispatches to statement-specific parsers.
// Statement parsers must return syntax and lexer errors directly so Parse can
// fail at the earliest invalid token while preserving partial-node Loc data.
func (p *Parser) parseStmt() (nodes.StmtNode, error) {
	switch p.cur.Type {
	case kwSELECT, kwWITH, '(':
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
