package parser

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/mysql/ast"
)

// parseBeginEndBlock parses a BEGIN...END compound statement block.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/begin-end.html
//
//	[begin_label:] BEGIN
//	    [statement_list]
//	END [end_label]
//
// The label, if any, must already be consumed by the caller and passed as labelName.
func (p *Parser) parseBeginEndBlock(labelName string, labelStart int) (*nodes.BeginEndBlock, error) {
	start := labelStart
	if labelName == "" {
		start = p.pos()
	}

	p.advance() // consume BEGIN

	stmt := &nodes.BeginEndBlock{
		Loc:   nodes.Loc{Start: start},
		Label: labelName,
	}

	// Register the label in the *current (outer)* scope so LEAVE from
	// within the body can find it via the parent chain.
	if labelName != "" {
		if err := p.declareLabel(labelName, labelBegin, stmt, labelStart); err != nil {
			return nil, err
		}
	}
	// Open a child scope for the body's declarations.
	if p.procScope != nil {
		p.pushScope(scopeBlock)
		defer p.popScope()
	}

	// Parse statements until END, enforcing MySQL's declaration-ordering rule:
	//   DECLARE var/condition  →  DECLARE cursor  →  DECLARE handler  →  stmts
	// Each declaration kind may move the phase forward but not backward.
	// Container-verified against MySQL 8.0 (2026-04-20); mirrors ERR 1337/1338.
	const (
		phaseInit = iota
		phaseSawVar
		phaseSawCursor
		phaseSawHandler
		phaseSawStmt
	)
	phase := phaseInit
	for p.cur.Type != kwEND && p.cur.Type != tokEOF {
		sStart := p.pos()
		s, err := p.parseCompoundStmtOrStmt()
		if err != nil {
			return nil, err
		}
		if err := checkDeclarePhase(&phase, s, sStart); err != nil {
			return nil, err
		}
		stmt.Stmts = append(stmt.Stmts, s)

		// Consume optional semicolon between statements
		if p.cur.Type == ';' {
			p.advance()
		}
	}

	// Consume END
	if _, err := p.expect(kwEND); err != nil {
		return nil, err
	}

	// Optional end_label. Must match the begin label (case-insensitive, per
	// MySQL 8.0 — container-verified 2026-04-20).
	if p.isLabelIdentToken() && p.cur.Type != tokEOF {
		endLabelTok := p.advance()
		stmt.EndLabel = endLabelTok.Str
		if err := checkLabelMatch(labelName, stmt.EndLabel, endLabelTok.Loc); err != nil {
			return nil, err
		}
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// checkDeclarePhase enforces MySQL's DECLARE ordering rule inside a
// BEGIN...END block. See parseBeginEndBlock for the phase state machine.
// MySQL surfaces ERR 1337 / 1338; omni uses its own wording.
func checkDeclarePhase(phase *int, s nodes.Node, sPos int) error {
	const (
		phaseInit = iota
		phaseSawVar
		phaseSawCursor
		phaseSawHandler
		phaseSawStmt
	)
	switch s.(type) {
	case *nodes.DeclareVarStmt, *nodes.DeclareConditionStmt:
		switch {
		case *phase == phaseSawStmt:
			return &ParseError{
				Message:  "variable or condition declaration after regular statement",
				Position: sPos,
			}
		case *phase > phaseSawVar:
			return &ParseError{
				Message:  "variable or condition declaration after cursor or handler declaration",
				Position: sPos,
			}
		}
		*phase = phaseSawVar
	case *nodes.DeclareCursorStmt:
		switch {
		case *phase == phaseSawStmt:
			return &ParseError{
				Message:  "cursor declaration after regular statement",
				Position: sPos,
			}
		case *phase > phaseSawCursor:
			return &ParseError{
				Message:  "cursor declaration after handler declaration",
				Position: sPos,
			}
		}
		*phase = phaseSawCursor
	case *nodes.DeclareHandlerStmt:
		if *phase > phaseSawHandler {
			return &ParseError{
				Message:  "handler declaration after regular statement",
				Position: sPos,
			}
		}
		*phase = phaseSawHandler
	default:
		*phase = phaseSawStmt
	}
	return nil
}

// containsReturn reports whether a stored-function body's AST contains at
// least one *ReturnStmt anywhere. This mirrors MySQL 8.0's actual CREATE-
// time check (sp_head sets HAS_RETURN flag when a RETURN is parsed; at
// CREATE the parser raises ERR 1320 "No RETURN found in FUNCTION" iff the
// flag is unset). MySQL does NOT do path analysis at CREATE — functions
// like `IF x THEN RETURN 1; END IF` (no else-RETURN) are accepted and the
// missing-return condition surfaces only at runtime if the path is taken.
// SIGNAL/RESIGNAL do NOT substitute for RETURN in MySQL's check.
//
// Container-verified 2026-04-21 via TestRoutineAlignment.
func containsReturn(s nodes.Node) bool {
	if s == nil {
		return false
	}
	switch n := s.(type) {
	case *nodes.ReturnStmt:
		return true
	case *nodes.BeginEndBlock:
		return containsReturnList(n.Stmts)
	case *nodes.IfStmt:
		if containsReturnList(n.ThenList) {
			return true
		}
		for _, ei := range n.ElseIfs {
			if containsReturnList(ei.ThenList) {
				return true
			}
		}
		return containsReturnList(n.ElseList)
	case *nodes.CaseStmtNode:
		for _, w := range n.Whens {
			if containsReturnList(w.ThenList) {
				return true
			}
		}
		return containsReturnList(n.ElseList)
	case *nodes.WhileStmt:
		return containsReturnList(n.Stmts)
	case *nodes.LoopStmt:
		return containsReturnList(n.Stmts)
	case *nodes.RepeatStmt:
		return containsReturnList(n.Stmts)
	case *nodes.DeclareHandlerStmt:
		return containsReturn(n.Stmt)
	default:
		return false
	}
}

func containsReturnList(stmts []nodes.Node) bool {
	for _, s := range stmts {
		if containsReturn(s) {
			return true
		}
	}
	return false
}

// checkLabelMatch enforces MySQL's end-label matching rule: when an end label
// is present, a matching begin label must also be present and the two names
// must be equal case-insensitively. MySQL surfaces ERR 1310 ("End-label X
// without match"); omni uses its own wording.
func checkLabelMatch(beginLabel, endLabel string, endLabelPos int) error {
	if endLabel == "" {
		// No end label: always fine.
		return nil
	}
	if beginLabel == "" {
		return &ParseError{
			Message:  fmt.Sprintf("end label %q without matching begin label", endLabel),
			Position: endLabelPos,
		}
	}
	if !strings.EqualFold(beginLabel, endLabel) {
		return &ParseError{
			Message:  fmt.Sprintf("end label %q does not match begin label %q", endLabel, beginLabel),
			Position: endLabelPos,
		}
	}
	return nil
}

// parseCompoundStmtOrStmt tries to parse compound statements (DECLARE, labeled blocks,
// flow control) and falls back to regular statement parsing.
func (p *Parser) parseCompoundStmtOrStmt() (nodes.Node, error) {
	// DECLARE
	if p.cur.Type == kwDECLARE {
		return p.parseDeclareStmt()
	}

	// BEGIN ... END (without label)
	if p.cur.Type == kwBEGIN {
		return p.parseBeginEndBlock("", 0)
	}

	// IF statement
	if p.cur.Type == kwIF {
		return p.parseIfStmt()
	}

	// CASE statement
	if p.cur.Type == kwCASE {
		return p.parseCaseStmt()
	}

	// WHILE loop (without label)
	if p.cur.Type == kwWHILE {
		return p.parseWhileStmt("", 0)
	}

	// REPEAT loop (without label)
	if p.cur.Type == kwREPEAT {
		return p.parseRepeatStmt("", 0)
	}

	// LOOP (without label)
	if p.cur.Type == kwLOOP {
		return p.parseLoopStmt("", 0)
	}

	// LEAVE
	if p.cur.Type == kwLEAVE {
		return p.parseLeaveStmt()
	}

	// ITERATE
	if p.cur.Type == kwITERATE {
		return p.parseIterateStmt()
	}

	// RETURN
	if p.cur.Type == kwRETURN {
		return p.parseReturnStmt()
	}

	// OPEN cursor
	if p.cur.Type == kwOPEN {
		return p.parseOpenCursorStmt()
	}

	// FETCH cursor
	if p.cur.Type == kwFETCH {
		return p.parseFetchCursorStmt()
	}

	// CLOSE cursor
	if p.cur.Type == kwCLOSE {
		return p.parseCloseCursorStmt()
	}

	// Label: check for identifier followed by ':'
	if p.isLabelIdentToken() && p.cur.Type != kwEND {
		next := p.peekNext()
		if next.Type == ':' {
			name := p.cur.Str
			nameStart := p.pos()
			p.advance() // consume identifier
			p.advance() // consume ':'
			// Labeled BEGIN...END
			if p.cur.Type == kwBEGIN {
				return p.parseBeginEndBlock(name, nameStart)
			}
			// Labeled WHILE
			if p.cur.Type == kwWHILE {
				return p.parseWhileStmt(name, nameStart)
			}
			// Labeled REPEAT
			if p.cur.Type == kwREPEAT {
				return p.parseRepeatStmt(name, nameStart)
			}
			// Labeled LOOP
			if p.cur.Type == kwLOOP {
				return p.parseLoopStmt(name, nameStart)
			}
			return nil, &ParseError{
				Message:  "expected BEGIN, WHILE, REPEAT, or LOOP after label",
				Position: p.cur.Loc,
			}
		}
	}

	// Regular statement
	return p.parseStmt()
}

// parseDeclareStmt parses DECLARE statements inside compound blocks.
// Dispatches to variable, condition, handler, or cursor declarations.
//
// DECLARE var [, var] ... type [DEFAULT value]
// DECLARE condition_name CONDITION FOR condition_value
// DECLARE handler_action HANDLER FOR condition_value [, ...] stmt
// DECLARE cursor_name CURSOR FOR select_stmt
func (p *Parser) parseDeclareStmt() (nodes.Node, error) {
	start := p.pos()
	p.advance() // consume DECLARE

	// Check for handler: CONTINUE | EXIT | UNDO
	if p.cur.Type == kwCONTINUE || p.cur.Type == kwEXIT || p.cur.Type == kwUNDO {
		return p.parseDeclareHandlerStmt(start)
	}

	// We need to look ahead to determine if it's CONDITION, CURSOR, or variable.
	// Parse the first identifier name.
	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	// DECLARE name CONDITION FOR ...
	if p.cur.Type == kwCONDITION {
		return p.parseDeclareConditionStmt(start, name)
	}

	// DECLARE name CURSOR FOR ...
	if p.cur.Type == kwCURSOR {
		return p.parseDeclareCursorStmt(start, name)
	}

	// DECLARE var [, var] ... type [DEFAULT value]
	return p.parseDeclareVarStmt(start, name)
}

// parseDeclareVarStmt parses a DECLARE variable statement.
// The first variable name has already been consumed.
//
//	DECLARE var_name [, var_name] ... type [DEFAULT value]
func (p *Parser) parseDeclareVarStmt(start int, firstName string) (*nodes.DeclareVarStmt, error) {
	stmt := &nodes.DeclareVarStmt{
		Loc:   nodes.Loc{Start: start},
		Names: []string{firstName},
	}

	// Additional variable names separated by commas
	for p.cur.Type == ',' {
		p.advance()
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Names = append(stmt.Names, name)
	}

	// Data type
	dt, err := p.parseDataType()
	if err != nil {
		return nil, err
	}
	stmt.TypeName = dt

	// Optional DEFAULT value
	if p.cur.Type == kwDEFAULT {
		p.advance()
		val, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if !p.collectMode() && val == nil {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Default = val
	}

	stmt.Loc.End = p.pos()
	// Static-validation insert: each declared name claims the var slot in
	// the current scope. Duplicate within same scope (or shadowing of a
	// param at the top scope) is rejected per ER_SP_DUP_VAR.
	for _, n := range stmt.Names {
		if err := p.declareVar(n, stmt, start); err != nil {
			return nil, err
		}
	}
	return stmt, nil
}

// parseDeclareConditionStmt parses DECLARE condition_name CONDITION FOR condition_value.
// The condition name has already been consumed.
//
//	DECLARE condition_name CONDITION FOR condition_value
//	condition_value: SQLSTATE [VALUE] sqlstate_value | mysql_error_code
func (p *Parser) parseDeclareConditionStmt(start int, name string) (*nodes.DeclareConditionStmt, error) {
	p.advance() // consume CONDITION

	// FOR
	if _, err := p.expect(kwFOR); err != nil {
		return nil, err
	}

	stmt := &nodes.DeclareConditionStmt{
		Loc:  nodes.Loc{Start: start},
		Name: name,
	}

	// condition_value
	condVal, err := p.parseHandlerConditionValue()
	if err != nil {
		return nil, err
	}
	stmt.ConditionValue = condVal

	stmt.Loc.End = p.pos()
	if err := p.declareCondition(name, stmt, start); err != nil {
		return nil, err
	}
	return stmt, nil
}

// parseDeclareHandlerStmt parses DECLARE handler_action HANDLER FOR condition_value [, ...] stmt.
//
//	DECLARE {CONTINUE|EXIT|UNDO} HANDLER FOR condition_value [, condition_value] ... statement
func (p *Parser) parseDeclareHandlerStmt(start int) (*nodes.DeclareHandlerStmt, error) {
	stmt := &nodes.DeclareHandlerStmt{
		Loc: nodes.Loc{Start: start},
	}

	// handler_action
	switch p.cur.Type {
	case kwCONTINUE:
		stmt.Action = "CONTINUE"
	case kwEXIT:
		stmt.Action = "EXIT"
	case kwUNDO:
		stmt.Action = "UNDO"
	}
	p.advance()

	// HANDLER
	if _, err := p.expect(kwHANDLER); err != nil {
		return nil, err
	}

	// FOR
	if _, err := p.expect(kwFOR); err != nil {
		return nil, err
	}

	// condition_value [, condition_value] ...
	for {
		condVal, err := p.parseHandlerConditionValue()
		if err != nil {
			return nil, err
		}
		stmt.Conditions = append(stmt.Conditions, condVal)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	// Handler body lives in a HANDLER_SCOPE: outer vars/conditions/cursors
	// are visible via parent chain, but outer labels are not (label barrier).
	if p.procScope != nil {
		p.pushScope(scopeHandlerBody)
		defer p.popScope()
	}
	body, err := p.parseCompoundStmtOrStmt()
	if err != nil {
		return nil, err
	}
	stmt.Stmt = body

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseDeclareCursorStmt parses DECLARE cursor_name CURSOR FOR select_statement.
// The cursor name has already been consumed.
//
//	DECLARE cursor_name CURSOR FOR select_statement
func (p *Parser) parseDeclareCursorStmt(start int, name string) (*nodes.DeclareCursorStmt, error) {
	p.advance() // consume CURSOR

	// FOR
	if _, err := p.expect(kwFOR); err != nil {
		return nil, err
	}

	// select_statement
	sel, err := p.parseSelectStmt()
	if err != nil {
		return nil, err
	}

	stmt := &nodes.DeclareCursorStmt{
		Loc:    nodes.Loc{Start: start, End: p.pos()},
		Name:   name,
		Select: sel,
	}
	if err := p.declareCursor(name, stmt, start); err != nil {
		return nil, err
	}
	return stmt, nil
}

// parseHandlerConditionValue parses a handler condition value, returning
// a tagged HandlerCondValue so semantic analysis can distinguish SQLSTATE
// literals from user condition names with the same text.
//
//	condition_value:
//	    mysql_error_code
//	  | SQLSTATE [VALUE] sqlstate_value
//	  | condition_name
//	  | SQLWARNING
//	  | NOT FOUND
//	  | SQLEXCEPTION
func (p *Parser) parseHandlerConditionValue() (nodes.HandlerCondValue, error) {
	// SQLSTATE [VALUE] 'sqlstate_value'
	if p.cur.Type == kwSQLSTATE {
		p.advance()
		// Optional VALUE
		if p.cur.Type == kwVALUE {
			p.advance()
		}
		if p.cur.Type != tokSCONST {
			return nodes.HandlerCondValue{}, &ParseError{
				Message:  "expected SQLSTATE value string",
				Position: p.cur.Loc,
			}
		}
		val := p.cur.Str
		p.advance()
		return nodes.HandlerCondValue{Kind: nodes.HandlerCondSQLState, Value: val}, nil
	}

	// SQLWARNING
	if p.cur.Type == kwSQLWARNING {
		p.advance()
		return nodes.HandlerCondValue{Kind: nodes.HandlerCondSQLWarning}, nil
	}

	// SQLEXCEPTION
	if p.cur.Type == kwSQLEXCEPTION {
		p.advance()
		return nodes.HandlerCondValue{Kind: nodes.HandlerCondSQLException}, nil
	}

	// NOT FOUND
	if p.cur.Type == kwNOT {
		p.advance()
		if _, err := p.expect(kwFOUND); err != nil {
			return nodes.HandlerCondValue{}, err
		}
		return nodes.HandlerCondValue{Kind: nodes.HandlerCondNotFound}, nil
	}

	// mysql_error_code (integer literal)
	if p.cur.Type == tokICONST {
		val := p.cur.Str
		p.advance()
		return nodes.HandlerCondValue{Kind: nodes.HandlerCondErrorCode, Value: val}, nil
	}

	// condition_name (identifier)
	name, _, err := p.parseIdentifier()
	if err != nil {
		return nodes.HandlerCondValue{}, err
	}
	return nodes.HandlerCondValue{Kind: nodes.HandlerCondName, Value: name}, nil
}

// isCompoundTerminator checks if the current token is a keyword that terminates
// a statement list inside a compound block (END, ELSEIF, ELSE, UNTIL, WHEN).
func (p *Parser) isCompoundTerminator() bool {
	switch p.cur.Type {
	case kwEND, kwELSEIF, kwELSE, kwUNTIL, kwWHEN, tokEOF:
		return true
	}
	return false
}

// parseCompoundStmtList parses a list of statements inside a compound block,
// stopping at any compound terminator keyword.
func (p *Parser) parseCompoundStmtList() ([]nodes.Node, error) {
	var stmts []nodes.Node
	for !p.isCompoundTerminator() {
		s, err := p.parseCompoundStmtOrStmt()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, s)
		// Consume optional semicolon between statements
		if p.cur.Type == ';' {
			p.advance()
		}
	}
	return stmts, nil
}

// parseIfStmt parses an IF/ELSEIF/ELSE/END IF compound statement.
//
//	IF search_condition THEN statement_list
//	  [ELSEIF search_condition THEN statement_list] ...
//	  [ELSE statement_list]
//	END IF
func (p *Parser) parseIfStmt() (*nodes.IfStmt, error) {
	start := p.pos()
	p.advance() // consume IF

	// Parse condition
	cond, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if cond == nil {
		return nil, p.syntaxErrorAtCur()
	}

	// THEN
	if _, err := p.expect(kwTHEN); err != nil {
		return nil, err
	}

	// Parse THEN statement list
	thenList, err := p.parseCompoundStmtList()
	if err != nil {
		return nil, err
	}

	stmt := &nodes.IfStmt{
		Loc:      nodes.Loc{Start: start},
		Cond:     cond,
		ThenList: thenList,
	}

	// Parse ELSEIF clauses
	for p.cur.Type == kwELSEIF {
		eiStart := p.pos()
		p.advance() // consume ELSEIF

		eiCond, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if eiCond == nil {
			return nil, p.syntaxErrorAtCur()
		}

		if _, err := p.expect(kwTHEN); err != nil {
			return nil, err
		}

		eiThenList, err := p.parseCompoundStmtList()
		if err != nil {
			return nil, err
		}

		stmt.ElseIfs = append(stmt.ElseIfs, &nodes.ElseIf{
			Loc:      nodes.Loc{Start: eiStart, End: p.pos()},
			Cond:     eiCond,
			ThenList: eiThenList,
		})
	}

	// Parse optional ELSE clause
	if p.cur.Type == kwELSE {
		p.advance() // consume ELSE

		elseList, err := p.parseCompoundStmtList()
		if err != nil {
			return nil, err
		}
		stmt.ElseList = elseList
	}

	// END IF
	if _, err := p.expect(kwEND); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwIF); err != nil {
		return nil, err
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseCaseStmt parses a CASE statement (not expression) in stored programs.
//
//	CASE case_value
//	  WHEN when_value THEN statement_list
//	  [WHEN when_value THEN statement_list] ...
//	  [ELSE statement_list]
//	END CASE
//
//	CASE
//	  WHEN search_condition THEN statement_list
//	  [WHEN search_condition THEN statement_list] ...
//	  [ELSE statement_list]
//	END CASE
func (p *Parser) parseCaseStmt() (*nodes.CaseStmtNode, error) {
	start := p.pos()
	p.advance() // consume CASE

	stmt := &nodes.CaseStmtNode{
		Loc: nodes.Loc{Start: start},
	}

	// Simple CASE: if next token is not WHEN, parse operand
	if p.cur.Type != kwWHEN {
		operand, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if !p.collectMode() && operand == nil {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Operand = operand
	}

	// Parse WHEN clauses
	for p.cur.Type == kwWHEN {
		whenStart := p.pos()
		p.advance() // consume WHEN

		whenCond, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if whenCond == nil {
			return nil, p.syntaxErrorAtCur()
		}

		if _, err := p.expect(kwTHEN); err != nil {
			return nil, err
		}

		whenThenList, err := p.parseCompoundStmtList()
		if err != nil {
			return nil, err
		}

		stmt.Whens = append(stmt.Whens, &nodes.CaseStmtWhen{
			Loc:      nodes.Loc{Start: whenStart, End: p.pos()},
			Cond:     whenCond,
			ThenList: whenThenList,
		})
	}

	// Optional ELSE clause
	if p.cur.Type == kwELSE {
		p.advance() // consume ELSE

		elseList, err := p.parseCompoundStmtList()
		if err != nil {
			return nil, err
		}
		stmt.ElseList = elseList
	}

	// END CASE
	if _, err := p.expect(kwEND); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwCASE); err != nil {
		return nil, err
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseWhileStmt parses a WHILE loop compound statement.
//
//	[begin_label:] WHILE search_condition DO
//	  statement_list
//	END WHILE [end_label]
func (p *Parser) parseWhileStmt(labelName string, labelStart int) (*nodes.WhileStmt, error) {
	start := labelStart
	if labelName == "" {
		start = p.pos()
	}
	p.advance() // consume WHILE

	// Parse condition
	cond, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if cond == nil {
		return nil, p.syntaxErrorAtCur()
	}

	// DO
	if _, err := p.expect(kwDO); err != nil {
		return nil, err
	}

	// Stub for label-registration target node identity; declared after stmts
	// is irrelevant since we register before walking the body.
	stmt := &nodes.WhileStmt{Loc: nodes.Loc{Start: start}, Label: labelName, Cond: cond}
	if labelName != "" {
		if err := p.declareLabel(labelName, labelLoop, stmt, labelStart); err != nil {
			return nil, err
		}
	}

	// Parse statement list until END
	stmts, err := p.parseCompoundStmtList()
	if err != nil {
		return nil, err
	}

	// END WHILE
	if _, err := p.expect(kwEND); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwWHILE); err != nil {
		return nil, err
	}

	// Optional end_label
	var endLabel string
	if p.isLabelIdentToken() && p.cur.Type != tokEOF && p.cur.Type != ';' {
		endLabelTok := p.advance()
		endLabel = endLabelTok.Str
		if err := checkLabelMatch(labelName, endLabel, endLabelTok.Loc); err != nil {
			return nil, err
		}
	}

	stmt.EndLabel = endLabel
	stmt.Stmts = stmts
	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseRepeatStmt parses a REPEAT loop compound statement.
//
//	[begin_label:] REPEAT
//	  statement_list
//	UNTIL search_condition
//	END REPEAT [end_label]
func (p *Parser) parseRepeatStmt(labelName string, labelStart int) (*nodes.RepeatStmt, error) {
	start := labelStart
	if labelName == "" {
		start = p.pos()
	}
	p.advance() // consume REPEAT

	stmt := &nodes.RepeatStmt{Loc: nodes.Loc{Start: start}, Label: labelName}
	if labelName != "" {
		if err := p.declareLabel(labelName, labelLoop, stmt, labelStart); err != nil {
			return nil, err
		}
	}

	// Parse statement list until UNTIL
	stmts, err := p.parseCompoundStmtList()
	if err != nil {
		return nil, err
	}

	// UNTIL
	if _, err := p.expect(kwUNTIL); err != nil {
		return nil, err
	}

	// Parse condition
	cond, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if cond == nil {
		return nil, p.syntaxErrorAtCur()
	}

	// END REPEAT
	if _, err := p.expect(kwEND); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwREPEAT); err != nil {
		return nil, err
	}

	// Optional end_label
	var endLabel string
	if p.isLabelIdentToken() && p.cur.Type != tokEOF && p.cur.Type != ';' {
		endLabelTok := p.advance()
		endLabel = endLabelTok.Str
		if err := checkLabelMatch(labelName, endLabel, endLabelTok.Loc); err != nil {
			return nil, err
		}
	}

	stmt.Stmts = stmts
	stmt.Cond = cond
	stmt.EndLabel = endLabel
	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseLoopStmt parses a LOOP compound statement.
//
//	[begin_label:] LOOP
//	  statement_list
//	END LOOP [end_label]
func (p *Parser) parseLoopStmt(labelName string, labelStart int) (*nodes.LoopStmt, error) {
	start := labelStart
	if labelName == "" {
		start = p.pos()
	}
	p.advance() // consume LOOP

	stmt := &nodes.LoopStmt{Loc: nodes.Loc{Start: start}, Label: labelName}
	if labelName != "" {
		if err := p.declareLabel(labelName, labelLoop, stmt, labelStart); err != nil {
			return nil, err
		}
	}

	// Parse statement list until END
	stmts, err := p.parseCompoundStmtList()
	if err != nil {
		return nil, err
	}

	// END LOOP
	if _, err := p.expect(kwEND); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwLOOP); err != nil {
		return nil, err
	}

	// Optional end_label
	var endLabel string
	if p.isLabelIdentToken() && p.cur.Type != tokEOF && p.cur.Type != ';' {
		endLabelTok := p.advance()
		endLabel = endLabelTok.Str
		if err := checkLabelMatch(labelName, endLabel, endLabelTok.Loc); err != nil {
			return nil, err
		}
	}

	stmt.EndLabel = endLabel
	stmt.Stmts = stmts
	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseLeaveStmt parses a LEAVE label statement.
//
//	LEAVE label
func (p *Parser) parseLeaveStmt() (*nodes.LeaveStmt, error) {
	start := p.pos()
	p.advance() // consume LEAVE

	label, _, err := p.parseLabelIdent()
	if err != nil {
		return nil, err
	}

	return &nodes.LeaveStmt{
		Loc:   nodes.Loc{Start: start, End: p.pos()},
		Label: label,
	}, nil
}

// parseIterateStmt parses an ITERATE label statement. ITERATE only targets
// loop labels (WHILE / LOOP / REPEAT) — not BEGIN-style labels.
//
//	ITERATE label
func (p *Parser) parseIterateStmt() (*nodes.IterateStmt, error) {
	start := p.pos()
	p.advance() // consume ITERATE

	label, _, err := p.parseLabelIdent()
	if err != nil {
		return nil, err
	}

	return &nodes.IterateStmt{
		Loc:   nodes.Loc{Start: start, End: p.pos()},
		Label: label,
	}, nil
}

// parseReturnStmt parses a RETURN expr statement.
//
//	RETURN expr
func (p *Parser) parseReturnStmt() (*nodes.ReturnStmt, error) {
	start := p.pos()
	// MySQL allows RETURN only inside functions. Procedures, triggers, and
	// events use LEAVE / SIGNAL to terminate. The isFunction flag is set on
	// the outermost scope by parseCreateFunctionStmt and inherited by inner
	// scopes via pushScope.
	if p.procScope != nil && !p.procScope.isFunction {
		return nil, &ParseError{
			Message:  "RETURN is only allowed inside a function body",
			Position: start,
		}
	}
	p.advance() // consume RETURN

	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if !p.collectMode() && expr == nil {
		return nil, p.syntaxErrorAtCur()
	}

	return &nodes.ReturnStmt{
		Loc:  nodes.Loc{Start: start, End: p.pos()},
		Expr: expr,
	}, nil
}

// parseOpenCursorStmt parses an OPEN cursor_name statement.
//
//	OPEN cursor_name
func (p *Parser) parseOpenCursorStmt() (*nodes.OpenCursorStmt, error) {
	start := p.pos()
	p.advance() // consume OPEN

	nameStart := p.pos()
	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	if err := p.requireCursor(name, nameStart); err != nil {
		return nil, err
	}
	return &nodes.OpenCursorStmt{
		Loc:  nodes.Loc{Start: start, End: p.pos()},
		Name: name,
	}, nil
}

// requireCursor verifies that the named cursor is declared in scope. Used
// by OPEN / FETCH / CLOSE.
func (p *Parser) requireCursor(name string, pos int) error {
	if p.procScope == nil {
		return nil
	}
	if p.lookupCursor(name) == nil {
		return &ParseError{
			Message:  "undeclared cursor: " + name,
			Position: pos,
		}
	}
	return nil
}

// parseFetchCursorStmt parses a FETCH cursor statement.
//
//	FETCH [[NEXT] FROM] cursor_name INTO var_name [, var_name] ...
func (p *Parser) parseFetchCursorStmt() (*nodes.FetchCursorStmt, error) {
	start := p.pos()
	p.advance() // consume FETCH

	// Optional NEXT
	if p.cur.Type == kwNEXT {
		p.advance()
	}

	// Optional FROM
	if p.cur.Type == kwFROM {
		p.advance()
	}

	// cursor_name
	nameStart := p.pos()
	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	if err := p.requireCursor(name, nameStart); err != nil {
		return nil, err
	}

	// INTO
	if _, err := p.expect(kwINTO); err != nil {
		return nil, err
	}

	// var_name [, var_name] ... — lvalue context (assignment target)
	var vars []string
	for {
		v, _, err := p.parseLvalueIdent()
		if err != nil {
			return nil, err
		}
		vars = append(vars, v)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	return &nodes.FetchCursorStmt{
		Loc:  nodes.Loc{Start: start, End: p.pos()},
		Name: name,
		Into: vars,
	}, nil
}

// parseCloseCursorStmt parses a CLOSE cursor_name statement.
//
//	CLOSE cursor_name
func (p *Parser) parseCloseCursorStmt() (*nodes.CloseCursorStmt, error) {
	start := p.pos()
	p.advance() // consume CLOSE

	nameStart := p.pos()
	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	if err := p.requireCursor(name, nameStart); err != nil {
		return nil, err
	}
	return &nodes.CloseCursorStmt{
		Loc:  nodes.Loc{Start: start, End: p.pos()},
		Name: name,
	}, nil
}

