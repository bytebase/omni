package parser

import (
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

	// Parse statements until END
	for p.cur.Type != kwEND && p.cur.Type != tokEOF {
		s, err := p.parseCompoundStmtOrStmt()
		if err != nil {
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

	// Optional end_label
	if p.isIdentToken() && p.cur.Type != tokEOF {
		stmt.EndLabel = p.cur.Str
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
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
	if p.isIdentToken() && p.cur.Type != kwEND {
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
		stmt.Default = val
	}

	stmt.Loc.End = p.pos()
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

	// handler statement body
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

	return &nodes.DeclareCursorStmt{
		Loc:    nodes.Loc{Start: start, End: p.pos()},
		Name:   name,
		Select: sel,
	}, nil
}

// parseHandlerConditionValue parses a handler condition value.
//
//	condition_value:
//	    mysql_error_code
//	  | SQLSTATE [VALUE] sqlstate_value
//	  | condition_name
//	  | SQLWARNING
//	  | NOT FOUND
//	  | SQLEXCEPTION
func (p *Parser) parseHandlerConditionValue() (string, error) {
	// SQLSTATE [VALUE] 'sqlstate_value'
	if p.cur.Type == kwSQLSTATE {
		p.advance()
		// Optional VALUE
		if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "VALUE") {
			p.advance()
		}
		if p.cur.Type != tokSCONST {
			return "", &ParseError{
				Message:  "expected SQLSTATE value string",
				Position: p.cur.Loc,
			}
		}
		val := p.cur.Str
		p.advance()
		return val, nil
	}

	// SQLWARNING
	if p.cur.Type == kwSQLWARNING {
		p.advance()
		return "SQLWARNING", nil
	}

	// SQLEXCEPTION
	if p.cur.Type == kwSQLEXCEPTION {
		p.advance()
		return "SQLEXCEPTION", nil
	}

	// NOT FOUND
	if p.cur.Type == kwNOT {
		p.advance()
		if _, err := p.expect(kwFOUND); err != nil {
			return "", err
		}
		return "NOT FOUND", nil
	}

	// mysql_error_code (integer literal)
	if p.cur.Type == tokICONST {
		val := p.cur.Str
		p.advance()
		return val, nil
	}

	// condition_name (identifier)
	name, _, err := p.parseIdentifier()
	if err != nil {
		return "", err
	}
	return name, nil
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

	// DO
	if _, err := p.expect(kwDO); err != nil {
		return nil, err
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
	if p.isIdentToken() && p.cur.Type != tokEOF && p.cur.Type != ';' {
		endLabel = p.cur.Str
		p.advance()
	}

	return &nodes.WhileStmt{
		Loc:      nodes.Loc{Start: start, End: p.pos()},
		Label:    labelName,
		EndLabel: endLabel,
		Cond:     cond,
		Stmts:    stmts,
	}, nil
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

	// END REPEAT
	if _, err := p.expect(kwEND); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwREPEAT); err != nil {
		return nil, err
	}

	// Optional end_label
	var endLabel string
	if p.isIdentToken() && p.cur.Type != tokEOF && p.cur.Type != ';' {
		endLabel = p.cur.Str
		p.advance()
	}

	return &nodes.RepeatStmt{
		Loc:      nodes.Loc{Start: start, End: p.pos()},
		Label:    labelName,
		EndLabel: endLabel,
		Stmts:    stmts,
		Cond:     cond,
	}, nil
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
	if p.isIdentToken() && p.cur.Type != tokEOF && p.cur.Type != ';' {
		endLabel = p.cur.Str
		p.advance()
	}

	return &nodes.LoopStmt{
		Loc:      nodes.Loc{Start: start, End: p.pos()},
		Label:    labelName,
		EndLabel: endLabel,
		Stmts:    stmts,
	}, nil
}

// parseLeaveStmt parses a LEAVE label statement.
//
//	LEAVE label
func (p *Parser) parseLeaveStmt() (*nodes.LeaveStmt, error) {
	start := p.pos()
	p.advance() // consume LEAVE

	label, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	return &nodes.LeaveStmt{
		Loc:   nodes.Loc{Start: start, End: p.pos()},
		Label: label,
	}, nil
}

// parseIterateStmt parses an ITERATE label statement.
//
//	ITERATE label
func (p *Parser) parseIterateStmt() (*nodes.IterateStmt, error) {
	start := p.pos()
	p.advance() // consume ITERATE

	label, _, err := p.parseIdentifier()
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
	p.advance() // consume RETURN

	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
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

	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	return &nodes.OpenCursorStmt{
		Loc:  nodes.Loc{Start: start, End: p.pos()},
		Name: name,
	}, nil
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
	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	// INTO
	if _, err := p.expect(kwINTO); err != nil {
		return nil, err
	}

	// var_name [, var_name] ...
	var vars []string
	for {
		v, _, err := p.parseIdentifier()
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

	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	return &nodes.CloseCursorStmt{
		Loc:  nodes.Loc{Start: start, End: p.pos()},
		Name: name,
	}, nil
}

