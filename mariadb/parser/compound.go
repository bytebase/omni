package parser

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/mariadb/ast"
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
//
// An unlabeled block may be MariaDB's BEGIN NOT ATOMIC, which holds what any block holds
// and is also a statement of its own; a labeled one may not.
//
// Ref: sp_labeled_block, sp_unlabeled_block and sp_unlabeled_block_not_atomic in MariaDB
// 11.8 sql/sql_yacc.yy.
func (p *Parser) parseBeginEndBlock(labelName string, labelStart int) (*nodes.BeginEndBlock, error) {
	p.routineDepth++
	defer func() { p.routineDepth-- }()
	start := labelStart
	if labelName == "" {
		start = p.pos()
	}

	p.advance() // consume BEGIN

	stmt := &nodes.BeginEndBlock{
		Loc:   nodes.Loc{Start: start},
		Label: labelName,
	}
	// not is NOT_SYM, or NOT2_SYM, which `!` is only under HIGH_NOT_PRECEDENCE.
	if labelName == "" && p.cur.Type == kwNOT {
		p.advance()
		if !p.atWord("ATOMIC") {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		stmt.NotAtomic = true
	}
	defer p.pushSpScope()()

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

		if err := p.expectStoredProgramStmtTerminator(); err != nil {
			return nil, err
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

	stmt.Loc.End = p.prev.End
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
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()

	// Label: check for identifier followed by ':', ahead of the statements a word can start,
	// so `close: LOOP` is a label and not CLOSE.
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
			// Labeled FOR
			if p.cur.Type == kwFOR {
				return p.parseForStmt(name, nameStart)
			}
			return nil, &ParseError{
				Message:  "expected BEGIN, WHILE, REPEAT, LOOP, or FOR after label",
				Position: p.cur.Loc,
			}
		}
	}

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

	// FOR (without label)
	if p.cur.Type == kwFOR {
		return p.parseForStmt("", 0)
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
		// Only a function returns: a procedure, a trigger and an event answer 1313 while parsing.
		if !p.routine.function {
			return nil, &ParseError{Message: "RETURN is only allowed in a FUNCTION", Position: p.cur.Loc}
		}
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

	// Regular statement
	start := p.pos()
	stmt, err := p.parseStmt()
	if err != nil {
		return nil, err
	}
	if err := p.refuseResultSet(stmt, start); err != nil {
		return nil, err
	}
	return stmt, nil
}

// parseDeclareStmt parses DECLARE statements inside compound blocks.
// Dispatches to variable, condition, handler, or cursor declarations.
//
// DECLARE var [, var] ... type [DEFAULT value]
// DECLARE condition_name CONDITION FOR condition_value
// DECLARE handler_action HANDLER FOR condition_value [, ...] stmt
// DECLARE cursor_name CURSOR [(param, ...)] FOR select_stmt
func (p *Parser) parseDeclareStmt() (nodes.Node, error) {
	start := p.pos()
	p.advance() // consume DECLARE

	// Check for handler: CONTINUE | EXIT | UNDO
	if p.cur.Type == kwCONTINUE || p.cur.Type == kwEXIT || p.cur.Type == kwUNDO {
		return p.parseDeclareHandlerStmt(start)
	}

	// We need to look ahead to determine if it's CONDITION, CURSOR, or variable.
	// Parse the first identifier name.
	name, namePos, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	// DECLARE name CONDITION FOR ...
	if p.cur.Type == kwCONDITION {
		return p.parseDeclareConditionStmt(start, name, namePos)
	}

	// DECLARE name CURSOR FOR ...
	if p.cur.Type == kwCURSOR {
		return p.parseDeclareCursorStmt(start, name, namePos)
	}

	// DECLARE var [, var] ... type [DEFAULT value]
	return p.parseDeclareVarStmt(start, name, namePos)
}

// parseDeclareVarStmt parses a DECLARE variable statement.
// The first variable name has already been consumed.
//
//	DECLARE var_name [, var_name] ... type [DEFAULT value]
//
// MariaDB's type may also be ROW(name type, ...), TYPE OF a variable or a column, or ROW TYPE
// OF a table or a cursor. A block declares a name once, and the names are declared after the
// declaration, so its TYPE OF or DEFAULT does not see them.
//
// Ref: sp_decl_variable_list and sp_decl_variable_list_anchored in MariaDB 11.8
// sql/sql_yacc.yy.
func (p *Parser) parseDeclareVarStmt(start int, firstName string, firstPos int) (*nodes.DeclareVarStmt, error) {
	stmt := &nodes.DeclareVarStmt{
		Loc:   nodes.Loc{Start: start},
		Names: []string{firstName},
	}
	positions := []int{firstPos}

	// Additional variable names separated by commas
	for p.cur.Type == ',' {
		p.advance()
		name, pos, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Names = append(stmt.Names, name)
		positions = append(positions, pos)
	}

	// Data type
	dt, row, anchor, err := p.parseSpVarType(true)
	if err != nil {
		return nil, err
	}
	stmt.TypeName, stmt.RowFields, stmt.Anchor = dt, row, anchor

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

	for i, name := range stmt.Names {
		if err := p.declareSpVar(name, positions[i]); err != nil {
			return nil, err
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseSpVarType reads the type of a stored program variable or parameter: a data type, or
// one of MariaDB's ROW(name type, ...), TYPE OF and ROW TYPE OF.
//
// A variable's TYPE OF and ROW TYPE OF name one to three parts, and TYPE OF one part is a
// variable, which must be declared. A parameter's TYPE OF names a column, two or three
// parts, and its ROW TYPE OF a table, one or two.
//
// Ref: sp_decl_variable_list, sp_decl_variable_list_anchored,
// optionally_qualified_column_ident, sp_param_init_vars and sp_param_anchored in MariaDB 11.8
// sql/sql_yacc.yy.
func (p *Parser) parseSpVarType(variable bool) (*nodes.DataType, []*nodes.RowField, *nodes.TypeAnchor, error) {
	switch {
	case p.cur.Type == kwROW && p.peekNext().Type == '(':
		row, err := p.parseRowTypeBody()
		return nil, row, nil, err
	case p.cur.Type == kwROW && p.peekNext().Type == kwTYPE,
		p.cur.Type == kwTYPE && p.peekNext().Type == kwOF:
		anchor := &nodes.TypeAnchor{Loc: nodes.Loc{Start: p.pos()}}
		if _, ok := p.match(kwROW); ok {
			anchor.Row = true
		}
		p.advance() // TYPE
		if _, err := p.expect(kwOF); err != nil {
			return nil, nil, nil, err
		}
		minParts, maxParts := 1, 3
		if !variable {
			minParts, maxParts = 2, 3
			if anchor.Row {
				minParts, maxParts = 1, 2
			}
		}
		firstPos := p.cur.Loc
		for {
			name, _, err := p.parseIdentifier()
			if err != nil {
				return nil, nil, nil, err
			}
			anchor.Names = append(anchor.Names, name)
			if len(anchor.Names) == maxParts || p.cur.Type != '.' {
				break
			}
			p.advance()
		}
		if len(anchor.Names) < minParts {
			return nil, nil, nil, p.syntaxErrorAtCur()
		}
		anchor.Loc.End = p.prev.End
		if !anchor.Row && len(anchor.Names) == 1 {
			if err := p.checkSpVarDeclared(anchor.Names[0], firstPos); err != nil {
				return nil, nil, nil, err
			}
		}
		return nil, nil, anchor, nil
	}
	dt, err := p.parseDataType()
	return dt, nil, nil, err
}

// parseRowTypeBody reads ROW(name type, ...), a field named once each (1060).
//
// Ref: row_type_body and row_field_definition_list in MariaDB 11.8 sql/sql_yacc.yy.
func (p *Parser) parseRowTypeBody() ([]*nodes.RowField, error) {
	p.advance() // ROW
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	var fields []*nodes.RowField
	seen := map[string]bool{}
	for {
		start := p.pos()
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		dt, err := p.parseDataType()
		if err != nil {
			return nil, err
		}
		key := strings.ToLower(name)
		if seen[key] {
			return nil, &ParseError{Message: "Duplicate column name '" + name + "'", Position: start}
		}
		seen[key] = true
		fields = append(fields, &nodes.RowField{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: name, TypeName: dt})
		if _, ok := p.match(','); !ok {
			break
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return fields, nil
}

// parseDeclareConditionStmt parses DECLARE condition_name CONDITION FOR condition_value.
// The condition name has already been consumed.
//
//	DECLARE condition_name CONDITION FOR condition_value
//	condition_value: SQLSTATE [VALUE] sqlstate_value | mysql_error_code
//
// Ref: sp_decl_body and sp_cond in MariaDB 11.8 sql/sql_yacc.yy.
func (p *Parser) parseDeclareConditionStmt(start int, name string, namePos int) (*nodes.DeclareConditionStmt, error) {
	p.advance() // consume CONDITION

	// FOR
	if _, err := p.expect(kwFOR); err != nil {
		return nil, err
	}

	stmt := &nodes.DeclareConditionStmt{
		Loc:  nodes.Loc{Start: start},
		Name: name,
	}

	// condition_value: a number or an SQLSTATE, not a class or another condition's name.
	if p.cur.Type != kwSQLSTATE && p.cur.Type != tokICONST && p.cur.Type != tokFCONST {
		return nil, p.syntaxErrorAtCur()
	}
	condVal, err := p.parseHandlerConditionValue()
	if err != nil {
		return nil, err
	}
	stmt.ConditionValue = condVal
	if err := p.declareSpCondition(name, condVal, namePos); err != nil {
		return nil, err
	}

	stmt.Loc.End = p.prev.End
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
		condPos := p.pos()
		condVal, err := p.parseHandlerConditionValue()
		if err != nil {
			return nil, err
		}
		if err := p.declareSpHandlerCondition(condVal, condPos); err != nil {
			return nil, err
		}
		stmt.Conditions = append(stmt.Conditions, condVal)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	body, err := p.parseCompoundStmtOrStmt()
	if err != nil {
		return nil, err
	}
	stmt.Stmt = body

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseDeclareCursorStmt parses DECLARE cursor_name CURSOR FOR select_statement.
// The cursor name has already been consumed.
//
//	DECLARE cursor_name CURSOR [(param [, param] ...)] FOR select_statement
//
// MariaDB's cursor takes parameters, each [IN] name type with no DEFAULT, named once each
// (1330); OUT and INOUT parse and are refused with 1235. They are visible to the cursor's
// query alone.
//
// Ref: sp_decl_body, opt_parenthesized_cursor_formal_parameters and sp_fdparams_no_default in
// MariaDB 11.8 sql/sql_yacc.yy.
func (p *Parser) parseDeclareCursorStmt(start int, name string, namePos int) (*nodes.DeclareCursorStmt, error) {
	p.advance() // consume CURSOR

	stmt := &nodes.DeclareCursorStmt{
		Loc:  nodes.Loc{Start: start},
		Name: name,
	}

	if p.cur.Type == '(' {
		params, err := p.parseCursorParams()
		if err != nil {
			return nil, err
		}
		stmt.Params = params
	}

	// FOR
	if _, err := p.expect(kwFOR); err != nil {
		return nil, err
	}

	// select_statement, which takes no INTO (1323)
	queryStart := p.pos()
	sel, err := p.parseSelectStmt()
	if err != nil {
		return nil, err
	}
	if selectHasInto(sel) {
		return nil, &ParseError{Message: "Cursor SELECT must not have INTO", Position: queryStart}
	}
	stmt.Select = sel
	if err := p.declareSpCursor(name, namePos); err != nil {
		return nil, err
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseCursorParams reads a cursor's parenthesized parameters in a scope of their own.
func (p *Parser) parseCursorParams() ([]*nodes.FuncParam, error) {
	p.advance() // (
	defer p.pushSpScope()()
	var params []*nodes.FuncParam
	for {
		param, err := p.parseSpParam(true)
		if err != nil {
			return nil, err
		}
		if param.Direction == "OUT" || param.Direction == "INOUT" {
			return nil, &ParseError{Message: "This version of MariaDB doesn't yet support 'OUT/INOUT cursor parameter'", Position: param.Loc.Start}
		}
		params = append(params, param)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return params, nil
}

// parseSpParam reads one parameter of a routine or a cursor, [IN | OUT | INOUT] name type,
// and declares it in the current scope. MariaDB's type may be ROW(...), TYPE OF a column or
// ROW TYPE OF a table.
//
// Ref: sp_param_name_and_mode, sp_param_init_vars and sp_param_anchored in MariaDB 11.8
// sql/sql_yacc.yy.
func (p *Parser) parseSpParam(withDirection bool) (*nodes.FuncParam, error) {
	start := p.pos()
	param := &nodes.FuncParam{Loc: nodes.Loc{Start: start}}

	if withDirection {
		switch p.cur.Type {
		case kwIN:
			param.Direction = "IN"
			p.advance()
		case kwOUT:
			param.Direction = "OUT"
			p.advance()
		case kwINOUT:
			param.Direction = "INOUT"
			p.advance()
		}
	}

	name, namePos, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	param.Name = name

	dt, row, anchor, err := p.parseSpVarType(false)
	if err != nil {
		return nil, err
	}
	param.TypeName, param.RowFields, param.Anchor = dt, row, anchor
	if err := p.declareSpParam(name, namePos); err != nil {
		return nil, err
	}

	param.Loc.End = p.prev.End
	return param, nil
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
		val, err := p.parseSQLStateValue()
		if err != nil {
			return nodes.HandlerCondValue{}, err
		}
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

	// mysql_error_code is ulong_num, which takes a decimal and reads its whole part,
	// and the action refuses 0: mariadb:11.8.8 parses FOR 1.5 and FOR
	// 99999999999999999999, and refuses FOR 0 and FOR 0.5 while it parses.
	//
	// Ref: sp_cond in MariaDB 11.8 sql/sql_yacc.yy.
	if p.cur.Type == tokICONST || p.cur.Type == tokFCONST {
		if p.cur.Ival == 0 {
			return nodes.HandlerCondValue{}, &ParseError{Message: "Incorrect CONDITION value: '0'", Position: p.cur.Loc}
		}
		val := leadingDigits(p.cur.Str)
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
// stopping at any compound terminator keyword. The list holds one statement at least.
//
// Ref: sp_proc_stmts1, which MariaDB 11.8 sql/sql_yacc.yy reads for the statements of IF,
// CASE, LOOP, WHILE, REPEAT and FOR.
func (p *Parser) parseCompoundStmtList() ([]nodes.Node, error) {
	var stmts []nodes.Node
	if p.isCompoundTerminator() && !p.completing {
		return nil, p.syntaxErrorAtCur()
	}
	for !p.isCompoundTerminator() {
		s, err := p.parseCompoundStmtOrStmt()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, s)
		if err := p.expectStoredProgramStmtTerminator(); err != nil {
			return nil, err
		}
	}
	return stmts, nil
}

func (p *Parser) expectStoredProgramStmtTerminator() error {
	if p.cur.Type == ';' {
		p.advance()
		return nil
	}
	return p.syntaxErrorAtCur()
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
			Loc:      nodes.Loc{Start: eiStart, End: p.prev.End},
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

	stmt.Loc.End = p.prev.End
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
			Loc:      nodes.Loc{Start: whenStart, End: p.prev.End},
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

	stmt.Loc.End = p.prev.End
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

	stmt := &nodes.WhileStmt{Loc: nodes.Loc{Start: start}, Label: labelName, Cond: cond}

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
	stmt.Loc.End = p.prev.End
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
	stmt.Loc.End = p.prev.End
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
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseForStmt parses MariaDB's FOR loop over the rows of a query.
//
//	[begin_label:] FOR record IN [REVERSE] (query) DO
//	  statement_list
//	END FOR [end_label]
//
// The record is a variable of the loop's body. The query takes no INTO (1323). The loop
// body holds one statement at least.
//
// omni reads no integer range (IN low..high) and no cursor (IN cursor_name): a name in
// their bounds is a variable or a cursor the blocks around must declare, and omni does not
// check a name an expression reads.
//
// Ref: sp_unlabeled_control, sp_labeled_control, sp_for_loop_bounds and
// for_loop_statements in MariaDB 11.8 sql/sql_yacc.yy.
func (p *Parser) parseForStmt(labelName string, labelStart int) (*nodes.ForStmt, error) {
	start := labelStart
	if labelName == "" {
		start = p.pos()
	}
	p.advance() // consume FOR

	stmt := &nodes.ForStmt{Loc: nodes.Loc{Start: start}, Label: labelName}
	name, namePos, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Var = name
	if _, err := p.expect(kwIN); err != nil {
		return nil, err
	}
	if _, ok := p.match(kwREVERSE); ok {
		stmt.Reverse = true
	}
	if p.cur.Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	queryStart := p.pos()
	sel, err := p.parseSelectStmt()
	if err != nil {
		return nil, err
	}
	if selectHasInto(sel) {
		return nil, &ParseError{Message: "Cursor SELECT must not have INTO", Position: queryStart}
	}
	stmt.Select = sel
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwDO); err != nil {
		return nil, err
	}

	closeScope := p.pushSpScope()
	if err := p.declareSpVar(name, namePos); err != nil {
		closeScope()
		return nil, err
	}
	stmts, err := p.parseCompoundStmtList()
	closeScope()
	if err != nil {
		return nil, err
	}

	// END FOR
	if _, err := p.expect(kwEND); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwFOR); err != nil {
		return nil, err
	}

	// Optional end_label
	if p.isLabelIdentToken() && p.cur.Type != tokEOF && p.cur.Type != ';' {
		endLabelTok := p.advance()
		stmt.EndLabel = endLabelTok.Str
		if err := checkLabelMatch(labelName, stmt.EndLabel, endLabelTok.Loc); err != nil {
			return nil, err
		}
	}

	stmt.Stmts = stmts
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseCompoundStatement parses compound_statement, a compound statement standing as a
// statement of its own: BEGIN NOT ATOMIC ... END, an unlabeled LOOP, or a FOR loop. It is a
// stored program of its own, which MariaDB runs as an anonymous block.
//
// IF, CASE, WHILE and REPEAT stand there too, but omni leaves them to a stored program: the
// name in their condition is a variable the server holds to a declaration (1327 for `IF a
// THEN`), and omni does not check a name an expression reads.
//
// Ref: directly_executable_statement, compound_statement and sp_proc_stmt_compound_ok in
// MariaDB 11.8 sql/sql_yacc.yy.
func (p *Parser) parseCompoundStatement() (nodes.Node, error) {
	p.routineDepth++
	defer func() { p.routineDepth-- }()
	defer p.enterStoredProgram()()
	outerRoutine := p.routine
	p.routine = routineScope{}
	defer func() { p.routine = outerRoutine }()
	defer p.pushSpScope()()

	switch p.cur.Type {
	case kwBEGIN:
		return p.parseBeginEndBlock("", 0)
	case kwLOOP:
		return p.parseLoopStmt("", 0)
	case kwFOR:
		return p.parseForStmt("", 0)
	}
	return nil, p.syntaxErrorAtCur()
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
		Loc:   nodes.Loc{Start: start, End: p.prev.End},
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
		Loc:   nodes.Loc{Start: start, End: p.prev.End},
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
	if !p.collectMode() && expr == nil {
		return nil, p.syntaxErrorAtCur()
	}

	return &nodes.ReturnStmt{
		Loc:  nodes.Loc{Start: start, End: p.prev.End},
		Expr: expr,
	}, nil
}

// parseOpenCursorStmt parses an OPEN cursor_name statement.
//
//	OPEN cursor_name
//
// omni reads no arguments for MariaDB's cursor parameters (OPEN name(expr, ...)): the server
// holds the name to a cursor the blocks around declare, which omni does not check.
func (p *Parser) parseOpenCursorStmt() (*nodes.OpenCursorStmt, error) {
	start := p.pos()
	p.advance() // consume OPEN

	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	return &nodes.OpenCursorStmt{
		Loc:  nodes.Loc{Start: start, End: p.prev.End},
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
		Loc:  nodes.Loc{Start: start, End: p.prev.End},
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
		Loc:  nodes.Loc{Start: start, End: p.prev.End},
		Name: name,
	}, nil
}
