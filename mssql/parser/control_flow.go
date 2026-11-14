// Package parser - control_flow.go implements T-SQL control flow statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseIfStmt parses an IF...ELSE statement.
//
// BNF: mssql/parser/bnf/if-else-transact-sql.bnf
//
//	IF boolean_expression
//	    { sql_statement | statement_block }
//	[ ELSE
//	    { sql_statement | statement_block } ]
//
// ifStatement (TSql160.g:13541): the condition is a boolean expression
// (procIsCondition), and ELSE takes a statement.
func (p *Parser) parseIfStmt() (*nodes.IfStmt, error) {
	loc := p.pos()
	p.advance() // consume IF

	// Completion: after IF → expression context
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		p.addRuleCandidate("func_name")
		return nil, errCollecting
	}

	stmt := &nodes.IfStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	cond, err := p.parseProcCondition("IF")
	if err != nil {
		return nil, err
	}
	stmt.Condition = cond
	then, err := p.parseStmt()
	if err != nil {
		return nil, err
	}
	if isNilNode(then) {
		return nil, p.newParseError(p.cur.Loc, "expected statement after IF condition")
	}
	stmt.Then = then

	// ELSE. Semicolons may end the THEN statement before it, as in
	// IF c SELECT 1; ELSE SELECT 2, which SQL Server accepts and a MERGE in
	// the THEN branch requires.
	p.skipSemicolonsBeforeElse()
	if _, ok := p.match(kwELSE); ok {
		elseStmt, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		if isNilNode(elseStmt) {
			return nil, p.newParseError(p.cur.Loc, "expected statement after ELSE")
		}
		stmt.Else = elseStmt
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// skipSemicolonsBeforeElse consumes the semicolons at the current token when
// ELSE follows them and reports whether it did. Otherwise it leaves the token
// stream where it was, so the semicolons end the IF statement.
func (p *Parser) skipSemicolonsBeforeElse() bool {
	if p.cur.Type != ';' {
		return false
	}
	lexPos, lexStart, lexErr := p.lexer.pos, p.lexer.start, p.lexer.Err
	cur, prev, nextBuf, hasNext, collecting := p.cur, p.prev, p.nextBuf, p.hasNext, p.collecting
	for p.cur.Type == ';' {
		p.advance()
	}
	if p.cur.Type == kwELSE {
		return true
	}
	p.lexer.pos, p.lexer.start, p.lexer.Err = lexPos, lexStart, lexErr
	p.cur, p.prev, p.nextBuf, p.hasNext, p.collecting = cur, prev, nextBuf, hasNext, collecting
	return false
}

// parseWhileStmt parses a WHILE loop.
//
// BNF: mssql/parser/bnf/while-transact-sql.bnf
//
//	WHILE boolean_expression
//	    { sql_statement | statement_block | BREAK | CONTINUE }
//
// whileStatement (TSql160.g:13594). BREAK and CONTINUE are allowed only in
// its body (procWhileDepth).
func (p *Parser) parseWhileStmt() (*nodes.WhileStmt, error) {
	loc := p.pos()
	p.advance() // consume WHILE

	// Completion: after WHILE → expression context
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		p.addRuleCandidate("func_name")
		return nil, errCollecting
	}

	stmt := &nodes.WhileStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	cond, err := p.parseProcCondition("WHILE")
	if err != nil {
		return nil, err
	}
	stmt.Condition = cond
	p.procWhileDepth++
	body, err := p.parseStmt()
	p.procWhileDepth--
	if err != nil {
		return nil, err
	}
	if isNilNode(body) {
		return nil, p.newParseError(p.cur.Loc, "expected statement after WHILE condition")
	}
	stmt.Body = body

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseBeginStmt parses BEGIN...END, BEGIN TRY...END TRY BEGIN CATCH...END CATCH,
// or BEGIN TRAN/TRANSACTION.
//
// BNF: mssql/parser/bnf/begin-end-transact-sql.bnf
//
//	BEGIN [ ; ]
//	    { sql_statement | statement_block }
//	END [ ; ]
func (p *Parser) parseBeginStmt() (nodes.StmtNode, error) {
	loc := p.pos()

	// Check for BEGIN TRY
	next := p.peekNext()
	if next.Type == kwTRY {
		stmt, err := p.parseTryCatchStmt()
		return stmt, err
	}

	// Check for BEGIN DISTRIBUTED TRAN/TRANSACTION
	if next.Type == kwDISTRIBUTED {
		stmt, err := p.parseBeginDistributedTransStmt()
		return stmt, err
	}

	// Check for BEGIN TRAN/TRANSACTION
	if next.Type == kwTRAN || next.Type == kwTRANSACTION {
		stmt, err := p.parseBeginTransStmt()
		return stmt, err
	}

	// Check for BEGIN CONVERSATION TIMER (service broker)
	if next.Type == kwCONVERSATION {
		p.advance() // consume BEGIN
		p.advance() // consume CONVERSATION
		if p.cur.Type == kwTIMER {
			p.advance() // consume TIMER
			stmt, err := p.parseBeginConversationTimerStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc = nodes.Loc{Start: loc, End: -1}
			stmt.Loc.End = p.prevEnd()
			return stmt, nil
		}
		// Not a valid statement; fall through to BEGIN...END block parsing
		// (BEGIN already consumed above, skip the p.advance() below)
		goto parseBlock
	}

	// Check for BEGIN DIALOG [CONVERSATION] (service broker)
	if next.Type == kwDIALOG {
		p.advance() // consume BEGIN
		p.advance() // consume DIALOG
		// optionally consume CONVERSATION
		if p.cur.Type == kwCONVERSATION {
			p.advance()
		}
		stmt, err := p.parseBeginConversationStmt()
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = loc
		return stmt, nil
	}

	p.advance() // consume BEGIN

	// Completion: after BEGIN → TRANSACTION, TRY, plus statement keywords
	if p.collectMode() {
		p.addTokenCandidate(kwTRANSACTION)
		p.addTokenCandidate(kwTRY)
		for _, t := range topLevelKeywords {
			p.addTokenCandidate(t)
		}
		return nil, errCollecting
	}

parseBlock:
	// Parse statements until END: at least one, and the END
	// (beginEndBlockStatement, TSql160.g:16146).
	stmts, _, err := p.parseProcStatementList()
	if err != nil {
		return nil, err
	}
	if len(stmts) == 0 {
		return nil, p.syntaxErrorAtCur()
	}
	if _, err := p.expect(kwEND); err != nil {
		return nil, err
	}

	return &nodes.BeginEndStmt{
		Stmts: &nodes.List{Items: stmts},
		Loc:   nodes.Loc{Start: loc, End: p.prevEnd()},
	}, nil
}

// parseTryCatchStmt parses BEGIN TRY...END TRY BEGIN CATCH...END CATCH.
//
// BNF: mssql/parser/bnf/try-catch-transact-sql.bnf
//
//	BEGIN TRY
//	    { sql_statement | statement_block }
//	END TRY
//	BEGIN CATCH
//	    [ { sql_statement | statement_block } ]
//	END CATCH
//	[ ; ]
//
// tryCatchStatement (TSql160.g:16066): each keyword is required, the TRY
// block holds a statement, and the CATCH block may be empty, though SQL
// Server refuses one that holds only semicolons (Msg 102). A THROW without
// arguments belongs in the CATCH block (procCatchDepth).
func (p *Parser) parseTryCatchStmt() (*nodes.TryCatchStmt, error) {
	loc := p.pos()
	p.advance() // consume BEGIN
	p.advance() // consume TRY

	stmt := &nodes.TryCatchStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Completion: after BEGIN TRY → statement keywords
	if p.collectMode() {
		for _, t := range topLevelKeywords {
			p.addTokenCandidate(t)
		}
		return nil, errCollecting
	}

	// TRY block
	tryStmts, _, err := p.parseProcStatementList()
	if err != nil {
		return nil, err
	}
	if len(tryStmts) == 0 {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.TryBlock = &nodes.List{Items: tryStmts}

	// END TRY BEGIN CATCH
	for _, kw := range []int{kwEND, kwTRY, kwBEGIN, kwCATCH} {
		if _, err := p.expect(kw); err != nil {
			return nil, err
		}
	}

	// Completion: after BEGIN CATCH → statement keywords
	if p.collectMode() {
		for _, t := range topLevelKeywords {
			p.addTokenCandidate(t)
		}
		return nil, errCollecting
	}

	// CATCH block
	p.procCatchDepth++
	catchStmts, semicolon, err := p.parseProcStatementList()
	p.procCatchDepth--
	if err != nil {
		return nil, err
	}
	if len(catchStmts) == 0 && semicolon {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.CatchBlock = &nodes.List{Items: catchStmts}

	// END CATCH
	if _, err := p.expect(kwEND); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwCATCH); err != nil {
		return nil, err
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseReturnStmt parses a RETURN statement.
//
// BNF: mssql/parser/bnf/return-transact-sql.bnf
//
//	RETURN [ integer_expression ]
//
// returnStatement (TSql160.g:14372). SQL Server takes a value only in a
// procedure or function body (Msg 178 elsewhere, triggers included).
func (p *Parser) parseReturnStmt() (*nodes.ReturnStmt, error) {
	loc := p.pos()
	p.advance() // consume RETURN

	// Completion: after RETURN → expression context
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		p.addRuleCandidate("func_name")
		return nil, errCollecting
	}

	stmt := &nodes.ReturnStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Optional return expression
	if p.cur.Type != ';' && p.cur.Type != tokEOF && p.cur.Type != kwEND &&
		p.cur.Type != kwELSE && p.cur.Type != kwGO {
		valLoc := p.pos()
		var err error
		stmt.Value, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
		if stmt.Value != nil && p.procRoutineDepth == 0 {
			return nil, p.newParseError(valLoc, "a RETURN statement with a return value cannot be used in this context")
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseBreakStmt parses a BREAK statement.
//
// BNF: mssql/parser/bnf/break-transact-sql.bnf
//
//	BREAK
//
// SQL Server refuses BREAK outside a WHILE (Msg 135).
func (p *Parser) parseBreakStmt() (*nodes.BreakStmt, error) {
	loc := p.pos()
	if p.procWhileDepth == 0 {
		return nil, p.newParseError(loc, "cannot use a BREAK statement outside the scope of a WHILE statement")
	}
	p.advance() // consume BREAK
	return &nodes.BreakStmt{
		Loc: nodes.Loc{Start: loc, End: p.prevEnd()},
	}, nil
}

// parseContinueStmt parses a CONTINUE statement.
//
// BNF: mssql/parser/bnf/continue-transact-sql.bnf
//
//	CONTINUE
//
// SQL Server refuses CONTINUE outside a WHILE (Msg 136).
func (p *Parser) parseContinueStmt() (*nodes.ContinueStmt, error) {
	loc := p.pos()
	if p.procWhileDepth == 0 {
		return nil, p.newParseError(loc, "cannot use a CONTINUE statement outside the scope of a WHILE statement")
	}
	p.advance() // consume CONTINUE
	return &nodes.ContinueStmt{
		Loc: nodes.Loc{Start: loc, End: p.prevEnd()},
	}, nil
}

// parseGotoStmt parses a GOTO label statement.
//
// BNF: mssql/parser/bnf/goto-transact-sql.bnf
//
//	GOTO label
//
// gotoStatement (TSql160.g:15847): the label is an undelimited identifier.
func (p *Parser) parseGotoStmt() (*nodes.GotoStmt, error) {
	loc := p.pos()
	p.advance() // consume GOTO

	if !p.isIdentLike() || p.isDelimited(p.cur) {
		return nil, p.syntaxErrorAtCur()
	}
	label := p.cur.Str
	p.advance()

	return &nodes.GotoStmt{
		Label: label,
		Loc:   nodes.Loc{Start: loc, End: p.prevEnd()},
	}, nil
}

// parseWaitForStmt parses a WAITFOR statement.
//
// BNF: mssql/parser/bnf/waitfor-transact-sql.bnf
//
//	WAITFOR
//	{
//	    DELAY 'time_to_pass'
//	  | TIME 'time_to_execute'
//	  | [ ( receive_statement ) | ( get_conversation_group_statement ) ]
//	    [ , TIMEOUT timeout ]
//	}
func (p *Parser) parseWaitForStmt() (*nodes.WaitForStmt, error) {
	loc := p.pos()
	p.advance() // consume WAITFOR

	// Completion: after WAITFOR → DELAY, TIME
	if p.collectMode() {
		p.addTokenCandidate(kwDELAY)
		p.addTokenCandidate(kwTIME)
		return nil, errCollecting
	}

	stmt := &nodes.WaitForStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	if p.cur.Type == kwDELAY {
		stmt.WaitType = "DELAY"
		p.advance()
		var err error
		stmt.Value, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
	} else if p.cur.Type == kwTIME {
		stmt.WaitType = "TIME"
		p.advance()
		var err error
		stmt.Value, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
	} else if p.cur.Type == '(' {
		// Parenthesized form: WAITFOR ( receive_statement | get_conversation_group_statement )
		// TSql160.g:14175 waitForStatement, 14207 waitforInnerStatement: the
		// statement is RECEIVE or GET CONVERSATION GROUP, and the timeout
		// after the comma is signedIntegerOrVariable.
		p.advance() // consume '('
		if p.cur.Type != kwRECEIVE && !(p.cur.Type == kwGET && p.peekNext().Type == kwCONVERSATION) {
			return nil, p.syntaxErrorAtCur()
		}
		innerStmt, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		stmt.InnerStmt = innerStmt
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		if _, ok := p.match(','); ok {
			if _, err := p.expect(kwTIMEOUT); err != nil {
				return nil, err
			}
			signed := (p.cur.Type == '-' || p.cur.Type == '+') && p.peekNext().Type == tokICONST
			if !signed && p.cur.Type != tokICONST && p.cur.Type != tokVARIABLE {
				return nil, p.syntaxErrorAtCur()
			}
			stmt.Timeout, err = p.parseUnary()
			if err != nil {
				return nil, err
			}
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseProcCondition parses the boolean expression of IF or WHILE.
func (p *Parser) parseProcCondition(what string) (nodes.ExprNode, error) {
	p.enterSearchCondition()
	condLoc := p.pos()
	cond, err := p.parseExpr()
	p.leaveSearchCondition()
	if err != nil {
		return nil, err
	}
	if cond == nil {
		return nil, p.newParseError(p.cur.Loc, "expected condition after "+what)
	}
	if !procIsCondition(cond) {
		return nil, p.newParseError(condLoc, "an expression of non-boolean type specified in a context where a condition is expected")
	}
	return cond, nil
}

// procIsCondition reports whether e is a boolean expression: a comparison, a
// predicate, AND, OR or NOT of them, or one in parentheses. SQL Server
// refuses any other expression where a condition is expected (Msg 4145),
// so IF (a) and IF (SELECT 1) are errors. UPDATE (column) in a trigger and
// TSEQUAL are the functions that are predicates, and omni reads a graph
// MATCH, which it does not parse yet, as a call of MATCH.
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/queries/search-condition-transact-sql
func procIsCondition(e nodes.ExprNode) bool {
	switch n := e.(type) {
	case *nodes.BinaryExpr:
		switch n.Op {
		case nodes.BinOpEq, nodes.BinOpNeq, nodes.BinOpLt, nodes.BinOpGt,
			nodes.BinOpLte, nodes.BinOpGte, nodes.BinOpNotLt, nodes.BinOpNotGt,
			nodes.BinOpAnd, nodes.BinOpOr:
			return true
		}
		return false
	case *nodes.UnaryExpr:
		return n.Op == nodes.UnaryNot
	case *nodes.ParenExpr:
		return procIsCondition(n.Expr)
	case *nodes.BetweenExpr, *nodes.InExpr, *nodes.LikeExpr, *nodes.IsExpr,
		*nodes.ExistsExpr, *nodes.FullTextPredicate, *nodes.SubqueryComparisonExpr:
		return true
	case *nodes.FuncCallExpr:
		if n.Name == nil || n.Name.Schema != "" || n.Name.Database != "" || n.Name.Server != "" {
			return false
		}
		switch strings.ToUpper(n.Name.Object) {
		case "UPDATE", "TSEQUAL", "MATCH":
			return true
		}
	}
	return false
}

// parseProcStatementList parses the statements of a BEGIN...END, TRY or CATCH
// block up to the END that closes it, with semicolons between them. It
// reports whether it skipped a semicolon.
func (p *Parser) parseProcStatementList() ([]nodes.Node, bool, error) {
	var stmts []nodes.Node
	semicolon := false
	for p.cur.Type != kwEND && p.cur.Type != tokEOF {
		if p.cur.Type == ';' {
			semicolon = true
			p.advance()
			continue
		}
		s, err := p.parseStmt()
		if err != nil {
			return nil, false, err
		}
		if isNilNode(s) {
			break
		}
		stmts = append(stmts, s)
	}
	return stmts, semicolon, nil
}

// procIsWord reports whether the current token is the word w (upper case):
// the keyword w, or an undelimited identifier spelled w when w is not a
// keyword.
func (p *Parser) procIsWord(w string) bool {
	if kw := lookupKeyword(w); kw != tokIDENT {
		return p.cur.Type == kw
	}
	return p.cur.Type == tokIDENT && !p.isDelimited(p.cur) && strings.ToUpper(p.cur.Str) == w
}

// procExpectWord consumes the word w or fails at the current token.
func (p *Parser) procExpectWord(w string) error {
	if !p.procIsWord(w) {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

// procIsBinary reports whether a numeric token is a binary constant (0x...).
func procIsBinary(tok Token) bool {
	return len(tok.Str) > 1 && tok.Str[0] == '0' && (tok.Str[1] == 'x' || tok.Str[1] == 'X')
}

// procIsInteger reports whether tok is an integer constant, not a binary one.
func procIsInteger(tok Token) bool {
	return tok.Type == tokICONST && !procIsBinary(tok)
}

// procIsBinaryConst reports whether tok is a binary constant.
func procIsBinaryConst(tok Token) bool {
	return (tok.Type == tokICONST || tok.Type == tokFCONST) && procIsBinary(tok)
}

// procIsString reports whether tok is a string literal.
func procIsString(tok Token) bool {
	return tok.Type == tokSCONST || tok.Type == tokNSCONST
}

// procValueKinds selects the constants parseProcValue takes, as the
// grammar's small value rules (signedIntegerOrVariable, stringOrVariable,
// binaryOrVariable, ...) do: a signed integer, an integer, a string, a binary
// constant, a local variable, NULL.
type procValueKinds struct {
	signedInt, integer, str, binary, variable, null bool
}

// parseProcValue parses a value of the given kinds or fails at the current
// token.
func (p *Parser) parseProcValue(k procValueKinds) (nodes.ExprNode, error) {
	switch {
	case k.signedInt && p.cur.Type == '-':
		loc := p.pos()
		p.advance() // consume '-'
		if !procIsInteger(p.cur) {
			return nil, p.syntaxErrorAtCur()
		}
		operand, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}
		return &nodes.UnaryExpr{Op: nodes.UnaryMinus, Operand: operand, Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil
	case (k.signedInt || k.integer) && procIsInteger(p.cur),
		k.str && procIsString(p.cur),
		k.binary && procIsBinaryConst(p.cur),
		k.variable && p.cur.Type == tokVARIABLE,
		k.null && p.cur.Type == kwNULL:
		return p.parsePrimary()
	}
	return nil, p.syntaxErrorAtCur()
}

// parseProcConstantOrIdentifier parses possibleNegativeConstantOrIdentifier
// (TSql160.g:25150): a literal (TSql160.g:33382), a minus sign before a
// number or money constant, or an identifier.
func (p *Parser) parseProcConstantOrIdentifier() (nodes.ExprNode, error) {
	switch p.cur.Type {
	case tokICONST, tokFCONST, tokSCONST, tokNSCONST, tokMONEY,
		kwNULL, tokVARIABLE, tokSYSVARIABLE:
		return p.parsePrimary()
	case '-':
		loc := p.pos()
		p.advance() // consume '-'
		if (p.cur.Type != tokICONST && p.cur.Type != tokFCONST && p.cur.Type != tokMONEY) || procIsBinary(p.cur) {
			return nil, p.syntaxErrorAtCur()
		}
		operand, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}
		return &nodes.UnaryExpr{Op: nodes.UnaryMinus, Operand: operand, Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil
	}
	if p.isIdentLike() {
		tok := p.advance()
		return &nodes.ColumnRef{Column: tok.Str, Loc: nodes.Loc{Start: tok.Loc, End: p.prevEnd()}}, nil
	}
	return nil, p.syntaxErrorAtCur()
}

// parseProcCursorRef parses cursorId (TSql160.g:14652): GLOBAL and a name, a
// name, or a cursor variable, one of them required.
func (p *Parser) parseProcCursorRef() (string, bool, error) {
	if p.cur.Type == tokVARIABLE {
		name := p.cur.Str
		p.advance()
		return name, false, nil
	}
	global := false
	if p.cur.Type == kwGLOBAL && p.peekNext().Type != tokEOF && p.peekNext().Type != ';' {
		p.advance() // consume GLOBAL
		global = true
	}
	name, ok := p.parseIdentifier()
	if !ok {
		return "", false, p.syntaxErrorAtCur()
	}
	return name, global, nil
}
