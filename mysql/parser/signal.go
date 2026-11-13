package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mysql/ast"
)

// parseSignalStmt parses a SIGNAL statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/signal.html
//
//	SIGNAL condition_value
//	    [SET signal_information_item
//	    [, signal_information_item] ...]
//
//	condition_value: {
//	    SQLSTATE [VALUE] sqlstate_value
//	  | condition_name
//	}
//
//	signal_information_item:
//	    condition_information_item_name = simple_value_specification
//
//	condition_information_item_name: {
//	    CLASS_ORIGIN
//	  | SUBCLASS_ORIGIN
//	  | MESSAGE_TEXT
//	  | MYSQL_ERRNO
//	  | CONSTRAINT_CATALOG
//	  | CONSTRAINT_SCHEMA
//	  | CONSTRAINT_NAME
//	  | CATALOG_NAME
//	  | SCHEMA_NAME
//	  | TABLE_NAME
//	  | COLUMN_NAME
//	  | CURSOR_NAME
//	}
func (p *Parser) parseSignalStmt() (*nodes.SignalStmt, error) {
	start := p.pos()
	p.advance() // consume SIGNAL

	stmt := &nodes.SignalStmt{Loc: nodes.Loc{Start: start}}

	// Parse condition_value (required for SIGNAL)
	condVal, err := p.parseSignalConditionValue()
	if err != nil {
		return nil, err
	}
	stmt.ConditionValue = condVal

	// Optional SET clause
	if p.cur.Type == kwSET {
		items, err := p.parseSignalSetClause()
		if err != nil {
			return nil, err
		}
		stmt.SetItems = items
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseResignalStmt parses a RESIGNAL statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/resignal.html
//
//	RESIGNAL [condition_value]
//	    [SET signal_information_item
//	    [, signal_information_item] ...]
//
//	condition_value: {
//	    SQLSTATE [VALUE] sqlstate_value
//	  | condition_name
//	}
//
//	signal_information_item:
//	    condition_information_item_name = simple_value_specification
//
//	condition_information_item_name: {
//	    CLASS_ORIGIN
//	  | SUBCLASS_ORIGIN
//	  | MESSAGE_TEXT
//	  | MYSQL_ERRNO
//	  | CONSTRAINT_CATALOG
//	  | CONSTRAINT_SCHEMA
//	  | CONSTRAINT_NAME
//	  | CATALOG_NAME
//	  | SCHEMA_NAME
//	  | TABLE_NAME
//	  | COLUMN_NAME
//	  | CURSOR_NAME
//	}
func (p *Parser) parseResignalStmt() (*nodes.ResignalStmt, error) {
	start := p.pos()
	p.advance() // consume RESIGNAL

	stmt := &nodes.ResignalStmt{Loc: nodes.Loc{Start: start}}

	// Optional condition_value
	// condition_value starts with SQLSTATE or an identifier (condition_name)
	// but NOT SET (which starts the SET clause)
	if p.cur.Type == kwSQLSTATE {
		condVal, err := p.parseSignalConditionValue()
		if err != nil {
			return nil, err
		}
		stmt.ConditionValue = condVal
	} else if p.cur.Type == tokIDENT {
		// condition_name — a plain identifier that is not SET
		condVal, err := p.parseSignalConditionValue()
		if err != nil {
			return nil, err
		}
		stmt.ConditionValue = condVal
	}

	// Optional SET clause
	if p.cur.Type == kwSET {
		items, err := p.parseSignalSetClause()
		if err != nil {
			return nil, err
		}
		stmt.SetItems = items
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseGetDiagnosticsStmt parses a GET DIAGNOSTICS statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/get-diagnostics.html
//
//	GET [CURRENT | STACKED] DIAGNOSTICS {
//	    statement_information_item
//	    [, statement_information_item] ...
//	  | CONDITION condition_number
//	    condition_information_item
//	    [, condition_information_item] ...
//	}
//
//	statement_information_item:
//	    target = statement_information_item_name
//
//	statement_information_item_name: {
//	    NUMBER
//	  | ROW_COUNT
//	}
//
//	condition_information_item:
//	    target = condition_information_item_name
//
//	condition_information_item_name: {
//	    CLASS_ORIGIN
//	  | SUBCLASS_ORIGIN
//	  | RETURNED_SQLSTATE
//	  | MESSAGE_TEXT
//	  | MYSQL_ERRNO
//	  | CONSTRAINT_CATALOG
//	  | CONSTRAINT_SCHEMA
//	  | CONSTRAINT_NAME
//	  | CATALOG_NAME
//	  | SCHEMA_NAME
//	  | TABLE_NAME
//	  | COLUMN_NAME
//	  | CURSOR_NAME
//	}
//
//	target:
//	    user_variable | local_variable | param
func (p *Parser) parseGetDiagnosticsStmt() (*nodes.GetDiagnosticsStmt, error) {
	start := p.pos()
	p.advance() // consume GET

	stmt := &nodes.GetDiagnosticsStmt{Loc: nodes.Loc{Start: start}}

	// Optional CURRENT | STACKED
	if p.cur.Type == kwCURRENT {
		p.advance()
	} else if p.cur.Type == kwSTACKED {
		p.advance()
		stmt.Stacked = true
	}

	// DIAGNOSTICS
	if _, err := p.expect(kwDIAGNOSTICS); err != nil {
		return nil, err
	}

	// CONDITION or statement_information_item
	if p.cur.Type == kwCONDITION {
		p.advance() // consume CONDITION
		stmt.StatementInfo = false

		// condition_number
		condNum, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if !p.collectMode() && condNum == nil {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.ConditionNumber = condNum

		// condition_information_item [, condition_information_item] ...
		items, err := p.parseDiagnosticsItems(false)
		if err != nil {
			return nil, err
		}
		stmt.Items = items
	} else {
		// statement_information_item [, statement_information_item] ...
		stmt.StatementInfo = true
		items, err := p.parseDiagnosticsItems(true)
		if err != nil {
			return nil, err
		}
		stmt.Items = items
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseSignalConditionValue parses a signal condition value.
//
//	condition_value: {
//	    SQLSTATE [VALUE] sqlstate_value
//	  | condition_name
//	}
func (p *Parser) parseSignalConditionValue() (string, error) {
	// SQLSTATE [VALUE] 'sqlstate_value'
	if p.cur.Type == kwSQLSTATE {
		return p.parseSQLStateValue()
	}

	// condition_name (identifier)
	name, _, err := p.parseIdent()
	if err != nil {
		return "", err
	}
	return name, nil
}

// parseSignalSetClause parses SET signal_information_item [, signal_information_item] ...
func (p *Parser) parseSignalSetClause() ([]*nodes.SignalInfoItem, error) {
	p.advance() // consume SET

	var items []*nodes.SignalInfoItem
	seen := map[string]bool{}
	for {
		at := p.cur.Loc
		item, err := p.parseSignalInfoItem()
		if err != nil {
			return nil, err
		}
		// Each item comes once (ER_DUP_SIGNAL_SET).
		if seen[item.Name] {
			return nil, &ParseError{Message: "Duplicate condition information item '" + item.Name + "'", Position: at}
		}
		seen[item.Name] = true
		items = append(items, item)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	return items, nil
}

// parseSignalInfoItem parses a single signal information item.
//
//	signal_information_item:
//	    condition_information_item_name = simple_value_specification
func (p *Parser) parseSignalInfoItem() (*nodes.SignalInfoItem, error) {
	start := p.pos()

	// signal_condition_information_item_name: one of signalItemNames.
	name := p.signalItemWord()
	if !signalItemNames[name] {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()

	// '='
	if _, err := p.expect('='); err != nil {
		return nil, err
	}

	// simple_value_specification (expression)
	val, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, p.syntaxErrorAtCur()
	}

	return &nodes.SignalInfoItem{
		Loc:   nodes.Loc{Start: start, End: p.prev.End},
		Name:  name,
		Value: val,
	}, nil
}

// parseDiagnosticsItems parses target = item_name [, target = item_name] ...
func (p *Parser) parseDiagnosticsItems(statement bool) ([]*nodes.DiagnosticsItem, error) {
	var items []*nodes.DiagnosticsItem
	for {
		item, err := p.parseDiagnosticsItem(statement)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	return items, nil
}

// parseDiagnosticsItem parses a single diagnostics item (target = item_name).
func (p *Parser) parseDiagnosticsItem(statement bool) (*nodes.DiagnosticsItem, error) {
	start := p.pos()

	// target: user_variable (@var) | local_variable (identifier) | param
	var target nodes.Node
	if p.isVariableRef() {
		t, err := p.parseVariableRef()
		if err != nil {
			return nil, err
		}
		target = t
	} else if p.isIdentToken() {
		name, _, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		target = &nodes.ColumnRef{
			Loc:    nodes.Loc{Start: start, End: p.prev.End},
			Column: name,
		}
	} else {
		return nil, &ParseError{
			Message:  "expected variable or identifier as diagnostics target",
			Position: p.cur.Loc,
		}
	}

	// '='
	if _, err := p.expect('='); err != nil {
		return nil, err
	}

	// statement_information_item_name or condition_information_item_name.
	name := p.signalItemWord()
	if (statement && !statementInfoNames[name]) || (!statement && !conditionInfoNames[name]) {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()

	return &nodes.DiagnosticsItem{
		Loc:    nodes.Loc{Start: start, End: p.prev.End},
		Target: target,
		Name:   name,
	}, nil
}

// signalItemNames, statementInfoNames and conditionInfoNames are the words
// signal_condition_information_item_name, statement_information_item_name and
// condition_information_item_name spell in MySQL 8.0 sql/sql_yacc.yy.
var (
	signalItemNames = map[string]bool{
		"CLASS_ORIGIN": true, "SUBCLASS_ORIGIN": true, "CONSTRAINT_CATALOG": true,
		"CONSTRAINT_SCHEMA": true, "CONSTRAINT_NAME": true, "CATALOG_NAME": true,
		"SCHEMA_NAME": true, "TABLE_NAME": true, "COLUMN_NAME": true, "CURSOR_NAME": true,
		"MESSAGE_TEXT": true, "MYSQL_ERRNO": true,
	}
	statementInfoNames = map[string]bool{"NUMBER": true, "ROW_COUNT": true}
	conditionInfoNames = map[string]bool{
		"CLASS_ORIGIN": true, "SUBCLASS_ORIGIN": true, "CONSTRAINT_CATALOG": true,
		"CONSTRAINT_SCHEMA": true, "CONSTRAINT_NAME": true, "CATALOG_NAME": true,
		"SCHEMA_NAME": true, "TABLE_NAME": true, "COLUMN_NAME": true, "CURSOR_NAME": true,
		"MESSAGE_TEXT": true, "MYSQL_ERRNO": true, "RETURNED_SQLSTATE": true,
	}
)

// parseSQLStateValue reads `SQLSTATE [VALUE] 'value'`, p.cur being SQLSTATE.
// The action refuses a value that is not five digits or capital letters, and
// one of class 00, the completion condition (ER_SP_BAD_SQLSTATE).
//
// Ref: sqlstate in MySQL 8.0 sql/sql_yacc.yy, is_sqlstate_valid and
// is_sqlstate_completion.
func (p *Parser) parseSQLStateValue() (string, error) {
	p.advance() // consume SQLSTATE
	if p.cur.Type == kwVALUE {
		p.advance()
	}
	if p.cur.Type != tokSCONST {
		return "", &ParseError{
			Message:  "expected SQLSTATE value string",
			Position: p.cur.Loc,
		}
	}
	val := p.cur.Str
	if !isSQLStateValid(val) || strings.HasPrefix(val, "00") {
		return "", &ParseError{Message: "Bad SQLSTATE: '" + val + "'", Position: p.cur.Loc}
	}
	p.advance()
	return val, nil
}

// isSQLStateValid reports whether s is five characters, each a digit or a
// capital letter.
func isSQLStateValid(s string) bool {
	if len(s) != 5 {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !(c >= '0' && c <= '9') && !(c >= 'A' && c <= 'Z') {
			return false
		}
	}
	return true
}

// signalItemWord returns the upper-cased word at the current token when it is
// written bare, as the grammar's item-name keywords are, and "" otherwise.
func (p *Parser) signalItemWord() string {
	if p.cur.Type != tokIDENT && p.cur.Type < 700 {
		return ""
	}
	text := p.inputText(p.cur.Loc, p.cur.End)
	if strings.HasPrefix(text, "`") || strings.HasPrefix(text, "@") {
		return ""
	}
	return strings.ToUpper(text)
}
